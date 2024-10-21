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

    /*=====  End of Properties  ======*/
    
</script>

<section class="dungeon-taxonomy-content">
    <header id="taxonomy-header">
        <h4>
            {taxonomy_tags.Taxonomy.Name}
        </h4>
    </header>
    <TagGroup 
        dungeon_tags={taxonomy_tags.Tags}
    />
</section>