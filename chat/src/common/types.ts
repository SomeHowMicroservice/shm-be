import { Document, type ObjectId } from "mongoose";

export interface IMessage extends Document {
  _id: ObjectId;
  conversationId: ObjectId;
  senderId: ObjectId;
  content?: string;
  image?: ObjectId;
  readBy: ObjectId[];
  createdAt: Date;
}

export interface IImage extends Document {
  _id: ObjectId;
  url: string;
  fileId: string;
}

export interface IConversation extends Document {
  _id: ObjectId;
  participants: ObjectId[];
  lastMessage?: ObjectId;
}