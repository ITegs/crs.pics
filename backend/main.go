package main

import (
	"fmt"
	"github.com/ITegs/crs.pics/apiserver"
	"github.com/ITegs/crs.pics/cloudprovider"
	"github.com/ITegs/crs.pics/database"
)

func main() {
	fmt.Println("Welcome to crs.pics")

	db := database.NewDB()
	cp := cloudprovider.NewCloudProvider()

	api := apiserver.NewApiServer(db, cp)

	api.Serve()
}
