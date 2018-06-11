
## Go GXPlatform

Official golang implementation of the GXPlatform protocol.

## Building the source

Building gxp requires both a Go (version 1.7 or later) and a C compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make gxp   or  make all

## Executables

The go-gxplatform project comes with several wrappers/executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **`gxp`** | Our main GXP CLI client. It is the entry point into the GXP network (main-, test- or private net), capable of running as a full node (default) archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the GXP network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `gxp --help`|

## Running gxp

Going through all the possible command line flags is out of scope here, but we've
enumerated a few common parameter combos to get you up to speed quickly on how you can run your
own GXP instance.

### Full node on the main GXPlatform network

By far the most common scenario is people wanting to simply interact with the GXPlatform network:
create accounts; transfer funds; deploy and interact with contracts. For this particular use-case
the user doesn't care about years-old historical data, so we can fast-sync quickly to the current
state of the network. To do so:

```
$ gxp console
```

This command will:

 * Start gxp in fast sync mode (default, can be changed with the `--syncmode` flag), causing it to
   download more data in exchange for avoiding processing the entire history of the GXPlatform network,
   which is very CPU intensive.

### Full node on the GXPlatform test network

Transitioning towards developers, if you'd like to play around with creating GXPlatform contracts, you
almost certainly would like to do that without any real money involved until you get the hang of the
entire system. In other words, instead of attaching to the main network, you want to join the **test**
network with your node, which is fully equivalent to the main network, but with play-GXP only.

```
$ gxp --testnet console
```

The `console` subcommand have the exact same meaning as above and they are equally useful on the
testnet too. Please see above for their explanations if you've skipped to here.

Specifying the `--testnet` flag however will reconfigure your GXP instance a bit:

 * Instead of using the default data directory (`~/.gxp` on Linux for example), GXP will nest
   itself one level deeper into a `testnet` subfolder (`~/.gxp/testnet` on Linux). Note, on OSX
   and Linux this also means that attaching to a running testnet node requires the use of a custom
   endpoint since `gxp attach` will try to attach to a production node endpoint by default. E.g.
   `gxp attach <datadir>/testnet/gxp.ipc`. Windows users are not affected by this.
 * Instead of connecting the main GXPlatform network, the client will connect to the test network,
   which uses different P2P bootnodes, different network IDs and genesis states.
   
*Note: Although there are some internal protective measures to prevent transactions from crossing
over between the main network and test network, you should make sure to always use separate accounts
for play-money and real-money. Unless you manually move accounts, GXP will by default correctly
separate the two networks and will not make any accounts available between them.*

```

### Configuration (PoW)

As an alternative to passing the numerous flags to the `gxp` binary, you can also pass a configuration file via:

```
$ gxp --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to export your existing configuration:

```
$ gxp --your-favourite-flags dumpconfig
```


Do not forget `--rpcaddr 0.0.0.0`, if you want to access RPC from other containers and/or hosts. By default, `gxp` binds to the local interface and RPC endpoints is not accessible from the outside.

### Programatically interfacing GXP nodes

