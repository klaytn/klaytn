// Copyright 2018 The klaytn Authors
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

/*
kscn is the command-line client for Klaytn ServiceChain Node.

kscn has the node type of "cn" internally and following commands and options are available.

COMMANDS:
   account     Manage accounts
   attach      Start an interactive JavaScript environment (connect to node)
   console     Start an interactive JavaScript environment
   dumpconfig  Show configuration values
   init        Bootstrap and initialize a new genesis block
   version     Show version number
   help, h     Shows a list of commands or help for one command

KLAY OPTIONS:
  --dbtype value                           Blockchain storage database type ("leveldb", "badger") (default: "leveldb")
  --datadir "/Users/andylee/Library/KSCN"  Data directory for the databases and keystore
  --keystore                               Directory for the keystore (default = inside the datadir)
  --identity value                         Custom node name
  --syncmode "full"                        Blockchain sync mode (only "full" is supported)
  --gcmode value                           Blockchain garbage collection mode ("full", "archive") (default: "full")
  --lightkdf                               Reduce key-derivation RAM & CPU usage at some expense of KDF strength
  --srvtype value                          json rpc server type ("http", "fasthttp") (default: "fasthttp")
  --extradata value                        Block extra data set by the work (default = client version)
  --config value                           TOML configuration file

SERVICECHAIN OPTIONS:
  --scsigner value            Public address for signing blocks in the service chain (default = first account created) (default: "0")
  --childchainindexing        Enables storing transaction hash of child chain transaction for fast access to child chain data
  --mainbridge                Enable main bridge service for service chain
  --mainbridgeport value      main bridge listen port (default: 50505)
  --subbridge                 Enable sub bridge service for service chain
  --subbridgeport value       sub bridge listen port (default: 50506)
  --chaintxperiod value       The period to make and send a chain transaction to the parent chain (default: 1)
  --chaintxlimit value        Number of service chain transactions stored for resending (default: 100)
  --parentchainid value       parent chain ID (default: 8217)
  --vtrecovery                Enable value transfer recovery (default: false)
  --vtrecoveryinterval value  Set the value transfer recovery interval (seconds) (default: 60)
  --scnewaccount              Enable account creation for the service chain (default: false). If set true, generated account can't be synced with the parent chain.
  --scconsensus value         Set the service chain consensus ("istanbul", "clique") (default: "istanbul")

ACCOUNT OPTIONS:
  --unlock value    Comma separated list of accounts to unlock
  --password value  Password file to use for non-interactive password input

TXPOOL OPTIONS:
  --txpool.nolocals                     Disables price exemptions for locally submitted transactions
  --txpool.journal value                Disk journal for local transaction to survive node restarts (default: "transactions.rlp")
  --txpool.journal-interval value       Time interval to regenerate the local transaction journal (default: 1h0m0s)
  --txpool.pricelimit value             Minimum gas price limit to enforce for acceptance into the pool (default: 1)
  --txpool.pricebump value              Price bump percentage to replace an already existing transaction (default: 10)
  --txpool.exec-slots.account value     Number of executable transaction slots guaranteed per account (default: 16)
  --txpool.exec-slots.all value         Maximum number of executable transaction slots for all accounts (default: 4096)
  --txpool.nonexec-slots.account value  Maximum number of non-executable transaction slots permitted per account (default: 64)
  --txpool.nonexec-slots.all value      Maximum number of non-executable transaction slots for all accounts (default: 1024)
  --txpool.lifetime value               Maximum amount of time non-executable transaction are queued (default: 5m0s)

DATABASE OPTIONS:
  --db.leveldb.cache-size value        Size of in-memory cache in LevelDB (MiB) (default: 768)
  --db.no-partitioning                 Disable partitioned databases for persistent storage
  --db.num-statetrie-partitions value  Number of internal partitions of state trie partition. Should be power of 2 (default: 4)
  --db.leveldb.compression value       Determines the compression method for LevelDB. 0=AllNoCompression, 1=ReceiptOnlySnappyCompression, 2=StateTrieOnlyNoCompression, 3=AllSnappyCompression (default: 0)
  --db.leveldb.no-buffer-pool          Disables using buffer pool for LevelDB's block allocation
  --db.no-parallel-write               Disables parallel writes of block data to persistent database
  --sendertxhashindexing               Enables storing mapping information of senderTxHash to txHash

STATE OPTIONS:
  --statedb.use-cache           Enables caching of state objects in stateDB
  --state.cache-size value      Size of in-memory cache of the global state (in MiB) to flush matured singleton trie nodes to disk (default: 512)
  --state.block-interval value  An interval in terms of block number to commit the global state to disk (default: 128)

CACHE OPTIONS:
  --cache.type value              Cache Type: 0=LRUCache, 1=LRUShardCache, 2=FIFOCache (default: 2)
  --cache.scale value             Scale of cache (cache size = preset size * scale of cache(%)) (default: 0)
  --cache.level value             Set the cache usage level ('saving', 'normal', 'extreme')
  --cache.memory value            Set the physical RAM size (GB, Default: 16GB) (default: 0)
  --cache.writethrough            Enables write-through writing to database and cache for certain types of cache.
  --statedb.use-txpool-cache      Enables caching of nonce and balance for txpool.
  --state.trie-cache-limit value  Memory allowance (MB) to use for caching trie nodes in memory (default: 4096)

CONSENSUS OPTIONS:
  --rewardbase value  Public address for block consensus rewards (default = first account created) (default: "0")

NETWORKING OPTIONS:
  --bootnodes value        Comma separated kni URLs for P2P discovery bootstrap
  --port value             Network listening port (default: 32323)
  --subport value          Network sub listening port (default: 32324)
  --multichannel           Create a dedicated channel for block propagation
  --maxconnections value   Maximum number of physical connections. All single channel peers can be maxconnections peers. All multi channel peers can be maxconnections/2 peers. (network disabled if set to 0) (default: 10)
  --maxpendpeers value     Maximum number of pending connection attempts (defaults used if set to 0) (default: 0)
  --targetgaslimit value   Target gas limit sets the artificial target gas floor for the blocks to mine (default: 4712388)
  --nat value              NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>) (default: "any")
  --nodiscover             Disables the peer discovery mechanism (manual peer addition)
  --rwtimerwaittime value  Wait time the rw timer waits for message writing (default: 15s)
  --rwtimerinterval value  Interval of using rw timer to check if it works well (default: 1000)
  --netrestrict value      Restricts network communication to the given IP network (CIDR masks)
  --nodekey value          P2P node key file
  --nodekeyhex value       P2P node key as hex (for testing)
  --networkid value        Network identifier (integer, 1=MainNet (Not yet launched), 1000=Aspen, 1001=Baobab) (default: 8217)

METRICS OPTIONS:
  --metrics               Enable metrics collection and reporting
  --prometheus            Enable prometheus exporter
  --prometheusport value  Prometheus exporter listening port (default: 61001)

VIRTUAL MACHINE OPTIONS:
  --vmdebug      Record information useful for VM and contract debugging
  --vmlog value  Set the output target of vmlog precompiled contract (0: no output, 1: file, 2: stdout, 3: both) (default: 0)

API AND CONSOLE OPTIONS:
  --rpc                  Enable the HTTP-RPC server
  --rpcaddr value        HTTP-RPC server listening interface (default: "localhost")
  --rpcport value        HTTP-RPC server listening port (default: 8551)
  --rpccorsdomain value  Comma separated list of domains from which to accept cross origin requests (browser enforced)
  --rpcvhosts value      Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default: "localhost")
  --rpcapi value         API's offered over the HTTP-RPC interface
  --ipcdisable           Disable the IPC-RPC server
  --ipcpath              Filename for IPC socket/pipe within the datadir (explicit paths escape it)
  --ws                   Enable the WS-RPC server
  --wsaddr value         WS-RPC server listening interface (default: "localhost")
  --wsport value         WS-RPC server listening port (default: 8552)
  --wsapi value          API's offered over the WS-RPC interface
  --wsorigins value      Origins from which to accept websockets requests
  --grpc                 Enable the gRPC server
  --grpcaddr value       gRPC server listening interface (default: "localhost")
  --grpcport value       gRPC server listening port (default: 8553)
  --jspath loadScript    JavaScript root path for loadScript (default: ".")
  --exec value           Execute JavaScript statement
  --preload value        Comma separated list of JavaScript files to preload into the console

LOGGING AND DEBUGGING OPTIONS:
  --verbosity value         Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail (default: 3)
  --vmodule value           Per-module verbosity: comma-separated list of <pattern>=<level> (e.g. klay/*=5,p2p=4)
  --backtrace value         Request a stack trace at a specific logging statement (e.g. "block.go:271")
  --debug                   Prepends log messages with call-site location (file and line number)
  --pprof                   Enable the pprof HTTP server
  --pprofaddr value         pprof HTTP server listening interface (default: "127.0.0.1")
  --pprofport value         pprof HTTP server listening port (default: 6060)
  --memprofile value        Write memory profile to the given file
  --memprofilerate value    Turn on memory profiling with the given rate (default: 524288)
  --blockprofilerate value  Turn on block profiling with the given rate (default: 0)
  --cpuprofile value        Write CPU profile to the given file
  --trace value             Write execution trace to the given file

MISC OPTIONS:
  --genkey value  generate a node private key and write to given filename
  --writeaddress  write out the node's public key which is given by "--nodekeyfile" or "--nodekeyhex"
  --help, -h      show help

*/
package main
