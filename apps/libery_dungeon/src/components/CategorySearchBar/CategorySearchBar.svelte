<script>
    import { createEventDispatcher } from "svelte";
    import { searchCategories } from "@models/WorkManagers";
    import { onMount } from "svelte";
    import { current_category } from "@stores/categories_tree";
    import { current_cluster } from "@stores/clusters";
    import { browser } from "$app/environment";
    import { emitLabeledError } from "@libs/LiberyFeedback/lf_utils";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";

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

        /**
         * Requests categories with a similar name(fuzzy logic) to the value of search_query.
         * Dispatches a "search-results" event with the results and returns true, if the search_query was valid
         * and yields results; otherwise returns false.
         * @returns {Promise<boolean>}
         */
        const handleSearch = async () => {
            if ($current_category === null) {
                throw new Error("On components/CategorySearchBar.svelte: Attempted to search categories without a current category. which is required to retrieve the cluster id");
            }

            const error_context = new VariableEnvironmentContextError("In CategorySearchBar.handleSearch")

            error_context.addVariable("search_query", search_query);

            if (search_query === "") {

                const labeled_err = new LabeledError(
                    error_context,
                    "Search query is empty, press Esc to cancle the search.",
                    lf_errors.ERR_HUMAN_ERROR
                )
                emitLabeledError(labeled_err);

                return false;
            }

            const search_results = await searchCategories(search_query, $current_category.ClusterUUID, $current_cluster.DownloadCategoryID);

            if (search_results.length > 0) {
                search_bar.blur();

                search_event_dispatcher("search-results", {
                    results: search_results,
                    search_query,
                });

                return true;
            }

            const labeled_err = new LabeledError(
                error_context,
                `No categoty matches '${search_query}'`,
                lf_errors.ERR_NO_CONTENT
            );
            emitLabeledError(labeled_err);

            return false;
        }

        /**
         * Handles the search bar keydown event.
         * @param {KeyboardEvent} e
         * @returns {Promise<void>}
         */
        const handleSearchBarKeydown = async e => {
            const command_was_handled = await handleSearchBarCommands(e);

            if (command_was_handled) {
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
         * @returns {Promise<boolean>}
         */
        const handleSearchBarCommands = async e => {
            let event_handled = false;

            if (e.key === "Escape") {
                e.preventDefault();
                search_bar.blur();
                event_handled = true;
            }

            if (e.key === "e" && e.ctrlKey || e.key === "Enter") {
                if (search_query === "") {
                    search_bar.blur();

                    return true;
                }                

                let search_query_valid = search_bar.checkValidity();

                if (!search_query_valid) {
                    search_bar.reportValidity();
                    return event_handled;
                }

                e.preventDefault();
                event_handled = await handleSearch();
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