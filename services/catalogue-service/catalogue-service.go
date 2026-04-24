package catalogueService

import (
	"encoding/json"
	"errors"
	"fmt"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

type CatalogueService struct {
	neo4jDriver *neo4j.Driver
	jwtSecret   string
}

// ErrPatchValidation marks errors from PATCH reference validation (unknown supplier/
// category/property UID, missing property UID). The handler uses errors.Is to surface 400.
var ErrPatchValidation = errors.New("patch validation failed")

type ICatalogueService interface {
	GetCatalogueCategoriesByParentPath(parentPath string) (categories []models.CatalogueCategory, err error)
	GetCatalogueItems(search string, categoryUid string, pageSize int, page int, filering *[]helpers.ColumnFilter, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.CatalogueItemSimple], err error)
	GetCatalogueItemWithDetailsByUid(uid string) (catalogueItem models.CatalogueItem, err error)
	GetCatalogueCategoryWithDetailsByUid(uid string) (catalogueItem models.CatalogueCategory, err error)
	GetCatalogueCategoryWithDetailsForCopyByUid(uid string) (result models.CatalogueCategory, err error)
	GetCatalogueCategoryImageByUid(uid string) (imageBase64 string, err error)
	GetCatalogueItemImageByUid(uid string) (imageBase64 string, err error)
	UpdateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error)
	PatchCatalogueCategory(uid string, fields *models.PatchCatalogueCategoryFields, userUID string) (result models.CatalogueCategory, err error)
	CreateCatalogueCategoryGroup(categoryUID string, fields *models.CreateCatalogueCategoryGroupFields, userUID string) (result models.CatalogueCategoryPropertyGroup, err error)
	GetCatalogueCategoryGroup(categoryUID, groupUID string) (result models.CatalogueCategoryPropertyGroup, err error)
	PatchCatalogueCategoryGroup(categoryUID, groupUID string, fields *models.PatchCatalogueCategoryGroupFields, userUID string) (result models.CatalogueCategoryPropertyGroup, err error)
	DeleteCatalogueCategoryGroup(categoryUID, groupUID, userUID string) (err error)
	CreateCatalogueCategoryProperty(categoryUID, groupUID string, fields *models.CreateCatalogueCategoryPropertyFields, userUID string) (result models.CatalogueCategoryProperty, err error)
	GetCatalogueCategoryProperty(categoryUID, propertyUID string) (result models.CatalogueCategoryProperty, err error)
	PatchCatalogueCategoryProperty(categoryUID, propertyUID string, fields *models.PatchCatalogueCategoryPropertyFields, userUID string) (result models.CatalogueCategoryProperty, err error)
	DeleteCatalogueCategoryProperty(categoryUID, propertyUID, userUID string) (err error)
	CreateCatalogueCategoryPhysicalProperty(categoryUID string, fields *models.CreateCatalogueCategoryPhysicalPropertyFields, userUID string) (result models.CatalogueCategoryProperty, err error)
	GetCatalogueCategoryPhysicalProperty(categoryUID, propertyUID string) (result models.CatalogueCategoryProperty, err error)
	PatchCatalogueCategoryPhysicalProperty(categoryUID, propertyUID string, fields *models.PatchCatalogueCategoryPhysicalPropertyFields, userUID string) (result models.CatalogueCategoryProperty, err error)
	DeleteCatalogueCategoryPhysicalProperty(categoryUID, propertyUID, userUID string) (err error)
	CreateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error)
	DeleteCatalogueCategory(uid string) (err error)
	GetUnitsCodebook() (result []codebookModels.Codebook, err error)
	GetPropertyTypesCodebook() (result []codebookModels.Codebook, err error)
	CopyCatalogueCategoryRecursive(originalUID string) (newUID string, err error)
	GetCatalogueCategoriesRecursiveByParentUID(parentUID string) (categories []models.CatalogueCategoryTreeItem, err error)
	GetManufacturersCodebook(searchString string, limit int) (result []codebookModels.Codebook, err error)
	GetCatalogueCategoriesCodebook(searchString string, limit int) (result []codebookModels.Codebook, err error)
	GetCatalogueCategoriesCodebookTree(search string) (result []codebookModels.CodebookTreeItem, err error)
	CreateNewCatalogueItem(catalogueItem *models.CatalogueItem, userUID string) (result models.CatalogueItem, err error)
	GetCatalogueCategoryPropertiesByUid(uid string, itemUID *string) (properties []models.CatalogueItemDetail, err error)
	GetCatalogueCategoryPhysicalItemPropertiesByUid(uid string) (properties []models.CatalogueItemDetail, err error)
	UpdateCatalogueItem(catalogueItem *models.CatalogueItem, userUID string) (result models.CatalogueItem, err error)
	PatchCatalogueItem(uid string, fields *models.PatchCatalogueItemFields, userUID string) (result models.CatalogueItem, err error)
	DeleteCatalogueItem(uid string, userUID string) (err error)
	GetCatalogueItemStatistics(uid string) (result []models.CatalogueStatistics, err error)
	CatalogueItemsOverallStatistics() (result []models.CatalogueStatistics, err error)
	GetCatalogueServiceTypeByUid(uid string) (result models.CatalogueServiceType, err error)
	GetCatalogueServiceTypes() (result []models.CatalogueServiceType, err error)
	CreateCatalogueServiceType(catalogueServiceType *models.CatalogueServiceType, userUID string) (result models.CatalogueServiceType, err error)
	UpdateCatalogueServiceType(catalogueServiceType *models.CatalogueServiceType, userUID string) (result models.CatalogueServiceType, err error)
	DeleteCatalogueServiceType(uid string, userUID string) (err error)
	IsCatalogueNumberUnique(catalogueNumber string, excludeUid string) (bool, error)
}

