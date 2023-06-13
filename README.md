# db-unit-extractor

**db-unit-extractor** is a database extractor for unit testing. Rather than a massive data extractor, its goal is to extract a data-set for unit testing. 

Sometimes, it is necessary to write some integration tests in order to assure that components may work together. A database is a middleware that is often mocked in unit test and tests are neglited in favor of black box tests. Some frameworks like [dbunit](https://www.dbunit.org/) and [h2](https://www.h2database.com) helps creating those integration tests accessing a database. Although, they lack a consistent tool for generating a data-set to each test scenario. This tool was made to support testers on creating data-sets by extracting data of a specific set of records.

Supported [operating systems](https://en.wikipedia.org/wiki/Operating_system) are [Linux](https://en.wikipedia.org/wiki/Linux), [Darwin](https://en.wikipedia.org/wiki/Darwin_(operating_system)) and [Windows](https://en.wikipedia.org/wiki/Microsoft_Windows).

## Index

 1. [Data-set schema](#data-set-schema)
    1. [Converters](#converters)
        1. [Date/Time ISO 8601](#datetime-iso-8601)
        2. [BLOB](#blob)
    2. [Tables](#tables)
        1. [Filtering](#filtering)
        2. [Fetch columns](#fetch-columns)
        3. [Ignore columns](#ignore-columns)
        4. [Dynamic filter - command line parameter](#dynamic-filter---command-line-parameter)
        5. [Dynamic filter - referenced table](#dynamic-filter---referenced-table)
        6. [Dynamic filter - multivalued referenced table](#dynamic-filter---multivalued-referenced-table)
 2. [Database reader](#database-reader)
    1. [Oracle](#oracle)
 3. [File writer](#file-writer)
    1. [Console](#console)
    2. [XML](#xml)
 4. [Command line application](#command-line-application)
 5. [Update program](#update-program)
 6. [Development](#development)
    1. [Integration tests](#integration-tests)
    2. [Profiling](#profiling)
    3. [Release](#release)
 7. [Contributing](#contributing)
 8. [License](#license)

## Data-set schema

A data-set schema is a yaml file with instructions of how a data-set will be created, as which tables and which columns will be part of the query, columns that will filter the query and data converters.

### Converters

A converter takes a column and formats its output in order to be suitable for consumers like dbunit. There two converters that are loaded by default and it's not necessary to declare them in the schema file.

Converters are defined by the key `converters` and handles an array of converter ids.

```yaml
---
converters:
  - conv-id-1
  - conv-id-2
```

#### Date/Time ISO 8601

Converts a date/time/timestamp data in a formatted [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) text.

```yaml
---
converters:
  - date-time-iso8601
```

#### BLOB

Converts a blob data in a [base64](https://en.wikipedia.org/wiki/Base64) encoded text.

```yaml
---
converters:
  - blob
```

### Tables

The tables key is the entry point to data-set definition. Indeed, it's where you tell which tables will be queried in order to build the data-set.

Tables are defined by the key `tables` and handles an array of table keys.

```yaml
tables:
  - name: table_name
```

In the sample above, an sql like `select * from table_name` will be sent to database and all records and all columns of table_name will be fetched.

#### Filtering

Most of times you'll probably want filter a query. In the case you want just a specific record, you pass a filter with name and value. so, to fetch only the record with id 12345 you do like bellow.

```yaml
tables:
  - name: table_name
    filters:
      - name: id
        value: 12345
```

It'll query `select * from table_name where id = '12345'`.

#### Fetch columns

If you don't want to fetch all columns and just select a few of them, you may set which ones.

```yaml
tables:
  - name: table_name
    filters:
      - name: id
        value: 12345
    columns:
      - id
      - column_2
      - column_5
      - column_9
```

It'll query `select id, column_2, column_5, column_9 from table_name where id = '12345'`.

#### Ignore columns

Sometimes, instead of setting which columns you want to get in a query, you may want to set which ones you don't. Following the samples above, imagine `table_name` has this set of columns: id, column_1, column_2, column_3, column_4, column_5, column_6, column_7, column_8, column_9. And you want to get all but columns column_3 and column_4.

```yaml
tables:
  - name: table_name
    filters:
      - name: id
        value: 12345
    ignore:
      - column_3
      - column_4
```

**Important to note that columns and ignore are excludent!** You cannot set both in a table.

#### Dynamic filter - command line parameter

So far, we've seen static references on filtering. A better approach would be using dynamic filters. Imagine you have a lot of scenarios to the same data-set. Instead of creating many schema files you just need to parameterize the filter.

```yaml
tables:
  - name: table_name
    filters:
      - name: id
        value: ${table_name_pk}
```

Now, since filter named `id` has a dynamic value, you must set `table_name_pk` parameter to command line.

#### Dynamic filter - referenced table

Unlike the previous section, you may have situations where you don't know which value a parameter may have because you depend of a resulted query. In the previous samples, we had just one table. So, imagine we have two tables: customers and orders. And we get the customer id from command line and the orders of the employee from que queried customer.

```yaml
tables:
  - name: customers
    filters:
      - name: id
        value: ${customer_id}
  - name: orders
    filters:
      - name: id
        value: ${customers.id}
```

Above, we see that, a query to customers will be made and the result will be used to query orders. So, a customer with `id` comming from command line and an order with `customer_id` comming from the result of querying table customers. The principle is referenced table name followed by a dot and column name: `${table_name.column_name}`

**Importnat!** The order of tables doesn't matter. They are ordered at runtime. Setting orders before customers will make no difference.

#### Dynamic filter - multivalued referenced table

Multivalued filters are useful when you need to reference a table that returns more than one record. Following the last sample, now we need to fetch all products ordered by the customer. As we known that a customer may have ordered more than once, it is necessary to filter products by an array of order ids.

```yaml
tables:
  - name: customers
    filters:
      - name: id
        value: ${customer_id}
  - name: orders
    filters:
      - name: id
        value: ${customers.id}
  - name: orders_products
    filters:
      - name: order_id
        value: ${orders.id[@]}
  - name: products
    filters:
      - name: id
        value: ${orders.product_id[@]}
```

Above you might have noticed the use of `[@]` suffix, that means: this reference is multivalued. At the end, our data-set will have a customer with many orders with many products.

## Database reader

Database reader is a component that handles data recovering. In the next subsections you'll see what database systems are supported by this project.

### Oracle

This reader handles data and metadata recovering from [Oracle](https://www.oracle.com) databases. Integrations tests were made in Oracle 21c. Although it is supposed to work on older versions.

## File writer

File writer is a component that handles data writing. It takes records from data reader and outputs to an arbitrary file type. In the next subsections you'll see what file writers are supported by this project.

### Console

This writer sends records to the standard output.

### XML

This writer sends records to an XML file.

## Command line application

Data-set extractions are made through a command line application named `db-unit-extractor`.

```
Database extractor for unit testing.
Go to https://github.com/aureliano/db-unit-extractor/issues in order to report a bug or make any suggestion.

Usage:
  db-unit-extractor [flags]
  db-unit-extractor [command]

Available Commands:
  extract     Extract data-set from database
  help        Help about any command
  update      Update this program

Flags:
  -h, --help      help for db-unit-extractor
  -v, --version   Print db-unit-extractor version

Use "db-unit-extractor [command] --help" for more information about a command.
```

In order to see some samples on how to extract data execute this command `db-unit-extractor extract --help`

```
Extract data-set from a database to any supported file.

Usage:
  db-unit-extractor extract [flags]

Examples:
  # Extract data-set from Oracle and write to the console.
  db-unit-extractor extract -s /path/to/schema.yml -n oracle://usr:pwd@127.0.0.1:1521/test

  # Pass parameter expected in schema file.
  db-unit-extractor extract -s /path/to/schema.yml -n oracle://usr:pwd@127.0.0.1:1521/test -r customer_id=4329

  # Write to xml file too.
  db-unit-extractor extract -s /path/to/schema.yml -n oracle://usr:pwd@127.0.0.1:1521/test -r customer_id=4329 -t xml

  # Format xml output.
  db-unit-extractor extract -s /path/to/schema.yml -n oracle://usr:pwd@127.0.0.1:1521/test -r customer_id=4329 -t xml -f

Flags:
  -n, --data-source-name string   Data source name (aka connection string: <driver>://<username>:<password>@<host>:<port>/<database>).
  -d, --directory string          Output directory. (default ".")
  -f, --formatted-output          Whether the output should be formatted.
  -h, --help                      help for extract
      --max-idle-conn int         Set the maximum number of concurrently idle connections (default 2)
      --max-open-conn int         Set the maximum number of concurrently open connections (default 3)
  -t, --output-type stringArray   Extracted data output format type. Expected: [console xml] (default [console])
  -r, --references stringArray    Expected input parameter in 'schema' file. Expected: name=value
  -s, --schema string             Path to the file with the data schema to be extracted.
```

## Update program

If you wanna stay up to date, you may call `db-unit-extractor update` and a new verion - if it is not the edge - will be installed.

## Development

### Integration tests

Go to the integration tests [documentation](./test/integration/README.md) to get detailed information and how to execute them.

### Profiling

Identify performance problems with code profilers. You may enable CPU profiling and memory profiling by exporting CPU_PROFILE and MEM_PROFILE environment variables. They must point to a file (non existing directories won't be created). In the sample bellow, both profilers are enabled.

```shell
export CPU_PROFILE=/tmp/db-unit-extractor/cpu.prof
export MEM_PROFILE=/tmp/db-unit-extractor/mem.prof

# Builb program.
make clean snapshot

# Launch program with profilers enabled.
dist/db-unit-extractor_linux_amd64_v1/db-unit-extractor \
    extract \
      -s /path/to/schema.yml \
      -n oracle://usr:pwd@127.0.0.1:1521/test \
      -r customer_id=4329 \
      -t xml \
      -f

# CPU profiling.
go tool pprof dist/db-unit-extractor_linux_amd64_v1/db-unit-extractor /tmp/db-unit-extractor/cpu.prof

# Memory profiling.
go tool pprof dist/db-unit-extractor_linux_amd64_v1/db-unit-extractor /tmp/db-unit-extractor/mem.prof
```

### Release

Programs are released under semantic versioning - [semver](https://semver.org).

 > Given a version number MAJOR.MINOR.PATCH, increment the:
 > 1. MAJOR version when you make incompatible API changes
 > 2. MINOR version when you add functionality in a backward compatible manner
 > 3. PATCH version when you make backward compatible bug fixes

Before you make sure that continuous integration pipeline isn't broken. The pipeline execute unit tests and code linters to check code compliance.

Beyond, execute [integration tests](./test/integration/README.md) and make sure none of them is broken.

When all branches and pull requests of bug fixes and new features are merged, just create a tag following semantic versioning. Bellow is a sample of the first version tag creation.

```shell
# Create an annotated tag.
git tag -a v1.0.0 -m "Database extractor for unit testing. Supports Oracle database reader and console and xml file writers."

# Publish tag.
git push origin v1.0.0
```

After a tag is published a [pipeline](./.github/workflows/release.yml) is triggered and create a release based on the created tag. Releases are listed [here](https://github.com/aureliano/db-unit-extractor/releases).

## Contributing

Please feel free to submit issues, fork the repository and send pull requests! But first, read [this guide](./CONTRIBUTING.md) in order to get orientations on how to contribute the best way.

## License

This project is licensed under the terms of the MIT license found in the [LICENSE](./LICENSE) file.
