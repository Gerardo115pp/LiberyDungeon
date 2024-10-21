<script>
    import { getClusterTags } from "@models/DungeonTags";
    import { cluster_tags, last_cluster_domain, refreshClusterTagsNoCheck } from "@stores/dungeons_tags";
    import TagTaxonomyCreator from "../TagTaxonomyComponents/TagTaxonomyCreator.svelte";
    import { onMount } from "svelte";
    import { current_cluster } from "@stores/clusters";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import TaxonomyTags from "../TagTaxonomyComponents/TaxonomyTags.svelte";

    
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
    <section id="tag-taxonomy-creator-section" class="dctt-section">
        <TagTaxonomyCreator />
    </section>
    <article id="dctt-cluster-user-tags"
        class="dungeon-scroll dctt-section"
    >
        {#if cluster_tags_checked}
            {#each $cluster_tags as taxonomy_tags}
                {#if !taxonomy_tags.Taxonomy.IsInternal}
                    <TaxonomyTags 
                        taxonomy_tags={taxonomy_tags}
                        enable_tag_creation
                    />
                {/if} 
            {/each}
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
        display: flex;
        height: 35cqh;
        flex-direction: column;
        row-gap: var(--spacing-2);
        overflow-y: auto;
    }
</style>