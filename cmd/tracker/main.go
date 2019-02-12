package main

import (
	"context"
	"flag"
	"log"
	"os"

	"gitlab.com/easywork.me/backend/httpservice"
	"gitlab.com/easywork.me/backend/storage"
)

var (
	addr       = flag.String("addr", env("EW_ADDR", ":8000"), "Addr of service")
	mongodbURI = flag.String("mongo", env("EW_MONGODB_URI", "mongodb://localhost:27017"), "URI of mongodb")

	isDebug = flag.Bool("debug", false, "Start service in debug mode")
)

type Config struct {
	IsDebug bool

	Server httpservice.Config

	MongoDB struct {
		URI string
	}
}

func main() {
	flag.Parse()

	var cfg Config
	cfg.Server.Addr = *addr
	cfg.MongoDB.URI = *mongodbURI
	cfg.IsDebug = *isDebug

	if cfg.IsDebug {
		log.Printf("Config, %v", cfg)
	}

	s, err := storage.NewMongoDB(context.Background(), cfg.MongoDB.URI)
	if err != nil {
		log.Fatalf("Error on get storage: %v", err)
	}
	err = s.Init(context.Background())
	if err != nil {
		log.Fatalf("Error on init storage: %v", err)
	}

	err = httpservice.Start(cfg.Server, cfg.IsDebug, s)
	if err != nil {
		log.Fatalf("Error on start service: %v", err)
	}
}

func env(key, defValue string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return defValue
}
