<script lang="ts">
    export let step = 0;
    export let userExists = false;

    let username: string;
    let password: string;

    function nextStep() {
        if(step == 0) {
            if(isValidUsername(username)) {
                step = 1;
            } else {
                alert("Bitte geben Sie einen gültigen Nutzernamen ein.");
            }
        } else if(step == 1) {
            if(userExists) {

            } else {
                if(password && password.length >= 6) {
                    step = 2;
                } else {
                    alert("Bitte geben Sie ein gültiges Passwort ein (mindestens 6 Zeichen).");
                }
            }
        }
    }

    function isValidUsername(username: string) {
        // Check if the username is valid (e.g., not empty, no special characters)
        return username && /^[a-zA-Z0-9_]+$/.test(username) && username.length >= 3 && username.length <= 20;
    }
</script>

<svelte:head>
    <title>JHGambling - Login</title>
</svelte:head>

<div class="app">
    <div class="center-content">
        <div class="head">
            JHGambling
        </div>
        <div class="form">
            <div class="form-title">
                {#if step == 0}
                Anmelden / Registrieren
                {:else}
                {#if userExists}
                Anmelden
                {:else}
                Registrieren
                {/if}
                {/if}
            </div>
            <div class="form-body">
                {#if step == 0}
                    <input type="text" placeholder="Nutzername" bind:value={username}>
                    <button on:click={nextStep}>Weiter</button>
                {:else if step == 1}
                    {#if userExists}
                        <input type="password" placeholder="Passwort" bind:value={password}>
                        <button on:click={nextStep}>Einloggen</button>
                    {:else}
                        <input type="password" placeholder="Passwort" bind:value={password}>
                        <button on:click={nextStep}>Weiter</button>
                    {/if}
                {:else if step == 2}
                    <input type="text" placeholder="Anzeigename">
                    <!-- Hier noch irgendwas für Profilbild -->
                    <button on:click={nextStep}>Registrieren</button>
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

                &:hover {
                    background-color: rgb(223, 74, 82);
                }
            }
        }
    }
</style>
