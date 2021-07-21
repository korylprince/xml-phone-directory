package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SQLDSN    string `default:"root@/asterisk"`
	TLSCert   string
	TLSKey    string
	HTTPAddr  string `default:":9080"`
	HTTPSAddr string `default:":9443"`
}

func main() {
	config := new(Config)
	if err := envconfig.Process("", config); err != nil {
		log.Fatalln("Could not process configuration from environment:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/yealink.xml", config.YealinkHandler)
	mux.HandleFunc("/grandstream.xml", config.GrandstreamHandler)
	middleware := handlers.CombinedLoggingHandler(os.Stdout, handlers.CompressHandler(mux))

	if config.HTTPAddr != "" {
		go func() {
			log.Println("INFO: listening on", config.HTTPAddr)
			err := http.ListenAndServe(config.HTTPAddr, middleware)
			if err != nil {
				log.Fatalln("ERROR: could not start HTTP listener:", err)
			}
		}()
	}

	if config.TLSCert != "" && config.TLSKey != "" && config.HTTPSAddr != "" {
		go func() {
			log.Println("INFO: listening on", config.HTTPSAddr)
			err := http.ListenAndServeTLS(config.HTTPSAddr, config.TLSCert, config.TLSKey, middleware)
			if err != nil {
				log.Fatalln("ERROR: could not start HTTPS listener:", err)
			}
		}()
	}

	select {}
}