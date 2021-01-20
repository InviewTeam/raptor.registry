package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/internal/db"
	"gitlab.com/inview-team/raptor_team/registry/internal/rabbitmq"
	"gitlab.com/inview-team/raptor_team/registry/internal/server"
)

var (
	addr = os.Getenv("REGISTRY_ADDR")

	db_host       = os.Getenv("MONGO_HOST")
	db_user       = os.Getenv("MONGO_USER")
	db_pswd       = os.Getenv("MONGO_PWD")
	db_database   = os.Getenv("MONGO_DBNAME")
	db_collection = os.Getenv("MONGO_COLL")

	rmq_addr = os.Getenv("RABBITMQ_ADDR")
	rmq_exch = os.Getenv("RABBITMQ_EXCHANGE")
	rmq_key  = os.Getenv("RABBITMQ_KEY")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongo, err := db.New(db_host, db_user, db_pswd, db_database, db_collection, ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to MongoDB: %w", err))
	}
	rmqConf := rabbitmq.Conf{
		Address:  rmq_addr,
		Exchange: rmq_exch,
		Key:      rmq_key,
	}
	pub := rabbitmq.NewPublisher(&rmqConf)
	err = pub.Connect()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect publisher to RabbitMQ: %w", err))
	}
	defer pub.Close()

	registry := registry.New(mongo, pub)

	srv := server.New(addr, registry)

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
		log.Printf("server started at %s", addr)
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
