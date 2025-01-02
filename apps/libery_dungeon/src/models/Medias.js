import { 
    GetMediaByUUIDRequest,
    getMediaUrl,
    PatchRenameMediasRequest,
    PatchMediaRenameRequest,
    GetMediaIdentityByUUIDRequest,
    getSharedMediaLink
} from "@libs/DungeonsCommunication/services_requests/media_requests";
import { GetContentTaggedRequest } from "@libs/DungeonsCommunication/services_requests/metadata_requests/dungeon_tags_requests";
import { GetSharedMediaTokenRequest } from "@libs/DungeonsCommunication/services_requests/categories_requests";

const DEFAULT_IMAGE_WIDTH = 307;

/**
* @typedef {Object} MediaParams
 * @property {string} uuid the 40 character identifier of the media resource
 * @property {string} name the name of the media resource
 * @property {string} last_seen the date the media resource was last seen
 * @property {string} main_category the main category of the media resource
 * @property {string} type the type of the media resource, either IMAGE or VIDEO
 * @property {number} downloaded_from the id of the download batch that this media resource was downloaded from
*/

/**
 * A media params extension with all the necessary properties from category and category cluster to retrieve a media file.
* @typedef {{media: MediaParams} & import("@models/Categories").CategoryWeakIdentityParams} MediaIdentityParams
*/

export class Media {

    /**
     * The name without extension of the media resource
     * @type {string}
     */
    #media_name;

    /**
     * The file extension of the media resource
     * @type {string}
     */
    #file_extension;

    /**
     * The category path of the media resource
     * @type {string}
     */
    #category_path;

    /**
     * A token created by the pandasworld server. Can be used to share the media outside of the category cluster or even the platform.
     * Has to be created by a user with proper grants.
     * @type {string}
     */
    #shared_media_token;

    /**
     * An Dungeon media resource.
     * @param {MediaParams} param0
     */
    constructor({uuid, name, last_seen, main_category, type, downloaded_from}, category_path = "") {
        /** @type {string} the 40 character identifier of the media resource */
        this.uuid = uuid;
        /** @type {string} the name of the media resource */
        this.name = name;
        /** @type {string} the date the media resource was last seen */
        this.last_seen = last_seen;
        /** @type {string} the type of the media resource, either IMAGE or V */
        this.type = type;
        /** @type {string} the main category of the media resource */
        this.main_category = main_category
        /** @type {number} the id of the download batch that this media resource was downloaded from */
        this.downloaded_from = downloaded_from;
        /** @type {string} the path of the category that the media resource belongs to */
        this.#category_path = category_path;

        this.#media_name = "";
        this.#file_extension = "";

        this.#shared_media_token = "";

        this.#setMediaName(name);
    }

    /**
     * Returns the file extension of the media resource.
     * @returns {string}
     */
    get FileExtension() {
        return this.#file_extension;
    }

    /**
     * Returns the fullpath of the media resource.
     * @returns {string}
     */
    get Fullpath() {
        let fullpath = this.#category_path;

        if (!fullpath.endsWith("/")) {
            fullpath += "/";
        }

        fullpath += this.name;

        return fullpath;
    }

