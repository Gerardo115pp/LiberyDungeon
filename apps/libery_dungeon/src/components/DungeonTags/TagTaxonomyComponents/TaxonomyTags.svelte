<script>
    import TagGroup from "../Tags/TagGroup.svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The Taxonomy tags composition class.
         * @type {import("@models/DungeonTags").TaxonomyTags}
         */
        export let taxonomy_tags;

        /**
         * Whether to allow the user to create new tags.
         * @type {boolean}
        */
        export let enable_tag_creation = false;

    /*=====  End of Properties  ======*/
    
</script>

<section class="dungeon-taxonomy-content">
    <header class="taxonomy-header">
        <h4>
            {taxonomy_tags.Taxonomy.Name}
        </h4>
    </header>
    <TagGroup 
        dungeon_tags={taxonomy_tags.Tags}
        enable_tag_creator={enable_tag_creation}
        on:tag-selected
    />
</section>

<style>
    .dungeon-taxonomy-content {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-1);
    }
</style>