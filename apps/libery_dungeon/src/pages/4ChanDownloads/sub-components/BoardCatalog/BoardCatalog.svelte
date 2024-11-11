<script>
    import { last_started_download, download_progress, last_enqueued_download } from "@stores/downloads";
    import { ChanCatalogThread, getBoardCatalog } from "@models/4Chan";
    import CatalogThread from "./CatalogThread.svelte";
    import DownloadNameModal from "./DownloadNameModal.svelte";
    import { DownloadProgress, getCurrentDownloadUUID } from "@models/Downloads";
    import { onDestroy, onMount, tick } from "svelte";
    import { selected_thread_id } from "@pages/4ChanDownloads/app_page_store";
    import { browser } from "$app/environment";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {ChanCatalogThread[]} */
        let catalog_threads = [];

        /** 
         * @type {string} the selected board's name
        */
        export let board_name;

        /** @type {boolean} whether the download name modal is visible */
        let is_download_name_modal_visible = false;

        /**
        * @type {HTMLDivElement} the catalog threads element
        */
        let catalog_threads_element;

        /** 
         * @type {ChanCatalogThread | undefined} the thread that is been downloaded 
         */
        let downloading_thread;

        /**
         * Selected categories cluster
         * @type {import("@models/CategoriesClusters").CategoriesCluster}
         */
        let selected_categories_cluster;

        $: refreshCatalog(board_name);

        
        /*----------  Unsubscribers  ----------*/
        
        let download_progress_unsubscriber = () => {};
    
    
    /*=====  End of Properties  ======*/
    
    onMount(() => {
        download_progress_unsubscriber = download_progress.subscribe(onProgressChange);
    });

    onDestroy(() => {
        download_progress_unsubscriber();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        const cancelDownload = () => {
            is_download_name_modal_visible = false;
            downloading_thread = undefined;
        }

        /**
         * @param {string} download_name the name of the category where content will be downloaded to
         * @returns {Promise<string | null>} the uuid of the download, used to track the download progress
        */
        const downloadThread = async (download_name) => {
            if (downloading_thread === undefined) {
                console.error("downloading thread is undefined");
                return null;
            }

            last_enqueued_download.set(downloading_thread.uuid);
            await tick();

            let download_uuid = await downloading_thread.download(download_name, selected_categories_cluster.DownloadCategoryID, selected_categories_cluster.UUID);

            if (download_uuid === null) {
                console.error("error downloading thread");
                return null;
            }

            last_started_download.set(download_uuid);

            downloading_thread = undefined;

            return download_uuid;
        }

        /**
         * Handles the event of a thread being selected for download
         * @param {CustomEvent<{thread: ChanCatalogThread}>} e
         */
        const handleDownloadThread = e => {
            is_download_name_modal_visible = true;

            downloading_thread = e.detail.thread;
        }

        /**
         * Handles the event of the download name being ready
         * @param {CustomEvent<{download_name: string, cluster: import("@models/CategoriesClusters").CategoriesCluster}>} event
         */
        const handleDownloadNameReady = event => {
            is_download_name_modal_visible = false;
            const { download_name, cluster} = event.detail;

            selected_categories_cluster = cluster;

            downloadThread(download_name)
        }

        /**
         * Is invoked when the download progress changes object changes. Is used to request a new download uuid whenever 
         * the new_download.completed property is true
         * @param {DownloadProgress} new_progress
         */
        const onProgressChange = async new_progress => {
            console.debug("progress change", new_progress);
            if (!new_progress.Completed) {
                return;
            }

            console.debug("requesting new download uuid");
            let new_download_uuid = await getCurrentDownloadUUID();

            if (new_download_uuid === null) {
                console.debug("an error likely happened while requesting a new download uuid");
            }

            console.debug(`new download uuid: ${new_download_uuid}`);
            last_started_download.set(new_download_uuid);
        }

        /**
         * Refreshes the catalog threads
         * @param {string} board_name
         */
        async function refreshCatalog (board_name) {
            if (!browser) return;

            scrollCatalogToTop();

            catalog_threads = await getBoardCatalog(board_name);

            console.log("remounted");
            if ($selected_thread_id !== null) {
                console.log("scrolling to thread", $selected_thread_id);
                scrollToThread($selected_thread_id);
            }
        }

        const scrollCatalogToTop = () => {
            if (catalog_threads_element === undefined || !browser) {
                return;
            }

            console.log("scrolling to top");
            catalog_threads_element.scrollTo({ top: 0, left: 0, behavior: "smooth" });
        }

        /**
         * Scrolls to the thread with the given id
         * @param {string} thread_id
         */
        const scrollToThread = async (thread_id) => {
            if (!browser) return;

            await tick();

            const thread_element_selector = `catalog-thread-${thread_id}`;

            console.debug(`scrolling to thread ${thread_id} with selector ${thread_element_selector}`);

            let thread_element = document.getElementById(thread_element_selector);


            console.log("scrolling to thread", thread_element);            
            if (thread_element === null) {
                return;
            }

            thread_element.scrollIntoView({ behavior: "smooth" });
        }
    
    /*=====  End of Methods  ======*/
</script>

{#if is_download_name_modal_visible}
    <DownloadNameModal 
        on:modal-close={cancelDownload}
        on:download-name-ready={handleDownloadNameReady}
    />
{/if}
<div bind:this={catalog_threads_element} id="board-thread-catalog" class="libery-scroll">
    <ul id="catalog-threads">
        {#each catalog_threads as ct}
            <CatalogThread on:thread-download={handleDownloadThread} catalog_thread={ct} />
        {/each}
    </ul>
</div>

<style>
    #board-thread-catalog {
        overflow: auto;
        container-type: size;
        width: 100%;
        height: 100%;
    }

    #catalog-threads {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 32.5%));
        grid-auto-rows: max-content; 
        list-style: none;
        margin: 0;
        padding: 0;
        gap: var(--vspacing-1);
    }

    /* if the container(board component) is smaller than 301px */
    @container (max-width: 460px) {
        #catalog-threads {
            grid-template-columns: 1fr;
        }
    }
</style>