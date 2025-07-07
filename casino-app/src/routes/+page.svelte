<script lang="ts">
    import { CasinoClient } from "casino-sdk";
    import { goto } from "$app/navigation";
    import { WS_URL } from "$lib/config";
    import { onMount } from "svelte";

    let client: CasinoClient = new CasinoClient(WS_URL);
    const isAuthenticated = client.casino.isAuthenticatedStore;
    client.connect();

    function play() {
        if (!$isAuthenticated) {
            goto("/user/login");
        } else {
            goto("/games");
        }
    }
</script>

<svelte:head>
    <title>JHGambling Home</title>
</svelte:head>

<div class="app">
    <div class="head">
        <div class="left head-section">
            <img src="/logo.png" alt="" class="logo" />
            <div class="title">JHGambling</div>
        </div>
        <div class="right head-section">
            <button on:click={play}>Spielen</button>
        </div>
    </div>

    <div class="banner"></div>
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

    .head {
        width: calc(100% - 2rem);
        height: 4rem;

        display: flex;
        justify-content: space-between;
        align-items: center;

        padding: 0 1rem;

        background-color: #222226;
    }

    .head-section {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 1rem;
    }

    .head-section.left {
        .logo {
            height: 2.5rem;
        }

        .title {
            font-family: var(--text-heading-family);
            font-size: var(--text-heading-size);
            font-weight: var(--text-heading-weight);
            color: #ffffff;
        }
    }

    .head-section.right {
        button {
            height: 2.5rem;

            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0rem 1rem;
            background-color: rgb(204, 52, 59);
            border: none;
            border-radius: 0.5rem;
            color: white;
            cursor: pointer;

            &:hover {
                background-color: rgb(223, 74, 82);
            }

            font-family: var(--text-main-family);
            font-size: var(--text-main-size);
            font-weight: var(--text-main-weight);
        }
    }

    .banner {
        width: 100%;
        aspect-ratio: 16 / 3;
        min-height: 10rem;

        background-color: #ffffff00;
        background-image: url("https://www.kayak.de/rimg/himg/16/96/c6/revato-177264-13657386-296801.jpg?width=1366&height=768&crop=true");
        background-position: center;
        background-size: cover;
    }
</style>
