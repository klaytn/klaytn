// Modifications Copyright 2019 The klaytn Authors
// Copyright 2018 The go-ethereum Authors
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
// This file is derived from build/update-license.go (2018/06/04).
// Modified and improved for the klaytn development.

// +build none

/*
This command generates GPL license headers on top of all source files.
You can run it once per month, before cutting a release or just
whenever you feel like it.

	go run update-license.go

All authors (people who have contributed code) are listed in the
AUTHORS file. The author names are mapped and deduplicated using the
.mailmap file. You can use .mailmap to set the canonical name and
address for each author. See git-shortlog(1) for an explanation of the
.mailmap format.

Please review the resulting diff to check whether the correct
copyright assignments are performed.
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	// only files with these extensions will be considered
	extensions = []string{".go", ".js", ".qml"}

	// paths with any of these prefixes will be skipped
	skipPrefixes = []string{
		// boring stuff
		"vendor/", "tests/testdata/", "build/",
		// don't relicense vendored sources
		"cmd/internal/browser",
		"cmd/homi",
		"consensus/ethash/xor.go",
		"crypto/bn256/",
		"crypto/ecies/",
		"crypto/secp256k1",
		//"crypto/secp256k1/curve.go",
		"crypto/sha3/",
		"console/jsre/deps",
		//"log/",
		"common/bitutil/bitutil",
		// don't license generated files
		"contracts/chequebook/contract/code.go",
		"log/term",
		"node/cn/tracers",
		"metrics/exp",
		"metrics/influxdb",
		"metrics/librato",
	}

	skipSuffixes = []string{
		"doc.go",
	}

	externalLicencePrefixes = []string{
		"log", "metrics",
	}

	// paths with this prefix are licensed as GPL. all other files are LGPL.
	gplPrefixes = []string{"cmd/"}

	// this regexp must match the entire license comment at the
	// beginning of each file.
	ethereumLicenseCommentRE   = regexp.MustCompile(`//\s*(Copyright .* (The go-ethereum|AMIS Technologies)).*?\n(?://.*?\n)*//.*>\.`)
	klaytnLicenseCommentRE     = regexp.MustCompile(`//\s*(Copyright .* The klaytn).*?\n(?://.*?\n)?// This file is part of the klaytn library.*\n(?://.*?\n)*//.*>\.`)
	mixedLicenseCommentRE      = regexp.MustCompile(`//\s*(.*Copyright .* The klaytn).*?\n//\s*(Copyright .* (The go-ethereum|AMIS Technologies)).*?\n(?://.*?\n)*//.*klaytn development\.`)
	externalLicenceCommentRE   = regexp.MustCompile(`//\s*(Copyright .* The klaytn).*?\n(?://.*?\n)*//.*See LICENSE in the top directory for the original copyright and license\.`)
	externalGoLicenceCommentRE = regexp.MustCompile(`//\s*(Copyright .* The Go Authors).*?\n(?://.*?\n)*//.*license that can be found in the LICENSE file\.`)

	// this text appears at the start of AUTHORS
	authorsFileHeader = "# This is the official list of go-ethereum authors for copyright purposes.\n\n"

	// ethereumDir is original path referenced by klaytn
	ethereumDir = map[string]string{
		"metrics/meter.go":                                   "metrics/meter.go",
		"metrics/log.go":                                     "metrics/log.go",
		"metrics/runtime_cgo.go":                             "metrics/runtime_cgo.go",
		"metrics/counter_test.go":                            "metrics/counter_test.go",
		"metrics/histogram_test.go":                          "metrics/histogram_test.go",
		"metrics/sample_test.go":                             "metrics/sample_test.go",
		"metrics/metrics.go":                                 "metrics/metrics.go",
		"metrics/gauge.go":                                   "metrics/gauge.go",
		"metrics/disk.go":                                    "metrics/disk.go",
		"metrics/writer_test.go":                             "metrics/writer_test.go",
		"metrics/gauge_float64.go":                           "metrics/gauge_float64.go",
		"metrics/healthcheck.go":                             "metrics/healthcheck.go",
		"metrics/runtime.go":                                 "metrics/runtime.go",
		"metrics/init_test.go":                               "metrics/init_test.go",
		"metrics/metrics_test.go":                            "metrics/metrics_test.go",
		"metrics/resetting_timer.go":                         "metrics/resetting_timer.go",
		"metrics/opentsdb_test.go":                           "metrics/opentsdb_test.go",
		"metrics/graphite.go":                                "metrics/graphite.go",
		"metrics/graphite_test.go":                           "metrics/graphite_test.go",
		"metrics/ewma_test.go":                               "metrics/ewma_test.go",
		"metrics/timer.go":                                   "metrics/timer.go",
		"metrics/registry.go":                                "metrics/registry.go",
		"metrics/histogram.go":                               "metrics/histogram.go",
		"metrics/timer_test.go":                              "metrics/timer_test.go",
		"metrics/runtime_no_gccpufraction.go":                "metrics/runtime_no_gccpufraction.go",
		"metrics/librato/client.go":                          "metrics/librato/client.go",
		"metrics/librato/librato.go":                         "metrics/librato/librato.go",
		"metrics/influxdb/influxdb.go":                       "metrics/influxdb/influxdb.go",
		"metrics/json_test.go":                               "metrics/json_test.go",
		"metrics/writer.go":                                  "metrics/writer.go",
		"metrics/sample.go":                                  "metrics/sample.go",
		"metrics/ewma.go":                                    "metrics/ewma.go",
		"metrics/syslog.go":                                  "metrics/syslog.go",
		"metrics/debug.go":                                   "metrics/debug.go",
		"metrics/runtime_gccpufraction.go":                   "metrics/runtime_gccpufraction.go",
		"metrics/gauge_float64_test.go":                      "metrics/gauge_float64_test.go",
		"metrics/registry_test.go":                           "metrics/registry_test.go",
		"metrics/disk_nop.go":                                "metrics/disk_nop.go",
		"metrics/disk_linux.go":                              "metrics/disk_linux.go",
		"metrics/gauge_test.go":                              "metrics/gauge_test.go",
		"metrics/runtime_test.go":                            "metrics/runtime_test.go",
		"metrics/meter_test.go":                              "metrics/meter_test.go",
		"metrics/counter.go":                                 "metrics/counter.go",
		"metrics/exp/exp.go":                                 "metrics/exp/exp.go",
		"metrics/json.go":                                    "metrics/json.go",
		"metrics/resetting_timer_test.go":                    "metrics/resetting_timer_test.go",
		"metrics/runtime_no_cgo.go":                          "metrics/runtime_no_cgo.go",
		"metrics/opentsdb.go":                                "metrics/opentsdb.go",
		"metrics/debug_test.go":                              "metrics/debug_test.go",
		"metrics/prometheus/prometheusmetrics.go":            "metrics/prometheus/prometheusmetrics.go",
		"cmd/kbn/main.go":                                    "cmd/bootnode/main.go",
		"cmd/kbn/node.go":                                    "node/node.go",
		"cmd/kcn/main.go":                                    "cmd/geth/main.go",
		"cmd/ken/main.go":                                    "cmd/geth/main.go",
		"cmd/kpn/main.go":                                    "cmd/geth/main.go",
		"cmd/kscn/main.go":                                   "cmd/geth/main.go",
		"cmd/kspn/main.go":                                   "cmd/geth/main.go",
		"cmd/ksen/main.go":                                   "cmd/geth/main.go",
		"cmd/klay/main.go":                                   "cmd/geth/main.go",
		"cmd/utils/cmd.go":                                   "cmd/utils/cmd.go",
		"cmd/utils/flags.go":                                 "cmd/utils/flags.go",
		"cmd/utils/app.go":                                   "cmd/utils/flags.go",
		"cmd/utils/customflags.go":                           "cmd/utils/customflags.go",
		"cmd/utils/nodecmd/accountcmd.go":                    "cmd/geth/accountcmd.go",
		"cmd/utils/nodecmd/accountcmd_test.go":               "cmd/geth/accountcmd_test.go",
		"cmd/utils/nodecmd/chaincmd.go":                      "cmd/geth/chaincmd.go",
		"cmd/utils/nodecmd/consolecmd.go":                    "cmd/geth/consolecmd.go",
		"cmd/utils/nodecmd/consolecmd_test.go":               "cmd/geth/consolecmd_test.go",
		"cmd/utils/nodecmd/defaultcmd.go":                    "cmd/geth/main.go",
		"cmd/utils/nodecmd/dumpconfigcmd.go":                 "cmd/geth/config.go",
		"cmd/utils/nodecmd/genesis_test.go":                  "cmd/geth/genesis_test.go",
		"cmd/utils/nodecmd/flags_test.go":                    "cmd/geth/genesis_test.go",
		"cmd/utils/nodecmd/nodeflags.go":                     "cmd/geth/main.go",
		"cmd/utils/nodecmd/run_test.go":                      "cmd/geth/run_test.go",
		"cmd/utils/nodecmd/versioncmd.go":                    "cmd/geth/misccmd.go",
		"cmd/utils/testcmd.go":                               "cmd/cmdtest/test_cmd.go",
		"cmd/utils/usage.go":                                 "cmd/geth/usage.go",
		"cmd/p2psim/main.go":                                 "cmd/p2psim/main.go",
		"cmd/abigen/main.go":                                 "cmd/abigen/main.go",
		"cmd/evm/compiler.go":                                "cmd/evm/compiler.go",
		"cmd/evm/runner.go":                                  "cmd/evm/runner.go",
		"cmd/evm/disasm.go":                                  "cmd/evm/disasm.go",
		"cmd/evm/staterunner.go":                             "cmd/evm/staterunner.go",
		"cmd/evm/internal/compiler/compiler.go":              "cmd/evm/internal/compiler/compiler.go",
		"cmd/evm/main.go":                                    "cmd/evm/main.go",
		"cmd/evm/json_logger.go":                             "cmd/evm/json_logger.go",
		"cmd/grpc-contract/internal/impl/contract.go":        "sol2proto/types/grpc/contract.go",
		"cmd/grpc-contract/internal/impl/mapping.go":         "sol2proto/types/grpc/mapping.go",
		"cmd/grpc-contract/internal/impl/method.go":          "sol2proto/types/grpc/types.go",
		"cmd/grpc-contract/main.go":                          "sol2proto/types/main.go",
		"consensus/istanbul/config.go":                       "quorum/consensus/istanbul/config.go",
		"consensus/istanbul/backend.go":                      "quorum/consensus/istanbul/backend.go",
		"consensus/istanbul/core/handler.go":                 "quorum/consensus/istanbul/core/handler.go",
		"consensus/istanbul/core/backlog.go":                 "quorum/consensus/istanbul/core/backlog.go",
		"consensus/istanbul/core/types.go":                   "quorum/consensus/istanbul/core/types.go",
		"consensus/istanbul/core/preprepare.go":              "quorum/consensus/istanbul/core/preprepare.go",
		"consensus/istanbul/core/core.go":                    "quorum/consensus/istanbul/core/core.go",
		"consensus/istanbul/core/request.go":                 "quorum/consensus/istanbul/core/request.go",
		"consensus/istanbul/core/commit.go":                  "quorum/consensus/istanbul/core/commit.go",
		"consensus/istanbul/core/message_set.go":             "quorum/consensus/istanbul/core/message_set.go",
		"consensus/istanbul/core/events.go":                  "quorum/consensus/istanbul/core/events.go",
		"consensus/istanbul/core/final_committed.go":         "quorum/consensus/istanbul/core/final_committed.go",
		"consensus/istanbul/core/roundstate.go":              "quorum/consensus/istanbul/core/roundstate.go",
		"consensus/istanbul/core/prepare.go":                 "quorum/consensus/istanbul/core/prepare.go",
		"consensus/istanbul/core/errors.go":                  "quorum/consensus/istanbul/core/errors.go",
		"consensus/istanbul/core/roundchange.go":             "quorum/consensus/istanbul/core/roundchange.go",
		"consensus/istanbul/validator/validator.go":          "quorum/consensus/istanbul/validator/validator.go",
		"consensus/istanbul/validator/default_test.go":       "quorum/consensus/istanbul/validator/default_test.go",
		"consensus/istanbul/validator/default.go":            "quorum/consensus/istanbul/validator/default.go",
		"consensus/istanbul/validator/weighted.go":           "quorum/consensus/istanbul/validator/default.go",
		"consensus/istanbul/validator.go":                    "quorum/consensus/istanbul/validator.go",
		"consensus/istanbul/types.go":                        "quorum/consensus/istanbul/types.go",
		"consensus/istanbul/backend/handler.go":              "quorum/consensus/istanbul/backend/handler.go",
		"consensus/istanbul/backend/backend.go":              "quorum/consensus/istanbul/backend/backend.go",
		"consensus/istanbul/backend/engine.go":               "quorum/consensus/istanbul/backend/engine.go",
		"consensus/istanbul/backend/api.go":                  "quorum/consensus/istanbul/backend/api.go",
		"consensus/istanbul/backend/snapshot.go":             "quorum/consensus/istanbul/backend/snapshot.go",
		"consensus/istanbul/events.go":                       "quorum/consensus/istanbul/events.go",
		"consensus/istanbul/utils.go":                        "quorum/consensus/istanbul/utils.go",
		"consensus/istanbul/errors.go":                       "quorum/consensus/istanbul/errors.go",
		"consensus/protocol.go":                              "quorum/consensus/protocol.go",
		"consensus/consensus.go":                             "consensus/consensus.go",
		"consensus/gxhash/sealer.go":                         "consensus/ethash/sealer.go",
		"consensus/gxhash/algorithm.go":                      "consensus/ethash/algorithm.go",
		"consensus/gxhash/gxhash_test.go":                    "consensus/ethash/ethash_test.go",
		"consensus/gxhash/consensus.go":                      "consensus/ethash/consensus.go",
		"consensus/gxhash/algorithm_test.go":                 "consensus/ethash/algorithm_test.go",
		"consensus/gxhash/gxhash.go":                         "consensus/ethash/ethash.go",
		"consensus/gxhash/consensus_test.go":                 "consensus/ethash/consensus_test.go",
		"consensus/clique/api.go":                            "go-ethereum/consensus/clique/api.go",
		"consensus/clique/clique.go":                         "go-ethereum/consensus/clique/clique.go",
		"consensus/clique/snapshot.go":                       "go-ethereum/consensus/clique/snapshot.go",
		"crypto/signature_cgo.go":                            "crypto/signature_cgo.go",
		"crypto/sha3/register.go":                            "crypto/sha3/register.go",
		"crypto/sha3/xor_unaligned.go":                       "crypto/sha3/xor_unaligned.go",
		"crypto/sha3/sha3_test.go":                           "crypto/sha3/sha3_test.go",
		"crypto/sha3/xor_generic.go":                         "crypto/sha3/xor_generic.go",
		"crypto/sha3/shake.go":                               "crypto/sha3/shake.go",
		"crypto/sha3/keccakf_amd64.go":                       "crypto/sha3/keccakf_amd64.go",
		"crypto/sha3/xor.go":                                 "crypto/sha3/xor.go",
		"crypto/sha3/doc.go":                                 "crypto/sha3/doc.go",
		"crypto/sha3/sha3.go":                                "crypto/sha3/sha3.go",
		"crypto/sha3/keccakf.go":                             "crypto/sha3/keccakf.go",
		"crypto/sha3/hashes.go":                              "crypto/sha3/hashes.go",
		"crypto/ecies/params.go":                             "crypto/ecies/params.go",
		"crypto/ecies/ecies_test.go":                         "crypto/ecies/ecies_test.go",
		"crypto/ecies/ecies.go":                              "crypto/ecies/ecies.go",
		"crypto/secp256k1/curve.go":                          "crypto/secp256k1/curve.go",
		"crypto/secp256k1/panic_cb.go":                       "crypto/secp256k1/panic_cb.go",
		"crypto/secp256k1/secp256.go":                        "crypto/secp256k1/secp256.go",
		"crypto/secp256k1/secp256_test.go":                   "crypto/secp256k1/secp256_test.go",
		"crypto/bn256/bn256_fast.go":                         "crypto/bn256/bn256_fast.go",
		"crypto/bn256/google/optate.go":                      "crypto/bn256/google/optate.go",
		"crypto/bn256/google/example_test.go":                "crypto/bn256/google/example_test.go",
		"crypto/bn256/google/bn256_test.go":                  "crypto/bn256/google/bn256_test.go",
		"crypto/bn256/google/twist.go":                       "crypto/bn256/google/twist.go",
		"crypto/bn256/google/constants.go":                   "crypto/bn256/google/constants.go",
		"crypto/bn256/google/curve.go":                       "crypto/bn256/google/curve.go",
		"crypto/bn256/google/gfp2.go":                        "crypto/bn256/google/gfp2.go",
		"crypto/bn256/google/gfp6.go":                        "crypto/bn256/google/gfp6.go",
		"crypto/bn256/google/bn256.go":                       "crypto/bn256/google/bn256.go",
		"crypto/bn256/google/gfp12.go":                       "crypto/bn256/google/gfp12.go",
		"crypto/bn256/google/main_test.go":                   "crypto/bn256/google/main_test.go",
		"crypto/bn256/cloudflare/optate.go":                  "crypto/bn256/cloudflare/optate.go",
		"crypto/bn256/cloudflare/gfp_amd64.s":                "crypto/bn256/cloudflare/gfp_amd64.s",
		"crypto/bn256/cloudflare/gfp_decl.go":                "crypto/bn256/cloudflare/gfp_decl.go",
		"crypto/bn256/cloudflare/example_test.go":            "crypto/bn256/cloudflare/example_test.go",
		"crypto/bn256/cloudflare/gfp_generic.go":             "crypto/bn256/cloudflare/gfp_generic.go",
		"crypto/bn256/cloudflare/bn256_test.go":              "crypto/bn256/cloudflare/bn256_test.go",
		"crypto/bn256/cloudflare/twist.go":                   "crypto/bn256/cloudflare/twist.go",
		"crypto/bn256/cloudflare/constants.go":               "crypto/bn256/cloudflare/constants.go",
		"crypto/bn256/cloudflare/lattice.go":                 "crypto/bn256/cloudflare/lattice.go",
		"crypto/bn256/cloudflare/gfp_test.go":                "crypto/bn256/cloudflare/gfp_test.go",
		"crypto/bn256/cloudflare/curve.go":                   "crypto/bn256/cloudflare/curve.go",
		"crypto/bn256/cloudflare/gfp2.go":                    "crypto/bn256/cloudflare/gfp2.go",
		"crypto/bn256/cloudflare/gfp.go":                     "crypto/bn256/cloudflare/gfp.go",
		"crypto/bn256/cloudflare/gfp6.go":                    "crypto/bn256/cloudflare/gfp6.go",
		"crypto/bn256/cloudflare/lattice_test.go":            "crypto/bn256/cloudflare/lattice_test.go",
		"crypto/bn256/cloudflare/bn256.go":                   "crypto/bn256/cloudflare/bn256.go",
		"crypto/bn256/cloudflare/gfp12.go":                   "crypto/bn256/cloudflare/gfp12.go",
		"crypto/bn256/cloudflare/main_test.go":               "crypto/bn256/cloudflare/main_test.go",
		"crypto/bn256/bn256_fuzz.go":                         "crypto/bn256/bn256_fuzz.go",
		"crypto/bn256/bn256_slow.go":                         "crypto/bn256/bn256_slow.go",
		"crypto/signature_test.go":                           "crypto/signature_test.go",
		"crypto/crypto.go":                                   "crypto/crypto.go",
		"crypto/crypto_test.go":                              "crypto/crypto_test.go",
		"datasync/fetcher/metrics.go":                        "eth/fetcher/metrics.go",
		"datasync/fetcher/fetcher_test.go":                   "eth/fetcher/fetcher_test.go",
		"datasync/fetcher/fetcher.go":                        "eth/fetcher/fetcher.go",
		"datasync/downloader/downloader.go":                  "eth/downloader/downloader.go",
		"datasync/downloader/metrics.go":                     "eth/downloader/metrics.go",
		"datasync/downloader/queue.go":                       "eth/downloader/queue.go",
		"datasync/downloader/modes.go":                       "eth/downloader/modes.go",
		"datasync/downloader/types.go":                       "eth/downloader/types.go",
		"datasync/downloader/downloader_test.go":             "eth/downloader/downloader_test.go",
		"datasync/downloader/events.go":                      "eth/downloader/events.go",
		"datasync/downloader/testchain_test.go":              "eth/downloader/testchain_test.go",
		"datasync/downloader/api.go":                         "eth/downloader/api.go",
		"datasync/downloader/peer.go":                        "eth/downloader/peer.go",
		"datasync/downloader/statesync.go":                   "eth/downloader/statesync.go",
		"interfaces.go":                                      "interfaces.go",
		"tests/vm_test_util.go":                              "tests/vm_test_util.go",
		"tests/gen_stenv.go":                                 "tests/gen_stenv.go",
		"tests/gen_tttransaction.go":                         "tests/gen_tttransaction.go",
		"tests/gen_btheader.go":                              "tests/gen_btheader.go",
		"tests/state_test_util.go":                           "tests/state_test_util.go",
		"tests/gen_sttransaction.go":                         "tests/gen_sttransaction.go",
		"tests/transaction_test.go":                          "tests/transaction_test.go",
		"tests/init_test.go":                                 "tests/init_test.go",
		"tests/gen_vmexec.go":                                "tests/gen_vmexec.go",
		"tests/block_test.go":                                "tests/block_test.go",
		"tests/transaction_test_util.go":                     "tests/transaction_test_util.go",
		"tests/state_test.go":                                "tests/state_test.go",
		"tests/rlp_test.go":                                  "tests/rlp_test.go",
		"tests/vm_test.go":                                   "tests/vm_test.go",
		"tests/block_test_util.go":                           "tests/block_test_util.go",
		"tests/rlp_test_util.go":                             "tests/rlp_test_util.go",
		"tests/init.go":                                      "tests/init.go",
		"utils/build/env.go":                                 "internal/build/env.go",
		"utils/build/pgp.go":                                 "internal/build/pgp.go",
		"utils/build/util.go":                                "internal/build/util.go",
		"utils/build/archive.go":                             "internal/build/archive.go",
		"storage/database/leveldb_database.go":               "ethdb/database.go",
		"storage/database/interface.go":                      "ethdb/interface.go",
		"storage/database/database_test.go":                  "ethdb/database_test.go",
		"storage/database/memory_database.go":                "ethdb/memory_database.go",
		"storage/database/schema.go":                         "core/rawdb/schema.go",
		"storage/statedb/encoding.go":                        "trie/encoding.go",
		"storage/statedb/secure_trie.go":                     "trie/secure_trie.go",
		"storage/statedb/derive_sha.go":                      "core/types/derive_sha.go",
		"storage/statedb/sync.go":                            "trie/sync.go",
		"storage/statedb/sync_test.go":                       "trie/sync_test.go",
		"storage/statedb/proof_test.go":                      "trie/proof_test.go",
		"storage/statedb/database.go":                        "trie/database.go",
		"storage/statedb/iterator_test.go":                   "trie/iterator_test.go",
		"storage/statedb/encoding_test.go":                   "trie/encoding_test.go",
		"storage/statedb/proof.go":                           "trie/proof.go",
		"storage/statedb/iterator.go":                        "trie/iterator.go",
		"storage/statedb/node_test.go":                       "trie/node_test.go",
		"storage/statedb/secure_trie_test.go":                "trie/secure_trie_test.go",
		"storage/statedb/trie.go":                            "trie/trie.go",
		"storage/statedb/trie_test.go":                       "trie/trie_test.go",
		"storage/statedb/hasher.go":                          "trie/hasher.go",
		"storage/statedb/node.go":                            "trie/node.go",
		"storage/statedb/errors.go":                          "trie/errors.go",
		"blockchain/genesis_test.go":                         "core/genesis_test.go",
		"blockchain/helper_test.go":                          "core/helper_test.go",
		"blockchain/error.go":                                "core/error.go",
		"blockchain/tx_list_test.go":                         "core/tx_list_test.go",
		"blockchain/evm.go":                                  "core/evm.go",
		"blockchain/state_processor.go":                      "core/state_processor.go",
		"blockchain/types/log.go":                            "core/types/log.go",
		"blockchain/types/transaction_signing.go":            "core/types/transaction_signing.go",
		"blockchain/types/gen_header_json.go":                "core/types/gen_header_json.go",
		"blockchain/types/derive_sha.go":                     "core/types/derive_sha.go",
		"blockchain/types/transaction_test.go":               "core/types/transaction_test.go",
		"blockchain/types/log_test.go":                       "core/types/log_test.go",
		"blockchain/types/block_test.go":                     "core/types/block_test.go",
		"blockchain/types/gen_log_json.go":                   "core/types/gen_log_json.go",
		"blockchain/types/transaction.go":                    "core/types/transaction.go",
		"blockchain/types/receipt.go":                        "core/types/receipt.go",
		"blockchain/types/bloom_test.go":                     "core/types/bloom9_test.go",
		"blockchain/types/gen_tx_json.go":                    "core/types/gen_tx_json.go",
		"blockchain/types/transaction_signing_test.go":       "core/types/transaction_signing_test.go",
		"blockchain/types/bloom.go":                          "core/types/bloom9.go",
		"blockchain/types/gen_receipt_json.go":               "core/types/gen_receipt_json.go",
		"blockchain/types/block.go":                          "core/types/block.go",
		"blockchain/types/contract_ref.go":                   "core/state_transition.go",
		"blockchain/asm/compiler.go":                         "core/asm/compiler.go",
		"blockchain/asm/lex_test.go":                         "core/asm/lex_test.go",
		"blockchain/asm/lexer.go":                            "core/asm/lexer.go",
		"blockchain/asm/asm_test.go":                         "core/asm/asm_test.go",
		"blockchain/asm/asm.go":                              "core/asm/asm.go",
		"blockchain/asm/doc.go":                              "core/asm/asm.go",
		"blockchain/state_transition.go":                     "core/state_transition.go",
		"blockchain/chain_makers.go":                         "core/chain_makers.go",
		"blockchain/tx_pool_test.go":                         "core/tx_pool_test.go",
		"blockchain/metrics.go":                              "eth/metrics.go",
		"blockchain/block_validator_test.go":                 "core/block_validator_test.go",
		"blockchain/tx_journal.go":                           "core/tx_journal.go",
		"blockchain/genesis_alloc.go":                        "core/genesis_alloc.go",
		"blockchain/chain_indexer_test.go":                   "core/chain_indexer_test.go",
		"blockchain/types.go":                                "core/types.go",
		"blockchain/bad_blocks.go":                           "core/blocks.go",
		"blockchain/tx_list.go":                              "core/tx_list.go",
		"blockchain/gen_genesis_account.go":                  "core/gen_genesis_account.go",
		"blockchain/events.go":                               "core/events.go",
		"blockchain/gaspool.go":                              "core/gaspool.go",
		"blockchain/state/managed_state_test.go":             "core/state/managed_state_test.go",
		"blockchain/state/journal.go":                        "core/state/journal.go",
		"blockchain/state/statedb_test.go":                   "core/state/statedb_test.go",
		"blockchain/state/sync.go":                           "core/state/sync.go",
		"blockchain/state/managed_state.go":                  "core/state/managed_state.go",
		"blockchain/state/database.go":                       "core/state/database.go",
		"blockchain/state/state_object.go":                   "core/state/state_object.go",
		"blockchain/state/state_test.go":                     "core/state/state_test.go",
		"blockchain/state/dump.go":                           "core/state/dump.go",
		"blockchain/state/main_test.go":                      "core/state/main_test.go",
		"blockchain/state/statedb.go":                        "core/state/statedb.go",
		"blockchain/bloombits/scheduler_test.go":             "core/bloombits/scheduler_test.go",
		"blockchain/bloombits/matcher.go":                    "core/bloombits/matcher.go",
		"blockchain/bloombits/generator.go":                  "core/bloombits/generator.go",
		"blockchain/bloombits/matcher_test.go":               "core/bloombits/matcher_test.go",
		"blockchain/bloombits/generator_test.go":             "core/bloombits/generator_test.go",
		"blockchain/bloombits/scheduler.go":                  "core/bloombits/scheduler.go",
		"blockchain/blockchain.go":                           "core/blockchain.go",
		"blockchain/vm/memory.go":                            "core/vm/memory.go",
		"blockchain/vm/opcodes.go":                           "core/vm/opcodes.go",
		"blockchain/vm/analysis.go":                          "core/vm/analysis.go",
		"blockchain/vm/gas_table_test.go":                    "core/vm/gas_table_test.go",
		"blockchain/vm/gas_table.go":                         "core/vm/gas_table.go",
		"blockchain/vm/evm.go":                               "core/vm/evm.go",
		"blockchain/vm/gas.go":                               "core/vm/gas.go",
		"blockchain/vm/intpool_test.go":                      "core/vm/intpool_test.go",
		"blockchain/vm/logger.go":                            "core/vm/logger.go",
		"blockchain/vm/int_pool_verifier_empty.go":           "core/vm/int_pool_verifier_empty.go",
		"blockchain/vm/runtime/env.go":                       "core/vm/runtime/env.go",
		"blockchain/vm/runtime/runtime.go":                   "core/vm/runtime/runtime.go",
		"blockchain/vm/runtime/runtime_example_test.go":      "core/vm/runtime/runtime_example_test.go",
		"blockchain/vm/runtime/doc.go":                       "core/vm/runtime/doc.go",
		"blockchain/vm/runtime/runtime_test.go":              "core/vm/runtime/runtime_test.go",
		"blockchain/vm/interface.go":                         "core/vm/interface.go",
		"blockchain/vm/analysis_test.go":                     "core/vm/analysis_test.go",
		"blockchain/vm/instructions.go":                      "core/vm/instructions.go",
		"blockchain/vm/gen_structlog.go":                     "core/vm/gen_structlog.go",
		"blockchain/vm/contracts.go":                         "core/vm/contracts.go",
		"blockchain/vm/memory_table.go":                      "core/vm/memory_table.go",
		"blockchain/vm/instructions_test.go":                 "core/vm/instructions_test.go",
		"blockchain/vm/stack.go":                             "core/vm/stack.go",
		"blockchain/vm/common.go":                            "core/vm/common.go",
		"blockchain/vm/stack_table.go":                       "core/vm/stack_table.go",
		"blockchain/vm/interpreter.go":                       "core/vm/interpreter.go",
		"blockchain/vm/intpool.go":                           "core/vm/intpool.go",
		"blockchain/vm/jump_table.go":                        "core/vm/jump_table.go",
		"blockchain/vm/contract.go":                          "core/vm/contract.go",
		"blockchain/vm/contracts_test.go":                    "core/vm/contracts_test.go",
		"blockchain/vm/logger_test.go":                       "core/vm/logger_test.go",
		"blockchain/vm/errors.go":                            "core/vm/errors.go",
		"blockchain/vm/logger_json.go":                       "cmd/evm/json_logger.go",
		"blockchain/blockchain_test.go":                      "core/blockchain_test.go",
		"blockchain/mkalloc.go":                              "core/mkalloc.go",
		"blockchain/headerchain.go":                          "core/headerchain.go",
		"blockchain/chain_indexer.go":                        "core/chain_indexer.go",
		"blockchain/tx_pool.go":                              "core/tx_pool.go",
		"blockchain/block_validator.go":                      "core/block_validator.go",
		"blockchain/bench_test.go":                           "core/bench_test.go",
		"blockchain/genesis.go":                              "core/genesis.go",
		"blockchain/chain_makers_test.go":                    "core/chain_makers_test.go",
		"blockchain/tx_cacher.go":                            "core/tx_cacher.go",
		"blockchain/gen_genesis.go":                          "core/gen_genesis.go",
		"networks/p2p/discover/table_test.go":                "p2p/discover/table_test.go",
		"networks/p2p/discover/ntp.go":                       "p2p/discover/ntp.go",
		"networks/p2p/discover/database.go":                  "p2p/discover/database.go",
		"networks/p2p/discover/udp_test.go":                  "p2p/discover/udp_test.go",
		"networks/p2p/discover/database_test.go":             "p2p/discover/database_test.go",
		"networks/p2p/discover/udp.go":                       "p2p/discover/udp.go",
		"networks/p2p/discover/node_test.go":                 "p2p/discover/node_test.go",
		"networks/p2p/discover/node.go":                      "p2p/discover/node.go",
		"networks/p2p/discover/table.go":                     "p2p/discover/table.go",
		"networks/p2p/discover/discover_storage_kademlia.go": "p2p/discover/table.go",
		"networks/p2p/server.go":                             "p2p/server.go",
		"networks/p2p/metrics.go":                            "p2p/metrics.go",
		"networks/p2p/rlpx_test.go":                          "p2p/rlpx_test.go",
		"networks/p2p/dial_test.go":                          "p2p/dial_test.go",
		"networks/p2p/rlpx.go":                               "p2p/rlpx.go",
		"networks/p2p/nat/natpmp.go":                         "p2p/nat/natpmp.go",
		"networks/p2p/nat/natupnp.go":                        "p2p/nat/natupnp.go",
		"networks/p2p/nat/nat.go":                            "p2p/nat/nat.go",
		"networks/p2p/nat/nat_test.go":                       "p2p/nat/nat_test.go",
		"networks/p2p/nat/natupnp_test.go":                   "p2p/nat/natupnp_test.go",
		"networks/p2p/message.go":                            "p2p/message.go",
		"networks/p2p/protocol.go":                           "p2p/protocol.go",
		"networks/p2p/message_test.go":                       "p2p/message_test.go",
		"networks/p2p/simulations/mocker_test.go":            "p2p/simulations/mocker_test.go",
		"networks/p2p/simulations/simulation.go":             "p2p/simulations/simulation.go",
		"networks/p2p/simulations/http_test.go":              "p2p/simulations/http_test.go",
		"networks/p2p/simulations/network_test.go":           "p2p/simulations/network_test.go",
		"networks/p2p/simulations/mocker.go":                 "p2p/simulations/mocker.go",
		"networks/p2p/simulations/events.go":                 "p2p/simulations/events.go",
		"networks/p2p/simulations/pipes/pipes.go":            "p2p/simulations/pipes/pipes.go",
		"networks/p2p/simulations/adapters/docker.go":        "p2p/simulations/adapters/docker.go",
		"networks/p2p/simulations/adapters/types.go":         "p2p/simulations/adapters/types.go",
		"networks/p2p/simulations/adapters/exec.go":          "p2p/simulations/adapters/exec.go",
		"networks/p2p/simulations/adapters/inproc.go":        "p2p/simulations/adapters/inproc.go",
		"networks/p2p/simulations/adapters/ws.go":            "p2p/simulations/adapters/ws.go",
		"networks/p2p/simulations/adapters/ws_test.go":       "p2p/simulations/adapters/ws_test.go",
		"networks/p2p/simulations/adapters/inproc_test.go":   "p2p/simulations/adapters/inproc_test.go",
		"networks/p2p/simulations/network.go":                "p2p/simulations/network.go",
		"networks/p2p/simulations/http.go":                   "p2p/simulations/http.go",
		"networks/p2p/simulations/examples/ping-pong.go":     "p2p/simulations/examples/ping-pong.go",
		"networks/p2p/peer_error.go":                         "p2p/peer_error.go",
		"networks/p2p/peer.go":                               "p2p/peer.go",
		"networks/p2p/netutil/error.go":                      "p2p/netutil/error.go",
		"networks/p2p/netutil/net.go":                        "p2p/netutil/net.go",
		"networks/p2p/netutil/net_test.go":                   "p2p/netutil/net_test.go",
		"networks/p2p/netutil/toobig_notwindows.go":          "p2p/netutil/toobig_notwindows.go",
		"networks/p2p/netutil/error_test.go":                 "p2p/netutil/error_test.go",
		"networks/p2p/netutil/toobig_windows.go":             "p2p/netutil/toobig_windows.go",
		"networks/p2p/peer_test.go":                          "p2p/peer_test.go",
		"networks/p2p/dial.go":                               "p2p/dial.go",
		"networks/p2p/server_test.go":                        "p2p/server_test.go",
		"networks/rpc/ipc_unix.go":                           "rpc/ipc_unix.go",
		"networks/rpc/ipc_windows.go":                        "rpc/ipc_windows.go",
		"networks/rpc/subscription.go":                       "rpc/subscription.go",
		"networks/rpc/utils_test.go":                         "rpc/utils_test.go",
		"networks/rpc/server.go":                             "rpc/server.go",
		"networks/rpc/http_test.go":                          "rpc/http_test.go",
		"networks/rpc/types_test.go":                         "rpc/types_test.go",
		"networks/rpc/types.go":                              "rpc/types.go",
		"networks/rpc/client.go":                             "rpc/client.go",
		"networks/rpc/ipc.go":                                "rpc/ipc.go",
		"networks/rpc/subscription_test.go":                  "rpc/subscription_test.go",
		"networks/rpc/json_test.go":                          "rpc/json_test.go",
		"networks/rpc/http.go":                               "rpc/http.go",
		"networks/rpc/inproc.go":                             "rpc/inproc.go",
		"networks/rpc/utils.go":                              "rpc/utils.go",
		"networks/rpc/websocket.go":                          "rpc/websocket.go",
		"networks/rpc/client_example_test.go":                "rpc/client_example_test.go",
		"networks/rpc/client_test.go":                        "rpc/client_test.go",
		"networks/rpc/endpoints.go":                          "rpc/endpoints.go",
		"networks/rpc/json.go":                               "rpc/json.go",
		"networks/rpc/server_test.go":                        "rpc/server_test.go",
		"networks/rpc/errors.go":                             "rpc/errors.go",
		"networks/grpc/gClient.go":                           "rpc/json.go",
		"networks/grpc/gServer.go":                           "rpc/http.go",
		"common/mclock/mclock.go":                            "common/mclock/mclock.go",
		"common/hexutil/hexutil.go":                          "common/hexutil/hexutil.go",
		"common/hexutil/hexutil_test.go":                     "common/hexutil/hexutil_test.go",
		"common/hexutil/json_example_test.go":                "common/hexutil/json_example_test.go",
		"common/hexutil/json_test.go":                        "common/hexutil/json_test.go",
		"common/hexutil/json.go":                             "common/hexutil/json.go",
		"common/bitutil/compress_test.go":                    "common/bitutil/compress_test.go",
		"common/bitutil/bitutil_test.go":                     "common/bitutil/bitutil_test.go",
		"common/bitutil/compress.go":                         "common/bitutil/compress.go",
		"common/bitutil/bitutil.go":                          "common/bitutil/bitutil.go",
		"common/size_test.go":                                "common/size_test.go",
		"common/types_test.go":                               "common/types_test.go",
		"common/size.go":                                     "common/size.go",
		"common/types.go":                                    "common/types.go",
		"common/format.go":                                   "common/format.go",
		"common/math/integer_test.go":                        "common/math/integer_test.go",
		"common/math/integer.go":                             "common/math/integer.go",
		"common/math/big.go":                                 "common/math/big.go",
		"common/math/big_test.go":                            "common/math/big_test.go",
		"common/debug.go":                                    "common/debug.go",
		"common/utils.go":                                    "common/test_utils.go",
		"common/fdlimit/fdlimit_test.go":                     "common/fdlimit/fdlimit_test.go",
		"common/fdlimit/fdlimit_unix.go":                     "common/fdlimit/fdlimit_unix.go",
		"common/fdlimit/fdlimit_windows.go":                  "common/fdlimit/fdlimit_windows.go",
		"common/fdlimit/fdlimit_freebsd.go":                  "common/fdlimit/fdlimit_freebsd.go",
		"common/path.go":                                     "common/path.go",
		"common/compiler/solidity.go":                        "common/compiler/solidity.go",
		"common/compiler/solidity_test.go":                   "common/compiler/solidity_test.go",
		"common/bytes.go":                                    "common/bytes.go",
		"common/bytes_test.go":                               "common/bytes_test.go",
		"common/big.go":                                      "common/big.go",
		"common/main_test.go":                                "common/main_test.go",
		"accounts/accounts.go":                               "accounts/accounts.go",
		"accounts/hd.go":                                     "accounts/hd.go",
		"accounts/keystore/watch_fallback.go":                "accounts/keystore/watch_fallback.go",
		"accounts/keystore/account_cache_test.go":            "accounts/keystore/account_cache_test.go",
		"accounts/keystore/presale.go":                       "accounts/keystore/presale.go",
		"accounts/keystore/keystore_passphrase_test.go":      "accounts/keystore/keystore_passphrase_test.go",
		"accounts/keystore/key.go":                           "accounts/keystore/key.go",
		"accounts/keystore/file_cache.go":                    "accounts/keystore/file_cache.go",
		"accounts/keystore/keystore_test.go":                 "accounts/keystore/keystore_test.go",
		"accounts/keystore/keystore_passphrase.go":           "accounts/keystore/keystore_passphrase.go",
		"accounts/keystore/account_cache.go":                 "accounts/keystore/account_cache.go",
		"accounts/keystore/keystore_plain_test.go":           "accounts/keystore/keystore_plain_test.go",
		"accounts/keystore/keystore_plain.go":                "accounts/keystore/keystore_plain.go",
		"accounts/keystore/keystore_wallet.go":               "accounts/keystore/keystore_wallet.go",
		"accounts/keystore/keystore.go":                      "accounts/keystore/keystore.go",
		"accounts/keystore/watch.go":                         "accounts/keystore/watch.go",
		"accounts/abi/error.go":                              "accounts/abi/error.go",
		"accounts/abi/event.go":                              "accounts/abi/event.go",
		"accounts/abi/argument.go":                           "accounts/abi/argument.go",
		"accounts/abi/pack.go":                               "accounts/abi/pack.go",
		"accounts/abi/type.go":                               "accounts/abi/type.go",
		"accounts/abi/numbers.go":                            "accounts/abi/numbers.go",
		"accounts/abi/unpack_test.go":                        "accounts/abi/unpack_test.go",
		"accounts/abi/event_test.go":                         "accounts/abi/event_test.go",
		"accounts/abi/abi.go":                                "accounts/abi/abi.go",
		"accounts/abi/unpack.go":                             "accounts/abi/unpack.go",
		"accounts/abi/method.go":                             "accounts/abi/method.go",
		"accounts/abi/abi_test.go":                           "accounts/abi/abi_test.go",
		"accounts/abi/type_test.go":                          "accounts/abi/type_test.go",
		"accounts/abi/reflect.go":                            "accounts/abi/reflect.go",
		"accounts/abi/bind/backend.go":                       "accounts/abi/bind/backend.go",
		"accounts/abi/bind/auth.go":                          "accounts/abi/bind/auth.go",
		"accounts/abi/bind/util_test.go":                     "accounts/abi/bind/util_test.go",
		"accounts/abi/bind/backends/simulated.go":            "accounts/abi/bind/backends/simulated.go",
		"accounts/abi/bind/bind_test.go":                     "accounts/abi/bind/bind_test.go",
		"accounts/abi/bind/util.go":                          "accounts/abi/bind/util.go",
		"accounts/abi/bind/template.go":                      "accounts/abi/bind/template.go",
		"accounts/abi/bind/topics.go":                        "accounts/abi/bind/topics.go",
		"accounts/abi/bind/bind.go":                          "accounts/abi/bind/bind.go",
		"accounts/abi/bind/base.go":                          "accounts/abi/bind/base.go",
		"accounts/abi/numbers_test.go":                       "accounts/abi/numbers_test.go",
		"accounts/abi/pack_test.go":                          "accounts/abi/pack_test.go",
		"accounts/hd_test.go":                                "accounts/hd_test.go",
		"accounts/url.go":                                    "accounts/url.go",
		"accounts/url_test.go":                               "accounts/url_test.go",
		"accounts/manager.go":                                "accounts/manager.go",
		"accounts/errors.go":                                 "accounts/errors.go",
		"work/worker.go":                                     "miner/worker.go",
		"work/remote_agent.go":                               "miner/remote_agent.go",
		"work/unconfirmed.go":                                "miner/unconfirmed.go",
		"work/agent.go":                                      "miner/agent.go",
		"work/unconfirmed_test.go":                           "miner/unconfirmed_test.go",
		"work/work.go":                                       "miner/miner.go",
		"params/config.go":                                   "params/config.go",
		"params/version.go":                                  "params/version.go",
		"params/gas_table.go":                                "params/gas_table.go",
		"params/denomination.go":                             "params/denomination.go",
		"params/network_params.go":                           "params/network_params.go",
		"params/protocol_params.go":                          "params/protocol_params.go",
		"params/computation_cost_params.go":                  "params/protocol_params.go",
		"params/bootnodes.go":                                "params/bootnodes.go",
		"api/backend.go":                                     "internal/ethapi/backend.go",
		"api/addrlock.go":                                    "internal/ethapi/addrlock.go",
		"api/api_private_account.go":                         "internal/ethapi/api.go",
		"api/api_private_debug.go":                           "internal/ethapi/api.go",
		"api/api_public_account.go":                          "internal/ethapi/api.go",
		"api/api_public_blockchain.go":                       "internal/ethapi/api.go",
		"api/api_public_cypress.go":                          "internal/ethapi/api.go",
		"api/api_public_debug.go":                            "internal/ethapi/api.go",
		"api/api_public_klay.go":                             "internal/ethapi/api.go",
		"api/api_public_net.go":                              "internal/ethapi/api.go",
		"api/api_public_transaction_pool.go":                 "internal/ethapi/api.go",
		"api/api_public_tx_pool.go":                          "internal/ethapi/api.go",
		"api/debug/trace.go":                                 "internal/debug/trace.go",
		"api/debug/flags.go":                                 "internal/debug/flags.go",
		"api/debug/loudpanic_fallback.go":                    "internal/debug/loudpanic_fallback.go",
		"api/debug/loudpanic.go":                             "internal/debug/loudpanic.go",
		"api/debug/api.go":                                   "internal/debug/api.go",
		"api/debug/trace_fallback.go":                        "internal/debug/trace_fallback.go",
		"log/handler.go":                                     "log/handler.go",
		"log/handler_glog.go":                                "log/handler_glog.go",
		"log/handler_go13.go":                                "log/handler_go13.go",
		"log/handler_go14.go":                                "log/handler_go14.go",
		"log/log15_logger.go":                                "log/logger.go",
		"log/format.go":                                      "log/format.go",
		"log/handler_syslog.go":                              "log/syslog.go",
		"log/term/terminal_openbsd.go":                       "log/term/terminal_openbsd.go",
		"log/term/terminal_freebsd.go":                       "log/term/terminal_freebsd.go",
		"log/term/terminal_linux.go":                         "log/term/terminal_linux.go",
		"log/term/terminal_appengine.go":                     "log/term/terminal_appengine.go",
		"log/term/terminal_netbsd.go":                        "log/term/terminal_netbsd.go",
		"log/term/terminal_notwindows.go":                    "log/term/terminal_notwindows.go",
		"log/term/terminal_solaris.go":                       "log/term/terminal_solaris.go",
		"log/term/terminal_darwin.go":                        "log/term/terminal_darwin.go",
		"log/term/terminal_windows.go":                       "log/term/terminal_windows.go",
		"ser/rlp/raw.go":                                     "rlp/raw.go",
		"ser/rlp/decode_test.go":                             "rlp/decode_test.go",
		"ser/rlp/encode.go":                                  "rlp/encode.go",
		"ser/rlp/typecache.go":                               "rlp/typecache.go",
		"ser/rlp/raw_test.go":                                "rlp/raw_test.go",
		"ser/rlp/encoder_example_test.go":                    "rlp/encoder_example_test.go",
		"ser/rlp/decode_tail_test.go":                        "rlp/decode_tail_test.go",
		"ser/rlp/decode.go":                                  "rlp/decode.go",
		"ser/rlp/encode_test.go":                             "rlp/encode_test.go",
		"node/config.go":                                     "node/config.go",
		"node/utils_test.go":                                 "node/utils_test.go",
		"node/node_example_test.go":                          "node/node_example_test.go",
		"node/service.go":                                    "node/service.go",
		"node/cn/filters/filter.go":                          "eth/filters/filter.go",
		"node/cn/filters/api.go":                             "eth/filters/api.go",
		"node/cn/filters/filter_system_test.go":              "eth/filters/filter_system_test.go",
		"node/cn/filters/api_test.go":                        "eth/filters/api_test.go",
		"node/cn/filters/filter_system.go":                   "eth/filters/filter_system.go",
		"node/cn/handler.go":                                 "eth/handler.go",
		"node/cn/api_backend.go":                             "eth/api_backend.go",
		"node/cn/api_sc_backend.go":                          "eth/api_backend.go",
		"node/cn/api_tracer.go":                              "eth/api_tracer.go",
		"node/cn/config.go":                                  "eth/config.go",
		"node/cn/backend.go":                                 "eth/backend.go",
		"node/cn/sc_backend.go":                              "eth/backend.go",
		"node/cn/gasprice/gasprice.go":                       "eth/gasprice/gasprice.go",
		"node/cn/metrics.go":                                 "eth/metrics.go",
		"node/cn/sync.go":                                    "eth/sync.go",
		"node/cn/protocol.go":                                "eth/protocol.go",
		"node/cn/gen_config.go":                              "eth/gen_config.go",
		"node/cn/bloombits.go":                               "eth/bloombits.go",
		"node/cn/api.go":                                     "eth/api.go",
		"node/cn/api_sc.go":                                  "eth/api.go",
		"node/cn/peer.go":                                    "eth/peer.go",
		"node/cn/api_test.go":                                "eth/api_test.go",
		"node/sc/bridge_addr_journal.go":                     "core/tx_journal.go",
		"node/sc/sorted_map_list.go":                         "core/tx_list.go",
		"node/sc/bridge_tx_journal.go":                       "core/tx_journal.go",
		"node/sc/bridge_tx_list.go":                          "core/tx_list.go",
		"node/sc/bridge_tx_pool.go":                          "core/tx_pool.go",
		"node/sc/bridgepeer.go":                              "eth/peer.go",
		"node/sc/config.go":                                  "eth/config.go",
		"node/sc/mainbridge.go":                              "eth/backend.go",
		"node/sc/metrics.go":                                 "eth/metrics.go",
		"node/sc/protocol.go":                                "eth/protocol.go",
		"node/sc/subbridge.go":                               "eth/backend.go",
		"node/service_test.go":                               "node/service_test.go",
		"node/defaults.go":                                   "node/defaults.go",
		"node/api.go":                                        "node/api.go",
		"node/node_test.go":                                  "node/node_test.go",
		"node/config_test.go":                                "node/config_test.go",
		"node/node.go":                                       "node/node.go",
		"node/errors.go":                                     "node/errors.go",
		"event/subscription.go":                              "event/subscription.go",
		"event/example_test.go":                              "event/example_test.go",
		"event/event.go":                                     "event/event.go",
		"event/feed.go":                                      "event/feed.go",
		"event/event_test.go":                                "event/event_test.go",
		"event/example_feed_test.go":                         "event/example_feed_test.go",
		"event/subscription_test.go":                         "event/subscription_test.go",
		"event/example_scope_test.go":                        "event/example_scope_test.go",
		"event/feed_test.go":                                 "event/feed_test.go",
		"event/example_subscription_test.go":                 "event/example_subscription_test.go",
		"event/filter/generic_filter.go":                     "event/filter/generic_filter.go",
		"event/filter/filter.go":                             "event/filter/filter.go",
		"client/klay_client.go":                              "ethclient/ethclient.go",
		"client/klay_client_test.go":                         "ethclient/ethclient_test.go",
		"client/signer.go":                                   "ethclient/signer.go",
		"client/bridge_client.go":                            "ethclient/ethclient.go",
		"console/prompter.go":                                "console/prompter.go",
		"console/console.go":                                 "console/console.go",
		"console/bridge.go":                                  "console/bridge.go",
		"console/console_test.go":                            "console/console_test.go",
		"console/web3ext/web3ext.go":                         "internal/web3ext/web3ext.go",
		"console/jsre/completion.go":                         "internal/jsre/completion.go",
		"console/jsre/completion_test.go":                    "internal/jsre/completion_test.go",
		"console/jsre/jsre.go":                               "internal/jsre/jsre.go",
		"console/jsre/jsre_test.go":                          "internal/jsre/jsre_test.go",
		"console/jsre/pretty.go":                             "internal/jsre/pretty.go",
		"console/jsre/deps/deps.go":                          "internal/jsre/deps/deps.go",
		"console/jsre/deps/bindata.go":                       "internal/jsre/deps/bindata.go",
	}
)

