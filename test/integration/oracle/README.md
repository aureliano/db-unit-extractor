# Oracle integration test

Execute integration test against an Oracle database.

## Environment

First of all, you must prepare an Oracle database with some data. To do that, a container is started with some tables and data. This environment is set to launch an [Oracle Linux container with Oracle XE 21c](https://container-registry.oracle.com). Database is fed with the Human Resources schema, available at [oracle-samples](https://github.com/oracle-samples/db-sample-schemas) repository.

### Setup and launch

In order to setup and launch Oracle, execute the [config.sh](config.sh) shell script. Oracle will listen at port 1521 and the Enterprise Manager Database Express can be accessed by https://localhost:5500 - user system and password admin.

## Test cases - Human Resources

Bellow there is an example of a test case execution. Although you better execute [execute-test-cases.sh](./execute-test-cases.sh). which has all test scenarios.

### Formatted employees data-set

```sh
dist/db-unit-extractor_linux_amd64_v1/db-unit-extractor extract \
    -n oracle://sys:admin@localhost:1521/xe \
    -s test/integration/oracle/models/employees-ds-model.yml \
    -t xml -t sql -f \
    -r employee_id=200 \
    -d /tmp/db-unit-extractor/integration-tests/oracle
```

#### XML data-set

```sh
go run test/integration/*.go \
    test/integration/oracle/expectations/employees-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/employees-ds-model.xml
```

#### SQL data-set

```sh
go run test/integration/*.go \
    test/integration/oracle/expectations/employees-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/employees-ds-model.sql
```

### Automation

```sh
# Setup Oracle XE 21c.
./config.sh

# Execute integration tests.
./execute-test-cases.sh
```
