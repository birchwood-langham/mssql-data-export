package dataexport

import (
	"errors"
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
	"os/user"
)

var (
	errInvalidHost = errors.New("Invalid hostname provided")
	//errInvalidDatabase     = errors.New("Invalid database provided")
	errInvalidTablesConfig = errors.New("Tables configuration is required, but has not been provided, see -t flag for help")
	errInvalidOutputType   = errors.New("Output Type must be CSV or SQL")
)

// Config represents application configurations that should be set with commandline flags
type Config struct {
	host              string
	port              int
	catalog           string
	user              string
	password          string
	TablesFile        string
	EncryptionLibrary string
	OutputPath        string
	OutputType        string
	Secret            string
	askSecret         bool
}

func (c *Config) Host() string {
	return c.host
}

func (c *Config) User() string {
	return c.user
}

// Setup initializes the flags the application supports inline with the Config struct
func (c *Config) Setup() error {
	pwd, err := os.Getwd()

	if err != nil {
		return err
	}

	flag.StringVar(&c.host, "H", "", "Required. MS SQL Server host to connect to")
	flag.IntVar(&c.port, "p", 1433, "TCP/IP port to connect to")
	flag.StringVar(&c.catalog, "c", "", "Required: Initial Catalog/Database to connect to")
	flag.StringVar(&c.user, "U", "", "Required if not using Integrated Security; User to connect as")
	flag.StringVar(&c.password, "P", "", "Required if not using Integrated Security; The password to use for authentication")
	flag.StringVar(&c.TablesFile, "t", "", "Required. File listing tables to export, one table per line in <db>.<schema>.<table> format")
	flag.StringVar(&c.EncryptionLibrary, "e", "", "File listing columns to encrypt in each table, one column per line in <db>.<schema>.<table>;<column> format, only include columns that need to be encrypted")
	flag.StringVar(&c.OutputPath, "o", pwd, "Directory to write output to (default: current directory)")
	flag.StringVar(&c.OutputType, "T", "CSV", "Output file type, accepted values are CSV (default) or SQL")
	flag.BoolVar(&c.askSecret, "s", false, "Prompt for secret to use when creating hash")

	return nil
}

// Validate checks to make sure the configuration has valid values for the flags that are set
func (c *Config) Validate() (bool, error) {
	if c.host == "" {
		flag.PrintDefaults()
		return false, errInvalidHost
	}

	//if c.catalog == "" {
	//	flag.PrintDefaults()
	//	return false, errInvalidDatabase
	//}

	if c.user == "" {
		u, err := user.Current()

		if err != nil {
			return false, err
		}

		c.user = u.Username
	}

	if c.password == "" {
		// we should attempt to prompt for a password
		fmt.Printf("Enter database password for user %s: ", c.user)
		pass, _ := gopass.GetPasswdMasked()
		c.password = string(pass)
	}

	if c.TablesFile == "" {
		flag.PrintDefaults()
		return false, errInvalidTablesConfig
	}

	if c.OutputType != "CSV" && c.OutputType != "SQL" {
		flag.PrintDefaults()
		return false, errInvalidOutputType
	}

	if c.askSecret {
		fmt.Printf("Enter salt to use when encrypting data: ")
		pass, _ := gopass.GetPasswdMasked()
		c.Secret = string(pass)
	}

	return true, nil
}

// GetConnectionString returns the appropriate SQL Server connection string from the commandline parameters provided
func (c *Config) GetConnectionString() string {
	return fmt.Sprintf("Server=%s;Port=%d;Database=%s;User Id=%s;Password=%s", c.host, c.port, c.catalog, c.user, c.password)
}
