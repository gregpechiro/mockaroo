package cli

var FILE = `package main

import (
	"encoding/json"
    {{ if ne .fullPackage "main" }}"{{ .fullPackage }}"{{ end }}
	"github.com/gregpechiro/mockaroo"
)

func main() {

	s := {{ if ne .fullPackage "main" }}{{ .package }}.{{ end }}{{ .struct }}{}

	mockTypes := mockaroo.NewMockTypes("{{ .fullPackage }}", "{{ .struct }}", &s)
	b := mockaroo.GetData(mockTypes.MTypes, {{ .count }})

	var ss []{{ if ne .fullPackage "main" }}{{ .package }}.{{ end }}{{ .struct }}
	if err := json.Unmarshal(b, &ss); err != nil {
		panic(err)
	}

	mockTypes.Setup.Vars = ss

	mockaroo.GenVars(mockTypes)

}
`
