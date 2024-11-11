import { getMediaUrl, PatchRenameMediasRequest } from "@libs/DungeonsCommunication/services_requests/media_requests";

const DEFAULT_IMAGE_WIDTH = 307;

/**
* @typedef {Object} MediaParams
 * @property {string} uuid the 40 character identifier of the media resource
 * @property {string} name the name of the media resource
 * @property {string} last_seen the date the media resource was last seen
 * @property {string} main_category the main category of the media resource
 * @property {'IMAGE' | 'VIDEO'} type the type of the media resource, either IMAGE or VIDEO
 * @property {number} downloaded_from the id of the download batch that this media resource was downloaded from
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
        /** @type {'IMAGE' | 'VIDEO'} the type of the media resource, either IMAGE or V */
        this.type = type;
        /** @type {string} the main category of the media resource */
        this.main_category = main_category
        /** @type {number} the id of the download batch that this media resource was downloaded from */
        this.downloaded_from = downloaded_from;
        /** @type {string} the path of the category that the media resource belongs to */
        this.#category_path = category_path;

        this.#media_name = "";
        this.#file_extension = "";


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
        let target_width = DEFAULT_IMAGE_WIDTH;
        if (globalThis.innerWidth != null) {
            // window context
            target_width = globalThis.innerWidth * viewport_percentage; 
        }

        target_width = Math.round(target_width);
        
        return getMediaUrl(this.#category_path, this.name, true, false, target_width);
    }

    /**
     * Whether the media source is a video.
     * @returns {boolean}
     */
    isVideo = () => {
        return this.type === media_types.VIDEO
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