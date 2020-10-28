# This code calls db migration.

# BIN file
KLAYTN_BIN=~/klaytn/bin

# src DB
SRC_DB_TYPE=LevelDB
DATA_DIR=~/klaytn/data  # klatyn data dir
SRC_DB=misc             # leave empty if it srcDB is singleDB
SRC_DB_DIR=$DATA_DIR/klay/chaindata/$SRC_DB

# Dst DynamoDB
DST_DB_TYPE=DynamoDBS3
DST_DB_DIR=~/klaytn/db_migration/dst
DST_TABLENAME=db-migration
DST_RCU=100
DST_WCU=100

# set this value if you are using DynamoDB
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=

# call db migration
$KLAYTN_BIN/ken db-migration start \
  --db.single --db.dst.single \
  --dbtype $SRC_DB_TYPE --dst.dbtype $DST_DB_TYPE \
  --datadir $SRC_DB_DIR  \
  --dst.datadir $DST_DB_DIR --db.dst.dynamo.tablename $DST_TABLENAME \
  --db.dst.dynamo.is-provisioned --db.dst.dynamo.read-capacity $DST_RCU --db.dst.dynamo.write-capacity $DST_WCU \
   &> logs-$SRC_DB.out &
