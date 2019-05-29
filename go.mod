module github.com/deps-cloud/dis

go 1.12

require (
	github.com/deps-cloud/des v0.1.2
	github.com/deps-cloud/dts v0.0.0-20190527142220-7e296e5aee1b
	github.com/deps-cloud/rds v0.0.6
	github.com/spf13/cobra v0.0.4
	google.golang.org/grpc v1.21.0
	gopkg.in/src-d/go-git.v4 v4.11.0
)

replace github.com/deps-cloud/dts v0.0.0-20190527142220-7e296e5aee1b => ../dts
