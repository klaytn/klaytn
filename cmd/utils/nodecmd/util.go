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
	"os"

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

const (
	DECODE_EXTRA = "decode-extra"
	DECODE_VOTE  = "decode-vote"
	DECODE_GOV   = "decode-gov"
	DECRYPT_KEY  = "decrypt-keystore"
)

var pgmInput = fmt.Sprintf(`
  %17s  <header file (json format)>
  %17s  <bytes>
  %17s  <bytes>
  %17s  <keystore path> <password>
`, DECODE_EXTRA, DECODE_VOTE, DECODE_GOV, DECRYPT_KEY)

var ErrInvalidParam = errors.New("Invalid length of parameter")

var UtilCommand = cli.Command{
	Action:    utils.MigrateFlags(decode),
	Name:      "util",
	Usage:     "offline utility" + pgmInput,
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
	fmt.Printf("Usage: %s util <command> <args>", os.Args[0])
	fmt.Println(pgmInput)
	return ErrInvalidParam
}

func validateInput(ctx *cli.Context, decodeTyp string) error {
	switch decodeTyp {
	case DECODE_EXTRA, DECODE_VOTE, DECODE_GOV:
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

func decode(ctx *cli.Context) error {
	decodeTyp := ctx.Args().Get(0)
	if err := validateInput(ctx, decodeTyp); err != nil {
		return err
	}
	var (
		m   map[string]interface{}
		err error
	)
	switch decodeTyp {
	case DECODE_VOTE:
		data := ctx.Args().Get(1)
		m, err = decodeVote(hex2Bytes(data))
		if err != nil {
			return err
		}
	case DECODE_EXTRA:
		headerFile := ctx.Args().Get(1)
		m, err = decodeExtra(headerFile)
		if err != nil {
			return err
		}
	case DECODE_GOV:
		data := ctx.Args().Get(1)
		m, err = decodeGov(hex2Bytes(data))
		if err != nil {
			return err
		}
	case DECRYPT_KEY:
		keystorePath, passwd := ctx.Args().Get(1), ctx.Args().Get(2)
		m, err = extractKeypair(keystorePath, passwd)
		if err != nil {
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

func extractKeypair(keystorePath, passwd string) (map[string]interface{}, error) {
	keyjson, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyjson, passwd)
	if err != nil {
		return nil, err
	}
	addr := key.GetAddress().String()
	pubkey := key.GetPrivateKey().PublicKey
	privkey := key.GetPrivateKey()
	m := make(map[string]interface{})
	m["addr"] = addr
	m["privkey"] = hex.EncodeToString(crypto.FromECDSA(privkey))
	m["pubkey"] = hex.EncodeToString(crypto.FromECDSAPub(&pubkey))
	return m, nil
}

func decodeGov(bytes []byte) (map[string]interface{}, error) {
	var b []byte
	m := make(map[string]interface{})
	if err := rlp.DecodeBytes(bytes, &b); err == nil {
		if err := json.Unmarshal(b, &m); err == nil {
			return m, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
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

func decodeExtra(headerFile string) (map[string]interface{}, error) {
	header, hash, err := parseHeaderFile(headerFile)
	if err != nil {
		return nil, err
	}
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	m := make(map[string]interface{})
	m["validators"] = validators
	m["seal"] = hexutil.Encode(istanbulExtra.Seal)
	m["committedSeal"] = cSeals
	m["validatorSize"] = len(validators)
	m["committedSealSize"] = len(cSeals)
	m["proposer"] = proposer.String()
	return m, nil
}

func decodeVote(bytes []byte) (map[string]interface{}, error) {
	vote := new(governance.GovernanceVote)
	m := make(map[string]interface{})
	if err := rlp.DecodeBytes(bytes, &vote); err == nil {
		m["validator"] = vote.Validator.String()
		m["key"] = vote.Key
		switch governance.GovernanceKeyMap[vote.Key] {
		case params.GovernanceMode, params.MintingAmount, params.MinimumStake, params.Ratio, params.Kip82Ratio:
			m["value"] = string(vote.Value.([]uint8))
		case params.GoverningNode, params.GovParamContract:
			m["value"] = common.BytesToAddress(vote.Value.([]uint8)).String()
		case params.Epoch, params.CommitteeSize, params.UnitPrice, params.DeriveShaImpl, params.StakeUpdateInterval,
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
			if v, ok := vote.Value.([]uint8); ok {
				m["value"] = common.BytesToAddress(v)
			} else if addresses, ok := vote.Value.([]interface{}); ok {
				if len(addresses) == 0 {
					return nil, governance.ErrValueTypeMismatch
				}
				var nodeAddresses []common.Address
				for _, item := range addresses {
					if in, ok := item.([]uint8); !ok || len(in) != common.AddressLength {
						return nil, governance.ErrValueTypeMismatch
					}
					nodeAddresses = append(nodeAddresses, common.BytesToAddress(item.([]uint8)))
				}
				m["value"] = nodeAddresses
			} else {
				return nil, governance.ErrValueTypeMismatch
			}
		}
		return m, nil
	} else {
		return nil, err
	}
}
