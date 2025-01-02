package client

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
}

const (
	defaultIndexChunkSize  = 1000
	defaultPubKeyChunkSize = 75
)

var log = logging.Logger("monitor-service")

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}
}

type bytesHexStr []byte

func (s *bytesHexStr) UnmarshalText(b []byte) error {
	if s == nil {
		return fmt.Errorf("cannot unmarshal bytes into nil")
	}
	if len(b) >= 2 && b[0] == '0' && b[1] == 'x' {
		b = b[2:]
	}
	out := make([]byte, len(b)/2)
	hex.Decode(out, b)
	*s = out
	return nil
}

type uint64Str uint64

func (s *uint64Str) UnmarshalJSON(b []byte) error {
	return Uint64Unmarshal((*uint64)(s), b)
}

func Uint64Unmarshal(v *uint64, b []byte) error {
	if v == nil {
		return errors.New("nil dest in uint64 decoding")
	}
	if len(b) == 0 {
		return errors.New("empty uint64 input")
	}
	if b[0] == '"' || b[0] == '\'' {
		if len(b) == 1 || b[len(b)-1] != b[0] {
			return errors.New("uneven/missing quotes")
		}
		b = b[1 : len(b)-1]
	}
	n, err := strconv.ParseUint(string(b), 0, 64)
	if err != nil {
		return err
	}
	*v = n
	return nil
}

type StandardProposerDutiesResponse struct {
	DependentRoot string                 `json:"dependent_root"`
	Data          []StandardProposerDuty `json:"data"`
}

type StandardProposerDuty struct {
	Pubkey         string    `json:"pubkey"`
	ValidatorIndex uint64Str `json:"validator_index"`
	Slot           uint64Str `json:"slot"`
}

func (c *Client) GetEpochProposer(epoch uint64) (*StandardProposerDutiesResponse, error) {
	return utils.Retry(func() (*StandardProposerDutiesResponse, error) {
		proposerResp, err := c.get(fmt.Sprintf("%s/eth/v1/validator/duties/proposer/%d", c.endpoint, epoch))
		if err != nil {
			log.Warnf("error retrieving proposer duties for epoch %v: %s", epoch, err)
			return nil, err
		}

		var parsedProposerResponse StandardProposerDutiesResponse
		err = json.Unmarshal(proposerResp, &parsedProposerResponse)
		if err != nil {
			return nil, fmt.Errorf("error parsing proposer duties: %s", err)
		}

		return &parsedProposerResponse, nil
	}, utils.DefaultRetryConfig)
}

type StandardBeaconHeaderResponse struct {
	Data struct {
		Root   string `json:"root"`
		Header struct {
			Message struct {
				Slot          uint64Str `json:"slot"`
				ProposerIndex uint64Str `json:"proposer_index"`
				ParentRoot    string    `json:"parent_root"`
				StateRoot     string    `json:"state_root"`
				BodyRoot      string    `json:"body_root"`
			} `json:"message"`
			Signature string `json:"signature"`
		} `json:"header"`
	} `json:"data"`
	Finalized bool `json:"finalized"`
}

func (c *Client) GetSlotHeader(slot uint64) (*StandardBeaconHeaderResponse, error) {
	return utils.Retry(func() (*StandardBeaconHeaderResponse, error) {
		resHeaders, err := c.get(fmt.Sprintf("%s/eth/v1/beacon/headers/%d", c.endpoint, slot))
		if err != nil {
			log.Warnf("error retrieving headers at slot %v: %s", slot, err)
			return nil, err
		}
		var parsedHeaders StandardBeaconHeaderResponse
		err = json.Unmarshal(resHeaders, &parsedHeaders)
		if err != nil {
			return nil, fmt.Errorf("error parsing header-response at slot %v: %s", slot, err)
		}
		return &parsedHeaders, nil
	}, utils.DefaultRetryConfig)
}

func (c *Client) GetLatestSlot() (uint64, error) {
	return utils.Retry(func() (uint64, error) {
		resHeaders, err := c.get(fmt.Sprintf("%s/eth/v1/beacon/headers/head", c.endpoint))
		if err != nil {
			log.Warnf("error retrieving headers: %s", err)
			return 0, err
		}
		var parsedHeaders StandardBeaconHeaderResponse
		err = json.Unmarshal(resHeaders, &parsedHeaders)
		if err != nil {
			return 0, fmt.Errorf("error parsing header-response: %s", err)
		}
		return uint64(parsedHeaders.Data.Header.Message.Slot), nil
	}, utils.DefaultRetryConfig)
}

