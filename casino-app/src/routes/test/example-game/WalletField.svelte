<script lang="ts">
    import type { CasinoClient } from "casino-sdk";
    export let client: CasinoClient;

    let wallet = client.casino.wallet.store;

    let changeAmount = 0;

    function addAmount() {
        client.casino.wallet.addFunds(changeAmount);
    }

    function removeAmount() {
        client.casino.wallet.removeFunds(changeAmount);
    }
</script>

<div class="field">
    <div class="balance">
        <div class="label">Networth</div>
        <div class="value">{($wallet.NetworthCents / 100).toFixed(2)}$</div>
    </div>

    <div class="input-row">
        <input
            type="range"
            bind:value={changeAmount}
            min="0"
            max="25000"
            step="50"
        />
        <input type="number" bind:value={changeAmount} min="0" />
    </div>
    <div class="button-row">
        <button class="red-button" on:click={removeAmount}
            >Remove {(changeAmount / 100).toFixed(2)}$</button
        >
        <button class="green-button" on:click={addAmount}
            >Add {(changeAmount / 100).toFixed(2)}$</button
        >
    </div>
</div>

<style lang="scss">
    .field {
        width: calc(100% - 2rem);
        height: calc(100% - 2rem);

        padding: 1rem;
        grid-column: span 1;

        background-color: #ffffff;
        border-radius: 0.5rem;

        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: space-between;
        gap: 0.5rem;
    }

    .balance {
        width: 100%;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        gap: 1rem;

        font-family: var(--text-heading-family);
        font-size: var(--text-main-size);
        font-weight: var(--text-heading-weight);
        color: #000000;

        .label {
            width: 50%;
            text-align: right;
            color: #777777;
        }
        .value {
            width: 50%;
            text-align: left;
            text-wrap: wrap;
            overflow-wrap: break-word;
        }
    }

    .input-row {
        width: 100%;

        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;

        input {
            width: 100%;
        }
    }

    .button-row {
        width: 100%;

        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;

        button {
            width: calc(100% - 1.5rem);
            padding: 0.75rem;

            font-family: var(--text-main-family);
            font-size: var(--text-main-size);
            font-weight: var(--text-main-weight);

            cursor: pointer;
            transition: 0.1s;

            outline: none;
            border: none;
            border-radius: 0.5rem;

            &.red-button {
                background-color: #cc3553;
                color: #ffffff;

                &:active {
                    background-color: #d93858;
                }
            }

            &.green-button {
                background-color: #35cc90;

                &:active {
                    background-color: #38d999;
                }
            }
        }
    }
</style>
