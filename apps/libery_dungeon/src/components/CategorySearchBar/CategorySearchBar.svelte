<script>
    import { createEventDispatcher } from "svelte";
    import { searchCategories } from "@models/WorkManagers";
    import { onMount } from "svelte";
    import { current_category } from "@stores/categories_tree";
    import { current_cluster } from "@stores/clusters";
    import { browser } from "$app/environment";

    /*=============================================
    =            Properties            =
    =============================================*/

        /**
         * The search bar input
         * @type {HTMLInputElement}
         */
        let search_bar;

        /**
         * The search query
         * @type {string}
         */
        export let search_query = "";
    
        /** @type {boolean} whether the search bar is focused or not */
        export let autofocus = false;

        const search_event_dispatcher = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    onMount(() => {
        if (autofocus && browser) {
            focusSearchBar();
            // setTimeout(() => {
            // }, 400);
        }
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        export const focusSearchBar = () => {
            if (search_bar === null) return;

            const search_bar_style =  window.getComputedStyle(search_bar);

            if (search_bar_style.visibility === "visible") {
                search_bar.focus();
            }
        }

        const handleSearch = async () => {
            if ($current_category === null) {
                throw new Error("On components/CategorySearchBar.svelte: Attempted to search categories without a current category. which is required to retrieve the cluster id");
            }

            search_bar.blur();

            if (search_query === "") {
                return;
            }

            const response = await searchCategories(search_query, $current_category.ClusterUUID, $current_cluster.DownloadCategoryID);

            const search_results = response.data;

            search_event_dispatcher("search-results", {
                results: search_results,
                search_query,
            });

            search_field_data.clear();
        }

        /**
         * Handles the search bar keydown event.
         * @param {KeyboardEvent} e
         */
        const handleSearchBarKeydown = e => {
            if (handleSearchBarCommands(e)) {
                return;
            }
        }

        /**
         * Handles the search bar keyup event.
         * @param {KeyboardEvent} e
         */
        const handleSearchBarKeyup = e => {
        }

        /**
         * Handles the search bar commands. The HotkeyBinder ignores events that occur on input and textarea elements. 
         * Returns whether the event was handled.
         * @param {KeyboardEvent} e
         * @returns {boolean}
         */
        const handleSearchBarCommands = e => {
            let event_handled = false;

            if (e.key === "Escape") {
                e.preventDefault();
                search_bar.blur();
                event_handled = true;
            }

            if (e.key === "e" && e.ctrlKey || e.key === "Enter") {
                let search_query_valid = search_bar.checkValidity();

                if (!search_query_valid) {
                    search_bar.reportValidity();
                    return;
                }

                e.preventDefault();
                handleSearch();
                event_handled = true;
            }

            return event_handled;
        }
    
    /*=====  End of Methods  ======*/

</script>

<label class="category-search-bar-wrapper dungeon-input">
    <span class="dungeon-label">
        search
    </span>
    <input 
        bind:this={search_bar}
        bind:value={search_query}
        type="text"
        id="category-search-bar-input"
        placeholder="category name"
        autocomplete="off"
        on:keydown={handleSearchBarKeydown}
        on:keyup={handleSearchBarKeyup}
        on:keypress={handleSearchBarKeyup}
        required
    >
</label>

<style>
    span.dungeon-label {
        color: var(--main-dark-color-5);
    }
</style>