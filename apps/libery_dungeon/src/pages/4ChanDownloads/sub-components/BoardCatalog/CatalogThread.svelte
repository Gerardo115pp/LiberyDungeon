<script>
    import { ChanCatalogThread } from "@models/4Chan";
    import LazyLoader from "@components/LazyLoader/LazyLoader.svelte";
    import { DownloadProgressTransmisor, getProxyMediaUrl } from "@libs/HttpRequests";
    import { DownloadProgress } from "@models/Downloads";
    import { createEventDispatcher, onDestroy } from "svelte";
    import { onMount } from "svelte";
    import { 
        last_started_download,
        download_progress,
        last_enqueued_download
    } from "@stores/downloads";
    import { selected_thread, selected_thread_id } from "@pages/4ChanDownloads/app_page_store";
    import { layout_properties } from "@stores/layout";
    
    /*=============================================
    =            properties            =
    =============================================*/
    
        /** @type {ChanCatalogThread} */
        export let catalog_thread;

        /** @type {boolean} whether the thread download is enqueued */
        let download_enqueued = false;

        /** @type {DownloadProgress} the current download progress */
        let thread_download_progress;

        // thread_download_progress = new DownloadProgress({ download_uuid: catalog_thread.uuid, total_files: 100, downloaded_files: 10, completed: true });

        /** @type {DownloadProgressTransmisor} the websocket transmisor that will receive the download progress */
        let download_progress_transmisor;

        /** @type {import('@models/4Chan').DownloadRegister} information the last time this thread was downloaded, if the thread has never been downloaded, this will be undefined */
        let download_register;
    
        const dispatcher = createEventDispatcher();

        let has_download_status = false;

        $: has_download_status = hasDownloadStatus(), thread_download_progress, download_register, download_enqueued;

        $: resetComponentValues(), catalog_thread;
        
        /*----------  Unsubscribers  ----------*/
        
        let download_uuid_unsubscriber = () => {};

        let download_enqueued_unsubscriber = () => {};
    
    /*=====  End of properties  ======*/
    

    onMount(() => {
        getDownloadRegister();

        download_uuid_unsubscriber = last_started_download.subscribe(checkDownloadStarted);

        download_enqueued_unsubscriber = last_enqueued_download.subscribe(checkDownloadEnqueued);
    });

    onDestroy(() => {
        download_uuid_unsubscriber();
        download_enqueued_unsubscriber();

        if (download_progress_transmisor !== undefined) {
            download_progress_transmisor.disconnect();
        }
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * It suscribes to the download_uuid store, if the download_uuid is the same as the thread's uuid, it calls the listenDownloadProgress method
         * @param download_uuid
         */
        const checkDownloadStarted = download_uuid => {
            if (download_uuid === catalog_thread.uuid) {
                console.debug(`thread ${catalog_thread.uuid} is downloading`);
                download_enqueued = false;
                listenDownloadProgress();
            }
        }

        const checkDownloadEnqueued = enqueued_download => {
            if (enqueued_download === catalog_thread.uuid) {
                download_enqueued = true;
                // download_enqueued_unsubscriber();
            }
        }

        const emitDownloadEvent = () => {
            dispatcher("thread-download", {
                thread: catalog_thread
            });
        }

        const setSelectedThread = () => {
            selected_thread_id.set(catalog_thread.uuid);
            selected_thread.set(catalog_thread);
        }

        const getDownloadRegister = async () => {
            let new_download_register = await catalog_thread.downloadRegister();

            if (!new_download_register.exists) {
                return;
            }

            console.log("download register", new_download_register);

            download_register = new_download_register;
        }

        /**
         * whether the thread has any download information to show 
         * on the download status bar
         * @returns {boolean} 
         */
        function hasDownloadStatus() {
            return (download_register !== undefined && download_register.exists) || thread_download_progress !== undefined || download_enqueued;
        }

        /**
         * handles clicks on the thread's li element and sets the catalog_thread on the selected_thread store if the click was not on a button
         * @param {MouseEvent} event 
         */
        const handleThreadClick = event => {
            console.log("click", event.target.tagName);
            if (event.target.tagName?.toLowerCase() === "img") {
                setSelectedThread();
            }
        }

        const listenDownloadProgress = () => {
            console.debug(`Setting up download progress transmisor for thread ${catalog_thread.uuid}`);
            download_progress_transmisor = new DownloadProgressTransmisor(catalog_thread.uuid);

            download_progress_transmisor.download_progress_callback = updateDownloadProgress;

            download_progress_transmisor.connect();
            console.debug("Transmisor connected");  
        }

        function resetComponentValues() {
            thread_download_progress = undefined;
            download_register = undefined;
            download_enqueued = false;

            if (download_progress_transmisor !== undefined) {
                download_progress_transmisor.disconnect();
            }   
        }

        /**
         * Its called by the download_progress_transmisor when it receives a new download progress update
         * @param new_download_progress
         */
        const updateDownloadProgress = (new_download_progress) => {
            console.debug(`Updating download progress for thread ${catalog_thread.uuid}`, new_download_progress);
            thread_download_progress = new_download_progress;
            download_progress.set(thread_download_progress);

            if (thread_download_progress.Completed) {
                console.debug(`Download for thread ${catalog_thread.uuid} completed`);
                
                getDownloadRegister();
            }
        }
    
    /*=====  End of Methods  ======*/

</script>

<li on:click={handleThreadClick} id="catalog-thread-{catalog_thread.uuid}" class="catalog-thread-li">
    <div id="ct-image-wrapper">
        <LazyLoader className="ct-image-loader image-interact" image_url={getProxyMediaUrl(catalog_thread.image_url)}>
            <img slot="lazy-wrapper-image" let:image_src  src={image_src} alt="{catalog_thread.uuid} thread" style:width="{catalog_thread.teaser_thumb_width}px" style:height="{catalog_thread.teaser_thumb_height}px">
        </LazyLoader>
    </div>
    <div id="ct-teaser-wrapper" class:adebug={false}>
        <div class="ct-tw-description libery-scroll">
            <p>
                <strong class="ct-tw-creation-date">Created on: <i>{catalog_thread.getCreationDate()}</i></strong>
                <br/>
                {@html catalog_thread.description}
            </p>
        </div>
        <div class="ct-tw-status-bar" class:adebug={false}>
            <class class="ct-tw-sb-download-status" style:visibility={has_download_status ? "visible" : "hidden"}>
                {#if thread_download_progress !== undefined && !thread_download_progress.Completed}
                    <span class="ct-tw-sb-ds-download-label"><i>
                        Downloading
                    </i></span>
                    <div class="ct-tw-sb-ds-download-bar-wrapper">
                        <div class="ct-tw-sb-ds-download-bar">
                            <div class="ds-download-bar-progress" style:width="{thread_download_progress.percentComplete()}%"></div>
                        </div>
                    </div>
                {:else if download_register !== undefined && download_register !== null}
                    <span class="ct-tw-sb-ds-download-information">
                        <i>
                            Downloaded {download_register.download_count} medias
                        </i>
                    </span>
                {:else if download_enqueued}
                    <span class="ct-tw-sb-ds-download-information"><i>
                        Enqueued
                    </i></span>
                {/if}
            </class>
            <div class="ct-tw-sb-status-details">
                <span>{!layout_properties.IS_MOBILE ? "Replies" : "R"}: <mark>{catalog_thread.responses}</mark></span>
                <span>{!layout_properties.IS_MOBILE ? "Images" : "M"}: <mark>{catalog_thread.images}</mark></span>
            </div>
            <div class="ct-tw-sb-controls">
                <button id="ct-tw-sb-download-btn" on:click={emitDownloadEvent}>
                    <svg viewBox="0 0 50 50">
                        <path class="download-icon-arrow" d="M25 1V41l-12 -11M25 41l12 -11"/>
                        <path class="download-icon-bottom" d="M1 49H49"/>
                    </svg>
                </button>
            </div>
        </div>
    </div>
</li>

<style>
    li.catalog-thread-li {
        display: grid;
        width: 100%;
        height: 32cqh;
        container-type: size;
        grid-template-columns: 1fr 2fr;
        grid-template-rows: 1fr;
        grid-template-areas: "ct-image-wrapper ct-teaser-wrapper";
        background: var(--grey-8);
        border-radius: var(--border-radius);
        padding: var(--vspacing-2);
        transition: background 0.2s ease-in-out;
    }

    div#ct-image-wrapper {
        grid-area: ct-image-wrapper;
        display: flex;
        background: var(--grey-8);
        background: hsl(from var(--grey-8) h s 8%);
        justify-content: center;
        width: 100%;
        height: 100cqh;
        align-items: center;
        overflow: hidden;
        border-radius: var(--border-radius);
    }

    div#ct-image-wrapper img {
        cursor: pointer;
        max-width: 100%;
        object-fit: contain;
    }

    #ct-teaser-wrapper {
        grid-area: ct-teaser-wrapper;
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        grid-template-rows: 52cqh 37cqh;
        justify-items: center;
        width: 100%;
        height: 100%;
        padding: var(--vspacing-1) var(--vspacing-2);
        row-gap: var(--vspacing-2);
    }

    .ct-tw-description {
        overflow: auto;
        width: 100%;
        height: 100%;
        grid-column: 1 / span 2;
        grid-row: 1 / span 1;
    }

    .ct-tw-description p {
        font-size: var(--font-size-1);
        width: 100%;
        height: 100%;
        white-space: pre-wrap;
        overflow-wrap: anywhere;
        word-break: break-word;
    }

    .ct-tw-creation-date {
        font-family: var(--font-decorative);
        font-size: var(--font-size-1);
        color: var(--main-dark);
    }

    .ct-tw-creation-date i {
        color: var(--grey-4);
    }

    .ct-tw-status-bar {
        display:  grid;
        grid-template-columns: 50% 50%;
        grid-template-rows: 20% 70%;
        width: 100%;
        grid-column: 1 / span 2;
        grid-row: 2 / span 1;
        row-gap: calc(0.5 * var(--vspacing-1));
    }

    .ct-tw-sb-download-status {
        display: flex;
        font-family: var(--font-decorative);
        grid-column: 1 / -1;
        grid-row: 1 / span 1;
        column-gap: var(--vspacing-1);
    }

    .ct-tw-sb-download-status span {
        width: min(max-content, 27%);
        color: var(--success);
        font-size: var(--font-size-1);
    }

    .ct-tw-sb-ds-download-bar-wrapper {
        display: flex;
        width: 70%;
        height: 100%;
        align-items: flex-end;
        padding: 2px 0
    }

    .ct-tw-sb-ds-download-bar {
        width: 100%;
        height: 2cqh;
        background: var(--grey-4);
        border-radius: var(--border-radius);
    }

    .ds-download-bar-progress {
        height: 100%;
        background: var(--success-3);
        border-radius: var(--border-radius);
    }

    .ct-tw-sb-ds-download-information {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--main-dark-color-7);
        font-size: var(--font-size-1);
    }

    .ct-tw-sb-status-details {
        display: flex;
        flex-direction: column;
        justify-content: center;
    }
    
    .ct-tw-sb-status-details span {
        font-family: var(--font-decorative);
        color: var(--main-dark);
    }

    .ct-tw-sb-status-details span mark {
        font-family: var(--font-read);
        color: var(--main-7);    
        background: transparent;
    }

    .ct-tw-sb-controls {
        display: grid;
        grid-template-columns: 33% 33% 33%;
        justify-content: center;
        align-items: center;
    }

    .ct-tw-sb-controls button {
        width: 70%;
        height: 70%;
        background: hsl(from var(--grey-8) h s 10%);
        border: none;
        transition: all 0.2s ease-in-out;
        padding: var(--vspacing-1);
    }

    @media (pointer: fine) {

        .ct-tw-sb-controls button:hover {
            background: hsl(from var(--grey-8) h s 15%);
        }

        li.catalog-thread-li:has(button:hover) {
            background: var(--grey-8);
        }
    }

    .ct-tw-sb-controls button svg {
        width: 100%;
        height: 100%;
    }

    .ct-tw-sb-controls button svg path {
        stroke: var(--main-dark);
        fill: none;
        stroke-width: 2px;
        stroke-linecap: round;
        stroke-linejoin: round;
    }

    .ct-tw-sb-controls button svg path.download-icon-bottom {
        stroke: var(--main);
    }

    #ct-tw-sb-download-btn {
        grid-column: 3 / span 1;
    }

    @container (max-width: 450px) {
        #ct-teaser-wrapper {
            padding: var(--vspacing-1) 0 0 var(--vspacing-2);
        }

        .ct-tw-sb-download-status span {
            font-size: var(--font-size-2);
            font-weight: bold;
        }

        .ct-tw-status-bar {
            grid-template-columns: 23% 77%;
        }
    }

    @supports not (hsl(from var(--grey-8) h s 10%)) {
        .ct-tw-sb-controls button {
            background: var(--grey);
        }
    }

</style>