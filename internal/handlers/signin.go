package handlers

import (
	"fmt"
	"net/http"
)

func SignIn(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Hello from sign in")
}
