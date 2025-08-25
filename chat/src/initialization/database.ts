import mongoose from "mongoose";

const initDb = async (uri: string) => {
  try {
    mongoose.set("strictQuery", true);
    await mongoose.connect(uri);

    process.on("SIGINT", async () => {
      await mongoose.connection.close();
      console.log("đã đóng kết nối MongoDB");
      process.exit(0);
    });
  } catch (err) {
    console.log("kết nối MongoDB thất bại: ", err);
    process.exit(1);
  }
};

export default initDb;
