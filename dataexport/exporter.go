package dataexport

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// this library is required only for the Exporter to export MSSQL data
	"github.com/denisenkom/go-mssqldb"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"strings"
)

// Exporter is a library for export SQL server data to text files and encrypts specified fields to anonymise data
type Exporter struct {
	Db         *sql.DB
	Separator  string
	OutputDir  string
	Library    EncryptedColumnLibrary
	columnData []interface{}
	Secret     string
}

// ExportCsv queries the data in the table specified and writes the data to output directory in a CSV format
func (e *Exporter) ExportCsv(table string, filter string) (int64, error) {
	table = strings.ToLower(strings.TrimSpace(table))
	filter = strings.ToLower(strings.TrimSpace(filter))

	query := fmt.Sprintf("select * from %s", table)

	if filter != "" {
		query += fmt.Sprintf(" where %s", filter)
	}

	result, err := e.Db.Query(query)

	if err != nil {
		return int64(0), err
	}

	// put it after the error check because it the query may not return any rows
	defer result.Close()

	rows := int64(0)
	columns, err := result.Columns()

	if err != nil {
		return rows, err
	}

	header := columnHeaders(columns, e.Separator)
	e.initializeColumns(columns)

	outputFile, err := os.Create(e.OutputDir + string(os.PathSeparator) + table + ".csv")
	if err != nil {
		return rows, err
	}

	// again we defer the close after we have checked for errors as we can't be certain the file has been created until this point.
	defer outputFile.Close()

	_, err = outputFile.WriteString(header + "\n")

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

func (e *Exporter) createCsvOutputString(table string, columns []string) string {
	output := ""
	columnCount := len(e.columnData)

	for i := 0; i < columnCount; i++ {
		encrypt, err := e.Library.Exists(table, columns[i])

		if err != nil {
			return ""
		}

		output += e.formatValue(e.columnData[i].(*interface{}), encrypt, false)

		if i < columnCount-1 {
			output += e.Separator
		}
	}

	return fmt.Sprintf("%s\n", output)
}

func (e *Exporter) createSQLOutputString(table string, columns []string) string {
	columnValues := ""
	columnCount := len(e.columnData)

	columnNames := columnHeaders(columns, ",") // for SQL files we will always use , as a separator regardless of what is provided

	for i := 0; i < columnCount; i++ {
		encrypt, err := e.Library.Exists(table, columns[i])

		if err != nil {
			return ""
		}

		columnValues += e.formatValue(e.columnData[i].(*interface{}), encrypt, true)

		if i < columnCount-1 {
			columnValues += ","
		}
	}

	return fmt.Sprintf("insert into %s (%s) values (%s)\n", table, columnNames, columnValues)
}

func (e *Exporter) formatValue(value *interface{}, encrypt bool, sqlOutput bool) string {
	switch v := (*value).(type) {
	case nil:
		return "null"
	case bool:
		return "true"
	case time.Time:
		if sqlOutput {
			return v.Format("'2006-01-02 15:04:05.000'")
		}

		return v.Format("2006-01-02 15:04:05.000")
	case []byte:
		// first check to see if it's a unique identifier field
		var uid mssql.UniqueIdentifier
		err := uid.Scan(*value)

		// Scan didn't produce an error so the field must be a uniqueidentifier type
		if err == nil {
			// we've just read a unique id so just output the data
			return fmt.Sprintf("%s", uid.String())
		}

		// it's not a uniqueidentifier type field so therefore try converting it to a string
		// and encrypt it if necessary
		if encrypt {
			if sqlOutput {
				return fmt.Sprintf("'%s'", Encrypt(string(v), e.Secret))
			}
			return fmt.Sprintf("%s", Encrypt(string(v), e.Secret))
		} else {
			if sqlOutput {
				// we need to escape any single quotes because we're not encrypting it
				return fmt.Sprintf("'%s'", strings.Replace(string(v), "'", "''", -1))
			}

			return fmt.Sprintf("%s", v)
		}
	case int, int8, int16, int32, int64:
		if encrypt {
			return Encrypt(fmt.Sprintf("%d", v), e.Secret)
		} else {
			return fmt.Sprintf("%d", v)
		}
	case float32, float64:
		if encrypt {
			return Encrypt(fmt.Sprintf("%f", v), e.Secret)
		} else {
			return fmt.Sprintf("%f", v)
		}
	case string:
		if encrypt {
			if sqlOutput {
				return fmt.Sprintf("'%s'", Encrypt(v, e.Secret))
			}
			return fmt.Sprintf("%s", Encrypt(v, e.Secret))
		} else {
			if sqlOutput {
				return fmt.Sprintf("'%s'", strings.Replace(v, "'", "''", -1))
			}

			return fmt.Sprintf("%s", v)
		}
	default:
		if encrypt {
			return Encrypt(fmt.Sprintf("%v", v), e.Secret)
		} else {
			return fmt.Sprintf("%s", v)
		}
	}
}

func (e *Exporter) initializeColumns(columns []string) {

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
		header += columns[i]
		if i < columnCount-1 {
			header += separator
		}
	}

	return header
}

// ExportSQL queries the data in the table specified and writes the data as insert statements to the output directory specified
func (e *Exporter) ExportSQL(table string, filter string) (int64, error) {
	table = strings.ToLower(strings.TrimSpace(table))
	filter = strings.ToLower(strings.TrimSpace(filter))

	query := fmt.Sprintf("select * from %s", table)

	if filter != "" {
		query += fmt.Sprintf(" where %s", filter)
	}

	log.Printf("Executing query: %s\n", query)

	result, err := e.Db.Query(query)

	if err != nil {
		return 0, err
	}

	defer result.Close()

	rows := int64(0)

	columns, err := result.Columns()

	if err != nil {
		return rows, err
	}

	e.initializeColumns(columns)

	outputFile, err := os.Create(e.OutputDir + string(os.PathSeparator) + table + ".sql")
	if err != nil {
		return rows, err
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(fmt.Sprintf("set identity_insert %s on\n", table))

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

	_, err = outputFile.WriteString(fmt.Sprintf("set identity_insert %s off\n", table))

	if err != nil {
		return rows, err
	}

	return rows, nil
}
