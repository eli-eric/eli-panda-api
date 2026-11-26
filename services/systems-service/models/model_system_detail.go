package models

import "time"

// CodebookRef is a nullable-code codebook projection used by the system detail endpoints.
type CodebookRef struct {
	UID  string  `json:"uid"`
	Name string  `json:"name"`
	Code *string `json:"code,omitempty"`
}

// EmployeeRef is an employee projection; Name is fullName with a "last first" fallback.
type EmployeeRef struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type SystemPathItem struct {
	UID         string  `json:"uid"`
	Name        string  `json:"name"`
	SystemLevel *string `json:"systemLevel"`
}

// SystemDetail is the response of GET/PATCH /v1/system/{uid}. It is a superset of the
// fields the legacy models.System detail returned, so existing consumers keep working.
type SystemDetail struct {
	UID                    string                    `json:"uid"`
	Name                   string                    `json:"name"`
	Description            *string                   `json:"description"`
	SystemCode             *string                   `json:"systemCode"`
	SystemAlias            *string                   `json:"systemAlias"`
	SystemLevel            *string                   `json:"systemLevel"`
	LastUpdateTime         *time.Time                `json:"lastUpdateTime"`
	ParentUID              *string                   `json:"parentUid"`
	ParentPath             []SystemPathItem          `json:"parentPath"`
	SystemType             *CodebookRef              `json:"systemType"`
	Location               *CodebookRef              `json:"location"`
	Zone                   *CodebookRef              `json:"zone"`
	Importance             *CodebookRef              `json:"importance"`
	Attribute              *CodebookRef              `json:"attribute"`
	Responsible            *EmployeeRef              `json:"responsible"`
	ResponsibleTeam        *CodebookRef              `json:"responsibleTeam"`
	Owner                  *EmployeeRef              `json:"owner"`
	Operators              []EmployeeRef             `json:"operators"`
	MaintainedBy           []EmployeeRef             `json:"maintainedBy"`
	MiniImageUrl           []string                  `json:"miniImageUrl"`
	MinimalSpareParstCount *float64                  `json:"minimalSpareParstCount"`
	SparePartsCoverageSum  *float64                  `json:"sparePartsCoverageSum"`
	SpCoverage             *float64                  `json:"sp_coverage"`
	PhysicalItem           *SystemDetailPhysicalItem `json:"physicalItem"`
	Statistics             SystemDetailStatistics    `json:"statistics"`
}

type SystemDetailPhysicalItem struct {
	UID             string                     `json:"uid"`
	Name            string                     `json:"name"`
	EUN             *string                    `json:"eun"`
	SerialNumber    *string                    `json:"serialNumber"`
	Notes           *string                    `json:"notes"`
	Price           *float64                   `json:"price"`
	Currency        *string                    `json:"currency"`
	PrintEUN        *bool                      `json:"printEUN"`
	ItemUsage       *CodebookRef               `json:"itemUsage"`
	ConditionStatus *CodebookRef               `json:"conditionStatus"`
	CatalogueItem   *SystemDetailCatalogueItem `json:"catalogueItem"`
}

type SystemDetailCatalogueItem struct {
	UID             string       `json:"uid"`
	Name            string       `json:"name"`
	CatalogueNumber *string      `json:"catalogueNumber"`
	Category        *CodebookRef `json:"category"`
	Supplier        *CodebookRef `json:"supplier"`
}

type SystemDetailStatistics struct {
	SubsystemsCount        int64    `json:"subsystemsCount"`
	DirectSubsystemsCount  int64    `json:"directSubsystemsCount"`
	SparePartsCount        int64    `json:"sparePartsCount"`
	SparesForCount         int64    `json:"sparesForCount"`
	MinimalSpareParstCount *float64 `json:"minimalSpareParstCount"`
	SparePartsCoverageSum  *float64 `json:"sparePartsCoverageSum"`
	SpCoverage             *float64 `json:"sp_coverage"`
}

// UIDRef is the {uid} reference used in PATCH bodies.
type UIDRef struct {
	UID string `json:"uid"`
}

// Optional signals presence of a field in a PATCH payload. A nil *Optional pointer means the
// JSON key was absent; a non-nil pointer with Value == nil means explicit null (clear).
type Optional[T any] struct {
	Value *T
}

// PatchToken carries the optimistic-concurrency token. Present=false means the client sent no
// token (no check); Present=true with Value=nil means the client expects a legacy null token.
type PatchToken struct {
	Present bool
	Value   *time.Time
}

// PatchSystemFields is the parsed PATCH /v1/system/{uid} body.
type PatchSystemFields struct {
	Token                  PatchToken
	Name                   *string
	Description            *Optional[string]
	SystemCode             *Optional[string]
	SystemAlias            *Optional[string]
	SystemLevel            *Optional[string]
	SystemType             *Optional[UIDRef]
	Location               *Optional[UIDRef]
	Zone                   *Optional[UIDRef]
	Importance             *Optional[UIDRef]
	Attribute              *Optional[UIDRef]
	Responsible            *Optional[UIDRef]
	ResponsibleTeam        *Optional[UIDRef]
	Owner                  *Optional[UIDRef]
	Operators              *[]UIDRef
	MaintainedBy           *[]UIDRef
	MinimalSpareParstCount *Optional[float64]
}

