package main

import (
	"fmt"
	"github.com/ITegs/crs.pics/apiserver"
	"github.com/ITegs/crs.pics/database"
)

func main() {
	fmt.Println("Welcome to crs.pics")

	db := database.NewDB()
	api := apiserver.NewApiServer(db)

	api.Serve()
}
