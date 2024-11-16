import { writable } from "svelte/store";

/*=============================================
=            Component state            =
=============================================*/

    /**
     * @type {import('svelte/store').Writable<boolean>} whether the media upload tool is mounted or not.
     * @description This is currently toggled on by the NavbarPanel of MediaExplorer and toggled off by the  MediaUploadTool itself.
    */
    export const media_upload_tool_mounted = writable(false);

    /**
     * @type {import('svelte/store').Writable<boolean>} whether the Category creation tool is mounted or not.
    */
    export const category_creation_tool_mounted = writable(false);

    /**
     * Whether the Media transactions management tool is mounted or not.
     * @type {import('svelte/store').Writable<boolean>}
     */
    export const media_transactions_tool_mounted = writable(false);

    /**
     * Whether the Category Tagger tool is mounted or not.
     * @type {import('svelte/store').Writable<boolean>}
     */
    export const category_tagger_tool_mounted = writable(false);

    /**
     * Whether the current category configuration component is mounted or not.
     * @type {import('svelte/store').Writable<boolean>}
     */
    export const current_category_configuration_mounted = writable(true);

    /**
     * Whether the category search bar is been focused or not.
     * @type {typeof media_upload_tool_mounted}
     */
    export const category_search_focused = writable(false);
    
    
    /*----------  Category search  ----------*/
    
        /**
         * The results of the category search.
         * @type {import('svelte/store').Writable<import('@models/Categories').Category[]>}
         * @default []
         */
        export const category_search_results = writable([]);
        
        /**
         * The term used to search for categories.
         * @type {import('svelte/store').Writable<string>}
         * @default ""
         */
        export const category_search_term = writable("");
        
        export const resetCategorySearchState = () => {
            category_search_results.set([]);
            category_search_term.set("");
            category_search_focused.set(false);
        }

/*=====  End of Component state  ======*/



