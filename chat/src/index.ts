import { Server, ServerCredentials } from "@grpc/grpc-js";
import { ChatServiceService } from "./protobuf/chat/chat";
import { grpcController } from "./controller/grpc.controller";
import config from "./config/config";

const server = new Server();
server.addService(ChatServiceService, grpcController)

server.bindAsync(`0.0.0.0:${config.port}`, ServerCredentials.createInsecure(), (err, port) => {
  if (err) {
    console.error("Server failed to bind:", err);
    return;
  }
  console.log(`🚀 gRPC ChatService running at 0.0.0.0:${port}`);
});