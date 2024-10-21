import { VariableEnvironmentContextError, LabeledError } from "@libs/LiberyFeedback/lf_models";
import { categories_server } from "../services";
import { HttpResponse, attributesToJson } from "../base";
import { dungeon_http_errors } from "../errors";


/**
 * Sends a request to retrieve all the categories clusters in the system.
 */
export class GetCategoriesClustersRequest {
    do = async () => {
        const response = await fetch(`${categories_server}/clusters`);

        let data = null;

        if (response.status <= 200 && response.status < 300) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/** 
 * Requests the update of cluster data, takes a cluster object as a parameter. any value can be overwritten except the uuid. different values in the passed object overwrite the current values in the cluster.
 */    
export class PatchClusterRequest {

    static endpoint = `${categories_server}/clusters`;

    /**
     * @param {import('@models/CategoriesClusters').CategoriesCluster} modified_cluster
     */
    constructor(modified_cluster) {
        this.uuid = modified_cluster.UUID;
        this.name = modified_cluster.Name;
        this.fs_path = modified_cluster.FSPath;
        this.filter_category = modified_cluster.DownloadCategoryID;
        this.root_category = modified_cluster.RootCategoryID;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        let response;
        let updated = false;

        try {
            response = await fetch(PatchClusterRequest.endpoint, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            updated = response.ok;
        } catch (error) {
            console.error("Error while updating cluster: ", error);
        }

        return new HttpResponse(response, updated);
    }
}

/**
* @typedef {Object} NewClusterDirectoryOption
 * @property {string} path
 * @property {string} name
*/

/**
 * Requests New cluster creation directory options.
 */
export class GetNewClusterDirectoryOptionsRequest {
    /**
     * @param {string} subdirectory - if empty the service's cluster root will be used
     */
    constructor(subdirectory) {
        this.subdirectory = encodeURI(subdirectory);
    }

    toJson = attributesToJson.bind(this);   

    /**
     * @returns {Promise<HttpResponse<NewClusterDirectoryOption[]>>}
     */
    do = async () => {
        let request_url = `${categories_server}/service-fs/new-cluster-options`;

        if (this.subdirectory != null || this.subdirectory !== "") {
            request_url += `?subdirectory=${this.subdirectory}`;
        }

        const response = await fetch(request_url);

        let data = null;

        if (response.status <= 200 && response.status < 300) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Requests the validation of a new cluster directory path.
 */
export class GetClusterDirectoryValidationRequest {
    /**
     * @param {string} unsafe_path
     */
    constructor(unsafe_path) {
        this.unsafe_path = unsafe_path;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("../base").ReasonedBooleanResponse>>}
     */
    do = async () => {
        const response = await fetch(`${categories_server}/service-fs/validate-path?unsafe_path=${this.unsafe_path}`);

        let data = null;

        if (response.status === 200) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Requests the cluster default root path
 */
export class GetClusterRootPathRequest {
    do = async () => {
        const endpoint = `${categories_server}/service-fs/cluster-root-default`;
        
        let response;
        
        try {
            response = await fetch(endpoint);
        } catch (error) {
            console.error("Error while getting cluster root path: ", error);
        }

        if (!response?.ok) {
            let variable_enviroment = new VariableEnvironmentContextError("In @libs/DungeonsCommunication/services_requests/categories_cluster_requests/GetClusterRootPathRequest after fetch");
            variable_enviroment.addVariable("endpoint", endpoint);
            variable_enviroment.addVariable("response.status", response?.status);


            const labeled_error = new LabeledError(variable_enviroment, "Failed to get cluster root path", dungeon_http_errors.ERR_TALKING_TO_ENDPOINT);

            throw labeled_error;
        }

        /**
         * @type {SingleStringResponse}
         */
        let data = null;

        if (response.status === 200) {
            data = await response.json();
        }

        return new HttpResponse(response, data?.response ?? "");
    }
}

/**
 * Requests the creation of a new cluster
 */
export class PostCreateClusterRequest {
    /**
     * @param {string} uuid
     * @param {string} name
     * @param {string} fs_path
     * @param {string} filter_category
     */
    constructor(uuid, name, fs_path, filter_category) {
        this.uuid = uuid;
        this.name = name;
        this.fs_path = fs_path;
        this.filter_category = filter_category;
        this.root_category = "";
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<CreatedCategoryCluster>>}
     * @typedef {Object} CreatedCategoryCluster
     * @property {string} uuid
     * @property {string} name
     * @property {string} fs_path
     * @property {string} filter_category
     * @property {string} root_category
     */
    do = async () => {
        const response = await fetch(`${categories_server}/clusters`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        let data = null;

        if (response.status === 201) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Deletes a cluster record, not its files in the filesystem.
 * @param {string} cluster_id
 */
export class DeleteClusterRecordRequest {
    constructor(cluster_id) {
        if (cluster_id == null) {
            throw new Error("cluster_id cannot be null or undefined");  
        }
        this.cluster = cluster_id;
    }

    toJson = attributesToJson.bind(this);   
    
    /**
     * @returns {Promise<HttpResponse<DeletionSuccesful>>}
     * @typedef {Object} DeletionSuccesful
     * @property {boolean} success
     */
    do = async () => {

        let deletion_state = {
            success: false
        }
        let response = null;

        try {
            response = await fetch(`${categories_server}/clusters/platform-data?cluster=${this.cluster}`, {
                method: "DELETE"
            });

            deletion_state.success = response.status === 204;   
        } catch (error) {
            console.error("Error while deleting cluster record: ", error);
        }

        return new HttpResponse(response, deletion_state);
    }
}

/**
 * Requests the resyncronization of a cluster branch by it's category uuid
 * @param {string} cluster_uuid
 * @param {string} from_category_uuid
 */
export class PutResyncClusterBranchRequest {

    static endpoint = `${categories_server}/service-fs/cluster-sync`;

    /**
     * @param {string} cluster_uuid
     * @param {string} from_category_uuid
     */
    constructor(cluster_uuid, from_category_uuid) {
        this.cluster_uuid = cluster_uuid;
        this.from_category_uuid = from_category_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        let response;
        let rsync_successful = false;
        
        try {
            response = await fetch(`${PutResyncClusterBranchRequest.endpoint}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });
        } catch (error) {
            console.error("Error while resyncing cluster branch: ", error);
        }

        if (response?.ok) {
            rsync_successful = true;
        }

        return new HttpResponse(response, rsync_successful);
    }
}

/*=============================================
=            App claims            =
=============================================*/

    export class GetCategoriesClusterSignAccessRequest {
        constructor(cluster_id) {
            this.cluster_id = cluster_id;
        }

        toJson = attributesToJson.bind(this);

        do = async () => {
            const response = await fetch(`${categories_server}/clusters/sign-access?cluster=${this.cluster_id}`);

            let data = null;

            if (response.status <= 200 && response.status < 300) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

/*=====  End of App claims  ======*/
