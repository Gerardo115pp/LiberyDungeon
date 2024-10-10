import { writable } from 'svelte/store';

/**
 * A Media changes manager used to edit the media gallery.
 * @type {import('svelte/store').Writable<import('@models/WorkManagers').MediaChangesEmitter | null>}
 * @default null
 */
export const me_gallery_changes_manager = writable(null);

/**
 * A list with all yanked medias.
 * @type {import('svelte/store').Writable<import('@models/Medias').Media[]>}
 */
export const me_gallery_yanked_medias = writable([]);
