import { get, writable } from 'svelte/store';
import { getTaggedMedias } from '@models/Medias';
import { arrayToParam } from '@libs/DungeonsCommunication/base';

const TAGGED_CONTENT_PAGE_SIZE = 100;

/**
 * The dungeon tags the medias displayed by the medias viewer must have.
 * @type {import('svelte/store').Writable<import('@models/DungeonTags').DungeonTag[]>}
 */
export const mv_filtering_tags = writable([]);

/**
 * The last page of tagged medias that has been loaded.
 * @type {import('svelte/store').Writable<number>}
 */
export const mv_last_tagged_content_page = writable(NaN);

/**
 * The last paginated response received from the server.
 * @type {import('@libs/DungeonsCommunication/dungeon_communication').PaginatedResponse<import('@models/Medias').MediaIdentity> | null}
 */
let last_paginated_response = null;

/**
 * The total amount of pages of tagged medias with the current page size.
 * @type {import('svelte/store').Writable<number>} 
 */
export const mv_tagged_content_total_pages = writable(NaN);

/**
 * The index corresponding to the active media on the tagged medias content list.
 * @type {import('svelte/store').Writable<number>}
 */
export const active_tag_content_media_index = writable(0);

/**
 * The medias tagged by the filtering tags.
 * @type {import('svelte/store').Writable<import('@models/Medias').MediaIdentity[]>}
 */
export const mv_tagged_content = writable([]);

/**
 * Whether the tag_mode is enabled. This store is used for streamlining reactivity but tagModeEnabled is the actual reliable source of truth.
 * @type {import('svelte/store').Writable<boolean>}
 */
export const mv_tag_mode_enabled = writable(false);

/**
 * Resets the state of the tagged content mode.
 * @returns {void}
 */
export const resetTaggedContentMode = () => {
    mv_filtering_tags.set([]);
    mv_last_tagged_content_page.set(NaN);
    mv_tagged_content_total_pages.set(NaN);
}

/**
 * Whether the Media viewer tag mode is enabled.
 * @returns {boolean}
 */
export const tagModeEnabled = () => {
    const tagged_content = get(mv_tagged_content);

    return tagged_content.length > 0;
}

/**
 * Changes the tagged content using a new group of tags.
 * @param {import('@models/DungeonTags').DungeonTag[]} new_tag_group 
 * @returns {Promise<void>}
 */
export const changeFilteringTags = async new_tag_group  => {
    mv_filtering_tags.set(new_tag_group);

    const tag_ids = new_tag_group.map(tag => tag.Id);

    const paginated_tagged_content = await getTaggedMedias(tag_ids, 1, TAGGED_CONTENT_PAGE_SIZE);

    if (paginated_tagged_content === null || paginated_tagged_content.content.length < 1) {
        console.error("In media_viewer_tag_mode.changeFilteringTags: The paginated tagged content was null for: ");
        console.log(new_tag_group);
        return;
    }

    console.log(`Paginated tagged content: `, paginated_tagged_content);

    last_paginated_response = paginated_tagged_content;

    mv_tagged_content.set(paginated_tagged_content.content);
    mv_tag_mode_enabled.set(true);
}