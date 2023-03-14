package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"panda/apigateway/ioutils"

	"github.com/google/uuid"
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

func GetNeo4jSingleRecordAndMapToStruct[T any](session neo4j.Session, query DatabaseQuery) (result T, err error) {
	resultMap, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(query.Query, query.Parameters)
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		rec, _ := record.Get(query.ReturnAlias)
		return rec, nil

	})

	if err == nil {
		result, err = MapStruct[T](resultMap.(map[string]interface{}))
	}

	return result, err
}

func GetNeo4jSingleRecordSingleValue[T any](session neo4j.Session, query DatabaseQuery) (result T, err error) {
	resultValue, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(query.Query, query.Parameters)
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		rec, _ := record.Get(query.ReturnAlias)
		return rec, nil

	})

	if err == nil {
		if resultValue != nil {
			result = resultValue.(T)
		}
	}

	return result, err
}

func WriteNeo4jAndReturnSingleValue[T any](session neo4j.Session, query DatabaseQuery) (result T, err error) {
	resultValue, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(query.Query, query.Parameters)
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		rec, _ := record.Get(query.ReturnAlias)
		return rec, nil

	})

	if err == nil {
		result = resultValue.(T)
	}

	return result, err
}

func WriteNeo4jAndReturnNothing(session neo4j.Session, query DatabaseQuery) (err error) {
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(query.Query, query.Parameters)

		return nil, err
	})

	return err
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

func LogDBHistory(session neo4j.Session, objectUID string, originObject any, newObject any, userUID string, action string) (uid string, err error) {

	originSystemJSON := ""

	if originObject != nil {

		originSystemBytes, err := json.Marshal(originObject)
		if err != nil {
			log.Println(err.Error())
			return uid, err
		}
		originSystemJSON = string(originSystemBytes)
	}

	newSystemBytes, err := json.Marshal(newObject)
	if err != nil {
		log.Println(err.Error())
		return uid, err
	}

	uid, err = WriteNeo4jAndReturnSingleValue[string](session, logHistoryQuery(objectUID, originSystemJSON, string(newSystemBytes), userUID, action))

	if err != nil {
		log.Println(err.Error())
	}

	return uid, err
}

func logHistoryQuery(objectUID string, originObjectJSON string, newObjectJSON string, userUID string, action string) (result DatabaseQuery) {

	result.Query = `
	MATCH(u:User{uid:$userUID})
	MATCH(s{uid:$objectUID})
	with u,s
	CREATE(h:History{uid: $uid})
	SET h.timestamp = datetime(), h.originObject = $originObjectJSON, h.newObject = $newObjectJSON, h.action = $action, h.objectType = labels(s)[0]
	with u,s,h
	CREATE(s)-[:HAS_HISTORY]->(h)
	CREATE(h)-[:DONE_BY_USER]->(u)	
	RETURN h.uid as result`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["objectUID"] = objectUID
	result.Parameters["originObjectJSON"] = originObjectJSON
	result.Parameters["newObjectJSON"] = newObjectJSON
	result.Parameters["userUID"] = userUID
	result.Parameters["uid"] = uuid.NewString()
	result.Parameters["action"] = action
	return result
}

func GetPaginationResult[T any](data []T, totalCount int64, err error) (result PaginationResult[T]) {

	//check for incoming errors
	if err == nil {

		//if there are no data we want to return empty array instead of null
		if data == nil {
			data = []T{}
		}

		result.Data = data
		result.TotalCount = totalCount

		return result
	}

	return result
}

func ProcessArrayResult[T any](data *[]T, err error) {

	//check for incoming errors
	if err == nil {
		//if there are no data we want to return empty array instead of null
		if *data == nil {
			*data = []T{}
		}
	}
}

type PaginationResult[T any] struct {
	TotalCount int64 `json:"totalCount"`
	Data       []T   `json:"data"`
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
