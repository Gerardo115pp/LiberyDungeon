<script>
    import viewport from "@components/viewport_actions/useViewportActions";
    import { active_media_index } from "@stores/media_viewer";
    import { createEventDispatcher, onMount, tick } from "svelte";
    import { Media } from "@models/Medias";
    import MediaGalleryThumbnail from "./MediaGalleryThumbnail.svelte";
    import { layout_properties } from "@stores/layout";
    import MediasGalleryExplorer from "./MediasGalleryExplorer.svelte";
    import { browser } from "$app/environment";
    
    /*=============================================
    =            Properties            =
    =============================================*/

        /** @type{Media[]} all the available medias */
        export let all_available_medias;

        /** @type{string} the current category path */
        export let category_path;

        /** @type {number} the amount of medias that will be mounted each time the user reaches the bottom of the scroll */
        const medias_page_size = 100;

        /** @type {boolean} whether or not the component loads more images on scroll limits reached */
        let load_more = false;

        /** @type {number} lower index loaded from the medias list into the medias array */
        let lower_index_loaded = -1;

        /** @type {number} upper index loaded from the medias list into the medias array */
        let upper_index_loaded = -1;

        /** @type {number} the index of the active media but converted to its index on displayed_medias */
        let converted_active_media_index = 0;

        /** @type {boolean} whether or not the click on a thumbnail should be dispatched as selection or deletion of the media */
        export let deletion_mode = false;

        /**
         * @typedef {Object} MediaGalleryStateHolder
         * @property {Media[]} all_available_medias
         * @property {string} category_path
         * @property {boolean} load_more
         * @property {number} lower_index_loaded
         * @property {number} upper_index_loaded
         * @property {number} converted_active_media_index
         * @property {Media[]} displayed_medias
        */

        /** @type {Media[]} this list of medias been displayed */
        let displayed_medias = [];
        $: loadMedias(), all_available_medias;


        /** 
         * when a remote gallery is opened, this object holds the state of the current gallery, when no remote gallery is opened, this object is null 
         * @type {MediaGalleryStateHolder | null}
         * @default null
        */
        let gallery_state = null;
        
        /*----------  Micro Categories Media Explorer  ----------*/
        
            let show_micro_categories_explorer = false;        

            const gallery_dispatcher = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        setTimeout(() => {
            const active_media = document.querySelector(".mg-gg-active-media");

            if (active_media !== null) {
                active_media.scrollIntoView({behavior: "instant", block: "center"});
                load_more = true;
            }
        }, 100)
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
    
        const loadMedias = (append_medias=false) => {
            if (all_available_medias === undefined || all_available_medias?.length === 0 || $active_media_index >= all_available_medias?.length || !browser) return;

            let slice_lower_limit = append_medias ? (lower_index_loaded - medias_page_size) : lower_index_loaded;

            let slice_upper_limit = append_medias ? upper_index_loaded : (upper_index_loaded + medias_page_size);

            if (displayed_medias.length === 0) {
                slice_lower_limit = $active_media_index - Math.floor(medias_page_size / 2);
                slice_upper_limit = Math.max(slice_lower_limit, 0) + medias_page_size;
            }

            slice_lower_limit = Math.max(0, slice_lower_limit);
            slice_upper_limit = Math.min(all_available_medias.length, slice_upper_limit);

            console.log(`slice_lower_limit: ${slice_lower_limit}\nslice_upper_limit: ${slice_upper_limit}\nappend_medias: ${append_medias}\nload_more: ${load_more}\ndisplayed_medias.length: ${displayed_medias.length}\nall_available_medias.length: ${all_available_medias.length}\npagesize: ${medias_page_size}\nconverted_active_media_index: ${converted_active_media_index}\nactive_media_index: ${$active_media_index}\n\n`);
            displayed_medias = all_available_medias.slice(slice_lower_limit, slice_upper_limit);
            converted_active_media_index = $active_media_index - slice_lower_limit;

            lower_index_loaded = slice_lower_limit;
            upper_index_loaded = slice_upper_limit;
        }

        /**
         * Handles image-selected event from the MediaGalleryThumbnail component.
         * @param {CustomEvent<ImageSelectedDetail>} e
         * @typedef {Object} ImageSelectedDetail
         * @property {Media} media
         */
        const handleThumbnailClick = e => {
            /** 
             * @type {Media}
             */
            let media = e.detail.media;

            gallery_dispatcher("thumbnail-click", media);
        }

        /**
         * Loads more medias when the first media enters the viewport after leaving it for the first time.
         * Usefull when the active_media_index is greater page_size
         * @param {Media} media_element the first media element
         */
        const handleFirstMediaEnters = async media_element => {
            if (!load_more) return
            
            loadMedias(true);

            await tick();

            const previous_first_media = document.querySelector(`#mg-gg-media-${media_element.uuid}`);

            if (previous_first_media !== null) {
                previous_first_media.scrollIntoView({behavior: "instant", block: "center"});
            }
        }

        const handleGalleryClose = () => {
            gallery_dispatcher("gallery-close");
        }

        const handleShowGalleryExplorer = () => {
            show_micro_categories_explorer = true;
        }

        /**
         * @param {CustomEvent<{content: Media[], fullpath: string}>} e
         */
        const handleOpenRemoteGallery = e => {
            saveCurrentGalleryState();

            resetComponentState(e.detail.content, e.detail.fullpath);

            console.log(gallery_state);
            show_micro_categories_explorer = false;
        }

        /**
         * @param {Media[]} available_medias
         * @param {string} new_path
         */
        const resetComponentState = (available_medias, new_path) => {
            all_available_medias = available_medias;
            category_path = new_path;
            load_more = false;
            lower_index_loaded = -1;
            upper_index_loaded = -1;
            converted_active_media_index = 0;
            displayed_medias = [];     
        }

        const resetCurrentGalleryState = () => {
            if (gallery_state === null) return;

            all_available_medias = gallery_state.all_available_medias;
            category_path = gallery_state.category_path;
            load_more = gallery_state.load_more;
            lower_index_loaded = gallery_state.lower_index_loaded;
            upper_index_loaded = gallery_state.upper_index_loaded;
            converted_active_media_index = gallery_state.converted_active_media_index;
            displayed_medias = gallery_state.displayed_medias;

            gallery_state = null;
        }

        const saveCurrentGalleryState = () => {
            if (gallery_state !== null) return;

            gallery_state = {
                all_available_medias,
                category_path,
                load_more,
                lower_index_loaded,
                upper_index_loaded,
                converted_active_media_index,
                displayed_medias
            }
        }
    
    /*=====  End of Methods  ======*/

