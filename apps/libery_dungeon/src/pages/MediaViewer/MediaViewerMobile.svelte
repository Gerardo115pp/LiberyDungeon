<script>
    import { media_changes_manager, skip_deleted_medias } from "@stores/media_viewer";
    import { category_cache, getCategoryTree } from "@models/Categories";
    import { current_category, categories_tree } from "@stores/categories_tree";
    import MediasGallery from "./sub-components/MediaGallery/MediasGallery.svelte";
    import VideoController from "@components/VideoController/VideoController.svelte";
    import { media_change_types, MediaChangesManager } from "@models/WorkManagers";
    import { getMediaUrl } from "@libs/DungeonsCommunication/services_requests/media_requests";
    import { media_types } from "@models/Medias";
    import { onDestroy, onMount, tick } from "svelte";
    import { pinch, swipe } from 'svelte-gestures';
    import { navbar_hidden } from "@stores/layout";
    import { pop } from "svelte-spa-router";
    import { 
        active_media_index,
        active_media_change,
        random_media_navigation,
        previous_media_index,
    } from "@stores/media_viewer";
    import { fly } from "svelte/transition";

    /* ------ Exports ------ */

    export let params; 

    /* ------ Refs ------ */

    /** @type{HTMLDivElement} media control wrapper */
    let media_control_wrapper;

    /** @type{HTMLVideoElement} video element passed to the video control */
    let video_element;

    
    /*=============================================
    =            State            =
    =============================================*/
    
        let { category_id: url_category_id } = params;

        if (params.media_index !== null) {
                active_media_index.set(parseInt(params.media_index));
        }
    
        let active_media_index_unsubscriber;

        /**
         * @type {number} the zoom level of the current media
         * @default 1
        */
        let media_zoom = 1;

        /** @type{boolean} wether or not to show the media gallery */
        let show_media_gallery = false;

    /*=====  End of State  ======*/
    
    onMount(async () => {
        console.log("Media viewer mobile mounted");
        navbar_hidden.set(true);

        if ($current_category === null) {
            console.log(`current category is null, fetching category tree for ${url_category_id}`);
            getCategoryTree(url_category_id);
        }

        // Attempts to get a cached media index for the current_category, and if the result is different from the current active media index, then it sets the active media index to the cached one
        // the url media_index params has priority over the cached index, so if params.media_index is not null, then no cached index will be used
        let cached_category_index = await category_cache.getCategoryIndex(url_category_id);
        if (cached_category_index !== $active_media_index && params.media_index === null) {
            cached_category_index = Math.max(0, Math.min(cached_category_index, $current_category.content.length - 1));
            active_media_index.set(cached_category_index);
        }

        // Sets the active media change every time the active media index changes
        active_media_index_unsubscriber = active_media_index.subscribe(new_index => {
            if ($current_category === null) return;

            if (new_index >= $current_category.content.length || new_index < 0) {
                return;
            }

            let media_change = $media_changes_manager.getMediaChangeType($current_category.content[new_index].uuid);
            active_media_change.set(media_change);
        });

    });

    onDestroy(() => {
        category_cache.addCategoryIndex(url_category_id, $active_media_index);
        navbar_hidden.set(false);
        

        resetComponentSettings(); 

        active_media_index_unsubscriber();
    });

    const automoveMedia = () => console.log("automoveMedia not implemented yet");

    function exitFullScreen() {
        if (!document.fullscreenElement) return;

        if (document.exitFullscreen) {
            document.exitFullscreen();
        } else if (document.webkitExitFullscreen) { // Safari
            document.webkitExitFullscreen();
        } else if (document.msExitFullscreen) { // IE11
            document.msExitFullscreen();
        }
    }

    const getNextNotDeletedMediaIndex = (from_index, forward) => {
        let next_index = from_index;
        const media_count = $current_category.content.length;
        let media_change;
        let media_uuid;

        if ((from_index === (media_count - 1) && forward) || (from_index === 0 && !forward)) return from_index;

        do {
            next_index += forward ? 1 : -1;
            media_uuid = $current_category.content[next_index]?.uuid;
            media_change = $media_changes_manager.getMediaChangeType(media_uuid);
        } while(media_change === media_change_types.DELETED && next_index > 0 && next_index <= (media_count - 1));

        next_index = Math.max(0, Math.min(next_index, media_count - 1));
        media_uuid = $current_category.content[next_index]?.uuid;
        media_change = $media_changes_manager.getMediaChangeType(media_uuid);

        return media_change === media_change_types.DELETED ? from_index : next_index;
    }

    const handleMediaZoom = e => {
        let { center, scale: zoom } = e.detail;

        media_zoom = zoom;

        media_control_wrapper.style.transformOrigin = `${center.x}px ${center.y}px`;
        media_control_wrapper.style.transform = `scale(${media_zoom})`;
    }
    /**
     * 
     * @param {"left"|"right"} direction
     */
    const handleMediaNavigation = async (direction) => {

        if (direction !== "left" && direction !== "right") return;

        if ($random_media_navigation) {
            return handleRandomMediaNavigation(direction);
        }

        let new_index = $active_media_index;

        new_index += direction === "right" ? -1 : 1;
        new_index = Math.max(0, Math.min(new_index, $current_category.content.length - 1));

        if(media_change_types.DELETED === $media_changes_manager.getMediaChangeType($current_category.content[new_index].uuid) && $skip_deleted_medias) {
            let not_deleted_new_index = getNextNotDeletedMediaIndex(new_index, key_combo !== "a");
            new_index = (not_deleted_new_index === new_index) ? $active_media_index : not_deleted_new_index;
        }

        automoveMedia();

        active_media_index.set(new_index);
        history.replaceState(null, "", `#/media-viewer-mobile/${$current_category.uuid}/${$active_media_index}`);

        await tick();
        
        resetMediaConfigs(true);
    }

    const handleMediaSwipe = e => {
        let { direction } = e.detail;
        
        if (direction === "left" || direction === "right") {
            handleMediaNavigation(direction);
        }

        if (direction === "bottom") {
            rejectMedia();   
        }

    }

    const handleRandomMediaNavigation = async direction => {
        let new_index = $active_media_index;

        

        new_index = Math.floor(Math.random() * $current_category.content.length);

        while (new_index === $active_media_index) {
            new_index = Math.floor(Math.random() * $current_category.content.length);
        }

        if (direction === "right") {
            new_index = $previous_media_index;
        }

        previous_media_index.set($active_media_index);

        active_media_index.set(new_index);
        history.replaceState(null, "", `#/media-viewer/${$current_category.uuid}/${$active_media_index}`);

        await tick();
        
        resetMediaConfigs(true);
    }
    
    const handleGoBack = async () => {
        await $media_changes_manager.commitChanges($current_category.uuid);
        $categories_tree.updateCurrentCategory();
        media_changes_manager.set(new MediaChangesManager());

        pop();

        exitFullScreen();
    }

    const handleShowMediaGallery = e => {
        let { direction } = e.detail;

        if (direction !== "top" ) return;

        show_media_gallery = true;
    }

    const handleThumbnailClick = e => {
        const media_selected = e.detail;
        
        if (media_selected === undefined || media_selected === null) return;

        const media_index = $current_category.content.findIndex(media => media.uuid === media_selected.uuid);

        if (media_index === -1) return;

        active_media_index.set(media_index);
        saveActiveMediaToRoute();
        resetMediaConfigs(true);
        show_media_gallery = false;
    }

    const resetMediaConfigs = () => {
        media_zoom = 1;
        media_control_wrapper.style.transformOrigin = `center`;
        media_control_wrapper.style.transform = `scale(${media_zoom})`;
    }

    const resetComponentSettings = () => {
        // $random_media_navigation = false;
        // $skip_deleted_medias = false;

        active_media_index.set(0);
        previous_media_index.set(0);
    }

    const rejectMedia = () => {
        const current_media = $current_category.content[$active_media_index];
        let not_deleted_media_index = $active_media_index;

        if ($active_media_change !== media_change_types.DELETED) {
            $media_changes_manager.stageMediaDeletion(current_media);
            not_deleted_media_index = getNextNotDeletedMediaIndex($active_media_index, true);
        } else {
            $media_changes_manager.unstageMediaDeletion(current_media.uuid);
        }
        
        if (not_deleted_media_index !== $active_media_index) {
            active_media_index.set(not_deleted_media_index);
            history.replaceState(null, "", `#/media-viewer/${$current_category.uuid}/${$active_media_index}`);
        } else {
            active_media_change.set($media_changes_manager.getMediaChangeType(current_media.uuid));
        }

    }

    const saveActiveMediaToRoute = () => {
        history.replaceState(null, "", `#/media-viewer/${$current_category.uuid}/${$active_media_index}`);
    }

    const toggleRandomMediaNavigation = () => {
        $random_media_navigation = !$random_media_navigation;        
    }

