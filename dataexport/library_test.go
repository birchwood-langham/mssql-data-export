package dataexport

import "testing"

func TestLibrary(t *testing.T) {
	test_library_file := "library_test.lib"
	library := new(EncryptedColumnLibrary)

	// first test the Parse method to read the file
	count, err := library.Parse(test_library_file)

	if err != nil {
		t.Fatalf("Could not parse test library file: %s", test_library_file)
	}

	if count != 2 {
		t.Logf("Parse library file failed, expected = %d, parse count = %d", 2, count)
		t.Fail()
	}

	// now test exists
	table_name := "testdb.dbo.test_table"
	test_column1 := "firstname"
	test_column2 := "lastname"
	test_column3 := "Address1"

	found, err := library.Exists(table_name, test_column1)
	if err != nil {
		t.Fatalf("Could not perform lookup table %s, column %s from library file: %s", table_name, test_column1, test_library_file)
	}

	if !found {
		t.Logf("EncryptedColumnLibrary could not find table: %s, column: %s", table_name, test_column1)
		t.Fail()
	}

	found, err = library.Exists(table_name, test_column2)
	if err != nil {
		t.Fatalf("Could not perform lookup table %s, column %s from library file: %s", table_name, test_column2, test_library_file)
	}

	if !found {
		t.Logf("EncryptedColumnLibrary could not find table: %s, column: %s", table_name, test_column2)
		t.Fail()
	}

	found, err = library.Exists(table_name, test_column3)
	if err != nil {
		t.Fatalf("Could not perform lookup table %s, column %s from library file: %s", table_name, test_column3, test_library_file)
	}

	if found {
		t.Logf("EncryptedColumnLibrary found something that shouldn't be there - table: %s, column: %s", table_name, test_column3)
		t.Fail()
	}
}
