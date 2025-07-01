package roomcardsservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/room-cards-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Room card status constants
const (
	DIRTY_MODE           = "DIRTY_MODE"
	CLEAN_MODE          = "CLEAN_MODE"
	IN_PREPARATION_MODE = "IN_PREPARATION_MODE"
)

type RoomCardsService struct {
	neo4jDriver *neo4j.Driver
}

type IRoomCardsService interface {
	GetLayoutRoomCardInfo(locationCode string) (models.LayoutRoomCard, error)
}

func NewRoomCardsService(driver *neo4j.Driver) IRoomCardsService {
	return &RoomCardsService{neo4jDriver: driver}
}

func (svc *RoomCardsService) GetLayoutRoomCardInfo(locationCode string) (result models.LayoutRoomCard, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	// Use the database query to get room card by location code
	dbQuery := GetRoomCardsByLocationCodeQuery(locationCode)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.LayoutRoomCard](session, dbQuery)

	if err == nil {
		// Add status color logic
		result.StatusColor = getStatusColor(result.Status)
	}

	return result, err
}

// getStatusColor returns the hex color based on the status
func getStatusColor(status string) string {
	switch status {
	case DIRTY_MODE:
		return "#fecaca" 
	case CLEAN_MODE:
		return "#d9f99d"
	case IN_PREPARATION_MODE:
		return "#fdba74" 
	default:
		return "#808080" 
	}
}