</script>

<main id="libery-dungeon-media-viewer-mobile" style:position="relative" on:swipe={handleShowMediaGallery} use:swipe={{timeframe: 500, minSwipeDistance: 120}}>
    <header id="ldmvm-media-top-bar" class:adebug={false}>
        <button on:click={handleGoBack} id="ldmvm-mtb-exit-btn" type="button" title="exit media viewer">
            <svg viewBox="0 0 24 24">
                <path d="M16 22L8 12L16 2"/>
            </svg>
        </button>
        <button id="ldmvm-mtb-random-navigation-btn" on:click={toggleRandomMediaNavigation} class:random-nav-enabled={$random_media_navigation}>
            <svg viewBox="0 0 24 24" >
                <path d="M3.2 8A1.8 .8 0.0 0 1 20.8 8M2.8 4.4L3.2 8M3.2 16A1.8 0.8 0.0 0 0 20.8 16M20.6 20.4L20.8 16"/>
            </svg>
        </button>
        <h2 id="ldmvm-mtb-media-counter">
            {#if $current_category !== null}
                (<span>{$active_media_index + 1}</span> / <span>{$current_category.content.length}</span>)
            {/if}
        </h2>
    </header>
    {#if $current_category !== null}
        <div on:swipe={handleMediaSwipe} id="ldmvm-media-wrapper" use:swipe>
            <div id="media-control-wrapper"
                bind:this={media_control_wrapper}
                on:pinch={handleMediaZoom}
                use:pinch
            >
                {#if $current_category.content[$active_media_index]?.type === media_types.IMAGE}
                    <img class="mw-media-element-display" src="{getMediaUrl($current_category.FullPath, $current_category.content[$active_media_index].name, false, true)}" alt="displayed media">
                {:else}
                    <video class="mw-media-element-display" 
                        src="{getMediaUrl($current_category.FullPath, $current_category.content[$active_media_index].name)}" 
                        bind:this={video_element}
                        muted autoplay loop
                    >
                        <track kind="caption"/>
                    </video>
                {/if}
            </div>
        </div>
    {/if}
    {#if $current_category?.content.length > 0 && show_media_gallery}            
        <div id="ldmvm-media-gallery-wrapper" in:fly={{delay: 150, duration: 800, y: 200, opacity: 0}}>
            <MediasGallery on:thumbnail-click={handleThumbnailClick} on:gallery-close={() => show_media_gallery = false} all_available_medias={$current_category.content} category_path={$current_category.FullPath}/>
        </div>
    {/if}
    {#if video_element !== undefined && video_element !== null}
        <div id="mvm-mw-video-controller-wrapper">
            <VideoController video_element={video_element} auto_hide={false}/>
        </div>
    {/if}
</main>

<style>
    #libery-dungeon-media-viewer-mobile {
        height: 100dvh;
        overflow: hidden;
    }

    
    /*=============================================
    =            Top bar        =
    =============================================*/
    
    
        #ldmvm-media-top-bar {
            position: absolute;
            container-type: size;
            display: flex;
            width: 100vw;
            height: 10vh;
            background: transparent;
            top: 0;
            left: 0;
            justify-content: space-between;
            align-items: center;
            z-index: var(--z-index-t-1);
            padding: 0 var(--vspacing-2);
        }

        #ldmvm-mtb-exit-btn {
            width: 10cqw;
            height: 50cqh;
            background: transparent;
            border: none;
            outline: none;
            padding: 0;
        }

        #ldmvm-media-top-bar svg {
            width: 100%;
            height: 100%;
            overflow: visible;
        }

        #ldmvm-media-top-bar svg path {
            stroke: var(--grey-1);
            fill: none;
            stroke-width: 2;
            stroke-linecap: round;
            stroke-linejoin: round;
        }

        #ldmvm-mtb-random-navigation-btn {
            width: 10cqw;
            height: 50cqh;
            background: transparent;
            border: none;
            outline: none;
            padding: 0;
        }

        #ldmvm-mtb-random-navigation-btn.random-nav-enabled svg path {
            stroke: var(--success-3);
        }

        #ldmvm-mtb-media-counter {
            font-family: var(--font-read);
            line-height: .2;
        }

        #ldmvm-mtb-media-counter span {
            font-family: var(--font-decorative);
            color: var(--main-dark);
        }
    
    /*=====  End of Top bar  ======*/

    #ldmvm-media-wrapper {
        position: absolute;
        width: 100vw;
        max-height: 100dvh;
        top: 50%;
        left: 0;
        transform-origin: center;
        transform: translateY(-50%);
    }

    /* #media-control-wrapper video {
        pointer-events: none;
    } */

    
    /*=============================================
    =            Media Gallery            =
    =============================================*/
    
    #ldmvm-media-gallery-wrapper {
        position: absolute;
        width: 100vw;
        height: 60vh;
        bottom: 0;
        left: 0;
    }
    
    /*=====  End of Media Gallery  ======*/
    
    

    .mw-media-element-display {
        width: 100%;
        height: 100%;
        object-fit: contain;
    }

    #mvm-mw-video-controller-wrapper {
        position: absolute;
        width: max(13vw, 300px);
        left: 50%;
        bottom: 10%;
        transform: translateX(-50%);
    }
</style>