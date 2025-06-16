import { CasinoClient } from "./client";
import {
    AuthAuthenticatePacket,
    AuthAuthenticateResponsePacket,
    AuthLoginPacket,
    AuthLoginResponsePacket,
    AuthRegisterPacket,
    AuthRegisterResponsePacket,
} from "./types/packets";

export class Auth {
    public isAuthenticated: boolean = false;
    public authenticatedAs: number = 0;
    public authenticationExpiresAt: Date = new Date(0);

    constructor(private client: CasinoClient) {}

    public async authFromLocalStorage(): Promise<boolean> {
        let token = localStorage.getItem("casino-token");
        if (token) {
            return this.authenticate(token);
        } else {
            return false;
        }
    }

    public async register(
        username: string,
        password: string,
        displayName: string,
    ): Promise<{ success: boolean; userAlreadyTaken: boolean }> {
        const response = (
            await this.client.socket.request("auth/register", {
                username,
                password,
                displayName,
            } as AuthRegisterPacket)
        ).payload as AuthRegisterResponsePacket;

        if (!response.success) {
            return {
                success: false,
                userAlreadyTaken: response.userAlreadyExists,
            };
        }

        const authSuccess = await this.authenticate(response.token || "");
        return {
            success: authSuccess,
            userAlreadyTaken: false,
        };
    }

    public async login(
        username: string,
        password: string,
    ): Promise<{ success: boolean; userNotFound: boolean }> {
        const response = (
            await this.client.socket.request("auth/login", {
                username,
                password,
            } as AuthLoginPacket)
        ).payload as AuthLoginResponsePacket;

        if (!response.success) {
            return {
                success: false,
                userNotFound: response.userDoesNotExist,
            };
        }

        const authSuccess = await this.authenticate(response.token || "");
        return {
            success: authSuccess,
            userNotFound: false,
        };
    }

    public async authenticate(token: string): Promise<boolean> {
        const response = (
            await this.client.socket.request("auth/authenticate", {
                token,
            } as AuthAuthenticatePacket)
        ).payload as AuthAuthenticateResponsePacket;

        if (response.success) {
            this.isAuthenticated = true;
            this.authenticatedAs = response.userID;
            this.authenticationExpiresAt = new Date(response.expiresAt);
            localStorage.setItem("casino-token", token);
            return true;
        } else {
            return false;
        }
    }
}
