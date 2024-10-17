import { MediaChangesEmitter } from '@models/WorkManagers';
import { writable } from 'svelte/store';

/**
 * A Media changes manager used to edit the media gallery.
 * @type {import('svelte/store').Writable<MediaChangesEmitter>| null>}
 * @default null
 */
export const me_gallery_changes_manager = writable(null);

/**
 * A list with all yanked medias.
 * @type {import('svelte/store').Writable<import('@models/Medias').Media[]>}
 */
export const me_gallery_yanked_medias = writable([]);

/**
 * Sets the media explorer gallery to a clean state.
 * @returns {void}
 */
export const meGalleryReset = () => {
    me_gallery_changes_manager.set(new MediaChangesEmitter());
    me_gallery_yanked_medias.set([]);
}