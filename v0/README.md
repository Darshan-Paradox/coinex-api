# Coinex API Backend

This is a backend written in Go+Gin to provide price details of cryptocurrencies. These APIs are wrapper around different APIs that provide data.

PostgreSQL database is used for caching API response to minimise the calls to external APIs.

# How to run

- Download the project using `git clone`
- Run `export $(cat .env | xargs)` to export .env file
- Run `go run ./src/main` or `go build` to run the server

--OR--

- Download the project using `git clone`
- Run `zsh build.sh` or `chmod +x build.sh && ./build.sh`
