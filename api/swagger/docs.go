// Code generated by swaggo/swag. DO NOT EDIT.

package swagger

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
        "/api/v1/chats": {
            "get": {
                "description": "Prints all chats id and name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chats"
                ],
                "summary": "Get all chats",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/v1.Chat"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates new chat, runs it in background and returns chat ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chats"
                ],
                "summary": "Create chat",
                "parameters": [
                    {
                        "description": "Chat name",
                        "name": "chat",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.CreateChatRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.CreateChatResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/v1.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/api/v1/chats/active": {
            "get": {
                "description": "Prints all active chats id and name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chats"
                ],
                "summary": "Get all active chats",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/v1.Chat"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/api/v1/chats/{id}/connect": {
            "get": {
                "description": "Upgrades http connection to websocket",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chats"
                ],
                "summary": "Upgrade http connection",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Chat ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/api/v1/chats/{id}/messages": {
            "get": {
                "description": "Prints 100 messages from chat",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chats"
                ],
                "summary": "Get chat message history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Chat ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/v1.Message"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/v1.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/v1.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/api/v1/users": {
            "post": {
                "description": "Creates nes User, return his ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Create New User",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.CreateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/v1.CreateUserResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/v1.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/api/v1/users/login": {
            "post": {
                "description": "Login User, returns token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Login User",
                "parameters": [
                    {
                        "description": "User login data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.LoginUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.LoginUserResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/v1.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.Chat": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "UUID"
                },
                "name": {
                    "type": "string",
                    "example": "Best chat name!"
                }
            }
        },
        "v1.CreateChatRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "Best chat name!"
                }
            }
        },
        "v1.CreateChatResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "UUID"
                }
            }
        },
        "v1.CreateUserRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "secret1234!"
                },
                "username": {
                    "type": "string",
                    "example": "Username"
                }
            }
        },
        "v1.CreateUserResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "UUID"
                }
            }
        },
        "v1.HTTPError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "error description"
                }
            }
        },
        "v1.LoginUserRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "secret1234!"
                }
            }
        },
        "v1.LoginUserResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "JWT access token"
                },
                "refresh_token": {
                    "type": "string",
                    "example": "JWT refresh token"
                }
            }
        },
        "v1.Message": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "string",
                    "example": "UUID"
                },
                "author_name": {
                    "type": "string",
                    "example": "Username"
                },
                "chat_id": {
                    "type": "string",
                    "example": "UUID"
                },
                "content": {
                    "type": "string",
                    "example": "Hello world!"
                },
                "created_at": {
                    "type": "string",
                    "example": "2022-05-01T00:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "UUID"
                },
                "is_edited": {
                    "type": "boolean",
                    "example": false
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
