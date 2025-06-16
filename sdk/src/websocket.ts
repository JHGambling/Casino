import { ConnectionEvent, ConnectionStatus } from "./types/ws";
import { WebsocketPacket } from "./types/packet";

interface WebSocketClientOptions {
    url: string;
    autoReconnect?: boolean;
    reconnectInterval?: number;
    maxReconnectAttempts?: number;
    debug?: boolean;
}

export class WebSocketClient {
    private ws: WebSocket | null = null;
    private status: ConnectionStatus = ConnectionStatus.DISCONNECTED;
    private reconnectAttempts = 0;
    private eventListeners: Map<string, Function[]> = new Map();
    private nextNonce = 1;

    // Configuration
    private url: string;
    private autoReconnect: boolean;
    private reconnectInterval: number;
    private maxReconnectAttempts: number;
    private debug: boolean;

    constructor(options: WebSocketClientOptions) {
        this.url = options.url;
        this.autoReconnect = options.autoReconnect !== false;
        this.reconnectInterval = options.reconnectInterval || 1000;
        this.maxReconnectAttempts = options.maxReconnectAttempts || 50;
        this.debug = options.debug || false;
    }

    /**
     * Establish a WebSocket connection to the server
     */
    public connect(): void {
        if (
            this.status === ConnectionStatus.CONNECTED ||
            this.status === ConnectionStatus.CONNECTING
        ) {
            this.log("Already connected or connecting");
            return;
        }

        this.status = ConnectionStatus.CONNECTING;
        this.log(`Connecting to ${this.url}`);

        try {
            this.ws = new WebSocket(this.url);
            this.configureWebSocket();
        } catch (error) {
            this.handleError(
                error instanceof Error ? error : new Error(String(error)),
            );
        }
    }

    /**
     * Disconnect the WebSocket
     */
    public disconnect(): void {
        if (
            this.ws &&
            (this.status === ConnectionStatus.CONNECTED ||
                this.status === ConnectionStatus.CONNECTING)
        ) {
            this.log("Disconnecting");

            // Prevent auto reconnect when explicitly disconnected
            this.autoReconnect = false;
            this.ws.close();
        }
    }

    /**
     * Send a packet to the server
     */
    public send(type: string, payload: any = {}): number {
        if (this.status !== ConnectionStatus.CONNECTED) {
            throw new Error("Cannot send message: WebSocket is not connected");
        }

        const nonce = this.getNextNonce();
        const packet: WebsocketPacket = {
            type,
            payload,
            nonce,
        };

        this.ws!.send(JSON.stringify(packet));
        return nonce;
    }

    /**
     * Get the current connection status
     */
    public getStatus(): ConnectionStatus {
        return this.status;
    }

    /**
     * Register an event listener
     */
    public on(event: ConnectionEvent.CONNECTED, callback: () => void): void;
    public on(event: ConnectionEvent.DISCONNECTED, callback: () => void): void;
    public on(
        event: ConnectionEvent.ERROR,
        callback: (error: Error) => void,
    ): void;
    public on(
        event: ConnectionEvent.MESSAGE,
        callback: (packet: WebsocketPacket) => void,
    ): void;
    public on(
        event: ConnectionEvent.RECONNECTING,
        callback: (attempt: number) => void,
    ): void;
    public on(event: string, callback: Function): void {
        if (!this.eventListeners.has(event)) {
            this.eventListeners.set(event, []);
        }

        this.eventListeners.get(event)!.push(callback);
    }

    /**
     * Remove an event listener
     */
    public off(event: string, callback: Function): void {
        if (!this.eventListeners.has(event)) return;

        const listeners = this.eventListeners.get(event)!;
        const index = listeners.indexOf(callback);

        if (index !== -1) {
            listeners.splice(index, 1);
        }
    }

    /**
     * Configure WebSocket event handlers
     */
    private configureWebSocket(): void {
        if (!this.ws) return;

        this.ws.onopen = () => {
            this.status = ConnectionStatus.CONNECTED;
            this.reconnectAttempts = 0;
            this.log("Connection established");
            this.emit(ConnectionEvent.CONNECTED);
        };

        this.ws.onclose = () => {
            const wasConnected = this.status === ConnectionStatus.CONNECTED;
            this.status = ConnectionStatus.DISCONNECTED;

            if (wasConnected) {
                this.log("Connection closed");
                this.emit(ConnectionEvent.DISCONNECTED);
            }

            if (this.autoReconnect) {
                this.attemptReconnect();
            }
        };

        this.ws.onerror = (event) => {
            this.log("WebSocket error", event);
            this.handleError(new Error("WebSocket error"));
        };

        this.ws.onmessage = (event) => {
            this.handleMessage(event);
        };
    }

    /**
     * Process incoming WebSocket messages
     */
    private handleMessage(event: MessageEvent): void {
        try {
            const packet = JSON.parse(event.data) as WebsocketPacket;
            this.log("Received packet", packet);
            this.emit(ConnectionEvent.MESSAGE, packet);
        } catch (error) {
            this.log("Error parsing message", error);
            this.handleError(new Error("Failed to parse incoming message"));
        }
    }

    /**
     * Handle WebSocket errors
     */
    private handleError(error: Error): void {
        this.log("Error", error);
        this.emit(ConnectionEvent.ERROR, error);
    }

    /**
     * Attempt to reconnect to the server
     */
    private attemptReconnect(): void {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            this.log("Max reconnect attempts reached");
            return;
        }

        this.status = ConnectionStatus.RECONNECTING;
        this.reconnectAttempts++;

        this.log(
            `Reconnecting: attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`,
        );
        this.emit(ConnectionEvent.RECONNECTING, this.reconnectAttempts);

        setTimeout(() => {
            this.connect();
        }, this.reconnectInterval);
    }

    /**
     * Emit an event to registered listeners
     */
    private emit(event: string, ...args: any[]): void {
        if (!this.eventListeners.has(event)) return;

        for (const listener of this.eventListeners.get(event)!) {
            try {
                listener(...args);
            } catch (error) {
                console.error("Error in event listener", error);
            }
        }
    }

    /**
     * Get the next nonce value
     */
    private getNextNonce(): number {
        return this.nextNonce++;
    }

    /**
     * Internal logging function
     */
    private log(...args: any[]): void {
        if (this.debug) {
            console.log("[WebSocketClient]", ...args);
        }
    }
}
