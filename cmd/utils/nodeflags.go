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

package utils

import (
	"github.com/klaytn/klaytn/api/debug"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// TODO-Klaytn: Check whether all flags are registered in utils.FlagGroups

func AllNodeFlags() []cli.Flag {
	nodeFlags := []cli.Flag{}
	nodeFlags = append(nodeFlags, CommonNodeFlags...)
	nodeFlags = append(nodeFlags, CommonRPCFlags...)
	nodeFlags = append(nodeFlags, ConsoleFlags...)
	nodeFlags = append(nodeFlags, debug.Flags...)
	nodeFlags = append(nodeFlags, ChainDataFetcherFlags...)
	nodeFlags = union(nodeFlags, SnapshotFlags)
	nodeFlags = union(nodeFlags, DBMigrationSrcFlags)
	nodeFlags = union(nodeFlags, DBMigrationDstFlags)
	nodeFlags = union(nodeFlags, BNFlags)
	nodeFlags = union(nodeFlags, KCNFlags)
	nodeFlags = union(nodeFlags, KPNFlags)
	nodeFlags = union(nodeFlags, KENFlags)
	nodeFlags = union(nodeFlags, KSCNFlags)
	nodeFlags = union(nodeFlags, KSPNFlags)
	nodeFlags = union(nodeFlags, KSENFlags)
	return nodeFlags
}

func contains(list []cli.Flag, item cli.Flag) bool {
	for _, flag := range list {
		if flag.Names()[0] == item.Names()[0] {
			return true
		}
	}
	return false
}

func union(list1, list2 []cli.Flag) []cli.Flag {
	for _, item := range list2 {
		if !contains(list1, item) {
			list1 = append(list1, item)
		}
	}
	return list1
}

// All flags used for each node type
func KcnNodeFlags() []cli.Flag {
	return append(CommonNodeFlags, KCNFlags...)
}

func KpnNodeFlags() []cli.Flag {
	return append(CommonNodeFlags, KPNFlags...)
}

func KenNodeFlags() []cli.Flag {
	return append(CommonNodeFlags, KENFlags...)
}

func KscnNodeFlags() []cli.Flag {
	return append(CommonNodeFlags, KSCNFlags...)
}

func KspnNodeFlags() []cli.Flag {
	return append(CommonNodeFlags, KSPNFlags...)
}

func KsenNodeFlags() []cli.Flag {
	return append(CommonNodeFlags, KSENFlags...)
}

func BNAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, BNFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, CommonRPCFlags...)
	return flags
}

func KcnAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KcnNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, DBMigrationDstFlags...)
	return flags
}

func KpnAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KpnNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, DBMigrationDstFlags...)
	flags = append(flags, ChainDataFetcherFlags...)
	return flags
}

func KenAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KenNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, DBMigrationDstFlags...)
	flags = append(flags, ChainDataFetcherFlags...)
	return flags
}

func KscnAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KscnNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	return flags
}

func KspnAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KspnNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, ChainDataFetcherFlags...)
	return flags
}

func KsenAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KsenNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, ChainDataFetcherFlags...)
	return flags
}

var ConsoleFlags = []cli.Flag{
	altsrc.NewStringFlag(JSpathFlag),
	altsrc.NewStringFlag(ExecFlag),
	altsrc.NewStringFlag(PreloadJSFlag),
}

