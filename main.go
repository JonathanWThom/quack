package main

import (
	"github.com/joho/godotenv"
	"github.com/jonathanwthom/quack/cmd"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmd.Execute()
}
