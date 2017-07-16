package main

import (
	"bufio"
	"database/sql"
	"flag"
	"log"
	"os"

	"birchwoodlangham.com/mssql-data-export/dataexport"
)

func main() {

	c := dataexport.Config{}
	c.Setup()

	flag.Parse()
	_, err := c.Validate()

	if err != nil {
		log.Fatalf("%s\n", err)
	}

	err = run(c)

	if err != nil {
		log.Fatalf("%s\n", err)
	}
}

func run(c dataexport.Config) error {
	tables, err := os.Open(c.TablesFile)
	if err != nil {
		return err
	}

	library := dataexport.EncryptedColumnLibrary{}

	_, err = library.Parse(c.EncryptionLibrary)

	if err != nil {
		return err
	}

	// open the Database
	db, err := sql.Open("mssql", c.GetConnectionString())

	if err != nil {
		return err
	}

	exporter := dataexport.DataExporter{
		Db:        db,
		Separator: dataexport.Char([]byte(";")[0]),
		OutputDir: c.OutputPath,
		Library:   library,
		Secret:    "test",
	}

	scanner := bufio.NewScanner(tables)
	for scanner.Scan() {
		table := scanner.Text()

		log.Printf("Exporting table %s", table)
		if c.OutputType == "CSV" {
			exporter.ExportCsv(table)
		} else {
			exporter.ExportSQL(table)
		}
	}

	return nil
}
