<script lang="ts">
    import GameCard from "./GameCard.svelte";
</script>

<div class="game-page">
    <div class="game">
        
    </div>
    
    <div class="sidebar"></div>
</div>

<style lang="scss">
    .game-page {
        width: calc(100vw - 2rem);
        height: calc(100vh - 4rem - 2rem);
        
        padding: 1rem;
        overflow: hidden;

        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
    }
</style>
