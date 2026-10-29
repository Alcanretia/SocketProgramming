import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// `npm run dev` serves the UI on :5173 and forwards the WebSocket to the Go
// server on :8080. Start the Go server with `-dev` so it accepts that origin.
export default defineConfig({
  plugins: [react()],
  server: {
    host: true,
    proxy: { "/ws": { target: "ws://localhost:8080", ws: true } },
  },
});
