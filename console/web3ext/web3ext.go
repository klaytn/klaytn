// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from internal/web3ext/web3ext.go (2018/06/04).
// Modified and improved for the klaytn development.

package web3ext

var Modules = map[string]string{
	"admin":            Admin_JS,
	"debug":            Debug_JS,
	"klay":             Klay_JS,
	"net":              Net_JS,
	"personal":         Personal_JS,
	"rpc":              RPC_JS,
	"txpool":           TxPool_JS,
	"istanbul":         Istanbul_JS,
	"mainbridge":       MainBridge_JS,
	"subbridge":        SubBridge_JS,
	"clique":           CliqueJs,
	"governance":       Governance_JS,
	"bootnode":         Bootnode_JS,
	"chaindatafetcher": ChainDataFetcher_JS,
	"eth":              Eth_JS,
}

const Eth_JS = `
web3._extend({
	property: 'eth',
	methods: [
		new web3._extend.Method({
			name: 'chainId',
			call: 'eth_chainId',
			params: 0
		}),
		new web3._extend.Method({
			name: 'sign',
			call: 'eth_sign',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Method({
			name: 'resend',
			call: 'eth_resend',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, web3._extend.utils.fromDecimal, web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Method({
			name: 'signTransaction',
			call: 'eth_signTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'estimateGas',
			call: 'eth_estimateGas',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputCallFormatter, web3._extend.formatters.inputBlockNumberFormatter],
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Method({
			name: 'submitTransaction',
			call: 'eth_submitTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'fillTransaction',
			call: 'eth_fillTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'getHeaderByNumber',
			call: 'eth_getHeaderByNumber',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getHeaderByHash',
			call: 'eth_getHeaderByHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getBlockByNumber',
			call: 'eth_getBlockByNumber',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, function (val) { return !!val; }]
		}),
		new web3._extend.Method({
			name: 'getBlockByHash',
			call: 'eth_getBlockByHash',
			params: 2,
			inputFormatter: [null, function (val) { return !!val; }]
		}),
		new web3._extend.Method({
			name: 'getRawTransaction',
			call: 'eth_getRawTransactionByHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getRawTransactionFromBlock',
			call: function(args) {
				return (web3._extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? 'eth_getRawTransactionByBlockHashAndIndex' : 'eth_getRawTransactionByBlockNumberAndIndex';
			},
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, web3._extend.utils.toHex]
		}),
		new web3._extend.Method({
			name: 'getProof',
			call: 'eth_getProof',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null, web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'createAccessList',
			call: 'eth_createAccessList',
			params: 2,
			inputFormatter: [null, web3._extend.formatters.inputBlockNumberFormatter],
		}),
		new web3._extend.Method({
			name: 'feeHistory',
			call: 'eth_feeHistory',
			params: 3,
			inputFormatter: [null, web3._extend.formatters.inputBlockNumberFormatter, null]
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'pendingTransactions',
			getter: 'eth_pendingTransactions',
			outputFormatter: function(txs) {
				var formatted = [];
				for (var i = 0; i < txs.length; i++) {
					formatted.push(web3._extend.formatters.outputTransactionFormatter(txs[i]));
					formatted[i].blockHash = null;
				}
				return formatted;
			}
		}),
		new web3._extend.Property({
			name: 'maxPriorityFeePerGas',
			getter: 'eth_maxPriorityFeePerGas',
			outputFormatter: web3._extend.utils.toBigNumber
		}),
	]
});
`

const ChainDataFetcher_JS = `
web3._extend({
	property: 'chaindatafetcher',
	methods: [
		new web3._extend.Method({
			name: 'startFetching',
			call: 'chaindatafetcher_startFetching',
			params: 0
		}),
		new web3._extend.Method({
			name: 'stopFetching',
			call: 'chaindatafetcher_stopFetching',
			params: 0
		}),
		new web3._extend.Method({
			name: 'startRangeFetching',
			call: 'chaindatafetcher_startRangeFetching',
			params: 3
		}),
		new web3._extend.Method({
			name: 'stopRangeFetching',
			call: 'chaindatafetcher_stopRangeFetching',
			params: 0
		}),
		new web3._extend.Method({
			name: 'readCheckpoint',
			call: 'chaindatafetcher_readCheckpoint',
			params: 0
		}),
		new web3._extend.Method({
			name: 'status',
			call: 'chaindatafetcher_status',
			params: 0
		}),
		new web3._extend.Method({
			name: 'writeCheckpoint',
			call: 'chaindatafetcher_writeCheckpoint',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getConfig',
			call: 'chaindatafetcher_getConfig',
			params: 0
		})
	],
	properties: []
});
`

