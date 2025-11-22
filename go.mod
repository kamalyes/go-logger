module github.com/kamalyes/go-logger

go 1.23.0

require (
	github.com/stretchr/testify v1.11.1
	go.opentelemetry.io/otel/trace v1.38.0
)

require go.opentelemetry.io/otel v1.38.0 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kamalyes/go-toolbox v0.11.76
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kamalyes/go-toolbox => ../go-toolbox