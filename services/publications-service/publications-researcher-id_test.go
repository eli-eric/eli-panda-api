package publicationsservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResearcherIDYear(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantYear   int
		wantParsed bool
	}{
		// Vintages taken from real ELI researcher records.
		{name: "single letter prefix", id: "E-9444-2015", wantYear: 2015, wantParsed: true},
		{name: "three letter prefix", id: "GZZ-7943-2022", wantYear: 2022, wantParsed: true},
		{name: "recent", id: "HKH-1227-2023", wantYear: 2023, wantParsed: true},
		{name: "lowercase is normalized", id: "aaa-1673-2020", wantYear: 2020, wantParsed: true},
		{name: "surrounding space", id: "  I-6474-2015 ", wantYear: 2015, wantParsed: true},
		{name: "orcid shaped", id: "0000-0002-1825-0097"},
		{name: "four letter prefix", id: "ABCD-1234-2020"},
		{name: "empty", id: ""},
		{name: "not an id at all", id: "not-an-id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			year, parsed := researcherIDYear(test.id)
			assert.Equal(t, test.wantParsed, parsed)
			assert.Equal(t, test.wantYear, year)
		})
	}
}

func TestNewestResearcherID(t *testing.T) {
	tests := []struct {
		name     string
		ids      []string
		want     string
		wantSure bool
	}{
		{
			name:     "picks the latest vintage",
			ids:      []string{"E-9444-2015", "HKH-1227-2023", "AAA-1673-2020"},
			want:     "HKH-1227-2023",
			wantSure: true,
		},
		{
			name:     "single id is trivially newest",
			ids:      []string{"G-5767-2014"},
			want:     "G-5767-2014",
			wantSure: true,
		},
		{
			name:     "normalizes before returning",
			ids:      []string{"e-9444-2015", " hkh-1227-2023 "},
			want:     "HKH-1227-2023",
			wantSure: true,
		},
		{
			name: "a tie has no defensible winner",
			ids:  []string{"E-1111-2022", "GZZ-7943-2022"},
		},
		{
			name:     "duplicates of one id are not a tie",
			ids:      []string{"E-1111-2022", "e-1111-2022 "},
			want:     "E-1111-2022",
			wantSure: true,
		},
		{
			name:     "unreadable ids are skipped, not ranked",
			ids:      []string{"garbage", "H-3158-2014"},
			want:     "H-3158-2014",
			wantSure: true,
		},
		{name: "nothing readable", ids: []string{"garbage", ""}},
		{name: "empty list", ids: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, sure := newestResearcherID(test.ids)
			assert.Equal(t, test.wantSure, sure)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestResolvePrimaryResearcherID(t *testing.T) {
	all := []string{"E-9444-2015", "HKH-1227-2023"}

	t.Run("empty field takes the newest known id", func(t *testing.T) {
		assert.Equal(t, "HKH-1227-2023",
			resolvePrimaryResearcherID("", "E-9444-2015", all, false))
	})

	t.Run("empty field falls back to the confirmed id when newest is ambiguous", func(t *testing.T) {
		tied := []string{"E-1111-2022", "GZZ-7943-2022"}
		assert.Equal(t, "GZZ-7943-2022",
			resolvePrimaryResearcherID("", "GZZ-7943-2022", tied, false))
	})

	t.Run("an older import never demotes a current id on its own", func(t *testing.T) {
		assert.Equal(t, "HKH-1227-2023",
			resolvePrimaryResearcherID("HKH-1227-2023", "E-9444-2015", all, false))
	})

	t.Run("a newer id still waits to be asked for", func(t *testing.T) {
		assert.Equal(t, "E-9444-2015",
			resolvePrimaryResearcherID("E-9444-2015", "HKH-1227-2023", all, false))
	})

	t.Run("promotion is honoured when requested", func(t *testing.T) {
		assert.Equal(t, "HKH-1227-2023",
			resolvePrimaryResearcherID("E-9444-2015", "HKH-1227-2023", all, true))
	})
}
