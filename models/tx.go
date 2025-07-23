package models

import (
	"encoding/json"

	"github.com/bsv-blockchain/go-bn/internal/util"
	"github.com/bsv-blockchain/go-bt/v2"
	"github.com/bsv-blockchain/go-bt/v2/bscript"
	"github.com/bsv-blockchain/go-bt/v2/sighash"
)

// Output model.
type Output struct {
	BestBlock     string
	Confirmations uint32
	Coinbase      bool

	*bt.Output
}

type AddToConsensusBlacklistArgs struct {
	Funds []Fund `json:"funds"`
}

type AddToConfiscationTxIdWhitelistArgs struct {
	ConfiscationTransactions []ConfiscationTransactionDetails `json:"confiscationTxs"`
}

type ConfiscationTransactionDetails struct {
	ConfiscationTransaction ConfiscationTransaction `json:"confiscationTx"`
}

type ConfiscationTransaction struct {
	EnforceAtHeight int64  `json:"enforceAtHeight"`
	Hex             string `json:"hex"`
}

// Fund represents a fund to freeze or unfreeze
type Fund struct {
	TxOut                      TxOut     `json:"txOut"`
	EnforceAtHeight            []Enforce `json:"enforceAtHeight"`
	PolicyExpiresWithConsensus bool      `json:"policyExpiresWithConsensus"`
}

// Enforce represents block start and end to enforce freezing
type Enforce struct {
	Start int `json:"start"`
	Stop  int `json:"stop"`
}

// Txout represents transaction output
type TxOut struct {
	TxId string `json:"txId"`
	Vout int    `json:"vout"`
}

// NodeJSON return node json variant.
func (o *Output) NodeJSON() interface{} {
	return o
}

// UnmarshalJSON unmarshal response.
func (o *Output) UnmarshalJSON(b []byte) error {
	oj := struct {
		BestBlock     string `json:"bestblock"`
		Confirmations uint32 `json:"confirmations"`
		Coinbase      bool   `json:"coinbase"`
	}{}

	if err := json.Unmarshal(b, &oj); err != nil {
		return err
	}

	var out bt.Output
	if err := json.Unmarshal(b, out.NodeJSON()); err != nil {
		return err
	}

	o.BestBlock = oj.BestBlock
	o.Confirmations = oj.Confirmations
	o.Coinbase = oj.Coinbase
	*o.Output = out

	return nil
}

// OutputSetInfo model.
type OutputSetInfo struct {
	Height         uint32  `json:"height"`
	BestBlock      string  `json:"bestblock"`
	Transactions   uint32  `json:"transactions"`
	OutputCount    uint32  `json:"txouts"`
	BogoSize       uint32  `json:"bogosize"`
	HashSerialised string  `json:"hash_serialized"`
	DiskSize       uint32  `json:"disk_size"`
	TotalAmount    float64 `json:"total_amount"`
}

// OptsOutput options.
type OptsOutput struct {
	IncludeMempool bool
}

// Args convert struct into optional positional arguments.
func (o *OptsOutput) Args() []interface{} {
	return []interface{}{o.IncludeMempool}
}

// ParamsCreateRawTransaction model.
type ParamsCreateRawTransaction struct {
	Outputs []*bt.Output
	mainnet bool
}

// Args convert struct into optional positional arguments.
func (p *ParamsCreateRawTransaction) Args() []interface{} {
	outputs := make(map[string]float64, len(p.Outputs))
	for _, o := range p.Outputs {
		pkh, err := o.LockingScript.PublicKeyHash()
		if err != nil {
			outputs["invalid locking script"] = util.SatoshisToBSV(int64(o.Satoshis))
			continue
		}
		addr, err := bscript.NewAddressFromPublicKeyHash(pkh, p.mainnet)
		if err != nil {
			outputs["invalid locking script"] = util.SatoshisToBSV(int64(o.Satoshis))
		}
		outputs[addr.AddressString] = util.SatoshisToBSV(int64(o.Satoshis))
	}

	return []interface{}{outputs}
}

// SetIsMainnet set request is in mainnet context.
func (p *ParamsCreateRawTransaction) SetIsMainnet(b bool) {
	p.mainnet = b
}

// FundRawTransaction model.
type FundRawTransaction struct {
	Fee            uint64
	ChangePosition int `json:"changeposition"`
	Tx             *bt.Tx
}

