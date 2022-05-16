// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"sort"

	"github.com/klaytn/klaytn/api/debug"
	"gopkg.in/urfave/cli.v1"
)

const uncategorized = "MISC" // Uncategorized flags will belong to this group

// FlagGroup is a collection of flags belonging to a single topic.
type FlagGroup struct {
	Name  string
	Flags []cli.Flag
}

// TODO-Klaytn: consider changing the type of FlagGroups to map
// FlagGroups categorizes flags into groups to print structured help.
var FlagGroups = []FlagGroup{
	{
		Name: "KLAY",
		Flags: []cli.Flag{
			DbTypeFlag,
			DataDirFlag,
			KeyStoreDirFlag,
			IdentityFlag,
			SyncModeFlag,
			GCModeFlag,
			LightKDFFlag,
			SrvTypeFlag,
			ExtraDataFlag,
			ConfigFileFlag,
			OverwriteGenesisFlag,
			StartBlockNumberFlag,
			BlockGenerationIntervalFlag,
			BlockGenerationTimeLimitFlag,
			OpcodeComputationCostLimitFlag,
		},
	},
	{
		Name: "ACCOUNT",
		Flags: []cli.Flag{
			UnlockedAccountFlag,
			PasswordFileFlag,
		},
	},
	{
		Name: "TXPOOL",
		Flags: []cli.Flag{
			TxPoolNoLocalsFlag,
			TxPoolAllowLocalAnchorTxFlag,
			TxPoolDenyRemoteTxFlag,
			TxPoolJournalFlag,
			TxPoolJournalIntervalFlag,
			TxPoolPriceLimitFlag,
			TxPoolPriceBumpFlag,
			TxPoolExecSlotsAccountFlag,
			TxPoolExecSlotsAllFlag,
			TxPoolNonExecSlotsAccountFlag,
			TxPoolNonExecSlotsAllFlag,
			TxPoolLifetimeFlag,
			TxPoolKeepLocalsFlag,
			TxResendIntervalFlag,
			TxResendCountFlag,
			TxResendUseLegacyFlag,
		},
	},
	{
		Name: "DATABASE",
		Flags: []cli.Flag{
			LevelDBCacheSizeFlag,
			SingleDBFlag,
			NumStateTrieShardsFlag,
			LevelDBCompressionTypeFlag,
			LevelDBNoBufferPoolFlag,
			DynamoDBTableNameFlag,
			DynamoDBRegionFlag,
			DynamoDBIsProvisionedFlag,
			DynamoDBReadCapacityFlag,
			DynamoDBWriteCapacityFlag,
			NoParallelDBWriteFlag,
			SenderTxHashIndexingFlag,
			DBNoPerformanceMetricsFlag,
		},
	},
	{
		Name: "DATABASE SYNCER",
		Flags: []cli.Flag{
			EnableDBSyncerFlag,
			DBHostFlag,
			DBPortFlag,
			DBNameFlag,
			DBUserFlag,
			DBPasswordFlag,
			EnabledLogModeFlag,
			MaxIdleConnsFlag,
			MaxOpenConnsFlag,
			ConnMaxLifeTimeFlag,
			BlockSyncChannelSizeFlag,
			DBSyncerModeFlag,
			GenQueryThreadFlag,
			InsertThreadFlag,
			BulkInsertSizeFlag,
			EventModeFlag,
			MaxBlockDiffFlag,
		},
	},
	{
		Name: "CHAINDATAFETCHER",
		Flags: []cli.Flag{
			EnableChainDataFetcherFlag,
			ChainDataFetcherMode,
			ChainDataFetcherNoDefault,
			ChainDataFetcherNumHandlers,
			ChainDataFetcherJobChannelSize,
			ChainDataFetcherChainEventSizeFlag,
			ChainDataFetcherKASDBHostFlag,
			ChainDataFetcherKASDBPortFlag,
			ChainDataFetcherKASDBNameFlag,
			ChainDataFetcherKASDBUserFlag,
			ChainDataFetcherKASDBPasswordFlag,
			ChainDataFetcherKASCacheUse,
			ChainDataFetcherKASCacheURLFlag,
			ChainDataFetcherKASXChainIdFlag,
			ChainDataFetcherKASBasicAuthParamFlag,
			ChainDataFetcherKafkaBrokersFlag,
			ChainDataFetcherKafkaTopicEnvironmentFlag,
			ChainDataFetcherKafkaTopicResourceFlag,
			ChainDataFetcherKafkaReplicasFlag,
			ChainDataFetcherKafkaPartitionsFlag,
			ChainDataFetcherKafkaMaxMessageBytesFlag,
			ChainDataFetcherKafkaSegmentSizeBytesFlag,
			ChainDataFetcherKafkaRequiredAcksFlag,
			ChainDataFetcherKafkaMessageVersionFlag,
			ChainDataFetcherKafkaProducerIdFlag,
		},
	},
	{
		Name: "DATABASE MIGRATION",
		Flags: []cli.Flag{
			DstDbTypeFlag,
			DstDataDirFlag,
			DstSingleDBFlag,
			DstLevelDBCompressionTypeFlag,
			DstNumStateTrieShardsFlag,
			DstDynamoDBTableNameFlag,
			DstDynamoDBRegionFlag,
			DstDynamoDBIsProvisionedFlag,
			DstDynamoDBReadCapacityFlag,
			DstDynamoDBWriteCapacityFlag,
		},
	},
	{
		Name: "STATE",
		Flags: []cli.Flag{
			TrieMemoryCacheSizeFlag,
			TrieBlockIntervalFlag,
			TriesInMemoryFlag,
		},
	},
	{
		Name: "CACHE",
		Flags: []cli.Flag{
			CacheTypeFlag,
			CacheScaleFlag,
			CacheUsageLevelFlag,
			MemorySizeFlag,
			TrieNodeCacheTypeFlag,
			NumFetcherPrefetchWorkerFlag,
			UseSnapshotForPrefetchFlag,
			TrieNodeCacheLimitFlag,
			TrieNodeCacheSavePeriodFlag,
			TrieNodeCacheRedisEndpointsFlag,
			TrieNodeCacheRedisClusterFlag,
			TrieNodeCacheRedisPublishBlockFlag,
			TrieNodeCacheRedisSubscribeBlockFlag,
		},
	},
	{
		Name: "CONSENSUS",
		Flags: []cli.Flag{
			ServiceChainSignerFlag,
			RewardbaseFlag,
		},
	},
	{
		Name: "NETWORKING",
		Flags: []cli.Flag{
			BootnodesFlag,
			ListenPortFlag,
			SubListenPortFlag,
			MultiChannelUseFlag,
			MaxConnectionsFlag,
			MaxPendingPeersFlag,
			TargetGasLimitFlag,
			NATFlag,
			NoDiscoverFlag,
			RWTimerWaitTimeFlag,
			RWTimerIntervalFlag,
			NetrestrictFlag,
			NodeKeyFileFlag,
			NodeKeyHexFlag,
			NetworkIdFlag,
			BaobabFlag,
			CypressFlag,
		},
	},
	{
		Name: "METRICS",
		Flags: []cli.Flag{
			MetricsEnabledFlag,
			PrometheusExporterFlag,
			PrometheusExporterPortFlag,
		},
	},
	{
		Name: "VIRTUAL MACHINE",
		Flags: []cli.Flag{
			VMEnableDebugFlag,
			VMLogTargetFlag,
			VMTraceInternalTxFlag,
		},
	},
	{
		Name: "API AND CONSOLE",
		Flags: []cli.Flag{
			RPCEnabledFlag,
			RPCListenAddrFlag,
			RPCPortFlag,
			RPCCORSDomainFlag,
			RPCVirtualHostsFlag,
			RPCApiFlag,
			RPCGlobalGasCap,
			RPCGlobalEthTxFeeCapFlag,
			RPCConcurrencyLimit,
			RPCNonEthCompatibleFlag,
			IPCDisabledFlag,
			IPCPathFlag,
			WSEnabledFlag,
			WSListenAddrFlag,
			WSPortFlag,
			WSApiFlag,
			WSAllowedOriginsFlag,
			GRPCEnabledFlag,
			GRPCListenAddrFlag,
			GRPCPortFlag,
			JSpathFlag,
			ExecFlag,
			PreloadJSFlag,
			MaxRequestContentLengthFlag,
			APIFilterGetLogsDeadlineFlag,
			APIFilterGetLogsMaxItemsFlag,
		},
	},
	{
		Name:  "LOGGING AND DEBUGGING",
		Flags: debug.Flags,
	},
	{
		Name: "SERVICECHAIN",
		Flags: []cli.Flag{
			ChildChainIndexingFlag,
			MainBridgeFlag,
			MainBridgeListenPortFlag,
			SubBridgeFlag,
			SubBridgeListenPortFlag,
			AnchoringPeriodFlag,
			SentChainTxsLimit,
			ParentChainIDFlag,
			VTRecoveryFlag,
			VTRecoveryIntervalFlag,
			ServiceChainAnchoringFlag,
			ServiceChainNewAccountFlag,
			ServiceChainParentOperatorTxGasLimitFlag,
			ServiceChainChildOperatorTxGasLimitFlag,
			KASServiceChainAnchorFlag,
			KASServiceChainAnchorPeriodFlag,
			KASServiceChainAnchorUrlFlag,
			KASServiceChainAnchorOperatorFlag,
			KASServiceChainAccessKeyFlag,
			KASServiceChainSecretKeyFlag,
			KASServiceChainXChainIdFlag,
			KASServiceChainAnchorRequestTimeoutFlag,
		},
	},
	{
		Name: "MISC",
		Flags: []cli.Flag{
			GenKeyFlag,
			WriteAddressFlag,
			AutoRestartFlag,
			RestartTimeOutFlag,
			DaemonPathFlag,
			KESNodeTypeServiceFlag,
			SnapshotFlag,
			SnapshotCacheSizeFlag,
		},
	},
}

