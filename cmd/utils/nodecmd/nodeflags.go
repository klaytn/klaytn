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
	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/cmd/utils"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

// TODO-Klaytn: Check whether all flags are registered in utils.FlagGroups

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

func KcnAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KcnNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, DBMigrationFlags...)
	return flags
}

func KpnAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KpnNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, DBMigrationFlags...)
	flags = append(flags, ChainDataFetcherFlags...)
	return flags
}

func KenAppFlags() []cli.Flag {
	flags := append([]cli.Flag{}, KenNodeFlags()...)
	flags = append(flags, CommonRPCFlags...)
	flags = append(flags, ConsoleFlags...)
	flags = append(flags, debug.Flags...)
	flags = append(flags, DBMigrationFlags...)
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

// Common flags that configure the node
var CommonNodeFlags = []cli.Flag{
	utils.ConfFlag,
	altsrc.NewBoolFlag(utils.NtpDisableFlag),
	altsrc.NewStringFlag(utils.NtpServerFlag),
	altsrc.NewStringFlag(utils.BootnodesFlag),
	altsrc.NewStringFlag(utils.IdentityFlag),
	altsrc.NewStringFlag(utils.UnlockedAccountFlag),
	altsrc.NewStringFlag(utils.PasswordFileFlag),
	altsrc.NewStringFlag(utils.DbTypeFlag),
	utils.NewWrappedDirectoryFlag(utils.DataDirFlag),
	utils.NewWrappedDirectoryFlag(utils.ChainDataDirFlag),
	altsrc.NewBoolFlag(utils.OverwriteGenesisFlag),
	altsrc.NewUint64Flag(utils.StartBlockNumberFlag),
	utils.NewWrappedDirectoryFlag(utils.KeyStoreDirFlag),
	altsrc.NewBoolFlag(utils.TxPoolNoLocalsFlag),
	altsrc.NewBoolFlag(utils.TxPoolAllowLocalAnchorTxFlag),
	altsrc.NewBoolFlag(utils.TxPoolDenyRemoteTxFlag),
	altsrc.NewStringFlag(utils.TxPoolJournalFlag),
	altsrc.NewDurationFlag(utils.TxPoolJournalIntervalFlag),
	altsrc.NewUint64Flag(utils.TxPoolPriceLimitFlag),
	altsrc.NewUint64Flag(utils.TxPoolPriceBumpFlag),
	altsrc.NewUint64Flag(utils.TxPoolExecSlotsAccountFlag),
	altsrc.NewUint64Flag(utils.TxPoolExecSlotsAllFlag),
	altsrc.NewUint64Flag(utils.TxPoolNonExecSlotsAccountFlag),
	altsrc.NewUint64Flag(utils.TxPoolNonExecSlotsAllFlag),
	altsrc.NewDurationFlag(utils.TxPoolLifetimeFlag),
	altsrc.NewBoolFlag(utils.TxPoolKeepLocalsFlag),
	utils.NewWrappedTextMarshalerFlag(utils.SyncModeFlag),
	altsrc.NewStringFlag(utils.GCModeFlag),
	altsrc.NewBoolFlag(utils.LightKDFFlag),
	altsrc.NewBoolFlag(utils.SingleDBFlag),
	altsrc.NewUintFlag(utils.NumStateTrieShardsFlag),
	altsrc.NewIntFlag(utils.LevelDBCompressionTypeFlag),
	altsrc.NewBoolFlag(utils.LevelDBNoBufferPoolFlag),
	altsrc.NewBoolFlag(utils.RocksDBSecondaryFlag),
	altsrc.NewUint64Flag(utils.RocksDBCacheSizeFlag),
	altsrc.NewBoolFlag(utils.RocksDBDumpMallocStatFlag),
	altsrc.NewStringFlag(utils.RocksDBCompressionTypeFlag),
	altsrc.NewStringFlag(utils.RocksDBBottommostCompressionTypeFlag),
	altsrc.NewStringFlag(utils.RocksDBFilterPolicyFlag),
	altsrc.NewBoolFlag(utils.RocksDBDisableMetricsFlag),
	altsrc.NewBoolFlag(utils.DBNoPerformanceMetricsFlag),
	altsrc.NewStringFlag(utils.DynamoDBTableNameFlag),
	altsrc.NewStringFlag(utils.DynamoDBRegionFlag),
	altsrc.NewBoolFlag(utils.DynamoDBIsProvisionedFlag),
	altsrc.NewInt64Flag(utils.DynamoDBReadCapacityFlag),
	altsrc.NewInt64Flag(utils.DynamoDBWriteCapacityFlag),
	altsrc.NewBoolFlag(utils.DynamoDBReadOnlyFlag),
	altsrc.NewIntFlag(utils.LevelDBCacheSizeFlag),
	altsrc.NewBoolFlag(utils.NoParallelDBWriteFlag),
	altsrc.NewBoolFlag(utils.SenderTxHashIndexingFlag),
	altsrc.NewIntFlag(utils.TrieMemoryCacheSizeFlag),
	altsrc.NewUintFlag(utils.TrieBlockIntervalFlag),
	altsrc.NewUint64Flag(utils.TriesInMemoryFlag),
	altsrc.NewIntFlag(utils.CacheTypeFlag),
	altsrc.NewIntFlag(utils.CacheScaleFlag),
	altsrc.NewStringFlag(utils.CacheUsageLevelFlag),
	altsrc.NewIntFlag(utils.MemorySizeFlag),
	altsrc.NewStringFlag(utils.TrieNodeCacheTypeFlag),
	altsrc.NewIntFlag(utils.NumFetcherPrefetchWorkerFlag),
	altsrc.NewBoolFlag(utils.UseSnapshotForPrefetchFlag),
	altsrc.NewIntFlag(utils.TrieNodeCacheLimitFlag),
	altsrc.NewDurationFlag(utils.TrieNodeCacheSavePeriodFlag),
	altsrc.NewStringSliceFlag(utils.TrieNodeCacheRedisEndpointsFlag),
	altsrc.NewBoolFlag(utils.TrieNodeCacheRedisClusterFlag),
	altsrc.NewBoolFlag(utils.TrieNodeCacheRedisPublishBlockFlag),
	altsrc.NewBoolFlag(utils.TrieNodeCacheRedisSubscribeBlockFlag),
	altsrc.NewIntFlag(utils.ListenPortFlag),
	altsrc.NewIntFlag(utils.SubListenPortFlag),
	altsrc.NewBoolFlag(utils.MultiChannelUseFlag),
	altsrc.NewIntFlag(utils.MaxConnectionsFlag),
	altsrc.NewIntFlag(utils.MaxRequestContentLengthFlag),
	altsrc.NewIntFlag(utils.MaxPendingPeersFlag),
	altsrc.NewUint64Flag(utils.TargetGasLimitFlag),
	altsrc.NewStringFlag(utils.NATFlag),
	altsrc.NewBoolFlag(utils.NoDiscoverFlag),
	altsrc.NewDurationFlag(utils.RWTimerWaitTimeFlag),
	altsrc.NewUint64Flag(utils.RWTimerIntervalFlag),
	altsrc.NewStringFlag(utils.NetrestrictFlag),
	altsrc.NewStringFlag(utils.NodeKeyFileFlag),
	altsrc.NewStringFlag(utils.NodeKeyHexFlag),
	altsrc.NewBoolFlag(utils.VMEnableDebugFlag),
	altsrc.NewIntFlag(utils.VMLogTargetFlag),
	altsrc.NewBoolFlag(utils.VMTraceInternalTxFlag),
	altsrc.NewUint64Flag(utils.NetworkIdFlag),
	altsrc.NewStringFlag(utils.RPCCORSDomainFlag),
	altsrc.NewStringFlag(utils.RPCVirtualHostsFlag),
	altsrc.NewBoolFlag(utils.RPCNonEthCompatibleFlag),
	altsrc.NewBoolFlag(utils.MetricsEnabledFlag),
	altsrc.NewBoolFlag(utils.PrometheusExporterFlag),
	altsrc.NewIntFlag(utils.PrometheusExporterPortFlag),
	altsrc.NewStringFlag(utils.ExtraDataFlag),
	altsrc.NewStringFlag(utils.SrvTypeFlag),
	altsrc.NewBoolFlag(utils.AutoRestartFlag),
	altsrc.NewDurationFlag(utils.RestartTimeOutFlag),
	altsrc.NewStringFlag(utils.DaemonPathFlag),
	altsrc.NewStringFlag(utils.ConfigFileFlag),
	altsrc.NewIntFlag(utils.APIFilterGetLogsMaxItemsFlag),
	altsrc.NewDurationFlag(utils.APIFilterGetLogsDeadlineFlag),
	altsrc.NewUint64Flag(utils.OpcodeComputationCostLimitFlag),
	altsrc.NewBoolFlag(utils.SnapshotFlag),
	altsrc.NewIntFlag(utils.SnapshotCacheSizeFlag),
	altsrc.NewBoolTFlag(utils.SnapshotAsyncGen),
}

// Common RPC flags
var CommonRPCFlags = []cli.Flag{
	altsrc.NewBoolFlag(utils.RPCEnabledFlag),
	altsrc.NewStringFlag(utils.RPCListenAddrFlag),
	altsrc.NewIntFlag(utils.RPCPortFlag),
	altsrc.NewStringFlag(utils.RPCApiFlag),
	altsrc.NewUint64Flag(utils.RPCGlobalGasCap),
	altsrc.NewDurationFlag(utils.RPCGlobalEVMTimeoutFlag),
	altsrc.NewFloat64Flag(utils.RPCGlobalEthTxFeeCapFlag),
	altsrc.NewBoolFlag(utils.WSEnabledFlag),
	altsrc.NewStringFlag(utils.WSListenAddrFlag),
	altsrc.NewIntFlag(utils.WSPortFlag),
	altsrc.NewBoolFlag(utils.GRPCEnabledFlag),
	altsrc.NewStringFlag(utils.GRPCListenAddrFlag),
	altsrc.NewIntFlag(utils.GRPCPortFlag),
	altsrc.NewIntFlag(utils.RPCConcurrencyLimit),
	altsrc.NewStringFlag(utils.WSApiFlag),
	altsrc.NewStringFlag(utils.WSAllowedOriginsFlag),
	altsrc.NewIntFlag(utils.WSMaxSubscriptionPerConn),
	altsrc.NewInt64Flag(utils.WSReadDeadLine),
	altsrc.NewInt64Flag(utils.WSWriteDeadLine),
	altsrc.NewIntFlag(utils.WSMaxConnections),
	altsrc.NewBoolFlag(utils.IPCDisabledFlag),
	utils.NewWrappedDirectoryFlag(utils.IPCPathFlag),
	altsrc.NewIntFlag(utils.RPCReadTimeout),
	altsrc.NewIntFlag(utils.RPCWriteTimeoutFlag),
	altsrc.NewIntFlag(utils.RPCIdleTimeoutFlag),
	altsrc.NewIntFlag(utils.RPCExecutionTimeoutFlag),
	altsrc.NewBoolFlag(utils.UnsafeDebugDisableFlag),
}

var KCNFlags = []cli.Flag{
	altsrc.NewStringFlag(utils.RewardbaseFlag),
	altsrc.NewBoolFlag(utils.CypressFlag),
	altsrc.NewBoolFlag(utils.BaobabFlag),
	altsrc.NewInt64Flag(utils.BlockGenerationIntervalFlag),
	altsrc.NewDurationFlag(utils.BlockGenerationTimeLimitFlag),
}

var KPNFlags = []cli.Flag{
	altsrc.NewUint64Flag(utils.TxResendIntervalFlag),
	altsrc.NewIntFlag(utils.TxResendCountFlag),
	altsrc.NewBoolFlag(utils.TxResendUseLegacyFlag),
	altsrc.NewBoolFlag(utils.CypressFlag),
	altsrc.NewBoolFlag(utils.BaobabFlag),
	altsrc.NewBoolFlag(utils.TxPoolSpamThrottlerDisableFlag),
}

var KENFlags = []cli.Flag{
	altsrc.NewStringFlag(utils.ServiceChainSignerFlag),
	altsrc.NewBoolFlag(utils.CypressFlag),
	altsrc.NewBoolFlag(utils.BaobabFlag),
	altsrc.NewBoolFlag(utils.ChildChainIndexingFlag),
	altsrc.NewBoolFlag(utils.MainBridgeFlag),
	altsrc.NewIntFlag(utils.MainBridgeListenPortFlag),
	altsrc.NewBoolFlag(utils.KESNodeTypeServiceFlag),
	// DBSyncer
	altsrc.NewBoolFlag(utils.EnableDBSyncerFlag),
	altsrc.NewStringFlag(utils.DBHostFlag),
	altsrc.NewStringFlag(utils.DBPortFlag),
	altsrc.NewStringFlag(utils.DBNameFlag),
	altsrc.NewStringFlag(utils.DBUserFlag),
	altsrc.NewStringFlag(utils.DBPasswordFlag),
	altsrc.NewBoolFlag(utils.EnabledLogModeFlag),
	altsrc.NewIntFlag(utils.MaxIdleConnsFlag),
	altsrc.NewIntFlag(utils.MaxOpenConnsFlag),
	altsrc.NewDurationFlag(utils.ConnMaxLifeTimeFlag),
	altsrc.NewIntFlag(utils.BlockSyncChannelSizeFlag),
	altsrc.NewStringFlag(utils.DBSyncerModeFlag),
	altsrc.NewIntFlag(utils.GenQueryThreadFlag),
	altsrc.NewIntFlag(utils.InsertThreadFlag),
	altsrc.NewIntFlag(utils.BulkInsertSizeFlag),
	altsrc.NewStringFlag(utils.EventModeFlag),
	altsrc.NewUint64Flag(utils.MaxBlockDiffFlag),
	altsrc.NewUint64Flag(utils.TxResendIntervalFlag),
	altsrc.NewIntFlag(utils.TxResendCountFlag),
	altsrc.NewBoolFlag(utils.TxResendUseLegacyFlag),
}

var KSCNFlags = []cli.Flag{
	altsrc.NewStringFlag(utils.RewardbaseFlag),
	altsrc.NewInt64Flag(utils.BlockGenerationIntervalFlag),
	altsrc.NewDurationFlag(utils.BlockGenerationTimeLimitFlag),
	altsrc.NewStringFlag(utils.ServiceChainSignerFlag),
	altsrc.NewUint64Flag(utils.AnchoringPeriodFlag),
	altsrc.NewUint64Flag(utils.SentChainTxsLimit),
	altsrc.NewBoolFlag(utils.MainBridgeFlag),
	altsrc.NewIntFlag(utils.MainBridgeListenPortFlag),
	altsrc.NewBoolFlag(utils.ChildChainIndexingFlag),
	altsrc.NewBoolFlag(utils.SubBridgeFlag),
	altsrc.NewIntFlag(utils.SubBridgeListenPortFlag),
	altsrc.NewIntFlag(utils.ParentChainIDFlag),
	altsrc.NewBoolFlag(utils.VTRecoveryFlag),
	altsrc.NewUint64Flag(utils.VTRecoveryIntervalFlag),
	altsrc.NewBoolFlag(utils.ServiceChainNewAccountFlag),
	altsrc.NewBoolFlag(utils.ServiceChainAnchoringFlag),
	altsrc.NewUint64Flag(utils.ServiceChainParentOperatorTxGasLimitFlag),
	altsrc.NewUint64Flag(utils.ServiceChainChildOperatorTxGasLimitFlag),
	// KAS
	altsrc.NewBoolFlag(utils.KASServiceChainAnchorFlag),
	altsrc.NewUint64Flag(utils.KASServiceChainAnchorPeriodFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAnchorUrlFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAnchorOperatorFlag),
	altsrc.NewStringFlag(utils.KASServiceChainSecretKeyFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAccessKeyFlag),
	altsrc.NewStringFlag(utils.KASServiceChainXChainIdFlag),
	altsrc.NewDurationFlag(utils.KASServiceChainAnchorRequestTimeoutFlag),
}

var KSPNFlags = []cli.Flag{
	altsrc.NewUint64Flag(utils.TxResendIntervalFlag),
	altsrc.NewIntFlag(utils.TxResendCountFlag),
	altsrc.NewBoolFlag(utils.TxResendUseLegacyFlag),
	altsrc.NewBoolFlag(utils.TxPoolSpamThrottlerDisableFlag),
	altsrc.NewStringFlag(utils.ServiceChainSignerFlag),
	altsrc.NewUint64Flag(utils.AnchoringPeriodFlag),
	altsrc.NewUint64Flag(utils.SentChainTxsLimit),
	altsrc.NewBoolFlag(utils.MainBridgeFlag),
	altsrc.NewIntFlag(utils.MainBridgeListenPortFlag),
	altsrc.NewBoolFlag(utils.ChildChainIndexingFlag),
	altsrc.NewBoolFlag(utils.SubBridgeFlag),
	altsrc.NewIntFlag(utils.SubBridgeListenPortFlag),
	altsrc.NewIntFlag(utils.ParentChainIDFlag),
	altsrc.NewBoolFlag(utils.VTRecoveryFlag),
	altsrc.NewUint64Flag(utils.VTRecoveryIntervalFlag),
	altsrc.NewBoolFlag(utils.ServiceChainNewAccountFlag),
	altsrc.NewBoolFlag(utils.ServiceChainAnchoringFlag),
	altsrc.NewUint64Flag(utils.ServiceChainParentOperatorTxGasLimitFlag),
	altsrc.NewUint64Flag(utils.ServiceChainChildOperatorTxGasLimitFlag),
	// KAS
	altsrc.NewBoolFlag(utils.KASServiceChainAnchorFlag),
	altsrc.NewUint64Flag(utils.KASServiceChainAnchorPeriodFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAnchorUrlFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAnchorOperatorFlag),
	altsrc.NewStringFlag(utils.KASServiceChainSecretKeyFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAccessKeyFlag),
	altsrc.NewStringFlag(utils.KASServiceChainXChainIdFlag),
	altsrc.NewDurationFlag(utils.KASServiceChainAnchorRequestTimeoutFlag),
}

var KSENFlags = []cli.Flag{
	altsrc.NewStringFlag(utils.ServiceChainSignerFlag),
	altsrc.NewBoolFlag(utils.ChildChainIndexingFlag),
	altsrc.NewBoolFlag(utils.MainBridgeFlag),
	altsrc.NewIntFlag(utils.MainBridgeListenPortFlag),
	altsrc.NewBoolFlag(utils.SubBridgeFlag),
	altsrc.NewIntFlag(utils.SubBridgeListenPortFlag),
	altsrc.NewUint64Flag(utils.AnchoringPeriodFlag),
	altsrc.NewUint64Flag(utils.SentChainTxsLimit),
	altsrc.NewIntFlag(utils.ParentChainIDFlag),
	altsrc.NewBoolFlag(utils.VTRecoveryFlag),
	altsrc.NewUint64Flag(utils.VTRecoveryIntervalFlag),
	altsrc.NewBoolFlag(utils.ServiceChainAnchoringFlag),
	altsrc.NewBoolFlag(utils.KESNodeTypeServiceFlag),
	altsrc.NewUint64Flag(utils.ServiceChainParentOperatorTxGasLimitFlag),
	altsrc.NewUint64Flag(utils.ServiceChainChildOperatorTxGasLimitFlag),
	// KAS
	altsrc.NewBoolFlag(utils.KASServiceChainAnchorFlag),
	altsrc.NewUint64Flag(utils.KASServiceChainAnchorPeriodFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAnchorUrlFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAnchorOperatorFlag),
	altsrc.NewStringFlag(utils.KASServiceChainSecretKeyFlag),
	altsrc.NewStringFlag(utils.KASServiceChainAccessKeyFlag),
	altsrc.NewStringFlag(utils.KASServiceChainXChainIdFlag),
	altsrc.NewDurationFlag(utils.KASServiceChainAnchorRequestTimeoutFlag),
	// DBSyncer
	altsrc.NewBoolFlag(utils.EnableDBSyncerFlag),
	altsrc.NewStringFlag(utils.DBHostFlag),
	altsrc.NewStringFlag(utils.DBPortFlag),
	altsrc.NewStringFlag(utils.DBNameFlag),
	altsrc.NewStringFlag(utils.DBUserFlag),
	altsrc.NewStringFlag(utils.DBPasswordFlag),
	altsrc.NewBoolFlag(utils.EnabledLogModeFlag),
	altsrc.NewIntFlag(utils.MaxIdleConnsFlag),
	altsrc.NewIntFlag(utils.MaxOpenConnsFlag),
	altsrc.NewDurationFlag(utils.ConnMaxLifeTimeFlag),
	altsrc.NewIntFlag(utils.BlockSyncChannelSizeFlag),
	altsrc.NewStringFlag(utils.DBSyncerModeFlag),
	altsrc.NewIntFlag(utils.GenQueryThreadFlag),
	altsrc.NewIntFlag(utils.InsertThreadFlag),
	altsrc.NewIntFlag(utils.BulkInsertSizeFlag),
	altsrc.NewStringFlag(utils.EventModeFlag),
	altsrc.NewUint64Flag(utils.MaxBlockDiffFlag),
	altsrc.NewUint64Flag(utils.TxResendIntervalFlag),
	altsrc.NewIntFlag(utils.TxResendCountFlag),
	altsrc.NewBoolFlag(utils.TxResendUseLegacyFlag),
}

var DBMigrationFlags = []cli.Flag{
	altsrc.NewStringFlag(utils.DstDbTypeFlag),
	utils.NewWrappedDirectoryFlag(utils.DstDataDirFlag),
	altsrc.NewBoolFlag(utils.DstSingleDBFlag),
	altsrc.NewIntFlag(utils.DstLevelDBCompressionTypeFlag),
	altsrc.NewUintFlag(utils.DstNumStateTrieShardsFlag),
	altsrc.NewStringFlag(utils.DstDynamoDBTableNameFlag),
	altsrc.NewStringFlag(utils.DstDynamoDBRegionFlag),
	altsrc.NewBoolFlag(utils.DstDynamoDBIsProvisionedFlag),
	altsrc.NewInt64Flag(utils.DstDynamoDBReadCapacityFlag),
	altsrc.NewInt64Flag(utils.DstDynamoDBWriteCapacityFlag),
}

var ChainDataFetcherFlags = []cli.Flag{
	altsrc.NewBoolFlag(utils.EnableChainDataFetcherFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherMode),
	altsrc.NewBoolFlag(utils.ChainDataFetcherNoDefault),
	altsrc.NewIntFlag(utils.ChainDataFetcherNumHandlers),
	altsrc.NewIntFlag(utils.ChainDataFetcherJobChannelSize),
	altsrc.NewIntFlag(utils.ChainDataFetcherChainEventSizeFlag),
	altsrc.NewIntFlag(utils.ChainDataFetcherMaxProcessingDataSize),
	altsrc.NewStringFlag(utils.ChainDataFetcherKASDBHostFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKASDBPortFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKASDBNameFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKASDBUserFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKASDBPasswordFlag),
	altsrc.NewBoolFlag(utils.ChainDataFetcherKASCacheUse),
	altsrc.NewStringFlag(utils.ChainDataFetcherKASCacheURLFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKASXChainIdFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKASBasicAuthParamFlag),
	altsrc.NewInt64Flag(utils.ChainDataFetcherKafkaReplicasFlag),
	altsrc.NewStringSliceFlag(utils.ChainDataFetcherKafkaBrokersFlag),
	altsrc.NewIntFlag(utils.ChainDataFetcherKafkaPartitionsFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKafkaTopicResourceFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKafkaTopicEnvironmentFlag),
	altsrc.NewInt64Flag(utils.ChainDataFetcherKafkaMaxMessageBytesFlag),
	altsrc.NewIntFlag(utils.ChainDataFetcherKafkaSegmentSizeBytesFlag),
	altsrc.NewIntFlag(utils.ChainDataFetcherKafkaRequiredAcksFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKafkaMessageVersionFlag),
	altsrc.NewStringFlag(utils.ChainDataFetcherKafkaProducerIdFlag),
}
