<script>
    import { createEventDispatcher, onDestroy } from "svelte";
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

        
        /*----------  Search preferences  ----------*/

            /**
             * The name of the search preference storage key for
             * `search_includes_download_category`
             * @type {string}
             */
            const SEARCH_PREF_STORAGE_KEY = "csb-preferences"; // Abbreviation for Category Search Bar - Search Includes Download Category
        
            /**
             * Whether to include the download target of the current cluster or not.
             * @type {boolean}
             * @default false 
             */ 
            let search_includes_download_category = false;

        /**
         * Whether the component is in debug mode. this will make state information and methods available through the
         * navigator's dev tools.
         * @type {boolean}
         */
        const debug_mode = true;

        const search_event_dispatcher = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    onMount(() => {
        if (autofocus && browser) {
            focusSearchBar();
            // setTimeout(() => {
            // }, 400);
        }

        if (debug_mode) {
            debugCSB__attachDebugMethods();
        }

        globalThis.requestIdleCallback(loadSearchPreferencesFromSession);
    });

    onDestroy(() => {
        if (browser) {
            saveSearchPreferencesToSession();
        }
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
        /*=============================================
        =            Debug            =
        =============================================*/

                /**
                 * Returns the object where all the category search bar debug state is stored.
                 * @returns {Object}
                 */
                const debugCSB__getComponentDebugObject = () => {
                    if (!browser || !debug_mode) return {};

                    const DEBUG_STATE_NAME = "category_search_bar";

                    // @ts-ignore
                    if (!globalThis[DEBUG_STATE_NAME]) {
                        // @ts-ignore
                        globalThis[DEBUG_STATE_NAME] = {

                        };   
                    }

                    // @ts-ignore
                    return globalThis[DEBUG_STATE_NAME];
                }
        
                /**
                 * Attaches debug methods to the globalThis object for debugging purposes.
                 * @returns {void}
                 */
                const debugCSB__attachDebugMethods = () => {
                    if (!browser || !debug_mode) return;

                    const meg_gallery_debug_state = debugCSB__getComponentDebugObject();

                    // @ts-ignore - for debugging purposes we do not care whether the globalThis object has the method name. same reason for all other ts-ignore in this function.
                    meg_gallery_debug_state.printComponentState = debugCSB__printComponentState;

                    // @ts-ignore
                    meg_gallery_debug_state.Page = null;

                    // @ts-ignore - state retrieval functions.
                    meg_gallery_debug_state.State = {
                        getSearchPreferences,
                    }

                    // @ts-ignore - Internal method references.
                    meg_gallery_debug_state.Methods = {
                        loadSearchPreferencesFromSession,
                        saveSearchPreferencesToSession,
                    }
                }

                /**
                 * Prints the whole gallery state to the console.
                 * @returns {void}
                 */
                const debugCSB__printComponentState = () => {
                    console.log("%cCategorySearchBar State", "color: green; font-weight: bold;");
                    console.group("Properties");
                    console.log(`search_query: ${search_query}`);
                    console.log(`search_bar: %o`, search_bar);
                    console.groupEnd();
                    console.group("Search Preferences");
                    console.log(`search_preferences: %O`, getSearchPreferences());
                    console.groupEnd();
                }

                /**
                 * Attaches an arbitrary object as a globalThis.meg_timeline_states.<group_name>{...timestamp -> object }.
                 * @param {string} group_name
                 * @param {object} object_to_snapshot
                 * @returns {void}
                 */
                const debugCSB__attachSnapshot = (group_name, object_to_snapshot) => {
                    if (!browser || !debug_mode) return;

                    const stack = new Error().stack;
                    const datetime_obj = new Date();
                    const timestamp = `${datetime_obj.toISOString()}-${datetime_obj.getTime()}`;

                    const snapshot = {
                        timestamp,
                        stack,
                        object_to_snapshot,
                    }

                    const debug_object = debugCSB__getComponentDebugObject();

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

        /**
         * Applies the given search options object.
         * @param {import('./category_search_bar').CategorySearchBar_Preferences} search_preferences
         * @returns {void}
         */
        const applySearchPreferences = search_preferences => {
            search_includes_download_category = search_preferences.include_download_target;
        }

        export const focusSearchBar = () => {
            if (search_bar === null) return;

            const search_bar_style =  window.getComputedStyle(search_bar);

            if (search_bar_style.visibility === "visible") {
                search_bar.focus();
            }
        }

        /**
         * Returns the current search preferences.
         * @returns {import('./category_search_bar').CategorySearchBar_Preferences}
         */
        const getSearchPreferences = () => {
            return {
                include_download_target: search_includes_download_category,
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

            let ignore_categories = [];

            if (!search_includes_download_category) {
                ignore_categories.push($current_cluster.DownloadCategoryID);
            }
            

            let search_results = await searchCategories(search_query, $current_category.ClusterUUID, ...ignore_categories);

            search_results = search_results.filter(c => c.UUID !== $current_category.uuid);

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

        /**
         * handles the search form submit event.
         * @type {import('svelte/elements').EventHandler<SubmitEvent, HTMLFormElement>}
         */
        const handleSearchBarSubmitEvent = e => {
            e.preventDefault();
        }

        /**
         * Loads the search option preferences from Session storage.
         * @returns {void}
         */
        const loadSearchPreferencesFromSession = () => {
            let storage_preferences_string = sessionStorage.getItem(SEARCH_PREF_STORAGE_KEY);

            if (storage_preferences_string === null) return;

            /**
             * @type {import('./category_search_bar').CategorySearchBar_Preferences}
             */
            let new_preferences;

            try {
                new_preferences = /** @type {typeof new_preferences} */ JSON.parse(storage_preferences_string)
            } catch (e) {
                console.error("In @components/CategorySearchBar/CategorySearchBar.loadSearchPreferencesFromSession: Error parsing search preferences from session storage %O", e); 
                return;
            }

            console.debug(`CategorySearchBar.${loadSearchPreferencesFromSession.constructor.name}: Loaded search preferences from session storage: %O`, new_preferences);

            applySearchPreferences(new_preferences);

            return;
        }

        /**
         * Saves the search preferences to session storage.
         * @returns {void}
         */
        const saveSearchPreferencesToSession = () => {
            /**
             * @type {import('./category_search_bar').CategorySearchBar_Preferences}
             */
            let search_preferences = getSearchPreferences();

            /**
             * @type {string}
             */
            let search_preferences_string;

            try {
                search_preferences_string = JSON.stringify(search_preferences);
            } catch (e) {
                console.error("In @components/CategorySearchBar/CategorySearchBar.saveSearchPreferencesToSession: Error stringifying search preferences %O", e);
                return;
            }

            sessionStorage.setItem(SEARCH_PREF_STORAGE_KEY, search_preferences_string);

            console.debug(`CategorySearchBar.${saveSearchPreferencesToSession.constructor.name}: Saved search preferences to session storage`);
        }
    
    /*=====  End of Methods  ======*/

</script>

<search id="category-search-bar">
    <form action="none"
        class='category-search-bar-wrapper'
        on:submit={handleSearchBarSubmitEvent}
    >
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
        <fieldset class="csb-search-modifiers">
            <label class="csb-search-option--include-cluster-download-target dungeon-input">
                <span class="dungeon-label">
                    Include Filter Category
                </span>
                <input
                    type="checkbox" 
                    id="csb-include-download-target"
                    bind:checked={search_includes_download_category}
                >
            </label>
        </fieldset>
    </form>
</search>

<style>
    search#category-search-bar {
        display: contents;
    }

    form.category-search-bar-wrapper {
        display: flex;
        font-size: var(--font-size-1);
        flex-wrap: wrap;
        gap: 1em;
    }

    fieldset.csb-search-modifiers {
        display: flex;
        align-items: center;
        gap: .5em;

    }

    fieldset.csb-search-modifiers > label.dungeon-input {
        border-width: 0;
        border-color: transparent;

        & > span.dungeon-label {
            font-size: .7em;
        }

        & input[type="checkbox"]::before {
            font-size: .8em;
        }
    }

    span.dungeon-label {
        color: var(--main-dark-color-5);
    }
</style>