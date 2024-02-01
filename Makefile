build:
	env GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o ./build/supervisor-agent

build-arm:
	env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-extldflags '-static'" -o ./build/supervisor-agent

build-old:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" -o ./build/supervisor-agent-old