As a developer, sooner rather than later you'll want to start interacting with GXP and the GXPlatform
network via your own programs and not manually through the console. To aid this, GXP has built-in
support for a JSON-RPC based APIs. These can be
exposed via HTTP, WebSockets and IPC (unix sockets on unix based platforms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by GXP, whereas the HTTP
and WS interfaces need to manually be enabled and only expose a subset of APIs due to security reasons.
These can be turned on/off and configured as you'd expect.

HTTP based JSON-RPC API options:

  * `--rpc` Enable the HTTP-RPC server
  * `--rpcaddr` HTTP-RPC server listening interface (default: "localhost")
  * `--rpcport` HTTP-RPC server listening port (default: 8545)
  * `--rpcapi` API's offered over the HTTP-RPC interface (default: "gxp,net,web3")
  * `--rpccorsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--wsaddr` WS-RPC server listening interface (default: "localhost")
  * `--wsport` WS-RPC server listening port (default: 8546)
  * `--wsapi` API's offered over the WS-RPC interface (default: "gxp,net,web3")
  * `--wsorigins` Origins from which to accept websockets requests
  * `--ipcdisable` Disable the IPC-RPC server
  * `--ipcapi` API's offered over the IPC-RPC interface (default: "admin,gxp,miner,net,personal,txpool,web3")
  * `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to connect
via HTTP, WS or IPC to a GXP node configured with the above flags and you'll need to speak [JSON-RPC](http://www.jsonrpc.org/specification)
on all transports. You can reuse the same connection for multiple requests!

### Operating a private network

Maintaining your own private network is more involved as a lot of configurations taken for granted in
the official networks need to be manually set up.

#### Defining the private genesis state

First, you'll need to create the genesis state of your networks, which all nodes need to be aware of
and agree upon. This consists of a small JSON file (e.g. call it `genesis.json`):

```json
{
  "config": {
        "chainId": 0,
        "homesteadBlock": 0
    },
  "alloc"      : {},
  "coinbase"   : "0x0000000000000000000000000000000000000000",
  "difficulty" : "0x20000",
  "extraData"  : "",
  "gasLimit"   : "0x2fefd8",
  "nonce"      : "0x0000000000000042",
  "mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp"  : "0x00"
}
```

The above fields should be fine for most purposes, although we'd recommend changing the `nonce` to
some random value so you prevent unknown remote nodes from being able to connect to you. If you'd
like to pre-fund some accounts for easier testing, you can populate the `alloc` field with account
configs:

```json
"alloc": {
  "0x0000000000000000000000000000000000000001": {"balance": "111111111"},
  "0x0000000000000000000000000000000000000002": {"balance": "222222222"}
}
```

With the genesis state defined in the above JSON file, you'll need to initialize **every** GXP node
with it prior to starting it up to ensure all blockchain parameters are correctly set:

```
$ gxp init path/to/genesis.json
```

### Configuration (istanbul-BFT)
When --nodes --verbose flags are given, a static-nodes.json template as well as the validators' node keys, public keys, and addresses are generated.
When --save flag is given, all generated configs will be saved. Use BFT when --bft flag is given.

Note: the generated static-nodes.json template are set with IP 0.0.0.0, please make according change to match your environment.
```
./istanbul setup --num 4 --nodes --verbose --bft
validators
{
	"Address": "0x7871bce4a383c0c01c81e0f7e1798d97b49cd64c",
	"Nodekey": "81a74dc939a2e023a3743396cc9beb04cc092e11aceadf07e5f5e4299bb9a8c6",
	"NodeInfo": "enode://fd6ea50203d1603ab2e2d26b48dc92feba58c1867e6abfbaf2193739fa23feddc61d2a7eb788ea3ec6f09ba86299df1ead14ddb3a3a43856203e0961631553ed@0.0.0.0:30303?discport=0"
}
{
	"Address": "0x7e1b0308c8baaf40f4614f5ca94aa0d1acfd645c",
	"Nodekey": "6e29d0af0f10f73dcd1547cf129ed99730188d2ea8b4eeaf2933ec15c13e9fef",
	"NodeInfo": "enode://ad0c9c871b4fd631e6a92442ec22b8b39108724bc2f48bffefdbc8afa776c5348c7a724cf0b003edf6bfe60fd36c59221d11478698c4dd093836198ca1c870c5@0.0.0.0:30303?discport=0"
}
{
	"Address": "0x9148c161d204c17bf0ff0df222242bc822663ffe",
	"Nodekey": "2fcb8cb40a2cc8262140566dd5aa8fb38e455263707ddac1be1c71becf884b79",
	"NodeInfo": "enode://2da0b0abfe22814790f4bd5925bcb4fb89c950bec9612ecff773417945c3899ca64094d331baef5196c6ba2744163991530978d50f05225b33c4a7d81eaee7b1@0.0.0.0:30303?discport=0"
}
{
	"Address": "0xf13783ba2087067663413cc2561262bdd892b357",
	"Nodekey": "e49b505d4a6d2b29d3411b5568eba9d4da0c48499cd6c30c9aae3db3822cc322",
	"NodeInfo": "enode://315a62772a05386f15d66ea7a39b1500d3749fdeea847811b4f8da842049c1ec7ccfe6db0f518415d21d713d92448c251793ebea9e1a6baf9b03502964737670@0.0.0.0:30303?discport=0"
}

static-nodes.json
[
	"enode://fd6ea50203d1603ab2e2d26b48dc92feba58c1867e6abfbaf2193739fa23feddc61d2a7eb788ea3ec6f09ba86299df1ead14ddb3a3a43856203e0961631553ed@0.0.0.0:30303?discport=0",
	"enode://ad0c9c871b4fd631e6a92442ec22b8b39108724bc2f48bffefdbc8afa776c5348c7a724cf0b003edf6bfe60fd36c59221d11478698c4dd093836198ca1c870c5@0.0.0.0:30303?discport=0",
	"enode://2da0b0abfe22814790f4bd5925bcb4fb89c950bec9612ecff773417945c3899ca64094d331baef5196c6ba2744163991530978d50f05225b33c4a7d81eaee7b1@0.0.0.0:30303?discport=0",
	"enode://315a62772a05386f15d66ea7a39b1500d3749fdeea847811b4f8da842049c1ec7ccfe6db0f518415d21d713d92448c251793ebea9e1a6baf9b03502964737670@0.0.0.0:30303?discport=0"
]


genesis.json
{
    "config": {
        "chainId": 2017,
        "homesteadBlock": 1,
        "eip155Block": 3,
        "istanbul": {
            "epoch": 30000,
            "policy": 0
        },
        "isBFT": true
    },
    "nonce": "0x0",
    "timestamp": "0x5b1e64b0",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f89af854947871bce4a383c0c01c81e0f7e1798d97b49cd64c947e1b0308c8baaf40f4614f5ca94aa0d1acfd645c949148c161d204c17bf0ff0df222242bc822663ffe94f13783ba2087067663413cc2561262bdd892b357b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0",
    "gasLimit": "0x82f79cd9000",
    "difficulty": "0x1",
    "mixHash": "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
    "coinbase": "0x0000000000000000000000000000000000000000",
    "alloc": {
        "7871bce4a383c0c01c81e0f7e1798d97b49cd64c": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "7e1b0308c8baaf40f4614f5ca94aa0d1acfd645c": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "9148c161d204c17bf0ff0df222242bc822663ffe": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        },
        "f13783ba2087067663413cc2561262bdd892b357": {
            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
        }
    },
    "number": "0x0",
    "gasUsed": "0x0",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}

