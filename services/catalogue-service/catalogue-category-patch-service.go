package catalogueService

import (
	"errors"
	"fmt"

	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/google/uuid"
)

// fetchCategoryGroup loads a single group by UID scoped to a category. Returns
// ERR_NOT_FOUND if either UID is unknown or the group belongs to a different category.
// Does NOT rely on GetCatalogueCategoryWithDetailsByUid because that read filters out
// groups with no properties — we need to see freshly-created empty groups too.
func (svc *CatalogueService) fetchCategoryGroup(categoryUID, groupUID string) (models.CatalogueCategoryPropertyGroup, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	g, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.CatalogueCategoryPropertyGroup](session, GetCatalogueCategoryGroupByUidsQuery(categoryUID, groupUID))
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return g, helpers.ERR_NOT_FOUND
		}
		return g, err
	}
	return g, nil
}

// listCategoryGroups loads every group under a category (including those with no
// properties). Needed by PATCH group's lazy-seed step to see the complete sibling set.
func (svc *CatalogueService) listCategoryGroups(categoryUID string) ([]models.CatalogueCategoryPropertyGroup, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	return helpers.GetNeo4jArrayOfNodes[models.CatalogueCategoryPropertyGroup](session, ListCatalogueCategoryGroupsQuery(categoryUID))
}

func (svc *CatalogueService) CreateCatalogueCategoryGroup(categoryUID string, fields *models.CreateCatalogueCategoryGroupFields, userUID string) (result models.CatalogueCategoryPropertyGroup, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	// Resolve the target order upfront: payload value wins; otherwise compute
	// max(siblings.order)+10 via a lightweight read so the write is a single round-trip.
	var order int
	if fields.Order != nil {
		order = *fields.Order
	} else {
		nextOrder, qerr := helpers.GetNeo4jSingleRecordSingleValue[int64](session, NextGroupOrderQuery(categoryUID))
		if qerr != nil {
			if errors.Is(qerr, helpers.ERR_NO_ROWS) {
				return result, helpers.ERR_NOT_FOUND
			}
			return result, qerr
		}
		order = int(nextOrder)
	}

	newUID := uuid.NewString()
	query := CreateCatalogueCategoryGroupQuery(categoryUID, newUID, userUID, fields.Name, order)
	returnedUID, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return result, helpers.ERR_NOT_FOUND
		}
		return result, err
	}
	if returnedUID == "" {
		return result, helpers.ERR_NOT_FOUND
	}

	result = models.CatalogueCategoryPropertyGroup{
		UID:   newUID,
		Name:  fields.Name,
		Order: &order,
	}
	return result, nil
}

func (svc *CatalogueService) PatchCatalogueCategoryGroup(categoryUID, groupUID string, fields *models.PatchCatalogueCategoryGroupFields, userUID string) (result models.CatalogueCategoryPropertyGroup, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	group, err := svc.fetchCategoryGroup(categoryUID, groupUID)
	if err != nil {
		return result, err
	}

	// Lazy order seed: if the caller is setting order but some siblings still lack
	// explicit order values (legacy data), renumber all unseeded siblings before applying
	// the requested value so the resulting ordering is deterministic.
	if fields.Order != nil {
		siblings, lerr := svc.listCategoryGroups(categoryUID)
		if lerr != nil {
			return result, lerr
		}
		allSeeded := true
		for _, g := range siblings {
			if g.Order == nil {
				allSeeded = false
				break
			}
		}
		if !allSeeded {
			_, serr := helpers.WriteNeo4jAndReturnSingleValue[int64](session, SeedCategoryGroupOrdersQuery(categoryUID))
			if serr != nil && !errors.Is(serr, helpers.ERR_NO_ROWS) {
				return result, serr
			}
			// Re-fetch to get post-seed order value for accurate change diff.
			seeded, ferr := svc.fetchCategoryGroup(categoryUID, groupUID)
			if ferr == nil {
				group = seeded
			}
		}
	}

	query := PatchCatalogueCategoryGroupQuery(categoryUID, groupUID, fields, &group, userUID)
	returnedUID, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return result, helpers.ERR_NOT_FOUND
		}
		return result, err
	}
	if returnedUID == "" {
		return result, helpers.ERR_NOT_FOUND
	}

	return svc.fetchCategoryGroup(categoryUID, groupUID)
}

