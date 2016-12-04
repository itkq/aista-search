package main

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/route"
	"aista-search/session"
	"aista-search/view"
	"aista-search/view/plugin"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	config.LoadEnv()
	db.Connect()
	session.Configure()

	view.Configure()
	view.LoadPlugins(
		plugin.FormattedTime(),
		plugin.EpisodeStatus(),
		plugin.ImagePath(),
		plugin.ThumbnailPath(),
	)

	router := route.New()
	sock := "/tmp/app.sock"
	os.Remove(sock)
	l, err := net.Listen("unix", sock)
	if err != nil {
		fmt.Println("%s\n", err)
		return
	}
	os.Chmod(sock, 0777)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down", sig)
		l.Close()
		os.Exit(0)
	}(sigc)

	err = http.Serve(l, router)
	if err != nil {
		panic(err)
	}
}
