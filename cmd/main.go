package main

import (
        "github.com/gorilla/mux"

        "context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

import "C"

type Config struct {
	ListenAddress string `json:"listen-address"`
	Path string
}

var (
	config    Config
)

func init() {
	flag.StringVar(&config.ListenAddress, "listen-address", "127.0.0.1:8069", "http server listen address")
	flag.StringVar(&config.Path, "config-path", "/etc/knd-media/config.gsl", "path to Glottie config")
	flag.Parse()

	//gltConfigBytes, err := ioutil.ReadFile(config.ConfigPath)
	//if err != nil {
	//	log.Fatalln("failed to read knd-media config, error:", err)
	//}
	//gltConfig = string(gltConfigBytes)
}

func main() {

	router := http.NewServeMux()
	router.Handle("/media",      mediaHandler(config))
	router.Handle("/media/task", taskHandler(config))

	server := http.Server{
		Handler:      logger(router), // todo
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

	log.Println("Knowdy Media Processor is ready to handle requests at:", config.ListenAddress)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s, err: %v\n",
		server.Addr, err)
	}

	<-done
	log.Println("Knowdy Media Processor stopped")
}

func mediaHandler(gp Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		t, ok := r.URL.Query()["t"]
		if !ok || len(t) < 1 {
			http.Error(w, "URL param t is missing", http.StatusBadRequest)
			return
		}
		var lang = "en"
		{
		    langs, ok := r.URL.Query()["lang"]
		    if !ok || len(langs) < 1 {
			http.Error(w, "URL param lang is missing", http.StatusBadRequest)
			return
		    }
		    lang = langs[0]
		}
		// TODO job
		var result = "OK"		
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = io.WriteString(w, result)
	})
}

func mediaTaskHandler(gp Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cs, ok := r.URL.Query()["lang"]
		if !ok || len(cs) < 1 {
			http.Error(w, "URL param lang is missing", http.StatusBadRequest)
			log.Println("no language coding system specified")
			return
		}
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		// TODO
		var result = "OK"
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, result)
	})
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path, r.URL.Query(), r.RemoteAddr, r.UserAgent())
		h.ServeHTTP(w, r)
	})
}
