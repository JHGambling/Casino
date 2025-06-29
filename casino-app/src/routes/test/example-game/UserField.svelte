<script lang="ts">
    import type { CasinoClient } from "casino-sdk";
    export let client: CasinoClient;

    let user = client.casino.user.store;
</script>

<div class="field">
    <div class="row">
        <div class="label">ID</div>
        <div class="value">{$user.ID}</div>
    </div>
    <div class="row">
        <div class="label">Username</div>
        <div class="value">{$user.Username}</div>
    </div>
    <div class="row">
        <div class="label">DisplayName</div>
        <div class="value">{$user.DisplayName}</div>
    </div>
    <div class="row">
        <div class="label">IsAdmin</div>
        <div class="value">{$user.IsAdmin}</div>
    </div>
    <div class="row">
        <div class="label">JoinedAt</div>
        <div class="value">{$user.JoinedAt}</div>
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

    .row {
        width: 100%;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        gap: 1rem;

        font-family: var(--text-main-family);
        font-size: var(--text-main-size);
        font-weight: var(--text-main-weight);
        color: #000000;

        .label {
            width: 40%;
            text-align: right;
            color: #777777;
        }
        .value {
            width: 60%;
            text-align: left;
            text-wrap: wrap;
            overflow-wrap: break-word;
        }
    }
</style>
