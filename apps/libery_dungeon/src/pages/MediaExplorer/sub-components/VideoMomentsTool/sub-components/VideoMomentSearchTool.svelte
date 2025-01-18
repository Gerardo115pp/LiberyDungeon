<script>
    import { dungeonSearch } from '@libs/DungeonSearchUtils';

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The video moments array to search in.
         * @type {import('@models/Metadata').VideoMomentIdentity[]} 
         */ 
        export let the_video_moments;

        /**
         * The moments search input.
         * @type {HTMLInputElement}
         */
        let the_search_input;

        /**
         * Whether the search query is valid.
         * @type {boolean}
         */
        let search_query_valid = false;

        /**
         * The search query.
         * @type {string}
         */
        let search_query = "";
        
        /*----------  Event handlers  ----------*/
        
            /**
             * called when a new set of results is generated for a search query
             * the user typed in the search bar.
             * @type {(search_results: import('@models/Metadata').VideoMomentIdentity[]) => void}
             */ 
            export let onSearchResults = (search_results) => {};

        
    /*=====  End of Properties  ======*/

    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * triggers a search with the current query value and calls the search results callback.
         * @param {string} search_query
         */
        const handleMomentSearchProcess = (search_query) => {
            search_query = search_query.trim();

            const search_results = searchMoments(the_video_moments, search_query)

            onSearchResults(search_results);
        }

        /**
         * Handles the keydown event from the_search_input.
         * @param {KeyboardEvent} event
         */
        const handleKeyDown = event => {
            if (the_search_input == null) {
                return
            }

            if (event.key === "Escape") the_search_input.blur();

            if (event.key === "Enter") {
                event.preventDefault();
                
                if (!search_query_valid) {
                    the_search_input.reportValidity();
                    return;
                }

                handleMomentSearchProcess(search_query);

                the_search_input.blur();
            }
        }

        /**
         * Handles the keyup event from the_search_input.
         * @param {KeyboardEvent} event
         */
        const handleKeyUp = event => {
            if (the_search_input == null) return;

            event.preventDefault();

            if (the_search_input.validationMessage !== "") {
                the_search_input.setCustomValidity("");
            }

            search_query_valid = the_search_input.checkValidity();
        } 

        /**
         * Returns moments in the given array that match the given query string.
         * @param {import('@models/Metadata').VideoMomentIdentity[]} search_pool
         * @param {string} search_query
         * @returns {import('@models/Metadata').VideoMomentIdentity[]}
         */
        const searchMoments = (search_pool, search_query) => {
            const search_results = dungeonSearch(search_pool, search_query, {
                boost_exact_inclusion: true,
                case_insensitive: true,
                minimum_similarity: 0.8
            });

            return search_results;
        }
    
    /*=====  End of Methods  ======*/
    
    
</script>

<div id="diviext-moments-search-bar-wrapper">
    <label class="dungeon-input">
        <span class="dungeon-label">
            Search for moments
        </span>
        <input 
            id="diviext-moments-search-bar"
            bind:this={the_search_input}
            type="text"
            bind:value={search_query}
            on:keydown={handleKeyDown}
            on:keyup={handleKeyUp}
            spellcheck="true"
            required
        />
    </label>
</div>

<style>
    #diviext-moments-search-bar-wrapper {
        & label.dungeon-input {
            font-size: inherit;
            padding-inline: 1em;
            padding-block: .4em;
            line-height: 1;
        }
        
        & label.dungeon-input input {
            font-size: inherit;
            color: var(--grey-3);
        }
    }

    @supports (color: rgb( from white r g b / 1)) {
        #diviext-moments-search-bar-wrapper label {
            padding: var(--spacing-1) var(--spacing-2);
            background: hsl(from var(--grey) h s l / 0.8);
        }
    }
</style>