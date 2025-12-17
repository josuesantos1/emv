package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type AuthorizationRequest struct {
	Pan          string    `json:"pan"`
	DataValidade time.Time `json:"data_validade"`
	CVM          string    `json:"cvm"`
}

type AuthorizationResponse struct {
	Approved  bool      `json:"approved"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/authorize", authorizeHandler)

	port := ":8080"
	fmt.Printf("Mock server Acquirer running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	approved := rand.Intn(100) < 70

	response := AuthorizationResponse{
		Approved:  approved,
		Timestamp: time.Now(),
	}

	if approved {
		response.Message = "Transaction approved by acquirer"
	} else {
		response.Message = "Transaction declined by acquirer"
	}

	log.Printf("Authorization request: PAN=%s, Approved=%v\n", maskPan(req.Pan), approved)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(AuthorizationResponse{
		Approved:  false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func maskPan(pan string) string {
	if len(pan) < 10 {
		return pan
	}
	return pan[:6] + "******" + pan[len(pan)-4:]
}