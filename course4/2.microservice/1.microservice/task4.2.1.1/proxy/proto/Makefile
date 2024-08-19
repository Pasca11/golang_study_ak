compile_auth:
	protoc -I proto proto/auth.proto --go_out=./gen/auth/ --go_opt=paths=source_relative --go-grpc_out=./gen/auth --go-grpc_opt=paths=source_relative
compile_geo:
	protoc -I proto proto/geo.proto --go_out=./gen/geo/ --go_opt=paths=source_relative --go-grpc_out=./gen/geo/ --go-grpc_opt=paths=source_relative
