package main

type User struct {
	FullName  string
	FirstName string
	LastName  string
	Address   Address
}

type Address struct {
	State  string
	Street string
	City   string
	zip    string
}
