package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	domain "github.com/josuesantos1/emv/internal/domain"
	"github.com/josuesantos1/emv/internal/gateway"
	"github.com/josuesantos1/emv/internal/handlers"
	"github.com/josuesantos1/emv/internal/logger"
	"github.com/josuesantos1/emv/pkg/tlv"
)

func main() {
	transactionLogger := logger.NewJSONLogger("transactions.json")
	gw := gateway.NewHTTPGateway("http://localhost:8080")
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("EMV Transaction Processor")
	fmt.Println("=========================")
	fmt.Println("Enter TLV hex data (or 'exit' to quit)")
	fmt.Println("Exemple: 5A0845395787636214865F2404251200009F340400000000")
	fmt.Println()

	for {
		fmt.Print("TLV> ")

		if !scanner.Scan() {
			break
		}

		raw := strings.TrimSpace(scanner.Text())

		if raw == "" {
			continue
		}

		if raw == "exit" || raw == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if err := processTransaction(raw, transactionLogger, gw); err != nil {
			fmt.Printf("Error: %v\n\n", err)
		}
	}
}

func processTransaction(raw string, transactionLogger *logger.JSONLogger, gw *gateway.HTTPGateway) error {
	data, err := hex.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("invalid hex data: %v", err)
	}

	parser := tlv.Parser{}
	tlvs, err := parser.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse TLV: %v", err)
	}

	transaction := &domain.Tlv{}
	if err := transaction.Populate(tlvs); err != nil {
		return fmt.Errorf("failed to populate transaction: %v", err)
	}

	result, err := handlers.ProcessTransaction(transaction, gw)
	if err != nil {
		return fmt.Errorf("transaction processing failed: %v", err)
	}

	if err := transactionLogger.Log(result); err != nil {
		log.Printf("Warning: Failed to log transaction: %v", err)
	}

	fmt.Println("\n========== TRANSACTION RESULT ==========")
	if result.Approved {
		fmt.Println("Status: APPROVED")
	} else {
		fmt.Println("Status: REJECTED")
	}
	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("PAN: %s\n", result.Pan)
	fmt.Printf("Expiry Date: %s\n", result.DataValidade.Format("01/2006"))
	fmt.Printf("CVM: %s\n", result.CVM)
	fmt.Printf("Timestamp: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Println("========================================")

	return nil
}
