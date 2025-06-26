<script>
    /*=============================================
    =            Imports            =
    =============================================*/
    
        import viewport from '@components/viewport_actions/useViewportActions';
        import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
        import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
        import { categories_tree, current_category } from '@stores/categories_tree';
        import { current_cluster } from '@stores/clusters';
        import { hotkeys_sheet_visible, layout_properties } from "@stores/layout";
        import { HOTKEYS_HIDDEN_GROUP, HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import { createEventDispatcher, onDestroy, onMount, tick } from 'svelte';
        import { me_gallery_changes_manager, me_gallery_yanked_medias, me_renaming_focused_media, meg_intersection_observer_event_names } from './me_gallery_state';
        import MeGalleryDisplayItem from './MEGalleryDisplayItem.svelte';
        import GridLoader from '@components/UI/Loaders/GridLoader.svelte';
        import CoverSlide from '@components/Animations/HoverEffects/CoverSlide.svelte';
        import { browser } from '$app/environment';
        import { media_change_types, MediaChangesEmitter } from '@models/WorkManagers';
        import { category_cache, InnerCategory } from '@models/Categories';
        import { LabeledError } from '@libs/LiberyFeedback/lf_models';
        import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
        import { confirmPlatformMessage, emitPlatformMessage } from '@libs/LiberyFeedback/lf_utils';
        import { pushState, replaceState } from '$app/navigation';
        import { page } from '$app/stores';
        import { Media, OrderedMedia } from '@models/Medias';
        import SequenceCreationTool from './SequenceCreationTool.svelte';
        import { ui_pandasworld_tag_references } from '@app/common/ui_references/dungeon_tags_references';
        import { ui_core_dungeon_references } from '@app/common/ui_references/core_ui_references';
        import { common_action_groups } from '@app/common/keybinds/CommonActionsName';
        import { SearchResultsWrapper } from '@common/keybinds/CommonActionWrappers';
        import { CursorMovementWASD } from '@app/common/keybinds/CursorMovement';
    
    /*=====  End of Imports  ======*/
    
    /*=============================================
    =            Properties            =
    =============================================*/
        
        /*=============================================
        =            Hotkeys            =
        =============================================*/

            const global_hotkeys_manager = getHotkeysManager();
        
            const hotkeys_context_name = "media_explorer_gallery"; 
        

        
        /*=====  End of Hotkeys  ======*/

        /** 
         * The all medias that are available to be displayed in the gallery(f.e. all the medias in a category). These are not 
         * are not necessarily the same as the ones displayed. but the list as provided by the MediaExplorer(parent component).
         * @type {import('@models/Medias').Media[]}
         */
        export let media_items = [];
        
        /**
         * A list of OrderedMedia objects with their respective positions in the media_items array which is ordered by the server either by name or if the media set represents some kind
         * of series(like a comic book, show, etc. e.g dr_house_s01e01.mp4) then the server attempts to infer the order of the 'episodes' based on the found indexers. long
         * story short, the order in which the media_items are provided is none trivial, it must be preserved.
         * @type {import('@models/Medias').OrderedMedia[]}
         */
        let ordered_medias = [];
        $: ordered_medias = media_items.map((mi, order_position) => new OrderedMedia(mi, order_position));

        /**
         * an html class that is added to all the media items in the gallery.
         * @type {string}
         */
        const media_item_html_class = "meg-gallery-grid-cell";
       
        
        /*----------  Intersection observer  ----------*/
        
            /**
             * The intersection observer passed to the gallery items.
             * @type {IntersectionObserver | null}
             */
            let gallery_item_intersection_observer = null;
        
            /**
             * Last media items orders that have entered the viewport.
             * @type {Set<number>}
             */
            let last_orders_unveiled = new Set();
        /*----------  State  ----------*/
        
            /**
             * The medias that ARE been displayed in the gallery. 
             * @type {import('@models/Medias').OrderedMedia[]}
             */
            export let active_medias = [];

            /**
             * Whether there are more proceeding medias to load.
             * Used to disable the ContentEndWatchdog.
             * @type {boolean}
             */
            let has_proceeding_medias = true;

            /** 
             * If true, enables the media explorer gallery hotkeys_context
             * @type {boolean}
             */
            export let enable_gallery_hotkeys = false;
            $: if (enable_gallery_hotkeys) {
                defineComponentHotkeys();
            }

            /**
             * Whether to magnify/zoom the keyboard focused media item.
             * @type {boolean}
             * @default false
             */
            let magnify_focused_media = false;

            /**
             * Whether to show the media titles in the gallery items.
             */
            let show_media_titles_mode = false;

            /**
             * Whether the gallery should attempt to recover the scroll position
             * and focus index from a previous session.
             * @type {boolean}
             */
            export let recovering_gallery_state = true;

            /**
             * Whether to enable the sequence creation tool.
             * @type {boolean}
             */
            let enable_sequence_creation_tool = false;

            /**
             * The last media that was added to the active_medias array. If it was added
             * by prepending, then it should be the first media on the batch that was added.
             * If it was added by appending, then it should be the last media on the batch that was added.
             * @type {import('@models/Medias').OrderedMedia | undefined}
             */
            let last_media_added_to_active_medias = undefined;

            
            /*----------  Performance regulation  ----------*/
            
                /**
                 * Whether to enable performance mode for the gallery. Completely controlled by
                 * it's managing function based on fuzzy logic.
                 * @type {boolean}
                 */
                let enable_gallery_performance_mode = false;
            
                /**
                 * Whether to enable heavy rendering of the medias in the gallery. This will force the media item in the gallery to render in higher quality and
                 * if they are videos, those will be loaded instead of their thumbnails(which is the default behavior) they will however only load when they are on 
                 * the viewport.
                 * @type {boolean}
                 */
                export let enable_gallery_heavy_rendering = false;

                /**
                 * Whether the process of regulating the active media load is running.
                 * @type {boolean}
                 */
                let regulating_active_medias_load = false;

            /*----------  Navigation  ----------*/

                /**
                 * The grid navigation wrapper of the Media Explorer Gallery.
                 * @type {CursorMovementWASD | null}
                 */
                let the_grid_navigation_wrapper = null;
            
                /** 
                 * Media selected index. This index corresponds to the medias position in the ordered_medias array. not just
                 * in the active_medias array.
                 * @type {number}
                 */
                let media_focus_index = 0;
            
                /**
                 * A list of functions that are called when the_grid_navigation_wrapper's Grid 
                 * changes. Managed by waitForGridSync function.
                 * @type {Array<() => void | Promise<void>>}
                 */
                let grid_navigation_change_callbacks = [];
            
            /*----------  Selections  ----------*/
            
                /**
                 * Whether focused media items should be auto selected(usually for media yanking). Enabled on keydown for media select
                 * key(originally space) and disabled on key up.
                 * @type {boolean}
                 */
                let auto_select_focused_media = false;

                /**
                 * Whether focused media should automatically staged for deletion.
                 * @type {boolean}
                 */
                let auto_stage_delete_focused_media = false;           

                /**
                 * Whether the auto select mode is enabled. 
                 * @type {boolean}
                 */
                let auto_select_mode_enabled = false;

                /**
                 * Whether the auto select mode should add or remove elements from selection.
                 * @type {boolean}
                 * @default true
                 */
                let auto_select_mode_adds = true;
                
            /*----------  Search media content  ----------*/
            
                /**
                 * The wrapper of search functionality for media content.
                 * @type {SearchResultsWrapper<import('@models/Medias').OrderedMedia> | null}
                 */
                let the_media_content_search_wrapper = null;

                /** 
                 * The query to search for in the media content.
                 * @type {string}
                 */
                let media_content_search_query = "";

                /**
                 * Whether the media search is active.
                 * @type {boolean}
                 */
                let capturing_media_content_search = false;
                
                /**
                 * The search results for the last query the
                 * user entered.
                 * @type {Set<string> | null}
                 */
                let media_content_search_results_lookup = null;
        
        /*----------  Masonry  ----------*/

            /**
             * Whether to use a masonry layout for the gallery.
             */
            export let use_masonry = false;
            $: handleMasonryLayoutChange(), use_masonry;
        

        let dispatch = createEventDispatcher();

        /**
         * Enables or disables(manually) the debug mode.
         * @type {boolean}
         */
        let debug_mode = true;

    /*=====  End of Properties  ======*/

    onMount(async () => {

        await defineGalleryState();
        defineComponentHotkeys();

        if (debug_mode) {
            debugMEG__attachDebugMethods();
        }
    });

    onDestroy(() => {
        if (!browser) return;
        resetGalleryState();
        console.debug("In MediaExplorerGallery.onDestroy: Gallery component was destroyed.");
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /*=============================================
        =            Hotkeys            =
        =============================================*/
            // Only call directly defineComponentHotkeys. All the other methods in this section should only be called by the global hotkeys manager.
        
            /**
             * Defines the hotkeys to interact with the media explorer gallery.
             */        
            const defineComponentHotkeys = () => {
                if (global_hotkeys_manager == null) {
                    console.error("In MediaExplorerGallery.defineComponentHotkeys: No hotkeys manager available.");
                    return;
                }
                
                if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    const hotkeys_context = new HotkeysContext();



                    hotkeys_context.register(["q"], handleGalleryExit, {
                        description: "<navigation>Sets focus on the category section without closing the gallery.",
                    });

                    hotkeys_context.register(["shift+g"], handleGalleryClose, {
                        description: "<navigation>Close the gallery and return focus to the category section.",
                    });

                    hotkeys_context.register(["shift+e"], handleOpenMedia, {
                        description: "<navigation>Open the media viewer on the focused media",
                    });

                    hotkeys_context.register(["e", "del"], handleMediaStageDeletionMode, {
                        description: "<content>Selects the focused media to be deleted on gallery close(hitting the g key).",
                        mode: "keydown"
                    });

                    hotkeys_context.register(["e", "del"], handleMediaStageDeletionMode, {
                        description: "<content>Selects the focused media to be deleted on gallery close(hitting the g key).",
                        mode: "keyup"
                    });

                    hotkeys_context.register(["space"], handleMediaSelectAddMode, {
                        description: "<content>Selects the focused media.",
                        mode: "keydown",
                    });

                    hotkeys_context.register(["space"], handleMediaSelectAddMode, {
                        description: `<${common_action_groups.HIDDEN}> hidden`,
                        mode: "keyup"
                    });

                    hotkeys_context.register(["y y"], handleSetYankedMedias, {
                        description: "<content>Yanks the selected medias. previous yanked medias returned to normal state.",
                    });

                    hotkeys_context.register(["y a"], handleAppendToYankedMedias, {
                        description: "<content>Appends the selected medias to the yanked medias list.",
                    });

                    hotkeys_context.register(["p"], handleYankPaste, {
                        description: "<content>Pastes the yanked medias if they were not yanked from the current category.",
                    });

                    hotkeys_context.register(["alt+n"], handleGalleryReset, {
                        description: "<content>Deselects all selected medias and restores deleted medias that have not been committed(they are committed when you close the gallery).",
                    });

                    hotkeys_context.register(["c s"], handleEnableSequenceCreationTool, {
                        description: "<tools>Enable the sequence creation tool.",
                    });

                    hotkeys_context.register(["c c"], handleRenameCurrentMediaHotkey, {
                        description: `${common_action_groups.CONTENT}Allows you to rename the focused ${ui_core_dungeon_references.MEDIA.EntityName}`,
                        await_execution: false,
                        mode: "keyup"
                    })

                    hotkeys_context.register(["m"], handleMagnifyFocusedMedia, {
                        description: "<behavior>Toggle the magnify focused media mode. when enabled, the focused media item will be scale up and be contained so the entier media can be seen.",
                    });

                    hotkeys_context.register(["v"], handleToggleHeavyRendering, {
                        description: "<behavior>Toggle the heavy rendering of the medias in the gallery. When enabled, the medias will be loaded in higher quality and videos will be loaded instead of their thumbnails.",
                    });

                    hotkeys_context.register(["t"], handleShowTitlesMode, {
                        description: `<${HOTKEYS_GENERAL_GROUP}> Turns media titles on and off.`,
                    });

                    hotkeys_context.register(["?"], () => hotkeys_sheet_visible.set(!$hotkeys_sheet_visible), {
                        description: `<${HOTKEYS_GENERAL_GROUP}> Toggle the hotkeys cheat sheet.`,
                    });

                    setSearchResultsWrapper(hotkeys_context);

                    // hotkeys_context.register(["w", "a", "s", "d"], handleGalleryGridMovement, {
                    //     description: `<navigation> Move the focus on the gallery grid.`,
                    // });

                    // hotkeys_context.register(["shift+s"], handleMovetoLastActiveMedia, {
                    //     description: `<navigation> Move the focus to the last loaded media in the gallery.`,
                    // });

                    // hotkeys_context.register(["shift+w"], handleMovetoFirstActiveMedia, {
                    //     description: `<navigation> Move the focus to the first loaded media in the gallery.`,
                    // });

                    hotkeys_context.register(["\\d g"], handleJumpToMediaOrder, {
                        description: "<navigate>Jump to the media in the \\d position.",
                    });

                    setGridNavigationWrapper(hotkeys_context);

                    global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
                }

                if (enable_gallery_hotkeys && global_hotkeys_manager.ContextName !== hotkeys_context_name) {
                    global_hotkeys_manager.loadContext(hotkeys_context_name);
                };
            }

            /**
             * Appends selected medias to the yanked medias list.
             */
            const handleAppendToYankedMedias = () => {
                if ($me_gallery_changes_manager === null || !($me_gallery_changes_manager instanceof MediaChangesEmitter)) {
                    console.warn("No changes manager available to append the selected medias to the yanked medias list.");
                    return;
                }

                let yanked_medias = $me_gallery_yanked_medias;
                const original_yanked_medias_count = yanked_medias.length;
                let selected_medias = $me_gallery_changes_manager.MovedMedias;

                $me_gallery_changes_manager.clearAllMoveChanges();

                // FIXME: This 
                // @ts-ignore - this could be an issue but it's complex so, i will Fix it later.
                me_gallery_yanked_medias.set([
                    ...yanked_medias,
                    ...selected_medias
                ]);

                const total_yanked_medias_count = $me_gallery_yanked_medias.length;

                emitPlatformMessage(`Appended ${selected_medias.length} medias to existing ${original_yanked_medias_count} yanked medias. There are now ${total_yanked_medias_count} yanked medias.`);
            }
            
            /**
             * Enables the media sequence creation tool.
             */
            const handleEnableSequenceCreationTool = () => {
                enable_sequence_creation_tool = true;
            }

            /**
             * Enables the focused media renaming state.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleRenameCurrentMediaHotkey = (event, hotkey) => {
                toggleRenamingFocusedMediaState();
            }

            /**
             * Toggles the display of the media titles in the gallery items.
             */
            const handleShowTitlesMode = () => {
                toggleMediaTitlesMode();
            }

            /**
             * Closes the gallery and exits the hotkeys context.
             */
            const handleGalleryClose = async () => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("No changes manager available to close the gallery.");
                    return;
                }

                if ($me_gallery_changes_manager.ChangesAmount > 0) {
                    let amount_deleted_medias = $me_gallery_changes_manager.DeletedMedias.length;

                    await processMediaChanges();

                    emitPlatformMessage(`Deleted ${amount_deleted_medias} medias.`);
                }

                emitCloseGallery();
            }


            /**
             * Handle Gallery Exit. This exists the gallery hotkeys control but doesn't close the gallery.
             */
            const handleGalleryExit = () => {
                disableGalleryHotkeys();
            }
            
            /**
             * Handles the WASD Gallery grid movement.
             * @param {KeyboardEvent} hotkey_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             * @deprecated
             */
            const handleGalleryGridMovement = async (hotkey_event, hotkey) => {
                const media_count = active_medias.length;
                const active_media_range = getLoadedMediaRange();
                const media_per_row = Math.floor(getMediaItemsPerRow() || 6);
                const row_count = Math.ceil(media_count / media_per_row);

                let new_focus_index = media_focus_index;

                switch (hotkey.KeyCombo) {
                    case "w": // up
                        new_focus_index -= media_per_row;

                        if (new_focus_index < 0) { // overflow
                            if (shouldMoveHotkeysContext(true, false, false, false)) return;
                        }   

                        new_focus_index = new_focus_index < 0 ? ((row_count - 1) * media_per_row) + media_focus_index : new_focus_index;
                        new_focus_index = new_focus_index >= active_media_range.end ? new_focus_index - media_per_row : new_focus_index;
                        break;
                    case "a": // left
                        new_focus_index = media_focus_index - 1;

                        if (new_focus_index < 0) { // overflow
                            if (shouldMoveHotkeysContext(false, false, false, true)) return;
                        }

                        new_focus_index = new_focus_index < 0 ? media_count - 1 : new_focus_index;

                        break;
                    case "s": // down
                        new_focus_index += media_per_row;

                        if (new_focus_index > active_media_range.end) { // overflow
                            if (shouldMoveHotkeysContext(false, false, true, false)) return;
                        }

                        new_focus_index = new_focus_index > active_media_range.end ? media_focus_index - ((row_count - 1) * media_per_row) : new_focus_index;
                        new_focus_index = new_focus_index < 0 ? new_focus_index + media_per_row : new_focus_index;
                        break;
                    case "d": // right
                        new_focus_index = media_focus_index + 1;

                        if (new_focus_index > active_media_range.end) { // overflow
                            if (shouldMoveHotkeysContext(false, true, false, false)) return;
                        }

                        new_focus_index = new_focus_index > active_media_range.end ? 0 : new_focus_index;

                        break;
                }

                let old_focus_index = media_focus_index;
                media_focus_index = new_focus_index;

                console.log(`media_focus_index: ${media_focus_index}`);

                await manageInfiniteScroll();

                if (auto_select_focused_media) {
                    // selectFocusedMedia();
                } else if (auto_stage_delete_focused_media) {
                    let delete_range = (hotkey.KeyCombo === "w" || hotkey.KeyCombo === "s") && Math.abs(old_focus_index - new_focus_index) <= media_per_row + 1;

                    if (delete_range) {
                        let start_index = Math.min(old_focus_index, new_focus_index);
                        let end_index = Math.max(old_focus_index, new_focus_index);

                        // stageMediaRangeDeletion(start_index, end_index);
                    } else {
                        stageFocusedMediaDeletion();
                    }
                }
            }

            /**
             * Handles the movement to the last media in active_medias.
             */
            const handleMovetoLastActiveMedia = () => media_focus_index = active_medias.length - 1;

            /**
             * Handles the movement to the first media in active_medias.
             */
            const handleMovetoFirstActiveMedia = () => media_focus_index = 0;
            

            /**
             * Toggles the magnify focused media mode.
             */
            const handleMagnifyFocusedMedia = () => {
                magnify_focused_media = !magnify_focused_media;
            }

            /**
             * Opens the media viewer on the focused media item.
             */
            const handleOpenMedia = () => {
                openMediaViewerOnIndex(media_focus_index);
            }

            /**
             * Handles the deletion of the focused media item.
             * @param {KeyboardEvent} hotkey_event
             */
            const handleMediaStageDeletionMode = hotkey_event => {
                const focused_media = getFocusedMedia();
                if (focused_media == null) {
                    console.warn("In MediaExplorerGallery.handleMediaStageDeletionMode: No focused media available to select.");
                    return;
                }

                const is_keydown = hotkey_event.type === "keydown";
                const is_keyup = hotkey_event.type === "keyup";
                const is_unknown_event = !(is_keydown || is_keyup);

                if (is_unknown_event) {
                    console.warn("In MediaExplorerGallery.handleMediaStageDeletionMode: Unknown hotkey event type. Expected 'keydown' or 'keyup'.");
                    return;
                }

                if (is_keyup) {
                    return disableAutoSelectMode();
                }

                if (!auto_select_mode_enabled) {
                    determineAutoSelectActivation(hotkey_event, false);
                } else {
                    return;
                }

                auto_select_mode_adds = !isMediaStagedForDeletion(focused_media);

                setMediaDeletionState(focused_media, auto_select_mode_adds);
            }

            /**
             * Handles the selection of the focused media item.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleMediaSelectAddMode = hotkey_event => {
                const focused_media = getFocusedMedia();
                if (focused_media == null) {
                    console.warn("In MediaExplorerGallery.handleMediaSelectAddMode: No focused media available to select.");
                    return;
                }

                const is_keydown = hotkey_event.type === "keydown";
                const is_keyup = hotkey_event.type === "keyup";
                const is_unknown_event = !(is_keydown || is_keyup);

                if (is_unknown_event) {
                    console.warn("In MediaExplorerGallery.handleMediaSelectAddMode: Unknown hotkey event type. Expected 'keydown' or 'keyup'.");
                    return;
                }

                if (is_keyup) {
                    return disableAutoSelectMode();
                }

                if (!auto_select_mode_enabled) {
                    determineAutoSelectActivation(hotkey_event, true);
                } else {
                    return;
                }

                auto_select_mode_adds = !isMediaYanked(focused_media);

                setMediaSelectedState(focused_media, auto_select_mode_adds);
            }

            /**
             * Handles the yanking of the selected medias. the actual yanking only happens when the user pastes the yanked medias.
             */
            const handleSetYankedMedias = () => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("No changes manager available to stage the yanking of the selected medias.");
                    return;
                }
                let selected_medias = $me_gallery_changes_manager.MovedMedias;

                $me_gallery_changes_manager.clearAllMoveChanges();

                // FIXME: Implement an interface for bothe MediaChangesEmitter and MediasChangesManager to accept something like a MediaLike element.
                // @ts-ignore - this could be an issue but it's complex so, i will Fix it later.
                me_gallery_yanked_medias.set(selected_medias);

                emitPlatformMessage(`Yanked ${selected_medias.length} medias.`);
            }

            /**
             * Emits an event to paste the yanked medias to the MediaExplorer, who handles the actual pasting.
             */
            const handleYankPaste = () => {
                dispatch("paste-yanked-medias");
            }

            /**
             * handles gallery changes reset.
             */
            const handleGalleryReset = async () => {
                if ($me_gallery_changes_manager === null) return;

                const changes_count = $me_gallery_changes_manager.ChangesAmount;

                const medias_per_row = getMediaItemsPerRow();

                if (medias_per_row == null) {
                    console.warn("In MediaExplorerGallery.handleGalleryReset: No media items per row available. Cannot determine if confirmation is required.");
                    return;
                }

                const requires_confirmation = changes_count >= (1.5 * medias_per_row);

                if (requires_confirmation) {
                    const user_confirmed = await confirmPlatformMessage({
                        message_title: "Reset selection",
                        question_message: `Do you really want to reset the selection of ${changes_count}?`,
                        auto_focus_cancel: true,
                        danger_level: 1
                    });

                    if (user_confirmed !== 1) {
                        console.log("Confirmation cancelled.");
                        return;
                    }
                }

                $me_gallery_changes_manager.clearAllChanges();
            }

            /**
             * Toggles the heavy rendering of the medias in the gallery.
             */
            const handleToggleHeavyRendering = () => {
                if (enable_gallery_performance_mode) {
                    emitPlatformMessage("Heavy rendering is disabled because the amount of medias loaded would likely kill your CPU with heavy rendering.");
                    return;
                }

                enable_gallery_heavy_rendering = !enable_gallery_heavy_rendering;
            }

            /*----------  Navigation hotkey handlers  ----------*/
            
                /**
                 * Jumps media focus to a specific media position by vim-motion.
                 * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
                 */
                const handleJumpToMediaOrder = async (event, hotkey) => {
                    if (!hotkey.WithVimMotion || !hotkey.HasMatch) return;

                    let requested_order = hotkey.MatchMetadata?.MotionMatches[0];

                    if (requested_order == null || isNaN(requested_order)) {
                        console.log("Corrupted hotkey", hotkey);
                        console.warn("Invalid media order requested.");
                        return;
                    }

                    requested_order -= 1 // 1-based to 0-based index.

                    const previous_media_focus_index = media_focus_index;

                    await setMediaFocusIndex(requested_order, true);

                    if (auto_select_mode_enabled) {
                        await tick();

                        handleCursorSelection(previous_media_focus_index, media_focus_index);
                    }
                }

        /*=====  End of Hotkeys  ======*/
        
        /*=============================================
        =            Debug            =
        =============================================*/

                /**
                 * Returns the object where all the gallery debug state is stored.
                 * @returns {Object}
                 */
                const debugMEG__getGalleryState = () => {
                    if (!browser || !debug_mode) return {};

                    const GALLERY_DEBUG_STATE_NAME = "meg_gallery_debug_state";

                    // @ts-ignore
                    if (!globalThis[GALLERY_DEBUG_STATE_NAME]) {
                        // @ts-ignore
                        globalThis[GALLERY_DEBUG_STATE_NAME] = {

                        };   
                    }

                    // @ts-ignore
                    return globalThis[GALLERY_DEBUG_STATE_NAME];
                }
        
                /**
                 * Attaches debug methods to the globalThis object for debugging purposes.
                 * @returns {void}
                 */
                const debugMEG__attachDebugMethods = () => {
                    if (!browser || !debug_mode) return;

                    const meg_gallery_debug_state = debugMEG__getGalleryState();

                    // @ts-ignore - for debugging purposes we do not care whether the globalThis object has the method name. same reason for all other ts-ignore in this function.
                    meg_gallery_debug_state.printGalleryState = debugMEG__printGalleryState;

                    // @ts-ignore
                    meg_gallery_debug_state.Page = $page;

                    // @ts-ignore - state retrieval functions.
                    meg_gallery_debug_state.State = {
                        getActiveMedias: () => active_medias,
                        getGridNavigationWrapper: () => the_grid_navigation_wrapper,
                        getLastMediaAddedToActiveMedias: () => last_media_added_to_active_medias,
                        getMediaFocusIndex: () => media_focus_index,
                        getOrderedMedias: () => ordered_medias,
                        getChangesManager: () => $me_gallery_changes_manager,
                    }

                    // @ts-ignore - Internal method references.
                    meg_gallery_debug_state.Methods = {
                        getMediaItemsPerRow,
                        getOptimalGalleryRowCount,
                        getOptimalMediaBatchSize,
                        getLoadedMediaRange,
                        getFocusedMediaElement,
                        getMediaDisplayIndexByOrder,
                        getMediaItemElementByOrder,
                        getFocusedMedia,
                    }
                }

                /**
                 * Prints the whole gallery state to the console.
                 * @returns {void}
                 */
                const debugMEG__printGalleryState = () => {
                    console.log("%cMediaExplorerGallery State", "color: green; font-weight: bold;");
                    console.group("Properties");
                    console.log(`active_medias.length: ${active_medias.length}`);
                    console.log("active_medias:", active_medias);
                    console.log(`ordered_medias.length: ${ordered_medias.length}`);
                    console.log("ordered_medias:", ordered_medias);
                    console.log(`media_focus_index: ${media_focus_index}`);
                    console.groupEnd();
                    console.group("Navigation");
                    console.log(`has_proceeding_medias: ${has_proceeding_medias}`);
                    console.log("last_media_added_to_active_medias: %O", last_media_added_to_active_medias);
                    console.log(`enable_gallery_hotkeys: ${enable_gallery_hotkeys}`);
                    console.log("the_grid_navigation_wrapper: %O", the_grid_navigation_wrapper);
                    console.groupEnd();
                    console.group("Performance");
                    console.log(`enable_gallery_performance_mode: ${enable_gallery_performance_mode}`);
                    console.log(`regulating_active_medias_load: ${regulating_active_medias_load}`);
                }

                /**
                 * Attaches an arbitrary object as a globalThis.meg_timeline_states.<group_name>{...timestamp -> object }.
                 * @param {string} group_name
                 * @param {object} object_to_snapshot
                 * @returns {void}
                 */
                const debugMEG__attachSnapshot = (group_name, object_to_snapshot) => {
                    if (!browser || !debug_mode) return;

                    const stack = new Error().stack;
                    const datetime_obj = new Date();
                    const timestamp = `${datetime_obj.toISOString()}-${datetime_obj.getTime()}`;

                    const snapshot = {
                        timestamp,
                        stack,
                        object_to_snapshot,
                    }

                    const debug_object = debugMEG__getGalleryState();

                    // @ts-ignore - that meg_timeline_states exists on globalThis if not, create it.
                    if (!debug_object.timeline_states) {
                        // @ts-ignore
                        debug_object.timeline_states = {};
                    }

                    // @ts-ignore
                    if (!debug_object.timeline_states[group_name]) {
                        // @ts-ignore
                        debug_object.timeline_states[group_name] = [];
                    }

                    // @ts-ignore
                    debug_object.timeline_states[group_name].push(snapshot);
                }
        
        /*=====  End of Debug  ======*/
        
        /*=============================================
        =            Navigation            =
        =============================================*/

            /**
             * Applies cursor navigation selection if the auto select mode is enabled.
             * @param {number} previous_index
             * @param {number} current_index
             * @returns {void}
             */
            const applyCursorNavigationSelection = (previous_index, current_index) => {
                if (!auto_select_mode_enabled) return;

                handleCursorSelection(previous_index, current_index);
            }

            /**
             * Drops the grid navigation wrapper if it exists.
             */
            const dropGridNavigationWrapper = () => {
                if (the_grid_navigation_wrapper != null) {
                    the_grid_navigation_wrapper.destroy();
                }
            }

            /**
             * Determines whether the cursors should be moved from an to one of the edges of
             * the gallery grid.
             * @param {import('@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils').GridWrappedValue} cursor_wrapped_value
             * @param {number} current_focus_index
             * @returns {boolean}
             */
            const cursorMovementHasOverflowedContent = (cursor_wrapped_value, current_focus_index) => {
                // TODO: manage bottom overflow from the last row when there are no more medias that can be loaded.
                return (
                    (current_focus_index === 0 && cursor_wrapped_value.overflowed_left) ||
                    (current_focus_index === (ordered_medias.length - 1) && cursor_wrapped_value.overflowed_right)
                );
            }
        
            /**
             * Returns the grid selectors for the  Media Explorer Gallery navigation gird.
             * @returns {import('@common/interfaces/common_actions').GridSelectors}
             */
            const getGridSelectors = () => {
                const grid_parent_selector = `#meg-gallery`;

                return {
                    grid_parent_selector,
                    grid_member_selector: `> .${media_item_html_class}`,
                }
            }

            /**
             * Handles the Cursor update event emitted by the_grid_navigation_wrapper.
             * @type {import("@common/keybinds/CursorMovement").CursorPositionCallback}
             */
            const handleCursorUpdate = (cursor_wrapped_value) => {
                const lost_hotkey_control = shouldMoveHotkeysContext(
                    cursor_wrapped_value.overflowed_top,
                    cursor_wrapped_value.overflowed_right,
                    cursor_wrapped_value.overflowed_bottom,
                    cursor_wrapped_value.overflowed_left
                );

                if (lost_hotkey_control) {
                    return true;
                }

                console.debug(`In MediaExplorerGallery.handleCursorUpdate: media_focus_index: ${media_focus_index}, cursor_wrapped_value:`, cursor_wrapped_value);

                if (cursorMovementHasOverflowedContent(cursor_wrapped_value, media_focus_index)) {
                    return handleCursorUpdateFromEdgeToEdgeOfContent(cursor_wrapped_value);
                }

                const media_index_in_displayed_array = cursor_wrapped_value.value;

                const focused_media = getMediaByDisplayIndex(media_index_in_displayed_array);

                if (focused_media == null) {
                    console.error(`In @pages/MediaExplorer/sub-components/MediaExplorerGallery/MediaExplorerGallery.handleCursorUpdate: No media found at active_medias`);
                    return;
                }

                applyCursorNavigationSelection(media_focus_index, focused_media.Order);

                setMediaFocusIndex(focused_media.Order, false);
            }

            /**
             * callback passed to the grid navigation wrapper. is triggered after it's grid sequence 
             * is reconstructed.
             * @type {import("@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils").GridSequenceReadyCallback}
             */
            const handleGridSequenceReady = async new_sequence => {
                if (grid_navigation_change_callbacks.length === 0) return;

                for (const callback of grid_navigation_change_callbacks) {
                    if (callback.constructor.name === "AsyncFunction") {
                        await callback();
                    } else if (callback.constructor.name === "Function") {
                        callback();
                    }
                }

                console.debug(`In MediaExplorerGallery.handleGridSequenceReady: Grid sequence ready. Executed ${grid_navigation_change_callbacks.length} callbacks.`);

                grid_navigation_change_callbacks = [];
            }

            /**
             * called by handleCursorUpdate exclusively. handles the cursor update when the cursor overflows from one of the edges of the content. To
             * be clear, this would be when 'from c = 0 and movement direction is left' or 'from c = last_index and movement direction is right'.
             * @type {import("@common/keybinds/CursorMovement").CursorPositionCallback}
             */
            const handleCursorUpdateFromEdgeToEdgeOfContent = (cursor_wrapped_value) => {
                if (the_grid_navigation_wrapper == null) {
                    console.warn("In MediaExplorerGallery.handleCursorUpdateFromEdgeToEdgeOfContent: No grid navigation wrapper.");
                    return;
                }

                const highest_order_media = getMediaByOrder(ordered_medias.length - 1);
                const lowest_order_media = getMediaByOrder(0);

                if (highest_order_media == null || lowest_order_media == null) {
                    console.log("highest_order_media:", highest_order_media, "lowest_order_media:", lowest_order_media);
                    console.error("In MediaExplorerGallery.handleCursorUpdateFromEdgeToEdgeOfContent: No highest order media available. This is impossible.");
                    return;
                }

                const should_jump_from_first_to_last = media_focus_index === lowest_order_media.Order && cursor_wrapped_value.overflowed_left;
                const should_jump_from_last_to_first = media_focus_index === highest_order_media.Order && cursor_wrapped_value.overflowed_right;

                if (!(should_jump_from_first_to_last || should_jump_from_last_to_first)) {
                    console.warn("In MediaExplorerGallery.handleCursorUpdateFromEdgeToEdgeOfContent: No edge overflow detected. Nothing to do.");
                    return;
                }

                const correct_media = should_jump_from_first_to_last ? highest_order_media : lowest_order_media;

                applyCursorNavigationSelection(media_focus_index, correct_media.Order);

                setMediaFocusIndex(correct_media.Order, true);
            }

            /**
             * Returns whether the current navigation grid in sync with the active_medias array.
             * This means that the grid doesn't have the same amount of items as the active_medias array.
             * @returns {boolean}
             */
            const isNavigationGridOutdated = () => {
                if (the_grid_navigation_wrapper == null) {
                    console.warn("No grid navigation wrapper available to check if the navigation grid is outdated.");
                    return true;
                }

                const active_medias_length = active_medias.length;
                const grid_2D_sequence_length = the_grid_navigation_wrapper.MovementController.Grid.SequenceLength;

                return active_medias_length !== grid_2D_sequence_length;
            }

            /**
             * Returns a promise that is resolved when the grid navigation wrapper finishes reconstructing the grid.
             * If the navigation grid is in sync with the active_medias array, the promise is resolved immediately.
             * Additionally, a timeout parameter(in milliseconds) can be passed to. If the grid doesn't update in that time, the promise is rejected.
             * @param {number} [timeout=-1]
             * @returns {Promise<void>}
             */
            const waitForGridSync = (timeout = -1) => {
                return new Promise((resolve, reject) => {
                    const grid_outdated = isNavigationGridOutdated();
                    if (!grid_outdated)  {
                        resolve();
                        return;
                    }

                    if (timeout > 0) {
                        setTimeout(() => {
                            reject(new Error("Grid navigation wrapper didn't update in time."));
                        }, timeout);
                    }

                    grid_navigation_change_callbacks.push(() => resolve());
                });
            }

            /**
             * Sets the grid navigation wrapper required data.
             * @param {import("@libs/LiberyHotkeys/hotkeys_context").default} hotkeys_context
             */
            const setGridNavigationWrapper = (hotkeys_context) => {
                if (!browser) return;

                if (the_grid_navigation_wrapper != null) {
                    the_grid_navigation_wrapper.destroy();
                }

                const grid_selectors = getGridSelectors();

                const matching_elements_count = document.querySelectorAll(grid_selectors.grid_parent_selector).length;

                if (matching_elements_count !== 1) {
                    throw new Error(`Expected 1 element matching the grid member selector "${grid_selectors.grid_member_selector}" but found ${matching_elements_count}.`);
                }
 
                const initial_cursor_position = getFocusedMediaDisplayIndex();

                console.log(`In MediaExplorerGallery.setGridNavigationWrapper: Initial cursor position is ${initial_cursor_position}.`);

                the_grid_navigation_wrapper = new CursorMovementWASD(grid_selectors.grid_parent_selector, handleCursorUpdate, {
                    initial_cursor_position: initial_cursor_position,
                    sequence_item_name: ui_core_dungeon_references.MEDIA.EntityName,
                    sequence_item_name_plural: ui_core_dungeon_references.MEDIA.EntityNamePlural,
                    grid_member_selector: grid_selectors.grid_member_selector,
                    goto_item_finalizer: "", // Disable goto item handler of the CursorMovementWASD wrapper.
                    on_mutation_cursor_correction_callback: translateCursorAfterContentMutation,
                });

                the_grid_navigation_wrapper.MovementController.setGridSequenceReadyCallback(handleGridSequenceReady);

                the_grid_navigation_wrapper.setup(hotkeys_context);

                console.debug("In MediaExplorerGallery.setGridNavigationWrapper: active_medias.length:", active_medias.length);
                console.debug("row_count:", the_grid_navigation_wrapper.MovementController.Grid.length);
                console.log("Grid navigation wrapper setup FINISHED.")
            }

            /**
             * Correct the current cursor of the the_grid_navigation_wrapper(CursorMovementWASD) to match the 
             * focused media when the gallery's content changes.
             * @type {import("@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils").CurrentCursorCorrectionProvider}
             */
            const translateCursorAfterContentMutation = (cursor_value) => {
                return getFocusedMediaDisplayIndex();
            }

            /**
             * Updates the navigation grid.
             * @returns {Promise<void>}
             */
            async function updateNavigationGrid () {
                if (the_grid_navigation_wrapper == null) {
                    console.warn("No grid navigation wrapper available to update the navigation grid.");
                    return;
                }

                await tick();

                the_grid_navigation_wrapper.MovementController.scanGridMembers();

                if (media_focus_index !== 0) {
                    const display_index = getMediaDisplayIndexByOrder(media_focus_index);
                    
                    the_grid_navigation_wrapper.updateCursorPositionSilently(display_index);
                }
            }
        
        /*=====  End of Navigation  ======*/
        
        /*=============================================
        =            Selection            =
        =============================================*/

            /**
             * Activates the passed selection mode. Only one may be active at a time.
             * @param {boolean} active_yank_select_mode
             * @param {boolean} active_deletion_select_mode
             * @returns {void}
             */
            const activateAutoSelectMode = (active_yank_select_mode, active_deletion_select_mode) => {
                auto_select_focused_media = active_yank_select_mode && !active_deletion_select_mode;
                auto_stage_delete_focused_media = active_deletion_select_mode && !active_yank_select_mode;

                // @ts-ignore - ts is stupid. XOR works well with boolean values.
                auto_select_mode_enabled = !!(auto_select_focused_media ^ auto_stage_delete_focused_media);
                console.debug(`In MediaExplorerGallery.activateAutoSelectMode: ${getSelectionState()}`);
            }

            /**
             * returns the current selection state as human readable text.
             * @returns {string}
             */
            const getSelectionState = () => {
                // TODO: move this to the debug section.
                return `
                    \nauto_select_focused_media: ${auto_select_focused_media}, 
                    \nauto_stage_delete_focused_media: ${auto_stage_delete_focused_media}, 
                    \nauto_select_mode_enabled: ${auto_select_mode_enabled}
                    \nauto_select_mode_adds: ${auto_select_mode_adds},
                `;
            }

            /**
             * Disables the auto select mode.
             * @returns {void}
             */
            const disableAutoSelectMode = () => {
                auto_select_focused_media = false;
                auto_stage_delete_focused_media = false;
                auto_select_mode_enabled = false;
                auto_select_mode_adds = true;
            }

            /**
             * Determines whether the keyboard event should activate the auto select mode. The 
             * exact selection mode is determined by the parameter passed to this function. The
             * function will not make an attempt to confirm whether the key represents a valid 
             * selection mode, as it holds no such knowledge, therefore this is left up to the 
             * caller.
             * @param {KeyboardEvent} keyboard_event
             * @param {boolean} [is_yank_select_mode=true]
             * @returns {boolean} - Whether the auto select has been activated.
             */
            const determineAutoSelectActivation = (keyboard_event, is_yank_select_mode = true) => {

                let is_activated = !auto_select_mode_enabled;

                if ( is_activated ) {
                    activateAutoSelectMode(is_yank_select_mode, !is_yank_select_mode);
                }

                console.debug(`In MediaExplorerGallery.determineAutoSelectActivation: Auto select mode activated: ${is_activated}. Yank select mode: ${is_yank_select_mode ? 'selection' : 'deletion'}.`);

                return is_activated;
            }

            /**
             * Called navigation handlers. Adds or removes all the medias between
             * the old and new focus index(media order) to the active selection.
             * @param {number} old_focus_index
             * @param {number} new_focus_index
             * @returns {void}
             */
            const handleCursorSelection = (old_focus_index, new_focus_index) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.handleCursorSelection: No changes manager available to handle the cursor selection.");
                    return;
                }

                const no_diff = old_focus_index === new_focus_index;
                const auto_select_mode_disabled = !auto_select_mode_enabled;

                if (no_diff || auto_select_mode_disabled) return;

                const start_index = Math.min(old_focus_index, new_focus_index);
                const end_index = Math.max(old_focus_index, new_focus_index);

                console.debug(`In MediaExplorerGallery.handleCursorSelection: Selecting media range [${start_index}, ${end_index}]: ${getSelectionState()}`);

                if (auto_stage_delete_focused_media) {
                    stageMediaRangeDeletion(start_index, end_index, auto_select_mode_adds);
                } else {
                    selectMediaRange(start_index, end_index, auto_select_mode_adds);
                }
            }

            /**
             * handles the state of the changes emitter and it's subscribers.
             * @returns {void}
             */
            const handleChangesEmitterState = () => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.handleChangesEmitterState: No changes manager available to handle the changes emitter state.");
                    return;
                }

                $me_gallery_changes_manager.clearAllChangeSubscriptions();
            }

            /**
             * Returns whether the passed ordered media is yanked.
             * @param {import('@models/Medias').OrderedMedia} ordered_media
             * @return {boolean}
             */
            const isMediaYanked = (ordered_media) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.isMediaYanked: No changes manager available to check if the media is yanked.");
                    return false;
                }

                const current_media_change = $me_gallery_changes_manager.getMediaChangeType(ordered_media.uuid);

                return current_media_change === media_change_types.MOVED;
            }

            /**
             * Returns whether the passed ordered media is staged for deletion.
             * @param {import('@models/Medias').OrderedMedia} ordered_media
             * @returns {boolean}
             */
            const isMediaStagedForDeletion = (ordered_media) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.isMediaStagedForDeletion: No changes manager available to check if the media is staged for deletion.");
                    return false;
                }

                const current_media_change = $me_gallery_changes_manager.getMediaChangeType(ordered_media.uuid);

                return current_media_change === media_change_types.DELETED;
            }

            /**
             * Modifies the Media Selected state of the focused media item.
             * Whether the media is added or removed from the selection is determined by the
             * auto_select_mode_adds property.
             * @returns {void}
             */
            const modifyFocusedMediaSelectedState = () => {
                modifyMediaSelectedStateByOrder(media_focus_index)
            }

            /**
             * Modifies the Media Deletion state of the focused media item.
             * Whether the media is added or removed from the deletion is determined by the
             * auto_select_mode_adds property.
             * @returns {void}
             */
            const modifyFocusedMediaDeletionState = () => {
                modifyMediaDeletionStateByOrder(media_focus_index);
            }
        
            /**
             * Modifies the Selected state of the media by its order. Whether the media is added
             * or removed from the selection is determined by the auto_select_mode_adds property.
             * @param {number} media_order
             * @returns {void}
             */
            const modifyMediaSelectedStateByOrder = (media_order) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.addMediaToActiveSelectionByOrder: No changes manager available to add the media to the active selection.");
                    return;
                }

                let ordered_media = getMediaByOrder(media_order);

                if (ordered_media == null) {
                    console.warn(`In MediaExplorerGallery.addMediaToActiveSelectionByOrder: No media found at order ${media_order}.`);
                    return;
                }

                setMediaSelectedState(ordered_media, auto_select_mode_adds);
            }

            /**
             * modifies the deletion state of the media by its order. Whether the media is added
             * or removed from the deletion is determined by the auto_select_mode_adds property.
             * @param {number} media_order
             * @returns {void}
             */
            const modifyMediaDeletionStateByOrder = (media_order) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.addMediaToActiveSelectionByOrder: No changes manager available to add the media to the active selection.");
                    return;
                }

                let ordered_media = getMediaByOrder(media_order);

                if (ordered_media == null) {
                    console.warn(`In MediaExplorerGallery.addMediaToActiveSelectionByOrder: No media found at order ${media_order}.`);
                    return;
                }

                setMediaDeletionState(ordered_media, auto_select_mode_adds);
            }

            /**
             * Stages the focused media to be deleted.
             * @returns {void}
             */
            const stageFocusedMediaDeletion = () => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("No changes manager available to stage the deletion of the focused media item.");
                    return;
                }

                let focused_media = getFocusedMedia();
                if (focused_media === null) return;

                let is_delete_already_staged = isMediaStagedForDeletion(focused_media);

                if (is_delete_already_staged && auto_stage_delete_focused_media) return;

                toggleMediaDeletion(focused_media);
            }

            /**
             * Adds or removes a range of medias to the normal selection.
             * @param {number} start_index
             * @param {number} end_index
             * @param {boolean} add_to_selection
             * @returns {void}
             */
            const selectMediaRange = (start_index, end_index, add_to_selection) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("No changes manager available to stage the selection of the media range.");
                    return;
                }

                if (start_index < 0 || end_index >= ordered_medias.length || start_index > end_index) {
                    console.warn(`In MediaExplorerGallery.selectMediaRange: Invalid range ${start_index} - ${end_index}.`);
                    return;
                }

                let media_range = ordered_medias.slice(start_index, end_index + 1);

                for (let media_item of media_range) {
                    setMediaSelectedState(media_item, add_to_selection);
                }
            }

            /**
             * Stages the deletion of the given media inclusive range. so if 4 - 10 are passed, medias on positions 4 to 10 will be staged for deletion.
             * @param {number} start_index
             * @param {number} end_index
             * @param {boolean} add_to_selection
             */
            const stageMediaRangeDeletion = (start_index, end_index, add_to_selection) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("No changes manager available to stage the deletion of the media range.");
                    return;
                }

                if (start_index < 0 || end_index >= ordered_medias.length || start_index > end_index) {
                    console.warn(`In MediaExplorerGallery.stageMediaRangeDeletion: Invalid range ${start_index} - ${end_index}.`);
                    return;
                }

                let media_range = ordered_medias.slice(start_index, end_index + 1);

                for (let media_item of media_range) {
                    setMediaDeletionState(media_item, add_to_selection);
                }
            }

            /**
             * Sets the given media Selection state. If the overwrite flag is set, this function will
             * undo any other selection state(e.g Deletion). otherwise it will shy away if the media
             * has other selection states.
             * @param {import('@models/Medias').OrderedMedia} ordered_media
             * @param {boolean} add_to_selection 
             * @param {boolean} [overwrite=false] 
             * @returns {void}
             */
            const setMediaSelectedState = (ordered_media, add_to_selection, overwrite=false) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.setMediaSelectedState: No changes manager available to set the media selection state.");
                    return;
                }

                let is_media_deleted = isMediaStagedForDeletion(ordered_media);

                if (is_media_deleted) {
                    if (!overwrite) {
                        console.debug(`In MediaExplorerGallery.setMediaSelectedState: Media ${ordered_media.uuid} is staged for deletion and overwrite was not set. Silently cowardly ignoring the request.`);
                        return;
                    }

                    $me_gallery_changes_manager.unstageMediaDeletion(ordered_media.uuid);
                }

                let is_media_yanked = isMediaYanked(ordered_media);

                if (is_media_yanked === add_to_selection) return; // Nothing to do

                if (add_to_selection) {
                    // on this meg gallery, we don't actually commit move media changes, only delete changes. instead of committing the move changes, we used them as
                    // a select mechanism. if there are any move changes when committing, we will unstage them first.
                    let fake_inner_category = new InnerCategory({
                        name: "me-gallery-mock-category", 
                        uuid: "me-gallery-mock-category", 
                        fullpath: "me-gallery-mock-category",
                        category_thumbnail: "me-gallery-mock-category",
                    });

                    $me_gallery_changes_manager.stageMediaMove(ordered_media.Media, fake_inner_category);
                } else {
                    $me_gallery_changes_manager.unstageMediaMove(ordered_media.uuid);
                }
            }

            /**
             * Sets the given media Stage Deletion state. If the media is currently in the regular
             * selection set, then it will only overwrite said state if the overwrite flag is set.
             * otherwise it will simply fail silently.
             * @param {import('@models/Medias').OrderedMedia} ordered_media
             * @param {boolean} staged_for_deletion
             * @param {boolean} [overwrite=false] 
             * @returns {void}
             */
            const setMediaDeletionState = (ordered_media, staged_for_deletion, overwrite=false) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.setMediaDeletionState: No changes manager available to set the media deletion state.");
                    return;
                }

                let is_media_yanked = isMediaYanked(ordered_media);

                if (is_media_yanked) {
                    if (!overwrite) {
                        console.debug(`In MediaExplorerGallery.setMediaDeletionState: Media ${ordered_media.uuid} is yanked and overwrite was not set. Silently cowardly ignoring the request.`);
                        return;
                    }

                    $me_gallery_changes_manager.clearMediaChanges(ordered_media.uuid);
                }

                const was_already_staged_for_deletion = isMediaStagedForDeletion(ordered_media);

                if (was_already_staged_for_deletion === staged_for_deletion) return; // Nothing to do

                if (staged_for_deletion) {
                    $me_gallery_changes_manager.stageMediaDeletion(ordered_media.Media);
                } else {
                    $me_gallery_changes_manager.unstageMediaDeletion(ordered_media.uuid);
                }
            }

            /**
             * Toggles the deletion state of a given media.
             * @param {import('@models/Medias').OrderedMedia} ordered_media 
             */
            const toggleMediaDeletion = (ordered_media) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.toggleMediaDeletion: No changes manager available to stage the deletion of the media item.");
                    return;
                }

                let is_media_deleted = isMediaStagedForDeletion(ordered_media);

                setMediaDeletionState(ordered_media, !is_media_deleted);
            }

            /**
             * Toggles the select stat of a given media.
             * @param {import('@models/Medias').OrderedMedia} media_item 
             */
            const toggleMediaSelect = (media_item) => {
                if ($me_gallery_changes_manager === null) {
                    console.warn("In MediaExplorerGallery.toggleMediaSelect: No changes manager available to stage the selection of the media item.");
                    return;
                }

                let is_media_yanked = isMediaYanked(media_item);

                setMediaSelectedState(media_item, !is_media_yanked);
            }

        /*=====  End of Selection  ======*/
        
        /*=============================================
        =            Gallery Layout analysis            =
        =============================================*/
            /**
            *  NOTE: If in the future we end up supporting an additional method to determine
             * column count, outside of gridTemplateColumns, we should turn the whole, 
             * Determine CSS column count into a lib.
            */

            /**
             * Returns the amount of medias that fit in a single row of the gallery's grid.
             * @returns {number | undefined} - a float number most likely.
             */
            const getMediaItemsPerRow = () => {
                if (the_grid_navigation_wrapper == null || !the_grid_navigation_wrapper.MovementController.Grid.isUsable()) {
                    return getMediaItemsPerRowFromCSS();
                }

                const current_row = the_grid_navigation_wrapper.MovementController.Grid.getCurrentRow();

                return current_row.length
            }

            /**
             * Returns the optimal number of gallery rows to fit the screen.
             * @param {CSSStyleDeclaration} [meg_gallery_style]
             * @returns {number}
             */
            const getOptimalGalleryRowCount = (meg_gallery_style) => {
                const load_frequency_reduction_factor = 2; // This is a factor how often we load the gallery rows.
                if (load_frequency_reduction_factor <= 0) { // anti-dumb error assertion.
                    throw new Error("In MediaExplorerGallery.getOptimalGalleryRowCount: load_frequency_reduction_factor must be greater than 0.");
                }

                const FUZZY_ROW_COUNT = Math.round(5 * load_frequency_reduction_factor); // If no better value can be determined, use this.

                if (!browser) return FUZZY_ROW_COUNT;

                if (meg_gallery_style === undefined) {
                    meg_gallery_style = getGalleryCSSStyle();

                    if (meg_gallery_style === undefined) {
                        console.warn("In MediaExplorerGallery.getOptimalGalleryRows: No gallery CSS style available.");
                        return FUZZY_ROW_COUNT
                    }
                }

                let row_count = FUZZY_ROW_COUNT;

                if (isGalleryLayoutGrid(meg_gallery_style)) {
                    const row_height = getGridItemsHeightFromCSS(meg_gallery_style);

                    if (row_height === undefined) {
                        console.error(`In MediaExplorerGallery.getOptimalGalleryRows: No row height found in the gallery CSS style.`);
                        return FUZZY_ROW_COUNT;
                    }

                    const viewport_height = window.innerHeight

                    row_count = Math.ceil((viewport_height / row_height) * load_frequency_reduction_factor);
                }

                return row_count
            }

            /**
             * Returns the optimal batch size of medias to load into the gallery.
             * @param {CSSStyleDeclaration} [meg_gallery_style]
             * @returns {number}
             */
            const getOptimalMediaBatchSize = (meg_gallery_style) => {
                const DEFAULT_MEDIA_BATCH_SIZE = 30

                if (!browser) return DEFAULT_MEDIA_BATCH_SIZE; 

                const number_of_items_per_row = getMediaItemsPerRowFromCSS(meg_gallery_style);
                if (number_of_items_per_row === undefined) {
                    console.error(`In MediaExplorerGallery.getOptimalMediaBatchSize: No number of items per row found in the gallery CSS style.`);
                    return DEFAULT_MEDIA_BATCH_SIZE;
                }

                const optimal_row_count = getOptimalGalleryRowCount(meg_gallery_style);
                console.log(`optimal_row_count: ${optimal_row_count}\nnumber_of_items_per_row: ${number_of_items_per_row}`);

                return number_of_items_per_row * optimal_row_count;
            }
            
            /*----------  CSS column count  ----------*/
            
                /**
                 * Returns the ItemsPerRow metric from the gallery's CSSStyleDeclaration. this is 
                 * a fallback for getMediaItemsPerRow, used when the_grid_navigation_wrapper is not available.
                 * @param {CSSStyleDeclaration} [meg_gallery_style]
                 * @returns {number | undefined}
                 */
                const getMediaItemsPerRowFromCSS = (meg_gallery_style) => {
                    if (meg_gallery_style === undefined) {
                        meg_gallery_style = getGalleryCSSStyle();

                        if (meg_gallery_style === undefined) {
                            console.warn("In MediaExplorerGallery.getMediaItemsPerRowFromCSS: No gallery CSS style available.");
                            return undefined;
                        }
                    }

                    return getColumnCountFromGridElement(meg_gallery_style);
                }

                /**
                 * Returns the grid items height determined from the gallery's template columns.
                 * assumes the items have a 1/1 aspect ratio.
                 * @param {CSSStyleDeclaration} [meg_gallery_style]
                 * @returns {number | undefined}
                 */
                const getGridItemsHeightFromCSS = (meg_gallery_style) => {
                    if (meg_gallery_style === undefined) {
                        meg_gallery_style = getGalleryCSSStyle();

                        if (meg_gallery_style === undefined) {
                            console.warn("In MediaExplorerGallery.getGridItemsHeightFromCSS: No gallery CSS style available.");
                            return undefined;
                        }
                    }

                    const grid_template_columns = meg_gallery_style.gridTemplateColumns;

                    if (grid_template_columns === '') {
                        console.error(`In MediaExplorerGallery.getGridItemsHeightFromCSS: No grid template columns found in the element style.`);
                        return undefined;
                    }

                    const column_widths = parseInt(grid_template_columns.split(/\s+/)[0], 10);

                    if (isNaN(column_widths)) {
                        console.error(`In MediaExplorerGallery.getGridItemsHeightFromCSS: Invalid column width found in the grid template columns: ${grid_template_columns}`);
                        return undefined;
                    }

                    return column_widths; 
                }

                /**
                 * Returns the grid column count of the gallery's CSSStyleDeclaration.
                 * @param {CSSStyleDeclaration} element_style
                 * @returns {number | undefined}
                 */
                const getColumnCountFromGridElement = (element_style) => {
                    if (!isGalleryLayoutGrid(element_style)) {
                        console.error(`In MediaExplorerGallery.getColumnCountFromGridElement: recieved element style.display = '${element_style.display}' not 'gird'`);
                        return undefined;
                    }

                    const template_columns_string = element_style.gridTemplateColumns;

                    if (template_columns_string === '') {
                        console.error(`In MediaExplorerGallery.getColumnCountFromGridElement: No grid template columns found in the element style.`);
                        return undefined;
                    }

                    const template_columns = template_columns_string.trim().split(/\s+/);

                    return template_columns.length;
                }

                /**
                 * Returns the CSSStyleDeclaration of the meg gallery ul element.
                 * @returns {CSSStyleDeclaration | undefined}
                 */
                const getGalleryCSSStyle = () => {
                    if (!browser) {
                        console.warn("In MediaExplorerGallery.getGalleryCSSStyle: No browser environment available. Cannot get gallery CSS style.");
                        return undefined;
                    }

                    const grid_selectors = getGridSelectors();

                    const meg_gallery_element = document.querySelector(grid_selectors.grid_parent_selector);

                    if (meg_gallery_element ===null) {
                        console.error(`In MediaExplorerGallery.getGalleryCSSStyle: No meg gallery element found with selector "${grid_selectors.grid_parent_selector}".`);
                        return undefined;
                    }

                    return getComputedStyle(meg_gallery_element);
                }

                /**
                 * Returns whether the given CSSStyleDeclaration object has a display grid
                 * @param {CSSStyleDeclaration} meg_gallery_style
                 * @returns {boolean}
                 */
                const isGalleryLayoutGrid = (meg_gallery_style) => {
                    return meg_gallery_style.display === "grid";
                }
        
        /*=====  End of Gallery Layout analysis  ======*/
        
        /*=============================================
        =            Intersection observer            =
        =============================================*/
        
            /**
             * Adds an item to the gallery item intersection observer.
             * @param {Element} gallery_item_node - The gallery item node to observe.
             * @param {import('@models/Medias').OrderedMedia} ordered_media
             * @returns {void}
             */
            const addMEGItemToObservationList = (gallery_item_node, ordered_media) => {
                if (gallery_item_intersection_observer == null) {
                    console.warn("In MediaExplorerGallery.addMEGItemToObservationList: No intersection observer available to observe the gallery item.");
                    return;
                }

                gallery_item_intersection_observer.observe(gallery_item_node);
            }
                
            /**
             * Drops the gallery item intersection observer.
             * @returns {void}
             */
            const dropGalleryItemsIntersectionObserver = () => {
                if (gallery_item_intersection_observer == null) {
                    return;
                }

                gallery_item_intersection_observer.disconnect();
                gallery_item_intersection_observer = null;
            }

            /**
             * Defines the intersection observer. If it's already 
             * returns immediately and does nothing. The intersection
             * observer reference is gallery_item_intersection_observer.
             * @returns {void}
             */
            const initializeGalleryItemsIntesectionObserver = () => {
                if (gallery_item_intersection_observer != null) {
                    return;
                }

                /**
                 * The intersection observer options.
                 * @type {IntersectionObserverInit}
                 */
                const intersection_observer_options = {
                    threshold: 0.2,
                }

                gallery_item_intersection_observer = new IntersectionObserver(
                    megGalleryIntersectionObserverCallback, 
                    intersection_observer_options
                );
            }

            /**
             * Intersection observer callback. set by initializeGalleryItemsIntesectionObserver.
             * don't call it anywhere.
             * @type {IntersectionObserverCallback}
             */
            const megGalleryIntersectionObserverCallback = (entries, observer) => {
                console.debug(`In MediaExplorerGallery.megGalleryIntersectionObserverCallback: entires count = ${entries.length}`);

                /**
                 * @type {IntersectionObserverEntry[]}
                 */
                const visible_entries = [];

                entries.forEach(entry => {
                    if (!(entry.target instanceof HTMLElement)) return;

                    let event_name = meg_intersection_observer_event_names.VIEWPORT_LEAVE;

                    if (entry.isIntersecting) {
                        event_name = meg_intersection_observer_event_names.VIEWPORT_ENTER;
                        visible_entries.push(entry);
                    }

                    if (event_name !== '') {
                        const event = new CustomEvent(event_name);

                        entry.target.dispatchEvent(event);
                    }
                });

                if (visible_entries.length > 0) {
                    requestIdleCallback(() => {
                        registerItemsUnveiled(visible_entries)
                    });
                }
            }

            /**
             * Methods provided to gallery item. It adds a gallery item node to the observation list.
             * @type {import('./media_explorer_gallery').ObserveMEGalleryCallback}
             */
            const observeGalleryItemCallbackForItems = (gallery_item_node, ordered_media) => {
                if (gallery_item_intersection_observer == null) {
                    console.warn("In MediaExplorerGallery.observeGalleryItem: No intersection observer available to observe the gallery item.");
                    return;
                }  

                addMEGItemToObservationList(gallery_item_node, ordered_media);
            }

            /**
             * Methods provided to gallery items. It removes a gallery item node from the observation list.
             * @type {import('./media_explorer_gallery').UnobserveMEGalleryCallback}
             */
            const unobserveGalleryItemCallbackForItems = (gallery_item_node, ordered_media) => {
                if (gallery_item_intersection_observer == null) {
                    console.warn("In MediaExplorerGallery.unobserveGalleryItem: No intersection observer available to unobserve the gallery item.");
                    return;
                }

                removeMEGItemFromObservationList(gallery_item_node, ordered_media);
            }

            /**
             * Removes a gallery item node from the observation list.
             * @param {Element} gallery_item_node - The gallery item node to unobserve.
             * @param {import('@models/Medias').OrderedMedia} ordered_media
             * @returns {void}
             */
            const removeMEGItemFromObservationList = (gallery_item_node, ordered_media) => {
                if (gallery_item_intersection_observer == null) {
                    console.warn("In MediaExplorerGallery.removeMEGItemFromObservationList: No intersection observer available to unobserve the gallery item.");
                    return;
                }

                gallery_item_intersection_observer.unobserve(gallery_item_node);
            }

            /**
             * Registers an array of IntersectionObserverEntries to the last_orders_unveiled record.
             * @param {IntersectionObserverEntry[]} entries
             * @returns {void}
             */
            const registerItemsUnveiled = (entries) => {
                if (entries.length === 0) return;

                last_orders_unveiled.clear();

                for (let entry of entries) {
                    if (!(entry.target instanceof HTMLElement)) continue;

                    const entry_order = getMediaOrderOutOfElement(entry.target);
                    if (entry_order === null) continue;

                    last_orders_unveiled.add(entry_order);
                }
            }
        
        /*=====  End of Intersection observer  ======*/

        /*=============================================
        =            Performance Mode            =
        =============================================*/

            /**
             * Disables performance mode. This MOST not be called by anything else that is not 
             * performanceModeWatchdog.
             * @returns {void}
             */
            const disablePerformanceMode = () => {
                if (!enable_gallery_performance_mode) return;

                console.debug("In MediaExplorerGallery.disablePerformanceMode: Disabling performance mode.");

                enable_gallery_performance_mode = false;

                if (gallery_item_intersection_observer == null) {
                    initializeGalleryItemsIntesectionObserver();
                }
            }

            /**
             * Enables performance mode. This MOST not be called by anything else that is not
             * performanceModeWatchdog.
             * @returns {void}
             */
            const enablePerformanceMode = () => {
                if (enable_gallery_performance_mode) return;

                console.debug("In MediaExplorerGallery.enablePerformanceMode: Enabling performance mode.");

                enable_gallery_performance_mode = true;

                if (gallery_item_intersection_observer != null) {
                    dropGalleryItemsIntersectionObserver();
                }

                if (enable_gallery_heavy_rendering) {
                    enable_gallery_heavy_rendering = false;
                }
            }

            /**
             * Determines whether the gallery's performance mode should be enabled.
             * @returns {Promise<void>}
             */
            const performanceModeWatchdog = async () => {
                const OBSERVATION_ACTIVE_MEDIAS_THRESHOLD = 300;

                let new_performance_mode_value = false;

                if (active_medias.length >= OBSERVATION_ACTIVE_MEDIAS_THRESHOLD) {
                    new_performance_mode_value = true;
                }

                // Apply changes related to performance mode if it's values has changed.
                if (enable_gallery_performance_mode !== new_performance_mode_value) {
                    if (new_performance_mode_value) {
                        enablePerformanceMode();
                    } else {
                        disablePerformanceMode();
                    }
                }

                if (!regulating_active_medias_load) {
                    await tick();
                    regulateActiveMediasLoadExcess();
                }
            }
            
            /**
             * Regulates the active medias load excess. Ensures that the amount of medias displayed is 
             * managable. If it's found to be excessive, it will attempt to drop some of the medias keeping in mind 
             * continuity(not dropping medias from the middle series), preserving the focus media loaded and also the 
             * currently visible medias(the ones on the viewport), there last ones should contain the focused media, but
             * this may not be the case.
             * @returns {Promise<void>}
             */
            const regulateActiveMediasLoadExcess = async () => {
                // TODO: Remove debugMEG__attachSnapshot before merging to master.
                // TODO: Break down this function before merging to master.
                const FUZZY_MAX_MEDIA_LOAD = 250;
                if (active_medias.length <= FUZZY_MAX_MEDIA_LOAD) return;

                if (isNavigationGridOutdated()) {
                    await tick();
                    if (isNavigationGridOutdated()) {
                        await waitForGridSync(300); // Wait for the grid to sync or 300ms to pass.
                    }                    

                    if (isNavigationGridOutdated()) {
                        console.error(`In MediaExplorerGallery.regulateActiveMediasLoadExcess: The navigation grid is taking too long to update. Cowardly refusing to regulate active medias without accurate grid data.`);
                        return;
                    }
                }

                console.debug("%c"+"Regulating active medias", "color: orange; font-weight: bold; font-size: 46px");
                
                const pivot_media = getFocusedMedia();
                const last_media_added = last_media_added_to_active_medias;
                const grid_navigation_wrapper = the_grid_navigation_wrapper; // TS is retarded.

                const cannot_regulate_active_medias = pivot_media == null || last_media_added == null || grid_navigation_wrapper == null; 
                
                if (cannot_regulate_active_medias) {
                    console.debug("In MediaExplorerGallery.regulateActiveMediasLoadExcess: Cannot regulate active medias because one or more required variables are missing: pivot_media, last_media_added_to_active_medias, or the_grid_navigation_wrapper.");
                    return;
                }
                // debugMEG__attachSnapshot(regulateActiveMediasLoadExcess.name, {
                //     pivot_media: pivot_media,
                //     last_media_added: last_media_added,
                //     active_medias: [...active_medias],
                //     media_focus_index: media_focus_index,
                // });

                // We most not remove medias from the side they were added the last time.
                const addition_direction_right = last_media_added.Order > pivot_media.Order; 

                // NOTE: Only unmount full rows

                const ROW_OFFSET = 3; // The amount of rows to keep from the current row to the unmount direction.

                const current_grid_row_index = grid_navigation_wrapper.MovementController.Grid.CursorRow;
                const total_grid_rows = grid_navigation_wrapper.MovementController.Grid.length;

                const desired_grid_row_index = addition_direction_right ? current_grid_row_index - ROW_OFFSET : current_grid_row_index + ROW_OFFSET;

                if (desired_grid_row_index < 0 && desired_grid_row_index >= total_grid_rows) {
                    console.warn(`In MediaExplorerGallery.regulateActiveMediasLoadExcess: Could not offset the current grid row index by ${ROW_OFFSET} because it would go out of bounds.`);
                    console.debug(`Current grid row index: ${current_grid_row_index}\nTotal grid rows: ${total_grid_rows}\nDesired grid row index: ${desired_grid_row_index}`);
                    return;
                }
                // debugMEG__attachSnapshot(regulateActiveMediasLoadExcess.name, {
                //     addition_direction_right: addition_direction_right,
                //     current_grid_row_index: current_grid_row_index,
                //     total_grid_rows: total_grid_rows,
                //     desired_grid_row_index: desired_grid_row_index,
                // });

                const desired_row = grid_navigation_wrapper.MovementController.Grid.getRowAtIndex(desired_grid_row_index);
                if (desired_row == null) {
                    throw new Error(`In MediaExplorerGallery.regulateActiveMediasLoadExcess: Could not get the row at index ${desired_grid_row_index}. Even though last check indicated that it was in bounds.`);
                }

                let mount_pivot = addition_direction_right ? desired_row.MinIndex : 0;
                let mount_untill = addition_direction_right ? (active_medias.length - 1) : desired_row.MaxIndex;
                
                let regulated_active_medias = [];

                // debugMEG__attachSnapshot(regulateActiveMediasLoadExcess.name, {
                //     desired_row: desired_row,
                //     mount_pivot: mount_pivot,
                //     mount_untill: mount_untill,
                // });

                for (let h = mount_pivot; h <= mount_untill; h++) {
                    if (regulated_active_medias.length >= FUZZY_MAX_MEDIA_LOAD) {
                        console.error(`In MediaExplorerGallery.regulateActiveMediasLoadExcess: Logical error detected, the regulated active medias exceeded the maximum load of ${FUZZY_MAX_MEDIA_LOAD}.`);
                        break;
                    }

                    const media_item = active_medias[h];
                    if (media_item == null) continue;


                    regulated_active_medias.push(media_item);
                }

                regulating_active_medias_load = true;

                try {
                    setActiveMedias(regulated_active_medias);
                } catch {
                    console.error("In MediaExplorerGallery.regulateActiveMediasLoadExcess: Failed to set the active medias with the regulated ones.");
                    return;
                } finally {
                    regulating_active_medias_load = false;
                }

                // debugMEG__attachSnapshot(regulateActiveMediasLoadExcess.name, {
                //     regulated_active_medias: regulated_active_medias,
                //     active_medias_length: active_medias.length,
                // });
            
                await tick();

                await setMediaFocusIndex(pivot_media.Order, true);

                scrollCompensateToFocusedMedia(addition_direction_right);
            }
        
        /*=====  End of Performance Mode  ======*/
        
        /**
         * Adds an amount of media N items to the active_medias. it starts appending from an
         * index == active_medias[active_medias.length - 1].Order + 1 * if there are no new
         * medias to append, the function will return false otherwise true.
         * @param {number} amount 
         * @returns {boolean}
         */
        const appendMediaItems = (amount) => {
            if (media_items.length === active_medias.length) return false; // now that we add medias in both directions, this is not enough to determine if the medias append will overflow.

            const active_media_range = getLoadedMediaRange();
            if (active_media_range.end === media_items.length) return false;

            let start_index = active_media_range.end + 1;
            let end_index = Math.min(start_index + amount, media_items.length);

            console.debug(`Appending media items from ${start_index} to ${end_index}`)

            setActiveMedias([
                ...active_medias,
                ...ordered_medias.slice(start_index, end_index)
            ]);

            return true;
        }

        /**
         * Adds an initial batch of medias using sliceOrderedMedias with indexes 0
         * and Math.min(media_batch_size, ordered_medias.length)
         */
        const addInitialMediaItems = () => {
            const media_batch_size = getOptimalMediaBatchSize();
            
            setActiveMedias([], true);
            sliceOrderedMedias(0, Math.min(media_batch_size, ordered_medias.length));
        }

        /**
         * Defines the gallery state. It will initialize the gallery items intersection observer
         * and set the initial media items to be displayed in the gallery.
         */
        const defineGalleryState = async () => {
            if ($me_gallery_changes_manager === null) {
                me_gallery_changes_manager.set(new MediaChangesEmitter());
            }

            if (gallery_item_intersection_observer == null) {
                initializeGalleryItemsIntesectionObserver();
            }

            if (recovering_gallery_state) {
                // recoverGalleryFocusItem:
                // - Sets `media_focus_index` to the cached media index or the last viewed media index.
                // - Loads enough media items into `active_medias` to ensure the focused media is visible.
                // - Scrolls the focused media into view.
                // - Sets `recovering_gallery_state` to false after recovery is complete.
                await recoverGalleryFocusItem();
            } else {
                // addInitialMediaItems:
                // - Clears the `active_medias` array.
                // - Loads the first batch of media items (up to `media_batch_size`) into `active_medias`.
                addInitialMediaItems();
            }
        }

        /**
         * Disable Gallery hotkeys and emit an event so the parent component regain hotkeys control.
         * @param {boolean} quite - if true, will not emit an exit-hotkeys-context to component parent.
         */
        const disableGalleryHotkeys = (quite=false) => {
            enable_gallery_hotkeys = false;
            media_focus_index = 0;
            
            if (quite) return;
            
            dispatch("exit-hotkeys-context");
        }

        /**
         * Emits an event to close that should be handled by the parent component by closing the gallery.
         */
        const emitCloseGallery = () => {
            disableGalleryHotkeys(true);
            dispatch("close-gallery");
        }

        /**
         * Ensures the last row in the grid has as many items as it can. meaning that it will add
         * elements to make the last row full unless there are no more medias to add.
         * @returns {Promise<void>}
         */
        const ensureLastRowIsFull = async () => {
            if (the_grid_navigation_wrapper == null) {
                console.warn("In MediaExplorerGallery.ensureLastRowIsFull: No grid navigation wrapper available.");
                return;
            }

            const last_media_item = ordered_medias[ordered_medias.length - 1];

            isMediaOrderDisplayed(last_media_item.Order)

            const css_row_size = getMediaItemsPerRowFromCSS();
            if (css_row_size === undefined) {
                console.warn("In MediaExplorerGallery.handleMasonryLayoutChange: No CSS row size found.");
                return;
            }

            const last_gird_row = the_grid_navigation_wrapper.MovementController.Grid.getLastRow();
            if (last_gird_row === null) {
                console.warn("In MediaExplorerGallery.handleMasonryLayoutChange: No last grid row found.");
                return;
            }

            const missing_items = css_row_size - last_gird_row.length;

            if (missing_items <= 0) return;

            appendMediaItems(missing_items);
        }

        /**
         * Focuses the a media item by it's order. if the order is not in the active_medias, it drops the current
         * active medias an loads a batch that contains the media item. Assuming the media is inside the bounds of ordered_medias
         * otherwise it will throw an error.
         * @param {number} order
         * @returns {Promise<void>}
         */
        const focusMediaItemByOrder = async (order) => {
            const order_in_bounds = isMediaOrderInBounds(order);
            if (!order_in_bounds) {
                throw new Error(`Media order ${order} is out of bounds.`);
            }

            const order_displayed = isMediaOrderDisplayed(order);

            if (!order_displayed) {
                const containing_batch_loaded = await loadBatchWithMediaOrder(order);

                if (!containing_batch_loaded) {
                    throw new Error(`Media order ${order} is not in the active medias and could not be loaded.`);
                }               
            }

            media_focus_index = order;

            await tick();

            scrollToFocusedMedia("center", "instant");
        }

        /**
         * Returns the selected media or null if media_focus_index is out of bounds.
         * @returns {import('@models/Medias').OrderedMedia | null}
         */
        const getFocusedMedia = () => {
            /**
             * @type {import('@models/Medias').OrderedMedia | null}
             */
            let focused_media = ordered_medias[media_focus_index];

            if (focused_media === undefined) {
                console.warn(`No media found at index ${media_focus_index}`);
                focused_media = null;
            }

            return focused_media;
        }

        /**
         * Returns the Dom element that represents the focused media item in the gallery.
         * @returns {HTMLElement | null}
         */
        const getFocusedMediaElement = () => {
            return getMediaItemElementByOrder(media_focus_index);
        }

        /**
         * Returns a media by it's display index. meaning the index it has in the active_medias array.
         * @param {number} display_index
         * @returns {import('@models/Medias').OrderedMedia | null}
         */
        const getMediaByDisplayIndex = (display_index) => {
            if (display_index < 0 || display_index >= active_medias.length) {
                console.warn(`Display index ${display_index} is out of bounds for active_medias with length ${active_medias.length}`);
                return null;
            }

            return active_medias[display_index];
        }

        /**
         * Returns an ordered media by it's order.
         * @param {number} order
         * @returns {import('@models/Medias').OrderedMedia | null}
         */
        const getMediaByOrder = (order) => {
            const order_in_bounds = isMediaOrderInBounds(order);

            if (!order_in_bounds) {
                console.warn(`Media order ${order} is out of bounds for ordered_medias with length ${ordered_medias.length}`);
                return null;
            }

            return ordered_medias[order] || null;
        }

        /**
         * Returns the focused media item element in the gallery by it's order.
         * @param {number} order
         * @returns {HTMLElement | null}
         */
        const getMediaItemElementByOrder = (order) => {
            return document.querySelector(`.meg-gallery-item[data-media-order="${order}"]`);
        }

        /**
         * Returns the range of loaded medias in the active_medias array. meaning the {lowest order loaded, highest order loaded}
         * @typedef {Object} LoadedMediaRange
         * @property {number} start
         * @property {number} end
         * @return {LoadedMediaRange}
         */
        const getLoadedMediaRange = () => {
            if (active_medias.length === 0) {
                return {start: 0, end: 0};
            }

            let start = active_medias[0].Order;
            let end = active_medias[active_medias.length - 1].Order;

            return {start, end};
        }

        /**
         * Returns the index of the focused media in the active_medias(aka the currently displayed medias) array.
         * This is not the order of the media but it's sequential position on the gallery's grid.
         * @returns {number}
         */
        const getFocusedMediaDisplayIndex = () => {
            return getMediaDisplayIndexByOrder(media_focus_index);
        }

        /**
         * Returns the display index of a media item by it's order.
         * @param {number} order
         * @returns {number}
         */
        const getMediaDisplayIndexByOrder = (order) => {
            return active_medias.findIndex(media => media.Order === order);
        }

        /**
         * Extracts and returns a media order out of a an HTMLElement that is under the gallery container element.
         * @param {HTMLElement} element
         * @returns {number | null}
         */
        const getMediaOrderOutOfElement = element => {
            const media_order_data = element.dataset.mediaOrder;

            if (media_order_data === undefined) {
                return null;
            }

            const media_order = parseInt(media_order_data, 10);

            if (isNaN(media_order)) {
                console.error(`In MediaExplorerGallery.getMediaOrderOutOfElement: media_order_data "${media_order_data}" is not a valid number.`);
                return null;
            }

            return media_order;
        }

        /**
         * Handles the change between masonary and normal layout.
         * @returns {Promise<void>}
         */
        const handleMasonryLayoutChange = async () => {
            await updateNavigationGrid();

            await ensureLastRowIsFull();
        }

        /**
         * Returns the distance between to media orders.
         * @param {number} order_a
         * @param {number} order_b
         * @returns {number}
         */
        const getMediaOrderDistance = (order_a, order_b) => {
            if (!isMediaOrderInBounds(order_a) || !isMediaOrderInBounds(order_b)) {
                console.error(`In MediaExplorerGallery.getMediaOrderDistance: One of the orders is out of bounds. order_a: ${order_a}, order_b: ${order_b}`);
                return -1;
            }

            return Math.abs(order_a - order_b);
        }

        /**
         * Handles the click event on a media item.
         * @param {MouseEvent} event
         */
        const handleMediaItemClicked = (event) => {
            if (!(event.currentTarget instanceof HTMLElement)) return;

            let media_order_data = event.currentTarget.dataset.mediaOrder;

            if (media_order_data == null) {
                throw new Error("Media item clicked but had no data-media-order attribute.");
            }

            let media_order = parseInt(media_order_data);

            const active_media_range = getLoadedMediaRange();

            if (media_order == null || isNaN(media_order) || (media_order < 0 || media_order > active_media_range.end)) return;

            if (!event.altKey) {
                return openMediaViewerOnIndex(media_order);
            }

            media_focus_index = media_order;

            requestHotkeyControl();
        }

        /**
         * Hanldes the content end watch dog viewport enter event.
         */
        const handleContentEndWatchdogEnter = async () => {
            if (active_medias.length === 0) {
                addInitialMediaItems();
                return;
            }

            await loadProceedingMedias();

            manageContentEndWatchdogState();
        }

        /**
         * Handles the close event emitted from the Sequence Creation Tool.
         */
        const handleSequenceCreationToolClose = () => {
            enable_sequence_creation_tool = false;
        }

        /**
         * Returns whether the media order is in bounds.
         * @param {number} order
         * @returns {boolean}
         */
        const isMediaOrderInBounds = (order) => {
            return order >= 0 && order < media_items.length;
        }

        /**
         * Returns whether the media order is displayed. in other words,
         * whether the media in the given order is in the active_medias array.
         * @param {number} order
         */
        const isMediaOrderDisplayed = (order) => {
            if (active_medias.length === 0) {
                return false;
            }
            const active_media_range = getLoadedMediaRange();

            return order >= active_media_range.start && order <= active_media_range.end;
        }

        /**
         * Returns whether an array of ordered medias containes a different range of
         * medias from the one in the active_medias array.
         * @param {import('@models/Medias').OrderedMedia[]} other_medias
         * @returns {boolean}
         */
        const isActiveMediasDifferentFrom = (other_medias) => {
            const active_medias_length = active_medias.length;
            const other_medias_length = other_medias.length;

            if (active_medias_length !== other_medias_length) {
                return true;
            }

            const current_max_order = active_medias[active_medias_length - 1]?.Order || 0;
            const current_min_order = active_medias[0]?.Order || 0;
            const other_max_order = other_medias[other_medias_length - 1]?.Order || 0;
            const other_min_order = other_medias[0]?.Order || 0;

            return current_max_order !== other_max_order || current_min_order !== other_min_order;
        }

        /**
         * Loads medias preceding active_medias[0].Order. Returns a promise that resolves to true if more medias were loaded, false otherwise.
         * Throws an error if called when active_medias is empty as the propouse of this function is to fetch more medias for the infinite scroll.
         * @returns {Promise<boolean>}
         */
        const loadPrecedingMedias = async () => {
            if (active_medias.length === 0) {
                throw new Error("Attempted to load more medias when active_medias is empty. Use addInitialMediaItems or sliceOrderedMedias instead.");
            }
            const media_batch_size = getOptimalMediaBatchSize();

            const lowest_order_media = active_medias[0];
            
            if (lowest_order_media.Order === 0) return false; // nothing more can be loaded.

            
            let succesful_prepend = prependMediaItems(media_batch_size);
            
            if (!succesful_prepend) return false;
            

            await tick();
            
            const media_item_element = getMediaItemElementByOrder(lowest_order_media.Order);

            if (media_item_element != null) {
                // Prevent layout shift by scrolling to the first media item in the gallery.

                media_item_element.scrollIntoView({
                    block: "start",
                    behavior: "instant",
                });

                // setTimeout(() => {})

                media_item_element.scrollIntoView({
                    block: "center",
                    behavior: "smooth",
                }); // Simulate regular scrolling behavior.
            }
            
            return true;
        }

        /**
         * Loads proceeding after active_medias[active_medias.length - 1].Order. Returns a promise that resolves to true if more medias were loaded, false otherwise.
         * Throws an error if called when active_medias is empty as the propouse of this function is to fetch more medias for the infinite scroll.
         * @returns {Promise<boolean>}
         */
        const loadProceedingMedias = async () => {
            if (active_medias.length === 0) {
                throw new Error("Attempted to load more medias when active_medias is empty. Use addInitialMediaItems or sliceOrderedMedias instead.");
            }
            const media_batch_size = getOptimalMediaBatchSize();

            const highest_order_media_active = active_medias[active_medias.length - 1];

            if (highest_order_media_active.Order === media_items.length - 1) return false; // nothing more can be loaded.

            let succesful_append = appendMediaItems(media_batch_size);

            if (!succesful_append) return false;

            await tick();

            return true;
        }

        /**
         * Clears the active_medias and loads a new batch of ordered_medias that
         * contain the given ordered media by order. returns whether the
         * operation was successful or not.
         * @param {number} media_order
         * @returns {Promise<boolean>}
         */
        const loadBatchWithMediaOrder = async media_order => {
            if (media_order < 0 || media_order >= ordered_medias.length) {
                console.error(`In MediaExplorerGallery.loadBatchWithMediaOrder: media_order ${media_order} is out of bounds for ordered_medias with length ${ordered_medias.length}`);
                return false;
            }
            const media_batch_size = getOptimalMediaBatchSize();

            setActiveMedias([], true);

            await tick();

            let batches_needed = media_order > media_batch_size ? Math.ceil(media_order / media_batch_size) : 1;

            const container_batch_start_index = (batches_needed - 1) * media_batch_size;
            const container_batch_end_index = Math.min(container_batch_start_index + media_batch_size, media_items.length);

            sliceOrderedMedias(container_batch_start_index, container_batch_end_index);

            return true;
        }

        /**
         * called when the media_focus_index changes. Verifies more items should be loaded at the end or start of the active_medias. it does this by checking
         * if the media_focus_index is in the first or last row.
         * @returns {Promise<void>}
         */
        const manageInfiniteScroll = async () => {
            if (active_medias.length === 0) {
                console.warn("Refusing to manage infinite scroll when active_medias is empty.");
                return;
            }

            const media_per_row = getMediaItemsPerRow();
            if (media_per_row == null) {
                console.error("In MediaExplorerGallery.manageInfiniteScroll: No media per row data available.");
                return;
            }

            const active_media_range = getLoadedMediaRange();
            
            if (media_focus_index < active_media_range.start + media_per_row) {
                await loadPrecedingMedias();
            } else if (media_focus_index >= active_media_range.end - media_per_row) {
                await loadProceedingMedias();
            }

            manageContentEndWatchdogState();            
        }

        /**
         * Determines whether the content end watchdog should disabled.
         * @returns {boolean}
         */
        const manageContentEndWatchdogState = () => {
            let keep_watching = true;

            if (active_medias.length === media_items.length && active_medias.length === 0) {
                keep_watching = false;
            } 
            
            if (keep_watching) {
                let highest_order_media = active_medias[active_medias.length - 1];

                if (highest_order_media.Order === media_items.length - 1) {
                    keep_watching = false;
                }
            }

            has_proceeding_medias = keep_watching;

            return keep_watching;
        }

        /**
         * Manages state that has to be revaluated when the active_medias change. called by
         * setActiveMedias exclusively.
         * @param {import('@models/Medias').OrderedMedia[]} current_active_medias
         * @param {import('@models/Medias').OrderedMedia[]} new_active_medias
         */
        const manageActiveMediasState = (current_active_medias, new_active_medias) => {
            if (new_active_medias.length !== 0 && isActiveMediasDifferentFrom(new_active_medias)) {
                updateLastMediaAdded(new_active_medias, current_active_medias);
                
                handleChangesEmitterState();
            }
        }

        /**
         * Opens the media viewer in a given media index.
         * @param {number} media_index 
         */
        const openMediaViewerOnIndex = (media_index) => {
            const active_medias_range = getLoadedMediaRange();

            if (media_index < active_medias_range.start || media_index > active_medias_range.end) return;

            let media_item = ordered_medias[media_index];

            if (media_item == null) return;

            const page_state = $page.state;

            replaceState(location.href, {
                ...page_state,
                meg_gallery: {
                    media_index: media_index,
                }
            });

            dispatch("open-media-viewer", {
                media_item: media_item,
                media_index: media_index
            });
        }

        /**
         * Process current media changes in the gallery. in the media explorer gallery, move changes are used as a selection mechanism, we dont actually commit move changes.
         * so if there are any when this method is called, they will be unstaged.
         * @returns {Promise<void>}
         */
        const processMediaChanges = async () => {
            if ($current_category == null) {
                console.error("In MediaExplorerGallery.processMediaChanges: No current category available.");
                return;
            }

            
            if ($me_gallery_changes_manager === null) {
                console.warn("No changes manager available to process the current media changes.");
                return;
            }
    
            $me_gallery_changes_manager.clearAllMoveChanges();

            await $me_gallery_changes_manager.commitChanges($current_category.uuid);
        }

        /**
         * Prepends an amount of N media items to the active_medias. it starts prepending from an index == active_medias[0].Order - 1, if this index is less than 0, the function will return false.
         * This function will not work if active_medias is empty, for that scenario use addInitialMediaItems or sliceOrderedMedias.
         * @param {number} amount
         * @returns {boolean}
         */
        const prependMediaItems = (amount) => {
            if (media_items.length === active_medias.length) return false;

            const active_media_range = getLoadedMediaRange();

            if (active_media_range.start === 0) return false;

            let start_index = Math.max(active_media_range.start - amount, 0);
            let end_index = active_media_range.start;

            console.debug(`Prepending media items from ${start_index} to ${end_index}`);

            setActiveMedias([
                ...ordered_medias.slice(start_index, end_index),
                ...active_medias
            ]);

            return true;
        }

        /**
         * Resets the MEG gallery state.
         */
        const resetGalleryState = () => {
            if (global_hotkeys_manager == null) {
                console.error("In MediaExplorerGallery.resetGalleryState: No hotkeys manager available.");
                return;
            }
            console.debug("In MediaExplorerGallery.resetGalleryState: Resetting gallery state.");
            dropGridNavigationWrapper();

            dropGalleryItemsIntersectionObserver();
            
            media_focus_index = 0;
            global_hotkeys_manager.dropContext(hotkeys_context_name);

            me_gallery_changes_manager.set(null);
            me_renaming_focused_media.set(false);

        }

        /**
         * Attempts to recover the gallery focus from the previous session by reading category cache and loading enough medias to reach the focused media index.
         * Throws an error if called when active_medias is not empty.
         */
        const recoverGalleryFocusItem = async () => {
            if ($current_category == null) {
                console.error("In MediaExplorerGallery.recoverGalleryFocusItem: No current category available.");
                return;
            }

            if (category_cache == null) {
                console.error("In MediaExplorerGallery.recoverGalleryFocusItem: No category cache available.");
                return;
            }

            /** @type {number | undefined} */
            let cached_media_index = await category_cache.getCategoryIndex($current_category.uuid); 

            /**
             * @type {number | undefined}
             */
            // @ts-ignore - meg_gallery is defined in the page state.
            let page_state_store_index = $page.state?.meg_gallery?.media_index;

            cached_media_index = cached_media_index ?? page_state_store_index;

            if (cached_media_index == null) {
                return;
            }


            cached_media_index = Math.max(0, Math.min(cached_media_index, media_items.length - 1));

            if (cached_media_index >= media_items.length) {
                throw new Error(`Cached media index ${cached_media_index} is out of bounds for media_items with length ${media_items.length}`);
            }

            if (active_medias.length !== 0) {
                throw new Error("Attempted to recover the gallery focus item when active_medias is not empty.");
            }

            const batch_loaded = await loadBatchWithMediaOrder(cached_media_index);
           
            if (!batch_loaded) {
                console.error("Failed to load the batch with the cached media index.");
                return;
            }
            
            // media_focus_index = cached_media_index;
            setMediaFocusIndex(cached_media_index, true);
            
            await tick();

            let keyboard_selected_media = getMediaItemElementByOrder(media_focus_index);

            if (keyboard_selected_media != null) {
                keyboard_selected_media.scrollIntoView({
                    block: "center",
                    behavior: "instant"
                });
            }

            recovering_gallery_state = false;

            requestHotkeyControl();
        }

        /**
         * Requests hotkey control form the media explorer.
         * @returns {void}
         */
        const requestHotkeyControl = () => {
            if (!enable_gallery_hotkeys) {
                dispatch("request-hotkeys-control");
            }
        }

        /**
         * Scrolls the media focused item into view.
         * @param {"center" | "start" | "end"} [block="center"] - The block position to scroll to.
         * @param {"smooth" | "instant"} [behavior="smooth"] - The scroll behavior.
         * @returns {void}
         */
        const scrollToFocusedMedia = (block = "center", behavior = "smooth") => {
            const focused_media_element = getFocusedMediaElement();

            if (focused_media_element != null) {
                focused_media_element.scrollIntoView({
                    block,
                    behavior
                });
            }
        }

        /**
         * Scrolls to the focused media with a keyboard navigation 
         * simulation. This is useful to give the illusion of
         * continuous motion even when the cursor position is drastically
         * changed.
         * @param {boolean} [forward_motion=true] 
         * @returns {void}
         */
        const scrollCompensateToFocusedMedia = (forward_motion = true) => {
            const focused_media_element = getFocusedMediaElement();
            if (focused_media_element == null) {
                console.warn("In MediaExplorerGallery.scrollCompensateToFocusedMedia: No focused media element found");
                return;
            }

            focused_media_element.scrollIntoView({
                behavior: "instant",
                block: forward_motion ? "end" : "start",
            });

            focused_media_element.scrollIntoView({
                behavior: "smooth",
                block: "center",
            });
        }

        /**
         * Is called when a movement the overflows media_focus_index. And it determines whether an overflow to that
         * direction should change the hotkeys context. Tipically that would be because there is another component in that direction that can handle the same type of movement.
         * returns true if an action was taken, in that case the caller is expected to not take any further action, that includes updating any component property. if false, the caller is free to take whatever action it
         * deems adequate.
         * @param {boolean} up
         * @param {boolean} right
         * @param {boolean} down
         * @param {boolean} left
         * @returns {boolean}
         */
        const shouldMoveHotkeysContext = (up, right, down, left) => {
            if (up) {
                disableGalleryHotkeys();

                return true;
            }

            return false;
        }

        /**
         * Sets the media_focus_index to the given index. Remember that in the MEGallery component, when we talk
         * about media_focus_index, we are talking about the order of the media. Not its position in the 
         * active_medias array, meaning is not is sequential position in the displayed gallery UI. but in the current
         * content list(as ordered by the server).
         * @param {number} new_media_focus_index
         * @param {boolean} [update_cursor] Whether to update the cursor position in the_grid_navigation_wrapper, the update doesn't trigger the cursor position callback.
         * @returns {Promise<void>}
         */
        const setMediaFocusIndex = async (new_media_focus_index, update_cursor=false) => {
            const media_in_bounds = isMediaOrderInBounds(new_media_focus_index);

            if (!media_in_bounds) {
                throw new Error(`Media focus index ${new_media_focus_index} is out of bounds for media_items with length ${media_items.length}`);
            }

            if (new_media_focus_index !== media_focus_index) {
                const media_is_loaded = isMediaOrderDisplayed(new_media_focus_index);
    
                if (media_is_loaded) {
                    media_focus_index = new_media_focus_index;
                    
                    await manageInfiniteScroll();
                } else {
                    await focusMediaItemByOrder(new_media_focus_index);
                }
            }

            if (update_cursor && the_grid_navigation_wrapper != null) {
                const display_index = getMediaDisplayIndexByOrder(new_media_focus_index);

                the_grid_navigation_wrapper.updateCursorPositionSilently(display_index);
            }
        }

        /**
         * Sets the active_medias to the given array of ordered medias. It will make adjustments
         * to the galleries settings in order to ensure good performance unless its specified otherwise.
         * @param {import('@models/Medias').OrderedMedia[]} new_active_medias
         * @param {boolean} [skip_performance_check=false] - Use this flag for example if you are going to clear the array and immediately re-populate it.
         * @returns {void}
         */
        const setActiveMedias = (new_active_medias, skip_performance_check = false) => {
            manageActiveMediasState(active_medias, new_active_medias);
            
            active_medias = new_active_medias;

            if (!skip_performance_check) {
                performanceModeWatchdog();
            }
        }

        /**
         * If and only if the active_medias array is empty it adds a range/slice of ordered medias from ordered_medias to active_medias.
         * @param {number} start_index
         * @param {number} end_index
         */
        const sliceOrderedMedias = (start_index, end_index) => {
            if (active_medias.length > 0) {
                throw new Error("Active medias array is not empty but attempted to add. to keep consistency, you should either use appendMediaItems or prependMediaItems, both add medias from the start or end order of the active_medias array.");
            }

            if ((start_index < 0 || start_index > end_index) || (end_index > ordered_medias.length)) {
                throw new Error(`Invalid start_index(${start_index}) or end_index(${end_index}) for the ordered_medias with length ${ordered_medias.length}`);
            }

            setActiveMedias(ordered_medias.slice(start_index, end_index));
        }
        
        /*=============================================
        =            Search content feature            =
        =============================================*/
        
            /**
             * Focuses the given search result.
             * @param {import('@models/Medias').OrderedMedia} search_match
             * @returns {Promise<void>}
             */
            const focusMediaContentSearchMatch = async (search_match) => {
                const search_match_in_bounds = isMediaOrderInBounds(search_match.Order);
                if (!search_match_in_bounds) {
                    console.error(`In MediaExplorerGallery.focusMediaContentSearchMatch: search match ${search_match.Order} is out of bounds for ordered_medias with length ${ordered_medias.length}`);
                    return;
                }

                await setMediaFocusIndex(search_match.Order, true);

                toggleMediaTitlesMode(true);
            }

            /**
             * Handles the update of the search query label
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCaptureCallback}
             */
            const handleCategorySearchQueryUpdate = (event, captured_string) => {
                if (the_media_content_search_wrapper == null) return;

                capturing_media_content_search = true;
                media_content_search_query = captured_string;
            }

            /**
             * The search result handler bound to the_category_search_results_wrapper
             * @type {import('@common/keybinds/CommonActionWrappers').SearchResultsUpdateCallback<import('@models/Medias').OrderedMedia>}
             */
            const handleSearchMatchRecieved = search_match => {
                if (the_media_content_search_wrapper == null || search_match == null) return;

                if (!(search_match instanceof OrderedMedia)) {
                    console.error("search_match is not an OrderedMedia:", search_match);
                    throw new Error("The displayed categories should only contain InnerCategory instances")
                }

                const search_match_in_bounds = isMediaOrderInBounds(search_match.Order);

                if (!search_match_in_bounds) {
                    console.error(`In MediaExplorerGallery.handleSearchMatchRecieved: search match ${search_match.Order} is out of bounds for ordered_medias with length ${ordered_medias.length}`);
                    return;
                }

                capturing_media_content_search = false;

                if (the_media_content_search_wrapper.SearchResults.length > 0) {
                    media_content_search_results_lookup = new Set(the_media_content_search_wrapper.SearchResults.map((result => result.uuid)));
                }

                focusMediaContentSearchMatch(search_match);
            }

            /**
             * @type {import('@common/keybinds/CommonActionWrappers').SearchResultsWrapperSetup}
             */
            const setSearchResultsWrapper = hotkey_contenxt => {
                if ($current_category == null) {
                    console.error("In MediaExplorerGallery.setSearchResultsWrapper , current_category is null");
                    return;
                }

                the_media_content_search_wrapper = new SearchResultsWrapper(hotkey_contenxt, ordered_medias, handleSearchMatchRecieved, {
                    minimum_similarity: 0.7,
                    search_hotkey: ["f"],
                    ui_search_result_reference: ui_core_dungeon_references.MEDIA,
                    search_typing_hotkey_handler: handleCategorySearchQueryUpdate,
                    boost_exact_inclusion: true,
                    allow_member_similarity_checking: true,
                    no_results_callback: () => resetCategoryFiltersState()
                });

                the_media_content_search_wrapper.setItemToStringFunction(ordered_medias => ordered_medias.MediaName.toLowerCase());
            }          

            /**
             * Resets the state of the category name filter
             * @returns {void}
             */
            const resetCategoryFiltersState = () => {
                if (global_hotkeys_manager == null) return;

                media_content_search_query = "";
                capturing_media_content_search = false;
                media_content_search_results_lookup = null;
            }

        /*=====  End of Search content feature  ======*/

        /**
         * Toggles the renaming focused media state.
         * @returns {void}
         */
        const toggleRenamingFocusedMediaState = () => {
            me_renaming_focused_media.set(!$me_renaming_focused_media);
        }
        
        /**
         * toggles the media titles mode.
         * @param {boolean} [force_state] - ensures the desired state is set instead of toggling it.
         * @returns {void}
         */
        const toggleMediaTitlesMode = (force_state) => {
            const new_state = force_state ?? !show_media_titles_mode;

            show_media_titles_mode = new_state;
        }

        /**
         * Updates the last_media_added_to_active_medias, recieves the new and old values of active medias. 
         * Determines if the change occurred at the start or end of the active_medias. If it can't determine,
         * sets the last_media_added_to_active_medias to undefined.
         * @param {import('@models/Medias').OrderedMedia[]} new_active_medias
         * @param {import('@models/Medias').OrderedMedia[]} old_active_medias
         * @returns {void}
         */
        const updateLastMediaAdded = (new_active_medias, old_active_medias) => {
            /**
             *  @type {import('@models/Medias').OrderedMedia | undefined}
             */
            let updated_last_added_media = undefined;

            
            if (new_active_medias.length > 0) {
                // Keep in mind: if arr = [1,2,3,4], arr[4] is undefined(no overflow error thrown in js). so if
                // for example, both arrays are empty, all of these variables will be undefined, which is 
                // the default value of last_media_added_to_active_medias.
                const last_media_in_new = new_active_medias[new_active_medias.length - 1];
                const first_media_in_new = new_active_medias[0];
                const last_media_in_old = old_active_medias[old_active_medias.length - 1];
                const first_media_in_old = old_active_medias[0];

                if (old_active_medias.length === 0) {
                    updated_last_added_media = last_media_in_new;
                } else if (first_media_in_old.Order === first_media_in_new.Order && last_media_in_old.Order !== last_media_in_new.Order) {
                    // New media was added to the end of the active medias.
                    updated_last_added_media = last_media_in_new;
                } else if (last_media_in_old.Order === last_media_in_new.Order && first_media_in_old.Order !== first_media_in_new.Order) {
                    // New media was added to the start of the active medias.
                    updated_last_added_media = first_media_in_new;
                }
            }

            last_media_added_to_active_medias = updated_last_added_media;
        }

    /*=====  End of Methods  ======*/
