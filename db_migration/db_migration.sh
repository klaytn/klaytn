# This code calls db migration.

# To use the db migration, a klaytn client should not be running.
# You need to have access to the DB.
# To checkout migration status, `tail -f logs-body.out`

# BIN file
KLAYTN_BIN=~/klaytn/bin/ken

# src DB
SRC_DB_TYPE=LevelDB     # one of "LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3"
DATA_DIR=~/klaytn/data  # for localDB ("LevelDB", "BadgerDB", "MemoryDB"). neglected config for "DynamoDBS3"
SRC_DB=body             # result of `ls $DATA_DIR/klay/chaindata`
                        # * one of "body", "bridgeservice", "header", "misc", "receipts", "statetrie/0", "statetrie/1", "statetrie/2", "statetrie/3", "txlookup"
SRC_DB_DIR=$DATA_DIR/klay/chaindata/$SRC_DB

# Dst DynamoDB
DST_DB_TYPE=DynamoDBS3                # one of "LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3"
DST_DB_DIR=~/klaytn/db_migration/dst  # for localDB ("LevelDB", "BadgerDB", "MemoryDB"). neglected config for "DynamoDBS3"
DST_TABLENAME=db-migration            # for remoteDB ("DynamoDBS3"). neglected config for localDB
DST_RCU=100                             # neglected config for existing DB.
DST_WCU=4000                          # recommended to use auto-scaling up to 4000 while db migration

# set this value if you are using DynamoDB
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=

# call db migration
$KLAYTN_BIN db-migration start \
  --db.single --db.dst.single \
  --dbtype $SRC_DB_TYPE --dst.dbtype $DST_DB_TYPE \
  --datadir $SRC_DB_DIR  \
  --dst.datadir $DST_DB_DIR --db.dst.dynamo.tablename $DST_TABLENAME \
  --db.dst.dynamo.is-provisioned --db.dst.dynamo.read-capacity $DST_RCU --db.dst.dynamo.write-capacity $DST_WCU \
   &> logs-$(basename $SRC_DB).out &