// Common flags that configure the node
var CommonNodeFlags = []cli.Flag{
	ConfFlag,
	altsrc.NewBoolFlag(NtpDisableFlag),
	altsrc.NewStringFlag(NtpServerFlag),
	altsrc.NewPathFlag(DocRootFlag),
	altsrc.NewStringFlag(BootnodesFlag),
	altsrc.NewStringFlag(IdentityFlag),
	altsrc.NewStringFlag(UnlockedAccountFlag),
	altsrc.NewStringFlag(PasswordFileFlag),
	altsrc.NewStringFlag(DbTypeFlag),
	altsrc.NewPathFlag(DataDirFlag),
	altsrc.NewPathFlag(ChainDataDirFlag),
	altsrc.NewBoolFlag(OverwriteGenesisFlag),
	altsrc.NewUint64Flag(StartBlockNumberFlag),
	altsrc.NewPathFlag(KeyStoreDirFlag),
	altsrc.NewBoolFlag(TxPoolNoLocalsFlag),
	altsrc.NewBoolFlag(TxPoolAllowLocalAnchorTxFlag),
	altsrc.NewBoolFlag(TxPoolDenyRemoteTxFlag),
	altsrc.NewStringFlag(TxPoolJournalFlag),
	altsrc.NewDurationFlag(TxPoolJournalIntervalFlag),
	altsrc.NewUint64Flag(TxPoolPriceLimitFlag),
	altsrc.NewUint64Flag(TxPoolPriceBumpFlag),
	altsrc.NewUint64Flag(TxPoolExecSlotsAccountFlag),
	altsrc.NewUint64Flag(TxPoolExecSlotsAllFlag),
	altsrc.NewUint64Flag(TxPoolNonExecSlotsAccountFlag),
	altsrc.NewUint64Flag(TxPoolNonExecSlotsAllFlag),
	altsrc.NewDurationFlag(TxPoolLifetimeFlag),
	altsrc.NewBoolFlag(TxPoolKeepLocalsFlag),
	NewWrappedTextMarshalerFlag(SyncModeFlag),
	altsrc.NewStringFlag(GCModeFlag),
	altsrc.NewBoolFlag(LightKDFFlag),
	altsrc.NewBoolFlag(SingleDBFlag),
	altsrc.NewUintFlag(NumStateTrieShardsFlag),
	altsrc.NewIntFlag(LevelDBCompressionTypeFlag),
	altsrc.NewBoolFlag(LevelDBNoBufferPoolFlag),
	altsrc.NewBoolFlag(DBNoPerformanceMetricsFlag),
	altsrc.NewBoolFlag(RocksDBSecondaryFlag),
	altsrc.NewUint64Flag(RocksDBCacheSizeFlag),
	altsrc.NewBoolFlag(RocksDBDumpMallocStatFlag),
	altsrc.NewStringFlag(RocksDBCompressionTypeFlag),
	altsrc.NewStringFlag(RocksDBBottommostCompressionTypeFlag),
	altsrc.NewStringFlag(RocksDBFilterPolicyFlag),
	altsrc.NewBoolFlag(RocksDBDisableMetricsFlag),
	altsrc.NewIntFlag(RocksDBMaxOpenFilesFlag),
	altsrc.NewBoolFlag(RocksDBCacheIndexAndFilterFlag),
	altsrc.NewStringFlag(DynamoDBTableNameFlag),
	altsrc.NewStringFlag(DynamoDBRegionFlag),
	altsrc.NewBoolFlag(DynamoDBIsProvisionedFlag),
	altsrc.NewInt64Flag(DynamoDBReadCapacityFlag),
	altsrc.NewInt64Flag(DynamoDBWriteCapacityFlag),
	altsrc.NewBoolFlag(DynamoDBReadOnlyFlag),
	altsrc.NewIntFlag(LevelDBCacheSizeFlag),
	altsrc.NewBoolFlag(NoParallelDBWriteFlag),
	altsrc.NewBoolFlag(SenderTxHashIndexingFlag),
	altsrc.NewIntFlag(TrieMemoryCacheSizeFlag),
	altsrc.NewUintFlag(TrieBlockIntervalFlag),
	altsrc.NewUint64Flag(TriesInMemoryFlag),
	altsrc.NewBoolFlag(LivePruningFlag),
	altsrc.NewUint64Flag(LivePruningRetentionFlag),
	altsrc.NewIntFlag(CacheTypeFlag),
	altsrc.NewIntFlag(CacheScaleFlag),
	altsrc.NewStringFlag(CacheUsageLevelFlag),
	altsrc.NewIntFlag(MemorySizeFlag),
	altsrc.NewStringFlag(TrieNodeCacheTypeFlag),
	altsrc.NewIntFlag(NumFetcherPrefetchWorkerFlag),
	altsrc.NewBoolFlag(UseSnapshotForPrefetchFlag),
	altsrc.NewIntFlag(TrieNodeCacheLimitFlag),
	altsrc.NewDurationFlag(TrieNodeCacheSavePeriodFlag),
	altsrc.NewStringSliceFlag(TrieNodeCacheRedisEndpointsFlag),
	altsrc.NewBoolFlag(TrieNodeCacheRedisClusterFlag),
	altsrc.NewBoolFlag(TrieNodeCacheRedisPublishBlockFlag),
	altsrc.NewBoolFlag(TrieNodeCacheRedisSubscribeBlockFlag),
	altsrc.NewIntFlag(ListenPortFlag),
	altsrc.NewIntFlag(SubListenPortFlag),
	altsrc.NewBoolFlag(MultiChannelUseFlag),
	altsrc.NewIntFlag(MaxConnectionsFlag),
	altsrc.NewIntFlag(MaxRequestContentLengthFlag),
	altsrc.NewIntFlag(MaxPendingPeersFlag),
	altsrc.NewUint64Flag(TargetGasLimitFlag),
	altsrc.NewStringFlag(NATFlag),
	altsrc.NewBoolFlag(NoDiscoverFlag),
	altsrc.NewDurationFlag(RWTimerWaitTimeFlag),
	altsrc.NewUint64Flag(RWTimerIntervalFlag),
	altsrc.NewStringFlag(NetrestrictFlag),
	altsrc.NewStringFlag(NodeKeyFileFlag),
	altsrc.NewStringFlag(NodeKeyHexFlag),
	altsrc.NewBoolFlag(VMEnableDebugFlag),
	altsrc.NewIntFlag(VMLogTargetFlag),
	altsrc.NewBoolFlag(VMTraceInternalTxFlag),
	altsrc.NewUint64Flag(NetworkIdFlag),
	altsrc.NewBoolFlag(MetricsEnabledFlag),
	altsrc.NewBoolFlag(PrometheusExporterFlag),
	altsrc.NewIntFlag(PrometheusExporterPortFlag),
	altsrc.NewStringFlag(ExtraDataFlag),
	altsrc.NewStringFlag(SrvTypeFlag),
	altsrc.NewBoolFlag(AutoRestartFlag),
	altsrc.NewDurationFlag(RestartTimeOutFlag),
	altsrc.NewStringFlag(DaemonPathFlag),
	altsrc.NewStringFlag(ConfigFileFlag),
	altsrc.NewIntFlag(APIFilterGetLogsMaxItemsFlag),
	altsrc.NewDurationFlag(APIFilterGetLogsDeadlineFlag),
	altsrc.NewUint64Flag(OpcodeComputationCostLimitFlag),
	altsrc.NewBoolFlag(SnapshotFlag),
	altsrc.NewIntFlag(SnapshotCacheSizeFlag),
	altsrc.NewBoolFlag(SnapshotAsyncGen),
}

