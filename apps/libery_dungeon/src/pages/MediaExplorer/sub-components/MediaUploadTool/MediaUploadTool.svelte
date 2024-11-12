<script>
    import { MediaFile, MediaUploader } from "@libs/LiberyUploads/models";
    import { current_category, categories_tree } from "@stores/categories_tree";
    import MediaFileItem from "./MediaFileItem.svelte";
    import { media_upload_tool_mounted } from "@pages/MediaExplorer/app_page_store";
    import GridLoader from "@components/UI/Loaders/GridLoader.svelte";
    import { mediaFilesReadableSize } from "@libs/LiberyUploads/utils";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * @type {MediaFile[]}
         */
        let new_media_files = [];


        /**
         * The input element that will be used to select files
         * @type {HTMLInputElement}
         */
        let files_mount_point;

        /**
         * Whether the files are currently been read
         * @type {boolean}
         */
        let files_been_read = false;

        /**
         * The amount of medias that have been uploaded
         * @type {number}
         */
        let medias_uploaded = 0;
        
        /*=============================================
        =            Styles            =
        =============================================*/
        
            let files_over_mut_content = false

        /*=====  End of Styles  ======*/
        
    
    
    /*=====  End of Properties  ======*/
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Deletes a media file from the new_media_files array by its index
         * @param {CustomEvent<{media_index: number}>} e 
         */
        const deleteMediaIndex = e => {
            const { media_index } = e.detail;
            const filtered_files = [];

            for (let h = 0; h < new_media_files.length; h++) {
                if (h !== media_index) {
                    filtered_files.push(new_media_files[h]);
                }
            }

            new_media_files = filtered_files;
        }

        /**
         * Handles the file uploaded event from the MediaUploader
         * @param {import('@libs/LiberyUploads/models').MediaFile} file
         * @param {number} index
         */
        const handleFileUploaded = (file, index) => {
            medias_uploaded++;
        }

        /**
         * @param {DragEvent} event
         */
        const handleFilesDragOver = event => {
            event.preventDefault();

            files_over_mut_content = true;
        }

        /**
         * Handles the drop event, reads the files and adds them to the new_media_files array
         * @param {DragEvent} e 
         */
        const handleFilesDrop = e => {
            e.preventDefault();

            if (e.dataTransfer?.files == null) {
                console.error("handleFilesDrop: e.dataTransfer is null or has no files");
                return;
            }

            const files = e.dataTransfer.files;

            readFileList(files);
        }

        /**
         * Handles the change event of the file input, reads the files and adds them to the new_media_files array
         * @type {import('svelte/elements').ChangeEventHandler<HTMLInputElement>}
         */
        const handleNewFiles = e => {
            if (!(e.target instanceof HTMLInputElement)) {
                console.error("handleNewFiles: e.target is not an instance of HTMLInputElement");
                return;
            }

            /**
             * @type {HTMLInputElement}
             */
            const file_input = e.target;

            if (file_input.type !== "file" || file_input.files == null) {
                console.error("handleNewFiles: file_input is not an instance of HTMLInputElement or was not 'file' type. likely called from a wrong event");
                return;
            }

            const files = file_input.files; 

            readFileList(files);
        }

        /**
         * Reads a FileList and loads them on the new_media_files array.
         * @param {FileList} file_list
         */
        const readFileList = async (file_list) => {
            /** @type {MediaFile[]} */
            const opened_files = [];
            files_been_read = true;

            for (let f of file_list) {
                let new_media_file = new MediaFile(f);
             
                if (!new_media_file.IsFileTooLarge) {
                    await new_media_file.loadURL();
                }

                opened_files.push(new_media_file);
            }

            new_media_files = [...new_media_files, ...opened_files];
            files_been_read = false;
        }

        
        const endUpload = async () => {
            media_upload_tool_mounted.set(false);
        }

        const uploadFiles = async () => {
            if ($current_category == null) {
                console.error("In MediaUploadTool.uploadFiles: current_category is null");
                return;
            }
            
            if (new_media_files.length === 0) {
                return;
            }

            const media_uploader = new MediaUploader(new_media_files);

            media_uploader.onAllUploaded = endUpload;
            media_uploader.onFileUploaded = handleFileUploaded

            await media_uploader.startUpload($current_category.uuid);
        }
    
    
    /*=====  End of Methods  ======*/

</script>

