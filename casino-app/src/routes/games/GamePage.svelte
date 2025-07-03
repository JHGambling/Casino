<script lang="ts">
    import type { CasinoClient } from "casino-sdk";
    import { WS_URL } from "$lib/config";
    export let client: CasinoClient;

    let isAuthenticated = client.casino.isAuthenticatedStore;

    //export let gameURL: string = "https://jhgambling.github.io/slotty/";
    export let gameURL: string = "http://localhost:4173/index.html";
</script>

<div class="game-page">
    <div class="game">
        {#if $isAuthenticated}
            <iframe
                src="{gameURL}?token={client.auth.usedToken}&wsUrl={WS_URL}"
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

        background-color: #202126;

        iframe {
            width: 100%;
            height: 100%;
        }
    }
</style>
