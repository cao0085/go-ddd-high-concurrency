package idgen

import (
	"sync"
	"time"
)

const (
	epoch int64 = 1704067200000

	nodeBits     = 10
	sequenceBits = 12

	maxNodeID   = -1 ^ (-1 << nodeBits)     // 1023
	maxSequence = -1 ^ (-1 << sequenceBits) // 4095

	nodeShift      = sequenceBits
	timestampShift = nodeBits + sequenceBits
)

// IDGenerator generates unique snowflake IDs
type IDGenerator struct {
	mu       sync.Mutex
	nodeID   int64
	sequence int64
	lastTime int64
}

// NewIDGenerator creates a new IDGenerator with the given node ID
// nodeID must be between 0 and 1023
func NewIDGenerator(nodeID int64) (*IDGenerator, error) {
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, ErrInvalidNodeID
	}

	return &IDGenerator{
		nodeID:   nodeID,
		sequence: 0,
		lastTime: 0,
	}, nil
}

// Generate creates a new unique ID
func (g *IDGenerator) Generate() int64 {
	g.mu.Lock() // 奈秒等級
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == g.lastTime { // 同一毫秒內
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	id := ((now - epoch) << timestampShift) |
		(g.nodeID << nodeShift) |
		g.sequence

	return id
}
