package helpers

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

var jsonFieldNameRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

type FieldProjectionError struct {
	InvalidFields []string
	AllowedFields []string
}

func (e *FieldProjectionError) Error() string {
	return fmt.Sprintf("invalid fields: %s", strings.Join(e.InvalidFields, ","))
}

// ParseFieldsParam parses a comma-separated list of JSON field names.
// It trims whitespace, removes empty segments, and de-duplicates while preserving order.
func ParseFieldsParam(fieldsParam string) []string {
	if strings.TrimSpace(fieldsParam) == "" {
		return nil
	}

	parts := strings.Split(fieldsParam, ",")
	seen := make(map[string]struct{}, len(parts))
	result := make([]string, 0, len(parts))

	for _, raw := range parts {
		field := strings.TrimSpace(raw)
		if field == "" {
			continue
		}
		if _, ok := seen[field]; ok {
			continue
		}
		seen[field] = struct{}{}
		result = append(result, field)
	}

	return result
}

// ProjectJSONFields projects a struct (or pointer to struct) to a map that contains only
// requested top-level fields by their `json:"..."` tag names.
//
// Safety constraints:
// - Only top-level fields are supported (no dot paths)
// - Only exported struct fields with a non-empty json tag are selectable
// - Unknown/invalid fields cause an error with allowed/invalid lists
func ProjectJSONFields(value any, fields []string) (map[string]any, error) {
	if value == nil {
		return nil, fmt.Errorf("value is nil")
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields requested")
	}

	structValue := reflect.ValueOf(value)
	if structValue.Kind() == reflect.Pointer {
		if structValue.IsNil() {
			return nil, fmt.Errorf("value is nil")
		}
		structValue = structValue.Elem()
	}
	if structValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("value must be a struct or pointer to struct")
	}

	structType := structValue.Type()
	jsonNameToIndex := make(map[string]int, structType.NumField())
	allowed := make([]string, 0, structType.NumField())

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		jsonTag := field.Tag.Get("json")
		jsonName := strings.TrimSpace(strings.Split(jsonTag, ",")[0])
		if jsonName == "" || jsonName == "-" {
			continue
		}
		// Only allow identifiers to avoid weird tag contents
		if !jsonFieldNameRegexp.MatchString(jsonName) {
			continue
		}
		if _, exists := jsonNameToIndex[jsonName]; exists {
			// Ignore duplicates; should not happen in normal models.
			continue
		}
		jsonNameToIndex[jsonName] = i
		allowed = append(allowed, jsonName)
	}

	sort.Strings(allowed)

	invalid := make([]string, 0)
	result := make(map[string]any, len(fields))

	for _, requested := range fields {
		if strings.Contains(requested, ".") {
			invalid = append(invalid, requested)
			continue
		}
		if !jsonFieldNameRegexp.MatchString(requested) {
			invalid = append(invalid, requested)
			continue
		}

		index, ok := jsonNameToIndex[requested]
		if !ok {
			invalid = append(invalid, requested)
			continue
		}

		fieldValue := structValue.Field(index)
		result[requested] = fieldValue.Interface()
	}

	if len(invalid) > 0 {
		sort.Strings(invalid)
		return nil, &FieldProjectionError{InvalidFields: invalid, AllowedFields: allowed}
	}

	return result, nil
}
