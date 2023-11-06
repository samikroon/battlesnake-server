package resources

import (
	_ "embed"
)

//go:embed openapi.json
var OpenApiSpec []byte

//go:embed redoc.html
var Redoc []byte
