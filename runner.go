package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const validPIN = "1234" // <-- CHANGE THIS

func runHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Starting...")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Pin       string `json:"pin"`
		TargetURL string `json:"target_url"`
	}
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &body)

	if body.Pin != validPIN {
		http.Error(w, "Invalid PIN", http.StatusUnauthorized)
		return
	}

	cmd := exec.Command("go", "run", "main.go")
	cmd.Env = append(os.Environ(),
		"TEMP=C:\\Temp",
		"TMP=C:\\Temp",
		"GOCACHE=C:\\GoCache",
		fmt.Sprintf("TARGET_URL=%s", body.TargetURL),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed: %s", output), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Started successfully:\n" + string(output)))
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/run", runHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
