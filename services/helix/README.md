# helix
realtime tcp socket server for managing objects in the world.

Not tied to terrain, can scale out, and can interact with each other.

## Usage
generate protobuf interfaces (mac)
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
brew install protobuf
./tools/generate_api.sh
```
go run
```bash
go run main.go
```

## TODO
- udp 기반으로 설계하고 일부 주요 패킷에 ack 처리? (reliable udp)
