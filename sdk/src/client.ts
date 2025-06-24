import { Auth } from "./auth";
import { Database } from "./db";
import { ConnectionEvent, ConnectionStatus } from "./types/ws";
import { WebSocketClient } from "./websocket";

export class CasinoClient {
    public socket: WebSocketClient;

    public auth: Auth;
    public db: Database;

    private wasConnected: boolean = false;

    constructor(public url: string) {
        this.socket = new WebSocketClient({
            url: this.url,
            autoReconnect: true,
            debug: true,
        });

        this.socket.on(ConnectionEvent.CONNECTED, () => {
            if(!this.wasConnected) return;
            this.onConnect();
        });
        this.socket.on(ConnectionEvent.DISCONNECTED, () => {
            this.onDisconnect();
        });

        this.auth = new Auth(this);
        this.db = new Database(this);
    }

    public async connect() {
        this.socket.connect();
        await this.waitForConnect();
        await this.onConnect();
        this.wasConnected = true;
    }

    private async onConnect() {
        this.wasConnected = true;
        if (!this.auth.isAuthenticated) {
            console.log("Trying to authenticate from localstorage...");
            if (await this.auth.authFromLocalStorage()) {
                console.log("Authenticated from local storage!");
            }
        }
    }

    private async onDisconnect() {
        this.auth.revokeAuth();
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
}
