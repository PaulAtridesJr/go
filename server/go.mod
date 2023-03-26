module server

go 1.20

replace dummy => ../dummy

replace advanced => ../advanced

require (
	advanced v0.0.0-00010101000000-000000000000
	dummy v0.0.0-00010101000000-000000000000
)

require github.com/google/uuid v1.3.0 // indirect
