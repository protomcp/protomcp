module protomcp.org/protomcp

go 1.23.0

require (
	darvaza.org/core v0.17.4 // indirect
	protomcp.org/protomcp/pkg/generator v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/protobuf v1.36.6
)

replace (
	protomcp.org/protomcp/pkg/generator => ./pkg/generator
	protomcp.org/protomcp/pkg/protomcp => ./pkg/protomcp
)
