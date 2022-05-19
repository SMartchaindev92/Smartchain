// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package gwc

import (
	"github.com/SMartchaindev92/Smartchain/core"
	"github.com/SMartchaindev92/Smartchain/core/forkid"
	"github.com/SMartchaindev92/Smartchain/p2p/enode"
	"github.com/SMartchaindev92/Smartchain/rlp"
)

// ethEntry is the "gwc" ENR entry which advertises gwc protocol
// on the discovery network.
type ethEntry struct {
	ForkID forkid.ID // Fork identifier per EIP-2124

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e ethEntry) ENRKey() string {
	return "gwc"
}

// startEthEntryUpdate starts the ENR updater loop.
func (gwc *Gquantumchain) startEthEntryUpdate(ln *enode.LocalNode) {
	var newHead = make(chan core.ChainHeadEvent, 10)
	sub := gwc.blockchain.SubscribeChainHeadEvent(newHead)

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-newHead:
				ln.Set(gwc.currentEthEntry())
			case <-sub.Err():
				// Would be nice to sync with gwc.Stop, but there is no
				// good way to do that.
				return
			}
		}
	}()
}

func (gwc *Gquantumchain) currentEthEntry() *ethEntry {
	return &ethEntry{ForkID: forkid.NewID(gwc.blockchain.Config(), gwc.blockchain.Genesis().Hash(),
		gwc.blockchain.CurrentHeader().Number.Uint64())}
}
