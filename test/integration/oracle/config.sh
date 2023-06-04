#!/bin/bash

container_name='db-unit-extractor'
oracle_port='1521'
emde_port='5500'
oracle_password='admin'

# Access Enterprise Manager Database Express at localhost:$emde_port
# user: system
# password: $oracle_password
#
# No data volume is created because there is no need to keep test data. If you want it,
# add this parameter `-v $HOME/oradata:/opt/oracle/oradata`.
#
# Connect to Oracle via sqlplus within the container: docker exec -it $container_name sqlplus / as sysdba

docker run -d --rm --name $container_name \
  -p $oracle_port:1521 -p $emde_port:5500 \
  -e ORACLE_PWD=$oracle_password \
  -v ./scripts/startup:/opt/oracle/scripts/startup \
  container-registry.oracle.com/database/express:21.3.0-xe