```
istanbul의 setup 명령어로 만들어진 nodekey를 각 노드의 nodekey 화일에 저장한다.

```nodekey
81a74dc939a2e023a3743396cc9beb04cc092e11aceadf07e5f5e4299bb9a8c6
```

미리 생성된 계정이 있는 경우에는 genesis.json 파일에 alloc 설정에 계정에 대한 balance 값을 추가함

nodekey 화일과 static-nodes.json 화일을 datadir로 복사한다.
Genesis 블락을 다음의 명령어로 생성한다.
```
./gxp --datadir $DATAPATH init genesis.json

```
각 노드를 다음 명령어로 실행함. 실행시에 --mine 옵션으로 마이닝을 실행하고 --gasprice 0 옵션으로 gasprice를 0으로 세팅함.
```
gxp --datadir $DATAPATH --port 30303 --rpc --rpcaddr 0.0.0.0 --rpcport "8123" --rpccorsdomain "*"
--nodiscover --networkid 3900 --nat "any" --wsport "8546" --ws --wsaddr 0.0.0.0 --wsorigins="*"
--rpcapi "db,txpool,gxp,net,web3,miner,personal,admin,rpc" --mine --gasprice 0 console

```
동일 머신에서 수행시에는 --port , --rpcport --wsport 옵션을 다르게 설정하고 --networkid는 동일하게 설정함


## License

The go-gxplatform library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The go-gxplatform binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
