package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// just write to stdout for now.
  log.SetLevel(log.DebugLevel)
  log.SetOutput(os.Stdout)

  server := flag.String("addr", "127.0.0.1:8081", "listening addr")
  flag.Parse()

	err := run(*server)
	if err != nil {
		log.Fatal(err)
	}
}

// run starts a http.Server for the passed in address
// with all requests handled by echoServer.
func run(server string) error {

	l, err := net.Listen("tcp", server)
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	ps := NewProxyServer()
	ps.logf = log.Printf
	s := &http.Server{
		Handler: ps,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
