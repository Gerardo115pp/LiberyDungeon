import { goto } from "$app/navigation";

/**
 * Navigates the Media viewer on a specific media.
 * @param {string} category_uuid
 * @param {MediaViewerParams} [navigation_params] 
 * @typedef {Object} MediaViewerParams
 * @property {number} [media_index]
 * @property {boolean} [use_spa] - if true, will attempt to use svelte_spa_router.
 * @property {number} [time] - if the media is a video, 
 * @property {string} [media_uuid]
 */
export const navigateToMediaViewer = (category_uuid, navigation_params) => {

    let media_index = navigation_params?.media_index ?? -1;
    let use_spa = navigation_params?.use_spa !== undefined ? navigation_params.use_spa : true;
    let time = navigation_params?.time ?? NaN;
    let media_uuid = navigation_params?.media_uuid ?? "";

    console.log("media_index", media_index);
    console.log("use_spa", use_spa);
    console.log("time", time);
    console.log("media_uuid", media_uuid);

    let media_viewer_url = `/media-viewer/${category_uuid}`;

    if (media_index >= 0) {
        media_viewer_url += `/${media_index}`;
    }

    const url = new URL(media_viewer_url, globalThis.location.origin);

    if (media_uuid !== "") {
        url.searchParams.append("media_uuid", media_uuid)
    }

    if (!isNaN(time)) {
        url.searchParams.append("time", time.toString());
    }

    console.log(url.toString());

    if (use_spa && globalThis.innerWidth !== undefined) {
        goto(url);
        return;
    }

    
    window.location.replace(url);
    window.location.reload();
}