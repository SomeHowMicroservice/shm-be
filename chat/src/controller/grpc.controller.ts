import { SendMessageRequest, SendMessageResponse } from "../protobuf/chat/chat";
import type { ServerUnaryCall, sendUnaryData } from "@grpc/grpc-js";

export const grpcController = {
  sendMessage(
    call: ServerUnaryCall<SendMessageRequest, SendMessageResponse>,
    callback: sendUnaryData<SendMessageResponse>
  ) {
    console.log("Received: ", call.request.message);
    callback(null, { ok: true });
  },
};
