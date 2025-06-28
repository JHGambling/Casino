import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

export default defineConfig({
    plugins: [sveltekit()],
    // Ensure environment variables are correctly handled
    envPrefix: ["PUBLIC_", "VITE_"],
    // Define environment variable default values (fallbacks)
    define: {
        // You can add default values here if needed
        // 'import.meta.env.PUBLIC_WS_URL': JSON.stringify('wss://casino-host.stmbl.dev/ws')
    },
});
