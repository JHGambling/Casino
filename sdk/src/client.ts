import { Auth } from "./auth";
import { Database } from "./db";
import { ConnectionEvent, ConnectionStatus } from "./types/ws";
import { WebSocketClient } from "./websocket";

export class CasinoClient {
    public socket: WebSocketClient;

    public auth: Auth;
    public db: Database;

    constructor(public url: string) {
        this.socket = new WebSocketClient({
            url: this.url,
            autoReconnect: true,
            debug: true,
        });

        this.socket.on(ConnectionEvent.CONNECTED, () => {
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
    }

    private async onConnect() {
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
}
