import { get, writable } from 'svelte/store';
import { getTaggedMedias } from '@models/Medias';
import { TagContentCache } from '@models/WorkManagers';

const TAGGED_CONTENT_PAGE_SIZE = 10;


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

    /**
     * Sets a TaggerContentCacher as active setting all necessary variables and stores to the appropriate values from the given cacher.
     * @param {import('@models/WorkManagers').TaggedContentCacher} cacher
     * @returns {void}
     */
    const tagMode_activateContentCacher = cacher => {
        const cacher_filtering_tags = cacher.ContentTags;

        if (cacher_filtering_tags.length === 0) return;

        const new_tagged_content = cacher.getAllContent();

        console.log("new_tagged_content:", new_tagged_content);

        mv_filtering_tags.set(cacher_filtering_tags);
        mv_tagged_content.set(new_tagged_content);
        filtering_medias_content_cacher = cacher;
        mv_tag_mode_total_content.set(cacher.TotalMedias);
        mv_current_tagged_content_loaded = new_tagged_content.length;
        mv_tag_mode_enabled.set(true);
        active_tag_content_media_index.set(0);
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

        if (paginated_tagged_content === null || paginated_tagged_content.content.length < 1) {
            console.error("In media_viewer_tag_mode.changeFilteringTags: The paginated tagged content was null for: ");
            console.log(new_tag_group);
            return;
        }

        content_cacher = tag_mode_content_cache.createTagContentCacher(new_tag_group, paginated_tagged_content);

        tagMode_activateContentCacher(content_cacher);
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

        const paginated_tagged_content = await getTaggedContentPage(filtering_medias_content_cacher.ContentTags, next_page);

        if (paginated_tagged_content === null || paginated_tagged_content.content.length < 1) {
            console.error("In media_viewer_tag_mode.loadNextTaggedContentPage: The paginated tagged content was null for: ");
            console.log(filtering_medias_content_cacher.ContentTags);
            return false;
        }

        filtering_medias_content_cacher.catchPaginatedResponse(paginated_tagged_content);

        const new_tagged_content = filtering_medias_content_cacher.getAllContent();

        mv_tagged_content.set(new_tagged_content);
        mv_current_tagged_content_loaded = new_tagged_content.length;

        return true;
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

        while (new_media_index >= mv_current_tagged_content_loaded) {
            console.log("Loading next tagged content page");
            try {
                await loadNextTaggedContentPage();
            } catch (error) {
                console.error("In media_viewer_tag_mode.setActiveMediaIndex: Error while loading next tagged content page");
                console.error(error);
                return false;
            }
        }

        active_tag_content_media_index.set(new_media_index);

        return true;
    }


/*=====  End of Methods  ======*/

