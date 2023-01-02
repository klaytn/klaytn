// Modifications Copyright 2023 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package nodecmd

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"gopkg.in/urfave/cli.v1"
)

const pgmInput = `
<extra> <header file (json format)>
<vote>  <bytes>
<gov>   <bytes>
<key>   <keystore path> <password>
`

var ErrInvalidParam = errors.New("Invalid length of parameter")

var UtilCommand = cli.Command{
	Action:    utils.MigrateFlags(parse),
	Name:      "util",
	Usage:     "offline utility " + pgmInput,
	ArgsUsage: " ",
	Category:  "MISCELLANEOUS COMMANDS",
}

func hex2Bytes(s string) []byte {
	if data, err := hexutil.Decode(s); err == nil {
		return data
	} else {
		panic(err)
	}
}

func printUsage() error {
	fmt.Println("Usage: ./kxn util " + pgmInput)
	return ErrInvalidParam
}

func validateInput(ctx *cli.Context, parseTyp string) error {
	switch parseTyp {
	case "vote", "extra", "gov":
		if len(ctx.Args()) != 2 {
			return printUsage()
		}
	case "key":
		if len(ctx.Args()) != 3 {
			return printUsage()
		}
	}
	return nil
}

func parse(ctx *cli.Context) error {
	parseTyp, m := ctx.Args().Get(0), make(map[string]interface{})
	if err := validateInput(ctx, parseTyp); err != nil {
		return err
	}
	switch parseTyp {
	case "vote":
		data := ctx.Args().Get(1)
		if err := parseVote(m, hex2Bytes(data)); err != nil {
			return err
		}
	case "extra":
		headerFile := ctx.Args().Get(1)
		if err := parseExtra(m, headerFile); err != nil {
			return err
		}
	case "gov":
		data := ctx.Args().Get(1)
		if err := parseGov(m, hex2Bytes(data)); err != nil {
			return err
		}
	case "key":
		keystorePath, passwd := ctx.Args().Get(1), ctx.Args().Get(2)
		if err := extractKeypair(m, keystorePath, passwd); err != nil {
			return err
		}
	default:
		return printUsage()
	}
	prettyPrint(m)
	return nil
}

func prettyPrint(m map[string]interface{}) {
	if b, err := json.MarshalIndent(m, "", "  "); err == nil {
		fmt.Println(string(b))
	} else {
		panic(err)
	}
}

func extractKeypair(m map[string]interface{}, keystorePath, passwd string) error {
	keyjson, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		return err
	}
	key, err := keystore.DecryptKey(keyjson, passwd)
	if err != nil {
		return err
	}
	addr := key.GetAddress().String()
	pubkey := key.GetPrivateKey().PublicKey
	privkey := key.GetPrivateKey()
	m["addr"] = addr
	m["privkey"] = hex.EncodeToString(crypto.FromECDSA(privkey))
	m["pubkey"] = hex.EncodeToString(crypto.FromECDSAPub(&pubkey))
	return nil
}

func parseGov(m map[string]interface{}, bytes []byte) error {
	var b []byte
	if err := rlp.DecodeBytes(bytes, &b); err == nil {
		if err := json.Unmarshal(b, &m); err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func parseHeaderFile(headerFile string) (*types.Header, common.Hash, error) {
	header := new(types.Header)
	bytes, err := ioutil.ReadFile(headerFile)
	if err != nil {
		return nil, common.Hash{}, err
	}
	if err = json.Unmarshal(bytes, &header); err != nil {
		return nil, common.Hash{}, err
	}
	var hash common.Hash
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return header, hash, nil
}

func parseExtra(m map[string]interface{}, headerFile string) error {
	header, hash, err := parseHeaderFile(headerFile)
	if err != nil {
		return err
	}
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return err
	}

	validators := make([]string, len(istanbulExtra.Validators))
	for idx, addr := range istanbulExtra.Validators {
		validators[idx] = addr.String()
	}
	cSeals := make([]string, len(istanbulExtra.CommittedSeal))
	for idx, cSeal := range istanbulExtra.CommittedSeal {
		cSeals[idx] = hexutil.Encode(cSeal)
	}

	proposer, err := istanbul.GetSignatureAddress(hash.Bytes(), istanbulExtra.Seal)
	if err != nil {
		return err
	}
	m["validators"] = validators
	m["seal"] = hexutil.Encode(istanbulExtra.Seal)
	m["committedSeal"] = cSeals
	m["validatorSize"] = len(validators)
	m["committedSealSize"] = len(cSeals)
	m["proposer"] = proposer.String()
	return nil
}

func parseVote(m map[string]interface{}, bytes []byte) error {
	vote := new(governance.GovernanceVote)
	if err := rlp.DecodeBytes(bytes, &vote); err == nil {
		m["validator"] = vote.Validator.String()
		m["key"] = vote.Key
		switch governance.GovernanceKeyMap[vote.Key] {
		case params.GovernanceMode, params.MintingAmount, params.MinimumStake, params.Ratio:
			m["value"] = string(vote.Value.([]uint8))
		case params.GoverningNode:
			m["value"] = common.BytesToAddress(vote.Value.([]uint8)).String()
		case params.Epoch, params.CommitteeSize, params.UnitPrice, params.StakeUpdateInterval,
			params.ProposerRefreshInterval, params.ConstTxGasHumanReadable, params.Policy, params.Timeout,
			params.LowerBoundBaseFee, params.UpperBoundBaseFee, params.GasTarget, params.MaxBlockGasUsedForBaseFee, params.BaseFeeDenominator:
			v := vote.Value.([]uint8)
			v = append(make([]byte, 8-len(v)), v...)
			m["value"] = binary.BigEndian.Uint64(v)
		case params.UseGiniCoeff, params.DeferredTxFee:
			v := vote.Value.([]uint8)
			v = append(make([]byte, 8-len(v)), v...)
			if binary.BigEndian.Uint64(v) != uint64(0) {
				m["value"] = true
			} else {
				m["value"] = false
			}
		case params.AddValidator, params.RemoveValidator:
			m["value"] = common.BytesToAddress(vote.Value.([]uint8)).String()
		}
		return nil
	} else {
		return err
	}
}
