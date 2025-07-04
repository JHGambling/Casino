<script lang="ts">
    import { CasinoClient } from "casino-sdk";
    import { goto } from "$app/navigation";
    import { slide } from "svelte/transition";
    import { createEventDispatcher } from "svelte";

    export let client: CasinoClient;
    let showDropdown = false;

    export let showExitButton = true;
    let dispatch = createEventDispatcher();

    let user = client.casino.user.store;
    let wallet = client.casino.wallet.store;

    function toggleDropdown() {
        showDropdown = !showDropdown;
    }

    function handleClickOutside(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".user-menu-container")) {
            showDropdown = false;
        }
    }

    function logout() {
        if (client) {
            client.auth.revokeAuth();
            localStorage.removeItem("casino-token");
            goto("/user/login");
        }
    }
</script>

<div class="topbar">
    <div class="right tb-section">
        {#if showExitButton}
            <button on:click={() => dispatch("exit")}>⬅️ Exit</button>
        {/if}
        <img src="/logo.png" alt="" class="logo" />
        <div class="title">JHGambling</div>
    </div>
    <div class="left tb-section">
        <div class="wallet-button">
            <span class="networth"
                >{($wallet.NetworthCents / 100).toFixed(2)}$</span
            >
        </div>
        <div class="user-menu-container">
            <button class="user-button" on:click={toggleDropdown}>
                <div class="user-avatar">
                    {#if $user.Username}
                        <img
                            src={`https://api.dicebear.com/9.x/glass/svg?seed=${$user.Username}`}
                            alt="Avatar"
                        />
                    {:else}
                        ?
                    {/if}
                </div>
                <span class="username">{$user.DisplayName || "User"}</span>
            </button>
            {#if showDropdown}
                <div class="dropdown-menu" transition:slide={{ duration: 150 }}>
                    <div class="dropdown-header">
                        <strong>{$user.DisplayName}</strong>
                        <small>@{$user.Username}</small>
                    </div>
                    <button class="dropdown-item" on:click={logout}>
                        <span class="icon">🚪</span> Abmelden
                    </button>
                </div>
            {/if}
        </div>
    </div>
</div>

<svelte:window on:click={handleClickOutside} />

<style lang="scss">
    .topbar {
        width: calc(100vw - 2rem);
        height: 4rem;

        padding: 0rem 1rem;

        background-color: #201f1e00;

        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
    }

    .tb-section {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 1rem;
    }

    .tb-section.right {
        .logo {
            height: 2.5rem;
        }

        .title {
            font-family: var(--text-heading-family);
            font-size: var(--text-heading-size);
            font-weight: var(--text-heading-weight);
            color: #ffffff;
        }

        button {
            height: 2.5rem;

            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0rem 1rem;
            background: #202126;
            border: none;
            border-radius: 0.5rem;
            color: white;
            cursor: pointer;

            &:hover {
                background: #292a30;
            }

            font-family: var(--text-main-family);
            font-size: var(--text-main-size);
            font-weight: var(--text-main-weight);
        }
    }

    .tb-section.left {
        .wallet-button {
            height: 2.5rem;

            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0rem 1rem;
            background: #202126;
            border: none;
            border-radius: 0.5rem;
            color: white;
            cursor: pointer;

            &:hover {
                background: #292a30;
            }

            font-family: var(--text-main-family);
            font-size: var(--text-main-size);
            font-weight: var(--text-main-weight);
        }

        .user-menu-container {
            position: relative;
        }

        .user-button {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.5rem 1rem;
            background: #202126;
            border: none;
            border-radius: 0.5rem;
            color: white;
            cursor: pointer;

            &:hover {
                background: #292a30;
            }

            .user-avatar {
                width: 1.5rem;
                height: 1.5rem;
                border-radius: 50%;
                background: rgb(204, 52, 59);
                display: flex;
                align-items: center;
                justify-content: center;
                font-weight: bold;
                overflow: hidden;

                img {
                    width: 100%;
                    height: 100%;
                    object-fit: cover;
                }
            }

            .username {
                font-family: var(--text-main-family);
                font-size: var(--text-main-size);
                font-weight: var(--text-main-weight);
            }
        }

        .dropdown-menu {
            position: absolute;
            top: 100%;
            right: 0;
            margin-top: 0.5rem;
            background: #202126;
            border-radius: 0.5rem;
            width: 100%;
            min-width: 12rem;
            box-shadow: 0 4px 0.8rem rgba(0, 0, 0, 0.25);
            z-index: 10;
            overflow: hidden;

            .dropdown-header {
                padding: 0.75rem 1rem;
                border-bottom: 1px solid #292a30;
                display: flex;
                flex-direction: column;

                strong {
                    font-family: var(--text-main-family);
                    font-size: var(--text-main-size);
                    color: white;
                }

                small {
                    font-family: var(--text-main-family);
                    font-size: 0.8rem;
                    color: #aaa;
                    margin-top: 0.2rem;
                }
            }

            .dropdown-item {
                display: flex;
                align-items: center;
                width: 100%;
                text-align: left;
                padding: 0.75rem 1rem;
                background: none;
                border: none;
                color: white;
                cursor: pointer;
                font-family: var(--text-main-family);
                font-size: var(--text-main-size);
                gap: 0.5rem;

                .icon {
                    font-size: 1.2rem;
                }

                &:hover {
                    background: #292a30;
                }

                &:not(:last-child) {
                    border-bottom: 1px solid #292a30;
                }
            }
        }
    }
</style>
