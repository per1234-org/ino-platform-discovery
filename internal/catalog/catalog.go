// Package catalog contains code related to the inoplatforms catalog.
package catalog

import (
	"bytes"
	"encoding/csv"

	"github.com/arduino/go-paths-helper"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogentry"
)

// Type is the type of the catalog data.
type Type [][]string

// Load returns a Type object populated with the data from the inoplatforms catalog file.
func Load(path *paths.Path) (Type, error) {
	raw, err := path.ReadFile()
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(bytes.NewReader(raw))
	csvReader.Comma = '\t'
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return csvData[1:], nil
}

func (catalog Type) Write(path *paths.Path) error {
	// Add the spreadsheet heading row to the provided data.
	headingRow := catalogentry.New()
	for column := range catalogcolumn.EnumEnd {
		headingRow[column] = column.String()
	}
	spreadsheet := append([][]string{headingRow}, catalog...)

	var fileContent bytes.Buffer
	writer := csv.NewWriter(&fileContent)
	writer.Comma = '\t'
	if err := writer.WriteAll(spreadsheet); err != nil {
		return err
	}
	if writer.Error() != nil {
		return writer.Error()
	}

	err := path.WriteFile(fileContent.Bytes())
	if err != nil {
		return err
	}

	return nil
}
