package parlia

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	wiggleTimeBeforeFork       = 500 * time.Millisecond // Random delay (per signer) to allow concurrent signers
	fixedBackOffTimeBeforeFork = 200 * time.Millisecond
)

func (p *Parlia) delayForRamanujanFork(snap *Snapshot, header *types.Header) time.Duration {
	delay := time.Unix(int64(header.Time), 0).Sub(time.Now()) // nolint: gosimple
	if p.chainConfig.IsRamanujan(header.Number) {
		log.Info("=== debug delayForRamanujanFork")
		return delay
	} else {
		log.Info("=== debug not delayForRamanujanFork")
	}
	wiggle := time.Duration(len(snap.Validators)/2+1) * wiggleTimeBeforeFork
	return delay + time.Duration(fixedBackOffTimeBeforeFork) + time.Duration(rand.Int63n(int64(wiggle)))
}

func (p *Parlia) blockTimeForRamanujanFork(snap *Snapshot, header, parent *types.Header) uint64 {
	blockTime := parent.Time + p.config.Period
	if p.chainConfig.IsRamanujan(header.Number) {
		log.Info("=== debug prepare blockTimeForRamanujanFork", "headerTime", header.Time, "height", "diff", header.Difficulty, "backoff", backOffTime(snap, p.val))
		blockTime = blockTime + backOffTime(snap, p.val)
	} else {
		log.Info("=== debug prepare not blockTimeForRamanujanFork")
	}
	return blockTime
}

func (p *Parlia) blockTimeVerifyForRamanujanFork(snap *Snapshot, header, parent *types.Header) error {
	if p.chainConfig.IsRamanujan(header.Number) {
		if header.Time < parent.Time+p.config.Period+backOffTime(snap, header.Coinbase) {
			return consensus.ErrFutureBlock
		}
	} else {
		log.Info("=== debug not blockTimeVerifyForRamanujanFork")
	}
	return nil
}
