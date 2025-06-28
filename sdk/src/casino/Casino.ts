import { CasinoClient } from "../client";
import { User } from "./User";
import { Wallet } from "./Wallet";
import { ClientEvent } from "../types/events";

/**
 * Casino class provides a simplified interface to casino operations
 * including user management and wallet operations
 */
export class Casino {
    /**
     * The current user instance, if authenticated
     */
    private _currentUser: User | null = null;

    /**
     * Creates a new Casino instance
     *
     * @param client The CasinoClient instance
     */
    constructor(private client: CasinoClient) {
        // Listen for authentication events to update the current user
        this.client.on(ClientEvent.AUTH_SUCCESS, async () => {
            await this.refreshCurrentUser();
        });

        this.client.on(ClientEvent.AUTH_REVOKED, () => {
            this._currentUser = null;
        });
    }

    /**
     * Get the current authenticated user
     *
     * @returns The current user or null if not authenticated
     */
    public get currentUser(): User | null {
        return this._currentUser;
    }

    /**
     * Refreshes the current user information from the server
     *
     * @returns A promise that resolves with the current user or null if not authenticated
     */
    public async refreshCurrentUser(): Promise<User | null> {
        const userData = await this.client.users.getCurrentUser();
        if (!userData) return null;

        this._currentUser = new User(this.client, userData);
        return this._currentUser;
    }
}