// Create new security service instance
func NewCatalogueService(settings *config.Config, driver *neo4j.Driver) ICatalogueService {

	return &CatalogueService{neo4jDriver: driver, jwtSecret: settings.JwtSecret}
}

func (svc *CatalogueService) GetCatalogueCategoriesByParentPath(parentPath string) (categories []models.CatalogueCategory, err error) {

	// Open a new Session
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get all categories by parent path
	query := CatalogueCategoriesByParentPathQuery(parentPath)
	categories, err = helpers.GetNeo4jArrayOfNodes[models.CatalogueCategory](session, query)

	helpers.ProcessArrayResult(&categories, err)

	return categories, err
}

func (svc *CatalogueService) GetCatalogueCategoriesRecursiveByParentUID(parentUID string) (categories []models.CatalogueCategoryTreeItem, err error) {

	// Open a new Session
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get all categories by parent path
	query := CatalogueSubCategoriesByParentQuery(parentUID)
	categories, err = helpers.GetNeo4jArrayOfNodes[models.CatalogueCategoryTreeItem](session, query)

	helpers.ProcessArrayResult(&categories, err)

	return categories, err
}

func (svc *CatalogueService) GetCatalogueItems(search string, categoryUid string, pageSize int, page int, filering *[]helpers.ColumnFilter, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.CatalogueItemSimple], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get all categories by parent path
	query := CatalogueItemsFiltersPaginationQuery(search, categoryUid, pageSize*(page-1), pageSize, filering, sorting)

	items, err := helpers.GetNeo4jArrayOfNodes[models.CatalogueItemSimple](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CatalogueItemsFiltersTotalCountQuery(search, categoryUid, filering))

	// we have to process the result and set the value for the range type
	for i, item := range items {
		for idt, detail := range item.Details {
			if detail.Property.Type.Code == "range" {
				var rangeValue helpers.RangeFloat64Nullable
				stringData := (detail.Value).(string)
				errJson := json.Unmarshal([]byte(stringData), &rangeValue)
				if errJson == nil {
					items[i].Details[idt].Value = rangeValue
				}
			}
		}
	}

	result = helpers.GetPaginationResult(items, totalCount, err)

	return result, err
}

func (svc *CatalogueService) GetCatalogueCategoryItemsCountRecursive(categoryUID string) (result int64, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetCatalogueCategoryItemsCountRecursiveQuery(categoryUID)
	result, err = helpers.GetNeo4jSingleRecordSingleValue[int64](session, query)

	return result, err
}

