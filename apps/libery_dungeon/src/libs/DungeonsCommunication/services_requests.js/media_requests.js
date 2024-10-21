import { medias_server, collect_server } from "../services";

const media_fs_endpoint = `${medias_server}/medias-fs`;
const thumbnails_fs_endpoint = `${medias_server}/thumbnails-fs`;

/**
 * @param {string} category_path
 * @param {string} media_name
 * @returns {string}
 */
export const getMediaUrl = (category_path, media_name, use_thumbnail=false, use_mobile=false, target_width=0) => {
    category_path = category_path.replace(/(^\/|\/$)/g, "");

    let endpoint = use_thumbnail ? thumbnails_fs_endpoint : media_fs_endpoint;

    if (use_mobile && !use_thumbnail) {
        endpoint = `${endpoint}/mobile`;
    }

    endpoint = `${endpoint}/${category_path}/${media_name}`;

    endpoint =  endpoint.replace(/#/g, "%23");

    if (target_width > 0 && use_thumbnail) {
        endpoint += `?width=${target_width}`;
    }   

    return endpoint;    
}

/**
 * Returns the url for a media that is on the trashcan.
 * @param {string} media_name
 */
export const getTrashcanMediaUrl = (media_name, width=120) => {
    let endpoint = `${thumbnails_fs_endpoint}/libery-trashcan/`;

    let url_encoded_media_name = encodeURIComponent(media_name);

    endpoint += url_encoded_media_name;

    if (width > 0) {
        endpoint += `?width=${width}`;
    }

    return endpoint;
}

/**
* @param {string} media_url
 * @returns {string}
*/
export const getProxyMediaUrl = (media_url) => `${collect_server}/4chan-boards/boards/proxy-media?media_url=${media_url}`;