// Common RPC flags
var CommonRPCFlags = []cli.Flag{
	altsrc.NewBoolFlag(RPCEnabledFlag),
	altsrc.NewStringFlag(RPCListenAddrFlag),
	altsrc.NewIntFlag(RPCPortFlag),
	altsrc.NewStringFlag(RPCApiFlag),
	altsrc.NewUint64Flag(RPCGlobalGasCap),
	altsrc.NewFloat64Flag(RPCGlobalEthTxFeeCapFlag),
	altsrc.NewStringFlag(RPCCORSDomainFlag),
	altsrc.NewStringFlag(RPCVirtualHostsFlag),
	altsrc.NewBoolFlag(RPCNonEthCompatibleFlag),
	altsrc.NewDurationFlag(RPCGlobalEVMTimeoutFlag),
	altsrc.NewBoolFlag(WSEnabledFlag),
	altsrc.NewStringFlag(WSListenAddrFlag),
	altsrc.NewIntFlag(WSPortFlag),
	altsrc.NewBoolFlag(GRPCEnabledFlag),
	altsrc.NewStringFlag(GRPCListenAddrFlag),
	altsrc.NewIntFlag(GRPCPortFlag),
	altsrc.NewIntFlag(RPCConcurrencyLimit),
	altsrc.NewStringFlag(WSApiFlag),
	altsrc.NewStringFlag(WSAllowedOriginsFlag),
	altsrc.NewIntFlag(WSMaxSubscriptionPerConn),
	altsrc.NewInt64Flag(WSReadDeadLine),
	altsrc.NewInt64Flag(WSWriteDeadLine),
	altsrc.NewIntFlag(WSMaxConnections),
	altsrc.NewBoolFlag(IPCDisabledFlag),
	altsrc.NewPathFlag(IPCPathFlag),
	altsrc.NewIntFlag(RPCReadTimeout),
	altsrc.NewIntFlag(RPCWriteTimeoutFlag),
	altsrc.NewIntFlag(RPCIdleTimeoutFlag),
	altsrc.NewIntFlag(RPCExecutionTimeoutFlag),
	altsrc.NewBoolFlag(UnsafeDebugDisableFlag),
	altsrc.NewIntFlag(HeavyDebugRequestLimitFlag),
	altsrc.NewDurationFlag(StateRegenerationTimeLimitFlag),
	altsrc.NewStringFlag(RPCUpstreamArchiveENFlag),
}