func (svc *CatalogueService) GetCatalogueItemWithDetailsByUid(uid string) (result models.CatalogueItem, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueItemWithDetailsByUidQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.CatalogueItem](session, query)

	if err != nil {
		log.Error().Msgf("Error while getting catalogue item with details by uid: %s", uid)
	} else {
		//fitt we got the item with details(but only details/properties with a value)
		//now we need add all properties for the specific category and parent categories
		allProperties, err := svc.GetCatalogueCategoryPropertiesByUid(result.Category.UID, nil)
		if err == nil {
			//we have to iterate on all properties and check if we have this property in the result
			for _, property := range allProperties {
				//check if we have this property in the result
				found := false
				for i, detail := range result.Details {
					if detail.Property.UID == property.Property.UID {
						found = true

						if property.Property.Type.Code == "range" {
							var rangeValue helpers.RangeFloat64Nullable
							stringData := (detail.Value).(string)
							errJson := json.Unmarshal([]byte(stringData), &rangeValue)
							if errJson == nil {
								result.Details[i].Value = rangeValue
							}

						}

						break
					}
				}
				//if we dont have this property in the result we have to add it with empty value
				if !found {
					result.Details = append(result.Details, models.CatalogueItemDetail{Property: property.Property, PropertyGroup: property.PropertyGroup, Value: nil})
				}
			}
		}
	}

	return result, err
}

func (svc *CatalogueService) GetCatalogueCategoryWithDetailsByUid(uid string) (result models.CatalogueCategory, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoryWithDetailsQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.CatalogueCategory](session, query)

	return result, err
}

func (svc *CatalogueService) GetCatalogueCategoryWithDetailsForCopyByUid(uid string) (result models.CatalogueCategory, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoryWithDetailsForCopyQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.CatalogueCategory](session, query)

	return result, err
}

func (svc *CatalogueService) UpdateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get the original record from db to compare because of the delete
	originalItem, err := svc.GetCatalogueCategoryWithDetailsByUid(catalogueCategory.UID)
	if err == nil {
		//update category query
		query := UpdateCatalogueCategoryQuery(catalogueCategory, &originalItem)
		_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	}

	return err
}

func (svc *CatalogueService) PatchCatalogueCategory(uid string, fields *models.PatchCatalogueCategoryFields, userUID string) (result models.CatalogueCategory, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	original, err := svc.GetCatalogueCategoryWithDetailsByUid(uid)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return result, helpers.ERR_NOT_FOUND
		}
		return result, err
	}

	if fields.SystemType != nil && fields.SystemType.Value != nil && fields.SystemType.Value.UID != "" {
		if exists, nerr := svc.nodeExists("SystemType", fields.SystemType.Value.UID); nerr != nil {
			return result, nerr
		} else if !exists {
			return result, fmt.Errorf("%w: systemType not found: %s", ErrPatchValidation, fields.SystemType.Value.UID)
		}
	}

	query := PatchCatalogueCategoryQuery(uid, fields, &original, userUID)
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

	result, err = svc.GetCatalogueCategoryWithDetailsByUid(uid)
	return result, err
}

func (svc *CatalogueService) CreateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//use this hybrid update/create method
	query := UpdateCatalogueCategoryQuery(catalogueCategory, nil)
	_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return err
}

func (svc *CatalogueService) GetCatalogueCategoryImageByUid(uid string) (result string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoryImageByUidQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordSingleValue[string](session, query)

	return result, err
}

func (svc *CatalogueService) GetCatalogueItemImageByUid(uid string) (result string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueItemImageByUidQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordSingleValue[string](session, query)

	return result, err
}

func (svc *CatalogueService) DeleteCatalogueCategory(uid string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//we have to check if this category has some items - if so we cant delete category
	itemsCount, err := svc.GetCatalogueCategoryItemsCountRecursive(uid)
	if err == nil {
		if itemsCount > 0 {
			err = errors.New("category has related items")
		} else {
			query := DeleteCatalogueCategoryByUidQuery(uid)
			err = helpers.WriteNeo4jAndReturnNothing(session, query)
		}
	}

	return err
}

func (svc *CatalogueService) GetUnitsCodebook() (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetUnitsCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *CatalogueService) GetPropertyTypesCodebook() (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetPropertyTypesCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *CatalogueService) CopyCatalogueCategoryRecursive(originalUID string) (newUID string, err error) {

	//get existing category
	category, err := svc.GetCatalogueCategoryWithDetailsForCopyByUid(originalUID)
	if err == nil {
		category.Name += " copy"
		category.Code += "-copy"
		category.Image = ""

		//create new one from existing one
		err = svc.CreateCatalogueCategory(&category)

		// if err is null we can continue....
		if err == nil {
			newUID = category.UID
		}
		//iterate on all sub-categories(recusively) and do the same(copy) for each sub-category
		subCategories, err := svc.GetCatalogueCategoriesRecursiveByParentUID(originalUID)
		if err == nil {
			if len(subCategories) > 0 {
				childs := subCategories[0].Has_subcategory
				svc.createSubChildsRecusive(&childs, newUID)
			}
		}
	}

	return newUID, err
}

