import { CasinoClient } from "../client";
import { WalletModel } from "../models/WalletModel";

/**
 * Wallet class provides simplified access to wallet information and operations
 */
export class Wallet {
    /**
     * Creates a new Wallet instance
     *
     * @param client The CasinoClient instance
     * @param walletData The wallet data model
     */
    constructor(
        private client: CasinoClient,
        private walletData: WalletModel,
    ) {}

    /**
     * Get the wallet ID
     */
    public get id(): number {
        return this.walletData.ID;
    }

    /**
     * Get the user ID associated with this wallet
     */
    public get userId(): number {
        return this.walletData.UserID;
    }

    /**
     * Check if the user has received the starting bonus
     */
    public get hasReceivedStartingBonus(): boolean {
        return this.walletData.ReceivedStartingBonus;
    }

    /**
     * Get the networth in cents
     */
    public get networthCents(): number {
        return this.walletData.NetworthCents;
    }

    /**
     * Get the networth in dollars (formatted with 2 decimal places)
     */
    public get networth(): number {
        return this.walletData.NetworthCents / 100;
    }

    /**
     * Get the raw wallet data model
     */
    public get data(): WalletModel {
        return this.walletData;
    }

    /**
     * Update the wallet's networth
     *
     * @param newAmountCents The new networth amount in cents
     * @returns A promise that resolves when the update is complete
     */
    public async updateNetworth(newAmountCents: number): Promise<void> {
        const result = await this.client.wallets.update(this.id, {
            NetworthCents: newAmountCents,
        });

        if (result.err) {
            throw new Error(`Failed to update networth: ${result.err}`);
        }

        // Update local data
        this.walletData.NetworthCents = newAmountCents;
    }

    /**
     * Add funds to the wallet
     *
     * @param amountCents The amount to add in cents
     * @returns A promise that resolves when the update is complete
     */
    public async addFunds(amountCents: number): Promise<void> {
        const newAmount = this.walletData.NetworthCents + amountCents;
        await this.updateNetworth(newAmount);
    }

    /**
     * Remove funds from the wallet
     *
     * @param amountCents The amount to remove in cents
     * @returns A promise that resolves when the update is complete
     * @throws Error if the wallet doesn't have enough funds
     */
    public async removeFunds(amountCents: number): Promise<void> {
        if (this.walletData.NetworthCents < amountCents) {
            throw new Error("Insufficient funds");
        }

        const newAmount = this.walletData.NetworthCents - amountCents;
        await this.updateNetworth(newAmount);
    }

    /**
     * Mark that the user has received the starting bonus
     *
     * @returns A promise that resolves when the update is complete
     */
    public async markStartingBonusReceived(): Promise<void> {
        if (this.walletData.ReceivedStartingBonus) {
            return; // Already received, nothing to do
        }

        const result = await this.client.wallets.update(this.id, {
            ReceivedStartingBonus: true,
        });

        if (result.err) {
            throw new Error(
                `Failed to mark starting bonus as received: ${result.err}`,
            );
        }

        // Update local data
        this.walletData.ReceivedStartingBonus = true;
    }

    /**
     * Refresh the wallet data from the server
     *
     * @returns A promise that resolves when the refresh is complete
     */
    public async refresh(): Promise<void> {
        const result = await this.client.wallets.findById(this.id);
        if (result.err) {
            throw new Error(`Failed to refresh wallet data: ${result.err}`);
        }

        this.walletData = result.result as WalletModel;
    }

    /**
     * Format the networth as a string with currency symbol
     *
     * @param currency The currency symbol to use (default: $)
     * @returns The formatted networth string
     */
    public formatNetworth(currency: string = "$"): string {
        return `${currency}${this.networth.toFixed(2)}`;
    }
}