type StandardValidatorsResponse struct {
	Data []StandardValidatorEntry `json:"data"`
}
type StandardValidatorEntry struct {
	Index     uint64Str `json:"index"`
	Balance   uint64Str `json:"balance"`
	Status    string    `json:"status"`
	Validator struct {
		Pubkey                     string    `json:"pubkey"`
		WithdrawalCredentials      string    `json:"withdrawal_credentials"`
		EffectiveBalance           uint64Str `json:"effective_balance"`
		Slashed                    bool      `json:"slashed"`
		ActivationEligibilityEpoch uint64Str `json:"activation_eligibility_epoch"`
		ActivationEpoch            uint64Str `json:"activation_epoch"`
		ExitEpoch                  uint64Str `json:"exit_epoch"`
		WithdrawableEpoch          uint64Str `json:"withdrawable_epoch"`
	} `json:"validator"`
}

func (c *Client) GetSlotValidatorsByPubKey(slot uint64, validatorPubKeys []string) (map[string]*StandardValidatorEntry, error) {
	return utils.Retry(func() (map[string]*StandardValidatorEntry, error) {
		if len(validatorPubKeys) == 0 {
			return nil, errors.New("validatorPubKeys must not be empty")
		}

		if len(validatorPubKeys) > defaultPubKeyChunkSize {
			return c.chunkedValidatorsByPubKey(slot, validatorPubKeys)
		}

		validatorsResp, err := c.get(fmt.Sprintf("%s/eth/v1/beacon/states/%d/validators?id=%s", c.endpoint, slot, strings.Join(validatorPubKeys, ",")))
		if err != nil {
			log.Warnf("error retrieving validators for slot %v: %s", slot, err)
			return nil, err
		}

		parsedValidators := &StandardValidatorsResponse{}
		err = json.Unmarshal(validatorsResp, parsedValidators)
		if err != nil {
			return nil, fmt.Errorf("error parsing slot validators: %s", err)
		}

		res := make(map[string]*StandardValidatorEntry)
		for _, validator := range parsedValidators.Data {
			res[validator.Validator.Pubkey] = &validator
		}
		return res, nil
	}, utils.DefaultRetryConfig)
}

func (c *Client) chunkedValidatorsByPubKey(slot uint64, validatorPubKeys []string) (map[string]*StandardValidatorEntry, error) {
	res := make(map[string]*StandardValidatorEntry)
	for i := 0; i < len(validatorPubKeys); i += defaultPubKeyChunkSize {
		chunkStart := i
		chunkEnd := i + defaultPubKeyChunkSize
		if len(validatorPubKeys) < chunkEnd {
			chunkEnd = len(validatorPubKeys)
		}
		chunk := validatorPubKeys[chunkStart:chunkEnd]
		chunkRes, err := c.GetSlotValidatorsByPubKey(slot, chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain chunk: %s", err)
		}
		for k, v := range chunkRes {
			res[k] = v
		}
	}

	return res, nil
}

func (c *Client) GetSlotValidatorsByIndex(slot uint64, validatorIndices []uint64) (map[string]*StandardValidatorEntry, error) {
	return utils.Retry(func() (map[string]*StandardValidatorEntry, error) {
		if len(validatorIndices) == 0 {
			return nil, errors.New("validatorIndices must not be empty")
		}

		if len(validatorIndices) > defaultIndexChunkSize {
			return c.chunkedValidatorsByIndex(slot, validatorIndices)
		}

		validatorsResp, err := c.get(fmt.Sprintf("%s/eth/v1/beacon/states/%d/validators?id=%s", c.endpoint, slot, joinUint64(validatorIndices, ",")))
		if err != nil {
			log.Warnf("error retrieving validators for slot %v: %s", slot, err)
			return nil, err
		}

		parsedValidators := &StandardValidatorsResponse{}
		err = json.Unmarshal(validatorsResp, parsedValidators)
		if err != nil {
			return nil, fmt.Errorf("error parsing slot validators: %s", err)
		}

		res := make(map[string]*StandardValidatorEntry)
		for _, validator := range parsedValidators.Data {
			res[validator.Validator.Pubkey] = &validator
		}

		return res, nil
	}, utils.DefaultRetryConfig)
}

func (c *Client) chunkedValidatorsByIndex(slot uint64, validatorIndices []uint64) (map[string]*StandardValidatorEntry, error) {
	res := make(map[string]*StandardValidatorEntry)
	for i := 0; i < len(validatorIndices); i += defaultIndexChunkSize {
		chunkStart := i
		chunkEnd := i + defaultIndexChunkSize
		if len(validatorIndices) < chunkEnd {
			chunkEnd = len(validatorIndices)
		}
		chunk := validatorIndices[chunkStart:chunkEnd]
		chunkRes, err := c.GetSlotValidatorsByIndex(slot, chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain chunk: %s", err)
		}
		for k, v := range chunkRes {
			res[k] = v
		}
	}
	return res, nil
}

