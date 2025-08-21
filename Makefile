auth-go-gen:
	@protoc --go_out=gateway/protobuf/auth --go-grpc_out=gateway/protobuf/auth \
	        --go_out=auth/protobuf/auth --go-grpc_out=auth/protobuf/auth \
	        protobuf/auth.proto

user-go-gen:
	@protoc --go_out=gateway/protobuf/user --go-grpc_out=gateway/protobuf/user \
	        --go_out=user/protobuf/user --go-grpc_out=user/protobuf/user \
					--go_out=auth/protobuf/user --go-grpc_out=auth/protobuf/user \
					--go_out=product/protobuf/user --go-grpc_out=product/protobuf/user \
	        protobuf/user.proto

product-go-gen:
	@protoc --go_out=gateway/protobuf/product --go-grpc_out=gateway/protobuf/product \
	        --go_out=product/protobuf/product --go-grpc_out=product/protobuf/product \
	        protobuf/product.proto

post-go-gen:
	@protoc --go_out=gateway/protobuf/post --go-grpc_out=gateway/protobuf/post \
	        protobuf/post.proto