package publicationsservice

import (
	"testing"

	"panda/apigateway/services/publications-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeORCIDAcceptsEveryPublishedForm(t *testing.T) {
	canonical := "0000-0002-1825-0097"
	for _, raw := range []string{
		"0000-0002-1825-0097",
		"0000000218250097",
		"https://orcid.org/0000-0002-1825-0097",
		"http://orcid.org/0000-0002-1825-0097",
		"  https://ORCID.org/0000-0002-1825-0097/  ",
	} {
		value, valid := normalizeORCID(raw)
		require.Truef(t, valid, "expected %q to normalize", raw)
		assert.Equal(t, canonical, value)
	}

	// The checksum character may be X, but only in the final position.
	value, valid := normalizeORCID("0000-0002-1694-233X")
	require.True(t, valid)
	assert.Equal(t, "0000-0002-1694-233X", value)

	for _, raw := range []string{"", "not-an-orcid", "0000-0002-1825", "X000-0002-1825-0097"} {
		_, valid := normalizeORCID(raw)
		assert.Falsef(t, valid, "expected %q to be rejected", raw)
	}
}

func TestMatchWosAuthorsPrefersIdentifiersOverNames(t *testing.T) {
	researchers := []wosResearcherRecord{
		{UID: "r-1", FirstName: "Jana", LastName: "Nováková", ResearcherIDs: []string{"A-1234-2015"}},
		{UID: "r-2", FirstName: "Petr", LastName: "Svoboda", Orcids: []string{"https://orcid.org/0000-0002-1825-0097"}},
		{UID: "r-3", FirstName: "Eva", LastName: "Dvořáková"},
	}

	matches := matchWosAuthors([]models.WosAuthor{
		{WosDisplayName: "Nováková, Jana", WosResearcherID: "A-1234-2015"},
		{WosDisplayName: "Someone Else Entirely", Orcid: "0000-0002-1825-0097"},
		{WosDisplayName: "Eva Dvořáková"},
		{WosDisplayName: "Nobody Here"},
	}, researchers)

	require.Len(t, matches, 4)

	assert.Equal(t, "researcher-id", matches[0].Match.Kind)
	assert.Equal(t, "r-1", matches[0].Match.Candidates[0].Uid)

	// An ORCID hit wins even though the display name matches nobody, and the
	// identifier is echoed back so the UI can show what it matched on.
	assert.Equal(t, "orcid", matches[1].Match.Kind)
	require.Len(t, matches[1].Match.Candidates, 1)
	assert.Equal(t, "r-2", matches[1].Match.Candidates[0].Uid)
	assert.Equal(t, "0000-0002-1825-0097", matches[1].Orcid)

	// Without any identifier the match degrades to a name suggestion, which the
	// editor still has to confirm.
	assert.Equal(t, "name", matches[2].Match.Kind)
	assert.Equal(t, "r-3", matches[2].Match.Candidates[0].Uid)

	assert.Equal(t, "none", matches[3].Match.Kind)
	assert.Empty(t, matches[3].Match.Candidates)
}

func TestMatchWosAuthorsReportsAmbiguousORCIDRatherThanGuessing(t *testing.T) {
	researchers := []wosResearcherRecord{
		{UID: "r-1", FirstName: "Petr", LastName: "Svoboda", Orcids: []string{"0000-0002-1825-0097"}},
		{UID: "r-2", FirstName: "Petr", LastName: "Svoboda", Orcids: []string{"0000000218250097"}},
	}

	matches := matchWosAuthors(
		[]models.WosAuthor{{WosDisplayName: "Petr Svoboda", Orcid: "https://orcid.org/0000-0002-1825-0097"}},
		researchers,
	)

	require.Len(t, matches, 1)
	assert.Equal(t, "ambiguous", matches[0].Match.Kind)
	assert.Len(t, matches[0].Match.Candidates, 2)
}

func TestMatchWosAuthorsIgnoresMalformedORCID(t *testing.T) {
	researchers := []wosResearcherRecord{
		{UID: "r-1", FirstName: "Eva", LastName: "Dvořáková", Orcids: []string{"not-an-orcid"}},
	}

	matches := matchWosAuthors(
		[]models.WosAuthor{{WosDisplayName: "Eva Dvořáková", Orcid: "also-not-an-orcid"}},
		researchers,
	)

	require.Len(t, matches, 1)
	// The unusable identifiers are skipped on both sides and the name still matches.
	assert.Equal(t, "name", matches[0].Match.Kind)
	assert.Equal(t, "r-1", matches[0].Match.Candidates[0].Uid)
}
