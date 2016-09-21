package user

import "github.com/gregpechiro/mockaroo/test-structs/address"

// type User struct {
// 	Name   string
// 	Age    int
// 	Active bool
// }

type User struct {
	SomeMap        map[int]interface{}
	Name           string
	Age            int
	Active         bool
	HomeAddress    address.Address
	FavoritePlaces []address.Address
	Friends        []string
	Mother         string
}
