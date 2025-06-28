// Configuration file for environment variables
// This file centralizes access to all environment variables used in the application

// WebSocket server URL (default to production if not defined)
export const WS_URL = import.meta.env.PUBLIC_WS_URL || 'wss://casino-host.stmbl.dev/ws';

// You can add other environment variables as needed
// Example:
// export const API_URL = import.meta.env.PUBLIC_API_URL || 'https://api.example.com';
// export const DEBUG_MODE = import.meta.env.PUBLIC_DEBUG_MODE === 'true';

// Function to determine if running in development mode
export const isDevelopment = () => {
    return import.meta.env.DEV === true;
};

// Function to determine if running in production mode
export const isProduction = () => {
    return import.meta.env.PROD === true;
};

// Log configuration in development mode
if (isDevelopment()) {
    console.log('🔧 Configuration loaded:');
    console.log(`📡 WebSocket URL: ${WS_URL}`);
}
