@echo off

set HOME=%~dp0
set CONF=%HOME%\conf

call %CONF%\kpn-conf.cmd

REM Check if exist data directory
set "NOT_INIT="
IF NOT EXIST %KLAY_HOME% (
    set NOT_INIT=1
)
IF NOT EXIST %DATA_DIR% (
    set NOT_INIT=1
)

IF DEFINED NOT_INIT (
    echo "[ERROR] : kpn is not initiated, Initiate kpn with genesis file first."
    GOTO end
)

set OPTIONS=--nodiscover

IF DEFINED NETWORK_ID (
    set OPTIONS=%OPTIONS% --networkid %NETWORK_ID%
)

IF DEFINED DATA_DIR (
    set OPTIONS=%OPTIONS% --datadir %DATA_DIR%
)

IF DEFINED PORT (
    set OPTIONS=%OPTIONS% --port %PORT%
)

IF DEFINED SERVER_TYPE (
    set OPTIONS=%OPTIONS% --srvtype %SERVER_TYPE%
)

IF DEFINED VERBOSITY (
    set OPTIONS=%OPTIONS% --verbosity %VERBOSITY%
)

IF DEFINED TXPOOL_EXEC_SLOTS_ALL (
    set OPTIONS=%OPTIONS% --txpool.exec-slots.all %TXPOOL_EXEC_SLOTS_ALL%
)

IF DEFINED TXPOOL_NONEXEC_SLOTS_ALL (
    set OPTIONS=%OPTIONS% --txpool.nonexec-slots.all %TXPOOL_NONEXEC_SLOTS_ALL%
)

IF DEFINED TXPOOL_EXEC_SLOTS_ACCOUNT (
    set OPTIONS=%OPTIONS% --txpool.exec-slots.account %TXPOOL_EXEC_SLOTS_ACCOUNT%
)

IF DEFINED TXPOOL_NONEXEC_SLOTS_ACCOUNT (
    set OPTIONS=%OPTIONS% --txpool.nonexec-slots.account %TXPOOL_NONEXEC_SLOTS_ACCOUNT%
)

IF DEFINED TXPOOL_LIFE_TIME (
    set OPTIONS=%OPTIONS% --txpool.lifetime %TXPOOL_LIFE_TIME%
)

IF DEFINED SYNCMODE (
    set OPTIONS=%OPTIONS% --syncmode %SYNCMODE%
)

IF DEFINED MAXCONNECTIONS (
    set OPTIONS=%OPTIONS% --maxconnections %MAXCONNECTIONS%
)

IF DEFINED LDBCACHESIZE (
    set OPTIONS=%OPTIONS% --db.leveldb.cache-size %LDBCACHESIZE%
)

IF DEFINED RPC_ENABLE (
    IF %RPC_ENABLE%==1 (
        set OPTIONS=%OPTIONS% --rpc --rpcapi %RPC_API% --rpcport %RPC_PORT% --rpcaddr %RPC_ADDR% --rpccorsdomain ^
%RPC_CORSDOMAIN% --rpcvhosts %RPC_VHOSTS%
    )
)

IF DEFINED WS_ENABLE (
    IF %WS_ENABLE%==1 (
        set OPTIONS=%OPTIONS% --ws --wsapi %WS_API% --wsaddr %WS_ADDR% --wsport %WS_PORT% --wsorigins %WS_ORIGINS%
    )
)

IF DEFINED METRICS (
    IF %METRICS%==1 (
        set OPTIONS=%OPTIONS% --metrics
    )
)

IF DEFINED PROMETHEUS (
    IF %PROMETHEUS%==1 (
        set OPTIONS=%OPTIONS% --prometheus
    )
)

IF DEFINED DB_NO_PARALLEL_WRITE (
    IF %DB_NO_PARALLEL_WRITE%==1 (
        set OPTIONS=%OPTIONS% --db.no-parallel-write
    )
)

IF DEFINED MULTICHANNEL (
    IF %MULTICHANNEL%==1 (
        set OPTIONS=%OPTIONS% --multichannel
    )
)

IF DEFINED ADDITIONAL (
    IF NOT %ADDITIONAL%=="" (
        set OPTIONS=%OPTIONS% %ADDITIONAL%
    )
)

%HOME%\bin\kpn.exe %OPTIONS%

:end
@pause
