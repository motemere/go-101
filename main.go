package main

import (
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		log.WithError(err).Error("Failed to log to file, using default stdout")
		log.SetOutput(os.Stdout)
	}
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		logRequest(r)

		writeResponse([]byte("Hello World"), w)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		logRequest(r)

		writeResponse([]byte("OK"), w)
	})

	http.HandleFunc("/pages", func(w http.ResponseWriter, r *http.Request) {

		logRequest(r)

		pageNumber := r.URL.Query().Get("page")

		writeResponse([]byte("Page called: "+pageNumber), w)
	})

	port := "8088"

	log.WithFields(log.Fields{
		"port": port,
	}).Info("Server started")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func logRequest(r *http.Request) {
	log.WithFields(log.Fields{
		"URL":        r.URL.Path,
		"Params":     r.URL.Query(),
		"Method":     r.Method,
		"RemoteAddr": r.RemoteAddr,
		"UserAgent":  r.Header.Get("User-Agent"),
	}).Info("Request received")
}

func writeResponse(response []byte, w http.ResponseWriter) {
	_, err := w.Write(response)
	if err != nil {
		return
	}
}
