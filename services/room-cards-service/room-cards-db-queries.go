package roomcardsservice

import (
	"panda/apigateway/helpers"
)

// GetRoomCardsByLocationCodeQuery returns room cards for a specific location code
func GetRoomCardsByLocationCodeQuery(locationCode string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (location:Location {code: $locationCode})
	MATCH (location)-[:HAS_ROOM_CARD]->(roomCard:RoomCard)
	WHERE (roomCard.deleted IS NULL OR roomCard.deleted = false)
	OPTIONAL MATCH (roomCard)-[:HAS_OPERATIONAL_STATE]->(os:OperationalState)
	RETURN {
		uid: roomCard.uid,
		name: roomCard.name,
		location: {uid: location.uid, name: location.name, code: location.code},
		operationalState: case when os is not null then {uid: os.uid, name: os.name, code: os.code} else null end,
		purityClass: roomCard.purityClass,
		status: roomCard.status
	} as roomCard
	ORDER BY roomCard.name
	LIMIT 1`
	
	result.Parameters = make(map[string]interface{})
	result.Parameters["locationCode"] = locationCode
	result.ReturnAlias = "roomCard"

	return result
}

// GetSingleRoomCardQuery returns a single room card by UID
func GetSingleRoomCardQuery(uid string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH (roomCard:RoomCard {uid: $uid})
	MATCH (location:Location)-[:HAS_ROOM_CARD]->(roomCard)
	WHERE (roomCard.deleted IS NULL OR roomCard.deleted = false)
	RETURN {
		uid: roomCard.uid,
		name: roomCard.name,
		location: {uid: location.uid, name: location.name, code: location.code},
		purityClass: roomCard.purityClass,
		status: roomCard.status
	} as roomCard`
	
	result.Parameters = make(map[string]interface{})
	result.Parameters["uid"] = uid
	result.ReturnAlias = "roomCard"

	return result
}
