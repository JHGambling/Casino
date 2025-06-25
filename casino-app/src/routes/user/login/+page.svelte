<script lang="ts">
    import { onMount } from "svelte";
    import { CasinoClient } from "casino-sdk";
    import { goto } from "$app/navigation";

    export let step = 0;
    export let userExists = false;

    let username: string;
    let password: string;
    let displayName: string;
    let loading = false;

    // References to input elements for focus management
    let usernameInput: HTMLInputElement;
    let passwordInput: HTMLInputElement;
    let displayNameInput: HTMLInputElement;

    let client: CasinoClient;

    onMount(async () => {
        client = new CasinoClient("wss://casino-host.stmbl.dev/ws");
        await client.connect();

        if (client.auth.isAuthenticated) {
            // User is already authenticated, redirect to games
            goto("/games");
        } else {
            // Focus on the username input when the component mounts
            setTimeout(() => usernameInput?.focus(), 0);
        }
    });

    async function nextStep() {
        if (loading) return;
        loading = true;

        try {
            if (step == 0) {
                if (isValidUsername(username)) {
                    step = 1;
                    // Focus on password input after username is validated
                    setTimeout(() => passwordInput?.focus(), 0);
                } else {
                    alert("Bitte geben Sie einen gültigen Nutzernamen ein.");
                }
            } else if (step == 1) {
                if (userExists) {
                    let res = await client.auth.login(username, password);
                    if (res.success) {
                        // Handle successful login
                        console.log("Login successful:", client.auth.user);
                        goto("/games");
                    } else {
                        alert("Login fehlgeschlagen");
                    }
                } else {
                    if (password && password.length >= 6) {
                        step = 2;
                        // Focus on display name input after password is validated
                        setTimeout(() => displayNameInput?.focus(), 0);
                    } else {
                        alert(
                            "Bitte geben Sie ein gültiges Passwort ein (mindestens 6 Zeichen).",
                        );
                    }
                }
            } else if (step == 2) {
                let res = await client.auth.register(
                    username,
                    password,
                    displayName,
                );
                if (res.success) {
                    // Handle successful login
                    console.log("Register successful:", client.auth.user);
                    goto("/games");
                } else {
                    alert("Login fehlgeschlagen");
                }
            }
        } finally {
            loading = false;
        }
    }

    async function onUsernameInput() {
        if (username && username.length >= 3) {
            const exists = await client.auth.doesUserExist(username);
            userExists = exists;
        } else {
            userExists = false;
        }
    }

    function isValidUsername(username: string) {
        // Check if the username is valid (e.g., not empty, no special characters)
        return (
            username &&
            /^[a-zA-Z0-9_]+$/.test(username) &&
            username.length >= 3 &&
            username.length <= 20
        );
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Enter") {
            nextStep();
        }
    }
</script>

<svelte:head>
    <title>JHGambling - Login</title>
</svelte:head>

<div class="app">
    <div class="center-content">
        <div class="head">
            <img src="/logo.png" alt="" class="logo" />
        </div>
        <div class="form">
            <div class="form-title">
                {#if step == 0}
                    Anmelden / Registrieren
                {:else if userExists}
                    Anmelden
                {:else}
                    Registrieren
                {/if}
            </div>
            <div class="form-body">
                {#if step == 0}
                    <input
                        type="text"
                        placeholder="Nutzername"
                        bind:value={username}
                        bind:this={usernameInput}
                        on:input={onUsernameInput}
                        on:keydown={handleKeydown}
                    />
                    <button
                        on:click={nextStep}
                        disabled={!isValidUsername(username) || loading}
                    >
                        {#if loading}
                            <span class="spinner"></span>
                        {:else}
                            {userExists ? "Anmelden" : "Weiter"}
                        {/if}
                    </button>
                {:else if step == 1}
                    {#if userExists}
                        <input
                            type="password"
                            placeholder="Passwort"
                            bind:value={password}
                            bind:this={passwordInput}
                            on:keydown={handleKeydown}
                        />
                        <button on:click={nextStep} disabled={loading}>
                            {#if loading}
                                <span class="spinner"></span>
                            {:else}
                                Einloggen
                            {/if}
                        </button>
                    {:else}
                        <input
                            type="password"
                            placeholder="Passwort"
                            bind:value={password}
                            bind:this={passwordInput}
                            on:keydown={handleKeydown}
                        />
                        <button on:click={nextStep} disabled={loading}>
                            {#if loading}
                                <span class="spinner"></span>
                            {:else}
                                Weiter
                            {/if}
                        </button>
                    {/if}
                {:else if step == 2}
                    <input
                        type="text"
                        placeholder="Anzeigename"
                        bind:value={displayName}
                        bind:this={displayNameInput}
                        on:keydown={handleKeydown}
                    />
                    <!-- Hier noch irgendwas für Profilbild -->
                    <button on:click={nextStep} disabled={loading}>
                        {#if loading}
                            <span class="spinner"></span>
                        {:else}
                            Registrieren
                        {/if}
                    </button>
                {/if}
            </div>
        </div>
    </div>
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
        justify-content: center;
    }

    .center-content {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 0.25rem;
    }

    .head {
        font-family: var(--text-heading-family);
        font-size: var(--text-heading-size);
        font-weight: var(--text-heading-weight);
        color: #ffffff;

        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        gap: 1rem;

        padding: 1rem 10rem;
        background-color: #201f1e;
        border-radius: 1rem 1rem 0.25rem 0.25rem;

        img.logo {
            height: 8rem;
        }
    }

    .form {
        width: 100%;
        padding: 1rem 0rem;

        border-radius: 0.25rem 0.25rem 1rem 1rem;
        background-color: #201f1e;

        display: flex;
        flex-direction: column;
        align-items: center;

        .form-title {
            font-family: var(--text-main-family);
            font-size: var(--text-main-size);
            font-weight: var(--text-main-weight);
            color: #ffffff;

            margin-bottom: 0.5rem;
        }

        .form-body {
            width: calc(100% - 4rem);
            padding: 1rem 2rem;

            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;

            gap: 0.25rem;

            input {
                width: calc(100% - 2rem);

                padding: 1rem 1rem;
                border-radius: 0.5rem;
                background-color: #2c2b2a;
                border: none;
                outline: none;

                font-family: var(--text-main-family);
                font-size: var(--text-main-size);
                font-weight: var(--text-main-weight);
                color: #ffffff;
            }

            button {
                width: 100%;
                padding: 0.5rem;
                border-radius: 0.5rem;
                background-color: rgb(204, 52, 59);
                border: none;
                cursor: pointer;

                margin-top: 0.25rem;

                font-family: var(--text-main-family);
                font-size: var(--text-main-size);
                font-weight: var(--text-main-weight);
                color: #ffffff;

                display: flex;
                justify-content: center;
                align-items: center;

                &:hover {
                    background-color: rgb(223, 74, 82);
                }

                &:disabled {
                    cursor: not-allowed;
                    background-color: #2c2b2a;
                }

                .spinner {
                    display: inline-block;
                    width: 1rem;
                    height: 1rem;
                    border: 2px solid rgba(255, 255, 255, 0.3);
                    border-radius: 50%;
                    border-top-color: white;
                    animation: spin 1s ease-in-out infinite;
                }
            }
        }
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