func (svc *CatalogueService) createSubChildsRecusive(childs *[]models.CatalogueCategoryTreeItem, newParentUID string) {
	for _, child := range *childs {
		//get existing category
		category, err := svc.GetCatalogueCategoryWithDetailsForCopyByUid(child.UID)
		if err == nil {
			category.ParentUID = newParentUID
			err = svc.CreateCatalogueCategory(&category)
			if err == nil {
				if len(child.Has_subcategory) > 0 {
					svc.createSubChildsRecusive(&child.Has_subcategory, category.UID)
				}
			}
		}
	}
}

func (svc *CatalogueService) GetManufacturersCodebook(searchString string, limit int) (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := ManufacturersForAutocompletQuery(searchString, limit)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *CatalogueService) GetCatalogueCategoriesCodebook(searchString string, limit int) (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoriesForAutocompleteQuery(searchString, limit)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *CatalogueService) CreateNewCatalogueItem(catalogueItem *models.CatalogueItem, userUID string) (result models.CatalogueItem, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//create new item
	query := NewCatalogueItemQuery(catalogueItem, userUID)
	uid, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	if err == nil {
		//get created item
		result, err = svc.GetCatalogueItemWithDetailsByUid(uid)
	}

	return result, err
}

func (svc *CatalogueService) GetCatalogueCategoryPropertiesByUid(uid string, itemUID *string) (properties []models.CatalogueItemDetail, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoryPropertiesQuery(uid)
	properties, err = helpers.GetNeo4jArrayOfNodes[models.CatalogueItemDetail](session, query)

	//get item if we have itemUID
	if itemUID != nil && *itemUID != "" {
		//get item
		item, err := svc.GetCatalogueItemWithDetailsByUid(*itemUID)
		if err == nil {
			//iterate on all properties and set the value from item
			for i, property := range properties {
				for _, detail := range item.Details {
					if detail.Property.UID == property.Property.UID {
						properties[i].Value = detail.Value
						break
					}
				}
			}
		}
	}

	helpers.ProcessArrayResult(&properties, err)

	return properties, err
}

func (svc *CatalogueService) GetCatalogueCategoryPhysicalItemPropertiesByUid(uid string) (properties []models.CatalogueItemDetail, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoryPhysicalItemPropertiesQuery(uid)
	properties, err = helpers.GetNeo4jArrayOfNodes[models.CatalogueItemDetail](session, query)

	helpers.ProcessArrayResult(&properties, err)

	return properties, err
}

func (svc *CatalogueService) UpdateCatalogueItem(catalogueItem *models.CatalogueItem, userUID string) (result models.CatalogueItem, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get the original record from db to compare because of the delete
	originalItem, err := svc.GetCatalogueItemWithDetailsByUid(catalogueItem.Uid)

	if err == nil {

		if catalogueItem.LastUpdateTime != originalItem.LastUpdateTime {
			err = helpers.ERR_CONFLICT
		} else {

			//update category query
			query := UpdateCatalogueItemQuery(catalogueItem, &originalItem, userUID)
			_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

			if err == nil {
				//get updated item
				result, err = svc.GetCatalogueItemWithDetailsByUid(catalogueItem.Uid)
			}
		}
	}

	return result, err
}

func (svc *CatalogueService) PatchCatalogueItem(uid string, fields *models.PatchCatalogueItemFields, userUID string) (result models.CatalogueItem, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	originalItem, err := svc.GetCatalogueItemWithDetailsByUid(uid)
	if err != nil {
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return result, helpers.ERR_NOT_FOUND
		}
		return result, err
	}

	if !fields.LastUpdateTime.Equal(originalItem.LastUpdateTime) {
		return result, helpers.ERR_CONFLICT
	}

	if isCombinedCategoryAndDetailsPatch(fields, originalItem.Category.UID) {
		newCatProps, propErr := svc.GetCatalogueCategoryPropertiesByUid(fields.Category.UID, nil)
		if propErr != nil {
			return result, propErr
		}
		for _, p := range newCatProps {
			if findDetailByPropUID(originalItem.Details, p.Property.UID) == nil {
				originalItem.Details = append(originalItem.Details, models.CatalogueItemDetail{Property: p.Property, PropertyGroup: p.PropertyGroup, Value: nil})
			}
		}
	}

	if err := svc.validatePatchReferences(fields, &originalItem); err != nil {
		return result, err
	}

	query := PatchCatalogueItemQuery(uid, fields, &originalItem, userUID)
	returnedUID, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	if err != nil {
		// Empty result = MATCH on uid+lastUpdateTime yielded no rows, i.e. a concurrent
		// write raced us between the initial read and the PATCH execution.
		if errors.Is(err, helpers.ERR_NO_ROWS) {
			return result, helpers.ERR_CONFLICT
		}
		return result, err
	}
	if returnedUID == "" {
		return result, helpers.ERR_CONFLICT
	}

	result, err = svc.GetCatalogueItemWithDetailsByUid(uid)
	return result, err
}

