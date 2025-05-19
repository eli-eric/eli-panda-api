# Instructions for Generating RIV 2025 XML Export

## Overview

This document describes how to generate XML exports for the Czech RIV (Rejstřík informací o výsledcích) system, strictly following the 2025 specification (version 3.4.0). Only two types of publication are supported:

- **PeerReviewedArticle** (`mediaType == "PeerReviewedArticle"`) → RIV druh="J" (clanek-v-periodiku)
- **OtherArticle** (`mediaType == "OtherArticle"`) → RIV druh="D" (clanek-ve-sborniku)

All other types must be skipped in the export.

The XML structure, field mapping, and placeholder policy are based on the official document:  
_Export do RIVu_XML_2025_3.4.0_v1.pdf_ ([see isvavai.cz](https://www.isvavai.cz/)).

---

## 1. Filtering Logic

- **Include** only publications where:
  - `mediaType == "PeerReviewedArticle"` (journal article, RIV druh="J")
  - `mediaType == "OtherArticle"` (conference proceedings, RIV druh="D")
- **Skip** all other `mediaType` values.

---

## 2. Field Mapping

### Common Fields

| RIV XML Field             | Publication Model Field | Notes/Placeholder         |
| ------------------------- | ----------------------- | ------------------------- |
| identifikacni-kod         | Uid                     | Must follow RIV format    |
| duvernost-udaju           | (hardcoded "S")         | "S" = veřejně přístupné   |
| rok-uplatneni             | YearOfPublication       |                           |
| kontrolni-kod             | (placeholder)           | `__DOPLNIT__`             |
| druh                      | (see mapping above)     | "J" or "D"                |
| jazyk                     | Language                | 3-letter code             |
| nazev (eng)               | Title                   |                           |
| nazev (orig)              | Title (if not eng)      | Only if Language != "eng" |
| anotace (eng)             | Abstract                |                           |
| anotace (orig)            | Abstract (if not eng)   | Only if Language != "eng" |
| odkaz                     | WebLink                 |                           |
| doi                       | Doi                     |                           |
| autori                    | EliAuthors              | See below for structure   |
| klasifikace/obor          | OecdFord                |                           |
| klasifikace/klicove-slovo | Keywords                | Split by ";"              |
| navaznosti/projekt        | Grant or (placeholder)  | `__DOPLNIT__` if missing  |

### "PeerReviewedArticle" (J) Specific

- poddruh, periodikum, ISSN, eISSN, vydavatel/stat, rocnik, cislo, strany, kod-ut-isi, EID, zpusob-publikovani, termin-zverejneni, etc.

### "OtherArticle" (D) Specific

- sbornik, ISBN, ISSN, eISSN, forma-vydani, misto-vydani, nakladatel, akce, konani, strany, kod-ut-isi, EID, etc.

All fields required by the RIV 2025 spec for these types must be mapped, with `__DOPLNIT__` for any missing data.

---

## 3. Author Handling (`<autori>` and `<autor>`)

- For each name in `EliAuthors` (comma-separated):
  - Split into first name(s) and surname (by last space).
  - Generate `<autor>` with all required subfields:
    - `<jmeno>`: first name(s)
    - `<prijmeni>`: surname
    - `je-domaci="true"`
    - `<statni-prislusnost>CZ</statni-prislusnost>`
    - `<rodne-cislo>__DOPLNIT__</rodne-cislo>`
    - `<orcid>__DOPLNIT__</orcid>`
    - `<scopusid>__DOPLNIT__</scopusid>`
    - `<researcherid>__DOPLNIT__</researcherid>`
- Set `pocet-celkem` and `pocet-domacich` to the number of EliAuthors.

**Example:**

```xml
<autori pocet-celkem="2" pocet-domacich="2">
  <autor je-domaci="true">
    <jmeno>Jan</jmeno>
    <prijmeni>Novak</prijmeni>
    <statni-prislusnost>CZ</statni-prislusnost>
    <rodne-cislo>__DOPLNIT__</rodne-cislo>
    <orcid>__DOPLNIT__</orcid>
    <scopusid>__DOPLNIT__</scopusid>
    <researcherid>__DOPLNIT__</researcherid>
  </autor>
  <autor je-domaci="true">
    <jmeno>Petr</jmeno>
    <prijmeni>Svoboda</prijmeni>
    <statni-prislusnost>CZ</statni-prislusnost>
    <rodne-cislo>__DOPLNIT__</rodne-cislo>
    <orcid>__DOPLNIT__</orcid>
    <scopusid>__DOPLNIT__</scopusid>
    <researcherid>__DOPLNIT__</researcherid>
  </autor>
</autori>
```

---

## 4. Placeholders

- Use `__DOPLNIT__` for all missing or unavailable data, both in elements and attributes.

---

## 5. XML Structure

- Strictly follow the RIV 2025 structure for each type, including all required and conditional fields.
- Use the correct nesting and attribute structure as per the official PDF.

---

## 6. Implementation Steps

1. **Update Filtering**: In `GenerateRivXml`, filter publications by `mediaType` as above.
2. **Type Switch**: For each publication, switch on `mediaType` and generate the correct XML structure ("J" or "D").
3. **Field Mapping**: For each type, map fields as per the mapping table, using placeholders where needed.
4. **Author XML**: Implement a helper to split `EliAuthors` and generate `<autori>` and `<autor>` elements with all required subfields.
5. **Placeholders**: Ensure all missing or unavailable data is filled with `__DOPLNIT__`.
6. **Helpers**: Add helpers for XML escaping, splitting names, and formatting authors.
7. **Testing**: Add test cases for both types, including edge cases with missing data.

---

## 7. Mermaid Diagram: Filtering and Generation Flow

```mermaid
flowchart TD
    A[Start] --> B{mediaType}
    B -- PeerReviewedArticle --> C[Generate RIV "J" XML]
    B -- OtherArticle --> D[Generate RIV "D" XML]
    B -- else --> E[Skip]
    C & D --> F[Insert into <obsah>]
    F --> G[Finish XML]
```

---

## 8. Reference

- Official specification: _Export do RIVu_XML_2025_3.4.0_v1.pdf_ ([isvavai.cz](https://www.isvavai.cz/))
- This document is intended for developers and data managers responsible for RIV XML export.
