import { MediaChangesEmitter } from '@models/WorkManagers';
import { writable } from 'svelte/store';

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