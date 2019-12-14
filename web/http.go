package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func runHttpServer(stop chan os.Signal) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hellowWorld)
	mux.HandleFunc("/ws", handleWs)
	server := &http.Server{Addr: "127.0.0.1:8088", Handler: mux}

	go func() {
		log.Printf("Listening on 8088 \n")

		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// stop server
	<-stop

	log.Println("shutting done server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// say hello to root visit
func hellowWorld(w http.ResponseWriter, req *http.Request) {

}

func handleWs(w http.ResponseWriter, req *http.Request) {

}

func handleSpdy(w http.ResponseWriter, req *http.Request) {

}
