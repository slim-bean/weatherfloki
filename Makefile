


arm7:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags '-extldflags "-static"' -o weatherfloki weatherfloki.go

arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags '-extldflags "-static"' -o weatherfloki weatherfloki.go