# db-unit-extractor

**db-unit-extractor** is a database extractor for unit testing. Rather than a massive data extractor, its goal is to extract a data-set for unit testing. 

Sometimes, it is necessary to write some integration tests in order to assure that components may work together. A database is a middleware that is often mocked in unit test and tests are neglited in favor of black box tests. Some frameworks like [dbunit](https://www.dbunit.org/) and [h2](https://www.h2database.com) helps creating those integration tests accessing a database. Although, they lack a consistent tool for generating a data-set to each test scenario. This tool was made to support testers on creating data-sets by extracting data of a specific set of records.

## Command line application

## Dataset schema

## Database reader

## File writer

## Contributing
Please feel free to submit issues, fork the repository and send pull requests! But first, read [this guide](./CONTRIBUTING.md) in order to get orientations on how to contribute the best way.

## License
This project is licensed under the terms of the MIT license found in the [LICENSE](./LICENSE) file.
