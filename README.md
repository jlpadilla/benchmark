# benchmark
Benchmark database technologies.

## Using this project
1. **REQUIRED:** Setup the target database. Instructions below.
2. `go run main.go [numRecords]`

## Setup the target database locally
Currently only [PostgreSQL](https://www.postgresql.org/), will add others soon.

### Redisgraph
```
docker run -p 6379:6379 -it --rm redislabs/redisgraph
```

### PostgresSQL
https://www.postgresql.org/

1. Start PostgreQSL in a docker container
```
mkdir -p ${HOME}/postgres/data
docker run -d \
 --name dev-postgres \
 -e POSTGRES_PASSWORD=dev-pass! \
 -v ${HOME}/postgres/data/:/var/lib/postgresql/data \
 -p 5432:5432 \
 postgres
```
2. Login to the container
```
$ docker exec -it dev-postgres bash
# psql -h localhost -U postgres
```
3. Create a database named benchmark
```
postgres=# CREATE DATABASE benchmark;
```
4. Use database
```
\c benchmark
```