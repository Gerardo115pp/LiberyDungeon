<script>
    /*=============================================
    =            Imports            =
    =============================================*/
    
        import VideoController from "@components/VideoController/VideoController.svelte";
        import { MediaChangesManager, media_change_types } from "@models/WorkManagers";
        import MediaMovementsTool from "./sub-components/MediaMovementsTool/MediaMovementsTool.svelte";
        import { categories_tree, current_category } from "@stores/categories_tree";
        import global_hotkeys_manager from "@libs/LiberyHotkeys/libery_hotkeys";
        import { app_context_manager } from "@libs/AppContext/AppContextManager";
        import { getCategoryTree, category_cache } from "@models/Categories";
        import MediasGallery from "./sub-components/MediaGallery/MediasGallery.svelte";
        import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
        import { app_contexts } from "@libs/AppContext/app_contexts";
        import { hotkeys_sheet_visible, inCinemaMode, inDarkMode, layout_properties, navbar_hidden, toggleCinemaMode, toggleDarkMode } from "@stores/layout";
        import { onMount, onDestroy, tick } from "svelte";
        import { getMediaUrl } from "@libs/HttpRequests";
        import { replaceState } from "$app/navigation";
        import { media_types } from "@models/Medias";
        import { resetMediaViewerPageStore } from "./app_page_store";
        import { 
            active_media_index, 
            active_media_change,
            media_changes_manager,
            random_media_navigation,
            previous_media_index,
            skip_deleted_medias,
            auto_move_on,
            auto_move_category,
            media_viewer_hotkeys_context_name,
            automute_enabled
        } from "@stores/media_viewer";
        import { current_cluster, loadCluster } from "@stores/clusters";
        import MediaInformationPanel from "./sub-components/MediaInformation/MediaInformationPanel.svelte";
        import { goto } from "$app/navigation";
        import { page } from "$app/stores";
    import { current_user_identity } from "@stores/user";
    
    /*=====  End of Imports  ======*/
     
    /*=============================================
    =            App Context            =
    =============================================*/
    
        const app_context = app_context_manager.CurrentContext;
        const app_page_name = "media-viewer";

        if (app_context !== app_contexts.DUNGEONS) {
            app_context_manager.setAppContext(app_contexts.DUNGEONS, app_page_name, onComponentExit);
        } else {
            app_context_manager.addOnContextExit(app_page_name, onComponentExit)
        }
    
    /*=====  End of App Context  ======*/
    
    /*=============================================
    =            properties            =
    =============================================*/
    
        /*----------  Exports  ----------*/

        /**
         * A category id passed through the url
         * @type {string}
         */
        export let url_category_id;

        /**
         * The index of the media to display
         * @type {number}
         */
        export let url_media_index;

        /*----------  References  ----------*/

            /**
             * @type {HTMLVideoElement} that displays video medias
            */
            let video_element;

        /*=============================================
        =            State            =
        =============================================*/

            /** 
             * Whether the value for active_media_index has been determined. if this is false, any value active_media_index cannot
             * be considered valid.
             * @type {boolean}
             * @default false
             */
            let active_media_index_determined = false;
             

            /** @type {number} on each movement event, the medias will move by media_movement_factor * media_height */
            let media_movement_factor = 0.05

            let base_media_top_value = "50%";
            /** @type {number} the medias can only move an amount equal to media_movement_threshold * media_height*/
            let media_movement_threshold = 2;
            let media_zoom_factor = 0.1;
            let media_zoom = 1;

            /** @type {boolean} whether to show the media movement manager */
            let show_media_movement_manager = false;

            /** @type {boolean} whether to show the media gallery */
            let show_media_gallery = false;

            /** @type {boolean} whether to show the media information panel */
            let show_media_information_panel = false;

            /** @type {string} the name of the hotkeys context */
            const hotkeys_context_name = "media_viewer";

            /**  @type {boolean} Whether the media viewer is currently closing. 
             */
            let media_viewer_closing = false;

            /**
             * Whether to enable Cinema mode which disables user media transformations and enters a fullscreen mode also scaling the media to fit the entire available space.
             * @type {boolean} 
             */
            let cinema_mode = false;
        
        /*=====  End of State  ======*/
        
        /*----------  Unsubscriber  ----------*/
        
            let active_media_index_unsubscriber = () => {};
    
    /*=====  End of properties  ======*/
    
    onMount(async () => {
        if (!ensureCluster()) return;

        toggleDarkMode(true);

        if (!layout_properties.IS_MOBILE) {
            defineDesktopKeybinds();
        }

        if ($current_category === null) {
            let new_categories_tree = await getCategoryTree(url_category_id, current_category);
            categories_tree.set(new_categories_tree);
        }

        let determined_index = await determineActiveMediaIndexInitialValue();
        active_media_index.set(determined_index);
        active_media_index_determined = true;

        console.log("Opening media viewer")

        // Sets the active media change every time the active media index changes
        active_media_index_unsubscriber = active_media_index.subscribe(new_index => {
            if ($current_category === null) return;

            if (new_index >= $current_category.content.length || new_index < 0) {
                return;
            }

            let media_change = $media_changes_manager.getMediaChangeType($current_category.content[new_index].uuid);
            active_media_change.set(media_change);
        })
    });

    onDestroy(() => {
        onComponentExit();
        
        if (inCinemaMode()) {
            toggleMediaViewerCinemaMode();
        }
        
        toggleDarkMode(false);
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /*=============================================
        =            Keybinding            =
        =============================================*/
        
            const defineDesktopKeybinds = () => {
                if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    const hotkeys_context = new HotkeysContext();

                    hotkeys_context.register(["a", "d"], handleMediaNavigation, {
                        description: "<navigate>Navigate through the medias, A for previous, D for next", 
                    });

                    hotkeys_context.register(["w", "s"], handleMoveImageUpDown, {
                        description: "<media_modification>Move the image up or down, W for up, S for down."
                    });

                    hotkeys_context.register("r", e => {e.preventDefault(); random_media_navigation.set(!$random_media_navigation)}, {
                        description: "<navigation>Toggle random media navigation."
                    });

                    hotkeys_context.register("1", e => {e.preventDefault(); changeDisplayedMedia(0)}, {
                        description: "<navigation>Go to the first media."
                    });

                    hotkeys_context.register(["shift+a", "shift+d"], handleMediaZoom, {
                        description: "<media_modification>Zoom in or out the media, Shift+A for zoom out, Shift+D for zoom in."
                    });

                    hotkeys_context.register("alt+shift+d", e => {e.preventDefault(); resetMediaConfigs()}, {
                        description: "<media_modification>Reset the media zoom and position."
                    });

                    hotkeys_context.register(["q", "esc"], handleGoBack, {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Go back to the media explorer.`
                    });

                    hotkeys_context.register("i", handleShowMediaInformationPanel, {
                        description: "<media_modification>Toggle the media information panel."
                    });

                    
                    if ($current_user_identity.canPublicContentAlter()) {
                        hotkeys_context.register("t", e => show_media_movement_manager = !show_media_movement_manager, {
                            description: "<media_modification>Toggle the media movement manager."
                        });

                        hotkeys_context.register("e", rejectMedia, {
                            description: "<media_filtration>Reject/Delete the current media."
                        });

                        hotkeys_context.register("shift+alt+e", e => skip_deleted_medias.set(!$skip_deleted_medias), {
                            description: "<navigation>Toggle skipping deleted medias."
                        });
                    }



                    hotkeys_context.register("g", e => show_media_gallery = !show_media_gallery, {
                        description: "<navigation>Toggle the media gallery."
                    });

                    hotkeys_context.register("n", e => clearActiveMediaChanges(), {
                        description: "<media_modification>Clear all media modifiers(e.g Zoom, Vertical position).",
                    });

                    hotkeys_context.register(["l o"], handleDarkModeToggle, {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Lights on/off.`,
                    });

                    hotkeys_context.register(["f"], handleCinemaModeToggle, {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Cinema mode on/off. when on, the browser will enter fullscreen mode and the media will be scaled to fit the available space, all media transformations will be disabled.`,
                    });

                    hotkeys_context.register("?", e => hotkeys_sheet_visible.set(!$hotkeys_sheet_visible), {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Toggle the hotkeys cheatsheet.`,
                    });



                    global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);

                    media_viewer_hotkeys_context_name.set(hotkeys_context_name);
                }

                global_hotkeys_manager.loadContext(hotkeys_context_name);            
            }
            
            const clearActiveMediaChanges = () => {
                if ($active_media_change === media_change_types.NORMAL) return;

                const current_media = $current_category.content[$active_media_index];

                $media_changes_manager.clearMediaChanges(current_media.uuid);

                active_media_change.set(media_change_types.NORMAL);
            }

            /**
             * Changes the displayed media based on a new given index. 
             * @param {number} new_index
             */
            const changeDisplayedMedia = new_index => {
                if (new_index < 0 || new_index >= $current_category.content.length) {
                    return;
                }

                active_media_index.set(new_index);
                replaceState(`/media-viewer/${$current_category.uuid}/${$active_media_index}`);
            }

            /**
             * Handles the navigation between medias. If the random media navigation is enabled
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMediaNavigation = async (key_event, hotkey) => {
                let key_combo = hotkey.key_combo.toLowerCase();
                
                if ($random_media_navigation) {
                    return handleRandomMediaNavigation(key_event, key_combo);
                }


                let new_index = $active_media_index;

                new_index += key_combo === "a" ? -1 : 1;
                new_index = Math.max(0, Math.min(new_index, $current_category.content.length - 1));

                if(media_change_types.DELETED === $media_changes_manager.getMediaChangeType($current_category.content[new_index].uuid) && $skip_deleted_medias) {
                    let not_deleted_new_index = getNextNotDeletedMediaIndex(new_index, key_combo !== "a");
                    new_index = (not_deleted_new_index === new_index) ? $active_media_index : not_deleted_new_index;
                }

                automoveMedia();

                changeDisplayedMedia(new_index);

                await tick();
                
                resetMediaConfigs(true);
            }

            /**
             * Handles the random media navigation. If the key_combo is "a", then the previous media index will be used as the new index.
             * Called only from handleMediaNavigation
             * @param {KeyboardEvent} key_event
             * @param {string} key_combo
             */
            const handleRandomMediaNavigation = async (key_event, key_combo) => {
                let new_index = $active_media_index;

                

                new_index = Math.floor(Math.random() * $current_category.content.length);

                while (new_index === $active_media_index) {
                    new_index = Math.floor(Math.random() * $current_category.content.length);
                }

                if (key_combo === "a") {
                    new_index = $previous_media_index;
                }

                previous_media_index.set($active_media_index);

                active_media_index.set(new_index);
                replaceState(`#/media-viewer/${$current_category.uuid}/${$active_media_index}`);

                await tick();
                
                resetMediaConfigs(true);
            }

            /**
             * Handles the movement of the image up or down. The movement is based on the media height and the media_movement_factor.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMoveImageUpDown = (key_event, hotkey) => {
                if (cinema_mode) return;

                let key_combo = hotkey.key_combo.toLowerCase();

                let media_wrapper = document.getElementById("media-wrapper");

                let media_wrapper_style = window.getComputedStyle(media_wrapper);

                let new_media_top_position = parseInt(media_wrapper_style.top);
                let movement_amount = media_wrapper.clientHeight * media_movement_factor;
                movement_amount = media_zoom < 1 ? movement_amount : movement_amount * media_zoom; // make movement faster if zoomed in
                let movement_direction = key_combo === "w" ? -1 : 1;

                new_media_top_position += movement_amount * movement_direction;


                if (Math.abs(new_media_top_position) <= media_wrapper.clientHeight * media_movement_threshold) {
                    media_wrapper.style.top = `${new_media_top_position}px`;
                }
            }

            /**
             * Handles the zooming in or out of the media. The zooming is based on the media_zoom_factor.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMediaZoom = (key_event, hotkey) => {
                if (cinema_mode) return;

                let key_combo = hotkey.key_combo.toLowerCase();

                let media_element = document.querySelector(".mw-media-element-display");


                if (media_element === null) {
                    return;
                }

                let new_media_zoom = media_zoom;

                new_media_zoom += key_combo === "shift+a" ? -media_zoom_factor : media_zoom_factor;

                if (new_media_zoom >= 0.1 && new_media_zoom <= 3) {
                    media_zoom = new_media_zoom;
                    media_element.style.transform = `scale(${media_zoom})`;
                }
            }

            const handleShowMediaInformationPanel = () => {
                show_media_information_panel = !show_media_information_panel;
            }

            const handleGoBack = async () => {
                if (media_viewer_closing) return;
                media_viewer_closing = true;

                global_hotkeys_manager.unregisterCurrentContext(); // Prevents users from using hotkeys while the page is closing

                // change the active media index so that onComponentExit will save the correct updated index(saves the current value of active_media_index)
                let new_active_media_index = resolveNewMediaIndex($media_changes_manager, $active_media_index, $current_category.content);
                
                await $media_changes_manager.commitChanges($current_category.uuid, $current_category.ClusterUUID);
                await $categories_tree.updateCurrentCategory();
                media_changes_manager.set(new MediaChangesManager());

                // wait to set the active media index to the new value after async operations so the user doens't see weird image switching
                active_media_index.set(new_active_media_index);
            
                history.back();
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
                    replaceState(`#/media-viewer/${$current_category.uuid}/${$active_media_index}`);
                } else {
                    active_media_change.set($media_changes_manager.getMediaChangeType(current_media.uuid));
                }
            }

            /**
             * Returns the next media index that is not marked as deleted. optionally, it can skip the medias that are marked as moved if
             * the skip_moved parameter is set to true.
             * @param {number} from_index
             * @param {boolean} forward
             * @param {boolean} skip_moved
             * @returns {number}
             */
            const getNextNotDeletedMediaIndex = (from_index, forward, skip_moved=false) => {
                let next_index = from_index;
                const media_count = $current_category.content.length;
                let media_change;
                let media_uuid;

                // if the current media is the last media and forward was requested, or if the current media is the first media and backward was requested, 
                // then STAY IN THE SAME INDEX
                if ((from_index === (media_count - 1) && forward) || (from_index === 0 && !forward)) return from_index;

                let stop_loop = false;
                let is_out_of_bounds = false;

                do {
                    next_index += forward ? 1 : -1;
                    
                    is_out_of_bounds = next_index < 0 || next_index >= (media_count - 1);
                    if (is_out_of_bounds) {
                        stop_loop = true;
                        break;
                    }

                    media_uuid = $current_category.content[next_index]?.uuid;
                    media_change = $media_changes_manager.getMediaChangeType(media_uuid);
                    
                    stop_loop = media_change !== media_change_types.DELETED && (!skip_moved || media_change !== media_change_types.MOVED);
                } while(!stop_loop);

                next_index = Math.max(0, Math.min(next_index, media_count - 1));
                media_uuid = $current_category.content[next_index]?.uuid;
                media_change = $media_changes_manager.getMediaChangeType(media_uuid);

                return media_change === media_change_types.DELETED || (media_change === media_change_types.MOVED && skip_moved) ? from_index : next_index;
            }

            /**
             * Hanldes the toggling of the dark mode
             */
            const handleDarkModeToggle = () => {
                toggleDarkMode(false); // parameter `force_enable` instead of toggling it turns it on despite its current state. so we set force_enable to false
            }

            /**
             * Handles the toggling of the cinema mode. 
             */
            const handleCinemaModeToggle = () => {
                toggleMediaViewerCinemaMode();
            }
        
        /*=====  End of Keybinding  ======*/

        const automoveMedia = () => {
            if ($auto_move_on === false) return;

            let current_media = $current_category.content[$active_media_index];

            let current_media_change = $media_changes_manager.getMediaChangeType(current_media.uuid);

            if (current_media_change !== media_change_types.NORMAL) return;

            $media_changes_manager.stageMediaMove(current_media, $auto_move_category);
        }

        /**
         * Determines what active_media_index should be set initially. Either the one from a url parameter, the one from the cache.
         * @requires url_category_id - to load the cached index from the cache
         * @requires url_media_index - if this is valid, then the cached index will be ignored
         * @requires $current_category
         * @requires category_cache
         * @returns {number}
         */
        const determineActiveMediaIndexInitialValue = async () => {
            let determined_index = 0;
            let url_holds_index = url_media_index != null && url_media_index !== "";

            if (url_holds_index) {
                determined_index = parseInt(url_media_index);
                url_holds_index = !isNaN(determined_index) && determined_index >= 0 && determined_index < $current_category.content.length;
            } 
            
            if (!url_holds_index) {
                // Attempts to get a cached media index for the current_category, and if the result is different from the current active media index, then it sets the active media index to the cached one
                // the url media_index params has priority over the cached index, so if params.media_index is not null, then no cached index will be used
                let cached_category_index = await category_cache.getCategoryIndex(url_category_id);
                if (cached_category_index !== determined_index) {
                    cached_category_index = Math.max(0, Math.min(cached_category_index, $current_category.content.length - 1));
                    determined_index = cached_category_index;
                }
            }
        

            return determined_index;
        }

        const ensureCluster = () => {
            if ($current_cluster !== null) {
                return true;
            }

            const able_to_load_cluster = loadCluster();

            if (!able_to_load_cluster) {
                goto("/");
            }

            return able_to_load_cluster;
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

        function onComponentExit() {        
            category_cache.addCategoryIndex(url_category_id, $active_media_index);

            resetComponentSettings();

            active_media_index_unsubscriber();

            global_hotkeys_manager.dropContext(hotkeys_context_name);

            resetMediaViewerPageStore();
        }

        const resetComponentSettings = () => {
            active_media_index.set(0);
            auto_move_on.set(false);
            auto_move_category.set(null);

            resetMediaConfigs();
        }

        const resetMediaConfigs = local_reset => {
            let media_wrapper = document.getElementById("media-wrapper");
            let media_element = document.querySelector(".mw-media-element-display");

            if (media_wrapper === null || media_element === null) return;

            media_wrapper.style.top = base_media_top_value;
            media_element.style.transform = `scale(1)`;

            if (local_reset) {
                media_zoom = 1;
            }
        }

        /**
         * Recalculates the new active media index subtracting from it the amount of medias that are removed from the category(delete and moved medias) and are
         * before the current active media index. If the current active media index is one of the removed or moved medias, it will first update this index until it finds a valid one.
         * if the amount of removed medias is greater or equal to total medias, then the new active media index will be 0.
         * @param {MediaChangesManager} changes_manager
         * @param {number} current_active_index
         * @param {import('@models/Medias').Media[]} all_medias
         * @returns {number}
         */
        const resolveNewMediaIndex = (changes_manager, current_active_index, all_medias) => {
            let new_index = current_active_index;
            let removed_medias = changes_manager.DeletedMediasMap;
            let moved_medias = changes_manager.MovedMediasMap;

            let is_active_media_removed = changes_manager.getMediaChangeType(all_medias[current_active_index].uuid) !== media_change_types.NORMAL;

            if (is_active_media_removed) {
                new_index = getNextNotDeletedMediaIndex(current_active_index, true, true);
                if (new_index === current_active_index) {
                    return current_active_index;
                }
            }

            for (let h = 0; h < current_active_index; h++) {
                let media_uuid = all_medias[h].uuid;
                let media_change = changes_manager.getMediaChangeType(media_uuid);

                if (media_change !== media_change_types.NORMAL) {
                    new_index--;
                }
            }

            return new_index;
        }

        const saveActiveMediaToRoute = () => {
            replaceState(`#/media-viewer/${$current_category.uuid}/${$active_media_index}`);
        }

        const screenshotVideo = () => {
            if (video_element == null) {
                return;
            };

            let canvas = document.createElement("canvas");
         
            canvas.width = video_element.videoWidth;
            canvas.height = video_element.videoHeight;

            let ctx = canvas.getContext("2d");

            ctx.drawImage(video_element, 0, 0, canvas.width, canvas.height);

            let dataURL = canvas.toDataURL("image/png");

            // download the image as a file just for testing purposes
            let a = document.createElement("a");
            a.href = dataURL;
            a.download = "screenshot.png";
            a.click();
        }

        /**
         * Enters or exists fullscreen mode.
         * @param {boolean} enter
         */
        const toggleFullscreen = enter => {
            if (!document.fullscreenEnabled) return; 


            if (enter) {
                if (document.fullscreenElement !== null) return;
                console.log("Entering fullscreen mode")
                document.documentElement.requestFullscreen({navigationUI: "hide"});
            } else {
                if (document.fullscreenElement === null) return;

                document.exitFullscreen();
            }
        }

        const toggleMediaViewerCinemaMode = () => {
            toggleCinemaMode();

            cinema_mode = inCinemaMode();

            navbar_hidden.set(cinema_mode);

            toggleFullscreen(cinema_mode);

            if (!inDarkMode()) {
                toggleDarkMode(true);
            }

            resetMediaConfigs(true);
        }
    
    /*=====  End of Methods  ======*/

</script>

<main id="libery-dungeon-media-viewer"
    style:position="relative"
    class:cinema-mode={cinema_mode}
>
    {#if $current_category !== null && active_media_index_determined && $current_category.content[$active_media_index] !== undefined}
        <div id="media-wrapper">
            {#if $current_category.content[$active_media_index].type === media_types.IMAGE}
                <img class="mw-media-element-display" src="{getMediaUrl($current_category.FullPath, $current_category.content[$active_media_index].name)}" alt="displayed media">
            {:else}
                <video 
                    class="mw-media-element-display"
                    bind:this={video_element}  
                    src="{$current_category.content[$active_media_index].Url}" 
                    muted={$automute_enabled}
                    autoplay 
                    loop
                >
                    <track kind="caption"/>
                </video>
            {/if}
        </div>
        {#if video_element !== undefined && video_element !== null}
            <div id="mw-video-controller-wrapper">
                <VideoController 
                    video_element={video_element} 
                    media_uuid={$current_category.content[$active_media_index].uuid}
                    auto_hide
                    on:screenshot-video={screenshotVideo}
                />
            </div>
        {/if}
        <div id="ldmv-category-changes-manager">
            {#if $current_user_identity.canPublicContentAlter()}
                <MediaMovementsTool bind:is_component_visible={show_media_movement_manager}/>
            {/if}
        </div>
        <div id="ldmv-media-information-panel-wrapper">
            {#if show_media_information_panel}
                <MediaInformationPanel 
                    current_category_information={$current_category} 
                    current_cluster_information={$current_cluster} 
                    current_media_information={$current_category.content[$active_media_index]}
                />
            {/if}
        </div>
        {#if $current_category?.content.length > 0 && show_media_gallery}            
            <div id="ldmv-media-gallery-wrapper">
                <MediasGallery on:thumbnail-click={handleThumbnailClick} all_available_medias={$current_category.content} category_path={$current_category.FullPath}/>
            </div>
        {/if}
    {/if}
</main>

<style>
    :global(#libery-dungeon-content:has(main#libery-dungeon-media-viewer) ) {
        margin-top: 0 !important;
    }

    #libery-dungeon-media-viewer {
        height: 100vh;
        overflow: hidden;
    }

    #media-wrapper {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        transition: all .28s ease-in-out;
    }
    
    .mw-media-element-display {
        /* max-width: 90vw; */
        max-height: 98vh;
        transform-origin: center center;
        transition: all .28s ease-in-out;
    }

    #mw-video-controller-wrapper {
        position: absolute;
        width: max(13vw, 400px);
        left: 50%;
        bottom: 10%;
        transform: translateX(-50%);
    }

    #ldmv-category-changes-manager {
        position: absolute;
        top: var(--navbar-height);
        right: 0;
        width: 20vw;
        height: calc(100vh - var(--navbar-height));
        z-index: var(--z-index-t-1);
    }

    #ldmv-media-information-panel-wrapper {
        position: absolute;
        container-type: size;
        top: calc(var(--navbar-height) + var(--vspacing-2));
        width: clamp(300px, 45em, 40vw);
        height: calc(100dvh - calc(var(--navbar-height) * 1.5));
        font-size: var(--font-size-1);
        left: var(--vspacing-3);
        z-index: var(--z-index-t-1);
    }

    #ldmv-media-gallery-wrapper {
        position: absolute;
        top: 30%;
        left: 15%;
        width: 35vw;
        height: 50vh;
        z-index: var(--z-index-t-2);
    }

    
    /*=============================================
    =            Cinema mode            =
    =============================================*/
    
        main#libery-dungeon-media-viewer.cinema-mode {

            & #media-wrapper {
                position: static;
                display: grid;
                place-items: center;
                top: 0;
                transform: translate(0, 0);
                height: 100%;
                width: 100%;
            }                

            & video.mw-media-element-display, img.mw-media-element-display {
                width: 100%;
                height: auto;
                object-fit: contain;
                object-position: center;
                max-height: 100dvh;
            }
        }
    
    /*=====  End of Cinema mode  ======*/
    
    
</style>