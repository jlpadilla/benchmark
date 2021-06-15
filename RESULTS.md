# Benchmark Results
Document benhmark results using this project.

## Benchmark database operations without any relationships.

Operation                     | Redisgraph | PostgreSQL (1 table) | Postgre (100 tables)
---                           | ---        | ---                  | ---
Insert 100k                   | 2.7s       | 11s                  | 8s
Insert 500k                   | 14s        | 50s                  | 48s
Insert 1M                     | 29s        | 3m                   | 1m40s
Insert 2M                     | 1m34s      | <b>TOO SLOW </b>     |
With 100k, edit 1k            |            | 7s                   |
With 500k, edit 1k            |            | 1m10s                |
With 1M, edit 1k              |            | 1m42s                |
With 100k, delete 1k          | 15s initial</br>4s after initial | 750ms  |             
With 500k, delete 1k          | 1m         | 34s                  |
With 1M, delete 1k            | 2m         | 1m9s                 |
Query using index (100k)      |            | 8ms
Query non-indexed (100k)      |            | 45ms
Query distinct values (100k)  |            | 80ms
Query using index (500k)      |            | 44ms
Query non-indexed (500k)      |            | 629ms
Query distinct values (500k)  |            | 413ms
Query using index (1M)        |            | 42ms
Query non-indexed (1M)        |            | 2.17s
Query distinct values (1M)    |            | 2.50s

### PostgreSQL observations:
- Time to insert data is proportional to the total existing records. Theres only a marginal improvement from dividing the data in multiple tables.
- 

## Benchmark database operations including relationships

TBD

Operation              | Redisgraph | Postgre
---                    | ---        | ---