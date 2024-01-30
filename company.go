package main

type Company struct {
	Name               string
	Id                 string
	Managers           []User
	RootAccountManager User
}
