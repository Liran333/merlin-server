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
                            "$ref": "#/definitions/controller.reqToCreateModel"
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
        "/v1/session": {
            "put": {
                "description": "logout",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Session"
                ],
                "summary": "Logout",
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/controller.logoutInfo"
                        }
                    }
                }
            },
            "post": {
                "description": "login",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Session"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "body of login",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.reqToLogin"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/github_com_openmerlin_merlin-server_session_app.UserDTO"
                        }
                    }
                }
            }
        },
        "/v1/session/check": {
            "put": {
                "description": "check and refresh session",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Session"
                ],
                "summary": "CheckAndRefresh",
                "responses": {
                    "202": {
                        "description": "Accepted"
                    }
                }
            }
        },
        "/v1/session/clear": {
            "delete": {
                "description": "Clear session when it expired",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Session"
                ],
                "summary": "Clear",
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/v1/space/{id}": {
            "get": {
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
                            "$ref": "#/definitions/app.SpaceMetaDTO"
                        }
                    }
                }
            }
        },
        "/v1/user": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "get current sign-in user info",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get current user info",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "update user basic info",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update",
                "parameters": [
                    {
                        "description": "body of updating user",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.userBasicInfoUpdateRequest"
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
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "delete",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Request Delete User info",
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/v1/user/email/bind": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "bind user's email",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Bind User Email",
                "parameters": [
                    {
                        "description": "body of bind email info",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.bindEmailRequest"
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
        "/v1/user/email/send": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "send user's email verify code",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Send User Email Verify code",
                "parameters": [
                    {
                        "description": "body of bind email info",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.sendEmailRequest"
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
        "/v1/user/privacy": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "revoke",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "PrivacyRevoke",
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
        "/v1/user/token": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "list all platform tokens of user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "GetTokenInfo",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ResponseData"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "create a new platform token of user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "CreatePlatformToken",
                "parameters": [
                    {
                        "description": "body of create token",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.tokenCreateRequest"
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
                    "User"
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
        "/v1/user/token/{name}": {
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "delete a new platform token of user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "DeletePlatformToken",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
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
                    "User"
                ],
                "summary": "Get platform user info",
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
        "app.SpaceMetaDTO": {
            "type": "object",
            "properties": {
                "hardware": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "sdk": {
                    "type": "string"
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
        "controller.bindEmailRequest": {
            "type": "object",
            "required": [
                "code",
                "email"
            ],
            "properties": {
                "code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "controller.logoutInfo": {
            "type": "object",
            "properties": {
                "id_token": {
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
        "controller.reqToCreateModel": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "init_readme": {
                    "type": "boolean"
                },
                "license": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "visibility": {
                    "type": "string"
                }
            }
        },
        "controller.reqToLogin": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "redirect_uri": {
                    "type": "string"
                }
            }
        },
        "controller.sendEmailRequest": {
            "type": "object",
            "required": [
                "capt",
                "email"
            ],
            "properties": {
                "capt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "controller.tokenCreateRequest": {
            "type": "object",
            "required": [
                "name",
                "perm"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "perm": {
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
        "controller.userBasicInfoUpdateRequest": {
            "type": "object",
            "properties": {
                "avatar_id": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "revoke_delete": {
                    "type": "boolean"
                }
            }
        },
        "github_com_openmerlin_merlin-server_session_app.UserDTO": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                },
                "allow_request": {
                    "type": "boolean"
                },
                "avatar_id": {
                    "type": "string"
                },
                "created_at": {
                    "type": "integer"
                },
                "default_role": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "request_delete": {
                    "type": "boolean"
                },
                "request_delete_at": {
                    "type": "integer"
                },
                "type": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "integer"
                },
                "website": {
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
