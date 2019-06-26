# How to generate `klaytn.pb.go` from `klaytn.proto`

## 1. Install protobuf for Go
```
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

## 2. Generate a Go file from protobuf IDL
```
$ protoc -I=. --go_out=plugins=grpc:. klaytn.proto
```

## 3. Change the generated file

Because of version mismatch issue, we need to change
`proto.ProtoPackageIsVersion3` in the generated file `klaytn.pb.go` to
`proto.ProtoPackageIsVersion2`.

```
$ sed -i -e 's/ProtoPackageIsVersion3/ProtoPackageIsVersion2/g' klaytn.pb.go
```
