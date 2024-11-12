<script>
    import MediaFileIcon from "@components/UI/icons/media_file_icon.svelte";
    import GlitchyLoader from "@components/UI/Loaders/GrlichyLoader.svelte";
    import { MediaFile } from "@libs/LiberyUploads/models";
    import { onMount, onDestroy, createEventDispatcher } from "svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {MediaFile} */
        export let media_file;
        /** @type {number} */
        export let media_index;

        /**
         * Whether or not the media has been uploaded
         * @type {boolean}
        */
        let this_is_uploaded = false;

        /**
         * Whether the media upload failed
         * @type {boolean}
        */
        let this_upload_failed = false;


        
        /*=============================================
        =            Large File Uploads            =
        =============================================*/
        
            /**
             * How much of the file has been uploaded(0,1). Used for large files
             * @type {number}
             */
            let upload_progress = 0;

            /**
             * Whether the server has received the entire file. 
             * But is still processing it. Used for large files
             * @type {boolean}
             */
            let file_sent = false;
        
        
        /*=====  End of Large File Uploads  ======*/
        
        

        
        const media_deleted_emiter = createEventDispatcher();

    /*=====  End of Properties  ======*/

    onMount(() => {
        media_file.onUploaded = handleMediaFileLoaded;
        media_file.onChunkUploadProgress = handleUploadProgress;
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        const handleMediaDeleted = () => {
            media_deleted_emiter("media-deleted", {media_index});
        }

        /**
         * Handles the media file Loaded event.
        */
        const handleMediaFileLoaded = () => {
            this_is_uploaded = true;
        }

        /**
         * Handles the progress of the file upload
         * @param {number} chunks_uploaded
         * @param {number} total_chunks
        */
        const handleUploadProgress = (chunks_uploaded, total_chunks) => {
            upload_progress = chunks_uploaded / total_chunks;
            
            if ((total_chunks - chunks_uploaded) === 1) {
                file_sent = true;
            }
        }
    
    /*=====  End of Methodsk  ======*/
    
</script>

<li 
    style:position="relative" 
    data-id="{media_index}" 
    class="uploaded-media-file" 
    class:file-uploaded={this_is_uploaded}
    class:upload-failed={this_upload_failed}
    class:large-file={media_file.IsFileTooLarge}
>
    <div class="umf-overlay">
        <div class="umf-overlay-top-bar">
            <button class="umf-overlay-delete-file" on:click={handleMediaDeleted}>
                <svg viewBox="0 0 50 50">
                    <path d="M1 1L49 49M1 49L49 1"/>
                </svg>
            </button>
        </div>
        <div class="umf-overlay-status-bar">
            {#if !media_file.isImage()}
                 <h5 class="umf-media-type-label">
                    {#if media_file.IsFileTooLarge}
                        large video:(<i>{media_file.HumanReadableSize}</i>)
                    {:else}
                        video
                    {/if}
                 </h5>
            {/if}
        </div>
    </div>
    <div class="umf-upload-progress-overlay">
        <div id="umf-upload-progress-bar"
            style:scale="{upload_progress} 1"
            role="progressbar"
            aria-valuenow={upload_progress}
            aria-valuemin="0"
            aria-valuemax="1"
        >
            {#if file_sent}
                <GlitchyLoader 
                    loader_color="var(--accent-3)"
                    label={this_is_uploaded ? "Complete" : "Finishing"}
                />
            {/if}
        </div>
    </div>
    {#if !media_file.IsFileTooLarge}
        {#if media_file.isImage()}
            <img src={media_file.Src} alt={media_file.name} />
        {:else}
            <!-- svelte-ignore a11y-media-has-caption -->
            <video src={media_file.Src} controls={false}>
                <track kind="caption"/>
            </video>
        {/if}
    {:else}
        <MediaFileIcon 
        />
    {/if}
</li>

<style>
    .uploaded-media-file {
        display: flex;
        background: var(--grey);
        width: 100%;
        height: 100%;
        justify-content: center;
        align-items: center;
        border: 1px solid var(--main);
        border-radius: var(--border-radius);

        &.large-file {
            padding: var(--spacing-1);
        }
    }


    
    /*=============================================
    =            Control overlay            =
    =============================================*/
    
        .umf-overlay {
            position: absolute;
            container-type: size;
            display: grid;
            width: 100%;
            height: 100%;
            grid-template-columns: 1fr;
            grid-template-rows: 10% 80% 10%;
            top: 0;
            left: 0;
            padding: var(--vspacing-1);
            z-index: var(--z-index-t-2)
        }

        .uploaded-media-file.file-uploaded .umf-overlay {
            background: hsl(from var(--success-4) h s l / 0.3)
        }

        .uploaded-media-file.upload-failed .umf-overlay {
            background: hsl(from var(--danger) h s l / 0.3)
        }

        .umf-overlay-top-bar {
            display: flex;
            grid-area: 1 / 1 / span 1 / span 1;
            justify-content: flex-end;
            align-items: center;
            visibility: hidden;
            transition: all 0.2s ease-in-out;
            opacity: 0.2;
        }

        .umf-overlay:hover .umf-overlay-top-bar {
            visibility: visible;
            opacity: 1;
        }

        button.umf-overlay-delete-file {
            width: 10cqh;
            height: 10cqh;
            display: grid;
            background: var(--danger);
            border: none;
            padding: 0;
            place-items: center;
            transition: all 0.2s ease-in-out;
        }

        button.umf-overlay-delete-file:hover {
            background: var(--danger-4);
        }

        button.umf-overlay-delete-file svg {
            width: 50%;
            height: 50%;
            fill: none;
        }

        button.umf-overlay-delete-file svg path {
            stroke: var(--grey-1);
            stroke-width: 6px;
            stroke-linecap: round;
        }    

        .umf-overlay-status-bar {
            display: flex;
            grid-area: 3 / 1 / span 1 / span 1;
            justify-content: flex-end;
            align-items: center;
            padding: 0 var(--vspacing-1);
        }

        .umf-media-type-label {
            font-family: var(--font-decorative);
            font-size: var(--font-size-2);
            color: var(--main-dark);
            line-height: 1;

        }

        .large-file .umf-media-type-label {
            color: var(--main-dark-color-5);
            font-size: var(--font-size-1);

            & i {
                font-family: var(--font-read);
                font-size: var(--font-size-1);
                color: var(--main-dark-color-8);
            }
        }
    
    /*=====  End of Control overlay  ======*/
    
   
   /*=============================================
   =            Upload progress overlay            =
   =============================================*/
   
        .umf-upload-progress-overlay {
            display: none;
            position: absolute;
            container-type: size;
            width: 100%;
            height: 100%;
        }

        .large-file:not(.upload-failed) .umf-upload-progress-overlay {
            display: flex;
        }

        #umf-upload-progress-bar {
            display: flex;
            width: 100%;
            height: 100%;
            background: hsl(from var(--accent) h s l / 0.3);
            justify-content: center;
            align-items: center;
            transform-origin: left;
            transition: scale 0.2s ease;
            z-index: var(--z-index-t-1);
        }
   
   /*=====  End of Upload progress overlay  ======*/
   
    


    .uploaded-media-file img {
        max-width: 100%;
        max-height: 100%;
        object-fit: contain;
    }

    .uploaded-media-file video {
        max-width: 100%;
        max-height: 100%;
        object-fit: contain;
    }
</style>