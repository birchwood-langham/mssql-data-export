package dataexport

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// EncryptedColumnLibrary is a library of tables that will require encrypting and the columns in each table that should be encrypted when data is exported
type EncryptedColumnLibrary struct {
	library map[string][]string
}

// Parse will read the encrypted columns configuration file and create the library
func (e *EncryptedColumnLibrary) Parse(libraryFile string) (int, error) {
	if libraryFile == "" {
		return 0, errors.New("Encrypted Column Library file has not been provided, cannot parse library")
	}

	count := 0

	file, err := os.Open(libraryFile)

	if err != nil {
		return 0, err
	}

	// make sure we close the file once we have finished with it at the end of the method
	defer file.Close()

	e.library = make(map[string][]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ";")
		table, column := strings.ToLower(strings.TrimSpace(data[0])), strings.ToLower(strings.TrimSpace(data[1]))

		l, ok := e.library[table]

		if !ok { // we don't have this table in the library so create it and add the column data for this table
			e.library[table] = []string{column}
		} else { // we are already encrypting other columns for this table, so add this column to the list of encrypted columns for the table
			e.library[table] = append(l, column)
		}

		count++
	}

	return count, nil
}

// Exists checks the library to see if the table and column specified should be encrypted
func (e *EncryptedColumnLibrary) Exists(table string, column string) (bool, error) {
	_table := strings.ToLower(strings.TrimSpace(table))
	_column := strings.ToLower(strings.TrimSpace(column))

	if e.library == nil {
		return false, errors.New("The library has not been created, specify a library file and Parse it before using it.")
	}

	cols, ok := e.library[_table]
	if !ok {
		return false, nil
	}

	for _, v := range cols {
		if v == _column {
			return true, nil
		}
	}

	return false, nil
}
