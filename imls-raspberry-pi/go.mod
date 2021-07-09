module gsa.gov/18f/imls-raspberry-pi

go 1.16

replace gsa.gov/18f/wifi-hardware-search-cli v0.0.0 => ./cmd/wifi-hardware-search-cli

replace gsa.gov/18f/input-initial-configuration v0.0.0 => ./cmd/input-configuration

replace gsa.gov/18f/session-counter v0.0.0 => ./cmd/session-counting

replace gsa.gov/18f/log-event v0.0.0 => ./cmd/log-event

replace gsa.gov/18f/version v0.0.0 => ./internal/version

replace gsa.gov/18f/config v0.0.0 => ./internal/config

replace gsa.gov/18f/http v0.0.0 => ./internal/http

replace gsa.gov/18f/cryptopasta v0.0.0 => ./internal/cryptopasta

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ./internal/wifi-hardware-search

replace gsa.gov/18f/analysis v0.0.0 => ./internal/analysis

replace gsa.gov/18f/logwrapper v0.0.0 => ./internal/logwrapper

require (
	github.com/newrelic/go-agent/v3 v3.13.0 // indirect
	github.com/newrelic/go-agent/v3/integrations/nrlogrus v1.0.1 // indirect
	gsa.gov/18f/input-initial-configuration v0.0.0 // indirect
	gsa.gov/18f/log-event v0.0.0 // indirect
	gsa.gov/18f/session-counter v0.0.0 // indirect
	gsa.gov/18f/wifi-hardware-search-cli v0.0.0 // indirect
)
