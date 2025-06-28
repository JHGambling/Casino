import { CasinoClient } from "../client";
import { WalletModel } from "../models/WalletModel";
import { BaseTable } from "./BaseTable";

/**
 * WalletTable provides a specialized interface for interacting with wallet records
 */
export class WalletTable extends BaseTable<WalletModel> {
    /**
     * Creates a new WalletTable instance
     *
     * @param client The CasinoClient instance
     */
    constructor(client: CasinoClient) {
        super(client, "wallets");
    }

    /**
     * Get the current authenticated user's wallet
     *
     * @returns The current user's wallet or null if not authenticated
     */
    public async getCurrentWallet(): Promise<WalletModel | null> {
        if (!this.client.auth.isAuthenticated || !this.client.auth.user) {
            return null;
        }

        // If the user already has the wallet loaded
        if (this.client.auth.user.Wallet) {
            return this.client.auth.user.Wallet;
        }

        /*try {
            // Try to get the wallet ID from the user or fall back to user ID
            const walletId =
                this.client.auth.user.Wallet?.ID ||
                this.client.auth.authenticatedAs;

            // Fetch the wallet by ID
            const result = await this.findById(walletId);
            if (result.err) {
                console.error("Error fetching current wallet:", result.err);
                return null;
            }

            return result.result as WalletModel;
        } catch (error) {
            console.error("Exception fetching current wallet:", error);
            return null;
            }*/
        return null;
    }
}
