<script lang="ts">
    import { onMount } from "svelte";
    import { CasinoClient } from "casino-sdk";
    import { goto } from "$app/navigation";
    import TopBar from "./TopBar.svelte";

    let client: CasinoClient;

    onMount(async () => {
        client = new CasinoClient("wss://casino-host.stmbl.dev/ws");
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
    <TopBar {client} />
</div>

<style lang="scss">
    .app {
        width: 100vw;
        height: 100vh;

        position: fixed;
        top: 0;
        left: 0;

        background-color: #141312;

        display: flex;
        flex-direction: column;
        align-items: center;
    }
</style>
