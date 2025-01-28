package docs

import _ "embed"

//go:embed swagger.json
var spec []byte

func Spec() []byte {
	return spec
}
