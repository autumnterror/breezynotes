package main

import (
	"github.com/autumnterror/breezynotes/internal/redis/config"
	"log"
)

func main() {
	//TODO
	config.MustSetup()
	log.Println(config.Test())
}
