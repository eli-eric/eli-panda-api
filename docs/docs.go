// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Jiří Švácha",
            "email": "jiri.svacha@eli-beams.eu"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/general/{uid}/graph": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get graph by uid",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "general"
                ],
                "summary": "Get graph by uid",
                "parameters": [
                    {
                        "type": "string",
                        "description": "uid",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GraphResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/getuserbyazureidtoken": {
            "get": {
                "description": "Get user by azure id token",
                "tags": [
                    "Security"
                ],
                "summary": "Get user by azure id token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tenant ID",
                        "name": "tenantId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Azure ID Token",
                        "name": "azureIdToken",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserAuthInfo"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/v1/order": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update an order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Update an order",
                "parameters": [
                    {
                        "description": "Order object that needs to be updated",
                        "name": "order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.OrderDetail"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.OrderDetail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Insert a new order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Insert a new order",
                "parameters": [
                    {
                        "description": "Order object that needs to be inserted",
                        "name": "order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.OrderDetail"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.OrderDetail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/order/{uid}": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete an order by order UID",
                "tags": [
                    "Orders"
                ],
                "summary": "Delete an order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order UID",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/orders/eun-for-print": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get items for EUN print",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Get items for EUN print",
                "parameters": [
                    {
                        "type": "string",
                        "description": "EUNs",
                        "name": "euns",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ItemForEunPrint"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/physical-item/move": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Move physical item from one system to another",
                "tags": [
                    "Systems"
                ],
                "summary": "Move physical item",
                "parameters": [
                    {
                        "description": "Move physical item request model",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PhysicalItemMovement"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Return destination system UID"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/physical-item/replace": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Replace physical items between two systems",
                "tags": [
                    "Systems"
                ],
                "summary": "Replace physical item",
                "parameters": [
                    {
                        "description": "Move physical item request model",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PhysicalItemMovement"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Return destination system UID"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/system/system-code": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create new system with new unique system code based on system type and zone",
                "tags": [
                    "Systems"
                ],
                "summary": "Create new system with code",
                "parameters": [
                    {
                        "description": "System code request model",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SystemCodeRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.System"
                        }
                    },
                    "400": {
                        "description": "Bad request - missing required fields"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/systems/locations-flat": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get all locations flat list",
                "tags": [
                    "Systems"
                ],
                "summary": "Get all locations flat list",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Codebook"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/systems/move": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Move systems to another parent system",
                "tags": [
                    "Systems"
                ],
                "summary": "Move systems",
                "parameters": [
                    {
                        "description": "Move systems request model",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SystemsMovement"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Return destination system UID"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/systems/recalculate-spare-parts": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Recalculate spare parts for all systems",
                "tags": [
                    "Systems"
                ],
                "summary": "Recalculate spare parts",
                "responses": {
                    "204": {
                        "description": "No content"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/systems/reload": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get systems tree by UIDs",
                "tags": [
                    "Systems"
                ],
                "summary": "Get systems tree by UIDs",
                "parameters": [
                    {
                        "description": "Array of system tree UIDs",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.SystemTreeUid"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.System"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/systems/sync-locations-by-eun": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Sync system locations by EUNs",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Systems"
                ],
                "summary": "Sync system locations by EUNs",
                "parameters": [
                    {
                        "description": "EUN with location UID",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.EunLocation"
                            }
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No content"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/systems/system-types": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get all system types",
                "tags": [
                    "Systems"
                ],
                "summary": "Get all system types",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Codebook"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/v1/systems/zones": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get all zones",
                "tags": [
                    "Systems"
                ],
                "summary": "Get all zones",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Codebook"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.CatalogueCategoryProperty": {
            "type": "object",
            "properties": {
                "defaultValue": {
                    "type": "string"
                },
                "listOfValues": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "type": {
                    "$ref": "#/definitions/models.CatalogueCategoryPropertyType"
                },
                "uid": {
                    "type": "string"
                },
                "unit": {
                    "$ref": "#/definitions/models.Codebook"
                }
            }
        },
        "models.CatalogueCategoryPropertyType": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.CatalogueItem": {
            "type": "object",
            "properties": {
                "catalogueNumber": {
                    "type": "string"
                },
                "category": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "categoryUID": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "details": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.CatalogueItemDetail"
                    }
                },
                "lastUpdateTime": {
                    "type": "string"
                },
                "manufacturerNumber": {
                    "type": "string"
                },
                "manufacturerUrl": {
                    "type": "string"
                },
                "miniImageUrl": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "supplier": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.CatalogueItemDetail": {
            "type": "object",
            "properties": {
                "property": {
                    "$ref": "#/definitions/models.CatalogueCategoryProperty"
                },
                "propertyGroup": {
                    "type": "string"
                },
                "value": {}
            }
        },
        "models.Codebook": {
            "type": "object",
            "properties": {
                "additionalData": {
                    "type": "string"
                },
                "code": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.EunLocation": {
            "type": "object",
            "properties": {
                "eun": {
                    "type": "string"
                },
                "location_uid": {
                    "type": "string"
                }
            }
        },
        "models.GraphLink": {
            "type": "object",
            "properties": {
                "relationship": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "target": {
                    "type": "string"
                }
            }
        },
        "models.GraphNode": {
            "type": "object",
            "properties": {
                "label": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "properties": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.GraphResponse": {
            "type": "object",
            "properties": {
                "links": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.GraphLink"
                    }
                },
                "nodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.GraphNode"
                    }
                }
            }
        },
        "models.ItemForEunPrint": {
            "type": "object",
            "properties": {
                "catalogueNumber": {
                    "type": "string"
                },
                "eun": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "manufacturer": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                },
                "serialNumber": {
                    "type": "string"
                }
            }
        },
        "models.OrderDetail": {
            "type": "object",
            "properties": {
                "contractNumber": {
                    "type": "string"
                },
                "lastUpdateTime": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "notes": {
                    "type": "string"
                },
                "orderDate": {
                    "type": "string"
                },
                "orderLines": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.OrderLine"
                    }
                },
                "orderNumber": {
                    "type": "string"
                },
                "orderStatus": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "procurementResponsible": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "requestNumber": {
                    "type": "string"
                },
                "requestor": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "supplier": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.OrderLine": {
            "type": "object",
            "properties": {
                "catalogueNumber": {
                    "type": "string"
                },
                "catalogueUid": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "deliveredTime": {
                    "type": "string"
                },
                "eun": {
                    "type": "string"
                },
                "isDelivered": {
                    "type": "boolean"
                },
                "itemUsage": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "location": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "name": {
                    "type": "string"
                },
                "notes": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "serialNumber": {
                    "type": "string"
                },
                "system": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.PhysicalItem": {
            "type": "object",
            "properties": {
                "catalogueItem": {
                    "$ref": "#/definitions/models.CatalogueItem"
                },
                "currency": {
                    "type": "string"
                },
                "eun": {
                    "type": "string"
                },
                "itemUsage": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "orderNumber": {
                    "type": "string"
                },
                "orderUid": {
                    "type": "string"
                },
                "price": {},
                "serialNumber": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.PhysicalItemMovement": {
            "type": "object",
            "properties": {
                "condition": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "deleteSourceSystem": {
                    "type": "boolean"
                },
                "destinationSystemUid": {
                    "type": "string"
                },
                "itemUsage": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "location": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "parentSystemUid": {
                    "type": "string"
                },
                "sourceSystemUid": {
                    "type": "string"
                },
                "systemName": {
                    "type": "string"
                }
            }
        },
        "models.System": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "hasSubsystems": {
                    "type": "boolean"
                },
                "history": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.SystemHistory"
                    }
                },
                "importance": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "location": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "miniImageUrl": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "parentPath": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.SystemSimpleResponse"
                    }
                },
                "parentUid": {
                    "type": "string"
                },
                "physicalItem": {
                    "$ref": "#/definitions/models.PhysicalItem"
                },
                "responsible": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "sparesIn": {
                    "type": "integer"
                },
                "sparesOut": {
                    "type": "integer"
                },
                "statistics": {
                    "$ref": "#/definitions/models.SystemStatistics"
                },
                "subSystems": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.System"
                    }
                },
                "systemAlias": {
                    "type": "string"
                },
                "systemCode": {
                    "type": "string"
                },
                "systemLevel": {
                    "type": "string"
                },
                "systemType": {
                    "$ref": "#/definitions/models.Codebook"
                },
                "uid": {
                    "type": "string"
                },
                "zone": {
                    "$ref": "#/definitions/models.Codebook"
                }
            }
        },
        "models.SystemCodeRequest": {
            "type": "object",
            "properties": {
                "parentUid": {
                    "type": "string"
                },
                "systemTypeUid": {
                    "type": "string"
                },
                "zoneUid": {
                    "type": "string"
                }
            }
        },
        "models.SystemHistory": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "changedAt": {
                    "type": "string"
                },
                "changedBy": {
                    "type": "string"
                },
                "detail": {
                    "$ref": "#/definitions/models.SystemHistoryDetail"
                },
                "historyType": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.SystemHistoryDetail": {
            "type": "object",
            "properties": {
                "direction": {
                    "type": "string"
                },
                "systemName": {
                    "type": "string"
                },
                "systemUid": {
                    "type": "string"
                }
            }
        },
        "models.SystemSimpleResponse": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.SystemStatistics": {
            "type": "object",
            "properties": {
                "minimalSpareParstCount": {
                    "type": "number"
                },
                "sp_coverage": {
                    "type": "number"
                },
                "sparePartsCount": {
                    "type": "integer"
                },
                "sparePartsCoverageSum": {
                    "type": "number"
                },
                "subsystemsCount": {
                    "type": "integer"
                }
            }
        },
        "models.SystemTreeUid": {
            "type": "object",
            "properties": {
                "children": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.SystemTreeUid"
                    }
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "models.SystemsMovement": {
            "type": "object",
            "properties": {
                "systemsToMoveUids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "targetParentSystemUid": {
                    "type": "string"
                }
            }
        },
        "models.UserAuthInfo": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "facility": {
                    "type": "string"
                },
                "facilityCode": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "isEnabled": {
                    "type": "boolean"
                },
                "lastName": {
                    "type": "string"
                },
                "passwordHash": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "uid": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "JWT token. \u003cbr\u003e How to obtain: https://eli-eric.atlassian.net/wiki/spaces/CS/pages/948797504/How+to+get+PANDA+API+Token \u003cbr\u003e Add word Bearer before the token here.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:50000",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "PANDA REST API - localhost",
	Description:      "This is the REST API to the PANDA database. \\n This is the only place to access data from the PANDA database.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
