package connectors

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type CsvConnector struct {
	File string
}

func newCsvConnector(filename string, append bool, headers *[]string) (*CsvConnector, error) {
	mode := os.O_CREATE | os.O_WRONLY
	if append {
		mode |= os.O_APPEND
	} else {
		mode |= os.O_TRUNC
	}

	file, err := os.OpenFile(filename, mode, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	if headers != nil {
		headerLine := fmt.Sprintln(strings.Join(*headers, ","))
		if _, err := file.Write([]byte(headerLine)); err != nil {
			return nil, fmt.Errorf("could not write headers: %w", err)
		}
	}

	return &CsvConnector{
		File: filename,
	}, nil
}

func NewCsvConnector(filename string, append bool) (*CsvConnector, error) {
	return newCsvConnector(filename, append, nil)
}

func NewCsvConnectorWithHeaders(filename string, append bool, headers []string) (*CsvConnector, error) {
	return newCsvConnector(filename, append, &headers)
}

func (c *CsvConnector) Send(data []string) error {
	file, err := os.OpenFile(c.File, os.O_WRONLY|os.O_APPEND, 0644)
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
