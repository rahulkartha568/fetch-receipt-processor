package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var (
	receiptsMutex sync.Mutex
	receiptsMap   = make(map[string]Receipt)
)

func generateUniqueID() string {
	// Generate UUID and return as String
	id := uuid.New()
	return id.String()
}

func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points++
		}
	}
	//fmt.Printf("Rule 1: ", points)

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && totalFloat == math.Floor(totalFloat) {
		points += 50
	}
	//fmt.Printf("Rule 2: ", points)

	// Rule 3: 25 points if the total is a multiple of 0.25
	if totalFloat != 0 && math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}
	//fmt.Printf("Rule 3: ", points)

	// Rule 4: 5 points for every two items on the receipt
	points += len(receipt.Items) / 2 * 5
	//fmt.Printf("Rule 4: ", points)

	// Rule 5: If the trimmed length of the item description is a multiple of 3, calculate points
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		//fmt.Printf("len " , trimmedLen)
		if trimmedLen%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(priceFloat * 0.2))
			}
		}
	}
	//fmt.Printf("Rule 5: ", points)

	// Rule 6: 6 points if the day in the purchase date is odd
	purchaseDateTime, err := time.Parse("2006-01-02 15:04", receipt.PurchaseDate+" "+receipt.PurchaseTime)
	if err == nil && purchaseDateTime.Day()%2 != 0 {
		points += 6
	}
	//fmt.Printf("Rule 6 Day", purchaseDateTime.Day())
	//fmt.Printf("Rule 6: ", points)

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	//fmt.Printf("Rule 7:  Purchase Time", purchaseTime)
	//fmt.Printf("Rule 7: ", err)
	if err == nil && purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}
	//fmt.Printf("Rule 7: ", points)

	return points
}

func processReceiptsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming JSON receipt
	var receivedReceipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receivedReceipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the receipt
	receivedReceiptID := generateUniqueID()

	// Store the receipt in the map
	receiptsMutex.Lock()
	receiptsMap[receivedReceiptID] = receivedReceipt
	receiptsMutex.Unlock()

	// Respond with the generated ID
	response := struct {
		ID string `json:"id"`
	}{
		ID: receivedReceiptID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the receipt ID from the URL path
	id := r.URL.Path[len("/receipts/") : len(r.URL.Path)-len("/points")]

	// Look up the receipt by ID
	receiptsMutex.Lock()
	receipt, exists := receiptsMap[id]
	receiptsMutex.Unlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Calculate points for the receipt
	points := calculatePoints(receipt)

	// Respond with the points
	response := struct {
		Points int `json:"points"`
	}{
		Points: points,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Create a router using Gorilla Mux
	router := mux.NewRouter()

	// Define the /receipts/process endpoint for processing receipts
	router.HandleFunc("/receipts/process", processReceiptsHandler).Methods("POST")

	// Define the /receipts/{id}/points endpoint for getting points
	router.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

	// Start the HTTP server
	port := ":3000"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(port, router)
}
