package handlers

import (
	"fmt"
	"net/http"

	"github.com/sellnat77/Decisionator-api/internal/util"
)

func Healthcheck(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	for _, c := range req.Cookies() {
		fmt.Println(c)
	}
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	healthy := util.Healthcheck()

	if !healthy {
		http.Error(w, "Healthcheck failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "healthy!")
}
