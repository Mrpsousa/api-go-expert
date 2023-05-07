package main

import "github.com/mrpsousa/api/configs"

func main() {
	config, _ := configs.LoadConfig("../../")
	println(config.DBDriver)
}
