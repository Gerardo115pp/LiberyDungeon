<script>
    import { createEventDispatcher, onMount } from "svelte";
    import Input from "@components/Input/Input.svelte";
    import FieldData from "@libs/FieldData";
    import { layout_properties } from "@stores/layout";
    import { getAllCategoriesClusters } from "@models/CategoriesClusters";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        const dispatch = createEventDispatcher();

        /** @type {FieldData} the name of the category where content will be downloaded to */
        let download_name_field = new FieldData("download-name-field", /[a-zA-Z\s\d]{2,}/g, "Download name");

        /**
         * The the posible clusters where the thread can be downloaded to
         * @type {import("@models/CategoriesClusters").CategoriesCluster[]}
         */
        let categories_clusters = [];
    
    /*=====  End of Properties  ======*/
    
    onMount(async () => {
        categories_clusters = await getAllCategoriesClusters();
    });

    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * @param {MouseEvent} event
         */
        const handleBackgroundClick = event => {
            const { target, currentTarget } = event;

            if (target === currentTarget) {
                emitModalCloseEvent();
            }
        }

        const handleDownloadNameReady = () => {
            const download_name = download_name_field.getFieldValue();

            if (download_name === "") {
                return;
            }

            // @ts-ignore
            const cluster_uuid = document.querySelector("#dnm-content--cluster-select")?.value ?? "";

            if (cluster_uuid === "") {
                return;
            }

            const cluster = categories_clusters.find(cluster => cluster.UUID === cluster_uuid);

            dispatch("download-name-ready", {
                download_name,
                cluster
            });
        }

        const emitModalCloseEvent = () => {
            dispatch("modal-close");
        }

    /*=====  End of Methods  ======*/

</script>

<aside on:click={handleBackgroundClick} id="download-name-modal-wrapper" class:adebug={false}>
    <div id="download-name-modal">
        <div id="dnm-modal-controls">
            <button on:click={emitModalCloseEvent} id="dnm-mc-close-modal-btn">
                <svg viewBox="0 0 50 50">
                    <path d="M1 1L49 49M49 1L1 49"/>
                </svg>
            </button>
        </div>
        <div id="dnm-content">
            <Input 
                field_data={download_name_field}
                isSquared={true}
                input_label="New category:"
                border_color="var(--main-6)"
                input_padding="var(--vspacing-1)"
                input_background="var(--grey-9)"
                onEnterPressed={handleDownloadNameReady}
                autofocus
            />
            <label id="dnm-content--cluster">
                <span>
                    You can change the cluster where the content will be downloaded to
                </span>
                <select name="" id="dnm-content--cluster-select">
                    {#each categories_clusters as cluster}
                        <option value="{cluster.UUID}">
                            {cluster.Name}
                        </option>
                    {/each}
                </select>
            </label>
        </div>
    </div>
</aside>

<style>
    #download-name-modal-wrapper {
        position: fixed;
        background: var(--grey-t);
        container-type: size;
        display: grid;
        top: var(--navbar-height);
        left: 0;
        width: 100dvw;
        height: calc(100dvh - var(--navbar-height));
        place-items: center;
        z-index: var(--z-index-t-2);
    }

    #download-name-modal {
        display: grid;
        width: min(500px, 40cqw);
        height: min(300px, 50cqh);
        grid-template-columns: 1fr;
        grid-template-rows: 10% 90%;
        background: var(--grey-9);
        border-radius: var(--border-radius);
        padding: var(--vspacing-2);
    }

    #dnm-modal-controls {
        grid-row: 1 / span 1;
        display: flex;
        width: 100%;
        justify-content: flex-end;
        align-items: center;
        padding: var(--vspacing-1);
    }

    #dnm-mc-close-modal-btn {
        background: none;
        width: 3%;
        padding: 0;
        margin: 0;
        border: none;
    }

    #dnm-mc-close-modal-btn svg {
        width: 100%;
        height: 100%;
    }

    #dnm-mc-close-modal-btn svg path {
        stroke: var(--danger-4);
        stroke-width: 4px;
    }    

    #dnm-content {
        grid-row: 2 / span 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        gap: var(--vspacing-3);
    }

    label#dnm-content--cluster {
        display: flex;
        flex-direction: column;
        gap: var(--vspacing-1);

        & span {
            display: block;
        } 

        & span::after {
            content: ": ";
        }

        & select {
            padding: var(--vspacing-1);
            background: var(--grey-9);
            border: 1px solid var(--main-6);
            border-radius: var(--border-radius);
            color: var(--main-2);
            font-family: var(--font-read);
        }
        
    }

    @container (max-width: 450px) {
        #download-name-modal {
            width: 90cqw;
            height: 50cqh;
        }
    }
</style>