const Bootnode_JS = `
web3._extend({
	property: 'bootnode',
	methods: [
		new web3._extend.Method({
			name: 'createUpdateNodeOnDB',
			call: 'bootnode_createUpdateNodeOnDB',
			params: 1
		}),
        new web3._extend.Method({
			name: 'createUpdateNodeOnTable',
			call: 'bootnode_createUpdateNodeOnTable',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getNodeFromDB',
			call: 'bootnode_getNodeFromDB',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getTableEntries',
			call: 'bootnode_getTableEntries'
		}),
		new web3._extend.Method({
			name: 'getTableReplacements',
			call: 'bootnode_getTableReplacements'
		}),
		new web3._extend.Method({
			name: 'deleteNodeFromDB',
			call: 'bootnode_deleteNodeFromDB',
			params: 1
		}),
        new web3._extend.Method({
			name: 'deleteNodeFromTable',
			call: 'bootnode_deleteNodeFromTable',
			params: 1
		}),
		new web3._extend.Method({
			name: 'name',
			call: 'bootnode_name',
			params: 0
		}),
		new web3._extend.Method({
			name: 'resolve',
			call: 'bootnode_resolve',
			params: 1
		}),
		new web3._extend.Method({
			name: 'lookup',
			call: 'bootnode_lookup',
			params: 1
		}),
		new web3._extend.Method({
			name: 'readRandomNodes',
			call: 'bootnode_readRandomNodes',
			params: 0
		}),
        new web3._extend.Method({
			name: 'getAuthorizedNodes',
			call: 'bootnode_getAuthorizedNodes',
			params: 0
		}),
        new web3._extend.Method({
			name: 'putAuthorizedNodes',
			call: 'bootnode_putAuthorizedNodes',
			params: 1
		}),
        new web3._extend.Method({
			name: 'deleteAuthorizedNodes',
			call: 'bootnode_deleteAuthorizedNodes',
			params: 1
		})
	],
	properties: []
});
`
const Governance_JS = `
web3._extend({
	property: 'governance',
	methods: [
		new web3._extend.Method({
			name: 'vote',
			call: 'governance_vote',
			params: 2
		}),
		new web3._extend.Method({
			name: 'itemsAt',
			call: 'governance_itemsAt',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'itemCacheFromDb',
			call: 'governance_itemCacheFromDb',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getStakingInfo',
			call: 'governance_getStakingInfo',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		})
	],
	properties: [
		new web3._extend.Property({
			name: 'showTally',
			getter: 'governance_showTally',
		}),
		new web3._extend.Property({
			name: 'totalVotingPower',
			getter: 'governance_totalVotingPower',
		}),
		new web3._extend.Property({
			name: 'myVotes',
			getter: 'governance_myVotes',
		}),
		new web3._extend.Property({
			name: 'myVotingPower',
			getter: 'governance_myVotingPower',
		}),
		new web3._extend.Property({
			name: 'chainConfig',
			getter: 'governance_chainConfig',
		}),	
		new web3._extend.Property({
			name: 'nodeAddress',
			getter: 'governance_nodeAddress',
		}),
		new web3._extend.Property({
			name: 'pendingChanges',
			getter: 'governance_pendingChanges',
		}),
		new web3._extend.Property({
			name: 'votes',
			getter: 'governance_votes',
		}),
		new web3._extend.Property({
			name: 'idxCache',
			getter: 'governance_idxCache',
		}),
		new web3._extend.Property({
			name: 'idxCacheFromDb',
			getter: 'governance_idxCacheFromDb',
		})
	]
});
`
const Admin_JS = `
web3._extend({
	property: 'admin',
	methods: [
		new web3._extend.Method({
			name: 'addPeer',
			call: 'admin_addPeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'removePeer',
			call: 'admin_removePeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'exportChain',
			call: 'admin_exportChain',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'importChain',
			call: 'admin_importChain',
			params: 1
		}),
		new web3._extend.Method({
			name: 'importChainFromString',
			call: 'admin_importChainFromString',
			params: 1
		}),
		new web3._extend.Method({
			name: 'sleepBlocks',
			call: 'admin_sleepBlocks',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startHTTP',
			call: 'admin_startHTTP',
			params: 5,
			inputFormatter: [null, null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'stopHTTP',
			call: 'admin_stopHTTP'
		}),
		// This method is deprecated.
		new web3._extend.Method({
			name: 'startRPC',
			call: 'admin_startRPC',
			params: 5,
			inputFormatter: [null, null, null, null, null]
		}),
		// This method is deprecated.
		new web3._extend.Method({
			name: 'stopRPC',
			call: 'admin_stopRPC'
		}),
		new web3._extend.Method({
			name: 'startWS',
			call: 'admin_startWS',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'stopWS',
			call: 'admin_stopWS'
		}),
		new web3._extend.Method({
			name: 'startStateMigration',
			call: 'admin_startStateMigration',
		}),
		new web3._extend.Method({
			name: 'stopStateMigration',
			call: 'admin_stopStateMigration',
		}),
		new web3._extend.Method({
			name: 'saveTrieNodeCacheToDisk',
			call: 'admin_saveTrieNodeCacheToDisk',
		}),
		new web3._extend.Method({
			name: 'setMaxSubscriptionPerWSConn',
			call: 'admin_setMaxSubscriptionPerWSConn',
			params: 1
		}),
		new web3._extend.Method({
			name: 'startSpamThrottler',
			call: 'admin_startSpamThrottler',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'stopSpamThrottler',
			call: 'admin_stopSpamThrottler',
		}),
		new web3._extend.Method({
			name: 'setSpamThrottlerWhiteList',
			call: 'admin_setSpamThrottlerWhiteList',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'getSpamThrottlerWhiteList',
			call: 'admin_getSpamThrottlerWhiteList',
		}),
		new web3._extend.Method({
			name: 'getSpamThrottlerThrottleList',
			call: 'admin_getSpamThrottlerThrottleList',
		}),
		new web3._extend.Method({
			name: 'getSpamThrottlerCandidateList',
			call: 'admin_getSpamThrottlerCandidateList',
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'nodeInfo',
			getter: 'admin_nodeInfo'
		}),
		new web3._extend.Property({
			name: 'peers',
			getter: 'admin_peers'
		}),
		new web3._extend.Property({
			name: 'datadir',
			getter: 'admin_datadir'
		}),
		new web3._extend.Property({
			name: 'stateMigrationStatus',
			getter: 'admin_stateMigrationStatus'
		}),
		new web3._extend.Property({
			name: 'spamThrottlerConfig',
			getter: 'admin_spamThrottlerConfig'
		}),
	]
});
`

