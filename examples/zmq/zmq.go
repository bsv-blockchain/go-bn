package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/bsv-blockchain/go-bc"
	"github.com/bsv-blockchain/go-bt/v2"

	"github.com/bsv-blockchain/go-bn/zmq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	z := zmq.NewNodeMQ(
		zmq.WithContext(ctx),
		zmq.WithHost("tcp://localhost:28332"),
		zmq.WithRaw(),
		zmq.WithErrorHandler(func(_ context.Context, err error) {
			log.Println("error found", err)
		}),
	)

	if err := z.Subscribe(zmq.TopicInvalidTx, func(_ context.Context, bb [][]byte) {
		log.Println("invalid tx", hex.EncodeToString(bb[1]))
	}); err != nil {
		panic(err)
	}

	if err := z.SubscribeRawTx(func(_ context.Context, tx *bt.Tx) {
		bb, err := json.Marshal(tx)
		if err != nil {
			panic(err)
		}

		log.Println("TX:", string(bb))
	}); err != nil {
		panic(err)
	}

	if err := z.SubscribeRawBlock(func(_ context.Context, blk *bc.Block) {
		bb, err := json.Marshal(blk)
		if err != nil {
			panic(err)
		}

		log.Println("Block:", string(bb))
	}); err != nil {
		panic(err)
	}

	if err := z.SubscribeHashBlock(func(_ context.Context, hash string) {
		log.Println("BLOCK HASH", hash)
	}); err != nil {
		panic(err)
	}

	log.Fatal(z.Connect())
}
