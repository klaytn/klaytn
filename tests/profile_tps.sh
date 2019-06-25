#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
BENCHNAME="BenchmarkValueTransfer"

BENCHTIME=${BENCHTIME:-"5s"}
BENCHCOUNT=${BENCHCOUNT:-1}
TXS_PER_BLOCK=${TXS_PER_BLOCK:-1000}

TMP=`mktemp`

cd $DIR/../tests

CMD="go test -cpuprofile cpu.out -run X -bench $BENCHNAME -benchtime $BENCHTIME"
echo "executing $CMD for $BENCHCOUNT times"
echo "" > $TMP
for i in `seq 1 $BENCHCOUNT`; do
    $CMD | tee -a $TMP
done
NS=`grep "ns/op" $TMP | awk 'BEGIN{total=0.0;count=0} {total+=$3;count++} END{printf("%f", total/count)}'`
TPS=$(echo "1.0 / $NS * 1000.0 * 1000.0 * 1000" | bc -l)
echo "TPS for a single machine = $TPS"
rm -rf $TMP

go tool pprof -http=localhost:6061 ./tests.test cpu.out
