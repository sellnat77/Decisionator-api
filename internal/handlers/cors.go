package handlers

import (
	"fmt"
	"net/http"

	"github.com/sellnat77/Decisionator-api/internal/util"
)

func enableCors(w *http.ResponseWriter) {
	corsAddr := util.GetEnvVar("CORS_URL")
	(*w).Header().Set("Access-Control-Allow-Origin", corsAddr)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
}

func Preflight(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "Header token is 2000000")
}
