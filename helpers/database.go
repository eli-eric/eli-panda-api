package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"panda/apigateway/ioutils"
	"panda/apigateway/services/codebook-service/models"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

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

// not used yet - its for the future
func CreateOrUpdateNodeQuery(node interface{}) (DatabaseQuery, error) {
	val := reflect.ValueOf(node)
	typ := reflect.TypeOf(node)

	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return DatabaseQuery{}, fmt.Errorf("expected a struct, got %s", typ.Kind())
	}

	// Build Cypher query and parameters
	var fields []string
	params := map[string]interface{}{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		//jsonTag := field.Tag.Get("json")
		neo4jTag := field.Tag.Get("neo4j")

		// Handle only `prop` fields for nodes
		if strings.HasPrefix(neo4jTag, "prop,") {
			propName := strings.TrimPrefix(neo4jTag, "prop,")
			fields = append(fields, fmt.Sprintf("%s: $%s", propName, propName))
			params[propName] = value.Interface()
		}

		// handle key property
		if strings.HasPrefix(neo4jTag, "key") {
			propName := strings.TrimPrefix(neo4jTag, "key,")
			params[propName] = value.Interface()
		}
	}

	// Create Cypher query
	query := fmt.Sprintf(`
	MERGE (n:%s {uid: $uid})
	SET n += {%s}
	RETURN n
	`, typ.Name(), strings.Join(fields, ", "))

	// Run the query
	return DatabaseQuery{
		Query:       query,
		Parameters:  params,
		ReturnAlias: "n",
	}, nil
}

// get single node by uid
func GetSingleNode(session neo4j.Session, node interface{}) (err error) {

	typ := reflect.TypeOf(node)
	uid := ""

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		uid = reflect.ValueOf(node).Elem().FieldByName("Uid").String()
	} else {
		uid = reflect.ValueOf(node).FieldByName("Uid").String()
	}

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %s", typ.Kind())
	}

	// Build Cypher query and parameters
	var fields []string
	var optionalMatchQueries []string

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		//jsonTag := field.Tag.Get("json")
		neo4jTag := field.Tag.Get("neo4j")

		// Handle only `prop` fields for nodes
		if strings.HasPrefix(neo4jTag, "prop,") {
			propName := strings.TrimPrefix(neo4jTag, "prop,")
			fields = append(fields, fmt.Sprintf("%s: n.%s", propName, propName))

		}

		// handle key property
		if strings.HasPrefix(neo4jTag, "key") {
			propName := strings.TrimPrefix(neo4jTag, "key,")
			fields = append(fields, fmt.Sprintf("%s: n.%s", propName, propName))
		}

		// handle optional match query and fields
		if strings.HasPrefix(neo4jTag, "rel,") {
			tagParts := strings.Split(neo4jTag, ",")
			if len(tagParts) < 5 {
				return fmt.Errorf("invalid 'rel' tag format: %s", neo4jTag)
			}

			// Extract relationship details
			targetNodeType := tagParts[1]
			relationshipType := tagParts[2]
			targetAlias := tagParts[4]

			optionalMatchQueries = append(optionalMatchQueries, fmt.Sprintf("OPTIONAL MATCH (n)-[:%s]->(%s:%s) ", relationshipType, targetAlias, targetNodeType))
			fields = append(fields, fmt.Sprintf("%s: case when %s is NOT NULL THEN {uid: %s.uid, name: %s.name} ELSE null END", targetAlias, targetAlias, targetAlias, targetAlias))
		}

	}

	// Create Cypher query
	query := fmt.Sprintf(`
	MATCH (n:%s {uid: $uid})
	%s
	RETURN {
			%s
	} as n
	`, typ.Name(),
		strings.Join(optionalMatchQueries, " "),
		strings.Join(fields, ","))

	// Run the query
	resultMap, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		result, err := tx.Run(query, map[string]interface{}{"uid": uid})
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		rec, _ := record.Get("n")
		return rec, nil
	})

	if err == nil {
		err = MapStructToInterface(resultMap.(map[string]interface{}), node)
	}

	return err
}

