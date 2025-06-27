import { MediaChangesEmitter } from '@models/WorkManagers';
import { writable } from 'svelte/store';

/**
 * Grid items intersection observer related event names.
 */
export const meg_intersection_observer_event_names = {
    VIEWPORT_ENTER: 'me-gallery-viewport-enter',
    VIEWPORT_LEAVE: 'me-gallery-viewport-leave',
}

/**
 * Media change state related event names.
 */
export const meg_media_change_state_event_names = {
    ALTERED_MEDIA_CHANGE_STATE: 'me-gallery-altered-media-change-state',
}

/**
 * Generates an ID for a MEGallery item that is a confirming CSS identifier.
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/id#syntax}
 * @param {import('@models/Medias').OrderedMedia} media
 * @returns {string}
 */
export const generateMEGalleryItemCSSIdentifier = (media) => {
    const generated_id = `meg-item-${media.Media.main_category}-${media.Media.uuid}`;

    return CSS.escape(generated_id);
}

/**
 * A Media changes manager used to edit the media gallery.
 * @type {import('svelte/store').Writable<MediaChangesEmitter | null>}
 * @default null
 */
export const me_gallery_changes_manager = writable(null);

/**
 * A list with all yanked medias.
 * @type {import('svelte/store').Writable<import('@models/Medias').OrderedMedia[]>}
 */
export const me_gallery_yanked_medias = writable([]);

/**
 * Whether the current focused media should be renamed.
 * @type {import('svelte/store').Writable<boolean>} 
 */
export const me_renaming_focused_media = writable(false);

/**
 * Sets the media explorer gallery to a clean state.
 * @returns {void}
 */
export const meGalleryReset = () => {
    me_gallery_changes_manager.set(new MediaChangesEmitter());
    me_gallery_yanked_medias.set([]);
}