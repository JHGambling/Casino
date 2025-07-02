<script lang="ts">
    import GameCard from "./GameCard.svelte";
</script>

<div class="games-container">
    <GameCard
        name="Slotty 1"
        imageUrl="/assets/images/slot_example_image.webp"
    />
    <GameCard
        name="Slotty 2"
        imageUrl="/assets/images/slot_example_image.webp"
    />
    <GameCard
        name="Slotty 3"
        imageUrl="/assets/images/slot_example_image.webp"
    />
    <GameCard
        name="Slotty 4"
        imageUrl="/assets/images/slot_example_image.webp"
    />
    <GameCard
        name="Blackjack Tisch 1"
        imageUrl="/assets/images/blackjack_example_image.webp"
    />
    <GameCard
        name="Rolly Roulette 1"
        imageUrl="/assets/images/roulette_example_image.png"
    />
    <GameCard
        name="Rolly Roulette 2"
        imageUrl="/assets/images/roulette_example_image.png"
    />
    <GameCard
        name="Russisch Roulette"
        imageUrl="/assets/images/russian_roulette_example_image.webp"
    />
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
