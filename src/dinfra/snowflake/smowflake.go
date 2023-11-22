package snowflake

import (
	"github.com/AlgerDu/go-dream/src/dinfra"

	"errors"
	"strconv"
	"sync"
	"time"
)

type (
	NodeOptions struct {
		Epoch    int64
		NodeBits uint8
		StepBits uint8
		Node     int64
	}

	Node struct {
		mu    sync.Mutex
		epoch time.Time
		time  int64
		node  int64
		step  int64

		nodeMax   int64
		nodeMask  int64
		stepMask  int64
		timeShift uint8
		nodeShift uint8
	}
)

func NewDefaultNodeOptions() *NodeOptions {
	return &NodeOptions{
		Epoch:    1661961600000, // 2022-09-01 00:00:00
		NodeBits: 10,
		StepBits: 12,
		Node:     1,
	}
}

func NewNode(options *NodeOptions) (*Node, error) {

	n := Node{}

	n.node = options.Node
	n.nodeMax = -1 ^ (-1 << options.NodeBits)
	n.nodeMask = n.nodeMax << options.StepBits
	n.stepMask = -1 ^ (-1 << options.StepBits)
	n.timeShift = options.NodeBits + options.StepBits
	n.nodeShift = options.StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(options.Epoch/1000, (options.Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

func (n *Node) Provide() dinfra.ID {

	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.epoch).Nanoseconds() / 1000000

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := dinfra.ID((now)<<n.timeShift |
		(n.node << n.nodeShift) |
		(n.step),
	)

	return r
}
