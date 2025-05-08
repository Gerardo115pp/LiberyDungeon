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
        import { me_gallery_changes_manager, me_gallery_yanked_medias, me_renaming_focused_media } from './me_gallery_state';
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
         * the amount of medias to display render at a time. If the end of the viewport is reached,
         * a new batch of exactly this amount(or the remaining available) of medias will be appended to the gallery.
         * @type {number}
         */
        export let media_batch_size = 30;

        /**
         * an html class that is added to all the media items in the gallery.
         * @type {string}
         */
        const media_item_html_class = "meg-gallery-grid-cell";
       
        /*----------  State  ----------*/
        
            /**
             * The medias that ARE been displayed in the gallery. 
             * @type {import('@models/Medias').OrderedMedia[]}
             */
            export let active_medias = [];

            /** 
             * Media selected index
             * @type {number}
             */
            let media_focus_index = 0;

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
             * Whether to enable heavy rendering of the medias in the gallery. This will force the media item in the gallery to render in higher quality and
             * if they are videos, those will be loaded instead of their thumbnails(which is the default behavior) they will however only load when they are on 
             * the viewport.
             * @type {boolean}
             */
            export let enable_gallery_heavy_rendering = false;

            /**
             * Whether to enable the sequence creation tool.
             * @type {boolean}
             */
            let enable_sequence_creation_tool = false;

            
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
        

        let dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        if (browser) {
            // @ts-ignore - globalThis is not defined in the but is always available despite the environment.
            globalThis.meg_gallery_page = $page;
        }

        defineComponentHotkeys();
        defineGalleryState();
    });

    onDestroy(() => {
        if (!browser) return;
        resetGalleryState();
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

                    hotkeys_context.register(["w", "a", "s", "d"], handleGalleryGridMovement, {
                        description: `<navigation> Move the focus on the gallery grid.`,
                    });

                    hotkeys_context.register(["shift+s"], handleMovetoLastActiveMedia, {
                        description: `<navigation> Move the focus to the last loaded media in the gallery.`,
                    });

                    hotkeys_context.register(["shift+w"], handleMovetoFirstActiveMedia, {
                        description: `<navigation> Move the focus to the first loaded media in the gallery.`,
                    });

                    hotkeys_context.register(["q"], handleGalleryExit, {
                        description: "<navigation>Sets focus on the category section without closing the gallery.",
                    });

                    hotkeys_context.register(["g"], handleGalleryClose, {
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

                    hotkeys_context.register(["space"], handleMediaSelectMode, {
                        description: "<content>Selects the focused media.",
                        mode: "keydown"
                    });

                    hotkeys_context.register(["space"], handleMediaSelectMode, {
                        description: "<content>Selects the focused media.",
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

                    global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
                }

                if (enable_gallery_hotkeys) {
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
             */
            const handleGalleryGridMovement = async (hotkey_event, hotkey) => {
                const media_count = active_medias.length;
                const active_media_range = getLoadedMediaRange();
                const media_per_row = Math.floor(getMediaItemsPerRow());
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
                    selectFocusedMedia();
                } else if (auto_stage_delete_focused_media) {
                    let delete_range = (hotkey.KeyCombo === "w" || hotkey.KeyCombo === "s") && Math.abs(old_focus_index - new_focus_index) <= media_per_row + 1;

                    if (delete_range) {
                        let start_index = Math.min(old_focus_index, new_focus_index);
                        let end_index = Math.max(old_focus_index, new_focus_index);

                        stageMediaRangeDeletion(start_index, end_index);
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
                hotkey_event.preventDefault();
                if ($me_gallery_changes_manager === null) {
                    console.warn("No changes manager available to stage the deletion of the focused media item.");
                    return;
                }
                
                if (hotkey_event.repeat || hotkey_event.type !== "keydown" && hotkey_event.type !== "keyup") return;

                let new_auto_stage_delete_focused_media = hotkey_event.type === "keydown";

                if (new_auto_stage_delete_focused_media) {
                    stageFocusedMediaDeletion();
                }

                auto_stage_delete_focused_media = new_auto_stage_delete_focused_media;
            }

            /**
             * Handles the selection of the focused media item.
             * @param {KeyboardEvent} hotkey_event
             */
            const handleMediaSelectMode = hotkey_event => {
                hotkey_event.preventDefault();
                if (hotkey_event.type !== "keydown" && hotkey_event.type !== "keyup") return;
                if (hotkey_event.repeat) return;
                
                auto_select_focused_media = hotkey_event.type === "keydown";

                if (auto_select_focused_media) {
                    selectFocusedMedia();
                }
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

                const medias_per_row = Math.floor(getMediaItemsPerRow());

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
                enable_gallery_heavy_rendering = !enable_gallery_heavy_rendering;
            }

        /*=====  End of Hotkeys  ======*/
    
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

            active_medias = [
                ...active_medias,
                ...ordered_medias.slice(start_index, end_index)
            ];

            return true;
        }

        /**
         * Adds an initial batch of medias using sliceOrderedMedias with indexes 0
         * and Math.min(media_batch_size, ordered_medias.length)
         */
        const addInitialMediaItems = () => {
            active_medias = [];
            sliceOrderedMedias(0, Math.min(media_batch_size, ordered_medias.length));
        }

        const defineGalleryState = () => {
            if ($me_gallery_changes_manager === null) {
                me_gallery_changes_manager.set(new MediaChangesEmitter());
            }

            if (recovering_gallery_state) {
                // recoverGalleryFocusItem:
                // - Sets `media_focus_index` to the cached media index or the last viewed media index.
                // - Loads enough media items into `active_medias` to ensure the focused media is visible.
                // - Scrolls the focused media into view.
                // - Sets `recovering_gallery_state` to false after recovery is complete.
                recoverGalleryFocusItem();
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

            const media_item_element = getMediaItemElementByOrder(order);

            if (media_item_element != null) {
                media_item_element.scrollIntoView({
                    block: "center",
                    behavior: "instant"
                });
            }
        }

        /**
         * Returns the amount of medias that fit in a single row of the gallery's grid.
         * @returns {number} - a float number most likely.
         */
        const getMediaItemsPerRow = () => {
            const meg_gallery_container = document.querySelector('#meg-gallery');
            const meg_gallery_item = document.querySelector(`.${media_item_html_class}`);

            if (meg_gallery_container == null || meg_gallery_item == null){
                if (meg_gallery_container == null) console.error("#meg-gallery not found");
                if (meg_gallery_item == null) console.error(`.${media_item_html_class} not found`);
                return 0;
            };

            const container_style = window.getComputedStyle(meg_gallery_container);

            const container_padding_left = parseFloat(container_style.paddingLeft);
            const container_padding_right = parseFloat(container_style.paddingRight);

            return (meg_gallery_container.clientWidth - (container_padding_left + container_padding_right)) / meg_gallery_item.clientWidth;
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
         * Returns the focused media item element in the gallery by it's order.
         * @param {number} order
         * @returns {HTMLElement | null}
         */
        const getMediaItemElementByOrder = (order) => {
            return document.querySelector(`.meg-gallery-item[data-media-order="${order}"]`);
        }

        /**
         * Returns the range of medias exising in the avtive_medias slice where start <= 0 and end >= active_medias.length.
         * start and end are garanteed to exist withing media_items bounds. 
         * @returns {LoadedMediaRange}
         * @typedef {Object} LoadedMediaRange
         * @property {number} start
         * @property {number} end
         */
        const getLoadedMediaRange = () => {
            let start = active_medias[0].Order;
            let end = active_medias[active_medias.length - 1].Order;

            return {start, end};
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
            const active_media_range = getLoadedMediaRange();

            return order >= active_media_range.start && order <= active_media_range.end;
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

            const highest_order_media = active_medias[active_medias.length - 1];

            if (highest_order_media.Order === media_items.length - 1) return false; // nothing more can be loaded.

            let succesful_append = appendMediaItems(media_batch_size);

            if (!succesful_append) return false;

            await tick();

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

            active_medias = [
                ...ordered_medias.slice(start_index, end_index),
                ...active_medias
            ];

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
            
            media_focus_index = cached_media_index;
            
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
         * Sets the current media as selected. if the media is already selected, it will be unselected.
         * @returns {void}
         */
        const selectFocusedMedia = () => {
            if ($me_gallery_changes_manager === null) {
                console.warn("No changes manager available to stage the selection of the focused media item.");
                return;
            }

            let focused_media = getFocusedMedia();
            if (focused_media === null) return;

            toggleMediaSelect(focused_media);
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

            active_medias = ordered_medias.slice(start_index, end_index);
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

                await focusMediaItemByOrder(search_match.Order);
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
                    boost_exact_inclusion: false,
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

            let is_delete_already_staged = $me_gallery_changes_manager.getMediaChangeType(focused_media.uuid) === media_change_types.DELETED;

            if (is_delete_already_staged && auto_stage_delete_focused_media) return;

            toggleMediaDeletion(focused_media);
        }

        /**
         * Stages the deletion of the given media inclusive range. so if 4 - 10 are passed, medias on positions 4 to 10 will be staged for deletion.
         * @param {number} start_index
         * @param {number} end_index
         */
        const stageMediaRangeDeletion = (start_index, end_index) => {
            if ($me_gallery_changes_manager === null) {
                console.warn("No changes manager available to stage the deletion of the media range.");
                return;
            }

            let media_range = ordered_medias.slice(start_index, end_index + 1);

            for (let media_item of media_range) {
                toggleMediaDeletion(media_item, true);
            }
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

            active_medias = [];

            await tick();

            let batches_needed = media_order > media_batch_size ? Math.ceil(media_order / media_batch_size) : 1;

            const container_batch_start_index = (batches_needed - 1) * media_batch_size;
            const container_batch_end_index = Math.min(container_batch_start_index + media_batch_size, media_items.length);

            sliceOrderedMedias(container_batch_start_index, container_batch_end_index);

            return true;
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
         * Toggles the select stat of a given media.
         * @param {import('@models/Medias').OrderedMedia} media_item 
         */
        const toggleMediaSelect = (media_item) => {
            if ($me_gallery_changes_manager === null) {
                console.warn("In MediaExplorerGallery.toggleMediaSelect: No changes manager available to stage the selection of the media item.");
                return;
            }

            let current_media_change = $me_gallery_changes_manager.getMediaChangeType(media_item.uuid);

            if (current_media_change === media_change_types.MOVED) {
                $me_gallery_changes_manager.unstageMediaMove(media_item.uuid);
                return;
            }

            // on this meg gallery, we don't actually commit move media changes, only delete changes. instead of committing the move changes, we used them as
            // a select mechanism. if there are any move changes when committing, we will unstage them first.
            let fake_inner_category = new InnerCategory({
                name: "me-gallery-mock-category", 
                uuid: "me-gallery-mock-category", 
                fullpath: "me-gallery-mock-category",
                category_thumbnail: "me-gallery-mock-category",
            });

            $me_gallery_changes_manager.stageMediaMove(media_item.Media, fake_inner_category);
        }

        /**
         * Toggles the deletion state of a given media.
         * @param {import('@models/Medias').OrderedMedia} ordered_media 
         * @param {boolean} [keep_deleted=true] - if true and the media is already staged for deletion, it will not be unstaged.
         */
        const toggleMediaDeletion = (ordered_media, keep_deleted=false) => {
            if ($me_gallery_changes_manager === null) {
                console.warn("In MediaExplorerGallery.toggleMediaDeletion: No changes manager available to stage the deletion of the media item.");
                return;
            }

            let current_media_change = $me_gallery_changes_manager.getMediaChangeType(ordered_media.uuid);

            if (current_media_change === media_change_types.DELETED && !keep_deleted) {
                $me_gallery_changes_manager.unstageMediaDeletion(ordered_media.uuid);
                return;
            }

            $me_gallery_changes_manager.stageMediaDeletion(ordered_media.Media);
        }

        /**
         * Toggles the renaming focused media state.
         * @returns {void}
         */
        const toggleRenamingFocusedMediaState = () => {
            me_renaming_focused_media.set(!$me_renaming_focused_media);
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
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        background: var(--grey);
        gap: 2px;
        padding: 4px;
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
            grid-template-columns: repeat(auto-fill, 15%);
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