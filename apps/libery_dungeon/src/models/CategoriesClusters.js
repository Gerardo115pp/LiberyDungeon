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
import { MediaIdentity } from "./Medias";
import { getVideoMoments, createVideoMoment } from "./Metadata";

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
    }

    
    /*=============================================
    =            Video moments            =
    =============================================*/

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

            this.#cluster_video_moments.set(media_uuid, fetched_video_moments);

            return fetched_video_moments;
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
    
    /*=====  End of Video moments  ======*/
    
    

    get DownloadCategoryID() {
        return this.#filter_category;
    }
    
    get FSPath() {
        return this.#fs_path;
    }

    /**
     * Returns a list of media identities that belong in this cluster.
     * @param {string[]} media_uuids
     * @returns {Promise<import("@models/Medias").MediaIdentity[]>}
     */
    async getClusterMedias(media_uuids) {
        const request = new GetMediaIdentitiesByUUIDsRequest(media_uuids, this.#uuid);

        /**
         *  @type {MediaIdentity[]} 
         */
        let media_identities = [];

        const response = await request.do();

        if (response.Ok && response.data != null) {
            media_identities = response.data.map(media_identity_params => new MediaIdentity(media_identity_params));
        }

        return media_identities;
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

