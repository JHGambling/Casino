<script lang="ts">
    import { onMount } from "svelte";
    import { CasinoClient, ClientEvent } from "casino-sdk";
    import { goto } from "$app/navigation";
    import TopBar from "./TopBar.svelte";
    import LoadingOverlay from "./LoadingOverlay.svelte";
    import GameCard from "./GameCard.svelte";

    let client: CasinoClient = new CasinoClient(
        "wss://casino-host.stmbl.dev/ws",
    );
    //let client = new CasinoClient("ws://localhost:9000/ws");

    onMount(async () => {
        // Listen for auth events
        client.on(ClientEvent.AUTH_FAIL, () => {
            goto("/user/login");
        });

        await client.connect();

        if (!client.auth.isAuthenticated) {
            goto("/user/login");
        }
    });
</script>

<svelte:head>
    <title>JHGambling</title>
</svelte:head>

<div class="app">
    <LoadingOverlay {client} />
    <TopBar {client} />

    <div class="games-container">
        <GameCard name="Slotty" imageUrl="/assets/images/slot_example_image.webp" />
        <GameCard name="Blackjack" imageUrl="/assets/images/blackjack_example_image.webp" />
        <GameCard name="Roulette" imageUrl="/assets/images/roulette_example_image.png" />
        <GameCard name="Russisch Roulette" imageUrl="/assets/images/russian_roulette_example_image.webp" />
    </div>
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

    .games-container {
        width: calc(100vw - 2rem);
        height: calc(100vh - 4rem - 2rem);

        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(25rem, 1fr));
        grid-auto-rows: 20rem;
        gap: 1rem;
        padding: 1rem;
        overflow-y: auto;
    }
</style>
