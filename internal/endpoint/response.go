package endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const supportedVersion = "2.0"

type response struct {
	Version string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *apiError       `json:"error,omitempty"`
}

func newResponse(id int64, result json.RawMessage) response {
	return response{
		Version: supportedVersion,
		ID:      id,
		Result:  result,
	}
}

func newErrResponse(id int64, err apiError) response {
	return response{
		Version: supportedVersion,
		ID:      id,
		Error:   &err,
	}
}

func writeResponse(w http.ResponseWriter, res response) {
	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(response)))
	if _, err := w.Write(response); err != nil {
		log.Printf("faile to write response: %s", err.Error())
	}
}
