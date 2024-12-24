import { get, writable } from 'svelte/store';
import { 
    getTaggedMedias,
    NULLISH_MEDIA,
    isNullishMedia
} from '@models/Medias';
import { TagContentCache } from '@models/WorkManagers';

const TAGGED_CONTENT_PAGE_SIZE = 100;


/*=============================================
=            Properties            =
=============================================*/

    /**
     * The dungeon tags the medias displayed by the medias viewer must have.
     * @type {import('svelte/store').Writable<import('@models/DungeonTags').DungeonTag[]>}
     */
    export const mv_filtering_tags = writable([]);

    /**
     * A cache for the tag mode content retrieved by different sets of filtering tags.
     * @type {TagContentCache}
     */
    const tag_mode_content_cache = new TagContentCache();

    /**
     * The index corresponding to the active media on the tagged medias content list.
     * @type {import('svelte/store').Writable<number>}
     */
    export const active_tag_content_media_index = writable(0);

    /**
     * The medias tagged by the filtering tags.
     * @type {import('svelte/store').Writable<import('@models/Medias').Media[]>}
     */
    export const mv_tagged_content = writable([]);

    /**
     * Association of filtering_tags and an active index. used to restore the same media when that cache is loaded again.
     * @type {Map<string, number>}
     */
    const filtering_tags_active_index_cache = new Map();

    /**
     * the content cacher for the current filtering tags.
     * @type {import('@models/WorkManagers').TaggedContentCacher | null}
     */
    let filtering_medias_content_cacher = null;

    /**
     * Whether the tag_mode is enabled. This store is used for streamlining reactivity but tagModeEnabled is the actual reliable source of truth.
     * @type {import('svelte/store').Writable<boolean>}
     */
    export const mv_tag_mode_enabled = writable(false);

    /**
     * The total amount of medias available for the current filtering tags
     * @type {import('svelte/store').Writable<number>}
     */
    export const mv_tag_mode_total_content = writable(NaN)

    /**
     * The currently loaded content.
     * @type {number}
     */
    let mv_current_tagged_content_loaded = 0;

/*=====  End of Properties  ======*/