func GetMultipleNodes[T any](session neo4j.Session, skip, limit int, searchText string) (result []T, totalCount int64, err error) {

	dbQuery := DatabaseQuery{}
	dbQuery.Parameters = make(map[string]interface{})

	typ := reflect.TypeOf(result)

	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return result, totalCount, fmt.Errorf("expected a struct, got %s", typ.Kind())
	}

	// Build Cypher query and parameters
	var fields []string
	searchFields := make(map[string]string) // key is the field name, value is the field type
	var optionalMatchQueries []string

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		//jsonTag := field.Tag.Get("json")
		neo4jTag := field.Tag.Get("neo4j")

		// Handle only `prop` fields for nodes
		if strings.HasPrefix(neo4jTag, "prop,") {
			propName := strings.TrimPrefix(neo4jTag, "prop,")
			fields = append(fields, fmt.Sprintf("%s: n.%s", propName, propName))
			searchFields[propName] = field.Type.String()
		}

		// handle key property
		if strings.HasPrefix(neo4jTag, "key") {
			propName := strings.TrimPrefix(neo4jTag, "key,")
			fields = append(fields, fmt.Sprintf("%s: n.%s", propName, propName))
			searchFields[propName] = field.Type.String()
		}

		// handle optional match query and fields
		if strings.HasPrefix(neo4jTag, "rel,") {
			tagParts := strings.Split(neo4jTag, ",")
			if len(tagParts) < 5 {
				return result, totalCount, fmt.Errorf("invalid 'rel' tag format: %s", neo4jTag)
			}

			// Extract relationship details
			targetNodeType := tagParts[1]
			relationshipType := tagParts[2]
			targetAlias := tagParts[4]

			optionalMatchQueries = append(optionalMatchQueries, fmt.Sprintf("OPTIONAL MATCH (n)-[:%s]->(%s:%s) ", relationshipType, targetAlias, targetNodeType))
			fields = append(fields, fmt.Sprintf("%s: case when %s is NOT NULL THEN {uid: %s.uid, name: %s.name} ELSE null END", targetAlias, targetAlias, targetAlias, targetAlias))
			searchFields[targetAlias] = "codebook"
		}

	}

	// create search query part
	searchQuery := ""
	if searchText != "" {
		dbQuery.Parameters["search"] = searchText
		// foreach field in the struct
		for propName, propType := range searchFields {
			if propType == "string" || propType == "*string" {
				if searchQuery == "" {
					searchQuery += fmt.Sprintf(" AND (toLower(n.%s) CONTAINS toLower($search)", propName)
				} else {
					searchQuery += fmt.Sprintf(" OR toLower(n.%s) CONTAINS toLower($search) ", propName)
				}
			}
		}
		searchQuery += ")"
	}

	// Create Cypher query
	query := fmt.Sprintf(`
	MATCH (n:%s) WHERE (n.deleted IS NULL OR n.deleted = false)
	%s
	%s		
	RETURN {
			%s			
	} as n ORDER BY n.updatedAt DESC SKIP %d LIMIT %d
	`,
		typ.Name(),
		searchQuery,
		strings.Join(optionalMatchQueries, " "),
		strings.Join(fields, ","),
		skip, limit)

	dbQuery.Query = query
	dbQuery.ReturnAlias = "n"
	// Run the query
	result, err = GetNeo4jArrayOfNodes[T](session, dbQuery)

	if err != nil {
		return result, totalCount, err
	}

	// Create Cypher query
	query = fmt.Sprintf(`
	MATCH (n:%s) WHERE (n.deleted IS NULL OR n.deleted = false)
	%s
	%s	
	RETURN count(n) as totalCount
	`,
		typ.Name(),
		searchQuery,
		strings.Join(optionalMatchQueries, " "),
	)

	dbQuery.Query = query
	dbQuery.ReturnAlias = "totalCount"
	// Run the query
	totalCount, err = GetNeo4jSingleRecordSingleValue[int64](session, dbQuery)

	return result, totalCount, err
}