func joinUint64(nums []uint64, sep string) string {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = fmt.Sprintf("%d", num)
	}

	return strings.Join(strs, sep)
}

type StandardFinalityCheckpointsResponse struct {
	Data struct {
		PreviousJustified struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"previous_justified"`
		CurrentJustified struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"current_justified"`
		Finalized struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"finalized"`
	} `json:"data"`
}

func (c *Client) GetFinalizedEpoch() (uint64, error) {
	return utils.Retry(func() (uint64, error) {
		finalityResp, err := c.get(fmt.Sprintf("%s/eth/v1/beacon/states/%s/finality_checkpoints", c.endpoint, "head"))
		if err != nil {
			log.Warnf("error retrieving finality checkpoints of head: %s", err)
			return 0, err
		}

		var parsedFinality StandardFinalityCheckpointsResponse
		err = json.Unmarshal(finalityResp, &parsedFinality)
		if err != nil {
			return 0, fmt.Errorf("error parsing finality checkpoints of head: %s", err)
		}

		var finalizedEpoch = uint64(parsedFinality.Data.Finalized.Epoch)
		if finalizedEpoch > 0 {
			finalizedEpoch--
		}
		return finalizedEpoch, nil
	}, utils.DefaultRetryConfig)
}

type StandardV2BlockResponse struct {
	Version             string         `json:"version"`
	ExecutionOptimistic bool           `json:"execution_optimistic"`
	Finalized           bool           `json:"finalized"`
	Data                AnySignedBlock `json:"data"`
}
type ProposerSlashing struct {
	SignedHeader1 struct {
		Message struct {
			Slot          uint64Str `json:"slot"`
			ProposerIndex uint64Str `json:"proposer_index"`
			ParentRoot    string    `json:"parent_root"`
			StateRoot     string    `json:"state_root"`
			BodyRoot      string    `json:"body_root"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"signed_header_1"`
	SignedHeader2 struct {
		Message struct {
			Slot          uint64Str `json:"slot"`
			ProposerIndex uint64Str `json:"proposer_index"`
			ParentRoot    string    `json:"parent_root"`
			StateRoot     string    `json:"state_root"`
			BodyRoot      string    `json:"body_root"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"signed_header_2"`
}

type AttesterSlashing struct {
	Attestation1 struct {
		AttestingIndices []uint64Str `json:"attesting_indices"`
		Signature        string      `json:"signature"`
		Data             struct {
			Slot            uint64Str `json:"slot"`
			Index           uint64Str `json:"index"`
			BeaconBlockRoot string    `json:"beacon_block_root"`
			Source          struct {
				Epoch uint64Str `json:"epoch"`
				Root  string    `json:"root"`
			} `json:"source"`
			Target struct {
				Epoch uint64Str `json:"epoch"`
				Root  string    `json:"root"`
			} `json:"target"`
		} `json:"data"`
	} `json:"attestation_1"`
	Attestation2 struct {
		AttestingIndices []uint64Str `json:"attesting_indices"`
		Signature        string      `json:"signature"`
		Data             struct {
			Slot            uint64Str `json:"slot"`
			Index           uint64Str `json:"index"`
			BeaconBlockRoot string    `json:"beacon_block_root"`
			Source          struct {
				Epoch uint64Str `json:"epoch"`
				Root  string    `json:"root"`
			} `json:"source"`
			Target struct {
				Epoch uint64Str `json:"epoch"`
				Root  string    `json:"root"`
			} `json:"target"`
		} `json:"data"`
	} `json:"attestation_2"`
}

type Attestation struct {
	AggregationBits string `json:"aggregation_bits"`
	Signature       string `json:"signature"`
	Data            struct {
		Slot            uint64Str `json:"slot"`
		Index           uint64Str `json:"index"`
		BeaconBlockRoot string    `json:"beacon_block_root"`
		Source          struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"source"`
		Target struct {
			Epoch uint64Str `json:"epoch"`
			Root  string    `json:"root"`
		} `json:"target"`
	} `json:"data"`
}

type Deposit struct {
	Proof []string `json:"proof"`
	Data  struct {
		Pubkey                string    `json:"pubkey"`
		WithdrawalCredentials string    `json:"withdrawal_credentials"`
		Amount                uint64Str `json:"amount"`
		Signature             string    `json:"signature"`
	} `json:"data"`
}

