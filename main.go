package main

import (
	"context"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
)

type BlockInfo struct {
	BlockNumber      uint64            `json:"blockNumber"`
	Status           string            `json:"status"`
	Timestamp        string            `json:"timestamp"`
	TransactionCount int               `json:"transactionCount"`
	Transactions     []string          `json:"transactions"`
	Withdrawals 	 int               `json:"withdrawals"`
	Details          BlockDetails      `json:"details"`
}

type BlockDetails struct {
	BlockHash  string `json:"blockHash"`
	ParentHash string `json:"parentHash"`
	StateRoot  string `json:"stateRoot"`
	Nonce      uint64 `json:"nonce"`
}

type TransactionInfo struct {
	TxHash      string `json:"txHash"`
	Status      string `json:"status"`
	BlockNumber uint64 `json:"blockNumber"`
	Timestamp   string `json:"timestamp"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	TxFee       string `json:"txFee"`
	GasPrice    string `json:"gasPrice"`
	InputData   string `json:"inputData"`
}

func connectToEthereumClient() (*ethclient.Client, error) {
	rpcURL := "https://endpoints.omniatech.io/v1/eth/sepolia/public"
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Printf("Failed to connect to Ethereum client: %v", err)
		return nil, err
	}
	return client, nil
}

func getBlockByNumber(client *ethclient.Client, blockNumber string) (*types.Block, error) {
	blockNumberInt := new(big.Int)
	blockNumberInt, ok := blockNumberInt.SetString(blockNumber, 10)
	if !ok {
		log.Printf("Invalid block number: %s", blockNumber)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid block number")
	}

	ctx := context.Background()
	block, err := client.BlockByNumber(ctx, blockNumberInt)
	if err != nil {
		log.Printf("Failed to retrieve block: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve block")
	}

	return block, nil
}

func GetBlock(c echo.Context) error {
	client, err := connectToEthereumClient()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to connect to Ethereum client")
	}

	blockNumber := c.Param("blockNumber")
	block, err := getBlockByNumber(client, blockNumber)
	if err != nil {
		return err
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

	withdrawals := len(block.Withdrawals())

	blockInfo := BlockInfo{
		BlockNumber:      block.Number().Uint64(),
		Status:           status,
		Timestamp:        timestamp.UTC().Format(time.RFC3339),
		TransactionCount: txCount,
		Transactions:     txHashes,
		Withdrawals:      withdrawals,
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

func GetTransaction(c echo.Context) error {
	client, err := connectToEthereumClient()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to connect to Ethereum client")
	}

	txHash := c.Param("txHash")

	ctx := context.Background()
	tx, _, err := client.TransactionByHash(ctx, common.HexToHash(txHash)) // _ = isPending
	if err != nil {
		log.Printf("Failed to retrieve transaction: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to retrieve transaction")
	}

	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		log.Printf("Failed to retrieve transaction receipt: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to retrieve transaction receipt")
	}

	block, err := client.BlockByHash(ctx, receipt.BlockHash)
	if err != nil {
		log.Printf("Failed to retrieve block: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to retrieve block")
	}

	timestamp := time.Unix(int64(block.Time()), 0)

	status := "fail"
	if receipt.Status == 1 {
		status = "success"
	}

	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Printf("Failed to retrieve sender: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to retrieve sender")
	}

	txFee := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), tx.GasPrice())

	txInfo := TransactionInfo{
		TxHash:      tx.Hash().Hex(),
		Status:      status,
		BlockNumber: receipt.BlockNumber.Uint64(),
		Timestamp:   timestamp.UTC().Format(time.RFC3339),
		From:        from.Hex(),
		To:          tx.To().Hex(),
		Value:       tx.Value().String(),
		TxFee:       txFee.String(),
		GasPrice:    tx.GasPrice().String(),
		InputData:   "0x" + common.Bytes2Hex(tx.Data()),
	}

	log.Printf("Transaction %s info retrieved", txHash)
	return c.JSON(http.StatusOK, txInfo)
}

func root(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	e := echo.New()

	e.GET("/", root)
	e.GET("/block/:blockNumber", GetBlock)
	e.GET("/tx/:txHash", GetTransaction)

	e.Start(":8080")
}