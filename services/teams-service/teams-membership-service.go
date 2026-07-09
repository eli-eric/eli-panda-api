package teamsservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/teams-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func (svc *TeamsService) AddMembers(teamUID, facilityCode, userUID string, userUids []string) (result models.TeamDetail, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	if _, err = svc.getTeam(session, teamUID, facilityCode); err != nil {
		return result, err
	}

	userUids = dedupeNonEmpty(userUids)
	if len(userUids) == 0 {
		return result, ErrValidation
	}

	if err = svc.validateMembers(session, facilityCode, userUids); err != nil {
		return result, err
	}

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		MergeTeamMembersQuery(teamUID, facilityCode, userUids),
		TeamMembershipAuditQuery(teamUID, userUID, buildMembershipChanges(userUids, nil)))
	if err != nil {
		return result, err
	}

	return svc.getTeamDetail(session, teamUID, facilityCode)
}

func (svc *TeamsService) RemoveMember(teamUID, facilityCode, userUID, memberUID string) (result models.TeamDetail, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	detail, err := svc.getTeamDetail(session, teamUID, facilityCode)
	if err != nil {
		return result, err
	}

	// idempotent: if the user is not a member, no write / no audit noise
	if !containsMember(detail.Members, memberUID) {
		return detail, nil
	}

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		RemoveTeamMembersQuery(teamUID, []string{memberUID}),
		TeamMembershipAuditQuery(teamUID, userUID, buildMembershipChanges(nil, []string{memberUID})))
	if err != nil {
		return result, err
	}

	return svc.getTeamDetail(session, teamUID, facilityCode)
}

func (svc *TeamsService) ReplaceMembers(teamUID, facilityCode, userUID string, userUids []string) (result models.TeamDetail, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	defer session.Close()

	if _, err = svc.getTeam(session, teamUID, facilityCode); err != nil {
		return result, err
	}

	target := dedupeNonEmpty(userUids)
	if len(target) > 0 {
		if err = svc.validateMembers(session, facilityCode, target); err != nil {
			return result, err
		}
	}

	currentUids, err := svc.getMemberUids(session, teamUID)
	if err != nil {
		return result, err
	}

	toAdd := sliceDifference(target, currentUids)
	toRemove := sliceDifference(currentUids, target)

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return svc.getTeamDetail(session, teamUID, facilityCode)
	}

	queries := make([]helpers.DatabaseQuery, 0, 3)
	if len(toRemove) > 0 {
		queries = append(queries, RemoveTeamMembersQuery(teamUID, toRemove))
	}
	if len(toAdd) > 0 {
		queries = append(queries, MergeTeamMembersQuery(teamUID, facilityCode, toAdd))
	}
	queries = append(queries, TeamMembershipAuditQuery(teamUID, userUID, buildMembershipChanges(toAdd, toRemove)))

	if err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, queries...); err != nil {
		return result, err
	}

	return svc.getTeamDetail(session, teamUID, facilityCode)
}

// validateMembers ensures every uid is a User in the facility; otherwise returns
// InvalidMembersError listing the offending uids (hard-fail, nobody is linked).
func (svc *TeamsService) validateMembers(session neo4j.Session, facilityCode string, userUids []string) error {
	valid, err := helpers.GetNeo4jSingleRecordSingleValue[[]interface{}](session, ValidateTeamMemberUidsQuery(facilityCode, userUids))
	if err != nil {
		return err
	}
	validSet := toStringSet(valid)
	var invalid []string
	for _, uid := range userUids {
		if _, ok := validSet[uid]; !ok {
			invalid = append(invalid, uid)
		}
	}
	if len(invalid) > 0 {
		return &InvalidMembersError{Uids: invalid}
	}
	return nil
}

func (svc *TeamsService) getTeamDetail(session neo4j.Session, uid, facilityCode string) (models.TeamDetail, error) {
	detail, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.TeamDetail](session, GetTeamByUIDQuery(uid, facilityCode))
	if err != nil {
		if isNoRecords(err) {
			return models.TeamDetail{}, ErrNotFound
		}
		return models.TeamDetail{}, err
	}
	if detail.Members == nil {
		detail.Members = []models.TeamMember{}
	}
	return detail, nil
}

func (svc *TeamsService) getMemberUids(session neo4j.Session, teamUID string) ([]string, error) {
	uids, err := helpers.GetNeo4jSingleRecordSingleValue[[]interface{}](session, GetTeamMemberUidsQuery(teamUID))
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(uids))
	for _, u := range uids {
		if s, ok := u.(string); ok {
			result = append(result, s)
		}
	}
	return result, nil
}

func buildMembershipChanges(added, removed []string) string {
	changes := []helpers.ChangeEntry{}
	if len(added) > 0 {
		changes = append(changes, helpers.ChangeEntry{Field: "membersAdded", Type: "membership", NewValue: added})
	}
	if len(removed) > 0 {
		changes = append(changes, helpers.ChangeEntry{Field: "membersRemoved", Type: "membership", OldValue: removed})
	}
	return helpers.MarshalChanges(changes)
}

func containsMember(members []models.TeamMember, uid string) bool {
	for _, m := range members {
		if m.UID == uid {
			return true
		}
	}
	return false
}

func dedupeNonEmpty(uids []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(uids))
	for _, u := range uids {
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		result = append(result, u)
	}
	return result
}

// sliceDifference returns elements in a that are not in b.
func sliceDifference(a, b []string) []string {
	bSet := map[string]struct{}{}
	for _, x := range b {
		bSet[x] = struct{}{}
	}
	var result []string
	for _, x := range a {
		if _, ok := bSet[x]; !ok {
			result = append(result, x)
		}
	}
	return result
}

func toStringSet(values []interface{}) map[string]struct{} {
	set := map[string]struct{}{}
	for _, v := range values {
		if s, ok := v.(string); ok {
			set[s] = struct{}{}
		}
	}
	return set
}
