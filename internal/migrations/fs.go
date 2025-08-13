package migrations

import (
	"embed"
)

//We are reserving this file for when we compile to binary
//Find all the sql files that are goint to existin the directory path here ->

//go:embed *.sql
var FS embed.FS