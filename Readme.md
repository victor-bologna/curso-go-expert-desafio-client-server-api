## Client/Server API

This repository contains a Go language HTTP API representing client and server sides to get actual currency of 1 Dollar(USD) to Brazil Reais (BRL).

The server/server.go application represents a localhost:8080 server with two endpoints:
- "/Cotacao": This endpoint must retrieve data from the server within 200 milliseconds, return to the user only the value of "bid" field and persist the whole data into SQLite3 (in memory database) within 10 milliseconds.
- "/VerifyCotacoes": The endpoint responsible to show a JSON list of all quotations requested while the server is alive.

The client/client.go is a client side app which requests the data from server/server.go app within the limit of 300 milliseconds then save the actual quotation value into a file called "cotacao.txt".

If the requests could not be accomplished within the time limit the system throws a error in the console.