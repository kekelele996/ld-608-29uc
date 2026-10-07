import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 20108,
    host: "0.0.0.0",
    // 本地开发代理；生产环境由 frontend/nginx.conf 的 location /api/ 反代。
    proxy: {
      "/api": {
        target: "http://localhost:21108",
        changeOrigin: true
      }
    }
  }
});