var BNFlags = []cli.Flag{
	altsrc.NewStringFlag(SrvTypeFlag),
	altsrc.NewPathFlag(DataDirFlag),
	altsrc.NewStringFlag(GenKeyFlag),
	altsrc.NewStringFlag(NodeKeyFileFlag),
	altsrc.NewStringFlag(NodeKeyHexFlag),
	altsrc.NewBoolFlag(WriteAddressFlag),
	altsrc.NewStringFlag(BNAddrFlag),
	altsrc.NewStringFlag(NATFlag),
	altsrc.NewStringFlag(NetrestrictFlag),
	altsrc.NewBoolFlag(MetricsEnabledFlag),
	altsrc.NewBoolFlag(PrometheusExporterFlag),
	altsrc.NewIntFlag(PrometheusExporterPortFlag),
	altsrc.NewStringFlag(AuthorizedNodesFlag),
	altsrc.NewUint64Flag(NetworkIdFlag),
}

var KCNFlags = []cli.Flag{
	altsrc.NewStringFlag(RewardbaseFlag),
	altsrc.NewBoolFlag(CypressFlag),
	altsrc.NewBoolFlag(BaobabFlag),
	altsrc.NewInt64Flag(BlockGenerationIntervalFlag),
	altsrc.NewDurationFlag(BlockGenerationTimeLimitFlag),
}

var KPNFlags = []cli.Flag{
	altsrc.NewUint64Flag(TxResendIntervalFlag),
	altsrc.NewIntFlag(TxResendCountFlag),
	altsrc.NewBoolFlag(TxResendUseLegacyFlag),
	altsrc.NewBoolFlag(CypressFlag),
	altsrc.NewBoolFlag(BaobabFlag),
	altsrc.NewBoolFlag(TxPoolSpamThrottlerDisableFlag),
}

