package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SQLDSN       string `default:"root@/asterisk"`
	TLSCert      string
	TLSKey       string
	HTTPAddr     string        `default:":9080"`
	HTTPSAddr    string        `default:":9443"`
	CacheTime    time.Duration `default:"1m"`
	CAPath       string
	HTTPOnlyCert bool   `default:"false"`
	cache        *Cache `ignored:"true"`
}

func main() {
	config := new(Config)
	if err := envconfig.Process("", config); err != nil {
		log.Fatalln("Could not process configuration from environment:", err)
	}

	config.cache = NewCache(config.CacheTime)

	httpMux := http.NewServeMux()
	if !config.HTTPOnlyCert {
		httpMux.HandleFunc("/yealink.xml", config.YealinkHandler)
		httpMux.HandleFunc("/grandstream/phonebook.xml", config.GrandstreamHandler)
		httpMux.Handle("/files/", http.FileServer(http.FS(files)))
	}
	if config.CAPath != "" {
		httpMux.HandleFunc("/ca.crt", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, config.CAPath)
		})
	}

	httpsMux := http.NewServeMux()
	httpsMux.HandleFunc("/yealink.xml", config.YealinkHandler)
	httpsMux.HandleFunc("/grandstream/phonebook.xml", config.GrandstreamHandler)
	httpsMux.Handle("/files/", http.FileServer(http.FS(files)))
	if config.CAPath != "" {
		httpsMux.HandleFunc("/ca.crt", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, config.CAPath)
		})
	}

	httpHandler := handlers.CombinedLoggingHandler(os.Stdout, handlers.CompressHandler(httpMux))
	httpsHandler := handlers.CombinedLoggingHandler(os.Stdout, handlers.CompressHandler(httpsMux))

	runHTTP := config.HTTPAddr != "" && (!config.HTTPOnlyCert || (config.HTTPOnlyCert && config.CAPath != ""))
	runHTTPS := config.TLSCert != "" && config.TLSKey != "" && config.HTTPSAddr != ""
	if !(runHTTP || runHTTPS) {
		log.Fatalln("Could not start: not configured to run in HTTP or HTTPS mode")
	}

	if runHTTP {
		go func() {
			log.Println("INFO: listening on", config.HTTPAddr)
			err := http.ListenAndServe(config.HTTPAddr, httpHandler)
			if err != nil {
				log.Fatalln("ERROR: could not start HTTP listener:", err)
			}
		}()
	}

	if runHTTPS {
		go func() {
			log.Println("INFO: listening on", config.HTTPSAddr)
			err := http.ListenAndServeTLS(config.HTTPSAddr, config.TLSCert, config.TLSKey, httpsHandler)
			if err != nil {
				log.Fatalln("ERROR: could not start HTTPS listener:", err)
			}
		}()
	}

	select {}
}
