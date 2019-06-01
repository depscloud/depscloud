module github.com/deps-cloud/gateway

go 1.12

require (
	github.com/deps-cloud/des v0.1.4
	github.com/deps-cloud/dts v0.0.3
	github.com/deps-cloud/rds v0.0.8
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/grpc-ecosystem/grpc-gateway v1.9.0
	github.com/spf13/cobra v0.0.4
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	google.golang.org/grpc v1.21.0
)

replace (
	github.com/deps-cloud/des v0.1.4 => ../des
	github.com/deps-cloud/dts v0.0.3 => ../dts
	github.com/deps-cloud/rds v0.0.8 => ../rds
)