var KENFlags = []cli.Flag{
	altsrc.NewStringFlag(ServiceChainSignerFlag),
	altsrc.NewBoolFlag(CypressFlag),
	altsrc.NewBoolFlag(BaobabFlag),
	altsrc.NewBoolFlag(ChildChainIndexingFlag),
	altsrc.NewBoolFlag(MainBridgeFlag),
	altsrc.NewIntFlag(MainBridgeListenPortFlag),
	altsrc.NewBoolFlag(KESNodeTypeServiceFlag),
	// DBSyncer
	altsrc.NewBoolFlag(EnableDBSyncerFlag),
	altsrc.NewStringFlag(DBHostFlag),
	altsrc.NewStringFlag(DBPortFlag),
	altsrc.NewStringFlag(DBNameFlag),
	altsrc.NewStringFlag(DBUserFlag),
	altsrc.NewStringFlag(DBPasswordFlag),
	altsrc.NewBoolFlag(EnabledLogModeFlag),
	altsrc.NewIntFlag(MaxIdleConnsFlag),
	altsrc.NewIntFlag(MaxOpenConnsFlag),
	altsrc.NewDurationFlag(ConnMaxLifeTimeFlag),
	altsrc.NewIntFlag(BlockSyncChannelSizeFlag),
	altsrc.NewStringFlag(DBSyncerModeFlag),
	altsrc.NewIntFlag(GenQueryThreadFlag),
	altsrc.NewIntFlag(InsertThreadFlag),
	altsrc.NewIntFlag(BulkInsertSizeFlag),
	altsrc.NewStringFlag(EventModeFlag),
	altsrc.NewUint64Flag(MaxBlockDiffFlag),
	altsrc.NewUint64Flag(TxResendIntervalFlag),
	altsrc.NewIntFlag(TxResendCountFlag),
	altsrc.NewBoolFlag(TxResendUseLegacyFlag),
}

var KSCNFlags = []cli.Flag{
	altsrc.NewStringFlag(RewardbaseFlag),
	altsrc.NewInt64Flag(BlockGenerationIntervalFlag),
	altsrc.NewDurationFlag(BlockGenerationTimeLimitFlag),
	altsrc.NewStringFlag(ServiceChainSignerFlag),
	altsrc.NewUint64Flag(AnchoringPeriodFlag),
	altsrc.NewUint64Flag(SentChainTxsLimit),
	altsrc.NewBoolFlag(MainBridgeFlag),
	altsrc.NewIntFlag(MainBridgeListenPortFlag),
	altsrc.NewBoolFlag(ChildChainIndexingFlag),
	altsrc.NewBoolFlag(SubBridgeFlag),
	altsrc.NewIntFlag(SubBridgeListenPortFlag),
	altsrc.NewIntFlag(ParentChainIDFlag),
	altsrc.NewBoolFlag(VTRecoveryFlag),
	altsrc.NewUint64Flag(VTRecoveryIntervalFlag),
	altsrc.NewBoolFlag(ServiceChainNewAccountFlag),
	altsrc.NewBoolFlag(ServiceChainAnchoringFlag),
	altsrc.NewUint64Flag(ServiceChainParentOperatorTxGasLimitFlag),
	altsrc.NewUint64Flag(ServiceChainChildOperatorTxGasLimitFlag),
	// KAS
	altsrc.NewBoolFlag(KASServiceChainAnchorFlag),
	altsrc.NewUint64Flag(KASServiceChainAnchorPeriodFlag),
	altsrc.NewStringFlag(KASServiceChainAnchorUrlFlag),
	altsrc.NewStringFlag(KASServiceChainAnchorOperatorFlag),
	altsrc.NewStringFlag(KASServiceChainSecretKeyFlag),
	altsrc.NewStringFlag(KASServiceChainAccessKeyFlag),
	altsrc.NewStringFlag(KASServiceChainXChainIdFlag),
	altsrc.NewDurationFlag(KASServiceChainAnchorRequestTimeoutFlag),
}