func DeleteNodeQuery(nodeUID string) (result DatabaseQuery) {
	result.Query = `
	MATCH (n {uid:$uid})
	DETACH DELETE n`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = nodeUID

	return result
}

func SoftDeleteNodeQuery(nodeUID string) (result DatabaseQuery) {
	result.Query = `
	MATCH (n {uid:$uid})
	SET n.deleted = true
	RETURN n`
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = nodeUID

	return result
}

func HistoryLogQuery(uid, action, userUID string) (result DatabaseQuery) {
	result.Query = `
	MATCH(u:User{uid:$userUID})
	MATCH(s{uid:$uid})
	with u,s
	CREATE(s)-[:WAS_UPDATED_BY{at: datetime(), action: $action}]->(u)
	RETURN true as result`

	result.ReturnAlias = "result"
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.Parameters["action"] = action
	result.Parameters["userUID"] = userUID

	return result
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

func WriteNeo4jReturnSingleRecordAndMapToStruct[T any](session neo4j.Session, query DatabaseQuery) (result T, err error) {

	resultMap, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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

// write transaction with multiple queries
func WriteNeo4jAndReturnNothingMultipleQueries(session neo4j.Session, queries ...DatabaseQuery) (err error) {

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		for _, query := range queries {
			_, err := tx.Run(query.Query, query.Parameters)

			if err != nil {
				log.Info().Msg(err.Error())
				return nil, err
			}
		}

		return nil, nil
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
			if itm != nil {
				mappedItem, _ := MapStruct[T](itm.(map[string]interface{}))
				txResults = append(txResults, mappedItem)
			}
		}
		return txResults, nil
	})

	if err == nil {
		resultArray = results.([]T)
	}

	return resultArray, err
}

