package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/inview-team/raptor_team/registry/server"
)

const addr = "0.0.0.0:1337"

func main() {
	srv := server.New(addr)
	done := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if err := srv.Stop(); err != nil {
			log.Fatal("server stopped with error: %w", err)
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
