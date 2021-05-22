package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/internal/config"
	"gitlab.com/inview-team/raptor_team/registry/internal/db"
	"gitlab.com/inview-team/raptor_team/registry/internal/rabbitmq"
	"gitlab.com/inview-team/raptor_team/registry/internal/server"
)

var (
	cfgFile string
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "path to config file")
}

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV")
	}

	var conf config.Settings
	if err := env.Parse(&conf); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongo, err := db.New(&conf, ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to MongoDB: %w", err))
	}

	pub := rabbitmq.NewPublisher(conf.RMQAddress)
	err = pub.Connect()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect publisher to RabbitMQ: %w", err))
	}
	defer pub.Close()

	registry := registry.New(&conf, mongo, pub)

	srv := server.New(conf.RegistryAddress, registry)

	done := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if err := srv.Stop(); err != nil {
			log.Fatal(fmt.Errorf("server stopped with error: %w", err))
			return
		}
	}()

	go func() {
		log.Printf("server started at %s", conf.RegistryAddress)
		errs <- srv.Start()
	}()

	select {
	case <-done:
		signal.Stop(done)
		return
	case err := <-errs:
		if err != nil {
			log.Fatal("server exited with error: %w", err)
		}
		return
	}
}
