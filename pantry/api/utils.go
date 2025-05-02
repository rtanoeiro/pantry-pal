package api

import (
	"net/http"
	"strings"
)

func respondWithError(writer http.ResponseWriter, code int, msg string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(code)
	_, errWriter := writer.Write([]byte(msg))

	if errWriter != nil {
		return
	}
}

func respondWithJSON(writer http.ResponseWriter, code int, data []byte) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(code)
	_, errWriter := writer.Write(data)

	if errWriter != nil {
		respondWithError(writer, http.StatusInternalServerError, errWriter.Error())
	}
}

func IdOrEmail(url string) string {
	// Add regex in the future for Email and ID
	if strings.Contains(url, "@") {
		return "Email"
	} else {
		return "ID"
	}
}
