<script lang="ts">
    import { onMount } from "svelte";
    import { CasinoClient, ClientEvent } from "casino-sdk";
    import { goto } from "$app/navigation";
    import { slide } from "svelte/transition";

    export let client: CasinoClient;
    let showDropdown = false;
    let username = "";
    let displayName = "";
    let balance = 0;

    $: if (client && client.auth && client.auth.user) {
        updateUserInfo();
    }

    onMount(() => {
        updateUserInfo();

        if (client) {
            // Listen for auth success to update user info
            client.on(ClientEvent.AUTH_SUCCESS, () => {
                updateUserInfo();
            });

            // Listen for auth revocation
            client.on(ClientEvent.AUTH_REVOKED, () => {
                username = "";
                displayName = "";
                balance = 0;
            });
        }
    });

    async function updateUserInfo() {
        if (client && client.auth && client.auth.user) {
            username = client.auth.user.Username;
            displayName = client.auth.user.DisplayName || username;
            // We could fetch the user's balance here if needed
        }
    }

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
        <img src="/logo.png" alt="" class="logo" />
        <div class="title">JHGambling</div>
    </div>
    <div class="left tb-section">
        <div class="user-menu-container">
            <button class="user-button" on:click={toggleDropdown}>
                <div class="user-avatar">
                    {#if username}
                        <img
                            src={`https://api.dicebear.com/9.x/glass/svg?seed=${username}`}
                            alt="Avatar"
                        />
                    {:else}
                        ?
                    {/if}
                </div>
                <span class="username">{displayName || "User"}</span>
            </button>
            {#if showDropdown}
                <div class="dropdown-menu" transition:slide={{ duration: 150 }}>
                    <div class="dropdown-header">
                        <strong>{displayName}</strong>
                        <small>@{username}</small>
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

        background-color: #201f1e;

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
            height: 3rem;
        }

        .title {
            font-family: var(--text-heading-family);
            font-size: var(--text-heading-size);
            font-weight: var(--text-heading-weight);
            color: #ffffff;
        }
    }

    .tb-section.left {
        .user-menu-container {
            position: relative;
        }

        .user-button {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.5rem 1rem;
            background: #2c2b2a;
            border: none;
            border-radius: 0.5rem;
            color: white;
            cursor: pointer;

            &:hover {
                background: #3c3b3a;
            }

            .user-avatar {
                width: 2rem;
                height: 2rem;
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

            .dropdown-arrow {
                font-size: 0.7rem;
                opacity: 0.7;
            }
        }

        .dropdown-menu {
            position: absolute;
            top: 100%;
            right: 0;
            margin-top: 0.5rem;
            background: #2c2b2a;
            border-radius: 0.5rem;
            width: 100%;
            min-width: 12rem;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
            z-index: 10;
            overflow: hidden;

            .dropdown-header {
                padding: 0.75rem 1rem;
                border-bottom: 1px solid #3c3b3a;
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
                    background: #3c3b3a;
                }

                &:not(:last-child) {
                    border-bottom: 1px solid #3c3b3a;
                }
            }
        }
    }
</style>
