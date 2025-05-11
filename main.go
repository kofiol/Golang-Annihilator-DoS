package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	concurrent      = 2000
	totalRequests   = 20000
	requestValueLen = 330000
)

type Payload struct {
	UserID string `json:"user_id"`
	Action string `json:"action"`
	Data   string `json:"data"`
}

func banner() {
	fmt.Println("\033[1;34m╔══════════════════════════════════════════════╗")
	fmt.Println("║              GOLANG ANNIHILATOR              ║")
	fmt.Println("╚══════════════════════════════════════════════╝\033[0m")
}

func main() {
	targetURL := os.Getenv("TARGET_URL")
	if targetURL == "" {
		fmt.Println("Error: TARGET_URL not set.")
		return
	}

	banner()

	requestValue := make([]byte, requestValueLen)
	for i := range requestValue {
		requestValue[i] = '0'
	}

	payloadTemplate := Payload{
		Action: "test",
		Data:   string(requestValue),
	}

	start := time.Now()
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrent)
	client := http.DefaultClient

	for i := 1; i <= totalRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(userID int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			payload := payloadTemplate
			payload.UserID = fmt.Sprintf("user_%d", userID)

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				fmt.Printf("\033[1;31m[ERROR] user_%d: %v\033[0m\n", userID, err)
				return
			}

			req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(payloadBytes))
			if err != nil {
				fmt.Printf("\033[1;31m[ERROR] user_%d: %v\033[0m\n", userID, err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Go-HttpClient/1.0")
			req.Header.Set("Accept", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("\033[1;31m[FAIL] user_%d: %v\033[0m\n", userID, err)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("\033[1;31m[FAIL] user_%d: HTTP %d — %s\033[0m\n", userID, resp.StatusCode, string(body))
			} else {
				fmt.Printf("\033[1;32m[SENT] user_%d: HTTP %d\033[0m\n", userID, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\033[1;34m\nFinished sending %d requests in %.2f seconds\033[0m\n", totalRequests, elapsed.Seconds())
}
