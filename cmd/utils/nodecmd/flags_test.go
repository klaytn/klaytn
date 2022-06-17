// Modifications Copyright 2019 The klaytn Authors
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
//
// This file is derived from cmd/geth/genesis_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package nodecmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	genesis = `{"config":{"chainId":2019,"istanbul":{"epoch":30,"policy":2,"sub":13},"unitPrice":25000000000,"deriveShaImpl":2,"governance":{"governingNode":"0xdddfb991127b43e209c2f8ed08b8b3d0b5843d36","governanceMode":"single","reward":{"mintingAmount":9600000000000000000,"ratio":"34/54/12","useGiniCoeff":false,"deferredTxFee":true,"stakingUpdateInterval":60,"proposerUpdateInterval":30,"minimumStake":5000000}}},"timestamp":"0x5ce33d6e","extraData":"0x0000000000000000000000000000000000000000000000000000000000000000f89af85494dddfb991127b43e209c2f8ed08b8b3d0b5843d3694195ba9cc787b00796a7ae6356e5b656d4360353794777fd033b5e3bcaad6006bc9f481ffed6b83cf5a94d473284239f704adccd24647c7ca132992a28973b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0","governanceData":null,"blockScore":"0x1","alloc":{"195ba9cc787b00796a7ae6356e5b656d43603537":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"},"777fd033b5e3bcaad6006bc9f481ffed6b83cf5a":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"},"d473284239f704adccd24647c7ca132992a28973":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"},"dddfb991127b43e209c2f8ed08b8b3d0b5843d36":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"},"f4316f69d9522667c0674afcd8638288489f0333":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"}},"number":"0x0","gasUsed":"0x0","parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`
)

const (
	FlagTypeBoolean = iota
	FlagTypeArgument
)

const (
	ErrorIncorrectUsage = iota
	ErrorInvalidValue
	ErrorFatal
	//TODO-Klaytn-Node fix the configuration to filter wrong input flags before the klay server is launched
	NonError // This error case expects an error, but currently it does not filter the wrong value.
)

var (
	commonThreeErrors = []string{"abcdefg", "1234567", "!@#$%^&"}
	commonTwoErrors   = []string{"abcdefg", "!@#$%^&"}
)

