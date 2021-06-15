# Benchmark Results
Document benhmark results using this project.

## Benchmark database operations without any relationships.

Operation              | Redisgraph | PostgreSQL
---                    | ---        | ---
Insert 100k            | 2.7s       | 11s
Insert 500k            | 14s        | 1m10s
Insert 1M              | 29s        | 3m
Insert 2M              |
With 100k, edit 1k     |
With 500k, edit 1k     |
With 1M, edit 1k       |
With 100k, delete 1k   |
With 500k, delete 1k   |
With 1M, delete 1k     |


## Benchmark database operations including relationships

TBD

Operation              | Redisgraph | Postgre
---                    | ---        | ---