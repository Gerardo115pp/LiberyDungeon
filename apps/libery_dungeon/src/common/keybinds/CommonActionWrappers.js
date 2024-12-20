import { HOTKEY_NULLISH_CAPTURE_HANDLER, HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
import { toggleHotkeysSheet } from "@stores/layout";
import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
import { global_hotkey_action_triggers, global_search_hotkeys } from "@app/config/hotkeys_config";
import { jaroWinkler } from "@libs/LiberySmetrics/jaro";
import { UIReference } from "@libs/LiberyFeedback/lf_models";
import { common_action_groups } from "./CommonActionsName";
import { HotkeyData } from "@libs/LiberyHotkeys/hotkeys";
import { linearCycleNavigationWrap } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";

/**
 * Adds the common show hotkeys table action to a given hotkeys context.
 * @param {HotkeysContext} hotkeys_context
 * @returns {void}
 */
export const wrapShowHotkeysTable = (hotkeys_context) => {
    hotkeys_context.register(global_hotkey_action_triggers.TOGGLE_HOTKEYS_SHEET, toggleHotkeysSheet, { 
        description: `<${HOTKEYS_GENERAL_GROUP}>Open/Close the hotkeys cheat sheet.`
    });
}

/**
 * @template T
* @callback SearchResultsUpdateCallback
 * @param {T} search_result
 * @returns {void}
*/

/**
 * Sets the search results wrapper required data.
* @callback SearchResultsWrapperSetup
 * @param {import("@libs/LiberyHotkeys/hotkeys_context").default} hotkeys_context
*/

/**
 * A class that abstracts search functionality using hotkeys. Inspired by the '/' search hotkey in Vim. Receives a HotkeysContext and a result update callback that takes T as a parameter and returns void.
 * By default converts T into searchable strings using the String object as a function(objects can override how String converts them into primitive string by implementing the [Symbol.toPrimitive] method).
 * @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Symbol/toPrimitive
 * @template T
 */
export class SearchResultsWrapper {

    /**
     * The search results array.
     * @type {T[]}
     */
    #search_results;

    /**
     * The searchable items.
     * @type {T[]}
     */
    #search_pool;

    /**
     * The current search result index.
     * @type {number}
     */
    #current_search_index;

    /**
     * The search results options.
     * @type {SearchResultsOptions}
     * @typedef {Object} SearchResultsOptions
     * @property {import('@libs/LiberyFeedback/lf_models').UIReference} ui_search_result_reference
     * @property {string[] | string} search_hotkey The hotkey to trigger a search.
     * @property {string[] | string} search_next_hotkey The hotkey to go to the next search result.
     * @property {string[] | string} search_previous_hotkey The hotkey to go to the previous search result.
     * @property {number} minimum_similarity The minimum similarity score for a search result to be considered a match. Default is 0.5.
     * @property {boolean} case_insensitive If true, the search string is converted to the same case as the search pool items.
     * @property {boolean} boost_exact_inclusion If true, strings that contain the exact match are(even if the string is much larger than the search query) are boosted.
     * @property {boolean} allow_member_similarity_checking If true, when a search string is considerably smaller than a item in the search pool. the search query is compared against the members(substrings splitted by a space) and the highest similarity is used. This is obviously slightly slower.
     * @property {() => void} no_results_callback
     * @property {import("@libs/LiberyHotkeys/hotkeys").HotkeyCaptureCallback} [search_typing_hotkey_handler] A callback that receives a search query as is been typed.
     */
    #search_options;

    /**
     * The search results update callback.
     * @type {SearchResultsUpdateCallback<T>}
     */
    #search_results_update_callback;

    /**
     * Item to string conversion function.
     * @type {function(T): string}
     */
    #item_to_string;

    /**
    * @param {HotkeysContext} hotkeys_context
    * @param {T[]} search_pool
    * @param {SearchResultsUpdateCallback<T>} search_results_update_callback
    * @param {SearchResultsOptionsParams} [search_options]
     * @typedef {Object} SearchResultsOptionsParams
     * @property {import('@libs/LiberyFeedback/lf_models').UIReference} [ui_search_result_reference]
     * @property {string[] | string} [search_hotkey] The hotkey to trigger a search.
     * @property {string[] | string} [search_next_hotkey] The hotkey to go to the next search result.
     * @property {string[] | string} [search_previous_hotkey] The hotkey to go to the previous search result
     * @property {number} [minimum_similarity] The minimum similarity score for a search result to be considered a match. Default is 0.5.
     * @property {boolean} [force_same_case=true]
     * @property {boolean} [boost_exact_inclusion=false]
     * @property {boolean} [allow_member_similarity_checking=true]
     * @property {() => void} [no_results_callback]
     * @property {import("@libs/LiberyHotkeys/hotkeys").HotkeyCaptureCallback} [search_typing_hotkey_handler] A callback that receives a search query as is been typed.
     */
    constructor(hotkeys_context, search_pool, search_results_update_callback, search_options) {
        if (search_options === undefined) {
            search_options = {};
        }

        search_options.ui_search_result_reference

        this.#search_options = { ...DEFAULT_SEARCH_OPTIONS, ...search_options };

        this.#search_results = [];

        this.#search_pool = search_pool;

        this.#current_search_index = 0;

        this.#search_results_update_callback = search_results_update_callback;

        this.#item_to_string = (item) => String(item);

        this.wrapSearchHotkeys(hotkeys_context);
    }

    /**
     * Creates a match array out of a given search string and a given search array.
     * @param {T[]} search_array
     * @param {string} search_string
     * @returns {MatchScore[]}
     * @typedef {Object} MatchScore
     * @property {number} similarity The similarity score. 0 is not a match at all, 1 is a perfect match.
     * @property {T} item The item that was matched.
     * @property {string} item_string The string representation of the item.
     */
    #createMatchArray(search_array, search_string) {
        /** @type {MatchScore[]} */
        let match_array = search_array.map((item) => {
            let item_string = this.getItemString(item);

            const similarity = this.#getStringSimilarity(item_string, search_string);

            return { 
                similarity,
                item,
                item_string
            };
        });

        match_array = match_array.filter((match) => match.similarity >= this.#search_options.minimum_similarity);

        return match_array;
    }

    /**
     * Drops the current search results, current search index.
     * @returns {void}
     */
    dropSearchState() {
        this.#search_results = [];
        this.#current_search_index = 0;
    }

    /**
     * Returns the similarity of a given string against a search string.
     * @param {string} item_string
     * @param {string} search_string
     * @returns {number}
     */
    #getStringSimilarity = (item_string, search_string) => {
        let similarity = 0;
        let similarity_boost = 0; 

        const use_member_similarity = this.#search_options.allow_member_similarity_checking && item_string.length > (1.5 * search_string.length);

        if (use_member_similarity) {

            for(const item_string_member of item_string.split(' ')) {
                if (item_string_member.length < (search_string.length * 0.7)) continue;

                const member_similarity = jaroWinkler(item_string_member, search_string)


                if (member_similarity > similarity) {
                    similarity = member_similarity;
                }
            }

        } else {
            similarity = jaroWinkler(item_string, search_string)
        }

        if (this.#search_options.boost_exact_inclusion && search_string.length > 3) {
            const includes_search_string = item_string.toLowerCase().includes(search_string.toLowerCase());

            if (includes_search_string) {
                similarity_boost = search_string.length * 0.1;
            }
        }

        similarity = Math.min(1, similarity + similarity_boost);


        return similarity
    }

    /**
     * Returns a string representation of the passed item.
     * @param {T} item
     * @returns {string}
     */
    getItemString = item => {
        let item_string = this.#item_to_string(item);
        
        if (this.#search_options.case_insensitive) {
            item_string = item_string.toLowerCase();
        }

        return item_string;
    }

    /**
     * handles a new search term captured by the search hotkey.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #searchHotkeyHandler = (event, hotkey) => {
        if (this.#search_pool.length == 0) return;

        if (!hotkey || hotkey.HotkeyType !== HotkeyData.HOTKEY_TYPE__CAPTURE || !hotkey.HasMatch) {
            console.error("Search hotkey handler called without a valid hotkey: ", hotkey);
            return;
        }

        let search_string = hotkey.MatchMetadata?.CaptureMatch;

        if (!search_string) {
            console.error("Search hotkey handler called without a valid search string: ", hotkey);
            this.#search_options.no_results_callback();
            return;
        }

        if (this.#search_options.case_insensitive) {
            search_string = search_string.toLowerCase();
        }

        /**
         * Whether the results update callback should be triggered
         */
        let found_new_results = false;

        if (search_string.length > 3) {
            found_new_results = this.search(this.#search_pool, search_string);
        } else if (search_string.length === 1) {
            found_new_results = this.searchPrefix(this.#search_pool, search_string);
        } else {
            found_new_results = this.searchInclude(this.#search_pool, search_string);
        }

        if (found_new_results) {
            this.#search_results_update_callback(this.#search_results[this.#current_search_index]);
        } else {
            this.#search_options.no_results_callback();
        }
    }

    /**
     * Handles the 'search next' hotkey.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #searchNextHotkeyHandler = (event, hotkey) => {
        let wrapped_value = linearCycleNavigationWrap(this.#current_search_index, this.#search_results.length - 1, 1);

        this.#current_search_index = wrapped_value.value;

        this.#search_results_update_callback(this.#search_results[this.#current_search_index]);
    }

    /**
     * Handles the 'search previous' hotkey.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #searchPreviousHotkeyHandler = (event, hotkey) => {
        let wrapped_value = linearCycleNavigationWrap(this.#current_search_index, this.#search_results.length - 1, -1);

        this.#current_search_index = wrapped_value.value;

        this.#search_results_update_callback(this.#search_results[this.#current_search_index]);
    }
     
    /**
     * Sets the function used to convert items to searchable strings.
     * @param {function(T): string} item_to_string
     * @returns {void}
     */
    setItemToStringFunction(item_to_string) {
        this.#item_to_string = item_to_string;
    }

    /**
     * a list of search results.
     * @type {T[]}
     */
    get SearchResults() {
        return this.#search_results;
    }

    /**
     * Searches a given array for a given search string. utilizes the Jaro-Winkler to get a similarity score 
     * and sorts the results by that similarity omits results that are below the minimum similarity threshold.
     * returns whether the search yielded any new results.
     * @param {T[]} search_array
     * @param {string} search_string
     * @returns {boolean}
     */
    search(search_array, search_string) {
        if (search_array.length == 0) {
            console.warn("Search array is empty.");
            return false;
        }

        /**
         * @type {T[]}
         */
        const new_search_results = [];

        const match_array = this.#createMatchArray(search_array, search_string);

        match_array.sort((a, b) => b.similarity - a.similarity);

        for (const match of match_array) {
            new_search_results.push(match.item);
        }

        if (new_search_results.length == 0) {
            console.warn(`No search results found for ${search_string}`);
            return false;
        }

        this.#search_results = new_search_results;
        this.#current_search_index = 0;

        return true;
    }

    /**
     * Searches a given array for a given search string. Checking only if the items are prefixed with the search string. Ideal for short search_strings(1-2 characters).
     * returns whether the search yielded any new results.
     * @param {T[]} search_array
     * @param {string} search_string
     * @returns {boolean}
     */
    searchPrefix(search_array, search_string) {
        if (search_array.length == 0) {
            console.warn("Search array is empty.");
            return false;
        }

        /**
         * @type {T[]}
         */
        const new_search_results = [];

        for (const item of search_array) {
            let item_string = this.getItemString(item);

            if (this.#search_options.case_insensitive) {
                item_string = item_string.toLowerCase();
            }

            if (item_string.startsWith(search_string)) {
                new_search_results.push(item);
            }
        }

        if (new_search_results.length == 0) {
            console.warn(`No search results found for ${search_string}`);
            return false;
        }

        this.#search_results = new_search_results;
        this.#current_search_index = 0;

        return true;
    }

    /**
     * Searches a given array for a given search string. Checks if the items include the search string and adding perfect matches at the top of the search result. Ideal for short-medium search strings(2-4 characters).
     * @param {T[]} search_array
     * @param {string} search_string
     * @returns {boolean}
     */
    searchInclude(search_array, search_string) {
        if (search_array.length == 0) {
            console.warn("Search array is empty.");
            return false;
        }

        /**
         * @type {T[]}
         */
        const new_search_results = [];

        for (const item of search_array) {
            let item_string = this.getItemString(item);

            if (this.#search_options.case_insensitive) {
                item_string = item_string.toLowerCase();
            }

            if (item_string === search_string) {
                new_search_results.unshift(item);
            } else if (item_string.includes(search_string)) {
                new_search_results.push(item);
            }
        }

        if (new_search_results.length == 0) {
            console.warn(`No search results found for ${search_string}`);
            return false;
        }

        this.#search_results = new_search_results;
        this.#current_search_index = 0;

        return true;
    }

    /**
     * Updates the search pool.
     * @param {T[]} search_pool
     * @returns {void}
     */
    updateSearchPool(search_pool) {
        this.#search_pool = search_pool;
    }

    /**
     * wraps the search functionality around a given hotkeys context.
     * @param {HotkeysContext} hotkeys_context
     * @returns {void}
     */
    wrapSearchHotkeys(hotkeys_context) {
        if (this.#search_options.search_hotkey != undefined || this.#search_options.search_hotkey !== "") {
            let search_hotkey_combo = Array.isArray(this.#search_options.search_hotkey) ? this.#search_options.search_hotkey[0] : this.#search_options.search_hotkey;

            search_hotkey_combo = `${search_hotkey_combo} \\s`; // creates a string capture hotkey.

            hotkeys_context.register(search_hotkey_combo, this.#searchHotkeyHandler.bind(this), {
                description: `${common_action_groups.NAVIGATION}Search for ${this.#search_options.ui_search_result_reference.EntityName}.`,
                capture_hotkey_callback: this.#search_options.search_typing_hotkey_handler,
            });
        }

        if (this.#search_options.search_next_hotkey != undefined || this.#search_options.search_next_hotkey !== "") {
            hotkeys_context.register(this.#search_options.search_next_hotkey, this.#searchNextHotkeyHandler.bind(this), {
                description: `${common_action_groups.NAVIGATION}Go to the next ${this.#search_options.ui_search_result_reference.EntityName}.`
            });
        }

        if (this.#search_options.search_previous_hotkey != undefined || this.#search_options.search_previous_hotkey !== "") {
            hotkeys_context.register(this.#search_options.search_previous_hotkey, this.#searchPreviousHotkeyHandler.bind(this), {
                description: `${common_action_groups.NAVIGATION}Go to the previous ${this.#search_options.ui_search_result_reference.EntityName}.`
            });
        }
    }
}

/**
 * Default options for the SearchResultsWrapper class.
 * @type {SearchResultsOptions}
 */
const DEFAULT_SEARCH_OPTIONS = {
    ui_search_result_reference: new UIReference("search result", "search results"),
    search_hotkey: global_search_hotkeys.SEARCH,
    search_next_hotkey: global_search_hotkeys.SEARCH_NEXT,
    search_previous_hotkey: global_search_hotkeys.SEARCH_PREVIOUS,
    minimum_similarity: 0.5,
    case_insensitive: true,
    boost_exact_inclusion: false,
    allow_member_similarity_checking: true,
    search_typing_hotkey_handler: HOTKEY_NULLISH_CAPTURE_HANDLER,
    no_results_callback: () => {}
}