package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", "127.0.0.1:8080", "Address to listen for HTTP")

type valueForKey func(string) *string

func main() {
	flag.Parse()

	http.HandleFunc("/", showMsg(fromEnv("MSG")))

	fmt.Printf("Listening on %s\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func fromEnv(key string) func() *string {
	return func() *string {
		v, ok := os.LookupEnv(key)
		if !ok {
			return nil
		}
		return &v
	}
}

func showMsg(f func() *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := f()
		if v == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(*v))
	}
}
