import { WebsocketPacket } from "./packet";

// Connection status enum
export enum ConnectionStatus {
    DISCONNECTED = "DISCONNECTED",
    CONNECTING = "CONNECTING",
    CONNECTED = "CONNECTED",
    RECONNECTING = "RECONNECTING",
}

// Connection events
export enum ConnectionEvent {
    CONNECTED = "connected",
    DISCONNECTED = "disconnected",
    MESSAGE = "message",
    ERROR = "error",
    RECONNECTING = "reconnecting",
}

// Event callback types
export type ConnectionEventCallback = () => void;
export type MessageEventCallback = (packet: WebsocketPacket) => void;
export type ErrorEventCallback = (error: Error) => void;
export type ReconnectingEventCallback = (attemptNumber: number) => void;
