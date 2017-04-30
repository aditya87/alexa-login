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
	modified    bool
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
		s.modified = true
		io.WriteString(w, "{}")
	} else {
		amazonURL := fmt.Sprintf("https://www.amazon.com/s/field-keywords=%s", s.lastKeyword)

		resp := struct {
			AmazonURL string `json:"url"`
			Modified  bool   `json:"modified"`
		}{
			AmazonURL: amazonURL,
			Modified:  s.modified,
		}

		s.modified = false

		body, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(body))
	}
}

func main() {
	port := os.Getenv("PORT")
	s := server{
		lastKeyword: "echo",
		modified:    false,
	}
	http.ListenAndServe(":"+port, &s)
}
