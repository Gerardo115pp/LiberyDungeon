<script>
    import { cluster_tags } from '@stores/dungeons_tags';
    import InlineTagGroup from './InlineTagGroup.svelte';
    
    /*=============================================
    =            Properties            =
    =============================================*/

        /**
         * The list of tags to display in the SelectedTags component.
         * @type {import('@models/DungeonTags').DungeonTag[]}
         */
        export let selected_tag_list = [];

    
        /**
         * The selected tags grouped by Taxonomy tag.
         * @type {import('../dungeon_tags').DungeonTags_GroupedTags<InlineTagGroup>[]}
         */ 
        let the_selected_tags;
        $: if (selected_tag_list) {
            transformSelectedTagsList();
        }

        /**
         * The index of the selected tag taxonomy 
         * @type {number}
         */
        let selected_tag_taxonomy_index = 0;

        /**
         * Maps tag taxonomy uuids to names.
         * @type {Map<string, string>}
         */
        const known_taxonomy_names = new Map();

    
    /*=====  End of Properties  ======*/

    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Returns the name of a taxonomy by it's uuid.
         * @param {string} taxonomy_uuid
         * @returns {string | null}
         */
        const getTaxonomyName = taxonomy_uuid => {
            let taxonomy_name = known_taxonomy_names.get(taxonomy_uuid);

            if (taxonomy_name !== undefined) {
                return taxonomy_name;
            }

            for (let taxonomy_tags of $cluster_tags) {
                if (taxonomy_tags.Taxonomy.UUID === taxonomy_uuid) {
                    taxonomy_name = taxonomy_tags.Taxonomy.Name;

                    known_taxonomy_names.set(taxonomy_uuid, taxonomy_name);

                    break;
                }
            }

            return taxonomy_name ?? null;
        }

        /**
         * Transforms the selected_tag_list to the structure of the_selected_tags when the selected_tags_list changes.
         * @returns {void}
         */
        function transformSelectedTagsList() {
            if (selected_tag_list.length === 0) {
                the_selected_tags = [];
            }

            /**
             * @type {Map<string, import('../dungeon_tags').DungeonTags_GroupedTags<InlineTagGroup>>}
             */
            const selected_tags_members = new Map();

            for (let selected_tag of selected_tag_list) {
                let tag_group = selected_tags_members.get(selected_tag.TaxonomyUUID);

                if (tag_group === undefined) {
                    const taxonomy_name = getTaxonomyName(selected_tag.TaxonomyUUID);

                    if (taxonomy_name === null) {
                        console.error(`In SelectedTags.transformSelectedTagsList: Taxonomy name for tag ${selected_tag.TaxonomyUUID} not found.`);
                        continue;
                    }

                    tag_group = {
                        taxonomy_name: taxonomy_name,
                        taxonomy_uuid: selected_tag.TaxonomyUUID,
                        tags: [],
                        component_ref: null,
                    };
                }

                tag_group.tags.push(selected_tag);

                selected_tags_members.set(selected_tag.TaxonomyUUID, tag_group);
            }

            the_selected_tags = Array.from(selected_tags_members.values());
        }
    
    /*=====  End of Methods  ======*/
    
    

</script>

<div class="selected-tags">
    {#if $$slots.headline}
        <header class="sc-header">
            <slot name="headline" />
            {#if $$slots.description}
                <p class="sc-header-description">
                    <slot name="description" />
                </p>
            {/if}
        </header>
    {/if}
    <article class="sc-tag-taxonomies">
        {#each the_selected_tags as selected_tag_group, h} 
            <InlineTagGroup
                the_grouped_tags={selected_tag_group}
            />
        {/each}
    </article>
</div>

<style>
    .selected-tags {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-2);
        color: var(--grey-1);
    }

    :global(header.sc-header > p.sc-header-description:has(> p)) {
        display: contents;
    }

    article.sc-tag-taxonomies {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-1);
    }
</style>