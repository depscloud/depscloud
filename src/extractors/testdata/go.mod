module github.com/deps-cloud/finch

go 1.12

replace github.com/gogo/protobuf v1.2.1 => github.com/gogo/protobuf v1.2.1

replace github.com/deps-cloud/rds v0.0.2 => ../rds

replace (
	github.com/sirupsen/logrus v1.3.0 => github.com/sirupsen/logrus v1.3.0
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd => golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
	google.golang.org/grpc v1.18.0 => google.golang.org/grpc v1.18.0
)

require (
	github.com/gogo/protobuf v1.2.1
	github.com/deps-cloud/rds v0.0.2
	github.com/sirupsen/logrus v1.3.0
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
	google.golang.org/grpc v1.18.0
)
