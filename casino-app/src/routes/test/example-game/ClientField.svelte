<script lang="ts">
    import type { CasinoClient } from "casino-sdk";

    export let client: CasinoClient;
    export let token: string;

    let isConnected = client.casino.isConnectedStore;
    let isAuthenticated = client.casino.isAuthenticatedStore;
</script>

<div class="field">
    <div class="row">
        <div class="value">Token: {token}</div>
    </div>
    <div class="row">
        <div class="value">Connected: {$isConnected ? "✅" : "❌"}</div>
        <div class="value">Authenticated: {$isAuthenticated ? "✅" : "❌"}</div>
    </div>
</div>

<style lang="scss">
    .field {
        width: calc(100% - 2rem);
        height: calc(100% - 2rem);

        padding: 1rem;

        grid-column: span 2;

        background-color: #ffffff;
        border-radius: 0.5rem;

        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: space-between;
        gap: 0.5rem;
    }

    .row {
        width: 100%;
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 0.5rem;

        font-family: var(--text-main-family);
        font-size: var(--text-main-size);
        font-weight: var(--text-main-weight);
        color: #000000;

        .value {
            width: 100%;
            text-wrap: wrap;
            overflow-wrap: break-word;
        }
    }
</style>
