package API_2_Concurrency

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func Post(wg *sync.WaitGroup, t *gin.Context) {
	defer wg.Done()
	Address := os.Getenv("Address")
	resp, err1 := http.Post(Address, "application/json", bytes.NewBuffer(jsonData))
	if err1 != nil {
		fmt.Println("Error:", err1)
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

type PostApiTest struct {
	Concurrency int `json:"-"`
}

var (
	wg       sync.WaitGroup
	post     map[string]interface{}
	jsonData []byte
	err      error
)

func Concurrency1(t *gin.Context) {

	if err = t.ShouldBindJSON(&post); err != nil {
		t.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	jsonData, err = json.Marshal(post)
	if err != nil {
		t.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}
	secret, _ := post["Concurrency"].(string)

	delete(post, "Concurrency")
	num, err1 := strconv.Atoi(secret)
	if err1 != nil {
		t.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert to integer"})
	}

	for i := 0; i <= (num - 1); i++ {
		wg.Add(1)
		go Post(&wg, t)

	}

	wg.Wait()

}
