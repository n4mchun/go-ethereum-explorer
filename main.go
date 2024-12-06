package main

import (
	"context"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
)

type BlockInfo struct {
	BlockNumber      uint64            `json:"blockNumber"`
	Status           string            `json:"status"`
	Timestamp        string            `json:"timestamp"`
	TransactionCount int               `json:"transactionCount"`
	Transactions     []string          `json:"transactions"`
	Details          BlockDetails      `json:"details"`
}

type BlockDetails struct {
	BlockHash  string `json:"blockHash"`
	ParentHash string `json:"parentHash"`
	StateRoot  string `json:"stateRoot"`
	Nonce      uint64 `json:"nonce"`
}

func getBlock(c echo.Context) error {
	rpcURL := "https://endpoints.omniatech.io/v1/eth/sepolia/public"

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Printf("Failed to connect to Ethereum client: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to connect to Ethereum client")
	}

	blockNumber := c.Param("blockNumber")
	blockNumberInt := new(big.Int)
	blockNumberInt, ok := blockNumberInt.SetString(blockNumber, 10)
	if !ok {
		log.Printf("Invalid block number: %s", blockNumber)
		return c.String(http.StatusBadRequest, "Invalid block number")
	}

	ctx := context.Background()
	block, err := client.BlockByNumber(ctx, blockNumberInt)
	if err != nil {
		log.Printf("Failed to retrieve block: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to retrieve block")
	}

	header := block.Header()
	status := "not finalized"
	if header.Difficulty.Cmp(big.NewInt(0)) == 0 {
		status = "finalized"
	}

	timestamp := time.Unix(int64(block.Time()), 0)

	txCount := len(block.Transactions())
	txHashes := []string{}
	for _, tx := range block.Transactions() {
		txHashes = append(txHashes, tx.Hash().Hex())
	}

	blockInfo := BlockInfo{
		BlockNumber:      block.Number().Uint64(),
		Status:           status,
		Timestamp:        timestamp.UTC().Format(time.RFC3339),
		TransactionCount: txCount,
		Transactions:     txHashes,
		Details: BlockDetails{
			BlockHash:  block.Hash().Hex(),
			ParentHash: header.ParentHash.Hex(),
			StateRoot:  header.Root.Hex(),
			Nonce:      header.Nonce.Uint64(),
		},
	}

	log.Printf("Block %s info retrieved", blockNumber)
	return c.JSON(http.StatusOK, blockInfo)
}

func root(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	e := echo.New()

	e.GET("/", root)
	e.GET("/block/:blockNumber", getBlock)

	e.Start(":8080")
}
