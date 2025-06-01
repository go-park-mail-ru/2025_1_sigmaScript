# Instruction

## generate protobuf
находясь в папке сервиса auth_service выполнить команду:
```bash
protoc --go_out=./api/auth_api_v1/ --go-grpc_out=./api/auth_api_v1/ --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./proto/*.proto 
```
