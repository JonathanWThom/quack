package main

import (
	"github.com/joho/godotenv"
	"github.com/jonathanwthom/quack/cmd"
)

func main() {
	// If .env file exists, use that, otherwise, use variables from OS
	_ = godotenv.Load()

	cmd.Execute()
}
