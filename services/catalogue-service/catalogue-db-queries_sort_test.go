package catalogueService

import (
	"panda/apigateway/helpers"
	"strings"
	"testing"
)

func TestBuildSortClause_EmptySorting(t *testing.T) {
	opt, ret, orderBy, params := buildCatalogueItemsSortClause(nil, nil)
	if opt != "" {
		t.Errorf("optionalMatchSQL not empty: %q", opt)
	}
	if len(ret) != 0 {
		t.Errorf("extraReturnVars not empty: %v", ret)
	}
	want := " ORDER BY lastUpdateTime DESC, itm.uid ASC "
	if orderBy != want {
		t.Errorf("orderBy = %q, want %q", orderBy, want)
	}
	if len(params) != 0 {
		t.Errorf("params not empty: %v", params)
	}
}

func TestBuildSortClause_Builtins(t *testing.T) {
	cases := []struct {
		id       string
		desc     bool
		wantExpr string
	}{
		{"name", false, "itm.name ASC"},
		{"name", true, "itm.name DESC"},
		{"catalogueNumber", false, "itm.catalogueNumber ASC"},
		{"partNumber", false, "itm.catalogueNumber ASC"},
		{"categoryName", true, "cat.name DESC"},
		{"description", false, "itm.description ASC"},
		{"manufacturerUrl", false, "itm.manufacturerUrl ASC"},
		{"supplier", false, "sname ASC"},
		{"lastUpdateTime", true, "lastUpdateTime DESC"},
		{"lastUpdateBy", false, "lastUpdateUser ASC"},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			opt, ret, orderBy, params := buildCatalogueItemsSortClause(
				&[]helpers.Sorting{{ID: c.id, DESC: c.desc}}, nil)
			if opt != "" {
				t.Errorf("optionalMatchSQL not empty: %q", opt)
			}
			if len(ret) != 0 {
				t.Errorf("extraReturnVars not empty: %v", ret)
			}
			want := " ORDER BY " + c.wantExpr + ", itm.uid ASC "
			if orderBy != want {
				t.Errorf("orderBy = %q, want %q", orderBy, want)
			}
			if len(params) != 0 {
				t.Errorf("params not empty: %v", params)
			}
		})
	}
}

func TestBuildSortClause_CustomProperty_Text(t *testing.T) {
	uid := "1adeae21-ab0f-40f3-868d-36f7d40210ca"
	propTypes := sortPropertyTypes{uid: "text"}
	opt, ret, orderBy, params := buildCatalogueItemsSortClause(
		&[]helpers.Sorting{{ID: uid, DESC: false}}, propTypes)

	if !strings.Contains(opt, "OPTIONAL MATCH (itm)-[pvSort0:HAS_CATALOGUE_PROPERTY]->(propSort0{uid: $sortPropKey0})") {
		t.Errorf("missing OPTIONAL MATCH; opt=%q", opt)
	}
	if !strings.Contains(opt, "toLower(toString(pvSort0.value)) AS sortValue0") {
		t.Errorf("missing text expression; opt=%q", opt)
	}
	if !strings.Contains(opt, "WITH itm, cat, sname,") {
		t.Errorf("missing narrowing WITH; opt=%q", opt)
	}
	if len(ret) != 1 || ret[0] != "sortValue0" {
		t.Errorf("extraReturnVars = %v", ret)
	}
	if orderBy != " ORDER BY sortValue0 ASC, itm.uid ASC " {
		t.Errorf("orderBy = %q", orderBy)
	}
	if params["sortPropKey0"] != uid {
		t.Errorf("params[sortPropKey0] = %v, want %v", params["sortPropKey0"], uid)
	}
}

func TestBuildSortClause_CustomProperty_Number(t *testing.T) {
	uid := "1adeae21-ab0f-40f3-868d-36f7d40210ca"
	propTypes := sortPropertyTypes{uid: "number"}
	opt, _, orderBy, _ := buildCatalogueItemsSortClause(
		&[]helpers.Sorting{{ID: uid, DESC: true}}, propTypes)

	if !strings.Contains(opt, "toFloat(pvSort0.value) AS sortValue0") {
		t.Errorf("missing number expression; opt=%q", opt)
	}
	if orderBy != " ORDER BY sortValue0 DESC, itm.uid ASC " {
		t.Errorf("orderBy = %q", orderBy)
	}
}

func TestBuildSortClause_CustomProperty_Range(t *testing.T) {
	uid := "1adeae21-ab0f-40f3-868d-36f7d40210ca"
	propTypes := sortPropertyTypes{uid: "range"}
	opt, _, _, _ := buildCatalogueItemsSortClause(
		&[]helpers.Sorting{{ID: uid, DESC: false}}, propTypes)

	if !strings.Contains(opt, "toFloat(apoc.convert.fromJsonMap(pvSort0.value).min) AS sortValue0") {
		t.Errorf("missing range expression; opt=%q", opt)
	}
}

func TestBuildSortClause_CustomProperty_List(t *testing.T) {
	uid := "1adeae21-ab0f-40f3-868d-36f7d40210ca"
	propTypes := sortPropertyTypes{uid: "list"}
	opt, _, _, _ := buildCatalogueItemsSortClause(
		&[]helpers.Sorting{{ID: uid, DESC: false}}, propTypes)

	if !strings.Contains(opt, "toLower(toString(pvSort0.value)) AS sortValue0") {
		t.Errorf("list type should use text expression; opt=%q", opt)
	}
}

