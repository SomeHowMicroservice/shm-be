import { model, Schema } from "mongoose";
import type { IMessage } from "../common/types";

const messageSchema = new Schema<IMessage>(
  {
    conversationId: {
      type: Schema.Types.ObjectId,
      ref: "Conversation",
      required: true,
    },
    senderId: {
      type: String,
      required: true,
      minLength: [36, "senderId phải đủ 36 ký tự"],
      maxLength: [36, "senderId phải đủ 36 ký tự"],
      trim: true,
    },
    content: {
      type: String,
      trim: true,
    },
    image: {
      type: Schema.Types.ObjectId,
      ref: "Image",
    },
    readBy: [
      {
        type: String,
        minLength: [36, "senderId phải đủ 36 ký tự"],
        maxLength: [36, "senderId phải đủ 36 ký tự"],
        trim: true,
      },
    ],
  },
  { timestamps: { createdAt: true, updatedAt: false } }
);

messageSchema.index({ conversationId: 1, createdAt: -1 });
messageSchema.index({ senderId: 1 });

export const MessageModel = model<IMessage>("Message", messageSchema);
