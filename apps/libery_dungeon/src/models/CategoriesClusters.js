import { 
    GetCategoriesClustersRequest, 
    PatchClusterRequest,
    GetNewClusterDirectoryOptionsRequest,
    GetClusterDirectoryValidationRequest,
    PostCreateClusterRequest,
    DeleteClusterRecordRequest,
    GetClusterRootPathRequest
 } from "@libs/DungeonsCommunication/services_requests/categories_cluster_requests";
import {
    GetMediaIdentitiesByUUIDsRequest
} from "@libs/DungeonsCommunication/services_requests/media_requests";
import { 
    MediaIdentity,
    getMediaIdentityByUUID
} from "./Medias";
import { 
    getVideoMoments,
    createVideoMoment,
    deleteVideoMoment,
    getAllClusterVideoMoments
} from "./Metadata";
import { 
    InnerCategory
} from "./Categories";
import { UUIDHistory } from "@models/WorkManagers";

export class CategoriesCluster {
    
    /*----------  Class properties  ----------*/
        /**
         * The uuid4 of the cluster. This should be client-side generated.
         * @type {string}
         */
        #uuid;

        /**
         * The name of the cluster.
         * @type {string}
         */
        #name;

        /**
         * The root of the cluster's file system.
         * @type {string}
         */
        #fs_path;

        /**
         * The uuid of the category where content is automatically downloaded for this cluster.
         * @type {string}
         */
        #filter_category;

        /**
         * The uuid of the category that acts as the root of the category tree for this cluster.
         * @type {string}
         */
        #root_category;

        /**
         * A map of media uuids -> video moments that exist on this cluster and have been fetched.
         * @type {Map<string, import('@models/Metadata').VideoMoment[]>}
         */
        #cluster_video_moments

        /**
         * A dictionary to retrieve cached.
         * @type {Map<string, import('@models/Medias').MediaIdentity>}
         */
        #media_identity_dictionary

        /**
         * The category usage history.
         * @type {import("@models/WorkManagers").UUIDHistory<import('@models/Categories').InnerCategory>}
         */
        #category_usage_history;

    /**
     * @param {CategoriesClusterParams} param0 
     * @typedef {Object} CategoriesClusterParams
     * @property {string} uuid
     * @property {string} name
     * @property {string} fs_path
     * @property {string} filter_category
     * @property {string} root_category
     */
    constructor({ uuid, name, fs_path, filter_category, root_category }) {
        if (uuid === undefined || name === undefined || fs_path === undefined || filter_category === undefined || root_category === undefined) {
            throw new Error('Missing required parameters.');
        }

        this.#uuid = uuid;
        this.#name = name;
        this.#fs_path = fs_path;
        this.#filter_category = filter_category;
        this.#root_category = root_category;

        this.#cluster_video_moments = new Map();
        this.#media_identity_dictionary = new Map();

        this.#category_usage_history = new UUIDHistory(30);
    }
    
    /*=============================================
    =            Category usage history            =
    =============================================*/
    
        /**
         * Adds or refreshes a category usage record.
         * @param {import('@models/Categories').InnerCategory} category
         */
        touchCategoryUsage = category => {
            this.#category_usage_history.Add(category);
        }

        /**
         * The category usage history of this cluster.
         * @returns {import("@models/WorkManagers").UUIDHistory<import('@models/Categories').InnerCategory>}
         */
        get CategoryUsageHistory() {
            return this.#category_usage_history;
        }
    
    /*=====  End of Category usage history  ======*/
    
    /*=============================================
    =            Media identities dictionary            =
    =============================================*/
    
        /**
         * Adds a media identity to the media identity dictionary.
         * @param {MediaIdentity} media_identity
         * @returns {void}
         */ 
        addMediaIdentityToDict = media_identity => {
            this.#media_identity_dictionary.set(media_identity.Media.uuid, media_identity);
        }

