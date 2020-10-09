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
	"gitlab.com/inview-team/raptor_team/registry/internal/server"
)

const (
	addr        = "0.0.0.0:1337"
	db_host     = "127.0.0.1:1338"
	db_user     = "user"
	db_pswd     = "password"
	db_database = "default"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	base, err := db.InitDB(db_host, db_user, db_pswd, db_database, ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connecto to MongoDB: %w", err))
	}
	mongo := db.New(base)
	registry := registry.New(mongo)

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
