<script>
    import { getClusterTags, getTaxonomyTagsByUUID } from "@models/DungeonTags";
    import { cluster_tags, last_cluster_domain, refreshClusterTagsNoCheck } from "@stores/dungeons_tags";
    import TagTaxonomyCreator from "../TagTaxonomyComponents/TagTaxonomyCreator.svelte";
    import { onMount } from "svelte";
    import { current_cluster } from "@stores/clusters";
    import { current_category } from "@stores/categories_tree";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import TaxonomyTags from "../TagTaxonomyComponents/TaxonomyTags.svelte";
    import ClusterPublicTags from "./sub-componenets/ClusterPublicTags.svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * Whether the cluster_tags correctness with respect to the current_cluster has been checked.
         * @type {boolean}
         */ 
        let cluster_tags_checked = false;


        let current_cluster_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        if ($current_cluster === null) {
            current_cluster_unsubscriber = current_cluster.subscribe(async (value) => {
                if (value !== null) {
                    await verifyLoadedClusterTags();
                    current_cluster_unsubscriber();
                }
            });
        } else {
            await verifyLoadedClusterTags();
        }
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Handles the tag-taxonomy-created event from the TagTaxonomyCreator component.
         */
        const handleTagTaxonomyCreated = async () => {
            await refreshClusterTagsNoCheck($current_cluster.UUID);
        }

        /**
         * Refreshes the content of a TagTaxonomy.
         * @param {CustomEvent<{taxonomy: string}>} event
        */
        const handleTaxonomyContentChanged = async event => {
            const taxonomy_uuid = event?.detail?.taxonomy;
            console.log("Refreshing: ", event.detail.taxonomy);
            if (taxonomy_uuid == null) return;

            const taxonomy_tags_index = $cluster_tags.findIndex(tag => tag.Taxonomy.UUID === taxonomy_uuid);
            console.log("Index: ", taxonomy_tags_index);

            let new_taxonomy_tags = await getTaxonomyTagsByUUID(taxonomy_uuid);

            if (new_taxonomy_tags === null) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.handleTaxonomyContentChanged");

                variable_environment.addVariable("taxonomy_tags_index", taxonomy_tags_index);
                variable_environment.addVariable("cluster_tags", $cluster_tags);

                const labeled_error = new LabeledError(variable_environment, "Failed to refresh the tag taxonomy content. Closing and opening the tool may solve the issue.", lf_errors.ERR_LOADING_ERROR);

                labeled_error.alert();
                return;
            }

            // Ensure the taxonomy is inserted at the same index as it was before.
            let new_cluster_tags = [];

            if (taxonomy_tags_index > 0) {
                new_cluster_tags = $cluster_tags.slice(0, taxonomy_tags_index);
            }

            new_cluster_tags.push(new_taxonomy_tags);

            if (taxonomy_tags_index < $cluster_tags.length - 1) {
                new_cluster_tags = new_cluster_tags.concat($cluster_tags.slice(taxonomy_tags_index + 1));
            }

            cluster_tags.set(new_cluster_tags);
        }
        
        /**
         * Handles the tag-selected event from the ClusterPublicTags component.
         * @param {CustomEvent<{item_id: string}>} event
         */
        const handleTagSelection = (event) => {
            const tag_id = event.detail.tag_id;

            console.log(`Tag selected: ${tag_id} -> ${$current_category.uuid}`);
        }
    
        /**
         * Verifies the cluster_tags correctness with respect to the current_cluster. And if necessary, updates the cluster_tags.
         * @requires cluster_tags_checked
         */    
        const verifyLoadedClusterTags = async () => {
            if ($last_cluster_domain === $current_cluster.UUID) {
                cluster_tags_checked = true;
            }

            let loaded = await refreshClusterTagsNoCheck($current_cluster.UUID);

            if (!loaded) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.verifyLoadedClusterTags");

                variable_environment.addVariable("current_cluster.UUID", $current_cluster.UUID);
                variable_environment.addVariable("cluster_tags", $cluster_tags);
                variable_environment.addVariable("last_cluster_domain", $last_cluster_domain);

                const labeled_error = new LabeledError(variable_environment, "Failed to load current dungeon tags. Please try again later.", lf_errors.ERR_LOADING_ERROR);

                labeled_error.alert();
            }

            cluster_tags_checked = true;
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<dialog open id="dungeon-category-tagger-tool" class="libery-dungeon-window">
    <section tabindex="-1" id="tag-taxonomy-creator-section" class="dctt-section">
        <TagTaxonomyCreator
            on:tag-taxonomy-created={handleTagTaxonomyCreated}
        />
    </section>
    <article id="dctt-cluster-user-tags"
        class="dungeon-scroll dctt-section"
    >
        {#if cluster_tags_checked}
            <ClusterPublicTags 
                on:tag-selected={handleTagSelection}
                on:taxonomy-content-change={handleTaxonomyContentChanged}
            />
        {/if}
    </article>
</dialog>

<style>
    #dungeon-category-tagger-tool {
        display: flex;
        width: clamp(400px, 70dvw, 1440px);
        height: calc(calc(100dvh - var(--navbar-height)) * 0.8);
        container-type: size;
        flex-direction: column;
        row-gap: var(--spacing-3);
        padding: var(--spacing-1);
        z-index: var(--z-index-t-1);
        outline: none;
    
        & > .dctt-section:not(:last-child) {
            border-bottom: var(--border-thin-grey-8);
        }
    }

    article#dctt-cluster-user-tags {
        height: 35cqh;
        overflow-y: auto;
        border-left: var(--border-thick-main);
        padding: var(--spacing-1) var(--spacing-2);
    }
</style>