<script lang="ts">
    export let name: string;
    export let imageUrl: string;
</script>

<div class="game-card" style:background-image="url('{imageUrl}')">
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

        background-color: #ff0000;
        background-image: url("{imageUrl}");
        background-size: cover; /* Ensure the image covers the entire card */
        background-position: center; /* Center the image */

        box-shadow: 0px 0.25rem 0.75rem #00000055;

        .info-bar {
            width: calc(100% - 2rem);
            padding: 1rem;
            border-radius: 0.75rem;

            overflow: hidden;
            text-align: center;

            background-color: #27272bcc;
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
                background-color: #35353add;
            }
        }
    }
</style>