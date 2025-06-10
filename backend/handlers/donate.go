package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golden-hearts/backend/mpesa"
)

type DonationRequest struct {
	ProjectID int    `json:"project_id"`
	Amount    int    `json:"amount"`
	Phone     string `json:"phone"`
}

func DonationsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DonationsHandler called for path: %s", r.URL.Path)
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var donation DonationRequest
		err := json.NewDecoder(r.Body).Decode(&donation)
		if err != nil || donation.Amount <= 0 || donation.Phone == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Call M-Pesa STK push
		err = mpesa.InitiateSTKPush(donation.Phone, donation.Amount, donation.ProjectID)
		if err != nil {
			http.Error(w, "Failed to initiate STK push", http.StatusInternalServerError)
			return
		}

		// Log donation as pending
		timestamp := time.Now().Format(time.RFC3339)
		logLine := fmt.Sprintf("%s - Project: %d - Phone: %s - Amount: %d - Status: pending\n",
			timestamp, donation.ProjectID, donation.Phone, donation.Amount)

		cwd, err := os.Getwd()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error getting the cwd\n%v", err)
			return
		}

		logFilePath := filepath.Join(cwd, "backend", "data", "donations.log")
		f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(logLine)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("STK Push sent"))
	})
}
