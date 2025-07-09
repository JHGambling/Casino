<script lang="ts">
    import { ClientEvent, type CasinoClient } from "casino-sdk";
    import { WS_URL } from "$lib/config";
    export let client: CasinoClient;

    let isAuthenticated = client.casino.isAuthenticatedStore;
    let hasFinishedLoading = false;

    //export let gameURL: string = "https://jhgambling.github.io/slotty/";
    export let gameURL: string = "http://localhost:5174/index.html";

    client.on(ClientEvent.GAME_FINISHED_LOADING, () => {
        hasFinishedLoading = true;
    });
</script>

<div class="game-page">
    <div class="game">
        {#if $isAuthenticated}
            {#if !hasFinishedLoading}
                <div class="loading-overlay">
                    <div class="spinner" />
                </div>
            {/if}
            <iframe
                src="{gameURL}?token={client.auth
                    .usedToken}&wsUrl={WS_URL}&session={client.session}&usesdk=1"
                title="Game"
                frameborder="0"
            ></iframe>
        {/if}
    </div>
</div>

<style lang="scss">
    .game-page {
        width: calc(100vw - 2rem);
        height: calc(100vh - 4rem - 1rem);

        padding: 1rem;
        padding-top: 0rem;
        overflow: hidden;

        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .game {
        width: 100%;
        flex-grow: 1;
        border-radius: 1rem;
        overflow: hidden;
        position: relative;

        background-color: #202126;

        iframe {
            width: 100%;
            height: 100%;
        }
    }

    .loading-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(32, 33, 38, 1);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 10;
    }

    .spinner {
        border: 5px solid rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        border-top-color: #fff;
        width: 50px;
        height: 50px;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
