<script lang="ts">
    import { onMount } from "svelte";
    import { CasinoClient } from "casino-sdk";
    import { goto } from "$app/navigation";
    import { WS_URL } from "$lib/config";

    export let step = 0;
    export let userExists = false;

    let username: string = "";
    let password: string = "";
    let displayName: string = "";
    let loading = false;

    // Validation messages
    let usernameError: string = "";
    let passwordError: string = "";
    let loginError: string = "";
    let registerError: string = "";

    // References to input elements for focus management
    let usernameInput: HTMLInputElement;
    let passwordInput: HTMLInputElement;
    let displayNameInput: HTMLInputElement;

    let client: CasinoClient;

    onMount(async () => {
        client = new CasinoClient(WS_URL);
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

        // Clear previous errors
        loginError = "";
        registerError = "";

        // Validate the current step
        if (step === 0 && !validateUsername()) {
            return;
        }

        if (step === 1 && !userExists && !validatePassword()) {
            return;
        }

        loading = true;

        try {
            if (step == 0) {
                step = 1;
                // Focus on password input after username is validated
                setTimeout(() => passwordInput?.focus(), 0);
            } else if (step == 1) {
                if (userExists) {
                    let res = await client.auth.login(username, password);
                    if (res.success) {
                        // Handle successful login
                        console.log("Login successful:", client.auth.user);
                        goto("/games");
                    } else {
                        loginError =
                            "Login fehlgeschlagen. Bitte überprüfen Sie Ihre Anmeldedaten.";
                    }
                } else {
                    step = 2;
                    // Focus on display name input after password is validated
                    setTimeout(() => displayNameInput?.focus(), 0);
                }
            } else if (step == 2) {
                if (!displayName || displayName.trim() === "") {
                    registerError = "Bitte geben Sie einen Anzeigenamen ein.";
                    return;
                }

                let res = await client.auth.register(
                    username,
                    password,
                    displayName,
                );
                if (res.success) {
                    // Handle successful registration
                    console.log("Register successful:", client.auth.user);
                    goto("/games");
                } else {
                    registerError = res.userAlreadyTaken
                        ? "Dieser Nutzername ist bereits vergeben."
                        : "Registrierung fehlgeschlagen. Bitte versuchen Sie es erneut.";
                }
            }
        } finally {
            loading = false;
        }
    }

    async function onUsernameInput() {
        validateUsername(false);

        if (username && username.length >= 3) {
            const exists = await client.auth.doesUserExist(username);
            userExists = exists;
        } else {
            userExists = false;
        }
    }

    function validateUsername(showError = true) {
        usernameError = "";

        if (!username || username.trim() === "") {
            if (showError)
                usernameError = "Bitte geben Sie einen Nutzernamen ein.";
            return false;
        }

        if (username.length < 3) {
            if (showError)
                usernameError =
                    "Der Nutzername muss mindestens 3 Zeichen lang sein.";
            return false;
        }

        if (username.length > 20) {
            if (showError)
                usernameError =
                    "Der Nutzername darf maximal 20 Zeichen lang sein.";
            return false;
        }

        if (!/^[a-zA-Z0-9_]+$/.test(username)) {
            if (showError)
                usernameError =
                    "Der Nutzername darf nur Buchstaben, Zahlen und Unterstriche enthalten.";
            return false;
        }

        return true;
    }

    function validatePassword(showError = true) {
        passwordError = "";

        if (!password || password.trim() === "") {
            if (showError) passwordError = "Bitte geben Sie ein Passwort ein.";
            return false;
        }

        if (password.length < 6) {
            if (showError)
                passwordError =
                    "Das Passwort muss mindestens 6 Zeichen lang sein.";
            return false;
        }

        return true;
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

    function onPasswordInput() {
        validatePassword(false);
    }

    function onDisplayNameInput() {
        registerError = "";
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
                    <div class="input-container">
                        <input
                            type="text"
                            placeholder="Nutzername"
                            bind:value={username}
                            bind:this={usernameInput}
                            on:input={onUsernameInput}
                            on:keydown={handleKeydown}
                            class:error={usernameError}
                        />
                        {#if usernameError}
                            <div class="error-message">{usernameError}</div>
                        {/if}
                    </div>
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
                        <div class="input-container">
                            <input
                                type="password"
                                placeholder="Passwort"
                                bind:value={password}
                                bind:this={passwordInput}
                                on:input={onPasswordInput}
                                on:keydown={handleKeydown}
                                class:error={passwordError}
                            />
                            {#if passwordError}
                                <div class="error-message">{passwordError}</div>
                            {/if}
                        </div>
                        {#if loginError}
                            <div class="error-message">{loginError}</div>
                        {/if}
                        <button on:click={nextStep} disabled={loading}>
                            {#if loading}
                                <span class="spinner"></span>
                            {:else}
                                Einloggen
                            {/if}
                        </button>
                    {:else}
                        <div class="input-container">
                            <input
                                type="password"
                                placeholder="Passwort"
                                bind:value={password}
                                bind:this={passwordInput}
                                on:input={onPasswordInput}
                                on:keydown={handleKeydown}
                                class:error={passwordError}
                            />
                            {#if passwordError}
                                <div class="error-message">{passwordError}</div>
                            {/if}
                        </div>
                        <button on:click={nextStep} disabled={loading}>
                            {#if loading}
                                <span class="spinner"></span>
                            {:else}
                                Weiter
                            {/if}
                        </button>
                    {/if}
                {:else if step == 2}
                    <div class="input-container">
                        <input
                            type="text"
                            placeholder="Anzeigename"
                            bind:value={displayName}
                            bind:this={displayNameInput}
                            on:input={onDisplayNameInput}
                            on:keydown={handleKeydown}
                            class:error={registerError}
                        />
                    </div>
                    {#if registerError}
                        <div class="error-message">{registerError}</div>
                    {/if}
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

            gap: 0.5rem;

            .input-container {
                width: 100%;
                position: relative;

                input {
                    width: calc(100% - 2rem);

                    padding: 1rem 1rem;
                    border-radius: 0.5rem;
                    background-color: #2c2b2a;
                    border: 2px solid transparent;
                    outline: none;

                    font-family: var(--text-main-family);
                    font-size: var(--text-main-size);
                    font-weight: var(--text-main-weight);
                    color: #ffffff;
                    transition: border-color 0.2s ease;

                    &.error {
                        border-color: rgb(255, 99, 107);
                    }
                }
            }

            .error-message {
                color: rgb(255, 99, 107);
                font-size: 0.9rem;
                margin-top: 0.2rem;
                width: 100%;
                text-align: left;

                font-family: var(--text-main-family);
                font-size: var(--text-main-size);
                font-weight: var(--text-main-weight);
            }

            button {
                width: 100%;
                padding: 0.5rem;
                border-radius: 0.5rem;
                background-color: rgb(204, 52, 59);
                border: none;
                cursor: pointer;

                margin-top: 0.5rem;

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
