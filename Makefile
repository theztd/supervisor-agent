build:
	env GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o ./build/supervisor-agent

build-arm:
	env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-extldflags '-static'" -o ./build/supervisor-agent
