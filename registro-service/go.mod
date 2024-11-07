module github.com/marcelogbrito/nats-centromedico/registro-service

go 1.23.0

replace github.com/marcelogbrito/nats-centromedico/shared v0.0.0-20241106202713-c17bdb35368e => ../shared

replace github.com/marcelogbrito/nats-centromedico/shared => ../shared

require (
	github.com/gorilla/mux v1.8.1
	github.com/marcelogbrito/nats-centromedico/shared v0.0.0-20241106202713-c17bdb35368e
	github.com/nats-io/nats.go v1.37.0
	github.com/nats-io/nuid v1.0.1
)

require filippo.io/edwards25519 v1.1.0 // indirect

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