</script>

<div id="meg-gallery-wrapper">
    {#if enable_sequence_creation_tool && enable_gallery_hotkeys}
        <SequenceCreationTool 
            unsequenced_medias={ordered_medias}
            on:close-sct={handleSequenceCreationToolClose}
        />
    {:else}
        <div id="meg-gallery-header-wrapper">
            <header id="meg-gallery-header">
                {#if $current_category != null}
                    <h2 id="meg-gh-title">
                        {ui_core_dungeon_references.CATEGORY.EntityName} gallery for
                        <span class="meg-gh-cn-category-name">{$current_category.name}</span>     
                    </h2>
                    <p class="meg-gh-paragraph">―</p>
                    <p class="meg-gh-paragraph">
                        {$current_category.content.length} total {ui_core_dungeon_references.MEDIA.EntityNamePlural}
                    </p>
                {/if}
            </header>
        </div>
        <div id="meg-gw-floating-controls-overlay">
            {#if capturing_media_content_search || media_content_search_query}
                <p id="meg-gw-media-content-filter">
                    /{media_content_search_query}
                </p>
            {/if}
        </div>
        <ul id="meg-gallery" class:masonry-layout={use_masonry}>
            {#if $me_gallery_changes_manager instanceof MediaChangesEmitter}
                {#each active_medias as ordered_media_item, h}
                    {@const is_media_keyboard_selected = enable_gallery_hotkeys && media_focus_index === ordered_media_item.Order}
                    {@const masonary_random_number = Math.random()}
                    <li class="meg-gallery-item {media_item_html_class}" 
                        class:keyboard-selected={is_media_keyboard_selected} 
                        class:masonry-item--small={use_masonry && (masonary_random_number < 0.3) }
                        class:masonry-item--medium={use_masonry && (masonary_random_number >= 0.3 && masonary_random_number < 0.6) }
                        class:masonry-item--large={use_masonry && (masonary_random_number >= 0.6) }
                        class:is-skeleton={recovering_gallery_state}
                        data-media-order={ordered_media_item.Order}
                        data-media-index={h}
                        on:click={handleMediaItemClicked}
                    >
                        <MeGalleryDisplayItem
                            ordered_media={ordered_media_item}
                            media_viewport_percentage={0.3}
                            is_keyboard_focused={is_media_keyboard_selected}
                            is_skeleton={recovering_gallery_state}
                            enable_magnify_on_keyboard_focus={magnify_focused_media}
                            check_container_limits
                            container_selector="#meg-gallery"
                            enable_heavy_rendering={enable_gallery_heavy_rendering}
                            enable_media_titles={show_media_titles_mode}
                            requestObserveMediaItem={observeGalleryItemCallbackForItems}
                            requestUnobserveMediaItem={unobserveGalleryItemCallbackForItems}
                            {use_masonry}
                        />
                        <CoverSlide 
                            component_id="meg-gallery-cover-slide-{h}"
                            ignore_parent_events={is_media_keyboard_selected || enable_gallery_hotkeys}
                        />
                    </li>        
                {/each}
            {/if}
        </ul>
        {#if has_proceeding_medias}
            <div class="meg-gallery-end-of-content-watchdog"
                on:viewportEnter={handleContentEndWatchdogEnter}
                use:viewport
            >
                <GridLoader />
            </div>
        {/if}
    {/if}
</div>

<style>

    #meg-gallery-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-2);
    }
    
    /*=============================================
    =            Gallery Header            =
    =============================================*/
    
        header#meg-gallery-header {
            display: flex;
            width: 100%;
            height: 80px;
            align-items: center;
            padding-inline: calc(var(--common-page-inline-padding) * 1.5);
            column-gap: var(--spacing-3);
            
            & h2, p {
                line-height: 1;
            }

            & h2#meg-gh-title {
                font-family: var(--font-read);
                color: var(--main);
                color: var(--grey-2);
                font-size: var(--font-size-h3);
            }


            & h2#meg-gh-title::first-letter {
                text-transform: uppercase;
            }

            & h2 span.meg-gh-cn-category-name {
                color: var(--main);
                font-size: calc(var(--font-size-h3) * 1.1);
                font-style: italic;
                font-weight: 600;
            }
        } 
    
    /*=====  End of Gallery Header  ======*/
    
   
    /*=============================================
    =            Control overlays            =
    =============================================*/
   
        :has(> #meg-gw-floating-controls-overlay) {
            position: relative;
        }
        
        #meg-gw-floating-controls-overlay {
            position: fixed;
            bottom: 0;
            left: 0;
            width: 100%;
            padding: var(--spacing-1);
            z-index: var(--z-index-t-7);

            & p#meg-gw-media-content-filter {
                padding-inline-start: var(--spacing-2);
                font-size: var(--font-size-2);
                color: var(--grey-1);
                line-height: 1.8;
                background: hsl(from var(--grey-9) h s l / 0.3);
                font-weight: 500;
            }
        }
   
    /*=====  End of Control overlays  ======*/
   
    

    #meg-gallery {
        display: grid;
        /* CRITICAL: grid-template columns most not change in different layouts, like the masonry.
         This is so loaded batch sizes, specially prepended ones, don't shift the focused media horizontally.
         It's not ideal but it is preferable to having this disorientating behavior. 
        */
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); 
        background: var(--grey);
        gap: var(--spacing-1);
        padding-block: 0;
        padding-inline: var(--spacing-2);
        list-style: none;
        /* padding: 0; */
        margin: 0;
    }    

    li.meg-gallery-item {
        position: relative;
        cursor: pointer;
        container-type: inline-size;
        background: var(--grey-9);
        width: 100%;
        height: 100%;
        z-index: var(--z-index-1);
    }

    li.meg-gallery-item.is-skeleton {
        height: 340px !important;
    }

    li.meg-gallery-item.keyboard-selected {
        outline: var(--main) solid 2px;
        z-index: var(--z-index-2);
    }

    #meg-gallery-wrapper .meg-gallery-end-of-content-watchdog {
        width: 100%;
        display: grid;        
        padding: var(--vspacing-3) 0;
        place-items: center;
    }
    
    /*=============================================
    =            Masonry            =
    =============================================*/
    
        #meg-gallery-wrapper:has(> .masnory-layout) {
            display: flex;
            flex-direction: column;
            width: 100%;
            align-items: center;
        }
        
        #meg-gallery.masonry-layout{
            --mesonry-item-small: 11;   
            --mesonry-item-medium: 15;
            --mesonry-item-large: 19;

            container-type: inline-size;
            width: 100%;
            /* grid-template-columns: repeat(auto-fill, 15%); */
            grid-auto-rows: 1.2vh;
            justify-content: center;
            gap: 12px;

            & .meg-gallery-item {
                break-inside: avoid;
            }

            & .meg-gallery-item.masonry-item--small {
                grid-row-end: span var(--mesonry-item-small);
            }

            & .meg-gallery-item.masonry-item--medium {
                grid-row-end: span var(--mesonry-item-medium); 
            }

            & .meg-gallery-item.masonry-item--large {
                grid-row-end: span var(--mesonry-item-large);
            }
        }
    
    /*=====  End of Masonry  ======*/
    
    
</style>