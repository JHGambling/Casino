<script lang="ts">
    import type { CasinoClient } from "casino-sdk";
    export let client: CasinoClient;

    let isAuthenticated = client.casino.isAuthenticatedStore;
</script>

<div class="game-page">
    <div class="game">
        {#if $isAuthenticated}
            <iframe
                src="http://localhost:5173/test/example-game?token={client.auth
                    .usedToken}"
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
