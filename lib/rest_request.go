package lib

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func ReadRequest(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func GetQueryInt(r *http.Request, key string, def int) int {
	value := r.URL.Query().Get(key)

	if value == "" {
		return def
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	return def
}
