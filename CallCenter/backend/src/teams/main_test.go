package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIndexHandler(t *testing.T) {

	// for _, env := range os.Environ() {
	// 	fmt.Println(env)
	// }

	fmt.Printf("\n NewVar:%s\n\n", os.Getenv("MONGO_SERVER_URI"))
	ts := httptest.NewServer(SetupServer())

	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", ts.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	responseData, _ := io.ReadAll(resp.Body)

	var apiResp = &ApiResponse{}
	err = json.Unmarshal(responseData, apiResp)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	//fmt.Printf("\napiResp: %s\n", apiResp.Data)

	// apiM, err := json.Marshal(apiResp.Data)
	// fmt.Printf("\nData: %s\n", apiM) // Data: [{"email":"prem@example.com","username":"prem"},{"email":"priya@example.com","username":"priya"}]

	if apiResp.Status != "success" {
		t.Fatalf("Expected hello world message, got %s", responseData)
	}

	// if string(responseData) != mockUserResp {
	// 	t.Fatalf("Expected hello world message, got %s", responseData)
	// }
}
