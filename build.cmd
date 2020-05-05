SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o mac_user ./cmd/main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o linux_user ./cmd/main.go