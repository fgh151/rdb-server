// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/auth": {
            "post": {
                "description": "Authenticate in admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email for login",
                        "name": "email",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Password for login",
                        "name": "password",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                }
            }
        },
        "/admin/cf": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "List cloud functions",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cloud functions",
                    "Admin"
                ],
                "summary": "List cloud functions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.CloudFunction"
                            }
                        }
                    }
                }
            }
        },
        "/admin/config": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "List configs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config manager",
                    "Admin"
                ],
                "summary": "List configs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Config"
                            }
                        }
                    }
                }
            }
        },
        "/admin/cron": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "List cron jobs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cron",
                    "Admin"
                ],
                "summary": "List cron jobs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.CronJob"
                            }
                        }
                    }
                }
            }
        },
        "/admin/ds": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "List data sources",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Data source",
                    "Admin"
                ],
                "summary": "List data sources",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.DataSource"
                            }
                        }
                    }
                }
            }
        },
        "/admin/push": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "List push messages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Push messages",
                    "Admin"
                ],
                "summary": "List push messages",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.PushMessage"
                            }
                        }
                    }
                }
            }
        },
        "/admin/topics": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "List topics",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Topic"
                ],
                "summary": "List topics",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Project"
                            }
                        }
                    }
                }
            }
        },
        "/admin/users": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "List users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "List users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    }
                }
            }
        },
        "/config/{id}": {
            "get": {
                "description": "Get config by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config manager"
                ],
                "summary": "Config",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Config id",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    }
                }
            }
        },
        "/dse/{id}": {
            "get": {
                "description": "Get data source by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Data source"
                ],
                "summary": "Get item",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Auth key",
                        "name": "db-key",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Project"
                        }
                    }
                }
            }
        },
        "/em/find/{topic}": {
            "get": {
                "description": "Search in topic",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity manager"
                ],
                "summary": "Search",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Topic name",
                        "name": "topic",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    }
                }
            }
        },
        "/em/list/{topic}": {
            "get": {
                "description": "List topic records",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity manager"
                ],
                "summary": "List",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Topic name",
                        "name": "topic",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    }
                }
            }
        },
        "/em/subscribe/{topic}/{key}": {
            "get": {
                "description": "Socket subscribe to topic",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity manager"
                ],
                "summary": "Subscribe",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Topic name",
                        "name": "topic",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Db key",
                        "name": "key",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    }
                }
            }
        },
        "/em/{topic}": {
            "post": {
                "description": "Create topic record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity manager"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Topic name",
                        "name": "topic",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    }
                }
            }
        },
        "/em/{topic}/{id}": {
            "delete": {
                "description": "Delete entity record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity manager"
                ],
                "summary": "Delete",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Topic name",
                        "name": "topic",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Topic record id",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    }
                }
            },
            "patch": {
                "description": "Update entity record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity manager"
                ],
                "summary": "Update",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Topic name",
                        "name": "topic",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Topic record id",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.CloudFunction": {
            "type": "object",
            "properties": {
                "container": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "params": {
                    "type": "string"
                },
                "project": {
                    "$ref": "#/definitions/models.Project"
                },
                "project_id": {
                    "type": "string"
                },
                "run_count": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.Config": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "project": {
                    "$ref": "#/definitions/models.Project"
                },
                "project_id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.CronJob": {
            "type": "object",
            "properties": {
                "function": {
                    "$ref": "#/definitions/models.CloudFunction"
                },
                "function_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "time_params": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.DataSource": {
            "type": "object",
            "properties": {
                "cache": {
                    "type": "boolean"
                },
                "dsn": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "project": {
                    "$ref": "#/definitions/models.Project"
                },
                "project_id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models.Project": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "key": {
                    "type": "string"
                },
                "origins": {
                    "type": "string"
                },
                "topic": {
                    "type": "string"
                }
            }
        },
        "models.PushMessage": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "payload": {
                    "type": "string"
                },
                "receivers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.UserDevice"
                    }
                },
                "sent": {
                    "type": "boolean"
                },
                "sent_at": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "topic": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "active": {
                    "type": "boolean"
                },
                "admin": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string"
                },
                "devices": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.UserDevice"
                    }
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_login": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "models.UserDevice": {
            "type": "object",
            "properties": {
                "device": {
                    "type": "string"
                },
                "device_token": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
