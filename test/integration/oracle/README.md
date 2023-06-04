# Oracle integration test
Execute integration test against an Oracle database.

## Environment
First of all, you must prepare an Oracle database with some data. To do that, a container is started with some tables and data. This environment is set to launch an [Oracle Linux container with Oracle XE 21c](https://container-registry.oracle.com). Database is fed with the Human Resources schema, available at [oracle-samples](https://github.com/oracle-samples/db-sample-schemas) repository.

## Setup and launch
In order to setup and launch Oracle, execute the [config.sh](config.sh) shell script. Oracle will listen at port 1521 and the Enterprise Manager Database Express can be accessed by https://localhost:5500 - user system and password admin.
