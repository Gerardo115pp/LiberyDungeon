<script>
    import TaxonomyTags from "../../TagTaxonomyComponents/TaxonomyTags.svelte";
    import { cluster_tags } from "@stores/dungeons_tags";
    
</script>

<div id="cpt-wrapper">
    {#each $cluster_tags as taxonomy_tags (taxonomy_tags.Taxonomy.UUID)}
        {#if !taxonomy_tags.Taxonomy.IsInternal}
            <TaxonomyTags 
                on:tag-selected
                on:taxonomy-content-change
                taxonomy_tags={taxonomy_tags}
                enable_tag_creation
            />
        {/if} 
    {/each}
</div>

<style>
    #cpt-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-2);
    }
</style>