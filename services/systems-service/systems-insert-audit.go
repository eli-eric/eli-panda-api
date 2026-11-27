package systemsService

import "fmt"

// insertAuditCypher writes the INSERT audit edge on system alias sys (user alias user in
// scope) with the system's initial values as change entries — the same shape POST /systems
// writes through the PATCH plan, so every create path reads "created … set X to Y" in history.
// Place it after the system's relationships are created; it reads them back. A non-empty
// where (e.g. "s._created") limits the audit to rows where it holds.
func insertAuditCypher(sys, user, where string) string {
	filter := ""
	if where != "" {
		filter = "WHERE " + where
	}
	codebook := func(rel string) string {
		return fmt.Sprintf(`head([(%s)-[:%s]->(x) | {uid: x.uid, name: x.name}])`, sys, rel)
	}
	employee := func(rel string) string {
		return fmt.Sprintf(`head([(%s)-[:%s]->(x) | {uid: x.uid, name: %s}])`, sys, rel, employeeNameCypher("x"))
	}
	str := func(field string) string {
		return fmt.Sprintf(`{field: "%[1]s", type: "string", oldValue: null, newValue: %[2]s.%[1]s}`, field, sys)
	}
	ref := func(field, value string) string {
		return fmt.Sprintf(`{field: "%s", type: "codebook", oldValue: null, newValue: %s}`, field, value)
	}
	return fmt.Sprintf(`
	WITH *
	CALL {
		WITH %[1]s, %[2]s
		WITH %[1]s, %[2]s %[3]s
		WITH %[1]s, %[2]s, [%[4]s] AS insertEntries
		CREATE (%[1]s)-[:WAS_UPDATED_BY{at: datetime(), action: "INSERT",
			changes: apoc.convert.toJson([e IN insertEntries WHERE e.newValue IS NOT NULL])}]->(%[2]s)
	}
	`, sys, user, filter, joinEntries(
		str("name"), str("systemCode"), str("systemLevel"), str("description"),
		ref("systemType", codebook("HAS_SYSTEM_TYPE")),
		ref("zone", codebook("HAS_ZONE")),
		ref("location", codebook("HAS_LOCATION")),
		ref("responsible", employee("HAS_RESPONSIBLE")),
	))
}

func joinEntries(entries ...string) string {
	out := ""
	for i, e := range entries {
		if i > 0 {
			out += ",\n\t\t\t"
		}
		out += e
	}
	return out
}
