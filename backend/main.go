package main

import (
	"flag"
	"fmt"
	"os"

	config "github.com/sriharivishnu/shopify-challenge/config"
	"github.com/sriharivishnu/shopify-challenge/server"
	db "github.com/sriharivishnu/shopify-challenge/services"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	config.PopulateConfig()

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	db.Init()
	server.Init()
}