// fetchCategoryProperty loads a property with its parent group UID, type, unit, and
// ordered fields so the PATCH flow can compute accurate change diffs.
func (svc *CatalogueService) fetchCategoryProperty(categoryUID, propertyUID string) (CategoryPropertyWithGroup, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	p, err := helpers.GetNeo4jSingleRecordAndMapToStruct[CategoryPropertyWithGroup](session, GetCatalogueCategoryPropertyByUidsQuery(categoryUID, propertyUID))
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return p, helpers.ERR_NOT_FOUND
		}
		return p, err
	}
	return p, nil
}

func (svc *CatalogueService) CreateCatalogueCategoryProperty(categoryUID, groupUID string, fields *models.CreateCatalogueCategoryPropertyFields, userUID string) (result models.CatalogueCategoryProperty, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	// Pre-validate type UID (required) + unit UID (optional) so an invalid codebook
	// reference surfaces as 400 ErrPatchValidation instead of silent zero-row.
	if fields.Type.UID == "" {
		return result, fmt.Errorf("%w: property type uid is required", ErrPatchValidation)
	}
	if exists, nerr := svc.nodeExists("CatalogueCategoryPropertyType", fields.Type.UID); nerr != nil {
		return result, nerr
	} else if !exists {
		return result, fmt.Errorf("%w: property type not found: %s", ErrPatchValidation, fields.Type.UID)
	}
	if fields.Unit != nil && fields.Unit.UID != "" {
		if exists, nerr := svc.nodeExists("Unit", fields.Unit.UID); nerr != nil {
			return result, nerr
		} else if !exists {
			return result, fmt.Errorf("%w: unit not found: %s", ErrPatchValidation, fields.Unit.UID)
		}
	}

	// Resolve the target order upfront (payload value wins; otherwise max(siblings)+10).
	var order int
	if fields.Order != nil {
		order = *fields.Order
	} else {
		next, qerr := helpers.GetNeo4jSingleRecordSingleValue[int64](session, NextPropertyOrderQuery(categoryUID, groupUID))
		if qerr != nil {
			if errors.Is(qerr, helpers.ERR_NO_ROWS) {
				return result, helpers.ERR_NOT_FOUND
			}
			return result, qerr
		}
		order = int(next)
	}

	newUID := uuid.NewString()
	query := CreateCatalogueCategoryPropertyQuery(categoryUID, groupUID, newUID, userUID, fields, order)
	returnedUID, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return result, helpers.ERR_NOT_FOUND
		}
		return result, err
	}
	if returnedUID == "" {
		return result, helpers.ERR_NOT_FOUND
	}

	fetched, err := svc.fetchCategoryProperty(categoryUID, newUID)
	if err != nil {
		return result, err
	}
	result = toCatalogueCategoryProperty(fetched)
	return result, nil
}

