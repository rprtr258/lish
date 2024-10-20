package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "embed"

	readability "github.com/go-shiori/go-readability"
)

//go:embed tokens.txt
var reveryToken string

var punctRegExp = regexp.MustCompile(`[.,:;?!#%*()\[\]\{\}\\|\/<>~"\-_]`)

func tokenize(text string) []string {
	replaced := punctRegExp.ReplaceAllString(text, " ")
	return strings.Fields(replaced)
}

func main() {
	wordCoords, err := parseModelFile("./models/fasttext-commoncrawl-50k.vec")
	if err != nil {
		log.Fatalln("Could not read model file:", err)
	}
	fmt.Println("Parsed model file.")

	docs, err := parseDocsFile("./corpus/docs.json", wordCoords)
	if err != nil {
		log.Fatalln("Could not read corpus file:", err)
	}
	fmt.Println("Parsed docs corpus.")

	// web server
	http.HandleFunc("/similar", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("token") != strings.TrimSpace(reveryToken) {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "not authorized")
			return
		}

		maxResultsString := req.URL.Query().Get("n")
		maxResults, err := strconv.Atoi(maxResultsString)
		if err != nil {
			maxResults = 10 // default maxResults is 10
		}

		query := req.URL.Query().Get("q")
		tokens := tokenize(query)

		if len(tokens) == 0 && req.Method == "POST" {
			// maybe POST with data?
			body, err := io.ReadAll(req.Body)
			if err == nil {
				tokens = tokenize(string(body))
			}
		}

		if len(tokens) == 0 {
			// maybe there's a URL?
			url := req.URL.Query().Get("url")
			if url != "" {
				article, err := readability.FromURL(url, 10*time.Second)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					io.WriteString(w, "failed to read URL")
				}

				tokens = tokenize(article.Title + " " + article.TextContent)
			}
		}

		if len(tokens) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "invalid query")
			return
		}

		docVector := documentVector(wordCoords, tokens)
		similarDocs := closestDocs(docs[:], docVector, maxResults)
		respBytes, err := serializeDocs(similarDocs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "error encoding JSON")
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(respBytes)
	})

	http.HandleFunc("/js/token.js", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		io.WriteString(w, fmt.Sprintf("const REVERY_TOKEN = `%s`;", strings.TrimSpace(reveryToken)))
	})
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(":9998", nil))
}
