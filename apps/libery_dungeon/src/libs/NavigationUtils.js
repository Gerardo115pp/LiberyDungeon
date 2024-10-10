import { goto } from "$app/navigation";

/**
 * Navigates the Media viewer on a specific media.
 * @param {string} category_uuid
 * @param {number} media_index
 * @param {boolean} use_spa - if true, will attempt to use svelte_spa_router.
 */
export const navigateToMediaViewer = (category_uuid, media_index=-1, use_spa=true) => {
    let media_viewer_url = `/media-viewer/${category_uuid}`;

    if (media_index >= 0) {
        media_viewer_url += `/${media_index}`;
    }


    if (use_spa && globalThis.innerWidth !== undefined) {
        goto(media_viewer_url);
        return;
    }

    
    window.location.replace(media_viewer_url);
    window.location.reload();
}