const Debug_JS = `
web3._extend({
	property: 'debug',
	methods: [
		new web3._extend.Method({
			name: 'printBlock',
			call: 'debug_printBlock',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getBlockRlp',
			call: 'debug_getBlockRlp',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setHead',
			call: 'debug_setHead',
			params: 1
		}),
		new web3._extend.Method({
			name: 'dumpBlock',
			call: 'debug_dumpBlock',
			params: 1
		}),
		new web3._extend.Method({
			name: 'dumpStateTrie',
			call: 'debug_dumpStateTrie',
			params: 1
		}),
		new web3._extend.Method({
			name: 'startWarmUp',
			call: 'debug_startWarmUp',
		}),
		new web3._extend.Method({
			name: 'startContractWarmUp',
			call: 'debug_startContractWarmUp',
			params: 1
		}),
		new web3._extend.Method({
			name: 'stopWarmUp',
			call: 'debug_stopWarmUp',
		}),
		new web3._extend.Method({
			name: 'startCollectingTrieStats',
			call: 'debug_startCollectingTrieStats',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'chaindbProperty',
			call: 'debug_chaindbProperty',
			params: 1,
			outputFormatter: console.log
		}),
		new web3._extend.Method({
			name: 'chaindbCompact',
			call: 'debug_chaindbCompact',
		}),
		new web3._extend.Method({
			name: 'metrics',
			call: 'debug_metrics',
			params: 1
		}),
		new web3._extend.Method({
			name: 'verbosity',
			call: 'debug_verbosity',
			params: 1
		}),
		new web3._extend.Method({
			name: 'verbosityByName',
			call: 'debug_verbosityByName',
			params: 2
		}),
		new web3._extend.Method({
			name: 'verbosityByID',
			call: 'debug_verbosityByID',
			params: 2
		}),
		new web3._extend.Method({
			name: 'vmodule',
			call: 'debug_vmodule',
			params: 1
		}),
		new web3._extend.Method({
			name: 'backtraceAt',
			call: 'debug_backtraceAt',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'stacks',
			call: 'debug_stacks',
			params: 0,
			outputFormatter: console.log
		}),
		new web3._extend.Method({
			name: 'freeOSMemory',
			call: 'debug_freeOSMemory',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'setGCPercent',
			call: 'debug_setGCPercent',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'memStats',
			call: 'debug_memStats',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'gcStats',
			call: 'debug_gcStats',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'startPProf',
			call: 'debug_startPProf',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'stopPProf',
			call: 'debug_stopPProf',
			params: 0
		}),
		new web3._extend.Method({
			name: 'isPProfRunning',
			call: 'debug_isPProfRunning',
			params: 0
		}),
		new web3._extend.Method({
			name: 'cpuProfile',
			call: 'debug_cpuProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startCPUProfile',
			call: 'debug_startCPUProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'stopCPUProfile',
			call: 'debug_stopCPUProfile',
			params: 0
		}),
		new web3._extend.Method({
			name: 'goTrace',
			call: 'debug_goTrace',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startGoTrace',
			call: 'debug_startGoTrace',
			params: 1
		}),
		new web3._extend.Method({
			name: 'stopGoTrace',
			call: 'debug_stopGoTrace',
			params: 0
		}),
		new web3._extend.Method({
			name: 'blockProfile',
			call: 'debug_blockProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setBlockProfileRate',
			call: 'debug_setBlockProfileRate',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeBlockProfile',
			call: 'debug_writeBlockProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'mutexProfile',
			call: 'debug_mutexProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setMutexProfileRate',
			call: 'debug_setMutexProfileRate',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeMutexProfile',
			call: 'debug_writeMutexProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeMemProfile',
			call: 'debug_writeMemProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'traceBlock',
			call: 'debug_traceBlock',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockFromFile',
			call: 'debug_traceBlockFromFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBadBlock',
			call: 'debug_traceBadBlock',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'standardTraceBadBlockToFile',
			call: 'debug_standardTraceBadBlockToFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'standardTraceBlockToFile',
			call: 'debug_standardTraceBlockToFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockByNumber',
			call: 'debug_traceBlockByNumber',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockByHash',
			call: 'debug_traceBlockByHash',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceTransaction',
			call: 'debug_traceTransaction',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'preimage',
			call: 'debug_preimage',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getBadBlocks',
			call: 'debug_getBadBlocks',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'storageRangeAt',
			call: 'debug_storageRangeAt',
			params: 5,
		}),
		new web3._extend.Method({
			name: 'getModifiedAccountsByNumber',
			call: 'debug_getModifiedAccountsByNumber',
			params: 2,
			inputFormatter: [null, null],
		}),
		new web3._extend.Method({
			name: 'getModifiedAccountsByHash',
			call: 'debug_getModifiedAccountsByHash',
			params: 2,
			inputFormatter:[null, null],
		}),
		new web3._extend.Method({
			name: 'getModifiedStorageNodesByNumber',
			call: 'debug_getModifiedStorageNodesByNumber',
			params: 4,
			inputFormatter: [null, null, null, null],
		}),
		new web3._extend.Method({
			name: 'setVMLogTarget',
			call: 'debug_setVMLogTarget',
			params: 1
		}),
	],
	properties: []
});
`

