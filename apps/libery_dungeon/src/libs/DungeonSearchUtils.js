import { jaroWinkler } from "./LiberySmetrics/jaro"


/*=============================================
=            Properties            =
=============================================*/

    /**
     * A set of options useful to customize search behavior
    * @typedef {Object} SearchResultsOptions
    * @property {number} [minimum_similarity] The minimum similarity score for a search result to be considered a match. Default is 0.5.
    * @property {boolean} [case_insensitive] If true, the search string is converted to the same case as the search pool items.
    * @property {boolean} [boost_exact_inclusion] If true, strings that contain the exact match are(even if the string is much larger than the search query) are boosted.
    * @property {boolean} [allow_member_similarity_checking] If true, when a search string is considerably smaller than an item in the search pool. the search query is compared against the members(substrings splitted by a space) and the highest similarity is used. This is obviously slightly slower.
    */

    /**
     * @template {{toString: () => string}} T
    * @typedef {Object} MatchScore
    * @property {number} similarity The similarity score. 0 is not a match at all, 1 is a perfect match.
    * @property {T} item The item that was matched.
    * @property {string} item_string The string representation of the item.
    */

    const DEFAULT_SEARCH_OPTIONS = {
        minimum_similarity: 0.5,
        case_insensitive: true,
        boost_exact_inclusion: false,
        allow_member_similarity_checking: true,
    }

/*=====  End of Properties  ======*/

/**
 * Creates a match array out of a given search string and a given search array.
 * @template {{toString: () => string}} T
 * @param {T[]} search_array
 * @param {string} search_string
 * @param {SearchResultsOptions} [options]
 * @returns {MatchScore<T>[]}
 */
const createMatchArray = (search_array, search_string, options) => {
    const search_options = options === undefined ? DEFAULT_SEARCH_OPTIONS : fillDefaultSearchOptions(options);

    if (search_options.case_insensitive) {
        search_string = search_string.toLowerCase();
    }

    /** @type {MatchScore<T>[]} */
    let match_array = search_array.map((item) => {
        let item_string = item.toString();

        if (search_options.case_insensitive) {
            item_string = item_string.toLowerCase();
        }

        const similarity = getStringSimilarity(item_string, search_string);

        return { 
            similarity,
            item,
            item_string
        };
    });

    match_array = match_array.filter((match) => match.similarity >= search_options.minimum_similarity);

    return match_array;
}

/**
 * Search for items in an item array using the most optimal method depending on the search string.
 * @template {{toString: () => string}} T
 * @param {T[]} search_array
 * @param {string} search_string
 * @param {SearchResultsOptions} [options]
 * @returns {T[]}
 */
export const dungeonSearch = (search_array, search_string, options) => {
    if (search_string === "") {
        return [];
    }

    const search_options = options === undefined ? DEFAULT_SEARCH_OPTIONS : fillDefaultSearchOptions(options);

    /**
     * @type {T[]}
     */
    let search_results = []

    if (search_string.length > 3) {
        search_results = searchSimilarity(search_array, search_string, search_options);
    } else if (search_string.length === 1) {
        search_results = searchPrefix(search_array, search_string, search_options);
    } else {
        search_results = searchInclude(search_array, search_string, search_options);
    }

    return search_results
}

/**
 * Fills the given options object with the ones on the DEFAULT_SEARCH_OPTIONS if they are missing.
 * @param {SearchResultsOptions} nullish_options
 * @returns {typeof DEFAULT_SEARCH_OPTIONS}
 */
const fillDefaultSearchOptions = nullish_options => {
    const the_full_options = {
        ...DEFAULT_SEARCH_OPTIONS,
        ...nullish_options
    }

    return the_full_options
}

/**
 * Returns the similarity of a given string against a search string.
 * @param {string} item_string
 * @param {string} search_string
 * @param {SearchResultsOptions} [options]
 * @returns {number}
 */