        /**
         * Returns a list of media identities that belong in this cluster.
         * @param {string[]} media_uuids
         * @returns {Promise<import("@models/Medias").MediaIdentity[]>}
         */
        async getClusterMedias(media_uuids) {
            const cached_identities = this.getStoredClusterMedias(media_uuids);

            if (cached_identities.length === media_uuids.length) {
                return cached_identities;
            }

            const missing_media_uuids = [];

            for (let uuid of media_uuids) {
                if (!this.#media_identity_dictionary.has(uuid)) {
                    missing_media_uuids.push(uuid);
                }
            }

            const request = new GetMediaIdentitiesByUUIDsRequest(missing_media_uuids, this.#uuid);

            /**
             *  @type {MediaIdentity[]} 
             */
            let new_media_identities = [];

            const response = await request.do();

            if (response.Ok && response.data != null) {
                new_media_identities = response.data.map(media_identity_params => new MediaIdentity(media_identity_params));

                new_media_identities.forEach(mi => {
                    this.addMediaIdentityToDict(mi);
                });
            }

            return new_media_identities;
        }

        /**
         * Returns the media identities on a media uuid array that exist on the dictionary 
         * @param {string[]} media_uuids 
         * @returns {MediaIdentity[]}
         */
        getStoredClusterMedias(media_uuids) {
            /**
             * @type {MediaIdentity[]}
             */
            const media_identities = [];

            for (let uuid of media_uuids) {
                const stored_identity = this.#media_identity_dictionary.get(uuid);

                if (stored_identity === undefined) {
                    continue;
                }

                media_identities.push(stored_identity);
            }

            return media_identities;
        }
    
    /*=====  End of Media identities dictionary  ======*/

    /*=============================================
    =            Video moments            =
    =============================================*/

        /**
         * Adds a video moment to a given cache.
         * @param {import('@models/Metadata').VideoMoment} video_moment
         * @param {Map<string, import('@models/Metadata').VideoMoment[]>} cache
         * @returns {void}
         */
        #addVideoMomentToCache = (video_moment, cache) => {
            if (!cache.has(video_moment.VideoUUID)) {
                cache.set(video_moment.VideoUUID, []);
            }

            const sibling_video_moments = /** @type {import('@models/Metadata').VideoMoment[]} */ (cache.get(video_moment.VideoUUID));

            sibling_video_moments.push(video_moment);
        }

        /**
         * Creates a new video moment for a media that exists on this cluster.
         * @param {string} media_uuid
         * @param {number} moment_time
         * @param {string} moment_title
         * @returns {Promise<import('@models/Metadata').VideoMoment | null>}
         */
        createVideoMoment = async (media_uuid, moment_time, moment_title) => {
            const new_video_moment = await createVideoMoment(
                media_uuid,
                this.#uuid,
                moment_time,
                moment_title
            );

            if (new_video_moment === null) return null;
            
            let existing_video_moments = this.#cluster_video_moments.get(media_uuid);

            if (existing_video_moments === undefined) {
                existing_video_moments = [new_video_moment];
            } else {
                existing_video_moments = [new_video_moment, ...existing_video_moments];
            }

            existing_video_moments = this.#sortVideoMoments(existing_video_moments);

            this.#cluster_video_moments.set(media_uuid, existing_video_moments);

            return new_video_moment;
        }

