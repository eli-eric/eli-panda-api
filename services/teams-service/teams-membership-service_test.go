package teamsservice

import (
	"encoding/json"
	"testing"

	"panda/apigateway/services/teams-service/models"

	"github.com/stretchr/testify/assert"
)

func TestDedupeNonEmpty(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{"nil", nil, []string{}},
		{"empty slice", []string{}, []string{}},
		{"drops empties", []string{"a", "", "b", ""}, []string{"a", "b"}},
		{"drops duplicates, keeps first order", []string{"a", "b", "a", "c", "b"}, []string{"a", "b", "c"}},
		{"all empty", []string{"", ""}, []string{}},
		{"clean passthrough", []string{"x", "y", "z"}, []string{"x", "y", "z"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, dedupeNonEmpty(tc.in))
		})
	}
}

func TestSliceDifference(t *testing.T) {
	cases := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{"disjoint", []string{"a", "b"}, []string{"c", "d"}, []string{"a", "b"}},
		{"full overlap", []string{"a", "b"}, []string{"a", "b"}, nil},
		{"partial overlap", []string{"a", "b", "c"}, []string{"b"}, []string{"a", "c"}},
		{"empty a", nil, []string{"a"}, nil},
		{"empty b", []string{"a", "b"}, nil, []string{"a", "b"}},
		{"preserves a order", []string{"c", "a", "b"}, []string{"a"}, []string{"c", "b"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, sliceDifference(tc.a, tc.b))
		})
	}
}

// TestReplaceMembersDiff exercises the exact add/remove computation ReplaceMembers performs,
// locking in replace-all semantics (target set replaces current set) without needing a DB.
func TestReplaceMembersDiff(t *testing.T) {
	cases := []struct {
		name       string
		current    []string
		target     []string
		wantAdd    []string
		wantRemove []string
	}{
		{"add all to empty", []string{}, []string{"a", "b"}, []string{"a", "b"}, nil},
		{"clear all", []string{"a", "b"}, []string{}, nil, []string{"a", "b"}},
		{"no change", []string{"a", "b"}, []string{"a", "b"}, nil, nil},
		{"swap one", []string{"a", "b"}, []string{"a", "c"}, []string{"c"}, []string{"b"}},
		{"grow", []string{"a"}, []string{"a", "b", "c"}, []string{"b", "c"}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			target := dedupeNonEmpty(tc.target)
			toAdd := sliceDifference(target, tc.current)
			toRemove := sliceDifference(tc.current, target)
			assert.Equal(t, tc.wantAdd, toAdd, "toAdd")
			assert.Equal(t, tc.wantRemove, toRemove, "toRemove")
		})
	}
}

func TestBuildMembershipChanges(t *testing.T) {
	t.Run("added only", func(t *testing.T) {
		entries := unmarshalChanges(t, buildMembershipChanges([]string{"a", "b"}, nil))
		assert.Len(t, entries, 1)
		assert.Equal(t, "membersAdded", entries[0].Field)
		assert.ElementsMatch(t, []interface{}{"a", "b"}, entries[0].NewValue)
		assert.Nil(t, entries[0].OldValue)
	})

	t.Run("removed only", func(t *testing.T) {
		entries := unmarshalChanges(t, buildMembershipChanges(nil, []string{"x"}))
		assert.Len(t, entries, 1)
		assert.Equal(t, "membersRemoved", entries[0].Field)
		assert.ElementsMatch(t, []interface{}{"x"}, entries[0].OldValue)
		assert.Nil(t, entries[0].NewValue)
	})

	t.Run("both", func(t *testing.T) {
		entries := unmarshalChanges(t, buildMembershipChanges([]string{"a"}, []string{"b"}))
		assert.Len(t, entries, 2)
	})

	t.Run("neither yields empty array", func(t *testing.T) {
		assert.Equal(t, "[]", buildMembershipChanges(nil, nil))
	})
}

func TestContainsMember(t *testing.T) {
	members := []models.TeamMember{{UID: "a"}, {UID: "b"}}
	assert.True(t, containsMember(members, "b"))
	assert.False(t, containsMember(members, "z"))
	assert.False(t, containsMember(nil, "a"))
}

type changeEntry struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"oldValue"`
	NewValue interface{} `json:"newValue"`
}

func unmarshalChanges(t *testing.T, s string) []changeEntry {
	t.Helper()
	var entries []changeEntry
	assert.NoError(t, json.Unmarshal([]byte(s), &entries))
	return entries
}