// PatchPhysicalItemFields is the parsed PATCH /v1/system/{uid}/physical-item body.
type PatchPhysicalItemFields struct {
	Token           PatchToken
	SerialNumber    *Optional[string]
	Notes           *Optional[string]
	ItemUsage       *Optional[UIDRef]
	ConditionStatus *Optional[UIDRef]
}

// CreateSystemFields is the parsed POST /v1/systems body: the PATCH fields plus an optional parent.
type CreateSystemFields struct {
	ParentUID *string
	Fields    PatchSystemFields
}

type CreateSubsystemRequest struct {
	Name        string `json:"name"`
	SystemLevel string `json:"systemLevel"`
}

type SystemCodeGenerateRequest struct {
	LastUpdateTime *time.Time `json:"lastUpdateTime"`
}

// SystemSparePartsResponse is the response of GET /v1/system/{uid}/spare-parts.
type SystemSparePartsResponse struct {
	MinimalSpareParstCount *float64           `json:"minimalSpareParstCount"`
	SparePartsCoverageSum  *float64           `json:"sparePartsCoverageSum"`
	SpCoverage             *float64           `json:"sp_coverage"`
	SpareParts             []SpareRelationRef `json:"spareParts"`
	SpareFor               []SpareRelationRef `json:"spareFor"`
}

type SpareRelationRef struct {
	Coverage *float64       `json:"coverage"`
	System   SpareSystemRef `json:"system"`
}

type SpareSystemRef struct {
	UID          string                   `json:"uid"`
	Name         string                   `json:"name"`
	SystemCode   *string                  `json:"systemCode"`
	SystemLevel  *string                  `json:"systemLevel"`
	ParentPath   []SystemSimpleResponse   `json:"parentPath"`
	Location     *CodebookRef             `json:"location"`
	PhysicalItem *SpareSystemPhysicalItem `json:"physicalItem"`
}

type SpareSystemPhysicalItem struct {
	UID           string                 `json:"uid"`
	EUN           *string                `json:"eun"`
	ItemUsage     *CodebookRef           `json:"itemUsage"`
	CatalogueItem *SpareSystemCatalogRef `json:"catalogueItem"`
}

type SpareSystemCatalogRef struct {
	UID             string  `json:"uid"`
	CatalogueNumber *string `json:"catalogueNumber"`
}

// PhysicalItemCatalogueProperties is the response of GET /v1/physical-item/{uid}/catalogue-properties.
type PhysicalItemCatalogueProperties struct {
	Groups                  []CataloguePropertyGroup `json:"groups"`
	HasOverriddenProperties bool                     `json:"hasOverriddenProperties"`
}

type CataloguePropertyGroup struct {
	Key        string                      `json:"key"`
	Name       string                      `json:"name"`
	Properties []ResolvedCatalogueProperty `json:"properties"`
}

type ResolvedCatalogueProperty struct {
	UID          string  `json:"uid"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Unit         *string `json:"unit"`
	Value        *string `json:"value"`
	ServiceValue *string `json:"serviceValue"`
	IsOverridden bool    `json:"isOverridden"`
}

// RawCatalogueProperty is one catalogue or service property value as read from the DB,
// before override resolution.
type RawCatalogueProperty struct {
	UID    string        `json:"uid"`
	Name   string        `json:"name"`
	Type   *string       `json:"type"`
	Unit   *string       `json:"unit"`
	Value  *string       `json:"value"`
	Groups []CodebookRef `json:"groups"`
}

type RawServiceItemProperties struct {
	Created    *time.Time             `json:"created"`
	Properties []RawCatalogueProperty `json:"properties"`
}

type RawPhysicalItemProperties struct {
	Catalogue    []RawCatalogueProperty     `json:"catalogue"`
	ServiceItems []RawServiceItemProperties `json:"serviceItems"`
}

// PhysicalItemOrders is the response of GET /v1/physical-item/{uid}/orders.
type PhysicalItemOrders struct {
	Order        *PhysicalItemOrder       `json:"order"`
	ServiceItems []PhysicalItemServiceRef `json:"serviceItems"`
}

type PhysicalItemOrder struct {
	UID         string     `json:"uid"`
	Name        string     `json:"name"`
	OrderDate   *time.Time `json:"orderDate"`
	IsDelivered *bool      `json:"isDelivered"`
}

type PhysicalItemServiceRef struct {
	UID         string             `json:"uid"`
	Name        string             `json:"name"`
	IsDelivered bool               `json:"isDelivered"`
	Created     *time.Time         `json:"created"`
	Order       *PhysicalItemOrder `json:"order"`
}

// ConflictResponse is the 409 body of the system PATCH family.
type ConflictResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
