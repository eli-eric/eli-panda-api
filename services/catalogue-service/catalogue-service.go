package catalogueService

import (
	"encoding/json"
	"errors"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

type CatalogueService struct {
	neo4jDriver *neo4j.Driver
	jwtSecret   string
}

type ICatalogueService interface {
	GetCatalogueCategoriesByParentPath(parentPath string) (categories []models.CatalogueCategory, err error)
	GetCatalogueItems(search string, categoryUid string, pageSize int, page int, filering *[]helpers.ColumnFilter, sorting *[]helpers.Sorting) (result helpers.PaginationResult[models.CatalogueItemSimple], err error)
	GetCatalogueItemWithDetailsByUid(uid string) (catalogueItem models.CatalogueItem, err error)
	GetCatalogueCategoryWithDetailsByUid(uid string) (catalogueItem models.CatalogueCategory, err error)
	GetCatalogueCategoryWithDetailsForCopyByUid(uid string) (result models.CatalogueCategory, err error)
	GetCatalogueCategoryImageByUid(uid string) (imageBase64 string, err error)
	GetCatalogueItemImageByUid(uid string) (imageBase64 string, err error)
	UpdateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error)
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
	DeleteCatalogueItem(uid string, userUID string) (err error)
	GetCatalogueItemStatistics(uid string) (result []models.CatalogueStatistics, err error)
	CatalogueItemsOverallStatistics() (result []models.CatalogueStatistics, err error)
	GetCatalogueServiceTypeByUid(uid string) (result models.CatalogueServiceType, err error)
	GetCatalogueServiceTypes() (result []models.CatalogueServiceType, err error)
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

	result, _, err = helpers.GetMultipleNodes[models.CatalogueServiceType](session, 0, 0, "")

	return result, err
}
