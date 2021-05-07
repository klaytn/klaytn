# This code calls db migration.

# To use the db migration, a klaytn client should not be running.
# You need to have access to the DB.
# To checkout migration status, `tail -f logs-body.out`

# BIN file
KLAYTN_BIN=~/klaytn/bin/ken

# Source DB
SRC_DB_TYPE=LevelDB     # one of "LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3"
SRC_DB_DIR=~/klaytn/data
SRC_DB_SHARDS=4         # should be 2^x

# Destination DB
DST_DB_TYPE=LevelDB     # one of "LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3"
DST_DB_DIR=~/klaytn/db_migration/dst # for only localDB ("LevelDB", "BadgerDB", "MemoryDB"). neglected config for "DynamoDBS3"
DST_DB_SHARDS=8         # should be 2^x

## For only DynamoDBS3
DST_TABLENAME=db-migration      # for remoteDB ("DynamoDBS3"). neglected config for localDB
DST_RCU=100                     # neglected config for existing DB.
DST_WCU=4000                    # recommended to use auto-scaling up to 4000 while db migration

# set this value if you are using DynamoDB
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=

# call db migration
$KLAYTN_BIN db-migration start \
  --dbtype $SRC_DB_TYPE \
  --datadir $SRC_DB_DIR  \
  --db.num-statetrie-shards $SRC_DB_SHARDS \
  --dst.dbtype $DST_DB_TYPE \
  --dst.datadir $DST_DB_DIR \
  --db.dst.num-statetrie-shards $DST_DB_SHARDS \
   &> dbmigration-logs.out &
