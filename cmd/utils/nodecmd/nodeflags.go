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
// This file is derived from cmd/geth/main.go (2018/06/04).
// Modified and improved for the klaytn development.

package nodecmd

import (
	"github.com/klaytn/klaytn/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

// Common flags that configure the node
var CommonNodeFlags = []cli.Flag{
	utils.BootnodesFlag,
	utils.IdentityFlag,
	utils.UnlockedAccountFlag,
	utils.PasswordFileFlag,
	utils.DbTypeFlag,
	utils.DataDirFlag,
	utils.KeyStoreDirFlag,
	utils.TxPoolNoLocalsFlag,
	utils.TxPoolJournalFlag,
	utils.TxPoolJournalIntervalFlag,
	utils.TxPoolPriceLimitFlag,
	utils.TxPoolPriceBumpFlag,
	utils.TxPoolExecSlotsAccountFlag,
	utils.TxPoolExecSlotsAllFlag,
	utils.TxPoolNonExecSlotsAccountFlag,
	utils.TxPoolNonExecSlotsAllFlag,
	utils.TxPoolLifetimeFlag,
	utils.TxPoolKeepLocalsFlag,
	utils.SyncModeFlag,
	utils.GCModeFlag,
	utils.LightKDFFlag,
	utils.StateDBCachingFlag,
	utils.NoPartitionedDBFlag,
	utils.NumStateTriePartitionsFlag,
	utils.LevelDBCompressionTypeFlag,
	utils.LevelDBNoBufferPoolFlag,
	utils.LevelDBCacheSizeFlag,
	utils.NoParallelDBWriteFlag,
	utils.DataArchivingBlockNumFlag,
	utils.SenderTxHashIndexingFlag,
	utils.TrieMemoryCacheSizeFlag,
	utils.TrieBlockIntervalFlag,
	utils.CacheTypeFlag,
	utils.CacheScaleFlag,
	utils.CacheUsageLevelFlag,
	utils.MemorySizeFlag,
	utils.CacheWriteThroughFlag,
	utils.TxPoolStateCacheFlag,
	utils.TrieCacheLimitFlag,
	utils.ListenPortFlag,
	utils.SubListenPortFlag,
	utils.MultiChannelUseFlag,
	utils.MaxConnectionsFlag,
	utils.MaxPendingPeersFlag,
	utils.TargetGasLimitFlag,
	utils.NATFlag,
	utils.NoDiscoverFlag,
	utils.RWTimerWaitTimeFlag,
	utils.RWTimerIntervalFlag,
	utils.NetrestrictFlag,
	utils.NodeKeyFileFlag,
	utils.NodeKeyHexFlag,
	utils.VMEnableDebugFlag,
	utils.VMLogTargetFlag,
	utils.NetworkIdFlag,
	utils.RPCCORSDomainFlag,
	utils.RPCVirtualHostsFlag,
	utils.MetricsEnabledFlag,
	utils.PrometheusExporterFlag,
	utils.PrometheusExporterPortFlag,
	utils.ExtraDataFlag,
	utils.SrvTypeFlag,
	utils.AutoRestartFlag,
	utils.RestartTimeOutFlag,
	utils.DaemonPathFlag,
	ConfigFileFlag,
}

// Common RPC flags
var CommonRPCFlags = []cli.Flag{
	utils.RPCEnabledFlag,
	utils.RPCListenAddrFlag,
	utils.RPCPortFlag,
	utils.RPCApiFlag,
	utils.WSEnabledFlag,
	utils.WSListenAddrFlag,
	utils.WSPortFlag,
	utils.GRPCEnabledFlag,
	utils.GRPCListenAddrFlag,
	utils.GRPCPortFlag,
	utils.WSApiFlag,
	utils.WSAllowedOriginsFlag,
	utils.IPCDisabledFlag,
	utils.IPCPathFlag,
}

var KCNFlags = []cli.Flag{
	utils.RewardbaseFlag,
	utils.CypressFlag,
	utils.BaobabFlag,
}

var KPNFlags = []cli.Flag{
	utils.TxResendIntervalFlag,
	utils.TxResendCountFlag,
	utils.TxResendUseLegacyFlag,
	utils.CypressFlag,
	utils.BaobabFlag,
}

var KENFlags = []cli.Flag{
	utils.ServiceChainSignerFlag,
	utils.CypressFlag,
	utils.BaobabFlag,
	utils.ChildChainIndexingFlag,
	utils.MainBridgeFlag,
	utils.MainBridgeListenPortFlag,
	// DBSyncer
	utils.EnableDBSyncerFlag,
	utils.DBHostFlag,
	utils.DBPortFlag,
	utils.DBNameFlag,
	utils.DBUserFlag,
	utils.DBPasswordFlag,
	utils.EnabledLogModeFlag,
	utils.MaxIdleConnsFlag,
	utils.MaxOpenConnsFlag,
	utils.ConnMaxLifeTimeFlag,
	utils.BlockSyncChannelSizeFlag,
	utils.DBSyncerModeFlag,
	utils.GenQueryThreadFlag,
	utils.InsertThreadFlag,
	utils.BulkInsertSizeFlag,
	utils.EventModeFlag,
	utils.MaxBlockDiffFlag,
	utils.TxResendIntervalFlag,
	utils.TxResendCountFlag,
	utils.TxResendUseLegacyFlag,
}

var KSCNFlags = []cli.Flag{
	utils.ServiceChainSignerFlag,
	utils.AnchoringPeriodFlag,
	utils.SentChainTxsLimit,
	utils.MainBridgeFlag,
	utils.MainBridgeListenPortFlag,
	utils.ChildChainIndexingFlag,
	utils.SubBridgeFlag,
	utils.SubBridgeListenPortFlag,
	utils.ParentChainIDFlag,
	utils.VTRecoveryFlag,
	utils.VTRecoveryIntervalFlag,
	utils.ServiceChainNewAccountFlag,
	utils.ServiceChainAnchoringFlag,
}

var KSPNFlags = []cli.Flag{
	utils.TxResendIntervalFlag,
	utils.TxResendCountFlag,
	utils.TxResendUseLegacyFlag,
	utils.ServiceChainSignerFlag,
	utils.AnchoringPeriodFlag,
	utils.SentChainTxsLimit,
	utils.MainBridgeFlag,
	utils.MainBridgeListenPortFlag,
	utils.ChildChainIndexingFlag,
	utils.SubBridgeFlag,
	utils.SubBridgeListenPortFlag,
	utils.ParentChainIDFlag,
	utils.VTRecoveryFlag,
	utils.VTRecoveryIntervalFlag,
	utils.ServiceChainNewAccountFlag,
	utils.ServiceChainAnchoringFlag,
}

var KSENFlags = []cli.Flag{
	utils.ServiceChainSignerFlag,
	utils.ChildChainIndexingFlag,
	utils.MainBridgeFlag,
	utils.MainBridgeListenPortFlag,
	utils.SubBridgeFlag,
	utils.SubBridgeListenPortFlag,
	utils.AnchoringPeriodFlag,
	utils.SentChainTxsLimit,
	utils.ParentChainIDFlag,
	utils.VTRecoveryFlag,
	utils.VTRecoveryIntervalFlag,
	utils.ServiceChainAnchoringFlag,
	// DBSyncer
	utils.EnableDBSyncerFlag,
	utils.DBHostFlag,
	utils.DBPortFlag,
	utils.DBNameFlag,
	utils.DBUserFlag,
	utils.DBPasswordFlag,
	utils.EnabledLogModeFlag,
	utils.MaxIdleConnsFlag,
	utils.MaxOpenConnsFlag,
	utils.ConnMaxLifeTimeFlag,
	utils.BlockSyncChannelSizeFlag,
	utils.DBSyncerModeFlag,
	utils.GenQueryThreadFlag,
	utils.InsertThreadFlag,
	utils.BulkInsertSizeFlag,
	utils.EventModeFlag,
	utils.MaxBlockDiffFlag,
	utils.TxResendIntervalFlag,
	utils.TxResendCountFlag,
	utils.TxResendUseLegacyFlag,
}
