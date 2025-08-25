const config = {
  serverPort: process.env.SERVER_PORT || 8085,
  serverHost: process.env.SERVER_HOST || "localhost",
  mongoUri: process.env.MONGO_URI || "mongodb://localhost:27017/shm_chat",
};

export default config;
