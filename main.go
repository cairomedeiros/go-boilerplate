package main

import (
	"fmt"

	"github.com/cairomedeiros/go-boilerplate/config"
	router "github.com/cairomedeiros/go-boilerplate/routes"
)

func main() {

	//initialize configs
	err := config.Init()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	router.Initialize()

}