</script>

<div id="medias-gallery">
    <menu id="mg-medias-controls">
        {#if layout_properties.IS_MOBILE}
            <div class="mg-mc-control-item">
                <button id="mg-mc-gallery-close-btn" on:click={handleGalleryClose}>
                    Close
                </button>
            </div>
        {/if}
        <div class="mg-mc-control-item">
            <span class="mg-mc-loaded-medias">
                Loaded medias: <strong>{displayed_medias?.length || '0'}</strong>
            </span>
        </div>
        {#if !show_micro_categories_explorer}
            <div class="mg-mc-control-item">
                <button id="mg-mc-media-deletion-toggle" class="mg-mc-gallery-btn" on:click={() => deletion_mode = !deletion_mode}>
                    Deletion mode: <strong>{deletion_mode ? 'ON' : 'OFF'}</strong>
                </button>
            </div>
        {/if}
        <div class="mg-mc-control-item">
            <button id="mg-mc-media-explorer-toggle" class="mg-mc-gallery-btn" on:click={handleShowGalleryExplorer}>
                Categories Explorer: <strong>{show_micro_categories_explorer ? 'ON' : 'OFF'}</strong>
            </button>
        </div>
        {#if gallery_state !== null}
            <div class="mg-mc-control-item">
                <button id="mg-mc-media-explorer-toggle" class="mg-mc-gallery-btn" on:click={resetCurrentGalleryState}>
                    Reset Gallery
                </button>
            </div>
        {/if}
    </menu>
    {#if !show_micro_categories_explorer}
        <ul id="mg-gallery-grid" class="libery-scroll">
            {#each displayed_medias as media, h}
                {#if h === 0}
                    <div id="mg-gg-media-{media.uuid}" class="mg-gg-media-wrapper" on:viewportEnter={() => handleFirstMediaEnters(media)}  use:viewport>
                        <MediaGalleryThumbnail media={media} category_path={category_path} deletion_mode={deletion_mode} on:image-selected={handleThumbnailClick}/>
                    </div>
                {:else if h !== displayed_medias.length - 1}
                    <div id="mg-gg-media-{media.uuid}" class="mg-gg-media-wrapper" class:mg-gg-active-media={h === converted_active_media_index && $active_media_index > 10} >
                        <MediaGalleryThumbnail media={media} category_path={category_path} deletion_mode={deletion_mode} on:image-selected={handleThumbnailClick}/>
                    </div>
                {:else}
                    <div id="mg-gg-media-{media.uuid}" class="mg-gg-media-wrapper"  on:viewportEnter={() => loadMedias()}  use:viewport>
                        <MediaGalleryThumbnail media={media} category_path={category_path} deletion_mode={deletion_mode} on:image-selected={handleThumbnailClick}/>
                    </div>
                {/if}
            {/each}
        </ul>
    {:else}
        <MediasGalleryExplorer />
    {/if}
</div>

<style>
    #medias-gallery {
        container-type: size;
        display: flex;
        width: 100%;
        height: 100%;
        background: linear-gradient(90deg, var(--grey) 0%, var(--grey) 9%, hsla(0, 0%, 5%, .1) 29%, var(--grey) 80%);
        flex-direction: column;
        gap: var(--vspacing-1);
        border-radius: var(--border-radius);
        overflow: hidden;
        backdrop-filter: blur(5px);
    }

    #mg-medias-controls {
        display: flex;
        height: 13%;
        width: 100%;
        margin: 0;
        align-items: center;
        gap: var(--vspacing-2);
        background: hsla(0, 0%, 5.9%, .6);
    }

    #mg-mc-gallery-close-btn {
        padding: var(--vspacing-1) var(--vspacing-2);
        background: transparent;
        border: 1px solid var(--main);
        color: var(--main);
        border-radius: var(--border-radius);
        transition: background .2s ease-in-out;
    }

    .mg-mc-gallery-btn {
        padding: var(--vspacing-1) var(--vspacing-2);
        background: transparent;
        border: 1px solid var(--main);
        color: var(--main);
        border-radius: var(--border-radius);
        transition: background .2s ease-in-out;

    }


    @media(pointer: fine) {
        .mg-mc-gallery-btn:hover {
            background: var(--main);
            color: var(--grey-7);
        }
    }

    #mg-gallery-grid {
        display: grid;
        width: 100%;
        height: 87%;
        grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        grid-auto-rows: 100px;
        gap: var(--vspacing-1);
        padding: 0 var(--vspacing-1);
        overflow: auto;
    }

    :global(#mg-gallery-grid .lazy-wrapper) {
        width: 100%;
        height: 100%;
        border: .5px solid var(--grey-4);
        border-radius: var(--border-radius);
        overflow: hidden;
    }

    @media only screen and (max-width: 768px) {
        #mg-gallery-grid {
            grid-template-columns: repeat(auto-fill, minmax(70px, 1fr));
            grid-auto-rows: 70px;
        }

        #mg-medias-controls {
            gap: var(--vspacing-3);
        }
    }
</style>