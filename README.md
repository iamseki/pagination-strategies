# Pagination Strategies

This project demonstrates various pagination strategies using a sample library database populated with books and reviews. It utilizes Docker Compose to set up a PostgreSQL container and provides scripts for data generation and table creation.

## Setup :wrench:

- `docker-compose up -d`
- `go run main.go keyset.go offset.go`

## Testing :cop:

- `go test`
- `go test -bench=.` (the performance difference were more notable in my machine on pagination with >= 10k rows)
  ```shell
  goos: linux
  goarch: amd64
  pkg: pagination-strategies
  cpu: AMD Ryzen 5 5500U with Radeon Graphics         
  
  BenchmarkKeysetForwardScan-12           1000000000               0.1313 ns/op
  BenchmarkOffsetScan-12                  1000000000               0.2836 ns/op
  
  PASS
  ok      pagination-strategies   6.095s
  ```

## Usage :scroll:

- Available pagination approaches: _offset_ || _keyset_
  ```txt
  ------== OFFSET ==------ 
  http://localhost:8080/books/offset?limit=10 // OFFSET 
  http://localhost:8080/books/offset?offset=10&limit=10&limit=10 // OFFSET NEXT  
  http://localhost:8080/books/offset?offset=0&limit=10&limit=10 // OFFSET PREVIOUS

  ------== KEYSET ==------
  http://localhost:8080/books/keyset?limit=10 // KEYSET
  http://localhost:8080/books/keyset?limit=10&nextPageToken=MTA%3D // KEYSET NEXT
  http://localhost:8080/books/keyset?limit=10&previousPageToken=MTE%3D // KEYSET PREVIOUS
  ```

# WIP

There are several strategies for implementing pagination in PostgreSQL, each with its own advantages and limitations. Here are the main approaches:

1. Offset-based Pagination:

This is the simplest method, where you specify an offset value (number of rows to skip) and a limit (number of rows to retrieve) in your query.
Advantages: Easy to implement and understand, efficient for small datasets.
Disadvantages: Performance degrades for large datasets as the database needs to scan through irrelevant rows, can lead to gaps or duplicates with concurrent modifications.

2. Keyset Pagination:

This method leverages an ordering column and specific values within that column to identify the boundaries of each page.
Advantages: Efficient for sorted datasets, avoids gaps and duplicates with concurrent modifications.
Disadvantages: Requires an appropriate ordering column and might not be suitable for all scenarios.

3. Window Functions:

This approach utilizes window functions like ROW_NUMBER() or NTILE() to assign unique identifiers or partition data into pages within the query itself.
Advantages: Flexible and efficient for various scenarios, can be combined with other techniques.
Disadvantages: Might require more complex SQL code compared to simpler methods.

4. Third-party Libraries:

Several open-source libraries and frameworks offer pagination functionalities specifically designed for PostgreSQL.
Advantages: Can provide additional features and abstractions, simplify implementation.
Disadvantages: Introduce external dependencies and might require additional learning curve.

Choosing the right strategy depends on your specific needs and considerations:

Dataset size and performance: For large datasets, cursor-based or keyset pagination might be more efficient.
Concurrency and data consistency: If dealing with frequent updates, consider methods that avoid gaps or duplicates.
Complexity and development time: Offset-based pagination is simpler to implement, while window functions or libraries might require more effort.
It's important to evaluate your specific use case and choose the pagination strategy that best balances efficiency, maintainability, and compatibility with your application requirements.