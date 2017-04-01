package http

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"weixin-provider/config"
)

func init() {
	configCommonRoutes()
	configProRoutes()
}

func Start() {

	addr := config.Config().Http.Listen
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	log.Println("http listening", addr)
	log.Fatalln(s.ListenAndServe())
}