var flagsWithValues = []struct {
	flag        string
	flagType    uint16
	values      []string
	wrongValues []string
	errors      []int
}{
	// TODO-Klaytn-Node the flag is not defined on any klay binaries
	//{
	//	flag:     "--networktype",
	//	flagType: FlagTypeArgument,
	//	// values: []string{"mn", "scn"},
	//	values: []string{},
	//	wrongValues: []string{"baobab", "abcdefg", "1234567", "!@#$%^&"},
	//	errors: []string{},
	//},
	{
		flag:        "--dbtype",
		flagType:    FlagTypeArgument,
		values:      []string{"LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3"},
		wrongValues: append(commonThreeErrors, "oracle"),
		errors:      []int{NonError, NonError, NonError, NonError},
	},
	{
		flag:        "--srvtype",
		flagType:    FlagTypeArgument,
		values:      []string{"http", "fasthttp"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--keystore",
		flagType:    FlagTypeArgument,
		values:      []string{""},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--networkid",
		flagType:    FlagTypeArgument,
		values:      []string{"1", "1000", "1001", "12312"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--identity",
		flagType:    FlagTypeArgument,
		values:      []string{"abc", "abde-", "oai121"},
		wrongValues: []string{},
		errors:      []int{},
	},
	//TODO-Klaytn-Node the flag is not defined on any klay binaries
	//{
	//	flag:        "--docroot",
	//	flagType:    FlagTypeBoolean,
	//	values:      []string{},
	//},
	{
		flag:        "--syncmode",
		flagType:    FlagTypeArgument,
		values:      []string{"full"}, //[]string{"fast", "full"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--gcmode",
		flagType:    FlagTypeArgument,
		values:      []string{"full", "archive"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorFatal, ErrorFatal, ErrorFatal},
	},
	{
		flag:     "--lightkdf",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--txpool.nolocals",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--txpool.journal",
		flagType:    FlagTypeArgument,
		values:      []string{"transactions.rlp"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--txpool.journal-interval",
		flagType:    FlagTypeArgument,
		values:      []string{"1h0m0s", "0h0m0s", "0h0m1s", "0h1m0s", "0.5h0.5m0.5s"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.pricelimit",
		flagType:    FlagTypeArgument,
		values:      []string{"1"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.pricebump",
		flagType:    FlagTypeArgument,
		values:      []string{"10"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.exec-slots.account",
		flagType:    FlagTypeArgument,
		values:      []string{"16"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.exec-slots.all",
		flagType:    FlagTypeArgument,
		values:      []string{"4096"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.nonexec-slots.account",
		flagType:    FlagTypeArgument,
		values:      []string{"64"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.nonexec-slots.all",
		flagType:    FlagTypeArgument,
		values:      []string{"1024"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	//TODO-Klaytn-Node the flag is not defined on any klay binaries
	//{
	//	flag:        "--txpool.keeplocals",
	//	flagType:    FlagTypeBoolean,
	//	values:      []string{},
	//},
	{
		flag:        "--txpool.lifetime",
		flagType:    FlagTypeArgument,
		values:      []string{"0h20m0s"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--db.single",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--db.num-statetrie-shards",
		flagType: FlagTypeArgument,
		//values:    []string{"1", "2"},
		values:      []string{"1"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--db.leveldb.cache-size",
		flagType:    FlagTypeArgument,
		values:      []string{"768"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--db.leveldb.compression",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "1", "2", "3"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--db.leveldb.no-buffer-pool",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--db.no-parallel-write",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--state.cache-size",
		flagType:    FlagTypeArgument,
		values:      []string{"64", "128", "256", "512"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--state.block-interval",
		flagType:    FlagTypeArgument,
		values:      []string{"64", "128", "256"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--cache.type",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "1", "2"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--cache.scale",
		flagType:    FlagTypeArgument,
		values:      []string{"100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--cache.level",
		flagType:    FlagTypeBoolean,
		values:      []string{"saving", "normal", "extreme"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--cache.memory",
		flagType:    FlagTypeArgument,
		values:      []string{"16"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--state.trie-cache-limit",
		flagType:    FlagTypeArgument,
		values:      []string{"512", "1024", "2048", "4096"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--sendertxhashindexing",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--childchainindexing",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--targetgaslimit",
		flagType:    FlagTypeArgument,
		values:      []string{"4712388"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--scsigner",
		flagType:    FlagTypeArgument,
		values:      []string{"0x777fd033b5e3bcaad6006bc9f481ffed6b83cf5a"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorFatal, ErrorFatal, ErrorFatal},
	},
	{
		flag:        "--rewardbase",
		flagType:    FlagTypeArgument,
		values:      []string{"0x777fd033b5e3bcaad6006bc9f481ffed6b83cf5a"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorFatal, ErrorFatal, ErrorFatal},
	},
	{
		flag:        "--extradata",
		flagType:    FlagTypeArgument,
		values:      []string{"0x0000000000000000000000000000000000000000000000000000000000000000f89af85494dddfb991127b43e209c2f8ed08b8b3d0b5843d3694195ba9cc787b00796a7ae6356e5b656d4360353794777fd033b5e3bcaad6006bc9f481ffed6b83cf5a94d473284239f704adccd24647c7ca132992a28973b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--txresend.interval",
		flagType:    FlagTypeArgument,
		values:      []string{"3", "5", "7"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txresend.max-count",
		flagType:    FlagTypeArgument,
		values:      []string{"1000", "2000"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--txresend.use-legacy",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--unlock",
		flagType:    FlagTypeArgument,
		values:      []string{"", "0x0", "0x777fd033b5e3bcaad6006bc9f481ffed6b83cf5a"},
		wrongValues: []string{"abcdefg", "!@#$%^&", "0x921jfinowaae333"},
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--password",
		flagType:    FlagTypeArgument,
		values:      []string{"abcd", "aoije091"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:     "--vmdebug",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--vmlog",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "1", "2", "3"},
		wrongValues: append(commonThreeErrors, "4"),
		errors:      []int{NonError, NonError, NonError, NonError},
	},
	{
		flag:     "--metrics",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--prometheus",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--prometheusport",
		flagType:    FlagTypeArgument,
		values:      []string{"61001"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:     "--rpc",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--rpcaddr",
		flagType:    FlagTypeArgument,
		values:      []string{"localhost", "123.123.123.123"},
		wrongValues: append(commonThreeErrors, "123.123.123.256"),
		errors:      []int{NonError, NonError, NonError, NonError},
	},
	{
		flag:        "--rpcport",
		flagType:    FlagTypeArgument,
		values:      []string{"8551"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:        "--rpccorsdomain",
		flagType:    FlagTypeArgument,
		values:      []string{"", "localhost", "123.123.123.123"},
		wrongValues: append(commonThreeErrors, "123.123.123.256"),
		errors:      []int{NonError, NonError, NonError, NonError},
	},
	{
		flag:        "--rpcvhosts",
		flagType:    FlagTypeArgument,
		values:      []string{"*"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--rpcapi",
		flagType:    FlagTypeArgument,
		values:      []string{"", "klay", "klay,personal,istanbul,debug,miner"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--ipcdisable",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--ipcpath",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--ws",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--wsaddr",
		flagType:    FlagTypeArgument,
		values:      []string{"localhost"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--wsport",
		flagType:    FlagTypeArgument,
		values:      []string{"8552"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--grpc",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--grpcaddr",
		flagType:    FlagTypeArgument,
		values:      []string{"localhost", "123.123.123.123"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--grpcport",
		flagType:    FlagTypeArgument,
		values:      []string{"8553"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:        "--wsapi",
		flagType:    FlagTypeArgument,
		values:      []string{"", "klay", "klay,personal,istanbul,debug,miner"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--wsorigins",
		flagType:    FlagTypeArgument,
		values:      []string{""},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--exec",
		flagType:    FlagTypeArgument,
		values:      []string{"klay.blockNumber", "klat.getBlock(0)", "governance.chainConfig.governance.reward.proposerUpdateInterval"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--preload",
		flagType:    FlagTypeArgument,
		values:      []string{"abc.js", "tmp.js", "tmp"},
		wrongValues: []string{},
		errors:      []int{},
	},
	//TODO-Klaytn-Node the flag is not defined on any klay binaries
	//{
	//	flag:        "--nodetype",
	//	flagType:    FlagTypeArgument,
	//	values:      []string{"cn", "pn", "en"},
	//  wrongValues: []string{},
	//  errors:      []int{},
	//},
	{
		flag:        "--maxconnections",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "30", "25000"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--maxpendpeers",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "30", "50"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--port",
		flagType:    FlagTypeArgument,
		values:      []string{"32323", "30303"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:        "--subport",
		flagType:    FlagTypeArgument,
		values:      []string{"32324", "32325", "32327"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:     "--multichannel",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--bootnodes",
		flagType:    FlagTypeArgument,
		values:      []string{"0xf4316f69d9522667c0674afcd8638288489f0333", "", "0xf4316f69d9522667c0674afcd8638288489f0333, d473284239f704adccd24647c7ca132992a28973"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--nodekey",
		flagType:    FlagTypeArgument,
		values:      []string{""},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--nodekeyhex",
		flagType:    FlagTypeArgument,
		values:      []string{"8da4ef21b864d2cc526dbdb2a120bd2874c36c9d0a1fb7f8c63d7f7a8b41de8f"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorFatal, ErrorFatal, ErrorFatal},
	},
	{
		flag:        "--nat",
		flagType:    FlagTypeArgument,
		values:      []string{"any", "none", "upnp", "pmp", "extip:127.0.0.1"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:     "--nodiscover",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--netrestrict",
		flagType:    FlagTypeArgument,
		values:      []string{},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--chaintxperiod",
		flagType:    FlagTypeArgument,
		values:      []string{"1", "5", "100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--chaintxlimit",
		flagType:    FlagTypeArgument,
		values:      []string{"100", "200"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--jspath",
		flagType:    FlagTypeArgument,
		values:      []string{".", "root/abc/efg"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--baobab",
		flagType: FlagTypeBoolean,
	},
	//TODO-Klaytn-Node the flag is not defined on any klay binaries
	//{
	//	flag:        "--bnaddr",
	//	flagType:    FlagTypeArgument,
	//	values:      []string{},
	//	wrongValues: []string{},
	//	errors:      []int{},
	//},
	//TODO-Klaytn-Node the flag is not defined on any klay binaries
	//{
	//	flag:        "--genkey",
	//	flagType:    FlagTypeArgument,
	//	values:      []string{},
	//	wrongValues: []string{},
	//	errors:      []int{},
	//},
	//TODO-Klaytn-Node the flag is not defined on any klay binaries
	//{
	//	flag:        "--writeaddress",
	//	flagType:    FlagTypeBoolean,
	//},
	{
		flag:     "--mainbridge",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--subbridge",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--mainbridgeport",
		flagType:    FlagTypeArgument,
		values:      []string{"50505", "23232"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:        "--subbridgeport",
		flagType:    FlagTypeArgument,
		values:      []string{"50505", "23232"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:     "--vtrecovery",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--vtrecoveryinterval",
		flagType:    FlagTypeArgument,
		values:      []string{"60", "100", "200"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--scnewaccount",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--dbsyncer",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--dbsyncer.db.host",
		flagType:    FlagTypeArgument,
		values:      []string{"localhost", "123.123.123.123", "127.0.0.1"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--dbsyncer.db.port",
		flagType:    FlagTypeArgument,
		values:      []string{"3306", "32323"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--dbsyncer.db.name",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--dbsyncer.db.user",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--dbsyncer.db.password",
		flagType:    FlagTypeArgument,
		values:      []string{"aboaise", "jaooao122!@", "18231@#!@412!"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:     "--dbsyncer.logmode",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--dbsyncer.db.max.idle",
		flagType:    FlagTypeArgument,
		values:      []string{"50", "100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.db.max.open",
		flagType:    FlagTypeArgument,
		values:      []string{"30", "50", "100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.db.max.lifetime",
		flagType:    FlagTypeArgument,
		values:      []string{"1h0m0s", "0h0m10s"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.block.channel.size",
		flagType:    FlagTypeArgument,
		values:      []string{"5"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.mode",
		flagType:    FlagTypeArgument,
		values:      []string{"single", "multi", "context"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--dbsyncer.genquery.th",
		flagType:    FlagTypeArgument,
		values:      []string{"50", "123"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.insert.th",
		flagType:    FlagTypeArgument,
		values:      []string{"30", "123"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.bulk.size",
		flagType:    FlagTypeArgument,
		values:      []string{"200"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.event.mode",
		flagType:    FlagTypeArgument,
		values:      []string{"block", "head"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--dbsyncer.max.block.diff",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "5", "100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--autorestart.daemon.path",
		flagType:    FlagTypeArgument,
		values:      []string{"~/klaytn/bin/kcnd", "~/klaytn/bin/kpnd", "~/klaytn/bin/kend"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--autorestart.enable",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--autorestart.timeout",
		flagType:    FlagTypeArgument,
		values:      []string{"10m", "60s", "1h"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
}

func testFlags(t *testing.T, flag string, value string, idx int) {
	datadir := tmpdir(t)
	defer os.RemoveAll(datadir)

	json := filepath.Join(datadir, "genesis.json")
	if err := ioutil.WriteFile(json, []byte(genesis), 0600); err != nil {
		t.Fatalf("test %d: failed to write genesis file: %v", idx, err)
	}

	runKlay(t, "klay-test-flag", "--verbosity", "0", "--datadir", datadir, "init", json).WaitExit()

	klay := runKlay(t, "klay-test-flag", "--datadir", datadir, flag, value)
	klay.ExpectExit()
}

func testWrongFlags(t *testing.T, flag string, value string, idx int, expectedError string) {
	datadir := tmpdir(t)
	defer os.RemoveAll(datadir)

	json := filepath.Join(datadir, "genesis.json")
	if err := ioutil.WriteFile(json, []byte(genesis), 0600); err != nil {
		t.Fatalf("test %d: failed to write genesis file: %v", idx, err)
	}

	runKlay(t, "klay-test-flag", "--verbosity", "0", "--datadir", datadir, "init", json).WaitExit()

	klay := runKlay(t, "klay-test-flag", "--datadir", datadir, flag, value)
	klay.ExpectRegexp(expectedError)
}

func TestFlags(t *testing.T) {
	expectedError := []string{
		"Incorrect Usage. flag provided but not defined: (.*)",
		"Incorrect Usage. invalid value (.*)",
		"Fatal:(.*)",
		"(.*)",
	}
	for idx, fwv := range flagsWithValues {
		switch fwv.flagType {
		case FlagTypeBoolean:
			t.Run("klay-test-flag"+fwv.flag, func(t *testing.T) {
				testFlags(t, fwv.flag, "", idx)
			})
		case FlagTypeArgument:
			for _, item := range fwv.values {
				t.Run("klay-test-flag"+fwv.flag+"-"+item, func(t *testing.T) {
					testFlags(t, fwv.flag, item, idx)
				})
			}
			for idx2, wrongItem := range fwv.wrongValues {
				t.Run("klay-test-flag"+fwv.flag+"-"+wrongItem, func(t *testing.T) {
					testWrongFlags(t, fwv.flag, wrongItem, idx, expectedError[fwv.errors[idx2]])
				})
			}
		}
	}
}
