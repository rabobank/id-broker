package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gomatbase/go-we"
	"github.com/rabobank/id-broker/cfg"
	"github.com/rabobank/id-broker/domain"
)

func dump(r *http.Request) {
	fmt.Printf("dumping %s request for URL: %s\n", r.Method, r.URL)
	fmt.Println("dumping request headers...")
	// Loop over header names
	for name, values := range r.Header {
		if name == "Authorization" {
			fmt.Printf(" %s: %s\n", name, "<redacted>")
		} else {
			// Loop over all values for the name.
			for _, value := range values {
				fmt.Printf(" %s: %s\n", name, value)
			}
		}
	}

	fmt.Println("dumping request body...")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
	} else if len(body) != 0 {
		var serviceInstance domain.ServiceInstance
		err = json.Unmarshal(body, &serviceInstance)
		if err != nil {
			fmt.Printf("failed to parse json object from request, error: %s\n", err)
			fmt.Println(string(body))
		} else {
			if jsonBytes, e := json.MarshalIndent(serviceInstance, "", "  "); e == nil {
				fmt.Println(string(jsonBytes))
			}
		}
	} else {
		fmt.Println("Empty request body...")
	}
	// Restore the io.ReadCloser to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(body))
}

func DebugFilter(_ http.Header, r we.RequestScope) error {
	if cfg.Debug {
		dump(r.Request())
	}

	return nil
}
