<script>
    import { browser } from "$app/environment";
    import { current_category } from "@stores/categories_tree";
    import { current_cluster } from "@stores/clusters";
    import { onDestroy, onMount } from "svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The list of previous categories. fetched from the cluster's history.
         * @type {import('@models/Categories').InnerCategory[]}
         */
        let categories_history_array = [];

        /**
         * Whether the category history has mutated since it was last fetched.
         * @type {boolean}
         * @default true
         */
        let categories_history_mutated = true; 

        /**
         * Whether the user has the mouse over the history button.
         * @type {boolean}
         */
        let history_button_hovered = false;

        /**
         * whether the categories history is pinned.
         * @type {boolean}
         */
        let categories_history_pinned = false;

        /**
         * Enables development mode debugging features.
         * @type {boolean}
         */
        const debug_mode = true;
        
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        if ($current_cluster === null) {
            throw new Error("In NavCategoryHistory.onMount: $current_cluster is null.");
        }

        updateCategoriesHistoryIfOutdated();

        attachCategoryHistoryEventListeners();


        debug__attachDebugMethods();
    });

    onDestroy(() => {
        removeCategoryHistoryEventListeners();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
            
        /*=============================================
        =            Debug            =
        =============================================*/

                /**
                 * Returns the object where all the debug state is stored.
                 * @returns {Object}
                 */
                const debug__getComponentDebugState = () => {
                    if (!browser || !debug_mode) return {};

                    const NAV_CATEGORY_HISTORY_DEBUG_STATE_NAME = "nav_category_history_debug_state";

                    // @ts-ignore
                    if (!globalThis[NAV_CATEGORY_HISTORY_DEBUG_STATE_NAME]) {
                        // @ts-ignore
                        globalThis[NAV_CATEGORY_HISTORY_DEBUG_STATE_NAME] = {

                        };   
                    }

                    // @ts-ignore
                    return globalThis[NAV_CATEGORY_HISTORY_DEBUG_STATE_NAME];
                }
        
                /**
                 * Attaches debug methods to the globalThis object for debugging purposes.
                 * @returns {void}
                 */
                const debug__attachDebugMethods = () => {
                    if (!debug_mode || !browser) return;

                    const me_component_debug_state = debug__getComponentDebugState();

                    // @ts-ignore - for debugging purposes we do not care whether the globalThis object has the method name. same reason for all other ts-ignore in this function.
                    me_component_debug_state.printDebugState = debug__printComponentState;

                    // @ts-ignore - state retrieval functions.
                    me_component_debug_state.State = {
                        CurrentCluster: () => $current_cluster,
                        CurrentCategory: () => $current_category,
                    }

                    // @ts-ignore - Internal method references.
                    me_component_debug_state.Methods = {
                        isCategoriesHistoryVisible,
                    }
                }

                /**
                 * Prints the whole gallery state to the console.
                 * @returns {void}
                 */
                const debug__printComponentState = () => {
                    console.log("%NavCategoryHistory State", "color: green; font-weight: bold;");
                    console.group("Properties");
                    console.log("categories_history_array: %O", categories_history_array);
                    console.log(`categories_history_mutated: ${categories_history_mutated}`);
                    console.log(`history_button_hovered: ${history_button_hovered}`);
                    console.log(`categories_history_pinned: ${categories_history_pinned}`);
                    console.groupEnd();
                }

                /**
                 * Attaches an arbitrary object as a globalThis.media_explorer_debug_state.<group_name>{...timestamp -> object }.
                 * @param {string} group_name
                 * @param {object} object_to_snapshot
                 * @returns {void}
                 */
                const debug__attachSnapshot = (group_name, object_to_snapshot) => {
                    if (!browser || !debug_mode) return;

                    const stack = new Error().stack;
                    const datetime_obj = new Date();
                    const timestamp = `${datetime_obj.toISOString()}-${datetime_obj.getTime()}`;

                    const snapshot = {
                        timestamp,
                        stack,
                        object_to_snapshot,
                    }

                    const debug_object = debug__getComponentDebugState();

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
        =            DOM Event handlers            =
        =============================================*/
        
            /**
             * Handles the mouse enter event on the history button.
             * @param {MouseEvent} event
             */ 
            const handleHistoryButtonMouseEnter = (event) => {
                history_button_hovered = true;

                updateCategoriesHistoryIfOutdated();
            }

            /**
             * Handles the mouse leave event on the history button.
             * @param {MouseEvent} event
             */
            const handleHistoryButtonMouseLeave = (event) => {
                history_button_hovered = false;
            }

            /**
             * Handles the click event on the history pin button.
             * @param {MouseEvent} event
             */
            const handleHistoryPinButtonClick = (event) => {
                categories_history_pinned = !categories_history_pinned;

                if (categories_history_pinned) {
                    updateCategoriesHistoryIfOutdated();
                }
            }

            /**
             * Handles the click event on a history item.
             * @param {MouseEvent} event
             */
            const handleHistoryItemClick = (event) => {
                const history_item = event.currentTarget;
                if (history_item == null || !(history_item instanceof HTMLElement)) {
                    console.error("In NavCategoryHistory.handleHistoryItemClick: history_item is null.");
                    return;
                }

                const history_element_uuid = history_item.dataset.historyIndex;
                if (history_element_uuid == null) {
                    console.error("In NavCategoryHistory.handleHistoryItemClick: history_index is null.");
                    return;
                }

                const inner_category = $current_cluster.CategoryUsageHistory.UUIDHistory.getElementByUUID(history_element_uuid);
                if (inner_category == null) {
                    console.error("In NavCategoryHistory.handleHistoryItemClick: inner_category is null for uuid <%s>.", history_element_uuid);
                    return;
                }

                // TODO: Find a way to ubiquitously propagate the selected category. And remember, the idea
                // is that the behavior of the button click is going to be left to the active app page.
                // Meaning that this likely should be an event. but the delivery of that event needs to be
                // carefully considered.

                console.debug("Clicked on:", inner_category);
            }

        /*=====  End of DOM Event handlers  ======*/

        /**
         * Attaches category history event listeners.
         * @returns {void}
         */
        const attachCategoryHistoryEventListeners = () => {
            $current_cluster.CategoryUsageHistory.UUIDHistory.addHistoryUpdatedListener(handleCategoryHistoryChange);
        }

        /**
         * Handles changes in the Category history.
         * @returns {void}
         */
        function handleCategoryHistoryChange() {
            console.debug("In NavCategoryHistory.handleCategoryHistoryChange: Category history has changed.");
            categories_history_mutated = true;

            const is_history_visible = isCategoriesHistoryVisible();
            if (is_history_visible) {
                updateCategoriesHistoryIfOutdated();
            }
        }

        /**
         * Returns whether the categories history is visible.
         * @returns {boolean}
         */
        const isCategoriesHistoryVisible = () => {
            return categories_history_pinned || history_button_hovered;
        }

        /**
         * Updates the categories history array.
         * @returns {void}
         */
        const updateCategoriesHistory = () => {
            if ($current_category == null || $current_cluster == null) {
                console.error("In NavCategoryHistory.updateCategoriesHistory: $current_category<%O> or $current_cluster<%O> are null.", $current_category, $current_cluster);
                return;
            }

            const new_history = $current_cluster.CategoryUsageHistory.UUIDHistory.toArray();
            new_history.shift(); // Remove the current category from the history

            categories_history_array = new_history;
            categories_history_mutated = false;
        }

        /**
         * Updates the categories history if required.
         * @returns {void}
         */
        const updateCategoriesHistoryIfOutdated = () => {
            if (categories_history_mutated) {
                updateCategoriesHistory();
            }
        }

        /**
         * Removes category history event listeners.
         * @returns {void}
         */
        const removeCategoryHistoryEventListeners = () => {
            if ($current_cluster === null) {
                console.debug("In NavCategoryHistory.removeCategoryHistoryEventListeners: $current_cluster is null.");
                return;
            }

            $current_cluster.CategoryUsageHistory.UUIDHistory.removeHistoryUpdatedListener(handleCategoryHistoryChange);
        }

    /*=====  End of Methods  ======*/
    
</script>

{#if $current_category != null && $current_cluster != null}
    <div id="ldn-nav-category-history">
        <button id="ldn-nch-history-pin-btn"
            on:mouseenter={handleHistoryButtonMouseEnter}
            on:mouseleave={handleHistoryButtonMouseLeave}
            on:click={handleHistoryPinButtonClick}
        >
            <h3 id="ldn-nch-hpb-current-category-name">
                {$current_category.name}
            </h3>
        </button>
        {#if categories_history_array.length > 0}
            <menu id="ldn-nch-category-history-container"
                class="dungeon-scroll"
                class:history-always-visible={categories_history_pinned}
            >
                {#each categories_history_array as history_item}
                    <li class="ldn-nch-history-record"
                        role="menuitem"
                        data-history-index={history_item.uuid}
                        on:click={handleHistoryItemClick}
                    >
                        <p class="ldn-nch-hr-category-name">
                            {history_item.name}
                        </p>
                    </li>
                {/each}
            </menu>
        {/if}
    </div>
{/if}

<style>
    #ldn-nav-category-history {
        --ldn-nav-ch-font-size: var(--font-size-1);

        position: relative;
        display: grid;
        place-items: center;
        transition: background 0.3s ease-out;

        &:hover button#ldn-nch-history-pin-btn, 
        &:has(menu#ldn-nch-category-history-container.history-always-visible) button#ldn-nch-history-pin-btn {
            background: var(--main-dark);
            color: var(--body-bg-color);
            z-index: var(--z-index-t-4);
        }

        &:hover menu#ldn-nch-category-history-container,
        & menu#ldn-nch-category-history-container.history-always-visible {
            visibility: visible;
            translate: 0;
        }
    }

    button#ldn-nch-history-pin-btn {
        width: 100%;
        background: var(--body-bg-color);
        padding-block: calc(var(--spacing-1) + 1px);
        padding-inline: var(--spacing-1);
        border-radius: var(--border-radius);
        transition: background 0.3s ease-out;
        color: var(--main-dark);

        &:hover {
            border-radius: var(--border-radius) var(--border-radius) 0 0;
        }

        & h3 {
            font-family: var(--font-read);
            font-size: var(--ldn-nav-ch-font-size);
            text-transform: none;
            line-height: normal;
            color: inherit;
        }
    }

    :global(#libery-dungeon-navbar.navbar-ethereal button#ldn-nch-history-pin-btn) {
        background-color: hsl(from var(--body-bg-color) h s l / 0.2);
    }
    
    /*=============================================
    =            Drop down menu            =
    =============================================*/
    
        menu#ldn-nch-category-history-container {
            --ldn-category-history-item-height: 48px;

            position: absolute;
            overflow-y: auto;
            width: 100%;
            height: min(calc(6 * var(--ldn-category-history-item-height)), 45dvh);
            background: hsl(from var(--grey-8) h s calc(l * 0.7) / 0.88);
            font-size: calc(var(--ldn-nav-ch-font-size) * 0.9);
            top: 100%;
            right: 0;
            translate: 0 -10%;
            visibility: hidden;
            overscroll-behavior-y: contain;
            transition: background 0.3s ease-out 0.2s, translate 0.2s ease-out;
            z-index: var(--z-index-t-3);
        }

        li.ldn-nch-history-record {
            height: 2.4em;
            line-height: 1;
            padding: 0 0;

            &:not(:last-child) {
                border-bottom: 1px solid var(--grey-9);
            }

            &:hover {
                background: hsl(from var(--main-9) h s l / 0.1);
            }

            &:has(p):hover {
                color: var(--main-2);
            }

            & > p {
                cursor: default;
                display: flex;
                width: 100%;
                font-size: inherit;
                font-weight: 600;
                height: 100%;
                align-items: center;
                padding: 0 var(--spacing-1);
            }
        }
    
    /*=====  End of Drop down menu  ======*/
    
    
</style>