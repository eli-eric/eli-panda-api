package cronservice

import (
	"testing"

	"panda/apigateway/services/cron-service/models"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLogin(t *testing.T) {
	cases := map[string]string{
		" Jan.Novak@ELI-Laser.eu ": "jan.novak@eli-laser.eu",
		"jan.novak@eli-beams.eu":   "jan.novak@eli-laser.eu",
		"jan.novak@eli-lasers.eu":  "jan.novak@eli-laser.eu",
		"jan.novak@external.cz":    "jan.novak@external.cz",
		"jnovak":                   "jnovak",
		"":                         "",
	}

	for input, expected := range cases {
		assert.Equal(t, expected, NormalizeLogin(input), input)
	}
}

func TestEmployeeLogins(t *testing.T) {
	tests := []struct {
		name     string
		employee models.SyncEliBeamlinesEmployee
		expected []string
	}{
		{
			name:     "login differs from email, login first",
			employee: models.SyncEliBeamlinesEmployee{LDAP_USER: "Jana.Nova@eli-laser.eu", EMAIL: "jana.stara@eli-beams.eu"},
			expected: []string{"jana.nova@eli-laser.eu", "jana.stara@eli-laser.eu"},
		},
		{
			name:     "same after normalization is listed once",
			employee: models.SyncEliBeamlinesEmployee{LDAP_USER: "jan.novak@eli-laser.eu", EMAIL: "Jan.Novak@eli-beams.eu"},
			expected: []string{"jan.novak@eli-laser.eu"},
		},
		{
			name:     "missing login falls back to email",
			employee: models.SyncEliBeamlinesEmployee{LDAP_USER: " ", EMAIL: "jan.novak@eli-lasers.eu"},
			expected: []string{"jan.novak@eli-laser.eu"},
		},
		{
			name:     "login without email",
			employee: models.SyncEliBeamlinesEmployee{LDAP_USER: "jan.novak@eli-laser.eu"},
			expected: []string{"jan.novak@eli-laser.eu"},
		},
		{
			name:     "nothing to pair with",
			employee: models.SyncEliBeamlinesEmployee{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, EmployeeLogins(tt.employee))
		})
	}
}

func TestSyncEliBeamlinesEmployeeQuery_Parameters(t *testing.T) {
	query := SyncEliBeamlinesEmployeeQuery(models.SyncEliBeamlinesEmployee{
		OSOBNI_CISLO: "1234",
		LDAP_USER:    " Jana.Nova@eli-laser.eu ",
		EMAIL:        "Jana.Stara@eli-beams.eu",
		OJ_KOD:       "300",
	}, "B")

	assert.Equal(t, "1234", query.Parameters["employeeNumber"])
	assert.Equal(t, "Jana.Nova@eli-laser.eu", query.Parameters["ldapUser"])
	assert.Equal(t, "300", query.Parameters["orgUnitCode"])
	assert.Equal(t, []string{"jana.nova@eli-laser.eu", "jana.stara@eli-laser.eu"}, query.Parameters["logins"])
}
