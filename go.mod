module github.com/yuanzhangcai/srsd

go 1.14

require (
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro v1.18.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.15.0 // indirect
	google.golang.org/grpc v1.26.0
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.1.0
