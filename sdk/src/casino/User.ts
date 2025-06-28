import { CasinoClient } from "../client";
import { UserModel } from "../models/UserModel";
import { Wallet } from "./Wallet";

/**
 * User class provides simplified access to user information and operations
 */
export class User {
    /**
     * The wallet associated with this user
     */
    private _wallet: Wallet | null = null;

    /**
     * Creates a new User instance
     *
     * @param client The CasinoClient instance
     * @param userData The user data model
     */
    constructor(
        private client: CasinoClient,
        private userData: UserModel,
    ) {}

    /**
     * Get the user's ID
     */
    public get id(): number {
        return this.userData.ID;
    }

    /**
     * Get the user's username
     */
    public get username(): string {
        return this.userData.Username;
    }

    /**
     * Get the user's display name
     */
    public get displayName(): string {
        return this.userData.DisplayName;
    }

    /**
     * Get the date when the user joined
     */
    public get joinedAt(): Date {
        return new Date(this.userData.JoinedAt);
    }

    /**
     * Check if the user is an admin
     */
    public get isAdmin(): boolean {
        return this.userData.IsAdmin;
    }

    /**
     * Get the raw user data model
     */
    public get data(): UserModel {
        return this.userData;
    }

    /**
     * Get the user's wallet
     *
     * @returns A promise that resolves with the user's wallet
     */
    public async getWallet(): Promise<Wallet> {
        if (this._wallet) {
            return this._wallet;
        }

        // If wallet data is already in the user model
        if (this.userData.Wallet) {
            this._wallet = new Wallet(this.client, this.userData.Wallet);
            return this._wallet;
        }

        // Try to fetch the wallet
        const walletData = await this.client.wallets.getCurrentWallet();
        if (walletData) {
            this._wallet = new Wallet(this.client, walletData);
            return this._wallet;
        }

        throw new Error("Could not retrieve wallet for this user");
    }

    /**
     * Update the user's display name
     *
     * @param newDisplayName The new display name
     * @returns A promise that resolves when the update is complete
     */
    public async updateDisplayName(newDisplayName: string): Promise<void> {
        const result = await this.client.users.update(this.id, {
            DisplayName: newDisplayName,
        });

        if (result.err) {
            throw new Error(`Failed to update display name: ${result.err}`);
        }

        // Update local data
        this.userData.DisplayName = newDisplayName;
    }

    /**
     * Refresh the user data from the server
     *
     * @returns A promise that resolves when the refresh is complete
     */
    public async refresh(): Promise<void> {
        const result = await this.client.users.findById(this.id);
        if (result.err) {
            throw new Error(`Failed to refresh user data: ${result.err}`);
        }

        this.userData = result.result as UserModel;
        this._wallet = null; // Clear cached wallet to force refresh
    }
}
