package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "summary": "Get list of songs",
                "parameters": [
                    {"name": "group", "in": "query", "type": "string", "description": "Group"},
                    {"name": "song", "in": "query", "type": "string", "description": "Song"}
                ],
                "responses": {
                    "200": {"description": "OK", "schema": {"type": "array", "items": {"$ref": "#/definitions/models.Song"}}}
                }
            },
            "post": {
                "summary": "Add a new song",
                "parameters": [
                    {"name": "song", "in": "body", "schema": {"$ref": "#/definitions/handlers.AddSongRequest"}, "required": true}
                ],
                "responses": {
                    "200": {"description": "OK", "schema": {"$ref": "#/definitions/models.Song"}}
                }
            }
        },
        "/songs/{id}": {
            "get": {"summary": "Get song by ID", "parameters": [{"name": "id", "in": "path", "type": "integer"}], "responses": {"200": {"description": "OK", "schema": {"$ref": "#/definitions/models.Song"}}}},
            "put": {"summary": "Update a song", "parameters": [{"name": "id", "in": "path", "type": "integer"}, {"name": "song", "in": "body", "schema": {"$ref": "#/definitions/models.Song"}}], "responses": {"200": {"description": "OK", "schema": {"$ref": "#/definitions/models.Song"}}}},
            "delete": {"summary": "Delete a song", "parameters": [{"name": "id", "in": "path", "type": "integer"}], "responses": {"200": {"description": "OK"}}}
        }
    },
    "definitions": {
        "handlers.AddSongRequest": {"type": "object", "properties": {"group": {"type": "string"}, "song": {"type": "string"}}},
        "models.Song": {
            "type": "object",
            "properties": {
                "id": {"type": "integer", "example": 1},
                "created_at": {"type": "string"},
                "group": {"type": "string", "example": "Muse"},
                "title": {"type": "string", "example": "Supermassive Black Hole"},
                "release_date": {"type": "string", "example": "16.07.2006"},
                "text": {"type": "string"},
                "link": {"type": "string"}
            }
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	SwaggerTemplate:  docTemplate,
	InfoInstanceName: "swagger",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
