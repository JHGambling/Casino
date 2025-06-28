import { Auth } from "./auth";
import { Database } from "./db";
import { UserTable } from "./db/UserTable";
import { WalletTable } from "./db/WalletTable";
import { ClientEvent } from "./types/events";
import { ConnectionEvent, ConnectionStatus } from "./types/ws";
import { WebSocketClient } from "./websocket";

export class CasinoClient {
    public socket: WebSocketClient;

    public auth: Auth;
    public db: Database;
    public users: UserTable;
    public wallets: WalletTable;

    private wasConnected: boolean = false;
    private eventListeners: Map<string, Function[]> = new Map();

    constructor(public url: string) {
        this.socket = new WebSocketClient({
            url: this.url,
            autoReconnect: true,
            debug: true,
        });

        this.socket.on(ConnectionEvent.CONNECTED, () => {
            if (!this.wasConnected) return;
            this.onConnect();
            this.emit(ClientEvent.CONNECT);
        });
        this.socket.on(ConnectionEvent.DISCONNECTED, () => {
            this.onDisconnect();
            this.emit(ClientEvent.DISCONNECT);
        });

        this.auth = new Auth(this);
        this.db = new Database(this);
        this.users = new UserTable(this);
        this.wallets = new WalletTable(this);
    }

    public async connect() {
        this.socket.connect();
        await this.waitForConnect();
        await this.onConnect();
        this.wasConnected = true;
        this.emit(ClientEvent.CONNECT);
    }

    private async onConnect() {
        this.wasConnected = true;
        if (!this.auth.isAuthenticated) {
            console.log("Trying to authenticate from localstorage...");
            if (await this.auth.authFromLocalStorage()) {
                console.log("Authenticated from local storage!");
                this.emit(ClientEvent.AUTH_SUCCESS, this.auth.authenticatedAs);
            }
        }
    }

    private async onDisconnect() {
        const wasAuthenticated = this.auth.isAuthenticated;
        this.auth.revokeAuth();
        if (wasAuthenticated) {
            this.emit(ClientEvent.AUTH_REVOKED);
        }
    }

    private async waitForConnect() {
        await new Promise<void>((resolve) => {
            if (this.socket.getStatus() === ConnectionStatus.CONNECTED) {
                resolve();
            } else {
                this.socket.on(ConnectionEvent.CONNECTED, resolve);
            }
        });
    }

    /**
     * Register an event listener
     */
    public on(event: ClientEvent, callback: Function): void {
        if (!this.eventListeners.has(event)) {
            this.eventListeners.set(event, []);
        }

        this.eventListeners.get(event)!.push(callback);
    }

    /**
     * Remove an event listener
     */
    public off(event: ClientEvent, callback: Function): void {
        if (!this.eventListeners.has(event)) return;

        const listeners = this.eventListeners.get(event)!;
        const index = listeners.indexOf(callback);

        if (index !== -1) {
            listeners.splice(index, 1);
        }
    }

    /**
     * Emit an event to registered listeners
     */
    private emit(event: ClientEvent, ...args: any[]): void {
        if (!this.eventListeners.has(event)) return;

        for (const listener of this.eventListeners.get(event)!) {
            try {
                listener(...args);
            } catch (error) {
                console.error(`Error in event listener for ${event}:`, error);
            }
        }
    }
}
