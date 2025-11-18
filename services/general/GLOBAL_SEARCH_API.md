# Global Search API Documentation

## Overview

The Global Search API allows you to search across multiple entity types (Systems, Orders, Catalogue Items) in the ELI-PANDA database using a single search endpoint.

## Endpoint

```
GET /v1/global-search
```

## Query Parameters

| Parameter  | Type   | Required | Description                                                     |
| ---------- | ------ | -------- | --------------------------------------------------------------- |
| searchText | string | Yes      | The text to search for across all entities                      |
| pagination | string | No       | JSON string with pagination info: `{"page": 1, "pageSize": 20}` |

## Response Format

```json
{
  "totalCount": 25,
  "data": [
    {
      "uid": "system-uuid-123",
      "name": "Main Laser System",
      "description": "Primary laser system for experiments",
      "nodeType": "System"
    },
    {
      "uid": "order-uuid-456",
      "name": "Equipment Purchase Order #2024-001",
      "description": "Order for new equipment",
      "nodeType": "Order"
    },
    {
      "uid": "catalogue-uuid-789",
      "name": "High-Power Laser Diode",
      "description": "Industrial grade laser diode component",
      "nodeType": "CatalogueItem"
    }
  ]
}
```

## Search Logic

The API searches in the following entity types and fields:

### Systems

- System name
- System code
- System description

### Orders

- Order name
- Order number
- Contract number
- Request number

### Catalogue Items

- Item name
- Catalogue number
- Item description

### Items (returns parent System/Order)

When Items match the search criteria, the API returns the related System or Order:

- Item name
- EUN (Equipment Unique Number)
- Serial number

## Features

1. **Full-text Search**: Case-insensitive search across multiple fields
2. **Cross-entity Results**: Single query returns results from all entity types
3. **Pagination**: Standard pagination support with configurable page size
4. **Relationship Traversal**: Items that match search criteria return their parent Systems/Orders

## Example Requests

### Basic Search

```
GET /v1/global-search?searchText=laser
```

### Search with Pagination

```
GET /v1/global-search?searchText=detector&pagination={"page":2,"pageSize":10}
```

### Search with Pagination (Custom Page Size)

```
GET /v1/global-search?searchText=optics&pagination={"page":1,"pageSize":25}
```

## Authentication

Requires valid JWT token with `ROLE_BASICS_VIEW` permission.

## Error Handling

### 400 Bad Request

```json
{
  "message": "Bad Request"
}
```

Returned when `searchText` parameter is missing or empty.

### 500 Internal Server Error

```json
{
  "message": "Internal Server Error"
}
```

Returned when database query fails or other server-side error occurs.

## Performance Notes

- The search is optimized for typical search patterns
- Results are ordered alphabetically by name
- Complex searches across large datasets may take longer
- Use pagination to manage large result sets efficiently
