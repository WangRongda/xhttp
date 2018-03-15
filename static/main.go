package main

import (
	"net/http"
	"os"
	"log"
	"flag"
	"github.com/WangRongda/xhttp"
)

func main() {
	rootPath := flag.String("root", os.Getenv("PWD"), "Static file server path")
	port := flag.String("port", "8778", "Server port")
	http.Handle("/", myHttp.Staticer{RootPath: *rootPath})
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
