<script>
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
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        if ($current_cluster === null) {
            throw new Error("In NavCategoryHistory.onMount: $current_cluster is null.");
        }

        updateCategoriesHistoryIfOutdated();

        $current_cluster.CategoryUsageHistory.addHistoryUpdatedListener(handleCategoryHistoryChange);
    });

    onDestroy(() => {
        if ($current_cluster !== null) {
            $current_cluster.CategoryUsageHistory.removeHistoryUpdatedListener(handleCategoryHistoryChange);
        }
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

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

                const inner_category = $current_cluster.CategoryUsageHistory.getElementByUUID(history_element_uuid);
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
         * Handles changes in the Category history.
         * @returns {void}
         */
        function handleCategoryHistoryChange() {
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

            const new_history = $current_cluster.CategoryUsageHistory.toArray();
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
        position: relative;
        display: grid;
        place-items: center;
        height: 100%;
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
        padding-block: var(--spacing-2);
        padding-inline: var(--spacing-4);
        border-radius: var(--border-radius);
        transition: background 0.3s ease-out;
        color: var(--main-dark);

        &:hover {
            border-radius: var(--border-radius) var(--border-radius) 0 0;
        }

        & h3 {
            font-family: var(--font-read);
            font-size: var(--font-size-2);
            text-transform: none;
            line-height: 1;
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
            top: 100%;
            right: 0;
            translate: 0 -10%;
            visibility: hidden;
            overscroll-behavior-y: contain;
            transition: background 0.3s ease-out 0.2s, translate 0.2s ease-out;
            z-index: var(--z-index-t-3);
        }

        li.ldn-nch-history-record {
            line-height: 1;
            height: 48px;
            padding: 0 var(--spacing-1);

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
                height: 100%;
                align-items: center;
                padding: 0 var(--spacing-1);
            }
        }
    
    /*=====  End of Drop down menu  ======*/
    
    
</style>