export const getStringSimilarity = (item_string, search_string, options) => {
    const search_options = options === undefined ? DEFAULT_SEARCH_OPTIONS : fillDefaultSearchOptions(options);

    let similarity = 0;
    let similarity_boost = 0; 

    const use_member_similarity = search_options.allow_member_similarity_checking && item_string.length > (1.5 * search_string.length);

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

    if (search_options.boost_exact_inclusion && search_string.length > 3) {
        const includes_search_string = item_string.toLowerCase().includes(search_string.toLowerCase());

        if (includes_search_string) {
            similarity_boost = search_string.length * 0.1;
        }
    }

    similarity = Math.min(1, similarity + similarity_boost);


    return similarity
}

/**
 * Searches a given array for a given search string. utilizes the Jaro-Winkler to get a similarity score 
 * and sorts the results by that similarity omits results that are below the minimum similarity threshold.
 * returns whether the search yielded any new results.
 * @template {{toString: () => string}} T
 * @param {T[]} search_array
 * @param {string} search_string
 * @param {SearchResultsOptions} [options]
 * @returns {T[]}
 */
export const searchSimilarity = (search_array, search_string, options) => {
    const search_options = options === undefined ? DEFAULT_SEARCH_OPTIONS : fillDefaultSearchOptions(options);

    if (search_array.length == 0) {
        console.warn("Search array is empty.");
        return [];
    }

    /**
     * @type {T[]}
     */
    const new_search_results = [];

    const match_array = createMatchArray(search_array, search_string, search_options);

    match_array.sort((a, b) => b.similarity - a.similarity);

    for (const match of match_array) {
        new_search_results.push(match.item);
    }

    if (new_search_results.length == 0) {
        console.warn(`No search results found for ${search_string}`);
        return [];
    }

    return new_search_results;
}

/**
 * Searches a given array for a given search string. Checks if the items include the search string and adding perfect matches at the top of the search result. Ideal for short-medium search strings(2-4 characters).
 * @template {{toString: () => string}} T
 * @param {T[]} search_array
 * @param {string} search_string
 * @param {SearchResultsOptions} [options]
 * @returns {T[]}
 */
export const searchInclude = (search_array, search_string, options) => {
    const search_options = options === undefined ? DEFAULT_SEARCH_OPTIONS : fillDefaultSearchOptions(options);

    if (search_array.length == 0) {
        console.warn("Search array is empty.");
        return [];
    }

    /**
     * @type {T[]}
     */
    const new_search_results = [];

    if (search_options.case_insensitive) {
        search_string = search_string.toLowerCase();
    }

    for (const item of search_array) {
        let item_string = item.toString();

        if (search_options.case_insensitive) {
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
        return [];
    }

    return new_search_results;
}

/**
 * Searches a given array for a given search string. Checking only if the items are prefixed with the search string. Ideal for short search_strings(1-2 characters).
 * returns whether the search yielded any new results.
 * @template {{toString: () => string}} T
 * @param {T[]} search_array
 * @param {string} search_string
 * @param {SearchResultsOptions} [options]
 * @returns {T[]}
 */
export const searchPrefix = (search_array, search_string, options) => {
    const search_options = options === undefined ? DEFAULT_SEARCH_OPTIONS : fillDefaultSearchOptions(options);

    if (search_array.length == 0) {
        return [];
    }

    if (search_options.case_insensitive) {
        search_string = search_string.toLowerCase();
    }

    /**
     * @type {T[]}
     */
    const new_search_results = [];

    for (const item of search_array) {
        let item_string = item.toString();

        if (search_options.case_insensitive) {
            item_string = item_string.toLowerCase();
        }

        if (item_string.startsWith(search_string)) {
            new_search_results.push(item);
        }
    }

    if (new_search_results.length == 0) {
        console.warn(`No search results found for ${search_string}`);
        return [];
    }

    return new_search_results;
}
