<script>
    import { getProxyMediaUrl } from "@libs/DungeonsCommunication/services_requests/media_requests";
    import { isUrlMediaFile, isUrlVideo, getMediaFilename } from "@libs/utils";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {string} */
        export let media_file_url;
    
        /** @type {boolean} */
        let is_video;

        $: validateMediaUrl(media_file_url);

    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/


        /**
         * Downloads the media file on the users device
         */
        const downloadMedia = () => {
            const media_filename = getMediaFilename(media_file_url);
            const media_url = getProxyMediaUrl(media_file_url);

            if (!media_url || !media_filename) return;

            const download_element = document.createElement("a");

            download_element.href = media_url;
            download_element.download = media_filename;

            download_element.click();
        }

        /**
         * @param {MouseEvent} e
         */
        const handleBackgroundClick = e => {
            if (e.currentTarget === e.target) {
                resetMediaUrl();
            }   
        }

        const handleMediaDoubleClick = () => {
            downloadMedia();
        }

        /**
         * sets the is_video property. if the url is not a valid media(not image or video) then it throws an error, except if the url is empty or undefined
         * @param {string} url
         * @returns {boolean}
         * @throws {Error} - if the url is not a valid media(not image or video)
         */
        function validateMediaUrl(url) {
            if (url === "" || url === undefined) {
                is_video = false;
                return false;
            }

            if (!isUrlMediaFile(url)) {
                throw new Error("the url is not a valid media(not image or video)");
            }

            is_video = isUrlVideo(url);

            return true;
        }

        const resetMediaUrl = () => {
            media_file_url = "";
        }
    
    /*=====  End of Methods  ======*/
</script>

{#if media_file_url !== undefined && media_file_url !== ""}
    <div on:click={handleBackgroundClick} id="media-preview-wrapper">
        <div id="media-preview-modal">
            <div id="mpm-modal-controls" class:adebug={false}>
                <button id="mpm-modal-close-btn" on:click={resetMediaUrl}>
                    <svg viewBox="0 0 50 50">
                        <path d="M1 1L49 49M49 1L1 49"/>
                    </svg>
                </button>
            </div>
            <div id="mpm-modal-media" on:dblclick={handleMediaDoubleClick} >
                {#if is_video}
                    <video src="{getProxyMediaUrl(media_file_url)}" muted controls>
                        <track kind="captions" />   
                    </video>
                {:else}
                    <img src="{getProxyMediaUrl(media_file_url)}" alt="media preview" loading="lazy">
                {/if}
            </div>
        </div>
    </div>
{/if}

<style>
    :global(:has(> #media-preview-wrapper)) {
        position: relative;
    }

    #media-preview-wrapper {
        position: absolute;
        top: 0;
        left: 0;
        container-type: size;
        container-name: media-preview;
        display: grid;
        background: var(--grey-t);
        width: 100%;
        height: 100%;
        z-index: var(--z-index-t-2);
        place-items: center;
    }

    #media-preview-modal {
        container-type: size;
        display: grid;
        width: 80cqw;
        height: 80cqh;
        background: var(--glass-gradient);
        grid-template-rows: 7% 90%;
        grid-template-columns: 1fr;
        border: .5px solid var(--grey-7);
        border-radius: var(--border-radius);
        gap: var(--vspacing-1);
    }

    #mpm-modal-controls {
        display: flex;
        justify-content: flex-end;
        align-items: center;
        padding: var(--vspacing-1);
    }

    #mpm-modal-close-btn {
        background: transparent;
        border: none;
        outline: none;
        max-width: 2cqw;
        max-height: 2cqw;
        padding: calc(0.5 * var(--vspacing-1));
    }

    #mpm-modal-close-btn svg {
        width: 100%;
        height: 100%;
    }

    #mpm-modal-close-btn svg path {
        stroke: var(--grey-1);
        stroke-width: 4px;
        stroke-linecap: round;
        stroke-linejoin: round;
        fill: none;
        transition: all .2s ease-in-out;
    }

    @media (pointer: fine) {
        #mpm-modal-close-btn:hover svg path {
            stroke: var(--main-dark);
        }
    }

    #mpm-modal-media {
        display: flex;
        justify-content: center;
        align-items: center;
        padding: var(--vspacing-1);
    }

    #mpm-modal-media img, #mpm-modal-media video {
        max-width: 100cqw;
        max-height: 80cqh;
        border-radius: var(--border-radius);
        object-fit: contain;
    }

    @container media-preview (max-width: 450px) {
        #media-preview-modal {
            width: 100cqw;
            height: 100cqh;
            grid-template-rows: 18% 80%;
            background: transparent;
            border: none;
        } 

        #mpm-modal-controls {
            padding: var(--vspacing-2);
        }

        #mpm-modal-close-btn {
            max-width: 13cqw;
            max-height: 13cqw;
        }
        
        #mpm-modal-media {
            flex-direction: column;
            justify-content: flex-start;
        }
    }

</style>