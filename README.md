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
│   ├── 📁 container/
│   │   └── 🐹 auth_container.go
│   ├── 📁 handler/
│   │   └── 🐹 grpc_handler.go
│   ├── 📁 initialization/
│   │   ├── 🐹 cache.go
│   │   ├── 🐹 grpc_client.go
│   │   └── 🐹 watermill.go
│   ├── 📁 mq/
│   │   ├── 🐹 consumer.go
│   │   └── 🐹 publisher.go
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
│   │   ├── 🐹 grpc_server.go
│   │   └── 🐹 server.go
│   ├── 📁 service/
│   │   ├── 🐹 auth_service.go
│   │   └── 🐹 auth_service_impl.go
│   ├── 📁 smtp/
│   │   ├── 📁 template/
│   │   │   └── 🌐 auth.html
│   │   ├── 🐹 smtp_service.go
│   │   └── 🐹 smtp_service_impl.go
│   ├── 📁 tmp/ 🚫 (auto-hidden)
│   ├── 📄 .git 🚫 (auto-hidden)
│   ├── 🚫 .gitignore
│   ├── ⚙️ Makefile
│   ├── 🐹 go.mod
│   ├── 🐹 go.sum
│   └── 🐹 main.go
├── 📁 chat/
│   ├── 📁 dist/ 🚫 (auto-hidden)
│   ├── 📁 node_modules/ 🚫 (auto-hidden)
│   ├── 📁 src/
│   │   ├── 📁 cloudinary/
│   │   │   ├── 📄 cloudinary-response.ts
│   │   │   ├── 📄 cloudinary.module.ts
│   │   │   ├── 📄 cloudinary.provider.ts
│   │   │   └── 📄 cloudinary.service.ts
│   │   ├── 📁 common/
│   │   │   ├── 📄 constants.ts
│   │   │   ├── 📄 error_handler.ts
│   │   │   ├── 📄 exceptions.ts
│   │   │   └── 📄 types.ts
│   │   ├── 📁 mq/
│   │   │   ├── 📄 consumer.service.ts
│   │   │   ├── 📄 mq.module.ts
│   │   │   ├── 📄 publisher.service.ts
│   │   │   └── 📄 retry.service.ts
│   │   ├── 📁 proto/
│   │   │   ├── 📄 chat.proto
│   │   │   └── 📄 user.proto
│   │   ├── 📁 protobuf/
│   │   │   ├── 📁 chat/
│   │   │   │   └── 📄 chat.ts
│   │   │   └── 📁 user/
│   │   │       └── 📄 user.ts
│   │   ├── 📁 repository/
│   │   │   ├── 📄 conversation.repository.ts
│   │   │   ├── 📄 image.repository.ts
│   │   │   └── 📄 message.repository.ts
│   │   ├── 📁 schema/
│   │   │   ├── 📄 conversation.schema.ts
│   │   │   ├── 📄 image.schema.ts
│   │   │   └── 📄 message.schema.ts
│   │   ├── 📄 app.controller.ts
│   │   ├── 📄 app.module.ts
│   │   ├── 📄 app.service.ts
│   │   └── 📄 main.ts
│   ├── 📁 test/
│   │   ├── 📄 app.e2e-spec.ts
│   │   └── 📄 jest-e2e.json
│   ├── 🔒 .env 🚫 (auto-hidden)
│   ├── 📄 .git 🚫 (auto-hidden)
│   ├── 🚫 .gitignore
│   ├── 📄 .prettierrc
│   ├── ⚙️ Makefile
│   ├── 📖 README.md
│   ├── 📄 eslint.config.mjs
│   ├── 📄 nest-cli.json
│   ├── 📄 package-lock.json
│   ├── 📄 package.json
│   ├── 📄 tsconfig.build.json 🚫 (auto-hidden)
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
│   ├── 📁 event/
│   │   ├── 🐹 client.go
│   │   └── 🐹 manager.go
│   ├── 📁 handler/
│   │   ├── 🐹 auth_handler.go
│   │   ├── 🐹 chat_handler.go
│   │   ├── 🐹 post_handler.go
│   │   ├── 🐹 product_handler.go
│   │   ├── 🐹 sse_handler.go
│   │   ├── 🐹 user_handler.go
│   │   └── 🐹 ws_handler.go
│   ├── 📁 initialization/
│   │   ├── 🐹 grpc_client.go
│   │   └── 🐹 watermill.go
│   ├── 📁 middleware/
│   │   ├── 🐹 jwt.go
│   │   └── 🐹 require.go
│   ├── 📁 mq/
│   │   └── 🐹 consumer.go
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
│   │   ├── 🐹 sse_router.go
│   │   ├── 🐹 user_router.go
│   │   └── 🐹 ws_router.go
│   ├── 📁 server/
│   │   ├── 🐹 http_server.go
│   │   └── 🐹 server.go
│   ├── 📁 socket/
│   │   ├── 🐹 client.go
│   │   └── 🐹 hub.go
│   ├── 📁 tmp/ 🚫 (auto-hidden)
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
│   │   │   │               │   ├── ☕ Base64UploadDto.java
│   │   │   │               │   └── ☕ ImageUploadedDto.java
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
│   │   └── 🐹 watermill.go
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
│   │   └── 🐹 publisher.go
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
│   │   ├── 🐹 grpc_server.go
│   │   └── 🐹 server.go
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
│   │   ├── 🐹 grpc_server.go
│   │   └── 🐹 server.go
│   ├── 📁 service/
│   │   ├── 🐹 user_service.go
│   │   └── 🐹 user_service_impl.go
│   ├── 📁 tmp/ 🚫 (auto-hidden)
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