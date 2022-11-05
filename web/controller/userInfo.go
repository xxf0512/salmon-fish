package controller

import "salmon-fish/service"

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName string
	Password  string
	IsAdmin   string
}

var Users []User

func init() {
	admin := User{LoginName: "admin", Password: "123456", IsAdmin: "T"}
	alice := User{LoginName: "ChainDesk", Password: "123456", IsAdmin: "T"}
	bob := User{LoginName: "alice", Password: "123456", IsAdmin: "F"}
	jack := User{LoginName: "bob", Password: "123456", IsAdmin: "F"}

	Users = append(Users, admin)
	Users = append(Users, alice)
	Users = append(Users, bob)
	Users = append(Users, jack)
}
