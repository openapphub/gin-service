// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/cache/clear": {
            "post": {
                "description": "Clear all cached items with a specific prefix",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cache"
                ],
                "summary": "Clear cache by prefix",
                "parameters": [
                    {
                        "description": "Cache key prefix",
                        "name": "prefix",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Cache cleared successfully",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    }
                }
            }
        },
        "/cache/invalidate": {
            "post": {
                "description": "Remove a specific key from the cache",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cache"
                ],
                "summary": "Invalidate cache for a specific key",
                "parameters": [
                    {
                        "description": "Invalidate Cache Info",
                        "name": "invalidate_info",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.InvalidateCacheInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Cache invalidated successfully",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    }
                }
            }
        },
        "/cache/refresh": {
            "post": {
                "description": "Refresh the cache for a specific key with a new duration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cache"
                ],
                "summary": "Refresh cache for a specific key",
                "parameters": [
                    {
                        "description": "Refresh Cache Info",
                        "name": "refresh_info",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.RefreshCacheInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Cache refreshed successfully",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    }
                }
            }
        },
        "/ping": {
            "post": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ping"
                ],
                "summary": "Ping test",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "Authenticate a user and return a token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Log in a user",
                "parameters": [
                    {
                        "description": "User Login Info",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UserLoginService"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User logged in successfully",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    }
                }
            }
        },
        "/user/logout": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Log out the currently authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Log out a user",
                "responses": {
                    "200": {
                        "description": "User logged out successfully",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    }
                }
            }
        },
        "/user/me": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get information about the currently logged-in user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get current user information",
                "responses": {
                    "200": {
                        "description": "User information retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    }
                }
            }
        },
        "/user/refresh": {
            "post": {
                "description": "Refresh the JWT access token using a refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Refresh JWT token",
                "parameters": [
                    {
                        "description": "Refresh Token",
                        "name": "refresh_token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "New access token",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    }
                }
            }
        },
        "/user/register": {
            "post": {
                "description": "Register a new user with the provided information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User Registration Info",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UserRegisterService"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User registered successfully",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/serializer.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.InvalidateCacheInput": {
            "type": "object",
            "required": [
                "method",
                "path"
            ],
            "properties": {
                "body": {
                    "type": "string"
                },
                "method": {
                    "type": "string",
                    "enum": [
                        "GET",
                        "POST"
                    ]
                },
                "path": {
                    "type": "string"
                }
            }
        },
        "api.RefreshCacheInput": {
            "type": "object",
            "required": [
                "duration",
                "method",
                "path"
            ],
            "properties": {
                "body": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer",
                    "minimum": 1
                },
                "method": {
                    "type": "string",
                    "enum": [
                        "GET",
                        "POST"
                    ]
                },
                "path": {
                    "type": "string"
                }
            }
        },
        "serializer.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "error": {
                    "type": "string"
                },
                "msg": {
                    "type": "string"
                },
                "token": {}
            }
        },
        "service.UserLoginService": {
            "type": "object",
            "required": [
                "password",
                "user_name"
            ],
            "properties": {
                "device_info": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 8
                },
                "user_name": {
                    "type": "string",
                    "maxLength": 30,
                    "minLength": 5
                }
            }
        },
        "service.UserRegisterService": {
            "type": "object",
            "required": [
                "nickname",
                "password",
                "password_confirm",
                "user_name"
            ],
            "properties": {
                "nickname": {
                    "type": "string",
                    "maxLength": 30,
                    "minLength": 2
                },
                "password": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 8
                },
                "password_confirm": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 8
                },
                "user_name": {
                    "type": "string",
                    "maxLength": 30,
                    "minLength": 5
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "openapphub API",
	Description:      "This is a sample server for openapphub.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