// OptsFundRawTransaction options.
type OptsFundRawTransaction struct {
	ChangeAddress          string   `json:"changeAddress,omitempty"`
	ChangePosition         int      `json:"changePosition,omitempty"`
	IncludeWatching        bool     `json:"includeWatching,omitempty"`
	LockUnspents           bool     `json:"lockUnspents,omitempty"`
	ReserveChangeKey       *bool    `json:"reserveChangeKey,omitempty"`
	FeeRate                uint64   `json:"feeRate,omitempty"`
	SubtractFeeFromOutputs []uint64 `json:"subtractFeeFromOutputs,omitempty"`
}

// Args convert struct into optional positional arguments.
func (o *OptsFundRawTransaction) Args() []interface{} {
	return []interface{}{o}
}

// SendRawTransaction model.
type SendRawTransaction struct {
	Hex string
	Tx  *bt.Tx
}

// PostProcess data.
func (s *SendRawTransaction) PostProcess() error {
	var err error
	s.Tx, err = bt.NewTxFromString(s.Hex)
	return err
}

// SignedRawTransaction model.
type SignedRawTransaction struct {
	Tx       *bt.Tx `json:"tx"`
	Complete bool   `json:"complete"`
	Errors   []struct {
		TxID            string `json:"txid"`
		Vout            int    `json:"vout"`
		UnlockingScript string `json:"scriptSig"`
		Sequence        uint32 `json:"sequence"`
		Error           string `json:"error"`
	} `json:"errors"`
}

// OptsSignRawTransaction options.
type OptsSignRawTransaction struct {
	From        bt.UTXOs
	PrivateKeys []string
	SigHashType sighash.Flag
}

// Args convert struct into optional positional arguments.
func (o *OptsSignRawTransaction) Args() []interface{} {
	aa := []interface{}{[]interface{}{}, []interface{}{}}
	if o.From != nil {
		aa[0] = o.From.NodeJSON()
	}
	if o.PrivateKeys != nil && len(o.PrivateKeys) > 0 {
		aa[1] = o.PrivateKeys
	}
	return append(aa, o.SigHashType.String())
}

// OptsSendRawTransaction options.
type OptsSendRawTransaction struct {
	AllowHighFees bool
	CheckFee      bool
}

// Args convert struct into optional positional arguments.
func (o *OptsSendRawTransaction) Args() []interface{} {
	return []interface{}{o.AllowHighFees, !o.CheckFee}
}

// ParamsSendRawTransactions params.
type ParamsSendRawTransactions struct {
	Hex                      string `json:"hex"`
	AllowHighFees            bool   `json:"allowhighfees"`
	DontCheckFee             bool   `json:"dontcheckfee"`
	ListUnconfirmedAncestors bool   `json:"listunconfirmedancestors"`
}

// SendRawTransactionsResponse response.
type SendRawTransactionsResponse struct {
	Known   []string `json:"known"`
	Evicted []string `json:"evicted"`
	Invalid []struct {
		TxID         string `json:"txid"`
		RejectCode   int    `json:"reject_code"`
		RejectReason string `json:"reject_reason"`
		CollidedWith []struct {
			TxID string `json:"txid"`
			Size uint64 `json:"size"`
			Hex  string `json:"hex"`
		} `json:"collidedWith"`
	} `json:"invalid"`
	Unconfirmed []struct {
		TxID      string `json:"txid"`
		Ancestors []struct {
			TxID string `json:"txid"`
			Vin  []struct {
				TxID string `json:"txid"`
				Vout uint32 `json:"vout"`
			} `json:"vin"`
		} `json:"ancestors"`
	} `json:"unconfirmed"`
}

// AddToConsensusBlacklistNotProcessed represents a not processed transaction output or confiscation transaction
type AddToConsensusBlacklistNotProcessed []struct {
	TxOut struct {
		TxId string `json:"txId"`
		Vout int    `json:"vout"`
	}
	Reason string `json:"reason"`
}

// AddToConsensusBlacklistResponse response
type AddToConsensusBlacklistResponse struct {
	NotProcessed AddToConsensusBlacklistNotProcessed `json:"notProcessed"`
}

type WhitelistConfiscationTransaction struct {
	TxId string `json:"txId"`
}

// AddToConfiscationTransactionWhitelistNotProcessed represents a not processed confiscation transaction
type AddToConfiscationTransactionWhitelistNotProcessed []struct {
	ConfiscationTransaction WhitelistConfiscationTransaction `json:"confiscationTx"`
	Reason                  string                           `json:"reason"`
}

// AddToConfiscationTransactionWhitelistResponse represents the response for adding confiscation transactions to the whitelist
type AddToConfiscationTransactionWhitelistResponse struct {
	NotProcessed AddToConfiscationTransactionWhitelistNotProcessed `json:"notProcessed"`
}