// the objects has to be the same type
// the object has to have neo4j struct Tags - look at the SystemForm for example
// there has to be existing update query (dbQuery param) with one strict alias for the updated node - updateNodeAlias
func AutoResolveObjectToUpdateQuery(dbQuery *DatabaseQuery, newObject any, originObject any, updateNodeAlias string) {

	newObj := reflect.TypeOf(newObject)
	oldObj := reflect.TypeOf(originObject)
	newValObj := reflect.ValueOf(newObject)
	oldValObj := reflect.ValueOf(originObject)

	if newObj == oldObj {
		for i := 0; i < newObj.NumField(); i++ {

			newField := newObj.Field(i)
			oldField := oldObj.Field(i)
			neo4jTags := strings.Split(newField.Tag.Get("neo4j"), ",")
			fieldType := newField.Type.String()

			if len(neo4jTags) > 0 {

				neo4jPropType := neo4jTags[0]
				if neo4jPropType == "prop" {
					neo4jPropName := neo4jTags[1]

					if fieldType == "string" {
						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name).String()
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name).String()

						if newValue != oldValue {
							dbQuery.Parameters[neo4jPropName] = strings.TrimSpace(newValue)
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}

					} else if fieldType == "*string" {
						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name)

						if newValue.IsNil() {
							dbQuery.Parameters[neo4jPropName] = nil
						} else if oldValue.IsNil() {
							dbQuery.Parameters[neo4jPropName] = strings.TrimSpace(newValue.Elem().String())
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						} else if oldValue.Elem().String() != newValue.Elem().String() {
							dbQuery.Parameters[neo4jPropName] = strings.TrimSpace(newValue.Elem().String())
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}

					} else if fieldType == "time.Time" {
						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name).Interface().(time.Time)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name).Interface().(time.Time)

						if newValue != oldValue {
							dbQuery.Parameters[neo4jPropName] = newValue.Local()
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}
					} else if fieldType == "*time.Time" {

						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name)

						if newValue.IsNil() {
							dbQuery.Parameters[neo4jPropName] = nil
						} else if oldValue.IsNil() {
							dbQuery.Parameters[neo4jPropName] = newValue.Elem().Interface().(time.Time).Local()
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						} else if oldValue.Elem().Interface().(time.Time) != newValue.Elem().Interface().(time.Time) {
							dbQuery.Parameters[neo4jPropName] = newValue.Elem().Interface().(time.Time).Local()
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}
					} else if fieldType == "int" {
						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name).Int()
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name).Int()

						if newValue != oldValue {
							dbQuery.Parameters[neo4jPropName] = newValue
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}
					} else if fieldType == "*int" {

						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name)

						if newValue.IsNil() && !oldValue.IsNil() {
							dbQuery.Parameters[neo4jPropName] = nil
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						} else if !newValue.IsNil() && oldValue.IsNil() {
							dbQuery.Parameters[neo4jPropName] = newValue.Elem().Int()
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						} else if !newValue.IsNil() && !oldValue.IsNil() && oldValue.Elem().Int() != newValue.Elem().Int() {
							dbQuery.Parameters[neo4jPropName] = newValue.Elem().Int()
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}
					} else if fieldType == "float64" {
						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name).Float()
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name).Float()

						if newValue != oldValue {
							dbQuery.Parameters[neo4jPropName] = newValue
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}
					} else if fieldType == "*float64" {

						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name)

						if newValue.IsNil() && !oldValue.IsNil() {
							dbQuery.Parameters[neo4jPropName] = nil
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						} else if !newValue.IsNil() && oldValue.IsNil() {
							dbQuery.Parameters[neo4jPropName] = newValue.Elem().Float()
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						} else if !newValue.IsNil() && !oldValue.IsNil() && oldValue.Elem().Float() != newValue.Elem().Float() {
							dbQuery.Parameters[neo4jPropName] = newValue.Elem().Float()
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}
						// else if array of strings
					} else if fieldType == "[]string" {

						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name).Interface().([]string)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name).Interface().([]string)

						if !reflect.DeepEqual(newValue, oldValue) {
							dbQuery.Parameters[neo4jPropName] = newValue
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v SET %[1]v.%[2]v=$%[2]v `, updateNodeAlias, neo4jPropName)
						}
					}

				} else if neo4jPropType == "rel" {
					neo4jLabel := neo4jTags[1]
					neo4jRelationType := neo4jTags[2]

					if fieldType == "*string" {

						neo4jID := neo4jTags[3]
						neo4jAlias := neo4jTags[4]

						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name)

						if !newValue.IsNil() && newValue.Elem().String() != "" && oldValue.IsNil() {
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[2]v:%[3]v{%[4]v:$%[5]v}) MERGE(%[1]v)-[:%[6]v]->(%[2]v) `, updateNodeAlias, neo4jAlias, neo4jLabel, neo4jID, newField.Name, neo4jRelationType)
							dbQuery.Parameters[newField.Name] = newValue.Elem().String()
						} else if !newValue.IsNil() && newValue.Elem().String() != "" && !oldValue.IsNil() && newValue.Elem().String() != oldValue.Elem().String() {
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[1]v)-[r%[2]v:%[3]v]->(%[2]v) delete r%[2]v `, updateNodeAlias, neo4jAlias, neo4jRelationType)
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[2]v:%[3]v{%[4]v:$%[5]v}) MERGE(%[1]v)-[:%[6]v]->(%[2]v) `, updateNodeAlias, neo4jAlias, neo4jLabel, neo4jID, newField.Name, neo4jRelationType)
							dbQuery.Parameters[newField.Name] = newValue.Elem().String()
						} else if (newValue.IsNil() || newValue.Elem().String() == "") && !oldValue.IsNil() {
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[1]v)-[r%[2]v:%[3]v]->(%[2]v) delete r%[2]v `, updateNodeAlias, neo4jAlias, neo4jRelationType)
						}
					} else if fieldType == "*models.Codebook" {

						neo4jAlias := neo4jTags[4]
						neo4jID := "uid"

						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name)
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name)

						if !newValue.IsNil() && oldValue.IsNil() {
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[2]v:%[3]v{%[4]v:$%[5]v}) MERGE(%[1]v)-[:%[6]v]->(%[2]v) `, updateNodeAlias, neo4jAlias, neo4jLabel, neo4jID, newField.Name, neo4jRelationType)
							dbQuery.Parameters[newField.Name] = newValue.Elem().Interface().(models.Codebook).UID
						} else if !newValue.IsNil() && !oldValue.IsNil() && newValue.Elem().Interface().(models.Codebook).UID != oldValue.Elem().Interface().(models.Codebook).UID {
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[1]v)-[r%[2]v:%[3]v]->(%[2]v) delete r%[2]v `, updateNodeAlias, neo4jAlias, neo4jRelationType)
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[2]v:%[3]v{%[4]v:$%[5]v}) MERGE(%[1]v)-[:%[6]v]->(%[2]v) `, updateNodeAlias, neo4jAlias, neo4jLabel, neo4jID, newField.Name, neo4jRelationType)
							dbQuery.Parameters[newField.Name] = newValue.Elem().Interface().(models.Codebook).UID
						} else if newValue.IsNil() && !oldValue.IsNil() {
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[1]v)-[r%[2]v:%[3]v]->(%[2]v) delete r%[2]v `, updateNodeAlias, neo4jAlias, neo4jRelationType)
						}
					} else if fieldType == "models.Codebook" {

						neo4jAlias := neo4jTags[4]
						neo4jID := "uid"

						newValue := reflect.Indirect(newValObj).FieldByName(newField.Name).Interface().(models.Codebook).UID
						oldValue := reflect.Indirect(oldValObj).FieldByName(oldField.Name).Interface().(models.Codebook).UID

						if newValue != oldValue {
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[1]v)-[r%[2]v:%[3]v]->(%[2]v) delete r%[2]v `, updateNodeAlias, neo4jAlias, neo4jRelationType)
							dbQuery.Query += fmt.Sprintf(`WITH %[1]v MATCH(%[2]v:%[3]v{%[4]v:$%[5]v}) MERGE(%[1]v)-[:%[6]v]->(%[2]v) `, updateNodeAlias, neo4jAlias, neo4jLabel, neo4jID, newField.Name, neo4jRelationType)
							dbQuery.Parameters[newField.Name] = newValue
						}
					}
				}
			}
		}
	}
}

