<script lang="ts">
    import { onMount } from "svelte";
    import { CasinoClient, ClientEvent } from "casino-sdk";
    import { goto } from "$app/navigation";
    import { WS_URL } from "$lib/config";
    import TopBar from "./TopBar.svelte";
    import LoadingOverlay from "./LoadingOverlay.svelte";
    import ListPage from "./ListPage.svelte";
    import GamePage from "./GamePage.svelte";

    let client: CasinoClient = new CasinoClient(WS_URL);

    let isInGameScreen = false;
    let selectedGame = "";

    const gameMap: { [key: string]: string } = {
        slotty: "https://casino-slotty.stmbl.dev",
        blackjack: "https://casino-blackjack.stmbl.dev",
        rolly: "https://casino-rolly.stmbl.dev/",
        catchfishy: "https://casino-catchfishy.stmbl.dev/",
    };

    onMount(async () => {
        // Listen for auth events
        client.on(ClientEvent.AUTH_FAIL, () => {
            console.log("## AUTH_FAIL ## 1");
            goto("/user/login");
        });

        await client.connect();

        if (!client.auth.isAuthenticated) {
            console.log("## AUTH_FAIL ## 2");
            goto("/user/login");
        }
    });
</script>

<svelte:head>
    <title>JHGambling</title>
</svelte:head>

<div class="app">
    <LoadingOverlay {client} />
    <TopBar
        {client}
        showExitButton={isInGameScreen}
        on:exit={() => (isInGameScreen = false)}
    />

    {#if isInGameScreen}
        <GamePage {client} gameURL={gameMap[selectedGame] || ""} />
    {:else}
        <ListPage
            on:select={(e) => {
                selectedGame = e.detail.id;
                isInGameScreen = true;
            }}
        />
    {/if}
</div>

<style lang="scss">
    .app {
        width: 100vw;
        height: 100vh;

        position: fixed;
        top: 0;
        left: 0;

        background-color: #18181b;

        display: flex;
        flex-direction: column;
        align-items: center;
    }
</style>
