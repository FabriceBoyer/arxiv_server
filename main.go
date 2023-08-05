package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"arxiv_server/arxiv"
	"arxiv_server/utils"

	"github.com/gorilla/mux"
)

var root_path = utils.GetEnv("DUMP_PATH", "./dump/")

var mgr = arxiv.ArxivMetadataManager{Root_path: root_path}

func main() {
	err := mgr.InitializeManager()
	if err != nil {
		panic(err)
	}

	fmt.Print("Serving requests\n")
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", http.FileServer(http.Dir("./static")))
	router.HandleFunc("/search", utils.ErrorHandler(handleSearch))
	router.HandleFunc("/id", utils.ErrorHandler(handlePage))

	log.Fatal(http.ListenAndServe(":9097", router))
}

func handleSearch(w http.ResponseWriter, r *http.Request) error {
	// id := r.URL.Query().Get("page")

	// titles, err := arxiv.SearchTitles(id)
	// if err != nil {
	// 	return err
	// }

	// err = json.NewEncoder(w).Encode(titles)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func handlePage(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("page")

	elm, err := mgr.GetIndexedArxivMetadata(id)
	if err != nil {
		return err
	}

	val, err := json.Marshal(elm)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(string(val)))
	if err != nil {
		return err
	}

	return nil
}
