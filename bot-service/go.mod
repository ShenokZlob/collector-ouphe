module github.com/ShenokZlob/collector-ouphe/bot-service

go 1.24.0

require (
	github.com/BlueMonday/go-scryfall v0.9.1
	github.com/ShenokZlob/collector-ouphe/pkg v0.0.0-20250517112458-b797b2576215
	github.com/go-telegram/bot v1.14.2
	github.com/go-telegram/fsm v0.2.0
	github.com/go-telegram/ui v0.5.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.8.0
	github.com/stretchr/testify v1.10.0
)

replace github.com/ShenokZlob/collector-ouphe/pkg => ../pkg

require (
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
