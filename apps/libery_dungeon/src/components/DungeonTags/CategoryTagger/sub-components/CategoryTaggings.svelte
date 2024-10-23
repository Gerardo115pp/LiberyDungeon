<script>
    import { current_category } from "@stores/categories_tree";
    import { cluster_tags } from "@stores/dungeons_tags";
    import { getEntityTaggings } from "@models/DungeonTags";
    import { onDestroy, onMount } from "svelte";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The current category's taggings.
         * @type {import("@models/DungeonTags").DungeonTagging[]}
         */ 
        export let current_category_taggings = [];
        $: globalThis.current_category_taggings = current_category_taggings;
        $: handleCategoryTaggingsChange(current_category_taggings);


        /**
         * A map of TagTaxonomy names -> DungeonTags.
         * @type {Map<string, import("@models/DungeonTags").DungeonTagging[]> | null}
         */
        let tag_taxonomy_map = null;
        $: globalThis.tag_taxonomy_map = tag_taxonomy_map;

    
    /*=====  End of Properties  ======*/

    onMount(async () => {
    });

    onDestroy(() => {
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Returns the tag taxonomy map from the taxonomies contained on the current category taggings. 
         * @param {import("@models/DungeonTags").DungeonTagging[]} taggings
         * @returns {Map<string, import("@models/DungeonTags").DungeonTagging[]>}
         */ 
        const getTagTaxonomyMap = taggings => {
            let new_tag_taxonomy_map = new Map();

            const taxonomy_uuid_to_name_lookup = new Map();

            for (let tagging of taggings) {
                let tag_taxonomy_name = taxonomy_uuid_to_name_lookup.get(tagging.Tag.TaxonomyUUID);

                if (tag_taxonomy_name == null) {
                    tag_taxonomy_name = getTagTaxonomyNameByUUID(tagging.Tag.TaxonomyUUID);
                    taxonomy_uuid_to_name_lookup.set(tagging.Tag.TaxonomyUUID, tag_taxonomy_name);
                }

                if (tag_taxonomy_name == null) {
                    console.warn(`Tag taxonomy with UUID '${tagging.Tag.TaxonomyUUID}' not found.`);
                    continue;
                }

                let tag_taxonomy_members = new_tag_taxonomy_map.get(tag_taxonomy_name) ?? [];

                tag_taxonomy_members.push(tagging);

                new_tag_taxonomy_map.set(tag_taxonomy_name, tag_taxonomy_members);
            }

            return new_tag_taxonomy_map;
        }

        /**
         * Returns a the name of a given tag taxonomy or null if it's not among the taxonomies in cluster_tags.
         * @param {string} tag_taxonomy_uuid
         * @returns {string | null}
         */
        const getTagTaxonomyNameByUUID = tag_taxonomy_uuid => {
            let tag_taxonomy_name = null;

            for (let tag_taxonomy_tags of $cluster_tags) {
                if (tag_taxonomy_tags.Taxonomy.UUID === tag_taxonomy_uuid) {
                    tag_taxonomy_name = tag_taxonomy_tags.Taxonomy.Name;
                    break;
                }
            }

            return tag_taxonomy_name;
        }

        /**
         * Handles changes to the cluster tags.
         * @param {import("@models/DungeonTags").DungeonTagging[]} new_taggings
         */
        function handleCategoryTaggingsChange(new_taggings) {
            if (new_taggings.length <= 0) return;

            const updated_tag_taxonomy_map = getTagTaxonomyMap(new_taggings);

            if (updated_tag_taxonomy_map.size <= 0) {
                console.warn("No tag taxonomies found for the current category.");
                return;
            }

            console.log("Updating taxonomy map:", updated_tag_taxonomy_map);

            tag_taxonomy_map = updated_tag_taxonomy_map;
        }

    /*=====  End of Methods  ======*/
    
</script>

{#if tag_taxonomy_map != null && $cluster_tags.length > 0} 
    <div id="cpt-current-category-tags">
        <header id="cpt-cct-header">
            <h4>
                Attributes for <span>{$current_category.name}</span>
            </h4>
        </header>
        {#each tag_taxonomy_map as [taxonomy_name, taxonomy_members]}
            <ol id="{$current_category.uuid}-attribute-{taxonomy_name}"
                class="current-category-attribute dungeon-tag-container"
            >
                <p class="dungeons-field-label">{taxonomy_name}</p>
                {#each taxonomy_members as dungeon_tagging}
                    <DeleteableItem
                        item_color="var(--grey)"
                        item_id={dungeon_tagging.TaggingID}
                        squared_style
                    >
                        <span class="cca-attribute-name">
                            {dungeon_tagging.Tag.Name}
                        </span>
                    </DeleteableItem>
                {/each}
            </ol>
        {/each}
    </div>
{/if}

<style>
    #cpt-current-category-tags {
        display: flex;
        flex-direction: column;
        align-items: center;
        row-gap: var(--spacing-1);

        & > * {
            width: 50%;
        }
    }

    header#cpt-cct-header {
        text-align: center;
        margin-bottom: var(--spacing-2);

        & > h4 {
            color: var(--grey-6);
        }

        & > h4 > span {
            color: var(--main);
        }        
    }

    ol.current-category-attribute {
        color: var(--grey-1);
        line-height: 1;
        align-items: center;

        & :not(p.dungeons-field-label) {
            font-size: var(--font-size-1);
        }

        & span.cca-attribute-name {
            display: flex;
            align-items: center;
            line-height: 1;
        }
    }
</style>