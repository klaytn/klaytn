# This code calls db migration.

# To use the db migration, a klaytn client should not be running.
# You need to have access to the DB.
# To checkout migration status, `tail -f logs-body.out`

# BIN file
KLAYTN_BIN=/Users/mini-admin/go/src/github.com/klaytn/klaytn/build/bin/ken

# src DB
SRC_DB_TYPE=LevelDB     # one of "LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3"
SRC_DB=body             # result of `ls $DATA_DIR/klay/chaindata`
                        # * one of "body", "bridgeservice", "header", "misc", "receipts", "statetrie/0", "statetrie/1", "statetrie/2", "statetrie/3", "txlookup"
SRC_DB_DIR=/Volumes/Samsung_T5/baobab_migrated/data/klay/chaindata

# Dst DynamoDB
DST_DB_TYPE=LevelDB                # one of "LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3"
DST_DB_DIR=/Volumes/Samsung_T5/baobab_migrated/data/klay/newchaindata  # for localDB ("LevelDB", "BadgerDB", "MemoryDB"). neglected config for "DynamoDBS3"
DST_TABLENAME=db-migration            # for remoteDB ("DynamoDBS3"). neglected config for localDB
DST_RCU=100                             # neglected config for existing DB.
DST_WCU=4000                          # recommended to use auto-scaling up to 4000 while db migration

# set this value if you are using DynamoDB
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=

# call db migration
$KLAYTN_BIN db-migration start \
  --dbtype $SRC_DB_TYPE \
  --datadir $SRC_DB_DIR  \
  --db.num-statetrie-shards 4 \
  --dst.dbtype $DST_DB_TYPE \
  --dst.datadir $DST_DB_DIR \
  --db.dst.num-statetrie-shards 8 \
#   &> logs-$(basename $SRC_DB).out &
