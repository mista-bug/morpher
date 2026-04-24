package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	port := ":6767"
	Routes(port)

	log.Printf("Running %s \n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error starting %s", port)
	}
}

func Routes(port string) {
	// health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Init(w, r)
	})
}

func Init(w http.ResponseWriter, r *http.Request) {

	var requestData map[string]any
	var fakeResults []map[string]any

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	countVal, ok := requestData["count"].(float64)
	if !ok {
		http.Error(w, "Invalid count", http.StatusInternalServerError)
		return
	}
	count := int(countVal)

	body, ok := requestData["body"].(map[string]any)
	if !ok {
		http.Error(w, "Invalid body", http.StatusInternalServerError)
		return
	}

	for i := 0; i < count; i++ {
		//parse keywords inside body
		fakeData, err := createFake(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fakeResults = append(fakeResults, fakeData)
	}

	sendResponse(fakeResults, w)
}

func sendResponse(fakeData any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fakeData)
}

func createFake(body map[string]any) (map[string]any, error) {
	fakeData := make(map[string]any)
	for k, v := range body {

		switch keywordValue := v.(type) {
		case string:
			fakeValue, err := parseKey(keywordValue)
			if err != nil {
				errorMessage := fmt.Sprintf("Keyword %s is not supported", keywordValue)
				return nil, fmt.Errorf(errorMessage)
			}

			fakeData[k] = fakeValue

		default:
			errorMessage := fmt.Sprintf("Keyword %s is not supported", keywordValue)
			return nil, fmt.Errorf(errorMessage)
		}
	}

	return fakeData, nil
}

func parseKey(key string) (any, error) {
	key = strings.ToLower(key)

	switch key {
	case "firstname":
		return gofakeit.FirstName(), nil
	case "lastname":
		return gofakeit.LastName(), nil
	case "number":
		return gofakeit.Number(0, 50), nil
	case "address":
		return gofakeit.Address(), nil
	default:
		// gofakeit.Generate()
		// gofakeit.FuncLookups()
		return nil, fmt.Errorf("Keyword %s is not supported", key)
	}

}
