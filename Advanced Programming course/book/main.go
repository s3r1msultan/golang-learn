package main

import (
	"log"
	"net/http"
)

const PORT = ":8080"

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/postrequest", postRequestHandler)

	log.Printf("Server is running on the %v port", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(res http.ResponseWriter, req *http.Request) {
	_, err := res.Write([]byte("I get your request"))
	if err != nil {
		return
	}
}

func postRequestHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.Header().Set("Allow", "POST")
		res.WriteHeader(http.StatusMethodNotAllowed)
		_, err := res.Write([]byte("Method is not allowed"))
		if err != nil {
			http.Error(res, "Error suddenly occured", http.StatusBadRequest)
			return
		}
		return
	}
	_, err := res.Write([]byte("I get your POST request"))
	if err != nil {
		return
	}
}
