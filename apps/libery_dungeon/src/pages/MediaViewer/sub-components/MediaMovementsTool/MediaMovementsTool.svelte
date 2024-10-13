<script>
    /*=============================================
    =            Imports            =
    =============================================*/
    
        import CategorySearchBar from "@components/CategorySearchBar/CategorySearchBar.svelte";
        import InnerCategoryView from "../InnerCategoryView.svelte";
        import { current_category } from "@stores/categories_tree";
        import { getShortCategoryTree } from "@models/Categories";
        import { media_change_types } from "@models/WorkManagers";
        import { InnerCategory } from "@models/Categories";
        import { createEventDispatcher, onDestroy } from "svelte";
        import { onMount, tick } from "svelte";
        import { 
            active_media_index, 
            active_media_change, 
            media_changes_manager,
            auto_move_on,
            auto_move_category,
            media_viewer_hotkeys_context_name
        } from "@stores/media_viewer";
        import global_hotkeys_manager from "@libs/LiberyHotkeys/libery_hotkeys";
        import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
        import QuickMovementsTools from "./QuickMovementsTools.svelte";
        import { create_subcategory } from "@pages/MediaViewer/app_page_store";
        import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import { browser } from "$app/environment";
    
    /*=====  End of Imports  ======*/
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {boolean} whether the component is visible */
        export let is_component_visible = false;

        /** @type{boolean} the state of the quick media move tool*/
        let show_quick_move_tool = false;

        /** @type{number} the index of the selected category on the quick move tool*/
        let quick_selected_category_index = 0;

        /** @type {number} the index of the focused search result */
        let focused_search_result_index = 0;

        /** @type{string} the name of the quick move tool hotkey context */
        const quick_move_context_name = "quick_media_move_tool";

        /** @type {string} the name of the hotkey context that is active when results from the search bar are shown and exits when one of the results is selected */
        const search_results_context_name = "search_results_context";

        /** @type { InnerCategory[] } */
        let root_short_categories = [];

        /** @type {boolean} wheter the category tree is showing the actual tree or the recently used categories, as in the categories that had medias moved to them on the current session */
        let tree_view_enabled = true;

        /** @type {InnerCategory[]} the recently used categories */
        let recently_used_categories = [];

        /** @type {InnerCategory[]} the search results */
        let search_results = [];
        
        /*----------  Unsuscribers  ----------*/
        
            let media_viewer_hotkeys_context_name_unsub = () => {};

    /*=====  End of Properties  ======*/
    
    onMount(async () => {
        console.log("MediaMovementsTool mounted");

        root_short_categories = await getShortCategoryTree(null, $current_category.ClusterUUID);

        if ($media_viewer_hotkeys_context_name !== undefined) {
            registerExtraHotkeys();
            console.log("MediaMovementsTool hotkeys context name already defined");
        } else {
            media_viewer_hotkeys_context_name_unsub = media_viewer_hotkeys_context_name.subscribe(v => {
                if (v === undefined) return;
                console.log("MediaMovementsTool hotkeys context name changed");
                registerExtraHotkeys();
                // media_viewer_hotkeys_context_name_unsub();
            })
        }
    });

    onDestroy(() => {
        if (!browser) return;

        global_hotkeys_manager.dropContext(quick_move_context_name);
        global_hotkeys_manager.dropContext(search_results_context_name);
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        const defineQuickMoveHotkeys = () => {
            const quick_move_tool_context = new HotkeysContext();

            quick_move_tool_context.register("d", () => {
                    quick_selected_category_index = Math.min(quick_selected_category_index + 1, $media_changes_manager.UsedCategories.length - 1);
                },
                {
                    description: "<navigation>Move the selection to the next category",
                    can_repeat: true
                }
            );

            quick_move_tool_context.register("a", () => {
                    quick_selected_category_index = Math.max(quick_selected_category_index - 1, 0);
                },
                {
                    description: "<navigation>Move the selection to the previous category",
                    can_repeat: true
                }
            );

            
            quick_move_tool_context.register("c", () => setQuickMoveToolState(false), {
                mode: "keyup",
                description: `<${HOTKEYS_GENERAL_GROUP}>Release the key to close the quick move tool`
            });

            global_hotkeys_manager.declareContext(quick_move_context_name, quick_move_tool_context);
        }

        const defineSearchResultsHotkeys = () => {
            const search_results_context = new HotkeysContext();

            if (global_hotkeys_manager.hasContext(search_results_context_name)) return;

            search_results_context.register("s", () => {
                    moveSearchResultFocus(1);
                }, {
                    description: "<navigation>Move the focus to the next search result"
            });

            search_results_context.register("w", () => {
                    moveSearchResultFocus(-1);
                }, {
                    description: "<navigation>Move the focus to the previous search result"
            });

            search_results_context.register("e", () => {
                    selectSearchResult(focused_search_result_index);

                    exitSearchResultsHotkeysContext();
                }, {
                    mode: "keyup",
                    description: "<media_modification>Select the focused search result and move the active media to it",
            });

            search_results_context.register("q", () => {
                    focusSearchBarComponent();
                }, {
                    mode: "keyup",
                    description: "<navigation>Focus the search bar"
            });

            // TODO: ⟱ ⟱ ⟱ This down below is DISGUSTING, I have to implement a cleaner way to create new subcategories.
            search_results_context.register("n", () => {
                create_subcategory.set(true);
            }, {
                mode: "keyup",
                description: "<media_modification>Creates a new subcategory on the focused category/search result"
            });

            search_results_context.register("\`", exitSearchResultsHotkeysContext, {
                mode: "keypress",
                description: "<navigation>Exit the search results tool",
            });


            global_hotkeys_manager.declareContext(search_results_context_name, search_results_context);
        }

        /**
         * Focuses the search bar input. Bounded from the CategorySearchBar component
         * @type {() => void}   
         */
        let focusSearchBarComponent = () => {}

        const exitSearchResultsHotkeysContext = () => {
            let max_iterations = 100;
            let iterations = 0;

            while (global_hotkeys_manager.ContextName === search_results_context_name && iterations < max_iterations) {
                global_hotkeys_manager.loadPreviousContext();
                iterations++;
            }
        }

        const handleCategoryItemSelected = async e => {
            e.stopPropagation();

            /** @type {InnerCategory} */
            let moved_to_category = e.detail.category;

            if (moved_to_category === undefined) {
                console.error(`Event detail: ${JSON.stringify(e.detail)}`);
                return;
            }

            let current_media = $current_category.content[$active_media_index];

            if ($active_media_change === media_change_types.MOVED) {
                $media_changes_manager.unstageMediaMove(current_media);
                active_media_change.set(media_change_types.NORMAL);
                
                await tick();
            }

            $media_changes_manager.stageMediaMove(current_media, moved_to_category);

            active_media_change.set(media_change_types.MOVED);

            recently_used_categories = $media_changes_manager.UsedCategories;
            tree_view_enabled = false;

            // Reset search results after a category is selected from them
            if (search_results.length > 0) {
                search_results = [];
            }       

            is_component_visible = false;   

            // Return to the media viewer hotkeys context
            if (global_hotkeys_manager.ContextName === search_results_context_name) {
                exitSearchResultsHotkeysContext();
            }
            
            // check if the category should be set to auto move
            if (e.detail.shift_key) {
                handleAutomoveCategorySelected(moved_to_category);
            }

        }

        /**
         * sets the category to medias will be moved by default and enables the automove mode
         * @param {InnerCategory} category_selected
         */
        const handleAutomoveCategorySelected = category_selected => {
            
            if (category_selected === undefined || category_selected === null) return;

            auto_move_category.set(category_selected);

            toggleAutoMoveState(true);
        }

        const handleSearchResults = e => {
            /**
             * @type {Object[]}
             */
            let new_search_results = e.detail.results;

            if (new_search_results.length === 0) return;

            search_results = new_search_results.map(result => new InnerCategory(result));

            focused_search_result_index = 0;

            defineSearchResultsHotkeys();

            global_hotkeys_manager.loadContext(search_results_context_name);
        }

        /**
         * Handles the event of a new subcategory being created
         * @param {CustomEvent<CategoryCreatedEventDetail>} e
         * @typedef {Object} CategoryCreatedEventDetail
         * @property {import('@models/Categories').Category} category
         */
        const handleSubcategoryCreated = e => {
            console.log("Subcategory created event received", e);

            if (e.detail.category === null) return;

            let new_subcategory = e.detail.category;
            let inner_category = new_subcategory.toInnerCategory();

            search_results = [inner_category];
        }

        const moveSearchResultFocus = steps => {
            focused_search_result_index = Math.max(Math.min(focused_search_result_index + steps, search_results.length - 1), 0);
            console.log("focused_search_result_index: ", focused_search_result_index);
            console.log("search_results: ", search_results);

            const focused_result = search_results[focused_search_result_index];

            const focused_result_element = document.getElementById(`inner-category-${focused_result.uuid}`);
            if (focused_result_element === null) return;

            focused_result_element.scrollIntoView({ behavior: "instant", block: "nearest" });
        }   

        /**
         * Registers the component extra hotkeys, if any of them are already registered then they are ignored and a warning is logged
        */
        const registerExtraHotkeys = () => {
            console.log(`MediaMovementsTool hotkeys context name: ${$media_viewer_hotkeys_context_name}`);
            if ($media_viewer_hotkeys_context_name === undefined) return; // means that probably another sub-component is been used


            // register the hotkeys for the quick move tool
            if (!global_hotkeys_manager.hasHotkey("c", "keydown")) {
                global_hotkeys_manager.registerHotkeyOnContext("c", () => setQuickMoveToolState(true), {
                    mode: "keydown",
                    description: "<media_modification>Press the key to open the quick move tool. Release to close it."
                });
            }
        }

        /**
         * Stages the active media to be moved to the provided category
         * @param {InnerCategory} selected_category
        */
        const stageMediaMoved = async selected_category => {
            if (selected_category === undefined || selected_category === null) return;

            let current_media = $current_category.content[$active_media_index];

            if ($active_media_change === media_change_types.MOVED) {
                $media_changes_manager.unstageMediaMove(current_media);
                active_media_change.set(media_change_types.NORMAL);
                
                await tick();
            }

            $media_changes_manager.stageMediaMove(current_media, selected_category);

            active_media_change.set(media_change_types.MOVED);

            recently_used_categories = $media_changes_manager.UsedCategories;
            tree_view_enabled = false;

            // Reset search results after a category is selected from them
            if (search_results.length > 0) {
                search_results = [];
            }
        }

        /**
         * Sets the quick move tool state, if state is the same as the current state then does nothing.
         * if the new state is true, then sets the quick move tool hotkeys context as the current context before setting the state.
         * if the new state is false, then sets the previous context back after setting the state.
         * @param {boolean} state the new state of the quick move tool
        */
        const setQuickMoveToolState = state => {
            if (state === show_quick_move_tool) return;

            if (state) {
                defineQuickMoveHotkeys();
                global_hotkeys_manager.loadContext(quick_move_context_name);
            } else {
                stageMediaMoved($media_changes_manager.UsedCategories[quick_selected_category_index]); // move the active media to the last category selected on the quick move tool

                quick_selected_category_index = 0;

                global_hotkeys_manager.loadPreviousContext();

                registerExtraHotkeys();
            }

            show_quick_move_tool = state;
        }

        const selectSearchResult = index => {
            if (index < 0 || index >= search_results.length) return;

            handleCategoryItemSelected({
                detail: {
                    category: search_results[index],
                    shift_key: false
                },
                stopPropagation: () => {} // we pretend to be a mouse event 
            });
        }

        const toggleAutoMoveState = (force_state=null) => {
            let new_state = force_state === null ? !$auto_move_on : force_state;

            new_state = $auto_move_category === null ? false : new_state;

            auto_move_on.set(new_state);
        }
    
    /*=====  End of Methods  ======*/

</script>

{#if show_quick_move_tool}
    <QuickMovementsTools used_categories={$media_changes_manager.UsedCategories} {quick_selected_category_index}/>
{/if}
<div id="mv-category-tree-model" class:adebug={false} class="libery-scroll" style:visibility={is_component_visible ? "visible" : "hidden"}>
    <div class="mv-ctm-category" id="mv-ctm-searh-bar-section">
        <div id="mv-ctm-search-bar-wrapper">
            {#if is_component_visible}
                <CategorySearchBar bind:focusSearchBar={focusSearchBarComponent}  autofocus on:search-results={handleSearchResults}/>
            {/if}
        </div>
    </div>
    <menu class="mv-ctm-category" id="mv-ctm-tree-mode-controls">
        {#if !tree_view_enabled}
            <button id="mv-ctm-tree-view-btn" on:click={() => tree_view_enabled = true} >Tree mode</button>
        {:else if recently_used_categories.length > 0}
            <button id="mv-ctm-recently-used-btn" on:click={() => tree_view_enabled = false}>
                Recently used
            </button>
        {/if}
    </menu>
    <ul class="mv-ctm-category" id="categories-tree-root">
        {#if search_results.length > 0}
            {#each search_results as category}
                <InnerCategoryView 
                    is_keyboard_selected={category.uuid === search_results[focused_search_result_index]?.uuid}
                    category={category}
                    show_inner_categories={false}
                    on:category-item-selected={handleCategoryItemSelected}
                    on:new-subcategory-created={handleSubcategoryCreated}
                />
            {/each}
        {:else if tree_view_enabled}
            {#each root_short_categories as category}
                <InnerCategoryView category={category} on:category-item-selected={handleCategoryItemSelected}/>
            {/each}
        {:else if recently_used_categories.length > 0}
            {#each recently_used_categories as category}
                <InnerCategoryView category={category} show_inner_categories={false} on:category-item-selected={handleCategoryItemSelected}/>
            {/each}
        {/if}
    </ul>
</div>

<style>
    #mv-category-tree-model {
        box-sizing: border-box;
        height: 100%;
        background: var(--grey-7);
        padding: var(--vspacing-3) 0;
        overflow-y: auto;
    }

    .mv-ctm-category {
        padding: var(--vspacing-3) 0;
    }

    .mv-ctm-category:not(:last-child) {
        border-bottom: 1px solid var(--grey-5);
    }

    #mv-ctm-searh-bar-section {
        padding-left: var(--vspacing-3);
        padding-right: var(--vspacing-3);
    }

    #mv-ctm-tree-mode-controls {
        display: grid;
        height: 15%;
        margin: 0;
        place-items: center;
    }

    #mv-ctm-tree-mode-controls button {
        padding: var(--vspacing-1) var(--vspacing-2);
        background: transparent;
        border: 1px solid var(--main);
        color: var(--main);
        border-radius: var(--border-radius);
    }

    #mv-category-tree-model ul {
        list-style: none;
        margin: 0;
    }
</style>