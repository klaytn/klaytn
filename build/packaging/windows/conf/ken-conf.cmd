REM Configuration file for the ken

set NETWORK_ID=8217

set PORT=32323

set SERVER_TYPE="fasthttp"
set SYNCMODE="full"
set VERBOSITY=3
set MAXCONNECTIONS=10

REM txpool options setting
set TXPOOL_EXEC_SLOTS_ALL=4096
set TXPOOL_NONEXEC_SLOTS_ALL=4096
set TXPOOL_EXEC_SLOTS_ACCOUNT=4096
set TXPOOL_NONEXEC_SLOTS_ACCOUNT=4096
set TXPOOL_LIFE_TIME="30m"

REM rpc options setting
set RPC_ENABLE=1 &:: if this is set, the following options will be used
set RPC_API="klay" &:: available apis: admin,debug,klay,eth,miner,net,personal,rpc,txpool,web3
set RPC_PORT=8551
set RPC_ADDR="0.0.0.0"
set RPC_CORSDOMAIN="*"
set RPC_VHOSTS="*"

REM ws options setting
set WS_ENABLE=1 &:: if this is set, the following options will be used
set WS_API="klay" &:: available apis: admin,debug,klay,eth,miner,net,personal,rpc,txpool,web3
set WS_ADDR="0.0.0.0"
set WS_PORT=8552
set WS_ORIGINS="*"

REM service chain options setting
set SC_BRIDGE=0 &:: if this is set, the following options will be used.
set SC_BRIDGE_PORT=50505
set SC_INDEXING=1

REM Setting 1 is to enable options, otherwise disabled.
set METRICS=1
set PROMETHEUS=1
set NO_DISCOVER=1
set DB_NO_PARALLEL_WRITE=0
set MULTICHANNEL=1

REM Raw options e.g) "--txpool.nolocals"
set ADDITIONAL=""

set KLAY_HOME=%homepath%\.ken

set DATA_DIR=%KLAY_HOME%\data
