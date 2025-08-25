import { Server, ServerCredentials } from "@grpc/grpc-js";
import config from "./config/config";
import initDb from "./initialization/database";
import { ChatServiceService } from "./protobuf/chat/chat";
import { grpcController } from "./controller/grpc.controller";

const main = async () => {
  await initDb(config.mongoUri);

  const server = new Server();
  server.addService(ChatServiceService, grpcController);
  server.bindAsync(
    `${config.serverHost}:${config.serverPort}`,
    ServerCredentials.createInsecure(),
    (err) => {
      if (err) {
        console.error("Kết nối tới phục vụ thất bại:", err);
        return;
      }
      console.log("Khởi chạy service thành công");
    }
  );
};

main();
