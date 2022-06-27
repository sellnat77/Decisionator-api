package handlers

import (
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world!")
}
