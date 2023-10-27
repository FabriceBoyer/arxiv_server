package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fabriceboyer/arxiv_server/arxiv"
	"github.com/spf13/viper"

	"github.com/fabriceboyer/common_go_utils/utils"
	"github.com/gorilla/mux"
)

var mgr = arxiv.ArxivMetadataManager{}

func main() {

	err := utils.SetupConfigPath(".")
	if err != nil {
		panic(err)
	}

	rootPath := viper.GetString("DUMP_PATH")
	mgr.Root_path = rootPath
	fmt.Printf("Root path: %s\n", rootPath)

	err = mgr.InitializeManager()
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
	id := r.URL.Query().Get("arxiv")

	elm, err := mgr.GetIndexedArxivMetadata(id)
	if err != nil {
		return err
	}

	val, err := json.MarshalIndent(elm, "", " ")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(string(val)))
	if err != nil {
		return err
	}

	return nil
}
