package models

import (
	"encoding/json"

	"github.com/bsv-blockchain/go-bn/internal/util"
	"github.com/bsv-blockchain/go-bn/models"
	"github.com/bsv-blockchain/go-bt/v2"
	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

// InternalDumpPrivateKey the true to form dumpprivkey response from the bitcoin node.
type InternalDumpPrivateKey struct {
	PrivateKey *primitives.PrivateKey
}

// UnmarshalJSON unmarshal the response.
func (i *InternalDumpPrivateKey) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	pk, err := primitives.PrivateKeyFromWif(s)
	if err != nil {
		return err
	}

	i.PrivateKey = pk
	return nil
}

// InternalTransaction the true to form transaction response from the bitcoin node.
type InternalTransaction struct {
	*models.Transaction
	Amount  float64 `json:"amount"`
	Fee     float64 `json:"fee"`
	Hex     string  `json:"hex"`
	Details []struct {
		Account   string  `json:"account"`
		Address   string  `json:"address"`
		Category  string  `json:"category"`
		Amount    float64 `json:"amount"`
		Label     string  `json:"label"`
		Vout      uint32  `json:"vout"`
		Fee       float64 `json:"fee"`
		Abandoned bool    `json:"abandoned"`
	} `json:"details"`
}

// PostProcess an RPC response.
func (i *InternalTransaction) PostProcess() error {
	i.Transaction.Amount = int64(util.BSVToSatoshis(i.Amount))
	i.Transaction.Fee = int64(util.BSVToSatoshis(i.Fee))

	i.Transaction.Details = make([]models.TransactionDetail, len(i.Details))
	for idx, detail := range i.Details {
		i.Transaction.Details[idx] = models.TransactionDetail{
			Account:   detail.Account,
			Abandoned: detail.Abandoned,
			Address:   detail.Address,
			Category:  detail.Category,
			Amount:    int64(util.BSVToSatoshis(detail.Amount)),
			Fee:       int64(util.BSVToSatoshis(detail.Fee)),
			Label:     detail.Label,
			Vout:      detail.Vout,
		}
	}
	var err error
	i.Tx, err = bt.NewTxFromString(i.Hex)
	return err
}
