package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golden-hearts/backend/mpesa"
)

type DonationRequest struct {
	Name      string `json:"name"`
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
		if err := json.NewDecoder(r.Body).Decode(&donation); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		//  Get access token
		token, err := mpesa.GetAccessToken()
		if err != nil {
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}

		// Simulate Paybill payment via C2B
		err = SimulateC2BPayment(token, donation)
		if err != nil {
			http.Error(w, "Failed to simulate payment", http.StatusInternalServerError)
			return
		}

		//  Log donation
		log := fmt.Sprintf("%s | %s donated KES %d to project %v\n", time.Now().Format(time.RFC3339), donation.Name, donation.Amount, donation.ProjectID)
		_ = appendToFile("donations.txt", log)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "Donation simulated successfully"})
	})
}

func SimulateC2BPayment(token string, d DonationRequest) error {
	url := "https://sandbox.safaricom.co.ke/mpesa/c2b/v1/simulate"

	payload := map[string]interface{}{
		"ShortCode":     "600000", // Sandbox shortcode
		"CommandID":     "CustomerPayBillOnline",
		"Amount":        d.Amount,
		"Msisdn":        d.Phone,     // Should be 254708374149 in sandbox
		"BillRefNumber": d.ProjectID, // Used to tag donation to project
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log response for debugging
	resBody, _ := io.ReadAll(resp.Body)
	fmt.Println("M-Pesa Response:", string(resBody))

	return nil
}

func appendToFile(filePath string, text string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		return err
	}
	return nil
}
