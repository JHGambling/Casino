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

        this.auth = new Auth(this);
        this.db = new Database(this);
    }

    public async connect() {
        this.socket.connect();
        await this.waitForConnect();

        if (await this.auth.authFromLocalStorage()) {
            console.log("Authenticated from local storage!");
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
}
