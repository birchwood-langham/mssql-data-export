package dataexport

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// this library is required only for the DataExporter to export MSSQL data
	_ "github.com/denisenkom/go-mssqldb"
)

// DataExporter is a library for export SQL server data to text files and encrypts specified fields to anonymise data
type DataExporter struct {
	Db         *sql.DB
	Separator  string
	OutputDir  string
	Library    EncryptedColumnLibrary
	columnData []interface{}
	Secret     string
}

// ExportCsv queries the data in the table specified and writes the data to output directory in a CSV format
func (e DataExporter) ExportCsv(table string) (int64, error) {
	result, err := e.Db.Query("select * from $1;", table)
	defer result.Close()

	if err != nil {
		return 0, err
	}

	rows := int64(0)
	columns, err := result.Columns()

	if err != nil {
		return rows, err
	}

	header := columnHeaders(columns, e.Separator)
	e.initializeColumns(columns)

	outputFile, err := os.Open(e.OutputDir + "/" + table + ".csv")
	defer outputFile.Close()

	if err != nil {
		return rows, err
	}

	_, err = outputFile.WriteString(header)

	if err != nil {
		return rows, err
	}

	for result.Next() {
		err = result.Scan(e.columnData...)

		if err != nil {
			return rows, err
		}

		outputFile.WriteString(e.createCsvOutputString(table, columns))

		rows++
	}

	return rows, nil
}

func (e *DataExporter) createCsvOutputString(table string, columns []string) string {
	output := ""
	columnCount := len(e.columnData)

	for i := 0; i < columnCount; i++ {
		encrypt, err := e.Library.Exists(table, columns[i])

		if err != nil {
			return ""
		}

		switch d := (e.columnData[i]).(type) {
		case nil:
			output += "null"
		case bool:
			output += "true"
		case time.Time:
			output += d.Format("'2000-01-01 00:00:00.000'")
		case []byte:
			if encrypt {
				output += fmt.Sprintf("'%s'", Encrypt(string(d), e.Secret))
			} else {
				output += fmt.Sprintf("'%s'", string(d))
			}
		default:
			if encrypt {
				output += Encrypt(fmt.Sprintf("%v", d), e.Secret)
			} else {
				output += fmt.Sprintf("%v", d)
			}
		}

		if i < columnCount-1 {
			output += e.Separator
		}
	}

	return fmt.Sprintf("%s\n", output)
}

func (e *DataExporter) createSQLOutputString(table string, columns []string) string {
	columnValues := ""
	columnCount := len(e.columnData)

	columnNames := columnHeaders(columns, e.Separator)

	for i := 0; i < columnCount; i++ {
		encrypt, err := e.Library.Exists(table, columns[i])

		if err != nil {
			return ""
		}

		switch d := (e.columnData[i]).(type) {
		case nil:
			columnValues += "null"
		case bool:
			columnValues += "true"
		case time.Time:
			columnValues += d.Format("'2000-01-01 00:00:00.000'")
		case []byte:
			if encrypt {
				columnValues += fmt.Sprintf("'%s'", Encrypt(string(d), e.Secret))
			} else {
				columnValues += fmt.Sprintf("'%s'", string(d))
			}
		default:
			if encrypt {
				columnValues += Encrypt(fmt.Sprintf("%v", d), e.Secret)
			} else {
				columnValues += fmt.Sprintf("%v", d)
			}

		}

		if i < columnCount-1 {
			columnValues += ","
		}
	}

	return fmt.Sprintf("insert into %s (%s) values (%s)\n", table, columnNames, columnValues)
}

func (e *DataExporter) initializeColumns(columns []string) {

	columnCount := len(columns)

	// lets initialize the column data array and get the column names to create the file header
	e.columnData = make([]interface{}, columnCount)
	for i := 0; i < columnCount; i++ {
		e.columnData[i] = new(interface{})
	}
}

func columnHeaders(columns []string, separator string) string {
	columnCount := len(columns)
	header := ""
	for i := 0; i < columnCount; i++ {
		if i < columnCount-1 {
			header += ","
		}

		header += columns[i]
	}

	return header
}

// ExportSQL queries the data in the table specified and writes the data as insert statements to the output directory specified
func (e *DataExporter) ExportSQL(table string) (int64, error) {
	result, err := e.Db.Query("select * from $1;", table)
	defer result.Close()

	if err != nil {
		return 0, err
	}

	rows := int64(0)

	columns, err := result.Columns()

	if err != nil {
		return rows, err
	}

	e.initializeColumns(columns)

	outputFile, err := os.Open(e.OutputDir + "/" + table + ".sql")
	defer outputFile.Close()

	if err != nil {
		return rows, err
	}

	_, err = outputFile.WriteString(fmt.Sprintf("set identity_insert %s on", table))

	if err != nil {
		return rows, err
	}

	for result.Next() {
		err = result.Scan(e.columnData...)

		if err != nil {
			return rows, err
		}

		_, err = outputFile.WriteString(e.createSQLOutputString(table, columns))

		if err != nil {
			return rows, err
		}

		rows++
	}

	_, err = outputFile.WriteString(fmt.Sprintf("set identity_insert %s off", table))

	if err != nil {
		return rows, err
	}

	return rows, nil
}
