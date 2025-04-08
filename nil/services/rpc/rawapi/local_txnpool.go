package rawapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/NilFoundation/nil/nil/internal/types"
	"github.com/NilFoundation/nil/nil/services/txnpool"
)

func (api *LocalShardApi) SendTransaction(ctx context.Context, encoded []byte) (txnpool.DiscardReason, error) {
	if api.txnpool == nil {
		return 0, errors.New("transaction pool is not available")
	}

	var extTxn types.ExternalTransaction
	if err := extTxn.UnmarshalSSZ(encoded); err != nil {
		return 0, fmt.Errorf("failed to decode transaction: %w", err)
	}

	reasons, err := api.txnpool.Add(ctx, extTxn.ToTransaction())
	if err != nil {
		return 0, err
	}
	return reasons[0], nil
}

func (api *LocalShardApi) GetTxpoolStatus(ctx context.Context) (uint64, error) {
	return uint64(api.txnpool.GetPendingLength()), nil
}

func (api *LocalShardApi) GetTxpoolContent(ctx context.Context) ([]*types.Transaction, error) {
	txns, err := api.txnpool.Peek(0)
	if err != nil {
		return nil, err
	}
	res := make([]*types.Transaction, len(txns))
	for i, txn := range txns {
		res[i] = txn.Transaction
	}
	return res, nil
}
