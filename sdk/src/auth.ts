import { CasinoClient } from "./client";
import { UserModel } from "./models/UserModel";
import {
    AuthAuthenticatePacket,
    AuthAuthenticateResponsePacket,
    AuthLoginPacket,
    AuthLoginResponsePacket,
    AuthRegisterPacket,
    AuthRegisterResponsePacket,
    DoesUserExistPacket,
    DoesUserExistResponsePacket,
} from "./types/packets";

export class Auth {
    public isAuthenticated: boolean = false;
    public authenticatedAs: number = 0;
    public authenticationExpiresAt: Date = new Date(0);

    public user: UserModel | null = null;

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

            this.user = await this.fetchUser();

            return true;
        } else {
            this.revokeAuth();

            return false;
        }
    }

    public async fetchUser(): Promise<UserModel | null> {
        if (!this.isAuthenticated) {
            return null;
        }

        const res = await this.client.db.performOperation(
            "users",
            "findByID",
            this.authenticatedAs,
            null,
        );

        return res.result as UserModel;
    }

    public revokeAuth() {
        this.user = null;
        this.isAuthenticated = false;
        this.authenticatedAs = -1;
        this.authenticationExpiresAt = new Date(0);
    }

    public async doesUserExist(username: string): Promise<boolean> {
        const response = (
            await this.client.socket.request("auth/does_user_exist", {
                username
            } as DoesUserExistPacket)
        ).payload as DoesUserExistResponsePacket;

        return response.userExists && response.success;
    }
}