type VoluntaryExit struct {
	Message struct {
		Epoch          uint64Str `json:"epoch"`
		ValidatorIndex uint64Str `json:"validator_index"`
	} `json:"message"`
	Signature string `json:"signature"`
}
type Eth1Data struct {
	DepositRoot  string    `json:"deposit_root"`
	DepositCount uint64Str `json:"deposit_count"`
	BlockHash    string    `json:"block_hash"`
}
type SyncAggregate struct {
	SyncCommitteeBits      string `json:"sync_committee_bits"`
	SyncCommitteeSignature string `json:"sync_committee_signature"`
}
type WithdrawalPayload struct {
	Index          uint64Str   `json:"index"`
	ValidatorIndex uint64Str   `json:"validator_index"`
	Address        bytesHexStr `json:"address"`
	Amount         uint64Str   `json:"amount"`
}
type ExecutionPayload struct {
	ParentHash    bytesHexStr   `json:"parent_hash"`
	FeeRecipient  bytesHexStr   `json:"fee_recipient"`
	StateRoot     bytesHexStr   `json:"state_root"`
	ReceiptsRoot  bytesHexStr   `json:"receipts_root"`
	LogsBloom     bytesHexStr   `json:"logs_bloom"`
	PrevRandao    bytesHexStr   `json:"prev_randao"`
	BlockNumber   uint64Str     `json:"block_number"`
	GasLimit      uint64Str     `json:"gas_limit"`
	GasUsed       uint64Str     `json:"gas_used"`
	Timestamp     uint64Str     `json:"timestamp"`
	ExtraData     bytesHexStr   `json:"extra_data"`
	BaseFeePerGas uint64Str     `json:"base_fee_per_gas"`
	BlockHash     bytesHexStr   `json:"block_hash"`
	Transactions  []bytesHexStr `json:"transactions"`
	// present only after capella
	Withdrawals []WithdrawalPayload `json:"withdrawals"`
	// present only after deneb
	BlobGasUsed   uint64Str `json:"blob_gas_used"`
	ExcessBlobGas uint64Str `json:"excess_blob_gas"`
}
type SignedBLSToExecutionChange struct {
	Message struct {
		ValidatorIndex     uint64Str   `json:"validator_index"`
		FromBlsPubkey      bytesHexStr `json:"from_bls_pubkey"`
		ToExecutionAddress bytesHexStr `json:"to_execution_address"`
	} `json:"message"`
	Signature bytesHexStr `json:"signature"`
}
type AnySignedBlock struct {
	Message struct {
		Slot          uint64Str `json:"slot"`
		ProposerIndex uint64Str `json:"proposer_index"`
		ParentRoot    string    `json:"parent_root"`
		StateRoot     string    `json:"state_root"`
		Body          struct {
			RandaoReveal      string             `json:"randao_reveal"`
			Eth1Data          Eth1Data           `json:"eth1_data"`
			Graffiti          string             `json:"graffiti"`
			ProposerSlashings []ProposerSlashing `json:"proposer_slashings"`
			AttesterSlashings []AttesterSlashing `json:"attester_slashings"`
			Attestations      []Attestation      `json:"attestations"`
			Deposits          []Deposit          `json:"deposits"`
			VoluntaryExits    []VoluntaryExit    `json:"voluntary_exits"`

			// not present in phase0 blocks
			SyncAggregate *SyncAggregate `json:"sync_aggregate,omitempty"`

			// not present in phase0/altair blocks
			ExecutionPayload *ExecutionPayload `json:"execution_payload"`

			// present only after capella
			SignedBLSToExecutionChange []*SignedBLSToExecutionChange `json:"bls_to_execution_changes"`

			// present only after deneb
			BlobKZGCommitments []bytesHexStr `json:"blob_kzg_commitments"`
		} `json:"body"`
	} `json:"message"`
	Signature bytesHexStr `json:"signature"`
}

// GetBlockBySlot When the slot is missed, ErrNotFound is returned
// So don't use utils.Retry
func (c *Client) GetBlockBySlot(slot uint64) (*StandardV2BlockResponse, error) {
	return utils.Retry(func() (*StandardV2BlockResponse, error) {
		resp, err := c.get(fmt.Sprintf("%s/eth/v2/beacon/blocks/%d", c.endpoint, slot))
		if err != nil {
			log.Warnf("error retrieving block data at slot %v: %s", slot, err)
			return nil, err
		}
		var parsedResponse StandardV2BlockResponse
		err = json.Unmarshal(resp, &parsedResponse)
		if err != nil {
			log.Warnf("error parsing block data at slot %v: %v", slot, err)
			return nil, fmt.Errorf("error parsing block-response at slot %v: %s", slot, err)
		}
		return &parsedResponse, nil
	}, utils.DefaultRetryConfig)
}

var ErrNotFound = errors.New("not found 404")

func (c *Client) get(url string) ([]byte, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("url: %v, error-response: %s", url, data)
	}

	// enhanced compatibility
	var apiRes APIResponse
	if err := json.Unmarshal(data, &apiRes); err == nil {
		if apiRes.Code == 404 {
			return nil, ErrNotFound
		}
	}

	return data, err
}

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
