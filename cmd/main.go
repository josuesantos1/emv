package main

import (
	"encoding/hex"
	"log"

	domain "github.com/josuesantos1/emv/internal/domain"
	"github.com/josuesantos1/emv/internal/gateway"
	"github.com/josuesantos1/emv/internal/handlers"
	"github.com/josuesantos1/emv/internal/logger"
	"github.com/josuesantos1/emv/pkg/tlv"
)

func main() {
	transactionLogger := logger.NewJSONLogger("transactions.json")

	raw := "5A0845395787636214865F2404251200009F340400000000"
	data, err := hex.DecodeString(raw)
	if err != nil {
		log.Fatalf("Failed to decode TLV data: %v", err)
	}

	parser := tlv.Parser{}
	tlvs, err := parser.Parse(data)
	if err != nil {
		log.Fatalf("Failed to parse TLV: %v", err)
	}

	transaction := &domain.Tlv{}
	if err := transaction.Populate(tlvs); err != nil {
		log.Fatalf("Failed to populate transaction: %v", err)
	}

	gw := gateway.NewHTTPGateway("http://localhost:8080")

	result, err := handlers.ProcessTransaction(transaction, gw)
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}

	if err := transactionLogger.Log(result); err != nil {
		log.Printf("Warning: Failed to log transaction: %v", err)
	}
}