const Klay_JS = `
var blockWithConsensusInfoCall = function (args) {
    return (web3._extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? "klay_getBlockWithConsensusInfoByHash" : "klay_getBlockWithConsensusInfoByNumber";
};

web3._extend({
	property: 'klay',
	methods: [
		new web3._extend.Method({
			name: 'clientVersion',
			call: 'klay_clientVersion',
		}),
		new web3._extend.Method({
			name: 'getBlockReceipts',
			call: 'klay_getBlockReceipts',
			params: 1,
			outputFormatter: function(receipts) {
				var formatted = [];
				for (var i = 0; i < receipts.length; i++) {
					formatted.push(web3._extend.formatters.outputTransactionReceiptFormatter(receipts[i]));
				}
				return formatted;
			}
		}),
		new web3._extend.Method({
			name: 'sign',
			call: 'klay_sign',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Method({
			name: 'resend',
			call: 'klay_resend',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, web3._extend.utils.fromDecimal, web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Method({
			name: 'signTransaction',
			call: 'klay_signTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'signTransactionAsFeePayer',
			call: 'klay_signTransactionAsFeePayer',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'sendTransactionAsFeePayer',
			call: 'klay_sendTransactionAsFeePayer',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'getCouncil',
			call: 'klay_getCouncil',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getCouncilSize',
			call: 'klay_getCouncilSize',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getCommittee',
			call: 'klay_getCommittee',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getCommitteeSize',
			call: 'klay_getCommitteeSize',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'gasPriceAt',
			call: 'klay_gasPriceAt',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter],
			outputFormatter: web3._extend.formatters.outputBigNumberFormatter
		}),
		new web3._extend.Method({
			name: 'accountCreated',
			call: 'klay_accountCreated'
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, web3._extend.formatters.inputDefaultBlockNumberFormatter],
		}),
		new web3._extend.Method({
			name: 'getAccount',
			call: 'klay_getAccount'
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, web3._extend.formatters.inputDefaultBlockNumberFormatter],
		}),
		new web3._extend.Method({
			name: 'getHeaderByNumber',
			call: 'klay_getHeaderByNumber',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getHeaderByHash',
			call: 'klay_getHeaderByHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getBlockWithConsensusInfo',
			call: blockWithConsensusInfoCall,
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getBlockWithConsensusInfoRange',
			call: 'klay_getBlockWithConsensusInfoByNumberRange',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'isContractAccount',
			call: 'klay_isContractAccount',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, web3._extend.formatters.inputDefaultBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'submitTransaction',
			call: 'klay_submitTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'getRawTransaction',
			call: 'klay_getRawTransactionByHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'estimateComputationCost',
			call: 'klay_estimateComputationCost',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputCallFormatter, web3._extend.formatters.inputDefaultBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getAccountKey',
			call: 'klay_getAccountKey',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, web3._extend.formatters.inputDefaultBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getRawTransactionFromBlock',
			call: function(args) {
				return (web3._extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? 'klay_getRawTransactionByBlockHashAndIndex' : 'klay_getRawTransactionByBlockNumberAndIndex';
			},
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, web3._extend.utils.toHex]
		}),
		new web3._extend.Method({
			name: 'isParallelDBWrite',
			call: 'klay_isParallelDBWrite',
		}),
		new web3._extend.Method({
			name: 'isSenderTxHashIndexingEnabled',
			call: 'klay_isSenderTxHashIndexingEnabled',
		}),
		new web3._extend.Method({
			name: 'getTransactionBySenderTxHash',
			call: 'klay_getTransactionBySenderTxHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getTransactionReceiptBySenderTxHash',
			call: 'klay_getTransactionReceiptBySenderTxHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getCypressCredit',
			call: 'klay_getCypressCredit',
		}),
		new web3._extend.Method({
			name: 'sha3',
			call: 'klay_sha3',
			params: 1,
			inputFormatter: [web3._extend.utils.toHex],
		}),
		new web3._extend.Method({
			name: 'encodeAccountKey',
			call: 'klay_encodeAccountKey',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'decodeAccountKey',
			call: 'klay_decodeAccountKey',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'createAccessList',
			call: 'klay_createAccessList',
			params: 2,
			inputFormatter: [null, web3._extend.formatters.inputBlockNumberFormatter],
		}),
		new web3._extend.Method({
			name: 'feeHistory',
			call: 'klay_feeHistory',
			params: 3,
			inputFormatter: [null, web3._extend.formatters.inputBlockNumberFormatter, null]
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'pendingTransactions',
			getter: 'klay_pendingTransactions',
			outputFormatter: function(txs) {
				var formatted = [];
				for (var i = 0; i < txs.length; i++) {
					formatted.push(web3._extend.formatters.outputTransactionFormatter(txs[i]));
					formatted[i].blockHash = null;
				}
				return formatted;
			}
		}),
        new web3._extend.Property({
            name : 'rewardbase',
            getter: 'klay_rewardbase'
        }),
        new web3._extend.Property({
            name : 'gasPrice',
            getter: 'klay_gasPrice',
            outputFormatter: web3._extend.formatters.outputBigNumberFormatter
        }),
		new web3._extend.Property({
			name: 'maxPriorityFeePerGas',
			getter: 'klay_maxPriorityFeePerGas',
			outputFormatter: web3._extend.utils.toBigNumber
		}),
	]
});
`

