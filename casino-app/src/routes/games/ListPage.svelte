<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import GameCard from "./GameCard.svelte";

    let dispatch = createEventDispatcher();

    function selectGame(id: string) {
        dispatch("select", { id });
    }
</script>

<div class="games-container">
    <GameCard
        name="Slotty"
        imageUrl="/assets/images/slot_example_image.webp"
        on:gameSelected={() => selectGame("slotty")}
    />
    <GameCard
        name="Blackjack"
        imageUrl="/assets/images/blackjack_example_image.webp"
        on:gameSelected={() => selectGame("blackjack")}
    />
    <GameCard
        name="Rolly Roulette"
        imageUrl="/assets/images/roulette_example_image.png"
        on:gameSelected={() => selectGame("rolly")}
    />
    <GameCard
        name="CatchFishy"
        imageUrl="https://i.imgflip.com/91hm6a.jpg"
        on:gameSelected={() => selectGame("catchfishy")}
    />
    <!--<GameCard
        name="Russisch Roulette"
        imageUrl="/assets/images/russian_roulette_example_image.webp"
        on:gameSelected={() => selectGame("russian_roulette")}
    />-->
</div>

<style lang="scss">
    .games-container {
        width: calc(100vw - 2rem);
        height: calc(100vh - 4rem - 2rem);

        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(20rem, 1fr));
        grid-auto-rows: 15rem;
        gap: 1rem;
        padding: 1rem;
        overflow-y: auto;
    }
</style>
