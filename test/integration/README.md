# Integration test
Execute integration test against supported database readers and file writers.

## Execution
Integration tests are consisted of executing the program, with a data-set schema model, and checking its output matching to an expected data-set. Suppose we have executed db-unit-extracto and got an XML file named `sample.xml`. Given an expected data-set, we can execute an integration test like this.

```shell
go run *.go \
    ${dbreader}/expectations/expectation.yml \
    sample.xml
```

It's supposed to have one directory to each database reader, in which we can find scripts to setup an environment and data-set model files and expectations files.

### Runner
There is an implementation - at test/integration - that, given an output from db-unit-extractor, executes a test that match the result with an expected data-set. Even not being mandatory building the runner, you may do that and run ti directly. Although, just call `go *.go expected-data-set resulted-data-set`.

## Database readers
Go to the specific directory of a database reader in order to get access to instructions of how to execute integration tests.