func LogDBHistory(session neo4j.Session, objectUID string, originObject any, newObject any, userUID string, action string) (uid string, err error) {

	originSystemJSON := ""

	if originObject != nil {

		originSystemBytes, err := json.Marshal(originObject)
		if err != nil {
			log.Info().Msg(err.Error())
			return uid, err
		}
		originSystemJSON = string(originSystemBytes)
	}

	newSystemBytes, err := json.Marshal(newObject)
	if err != nil {
		log.Info().Msg(err.Error())
		return uid, err
	}

	uid, err = WriteNeo4jAndReturnSingleValue[string](session, logHistoryQuery(objectUID, originSystemJSON, string(newSystemBytes), userUID, action))

	if err != nil {
		log.Info().Msg(err.Error())
	}

	return uid, err
}

func logHistoryQuery(objectUID string, originObjectJSON string, newObjectJSON string, userUID string, action string) (result DatabaseQuery) {

	result.Query = `
	MATCH(u:User{uid:$userUID})
	MATCH(s{uid:$objectUID})
	with u,s
	CREATE(h:History{uid: $uid})
	SET h.timestamp = datetime(), h.objectUID = $objectUID, h.originObject = $originObjectJSON, h.newObject = $newObjectJSON, h.action = $action, h.objectType = labels(s)[0]
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

func GetFullTextSearchString(searchString string) (result string) {

	searchString = strings.TrimSpace(searchString)
	searchString = strings.ReplaceAll(searchString, "/", " ")

	if searchString != "" {
		searchStrings := strings.Split(searchString, " ")
		for i, search := range searchStrings {
			if i == 0 {
				result += "*" + search + "*"
			} else {
				result += " +*" + search + "*"
			}
		}
	}

	return result
}

func GetFilterValue[T any](filters *[]ColumnFilter, filterID string) (result *T) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := f.Value.(T)
			return &value
		}
	}

	return nil

}

func GetFilterValueString(filters *[]ColumnFilter, filterID string) (result *string) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := strings.TrimSpace(f.Value.(string))
			return &value
		}
	}

	return nil
}

func GetFilterValueInt(filters *[]ColumnFilter, filterID string) (result *int) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := f.Value.(int)
			return &value
		}
	}

	return nil
}

func GetFilterValueBool(filters *[]ColumnFilter, filterID string) (result *bool) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := f.Value.(bool)
			return &value
		}
	}

	return nil
}

func GetFilterValueFloat64(filters *[]ColumnFilter, filterID string) (result *float64) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := f.Value.(float64)
			return &value
		}
	}

	return nil
}

func GetFilterValueTime(filters *[]ColumnFilter, filterID string) (result *time.Time) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := f.Value.(time.Time)
			return &value
		}
	}

	return nil
}

func GetFilterValueCodebook(filters *[]ColumnFilter, filterID string) (result *models.Codebook) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := f.Value.(map[string]interface{})
			uid := value["uid"].(string)
			name := value["name"].(string)
			return &models.Codebook{UID: uid, Name: name}
		}
	}

	return nil
}

func GetFilterValueListString(filters *[]ColumnFilter, filterID string) (result *[]string) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := f.Value.([]interface{})
			var result []string
			for _, v := range value {
				result = append(result, v.(string))
			}
			return &result
		}
	}

	return nil
}

func GetFilterValueRangeFloat64(filters *[]ColumnFilter, filterID string) (result *RangeFloat64Nullable) {

	if filters == nil {
		return nil
	}

	for _, f := range *filters {
		if f.Id == filterID {
			value := f.Value.(map[string]interface{})

			minValue := value["min"]
			maxValue := value["max"]
			var result = RangeFloat64Nullable{}

			if minValue != nil {
				min := minValue.(float64)
				result.Min = &min
			}

			if maxValue != nil {
				max := maxValue.(float64)
				result.Max = &max
			}

			return &result

		}
	}

	return nil
}

type RangeFloat64Nullable struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
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

type Sorting struct {
	ID   string
	DESC bool
}

func GetSortingDirectionString(desc bool) (result string) {
	if desc {
		return "DESC"
	}
	return "ASC"
}

type Filter struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type ColumnFilter struct {
	Id       string `json:"id"`
	Value    any    `json:"value"`
	Type     string `json:"type"`
	PropType string `json:"propType"` // could be "CATAOGUE_ITEM" or "PHYISICAL_ITEM"
}

type RelatedNodeLabelAmount struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type ConflictErrorResponse struct {
	ErrorMessage string                   `json:"errorMessage"`
	RelatedNodes []RelatedNodeLabelAmount `json:"relatedNodes,omitempty"`
}

const DB_LOG_CREATE string = "CREATE"
const DB_LOG_UPDATE string = "UPDATE"
const DB_LOG_DELETE string = "DELETE"

const CATALOGUE_CATEGORY_GENERAL_UID string = "97598f04-948f-4da5-95b6-b2a44e0076db"

const FACILITY_CODE_BEAMLINES = "B"
const FACILITY_CODE_ALPS = "A"
const FACILITY_CODE_NP = "N"

var ERR_INVALID_INPUT = errors.New("INVALID_INPUT")
var ERR_UNAUTHORIZED = errors.New("UNAUTHORIZED")
var ERR_NOT_FOUND = errors.New("NOT_FOUND")
var ERR_DELETE_RELATED_ITEMS = errors.New("DELETE_NOT_POSSIBLE_RELATED_ITEMS")
var ERR_CONFLICT = errors.New("CONFLICT")
var ERR_DUPLICATE_SYSTEM_CODE = errors.New("DUPLICATE_SYSTEM_CODE")
var ERR_MISSING_REQUIRED_FIELD = errors.New("MISSING_REQUIRED_FIELD")
