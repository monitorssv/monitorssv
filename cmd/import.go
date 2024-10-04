package main

import (
	"encoding/json"
	"fmt"
	"github.com/monitorssv/monitorssv/config"
	"github.com/monitorssv/monitorssv/eth1/client"
	"github.com/monitorssv/monitorssv/eth1/ssv"
	"github.com/monitorssv/monitorssv/store"
	"github.com/urfave/cli/v2"
	"math/big"
	"os"
)

type MerkleProof struct {
	Root string `json:"root"`
	Data []struct {
		Address string   `json:"address"`
		Amount  string   `json:"amount"`
		Proof   []string `json:"proof"`
	} `json:"data"`
}

var importCmd = &cli.Command{
	Name:  "import",
	Usage: "Import SSV Network’s Incentive Program Merkle Proofs",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "conf-path",
			Usage: "config.yaml path",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "merkle-proof-file-path",
			Usage: "merkle-proof json file path",
			Value: "",
		},
	},
	Action: func(ctx *cli.Context) error {
		cfg, err := config.InitConfig(ctx.String("conf-path"))
		if err != nil {
			log.Errorw("InitConfig", "err", err)
			return err
		}

		store, err := store.NewStore(cfg)
		if err != nil {
			log.Errorw("NewStore", "err", err)
			return err
		}

		eth1Client, err := client.NewEth1Client(cfg)
		if err != nil {
			log.Errorw("NewEth1Client", "err", err)
			return err
		}

		chainRoot, err := ssv.GetSSVRewardMerkleRootOnChain(cfg.Network, eth1Client.GetClient())
		if err != nil {
			log.Errorw("GetSSVRewardMerkleRootOnChain", "err", err)
			return err
		}

		merkleProofFilePath := ctx.String("merkle-proof-file-path")
		data, err := os.ReadFile(merkleProofFilePath)
		if err != nil {
			log.Errorw("ReadFile", "err", err)
			return err
		}

		var merkleProof MerkleProof
		err = json.Unmarshal(data, &merkleProof)
		if err != nil {
			log.Errorw("Unmarshal", "err", err)
			return err
		}

		root := merkleProof.Root
		if root == "" {
			log.Errorw("Empty root", "err", err)
			return err
		}
		if chainRoot != root {
			log.Errorw("does not match the merkleRoot on the chain", "chainRoot", chainRoot, "root", root)
			return err
		}

		for _, mp := range merkleProof.Data {
			amount, isOk := big.NewInt(0).SetString(mp.Amount, 10)
			if !isOk {
				log.Errorw("Failed to parse amount", "amount", mp.Amount)
				return err
			}

			err = store.CreateOrUpdateSSVReward(root, mp.Address, amount, toString(mp.Proof))
			if err != nil {
				log.Errorw("CreateOrUpdateSSVReward", "err", err)
				return err
			}
		}

		log.Infow("Import Success", "root", root)
		return nil
	},
}

func toString(proofs []string) string {
	var proofStr string
	for _, proof := range proofs {
		if proofStr == "" {
			proofStr = proof
			continue
		}
		proofStr = fmt.Sprintf("%s,%s", proofStr, proof)
	}

	return proofStr
}
