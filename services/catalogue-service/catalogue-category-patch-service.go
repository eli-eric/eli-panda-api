package catalogueService

import (
	"errors"

	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"

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

