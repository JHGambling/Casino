<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { WS_URL } from "$lib/config";

    import { CasinoClient } from "casino-sdk";

    let client: CasinoClient = new CasinoClient(WS_URL, {
        authenticateFromLocalStorage: false,
        token: $page.data.token,
    });

    onMount(async () => {
        console.log("Token:", $page.data.token);
        await client.connect();
    });
</script>

<svelte:head>
    <title>Example Game</title>
</svelte:head>

<div class="game">
    Token: {$page.data.token}
</div>

<style lang="scss">
    .game {
        width: 100vw;
        height: 100vh;

        position: fixed;
        top: 0;
        left: 0;

        background-color: #202126;

        display: flex;
        flex-direction: column;
        align-items: center;
    }
</style>