<aside id="media-upload-tool-wrapper">
    <div id="media-upload-tool">
        <h3 id="mut-title">
            Upload medias to this category
        </h3>
        <div id="mut-content"
            on:dragover={handleFilesDragOver} 
            on:drop={handleFilesDrop} 
            on:dragleave={() => files_over_mut_content = false}
            class:content_dragover={files_over_mut_content}
        >
            {#if new_media_files.length > 0}
                <ul id="uploaded-content">
                    {#each new_media_files as mf, h}
                        <MediaFileItem media_file={mf} media_index={h} on:media-deleted={deleteMediaIndex}/>
                    {/each}
                </ul>
                <div id="mut-status-bar">
                    <div id="mut-upload-content-controls">
                        <div id="mut-upload-content-btn" class="button-1-wrapper">
                            <button on:click={uploadFiles} class="button-1">
                                Upload
                            </button>
                        </div>
                        <div id="mut-add-medias-btn" class="button-g-wrapper">
                            <button on:click={() => files_mount_point?.click()} class="button-g">
                                Add more
                            </button>
                        </div>
                    </div>  
                    <div id="status-data-wrapper">
                        <h4 id="uploads-counter">
                            Uploaded {medias_uploaded} of {new_media_files.length}. Size (<span class="size-tag">~{mediaFilesReadableSize(new_media_files)}</span>)
                        </h4>
                    </div>
                </div>
            {:else if !files_been_read}
                <div id="no-medias-uploaded-open-dialog" class="button-1-wrapper">
                    <button on:click={() => files_mount_point?.click()} class="button-1">
                        Upload medias
                    </button>
                </div>
            {:else}
                <GridLoader />
            {/if}
        </div>
    </div>
    <input 
        id="files-mount-point"
        bind:this={files_mount_point}
        type="file" 
        on:change={handleNewFiles}
        accept="image/*, video/*"
        multiple
    >
</aside>

<style>
    #media-upload-tool-wrapper {
        display: grid;
        container-type: size;
        width: 100%;
        height: 100%;
        place-items: center;
        background: hsl(from var(--grey) h s l / 0.8);
    }
    
    #media-upload-tool {
        position: relative;
        background: var(--grey-9);
        width: 80cqw;
        height: 80cqh;
        z-index: var(--z-index-t-2);
        padding: var(--vspacing-1);
        
    }

    #media-upload-tool > h3 {
        height: 10%;
        font-family: var(--font-read);
        text-align: center;
        padding: var(--vspacing-2) 0;
    }
    
    #mut-content {
        width: 100%;
        container-type: size;
        height: 90%;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        border-radius: var(--border-radius);
    }

    #mut-content.content_dragover {
        border: 2px solid var(--main);
    }

    ul#uploaded-content {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(10cqw, 1fr));
        grid-auto-rows: 10cqw;
        gap: var(--vspacing-2);
        width: 100%;
        height: 85%;
        list-style: none;
        padding: var(--vspacing-2);
        margin: 0;
        overflow-y: auto;
    }
    
    #mut-status-bar {
        width: 100%;
        display: grid;
        height: 15%;
        grid-template-columns: repeat(4, 1fr);
    }

    #mut-upload-content-controls {
        grid-column: 2 / span 2;
        display: flex;
        flex-direction: row-reverse;
        justify-content: center;
        align-items: center;
        gap: var(--vspacing-2);
    }

    #status-data-wrapper {
        width: 100%;
        display: flex;
        grid-column: 4 / span 1;
        justify-content: center;
        align-items: center;

        & span.size-tag {
            font-family: var(--font-read);
            font-size: var(--font-size-2);
            color: var(--main-3);
        }
    }

    #uploads-counter {
        font-family: var(--font-decorative);
        font-size: var(--font-size-2);
        color: var(--grey-1);
    }

    #files-mount-point {
        display: none;
    }

    @media only screen and (max-width: 768px) {
        #media-upload-tool {
            width: 98cqw;
            height: 94cqh;
        }

        #uploaded-content {
            height: 75%;
            grid-template-columns: repeat(auto-fill, minmax(42cqw, 2fr));
            gap: var(--vspacing-1);
            grid-auto-rows: 48cqw;
        }

        #mut-status-bar {
            height: 25%;
            grid-template-columns: repeat(3, 1fr);
        }

        #mut-upload-content-controls {
            grid-column: 1 / span 2;
        }

        #status-data-wrapper {
            grid-column: 3 / span 1;
        }
    }
</style>