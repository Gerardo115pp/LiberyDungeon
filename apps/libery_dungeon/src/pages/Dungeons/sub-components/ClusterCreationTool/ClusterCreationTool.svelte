<script>
    import  CctPathSelectionStep from "./CCTPathSelectionStep.svelte";
    import CctClusterCreationStep from "./CCTClusterCreationStep.svelte";
    import { createEventDispatcher } from "svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The path were the new cluster will be created.
         * @type {string}
         */
        let new_cluster_path = ""; // An Elden Ring reference :D
    
        let dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Hanlde the Path selected event from the path selection step.
         * @param {CustomEvent<PathSelectedEventDetail>} e - The path selected event.
         * @typedef {Object} PathSelectedEventDetail
         * @property {string} path - The path selected.
         */
        const handlePathSelected = e => {
            new_cluster_path = e.detail.path;
        }

        /**
         * Handles the new cluster created event from the cluster creation step.
         * @param {CustomEvent} e - The new cluster created event.
        */
        const handleNewClusterCreated = e => {
            dispatch("cluster-created", {});
        }
    
    /*=====  End of Methods  ======*/
</script>

<div id="cluster-creation-tool" class:adebug={false}>
    {#if new_cluster_path === ""}
        <CctPathSelectionStep
            on:new-cluster-path-selected={handlePathSelected}
        />
    {:else}
        <CctClusterCreationStep
            new_cluster_path={new_cluster_path}
            on:new-cluster-created={handleNewClusterCreated}
        />
    {/if}
</div>

<style>
    #cluster-creation-tool {
        display: flex;
        width: 100%;
        height: min(100%, 600px);
        container-type: inline-size;
        background: var(--grey-9);
        flex-direction: column;
        justify-content: center;
        align-items: center;
        padding: var(--vspacing-3);
        gap: var(--vspacing-2);
        border-radius: var(--border-radius);
    }






</style>

