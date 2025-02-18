import { medias_server, collect_server, categories_server } from "../services";
import { HttpResponse, attributesToJson } from "../base";

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
 * The server returns a random image from a category. optionally if cache_seconds is passed as a positive integer, the server instructs the browser to cache the retrieved resource of the requested amount of time.
 * @param {string} category_id
 * @param {string} cluster_id
 * @param {number} [cache_seconds]
 * @returns {string}
 */
export const getRandomMediaUrl = (category_id, cluster_id, cache_seconds) => {
    if (globalThis.self == null) {
        console.error("Cannot use getRandomMediaUrl outside of a windowed context");
        return '';
    };

    if (!category_id || !cluster_id) {
        throw new Error("In getRandomMediaUrl: category_id and cluster_id are required.");
    }

    if (!URL.canParse(`${medias_server}/random-medias-fs`, location.origin)) {
        throw new Error(`url<${medias_server}/random-medias-fs> has a syntax problem`);
    }

    const random_media_url = new URL(`${medias_server}/random-medias-fs`, location.origin);

    random_media_url.searchParams.set("category_id", category_id);
    random_media_url.searchParams.set("cluster_id", cluster_id);

    if (cache_seconds) {
        random_media_url.searchParams.set("cache_seconds", cache_seconds.toString());
    }

    return random_media_url.toString();
}

/**
* @param {string} media_url
 * @returns {string}
*/
export const getProxyMediaUrl = (media_url) => `${collect_server}/4chan-boards/boards/proxy-media?media_url=${media_url}`;

/**
 * Returns a link to consume a share media.
 * @param {string} shared_media_token
 * @returns {string}
 */
export const getSharedMediaLink = (shared_media_token) => {
    return `${location.origin}${medias_server}/shared-content/shared-media?share_token=${shared_media_token}`;
}

/**
 * Renames a sequence of medias. gets a Sequence map, which is just an object of the form {media_uuid: new_name} where new_name does not include
 * a file extension. and a category_uuid were the medias are located.
 */
export class PatchRenameMediasRequest {
    static endpoint = `${medias_server}/medias-fs/sequence-rename`;

    /**
     * @param {string} category_uuid
     * @param {Object.<string, string>} sequence_members
     */
    constructor(category_uuid, sequence_members) {
        this.category_uuid = category_uuid;
        this.sequence_members = sequence_members;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        let response;
        let renamed = false;

        try {
            response = await fetch(PatchRenameMediasRequest.endpoint, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            renamed = response.ok;
        } catch (error) {
            console.error("Error while renaming medias: ", error);
            throw error;
        }

        return new HttpResponse(response, renamed);
    }
}

/**
 * Returns a media object by its uuid.
 */
export class GetMediaByUUIDRequest {

    static endpoint = `${medias_server}/medias/by-uuid`;

    /**
     * @param {string} media_uuid
     */
    constructor(media_uuid) {
        this.media_uuid = media_uuid;

        if (!URL.canParse(`${GetMediaByUUIDRequest.endpoint}?uuid=${media_uuid}`, location.origin)) {
            throw new Error(`url<${GetMediaByUUIDRequest.endpoint}> has a syntax problem`);
        }
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('@models/Medias').MediaParams | null>>}
     */
    do = async () => {
        
        let response;
        let media = null;

        const resource_url = new URL(GetMediaByUUIDRequest.endpoint, location.origin);

        resource_url.searchParams.set("uuid", this.media_uuid);

        try {
            response = await fetch(resource_url);
        } catch (error) {
            console.error("Error while fetching media by uuid: ", error);
            throw error;
        }

        if (response.ok) {
            media = await response.json();
        }

        return new HttpResponse(response, media);
    }
}

/**
 * Returns a media identity object by the media uuid. Requires a cluster sign access to be present as a http-only cookie. otherwise 403 is returned.
 */
export class GetMediaIdentityByUUIDRequest {

    static endpoint = `${medias_server}/medias/identity`;

    /**
     * @param {string} media_uuid
     */
    constructor(media_uuid) {
        this.media_uuid = media_uuid;

        if (!URL.canParse(`${GetMediaIdentityByUUIDRequest.endpoint}?uuid=${media_uuid}`, location.origin)) {
            throw new Error(`url<${GetMediaIdentityByUUIDRequest.endpoint}> has a syntax problem`);
        }
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('@models/Medias').MediaIdentityParams | null>>}
     */
    do = async () => {
        
        let response;
        let media_identity = null;

        const resource_url = new URL(GetMediaIdentityByUUIDRequest.endpoint, location.origin);

        resource_url.searchParams.set("uuid", this.media_uuid);

        try {
            response = await fetch(resource_url);
        } catch (error) {
            console.error("Error while fetching media identity by uuid: ", error);
            throw error;
        }

        if (response.ok) {
            media_identity = await response.json();
        }

        return new HttpResponse(response, media_identity);
    }
}

/**
 * Returns a list of media identities by their uuids. All medias must belong to the same category cluster(not necessarily the same category). Requires public_content_view grant.
 */
export class GetMediaIdentitiesByUUIDsRequest {

    static endpoint = `${categories_server}/medias/in-list`;

    /**
     * @param {string[]} media_uuids
     * @param {string} cluster_uuid
     */
    constructor(media_uuids, cluster_uuid) {
        this.media_uuids = media_uuids;
        this.cluster_uuid = cluster_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('@models/Medias').MediaIdentityParams[]>>}`
     */
    do = async () => {
        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {import('@models/Medias').MediaIdentityParams[]}
         */
        let media_identities = [];

        let media_uuids_param = this.media_uuids.join(",");
        
        const resource_url = new URL(GetMediaIdentitiesByUUIDsRequest.endpoint, location.origin);

        resource_url.searchParams.set("media_uuids", media_uuids_param);

        resource_url.searchParams.set("cluster_uuid", this.cluster_uuid);

        try {
            response = await fetch(resource_url);

            if (response.ok) {
                media_identities = await response.json();
            }
        } catch (error) {
            console.error("Error while fetching media identities by uuids: ", error);
            throw error;
        }

        return new HttpResponse(response, media_identities);
    }
}

/**
 * A request to rename a media resource by UUID. It takes the media uuid and the new name. if the name passed is an empty string the request class will not complain but the request will fail. if the same 
 * name as the media currently has is passed as the new name, then the server will return a 200 will not actually process the request. The user most have content altering permission to be able use the endpoint.
 */
export class PatchMediaRenameRequest {

    static endpoint = `${medias_server}/medias/rename`;

    /**
     * @param {string} new_name
     * @param {string} media_uuid
     */
    constructor(new_name, media_uuid) {
        this.new_name = new_name;
        this.media_uuid = media_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('../base').BooleanResponse>>}
     */
    do = async () => {
        
        /**
         * @type {import('../base').BooleanResponse}
         */
        const response_body = {
            response: false
        }

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(PatchMediaRenameRequest.endpoint, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            response_body.response = response.ok;
        } catch(error) {
            console.error("Error while renaming media: ", error);
            throw error;
        }

        return new HttpResponse(response, response_body);
    }
}