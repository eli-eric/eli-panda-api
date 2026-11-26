package systemsService

import (
	"encoding/json"
	"fmt"
	"sort"

	"panda/apigateway/services/systems-service/models"
)

const (
	noGroupKey   = "no-group"
	noGroupName  = "General"
	notAvailable = "N/A"
)

// formatPropertyValue renders a raw property value for display; Range values are stored as
// JSON {"min": x, "max": y}. Mirrors the former FE engine (useItemPropertiesData).
func formatPropertyValue(value *string, typeName string) string {
	if typeName == "Range" {
		if value == nil || *value == "" {
			return notAvailable + " - " + notAvailable
		}
		var r struct {
			Min interface{} `json:"min"`
			Max interface{} `json:"max"`
		}
		if err := json.Unmarshal([]byte(*value), &r); err != nil {
			return *value
		}
		return rangeBound(r.Min) + " - " + rangeBound(r.Max)
	}
	if value == nil || *value == "" {
		return notAvailable
	}
	return *value
}

func rangeBound(v interface{}) string {
	if v == nil {
		return notAvailable
	}
	if f, ok := v.(float64); ok {
		return fmt.Sprintf("%g", f)
	}
	return fmt.Sprint(v)
}

func propertyGroupsOf(p models.RawCatalogueProperty) []models.CodebookRef {
	if len(p.Groups) == 0 {
		return []models.CodebookRef{{UID: noGroupKey, Name: noGroupName}}
	}
	return p.Groups
}

func typeNameOf(p models.RawCatalogueProperty) string {
	if p.Type == nil || *p.Type == "" {
		return "Unknown"
	}
	return *p.Type
}

// latestServiceProperties keeps, per property uid, the value from the newest service item.
func latestServiceProperties(items []models.RawServiceItemProperties) []models.RawCatalogueProperty {
	sorted := append([]models.RawServiceItemProperties(nil), items...)
	sort.SliceStable(sorted, func(i, j int) bool {
		a, b := sorted[i].Created, sorted[j].Created
		if a == nil || b == nil {
			return a != nil
		}
		return a.After(*b)
	})
	seen := map[string]struct{}{}
	latest := []models.RawCatalogueProperty{}
	for _, item := range sorted {
		for _, p := range item.Properties {
			if _, ok := seen[p.UID]; ok {
				continue
			}
			seen[p.UID] = struct{}{}
			latest = append(latest, p)
		}
	}
	return latest
}

// resolveCatalogueProperties groups catalogue properties and applies the latest service
// overrides; service-only properties are added with a nil catalogue value.
func resolveCatalogueProperties(raw models.RawPhysicalItemProperties) models.PhysicalItemCatalogueProperties {
	groups := map[string]*models.CataloguePropertyGroup{}
	group := func(ref models.CodebookRef) *models.CataloguePropertyGroup {
		if g, ok := groups[ref.UID]; ok {
			return g
		}
		g := &models.CataloguePropertyGroup{Key: ref.UID, Name: ref.Name, Properties: []models.ResolvedCatalogueProperty{}}
		groups[ref.UID] = g
		return g
	}

	for _, p := range raw.Catalogue {
		value := formatPropertyValue(p.Value, typeNameOf(p))
		for _, ref := range propertyGroupsOf(p) {
			g := group(ref)
			g.Properties = append(g.Properties, models.ResolvedCatalogueProperty{
				UID: p.UID, Name: p.Name, Type: typeNameOf(p), Unit: p.Unit, Value: &value,
			})
		}
	}

	hasOverridden := false
	for _, p := range latestServiceProperties(raw.ServiceItems) {
		serviceValue := formatPropertyValue(p.Value, typeNameOf(p))
		for _, ref := range propertyGroupsOf(p) {
			g := group(ref)
			idx := -1
			for i := range g.Properties {
				if g.Properties[i].UID == p.UID {
					idx = i
					break
				}
			}
			if idx < 0 {
				g.Properties = append(g.Properties, models.ResolvedCatalogueProperty{
					UID: p.UID, Name: p.Name, Type: typeNameOf(p), Unit: p.Unit, ServiceValue: &serviceValue,
				})
				continue
			}
			existing := &g.Properties[idx]
			existing.ServiceValue = &serviceValue
			existing.IsOverridden = existing.Value == nil || *existing.Value != serviceValue
			hasOverridden = hasOverridden || existing.IsOverridden
		}
	}

	result := models.PhysicalItemCatalogueProperties{Groups: []models.CataloguePropertyGroup{}, HasOverriddenProperties: hasOverridden}
	for _, g := range groups {
		sort.SliceStable(g.Properties, func(i, j int) bool { return g.Properties[i].Name < g.Properties[j].Name })
		result.Groups = append(result.Groups, *g)
	}
	sort.SliceStable(result.Groups, func(i, j int) bool {
		a, b := result.Groups[i], result.Groups[j]
		if (a.Key == noGroupKey) != (b.Key == noGroupKey) {
			return a.Key == noGroupKey
		}
		return a.Name < b.Name
	})
	return result
}
