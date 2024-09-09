package main

import (
	"flag"
	"log/slog"
)

func main() {
	var (
		httpAddr   string
		configPath string
	)
	flag.StringVar(&httpAddr, "http", ":8080", "HTTP server address")
	flag.StringVar(&configPath, "config", "config.toml", "config path")
	flag.Parse()

	cfg, err := NewConfig(configPath)
	if err != nil {
		slog.Error("failed parsing confg", err, slog.String("path", configPath))
		return
	}
	
	storage := UniversalStorage{}
	
	if *cfg.S3.Type == "local" {
		localStore, err := NewLocalStorage(cfg.S3)
		if err != nil {
			slog.Error("failed initializing S3 storage", err)
			return
		}
		storage.local = localStore
	}
	
	if *cfg.S3.Type == "s3" {
		s3Store, err := NewS3Storage(cfg.S3)
		if err != nil {
			slog.Error("failed initializing S3 storage", err)
			return
		}
		storage.s3 = s3Store
	}

	
	mediaLib := NewMediaLibrary(&storage)
	authLib := NewAuthLibrary(cfg)
	slog.Info("started HTTP server", slog.String("address", httpAddr))
	err = StartServer(mediaLib, authLib, httpAddr)
	slog.Error("failed starting HTTP server", err)
}
