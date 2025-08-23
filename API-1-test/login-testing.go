package API_1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type PostApiTest struct {
}

func TestPostApi(t *gin.Context) {
	var test PostApiTest
	if err := t.ShouldBindJSON(&test); err != nil {
		t.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	jsonData, err := json.Marshal(test)
	if err != nil {
		t.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}
	Address := os.Getenv("Address")
	resp, err := http.Post(Address, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var result map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	t.JSON(resp.StatusCode, result)
}