    /**
     * Returns a width-resized url of the media resource. If not running on a browser context, the DEFAULT_IMAGE_WIDTH will be used. 
     * if the caller knows what the viewport width is, set globalThis.innerWidth to the viewport width.
     * @param {number} viewport_percentage a number between 0 and 1 representing the percentage of the viewport width the media will occupy.
     * @returns {string}
     */
    getResizedUrl(viewport_percentage) {
        if (viewport_percentage === 0) return this.Url;

        let target_width = DEFAULT_IMAGE_WIDTH;
        if (globalThis.innerWidth != null) {
            // window context
            target_width = globalThis.innerWidth * viewport_percentage; 
        }

        target_width = Math.round(target_width);
        
        return getMediaUrl(this.#category_path, this.name, true, false, target_width);
    }

    /**
     * Returns a shared link. The first time this function is called it will need to create a shared token which involves communication with the server.
     * @return {Promise<string | null>}
     */
    getSharedUrl = async () => {
        if (this.#shared_media_token === "") {
            const new_shared_media_token = await this.#requestSharedMediaToken();

            if (new_shared_media_token === null) {
                return null;
            }

            this.#shared_media_token = new_shared_media_token;
        }

        return getSharedMediaLink(this.#shared_media_token);
    }

    /**
     * Whether the media source is a video.
     * @returns {boolean}
     */
    isVideo = () => {
        return this.type === media_types.VIDEO
    }

    /**
     * Whether the media is an image.
     * @returns {boolean}
     */
    isImage = () => {
        return this.type === media_types.IMAGE;
    }

    /**
     * Whether the media sources is an animation. returns true if the media source is a video or gif.
     * TODO: find a way to detect animated webp.
     * @returns {boolean}
     */
    isAnimated = () => {
        let ext = this.FileExtension.toLowerCase();

        return this.isVideo() || ext === "gif";
    }

    /**
     * Whether the media source is an animated image.
     * @returns {boolean}
     */
    isAnimatedImage = () => {
        let ext = this.FileExtension.toLowerCase();

        return ext === "gif";
    }

    /**
     * Returns the name of the media without file extension.
     * @returns {string}
     */
    get MediaName() {
        return this.#media_name;
    } 

    /**
     * Returns the url of the media resource.
     * @returns {string}
     */
    get MobileUrl() {
        return getMediaUrl(this.#category_path, this.name, false, true);
    }

    /**
     * Requests a shared media token from the server.
     * @returns {Promise<string | null>}
     */
    #requestSharedMediaToken = async () => {
        /**
         * @type {string | null}
         */
        let shared_media_token = null;
         
        const request = new GetSharedMediaTokenRequest(this.uuid);

        let response = await request.do();

        if (response.Ok && response.data.response != "") {
            shared_media_token = response.data.response;
        }

        return shared_media_token;
    }

    /**
     * Renames the media. This changes the media record in the server, not just this instance.
     * @param {string} new_name
     * @returns {Promise<boolean>}
     */
    rename = async (new_name) => {
        if (new_name === this.name) return true;

        if (new_name === "" || !new_name) {
            throw new Error(`In @models/Medias.Media.rename: tried to rename with an invalid value <${new_name}>`);
        }

        let successful = await renameMedia(this, new_name);

        if (successful) {
            this.name = new_name;
        }

        // Reprocessing name fragments. 
        this.#media_name = "";
        this.#file_extension = "";

        this.#setMediaName(new_name);

        return successful
    }

    /**
     * Sets the media name and file extension.
     * @param {string} name full name of the media resource
     * @returns {void}
     */
    #setMediaName(name) {
        const name_fragments = name.split(".");
        const fragment_count = name_fragments.length;

        if (fragment_count < 2) {
            this.#media_name = name;
            this.#file_extension = "unknown";
            return;
        }
        
        const media_file_extension = name_fragments.pop();

        if (media_file_extension != null) {
            this.#file_extension = media_file_extension;
        }

        this.#media_name = name_fragments.join(".");
    }

    /**
     * Returns the url of the media resource.
     * @returns {string}
     */
    get Url() {
        return getMediaUrl(this.#category_path, this.name, false, false);
    }

    /**
     * Converts the media resource to MediaParams representation.
     * @returns {MediaParams}
     */
    toParams() {
        return {
            uuid: this.uuid,
            name: this.name,
            last_seen: this.last_seen,
            main_category: this.main_category,
            type: this.type,
            downloaded_from: this.downloaded_from
        }
    }
}

export class MediaIdentity {
    /**
     * The media object
     * @type {Media}
     */
    #the_media;

    /**
     * The category path
     * @type {string}
     */
    #category_path;
     
    /**
     * @type {string} 
     */
    #category_uuid;

    /**
     * @type {string}
     */
    #cluster_path;

    /**
     * @type {string}
     */
    #cluster_uuid;

    /**
     * @param {MediaIdentityParams} param0
     */
    constructor({media, category_path, category_uuid, cluster_path, cluster_uuid}) {
        this.#the_media = new Media(media, category_path);

        this.#category_path = category_path;
        this.#category_uuid = category_uuid;

        this.#cluster_path = cluster_path;
        this.#cluster_uuid = cluster_uuid;
    }

    /**
     * the uuid of the media.
     * @type {string}
     */
    get uuid() {
        return this.#the_media.uuid;
    }

    /**
     * the name of the media.
     * @type {string}
     */
    get name() {
        return this.#the_media.name;
    }

    /**
     * the main_category of the media.
     * @type {string}
     */
    get main_category() {
        return this.#the_media.main_category;
    }

    /**
     * the type of the media.
     * @type {string}
     */
    get type() {
        return this.#the_media.type;
    }

    /**
     * The media Object of the media identity.
     * @type {Media}
     */
    get Media() {
        return this.#the_media;
    }
}

export class OrderedMedia {
    /**
     * The media object
     * @type {Media}
     */
    #the_media;

    /**
     * The provided order of the media
     * @type {number}
     */
    #the_order;

    /**
     * 
     * @param {Media} media 
     * @param {number} order 
     */
    constructor(media, order) {
        this.#the_media = media;
        this.#the_order = order;
    }

    /**
     * Returns the media object.
     * @returns {Media}
     */
    get Media() {
        return this.#the_media;
    }

    /**
     * Returns the order of the media.
     * @returns {number}
     */
    get Order() {
        return this.#the_order;
    }

    /**
     * the uuid of the media object
     * @type {string}
     */
    get uuid() {
        return this.#the_media.uuid;
    }

    /**
     * the name of the media object
     * @type {string}
     */
    get name() {
        return this.#the_media.name;
    }

    /**
     * The media name, like `name` but without the file extension.
     * @type {string}
     */
    get MediaName() {
        return this.#the_media.MediaName;
    }

    /**
     * the last_seen of the media object
     * @type {string}
     */
    get last_seen() {
        return this.#the_media.last_seen;
    }

    /**
     * the main_category of the media object
     * @type {string}
     */
    get main_category() {
        return this.#the_media.main_category;
    }

    /**
     * the type of the media object
     * @type {string}
     */
    get type() {
        return this.#the_media.type;
    }

    /**
     * the fullpath of the media resource
     * @type {string}
     */
    get Fullpath() {
        return this.#the_media.Fullpath;
    }

    /**
     * Returns the file extension of the media resource.
     * @returns {string}
     */
    get FileExtension() {
        return this.#the_media.FileExtension;
    }
}

export const media_types = {
    IMAGE: "IMAGE",
    VIDEO: "VIDEO"
}



/* ----------------------------- Nullish models ----------------------------- */

    export const NULLISH_MEDIA = new Media({
        uuid: "",
        name: "",
        last_seen: "",
        main_category: "",
        type: "",
        downloaded_from: NaN
    });

    /**
     * Whether a media object is nullish.
     * @param {Media} media
     */
    export const isNullishMedia = (media) => {
        return Object.is(media, NULLISH_MEDIA);
    }

    export const NULLISH_MEDIA_IDENTITY = new MediaIdentity({
        media: NULLISH_MEDIA,
        category_path: "",
        category_uuid: "",
        cluster_path: "",
        cluster_uuid: ""
    });

    /**
     * Whether a media identity object is nullish.
     * @param {MediaIdentity} media_identity
     */
    export const isNullishMediaIdentity = (media_identity) => {
        return Object.is(media_identity, NULLISH_MEDIA_IDENTITY);
    }

/*=============================================
=            Model actions            =
=============================================*/

/**
 * Returns a media object by its uuid.
 * @param {string} media_uuid
 * @returns {Promise<Media | null>}
 */
export const getMediaByUUID = async (media_uuid) => {
    /**
     * @type {Media | null}
     */
    let media = null;

    const request = new GetMediaByUUIDRequest(media_uuid);

    const response = await request.do();

    if (response.Ok && response.data != null) {
        media = new Media(response.data);
    }

    return media;
}

/**
 * Returns a media identity object by the media uuid. Requires a cluster sign access to be present as a http-only cookie. otherwise the request will fail with a Forbidden status code.
 * @param {string} media_uuid
 * @returns {Promise<MediaIdentity | null>}
 */
export const getMediaIdentityByUUID = async (media_uuid) => {
    /**
     * @type {MediaIdentity | null}
     */
    let media_identity = null;

    const request = new GetMediaIdentityByUUIDRequest(media_uuid);

    const response = await request.do();

    if (response.Ok && response.data != null) {
        media_identity = new MediaIdentity(response.data);
    }

    return media_identity;
}

/**
 * Returns a list of media identites that matcha given list of tags by their ids. Returns a paginated response. Optionally a page and page_size parameters can be passed.
 * @param {number[]} tag_ids
 * @param {number} [page]
 * @param {number} [page_size]
 * @returns {Promise<import('@libs/DungeonsCommunication/dungeon_communication').PaginatedResponse<MediaIdentity> | null>}
 */
export const getTaggedMedias = async (tag_ids, page, page_size) => {
    const request = new GetContentTaggedRequest(tag_ids, page, page_size);

    /**
     * @type {import('@libs/DungeonsCommunication/dungeon_communication').PaginatedResponse<MediaIdentity> | null}
     */    
    let paginated_response = null;

    try {
        const response = await request.do();

        if (response.Ok && response.data != null) {
            const media_identities = response.data.content.map(media_identity_params => new MediaIdentity(media_identity_params))

            paginated_response = {
                ...response.data,
                content: media_identities
            }
        } 
    } catch (error) {
        console.error(`In @models/Medias/getTaggedMedias: while fetching tagged medias got ${error}`);
    }

    return paginated_response;
}

/**
 * Renames a single media. Returns true if the rename was successful.
 * @param {Media} media
 * @param {string} new_name
 * @returns {Promise<boolean>}
 */
export const renameMedia = async (media, new_name) => {
    if (new_name === "") {
        return false;
    }

    if (new_name === media.name) {
        return true;
    }

    const request = new PatchMediaRenameRequest(new_name, media.uuid);

    let rename_successful = false;

    const response = await request.do();

    if (response.Ok && response.data != null) {
        rename_successful = response.data.response;
    }

    return rename_successful;
}

/**
 * Renames a sequence of medias. gets a Sequence map, which is just an object of the form {media_uuid: new_name} where new_name does not include 
 * a file extension. and a category_uuid were the medias are located.
 * @param {Object<string, string>} sequence_map 
 * @param {string} category_uuid 
 */
export const sequenceRenameMedias = async (sequence_map, category_uuid) => {
    const request = new PatchRenameMediasRequest(category_uuid, sequence_map);
    let renamed = false;

    const response = await request.do();

    if (response.status === 204) {
        renamed = true;
    }

    return renamed;
}

/*=====  End of Model actions  ======*/

