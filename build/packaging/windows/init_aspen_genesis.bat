@echo off

set HOME=%~dp0
set CONF=%HOME%\conf

call %CONF%\klay-conf.cmd

IF NOT EXIST %KLAY_HOME% (
    mkdir %KLAY_HOME%
)

IF NOT EXIST %DATA_DIR% (
    mkdir %DATA_DIR%
)

echo "Init genesis for Klaytn aspen network"

copy %CONF%\aspen\static-nodes.json %DATA_DIR%\

%HOME%\bin\klay.exe init --datadir %DATA_DIR% %CONF%\aspen\genesis.json

@pause