// isCombinedCategoryAndDetailsPatch reports whether a PATCH swaps category AND also
// supplies details. In that case the service augments originalItem.Details with the new
// category's property set so detail UIDs from the target category pass validation
// and the Cypher builder can source DB-backed Name/Type.Code for audit entries.
func isCombinedCategoryAndDetailsPatch(fields *models.PatchCatalogueItemFields, oldCategoryUID string) bool {
	if fields.Details == nil || fields.Category == nil {
		return false
	}
	if fields.Category.UID == "" {
		return false
	}
	return fields.Category.UID != oldCategoryUID
}

// validatePatchReferences rejects references to non-existent supplier/category/property
// nodes before the main PATCH query runs, so an invalid UID produces a clear 400 instead
// of silently stripping relationships. Returned errors wrap ErrPatchValidation.
func (svc *CatalogueService) validatePatchReferences(fields *models.PatchCatalogueItemFields, originalItem *models.CatalogueItem) error {
	if fields.Supplier != nil && fields.Supplier.Value != nil && fields.Supplier.Value.UID != "" {
		if exists, err := svc.nodeExists("Supplier", fields.Supplier.Value.UID); err != nil {
			return err
		} else if !exists {
			return fmt.Errorf("%w: supplier not found: %s", ErrPatchValidation, fields.Supplier.Value.UID)
		}
	}

	if fields.Category != nil && fields.Category.UID != "" {
		if exists, err := svc.nodeExists("CatalogueCategory", fields.Category.UID); err != nil {
			return err
		} else if !exists {
			return fmt.Errorf("%w: category not found: %s", ErrPatchValidation, fields.Category.UID)
		}
	}

	if fields.Details != nil {
		knownProps := map[string]bool{}
		for _, d := range originalItem.Details {
			knownProps[d.Property.UID] = true
		}
		seen := map[string]bool{}
		for _, d := range *fields.Details {
			if d.Property.UID == "" {
				return fmt.Errorf("%w: detail.property.uid is required", ErrPatchValidation)
			}
			if seen[d.Property.UID] {
				return fmt.Errorf("%w: duplicate property UID in details: %s", ErrPatchValidation, d.Property.UID)
			}
			seen[d.Property.UID] = true
			if !knownProps[d.Property.UID] {
				return fmt.Errorf("%w: property %s is not part of the item's category; change the category first or send a valid property UID", ErrPatchValidation, d.Property.UID)
			}
		}
	}

	return nil
}

func (svc *CatalogueService) nodeExists(label, uid string) (bool, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	q := helpers.DatabaseQuery{
		Query:       fmt.Sprintf("MATCH(n:%s{uid: $uid}) RETURN count(n) as c", label),
		ReturnAlias: "c",
		Parameters:  map[string]interface{}{"uid": uid},
	}
	c, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, q)
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

func (svc *CatalogueService) DeleteCatalogueItem(uid string, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := DeleteCatalogueItemQuery(uid, userUID)
	itemsAffected, err := helpers.WriteNeo4jAndReturnSingleValue[int64](session, query)

	if itemsAffected > 0 {
		err = helpers.ERR_DELETE_RELATED_ITEMS
	}

	return err
}

