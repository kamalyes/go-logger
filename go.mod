module github.com/kamalyes/go-logger

go 1.25.0

require (
	github.com/kamalyes/go-toolbox v0.18.1
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel/trace v1.44.0
	google.golang.org/grpc v1.77.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/rogpeppe/go-internal v1.15.0 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.37.0 // indirect
)

// 本地开发替换
// replace github.com/kamalyes/go-toolbox => ../go-toolbox
