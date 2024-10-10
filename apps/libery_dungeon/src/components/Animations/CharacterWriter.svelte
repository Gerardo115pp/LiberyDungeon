<script>
    import { cubicOut } from "svelte/easing";
    import { fly } from "svelte/transition";
    
    export let text = "";
    let characters = text.split("");
    console.log(characters);

    export let delay = 300;
    export let duration = 1000;
    export let easing = cubicOut;

</script>

<div class="character-writer-wrapper">
    {#each characters as char, h}
        <span style:width="{char === ' ' ? '.5ch' : 'max-content'}" class="blast" in:fly={{x: -50, delay: delay*h, duration, easing, opacity: 0}}>
            {char}
        </span>
    {/each}
</div>

<style>
    .character-writer-wrapper {
        display: inline;
    }

    .blast {
        display: inline-block;
        inset: 0;
    }
</style>