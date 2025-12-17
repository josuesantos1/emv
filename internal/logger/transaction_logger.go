package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/josuesantos1/emv/internal/handlers"
)

type Log struct {
	ID           string    `json:"id"`
	Pan          string    `json:"pan"`
	DataValidade string    `json:"data_validade"`
	CVM          string    `json:"cvm"`
	Approved     bool      `json:"approved"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
}

type JSONLogger struct {
	filePath string
	mu       sync.Mutex
}

func NewJSONLogger(filePath string) *JSONLogger {
	return &JSONLogger{
		filePath: filePath,
	}
}

func (l *JSONLogger) Log(result *handlers.TransactionResult) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	logEntry := Log{
		ID:           generateID(),
		Pan:          result.Pan,
		DataValidade: result.DataValidade.Format("01/2006"),
		CVM:          result.CVM,
		Approved:     result.Approved,
		Message:      result.Message,
		Timestamp:    result.Timestamp,
	}

	var logs []Log

	if data, err := os.ReadFile(l.filePath); err == nil {
		if err := json.Unmarshal(data, &logs); err != nil {
			return fmt.Errorf("failed to unmarshal existing logs: %w", err)
		}
	}

	logs = append(logs, logEntry)

	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}

	if err := os.WriteFile(l.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write logs: %w", err)
	}

	return nil
}

func generateID() string {
	return fmt.Sprintf("TRX-%d", time.Now().UnixNano())
}
