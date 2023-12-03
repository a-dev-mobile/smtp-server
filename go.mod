module github.com/a-dev-mobile/smtp-server

go 1.20

require golang.org/x/exp v0.0.0-20231127185646-65229373498e

require gopkg.in/yaml.v3 v3.0.1 // indirect

require github.com/a-dev-mobile/common-lib v0.1.0

require (
	github.com/golang/protobuf v1.5.3 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231127180814-3a041ad873d4 // indirect
	google.golang.org/protobuf v1.31.0
)

require (
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

require (
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/grpc v1.59.0
)

// replace github.com/a-dev-mobile/common-lib => c:/DEV/MY/MY_GITHUB/common-lib/
