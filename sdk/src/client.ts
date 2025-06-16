import { WebSocketClient } from "./websocket";

export class Client {
    public socket: WebSocketClient;

    constructor(public url: string) {
        this.socket = new WebSocketClient({
            url: this.url,
            autoReconnect: true,
            debug: true,
        });
    }
}
