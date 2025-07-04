<script lang="ts">
    import { createEventDispatcher } from "svelte";

    export let name: string;
    export let imageUrl: string;

    let dispatch = createEventDispatcher();

    function handleClick() {
        dispatch("gameSelected", {});
    }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
    class="game-card"
    style:background-image="url('{imageUrl}')"
    on:click={handleClick}
>
    <div class="info-bar">
        <span>{name}</span>
    </div>
</div>

<style lang="scss">
    .game-card {
        width: 100%;
        height: 100%;

        display: flex;
        flex-direction: column;
        justify-content: flex-end; /* Ensure the info-bar stays at the bottom */

        //border: 2px solid #3b3b41;
        border-radius: 1rem;

        overflow: hidden;
        transition: transform 0.2s ease;

        background-color: #27292f;
        background-image: url("{imageUrl}");
        background-size: cover; /* Ensure the image covers the entire card */
        background-position: center; /* Center the image */

        box-shadow: 0px 0.25rem 0.75rem #00000055;

        .info-bar {
            width: calc(100% - 2rem);
            padding: 1rem;

            overflow: hidden;
            text-align: center;

            background-color: #202126;
            backdrop-filter: blur(0.2rem);

            transition: 0.2s ease;

            font-family: var(--text-main-family);
            font-size: var(--text-main-size);
            font-weight: var(--text-main-weight);
            color: #ffffff;
        }

        &:hover {
            //transform: scale(1.01);
            cursor: pointer;

            .info-bar {
                background-color: #292a30;
            }
        }
    }
</style>
