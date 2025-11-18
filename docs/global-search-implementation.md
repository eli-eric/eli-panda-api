# Global Full-Text Search - Complete Implementation Plan

> **Purpose:** Highly optimized global search endpoint for Command Palette UI
> **Date:** 2025-01-12
> **Status:** Design Phase

---

## 📋 Table of Contents

1. [Overview and Requirements](#overview-and-requirements)
2. [API Design](#api-design)
3. [TypeScript Definitions](#typescript-definitions)
4. [Database Indexes](#database-indexes)
5. [Neo4j Query Implementation](#neo4j-query-implementation)
6. [Go Service Implementation](#go-service-implementation)
7. [Performance Expectations](#performance-expectations)
8. [Key File References](#key-file-references)
9. [Implementation Steps](#implementation-steps)

---

## 🎯 Overview and Requirements

### Business Requirements
- **Use Case:** Command Palette - user starts typing, results appear immediately
- **Node Types:** System, CatalogueItem, Order
- **Response Time:** < 200ms (ideally < 100ms)
- **User Experience:** Minimal but sufficient context to identify records

### Technical Requirements
- Full-text search across all 3 node types
- Facility-based filtering (security)
- Relevance scoring
- Paginated results (limit per type)
- Optimized Neo4j queries with index utilization

---

## 🔌 API Design

### Endpoint

```
GET /v1/search
```

### Query Parameters

| Parameter | Type | Required | Default | Description |
|----------|-----|----------|---------|-------|
| `q` | string | ✅ | - | Search term (min 2 chars) |
| `types` | string | ❌ | "system,catalogueItem,order" | Comma-separated: system, catalogueItem, order |
| `limit` | int | ❌ | 10 | Max results per type (max: 50) |

### Usage Examples

```http
# Search across all types
GET /v1/search?q=laser

# Search only in systems and catalogue items
GET /v1/search?q=laser&types=system,catalogueItem&limit=15

# Search in orders
GET /v1/search?q=2024-123&types=order

# Search by EUN
GET /v1/search?q=ELI-BM-12345
```

### HTTP Status Codes

| Code | Meaning |
|------|--------|
| 200 | Success - results found (or empty) |
| 400 | Invalid request (missing q, search term too short) |
| 401 | Unauthorized (missing/invalid JWT) |
| 403 | Forbidden (missing required role) |
| 500 | Internal server error |

---

## 📦 TypeScript Definitions

### Response Interface

```typescript
/**
 * Global search response
 */
export interface SearchResponse {
  /** Search results grouped by type */
  results: SearchResult[];

  /** Original search term */
  searchTerm: string;

  /** Total number of results across all types */
  totalCount: number;

  /** Query execution time in milliseconds */
  executionTimeMs: number;

  /** Results breakdown by type */
  resultsByType: {
    system: number;
    catalogueItem: number;
    order: number;
  };
}

/**
 * Single search result
 */
export interface SearchResult {
  /** Unique identifier of the node */
  uid: string;

  /** Type of the result */
  type: 'system' | 'catalogueItem' | 'order';

  /** Display name */
  name: string;

  /** Field that matched the search term */
  matchedField: string;

  /** Value that matched (e.g., system code, catalogue number) */
  matchedValue: string;

  /** Relevance score (1-10, higher = more relevant) */
  relevanceScore: number;

  /** Type-specific context information */
  context: SystemContext | CatalogueItemContext | OrderContext;
}

/**
 * Context for System results
 */
export interface SystemContext {
  /** System code identifier */
  systemCode?: string;

  /** System type name */
  systemType?: string;

  /** Location path (e.g., "Building A > Floor 2 > Room 201") */
  locationPath?: string;

  /** Facility name */
  facilityName?: string;

  /** Hierarchical parent path */
  parentPath: ParentPathItem[];

  /** Whether system has subsystems */
  hasSubsystems: boolean;

  /** System status */
  status?: string;

  /** System level (TECHNOLOGY_UNIT, KEY_SYSTEMS, etc.) */
  systemLevel?: string;
}

/**
 * Parent path item in hierarchy
 */
export interface ParentPathItem {
  uid: string;
  name: string;
}

/**
 * Context for CatalogueItem results
 */
export interface CatalogueItemContext {
  /** Catalogue item number */
  catalogueNumber?: string;

  /** Category hierarchy from root to leaf */
  categoryHierarchy: string[];

  /** Manufacturer name */
  manufacturerName?: string;

  /** Supplier name */
  supplierName?: string;

  /** Item description (truncated) */
  description?: string;

  /** Number of physical items in stock */
  inStock: number;

  /** Image URL (mini) */
  imageUrl?: string;
}

/**
 * Context for Order results
 */
export interface OrderContext {
  /** Order number */
  orderNumber?: string;

  /** Order date in ISO format */
  orderDate?: string;

  /** Order status name */
  statusName?: string;

  /** Order status code */
  statusCode?: string;

  /** Supplier name */
  supplierName?: string;

  /** Requestor full name */
  requestorName?: string;

  /** Total number of items in order */
  totalItemsCount: number;

  /** Delivery status percentage (0-100) */
  deliveryStatus?: number;

  /** Request number */
  requestNumber?: string;
}
```

### Example Response

```typescript
// Example: searching for "laser"
const response: SearchResponse = {
  searchTerm: "laser",
  totalCount: 3,
  executionTimeMs: 45,
  resultsByType: {
    system: 1,
    catalogueItem: 2,
    order: 0
  },
  results: [
    {
      uid: "sys-123-456",
      type: "system",
      name: "Laser System Alpha",
      matchedField: "name",
      matchedValue: "Laser System Alpha",
      relevanceScore: 9.5,
      context: {
        systemCode: "ELI-BM-LS-001",
        systemType: "Laser System",
        locationPath: "Building A > Lab 3",
        facilityName: "ELI Beamlines",
        parentPath: [
          { uid: "parent-1", name: "Main Building" },
          { uid: "parent-2", name: "Laser Hall" }
        ],
        hasSubsystems: true,
        status: "active",
        systemLevel: "TECHNOLOGY_UNIT"
      }
    },
    {
      uid: "cat-789-012",
      type: "catalogueItem",
      name: "CO2 Laser Module",
      matchedField: "catalogueNumber",
      matchedValue: "CAT-2024-LS-042",
      relevanceScore: 8.2,
      context: {
        catalogueNumber: "CAT-2024-LS-042",
        categoryHierarchy: ["Equipment", "Lasers", "CO2 Lasers"],
        manufacturerName: "LaserTech Inc.",
        supplierName: "Scientific Supplies",
        description: "High-power CO2 laser module for...",
        inStock: 3,
        imageUrl: "https://..."
      }
    }
  ]
};
```

---

## 🗄️ Database Indexes

### Existing Indexes ✅

```cypher
// System full-text index
CREATE FULLTEXT INDEX searchIndexSystems
FOR (n:System) ON EACH [n.name, n.description, n.systemCode, n.systemAlias]
OPTIONS { indexConfig: { `fulltext.analyzer`: 'standard' }};

// Order full-text index
CREATE FULLTEXT INDEX searchIndexOrders
FOR (n:Order) ON EACH [n.name, n.orderNumber, n.requestNumber, n.contractNumber, n.notes]
OPTIONS { indexConfig: { `fulltext.analyzer`: 'standard' }};

// Order UID index
CREATE INDEX ordersUidIndex FOR (o:Order) ON o.uid;
```

**Files:**
- `db/neo4j/migrations/20230628124422_add_fulltext_systems.up.cypher`
- `db/neo4j/migrations/20230502135437_orders_search_indexes.up.cypher`

### New Indexes - To Be Created 🆕

**File:** `db/neo4j/migrations/20250112120000_add_global_search_indexes.up.cypher`

```cypher
// ============================================
// FULLTEXT INDEX FOR CATALOGUE ITEMS
// ============================================
CREATE FULLTEXT INDEX searchIndexCatalogueItems IF NOT EXISTS
FOR (n:CatalogueItem)
ON EACH [n.name, n.catalogueNumber, n.description]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'standard',
    `fulltext.eventually_consistent`: false
  }
};

// ============================================
// PROPERTY INDEXES FOR FAST FILTERING
// ============================================

// System property indexes
CREATE INDEX systemDeletedIndex IF NOT EXISTS
  FOR (s:System) ON (s.deleted);

CREATE INDEX systemUidIndex IF NOT EXISTS
  FOR (s:System) ON (s.uid);

// CatalogueItem property indexes
CREATE INDEX catalogueItemUidIndex IF NOT EXISTS
  FOR (ci:CatalogueItem) ON (ci.uid);

CREATE INDEX catalogueItemDeletedIndex IF NOT EXISTS
  FOR (ci:CatalogueItem) ON (ci.deleted);

CREATE INDEX catalogueItemNumberIndex IF NOT EXISTS
  FOR (ci:CatalogueItem) ON (ci.catalogueNumber);

// Order property indexes
CREATE INDEX orderDeletedIndex IF NOT EXISTS
  FOR (o:Order) ON (o.deleted);

CREATE INDEX orderNumberIndex IF NOT EXISTS
  FOR (o:Order) ON (o.orderNumber);

// Item property indexes (for EUN search)
CREATE INDEX itemEunIndex IF NOT EXISTS
  FOR (i:Item) ON (i.eun);

CREATE INDEX itemUidIndex IF NOT EXISTS
  FOR (i:Item) ON (i.uid);

// Facility code index (critical for security)
CREATE INDEX facilityCodeIndex IF NOT EXISTS
  FOR (f:Facility) ON (f.code);

// Location, Zone, SystemType indexes (for context data)
CREATE INDEX locationUidIndex IF NOT EXISTS
  FOR (l:Location) ON (l.uid);

CREATE INDEX zoneUidIndex IF NOT EXISTS
  FOR (z:Zone) ON (z.uid);

CREATE INDEX systemTypeUidIndex IF NOT EXISTS
  FOR (st:SystemType) ON (st.uid);
```

**File:** `db/neo4j/migrations/20250112120000_add_global_search_indexes.down.cypher`

```cypher
DROP INDEX searchIndexCatalogueItems IF EXISTS;
DROP INDEX systemDeletedIndex IF EXISTS;
DROP INDEX systemUidIndex IF EXISTS;
DROP INDEX catalogueItemUidIndex IF EXISTS;
DROP INDEX catalogueItemDeletedIndex IF EXISTS;
DROP INDEX catalogueItemNumberIndex IF EXISTS;
DROP INDEX orderDeletedIndex IF EXISTS;
DROP INDEX orderNumberIndex IF EXISTS;
DROP INDEX itemEunIndex IF EXISTS;
DROP INDEX itemUidIndex IF EXISTS;
DROP INDEX facilityCodeIndex IF EXISTS;
DROP INDEX locationUidIndex IF EXISTS;
DROP INDEX zoneUidIndex IF EXISTS;
DROP INDEX systemTypeUidIndex IF EXISTS;
```

### Index Performance Impact

| Index Type | Purpose | Performance Gain | Use Case |
|-----------|---------|------------------|----------|
| **FULLTEXT** | Text search with tokenization | 50-100x | Searching in name, description, catalogueNumber |
| **Property (uid)** | Direct lookup | 10-50x | WHERE node.uid = $uid |
| **Property (deleted)** | Boolean filtering | 5-20x | WHERE node.deleted = false |
| **Property (code/number)** | String exact match | 10-30x | WHERE node.orderNumber = $num |

### Expected Query Performance

| Dataset Size | Without Indexes | With Indexes | Improvement |
|--------------|----------------|--------------|-------------|
| 1,000 records | 50-100ms | 5-10ms | **10x** |
| 10,000 records | 500ms-1s | 10-20ms | **50x** |
| 100,000 records | 5-10s | 20-50ms | **200x** |
| 1,000,000 records | 50s+ | 50-150ms | **500x+** |

---

## 🔍 Neo4j Query Implementation

### Main Search Query

```cypher
// Global search query - optimized with UNION subqueries
CALL {
  // =========================================
  // SUBQUERY 1: SYSTEM SEARCH
  // =========================================
  CALL db.index.fulltext.queryNodes('searchIndexSystems', $fulltextSearch)
  YIELD node AS sys, score
  WHERE sys:System
    AND sys.deleted = false

  // Facility filtering (security)
  MATCH (sys)-[:BELONGS_TO_FACILITY]->(f:Facility {code: $facilityCode})

  // Context data - optimized with OPTIONAL MATCH
  OPTIONAL MATCH (parents)-[:HAS_SUBSYSTEM*1..10]->(sys)
  OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc:Location)
  OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone:Zone)
  OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st:SystemType)
  OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM]->(subsys)

  WITH sys, score, loc, zone, st,
       collect(DISTINCT {uid: parents.uid, name: parents.name}) as parentPath,
       count(DISTINCT subsys) > 0 as hasSubsystems

  RETURN
    sys.uid as uid,
    'system' as type,
    sys.name as name,
    CASE
      WHEN toLower(sys.systemCode) CONTAINS $search THEN 'systemCode'
      WHEN toLower(sys.name) CONTAINS $search THEN 'name'
      ELSE 'description'
    END as matchedField,
    coalesce(sys.systemCode, sys.name) as matchedValue,
    score as relevanceScore,
    {
      systemCode: sys.systemCode,
      systemType: st.name,
      locationPath: loc.code,
      facilityName: f.code,
      parentPath: [p IN parentPath | {uid: p.uid, name: p.name}],
      hasSubsystems: hasSubsystems,
      status: sys.status,
      systemLevel: sys.systemLevel
    } as context
  ORDER BY score DESC
  LIMIT $limit

  UNION

  // =========================================
  // SUBQUERY 2: CATALOGUE ITEM SEARCH
  // =========================================
  CALL db.index.fulltext.queryNodes('searchIndexCatalogueItems', $fulltextSearch)
  YIELD node AS ci, score
  WHERE ci:CatalogueItem
    AND ci.deleted = false

  // Context data
  OPTIONAL MATCH (ci)-[:BELONGS_TO_CATEGORY]->(cat:CatalogueCategory)
  OPTIONAL MATCH (cat)<-[:HAS_SUBCATEGORY*0..10]-(rootCat:CatalogueCategory)
    WHERE NOT (rootCat)<-[:HAS_SUBCATEGORY]-()
  OPTIONAL MATCH (ci)-[:HAS_MANUFACTURER]->(man:Manufacturer)
  OPTIONAL MATCH (ci)-[:HAS_SUPPLIER]->(sup:Supplier)
  OPTIONAL MATCH (ci)<-[:IS_BASED_ON]-(items:Item)

  WITH ci, score, cat, man, sup,
       collect(DISTINCT rootCat.name) + [cat.name] as categoryHierarchy,
       count(DISTINCT items) as stockCount

  RETURN
    ci.uid as uid,
    'catalogueItem' as type,
    ci.name as name,
    CASE
      WHEN toLower(ci.catalogueNumber) CONTAINS $search THEN 'catalogueNumber'
      WHEN toLower(ci.name) CONTAINS $search THEN 'name'
      ELSE 'description'
    END as matchedField,
    ci.catalogueNumber as matchedValue,
    score as relevanceScore,
    {
      catalogueNumber: ci.catalogueNumber,
      categoryHierarchy: categoryHierarchy,
      manufacturerName: man.name,
      supplierName: sup.name,
      description: substring(ci.description, 0, 200),
      inStock: stockCount,
      imageUrl: ci.miniImageUrl
    } as context
  ORDER BY score DESC
  LIMIT $limit

  UNION

  // =========================================
  // SUBQUERY 3: ORDER SEARCH
  // =========================================
  CALL db.index.fulltext.queryNodes('searchIndexOrders', $fulltextSearch)
  YIELD node AS o, score
  WHERE o:Order
    AND o.deleted = false

  // Facility filtering
  MATCH (o)-[:BELONGS_TO_FACILITY]->(f:Facility {code: $facilityCode})

  // Context data
  OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(status:OrderStatus)
  OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(sup:Supplier)
  OPTIONAL MATCH (o)-[:HAS_REQUESTOR]->(req:Employee)
  OPTIONAL MATCH (o)-[:HAS_ORDER_LINE]->(items:Item)

  WITH o, score, status, sup, req,
       count(DISTINCT items) as itemCount

  RETURN
    o.uid as uid,
    'order' as type,
    o.name as name,
    CASE
      WHEN toLower(o.orderNumber) CONTAINS $search THEN 'orderNumber'
      WHEN toLower(o.requestNumber) CONTAINS $search THEN 'requestNumber'
      ELSE 'name'
    END as matchedField,
    o.orderNumber as matchedValue,
    score as relevanceScore,
    {
      orderNumber: o.orderNumber,
      orderDate: toString(o.orderDate),
      statusName: status.name,
      statusCode: status.code,
      supplierName: sup.name,
      requestorName: req.firstName + ' ' + req.lastName,
      totalItemsCount: itemCount,
      deliveryStatus: o.deliveryStatus,
      requestNumber: o.requestNumber
    } as context
  ORDER BY score DESC
  LIMIT $limit
}
RETURN uid, type, name, matchedField, matchedValue, relevanceScore, context
ORDER BY relevanceScore DESC, type, name
```

### Query Parameters

```go
parameters := map[string]interface{}{
    "search":          strings.ToLower(searchTerm),
    "fulltextSearch":  helpers.GetFullTextSearchString(searchTerm),
    "facilityCode":    facilityCode,
    "limit":           limit,
}
```

### Query Optimization Techniques

1. **UNION Strategy**: Each subquery runs in parallel
2. **Early LIMIT**: Limit in each UNION branch, not at the end
3. **Index Hints**: CALL db.index.fulltext.queryNodes utilizes index automatically
4. **Shallow OPTIONAL MATCH**: Max 10 level depth for parent hierarchy
5. **WITH Staging**: Reduce intermediate results before final RETURN
6. **Property Index Coverage**: WHERE deleted = false utilizes property index

---

## 💻 Go Service Implementation

### File Structure

```
services/search-service/
├── search-service.go          # Service layer (business logic)
├── search-handlers.go         # HTTP handlers
├── search-routes.go           # Route mapping
├── search-db-queries.go       # Neo4j Cypher queries
└── models/
    └── search-models.go       # Data models
```

### 1. Service Interface & Implementation

**File:** `services/search-service/search-service.go`

```go
package searchService

import (
    "errors"
    "panda/apigateway/helpers"
    "panda/apigateway/services/search-service/models"
    "strings"
    "time"

    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
    "github.com/rs/zerolog/log"
)

type SearchService struct {
    neo4jDriver *neo4j.Driver
}

type ISearchService interface {
    GlobalSearch(searchTerm string, types []string, limit int, facilityCode string) (result models.SearchResponse, err error)
}

func NewSearchService(driver *neo4j.Driver) ISearchService {
    return &SearchService{neo4jDriver: driver}
}

func (svc *SearchService) GlobalSearch(searchTerm string, types []string, limit int, facilityCode string) (result models.SearchResponse, err error) {
    startTime := time.Now()

    // Validation
    if len(searchTerm) < 2 {
        return result, errors.New("search term must be at least 2 characters")
    }

    if len(searchTerm) > 100 {
        return result, errors.New("search term too long (max 100 characters)")
    }

    // Sanitize search term
    searchTerm = strings.TrimSpace(searchTerm)

    // Default types if not specified
    if len(types) == 0 {
        types = []string{"system", "catalogueItem", "order"}
    }

    // Validate limit
    if limit <= 0 {
        limit = 10
    }
    if limit > 50 {
        limit = 50
    }

    session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
    defer session.Close()

    query := GetGlobalSearchQuery(searchTerm, types, limit, facilityCode)
    results, err := helpers.GetNeo4jArrayOfNodes[models.SearchResult](session, query)

    if err != nil {
        log.Error().Err(err).Str("searchTerm", searchTerm).Msg("GlobalSearch query failed")
        return result, err
    }

    // Process results
    helpers.ProcessArrayResult(&results, err)

    executionTime := time.Since(startTime).Milliseconds()

    // Build response
    result = models.SearchResponse{
        SearchTerm:      searchTerm,
        Results:         results,
        TotalCount:      len(results),
        ExecutionTimeMs: executionTime,
        ResultsByType:   countResultsByType(results),
    }

    log.Info().
        Str("searchTerm", searchTerm).
        Int("totalResults", result.TotalCount).
        Int64("executionTimeMs", executionTime).
        Msg("GlobalSearch completed")

    return result, nil
}

func countResultsByType(results []models.SearchResult) models.ResultsByType {
    counts := models.ResultsByType{}
    for _, r := range results {
        switch r.Type {
        case "system":
            counts.System++
        case "catalogueItem":
            counts.CatalogueItem++
        case "order":
            counts.Order++
        }
    }
    return counts
}
```

### 2. Models

**File:** `services/search-service/models/search-models.go`

```go
package models

type SearchResponse struct {
    SearchTerm      string          `json:"searchTerm"`
    Results         []SearchResult  `json:"results"`
    TotalCount      int             `json:"totalCount"`
    ExecutionTimeMs int64           `json:"executionTimeMs"`
    ResultsByType   ResultsByType   `json:"resultsByType"`
}

type ResultsByType struct {
    System        int `json:"system"`
    CatalogueItem int `json:"catalogueItem"`
    Order         int `json:"order"`
}

type SearchResult struct {
    UID            string                 `json:"uid"`
    Type           string                 `json:"type"`
    Name           string                 `json:"name"`
    MatchedField   string                 `json:"matchedField"`
    MatchedValue   string                 `json:"matchedValue"`
    RelevanceScore float64                `json:"relevanceScore"`
    Context        map[string]interface{} `json:"context"`
}

// Context structures (can also be explicitly defined)
type SystemContext struct {
    SystemCode   *string          `json:"systemCode,omitempty"`
    SystemType   *string          `json:"systemType,omitempty"`
    LocationPath *string          `json:"locationPath,omitempty"`
    FacilityName *string          `json:"facilityName,omitempty"`
    ParentPath   []ParentPathItem `json:"parentPath"`
    HasSubsystems bool            `json:"hasSubsystems"`
    Status       *string          `json:"status,omitempty"`
    SystemLevel  *string          `json:"systemLevel,omitempty"`
}

type ParentPathItem struct {
    UID  string `json:"uid"`
    Name string `json:"name"`
}

type CatalogueItemContext struct {
    CatalogueNumber   *string  `json:"catalogueNumber,omitempty"`
    CategoryHierarchy []string `json:"categoryHierarchy"`
    ManufacturerName  *string  `json:"manufacturerName,omitempty"`
    SupplierName      *string  `json:"supplierName,omitempty"`
    Description       *string  `json:"description,omitempty"`
    InStock           int      `json:"inStock"`
    ImageUrl          *string  `json:"imageUrl,omitempty"`
}

type OrderContext struct {
    OrderNumber      *string  `json:"orderNumber,omitempty"`
    OrderDate        *string  `json:"orderDate,omitempty"`
    StatusName       *string  `json:"statusName,omitempty"`
    StatusCode       *string  `json:"statusCode,omitempty"`
    SupplierName     *string  `json:"supplierName,omitempty"`
    RequestorName    *string  `json:"requestorName,omitempty"`
    TotalItemsCount  int      `json:"totalItemsCount"`
    DeliveryStatus   *int     `json:"deliveryStatus,omitempty"`
    RequestNumber    *string  `json:"requestNumber,omitempty"`
}
```

### 3. Database Queries

**File:** `services/search-service/search-db-queries.go`

```go
package searchService

import (
    "fmt"
    "panda/apigateway/helpers"
    "strings"
)

func GetGlobalSearchQuery(searchTerm string, types []string, limit int, facilityCode string) helpers.DatabaseQuery {
    result := helpers.DatabaseQuery{}
    result.Parameters = make(map[string]interface{})
    result.Parameters["search"] = strings.ToLower(searchTerm)
    result.Parameters["fulltextSearch"] = helpers.GetFullTextSearchString(searchTerm)
    result.Parameters["facilityCode"] = facilityCode
    result.Parameters["limit"] = limit

    // Build UNION queries based on requested types
    var unions []string

    for _, nodeType := range types {
        switch nodeType {
        case "system":
            unions = append(unions, getSystemSearchSubquery())
        case "catalogueItem":
            unions = append(unions, getCatalogueItemSearchSubquery())
        case "order":
            unions = append(unions, getOrderSearchSubquery())
        }
    }

    result.Query = fmt.Sprintf(`
        CALL {
            %s
        }
        RETURN uid, type, name, matchedField, matchedValue, relevanceScore, context
        ORDER BY relevanceScore DESC, type, name
    `, strings.Join(unions, "\nUNION\n"))

    result.ReturnAlias = "uid"
    return result
}

func getSystemSearchSubquery() string {
    return `
    CALL db.index.fulltext.queryNodes('searchIndexSystems', $fulltextSearch)
    YIELD node AS sys, score
    WHERE sys:System AND sys.deleted = false

    MATCH (sys)-[:BELONGS_TO_FACILITY]->(f:Facility {code: $facilityCode})

    OPTIONAL MATCH (parents)-[:HAS_SUBSYSTEM*1..10]->(sys)
    OPTIONAL MATCH (sys)-[:HAS_LOCATION]->(loc:Location)
    OPTIONAL MATCH (sys)-[:HAS_ZONE]->(zone:Zone)
    OPTIONAL MATCH (sys)-[:HAS_SYSTEM_TYPE]->(st:SystemType)
    OPTIONAL MATCH (sys)-[:HAS_SUBSYSTEM]->(subsys)

    WITH sys, f, score, loc, zone, st,
         collect(DISTINCT {uid: parents.uid, name: parents.name}) as parentPath,
         count(DISTINCT subsys) > 0 as hasSubsystems

    RETURN
      sys.uid as uid,
      'system' as type,
      sys.name as name,
      CASE
        WHEN toLower(sys.systemCode) CONTAINS $search THEN 'systemCode'
        WHEN toLower(sys.name) CONTAINS $search THEN 'name'
        ELSE 'description'
      END as matchedField,
      coalesce(sys.systemCode, sys.name) as matchedValue,
      score as relevanceScore,
      {
        systemCode: sys.systemCode,
        systemType: st.name,
        locationPath: loc.code,
        facilityName: f.code,
        parentPath: [p IN parentPath | {uid: p.uid, name: p.name}],
        hasSubsystems: hasSubsystems,
        status: sys.status,
        systemLevel: sys.systemLevel
      } as context
    ORDER BY score DESC
    LIMIT $limit`
}

func getCatalogueItemSearchSubquery() string {
    return `
    CALL db.index.fulltext.queryNodes('searchIndexCatalogueItems', $fulltextSearch)
    YIELD node AS ci, score
    WHERE ci:CatalogueItem AND ci.deleted = false

    OPTIONAL MATCH (ci)-[:BELONGS_TO_CATEGORY]->(cat:CatalogueCategory)
    OPTIONAL MATCH (cat)<-[:HAS_SUBCATEGORY*0..10]-(rootCat:CatalogueCategory)
      WHERE NOT (rootCat)<-[:HAS_SUBCATEGORY]-()
    OPTIONAL MATCH (ci)-[:HAS_MANUFACTURER]->(man:Manufacturer)
    OPTIONAL MATCH (ci)-[:HAS_SUPPLIER]->(sup:Supplier)
    OPTIONAL MATCH (ci)<-[:IS_BASED_ON]-(items:Item)

    WITH ci, score, cat, man, sup,
         collect(DISTINCT rootCat.name) + [cat.name] as categoryHierarchy,
         count(DISTINCT items) as stockCount

    RETURN
      ci.uid as uid,
      'catalogueItem' as type,
      ci.name as name,
      CASE
        WHEN toLower(ci.catalogueNumber) CONTAINS $search THEN 'catalogueNumber'
        WHEN toLower(ci.name) CONTAINS $search THEN 'name'
        ELSE 'description'
      END as matchedField,
      ci.catalogueNumber as matchedValue,
      score as relevanceScore,
      {
        catalogueNumber: ci.catalogueNumber,
        categoryHierarchy: categoryHierarchy,
        manufacturerName: man.name,
        supplierName: sup.name,
        description: substring(ci.description, 0, 200),
        inStock: stockCount,
        imageUrl: ci.miniImageUrl
      } as context
    ORDER BY score DESC
    LIMIT $limit`
}

func getOrderSearchSubquery() string {
    return `
    CALL db.index.fulltext.queryNodes('searchIndexOrders', $fulltextSearch)
    YIELD node AS o, score
    WHERE o:Order AND o.deleted = false

    MATCH (o)-[:BELONGS_TO_FACILITY]->(f:Facility {code: $facilityCode})

    OPTIONAL MATCH (o)-[:HAS_ORDER_STATUS]->(status:OrderStatus)
    OPTIONAL MATCH (o)-[:HAS_SUPPLIER]->(sup:Supplier)
    OPTIONAL MATCH (o)-[:HAS_REQUESTOR]->(req:Employee)
    OPTIONAL MATCH (o)-[:HAS_ORDER_LINE]->(items:Item)

    WITH o, score, status, sup, req,
         count(DISTINCT items) as itemCount

    RETURN
      o.uid as uid,
      'order' as type,
      o.name as name,
      CASE
        WHEN toLower(o.orderNumber) CONTAINS $search THEN 'orderNumber'
        WHEN toLower(o.requestNumber) CONTAINS $search THEN 'requestNumber'
        ELSE 'name'
      END as matchedField,
      o.orderNumber as matchedValue,
      score as relevanceScore,
      {
        orderNumber: o.orderNumber,
        orderDate: toString(o.orderDate),
        statusName: status.name,
        statusCode: status.code,
        supplierName: sup.name,
        requestorName: req.firstName + ' ' + req.lastName,
        totalItemsCount: itemCount,
        deliveryStatus: o.deliveryStatus,
        requestNumber: o.requestNumber
      } as context
    ORDER BY score DESC
    LIMIT $limit`
}
```

### 4. Handlers

**File:** `services/search-service/search-handlers.go`

```go
package searchService

import (
    "net/http"
    "panda/apigateway/helpers"
    "strconv"
    "strings"

    "github.com/labstack/echo/v4"
    "github.com/rs/zerolog/log"
)

type SearchHandlers struct {
    searchService ISearchService
}

type ISearchHandlers interface {
    GlobalSearch() echo.HandlerFunc
}

func NewSearchHandlers(searchSvc ISearchService) ISearchHandlers {
    return &SearchHandlers{searchService: searchSvc}
}

// GlobalSearch godoc
// @Summary Global search across systems, catalogue items, and orders
// @Description Fast full-text search endpoint for command palette. Searches across System, CatalogueItem, and Order nodes with facility-based filtering.
// @Tags Search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search term (min 2 chars, max 100 chars)" minlength(2) maxlength(100)
// @Param types query string false "Comma-separated node types: system, catalogueItem, order (default: all)" example("system,catalogueItem")
// @Param limit query int false "Max results per type (default: 10, max: 50)" minimum(1) maximum(50) default(10)
// @Success 200 {object} models.SearchResponse "Search results"
// @Failure 400 {string} string "Invalid search term or parameters"
// @Failure 401 {string} string "Unauthorized - missing or invalid JWT token"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/search [get]
func (h *SearchHandlers) GlobalSearch() echo.HandlerFunc {
    return func(c echo.Context) error {
        // Get query parameters
        searchTerm := c.QueryParam("q")
        typesParam := c.QueryParam("types")
        limitParam := c.QueryParam("limit")

        // Validation
        if searchTerm == "" {
            return helpers.BadRequest("search term 'q' is required")
        }

        // Parse types
        var types []string
        if typesParam != "" {
            types = strings.Split(typesParam, ",")
            // Trim spaces
            for i, t := range types {
                types[i] = strings.TrimSpace(t)
            }
        }

        // Parse limit
        limit := 10
        if limitParam != "" {
            parsedLimit, err := strconv.Atoi(limitParam)
            if err != nil || parsedLimit < 1 {
                return helpers.BadRequest("limit must be a positive integer")
            }
            limit = parsedLimit
        }

        // Get facility code from context (set by auth middleware)
        facilityCode := c.Get("facilityCode").(string)

        // Execute search
        results, err := h.searchService.GlobalSearch(searchTerm, types, limit, facilityCode)

        if err != nil {
            log.Error().Err(err).Str("searchTerm", searchTerm).Msg("GlobalSearch failed")
            return helpers.BadRequest(err.Error())
        }

        return c.JSON(http.StatusOK, results)
    }
}
```

### 5. Routes

**File:** `services/search-service/search-routes.go`

```go
package searchService

import (
    "panda/apigateway/middlewares"
    "panda/apigateway/shared"

    "github.com/labstack/echo/v4"
)

func MapSearchRoutes(e *echo.Echo, h ISearchHandlers, jwtMiddleware echo.MiddlewareFunc) {
    // Global search endpoint
    e.GET("/v1/search",
        middlewares.Authorization(
            h.GlobalSearch(),
            shared.ROLE_SYSTEMS_VIEW, // Or create new ROLE_SEARCH
        ),
        jwtMiddleware,
    )
}
```

### 6. Service Registration

**Add to:** `services/init.go`

```go
// At the beginning of imports
searchService "panda/apigateway/services/search-service"

// In InitializeServicesAndMapRoutes() function add:

// search service
searchSvc := searchService.NewSearchService(neo4jDriver)
searchHandlers := searchService.NewSearchHandlers(searchSvc)
searchService.MapSearchRoutes(e, searchHandlers, jwtMiddleware)
log.Info().Msg("Search    service initialized successfully.")
```

---

## 📊 Performance Expectations

### Target Metrics

| Metric | Target | Acceptable | Critical |
|---------|--------|------------|----------|
| P50 (median) | < 50ms | < 100ms | < 200ms |
| P95 | < 100ms | < 200ms | < 500ms |
| P99 | < 200ms | < 500ms | < 1000ms |
| Error rate | < 0.1% | < 1% | < 5% |

### Performance Profiling

```go
// Logging in service
log.Info().
    Str("searchTerm", searchTerm).
    Int("totalResults", result.TotalCount).
    Int64("executionTimeMs", executionTime).
    Msg("GlobalSearch completed")
```

### Monitoring Queries

```cypher
// Check index usage
CALL db.indexes()
YIELD name, state, populationPercent, type
WHERE name STARTS WITH 'searchIndex'
RETURN name, state, populationPercent, type;

// Query profiling
PROFILE
CALL db.index.fulltext.queryNodes('searchIndexSystems', 'laser')
YIELD node, score
RETURN node.name, score
LIMIT 10;
```

---

## 🔗 Key File References

### Reference Implementation (For Understanding Patterns)

| File | Purpose | Important Parts |
|--------|------|----------------|
| `services/systems-service/systems-service.go` | Service layer pattern | Lines 282-293: GetSystemsWithSearchAndPagination |
| `services/systems-service/systems-handlers.go` | HTTP handlers pattern | Lines 307-335: GetSystemsWithSearchAndPagination handler |
| `services/systems-service/systems-routes.go` | Route mapping pattern | Line 13: Route with Authorization |
| `services/systems-service/systems-db-queries.go` | Query patterns | Lines 664-746: GetSystemsBySearchTextFullTextQuery |
| `services/systems-service/systems-db-queries.go` | Filter query pattern | Lines 332-382: GetSystemsSearchFilterQueryOnly |

### Integration

| File | Purpose |
|--------|------|
| `services/init.go` | Registration of all services |
| `server.go` | Main application entry point |
| `shared/roles.go` | Role constants for authorization |

### Helpers & Utilities

| File | Functionality |
|--------|---------------|
| `helpers/database.go` | Neo4j query helpers (`GetNeo4jArrayOfNodes`, etc.) |
| `helpers/interfaces.go` | Generic mapping utilities (`MapStruct`) |
| `helpers/echo.go` | HTTP error helpers (`BadRequest`) |

### Database

| File | Purpose |
|--------|------|
| `db/schema-simple.json` | Complete graph schema reference |
| `db/migrations.go` | Migration runner implementation |
| `db/neo4j/migrations/20230628124422_add_fulltext_systems.up.cypher` | Reference: System fulltext index |
| `db/neo4j/migrations/20230502135437_orders_search_indexes.up.cypher` | Reference: Order indexes |

### Existing Search Implementation (Inspiration)

| Endpoint | File | Query Pattern |
|----------|--------|---------------|
| Systems Search | `systems-service.go:282-293` | Fulltext + pagination |
| Catalogue Items | `catalogue-handlers.go:65-96` | Search with filters |
| Orders Search | `orders-handlers.go:56-86` | Search with sorting |

---

## 🚀 Implementation Steps

### Phase 1: Database Indexes ✅

1. **Create migrations:**
   ```bash
   # Up migration
   touch db/neo4j/migrations/20250112120000_add_global_search_indexes.up.cypher

   # Down migration
   touch db/neo4j/migrations/20250112120000_add_global_search_indexes.down.cypher
   ```

2. **Run migrations:**
   ```bash
   make run  # Migrations run automatically on startup
   ```

3. **Verify indexes:**
   ```cypher
   // In Neo4j Browser or cypher-shell
   SHOW INDEXES;

   // Expected to see:
   // - searchIndexSystems (FULLTEXT)
   // - searchIndexOrders (FULLTEXT)
   // - searchIndexCatalogueItems (FULLTEXT) ← new
   // + all property indexes
   ```

### Phase 2: Service Implementation 💻

1. **Create structure:**
   ```bash
   mkdir -p services/search-service/models
   touch services/search-service/search-service.go
   touch services/search-service/search-handlers.go
   touch services/search-service/search-routes.go
   touch services/search-service/search-db-queries.go
   touch services/search-service/models/search-models.go
   ```

2. **Implement in order:**
   - ✅ `models/search-models.go` - Data structures
   - ✅ `search-db-queries.go` - Neo4j queries
   - ✅ `search-service.go` - Business logic
   - ✅ `search-handlers.go` - HTTP handlers
   - ✅ `search-routes.go` - Route mapping

3. **Integrate into application:**
   - ✅ Modify `services/init.go` - add search service registration
   - ✅ (Optional) Add `ROLE_SEARCH` to `shared/roles.go`

### Phase 3: Testing 🧪

1. **Unit tests:**
   ```bash
   touch services/search-service/search-service_test.go
   go test ./services/search-service -v
   ```

2. **Integration tests:**
   - Test with real Neo4j database
   - Test all 3 node types
   - Test facility filtering
   - Test edge cases (empty results, special characters)

3. **Performance tests:**
   ```bash
   # Load testing with ab (Apache Bench) or k6
   ab -n 1000 -c 10 "http://localhost:50000/v1/search?q=laser"
   ```

### Phase 4: Documentation 📚

1. **Swagger:**
   ```bash
   make swagger
   ```

2. **Verify Swagger UI:**
   - Navigate to `http://localhost:50000/swagger/index.html`
   - Check `/v1/search` endpoint

3. **Client TypeScript types:**
   - Copy TypeScript interfaces from this document
   - Add to frontend project

### Phase 5: Monitoring & Optimization 📈

1. **Add metrics:**
   - Query execution time logging
   - Popular search terms tracking
   - Error rate monitoring

2. **Performance tuning:**
   - Analyze slow queries (>500ms)
   - Optimize OPTIONAL MATCH depth
   - Consider query caching for frequent queries

3. **Index maintenance:**
   ```cypher
   // Regularly check index health
   CALL db.indexes()
   YIELD name, state, populationPercent
   WHERE populationPercent < 100
   RETURN name, state, populationPercent;
   ```

---

## 📝 Checklist

### Database
- [ ] Up migration with indexes created
- [ ] Down migration created
- [ ] Migrations run successfully
- [ ] Indexes verified in Neo4j

### Code
- [ ] `search-service.go` implemented
- [ ] `search-handlers.go` implemented
- [ ] `search-routes.go` implemented
- [ ] `search-db-queries.go` implemented
- [ ] `models/search-models.go` implemented
- [ ] Integration in `services/init.go`
- [ ] Role constants (optional)

### Testing
- [ ] Unit tests written
- [ ] Integration tests written
- [ ] Performance tests performed
- [ ] Edge cases tested
- [ ] Error handling verified

### Documentation
- [ ] Swagger generated
- [ ] TypeScript types exported
- [ ] README updated (if exists)
- [ ] This design doc updated with lessons learned

### Production Readiness
- [ ] Logging implemented
- [ ] Monitoring setup
- [ ] Error tracking
- [ ] Performance baselines established
- [ ] Security review (facility filtering, injection prevention)

---

## 🎓 Lessons Learned & Notes

_After implementation, add notes about what worked well and what needed to be adjusted._

### Performance Notes
- TBD after implementation

### Gotchas
- TBD after implementation

### Future Improvements
- [ ] Implement caching for frequent queries
- [ ] Add highlights (matched text snippets)
- [ ] Implement fuzzy search
- [ ] Add search suggestions/autocomplete
- [ ] Analytics dashboard for popular searches

---

**Created:** 2025-01-12
**Last Update:** 2025-01-12
**Status:** 🟡 Design Phase - Ready for Implementation
**Owner:** Development Team
