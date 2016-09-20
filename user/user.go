package user

// type User struct {
// 	Name   string
// 	Age    int
// 	Active bool
// }

type User struct {
	Name           string
	Age            int
	Active         bool
	HomeAddress    Address
	FavoritePlaces []Address
	Friends        []string
	Mother         string
}

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}
