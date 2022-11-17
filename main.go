package main

import (
	"arxiv_server/arxiv"
	"arxiv_server/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	router.HandleFunc("/", handleHomePage)
	router.HandleFunc("/search/{id}", utils.ErrorHandler(handleSearch))
	router.HandleFunc("/id/{id}", utils.ErrorHandler(handlePage))

	log.Fatal(http.ListenAndServe(":9097", router))
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to arxiv server, please use API")
}

func handleSearch(w http.ResponseWriter, r *http.Request) error {
	// vars := mux.Vars(r)
	// key := vars["id"]

	// titles, err := arxiv.SearchTitles(key)
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
	//articleName := wikitext.URLToTitle(path.Base(r.URL.Path))
	vars := mux.Vars(r)
	id := vars["id"]

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
