<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { CasinoClient, ClientEvent } from "casino-sdk";
    import { fade } from "svelte/transition";

    export let client: CasinoClient;
    let isConnected = false;
    let fadeTransitionDuration = 300;

    function handleConnect() {
        isConnected = true;
    }

    function handleDisconnect() {
        isConnected = false;
    }

    onMount(() => {
        // Check initial connection status
        isConnected = client.socket.getStatus() === "CONNECTED";

        // Set up event listeners
        client.on(ClientEvent.CONNECT, handleConnect);
        client.on(ClientEvent.DISCONNECT, handleDisconnect);
    });

    onDestroy(() => {
        // Clean up event listeners
        client.off(ClientEvent.CONNECT, handleConnect);
        client.off(ClientEvent.DISCONNECT, handleDisconnect);
    });
</script>

{#if !isConnected}
    <div
        class="loading-overlay"
        transition:fade={{ duration: fadeTransitionDuration }}
    >
        <div class="spinner-container">
            <div class="spinner"></div>
            <div class="loading-text">Connecting...</div>
        </div>
    </div>
{/if}

<style lang="scss">
    .loading-overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100vw;
        height: 100vh;
        background-color: rgba(31, 33, 34, 0.3);
        backdrop-filter: blur(1rem);
        z-index: 1000;

        display: flex;
        justify-content: center;
        align-items: center;
    }

    .spinner-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1.5rem;
    }

    .spinner {
        width: 4rem;
        height: 4rem;
        border: 6px solid rgba(255, 255, 255, 0.1);
        border-radius: 50%;
        border-top: 6px solid #cc343b;
        animation: spin 0.75s infinite linear;
    }

    .loading-text {
        color: white;
        font-family: var(--text-main-family);
        font-size: var(--text-main-size);
        font-weight: var(--text-main-weight);
        letter-spacing: 0.05em;
    }

    @keyframes spin {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }
</style>
