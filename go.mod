module github.com/klaytn/klaytn

go 1.15

replace github.com/labstack/echo/v4 v4.2.0 => github.com/labstack/echo/v4 v4.9.0

require (
	github.com/Shopify/sarama v1.26.4
	github.com/VictoriaMetrics/fastcache v1.6.0
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf
	github.com/aristanetworks/goarista v0.0.0-20191001182449-186a6201b8ef
	github.com/aws/aws-sdk-go v1.34.28
	github.com/bt51/ntpclient v0.0.0-20140310165113-3045f71e2530
	github.com/cespare/cp v1.0.0
	github.com/clevergo/websocket v1.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.8.0
	github.com/dgraph-io/badger v1.6.0
	github.com/docker/docker v1.13.1
	github.com/edsrzf/mmap-go v1.0.0
	github.com/fatih/color v1.9.0
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-stack/stack v1.8.0
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/uint256 v1.2.0
	github.com/huin/goupnp v1.0.3-0.20220313090229-ca81a64b4204
	github.com/influxdata/influxdb v1.8.3
	github.com/jackpal/go-nat-pmp v1.0.2
	github.com/jinzhu/gorm v1.9.15
	github.com/julienschmidt/httprouter v1.2.0
	github.com/mattn/go-colorable v0.1.11
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/newrelic/go-agent/v3 v3.11.0
	github.com/otiai10/copy v1.0.1
	github.com/otiai10/mint v1.2.4 // indirect
	github.com/pbnjay/memory v0.0.0-20190104145345-974d429e7ae4
	github.com/pborman/uuid v0.0.0-20170612153648-e790cca94e6c
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/prometheus v2.1.0+incompatible
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563
	github.com/rjeczalik/notify v0.9.1
	github.com/robertkrimen/otto v0.0.0-20180506084358-03d472dc43ab
	github.com/rs/cors v1.7.0
	github.com/satori/go.uuid v1.2.0
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/urfave/cli v1.20.0
	github.com/valyala/fasthttp v1.34.0
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9
	golang.org/x/tools v0.1.0
	google.golang.org/grpc v1.32.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.40.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/fatih/set.v0 v0.1.0
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200619000410-60c24ae608a6
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0
	gotest.tools v2.2.0+incompatible
)