// CategorizeFlags classifies each flag into pre-defined flagGroups.
func CategorizeFlags(flags []cli.Flag) []FlagGroup {
	flagGroupsMap := make(map[string][]cli.Flag)
	isFlagAdded := make(map[string]bool) // Check duplicated flags

	// Find its group for each flag
	for _, flag := range flags {
		if isFlagAdded[flag.GetName()] {
			logger.Debug("a flag is added in the help description more than one time", "flag", flag.GetName())
			continue
		}

		// Find a group of the flag. If a flag doesn't belong to any groups, categorize it as a MISC flag
		group := flagCategory(flag, FlagGroups)
		flagGroupsMap[group] = append(flagGroupsMap[group], flag)
		isFlagAdded[flag.GetName()] = true
	}

	// Convert flagGroupsMap to a slice of FlagGroup
	flagGroups := []FlagGroup{}
	for group, flags := range flagGroupsMap {
		flagGroups = append(flagGroups, FlagGroup{Name: group, Flags: flags})
	}

	// Sort flagGroups in ascending order of name
	sortFlagGroup(flagGroups, uncategorized)

	return flagGroups
}

// sortFlagGroup sorts a slice of FlagGroup in ascending order of name,
// but an uncategorized group is exceptionally placed at the end.
func sortFlagGroup(flagGroups []FlagGroup, uncategorized string) []FlagGroup {
	sort.Slice(flagGroups, func(i, j int) bool {
		if flagGroups[i].Name == uncategorized {
			return false
		}
		if flagGroups[j].Name == uncategorized {
			return true
		}
		return flagGroups[i].Name < flagGroups[j].Name
	})

	// Sort flags in each group i ascending order of flag name.
	for _, group := range flagGroups {
		sort.Slice(group.Flags, func(i, j int) bool {
			return group.Flags[i].GetName() < group.Flags[j].GetName()
		})
	}

	return flagGroups
}

// flagCategory returns belonged group of the given flag.
func flagCategory(flag cli.Flag, fg []FlagGroup) string {
	for _, category := range fg {
		for _, flg := range category.Flags {
			if flg.GetName() == flag.GetName() {
				return category.Name
			}
		}
	}
	return uncategorized
}
