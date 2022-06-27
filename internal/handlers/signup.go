package handlers

import (
	"fmt"
	"net/http"
)

func SignUp(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Hello from sign up")
}
