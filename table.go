package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func (m *mainModel) loadParquetFile() error {
	tableModel, ok := m.tableWindow.model.(table.Model)
	if !ok {
		// This should never happen. At this state, prefer assuming that this should never happen instead of
		//handling this case later.
		panic("should not happen")
	}

	fr, err := local.NewLocalFileReader(m.selectedFile)
	if err != nil {
		return err
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, nil, 1)
	if err != nil {
		return err
	}
	defer pr.ReadStop()

	res, err := pr.ReadByNumber(5) // Read only the first 5 rows
	if err != nil {
		log.Println("Can't read", err)
	}

	jsonData, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't convert to json: %s\n", err)
		os.Exit(1)
	}

	// Convert JSON data to a slice of maps
	var rowsData []map[string]interface{}
	err = json.Unmarshal(jsonData, &rowsData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't unmarshal json: %s\n", err)
		os.Exit(1)
	}

	// Determine columns from the first row (if available)
	if len(rowsData) > 0 {
		columns := []table.Column{}
		colCount := 0
		for col := range rowsData[0] {
			columns = append(columns, table.Column{
				Title: col,
				Width: len(col) + 5,
			})
			colCount++
			if colCount >= 3 {
				break // Keep only the first 5 columns
			}
		}
		tableModel.SetColumns(columns)

		log.Printf("%v", columns)

		// Fill table rows
		rows := []table.Row{}
		for _, rowData := range rowsData {
			row := make(table.Row, len(columns))
			for i, col := range columns {
				row[i] = fmt.Sprintf("%v", rowData[col.Title])
			}
			rows = append(rows, row)
		}
		tableModel.SetRows(rows)

		log.Printf("%v", rows)

	}
	m.tableWindow.model = tableModel
	return nil
}
