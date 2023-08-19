package main

import (
	"github.com/globbie/knowdy-media/internal/app/interface/webgate"
	"github.com/globbie/knowdy-media/internal/app/storage/memstore"
	"github.com/globbie/knowdy-media/internal/app/usecases/upload"

	"github.com/gofrs/uuid"
        "context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	ListenAddress string `json:"listen-address"`
	Path string
}

var (
	config    Config
)

const (
	memStoreId = "Simple Storage 00"
)

func init() {
	flag.StringVar(&config.ListenAddress, "listen-address", "127.0.0.1:8082", "http server listen address")
	flag.StringVar(&config.Path, "config-path", "/etc/knd-media/config.gsl", "path to Glottie config")
	flag.Parse()	
}

func main() {
	// initialize uuid generator
	_ = uuid.Must(uuid.NewV4())

	var mem = memstore.New(memStoreId)

	fs := &upload.FileStorage{
		FileMetaSaver: mem,
		FileMetaQuery: mem,
	}
	wg := webgate.New(fs)

	server := http.Server{
		Handler:      wg.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
		Addr:         config.ListenAddress,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("shutting down Knowdy Media Processor...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalln("could not gracefully shutdown the server:", server.Addr)
		}
		close(done)
	}()

	log.Println("Knowdy Media Processor is ready to handle requests at:",
		config.ListenAddress)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s, err: %v\n",
		server.Addr, err)
	}

	<-done
	log.Println("Knowdy Media Processor stopped")
}