var KSPNFlags = []cli.Flag{
	altsrc.NewUint64Flag(TxResendIntervalFlag),
	altsrc.NewIntFlag(TxResendCountFlag),
	altsrc.NewBoolFlag(TxResendUseLegacyFlag),
	altsrc.NewBoolFlag(TxPoolSpamThrottlerDisableFlag),
	altsrc.NewStringFlag(ServiceChainSignerFlag),
	altsrc.NewUint64Flag(AnchoringPeriodFlag),
	altsrc.NewUint64Flag(SentChainTxsLimit),
	altsrc.NewBoolFlag(MainBridgeFlag),
	altsrc.NewIntFlag(MainBridgeListenPortFlag),
	altsrc.NewBoolFlag(ChildChainIndexingFlag),
	altsrc.NewBoolFlag(SubBridgeFlag),
	altsrc.NewIntFlag(SubBridgeListenPortFlag),
	altsrc.NewIntFlag(ParentChainIDFlag),
	altsrc.NewBoolFlag(VTRecoveryFlag),
	altsrc.NewUint64Flag(VTRecoveryIntervalFlag),
	altsrc.NewBoolFlag(ServiceChainNewAccountFlag),
	altsrc.NewBoolFlag(ServiceChainAnchoringFlag),
	altsrc.NewUint64Flag(ServiceChainParentOperatorTxGasLimitFlag),
	altsrc.NewUint64Flag(ServiceChainChildOperatorTxGasLimitFlag),
	// KAS
	altsrc.NewBoolFlag(KASServiceChainAnchorFlag),
	altsrc.NewUint64Flag(KASServiceChainAnchorPeriodFlag),
	altsrc.NewStringFlag(KASServiceChainAnchorUrlFlag),
	altsrc.NewStringFlag(KASServiceChainAnchorOperatorFlag),
	altsrc.NewStringFlag(KASServiceChainSecretKeyFlag),
	altsrc.NewStringFlag(KASServiceChainAccessKeyFlag),
	altsrc.NewStringFlag(KASServiceChainXChainIdFlag),
	altsrc.NewDurationFlag(KASServiceChainAnchorRequestTimeoutFlag),
}

var KSENFlags = []cli.Flag{
	altsrc.NewStringFlag(ServiceChainSignerFlag),
	altsrc.NewBoolFlag(ChildChainIndexingFlag),
	altsrc.NewBoolFlag(MainBridgeFlag),
	altsrc.NewIntFlag(MainBridgeListenPortFlag),
	altsrc.NewBoolFlag(SubBridgeFlag),
	altsrc.NewIntFlag(SubBridgeListenPortFlag),
	altsrc.NewUint64Flag(AnchoringPeriodFlag),
	altsrc.NewUint64Flag(SentChainTxsLimit),
	altsrc.NewIntFlag(ParentChainIDFlag),
	altsrc.NewBoolFlag(VTRecoveryFlag),
	altsrc.NewUint64Flag(VTRecoveryIntervalFlag),
	altsrc.NewBoolFlag(ServiceChainAnchoringFlag),
	altsrc.NewBoolFlag(KESNodeTypeServiceFlag),
	altsrc.NewUint64Flag(ServiceChainParentOperatorTxGasLimitFlag),
	altsrc.NewUint64Flag(ServiceChainChildOperatorTxGasLimitFlag),
	// KAS
	altsrc.NewBoolFlag(KASServiceChainAnchorFlag),
	altsrc.NewUint64Flag(KASServiceChainAnchorPeriodFlag),
	altsrc.NewStringFlag(KASServiceChainAnchorUrlFlag),
	altsrc.NewStringFlag(KASServiceChainAnchorOperatorFlag),
	altsrc.NewStringFlag(KASServiceChainSecretKeyFlag),
	altsrc.NewStringFlag(KASServiceChainAccessKeyFlag),
	altsrc.NewStringFlag(KASServiceChainXChainIdFlag),
	altsrc.NewDurationFlag(KASServiceChainAnchorRequestTimeoutFlag),
	// DBSyncer
	altsrc.NewBoolFlag(EnableDBSyncerFlag),
	altsrc.NewStringFlag(DBHostFlag),
	altsrc.NewStringFlag(DBPortFlag),
	altsrc.NewStringFlag(DBNameFlag),
	altsrc.NewStringFlag(DBUserFlag),
	altsrc.NewStringFlag(DBPasswordFlag),
	altsrc.NewBoolFlag(EnabledLogModeFlag),
	altsrc.NewIntFlag(MaxIdleConnsFlag),
	altsrc.NewIntFlag(MaxOpenConnsFlag),
	altsrc.NewDurationFlag(ConnMaxLifeTimeFlag),
	altsrc.NewIntFlag(BlockSyncChannelSizeFlag),
	altsrc.NewStringFlag(DBSyncerModeFlag),
	altsrc.NewIntFlag(GenQueryThreadFlag),
	altsrc.NewIntFlag(InsertThreadFlag),
	altsrc.NewIntFlag(BulkInsertSizeFlag),
	altsrc.NewStringFlag(EventModeFlag),
	altsrc.NewUint64Flag(MaxBlockDiffFlag),
	altsrc.NewUint64Flag(TxResendIntervalFlag),
	altsrc.NewIntFlag(TxResendCountFlag),
	altsrc.NewBoolFlag(TxResendUseLegacyFlag),
}

