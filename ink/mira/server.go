package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const dataFilePath = "./data/mira.txt"

func ensureDataFileExists() error {
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		dataFile, err := os.OpenFile(dataFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer dataFile.Close()

		// empty JSON array
		if _, err := dataFile.Write([]byte("[]")); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func getData(w http.ResponseWriter, r *http.Request) {
	dataFile, err := os.Open(dataFilePath)
	if err != nil {
		log.Println("file open on get", err.Error())
		w.WriteHeader(500)
		io.WriteString(w, "error reading file")
		return
	}
	defer dataFile.Close()

	io.Copy(w, dataFile)
}

func postData(w http.ResponseWriter, r *http.Request) {
	dataFile, err := os.OpenFile(dataFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("file open on post", err.Error())
		w.WriteHeader(500)
		io.WriteString(w, "error reading file")
		return
	}
	defer dataFile.Close()

	_, err = io.Copy(dataFile, r.Body)
	if err != nil {
		log.Println("copy from request:", err.Error())
		w.WriteHeader(500)
		return
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	indexFile, err := os.Open("./static/index.html")
	if err != nil {
		io.WriteString(w, "error reading index")
		return
	}
	defer indexFile.Close()

	io.Copy(w, indexFile)
}

func start() {
	if err := ensureDataFileExists(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handleHome)
	r.Methods("GET").Path("/data").HandlerFunc(getData)
	r.Methods("POST").Path("/data").HandlerFunc(postData)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8998",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	log.Println("Mira listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
