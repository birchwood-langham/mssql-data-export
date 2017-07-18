# MSSQL Data Exporter


This application exports data from a MS SQL Server database and allows specific columns within the tables specified to be encrypted on output so that data can be anonymised.
 
## Pre-requisites

To use, you will need to create a text file listing the tables you want to export in the format:

```
    <database>.<schema>.<table1>;<where clause>
    <database>.<schema>.<table2>;
    <database>.<schema>.<table3>;
```

If a where clause is included, this will be added to the query to filter the results, you do not need to include the where keyword, just the conditions that you want to filter with; e.g.

```SQL
Column1 = 'Test' AND Column2 > '2012-01-01'
```

You should create one table per line in the text file. The name is not important, and if you are exporting from one table with the default schema, you could just list each table name you want exporting because it will just substitute the table name into the query.

To specify the columns to encrypt, create another file listing the table and the column within the table to encrypt; separate the table name and the column with a semi-colon ;

```
    <database>.<schema>.<table1>;column1
    <database>.<schema>.<table1>;column2
    <database>.<schema>.<table2>;column3
    <database>.<schema>.<table2>;column7
    <database>.<schema>.<table2>;column9
    <database>.<schema>.<table3>;column2
```

You should create one table/column pair per row.

**You will also need to have read permission on all the database/tables you want to export.**

## Running the application

To export without securing any of the data:
 
```bash
mssql-data-export -H <database server> -t <path/to/tables.file> -o <path/to/output-directory>
```

The application will use your user name and prompt you for your password. If you don't specify an initial catalog, you will be connected to the default database for your SQL login.

If you want to specify the user to use, you can pass the user information with the -U flag, for example:

```bash
mssql-data-export -H <database server> -U sa -t <path/to/tables.file> -o <path/to/output-directory>
```

Or if the database server uses Windows Authentication:

```bash
mssql-data-export -H <database server> -U DOMAIN\Username -t <path/to/tables.file> -o <path/to/output-directory>
```

When prompted for a password, enter your Windows Login password.

If you want to specify the database to connect to, you can use the -c flag

```bash
mssql-data-export -H <database server> -c MYDB -t <path/to/tables.file> -o <path/to/output-directory>
```

To specify a file containing the list of fields to encrypt, use the -e flag

```bash
mssql-data-export -H <database server> -c MYDB -t <path/to/tables.file> -e <path/to/encrypt.file> -o <path/to/output-directory>
```

when the application writes the data to file, if it encounters any field that has been listed in the encrypted columns library file.

If you want to add extra protection, you may also provide a salt to use with the encryption with the -s flag.

```bash
mssql-data-export -H <database server> -c MYDB -t <path/to/tables.file> -e <path/to/encrypt.file> -o <path/to/output-directory> -s
```

This will cause the application to prompt you for a salt to use as part of the encryption.

## Output

By default, the exporter will write the data into a csv file format, using ; as the column separator. 

If you would rather have a SQL script that will generate the appropriate insert statements to insert the data into the tables you have exported from, i.e. in case you want to create test data that can be used, but need to keep production information secret, you can specify the output type with the -T flag.

```bash
mssql-data-export -H <database server> -c MYDB -t <path/to/tables.file> -e <path/to/encrypt.file> -o <path/to/output-directory> -T SQL
```

The -T flag will only accept 2 values CSV and SQL, CSV is the default value.

The application will generate one file per table in the output folder specified. You will need to make sure you have write permission in the folder you have chosen to output the data to.
