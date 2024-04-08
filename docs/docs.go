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
        "/createorder": {
            "post": {
                "description": "create an order",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CreateOrder"
                ],
                "summary": "Create order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user address",
                        "name": "userAddress",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "cpAddress",
                        "name": "cpAddress",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "num cpu",
                        "name": "numCPU",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "price cpu",
                        "name": "priCPU",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "num",
                        "name": "numGPU",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "price",
                        "name": "priGPU",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "num",
                        "name": "numStore",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "price",
                        "name": "priStore",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "num",
                        "name": "numMem",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "price",
                        "name": "priMem",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "duration",
                        "name": "duration",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "regist OK",
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
                    "get cp"
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
        "/listorder": {
            "get": {
                "description": "list an order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ListOrder"
                ],
                "summary": "List order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user or provider",
                        "name": "role",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "address",
                        "name": "address",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "list OK",
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
        "/listpay": {
            "get": {
                "description": "ListPay",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ListPay"
                ],
                "summary": "ListPay",
                "parameters": [
                    {
                        "type": "string",
                        "description": "address of an user",
                        "name": "addr",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "list pay OK",
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
        "/listtransfer": {
            "get": {
                "description": "List all transfers of an address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "List transfers"
                ],
                "summary": "List all transfers",
                "parameters": [
                    {
                        "type": "string",
                        "description": "address to show list",
                        "name": "address",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "list transfer OK",
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
        "/pay": {
            "post": {
                "description": "Pay to credit with a transfer's key",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pay"
                ],
                "summary": "Pay for credit",
                "parameters": [
                    {
                        "type": "string",
                        "description": "transfer key",
                        "name": "transkey",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "pay OK",
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
                        "description": "role of this caller",
                        "name": "role",
                        "in": "query",
                        "required": true
                    },
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
        "/refreshtransfer": {
            "post": {
                "description": "Refresh status of transfer of an address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Refresh Transfer"
                ],
                "summary": "RefreshTransfer status of transfer",
                "parameters": [
                    {
                        "type": "string",
                        "description": "address to refresh",
                        "name": "address",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "refresh OK",
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
        "/registcp": {
            "post": {
                "description": "Regist CP",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "RegistCP"
                ],
                "summary": "Regist CP",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "address",
                        "name": "address",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "endpoint",
                        "name": "endpoint",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "num cpu",
                        "name": "numCPU",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "price cpu",
                        "name": "priCPU",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "num",
                        "name": "numGPU",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "price",
                        "name": "priGPU",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "num",
                        "name": "numStore",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "price",
                        "name": "priStore",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "num",
                        "name": "numMem",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "price",
                        "name": "priMem",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "regist OK",
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
        "/transfer": {
            "post": {
                "description": "user transfer token to platform",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transfer"
                ],
                "summary": "Transfer token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "tx hash",
                        "name": "txHash",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "from addr",
                        "name": "from",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "to addr",
                        "name": "to",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "transfer value",
                        "name": "value",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "transfer OK",
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
	Host:             "localhost:8081",
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
