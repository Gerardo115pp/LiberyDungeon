<script>
    /*=============================================
    =            Imports            =
    =============================================*/
        import { 
            media_upload_tool_mounted,
            category_creation_tool_mounted,
            media_transactions_tool_mounted,
            category_tagger_tool_mounted,
            category_search_focused,
            category_search_results,
            category_search_term,
            resetCategorySearchState,
        } from "@pages/MediaExplorer/app_page_store";
        import CategoryTagger from "@components/DungeonTags/CategoryTagger/CategoryTagger.svelte";
        import MediaUploadTool from "./sub-components/MediaUploadTool/MediaUploadTool.svelte";
        import CreateNewCategoryTool from "./sub-components/CreateNewCategoryTool.svelte";
        import CategoryFolder from "@components/Categories/CategoryFolder.svelte";
        import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
        import { HOTKEYS_HIDDEN_GROUP, HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import { CursorMovementWASD } from "@app/common/keybinds/CursorMovement";
        import { CategoryLeaf, getCategory, getCategoryTree, InnerCategory, moveCategory } from "@models/Categories";
        import LiberyHeadline from "@components/UI/LiberyHeadline.svelte";
        import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
        import MediasIcon from "@components/Medias/MediasIcon.svelte";
        import { onDestroy, onMount, tick } from "svelte";
        import { hotkeys_sheet_visible, layout_properties } from "@stores/layout";
        import { app_contexts } from "@libs/AppContext/app_contexts";
        import { app_context_manager } from "@libs/AppContext/AppContextManager";
        import { 
            categories_tree,
            current_category,
            resetCategoriesTreeStore,
            navigateToParentCategory,
            yanked_category,
            navigateToUnconnectedCategory
        } from "@stores/categories_tree";
        import { current_cluster, loadCluster } from "@stores/clusters";
        import MediaExplorerGallery from "./sub-components/MediaExplorerGallery/MediaExplorerGallery.svelte";
        import { me_gallery_yanked_medias } from "./sub-components/MediaExplorerGallery/me_gallery_state";
        import { page } from "$app/stores";
        import { goto, replaceState } from "$app/navigation";
        import { browser } from "$app/environment";
        import { confirm_question_responses, LabeledError, UIReference, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
        import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
        import { MediaChangesManager, resyncClusterBranch } from "@models/WorkManagers";  
        import { confirmPlatformMessage, emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
        import CategoryNavigationBreadcrumps from "@components/Categories/CategoryNavigationBreadcrumps.svelte";
        import global_platform_events_manager, { PlatformEventContext } from "@libs/LiberyEvents/libery_events";
        import { platform_well_known_events } from "@libs/LiberyEvents/well_known_events";
        import TransactionsManagementTool from "@components/TransactionsManagementTool/TransactionsManagementTool.svelte";
        import { current_user_identity } from "@stores/user";
        import { SearchResultsWrapper, wrapShowHotkeysTable } from "@app/common/keybinds/CommonActionWrappers";
        import { use_category_folder_thumbnails } from "@app/config/ui_design";
        import { common_action_groups } from "@app/common/keybinds/CommonActionsName";

    /*=====  End of Imports  ======*/
  
    /*=============================================
    =            App Context            =
    =============================================*/
    
        const app_context = app_context_manager.CurrentContext;
        const app_page_name = "medias-explorer";

        if (app_context !== app_contexts.DUNGEONS) {
            app_context_manager.setAppContext(app_contexts.DUNGEONS, app_page_name, resetCategoriesTreeStore);
        } else {
            app_context_manager.addOnContextExit(app_page_name, resetCategoriesTreeStore)
        }
    
    /*=====  End of App Context  ======*/
    
    /*=============================================
    =            Properties            =
    =============================================*/

        let global_hotkeys_manager = getHotkeysManager();

        /**
         * @type {string}
         */
        export let category_id;
    
        /**
         * @type {import('svelte/store').Unsubscriber}
         */
        let empty_categories_subscriber;

        const hotkey_context_name = "categories_explorer";
        
        /**
         * The category that is currently been focused by the user through keyboard navigation.
         * @type {number}
         */
        let keyboard_focused_category = 0;

        /**
         * Used to hold the category position when the user changes hotkeys context.
         * @type {number}
         */
        let keyboard_focused_category_holder = 0;

        /**
         * The WASD grid navigator.
         * @type {CursorMovementWASD | null}
         */
        let the_wasd_dungeon_explorer_navigator = null;
        
        /*----------  Child component properties  ----------*/
        
            /**
             * Extra classes to be added to the category folder component.
             * @type {string}
             */
            const category_content_member_html_class = "dme-content-item";
        
        /*----------  State  ----------*/

            /**
             * Whether a category is been dragged over the parent category label
             * @type {boolean}
             */
            let category_over_parent_label = false;

            /**
             * Whether to display the current category as a gallery or just as the normal MediasIcon.
             * In any case, the MediaIcon is also displayed.
             * @type {boolean}
             * @default false - change to true to modify the gallery UI.
             */
            let media_display_as_gallery = false;

            /**
             * @type {boolean}
             */
            export let was_media_display_as_gallery;
            $: console.log("MediaExplorer.was_media_display_as_gallery:", was_media_display_as_gallery);

            /**
             * Whether to force the media gallery to recover the focus on the last viewed media.
             * @type {boolean}
             */
            let recover_gallery_focus = false;

            /**
             * Whether to enable gallery hotkeys
             * @type {boolean}
             */
            let enable_gallery_hotkeys = false;

            /**
             * Whether the a category is syncing with the file system. right now we can't know when the resync is done
             * so we just use this flag to prevent the user from spamming the resync button. In the future when we implement jobs.
             * we can use this flag to show a loading spinner or something.
             * @type {boolean}
             */
            let is_category_resyncing = false;
        
        /*----------  Category name filter  ----------*/

            /**
             * The category name filter
             * @type {string}
             */
            let category_name_search_query = "";

            /**
             * Whether the search hotkey has been triggered
             * @type {boolean}
             */
            let capturing_category_name_search = false;

            /**
             * @type {import('svelte/store').Unsubscriber}
             */
            let filtered_category_selected_unsubscriber;

            /**
             * The inner category search results hotkey wrapper.
             * @type {SearchResultsWrapper<import('@models/Categories').InnerCategory> | null}
             */
            let the_category_search_results_wrapper = null;

            /**
             * A lookup set for the search result from the_category_search_results_wrapper.SearchResults
             * @type {Set<string> | null}
             */
            let category_local_search_results_lookup = null;

        
        /*----------  UI references  ----------*/

            /**
             * A UI reference for the categories.
             * @type {UIReference}
             */
            const ui_category_reference = new UIReference("cell", "cells");

            /**
             * A UI reference for the media in the category.
             * @type {UIReference}
             */
            const ui_media_reference = new UIReference("media", "medias");
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        if (!ensureCluster()) {
            return;
        }

        if ($categories_tree === null) {
            if (category_id === undefined) {

                const url_match_array = $page.url.pathname.match(/[a-zA-Z\d]{40}/g);

                category_id = $current_cluster.RootCategoryID;

                if (url_match_array != null) {
                    category_id = url_match_array[0];
                }

                console.debug("Attempted to infer category_id from URL or current cluster. category_id:", category_id);
            }

            let new_category_tree = await getCategoryTree(category_id, current_category);

            if (new_category_tree === null) {
                let variable_enviroment = new VariableEnvironmentContextError("In MediaExplorer.onMount while getting category tree");

                variable_enviroment.addVariable("category_id", category_id);
                variable_enviroment.addVariable("current_category", current_category);
                variable_enviroment.addVariable("current_cluster", current_cluster);

                let labeled_err = new LabeledError(variable_enviroment, "Failed to get the category tree", lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();

                return;
            }

            categories_tree.set(new_category_tree);
        }

        empty_categories_subscriber = current_category.subscribe(handleEmptyCategories); // this has to run after defining the desktop keybinds.

        filtered_category_selected_unsubscriber = current_category.subscribe(resetCategoryFiltersOnCategoryChange);

        if (!$layout_properties.IS_MOBILE) {
            defineDesktopKeybinds();
        }

        definePlatformEventHandlers();

        // DELETE EVERYTHING BELOW THIS LINE

        // category_tagger_tool_mounted.set(true);
    });

    onDestroy(() => {
        if (!browser) return;

        if (global_hotkeys_manager != null) {
            global_hotkeys_manager.dropContext(hotkey_context_name);
        }

        if (empty_categories_subscriber != null) {
            empty_categories_subscriber();  
        }

        if (filtered_category_selected_unsubscriber != null) {
            filtered_category_selected_unsubscriber();
        }

        dropGridNavigationWrapper();

        global_platform_events_manager.dropContext(app_page_name);
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /*=============================================
        =            Keybinding            =
        =============================================*/
        
            const defineDesktopKeybinds = () => {
                if (global_hotkeys_manager == null) {
                    console.error("Global hotkeys manager is null");
                    return;
                }

                if (!global_hotkeys_manager.hasContext(hotkey_context_name)) {

                    const hotkeys_context = new HotkeysContext();
            
                    hotkeys_context.register(["q", "esc"], handleGoToParentCategory, {
                        description: "<navigation>Go to parent category",
                    });

                    setGridNavigationWrapper(hotkeys_context);
            
                    hotkeys_context.register(["enter", "e"], handleCategoryItemSelected, {
                        description: "<navigation>Enter/Select category",
                    });

                    if ($current_user_identity?.canSyncClusters()) {
                        hotkeys_context.register(["alt+r"], handleSyncCurrentCategory, {
                            description: "<content>Syncs the state of the current category in the database with the real state of the file system. Useful when the file system is modified externally(like if you add a file to a category or delete folders or any other change).",    
                        });
                    }
                          
                    if ($current_user_identity?.canUploadFiles()) {
                        hotkeys_context.register(["u"], () => {
                                media_upload_tool_mounted.set(!$media_upload_tool_mounted);
                            }, {
                                description: "<content>Open media upload tool",
                        });
                    }

                    if ($current_user_identity?.canPublicContentAlter()) {
                        hotkeys_context.register(["y y"], yankSelectedCategory, {
                            description: "<content>Copy the selected category",
                        });

                        hotkeys_context.register(["p"], handleYankPaste, {
                            description: "<content>Paste the copied category or media. the category most be a valid category UUID and must not be an ancestor of the current category. if no category is copied, and there are medias copied, the medias will be pasted",
                        }); 

                        hotkeys_context.register(["c c"], handleRenameCategory, {
                            mode: "keyup",
                            description: "<content>Rename category",
                        });

                        hotkeys_context.register(["n"], handleCreateNewCategory, {
                            description: "<content>Open create new category tool",
                        });
                    }

                    hotkeys_context.register(["b"], () => { category_search_focused.set(!$category_search_focused);}, {
                        description: "<navigation>Open category search bar to search for categories in the current dungeon.",
                        mode: "keyup",
                    });

                    hotkeys_context.register(["; m"],  handleEnableCategoryFolderThumbnails, {
                        description: `${common_action_groups.EXPERIMENTAL}Enable ${ui_category_reference.EntityNamePlural} thumbnails. This feature is experimental and is not optimized.`,
                    });

                    // hotkeys_context.register(["/"], handleCategoryNameFilter, {
                    //     description: "<navigation>Filter subcategories in the current category by name",
                    // });

                    // hotkeys_context.register(["f \\l"], handleFindStartsWith, {
                    //     description: "<navigation>Filter subcategories in the current category by name",
                    //     consider_time_in_sequence: true,
                    // });

                    if ($current_user_identity?.canViewTrashcan()) {
                        hotkeys_context.register(["="], handleEnableTransactionManagementToolState, {
                            description: "<tools>Open the transaction management tool. This tool allows you manage your trashcan"
                        });
                    }

                    setSearchResultsWrapper(hotkeys_context);

                    wrapShowHotkeysTable(hotkeys_context);

                    hotkeys_context.register(["g", "shift+g"], enableMediaGalleryMode, {
                        description: "<content>Opens the category content as a gallery. the Shift key will automatically focus on the last viewed media.",
                    });

                    hotkeys_context.register(["t"], handleEnableCategoryTaggerTool, {
                        description: "<tools>Opens the category tagger tool."
                    });
            
                    global_hotkeys_manager.declareContext(hotkey_context_name, hotkeys_context);

                }

                global_hotkeys_manager.loadContext(hotkey_context_name);
            }

            /**
             * Destroys the WASD grid navigation wrapper. Is required because the grid navigation wrapper sets a MutationObserver on the grid parent element
             * to handle the addition/removal of new elements.
             */
            const dropGridNavigationWrapper = () => {
                if (the_wasd_dungeon_explorer_navigator != null) {
                    the_wasd_dungeon_explorer_navigator.destroy();
                }
            }

            /**
             * opens the category content as a gallery blow the subcategories grid. if key_event.shiftKey is true, it will open with the focus on the last viewed media.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const enableMediaGalleryMode = (key_event, hotkey) => {
                let key_combo = hotkey.KeyCombo.toLowerCase();

                media_display_as_gallery = !media_display_as_gallery;

                if (!media_display_as_gallery) {
                    was_media_display_as_gallery = false;
                }

                console.log("Shift is pressed: ", key_event.shiftKey);

                if (key_event.shiftKey) {
                    recover_gallery_focus = true;
                }

                console.debug('Current_category:', $current_category);
            }

            const handleGoToParentCategory = async () => {

                if ($current_category.parent != "") {
                    return navigateToSelectedCategory($current_category.parent);
                }

                goto("/");
            }

            const handleYankPaste = () => {
                if ($yanked_category !== "") {
                    pasteYankedCategory();
                } else if ($me_gallery_yanked_medias.length > 0) {
                    pasteYankedMedia();
                }
            }

            const handleRenameCategory = () => {
                const focused_element = document.querySelector(".keyboard-focused");

                if (focused_element === null) {
                    return;
                }

                focused_element.dispatchEvent(new CustomEvent("rename-requested"));
            }

            /**
             * The handler for cursor update bounded to the CursorMovementWASD.
             * @type {import('@common/keybinds/CursorMovement').CursorPositionCallback}
             */
            const handleCategoryContentCursorUpdate = (cursor_position) => {

                const should_move_to_another_hotkeys_context = shouldMoveToAnotherHotkeysContext( 
                    cursor_position.overflowed_top, 
                    cursor_position.overflowed_right, 
                    cursor_position.overflowed_bottom, 
                    cursor_position.overflowed_left
                );

                if (should_move_to_another_hotkeys_context) return;

                keyboard_focused_category = cursor_position.value;
            }
             
            const handleCategoryItemSelected = async () => {
                const can_keyboard_focus_be_medias = $category_search_results.length === 0 && $current_category.content.length > 0;
                const displayed_categories = getDisplayCategories();
                const keyboard_index_in_displayed_categories = (keyboard_focused_category >= 0 && keyboard_focused_category < displayed_categories.length);
                const keyboard_index_is_medias = keyboard_focused_category === displayed_categories.length;

                if (can_keyboard_focus_be_medias && keyboard_index_is_medias) {
                    return enterMediaViewer();
                } else if (!keyboard_index_in_displayed_categories) {
                    throw new Error("Keyboard focused category index is out of bounds");
                }

                const item_selected = getFocusedCategory();

                if (item_selected == null) return;

                navigateToSelectedCategory(item_selected.uuid);
            }

            const handleCreateNewCategory = () => {
                category_creation_tool_mounted.set(!$category_creation_tool_mounted);
            }

            /**
             * Enables the experimental feature use_category_folder_thumbnails. The thumbnails used for categories are fullsized images, heavliy deteriorates performance.
             * Once this feature is optimized and tested, it will be enabled by default.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleEnableCategoryFolderThumbnails = (event, hotkey) => {
                use_category_folder_thumbnails.set(!$use_category_folder_thumbnails);
            }

            const handleSyncCurrentCategory = async () => {
                console.log("Syncing current category");
                if (is_category_resyncing) return;
                is_category_resyncing = true;

                let sync_successful = await resyncClusterBranch($current_category.uuid, $current_cluster.UUID);
                if (sync_successful) {
                    await $categories_tree.updateCurrentCategory();
                } else {
                    let variable_enviroment = new VariableEnvironmentContextError("In MediaExplorer.handleSyncCurrentCategory");
                    variable_enviroment.addVariable("current_category.uuid", $current_category.uuid);
                    variable_enviroment.addVariable("current_cluster.UUID", $current_cluster.UUID);

                    let labeled_err = new LabeledError(variable_enviroment, "Failed to sync the current category with the file system", lf_errors.ERR_PROCESSING_ERROR);

                    labeled_err.alert();
                }

                is_category_resyncing = false;  
            }

            const handleEnableTransactionManagementToolState = () => {
                media_transactions_tool_mounted.set(!$media_transactions_tool_mounted); 
            }

            /**
             * Opens the category tagger tool. Handled by hotkeys
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleEnableCategoryTaggerTool = (event, hotkey) => {
                category_tagger_tool_mounted.set(true);
            }

            /**
             * Handle the category name filter. Note, the category search feature and category name filter are different. the first one checks with metaphones and other string differentiation 
             * techniques to find the category that matches the search term and is handled by the Backend. The category name filter is much simpler, it filters the categories in the current category
             * by a string which is matched with the start of the category name.
             * 
             * The feature works as follows, we first of all create a new hotkey context that responds to all ascii characters, the callback for all of these hotkeys registers the key pressed
             * and appends it to the category search term. Also at the end we register an interval that checks if the user has stopped typing for an amount of time, if so
             * we load the previous hotkey context and remove the interval
             */
            // const handleCategoryNameFilter = () => {
            //     if (global_hotkeys_manager == null) return;

            //     resetCategoryFiltersState(true);

            //     if (!global_hotkeys_manager.hasContext(filter_hotkeys_context_name)) {
            //         const filter_hotkeys_context = new HotkeysContext();

            //         /**
            //          * @param {KeyboardEvent} e
            //          * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
            //          */
            //         const handleAllLetters = (e, hotkey) => {
            //             category_name_filter += e.key.toLowerCase();
                        
            //             last_typing_time = Date.now();
            //             category_filter_changed = true;
            //         }

            //         filter_hotkeys_context.register(["\\l"], handleAllLetters, {
            //                 description: `<${HOTKEYS_HIDDEN_GROUP}>used for typing category name filter`,
            //         });

            //         filter_hotkeys_context.register(["backspace"], () => {
            //                 category_name_filter = category_name_filter.slice(0, -1);

            //                 last_typing_time = Date.now();
            //                 category_filter_changed = true;
            //             }, {
            //                 description: `<${HOTKEYS_HIDDEN_GROUP}>used for deleting category name filter`,
            //         });

            //         filter_hotkeys_context.register(["enter"], () => {
            //                 filterCategoriesByName();
            //                 resetCategoryFiltersState(false);
            //             }, {
            //                 description: `<media_movements>Apply category name filter`,
            //         });

            //         global_hotkeys_manager.declareContext(filter_hotkeys_context_name, filter_hotkeys_context);
            //         global_hotkeys_manager.loadContext(filter_hotkeys_context_name);

            //         category_name_filter_interval_id = setInterval(() => {

            //             let current_time = Date.now();

            //             filterCategoriesByName();

            //             if (category_name_filter !== "" && (current_time - last_typing_time) > category_name_filter_interval) {
            //                 resetCategoryFiltersState(false);
            //             }

            //             category_filter_changed = false;
            //         }, Math.max(category_name_filter_interval/3, 300));
            //     }
            // }

            /**
             * Set keyboard_focused_category to the first category that starts with the entered character.
             * @param {KeyboardEvent} e
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             * @returns {void}
             */
            const handleFindStartsWith = (e, hotkey) => {
                const initial_character = e.key;

                for (let h = 0; h < $current_category.InnerCategories.length; h++) {
                    let lower_case_name = $current_category.InnerCategories[h].name.toLowerCase();

                    if (lower_case_name.startsWith(initial_character)) {
                        setKeyboardFocusedCategory(h);
                        return;
                    }
                }
            }

            /**
             * The search result handler bound to the_category_search_results_wrapper
             * @type {import('@common/keybinds/CommonActionWrappers').SearchResultsUpdateCallback<import('@models/Categories').InnerCategory>}
             */
            const handleFocuseSearchMatch = search_match => {
                if (the_category_search_results_wrapper == null || search_match == null) return;

                capturing_category_name_search = false;

                let inner_category_uuid = search_match.uuid;

                let displayed_categories = getDisplayCategories();

                for (let h = 0; h <= displayed_categories.length; h++) {
                    let iteration_category = displayed_categories[h];

                    if (!(iteration_category instanceof InnerCategory)) {
                        console.error("iteration_category is not an InnerCategory:", iteration_category);
                        throw new Error("The displayed categories should only contain InnerCategory instances")
                    }

                    if (iteration_category.uuid === inner_category_uuid) {
                        setKeyboardFocusedCategory(h);
                        break;
                    }
                }

                if (the_category_search_results_wrapper.SearchResults.length > 0) {
                    category_local_search_results_lookup = new Set(the_category_search_results_wrapper.SearchResults.map((result => result.uuid)));
                }
            }

            /**
             * Handles the update of the search query label
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCaptureCallback}
             */
            const handleCategorySearchQueryUpdate = (event, captured_string) => {
                if (the_category_search_results_wrapper == null) return;

                capturing_category_name_search = true;
                category_name_search_query = captured_string;
            }

            /**
             * @type {import('@common/keybinds/CursorMovement').GridNavigationWrapperSetup}
             */
            const setGridNavigationWrapper = (hotkey_context) => {
                if (!browser) return;

                if (the_wasd_dungeon_explorer_navigator != null) {
                    the_wasd_dungeon_explorer_navigator.destroy();
                }

                const grid_selectors = getCategoriesGridSelectors();

                // The grid parent selector should only match one element, never more.
                const matching_grid_parent_count = document.querySelectorAll(grid_selectors.grid_parent_selector).length;
                if (matching_grid_parent_count !== 1) {
                    throw new Error(`The grid parent selector<${grid_selectors.grid_parent_selector}> should match only one element, but it matched ${matching_grid_parent_count}`); 
                }

                the_wasd_dungeon_explorer_navigator = new CursorMovementWASD(grid_selectors.grid_parent_selector, handleCategoryContentCursorUpdate, {
                    grid_member_selector: grid_selectors.grid_member_selector,
                    sequence_item_name: "content item",
                    sequence_item_name_plural: "content items",
                    initial_cursor_position: keyboard_focused_category,
                });
                the_wasd_dungeon_explorer_navigator.setup(hotkey_context);

                // @ts-ignore
                globalThis.the_wasd_dungeon_explorer_navigator = the_wasd_dungeon_explorer_navigator;
            }

            /**
             * @type {import('@common/keybinds/CommonActionWrappers').SearchResultsWrapperSetup}
             */
            const setSearchResultsWrapper = hotkey_contenxt => {

                the_category_search_results_wrapper = new SearchResultsWrapper(hotkey_contenxt, $current_category.InnerCategories, handleFocuseSearchMatch, {
                    minimum_similarity: 0.7,
                    search_hotkey: ["f"],
                    ui_search_result_reference: ui_category_reference,
                    search_typing_hotkey_handler: handleCategorySearchQueryUpdate,
                    boost_exact_inclusion: false,
                    allow_member_similarity_checking: true,
                    no_results_callback: () => resetCategoryFiltersState()
                });

                the_category_search_results_wrapper.setItemToStringFunction(inner_category => inner_category.name.toLowerCase());
            }

        /*=====  End of Keybinding  ======*/
        
        /*=============================================
        =            PlatformEvents            =
        =============================================*/
        
            const definePlatformEventHandlers = () => {
                const dungeon_explorer_handlers = new PlatformEventContext();
                if (!global_platform_events_manager.hasContext(app_page_name)) {
                    dungeon_explorer_handlers.addEventHandler(platform_well_known_events.FS_CHANGED, handleFsChangedPlatformEvent);
                    
                    global_platform_events_manager.declareContext(app_page_name, dungeon_explorer_handlers);
                }

                global_platform_events_manager.loadContext(app_page_name);

            }

            /**
             * Reloads the current category when the platform event FS_CHANGED is received.
             * @param {import("@libs/DungeonsCommunication/transmissors/platform_events_transmisor").PlatformEventMessage<import('@libs/LiberyEvents/well_known_events').ClusterFsChangeEvent>} event
             */
            const handleFsChangedPlatformEvent = async (event) => {
                console.log("FS_CHANGED event received", event);

                const fs_change_message = event?.EventPayload.payload;

                if (fs_change_message == null) {
                    console.error("FS_CHANGED event received with no payload");
                    return;
                }

                if (fs_change_message.cluster_uuid === $current_cluster.UUID && (fs_change_message.medias_added + fs_change_message.medias_deleted + fs_change_message.medias_updated) !== 0) {
                    await $categories_tree.updateCurrentCategory();
                } 
            }
                
        
        /*=====  End of PlatformEvents  ======*/

        /**
         * Restores hotkeys control to the categories explorer.
         * @returns {void}
         */
        const disableGalleryHotkeys = () => {
            if (global_hotkeys_manager == null) return;

            loadKeyboardFocusedCategory();
            enable_gallery_hotkeys = false; 

            if (global_hotkeys_manager.ContextName !== hotkey_context_name) {
                global_hotkeys_manager.loadContext(hotkey_context_name);
            }
        }

        /**
         * Makes sure there is a current_cluster set, if there isn't then it will try to load the cluster from local storage if that also fails then it will redirect the user to '/'.
         * Note: this method must be replaced when we move from pure svelte to sveltekit.
         * @returns {boolean} Whether it successfully loaded the cluster or not
         */
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

        /**
         * Navigates to the media viewer for the current category
         * @param {number} [media_index]
         */
        const enterMediaViewer = (media_index) => {
            let href = `/${$layout_properties.IS_MOBILE ? 'media-viewer-mobile' : 'media-viewer'}/${$current_category.uuid}`;

            if (media_index != null) {
                href += `/${media_index}`;
            }

            goto(href); 
        }

        /**
         * Grants hotkeys control to the explorer gallery. saves the current category position and sets the keyboard_focused_category to -1. so that no
         * category appears as focused in the explorer UI.
         * @returns {void}
         */
        const enableGalleryHotkeys = () => {
            saveKeyboardFocusedCategory();
            setKeyboardFocusedCategory(-1);

            enable_gallery_hotkeys = true;
        }

        /**
         * Focuse a given category by its uuid. Returns true if the category was found and focused, false otherwise.
         * @param {string} category_uuid
         * @returns {boolean}
         */
        const focusCategoryByUUID = (category_uuid) => {
            const displayed_categories = getDisplayCategories();
            const category_index = displayed_categories.findIndex(inner_category => inner_category.uuid === category_uuid);

            if (category_index === -1) {
                return false;
            }

            setKeyboardFocusedCategory(category_index);
            return true
        }

        /**
         * Returns the categories that must be displayed in the current viewport. 
         * @returns {import('@models/Categories').InnerCategory[]}
         */
        const getDisplayCategories = () => {
            const displaying_search_results = $category_search_results.length > 0;

            /**
             * @type {import('@models/Categories').InnerCategory[]}
             */
            let displayed_categories = [];

            if (displaying_search_results) {
                displayed_categories = $category_search_results.map(category => category.toInnerCategory());
            } else {
                displayed_categories = $current_category.InnerCategories;
            }

            return displayed_categories;
        }

        /**
         * Returns the category focused by the user through keyboard navigation.
         * @returns {import('@models/Categories').InnerCategory | undefined}
         */
        const getFocusedCategory = () => {
            const displayed_categories = getDisplayCategories();

            return displayed_categories[keyboard_focused_category];
        }

        /**
         * Returns the categories grid list pair of selectors on a 2 sized array with index 0 being the grid-parent selector and the index 1 being the grid-member selector. 
         * This is used for the parameters of the CursorMovementWASD either on creation or afterwards for it's changeGridContainer method.
         * @returns {import('@common/interfaces/common_actions').GridSelectors}
         */
        const getCategoriesGridSelectors = () => {
            /**
             * @type {import('@common/interfaces/common_actions').GridSelectors}
             */
            const grid_selectors = {
                grid_parent_selector: "#libery-categories-explorer #category-content",
                grid_member_selector: `.${category_content_member_html_class}`,
            }


            return grid_selectors;
        }

        /**
         * Checks if a category is empty, meaning it has no subcategories and no content. If it is empty, then notify
         * the user that the category is empty and ask if they want to delete it. If so, go back to the parent category
         * and send a request to the server to delete the category.
         * @param {CategoryLeaf} category   
         */
        const handleEmptyCategories = async (category) => {
            const category_param_valid = !(category === null && category === undefined)
            const category_protected = category.uuid === $current_cluster.RootCategoryID || category.uuid === $current_cluster.DownloadCategoryID;  

            if (!category_param_valid || !category?.isEmpty() || category_protected) {
                return;
            }            

            const delete_category_choice = await confirmPlatformMessage({
                message_title: "Empty category",
                question_message: `The category '${category.name}' is empty, do you want to delete it? If there are files on the cateogry folder that are note registered on the system, only system record of this category will be deleted but the directory and the not-registered files will remain.`,
                cancel_label: "No",
                confirm_label: "Yes",
                danger_level: 0,
                auto_focus_cancel: false,
            });

            if (delete_category_choice === confirm_question_responses.CONFIRM) {
                const category_to_delete = category.uuid;
                await handleGoToParentCategory();
                $categories_tree.deleteChildCategory(category_to_delete);
            }
        }

        /**
         * Handles the close event from the explorer gallery.
         */
        const handleGalleryClose = () => {
            console.debug("Gallery close event received");
            disableGalleryHotkeys();
            media_display_as_gallery = false;
            was_media_display_as_gallery = false;
        }

        /**
         * Handles the drag enter event on the category name label.
         * @param {DragEvent} e
         */
        const handleParentDragEnter = e => {
            console.debug(`Drag enter on parent category label. parent: ${$current_category.parent}`);
            if ($current_category.parent != null) {
                e.preventDefault();
                e.stopPropagation();
                category_over_parent_label = true;

                if (e.dataTransfer != null) {
                    e.dataTransfer.dropEffect = "move";
                }
            }
        }

        /**
         * Handles the drag over event on the category name label.
         * @param {DragEvent} e
         */
        const handleParentDragOver = e => {
            if ($current_category.parent != null) {
                e.preventDefault();
                e.stopPropagation();
            }
        }

        /**
         * Handles the drag leave event on the category name label.
         * @param {DragEvent} e
         */
        const handleParentDragLeave = e => {
            if ($current_category.parent != null) {
                e.preventDefault();
                e.stopPropagation();
                category_over_parent_label = false;
            }
        }

        /**
         * Handles the drop event on the category name label. Attempts to move the dropped category to the parent of the current category.
         * @param {DragEvent} e
         */
        const handleDropCategoryOnParent = async e => {
            const dragged_category_uuid = e.dataTransfer?.getData("text/plain");
            category_over_parent_label = false;

            if ($current_category.parent != null && dragged_category_uuid != null) {
                e.preventDefault();
                e.stopPropagation();

                if (e.dataTransfer != null) {
                    e.dataTransfer.dropEffect = "move";
                }

                console.debug(`Dropped category '${dragged_category_uuid}' on parent category '${$current_category.parent}'`);

                let updated_dragged_category = await moveCategory(dragged_category_uuid, $current_category.parent);

                if (updated_dragged_category == null) {
                    return;
                }

                if ($current_category.ParentCategory != null) {
                    let updated_inner_category = updated_dragged_category.toInnerCategory();
                    $current_category.ParentCategory.addInnerCategory(updated_inner_category);
                }

                $categories_tree.updateCurrentCategory();
            }
        }

        /**
         * Handles the breadcrumb fragment selected event. Navigates to the selected category.
         * @param {CustomEvent<import('@models/Categories').Category>} event
         */
        const handleBreadcrumbSelected = async event => {
            let selected_category = event.detail;
            
            if (selected_category.UUID === $current_category.uuid || selected_category.UUID === "") {
                return;
            }

            navigateToSelectedCategory(selected_category.UUID);
        }

        /**
         * Handles the hotkey context control requested by the explorer gallery.
         */
        const handleGalleryRequestedControl = () => {
            media_display_as_gallery = true;
            shouldMoveToAnotherHotkeysContext(false, false, true, false);
        }

        /**
         * Hanldes the open-media-viewer event from the explorer gallery. 
         * @param {CustomEvent<OpenMediaViewerDetail>} event
         * @typedef {Object} OpenMediaViewerDetail
         * @property {number} media_index
         * @property {import('@models/Medias').Media} media_item
         */
        const handleGalleryOpenMediaViewer = event => {
            const { media_index } = event.detail;

            console.log("Opening media viewer for media index:", media_index);

            was_media_display_as_gallery = true;
            
            enterMediaViewer(media_index);
        }

        /**
         * Hanldes the clouser of the transaction management tool.
         */ 
        const handleTransactionManagementToolClose = () => {
            media_transactions_tool_mounted.set(false);
        }

        /**
         * Handles close-category-tagger event emitted by the CategoryTagger component.
         */
        const handleCategoryTaggerClose = () => {
            category_tagger_tool_mounted.set(false);
        }

        /**
         * Navigates to a selected category. Determines if the category is a child category in which case it will navigate to it using SPA navigation.
         * If the category is not a child category or parent, it will navigate it using the browser's navigation.
         * @param {string} category_uuid
         */
        const navigateToSelectedCategory = async category_uuid => {
            resetNavigationSuseptibleState();

            switch (true) {
                case $current_category.isParentCategoryUUID(category_uuid):
                    console.log("Navigating to parent category");
                    const current_category_uuid = $current_category.uuid;
                    await navigateToParentCategory();
                    await tick();
                    focusCategoryByUUID(current_category_uuid);
                    return 
                case $current_category.isChildCategoryUUID(category_uuid):
                    return navigateToChildCategory(category_uuid);
                default:
                    if (category_uuid === $current_category.uuid) return;
                    console.log("Navigating to unconnected category");
                    navigateToUnconnectedCategory(category_uuid);
            }                
        }

        /**
         * Navigates to a child category using spa navigation.
         * @param {string} category_uuid
         */
        const navigateToChildCategory = async (category_uuid) => {
            replaceState(`/dungeon-explorer/${category_uuid}`, $page.state);
            $categories_tree.navigateToLeaf(category_uuid);
        }

        /**
         * Checks the yanked_category store, if it's an empty string, it attempts to fetch the clipboard content and check if the content is a valid category UUID. if
         * it manages to get a category UUID, it will use moveCategory to move the category to the current category.
         * @returns {Promise<void>} 
         */
        const pasteYankedCategory = async () => {
            if (!navigator.userActivation?.isActive) {
                console.warn("Transient user activation required to paste the copied category");
                return;
            };

            let category_uuid = $yanked_category;

            if (category_uuid === "") {
                try {
                    category_uuid = await navigator.clipboard.readText();
                } catch (error) {
                    console.error("Error reading clipboard content:", error);
                    return;
                }

                if (typeof category_uuid !== 'string' || category_uuid.length !== 40) {
                    console.warn(`Clipboard content is not a valid category UUID: ${category_uuid}`);
                    return;
                }
            }

            let trusted_category = await getCategory(category_uuid); // If the server returns a category, the category is trusted because if absolutely exists.

            if (trusted_category == null) {
                console.warn(`Category with UUID(${category_uuid}) doesn't seem to exist`);
                return;
            }

            if ($current_category.isParentCategoryPath(trusted_category.FullPath)) {
                console.warn(`Category with UUID(${category_uuid}) is an ancestor of the current category`);
                return;
            }

            let updated_trusted_category = await moveCategory(category_uuid, $current_category.uuid);

            if (updated_trusted_category != null && updated_trusted_category.UUID === trusted_category.UUID) {
                $categories_tree.updateCurrentCategory();
                yanked_category.set("");
            } else {
                alert(`Error moving category with UUID(${category_uuid}) to the current category`);
            }
        }

        /**
         * Reads the medias in me_gallery_yanked_medias store, and if all of them dont belong to this category, it will create a new media_changes_manager, stage them to be moved
         * and pass the current category uuid as the destination category.
         */
        const pasteYankedMedia = async () => {
            if ($me_gallery_yanked_medias.length === 0) {
                let labeled_err = new LabeledError("In MediaExplorerGallery.handleYankPaste", "No medias to paste.", lf_errors.ERR_PROCESSING_ERROR);
                labeled_err.alert();
                return;
            }
            const all_yanked_medias_count = $me_gallery_yanked_medias.length;

            let not_all_medias_foreign = $me_gallery_yanked_medias.some(m => m.main_category === $current_category.uuid)

            if (not_all_medias_foreign) {
                let labeled_err = new LabeledError("In MediaExplorerGallery.handleYankPaste", "Medias, at least some, already belong to this category.", lf_errors.ERR_PROCESSING_ERROR);
                labeled_err.alert();   
                return;
            }

            /**
             * A hashmap of category_uuid -> medias[], where category_uuid is the source where the category comes from and medias[] 
             * is an array of medias come from that category.
             * @type {Object<string, import('@models/Medias').Media[]>}
             */
            let yanked_media_sources = {}

            $me_gallery_yanked_medias.forEach(media => {
                yanked_media_sources[media.main_category] = yanked_media_sources[media.main_category] ?? [];

                yanked_media_sources[media.main_category].push(media);
            });

            console.log("Yanked media sources:", yanked_media_sources);
            let media_changes_manager = new MediaChangesManager();

            for (let source_category_uuid of Object.keys(yanked_media_sources)) {
                let medias = yanked_media_sources[source_category_uuid];
                console.log(`Staging medias from category(${source_category_uuid}):`, medias);

                medias.forEach(media => {
                    media_changes_manager.stageMediaMove(media, $current_category.asInnerCategory());
                });

                await media_changes_manager.commitChanges(source_category_uuid);
                await $categories_tree.updateCategory(source_category_uuid);
                media_changes_manager.clearAllChanges();
            }

            await $categories_tree.updateCurrentCategory();

            me_gallery_yanked_medias.set([]);

            emitPlatformMessage(`Pasted ${all_yanked_medias_count} new medias in `+ $current_category.name);
        }

        /**
         * Resets the state of the category name filter
         * @returns {void}
         */
        const resetCategoryFiltersState = () => {
            if (global_hotkeys_manager == null) return;

            category_name_search_query = "";
            capturing_category_name_search = false;
            category_local_search_results_lookup = null;
        }

        /**
         * Resets the category filters state on current category change
         * @returns {void}
         */
        const resetCategoryFiltersOnCategoryChange = () => {
            if (the_category_search_results_wrapper == null) return;

            if (category_name_search_query !== "") {
                resetCategoryFiltersState();
            }

            the_category_search_results_wrapper.updateSearchPool($current_category.InnerCategories);            
        }

        /**
         * Resets all the state changes that should expire after a navigating to a different category.
         * @returns {void}
         */
        const resetNavigationSuseptibleState = () => {
            setKeyboardFocusedCategory(0);
            keyboard_focused_category_holder = 0; 
            
            resetCategoryFiltersState();
            resetCategorySearchState();

            if (media_display_as_gallery) {
                disableGalleryHotkeys();
                media_display_as_gallery = false;
            }
        }

        /**
         * Meant to be used in parity with the handleCategoryMove method. Is called when the keyboard_focused_category overflows. And it determines whether an overflow to that
         * direction should change the hotkeys context. tipically that would be because there is another component in that direction that can handle the same type of movement.
         * returns true if an action was taken, in that case the caller is expected to not take any further action. if false, the caller is free to take whatever action it
         * deems adequate.
         * @param {boolean} up
         * @param {boolean} right
         * @param {boolean} down
         * @param {boolean} left
         * @returns {boolean}
         */
        const shouldMoveToAnotherHotkeysContext = (up, right, down, left) => {
            if (down && media_display_as_gallery) {
                enableGalleryHotkeys();
                return true;
            }

            return false;
        }

        /**
         * Sets the keyboard_focused_category to the given index.
         * @param {number} index
         */
        const setKeyboardFocusedCategory = index => {
            if (the_wasd_dungeon_explorer_navigator !== null) {
                the_wasd_dungeon_explorer_navigator.updateCursorPosition(index);
            } else {
                keyboard_focused_category = index;
            }
        }

        /* ------------------------ Keyboard focuse save/load ----------------------- */

            /**
             * Saves the current keyboard focused category index.
             * @returns {void}
             */
            const saveKeyboardFocusedCategory = () => {
                keyboard_focused_category_holder = keyboard_focused_category;
            }

            /**
             * Loads the saved keyboard focused category index.
             * @returns {void}
             */
            const loadKeyboardFocusedCategory = () => {
                setKeyboardFocusedCategory(keyboard_focused_category_holder);
            }

        /* -------------------------------------------------------------------------- */

        /**
         * Yanks the selected category
         */
        const yankSelectedCategory = async () => {
            if (!navigator.userActivation?.isActive) {
                console.warn("Transient user activation required to copy the selected category");
                return;
            };

            const selected_inner_category = $current_category.InnerCategories[keyboard_focused_category];

            if (selected_inner_category == null) {
                alert("No category selected");
                return;
            }

            yanked_category.set(selected_inner_category.uuid);

            /**
             * @type {Promise<void>}
             */
            let written_promise;

            try {
                written_promise = navigator.clipboard.writeText(selected_inner_category.uuid);
            } catch (error) {
                console.error(`Error copying category UUID(${selected_inner_category.uuid}) to clipboard:`, error);
                return;
            }

            await written_promise;

            console.debug(`Copied category UUID(${selected_inner_category.uuid}) to clipboard`);    
        }
    
    /*=====  End of Methods  ======*/

</script>

<svelte:head>
    <title>
        {$current_category != null ? `${$current_category.name} - Dungeon Cell` : "Dungeon cells"}
    </title>
</svelte:head>
<main id="libery-categories-explorer"
    class:with-category-thumbnails={$use_category_folder_thumbnails}
    style:position="relative"
>
    {#if $category_creation_tool_mounted}
        <div class="fullwidth-modals" id="new-category-component-wrapper">
            <CreateNewCategoryTool />
        </div>
    {/if}
    {#if $media_upload_tool_mounted}
        <div class="fullwidth-modals" id="media-upload-component-wrapper">
            <MediaUploadTool />
        </div>
    {/if}
    {#if $media_transactions_tool_mounted}
        <div class="fullwidth-modals" id="media-transactions-component-wrapper">
            <TransactionsManagementTool
                transaction_management_tool_open={$media_transactions_tool_mounted}
                on:tool-closed={handleTransactionManagementToolClose}
            />
        </div>
    {/if}
    {#if $category_tagger_tool_mounted}
        <CategoryTagger
            on:close-category-tagger={handleCategoryTaggerClose}
        />
    {/if}
    <div id="lce-floating-controls-overlay">
        {#if category_name_search_query || capturing_category_name_search}
             <p id="lce-fco-category-filter-string">
                /{category_name_search_query}
             </p>
        {/if}
    </div>
    <header id="lce-upper-content">
        {#if $current_category != null}
            <CategoryNavigationBreadcrumps 
                category_leaf_item={$current_category}
                on:breadcrumb-selected={handleBreadcrumbSelected}
            />
        {/if}
        <button id="lce-uc-current-category-name" 
            class:parent-accepting-drop={category_over_parent_label}
            on:click={handleGoToParentCategory} 
            on:dragover={handleParentDragOver}
            on:dragenter={handleParentDragEnter}
            on:dragleave={handleParentDragLeave}
            on:drop={handleDropCategoryOnParent}
        >
            {#if $current_category !== null && $current_category !== undefined}
                <LiberyHeadline 
                    headline_tag="Cell"
                    headline_color="var(--grey-1)"
                    forced_font_size="calc(var(--font-size-{$layout_properties.IS_MOBILE ? '4' : 'h3'}) * 1.6)"
                    extra_props="style='awesome!'"
                    headline_text={$category_search_term === "" ? $current_category.name : `Search: ${$category_search_term}`}
                    force_bottom_lines
                />
            {/if}
        </button>
    </header>
    <ul id="category-content" class:debug={false}>
        {#if $category_search_results.length > 0}
            <!-- Search results -->
            {#each $category_search_results as result_category, h}
                <CategoryFolder 
                    category_item_class={category_content_member_html_class}
                    inner_category={result_category.toInnerCategory()} 
                    category_keyboard_focused={keyboard_focused_category === h}
                    is_ephemeral
                />
            {/each}
        {:else if $current_category !== null}
            <!-- Category content -->
            {#each $current_category?.InnerCategories as ic, h}
                <CategoryFolder 
                    highlight_category={category_local_search_results_lookup != null && category_local_search_results_lookup.has(ic.uuid)}
                    category_item_class={category_content_member_html_class}
                    inner_category={ic}
                    category_keyboard_focused={keyboard_focused_category === h}
                />
            {/each}
            {#if $current_category.content.length > 0}
                <MediasIcon 
                    extra_class={category_content_member_html_class}
                    images_count={$current_category.content.length} 
                    keyboard_focused={keyboard_focused_category === $current_category.InnerCategories.length}
                />
            {/if}
        {/if}
    </ul>
    {#if (media_display_as_gallery || was_media_display_as_gallery) && ($current_category?.content?.length ?? 0) > 0}
        <div id="media-gallery-wrapper">
            <MediaExplorerGallery 
                media_items={$current_category.content}
                enable_gallery_hotkeys={enable_gallery_hotkeys}
                recovering_gallery_state={was_media_display_as_gallery || recover_gallery_focus}
                on:exit-hotkeys-context={disableGalleryHotkeys}
                on:request-hotkeys-control={handleGalleryRequestedControl}
                on:close-gallery={handleGalleryClose}
                on:open-media-viewer={handleGalleryOpenMediaViewer}
                use_masonry={!enable_gallery_hotkeys && $current_category.content.length > 15}
            />
        </div>
    {/if}
</main>

<style>
    #libery-categories-explorer {
        position: relative;
        display: flex;
        flex-direction: column;
        row-gap: var(--vspacing-4);
        padding: var(--vspacing-3);
    }
    
    /*=============================================
    =            header            =
    =============================================*/

        header#lce-upper-content {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-3);
        }
    
        button#lce-uc-current-category-name {
            box-sizing: content-box;
            padding: var(--spacing-2) var(--spacing-1);
            margin: 0;
            cursor: pointer;
            transition: all 0.2s linear;

            & > * {
                pointer-events: none;
            }
        }

        @media(pointer: fine) {
            button#lce-uc-current-category-name:hover {
                background-color: var(--grey-9);
            }
        }

        button#lce-uc-current-category-name.parent-accepting-drop {
            background-color: var(--grey-9);
            border: 2px dotted var(--main-dark);
        }
    
    
    /*=====  End of header  ======*/
    
    /*=============================================
    =            Controls overlay            =
    =============================================*/
    
        #lce-floating-controls-overlay {
            position: fixed;
            bottom: 0;
            left: 0;
            width: 100%;
            padding: var(--vspacing-1);
        }
    
    /*=====  End of Controls overlay  ======*/
    
    
    
    .fullwidth-modals {
        position: fixed;
        top: var(--navbar-height);
        left: 0;
        width: 100%;
        height: calc(100dvh - var(--navbar-height));
        z-index: var(--z-index-t-1);
    }

    ul#category-content {
        --category-folder-size: 230px;

        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(var(--category-folder-size), 1fr));
        grid-auto-rows: var(--category-folder-size);
        list-style: none;
        padding: 0;
        gap: 2px;
    }

    
    /*=============================================
    =            Category thumbnails            =
    =============================================*/
    
    #libery-categories-explorer.with-category-thumbnails ul#category-content {
        grid-auto-rows: calc(var(--category-folder-size) * 1.7);
        gap: var(--spacing-2);
    }
    
    /*=====  End of Category thumbnails  ======*/
    
    


    @media only screen and (max-width: 768px) {
        #category-content {
            grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        }
    }
</style>