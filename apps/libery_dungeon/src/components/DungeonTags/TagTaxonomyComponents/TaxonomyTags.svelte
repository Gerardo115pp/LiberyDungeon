<script>
    import { createDungeonTag, deleteDungeonTag } from "@models/DungeonTags";
    import TagGroup from "../Tags/TagGroup.svelte";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { createEventDispatcher } from "svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The Taxonomy tags composition class.
         * @type {import("@models/DungeonTags").TaxonomyTags}
         */
        export let taxonomy_tags;
        $: console.log("TaxonomyTags: ", taxonomy_tags);

        /**
         * Whether to allow the user to create new tags.
         * @type {boolean}
         */
        export let enable_tag_creation = false;

        const dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the tag created event.
         * @param {CustomEvent<{tag_name: string}>} event
         */ 
        const handleTagCreated = async event => {
            /**
             * @type {import('@models/DungeonTags').DungeonTag | null}
             */
            let new_dungeon_tag = await createDungeonTag(event.detail.tag_name, taxonomy_tags.Taxonomy.UUID);

            if (new_dungeon_tag === null) {
                const variable_environment = new VariableEnvironmentContextError("In TaxonomyTags.handleTagCreated")
                variable_environment.addVariable("triggering_event", event);
                variable_environment.addVariable("taxonomy_tags.Taxonomy.UUID", taxonomy_tags.Taxonomy.UUID);

                const labeled_err = new LabeledError(variable_environment, `Failed to create tag '${event?.detail?.tag_name}'`);

                labeled_err.alert();
                return;
            }

            emitTaxonomyContentChange();
        }

        /**
         * Handles the tag deleted event.
         * @param {CustomEvent<{tag_id: string}>} event
         */
        const handleTagDeleted = async event => {
            
            let tag_id = event?.detail?.tag_id;

            if (tag_id == null) return;

            let tag_deleted = await deleteDungeonTag(tag_id);

            if (!tag_deleted) {
                const labeled_err = new LabeledError("In TaxonomyTags.handleTagDeleted", `Failed to delete tag with id '${tag_id}'`);
                labeled_err.alert();
                return;
            }

            emitTaxonomyContentChange();
        }

        /**
         * Emits an event that should be interpreted as 'the tag taxonomy content has changed'. The taxonomy emits an event with a detail.taxonomy, this
         * contains the tag taxonomy uuid. 
         */
        const emitTaxonomyContentChange = () => {
            dispatch("taxonomy-content-change", {taxonomy: taxonomy_tags.Taxonomy.UUID});
        }
    
    /*=====  End of Methods  ======*/
    
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
        on:tag-created={handleTagCreated}
        on:tag-deleted={handleTagDeleted}
    />
</section>

<style>
    .dungeon-taxonomy-content {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-1);
    }
</style>