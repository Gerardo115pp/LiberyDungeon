import { 
    GetCategoriesClustersRequest, 
    PatchClusterRequest,
    GetNewClusterDirectoryOptionsRequest,
    GetClusterDirectoryValidationRequest,
    PostCreateClusterRequest,
    DeleteClusterRecordRequest,
    GetClusterRootPathRequest
 } from "@libs/DungeonsCommunication/services_requests/categories_cluster_requests";

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

    constructor({ uuid, name, fs_path, filter_category, root_category }) {
        if (uuid === undefined || name === undefined || fs_path === undefined || filter_category === undefined || root_category === undefined) {
            throw new Error('Missing required parameters.');
        }

        this.#uuid = uuid;
        this.#name = name;
        this.#fs_path = fs_path;
        this.#filter_category = filter_category;
        this.#root_category = root_category;
    }

    get UUID() {
        return this.#uuid;
    }

    get Name() {
        return this.#name;
    }

    set Name(new_name) {
        this.#name = new_name;
    }
    
    get FSPath() {
        return this.#fs_path;
    }

    get DownloadCategoryID() {
        return this.#filter_category;
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
}

/**
 * Creates a new cluster.
 * @param {string} uuid a uuid4 string
 * @param {string} name the name of the cluster
 * @param {string} fs_path the path to the cluster's file system root
 * @param {string} filter_category the uuid of the category where content is automatically downloaded for this cluster
 * @returns {Promise<CategoriesCluster>} the newly created cluster
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
 * @returns {Promise<import("@libs/DungeonsCommunication/services_requests/categories_cluster_requests").NewClusterDirectoryOption>} 
 */
export const getNewClusterDirectoryOptions = async (subdirectory) => {
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

