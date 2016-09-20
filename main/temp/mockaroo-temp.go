package main

import (
	"encoding/json"
    
	"github.com/gregpechiro/mockaroo"
	"github.com/gregpechiro/mockaroo/cli"
)

func main() {

	s := User{}

	mockTypes := mockaroo.NewMockTypes("main", "User", &s)
	b := mockaroo.GetData(mockTypes.MTypes, 1)

	var ss []User
	if err := json.Unmarshal(b, &ss); err != nil {
		panic(err)
	}

	mockTypes.Setup.Vars = ss

	cli.GenVars(mockTypes)

}
