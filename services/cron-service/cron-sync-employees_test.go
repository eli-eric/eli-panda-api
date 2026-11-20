package cronservice

import (
	"strings"
	"testing"

	"panda/apigateway/helpers"
	"panda/apigateway/services/cron-service/models"
	"panda/apigateway/services/testsetup"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
)

// Integration tests against the local Neo4j (see AGENTS.md). Every test works on its own
// nodes (random employee number / logins) and removes them afterwards.

type syncFixture struct {
	t       *testing.T
	session neo4j.Session
	id      string
}

func newSyncFixture(t *testing.T) *syncFixture {
	t.Helper()
	session, err := helpers.NewNeo4jSession(testsetup.TestDriver)
	assert.NoError(t, err)

	f := &syncFixture{t: t, session: session, id: "test-sync-" + uuid.New().String()[:8]}
	t.Cleanup(func() {
		f.run(`MATCH (n) WHERE (n:Employee AND n.employeeNumber STARTS WITH $id) OR (n:User AND n.uid STARTS WITH $id) DETACH DELETE n`, nil)
		session.Close()
	})
	return f
}

func (f *syncFixture) run(query string, params map[string]interface{}) {
	f.t.Helper()
	if params == nil {
		params = map[string]interface{}{}
	}
	params["id"] = f.id
	assert.NoError(f.t, helpers.WriteNeo4jAndReturnNothing(f.session, helpers.DatabaseQuery{Query: query, Parameters: params}))
}

func (f *syncFixture) number(suffix string) string { return f.id + "-" + suffix }

// login builds a unique login in the given domain
func (f *syncFixture) login(name, domain string) string { return f.id + "." + name + "@" + domain }

func (f *syncFixture) createUser(key, username string) {
	f.run(`CREATE (:User {uid: $id + '-' + $key, username: $username, email: $id + '-' + $key + '@test.invalid'})`,
		map[string]interface{}{"key": key, "username": username})
}

func (f *syncFixture) createEmployee(number string) {
	f.run(`CREATE (:Employee {uid: $number, employeeNumber: $number})`, map[string]interface{}{"number": number})
}

func (f *syncFixture) link(number, userKey string) {
	f.run(`MATCH (e:Employee {employeeNumber: $number}), (u:User {uid: $id + '-' + $key}) MERGE (e)-[:HAS_USER]->(u)`,
		map[string]interface{}{"number": number, "key": userKey})
}

// linkedUsers returns the User keys linked to every Employee with the given number
func (f *syncFixture) linkedUsers(number string) []string {
	f.t.Helper()
	query := helpers.DatabaseQuery{
		Query:       `MATCH (e:Employee {employeeNumber: $number})-[:HAS_USER]->(u:User) RETURN collect(u.uid) AS result`,
		Parameters:  map[string]interface{}{"number": number},
		ReturnAlias: "result",
	}
	uids, err := helpers.GetNeo4jSingleRecordSingleValue[[]interface{}](f.session, query)
	assert.NoError(f.t, err)

	keys := []string{}
	for _, uid := range uids {
		keys = append(keys, strings.TrimPrefix(uid.(string), f.id+"-"))
	}
	return keys
}

func (f *syncFixture) sync(employee models.SyncEliBeamlinesEmployee) error {
	return syncEmployee(f.session, employee)
}

func TestSyncEmployee_PrefersLoginOverEmail(t *testing.T) {
	f := newSyncFixture(t)
	number := f.number("1")
	f.createUser("byLogin", f.login("jana.nova", "eli-laser.eu"))
	f.createUser("byEmail", f.login("jana.stara", "eli-laser.eu"))

	err := f.sync(models.SyncEliBeamlinesEmployee{
		OSOBNI_CISLO: number,
		LDAP_USER:    strings.ToUpper(f.login("jana.nova", "eli-laser.eu")),
		EMAIL:        f.login("jana.stara", "eli-beams.eu"),
		OJ_KOD:       "300",
	})

	assert.NoError(t, err)
	assert.Equal(t, []string{"byLogin"}, f.linkedUsers(number))
}

func TestSyncEmployee_FallsBackToNormalizedEmail(t *testing.T) {
	f := newSyncFixture(t)
	number := f.number("1")
	f.createUser("byEmail", f.login("jan.novak", "eli-laser.eu"))

	err := f.sync(models.SyncEliBeamlinesEmployee{OSOBNI_CISLO: number, EMAIL: f.login("jan.novak", "eli-lasers.eu")})

	assert.NoError(t, err)
	assert.Equal(t, []string{"byEmail"}, f.linkedUsers(number))
}

func TestSyncEmployee_StoresLdapUserAndOrgUnitCode(t *testing.T) {
	f := newSyncFixture(t)
	number := f.number("1")
	login := f.login("jan.novak", "eli-laser.eu")

	assert.NoError(t, f.sync(models.SyncEliBeamlinesEmployee{OSOBNI_CISLO: number, LDAP_USER: login, OJ_KOD: "300"}))

	query := helpers.DatabaseQuery{
		Query:       `MATCH (e:Employee {employeeNumber: $number}) RETURN e.ldapUser + '|' + e.orgUnitCode AS result`,
		Parameters:  map[string]interface{}{"number": number},
		ReturnAlias: "result",
	}
	stored, err := helpers.GetNeo4jSingleRecordSingleValue[string](f.session, query)
	assert.NoError(t, err)
	assert.Equal(t, login+"|300", stored)
}

func TestSyncEmployee_KeepsExistingLink(t *testing.T) {
	f := newSyncFixture(t)
	number := f.number("1")
	f.createUser("current", f.login("old.account", "eli-laser.eu"))
	f.createUser("byLogin", f.login("jan.novak", "eli-laser.eu"))
	f.createEmployee(number)
	f.link(number, "current")

	assert.NoError(t, f.sync(models.SyncEliBeamlinesEmployee{OSOBNI_CISLO: number, LDAP_USER: f.login("jan.novak", "eli-laser.eu")}))

	assert.Equal(t, []string{"current"}, f.linkedUsers(number))
}

func TestSyncEmployee_DoesNotLinkUserOfAnotherEmployee(t *testing.T) {
	f := newSyncFixture(t)
	number, otherNumber := f.number("1"), f.number("2")
	f.createUser("byLogin", f.login("jan.novak", "eli-laser.eu"))
	f.createUser("byEmail", f.login("jan.novak2", "eli-laser.eu"))
	f.createEmployee(otherNumber)
	f.link(otherNumber, "byLogin")

	err := f.sync(models.SyncEliBeamlinesEmployee{
		OSOBNI_CISLO: number,
		LDAP_USER:    f.login("jan.novak", "eli-laser.eu"),
		EMAIL:        f.login("jan.novak2", "eli-laser.eu"),
	})

	// the login account is taken: no link, and no fallback to the e-mail account either
	assert.NoError(t, err)
	assert.Empty(t, f.linkedUsers(number))
	assert.Equal(t, []string{"byLogin"}, f.linkedUsers(otherNumber))
}

func TestSyncEmployee_SkipsDuplicatedEmployeeNumber(t *testing.T) {
	f := newSyncFixture(t)
	number := f.number("1")
	f.createUser("byLogin", f.login("jan.novak", "eli-laser.eu"))
	f.createEmployee(number)
	f.createEmployee(number)

	err := f.sync(models.SyncEliBeamlinesEmployee{OSOBNI_CISLO: number, LDAP_USER: f.login("jan.novak", "eli-laser.eu")})

	assert.ErrorContains(t, err, "has 2 Employee nodes")
	assert.Empty(t, f.linkedUsers(number))
}
