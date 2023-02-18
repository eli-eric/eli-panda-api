package helpers

import (
	"fmt"
	"panda/apigateway/ioutils"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func NewNeo4jSession(driver neo4j.Driver) (neo4j.Session, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	var err error
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()
	return session, err
}

func GetNeo4jSingleRecord(session neo4j.Session, cypher string, params map[string]interface{}, returnAlias string) (interface{}, error) {
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return nil, fmt.Errorf("record not found")
		}

		rec, _ := record.Get(returnAlias)
		return rec, nil

	})

	return result, err
}

func GetNeo4jSingleRecordAndMapToStruct[T any](session neo4j.Session, query DatabaseQuery) (result T, err error) {
	resultMap, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(query.Query, query.Parameters)
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return nil, fmt.Errorf("record not found")
		}

		rec, _ := record.Get(query.ReturnAlias)
		return rec, nil

	})

	if err == nil {
		result, err = MapStruct[T](resultMap.(map[string]interface{}))
	}

	return result, err
}

func GetNeo4jArrayOfNodes[T any](session neo4j.Session, query DatabaseQuery) (resultArray []T, err error) {
	// Execute a query in a new Read Transaction
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		result, err := tx.Run(query.Query, query.Parameters)
		if err != nil {
			return nil, err
		}

		// Get a list of Movies from the Result
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var txResults []T
		for _, record := range records {
			itm, _ := record.Get(query.ReturnAlias)
			mappedItem, _ := MapStruct[T](itm.(map[string]interface{}))
			txResults = append(txResults, mappedItem)
		}
		return txResults, nil
	})

	if err == nil {
		resultArray = results.([]T)
	}

	return resultArray, err
}

func GetPaginationResult[T any](data []T, err error) (result PaginationResult[T]) {

	//check for incoming errors
	if err == nil {

		//if there are no data we want to return empty array instead of null
		if data == nil {
			data = []T{}
		}

		result.Data = data
		result.TotalCount = 100

		return result
	}

	return result
}

func ProcessArrayResult[T any](data *[]T, err error) {

	//check for incoming errors
	if err == nil {
		//if there are no data we want to return empty array instead of null
		if data == nil {
			*data = []T{}
		}
	}
}

type PaginationResult[T any] struct {
	TotalCount int `json:"totalCount"`
	Data       []T `json:"data"`
}

type DatabaseQuery struct {
	Query       string
	Parameters  map[string]interface{}
	ReturnAlias string
}

type Pagination struct {
	PageSize int `query:"pageSize"`
	Page     int `query:"page"`
}
