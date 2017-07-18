package main

import (
	"bufio"
	"database/sql"
	"flag"
	"log"
	"os"

	"birchwoodlangham.com/mssql-data-export/dataexport"
	"strings"
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
	log.Println("Starting data export")

	log.Printf("Opening tables file: %s\n", c.TablesFile)
	tables, err := os.Open(c.TablesFile)
	if err != nil {
		return err
	}
	defer tables.Close()

	log.Println("Creating Encrypted Column Library...")
	library := dataexport.EncryptedColumnLibrary{}

	if c.EncryptionLibrary != "" {
		log.Printf("Parsing encryption library file: %s\n", c.EncryptionLibrary)
		_, err = library.Parse(c.EncryptionLibrary)

		if err != nil {
			return err
		}
	}

	// open the Database
	log.Printf("Opening connection to SQL database server: %s, user: %s\n", c.Host(), c.User())

	db, err := sql.Open("mssql", c.GetConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	log.Printf("Creating data exporter\n")
	exporter := dataexport.Exporter{
		Db:        db,
		Separator: ";",
		OutputDir: c.OutputPath,
		Library:   library,
		Secret:    c.Secret,
	}

	scanner := bufio.NewScanner(tables)
	var count int64

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		// in case there's no ; to split the table with the filter condition, let's just add it so the
		// application can continue without raising an error
		if !strings.Contains(line, ";") {
			line += ";"
		}

		lineData := strings.Split(line, ";")

		table, filter := "", ""
		if len(lineData) > 1 {
			filter = strings.ToLower(strings.TrimSpace(lineData[1]))
		}

		table = strings.ToLower(strings.TrimSpace(lineData[0]))

		log.Printf("Exporting table %s\n", table)

		if c.OutputType == "CSV" {
			count, err = exporter.ExportCsv(table, filter)
		} else {
			count, err = exporter.ExportSQL(table, filter)
		}

		if err != nil {
			return err
		}

		log.Printf("Exported table: %s, %d rows exported", table, count)

	}

	return nil
}