const Net_JS = `
web3._extend({
	property: 'net',
	methods: [
		new web3._extend.Method({
			name: 'peerCountByType',
			call: 'net_peerCountByType',
			params: 0,
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'version',
			getter: 'net_version'
		}),
		new web3._extend.Property({
			name: 'networkID',
			getter: 'net_networkID'
		}),
	]
});
`

const Personal_JS = `
web3._extend({
	property: 'personal',
	methods: [
		new web3._extend.Method({
			name: 'importRawKey',
			call: 'personal_importRawKey',
			params: 2
		}),
		new web3._extend.Method({
			name: 'replaceRawKey',
			call: 'personal_replaceRawKey',
			params: 3
		}),
		new web3._extend.Method({
			name: 'sign',
			call: 'personal_sign',
			params: 3,
			inputFormatter: [null, web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Method({
			name: 'ecRecover',
			call: 'personal_ecRecover',
			params: 2
		}),
		new web3._extend.Method({
			name: 'openWallet',
			call: 'personal_openWallet',
			params: 2
		}),
		new web3._extend.Method({
			name: 'deriveAccount',
			call: 'personal_deriveAccount',
			params: 3
		}),
		new web3._extend.Method({
			name: 'sendValueTransfer',
			call: 'personal_sendValueTransfer',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, null]
		}),
		new web3._extend.Method({
			name: 'sendAccountUpdate',
			call: 'personal_sendAccountUpdate',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, null]
		}),
		new web3._extend.Method({
			name: 'signTransaction',
			call: 'personal_signTransaction',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, null]
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'listWallets',
			getter: 'personal_listWallets'
		}),
	]
})
`

