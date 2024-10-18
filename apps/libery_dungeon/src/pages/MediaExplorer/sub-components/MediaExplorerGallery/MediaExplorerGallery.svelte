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
        import { me_gallery_changes_manager, me_gallery_yanked_medias } from './me_gallery_state';
        import MeGalleryDisplayItem from './MEGalleryDisplayItem.svelte';
        import GridLoader from '@components/UI/Loaders/GridLoader.svelte';
        import CoverSlide from '@components/Animations/HoverEffects/CoverSlide.svelte';
        import { browser } from '$app/environment';
        import { media_change_types, MediaChangesEmitter } from '@models/WorkManagers';
        import { category_cache, InnerCategory } from '@models/Categories';
        import { LabeledError } from '@libs/LiberyFeedback/lf_models';
        import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
        import { emitPlatformMessage } from '@libs/LiberyFeedback/lf_utils';
        import { pushState, replaceState } from '$app/navigation';
        import { page } from '$app/stores';
        import { OrderedMedia } from '@models/Medias';
        import { active_media_index } from '@stores/media_viewer';
    import SequenceCreationTool from './SequenceCreationTool.svelte';
    
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
         * The media items to be displayed
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
         * the amount of medias to display render at a time. If the end of the viewport is reached, a new batch of exactly this amount of medias will be appended to the gallery.
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
             * Active medias that will be displayed in the gallery
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
             * Whether focused media items should be auto selected. Enabled on keydown for media select key(originally space) and disabled on key up.
             * @type {boolean}
             */
            let auto_select_focused_media = false;

            /**
             * Whether focused media should automatically staged
             * for deletion.
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
        
        /*----------  Masonry  ----------*/

            /**
             * Whether to use a masonry layout for the gallery.
            */
            export let use_masonry = false;
        

        let dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        if (browser) {
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

                    hotkeys_context.register(["n"], handleGalleryReset, {
                        description: "<content>Deselects all selected medias and restores deleted medias that have not been committed(they are committed when you close the gallery).",
                    });

                    hotkeys_context.register(["c s"], handleEnableSequenceCreationTool, {
                        description: "<tools>Enable the sequence creation tool.",
                    });

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
             * Toggles the display of the media titles in the gallery items.
             */
            const handleShowTitlesMode = () => {
                show_media_titles_mode = !show_media_titles_mode;
            }

            /**
             * Closes the gallery and exits the hotkeys context.
             */
            const handleGalleryClose = async () => {
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
                console.log(`Opening media viewer on media index ${media_focus_index}`);
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
                
                if (!hotkey_event.repeat) {
                    console.log('event:', hotkey_event);    
                }

                if (hotkey_event.repeat || hotkey_event.type !== "keydown" && hotkey_event.type !== "keyup") return;

                let new_auto_stage_delete_focused_media = hotkey_event.type === "keydown";
                console.log(`auto_stage_delete_focused_media: ${new_auto_stage_delete_focused_media}`);

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
                console.log(`Event type: ${hotkey_event.type}`);
                
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
            const handleGalleryReset = () => {
                if ($me_gallery_changes_manager === null) return;

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
         * Adds an amount of media N items to the active_medias. it start appending from an index == active_medias[active_medias.length - 1].Order + 1
         * if there are no new medias to append, the function will return false otherwise true.
         * @param {number} amount 
         * @returns {boolean}
         */
        const appendMediaItems = (amount) => {
            if (media_items.length === active_medias.length) return false; // now that we add medias in both directions, this is not enough to determine if the medias append will overflow.

            const active_media_range = getLoadedMediaRange();
            if (active_media_range.end === media_items.length) return false;

            console.log(`active_media_range:`, active_media_range);

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
         * Adds an initial batch of medias using sliceOrderedMedias with indexes 0 and Math.min(media_batch_size, ordered_medias.length)
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
                recoverGalleryFocusItem();
            } else {
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
        const getMediaItemByOrder = (order) => {
            return document.querySelector(`.meg-gallery-item[data-media-order="${order}"]`);
        }

        /**
         * Returns the range of medias exising in the avtive_medias slice where start <= 0 and end >= active_medias.length. start and end are garanteed to exist withing media_items bounds. 
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
            let media_order_data = event.currentTarget.dataset.mediaOrder;
            let media_order = parseInt(media_order_data);

            const active_media_range = getLoadedMediaRange();

            console.log(`Media item clicked with order ${media_order}`);

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
            console.log("Content end watchdog enter event.");
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

            console.log(`Lowest order media:`, lowest_order_media);
            
            let succesful_prepend = prependMediaItems(media_batch_size);
            console.log(`Prepending medias result ${succesful_prepend}`);
            
            if (!succesful_prepend) return false;
            

            await tick();
            
            const media_item_element = getMediaItemByOrder(lowest_order_media.Order);
            console.log(`Scrolling to media item with order ${lowest_order_media.Order}`, media_item_element);

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
                console.log(`Prepending medias as Media focus index ${media_focus_index} is in the first row.`);
                await loadPrecedingMedias();
            } else if (media_focus_index >= active_media_range.end - media_per_row) {
                console.log(`Appending medias as Media focus index ${media_focus_index} is in the last row.`);
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

            console.log(`storing media_index: ${media_index} in the page state.`);

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
         * @returns {promise<void>}
         */
        const processMediaChanges = async () => {
            if ($me_gallery_changes_manager === null) {
                console.warn("No changes manager available to process the current media changes.");
                return;
            }
    
            $me_gallery_changes_manager.clearAllMoveChanges();

            await $me_gallery_changes_manager.commitChanges($current_category.uuid, $current_cluster.UUID);
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
            media_focus_index = 0;
            global_hotkeys_manager.dropContext(hotkeys_context_name);

            me_gallery_changes_manager.set(null);
        }

        /**
         * Attempts to recover the gallery focus from the previous session by reading category cache and loading enough medias to reach the focused media index.
         * Throws an error if called when active_medias is not empty.
         */
        const recoverGalleryFocusItem = async () => {
            /** @type {number} */
            let cached_media_index = await category_cache.getCategoryIndex($current_category.uuid); 
            console.log("Indexed db cached media index:", cached_media_index);
            let page_state_store_index = $page.state?.meg_gallery?.media_index;
            console.log("History API cached media index:", page_state_store_index);

            cached_media_index = cached_media_index ?? page_state_store_index;

            if (cached_media_index == null) {
                console.log("No cached media index found in the page state.");
            }

            console.log(`kept media index: ${cached_media_index}`);

            cached_media_index = Math.max(0, Math.min(cached_media_index, media_items.length - 1));

            if (cached_media_index >= media_items.length) {
                throw new Error(`Cached media index ${cached_media_index} is out of bounds for media_items with length ${media_items.length}`);
            }

            if (active_medias.length !== 0) {
                throw new Error("Attempted to recover the gallery focus item when active_medias is not empty.");
            }

            console.log(`Focused media index: ${cached_media_index}`);

            let batches_needed = cached_media_index > media_batch_size ? Math.ceil(cached_media_index / media_batch_size) : 1;



            const container_batch_start_index = (batches_needed - 1) * media_batch_size;
            const container_batch_end_index = Math.min(container_batch_start_index + media_batch_size, media_items.length);

            sliceOrderedMedias(container_batch_start_index, container_batch_end_index);
            
            media_focus_index = cached_media_index;
            
            await tick();

            let keyboard_selected_media = document.querySelector(`.meg-gallery-item[data-media-order="${media_focus_index}"]`);

            console.log(`Keyboard selected media:`, keyboard_selected_media);

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
         * Toggles the select stat of a given media.
         * @param {import('@models/Medias').OrderedMedia} media_item 
         */
        const toggleMediaSelect = (media_item) => {
            let current_media_change = $me_gallery_changes_manager.getMediaChangeType(media_item.uuid);

            if (current_media_change === media_change_types.MOVED) {
                $me_gallery_changes_manager.unstageMediaMove(media_item.uuid);
                return;
            }

            // on this meg gallery, we don't actually commit move media changes, only delete changes. instead of committing the move changes, we used them as
            // a select mechanism. if there are any move changes when committing, we will unstage them first.
            let fake_inner_category = new InnerCategory({name: "me-gallery-mock-category", uuid: "me-gallery-mock-category"});
            $me_gallery_changes_manager.stageMediaMove(media_item.Media, fake_inner_category);
            console.log(`Media ${media_item.uuid} staged to be moved.`);
        }

        /**
         * Toggles the deletion state of a given media.
         * @param {import('@models/Medias').OrderedMedia} ordered_media 
         * @param {boolean} [keep_deleted=true] - if true and the media is already staged for deletion, it will not be unstaged.
         */
        const toggleMediaDeletion = (ordered_media, keep_deleted=false) => {
            let current_media_change = $me_gallery_changes_manager.getMediaChangeType(ordered_media.uuid);

            if (current_media_change === media_change_types.DELETED && !keep_deleted) {
                $me_gallery_changes_manager.unstageMediaDeletion(ordered_media.uuid);
                return;
            }

            $me_gallery_changes_manager.stageMediaDeletion(ordered_media.Media);
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
                            enable_video_titles={show_media_titles_mode}
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