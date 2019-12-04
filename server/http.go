package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// 主动关闭服务器
var server *http.Server

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/closed", sayBye)

	server = &http.Server{
		Addr:         ":9090",
		WriteTimeout: time.Second * 4,
		Handler:      mux,
	}

	go func() {
		// 接收退出信号
		<-quit
		if err := server.Close(); err != nil {
			log.Fatal("Close server:", err)
		}
	}()

	log.Println("Starting v3 httpserver")
	err := server.ListenAndServe()
	if err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			log.Fatal("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected", err)
		}
	}

	log.Fatal("Server exited")

}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	var Agents []string = header["User-Agent"]
	log.Println(Agents)
	w.Write([]byte("this is version 3"))
}

// 关闭http
func sayBye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bye bye ,shutdown the server"))
	err := server.Shutdown(nil)
	if err != nil {
		log.Fatal("shutdown the server err")
	}
}
