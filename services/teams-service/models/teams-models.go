package models

// Team is the core team node projection (name/code/description), facility-scoped.
type Team struct {
	UID         string `json:"uid" neo4j:"key,uid"`
	Name        string `json:"name" neo4j:"prop,name"`
	Code        string `json:"code" neo4j:"prop,code"`
	Description string `json:"description" neo4j:"prop,description"`
}

// TeamListItem is a team plus its member count, returned in the flat list endpoint.
type TeamListItem struct {
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	MemberCount int    `json:"memberCount"`
}

// TeamMember is a user summary as returned in a team's member list. Disabled users
// are included (membership persists) and flagged via IsEnabled.
type TeamMember struct {
	UID       string `json:"uid"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	IsEnabled bool   `json:"isEnabled"`
}

// TeamDetail is a team plus its full inline member list.
type TeamDetail struct {
	UID         string       `json:"uid"`
	Name        string       `json:"name"`
	Code        string       `json:"code"`
	Description string       `json:"description"`
	Members     []TeamMember `json:"members"`
}

// TeamCreateRequest is the POST body. Name is required; code/description optional.
type TeamCreateRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// TeamUpdateRequest is the PUT (full replace) body. Name is required; code/description
// are overwritten (description may be null to clear).
type TeamUpdateRequest struct {
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
}

// Optional signals presence of a field in a PATCH payload. A nil *Optional pointer
// means the JSON key was absent; a non-nil pointer with Value == nil means the key was
// explicitly null (clear operation).
type Optional[T any] struct {
	Value *T
}

// PatchTeamFields is the parsed PATCH body. Name uses a plain pointer (nil = absent);
// nullable fields use *Optional[T] (nil = absent; non-nil with Value=nil = explicit null clear).
type PatchTeamFields struct {
	Name        *string
	Code        *Optional[string]
	Description *Optional[string]
}

// TeamMembersRequest is the body for add-members (POST) and replace-members (PUT).
type TeamMembersRequest struct {
	UserUids []string `json:"userUids"`
}
