# benchmark
Benchmark database technologies.


### Setup PostgresSQL

mkdir -p ${HOME}/postgres/data

docker run -d \
 --name dev-postgres \
 -e POSTGRES_PASSWORD=dev-pass! \
 -v ${HOME}/postgres/data/:/var/lib/postgresql/data \
 -p 5432:5432 \
 postgres

Login to container
```
$ docker exec -it dev-postgres bash
# psql -h localhost -U postgres
```
Create database
```
postgres=# CREATE DATABASE benchmark;
```