var SnapshotFlags = []cli.Flag{
	altsrc.NewStringFlag(DbTypeFlag),
	altsrc.NewPathFlag(DataDirFlag),
	altsrc.NewPathFlag(ChainDataDirFlag),
	altsrc.NewBoolFlag(SingleDBFlag),
	altsrc.NewUintFlag(NumStateTrieShardsFlag),
	altsrc.NewStringFlag(DynamoDBTableNameFlag),
	altsrc.NewStringFlag(DynamoDBRegionFlag),
	altsrc.NewBoolFlag(DynamoDBIsProvisionedFlag),
	altsrc.NewInt64Flag(DynamoDBReadCapacityFlag),
	altsrc.NewInt64Flag(DynamoDBWriteCapacityFlag),
	altsrc.NewIntFlag(LevelDBCompressionTypeFlag),
	altsrc.NewBoolFlag(RocksDBSecondaryFlag),
	altsrc.NewUint64Flag(RocksDBCacheSizeFlag),
	altsrc.NewBoolFlag(RocksDBDumpMallocStatFlag),
	altsrc.NewStringFlag(RocksDBFilterPolicyFlag),
	altsrc.NewStringFlag(RocksDBCompressionTypeFlag),
	altsrc.NewStringFlag(RocksDBBottommostCompressionTypeFlag),
	altsrc.NewBoolFlag(RocksDBDisableMetricsFlag),
	altsrc.NewIntFlag(RocksDBMaxOpenFilesFlag),
	altsrc.NewBoolFlag(RocksDBCacheIndexAndFilterFlag),
}

var DBMigrationSrcFlags = []cli.Flag{
	altsrc.NewStringFlag(DbTypeFlag),
	altsrc.NewPathFlag(DataDirFlag),
	altsrc.NewBoolFlag(SingleDBFlag),
	altsrc.NewIntFlag(LevelDBCacheSizeFlag),
	altsrc.NewUintFlag(NumStateTrieShardsFlag),
	altsrc.NewStringFlag(DynamoDBTableNameFlag),
	altsrc.NewStringFlag(DynamoDBRegionFlag),
	altsrc.NewBoolFlag(DynamoDBIsProvisionedFlag),
	altsrc.NewInt64Flag(DynamoDBReadCapacityFlag),
	altsrc.NewInt64Flag(DynamoDBWriteCapacityFlag),
	altsrc.NewIntFlag(LevelDBCompressionTypeFlag),
	altsrc.NewBoolFlag(DBNoPerformanceMetricsFlag),
	altsrc.NewBoolFlag(RocksDBSecondaryFlag),
	altsrc.NewUint64Flag(RocksDBCacheSizeFlag),
	altsrc.NewBoolFlag(RocksDBDumpMallocStatFlag),
	altsrc.NewStringFlag(RocksDBFilterPolicyFlag),
	altsrc.NewStringFlag(RocksDBCompressionTypeFlag),
	altsrc.NewStringFlag(RocksDBBottommostCompressionTypeFlag),
	altsrc.NewBoolFlag(RocksDBDisableMetricsFlag),
	altsrc.NewIntFlag(RocksDBMaxOpenFilesFlag),
	altsrc.NewBoolFlag(RocksDBCacheIndexAndFilterFlag),
}

