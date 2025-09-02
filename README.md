### Đây là hệ thống phân tán trong cửa hàng thương mại điện tử😚

```
├── 📁 .git/ 🚫 (auto-hidden)
├── 📁 .vscode/ 🚫 (auto-hidden)
├── 📁 auth/
│   ├── 📁 common/
│   │   ├── 🐹 constants.go
│   │   ├── 🐹 errors.go
│   │   ├── 🐹 types.go
│   │   └── 🐹 utils.go
│   ├── 📁 config/
│   │   ├── 🐹 config.go
│   │   └── ⚙️ config.yaml
│   ├── 📁 consumers/
│   │   └── 🐹 send_email_consumer.go
│   ├── 📁 container/
│   │   └── 🐹 auth_container.go
│   ├── 📁 handler/
│   │   └── 🐹 grpc_handler.go
│   ├── 📁 initialization/
│   │   ├── 🐹 cache.go
│   │   ├── 🐹 grpc_client.go
│   │   └── 🐹 mq.go
│   ├── 📁 mq/
│   │   ├── 🐹 consumer.go
│   │   ├── 🐹 publisher.go
│   │   └── 🐹 retry.go
│   ├── 📁 proto/
│   │   ├── 📄 auth.proto
│   │   └── 📄 user.proto
│   ├── 📁 protobuf/
│   │   ├── 📁 auth/
│   │   │   ├── 🐹 auth.pb.go
│   │   │   └── 🐹 auth_grpc.pb.go
│   │   └── 📁 user/
│   │       ├── 🐹 user.pb.go
│   │       └── 🐹 user_grpc.pb.go
│   ├── 📁 repository/
│   │   ├── 🐹 cache_repository.go
│   │   └── 🐹 cache_repository_impl.go
│   ├── 📁 security/
│   │   └── 🐹 jwt.go
│   ├── 📁 server/
│   │   └── 🐹 grpc_server.go
│   ├── 📁 service/
│   │   ├── 🐹 auth_service.go
│   │   └── 🐹 auth_service_impl.go
│   ├── 📁 smtp/
│   │   ├── 📁 template/
│   │   │   └── 🌐 auth.html
│   │   ├── 🐹 smtp_service.go
│   │   └── 🐹 smtp_service_impl.go
│   ├── 📁 tmp/ 🚫 (auto-hidden)
│   ├── ⚙️ .air.toml
│   ├── 📄 .git 🚫 (auto-hidden)
│   ├── 🚫 .gitignore
│   ├── ⚙️ Makefile
│   ├── 🐹 go.mod
│   ├── 🐹 go.sum
│   └── 🐹 main.go
├── 📁 chat/
│   ├── 📁 node_modules/ 🚫 (auto-hidden)
│   ├── 📁 src/
│   │   ├── 📁 common/
│   │   │   └── 📄 types.ts
│   │   ├── 📁 config/
│   │   │   └── 📄 config.ts
│   │   ├── 📁 controller/
│   │   │   └── 📄 grpc.controller.ts
│   │   ├── 📁 initialization/
│   │   │   └── 📄 database.ts
│   │   ├── 📁 model/
│   │   │   ├── 📄 image.model.ts
│   │   │   └── 📄 message.model.ts
│   │   ├── 📁 proto/
│   │   │   └── 📄 chat.proto
│   │   ├── 📁 protobuf/
│   │   │   └── 📁 chat/
│   │   │       └── 📄 chat.ts
│   │   └── 📄 index.ts
│   ├── 🔒 .env 🚫 (auto-hidden)
│   ├── 📄 .git 🚫 (auto-hidden)
│   ├── 🚫 .gitignore
│   ├── ⚙️ Makefile
│   ├── 📖 README.md
│   ├── 🔒 bun.lock
│   ├── 📄 package.json
│   └── 📄 tsconfig.json
├── 📁 docs/
│   ├── 📄 bun_add.doc
│   ├── 📄 go_get.doc
│   └── 📄 go_install.doc
├── 📁 gateway/
│   ├── 📁 common/
│   │   ├── 🐹 api_response.go
│   │   ├── 🐹 constants.go
│   │   ├── 🐹 error_handler.go
│   │   ├── 🐹 errors.go
│   │   └── 🐹 types.go
│   ├── 📁 config/
│   │   ├── 🐹 app_config.go
│   │   ├── ⚙️ config.yaml
│   │   └── 🐹 cors_config.go
│   ├── 📁 container/
│   │   ├── 🐹 app_container.go
│   │   ├── 🐹 auth_container.go
│   │   ├── 🐹 chat_container.go
│   │   ├── 🐹 post_container.go
│   │   ├── 🐹 product_container.go
│   │   └── 🐹 user_container.go
│   ├── 📁 handler/
│   │   ├── 🐹 auth_handler.go
│   │   ├── 🐹 chat_handler.go
│   │   ├── 🐹 post_handler.go
│   │   ├── 🐹 product_handler.go
│   │   └── 🐹 user_handler.go
│   ├── 📁 initialization/
│   │   └── 🐹 grpc_client.go
│   ├── 📁 middleware/
│   │   ├── 🐹 jwt.go
│   │   └── 🐹 require.go
│   ├── 📁 proto/
│   │   ├── 📄 auth.proto
│   │   ├── 📄 chat.proto
│   │   ├── 📄 post.proto
│   │   ├── 📄 product.proto
│   │   └── 📄 user.proto
│   ├── 📁 protobuf/
│   │   ├── 📁 auth/
│   │   │   ├── 🐹 auth.pb.go
│   │   │   └── 🐹 auth_grpc.pb.go
│   │   ├── 📁 chat/
│   │   │   ├── 🐹 chat.pb.go
│   │   │   └── 🐹 chat_grpc.pb.go
│   │   ├── 📁 post/
│   │   │   ├── 🐹 post.pb.go
│   │   │   └── 🐹 post_grpc.pb.go
│   │   ├── 📁 product/
│   │   │   ├── 🐹 product.pb.go
│   │   │   └── 🐹 product_grpc.pb.go
│   │   └── 📁 user/
│   │       ├── 🐹 user.pb.go
│   │       └── 🐹 user_grpc.pb.go
│   ├── 📁 request/
│   │   ├── 🐹 auth_request.go
│   │   ├── 🐹 post_request.go
│   │   ├── 🐹 product_request.go
│   │   └── 🐹 user_request.go
│   ├── 📁 router/
│   │   ├── 🐹 auth_router.go
│   │   ├── 🐹 chat_router.go
│   │   ├── 🐹 post_router.go
│   │   ├── 🐹 product_router.go
│   │   └── 🐹 user_router.go
│   ├── 📁 server/
│   │   └── 🐹 http_server.go
│   ├── 📁 tmp/ 🚫 (auto-hidden)
│   ├── ⚙️ .air.toml
│   ├── 📄 .git 🚫 (auto-hidden)
│   ├── 🚫 .gitignore
│   ├── ⚙️ Makefile
│   ├── 🐹 go.mod
│   ├── 🐹 go.sum
│   └── 🐹 main.go
├── 📁 http/ 🚫 (auto-hidden)
├── 📁 post/
│   ├── 📁 .mvn/
│   │   └── 📁 wrapper/
│   │       └── 📄 maven-wrapper.properties
│   ├── 📁 src/
│   │   ├── 📁 main/
│   │   │   ├── 📁 java/
│   │   │   │   └── 📁 com/
│   │   │   │       └── 📁 service/
│   │   │   │           └── 📁 post/
│   │   │   │               ├── 📁 common/
│   │   │   │               │   └── ☕ SlugUtil.java
│   │   │   │               ├── 📁 config/
│   │   │   │               │   ├── ☕ ImageKitConfig.java
│   │   │   │               │   ├── ☕ JpaConfig.java
│   │   │   │               │   ├── ☕ RabbitMQConfig.java
│   │   │   │               │   └── ☕ RedisConfig.java
│   │   │   │               ├── 📁 controller/
│   │   │   │               │   └── ☕ GrpcController.java
│   │   │   │               ├── 📁 dto/
│   │   │   │               │   └── ☕ Base64UploadDto.java
│   │   │   │               ├── 📁 entity/
│   │   │   │               │   ├── ☕ BaseEntity.java
│   │   │   │               │   ├── ☕ ImageEntity.java
│   │   │   │               │   ├── ☕ PostEntity.java
│   │   │   │               │   └── ☕ TopicEntity.java
│   │   │   │               ├── 📁 exceptions/
│   │   │   │               │   ├── ☕ AlreadyExistsException.java
│   │   │   │               │   └── ☕ ResourceNotFoundException.java
│   │   │   │               ├── 📁 grpc_clients/
│   │   │   │               │   ├── ☕ BaseClient.java
│   │   │   │               │   ├── ☕ GrpcClientFactory.java
│   │   │   │               │   └── ☕ UserClient.java
│   │   │   │               ├── 📁 imagekit/
│   │   │   │               │   ├── ☕ ImageKitService.java
│   │   │   │               │   └── ☕ ImageKitServiceImpl.java
│   │   │   │               ├── 📁 mq/
│   │   │   │               │   ├── ☕ Consumer.java
│   │   │   │               │   └── ☕ Publisher.java
│   │   │   │               ├── 📁 redis/
│   │   │   │               │   ├── ☕ RedisService.java
│   │   │   │               │   └── ☕ RedisServiceImpl.java
│   │   │   │               ├── 📁 repository/
│   │   │   │               │   ├── ☕ ImageRepository.java
│   │   │   │               │   ├── ☕ PostRepository.java
│   │   │   │               │   └── ☕ TopicRepository.java
│   │   │   │               ├── 📁 service/
│   │   │   │               │   ├── ☕ PostService.java
│   │   │   │               │   └── ☕ PostServiceImpl.java
│   │   │   │               ├── 📁 specification/
│   │   │   │               │   └── ☕ PostSpecification.java
│   │   │   │               └── ☕ PostApplication.java
│   │   │   ├── 📁 proto/
│   │   │   │   ├── 📄 post.proto
│   │   │   │   └── 📄 user.proto
│   │   │   └── 📁 resources/
│   │   │       ├── 📄 application-sample.properties
│   │   │       └── 📄 application.properties
│   │   └── 📁 test/
│   │       └── 📁 java/
│   │           └── 📁 com/
│   │               └── 📁 service/
│   │                   └── 📁 post/
│   │                       └── ☕ PostApplicationTests.java
│   ├── 📁 target/ 🚫 (auto-hidden)
│   ├── 📄 .git 🚫 (auto-hidden)
│   ├── 📄 .gitattributes
│   ├── 🚫 .gitignore
│   ├── 📄 mvnw
│   ├── 🐚 mvnw.cmd
│   └── 📄 pom.xml
├── 📁 product/
│   ├── 📁 common/
│   │   ├── 🐹 constants.go
│   │   ├── 🐹 errors.go
│   │   ├── 🐹 types.go
│   │   └── 🐹 utils.go
│   ├── 📁 config/
│   │   ├── 🐹 config.go
│   │   └── ⚙️ config.yaml
│   ├── 📁 consumers/
│   │   ├── 🐹 delete_image_consumer.go
│   │   └── 🐹 upload_image_consumer.go
│   ├── 📁 container/
│   │   └── 🐹 product_container.go
│   ├── 📁 handler/
│   │   └── 🐹 grpc_handler.go
│   ├── 📁 imagekit/
│   │   ├── 🐹 imagekit_service.go
│   │   └── 🐹 imagekit_service_impl.go
│   ├── 📁 initialization/
│   │   ├── 🐹 database.go
│   │   ├── 🐹 grpc_client.go
│   │   └── 🐹 mq.go
│   ├── 📁 model/
│   │   ├── 🐹 category_model.go
│   │   ├── 🐹 color_model.go
│   │   ├── 🐹 image_model.go
│   │   ├── 🐹 inventory_model.go
│   │   ├── 🐹 product_model.go
│   │   ├── 🐹 size_model.go
│   │   ├── 🐹 tag_model.go
│   │   └── 🐹 variant_model.go
│   ├── 📁 mq/
│   │   ├── 🐹 consumer.go
│   │   ├── 🐹 publisher.go
│   │   └── 🐹 retry.go
│   ├── 📁 proto/
│   │   ├── 📄 product.proto
│   │   └── 📄 user.proto
│   ├── 📁 protobuf/
│   │   ├── 📁 product/
│   │   │   ├── 🐹 product.pb.go
│   │   │   └── 🐹 product_grpc.pb.go
│   │   └── 📁 user/
│   │       ├── 🐹 user.pb.go
│   │       └── 🐹 user_grpc.pb.go
│   ├── 📁 repository/
│   │   ├── 📁 category/
│   │   │   ├── 🐹 category_repository.go
│   │   │   └── 🐹 category_repository_impl.go
│   │   ├── 📁 color/
│   │   │   ├── 🐹 color_repository.go
│   │   │   └── 🐹 color_repository_impl.go
│   │   ├── 📁 image/
│   │   │   ├── 🐹 image_repository.go
│   │   │   └── 🐹 image_repository_impl.go
│   │   ├── 📁 inventory/
│   │   │   ├── 🐹 inventory_repository.go
│   │   │   └── 🐹 inventory_repository_impl.go
│   │   ├── 📁 product/
│   │   │   ├── 🐹 product_repository.go
│   │   │   └── 🐹 product_repository_impl.go
│   │   ├── 📁 size/
│   │   │   ├── 🐹 size_repository.go
│   │   │   └── 🐹 size_repository_impl.go
│   │   ├── 📁 tag/
│   │   │   ├── 🐹 tag_repository.go
│   │   │   └── 🐹 tag_repository_impl.go
│   │   └── 📁 variant/
│   │       ├── 🐹 variant_repository.go
│   │       └── 🐹 variant_repository_impl.go
│   ├── 📁 server/
│   │   └── 🐹 grpc_server.go
│   ├── 📁 service/
│   │   ├── 🐹 product_service.go
│   │   └── 🐹 product_service_impl.go
│   ├── 📁 tmp/ 🚫 (auto-hidden)
│   ├── 📄 .git 🚫 (auto-hidden)
│   ├── 🚫 .gitignore
│   ├── ⚙️ Makefile
│   ├── 🐹 go.mod
│   ├── 🐹 go.sum
│   └── 🐹 main.go
├── 📁 user/
│   ├── 📁 common/
│   │   ├── 🐹 errors.go
│   │   └── 🐹 utils.go
│   ├── 📁 config/
│   │   ├── 🐹 config.go
│   │   └── ⚙️ config.yaml
│   ├── 📁 container/
│   │   └── 🐹 user_container.go
│   ├── 📁 handler/
│   │   └── 🐹 grpc_handler.go
│   ├── 📁 initialization/
│   │   └── 🐹 database.go
│   ├── 📁 model/
│   │   ├── 🐹 address_model.go
│   │   ├── 🐹 measurement_model.go
│   │   ├── 🐹 profile_model.go
│   │   ├── 🐹 role_model.go
│   │   └── 🐹 user_model.go
│   ├── 📁 proto/
│   │   └── 📄 user.proto
│   ├── 📁 protobuf/
│   │   └── 📁 user/
│   │       ├── 🐹 user.pb.go
│   │       └── 🐹 user_grpc.pb.go
│   ├── 📁 repository/
│   │   ├── 📁 address/
│   │   │   ├── 🐹 address_repository.go
│   │   │   └── 🐹 address_repository_impl.go
│   │   ├── 📁 measurement/
│   │   │   ├── 🐹 measurement_repository.go
│   │   │   └── 🐹 measurement_repository_impl.go
│   │   ├── 📁 profile/
│   │   │   ├── 🐹 profile_repository.go
│   │   │   └── 🐹 profile_repository_impl.go
│   │   ├── 📁 role/
│   │   │   ├── 🐹 role_repository.go
│   │   │   └── 🐹 role_repository_impl.go
│   │   └── 📁 user/
│   │       ├── 🐹 user_repository.go
│   │       └── 🐹 user_repository_impl.go
│   ├── 📁 server/
│   │   └── 🐹 grpc_server.go
│   ├── 📁 service/
│   │   ├── 🐹 user_service.go
│   │   └── 🐹 user_service_impl.go
│   ├── 📁 tmp/ 🚫 (auto-hidden)
│   ├── ⚙️ .air.toml
│   ├── 📄 .git 🚫 (auto-hidden)
│   ├── 🚫 .gitignore
│   ├── ⚙️ Makefile
│   ├── 🐹 go.mod
│   ├── 🐹 go.sum
│   └── 🐹 main.go
├── 🚫 .gitignore
├── 📄 .gitmodules
├── 📖 README.md
├── 📄 go.work
└── 🐹 go.work.sum
```

---