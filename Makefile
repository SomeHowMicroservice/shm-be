auth:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protobufs/auth.proto

user:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    services/user/protobuf/user.proto

product:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    services/product/protobuf/product.proto

# post:
# 	@protoc --go_out=. --go_opt=paths=source_relative \
#     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#     services/post/protobuf/post.proto


auth-go-gen:
	@protoc --go_out=gateway/protobuf/auth --go-grpc_out=gateway/protobuf/auth \
	        --go_out=auth/protobuf --go-grpc_out=auth/protobuf \
	        protobuf/auth.proto

user-go-gen:
	@protoc --go_out=gateway/protobuf/user --go-grpc_out=gateway/protobuf/user \
	        --go_out=user/protobuf --go-grpc_out=user/protobuf \
	        protobuf/user.proto

product-go-gen:
	@protoc --go_out=gateway/protobuf/product --go-grpc_out=gateway/protobuf/product \
	        --go_out=product/protobuf --go-grpc_out=product/protobuf \
	        protobuf/product.proto

post-go-gen:
	@protoc --go_out=gateway/protobuf/post --go-grpc_out=gateway/protobuf/post \
	        protobuf/product.proto