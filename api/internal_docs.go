// Package api Code generated by swaggo/swag. DO NOT EDIT
package api

import "github.com/swaggo/swag"

const docTemplateinternal = `{
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
        "/v1/activity": {
            "post": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "add activities to DB",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "ActivityInternal"
                ],
                "summary": "AddActivity",
                "parameters": [
                    {
                        "description": "body of create activity app",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/activityapp.ReqToCreateActivity"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controller.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "code": {
                                            "type": "string"
                                        },
                                        "data": {
                                            "type": "object"
                                        },
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controller.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "code": {
                                            "type": "string"
                                        },
                                        "data": {},
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "delete all the record of an resource in the DB",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "ActivityInternal"
                ],
                "summary": "DeleteActivity",
                "parameters": [
                    {
                        "description": "body of delete activity app",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/activityapp.ReqToDeleteActivity"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controller.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "code": {
                                            "type": "string"
                                        },
                                        "data": {
                                            "type": "object"
                                        },
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controller.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "code": {
                                            "type": "string"
                                        },
                                        "data": {},
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/coderepo/permission/read": {
            "post": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "check if can read repo's sub-resource not the repo itsself",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Permission"
                ],
                "summary": "Read",
                "parameters": [
                    {
                        "description": "body of request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToCheckPermission"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/coderepo/permission/update": {
            "post": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "check if can create/update/delete repo's sub-resource not the repo itsself",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Permission"
                ],
                "summary": "Update",
                "parameters": [
                    {
                        "description": "body of request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToCheckPermission"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/computility/account": {
            "post": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "user joined computility org",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "ComputilityInternal"
                ],
                "summary": "ComputilityUserJoin",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToUserOrgOperate"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/computility/account/remove": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "user removed from computility org",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "ComputilityInternal"
                ],
                "summary": "ComputilityUserRemove",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToUserOrgOperate"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/computility/org/delete": {
            "delete": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "delete computility org",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "ComputilityInternal"
                ],
                "summary": "ComputilityOrgDelete",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToOrgDelete"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/model/{id}": {
            "get": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "get model info by id",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "ModelInternal"
                ],
                "summary": "GetById",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of model",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/model/{id}/label": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "reset label of model",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "ModelInternal"
                ],
                "summary": "ResetLabel",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of model",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToResetLabel"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/session/check": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "check and refresh session",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SessionInternal"
                ],
                "summary": "CheckAndRefresh",
                "parameters": [
                    {
                        "description": "body of new member",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/session.RequestToCheckAndRefresh"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controller.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "code": {
                                            "type": "string"
                                        },
                                        "data": {
                                            "$ref": "#/definitions/session.ResponseToCheckAndRefresh"
                                        },
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/session/clear": {
            "delete": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "Clear session when it expired",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SessionInternal"
                ],
                "summary": "Clear session by session id",
                "parameters": [
                    {
                        "description": "body of new member",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/session.RequestToClear"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space-app/": {
            "post": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "create space app",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceApp"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "description": "body of creating space app",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToCreateSpaceApp"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space-app/build/done": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "notify space app build is done",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceApp"
                ],
                "summary": "NotifyBuildIsDone",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToSetBuildIsDone"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space-app/build/started": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "notify space app building is started",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceApp"
                ],
                "summary": "NotifyBuildIsStarted",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToUpdateBuildInfo"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space-app/service/started": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "notify space app service is started",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceApp"
                ],
                "summary": "NotifyServiceIsStarted",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToUpdateServiceInfo"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space-app/status": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "notify space app status",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceApp"
                ],
                "summary": "NotifyUpdateStatus",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToSetStatus"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space/{id}": {
            "get": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "get space",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceInternal"
                ],
                "summary": "Get",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of space",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space/{id}/local_cmd": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "update space local cmd",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceInternal"
                ],
                "summary": "UpdateSpaceLocalCmd",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of space",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "local cmd to reproduce the space",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space/{id}/local_env_info": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "update space local env info",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceInternal"
                ],
                "summary": "UpdateSpaceLocalEnvInfo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of space",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "local env info to update local space env info",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/space/{id}/model": {
            "put": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "update space models relations",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SpaceInternal"
                ],
                "summary": "UpdateSpaceModels",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of space",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.ModeIds"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        },
        "/v1/user/token/verify": {
            "post": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "verify a platform token of user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "UserInternal"
                ],
                "summary": "Verify token",
                "parameters": [
                    {
                        "description": "body of token",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.tokenVerifyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "token"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "token"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "token"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "internal"
                        }
                    }
                }
            }
        },
        "/v1/user/{name}/platform": {
            "get": {
                "security": [
                    {
                        "Internal": []
                    }
                ],
                "description": "Get platform user info",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "UserInternal"
                ],
                "summary": "GetPlatformUser info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of the user",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "activityapp.ReqToCreateActivity": {
            "type": "object",
            "properties": {
                "owner": {
                    "type": "string"
                },
                "resource_index": {
                    "type": "string"
                },
                "resource_type": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "activityapp.ReqToDeleteActivity": {
            "type": "object",
            "properties": {
                "resource_index": {
                    "type": "string"
                },
                "resource_type": {
                    "type": "string"
                }
            }
        },
        "controller.ModeIds": {
            "type": "object",
            "properties": {
                "ids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "controller.ResponseData": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {},
                "msg": {
                    "type": "string"
                }
            }
        },
        "controller.reqToCheckPermission": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "controller.reqToCreateSpaceApp": {
            "type": "object",
            "properties": {
                "commit_id": {
                    "type": "string"
                },
                "space_id": {
                    "type": "string"
                }
            }
        },
        "controller.reqToOrgDelete": {
            "type": "object",
            "properties": {
                "org_name": {
                    "type": "string"
                }
            }
        },
        "controller.reqToResetLabel": {
            "type": "object",
            "properties": {
                "frameworks": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "license": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "task": {
                    "type": "string"
                }
            }
        },
        "controller.reqToSetBuildIsDone": {
            "type": "object",
            "properties": {
                "commit_id": {
                    "type": "string"
                },
                "logs": {
                    "type": "string"
                },
                "space_id": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "controller.reqToSetStatus": {
            "type": "object",
            "properties": {
                "commit_id": {
                    "type": "string"
                },
                "space_id": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "controller.reqToUpdateBuildInfo": {
            "type": "object",
            "properties": {
                "commit_id": {
                    "type": "string"
                },
                "log_url": {
                    "type": "string"
                },
                "space_id": {
                    "type": "string"
                }
            }
        },
        "controller.reqToUpdateServiceInfo": {
            "type": "object",
            "properties": {
                "app_url": {
                    "type": "string"
                },
                "commit_id": {
                    "type": "string"
                },
                "log_url": {
                    "type": "string"
                },
                "space_id": {
                    "type": "string"
                }
            }
        },
        "controller.reqToUserOrgOperate": {
            "type": "object",
            "properties": {
                "org_name": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                }
            }
        },
        "controller.tokenVerifyRequest": {
            "type": "object",
            "required": [
                "action",
                "token"
            ],
            "properties": {
                "action": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "session.RequestToCheckAndRefresh": {
            "type": "object",
            "properties": {
                "csrf_token": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "session_id": {
                    "type": "string"
                },
                "user_agent": {
                    "type": "string"
                }
            }
        },
        "session.RequestToClear": {
            "type": "object",
            "properties": {
                "session_id": {
                    "type": "string"
                }
            }
        },
        "session.ResponseToCheckAndRefresh": {
            "type": "object",
            "properties": {
                "csrf_token": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and api Bearer.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "Internal": {
            "description": "Type \"Internal\" followed by a space and internal token.",
            "type": "apiKey",
            "name": "TOKEN",
            "in": "header"
        }
    }
}`

// SwaggerInfointernal holds exported Swagger Info so clients can modify it
var SwaggerInfointernal = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "internal",
	SwaggerTemplate:  docTemplateinternal,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfointernal.InstanceName(), SwaggerInfointernal)
}
