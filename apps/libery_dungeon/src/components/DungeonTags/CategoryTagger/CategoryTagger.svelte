<script>
    import TagTaxonomyCreator from "../TagTaxonomyCreator.svelte";
</script>

<dialog open id="dungeon-category-tagger-tool" class="libery-dungeon-window">
    <section id="tag-taxonomy-creator-section">
        <TagTaxonomyCreator />
    </section>
</dialog>

<style>
    #dungeon-category-tagger-tool {
        width: clamp(400px, 70dvw, 1440px);
        height: calc(calc(100dvh - var(--navbar-height)) * 0.8);
        container-type: size;
        padding: var(--spacing-1);
        z-index: var(--z-index-t-1);
    }
</style>