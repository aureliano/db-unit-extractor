#!/bin/bash

base_dir="$(cd "$(dirname "${0}")/../../.." ; pwd -P)"
cd "${base_dir}"

echo "db-unit-extractor integration tests"
echo

echo "Oracle database"
echo

echo "Scenario: Unformatted employees data-set"
dist/db-unit-extractor_linux_amd64_v1/db-unit-extractor extract \
    -n oracle://sys:admin@localhost:1521/xe \
    -s test/integration/oracle/models/employees-ds-model.yml \
    -t xml -t sql \
    -r employee_id=200 \
    -d /tmp/db-unit-extractor/integration-tests/oracle

echo "  Check XML data-set"
go run test/integration/*.go \
    test/integration/oracle/expectations/employees-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/employees-ds-model.xml

if [ $? -ne 0 ];
then
    echo
    echo " [ERROR] Check XML data-set failed."
    exit 1
fi

echo "  Check SQL data-set"
go run test/integration/*.go \
    test/integration/oracle/expectations/employees-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/employees-ds-model.sql

if [ $? -ne 0 ];
then
    echo
    echo " [ERROR] Check SQL data-set failed."
    exit 1
fi

echo

echo "Scenario: Formatted employees data-set"

dist/db-unit-extractor_linux_amd64_v1/db-unit-extractor extract \
    -n oracle://sys:admin@localhost:1521/xe \
    -s test/integration/oracle/models/employees-ds-model.yml \
    -t xml -t sql -f \
    -r employee_id=200 \
    -d /tmp/db-unit-extractor/integration-tests/oracle

echo "  Check XML data-set"
go run test/integration/*.go \
    test/integration/oracle/expectations/employees-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/employees-ds-model.xml

if [ $? -ne 0 ];
then
    echo
    echo " [ERROR] Check XML data-set failed."
    exit 1
fi

echo "  Check SQL data-set"
go run test/integration/*.go \
    test/integration/oracle/expectations/employees-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/employees-ds-model.sql

if [ $? -ne 0 ];
then
    echo
    echo " [ERROR] Check SQL data-set failed."
    exit 1
fi

echo

echo "Scenario: Unformatted departments data-set"
dist/db-unit-extractor_linux_amd64_v1/db-unit-extractor extract \
    -n oracle://sys:admin@localhost:1521/xe \
    -s test/integration/oracle/models/departments-ds-model.yml \
    -t xml -t sql \
    -r department_id=90 \
    -d /tmp/db-unit-extractor/integration-tests/oracle

echo "  Check XML data-set"
go run test/integration/*.go \
    test/integration/oracle/expectations/departments-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/departments-ds-model.xml

if [ $? -ne 0 ];
then
    echo
    echo " [ERROR] Check XML data-set failed."
    exit 1
fi

echo "  Check SQL data-set"
go run test/integration/*.go \
    test/integration/oracle/expectations/departments-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/departments-ds-model.sql

if [ $? -ne 0 ];
then
    echo
    echo " [ERROR] Check SQL data-set failed."
    exit 1
fi

echo

echo "Scenario: Formatted departments data-set"
dist/db-unit-extractor_linux_amd64_v1/db-unit-extractor extract \
    -n oracle://sys:admin@localhost:1521/xe \
    -s test/integration/oracle/models/departments-ds-model.yml \
    -t xml -t sql -f \
    -r department_id=90 \
    -d /tmp/db-unit-extractor/integration-tests/oracle

echo "  Check XML data-set"
go run test/integration/*.go \
    test/integration/oracle/expectations/departments-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/departments-ds-model.xml

if [ $? -ne 0 ];
then
    echo
    echo " [ERROR] Check XML data-set failed."
    exit 1
fi

echo "  Check SQL data-set"
go run test/integration/*.go \
    test/integration/oracle/expectations/departments-ds-expectation.yml \
    /tmp/db-unit-extractor/integration-tests/oracle/departments-ds-model.sql

if [ $? -ne 0 ];
then
    echo
    echo " [ERROR] Check SQL data-set failed."
    exit 1
fi