// this template generates the license comment.
// its input is an info structure.
var klaytnLicenseT = template.Must(template.New("").Parse(`
// Copyright {{.Year}} The klaytn Authors
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
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.`[1:]))

// this template generates the license comment.
// its input is an info structure.
var mixedLicenseT = template.Must(template.New("").Parse(`
// Modifications Copyright {{.Year}} The klaytn Authors
{{.OtherLicence}}
//
// This file is derived from {{.File}} (2018/06/04).
// Modified and improved for the klaytn development.`[1:]))

// this template generates the license comment.
// its input is an info structure.
var externalLicenseT = template.Must(template.New("").Parse(`
// Copyright {{.Year}} The klaytn Authors
//
// This file is derived from {{.File}} (2018/06/04).
// See LICENSE in the top directory for the original copyright and license.`[1:]))

type info struct {
	file         string
	Year         int64
	LastCommit   int64 //unit: YYYY
	otherLicence string
}

//File get derived ethereum path
func (i info) File() string {
	etherdir, exists := ethereumDir[i.file]
	if !exists {
		fmt.Printf("Check ethereum dir : %s \n", i.file)
	}
	return etherdir
}

//get original ethereum licence
func (i info) OtherLicence() string {
	return i.otherLicence
}

func main() {
	var (
		files = getFiles()
		filec = make(chan string)
		infoc = make(chan *info, 20)
		wg    sync.WaitGroup
	)

	//writeAuthors(files)

	go func() {
		for _, f := range files {
			filec <- f
		}
		close(filec)
	}()
	for i := runtime.NumCPU(); i >= 0; i-- {
		// getting file info is slow and needs to be parallel.
		// it traverses git history for each file.
		wg.Add(1)
		go getInfo(filec, infoc, &wg)
	}
	go func() {
		wg.Wait()
		close(infoc)
	}()
	writeLicenses(infoc)
}