        /**
         * Clears the entire video moments cache of the cluster.
         * @returns {void}
         */
        #clearVideoMomentsCache = () => {
            this.#cluster_video_moments.clear();
        }

        /**
         * Deletes a given video moment.
         * @param {import('@models/Metadata').VideoMoment} video_moment
         * @returns {Promise<boolean>}
         */
        deleteVideoMoment = async video_moment => {
            const success = await deleteVideoMoment(video_moment.ID);

            if (!success) return false;

            this.#removeVideoMomentLocally(video_moment);

            return true;
        }

        /**
         * Returns all the video moments available for the given media uuid.
         * @param {string} media_uuid
         * @returns {Promise<import('@models/Metadata').VideoMoment[] | null>}
         */
        getMediaVideoMoments = async media_uuid => {
            const cached_video_moments = this.#cluster_video_moments.get(media_uuid);
            if (cached_video_moments != undefined) {
                return cached_video_moments;
            }

            const fetched_video_moments = await getVideoMoments(media_uuid, this.#uuid);

            if (fetched_video_moments.length === 0) {
                return null;
            }

            const sorted_video_moments = this.#sortVideoMoments(fetched_video_moments);

            this.#cluster_video_moments.set(media_uuid, sorted_video_moments);

            return sorted_video_moments;
        }

        /**
         * Returns all the items in the video content cache.
         * @returns {import('@models/Metadata').VideoMoment[]}
         */
        getVideoMomentsCacheItems = () => {
            /**
             * @type {import('@models/Metadata').VideoMoment[]}
             */
            let video_moments = [];

            for (let moments of this.#cluster_video_moments.values()) {
                video_moments = [...video_moments, ...moments];
            }

            return video_moments;
        }

        /**
         * Returns whether the Video moments cache is empty or not.
         * @returns {boolean}
         */
        isVideoMomentsCacheEmpty = () => {
            return this.#cluster_video_moments.size === 0;
        }

        /**
         * Loads all the video moments for available for the cluster an stores them on the 
         * the local cache.
         * @returns {Promise<void>}
         */
        loadAllVideoMoments = async () => {
            const cluster_moments = await getAllClusterVideoMoments(this.#uuid);

            this.#setVideoMoments(cluster_moments);
        }

        /**
         * Returns an ordered video moments cache from list of video moments.
         * @param {import('@models/Metadata').VideoMoment[]} video_moments
         * @returns {Map<string, import('@models/Metadata').VideoMoment[]>}
         */
        #orderVideoMoments = (video_moments) => {
            const new_cache = new Map();

            for (let h = 0; h < video_moments.length; h++) {
                this.#addVideoMomentToCache(video_moments[h], new_cache);
            }

            this.#sortCacheVideoMoments(new_cache);

            return new_cache;
        }

        /**
         * Removes a given video moment from the local cache.
         * @param {import('@models/Metadata').VideoMoment} video_moment
         * @returns {void}
         */
        #removeVideoMomentLocally = video_moment => {
            const the_video_moment_siblings = this.#cluster_video_moments.get(video_moment.VideoUUID);

            if (the_video_moment_siblings === undefined) return;

            let new_video_moments = the_video_moment_siblings.filter(vm => vm.ID !== video_moment.ID);

            this.#cluster_video_moments.set(video_moment.VideoUUID, new_video_moments);
        }

        /**
         * Returns a copy of the given list of video moment but sorted by moment
         * time.
         * @param {import('@models/Metadata').VideoMoment[]} video_moments
         * @returns {import('@models/Metadata').VideoMoment[]}
         */
        #sortVideoMoments = video_moments => {
            const sorted_video_moments = [...video_moments];

            sorted_video_moments.sort((a, b) => {
                return a.StartTime - b.StartTime;
            });

            return sorted_video_moments;
        }

        /**
         * Sorts the video moments of a given video moment cache in place.
         * @param {Map<string, import('@models/Metadata').VideoMoment[]>} cache
         * @returns {void}
         */
        #sortCacheVideoMoments = cache => {

            for (let [video_uuid, video_moments] of cache.entries()) {
                const sorted_copy = this.#sortVideoMoments(video_moments);

                cache.set(video_uuid, sorted_copy);
            }

            return;
        }


        /**
         * Sets a list of video moments to replace the current cache.
         * @param {import('@models/Metadata').VideoMoment[]} video_moments
         * @returns {void}
         */
        #setVideoMoments = video_moments => {
            const new_cache = this.#orderVideoMoments(video_moments);

            this.#cluster_video_moments = new_cache;
        }
    
    /*=====  End of Video moments  ======*/

    get DownloadCategoryID() {
        return this.#filter_category;
    }
    
    get FSPath() {
        return this.#fs_path;
    }

    get Name() {
        return this.#name;
    }

    set Name(new_name) {
        this.#name = new_name;
    }

    get RootCategoryID() {
        return this.#root_category;
    }

    toJSON = () => {
        return {
            uuid: this.#uuid,
            name: this.#name,
            fs_path: this.#fs_path,
            filter_category: this.#filter_category,
            root_category: this.#root_category
        }
    }

    get UUID() {
        return this.#uuid;
    }
}

