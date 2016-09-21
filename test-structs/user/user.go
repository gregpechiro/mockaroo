package user

import "github.com/gregpechiro/mockaroo/address"

// type User struct {
// 	Name   string
// 	Age    int
// 	Active bool
// }

type User struct {
	Name           string
	Age            int
	Active         bool
	HomeAddress    address.Address
	FavoritePlaces []address.Address
	Friends        []string
	Mother         string
}