const RPC_JS = `
web3._extend({
	property: 'rpc',
	methods: [],
	properties: [
		new web3._extend.Property({
			name: 'modules',
			getter: 'rpc_modules'
		}),
	]
});
`

const TxPool_JS = `
web3._extend({
	property: 'txpool',
	methods: [],
	properties:
	[
		new web3._extend.Property({
			name: 'content',
			getter: 'txpool_content'
		}),
		new web3._extend.Property({
			name: 'inspect',
			getter: 'txpool_inspect'
		}),
		new web3._extend.Property({
			name: 'status',
			getter: 'txpool_status',
			outputFormatter: function(status) {
				status.pending = web3._extend.utils.toDecimal(status.pending);
				status.queued = web3._extend.utils.toDecimal(status.queued);
				return status;
			}
		}),
	]
});
`

const Istanbul_JS = `
web3._extend({
	property: 'istanbul',
	methods:
	[
		new web3._extend.Method({
			name: 'getSnapshot',
			call: 'istanbul_getSnapshot',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSnapshotAtHash',
			call: 'istanbul_getSnapshotAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getValidators',
			call: 'istanbul_getValidators',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getValidatorsAtHash',
			call: 'istanbul_getValidatorsAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getDemotedValidators',
			call: 'istanbul_getDemotedValidators',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getDemotedValidatorsAtHash',
			call: 'istanbul_getDemotedValidatorsAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'discard',
			call: 'istanbul_discard',
			params: 1
		})
	],
	properties:
	[
		new web3._extend.Property({
			name: 'candidates',
			getter: 'istanbul_candidates'
		}),
		new web3._extend.Property({
			name: 'timeout',
			getter: 'istanbul_getTimeout'
		})
	]
});
`
const MainBridge_JS = `
web3._extend({
	property: 'mainbridge',
	methods:
	[
		new web3._extend.Method({
			name: 'getChildChainIndexingEnabled',
			call: 'mainbridge_getChildChainIndexingEnabled'
		}),
		new web3._extend.Method({
			name: 'convertChildChainBlockHashToParentChainTxHash',
			call: 'mainbridge_convertChildChainBlockHashToParentChainTxHash',
			params: 1
		}),
	],
    properties: [
		new web3._extend.Property({
			name: 'nodeInfo',
			getter: 'mainbridge_nodeInfo'
		}),
		new web3._extend.Property({
			name: 'peers',
			getter: 'mainbridge_peers'
		}),
	]
});
`
const SubBridge_JS = `
web3._extend({
	property: 'subbridge',
	methods:
	[
		new web3._extend.Method({
			name: 'addPeer',
			call: 'subbridge_addPeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'removePeer',
			call: 'subbridge_removePeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getReceiptFromParentChain',
			call: 'subbridge_getReceiptFromParentChain',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getAnchoringTxHashByBlockNumber',
			call: 'subbridge_getAnchoringTxHashByBlockNumber',
			params: 1
		}),
		new web3._extend.Method({
			name: 'registerOperator',
			call: 'subbridge_registerOperator',
			params: 2
		}),
		new web3._extend.Method({
			name: 'getOperators',
			call: 'subbridge_getRegisteredOperators',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getValueTransferOperatorThreshold',
			call: 'subbridge_getValueTransferOperatorThreshold',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setValueTransferOperatorThreshold',
			call: 'subbridge_setValueTransferOperatorThreshold',
			params: 2
		}),
		new web3._extend.Method({
			name: 'deployBridge',
			call: 'subbridge_deployBridge',
			params: 0
		}),
		new web3._extend.Method({
			name: 'subscribeBridge',
			call: 'subbridge_subscribeBridge',
			params: 2
		}),
		new web3._extend.Method({
			name: 'unsubscribeBridge',
			call: 'subbridge_unsubscribeBridge',
			params: 2
		}),
		new web3._extend.Method({
			name: 'KASAnchor',
			call: 'subbridge_kASAnchor',
			params: 1
		}),
		new web3._extend.Method({
			name: 'anchoring',
			call: 'subbridge_anchoring',
			params: 1
		}),
		new web3._extend.Method({
			name: 'registerBridge',
			call: 'subbridge_registerBridge',
			params: 2
		}),
		new web3._extend.Method({
			name: 'deregisterBridge',
			call: 'subbridge_deregisterBridge',
			params: 2
		}),
		new web3._extend.Method({
			name: 'registerToken',
			call: 'subbridge_registerToken',
			params: 4
		}),
		new web3._extend.Method({
			name: 'deregisterToken',
			call: 'subbridge_deregisterToken',
			params: 4
		}),
		new web3._extend.Method({
			name: 'convertRequestTxHashToHandleTxHash',
			call: 'subbridge_convertRequestTxHashToHandleTxHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getBridgeInformation',
			call: 'subbridge_getBridgeInformation',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getParentTransactionReceipt',
			call: 'subbridge_getParentTransactionReceipt',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setKLAYFee',
			call: 'subbridge_setKLAYFee',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setERC20Fee',
			call: 'subbridge_setERC20Fee',
			params: 3
		}),
		new web3._extend.Method({
			name: 'setFeeReceiver',
			call: 'subbridge_setFeeReceiver',
			params: 2
		}),
		new web3._extend.Method({
			name: 'getKLAYFee',
			call: 'subbridge_getKLAYFee',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getERC20Fee',
			call: 'subbridge_getERC20Fee',
			params: 2
		}),
		new web3._extend.Method({
			name: 'getFeeReceiver',
			call: 'subbridge_getFeeReceiver',
			params: 1
		}),
		new web3._extend.Method({
			name: 'lockParentOperator',
			call: 'subbridge_lockParentOperator'
		}),
		new web3._extend.Method({
			name: 'lockChildOperator',
			call: 'subbridge_lockChildOperator'
		}),
		new web3._extend.Method({
			name: 'unlockParentOperator',
			call: 'subbridge_unlockParentOperator',
			params: 2
		}),
		new web3._extend.Method({
			name: 'unlockChildOperator',
			call: 'subbridge_unlockChildOperator',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setParentOperatorFeePayer',
			call: 'subbridge_setParentOperatorFeePayer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setChildOperatorFeePayer',
			call: 'subbridge_setChildOperatorFeePayer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setParentBridgeOperatorGasLimit',
			call: 'subbridge_setParentBridgeOperatorGasLimit',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setChildBridgeOperatorGasLimit',
			call: 'subbridge_setChildBridgeOperatorGasLimit',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getParentBridgeContractBalance',
			call: 'subbridge_getParentBridgeContractBalance',
			params: 1,
			outputFormatter: web3._extend.formatters.outputBigNumberFormatter
		}),
		new web3._extend.Method({
			name: 'getChildBridgeContractBalance',
			call: 'subbridge_getChildBridgeContractBalance',
			params: 1,
			outputFormatter: web3._extend.formatters.outputBigNumberFormatter
		}),
	],
    properties: [
		new web3._extend.Property({
			name: 'nodeInfo',
			getter: 'subbridge_nodeInfo'
		}),
		new web3._extend.Property({
			name: 'peers',
			getter: 'subbridge_peers'
		}),
		new web3._extend.Property({
			name: 'parentOperator',
			getter: 'subbridge_getParentOperatorAddr'
		}),
			new web3._extend.Property({
			name: 'childOperator',
			getter: 'subbridge_getChildOperatorAddr'
		}),
		new web3._extend.Property({
			name: 'operators',
			getter: 'subbridge_getOperators'
		}),
		new web3._extend.Property({
			name: 'anchoringPeriod',
			getter: 'subbridge_getAnchoringPeriod'
		}),
		new web3._extend.Property({
			name: 'sendChainTxslimit',
			getter: 'subbridge_getSentChainTxsLimit'
		}),
		new web3._extend.Property({
			name: 'parentOperatorNonce',
			getter: 'subbridge_getParentOperatorNonce'
		}),
		new web3._extend.Property({
			name: 'childOperatorNonce',
			getter: 'subbridge_getChildOperatorNonce'
		}),
		new web3._extend.Property({
			name: 'parentOperatorBalance',
			getter: 'subbridge_getParentOperatorBalance',
			outputFormatter: web3._extend.formatters.outputBigNumberFormatter
		}),
		new web3._extend.Property({
			name: 'childOperatorBalance',
			getter: 'subbridge_getChildOperatorBalance',
			outputFormatter: web3._extend.formatters.outputBigNumberFormatter
		}),
		new web3._extend.Property({
			name: 'listBridge',
			getter: 'subbridge_listBridge'
		}),
		new web3._extend.Property({
			name: 'txPendingCount',
			getter: 'subbridge_txPendingCount'
		}),
		new web3._extend.Property({
			name: 'txPending',
			getter: 'subbridge_txPending'
		}),
		new web3._extend.Property({
			name: 'latestAnchoredBlockNumber',
			getter: 'subbridge_getLatestAnchoredBlockNumber'
		}),
		new web3._extend.Property({
			name: 'parentOperatorFeePayer',
			getter: 'subbridge_getParentOperatorFeePayer',
		}),
		new web3._extend.Property({
			name: 'childOperatorFeePayer',
			getter: 'subbridge_getChildOperatorFeePayer',
		}),
		new web3._extend.Property({
			name: 'parentBridgeOperatorGasLimit',
			getter: 'subbridge_getParentBridgeOperatorGasLimit',
		}),
		new web3._extend.Property({
			name: 'childBridgeOperatorGasLimit',
			getter: 'subbridge_getChildBridgeOperatorGasLimit',
		}),
	]
});
`
const CliqueJs = `
web3._extend({
	property: 'clique',
	methods: [
		new web3._extend.Method({
			name: 'getSnapshot',
			call: 'clique_getSnapshot',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSnapshotAtHash',
			call: 'clique_getSnapshotAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getSigners',
			call: 'clique_getSigners',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSignersAtHash',
			call: 'clique_getSignersAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'propose',
			call: 'clique_propose',
			params: 2
		}),
		new web3._extend.Method({
			name: 'discard',
			call: 'clique_discard',
			params: 1
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'proposals',
			getter: 'clique_proposals'
		}),
	]
});
`
