package catalog

//Schema for the
const WoTSchema = `
{
    "title": "WoT TD Schema - 16 October 2019",
    "description": "JSON Schema for validating TD instances against the TD model. TD instances can be with or without terms that have default values",
    "$schema ": "http://json-schema.org/draft-07/schema#",
    "definitions": {
        "thing-context-w3c-uri": {
            "type": "string",
            "enum": [
                "https://www.w3.org/2019/wot/td/v1"
            ]
        },
        "thing-context": {
            "oneOf": [{
                    "type": "array",
                    "items": {
                        "anyOf": [{
                                "$ref": "#/definitions/anyUri"
                            },
                            {
                                "type": "object"
                            }
                        ]
                    },
                    "contains": {
                        "$ref": "#/definitions/thing-context-w3c-uri"
                    }
                },
                {
                    "$ref": "#/definitions/thing-context-w3c-uri"
                }
            ]
        },
        "type_declaration": {
            "oneOf": [{
                    "type": "string"
                },
                {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            ]
        },
        "property_element": {
            "type": "object",
            "properties": {
                "@type": {
                    "$ref": "#/definitions/type_declaration"
                },
                "description": {
                    "$ref": "#/definitions/description"
                },
                "descriptions": {
                    "$ref": "#/definitions/descriptions"
                },
                "title": {
                    "$ref": "#/definitions/title"
                },
                "titles": {
                    "$ref": "#/definitions/titles"
                },
                "uriVariables": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/dataSchema"
                    }
                },
                "forms": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/form_element_property"
                    }
                },
                "observable": {
                    "type": "boolean"
                },
                "writeOnly": {
                    "type": "boolean"
                },
                "readOnly": {
                    "type": "boolean"
                },
                "oneOf": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dataSchema"
                    }
                },
                "unit": {
                    "type": "string"
                },
                "enum": {
                    "type": "array",
					"items":{
						"type": "object"
					},
                    "minItems": 1,
                    "uniqueItems": true
                },
                "format": {
                    "type": "string"
                },
                "const": {},
                "type": {
                    "type": "string",
                    "enum": [
                        "boolean",
                        "integer",
                        "number",
                        "string",
                        "object",
                        "array",
                        "null"
                    ]
                },
                "items": {
                    "oneOf": [{
                            "$ref": "#/definitions/dataSchema"
                        },
                        {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dataSchema"
                            }
                        }
                    ]
                },
                "maxItems": {
                    "type": "integer",
                    "minimum": 0
                },
                "minItems": {
                    "type": "integer",
                    "minimum": 0
                },
                "minimum": {
                    "type": "number"
                },
                "maximum": {
                    "type": "number"
                },
                "properties": {
                    "additionalProperties": {
                        "$ref": "#/definitions/dataSchema"
                    }
                },
                "required": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            },
            "required": [
                "forms"
            ],
            "additionalProperties": true
        },
        "action_element": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "descriptions": {
                    "$ref": "#/definitions/descriptions"
                },
                "title": {
                    "$ref": "#/definitions/title"
                },
                "titles": {
                    "$ref": "#/definitions/titles"
                },
                "uriVariables": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/dataSchema"
                    }
                },
                "@type": {
                    "$ref": "#/definitions/type_declaration"
                },
                "forms": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/form_element_action"
                    }
                },
                "input": {
                    "$ref": "#/definitions/dataSchema"
                },
                "output": {
                    "$ref": "#/definitions/dataSchema"
                },
                "safe": {
                    "type": "boolean"
                },
                "idempotent": {
                    "type": "boolean"
                }
            },
            "required": [
                "forms"
            ],
            "additionalProperties": true
        },
        "event_element": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "descriptions": {
                    "$ref": "#/definitions/descriptions"
                },
                "title": {
                    "$ref": "#/definitions/title"
                },
                "titles": {
                    "$ref": "#/definitions/titles"
                },
                "uriVariables": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/dataSchema"
                    }
                },
                "@type": {
                    "$ref": "#/definitions/type_declaration"
                },
                "forms": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/form_element_event"
                    }
                },
                "subscription": {
                    "$ref": "#/definitions/dataSchema"
                },
                "data": {
                    "$ref": "#/definitions/dataSchema"
                },
                "cancellation": {
                    "$ref": "#/definitions/dataSchema"
                },
                "type": {
                    "not": {}
                },
                "enum": {
                    "not": {}
                },
                "const": {
                    "not": {}
                }
            },
            "required": [
                "forms"
            ],
            "additionalProperties": true
        },
        "form_element_property": {
            "type": "object",
            "properties": {
                "href": {
                    "$ref": "#/definitions/anyUri"
                },
                "op": {
                    "oneOf": [{
                            "type": "string",
                            "enum": [
                                "readproperty",
                                "writeproperty",
                                "observeproperty",
                                "unobserveproperty"
                            ]
                        },
                        {
                            "type": "array",
                            "items": {
                                "type": "string",
                                "enum": [
                                    "readproperty",
                                    "writeproperty",
                                    "observeproperty",
                                    "unobserveproperty"
                                ]
                            }
                        }
                    ]
                },
                "contentType": {
                    "type": "string"
                },
                "security": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "scopes": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "subProtocol": {
                    "type": "string",
                    "enum": [
                        "longpoll",
                        "websub",
                        "sse"
                    ]
                },
                "response": {
                    "type": "object",
                    "properties": {
                        "contentType": {
                            "type": "string"
                        }
                    }
                }
            },
            "required": [
                "href"
            ],
            "additionalProperties": true
        },
        "form_element_action": {
            "type": "object",
            "properties": {
                "href": {
                    "$ref": "#/definitions/anyUri"
                },
                "op": {
                    "oneOf": [{
                            "type": "string",
                            "enum": [
                                "invokeaction"
                            ]
                        },
                        {
                            "type": "array",
                            "items": {
                                "type": "string",
                                "enum": [
                                    "invokeaction"
                                ]
                            }
                        }
                    ]
                },
                "contentType": {
                    "type": "string"
                },
                "security": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "scopes": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "subProtocol": {
                    "type": "string",
                    "enum": [
                        "longpoll",
                        "websub",
                        "sse"
                    ]
                },
                "response": {
                    "type": "object",
                    "properties": {
                        "contentType": {
                            "type": "string"
                        }
                    }
                }
            },
            "required": [
                "href"
            ],
            "additionalProperties": true
        },
        "form_element_event": {
            "type": "object",
            "properties": {
                "href": {
                    "$ref": "#/definitions/anyUri"
                },
                "op": {
                    "oneOf": [{
                            "type": "string",
                            "enum": [
                                "subscribeevent",
                                "unsubscribeevent"
                            ]
                        },
                        {
                            "type": "array",
                            "items": {
                                "type": "string",
                                "enum": [
                                    "subscribeevent",
                                    "unsubscribeevent"
                                ]
                            }
                        }
                    ]
                },
                "contentType": {
                    "type": "string"
                },
                "security": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "scopes": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "subProtocol": {
                    "type": "string",
                    "enum": [
                        "longpoll",
                        "websub",
                        "sse"
                    ]
                },
                "response": {
                    "type": "object",
                    "properties": {
                        "contentType": {
                            "type": "string"
                        }
                    }
                }
            },
            "required": [
                "href"
            ],
            "additionalProperties": true
        },
        "form_element_root": {
            "type": "object",
            "properties": {
                "href": {
                    "$ref": "#/definitions/anyUri"
                },
                "op": {
                    "oneOf": [{
                            "type": "string",
                            "enum": [
                                "readallproperties",
                                "writeallproperties",
                                "readmultipleproperties",
                                "writemultipleproperties"
                            ]
                        },
                        {
                            "type": "array",
                            "items": {
                                "type": "string",
                                "enum": [
                                    "readallproperties",
                                    "writeallproperties",
                                    "readmultipleproperties",
                                    "writemultipleproperties"
                                ]
                            }
                        }
                    ]
                },
                "contentType": {
                    "type": "string"
                },
                "security": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "scopes": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "subProtocol": {
                    "type": "string",
                    "enum": [
                        "longpoll",
                        "websub",
                        "sse"
                    ]
                },
                "response": {
                    "type": "object",
                    "properties": {
                        "contentType": {
                            "type": "string"
                        }
                    }
                }
            },
            "required": [
                "href"
            ],
            "additionalProperties": true
        },
        "description": {
            "type": "string"
        },
        "title": {
            "type": "string"
        },
        "descriptions": {
            "type": "object"
        },
        "titles": {
            "type": "object"
        },
        "dataSchema": {
            "type": "object",
            "properties": {
                "@type": {
                    "$ref": "#/definitions/type_declaration"
                },
                "description": {
                    "$ref": "#/definitions/description"
                },
                "title": {
                    "$ref": "#/definitions/title"
                },
                "descriptions": {
                    "$ref": "#/definitions/descriptions"
                },
                "titles": {
                    "$ref": "#/definitions/titles"
                },
                "writeOnly": {
                    "type": "boolean"
                },
                "readOnly": {
                    "type": "boolean"
                },
                "oneOf": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dataSchema"
                    }
                },
                "unit": {
                    "type": "string"
                },
                "enum": {
                    "type": "array",
					"items": {
                        "type": "object"
                    },
                    "minItems": 1,
                    "uniqueItems": true
                },
                "format": {
                    "type": "string"
                },
                "const": {},
                "type": {
                    "type": "string",
                    "enum": [
                        "boolean",
                        "integer",
                        "number",
                        "string",
                        "object",
                        "array",
                        "null"
                    ]
                },
                "items": {
                    "oneOf": [{
                            "$ref": "#/definitions/dataSchema"
                        },
                        {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dataSchema"
                            }
                        }
                    ]
                },
                "maxItems": {
                    "type": "integer",
                    "minimum": 0
                },
                "minItems": {
                    "type": "integer",
                    "minimum": 0
                },
                "minimum": {
                    "type": "number"
                },
                "maximum": {
                    "type": "number"
                },
                "properties": {
                    "additionalProperties": {
                        "$ref": "#/definitions/dataSchema"
                    }
                },
                "required": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "link_element": {
            "type": "object",
            "properties": {
                "anchor": {
                    "$ref": "#/definitions/anyUri"
                },
                "href": {
                    "$ref": "#/definitions/anyUri"
                },
                "rel": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            },
            "required": [
                "href"
            ],
            "additionalProperties": true
        },
        "securityScheme": {
            "oneOf": [{
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "nosec"
                            ]
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "basic"
                            ]
                        },
                        "in": {
                            "type": "string",
                            "enum": [
                                "header",
                                "query",
                                "body",
                                "cookie"
                            ]
                        },
                        "name": {
                            "type": "string"
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "cert"
                            ]
                        },
                        "identity": {
                            "type": "string"
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "digest"
                            ]
                        },
                        "qop": {
                            "type": "string",
                            "enum": [
                                "auth",
                                "auth-int"
                            ]
                        },
                        "in": {
                            "type": "string",
                            "enum": [
                                "header",
                                "query",
                                "body",
                                "cookie"
                            ]
                        },
                        "name": {
                            "type": "string"
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "bearer"
                            ]
                        },
                        "authorization": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "alg": {
                            "type": "string",
                            "enum": [
                                "MD5",
                                "ES256",
                                "ES512-256"
                            ]
                        },
                        "format": {
                            "type": "string",
                            "enum": [
                                "jwt",
                                "jwe",
                                "jws"
                            ]
                        },
                        "in": {
                            "type": "string",
                            "enum": [
                                "header",
                                "query",
                                "body",
                                "cookie"
                            ]
                        },
                        "name": {
                            "type": "string"
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "psk"
                            ]
                        },
                        "identity": {
                            "type": "string"
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "public"
                            ]
                        },
                        "identity": {
                            "type": "string"
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "oauth2"
                            ]
                        },
                        "authorization": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "token": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "refresh": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scopes": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        },
                        "flow": {
                            "type": "string",
                            "enum": [
                                "implicit",
                                "password",
                                "client",
                                "code"
                            ]
                        }
                    },
                    "required": [
                        "scheme",
                        "flow"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "apikey"
                            ]
                        },
                        "in": {
                            "type": "string",
                            "enum": [
                                "header",
                                "query",
                                "body",
                                "cookie"
                            ]
                        },
                        "name": {
                            "type": "string"
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                },
                {
                    "type": "object",
                    "properties": {
                        "@type": {
                            "$ref": "#/definitions/type_declaration"
                        },
                        "description": {
                            "$ref": "#/definitions/description"
                        },
                        "descriptions": {
                            "$ref": "#/definitions/descriptions"
                        },
                        "proxy": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "scheme": {
                            "type": "string",
                            "enum": [
                                "pop"
                            ]
                        },
                        "authorization": {
                            "$ref": "#/definitions/anyUri"
                        },
                        "format": {
                            "type": "string",
                            "enum": [
                                "jwt",
                                "jwe",
                                "jws"
                            ]
                        },
                        "alg": {
                            "type": "string",
                            "enum": [
                                "MD5",
                                "ES256",
                                "ES512-256"
                            ]
                        },
                        "in": {
                            "type": "string",
                            "enum": [
                                "header",
                                "query",
                                "body",
                                "cookie"
                            ]
                        },
                        "name": {
                            "type": "string"
                        }
                    },
                    "required": [
                        "scheme"
                    ]
                }
            ]
        },
        "anyUri": {
            "type": "string",
            "format": "iri-reference"
        }
    },
    "type": "object",
    "properties": {
        "id": {
            "type": "string",
            "format": "uri"
        },
        "title": {
            "$ref": "#/definitions/title"
        },
        "titles": {
            "$ref": "#/definitions/titles"
        },
        "properties": {
            "type": "object",
            "additionalProperties": {
                "$ref": "#/definitions/property_element"
            }
        },
        "actions": {
            "type": "object",
            "additionalProperties": {
                "$ref": "#/definitions/action_element"
            }
        },
        "events": {
            "type": "object",
            "additionalProperties": {
                "$ref": "#/definitions/event_element"
            }
        },
        "description": {
            "$ref": "#/definitions/description"
        },
        "descriptions": {
            "$ref": "#/definitions/descriptions"
        },
        "version": {
            "type": "object",
            "properties": {
                "instance": {
                    "type": "string"
                }
            },
            "required": [
                "instance"
            ]
        },
        "links": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/link_element"
            }
        },
        "forms": {
            "type": "array",
            "minItems": 1,
            "items": {
                "$ref": "#/definitions/form_element_root"
            }
        },
        "base": {
            "$ref": "#/definitions/anyUri"
        },
        "securityDefinitions": {
            "type": "object",
            "minProperties": 1,
            "additionalProperties": {
                "$ref": "#/definitions/securityScheme"
            }
        },
        "support": {
            "$ref": "#/definitions/anyUri"
        },
        "created": {
            "type": "string"
        },
        "modified": {
            "type": "string"
        },
        "security": {
            "type": "array",
            "minItems": 1,
            "items": {
                "type": "string"
            }
        },
        "@type": {
            "$ref": "#/definitions/type_declaration"
        },
        "@context": {
            "$ref": "#/definitions/thing-context"
        }
    },
    "required": [
        "title",
        "security",
        "securityDefinitions",
        "@context"
    ],
    "additionalProperties": true
}
`