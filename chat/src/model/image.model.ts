import { model, Schema } from "mongoose";
import type { IImage } from "../common/types";

const imageSchema = new Schema<IImage>({
  url: { type: String, required: true, trim: true },
  fileId: { type: String, required: false, trim: true },
});

export const ImageModel = model<IImage>("Image", imageSchema);
