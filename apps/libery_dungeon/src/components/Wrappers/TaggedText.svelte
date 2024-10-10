<script>
    export let tag_name;
    export let content_padding = 26;
    export let extra_props = "";
    export let vspacing = "var(--spacing-1)";
</script>

<div class="tagged-text-wrapper" style:row-gap={vspacing}>
    <div class="top-tag-wrapper">
        <i class="html-tag">{`<${tag_name}${ extra_props !== '' ? ' ' : '>'}`}</i>
        {#if extra_props !== ''}
            <i class="html-tag">{@html extra_props}&gt;</i>
        {/if}
    </div>
    <div class="wrapped-content" style:padding-left="{content_padding}px">
        <slot>---</slot>
    </div>
    <i class="html-tag">{`</${tag_name}>`}</i>
</div>

<style>
    .tagged-text-wrapper {
        display: flex;
        flex-direction: column;
    }

    .tagged-text-wrapper .top-tag-wrapper {
        display: flex;
        gap: var(--spacing-1);
    }

    .tagged-text-wrapper i {
        text-align: left;
        width: max-content;
    }


</style>