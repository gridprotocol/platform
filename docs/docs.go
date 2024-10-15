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
        "/": {
            "get": {
                "description": "welcome api",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Welcome"
                ],
                "summary": "welcome",
                "responses": {
                    "200": {
                        "description": "file id",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/allowance/": {
            "get": {
                "description": "check the allowance between an owner and a spender",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Allowance"
                ],
                "summary": "Check Allowance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "owner",
                        "name": "owner",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "spender",
                        "name": "spender",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "404": {
                        "description": "page not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/getcp/": {
            "get": {
                "description": "get a provider's info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get Provider Info"
                ],
                "summary": "get a provider",
                "parameters": [
                    {
                        "type": "string",
                        "description": "address of a provider",
                        "name": "address",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.CPInfo"
                        }
                    },
                    "404": {
                        "description": "page not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/getorder/": {
            "get": {
                "description": "get an order info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get Order"
                ],
                "summary": "Get Order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user",
                        "name": "user",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "cp",
                        "name": "cp",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "404": {
                        "description": "page not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/getorders/": {
            "get": {
                "description": "get all orders of an user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get Orders"
                ],
                "summary": "Get Orders",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user",
                        "name": "user",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "404": {
                        "description": "page not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ginfo/": {
            "get": {
                "description": "get global info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get GInfo"
                ],
                "summary": "Get Global Info",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "404": {
                        "description": "page not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/listcp/": {
            "get": {
                "description": "list all providers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Listcps"
                ],
                "summary": "List all providers",
                "parameters": [
                    {
                        "type": "string",
                        "description": "start",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "number",
                        "name": "num",
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
                                "$ref": "#/definitions/routes.CPInfo"
                            }
                        }
                    },
                    "404": {
                        "description": "page not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/node/": {
            "get": {
                "description": "Get a node of a cp with node id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get Node"
                ],
                "summary": "Node",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cp address",
                        "name": "cp",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "node id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "404": {
                        "description": "page not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/nodes/": {
            "get": {
                "description": "Get all nodes of this provider",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get Nodes"
                ],
                "summary": "Nodes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cp address",
                        "name": "cp",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "404": {
                        "description": "page not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/querycredit": {
            "get": {
                "description": "Query credit of someone",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "QueryCredit"
                ],
                "summary": "QueryCredit",
                "parameters": [
                    {
                        "type": "string",
                        "description": "address of this caller",
                        "name": "address",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "query OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "get version",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Version"
                ],
                "summary": "version",
                "responses": {
                    "200": {
                        "description": "version OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "routes.CPInfo": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "endpoint": {
                    "type": "string"
                },
                "name": {
                    "description": "provider name",
                    "type": "string"
                },
                "numCPU": {
                    "type": "string"
                },
                "numGPU": {
                    "type": "string"
                },
                "numMem": {
                    "type": "string"
                },
                "numStore": {
                    "type": "string"
                },
                "priCPU": {
                    "type": "string"
                },
                "priGPU": {
                    "type": "string"
                },
                "priMem": {
                    "type": "string"
                },
                "priStore": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8002",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "PLATFORM API",
	Description:      "This is the grid platform",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
