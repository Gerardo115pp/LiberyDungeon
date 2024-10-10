<script>
    import { slide } from "svelte/transition";

    export let transition = slide;
    export let key_value;
    export let duration = 300;
</script>

{#key key_value}
    <div class="transition-page-wrapper" in:transition={{ duration, delay: duration }} out:transition={{ duration }} >
        <slot/>
    </div>
{/key}

<style>
    .transition-page-wrapper {
        position: absolute;
        inset: 0;
    }
</style>