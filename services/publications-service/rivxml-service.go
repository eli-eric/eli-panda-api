package publicationsservice

import (
	"panda/apigateway/services/publications-service/models"
	"strconv"
	"strings"
)

// GenerateRivXml generates RIV 2025 XML for given publications.
// Chybějící hodnoty jsou označeny "__DOPLNIT__".
func GenerateRivXml(publications []models.Publication) (string, error) {
	var sb strings.Builder

	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<dodavka struktura="RIV25A">` + "\n")
	sb.WriteString(`  <zahlavi>` + "\n")
	sb.WriteString(`    <rozsah>` + "\n")
	sb.WriteString(`      <informacni-oblast>RIV</informacni-oblast>` + "\n")
	sb.WriteString(`      <obdobi-sberu>2025</obdobi-sberu>` + "\n")
	sb.WriteString(`      <predkladatel>` + "\n")
	sb.WriteString(`        <subjekt>__DOPLNIT__</subjekt>` + "\n")
	sb.WriteString(`      </predkladatel>` + "\n")
	sb.WriteString(`      <dodavatel>` + "\n")
	sb.WriteString(`        <subjekt>` + "\n")
	sb.WriteString(`          <kod>__DOPLNIT__</kod>` + "\n")
	sb.WriteString(`        </subjekt>` + "\n")
	sb.WriteString(`        <pracovnik-povereny-pripravou-dodavky>` + "\n")
	sb.WriteString(`          <osoba>` + "\n")
	sb.WriteString(`            <cele-jmeno>__DOPLNIT__</cele-jmeno>` + "\n")
	sb.WriteString(`            <kontakt>` + "\n")
	sb.WriteString(`            </kontakt>` + "\n")
	sb.WriteString(`          </osoba>` + "\n")
	sb.WriteString(`        </pracovnik-povereny-pripravou-dodavky>` + "\n")
	sb.WriteString(`      </dodavatel>` + "\n")
	sb.WriteString(`    </rozsah>` + "\n")
	sb.WriteString(`    <pruvodka cislo-jednaci="__DOPLNIT__" />` + "\n")
	sb.WriteString(`  </zahlavi>` + "\n")
	sb.WriteString(`  <obsah>` + "\n")

	for _, pub := range publications {
		// Only export supported types
		var druh string
		if pub.MediaType == "PeerReviewedArticle" {
			druh = "J"
		} else if pub.MediaType == "OtherArticle" {
			druh = "D"
		} else {
			continue
		}

		sb.WriteString("    <vysledek\n")
		sb.WriteString(`      identifikacni-kod="` + xmlEscape(pub.Uid) + `"` + "\n")
		sb.WriteString(`      duvernost-udaju="S"` + "\n")
		sb.WriteString(`      rok-uplatneni="` + xmlEscape(pub.YearOfPublication) + `"` + "\n")
		sb.WriteString(`      kontrolni-kod="__DOPLNIT__"` + "\n")
		sb.WriteString(`      druh="` + druh + `"` + "\n")
		sb.WriteString(`    >` + "\n")
		sb.WriteString(`      <jazyk>` + xmlEscape(pub.Language) + `</jazyk>` + "\n")
		sb.WriteString(`      <nazev jazyk="eng">` + xmlEscape(pub.Title) + `</nazev>` + "\n")
		if pub.Language != "eng" {
			sb.WriteString(`      <nazev jazyk="` + xmlEscape(pub.Language) + `">` + xmlEscape(pub.Title) + `</nazev>` + "\n")
		}
		sb.WriteString(`      <anotace jazyk="eng">` + xmlEscape(pub.Abstract) + `</anotace>` + "\n")
		if pub.Language != "eng" {
			sb.WriteString(`      <anotace jazyk="` + xmlEscape(pub.Language) + `">` + xmlEscape(pub.Abstract) + `</anotace>` + "\n")
		}
		sb.WriteString(`      <odkaz>` + xmlEscape(pub.WebLink) + `</odkaz>` + "\n")
		sb.WriteString(`      <doi>` + xmlEscape(pub.Doi) + `</doi>` + "\n")
		sb.WriteString(generateAutori(pub.EliAuthors, pub.EliAuthorsCount))
		sb.WriteString(`      <klasifikace>` + "\n")
		sb.WriteString(`        <obor postaveni="hlavni" ciselnik="OblastiOECD">` + valueOrPlaceholder(pub.OecdFord) + `</obor>` + "\n")
		// Keywords: split by ';' and output each as <klicove-slovo>
		keywords := splitAndTrim(pub.Keywords, ";")
		for _, kw := range keywords {
			sb.WriteString(`        <klicove-slovo jazyk="eng">` + xmlEscape(kw) + `</klicove-slovo>` + "\n")
		}
		sb.WriteString(`      </klasifikace>` + "\n")
		// navaznosti/projekt (grant)
		sb.WriteString(`      <navaznosti>` + "\n")
		sb.WriteString(`        <navaznost druh-vztahu="byl-dosazen-pri-reseni">` + "\n")
		if pub.Grant != nil && *pub.Grant != "" {
			sb.WriteString(`          <projekt identifikacni-kod="` + xmlEscape(*pub.Grant) + `" />` + "\n")
		} else {
			sb.WriteString(`          <projekt identifikacni-kod="__DOPLNIT__" />` + "\n")
		}
		sb.WriteString(`        </navaznost>` + "\n")
		sb.WriteString(`      </navaznosti>` + "\n")

		if druh == "J" {
			// Journal article
			sb.WriteString(`      <poddruh>__DOPLNIT__</poddruh>` + "\n")
			sb.WriteString(`      <periodikum>` + "\n")
			sb.WriteString(`        <nazev>` + xmlEscape(pub.LongJournalTitle) + `</nazev>` + "\n")
			sb.WriteString(`        <ISSN>` + valueOrPlaceholder(pub.Issn) + `</ISSN>` + "\n")
			sb.WriteString(`        <eISSN>` + valueOrPlaceholder(pub.EIssn) + `</eISSN>` + "\n")
			sb.WriteString(`        <vydavatel>` + "\n")
			if pub.PublishingCountry != nil && pub.PublishingCountry.Code != "" {
				sb.WriteString(`          <stat>` + xmlEscape(pub.PublishingCountry.Code) + `</stat>` + "\n")
			} else {
				sb.WriteString(`          <stat>__DOPLNIT__</stat>` + "\n")
			}
			sb.WriteString(`        </vydavatel>` + "\n")
			sb.WriteString(`      </periodikum>` + "\n")
			sb.WriteString(`      <rocnik>` + intToString(pub.Volume) + `</rocnik>` + "\n")
			sb.WriteString(`      <cislo>` + valueOrPlaceholderInt(pub.Issue) + `</cislo>` + "\n")
			sb.WriteString(`      <strany pocet="` + intToString(pub.PagesCount) + `">` + "\n")
			sb.WriteString(`        <rozsah>` + xmlEscape(pub.Pages) + `</rozsah>` + "\n")
			sb.WriteString(`      </strany>` + "\n")
			sb.WriteString(`      <kod-ut-isi>` + valueOrPlaceholder(pub.WosNumber) + `</kod-ut-isi>` + "\n")
			sb.WriteString(`      <EID>` + valueOrPlaceholder(pub.EidScopus) + `</EID>` + "\n")
			if pub.OpenAccessType != nil && pub.OpenAccessType.Code != "" {
				sb.WriteString(`      <zpusob-publikovani>` + xmlEscape(pub.OpenAccessType.Code) + `</zpusob-publikovani>` + "\n")
			} else {
				sb.WriteString(`      <zpusob-publikovani>__DOPLNIT__</zpusob-publikovani>` + "\n")
			}
			sb.WriteString(`      <termin-zverejneni>__DOPLNIT__</termin-zverejneni>` + "\n")
		} else if druh == "D" {
			// Conference proceedings
			sb.WriteString(`      <sbornik>` + "\n")
			sb.WriteString(`        <nazev>` + xmlEscape(pub.LongJournalTitle) + `</nazev>` + "\n")
			sb.WriteString(`        <ISBN>__DOPLNIT__</ISBN>` + "\n")
			sb.WriteString(`        <ISSN>` + valueOrPlaceholder(pub.Issn) + `</ISSN>` + "\n")
			sb.WriteString(`        <eISSN>` + valueOrPlaceholder(pub.EIssn) + `</eISSN>` + "\n")
			sb.WriteString(`        <forma-vydani>__DOPLNIT__</forma-vydani>` + "\n")
			sb.WriteString(`        <misto-vydani>__DOPLNIT__</misto-vydani>` + "\n")
			sb.WriteString(`        <nakladatel>` + "\n")
			sb.WriteString(`          <nazev>__DOPLNIT__</nazev>` + "\n")
			sb.WriteString(`        </nakladatel>` + "\n")
			sb.WriteString(`      </sbornik>` + "\n")
			sb.WriteString(`      <akce>` + "\n")
			sb.WriteString(`        <konani>` + "\n")
			sb.WriteString(`          <v-roce>` + xmlEscape(pub.YearOfPublication) + `</v-roce>` + "\n")
			sb.WriteString(`          <misto>__DOPLNIT__</misto>` + "\n")
			sb.WriteString(`        </konani>` + "\n")
			sb.WriteString(`        <ucastnici>` + "\n")
			sb.WriteString(`          <klasifikace-podle-statni-prislusnosti>__DOPLNIT__</klasifikace-podle-statni-prislusnosti>` + "\n")
			sb.WriteString(`        </ucastnici>` + "\n")
			sb.WriteString(`      </akce>` + "\n")
			sb.WriteString(`      <strany pocet="` + intToString(pub.PagesCount) + `">` + "\n")
			sb.WriteString(`        <rozsah>` + xmlEscape(pub.Pages) + `</rozsah>` + "\n")
			sb.WriteString(`      </strany>` + "\n")
			sb.WriteString(`      <kod-ut-isi>` + valueOrPlaceholder(pub.WosNumber) + `</kod-ut-isi>` + "\n")
			sb.WriteString(`      <EID>` + valueOrPlaceholder(pub.EidScopus) + `</EID>` + "\n")
		}
		sb.WriteString(`    </vysledek>` + "\n")
	}

	sb.WriteString(`  </obsah>` + "\n")
	sb.WriteString(`</dodavka>`)

	return sb.String(), nil
}

// generateAutori creates the <autori> block from a comma-separated list of full names.
func generateAutori(eliAuthors string, eliAuthorsCount int) string {
	authors := splitAndTrim(eliAuthors, ",")
	if len(authors) == 0 || (len(authors) == 1 && authors[0] == "") {
		return `      <autori pocet-celkem="0" pocet-domacich="0"></autori>` + "\n"
	}
	sb := strings.Builder{}
	sb.WriteString(`      <autori pocet-celkem="` + strconv.Itoa(len(authors)) + `" pocet-domacich="` + strconv.Itoa(len(authors)) + `">` + "\n")
	for _, fullName := range authors {
		first, last := splitName(fullName)
		sb.WriteString(`        <autor je-domaci="true">` + "\n")
		sb.WriteString(`          <jmeno>` + xmlEscape(first) + `</jmeno>` + "\n")
		sb.WriteString(`          <prijmeni>` + xmlEscape(last) + `</prijmeni>` + "\n")
		sb.WriteString(`          <statni-prislusnost>CZ</statni-prislusnost>` + "\n")
		sb.WriteString(`          <rodne-cislo>__DOPLNIT__</rodne-cislo>` + "\n")
		sb.WriteString(`          <orcid>__DOPLNIT__</orcid>` + "\n")
		sb.WriteString(`          <scopusid>__DOPLNIT__</scopusid>` + "\n")
		sb.WriteString(`          <researcherid>__DOPLNIT__</researcherid>` + "\n")
		sb.WriteString(`        </autor>` + "\n")
	}
	sb.WriteString(`      </autori>` + "\n")
	return sb.String()
}

// splitAndTrim splits a string by sep and trims whitespace from each part.
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var out []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// splitName splits a full name into first name(s) and surname (last word).
func splitName(fullName string) (string, string) {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "__DOPLNIT__", "__DOPLNIT__"
	}
	if len(parts) == 1 {
		return parts[0], "__DOPLNIT__"
	}
	return strings.Join(parts[:len(parts)-1], " "), parts[len(parts)-1]
}

// --- helpers ---

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func valueOrPlaceholder(ptr *string) string {
	if ptr != nil && *ptr != "" {
		return xmlEscape(*ptr)
	}
	return "__DOPLNIT__"
}

func valueOrPlaceholderInt(ptr *int) string {
	if ptr != nil {
		return strconv.Itoa(*ptr)
	}
	return "__DOPLNIT__"
}

func intToString(i int) string {
	if i == 0 {
		return "__DOPLNIT__"
	}
	return strconv.Itoa(i)
}