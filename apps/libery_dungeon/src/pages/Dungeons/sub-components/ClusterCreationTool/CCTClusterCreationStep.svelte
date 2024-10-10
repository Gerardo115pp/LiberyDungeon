<script>
    import { createEventDispatcher } from "svelte";
    import { createCluster } from "@models/CategoriesClusters";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The name of the new cluster.
         * @type {string}
         */
        let new_cluster_name = "";

        /**
         * The name of the category where new content will be downloaded.
         * @type {string}
         */
        let new_cluster_filter_category_name = "";

        /**
         * The directory where the new cluster will be created.
         * @type {string}
         */
        export let new_cluster_path = "";
        const cluster_path = new_cluster_path; // just to be safe.

        /**
         * Whether the cluster is ready to be created or not.
         * @type {boolean}
         */
        let cluster_ready = false;

        let dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the click event of the create cluster button.
         * @param {MouseEvent} e - The click event.
         */
        const handleCreateClusterBtnClick = async e => {
            let new_cluster_uuid = crypto.randomUUID();

            let new_cluster = await createCluster(new_cluster_uuid, new_cluster_name, new_cluster_path, new_cluster_filter_category_name);

            if (new_cluster === null) {
                console.error("Error creating cluster.");
                return;
            }

            dispatch("new-cluster-created", {
                new_cluster: new_cluster
            });
        }

        /**
         * Handles the keypress event on the cluster name input.
         * @param {KeyboardEvent} e - The keypress event.
         */
        const handleClusterNameInputKeypress = e => {
            cluster_ready = new_cluster_name !== "" && new_cluster_filter_category_name !== "";
        }

        /**
         * Handles the keypress event on the filter category name input.
         * @param {KeyboardEvent} e - The keypress event.
         */
        const handleFilterCategoryNameInputKeypress = e => {
            cluster_ready = new_cluster_name !== "" && new_cluster_filter_category_name !== "";
        }

    
    /*=====  End of Methods  ======*/
    
    
</script>

<div id="cct-cluster-creation-step">
    <header id="cct-ccs-instruction">
        <h2 id="ccs-instruction-label">
            Create a new dungeon.
        </h2>
        <p id="cct-ccs-instructions">
            We will create you'r new dungeon under '{cluster_path}'. Now just pick a name for it. And also choose a name for the filter category, this is where all the new content will be initially downloaded.
        </p>
    </header>
    <form action="none" id="cct-cluster-final-details">
        <label class="cct-ccs-cluster-field">
            <p class="field-tool-tip">
                Any name you want to give to your new dungeon. Except for ''.
            </p>
            <span>
                Name:
            </span>
            <input type="text" id="cct-ccs-cluster-name-input" bind:value={new_cluster_name} on:keypress={handleClusterNameInputKeypress}/>
        </label>
        <label class="cct-ccs-cluster-field">
            <p class="field-tool-tip">
                The name of the category where all the new content will be initially downloaded.
            </p>
            <span>
                Filter Category:
            </span>
            <input type="text" id="cct-ccs-cluster-filter-category-name-input" bind:value={new_cluster_filter_category_name}  on:keypress={handleFilterCategoryNameInputKeypress}/>
        </label>
        <button disabled={!cluster_ready} id="cct-ccs-create-cluster-button" type="button" class="dungeon-button-1" on:click={handleCreateClusterBtnClick}>
            Create
        </button>
    </form>
</div>

<style>
    #cct-cluster-creation-step {
        display: flex;
        flex-direction: column;
        justify-content: center;
        height: 100%;
        gap: var(--vspacing-4);
    }

    header#cct-ccs-instruction {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--vspacing-1);
    }

    h2#ccs-instruction-label {
        font-size: var(--font-size-h4);
        color: var(--main);
    }

    p#cct-ccs-instructions {
        font-size: var(--font-size-p-small);
        color: var(--grey-3);
        text-align: center;
    }

    form#cct-cluster-final-details {
        display: flex;
        flex-direction: column;
        gap: var(--vspacing-3);
    }

    label.cct-ccs-cluster-field {
        position: relative;
        display: flex;
        align-items: center;
        padding: var(--vspacing-1);
        gap: var(--vspacing-1);
        background: var(--grey-5);

        &:has(input:focus) {
            outline: var(--main) dashed 1px;
        }
            
    }

    p.field-tool-tip {
        position: absolute;
        bottom: 100%;
        font-size: var(--font-size-fineprint);
        background: var(--grey);
        color: var(--grey-3);
        visibility: hidden;
    }

    span {
        font-size: var(--font-size-p);
        color: var(--main);
    }

    input {
        flex-grow: 1;
        font-size: var(--font-size-p-small);
        padding: var(--vspacing-1);
        background: transparent;
        color: var(--main-dark);
        border: none;
        outline: none;
    }

    @media (pointer: fine) {
        label.cct-ccs-cluster-field:hover > p.field-tool-tip {
            visibility: visible;
        }
    }




</style>