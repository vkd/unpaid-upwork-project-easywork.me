package main

import (
	"flag"
	"log"
	"os"

	"gitlab.com/easywork.me/backend/config"
)

var (
	addr       = flag.String("addr", env("EW_ADDR", ":8000"), "Addr of service")
	mongodbURI = flag.String("mongo", env("EW_MONGODB_URI", "mongodb://localhost:27017"), "URI of mongodb")

	isDebug = flag.Bool("debug", false, "Start service in debug mode")
)

func main() {
	flag.Parse()

	var cfg config.Config
	cfg.Addr = *addr
	cfg.MongoDB.URI = *mongodbURI

	log.Printf("OK, %v", cfg)
}

func env(key, defValue string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return defValue
}
