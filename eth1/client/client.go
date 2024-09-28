package client

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/monitorssv/monitorssv/config"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"math/big"
	"time"
)

var timeout = 1 * time.Minute

type Eth1Client struct {
	client *ethclient.Client
}

func NewEth1Client(cfg *config.Config) (*Eth1Client, error) {
	client, err := ethclient.Dial(cfg.Eth1Rpc)
	if err != nil {
		return nil, err
	}
	return &Eth1Client{client: client}, nil
}

func (c *Eth1Client) GetClient() *ethclient.Client {
	return c.client
}

func (c *Eth1Client) BlockNumber() (uint64, error) {
	return utils.Retry(func() (uint64, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		return c.client.BlockNumber(ctx)
	}, utils.DefaultRetryConfig)
}

func (c *Eth1Client) CodeAt(account string) ([]byte, error) {
	return utils.Retry(func() ([]byte, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		return c.client.CodeAt(ctx, common.HexToAddress(account), nil)
	}, utils.DefaultRetryConfig)
}

func (c *Eth1Client) ChainId() (*big.Int, error) {
	return utils.Retry(func() (*big.Int, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		return c.client.ChainID(ctx)
	}, utils.DefaultRetryConfig)
}

func (c *Eth1Client) FilterLogs(q ethereum.FilterQuery) ([]types.Log, error) {
	return utils.Retry(func() ([]types.Log, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return c.client.FilterLogs(ctx, q)
	}, utils.DefaultRetryConfig)
}
