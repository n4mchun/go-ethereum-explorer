# Go Ethereum Explorer

This project is a simple Ethereum block explorer written in Go. It uses the `go-ethereum` package to connect to an Ethereum node and retrieve block and transaction information. The project also uses the `echo` web framework to create a RESTful API.

## Features

- Connects to an Ethereum node using the `go-ethereum` package.
- Retrieves block information by block number.
- Retrieves transaction information by transaction hash.
- Provides block details including block hash, parent hash, state root, nonce, and transactions.
- Provides transaction details including transaction hash, status, block number, timestamp, sender, receiver, value, transaction fee, gas price, and input data.

## Usage

### Prerequisites

- Go 1.23.4 (latest version)
- An Ethereum node RPC URL (e.g., from a service like Omnia)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/n4mchun/go-ethereum-explorer.git
    cd go-ethereum-explorer
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Run the application:

    ```sh
    go run main.go
    ```

### API Endpoints

#### Root Endpoint

- **URL:** `/`
- **Method:** `GET`
- **Description:** Returns a simple "Hello, World!" message.

#### Get Block Information

- **URL:** `/block/:blockNumber`
- **Method:** `GET`
- **Description:** Retrieves information about a specific block by its block number.
- **Parameters:**
  - `blockNumber` (path parameter): The block number to retrieve information for.
- **Response:**
  - `200 OK`: Returns block information in JSON format.
  - `400 Bad Request`: If the block number is invalid.
  - `500 Internal Server Error`: If there is an error connecting to the Ethereum client or retrieving the block.

#### Get Transaction Information

- **URL:** `/tx/:txHash`
- **Method:** `GET`
- **Description:** Retrieves information about a specific transaction by its hash.
- **Parameters:**
  - `txHash` (path parameter): The transaction hash to retrieve information for.
- **Response:**
  - `200 OK`: Returns transaction information in JSON format.
  - `500 Internal Server Error`: If there is an error connecting to the Ethereum client or retrieving the transaction.

### Example

To retrieve information about block number `123456`, you can use the following curl command:

```sh
curl http://localhost:8080/block/123456
```

To retrieve information about a transaction with hash `0x123...`, you can use the following curl command:

```sh
curl http://localhost:8080/tx/0x123...
```