func (svc *CatalogueService) GetCatalogueCategoriesCodebookTree(search string) (result []codebookModels.CodebookTreeItem, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoriesTreeQuery(search)
	categoriesTreeResult, err := helpers.GetNeo4jArrayOfNodes[codebookModels.CodebookTreeItemCatalogueCategory](session, query)

	if categoriesTreeResult != nil {
		result = make([]codebookModels.CodebookTreeItem, len(categoriesTreeResult))
		for i, item := range categoriesTreeResult {
			result[i] = codebookModels.CodebookTreeItem{UID: item.UID, Name: item.Name, Children: svc.convertCatalogueCategoriesTreeToCodebookTree(item.Children)}
		}
	} else {
		helpers.ProcessArrayResult(&categoriesTreeResult, err)
	}

	return result, err
}

// convert []codebookModels.CodebookTreeItemCatalogueCategory to []codebookModels.CodebookTreeItem recursively
func (svc *CatalogueService) convertCatalogueCategoriesTreeToCodebookTree(categoriesTree []codebookModels.CodebookTreeItemCatalogueCategory) (result []codebookModels.CodebookTreeItem) {

	result = make([]codebookModels.CodebookTreeItem, len(categoriesTree))
	for i, item := range categoriesTree {
		result[i] = codebookModels.CodebookTreeItem{UID: item.UID, Name: item.Name}
		if item.Children != nil {
			result[i].Children = svc.convertCatalogueCategoriesTreeToCodebookTree(item.Children)
		}
	}

	return result
}

func (svc *CatalogueService) GetCatalogueItemStatistics(uid string) (result []models.CatalogueStatistics, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueItemStatisticsQuery(uid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.CatalogueStatistics](session, query)

	return result, err
}

func (svc *CatalogueService) CatalogueItemsOverallStatistics() (result []models.CatalogueStatistics, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueItemsOverallStatisticsQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[models.CatalogueStatistics](session, query)

	return result, err
}

func (svc *CatalogueService) GetCatalogueServiceTypeByUid(uid string) (result models.CatalogueServiceType, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	result.Uid = uid

	err = helpers.GetSingleNode(session, &result)

	return result, err
}

func (svc *CatalogueService) GetCatalogueServiceTypes() (result []models.CatalogueServiceType, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	result, _, err = helpers.GetMultipleNodes[models.CatalogueServiceType](session, 0, 100, "")

	return result, err
}

func (svc *CatalogueService) CreateCatalogueServiceType(catalogueServiceType *models.CatalogueServiceType, userUID string) (result models.CatalogueServiceType, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	// Ensure UID is set
	if catalogueServiceType.Uid == "" {
		catalogueServiceType.Uid = uuid.New().String()
	}

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MERGE (n:CatalogueServiceType {uid: $uid}) SET n.updatedAt = datetime() WITH n `
	updateQuery.Parameters["uid"] = catalogueServiceType.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *catalogueServiceType, models.CatalogueServiceType{}, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(catalogueServiceType.Uid, "CREATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updateQuery,
		historyLog)

	return *catalogueServiceType, err
}

func (svc *CatalogueService) UpdateCatalogueServiceType(catalogueServiceType *models.CatalogueServiceType, userUID string) (result models.CatalogueServiceType, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	oldCatalogueServiceType, err := svc.GetCatalogueServiceTypeByUid(catalogueServiceType.Uid)

	if err != nil {
		return result, err
	}

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MERGE (n:CatalogueServiceType {uid: $uid}) SET n.updatedAt = datetime() WITH n `
	updateQuery.Parameters["uid"] = catalogueServiceType.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *catalogueServiceType, oldCatalogueServiceType, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(catalogueServiceType.Uid, "UPDATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updateQuery,
		historyLog)

	return *catalogueServiceType, err
}

func (svc *CatalogueService) DeleteCatalogueServiceType(uid string, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		helpers.SoftDeleteNodeQuery(uid),
		helpers.HistoryLogQuery(uid, "DELETE", userUID))

	return err
}

func (svc *CatalogueService) IsCatalogueNumberUnique(catalogueNumber string, excludeUid string) (bool, error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := helpers.DatabaseQuery{}
	query.Parameters = make(map[string]interface{})
	query.Query = `
		MATCH (ci:CatalogueItem)
		WHERE ci.catalogueNumber = $catalogueNumber
		  AND coalesce(ci.deleted, false) = false
		  AND ($excludeUid = '' OR ci.uid <> $excludeUid)
		RETURN count(ci) as count
	`
	query.Parameters["catalogueNumber"] = catalogueNumber
	query.Parameters["excludeUid"] = excludeUid
	query.ReturnAlias = "count"

	count, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, query)

	if err != nil {
		return false, err
	}

	return count == 0, nil
}
