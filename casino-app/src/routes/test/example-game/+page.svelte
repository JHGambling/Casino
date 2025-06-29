<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { WS_URL } from "$lib/config";

    import { CasinoClient } from "casino-sdk";
    import ClientField from "./ClientField.svelte";
    import UserField from "./UserField.svelte";
    import WalletField from "./WalletField.svelte";

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
    <ClientField {client} token={$page.data.token} />
    <UserField {client} />
    <WalletField {client} />
</div>

<style lang="scss">
    .game {
        width: calc(100vw - 4rem);
        height: calc(100vh - 4rem);

        position: fixed;
        top: 0;
        left: 0;

        padding: 2rem;

        background-color: #ffffff;
        background-image: url("/assets/images/textured-background-1.png");

        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(20rem, 1fr));
        grid-auto-rows: 15rem;
        gap: 1rem;
    }
</style>
