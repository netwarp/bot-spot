# Bot Spot v3

```bash

# Windows
go build -o bin/bot-spot.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/bot-spot

# Mac
GOOS=darwin GOARCH=amd64 go build -o bin/bot-spot
```

#### migrate from old db to new

```go
go test ./tools -run TestFromCloverToSqlite
```
