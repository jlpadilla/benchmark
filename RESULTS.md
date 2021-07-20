# Benchmark Results

## Benchmark database operations (excludes relationships)

Operation                     | Redisgraph on laptop | PostgreSQL Docker on laptop        | PostgreSQL on AWS    | SQLite (mem) | SQLite (disk)
---                           | ---                  | ---                                | ---                  | ---          | ---
Insert 100k                   | 2.7s                 | 11s (1 table)</br>8s (100 tables)  | 2.5s                 | 755 ms       | 1.2s
Insert 500k                   | 14s                  | 50s (1 table)</br>48s (100 tables) | 10.3s                | 4.1s         | 8.5s
Insert 1M                     | 29s                  | 3m (1 table)</br>1m40s (100 tables)| 21.3s                | 8.7s         | 17.9s
Insert 2M                     | 1m34s                | TOO SLOW!                          | 1m2s                 | 19.4s        | 38.7s
Update 100 (with 100k)        |                      | 741ms                              | 35s <b>TOO SLOW!</b> |              |
Update 100 (with 500k)        |                      | 3.7s                               |                      |              |
Update 1K (with 100k)         |                      | 7s                                 |                      |              |
Update 1K (with 500k)         |                      | 1m10s                              |                      |              |
Update 1K (with 1M)           |                      | 1m42s                              |                      |              |
Delete 100 (with 100k)        |                      | 735ms                              | 37s <b>TOO SLOW!</b> |              |
Delete 100 (with 500k)        |                      | 3.6s                               |                      |              |
Delete 1k (With 100k)         | 15s initial</br>4s after | 750ms                          |                      |              |
Delete 1k (With 500k)         | 1m                   | 34s                                |                      |              |
Delete 1k (With 1M)           | 2m                   | 1m9s                               |                      |              |
Query using index (100k)      |                      | 8ms                                | 41ms                 | 68.184µs     | 157.198µs
Query non-indexed (100k)      |                      | 45ms                               | 120ms                |              |
Query distinct values (100k)  |                      | 80ms                               | 119ms                | 155ms        | 168ms
Query using index (500k)      |                      | 44ms                               | 49ms                 | 80.834µs     | 142.538µs
Query non-indexed (500k)      |                      | 629ms                              | 389ms                |              |
Query distinct values (500k)  |                      | 413ms                              | 350ms                | 770ms        | 827ms
Query using index (1M)        |                      | 42ms                               | 58ms                 | 148.86µs     | 193.038µs
Query non-indexed (1M)        |                      | 2.17s                              | 860ms                |              |
Query distinct values (1M)    |                      | 2.50s                              | 659ms                | 1.55s        | 2.5s



## Benchmark database operations for relationships

TBD

Operation              | Redisgraph | Postgre
---                    | ---        | ---
