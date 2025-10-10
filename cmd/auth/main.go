package main

import (
	"github.com/autumnterror/breezynotes/internal/auth/config"
	"log"
)

func main() {
	//TODO
	config.MustSetup()
	log.Println(config.Test())
}
