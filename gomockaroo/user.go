package main

type User struct {
	SomeMap         map[int]interface{}
	FullName        string
	Age             int
	Race            string
	Active          bool
	HomeAddress     Address
	FavoritePlaces  []Address
	FriendsFullName []string
	MotherFullName  string
	IPAddress       string
	CostMoney       string
	Percentage      float32
}