/**
 * Creates a new cluster.
 * @param {string} uuid a uuid4 string
 * @param {string} name the name of the cluster
 * @param {string} fs_path the path to the cluster's file system root
 * @param {string} filter_category the uuid of the category where content is automatically downloaded for this cluster
 * @returns {Promise<CategoriesCluster | null>} the newly created cluster
 */
export const createCluster = async (uuid, name, fs_path, filter_category) => {
    let new_cluster = null;
    
    let request = new PostCreateClusterRequest(uuid, name, fs_path, filter_category);

    const response = await request.do();

    if (response.Created) {
        new_cluster = new CategoriesCluster(response.data);
    }

    return new_cluster;
}

/**
 * Deletes a cluster record leaving the file system untouched. returns true if successful.
 * @param {string} cluster_uuid
 * @returns {Promise<boolean>}
 */
export const deleteClusterRecord = async (cluster_uuid) => {
    let is_deleted = false;

    let request = new DeleteClusterRecordRequest(cluster_uuid);

    const response = await request.do();

    is_deleted = response.data.success;

    return is_deleted;
}

/**
 * Gets all the categories clusters.
 * @returns {Promise<CategoriesCluster[]>} A list of all the categories clusters.
 */
export const getAllCategoriesClusters = async () => {
    const request = new GetCategoriesClustersRequest();

    const response = await request.do();

    /** @type {CategoriesCluster[]} */
    const categories_clusters = [];

    for (let cluster of response.data) {
        let categories_cluster = new CategoriesCluster(cluster);

        categories_clusters.push(categories_cluster);
    }

    return categories_clusters;
}

/**
 * Requests the update of cluster data, takes a cluster object as a parameter. any value can be overwritten except the uuid. different values in the passed object overwrite the current values in the cluster.
 * @param {CategoriesCluster} cluster
 * @returns {Promise<boolean>}
 */
export const updateCluster = async (cluster) => {
    let is_updated = false;

    const request = new PatchClusterRequest(cluster);

    const response = await request.do();

    is_updated = response?.data ?? false;

    return is_updated;
}

/**
 * Gets the new cluster directory options. It takes a subdirectory, which if empty will default to the service's SERVICE_CLUSTERS_ROOT. If not empty
 * the subdirectory most be a path relative to the service's SERVICE_CLUSTERS_ROOT.
 * @param {string} subdirectory 
 * @returns {Promise<import("@libs/DungeonsCommunication/services_requests/categories_cluster_requests").NewClusterDirectoryOption[]>} 
 */
export const getNewClusterDirectoryOptions = async (subdirectory) => {
    /**
     * @type {import("@libs/DungeonsCommunication/services_requests/categories_cluster_requests").NewClusterDirectoryOption[]}
     */
    let directory_creation_options = [];
    const request = new GetNewClusterDirectoryOptionsRequest(subdirectory);

    const response = await request.do();

    if (response.Ok) {
        directory_creation_options = response.data;
    }

    return directory_creation_options;
}

/**
* @typedef {Object} PathAvailability
 * @property {boolean} is_path_valid
 * @property {string} reason
*/

/**
 * Validates a candidate directory for a new cluster.
 * @param {string} unsafe_path
 * @returns {Promise<PathAvailability>}
 */
export const validateClusterDirectory = async (unsafe_path) => {
    const path_availability = {
        is_path_valid: false,
        reason: ''
    }

    const request = new GetClusterDirectoryValidationRequest(unsafe_path);

    const response = await request.do();

    if (response.Ok) {
        path_availability.is_path_valid = response.data.response;
        path_availability.reason = response.data.reason;
    }

    return path_availability;
}

/**
 * Returns the value of the cluster's default root path.
 * @returns {Promise<string>}
 * @throws {import("@libs/LiberyFeedback/lf_models").LabeledError}  
 */
export const getClusterRootPath = async () => {
    let root_path = null;
    const request = new GetClusterRootPathRequest();
    
    const response = await request.do();

    if (response.Ok) {
        root_path = response.data;
    }

    return root_path;
}    