/*=============================================
=            Methods            =
=============================================*/
    // functions that start with 'tagMode_' were thought to be exported even if they are currently not. If you want to export a function that doesn't start with 'tagMode_', 
    // Create a wrapper function that calls the function you want to export. DO NOT EXPORT FUNCTIONS THAT DON'T START WITH 'tagMode_'

    /**
     * Sets a TaggerContentCacher as active setting all necessary variables and stores to the appropriate values from the given cacher.
     * @param {import('@models/WorkManagers').TaggedContentCacher} cacher
     * @returns {void}
     */
    const tagMode_activateContentCacher = cacher => {
        const cacher_filtering_tags = cacher.ContentTags;

        if (cacher_filtering_tags.length === 0) return;

        const new_tagged_content = cacher.getSequentialContent();

        console.log("new_tagged_content:", new_tagged_content);

        mv_filtering_tags.set(cacher_filtering_tags);
        mv_tagged_content.set(new_tagged_content);
        filtering_medias_content_cacher = cacher;
        mv_tag_mode_total_content.set(cacher.TotalMedias);
        mv_current_tagged_content_loaded = new_tagged_content.length;
        mv_tag_mode_enabled.set(true);
        active_tag_content_media_index.set(0);

        // @ts-ignore
        globalThis.content_cacher = cacher;
        // @ts-ignore
        globalThis.loadRandomTaggedContentPage = loadRandomTaggedContentPage;
    }

    /**
     * Changes the tagged content using a new group of tags.
     * @param {import('@models/DungeonTags').DungeonTag[]} new_tag_group 
     * @returns {Promise<void>}
     */
    export const tagMode_changeFilteringTags = async new_tag_group  => {
        if (new_tag_group.length === 0) {
            return tagMode_disableTagMode();
        }

        mv_filtering_tags.set(new_tag_group);

        let content_cacher = tag_mode_content_cache.getTagContentCacher(new_tag_group);

        if (content_cacher !== undefined) {
            return tagMode_activateContentCacher(content_cacher);
        }

        const paginated_tagged_content = await getTaggedContentPage(new_tag_group, 1);
        console.log("paginated_tagged_content:", paginated_tagged_content);

        if (paginated_tagged_content === null || paginated_tagged_content.content.length < 1) {
            console.error("In media_viewer_tag_mode.changeFilteringTags: The paginated tagged content was null for: ");
            console.log(new_tag_group);
            return;
        }

        content_cacher = tag_mode_content_cache.createTagContentCacher(new_tag_group, paginated_tagged_content);

        tagMode_activateContentCacher(content_cacher);
    }

    /**
     * Stores in active cache the given content page. Returns whether the content is now available on the active cache(filtering_medias_content_cacher).
     * If the page is already cached, it doesn't retrieve it again, just returns true immediately.
     * @param {number} page
     * @returns {Promise<boolean>}
     */
    const cacheTaggedContentPage = async page => {
        if (filtering_medias_content_cacher === null) return false;

        if (page > filtering_medias_content_cacher.TotalPages) return false;

        if (filtering_medias_content_cacher.hasPage(page)) return true;

        const paginated_tagged_content = await getTaggedContentPage(filtering_medias_content_cacher.ContentTags, page);

        if (paginated_tagged_content === null || paginated_tagged_content.content.length < 1) {
            console.error("In media_viewer_tag_mode.loadNextTaggedContentPage: The paginated tagged content was null for: ");
            console.log(filtering_medias_content_cacher.ContentTags);
            return false;
        }

        filtering_medias_content_cacher.catchPaginatedResponse(paginated_tagged_content);

        return true;
    }

    /**
     * resets the state of the media viewer tag mode but doesn't destroy the loaded cache.
     */
    export const tagMode_disableTagMode = () => {
        mv_filtering_tags.set([]);
        mv_tag_mode_enabled.set(false);
        mv_tagged_content.set([]);
        mv_tag_mode_total_content.set(NaN);
        mv_current_tagged_content_loaded = 0;
    }

    /**
     * Gets the requested page of tagged content with the given list of tags.
     * @param {import('@models/DungeonTags').DungeonTag[]} dungeon_tags
     * @param {number} page_num
     * @returns {Promise<import('@libs/DungeonsCommunication/dungeon_communication').PaginatedResponse<import('@models/Medias').MediaIdentity> | null>}
     */
    const getTaggedContentPage = async (dungeon_tags, page_num) => {
        const tag_ids = dungeon_tags.map(tag => tag.Id);

        const page_content = await getTaggedMedias(tag_ids, page_num, TAGGED_CONTENT_PAGE_SIZE);

        return page_content;
    } 

    /**
     * Whether the Media viewer tag mode is enabled.
     * @returns {boolean}
     */
    export const tagMode_isEnabled = () => {
        const tagged_content = get(mv_tagged_content);

        return tagged_content.length > 0;
    }

    /**
     * Whether the given content index is loaded.
     * @param {number} content_index
     * @param {import('@models/Medias').Media[]} [tagged_content] - To avoid having to call 'get' you can pass the tagged content array directly.
     * @returns {boolean}
     */
    const isContentIndexLoaded = (content_index, tagged_content) => {
        if (tagged_content == null) {
            tagged_content = get(mv_tagged_content);
        }

        return tagged_content[content_index] !== undefined && !isNullishMedia(tagged_content[content_index]);
    }

    /**
     * Loads the next tagged content page. Returns a true if successful
     * @returns {Promise<boolean>}
     */
    const loadNextTaggedContentPage = async () => {
        if (filtering_medias_content_cacher == null) {
            return false;
        }

        let next_page = filtering_medias_content_cacher.LastCachedPage;

        do {
            next_page++;
            if (next_page > 1000000) {
                throw new Error("In media_viewer_tag_mode.loadNextTaggedContentPage: Infinite loop detected"); // Theoretically this could be a valid next page. once we are sure this code is robust, remove this guard.
            }
        } while (next_page < filtering_medias_content_cacher.TotalPages && filtering_medias_content_cacher.hasPage(next_page))

        if (next_page > filtering_medias_content_cacher.TotalPages) {
            return false;
        }

        const cached_successfully = await cacheTaggedContentPage(next_page);

        if (!cached_successfully) {
            console.error("In media_viewer_tag_mode.loadNextTaggedContentPage: The next tagged content page couldn't be cached");
            return false;
        }

        const new_tagged_content = filtering_medias_content_cacher.getAllContent();

        mv_tagged_content.set(new_tagged_content);
        mv_current_tagged_content_loaded = new_tagged_content.length;

        return true;
    }

    /**
     * Loads a random tagged content page. Returns a true if successful. By random we are not talking about an unknown page, but any page event if it's not sequential. if there are gaps then the missing pages are
     * filled with NULLISH_MEDIA.
     * @param {number} page
     * @returns {Promise<boolean>}
     */
    const loadRandomTaggedContentPage = async page => {
        if (filtering_medias_content_cacher === null) {
            console.error("In media_viewer_tag_mode.loadRandomTaggedContentPage: The filtering medias content cacher is null");
            return false;
        }

        if (page > filtering_medias_content_cacher.TotalPages) {
            console.error("In media_viewer_tag_mode.loadRandomTaggedContentPage: The page is out of bounds");
            return false;
        }

        const page_state = filtering_medias_content_cacher.getPageState(page);

        if (page_state == null) {
            console.error("In media_viewer_tag_mode.loadRandomTaggedContentPage: The page state was null. which happens when the page num is out of bounds");
            return false;
        }

        const current_tagged_content = get(mv_tagged_content);
        let new_mv_tagged_content;

        if (filtering_medias_content_cacher.hasPage(page)) {

            if (isContentIndexLoaded(page_state.FirstIndex) && isContentIndexLoaded(page_state.LastIndex)) {
                return true
            };
        } else { // the pages is not on the cache.
            const cached_successfully = await cacheTaggedContentPage(page);

            if (!cached_successfully) {
                console.error("In media_viewer_tag_mode.loadRandomTaggedContentPage: The page couldn't be cached");
                return false;
            }
        }

        updateLoadedContent(); 

        // @ts-ignore
        globalThis.mv_tagged_content = new_mv_tagged_content; // TODO: Remove this line

        return true;
    }

    /**
     * Use when you cache a new page. just sets the new content on the mv_tagged_content array.
     * @returns {void}
     */
    const updateLoadedContent = () => {
        if (filtering_medias_content_cacher === null) {
            console.error("In media_viewer_tag_mode.updateLoadedContent: The filtering medias content cacher is null");
            return;
        }

        const new_mv_tagged_content = filtering_medias_content_cacher.getSequentialContent();

        mv_tagged_content.set(new_mv_tagged_content);
        mv_current_tagged_content_loaded = new_mv_tagged_content.length;
    }

    /**
     * Resets the state of the tagged content mode.
     * @returns {void}
     */
    export const tagMode_resetTaggedContentMode = () => {
        tagMode_disableTagMode();
        filtering_medias_content_cacher = null;
        tag_mode_content_cache.resetCache();
        filtering_tags_active_index_cache.clear();
    }

    /**
     * Sets the active tag content media index and takes care of any necessary extra content loading. Returns a promise that resolves to true if the media index now resolves to an actual media.
     * @param {number} new_media_index
     * @returns {Promise<boolean>}
     */
    export const tagMode_setActiveMediaIndex = async new_media_index => {
        // TODO: Address cases where the next Tag content needs to be loaded non linearly
        if (filtering_medias_content_cacher === null) {
            console.error("In media_viewer_tag_mode.setActiveMediaIndex: The filtering medias content cacher is null");
            return false;
        }

        if (new_media_index >= filtering_medias_content_cacher.TotalMedias || mv_current_tagged_content_loaded === 0) {
            console.error("In media_viewer_tag_mode.setActiveMediaIndex: The new media index is out of bounds");
            return false;
        }

        console.log(`new_media_index<${new_media_index}> > mv_current_tagged_content_loaded<${mv_current_tagged_content_loaded}>: ${new_media_index > mv_current_tagged_content_loaded}`);

        const current_tagged_content = filtering_medias_content_cacher.getSequentialContent();

        if (!isContentIndexLoaded(new_media_index, current_tagged_content)) {
            const page_with_content_index = filtering_medias_content_cacher.pageForContentIndex(new_media_index);
            
            const loaded_successfully = await loadRandomTaggedContentPage(page_with_content_index);

            if (!loaded_successfully) {
                console.error("In media_viewer_tag_mode.setActiveMediaIndex: The page couldn't be loaded");
                return false;
            }
        }

        active_tag_content_media_index.set(new_media_index);

        return true;
    }

    
/*=====  End of Methods  ======*/

