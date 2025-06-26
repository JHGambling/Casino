<script lang="ts">
    import { onMount } from "svelte";
    import { CasinoClient, ClientEvent } from "casino-sdk";
    import { goto } from "$app/navigation";
    import TopBar from "./TopBar.svelte";
    import LoadingOverlay from "./LoadingOverlay.svelte";
    import ListPage from "./ListPage.svelte";
    import GamePage from "./GamePage.svelte";

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

    <ListPage />
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
