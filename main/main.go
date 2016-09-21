package main

import (
	"encoding/json"
	"fmt"

	"github.com/gregpechiro/mockaroo"
	"github.com/gregpechiro/mockaroo/test-structs/user"
)

func main() {
	u := user.User{}
	mockTypes := mockaroo.NewMockTypes("github.com/gregpechiro/mockaroo/test-structs/user", "User", &u, true)
	fmt.Println()
	b1, _ := json.MarshalIndent(mockTypes.MTypes, "", "\t")
	fmt.Printf("%s\n", b1)
	fmt.Println()
	b := mockaroo.GetDataPretty(mockTypes.MTypes, 1)
	fmt.Printf("%s\n", b)
	// var ss []User
	// if err := json.Unmarshal(b, &ss); err != nil {
	// 	panic(err)
	// }
	//
	// mockTypes.Setup.Vars = ss
	//
	// cli.GenVars(mockTypes)
	fmt.Println()
	fmt.Println(mockTypes.Template)
}
