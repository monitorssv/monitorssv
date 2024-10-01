package ssv

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type Cluster struct {
	ClusterId   string
	Owner       common.Address
	OperatorIds []uint64
	ClusterInfo ISSVNetworkCoreCluster
}

func CalcClusterId(owner common.Address, operatorIds []uint64) string {
	var packed []byte
	packed = append(packed, owner.Bytes()...)
	for _, operatorId := range operatorIds {
		bigInt := new(big.Int).SetUint64(operatorId)
		bytes := bigInt.Bytes()
		paddedBytes := make([]byte, 32)
		copy(paddedBytes[32-len(bytes):], bytes)
		packed = append(packed, paddedBytes...)
	}

	hash := crypto.Keccak256(packed)
	return hex.EncodeToString(hash)
}

type ContractCreationInfo struct {
	ContractAddress string
	ContractCreator string
	TxHash          string
}

func GetContractCreator(address string, endpoint, apiKey string) (*ContractCreationInfo, error) {
	if len(address) == 0 {
		return nil, fmt.Errorf("no addresses")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var resp struct {
		Status  string
		Message string
		Result  json.RawMessage
	}
	err := requests.URL(endpoint).
		Client(&http.Client{Timeout: 15 * time.Second}).
		Path("/api").
		Param("module", "contract").
		Param("action", "getcontractcreation").
		Param("contractaddresses", address).
		Param("apikey", apiKey).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.Status != "1" || !strings.HasPrefix(resp.Message, "OK") {
		return nil, fmt.Errorf(
			"bad status %s (%s): %s",
			resp.Status,
			resp.Message,
			resp.Result,
		)
	}

	// Decode the response.
	var results []struct {
		ContractAddress string `json:"contractAddress"`
		ContractCreator string `json:"contractCreator"`
		TxHash          string `json:"txHash"`
	}
	if err := json.Unmarshal(resp.Result, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contract creation: %w", err)
	}
	if len(results) != 1 {
		return nil, fmt.Errorf(
			"failed to get contract creation: expected %d results, got %d",
			1,
			len(results),
		)
	}

	result := results[0]
	contractAddress, err := hex.DecodeString(result.ContractAddress[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode deployer address: %w", err)
	}
	deployerAddress, err := hex.DecodeString(result.ContractCreator[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode contract creator: %w", err)
	}
	txHash, err := hex.DecodeString(result.TxHash[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx hash: %w", err)
	}
	if len(contractAddress) != 20 || len(deployerAddress) != 20 || len(txHash) != 32 {
		return nil, fmt.Errorf("failed to decode contract creation")
	}
	if common.BytesToAddress(contractAddress).String() != address {
		return nil, fmt.Errorf("contract creation address mismatch")
	}

	return &ContractCreationInfo{
		ContractAddress: common.BytesToAddress(contractAddress).String(),
		ContractCreator: common.BytesToAddress(deployerAddress).String(),
		TxHash:          common.Bytes2Hex(txHash),
	}, nil
}