func TestBuildSortClause_MultiColumn(t *testing.T) {
	uid1 := "1adeae21-ab0f-40f3-868d-36f7d40210ca"
	uid2 := "2bdeae21-ab0f-40f3-868d-36f7d40210cb"
	propTypes := sortPropertyTypes{uid1: "number", uid2: "text"}
	opt, ret, orderBy, params := buildCatalogueItemsSortClause(&[]helpers.Sorting{
		{ID: "name", DESC: false},
		{ID: uid1, DESC: true},
		{ID: uid2, DESC: false},
	}, propTypes)

	if !strings.Contains(opt, "pvSort0:HAS_CATALOGUE_PROPERTY") {
		t.Errorf("missing pvSort0; opt=%q", opt)
	}
	if !strings.Contains(opt, "pvSort1:HAS_CATALOGUE_PROPERTY") {
		t.Errorf("missing pvSort1; opt=%q", opt)
	}
	if !strings.Contains(opt, "toFloat(pvSort0.value) AS sortValue0") {
		t.Errorf("missing sortValue0; opt=%q", opt)
	}
	if !strings.Contains(opt, "toLower(toString(pvSort1.value)) AS sortValue1") {
		t.Errorf("missing sortValue1; opt=%q", opt)
	}
	if len(ret) != 2 || ret[0] != "sortValue0" || ret[1] != "sortValue1" {
		t.Errorf("extraReturnVars = %v", ret)
	}
	want := " ORDER BY itm.name ASC, sortValue0 DESC, sortValue1 ASC, itm.uid ASC "
	if orderBy != want {
		t.Errorf("orderBy = %q, want %q", orderBy, want)
	}
	if params["sortPropKey0"] != uid1 || params["sortPropKey1"] != uid2 {
		t.Errorf("params = %v", params)
	}
}

func TestBuildSortClause_UnknownIdSkipped(t *testing.T) {
	opt, ret, orderBy, params := buildCatalogueItemsSortClause(
		&[]helpers.Sorting{{ID: "Inletflangesize", DESC: false}}, nil)

	if opt != "" {
		t.Errorf("optionalMatchSQL not empty: %q", opt)
	}
	if len(ret) != 0 {
		t.Errorf("extraReturnVars not empty: %v", ret)
	}
	if strings.Contains(orderBy, "Inletflangesize") {
		t.Errorf("unknown id leaked into orderBy: %q", orderBy)
	}
	if orderBy != " ORDER BY lastUpdateTime DESC, itm.uid ASC " {
		t.Errorf("orderBy = %q", orderBy)
	}
	if len(params) != 0 {
		t.Errorf("params not empty: %v", params)
	}
}

func TestBuildSortClause_UnknownPropertyUidSkipped(t *testing.T) {
	uid := "1adeae21-ab0f-40f3-868d-36f7d40210ca"
	propTypes := sortPropertyTypes{}
	opt, _, orderBy, params := buildCatalogueItemsSortClause(
		&[]helpers.Sorting{{ID: uid, DESC: false}}, propTypes)

	if opt != "" {
		t.Errorf("optionalMatchSQL not empty: %q", opt)
	}
	if orderBy != " ORDER BY lastUpdateTime DESC, itm.uid ASC " {
		t.Errorf("orderBy = %q", orderBy)
	}
	if len(params) != 0 {
		t.Errorf("params not empty: %v", params)
	}
}

func TestBuildSortClause_InjectionAttempt(t *testing.T) {
	cases := []string{
		"itm.name) RETURN 1 //",
		"x; DROP DATABASE neo4j",
		"sname DESC, itm.uid ASC) UNION MATCH (n) RETURN n",
	}
	for _, id := range cases {
		opt, _, orderBy, params := buildCatalogueItemsSortClause(
			&[]helpers.Sorting{{ID: id, DESC: false}}, nil)
		if strings.Contains(opt, id) || strings.Contains(orderBy, id) {
			t.Errorf("injection leaked for %q; opt=%q orderBy=%q", id, opt, orderBy)
		}
		if len(params) != 0 {
			t.Errorf("params not empty for %q: %v", id, params)
		}
	}
}

func TestPaginationQuery_AlwaysHasTiebreaker(t *testing.T) {
	q := CatalogueItemsFiltersPaginationQuery("", "", 0, 50, nil, nil, nil)
	if !strings.Contains(q.Query, "ORDER BY lastUpdateTime DESC, itm.uid ASC") {
		t.Errorf("missing tiebreaker on default sort; query=%q", q.Query)
	}
}

func TestPaginationQuery_SupplierSort_ProjectsSname(t *testing.T) {
	q := CatalogueItemsFiltersPaginationQuery("", "", 0, 50, nil,
		&[]helpers.Sorting{{ID: "supplier", DESC: false}}, nil)

	if !strings.Contains(q.Query, "RETURN itm, cat, lastUpdateTime, lastUpdateUser, sname") {
		t.Errorf("missing sname in RETURN; query=%q", q.Query)
	}
	if !strings.Contains(q.Query, "ORDER BY sname ASC, itm.uid ASC") {
		t.Errorf("missing supplier ORDER BY; query=%q", q.Query)
	}
}

func TestPaginationQuery_CustomPropertySort_BoundParam(t *testing.T) {
	uid := "1adeae21-ab0f-40f3-868d-36f7d40210ca"
	propTypes := sortPropertyTypes{uid: "number"}
	q := CatalogueItemsFiltersPaginationQuery("", "", 0, 50, nil,
		&[]helpers.Sorting{{ID: uid, DESC: false}}, propTypes)

	if !strings.Contains(q.Query, "$sortPropKey0") {
		t.Errorf("missing $sortPropKey0 in query; query=%q", q.Query)
	}
	if q.Parameters["sortPropKey0"] != uid {
		t.Errorf("sortPropKey0 param = %v, want %v", q.Parameters["sortPropKey0"], uid)
	}
	if strings.Contains(q.Query, uid) {
		t.Errorf("uid leaked into query string instead of being parameterized; query=%q", q.Query)
	}
}
