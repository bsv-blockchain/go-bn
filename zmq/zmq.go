package zmq

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/bsv-blockchain/go-bc"
	"github.com/bsv-blockchain/go-bt/v2"
	"github.com/go-zeromq/zmq4"
)

// Topic a subscription topic.
type Topic string

// Subscription topics.
const (
	TopicHashTx                  Topic = "hashtx"
	TopicHashBlock               Topic = "hashblock"
	TopicInvalidTx               Topic = "invalidtx"
	TopicDiscardFromMempool      Topic = "discardfrommempool"
	TopicRemovedFromMempoolBlock Topic = "removedfrommempoolblock"

	TopicRawTx    Topic = "rawtx"
	TopicRawBlock Topic = "rawblock"
)

type nodeMq struct {
	mu            sync.RWMutex
	conn          zmq4.Socket
	connected     bool
	onErrFn       ErrorFunc
	cfg           *nodeMqCfg
	subscriptions map[Topic]MessageFunc
}

// NodeMQ interfaces connecting and subscribing to a bitcoin node NodeMQ connection.
type NodeMQ interface {
	Connect() error
	Subscribe(topic Topic, fn MessageFunc) error
	SubscribeHashTx(fn HashFunc) error
	SubscribeHashBlock(fn HashFunc) error
	SubscribeDiscardFromMempool(fn DiscardFunc) error
	SubscribeRemovedFromMempoolBlock(fn DiscardFunc) error
	SubscribeRawTx(fn RawTxFunc) error
	SubscribeRawBlock(fn RawBlockFunc) error
	Unsubscribe(topic Topic) error
}

// NewNodeMQ build and return a new zmq.ZMQ configured via the provided opt funcs.
func NewNodeMQ(oo ...NodeMQOptFunc) NodeMQ {
	cfg := &nodeMqCfg{
		optionValue: "hash",
		errorFn:     defaultOnError,
		topics: map[Topic]bool{
			TopicHashBlock:               true,
			TopicHashTx:                  true,
			TopicDiscardFromMempool:      true,
			TopicRemovedFromMempoolBlock: true,
			TopicInvalidTx:               true,

			TopicRawTx:    false,
			TopicRawBlock: false,
		},
	}
	for _, o := range oo {
		o(cfg)
	}

	if cfg.zmqSocket == nil {
		cfg.zmqSocket = zmq4.NewSub(cfg.ctx, zmq4.WithID(zmq4.SocketIdentity("sub")))
	}

	return &nodeMq{
		cfg:           cfg,
		subscriptions: make(map[Topic]MessageFunc),
		onErrFn:       cfg.errorFn,
		conn:          cfg.zmqSocket,
	}
}

// Connect to the bitcoin node 0MQ.
func (n *nodeMq) Connect() error {
	if err := n.cfg.validate(); err != nil {
		return err
	}

	if err := n.conn.Dial(n.cfg.host); err != nil {
		return err
	}

	defer func() {
		if !n.connected {
			return
		}

		if err := n.conn.Close(); err != nil {
			n.onErrFn(context.Background(), err)
		}
		n.connected = false
	}()

	if err := n.conn.SetOption(zmq4.OptionSubscribe, n.cfg.optionValue); err != nil {
		return err
	}

	if n.cfg.raw {
		if err := n.conn.SetOption(zmq4.OptionSubscribe, "raw"); err != nil {
			return err
		}
	}

	for {
		msg, err := n.conn.Recv()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			n.onErrFn(context.Background(), err)
			continue
		}
		n.connected = true
		func() {
			n.mu.RLock()
			defer n.mu.RUnlock()

			if fn, ok := n.subscriptions[Topic(msg.Frames[0])]; ok {
				go fn(context.Background(), msg.Frames)
			}
		}()
	}
}

// Subscribe to a topic on a bitcoin node 0MQ.
func (n *nodeMq) Subscribe(topic Topic, fn MessageFunc) error {
	if ok := n.cfg.topics[topic]; !ok {
		return fmt.Errorf("%w: %s", ErrInvalidTopic, topic)
	}

	if !n.cfg.allowOverwrite {
		if err := func() error {
			n.mu.RLock()
			defer n.mu.RUnlock()
			if _, ok := n.subscriptions[topic]; ok {
				return fmt.Errorf("%w: %s", ErrAlreadySubscribed, topic)
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.subscriptions[topic] = fn
	return nil
}

// SubscribeHashTx subscribe to `hashtx` and receive its messages parsed.
func (n *nodeMq) SubscribeHashTx(fn HashFunc) error {
	return n.Subscribe(TopicHashTx, func(ctx context.Context, bb [][]byte) {
		fn(ctx, hex.EncodeToString(bb[1]))
	})
}

// SubscribeHashBlock subscribe to `hashblock` and receive its messages parsed.
func (n *nodeMq) SubscribeHashBlock(fn HashFunc) error {
	return n.Subscribe(TopicHashBlock, func(ctx context.Context, bb [][]byte) {
		fn(ctx, hex.EncodeToString(bb[1]))
	})
}

// SubscribeDiscardFromMempool subscribe to `discardfrommempool` and receive its messages parsed.
func (n *nodeMq) SubscribeDiscardFromMempool(fn DiscardFunc) error {
	return n.Subscribe(TopicDiscardFromMempool, func(ctx context.Context, bb [][]byte) {
		var d MempoolDiscard
		if err := json.Unmarshal(bb[1], &d); err != nil {
			n.onErrFn(ctx, err)
			return
		}
		fn(ctx, &d)
	})
}

// SubscribeRemovedFromMempoolBlock subscribe to `removedfrommempoolblock` and receive its messages parsed.
func (n *nodeMq) SubscribeRemovedFromMempoolBlock(fn DiscardFunc) error {
	return n.Subscribe(TopicRemovedFromMempoolBlock, func(ctx context.Context, bb [][]byte) {
		var d MempoolDiscard
		if err := json.Unmarshal(bb[1], &d); err != nil {
			n.onErrFn(ctx, err)
			return
		}
		fn(ctx, &d)
	})
}

// SubscribeRawTx subscribe to `rawtx` and receive its messages parsed.
func (n *nodeMq) SubscribeRawTx(fn RawTxFunc) error {
	return n.Subscribe(TopicRawTx, func(ctx context.Context, bb [][]byte) {
		tx, err := bt.NewTxFromBytes(bb[1])
		if err != nil {
			n.onErrFn(ctx, err)
			return
		}
		fn(ctx, tx)
	})
}

// SubscribeRawBlock subscribe to `rawblock` and receive its messages parsed.
func (n *nodeMq) SubscribeRawBlock(fn RawBlockFunc) error {
	return n.Subscribe(TopicRawBlock, func(ctx context.Context, bb [][]byte) {
		blk, err := bc.NewBlockFromBytes(bb[1])
		if err != nil {
			n.onErrFn(ctx, err)
			return
		}
		fn(ctx, blk)
	})
}

// Unsubscribe from a topic on the bitcoin node 0MQ.
func (n *nodeMq) Unsubscribe(topic Topic) error {
	if ok := n.cfg.topics[topic]; !ok {
		return fmt.Errorf("%w: %s", ErrInvalidTopic, topic)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	delete(n.subscriptions, topic)
	return nil
}

func defaultOnError(_ context.Context, err error) {
	fmt.Fprintln(os.Stderr, err)
}
