<script>
    import { getProxyMediaUrl } from "@libs/DungeonsCommunication/services_requests/media_requests";
    import { ChanThread, ChanThreadReply, getThreadContent } from "@models/4Chan";
    import { onMount } from "svelte";
    import ThreadReply from "./ThreadReply.svelte";
    import { selected_thread } from "@pages/4ChanDownloads/app_page_store";
    import ChanMediaPreview from "./ChanMediaPreview.svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        export let board_name;
        export let thread_id;

        /** @type {ChanThread} */
        let chan_thread;

        /** @type {string} the url of the media selected for preview*/
        let preview_media_url;

    /*=====  End of Properties  ======*/

    onMount(async () => {
        chan_thread = await refreshThread();
    });

    
    /*=============================================
    =            Methods            =
    =============================================*/

        const handleSelectMediaPreview = e => {
            const { media_file_url } = e.detail;

            setPreviewMediaUrl(media_file_url);
        }

        const refreshThread = async () => {
            let new_thread_data = await getThreadContent(board_name, thread_id);
            return new_thread_data;
        }

        const resetSelectedThread = () => {
            selected_thread.set(null);
        }

        const setPreviewMediaUrl = url => {
            preview_media_url = url;
        }
    
    /*=====  End of Methods  ======*/

</script>


<div id="thread-component-wrapper">
    {#if chan_thread !== undefined}
        <ChanMediaPreview bind:media_file_url={preview_media_url} />
    {/if}
    <div id="thread-controls">
        <button on:click={resetSelectedThread} id="go-back-btn">
            <svg viewBox="0 0 50 50">
                <path d="M37 1L1 25L37 49"/>
            </svg>
        </button>
    </div>
    {#if chan_thread !== undefined}
        <div id="thread-content-wrapper">
            <div id="tcw-op-post">
                <div id="tcw-op-post-status">
                    <span class="tcw-op-ps-date">
                        {chan_thread.getCreationDate()}
                    </span>
                </div>
                <div on:click={() => setPreviewMediaUrl(chan_thread.file)} id="tcw-op-image-wrapper" class="image-interact">
                    <img src="{getProxyMediaUrl(chan_thread.cover_image_url)}" alt="thread {thread_id} cover">
                </div>
                <div id="tcw-op-post-content">
                    <h1 id="tcw-op-post-title">
                        {chan_thread.title}
                    </h1>
                    <p id="tcw-op-post-text" class="libery-scroll">
                        {@html chan_thread.description}
                        <br>
                        <span id="tcw-op-post-id">
                            <strong>id: <i>{chan_thread.uuid.replace('pc', '')}</i></strong>
                        </span>
                    </p>
                </div>
            </div>
            <ul id="thread-replies" class="libery-scroll">
                {#each chan_thread.replies as reply}
                    <ThreadReply on:select-media-preview={handleSelectMediaPreview} thread_reply={reply} op_post_id={chan_thread.uuid}/>
                {/each}
            </ul>
        </div>
    {/if}
</div>

<style>
    #thread-component-wrapper {
        container-type: size;
        container-name: thread-content-component;
        width: 100%;
        height: 100%;
    }

    #thread-controls {
        width: 100%;
        height: 11cqh;
        container-type: size;
        display: flex;
        justify-content: flex-start;
        align-items: center;
        padding: var(--vspacing-2) var(--vspacing-3);
    }

    #go-back-btn {
        width: 60cqh;
        background: transparent;
        border-radius: var(--border-radius);
        border: none;
        outline: none;
        padding: var(--vspacing-1); 
        transition: all .2s ease-in-out;
    }

    @media (pointer: fine) {
        #go-back-btn:hover svg path {
            stroke: var(--main-dark);
        }
    }

    #go-back-btn svg {
        width: 100%;
        height: 100%;
    }

    #go-back-btn svg path {
        stroke: var(--grey-1);
        stroke-width: 4px;
        stroke-linecap: round;
        stroke-linejoin: round;
        fill: none;
        transition: all .2s ease-in-out;
    }

    #thread-content-wrapper {
        width: 100%;
        height: 95cqh;
        display: grid;
        grid-template-columns: 38% 60%;
        grid-template-rows: repeat(2, 1fr);
    }

    
    /*=============================================
    =            Op Post            =
    =============================================*/
    
        #tcw-op-post {
            grid-column: 1 / span 1;
            grid-row: 1 / span 1;
            display: grid;
            max-height: 45cqh;
            grid-template-columns: 45% 45%;
            grid-template-rows: 10% 90%;
            background: var(--grey);
            border-radius: var(--border-radius);
            padding: var(--vspacing-2);
            column-gap: var(--vspacing-2);
        }

        #tcw-op-post-status {
            grid-column: 2 / span 1;
            grid-row: 1 / span 1;
            display: flex;
            justify-content: flex-end;
            align-items: center;
        } 

        #tcw-op-post-status .tcw-op-ps-date {
            font-family: var(--font-decorative);
            font-size: var(--font-size-1);
            color: var(--grey-5);
        }

        #tcw-op-image-wrapper {
            grid-column: 1 / span 1;
            grid-row: 1 / span 2;
            background: var(--grey);
            display: flex;
            justify-content: center;
            border-radius: var(--border-radius);    
            align-items: center;
        }

        #tcw-op-image-wrapper img {
            max-height: 100%;
            max-width: 100%;
            object-fit: cover;
            border-radius: var(--border-radius);
        }

        #tcw-op-post-content {
            grid-column: 2 / span 1;
            grid-row: 2 / span 1;
            display: flex;
            flex-direction: column;
            justify-content: flex-start;
            align-items: flex-start;
            gap: var(--vspacing-1);
        }

        #tcw-op-post-title {
            font-family: var(--font-titles);
            font-size: var(--font-size-2);
            color: var(--main-dark);
            font-weight: lighter;
        }

        #tcw-op-post-text {
            overflow-y: auto;
            overflow-x: hidden;
            font-family: var(--font-read);
            font-size: var(--font-size-1);
            color: var(--main-1);
            height: 25cqh;
            white-space: pre-wrap;
            overflow-wrap: anywhere;
            word-break: break-word;
        }

        #tcw-op-post-id {
            font-family: var(--font-decorative);
            font-size: var(--font-size-1);
            color: var(--grey-4);
            line-height: 3;
        }
        
        #tcw-op-post-id i {
            color: var(--main);
        }

    
    
    /*=====  End of Op Post  ======*/
    
    


    #thread-replies {
        overflow-y: auto;
        max-width: 100%;
        display: flex;
        grid-column: 2 / span 1;
        grid-row: 1 / span 2;
        flex-direction: column;
        list-style: none;
        padding: var(--vspacing-1);
        row-gap: var(--vspacing-1); 
    }

    @container thread-content-component (max-width: 450px) {
        #thread-content-wrapper {
            grid-template-columns: 100%;
            grid-template-rows: 30% 70%;            
        }

        #thread-replies {
            grid-column: 1 / span 1;
            grid-row: 2 / span 1;
        }
    }

</style>