//skipFile returns whether the path updates the license
func skipFile(path string) bool {
	if strings.Contains(path, "/testdata/") {
		return true
	}
	for _, p := range skipPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	for _, p := range skipSuffixes {
		if strings.HasSuffix(path, p) {
			return true
		}
	}

	return false
}

//externalLiceceFile returns whether the path has external license
func externalLicenceFile(path string) bool {
	for _, p := range externalLicencePrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

//getFiles returns all klaytn files
func getFiles() []string {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", "HEAD")
	var files []string
	err := doLines(cmd, func(line string) {
		if skipFile(line) {
			return
		}
		ext := filepath.Ext(line)
		for _, wantExt := range extensions {
			if ext == wantExt {
				goto keep
			}
		}
		return
	keep:
		files = append(files, line)
	})
	if err != nil {
		log.Fatal("error getting files:", err)
	}
	return files
}

var authorRegexp = regexp.MustCompile(`\s*[0-9]+\s*(.*)`)

// gitAuthors returns git authors from files
func gitAuthors(files []string) []string {
	cmds := []string{"shortlog", "-s", "-n", "-e", "HEAD", "--"}
	cmds = append(cmds, files...)
	cmd := exec.Command("git", cmds...)
	var authors []string
	err := doLines(cmd, func(line string) {
		m := authorRegexp.FindStringSubmatch(line)
		if len(m) > 1 {
			authors = append(authors, m[1])
		}
	})
	if err != nil {
		log.Fatalln("error getting authors:", err)
	}
	return authors
}

// gitAuthors returns git authors from AUTHORS file
func readAuthors() []string {
	content, err := ioutil.ReadFile("AUTHORS")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln("error reading AUTHORS:", err)
	}
	var authors []string
	for _, a := range bytes.Split(content, []byte("\n")) {
		if len(a) > 0 && a[0] != '#' {
			authors = append(authors, string(a))
		}
	}
	// Retranslate existing authors through .mailmap.
	// This should catch email address changes.
	authors = mailmapLookup(authors)
	return authors
}