var DBMigrationDstFlags = []cli.Flag{
	altsrc.NewStringFlag(DstDbTypeFlag),
	altsrc.NewPathFlag(DstDataDirFlag),
	altsrc.NewBoolFlag(DstSingleDBFlag),
	altsrc.NewIntFlag(DstLevelDBCacheSizeFlag),
	altsrc.NewIntFlag(DstLevelDBCompressionTypeFlag),
	altsrc.NewUintFlag(DstNumStateTrieShardsFlag),
	altsrc.NewStringFlag(DstDynamoDBTableNameFlag),
	altsrc.NewStringFlag(DstDynamoDBRegionFlag),
	altsrc.NewBoolFlag(DstDynamoDBIsProvisionedFlag),
	altsrc.NewInt64Flag(DstDynamoDBReadCapacityFlag),
	altsrc.NewInt64Flag(DstDynamoDBWriteCapacityFlag),
	altsrc.NewUint64Flag(DstRocksDBCacheSizeFlag),
	altsrc.NewBoolFlag(DstRocksDBDumpMallocStatFlag),
	altsrc.NewBoolFlag(DstRocksDBDisableMetricsFlag),
	altsrc.NewBoolFlag(DstRocksDBSecondaryFlag),
	altsrc.NewStringFlag(DstRocksDBCompressionTypeFlag),
	altsrc.NewStringFlag(DstRocksDBBottommostCompressionTypeFlag),
	altsrc.NewStringFlag(DstRocksDBFilterPolicyFlag),
	altsrc.NewIntFlag(DstRocksDBMaxOpenFilesFlag),
	altsrc.NewBoolFlag(DstRocksDBCacheIndexAndFilterFlag),
}

var ChainDataFetcherFlags = []cli.Flag{
	altsrc.NewBoolFlag(EnableChainDataFetcherFlag),
	altsrc.NewStringFlag(ChainDataFetcherMode),
	altsrc.NewBoolFlag(ChainDataFetcherNoDefault),
	altsrc.NewIntFlag(ChainDataFetcherNumHandlers),
	altsrc.NewIntFlag(ChainDataFetcherJobChannelSize),
	altsrc.NewIntFlag(ChainDataFetcherChainEventSizeFlag),
	altsrc.NewIntFlag(ChainDataFetcherMaxProcessingDataSize),
	altsrc.NewStringFlag(ChainDataFetcherKASDBHostFlag),
	altsrc.NewStringFlag(ChainDataFetcherKASDBPortFlag),
	altsrc.NewStringFlag(ChainDataFetcherKASDBNameFlag),
	altsrc.NewStringFlag(ChainDataFetcherKASDBUserFlag),
	altsrc.NewStringFlag(ChainDataFetcherKASDBPasswordFlag),
	altsrc.NewBoolFlag(ChainDataFetcherKASCacheUse),
	altsrc.NewStringFlag(ChainDataFetcherKASCacheURLFlag),
	altsrc.NewStringFlag(ChainDataFetcherKASXChainIdFlag),
	altsrc.NewStringFlag(ChainDataFetcherKASBasicAuthParamFlag),
	altsrc.NewInt64Flag(ChainDataFetcherKafkaReplicasFlag),
	altsrc.NewStringSliceFlag(ChainDataFetcherKafkaBrokersFlag),
	altsrc.NewIntFlag(ChainDataFetcherKafkaPartitionsFlag),
	altsrc.NewStringFlag(ChainDataFetcherKafkaTopicResourceFlag),
	altsrc.NewStringFlag(ChainDataFetcherKafkaTopicEnvironmentFlag),
	altsrc.NewIntFlag(ChainDataFetcherKafkaMaxMessageBytesFlag),
	altsrc.NewIntFlag(ChainDataFetcherKafkaSegmentSizeBytesFlag),
	altsrc.NewIntFlag(ChainDataFetcherKafkaRequiredAcksFlag),
	altsrc.NewStringFlag(ChainDataFetcherKafkaMessageVersionFlag),
	altsrc.NewStringFlag(ChainDataFetcherKafkaProducerIdFlag),
}
