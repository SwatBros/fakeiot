package connectors

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CsvConnector struct {
	File string
}

func (c *CsvConnector) Send(data []string) error {
	file, err := os.OpenFile(c.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(data); err != nil {
		return fmt.Errorf("could not write to CSV: %w", err)
	}

	return nil
}