func mailmapLookup(authors []string) []string {
	if len(authors) == 0 {
		return nil
	}
	cmds := []string{"check-mailmap", "--"}
	cmds = append(cmds, authors...)
	cmd := exec.Command("git", cmds...)
	var translated []string
	err := doLines(cmd, func(line string) {
		translated = append(translated, line)
	})
	if err != nil {
		log.Fatalln("error translating authors:", err)
	}
	return translated
}

func writeAuthors(files []string) {
	merge := make(map[string]bool)
	// Add authors that Git reports as contributorxs.
	// This is the primary source of author information.
	for _, a := range gitAuthors(files) {
		merge[a] = true
	}
	// Add existing authors from the file. This should ensure that we
	// never lose authors, even if Git stops listing them. We can also
	// add authors manually this way.
	for _, a := range readAuthors() {
		merge[a] = true
	}
	// Write sorted list of authors back to the file.
	var result []string
	for a := range merge {
		result = append(result, a)
	}
	sort.Strings(result)
	content := new(bytes.Buffer)
	content.WriteString(authorsFileHeader)
	for _, a := range result {
		content.WriteString(a)
		content.WriteString("\n")
	}
	fmt.Println("writing AUTHORS")
	if err := ioutil.WriteFile("AUTHORS", content.Bytes(), 0644); err != nil {
		log.Fatalln(err)
	}
}

