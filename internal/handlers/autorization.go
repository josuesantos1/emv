package handlers

import (
	"time"

	domain "github.com/josuesantos1/emv/internal/domain"
)

type TransactionResult struct {
	Approved     bool
	Message      string
	Pan          string
	DataValidade time.Time
	CVM          string
	Timestamp    time.Time
}

type Gateway interface {
	Authorize(transaction *domain.Tlv) (bool, error)
}

func ProcessTransaction(transaction *domain.Tlv, gateway Gateway) (*TransactionResult, error) {
	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	authorized, err := gateway.Authorize(transaction)
	if err != nil {
		return nil, err
	}

	result := &TransactionResult{
		Approved:     authorized,
		Pan:          transaction.Pan,
		DataValidade: transaction.DataValidade,
		CVM:          transaction.CVM,
		Timestamp:    time.Now(),
	}

	if authorized {
		result.Message = "Transaction authorized successfully"
	} else {
		result.Message = "Transaction rejected by gateway"
	}

	return result, nil
}
