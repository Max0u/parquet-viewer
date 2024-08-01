package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/tool/parquet-tools/schematool"
)

func (m *mainModel) loadParquetFile() error {
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

	tree := schematool.CreateSchemaTree(pr.SchemaHandler.SchemaElements)
	log.Printf("%s\n", tree.OutputJsonSchema())

	// Get column names from the Parquet file schema
	// columns := []table.Column{}
	// for _, schemaElem := range pr.SchemaHandler.SchemaElements {
	// 	if schemaElem.GetNumChildren() == 0 { // leaf nodes represent actual columns
	// 		column := table.Column{
	// 			Title: schemaElem.Name,
	// 			Width: len(schemaElem.Name) + 5, // adding a bit of padding for display
	// 		}
	// 		log.Printf(schemaElem.Name)
	// 		columns = append(columns, column)
	// 	}
	// }
	// m.tableView.SetColumns(columns)

	res, err := pr.ReadByNumber(1)
	if err != nil {
		log.Println("Can't read", err)
	}

	// rowsData, err := json.Marshal(res)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Can't convert to json: %s\n", err)
	// 	os.Exit(1)
	// }
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
	// if len(rowsData) > 0 {
	// 	columns := []table.Column{}
	// 	for col := range rowsData[0] {
	// 		columns = append(columns, table.Column{
	// 			Title: col,
	// 			Width: len(col) + 5,
	// 		})
	// 	}
	// 	m.tableView.SetColumns(columns)

	// 	// Fill table rows
	// 	rows := []table.Row{}
	// 	for _, rowData := range rowsData {
	// 		row := make(table.Row, len(columns))
	// 		for i, col := range columns {
	// 			row[i] = fmt.Sprintf("%v", rowData[col.Title])
	// 		}
	// 		rows = append(rows, row)
	// 	}
	// 	m.tableView.SetRows(rows)
	// }
	return nil
}
