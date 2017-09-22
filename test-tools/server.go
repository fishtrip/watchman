package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/test.html", func(w http.ResponseWriter, r *http.Request) {
		var t interface{}

		fmt.Printf("Request, %v\n", r)

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&t)

		if err != nil {
			fmt.Printf("decode error, %s\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//n := rand.Intn(10)
		//fmt.Printf("Prepare to Sleep %d seconds\n", n)
		//time.Sleep(time.Duration(n) * time.Second)

		fmt.Printf("Request Body, %v\n", t)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
