// Modifications Copyright 2019 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
//
// This file is derived from core/tx_journal.go (2018/06/04).
// Modified and improved for the klaytn development.

package sc

import (
	"errors"
	"io"
	"os"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/node/sc/bridgepool"
	"github.com/klaytn/klaytn/rlp"
)

var (
	ErrNoActiveAddressJournal = errors.New("no active address journal")
	ErrDuplicatedJournal      = errors.New("duplicated journal is inserted")
	ErrEmptyBridgeAddress     = errors.New("empty bridge address is not allowed")
	ErrEmptyBridgeAlias       = errors.New("empty bridge Alias")
)

// bridgeAddrJournal is a rotating log of addresses with the aim of storing locally
// created addresses to allow deployed bridge contracts to survive node restarts.
type bridgeAddrJournal struct {
	path       string         // Filesystem path to store the addresses at
	writer     io.WriteCloser // Output stream to write new addresses into
	cache      map[common.Address]*BridgeJournal
	aliasCache map[string]common.Address
}

// newBridgeAddrJournal creates a new bridge addr journal to
func newBridgeAddrJournal(path string) *bridgeAddrJournal {
	return &bridgeAddrJournal{
		path:       path,
		cache:      make(map[common.Address]*BridgeJournal),
		aliasCache: make(map[string]common.Address),
	}
}

// load parses a address journal dump from disk, loading its contents into
// the specified pool.
func (journal *bridgeAddrJournal) load(add func(journal BridgeJournal) error) error {
	// Skip the parsing if the journal file doens't exist at all
	if _, err := os.Stat(journal.path); os.IsNotExist(err) {
		return nil
	}
	// Open the journal for loading any past addresses
	input, err := os.Open(journal.path)
	if err != nil {
		return err
	}
	defer input.Close()

	// Temporarily discard any journal additions (don't double add on load)
	journal.writer = new(bridgepool.DevNull)
	defer func() { journal.writer = nil }()

	// Inject all addresses from the journal into the pool
	stream := rlp.NewStream(input, 0)
	total, dropped := 0, 0

	var (
		failure        error
		firstDecodeErr = true
	)
	for {
		// Parse the next address and terminate on error
		addr := new(BridgeJournal)
		if !firstDecodeErr {
			addr.isLegacyBridgeJournal = true
		}
	DecodeAgainWithLegacyStruct:
		if err = stream.Decode(addr); err != nil {
			if err != io.EOF {
				failure = err
			}
			if err == ErrBridgeAliasFormatDecode && firstDecodeErr {
				input.Close()
				input, err = os.Open(journal.path)
				if err != nil {
					return err
				}
				firstDecodeErr = false
				stream.Reset(input, 0)
				goto DecodeAgainWithLegacyStruct
			}
			break
		}

		total++

		if err := add(*addr); err != nil {
			failure = err
			dropped++
		}
	}
	logger.Info("Loaded local bridge journal", "addrs", total, "dropped", dropped)

	return failure
}

// ChangeBridgeAlias changes oldBridgeAlias to newBridgeAlias
func (journal *bridgeAddrJournal) ChangeBridgeAlias(oldBridgeAlias, newBridgeAlias string) error {
	if addr, ok := journal.aliasCache[oldBridgeAlias]; ok {
		delete(journal.aliasCache, oldBridgeAlias)
		journal.aliasCache[newBridgeAlias] = addr
		journal.cache[addr].BridgeAlias = newBridgeAlias
		return nil
	}
	return ErrEmptyBridgeAlias
}

// insert adds the specified address to the local disk journal.
func (journal *bridgeAddrJournal) insert(bridgeAlias string, localAddress common.Address, remoteAddress common.Address) error {
	if len(bridgeAlias) != 0 && journal.aliasCache[bridgeAlias] != (common.Address{}) {
		return ErrDuplicatedJournal
	}
	if journal.cache[localAddress] != nil {
		return ErrDuplicatedJournal
	}

	if journal.writer == nil {
		return ErrNoActiveAddressJournal
	}
	empty := common.Address{}
	if localAddress == empty || remoteAddress == empty {
		return ErrEmptyBridgeAddress
	}
	// TODO-Klaytn-ServiceChain: support false paired
	item := BridgeJournal{
		bridgeAlias,
		localAddress,
		remoteAddress,
		false,
		false,
	}
	if err := rlp.Encode(journal.writer, &item); err != nil {
		return err
	}

	journal.cache[localAddress] = &item
	if len(bridgeAlias) != 0 {
		journal.aliasCache[bridgeAlias] = localAddress
	}
	return nil
}

// rotate regenerates the addresses journal based on the current contents of
// the address pool.
func (journal *bridgeAddrJournal) rotate(all []*BridgeJournal) error {
	// Close the current journal (if any is open)
	if journal.writer != nil {
		if err := journal.writer.Close(); err != nil {
			return err
		}
		journal.writer = nil
	}
	// Generate a new journal with the contents of the current pool
	replacement, err := os.OpenFile(journal.path+".new", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	journaled := 0
	for _, journal := range all {
		if err = rlp.Encode(replacement, journal); err != nil {
			replacement.Close()
			return err
		}
		journaled++
	}
	replacement.Close()

	// Replace the live journal with the newly generated one
	if err = os.Rename(journal.path+".new", journal.path); err != nil {
		return err
	}
	sink, err := os.OpenFile(journal.path, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	journal.writer = sink
	logger.Info("Regenerated local addr journal", "addrs", journaled, "accounts", len(all))

	return nil
}

// close flushes the addresses journal contents to disk and closes the file.
func (journal *bridgeAddrJournal) close() error {
	var err error

	if journal.writer != nil {
		err = journal.writer.Close()
		journal.writer = nil
	}
	return err
}
