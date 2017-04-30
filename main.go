package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type server struct {
	lastKeyword string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/keyword" && r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusUnprocessableEntity)
			return
		}

		var r struct {
			Keyword string
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		s.lastKeyword = r.Keyword
		io.WriteString(w, "{}")
	} else {
		amazonURL := fmt.Sprintf("https://www.amazon.com/s/field-keywords=%s", s.lastKeyword)
		resp, err := http.Get(amazonURL)
		if err != nil {
			http.Error(w, "Couldn't query amazon", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Couldn't parse response from amazon", http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(body))
	}
}

func main() {
	port := os.Getenv("PORT")
	s := server{
		lastKeyword: "echo",
	}
	http.ListenAndServe(":"+port, &s)
}