func (svc *CatalogueService) PatchCatalogueCategoryProperty(categoryUID, propertyUID string, fields *models.PatchCatalogueCategoryPropertyFields, userUID string) (result models.CatalogueCategoryProperty, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	original, err := svc.fetchCategoryProperty(categoryUID, propertyUID)
	if err != nil {
		return result, err
	}

	// Pre-validate move target, new type, new unit.
	if fields.GroupUID != nil && *fields.GroupUID != "" && *fields.GroupUID != original.GroupUID {
		exists, nerr := svc.nodeExists("CatalogueCategoryPropertyGroup", *fields.GroupUID)
		if nerr != nil {
			return result, nerr
		}
		if !exists {
			return result, fmt.Errorf("%w: target group not found: %s", ErrPatchValidation, *fields.GroupUID)
		}
		// Ensure the new group belongs to this category (flat URL enforcement).
		gr, gerr := svc.fetchCategoryGroup(categoryUID, *fields.GroupUID)
		if gerr != nil {
			return result, fmt.Errorf("%w: target group not under this category: %s", ErrPatchValidation, *fields.GroupUID)
		}
		_ = gr
	}
	if fields.Type != nil && fields.Type.UID != "" {
		exists, nerr := svc.nodeExists("CatalogueCategoryPropertyType", fields.Type.UID)
		if nerr != nil {
			return result, nerr
		}
		if !exists {
			return result, fmt.Errorf("%w: property type not found: %s", ErrPatchValidation, fields.Type.UID)
		}
	}
	if fields.Unit != nil && fields.Unit.Value != nil && fields.Unit.Value.UID != "" {
		exists, nerr := svc.nodeExists("Unit", fields.Unit.Value.UID)
		if nerr != nil {
			return result, nerr
		}
		if !exists {
			return result, fmt.Errorf("%w: unit not found: %s", ErrPatchValidation, fields.Unit.Value.UID)
		}
	}

	// Lazy order seed for the group the property currently lives in.
	if fields.Order != nil {
		siblings, lerr := helpers.GetNeo4jArrayOfNodes[models.CatalogueCategoryProperty](session, ListCatalogueCategoryPropertiesInGroupQuery(categoryUID, original.GroupUID))
		if lerr != nil && !errors.Is(lerr, helpers.ERR_NO_ROWS) {
			return result, lerr
		}
		allSeeded := true
		for _, p := range siblings {
			if p.Order == nil {
				allSeeded = false
				break
			}
		}
		if !allSeeded {
			_, serr := helpers.WriteNeo4jAndReturnSingleValue[int64](session, SeedCategoryPropertyOrdersQuery(categoryUID, original.GroupUID))
			if serr != nil && !errors.Is(serr, helpers.ERR_NO_ROWS) {
				return result, serr
			}
			if refreshed, ferr := svc.fetchCategoryProperty(categoryUID, propertyUID); ferr == nil {
				original = refreshed
			}
		}
	}

	query := PatchCatalogueCategoryPropertyQuery(categoryUID, propertyUID, fields, &original, userUID)
	returnedUID, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return result, helpers.ERR_NOT_FOUND
		}
		return result, err
	}
	if returnedUID == "" {
		return result, helpers.ERR_NOT_FOUND
	}

	fetched, err := svc.fetchCategoryProperty(categoryUID, propertyUID)
	if err != nil {
		return result, err
	}
	result = toCatalogueCategoryProperty(fetched)
	return result, nil
}

func (svc *CatalogueService) DeleteCatalogueCategoryProperty(categoryUID, propertyUID, userUID string) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	original, err := svc.fetchCategoryProperty(categoryUID, propertyUID)
	if err != nil {
		return err
	}

	query := DeleteCatalogueCategoryPropertyQuery(categoryUID, propertyUID, userUID, original.Name)
	returnedUID, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			still, checkErr := svc.nodeExists("CatalogueCategoryProperty", propertyUID)
			if checkErr != nil {
				return checkErr
			}
			if still {
				return helpers.ERR_DELETE_RELATED_ITEMS
			}
			return helpers.ERR_NOT_FOUND
		}
		return err
	}
	if returnedUID == "" {
		return helpers.ERR_NOT_FOUND
	}
	return nil
}

// toCatalogueCategoryProperty converts the internal group-aware shape into the public
// CatalogueCategoryProperty model used in API responses.
func toCatalogueCategoryProperty(p CategoryPropertyWithGroup) models.CatalogueCategoryProperty {
	out := models.CatalogueCategoryProperty{
		UID:          p.UID,
		Name:         p.Name,
		DefaultValue: p.DefaultValue,
		Order:        p.Order,
		ListOfValues: p.ListOfValues,
	}
	if p.Type != nil {
		out.Type = *p.Type
	}
	if p.Unit != nil {
		out.Unit = &codebookModels.Codebook{UID: p.Unit.UID, Name: p.Unit.Name}
	}
	return out
}

func (svc *CatalogueService) DeleteCatalogueCategoryGroup(categoryUID, groupUID, userUID string) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	group, err := svc.fetchCategoryGroup(categoryUID, groupUID)
	if err != nil {
		return err
	}

	query := DeleteCatalogueCategoryGroupQuery(categoryUID, groupUID, userUID, group.Name)
	returnedUID, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			// Either the group vanished mid-call (404) or the WHERE refs=0 gate filtered
			// it out because a CatalogueItem references one of its properties (409).
			// Distinguish by re-checking existence.
			still, checkErr := svc.nodeExists("CatalogueCategoryPropertyGroup", groupUID)
			if checkErr != nil {
				return checkErr
			}
			if still {
				return helpers.ERR_DELETE_RELATED_ITEMS
			}
			return helpers.ERR_NOT_FOUND
		}
		return err
	}
	if returnedUID == "" {
		return helpers.ERR_NOT_FOUND
	}
	return nil
}

