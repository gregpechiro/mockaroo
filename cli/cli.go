package cli

var FILE = `package main

import (
	"encoding/json"
    {{ if ne .fullPackage "main" }}"{{ .fullPackage }}"{{ end }}
	"github.com/gregpechiro/mockaroo"
)

func main() {
	{{ if gt .maxSlice 0 }}mockaroo.SLICERANGE = "[1-{{ .maxSlice }}]"
	mockaroo.JSONARRAYMAX = {{ .maxSlice }}{{ end }}

	s := {{ if ne .fullPackage "main" }}{{ .package }}.{{ end }}{{ .struct }}{}

	mockTypes := mockaroo.NewMockTypes("{{ .fullPackage }}", "{{ .struct }}", &s, {{ .match }})
	b := mockaroo.GetData(mockTypes.MTypes, {{ .count }})

	var ss []{{ if ne .fullPackage "main" }}{{ .package }}.{{ end }}{{ .struct }}
	if err := json.Unmarshal(b, &ss); err != nil {
		panic(err)
	}

	mockTypes.Setup.Vars = ss

	mockaroo.GenVars(mockTypes)

}
`
