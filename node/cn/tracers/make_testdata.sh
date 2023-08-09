#!/bin/sh

if [ $# -ne 2 ]; then
    echo "Usage: $0 <rpc_url> <txid>"
    echo ""
    echo "Outputs an testdata JSON to stdout"
    echo "Example:"
    echo "    $0 http://127.0.0.1:8551 0x85c29014bc3e11442e06507dc443110becfb63e69fe030d0d104f66a9f705db2 > testdata/call_tracer_hello.json"
    exit 1
fi

SCRIPT_DIR=$(dirname "$0")
RPC_URL="$1"
TXID=$2

PATH="$PATH:$SCRIPT_DIR/../../../build/bin" # For convenience
which ken 1>/dev/null 2>/dev/null
if [ $? -ne 0 ]; then
    echo "The 'ken' program is not found in \$PATH"
fi

ken attach \
    --preload $SCRIPT_DIR/make_testdata.js \
    --exec "makeTest('${TXID}')" \
    ${RPC_URL} \
    | sed '$ d'   # delete last line