func getInfo(files <-chan string, out chan<- *info, wg *sync.WaitGroup) {
	for file := range files {
		fmt.Println(file)
		stat, err := os.Lstat(file)
		if err != nil {
			fmt.Printf("ERROR %s: %v\n", file, err)
			continue
		}
		if !stat.Mode().IsRegular() {
			continue
		}
		if isGenerated(file) {
			continue
		}
		info, err := fileInfo(file)
		if err != nil {
			fmt.Printf("ERROR %s: %v\n", file, err)
			continue
		}
		if info.LastCommit <= 2018 {
			continue
		}
		out <- info
	}
	wg.Done()
}

// isGenerated returns whether the input file is an automatically generated file.
func isGenerated(file string) bool {
	fd, err := os.Open(file)
	if err != nil {
		return false
	}
	defer fd.Close()
	buf := make([]byte, 2048)
	n, _ := fd.Read(buf)
	buf = buf[:n]
	for _, l := range bytes.Split(buf, []byte("\n")) {
		if bytes.HasPrefix(l, []byte("// Code generated")) {
			return true
		}
	}
	return false
}

// fileInfo finds the lowest year in which the given file was committed.
func fileInfo(file string) (*info, error) {
	info := &info{file: file, Year: int64(time.Now().Year()), LastCommit: 0}
	cmd := exec.Command("git", "log", "--follow", "--find-renames=80", "--find-copies=80", "--pretty=format:%ai", "--", file)
	err := doLines(cmd, func(line string) {
		y, err := strconv.ParseInt(line[:4], 10, 64)
		if err != nil {
			fmt.Printf("cannot parse year: %q", line[:4])
		}
		if y < info.Year {
			info.Year = y
		}
		if y > info.LastCommit {
			info.LastCommit = y
		}
	})
	return info, err
}

// writeLicenses write the licenses in the files path in infos.
func writeLicenses(infos <-chan *info) {
	for i := range infos {
		writeLicense(i)
	}
}

// writeLicense write the license in the file path in info.
func writeLicense(info *info) {
	fi, err := os.Stat(info.file)
	if os.IsNotExist(err) {
		fmt.Println("skipping (does not exist)", info.file)
		return
	}
	if err != nil {
		log.Fatalf("error stat'ing %s: %v\n", info.file, err)
	}
	content, err := ioutil.ReadFile(info.file)
	if err != nil {
		log.Fatalf("error reading %s: %v\n", info.file, err)
	}
	// Construct new file content.
	buf := new(bytes.Buffer)

	if m := mixedLicenseCommentRE.FindIndex(content); m != nil {
		if e := ethereumLicenseCommentRE.FindIndex(content); e != nil {
			info.otherLicence = string(content[e[0]:e[1]])
			mixedLicenseT.Execute(buf, info)
			buf.Write(content[m[1]:])
		} else {
			return
		}
	} else if m := ethereumLicenseCommentRE.FindIndex(content); m != nil {
		info.otherLicence = string(content[m[0]:m[1]])
		mixedLicenseT.Execute(buf, info)
		buf.Write(content[m[1]:])
	} else if m := klaytnLicenseCommentRE.FindIndex(content); m != nil {
		klaytnLicenseT.Execute(buf, info)
		buf.Write(content[m[1]:])
	} else if m := externalLicenceCommentRE.FindIndex(content); m != nil {
		externalLicenseT.Execute(buf, info)
		buf.Write(content[m[1]:])
	} else {
		if externalLicenceFile(info.file) {
			externalLicenseT.Execute(buf, info)
		} else {
			klaytnLicenseT.Execute(buf, info)
		}
		buf.Write([]byte("\n\n"))
		buf.Write(content)
	}

	// Write it to the file.
	if bytes.Equal(content, buf.Bytes()) {
		fmt.Println("skipping (no changes)", info.file)
		return
	}
	if err := ioutil.WriteFile(info.file, buf.Bytes(), fi.Mode()); err != nil {
		log.Fatalf("error writing %s: %v", info.file, err)
	}
}

// dolines executes cmd
func doLines(cmd *exec.Cmd, f func(string)) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	s := bufio.NewScanner(stdout)
	for s.Scan() {
		f(s.Text())
	}
	if s.Err() != nil {
		return s.Err()
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%v (for %s)", err, strings.Join(cmd.Args, " "))
	}
	return nil
}
