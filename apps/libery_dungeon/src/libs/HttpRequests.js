import { attributesToJson, attributesToJsonExclusive, createUnsecureJWT, parseUnsecureJWT } from "@libs/utils";
import { DownloadProgress } from "@models/Downloads";
import { LabeledError, VariableEnvironmentContextError } from "./LiberyFeedback/lf_models";

export const pw_server = PW_SERVER;
export const base_domain = BASE_DOMAIN;
export const categories_server = CATEGORIES_SERVER;
export const medias_server = MEDIAS_SERVER;
export const collect_server = COLLECT_SERVER;
export const downloads_server = DOWNLOADS_SERVER;
export const metadata_server = METADATA_SERVER;
export const jd_server = JD_SERVER; 
export const users_server = USERS_SERVER;

const dungeon_http_errors = {
    ERR_TALKING_TO_ENDPOINT: "ERR_TALKING_TO_ENDPOINT",
}

const media_fs_endpoint = `${MEDIAS_SERVER}/medias-fs`;
const thumbnails_fs_endpoint = `${MEDIAS_SERVER}/thumbnails-fs`;

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

/**
 * A base class for all services responses
 * @template T
 */
class HttpResponse {
    /**
     * @type {number}
     */
    status;

    /**
     * @type {T}
     */
    data;

    /**
     * @type {Headers}
     */
    headers;

    /**
     * @param {Response} response
     * @param {T} data
     */
    constructor(response, data) {
        this.status = response.status;
        this.data = data;
        this.headers = response.headers;
    }

    /**
     * Returns true if the response status larger than or equal to 200 and less than 300
     * @returns {boolean}
     */
    get Ok() {
        return this.status >= 200 && this.status < 300;
    }
    
    /**
     * indicates that the request has succeeded and has led to the creation of a resource. The new resource, or a description and link to the new resource, is effectively created before
     * the response is sent back and the newly created items are returned in the body of the message, located at either the URL of the request, or at the URL in the value of the 
     * Location header.
     * @returns {boolean}
     */
    get Created() {
        return this.status === 201;
    }
}

/**
 * @typedef {Object} ReasonedBooleanResponse
 * @property {boolean} response
 * @property {string} reason
 */
/**
 * @typedef {Object} BooleanResponse
 * @property {boolean} response
 */
/**
 * @typedef {Object} SingleStringResponse
 * @property {string} response
*/

/**
 * Gets a leaf of the categories tree using the leaf id which is provided by the parents inner categories attribute. if you want the root of the tree, pass 'main' as the leaf id.
 * @date 10/20/2023 - 10:38:16 PM
 *
 * @export
 * @class GetCategoryDataRequest
 * @typedef {GetCategoryTreeLeafRequest}
 */
export class GetCategoryTreeLeafRequest {
    constructor(category_id, cluster_id) {
        this.category_id = category_id;
        this.cluster_id = cluster_id;   
    }

    toJson = attributesToJson.bind(this);

    
    /**
     * Sends the http request to the server and returns the response.
     * @returns {Promise<HttpResponse>}
     */
    do = async () => {
        let request_url = `${categories_server}/categories-tree?category_id=${this.category_id}`;

        if (this.cluster_id != null) {
            request_url += `&cluster_id=${this.cluster_id}`;
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
 * Gets a short version of the categories tree
 * @date 10/20/2023 - 10:38:16 PM
*/
export class getShortCategoryTreeRequest {
    constructor(category_id, cluster_id) {
        this.category_id = category_id;
        this.cluster_id = cluster_id;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {
        let request_url = `${categories_server}/categories-tree/short?category_id=${this.category_id}`;

        if (this.cluster_id != null) {
            request_url += `&cluster_id=${this.cluster_id}`;
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
 * Requests the only the data of a category, no content or children categories  
 */
export class GetCategoryRequest {
    constructor(category_id) {
        this.category_id = category_id;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {
        const response = await fetch(`${categories_server}/categories/data?category_id=${this.category_id}`);

        let data = null;

        if (response.status <= 200 && response.status < 300) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Returns a category from a given fullpath and cluster id
 */
export class GetCategoryByFullpathRequest {
    constructor(category_path, category_cluster) {
        this.category_path = encodeURI(category_path);
        this.category_cluster = category_cluster;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {
        const response = await fetch(`${categories_server}/categories/by-fullpath?category_path=${this.category_path}&category_cluster=${this.category_cluster}`);

        let data = null;

        if (response.ok) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Applies the changes made to the categories tree
*/
export class CommitCategoryTreeChangesRequest {
    constructor(category_id, rejected_medias, moved_medias) {
        this._category_id = category_id;
        this.rejected_medias = rejected_medias;
        this.moved_medias = moved_medias;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {
        const response = await fetch(`${categories_server}/categories?category_id=${this._category_id}`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        return new HttpResponse(response, {});
    }
}

/**
 * Renames a category, changes the full path of the category based on the new name and does the same for all subcategories full paths
 */
export class PatchRenameCategoryRequest {
    constructor(category_id, new_name) {
        this.category_id = category_id;
        this.new_name = new_name;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @typedef {Object} PatchRenameCategoryResponse
     * @property {string} uuid
     * @property {string} name
     * @property {string} parent
     * @property {string} fullpath
     * @property {string} cluster
     */
    

    /**
     * @returns {Promise<HttpResponse<PatchRenameCategoryResponse>>}
     */
    do = async () => {
        const response = await fetch(`${categories_server}/categories/rename`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        let data = null;

        if (response.status <= 200 && response.status < 300) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Moves a category to another parent category
 */
export class PatchMoveCategoryRequest {
    constructor(new_parent_category, moved_category) {
        this.new_parent_category = new_parent_category;
        this.moved_category = moved_category;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @typedef {Object} PatchMoveCategoryResponse
     * @property {string} uuid
     * @property {string} name
     * @property {string} parent
     * @property {string} fullpath
     * @property {string} cluster
     */

    /**
     * @returns {Promise<HttpResponse<PatchMoveCategoryResponse>>}
     */
    do = async () => {
        const response = await fetch(`${categories_server}/categories-tree/move`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        let data = null;

        if (response.status <= 200 && response.status < 300) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}

/**
 * Gets all the categories that match a search query
*/
export class GetCategorySearchResultsRequest {
    constructor(query, cluster_id, ignore) {
        this.query = query;
        this.cluster_id = cluster_id;
        this.ignore = ignore;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {

        let request_url = `${categories_server}/search?query=${this.query}&ignore=${this.ignore}`;

        if (this.cluster_id != null) {
            request_url += `&cluster_id=${this.cluster_id}`;
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
 * Sends a request to delete a category, it has an optional parameter `force` which if set to true, will delete the category even if it has medias inside or subcategories.
*/ 
export class DeleteCategoryRequest {
    constructor(category_id, force=false) {
        this.category_id = category_id;
        this.force = force;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {
        const response = await fetch(`${categories_server}/categories`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        return new HttpResponse(response, {});
    }
}

/**
 * Asks the server if a category name is available on a specific parent category
*/
export class GetCategoryNameAvailabilityRequest {
    constructor(category_name, parent_category_id) {
        this.category_name = category_name;
        this.parent_category_id = parent_category_id;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {
        const response = await fetch(`${categories_server}/categories/name-available?category_name=${this.category_name}&parent_id=${this.parent_category_id}`);

        let is_available = null;

        if (response.status === 200) {
            is_available = true;
        } else if (response.status === 409) {
            is_available = false;
        }

        return new HttpResponse(response, is_available);
    }
}

/**
 * Sends a request to create a new category
*/
export class PostCreateCategoryRequest {
    constructor(category_name, parent_category_id, parent_path, cluster) {
        this.name = category_name;
        this.parent = parent_category_id;
        this.parent_path = parent_path;
        this.cluster = cluster;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {
        let response;
        try {
            response = await fetch(`${categories_server}/categories`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });
        } catch (error) {
            console.log("Error while making a create category request: ", error);
            throw error;
        }

        let data = null;

        if (response.status >= 200 && response.status < 300) {
            data = await response.json();
        }

        return new HttpResponse(response, data);
    }
}


/*=============================================
=            Trashcan            =
=============================================*/

    export class GetTrashcanEntriesRequest {
        static endpoint = `${categories_server}/trashcan/entries`;

        /**
         * @returns {Promise<HttpResponse<import('@models/Trashcan').TrashcanTransactionEntryParams[]>>}
         */
        do = async () => {
            let response;

            try {
                response = await fetch(GetTrashcanEntriesRequest.endpoint);
            } catch (error) {
                console.error("Error while getting trashcan entries: ", error);
            }

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    export class GetTrashcanTransactionRequest {
        static endpoint = `${categories_server}/trashcan/transaction`;

        /**
         * @param {string} transaction_id - a timestamp with the format(in go's time format) '2006-01-02 15:04:05' that also serves as a transaction id
         */
        constructor(transaction_id) {
            this.transaction_id = transaction_id;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<import('@models/Trashcan').TrashcanTransactionParams>>}
         */
        do = async () => {
            let response;
            let data = {};

            try {  
                response = await fetch(`${GetTrashcanTransactionRequest.endpoint}?transaction_id=${this.transaction_id}`);

                if (response.ok) {
                    data = await response.json();
                }
            } catch (error) {
                console.error("Error while getting trashcan transaction: ", error); 
            }

            return new HttpResponse(response, data);
        }
            
    }

    export class PatchRestoreMediaRequest {
        static endpoint = `${categories_server}/trashcan/media/restore`;

        /**
         * @param {string} media_uuid - the uuid of the media to restore    
         * @param {string} transaction_id - the transaction id that the media was deleted in
         * @param {string} main_category - the category where the media will be restored
         */
        constructor(media_uuid, transaction_id, main_category) {
            this.media_uuid = media_uuid;
            this.transaction_id = transaction_id;
            this.main_category = main_category;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>>}
         */
        do = async () => {
            let response;
            let restored = false;

            try {
                response = await fetch(`${PatchRestoreMediaRequest.endpoint}?media_uuid=${this.media_uuid}&transaction_id=${this.transaction_id}&main_category=${this.main_category}`, {
                    method: "PATCH",    
                });

                restored = response.ok; 
            } catch (error) {
                console.error("Error while restoring media: ", error);
            }

            return new HttpResponse(response, restored);
        }
    }

    export class PatchRestoreTransactionRequest {
        static endpoint = `${categories_server}/trashcan/transaction/restore`;  

        /**
         * @param {string} transaction_id - the transaction id that will be restored
         * @param {string} main_category - the category where the transaction will be restored
         */
        constructor(transaction_id, main_category) {
            this.transaction_id = transaction_id;
            this.main_category = main_category;
        }   

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>>}
         */
        do = async () => {
            let response;
            let restored = false;

            try {
                response = await fetch(`${PatchRestoreTransactionRequest.endpoint}?transaction_id=${this.transaction_id}&main_category=${this.main_category}`, {
                    method: "PATCH",
                });

                restored = response.ok;
            } catch (error) {
                console.error("Error while restoring transaction: ", error);
            }

            return new HttpResponse(response, restored);
        }
    }

    export class DeleteTrashcanTransactionRequest {
        static endpoint = `${categories_server}/trashcan/transaction`;  

        /**
         * @param {string} transaction_id - the transaction id that will be deleted
         */
        constructor(transaction_id) {
            this.transaction_id = transaction_id;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>>}
         */
        do = async () => {
            let response;
            let deleted = false;

            try {
                response = await fetch(`${DeleteTrashcanTransactionRequest.endpoint}?transaction_id=${this.transaction_id}`, {
                    method: "DELETE",
                });

                deleted = response.ok;
            } catch (error) {
                console.error("Error while deleting transaction: ", error);
            }

            return new HttpResponse(response, deleted);
        }
    }

    export class DeleteEmptyTrashcanRequest {

        static endpoint = `${categories_server}/trashcan/empty`;

        /**
         * @returns {Promise<HttpResponse<boolean>>}
         */
        do = async () => {
            let response;
            let deleted = false;

            try {
                response = await fetch(DeleteEmptyTrashcanRequest.endpoint, {
                    method: "DELETE",
                });

                deleted = response.ok;
            } catch (error) {
                console.error("Error while deleting trashcan: ", error);
            }

            return new HttpResponse(response, deleted);
        }
    }

/*=====  End of Trashcan  ======*/

/*=============================================
=            Medias                           =
=============================================*/

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
            }

            return new HttpResponse(response, renamed);
        }
    }
    
    /*----------  Uploads  ----------*/
    
        /**
         * Sends a request to have an Upload Stream ticket signed and set as a cookie.
         */
        export class GetUploadStreamTicketRequest {
            static endpoint = `${medias_server}/upload-streams/stream-ticket`;

            /**
             * @param {string} upload_uuid
             * @param {int} total_medias
             * @param {string} category_uuid
             */
            constructor(upload_uuid, total_medias, category_uuid) {
                this.upload_uuid = upload_uuid;
                this.total_medias = total_medias;
                this.category_uuid = category_uuid;
            }

            toJson = attributesToJson.bind(this);

            /**
             * @returns {Promise<HttpResponse<BooleanResponse>}
             */
            do = async () => {
                let response;
                let ticket_signed_response = {
                    response: false
                }

                let request_url = GetUploadStreamTicketRequest.endpoint;
                request_url += `?upload_uuid=${this.upload_uuid}`;
                request_url += `&total_medias=${this.total_medias}`;
                request_url += `&category_uuid=${this.category_uuid}`;

                console.log("Requesting upload stream ticket: ", request_url);

                let url = new URL(request_url, location.origin);
                const request = new Request(url);

                try {
                    response = await fetch(request);
                } catch (error) {
                    console.error("Error while getting upload stream ticket: ", error);
                }

                if (response.ok) {
                    ticket_signed_response = await response.json();
                }
                
                return new HttpResponse(response, ticket_signed_response);    
            }
        }

        /**
        * @typedef {Object} GetChunkedUploadTicketParams
        * @property {string} upload_uuid - a client side generated uuid that represents the upload
        * @property {string} upload_filename - the name of the file that will be uploaded
        * @property {string} upload_size - the total size of the file that will be uploaded
        * @property {string} upload_chunks - the amount of chunks that the file will be split into
        * @property {string} category_uuid - the category where the file will be uploaded
        */

        /**
         * Sends a request to add a media to a category
         */
        export class PostMediaRequest {

            static endpoint = `${medias_server}/upload-streams/stream-fragment`;


            /**
             * @param {string} media_name
             * @param {Blob} media_blob
            */
            constructor(media_name, media_blob) {
                this.media_name = media_name;
                this.media_blob = media_blob;
            }

            toJson = attributesToJson.bind(this);

            do = async () => {
                const form_data = new FormData();

                form_data.append(this.media_name, this.media_blob);

                let url_string = PostMediaRequest.endpoint;

                const url = new URL(url_string, location.origin);
                
                const request = new Request(url, {
                    method: "POST",
                    body: form_data
                });

                const response = await fetch(request);

                return new HttpResponse(response, {});
            }
        }

        /**
         * Gets an upload ticket for a chunked upload.
         */
        export class GetChunkedUploadTicketRequest {
            static endpoint = `${medias_server}/upload-streams/chunked-ticket`;

            /**
             * @param {GetChunkedUploadTicketParams} param0
             */
            constructor({upload_uuid, upload_filename, upload_size, upload_chunks, category_uuid}) {
                this.upload_uuid = upload_uuid;
                this.upload_filename = upload_filename;
                this.upload_size = upload_size;
                this.upload_chunks = upload_chunks;
                this.category_uuid = category_uuid;
            }

            toJson = attributesToJson.bind(this);

            /**
             * Returns the upload ticket which is a jwt token. Must be passed in the Authorization header of each chunk upload request.
             * @returns {Promise<HttpResponse<string|null>>}
             */
            do = async () => {
                let upload_ticket = null;
                let response;

                let request_endpoint = GetChunkedUploadTicketRequest.endpoint;
                request_endpoint += `?upload_uuid=${this.upload_uuid}`
                request_endpoint += `&upload_filename=${this.upload_filename}`
                request_endpoint += `&upload_size=${this.upload_size}`
                request_endpoint += `&upload_chunks=${this.upload_chunks}`
                request_endpoint += `&category_uuid=${this.category_uuid}`


                try {
                    response = await fetch(request_endpoint, {
                        method: "GET",
                        headers: {
                            "Content-Type": "application/json"
                        }
                    });
                } catch (error) {
                    console.error("Error while getting upload ticket: ", error);
                }

                if (response.ok) {
                    /**
                     * @type {SingleStringResponse}
                     */
                    let data = await response.json();
                    upload_ticket = data.response;
                }

                return new HttpResponse(response, upload_ticket);
            }
        }

        /**
         * Sends a request to upload a chunk of a file
         */
        export class PostChunkedUploadRequest {
            static endpoint = `${medias_server}/upload-streams/chunked-upload`;

            /**
             * The upload ticket for the chunked upload
             * @type {string} 
             */
            #upload_ticket;

            /**
             * @param {string} upload_ticket
             * @param {Blob} chunk
             * @param {number} chunk_serial
             */
            constructor(upload_ticket, chunk, chunk_serial) {
                this.#upload_ticket = upload_ticket;
                this.chunk = chunk;
                this.chunk_serial = chunk_serial;
            }

            toJson = () => null;

            /**
             * @returns {Promise<HttpResponse<boolean>>}
             */
            do = async () => {
                const form_data = new FormData();

                form_data.append("chunk", this.chunk);

                const headers = new Headers();

                headers.append("Authorization", this.#upload_ticket);

                let response;

                try {
                    response = await fetch(`${PostChunkedUploadRequest.endpoint}?chunk_serial=${this.chunk_serial}`, {
                        method: "POST",
                        headers: headers,
                        body: form_data
                    });
                } catch (error) {
                    console.error("Error while uploading chunk: ", error);
                }

                let was_uploaded = response?.ok === true;

                return new HttpResponse(response, was_uploaded);
            }


        }

/*=====  End of Media  ======*/

/*=============================================
=            Clusters            =
=============================================*/

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
         * @returns {Promise<HttpResponse<ReasonedBooleanResponse>>}
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
                let variable_enviroment = new VariableEnvironmentContextError("In HttpRequests/GetClusterRootPathRequest after fetch");
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

/*=====  End of Clusters  ======*/

/*=============================================
=            Collect&Download requests            =
=============================================*/

    /**
     * Requests the 4chan tracked boards
    */
    export class GetChanBoardsRequest {
        constructor() {
        }
    
        toJson = attributesToJson.bind(this);
    
        do = async () => {
            const response = await fetch(`${collect_server}/4chan-boards/boards`);

            let data = null;

            if (response.status <= 200 && response.status < 300) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }
    
    /**
     * Requests a board's thread catalog
    */
    export class GetChanCatalogRequest {
        constructor(board_name) {
            this.board_name = board_name;
        }
    
        toJson = attributesToJson.bind(this);
    
        do = async () => {
            const response = await fetch(`${collect_server}/4chan-boards/board/catalog?board_name=${this.board_name}`);

            let data = null;

            if (response.status <= 200 && response.status < 300) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    /**
     * Requests a thread's contents
     * @param {string} thread_id
     * @param {string} board_name
     */
    export class GetChanThreadRequest {
        constructor(thread_id, board_name) {
            this.thread_id = thread_id;
            this.board_name = board_name;
        }

        toJson = attributesToJson.bind(this);

        do = async () => {
            const response = await fetch(`${collect_server}/4chan-threads/thread?thread_id=${this.thread_id}&board_name=${this.board_name}`);

            let data = null;

            if (response.status <= 200 && response.status < 300) {  
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    export class PostThreadDownloadRequest {
        /** 
         * @param {string} thread_uuid
         * @param {string} board_name
         * @param {string} target_category_name
         * @param {string} parent_uuid
        */
        constructor(thread_uuid, board_name, target_category_name, parent_uuid, cluster_uuid) {
            this.thread_uuid = thread_uuid;
            this.board_name = board_name;
            this.target_category_name = target_category_name;
            this.parent_uuid = parent_uuid;
            this.cluster_uuid = cluster_uuid;
        }


        toJson = attributesToJsonExclusive.bind(this);

        do = async () => {
            const response = await fetch(`${collect_server}/4chan-downloads/thread/images`, {   
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            }); 

            let data = {};

            if (response.status <= 200 && response.status < 300) {
                data = await response.json();
            }

            return new HttpResponse(response, {
                download_uuid: data.download_uuid
            });
        }


    }

    export class GetDownloadRegisterRequest {
        constructor(download_uuid) {
            this.download_uuid = download_uuid;
        }

        toJson = attributesToJson.bind(this);

        do = async () => {
            const response = await fetch(`${downloads_server}/download-history/download?download_uuid=${this.download_uuid}`);

            let data = null;

            if (response.status <= 200 && response.status < 300) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    /**
     * Requests the current download uuid, there is none then the download_uuid attribute will be an empty string
     */
    export class GetCurrentDownloadRequest {
        toJson = attributesToJson.bind(this);

        do = async () => {
            const response = await fetch(`${downloads_server}/downloads/current-download`);

            let data = null;

            if (response.status <= 200 && response.status < 300) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }
/*=====  End of Collect&Download requests  ======*/

/*=============================================
=            Websocket transmisor            =
=============================================*/
// TODO: move this to a separate file

    export class DownloadProgressTransmisor {
        constructor(download_uuid) {
            this.download_uuid = download_uuid;
            this.host = `wss://${base_domain}${downloads_server}/ws/download-progress?download_uuid=${this.download_uuid}`;
            this.socket = null;

            this.download_progress_callback = null;
            this.download_completed_callback = null;
        }

        connect = () => {
            this.socket = new WebSocket(this.host);
            this.socket.onopen = this.onOpen;
            this.socket.onmessage = this.onMessage;
            this.socket.onclose = this.onClose;
            this.socket.onerror = this.onError;
        }

        disconnect = () => {
            this.socket.close();
        }

        onMessage = message => {
            console.debug("received message from server: ", message)
            const data = JSON.parse(message.data);

            let download_progress;

            try {
                download_progress = new DownloadProgress(data);
            } catch (error) {
                console.error(error);
            }

            if (this.download_progress_callback !== null) {
                this.download_progress_callback(download_progress);
            }   
        }

        onOpen = () => {}

        onClose = () => {}

        onError = error => {}
    }

    export class PlatformEventsTransmisor {

        /**
         * The socket used for communication
         * @type {WebSocket}
         */
        #socket;

        /**
         * set by the caller, valid events received will be forwarded to this callback  
         * @param {PlatformEventMessage} event_message
         */
        #on_message_callback = event_message => console.log("Received event: ", event_message);
        
        constructor() {
            this.host = `wss://${base_domain}${jd_server}/platform-events/public/suscribe`;
            this.#socket = null;
        }

        connect = () => {
            this.#socket = new WebSocket(this.host);
            this.#socket.onopen = this.onOpen;
            this.#socket.onmessage = this.onMessage;
            this.#socket.onclose = this.onClose;
            this.#socket.onerror = this.onError;
        }

        disconnect = () => {
            if (this.#socket == null) return;
            
            this.#socket.close();
        }

        /**
         * @param {MessageEvent} message 
         */
        onMessage = message => {
            console.debug("received message from server: ", message)
            const data = JSON.parse(message.data);

            /**
             * @type {PlatformEventMessage}
             */
            let event_message;

            try {
                event_message = new PlatformEventMessage(data);
            } catch (error) {
                console.error(`Wrong event message format: ${error}`);
            }

            this.#on_message_callback(event_message);
        }

        onOpen = () => {}

        onClose = () => {}

        onError = error => {}

        /**
         * Sets the event that will handle new messages
         * @param {(event_message: PlatformEventMessage) => void} callback
         */
        setOnMessageCallback = callback => {
            if (callback?.constructor.name !== "Function" && callback?.constructor.name !== "AsyncFunction") {
                console.log("The callback must be a function, got: ", callback);
                throw new TypeError("The callback must be a function, got: ", callback);
            }

            this.#on_message_callback = callback;
        }
    }

    
    /*=============================================
    =            Messages            =
    =============================================*/
    
        /**
         * @template T
         */
        export class PlatformEventMessage {
            /**
             * Some endpoints require async processing, in which case they will return a UUID to track the event's status
             * @type {string}
             */
            #event_uuid;

            /**
             * Identify the type of effect this event has
             * @type {string}
             */
            #event_type;

            /**
             * Human readable message
             * @type {string}
             */
            #event_message;

            /**
             * Event payload encoded as a signed JWT token
             * @type {string}
             */
            #event_payload;

            /**
             * @param {PlatformEventMessageParams} param0
             * @typedef {Object} PlatformEventMessageParams
             * @property {string} uuid
             * @property {string} event_type
             * @property {string} event_message
             * @property {string} event_payload
             */
            constructor({uuid, event_type, event_message, event_payload}) {
                this.#event_uuid = uuid;
                this.#event_type = event_type;
                this.#event_message = event_message;
                this.#event_payload = event_payload;
            }

            get UUID() {
                return this.#event_uuid;
            }

            get EventType() {
                return this.#event_type;
            }

            get EventMessage() {
                return this.#event_message;
            }

            /**
             * Parses the event payload without verifying the signature
             * @returns {EventPayload<T>}
             * @template T
             * @typedef {Object} EventPayload
             * @property {Object} header
             * @property {T} payload
             */
            get EventPayload() {
                let parsed_payload = parseUnsecureJWT(this.#event_payload);

                return parsed_payload;
            }
        } 
    
    /*=====  End of Messages  ======*/
    
    

/*=====  End of Websocket transmisor  ======*/

/*=============================================
=            Metadata            =
=============================================*/

    export class GetWatchPointRequest {
        constructor(media_uuid) {
            this.media_uuid = media_uuid;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<WatchPointTime>>}
         * @typedef {Object} WatchPointTime
         * @property {number} start_time 
         */
        do = async () => {
            const response = await fetch(`${metadata_server}/watch-points?media_uuid=${this.media_uuid}`);

            let data = null;
            let http_response = new HttpResponse(response, data);

            if (!(response.status <= 200 && response.status < 300) || !response.headers.has("Content-Type")) {
                console.error("Error getting watch point: ", response);
                return http_response;
            }

            switch (response.headers.get("Content-Type")) {
                case "application/json":
                    data = await response.json();
                    http_response = new HttpResponse(response, data);
                    break;
                case "application/octet-stream":
                    let response_data = {};
                    data = await response.arrayBuffer();
                    let data_view = new DataView(data);
                    response_data.start_time = data_view.getUint32(0, true);
                    http_response = new HttpResponse(response, response_data);
                    break;
            }

            return http_response;
        }
    }

    export class PostWatchPointRequest {
        constructor(media_uuid, start_time) {
            this.media_uuid = media_uuid;
            this.start_time = start_time;
        }

        toJson = attributesToJson.bind(this);

        do = async () => {
            const response = await fetch(`${metadata_server}/watch-points`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            return new HttpResponse(response, {});
        }
    }

/*=====  End of Metadata  ======*/

/*=============================================
=            Users            =
=============================================*/

    export class GetIsInitialSetupRequest {
        static endpoint = `${users_server}/is-initial-setup`;

        /**
         * @returns {Promise<HttpResponse<BooleanResponse>}
         */
        do = async () => {
            const response = await fetch(GetIsInitialSetupRequest.endpoint);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    export class PostCreateInitialUserRequest {
        static endpoint = `${users_server}/users`;

        /**
         * @param {string} username
         * @param {string} secret
         * @param {string} initial_setup_secret
         */
        constructor(username, secret, initial_setup_secret) {
            this.username = username;
            this.secret = secret;
            this._initial_setup_secret = initial_setup_secret;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<SingleStringResponse>}
         */
        do = async () => {
            const response = await fetch(`${PostCreateInitialUserRequest.endpoint}?initial-setup-secret=${this._initial_setup_secret}`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    export class PostCreateUserRequest {
        static endpoint = `${users_server}/users`;

        /**
         * @param {string} username
         * @param {string} secret
         */
        constructor(username, secret) {
            this.username = username;
            this.secret = secret;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<SingleStringResponse>>}
         */
        do = async () => {
            const response = await fetch(PostCreateUserRequest.endpoint, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    export class GetUserSignAccessRequest {
        static endpoint = `${users_server}/user-auth`;

        /**
         * @param {string} username
         * @param {string} secret
         */
        constructor(username, secret) {
            this.username = username;
            this.secret = secret;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<UserSignAccessResponse>}
         * @typedef {Object} UserSignAccessResponse
         * @property {boolean} granted
         * @property {import('@models/Users').UserIdentityParams} user_data
         */
        do = async () => {
            let data = {
                granted: false,
                user_data: {}
            };

            const response = await fetch(`${GetUserSignAccessRequest.endpoint}?username=${this.username}&secret=${this.secret}`);


            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    export class GetUserAccessTokenValidationRequest {
        static endpoint = `${users_server}/user-auth/verify`;

        /**
         * Requests the validation of a user access token(present as an http only cookie) the request will return 200
         * regardless of the token's validity. a boolan response indicates whether the token is valid or not.
         * @returns {Promise<HttpResponse<BooleanResponse>}
         */
        do = async () => {
            const response = await fetch(GetUserAccessTokenValidationRequest.endpoint);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    /**
     * Requests the user's identity, relies on the user agent already having a valid access token stored.
     */
    export class GetUserIdentityRequest {
        static endpoint = `${users_server}/users/identity`;

        /**
         * @returns {Promise<HttpResponse<import('@models/Users').UserIdentityParams>}
         */
        do = async () => {
            const response = await fetch(GetUserIdentityRequest.endpoint);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    export class GetUserSignOutRequest {
        static endpoint = `${users_server}/user-auth/logout`;

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const response = await fetch(GetUserSignOutRequest.endpoint);

            let logged_out = false;

            if (response.ok) {
                logged_out = true;
            }

            return new HttpResponse(response, logged_out);
        }
    }

    /**
     * Returns all the users as user entries(only the username and the uuid). to use this endpoint,
     * the user token must have the 'read_users' grant.
     */
    export class GetAllUsersRequest {
        static endpoint = `${users_server}/users/read-all`;  

        /**
         * @returns {Promise<HttpResponse<import('@models/Users').UserEntry[]>}
         */
        do = async () => {
            const response = await fetch(GetAllUsersRequest.endpoint);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    export class GetAllRoleLabelsRequest {
        static endpoint = `${users_server}/roles/read-all`;

        /**
         * @returns {Promise<HttpResponse<string[]>}
         */
        do = async () => {
            const response = await fetch(GetAllRoleLabelsRequest.endpoint);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    /**
     * Returns all the grants which are just a string array. to use this endpoint, the user token must have the 'grant_option' grant, which only a
     * super admin has.
     */
    export class GetAllGrantsRequest {
        static endpoint = `${users_server}/roles/grant/read-all`;

        /**
         * @returns {Promise<HttpResponse<string[]>}
         */
        do = async () => {
            const response = await fetch(GetAllGrantsRequest.endpoint);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    /**
     * Registers a new grant to be used in by any role. to use this endpoint, the user token must have the 'grant_option' grant, which only a
     * super admin has.
     */
    export class PostCreateGrantRequest {
        static endpoint = `${users_server}/roles/grant`;

        /**
         * @param {string} new_grant
         */
        constructor(new_grant) {
            this.new_grant = new_grant;
        }


        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = `${PostCreateGrantRequest.endpoint}?new_grant=${this.new_grant}`;

            const response = await fetch(request_url, {
                method: "POST"
            });

            let created = false;

            if (response.status === 201) {
                created = true;
            }

            return new HttpResponse(response, created);
        }
    }

    /**
     * Links a grant to a role. to use this endpoint, the user token must have the 'grant_option' grant, which only a
     * super admin has.
     */
    export class PostLinkGrantToRoleRequest {
        static endpoint = `${users_server}/roles/add-grant`;

        /**
         * @param {string} role_label
         * @param {string} grant
         */
        constructor(role_label, grant) {
            this.role_label = role_label;
            this.grant = grant;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = `${PostLinkGrantToRoleRequest.endpoint}?role_label=${this.role_label}&grant=${this.grant}`;

            const response = await fetch(request_url, {
                method: "POST"
            });

            let linked = false;

            if (response.status === 201) {
                linked = true;
            }

            return new HttpResponse(response, linked);
        }
    }

    /**
     * Requests the role taxonomy of a given role label. requires the 'grant_option' grant.
     */
    export class GetRoleTaxonomyRequest {
        static endpoint = `${users_server}/roles/role`;

        /**
         * @param {string} role_label
         */
        constructor(role_label) {
            this.role_label = role_label;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {import('@models/Users').RoleTaxonomyParams}
         */
        do = async () => {
            const request_url = `${GetRoleTaxonomyRequest.endpoint}?role_label=${this.role_label}`;

            const response = await fetch(request_url);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    /**
     * Use this request when you want to know what grants will a newly created role inherit. Get all role taxonomies that are directly below a given role hierarchy.
     * for example, assume the system has roles with the following hierarchy: [0, 2, 3, 7, 8, 8 , 10]. If 4 is passed then it will return 2 taxonomies, with hierarchies 8. 
     * if 9 is passed then it will return 1 taxonomy with hierarchy 10. you can use the grants in these taxonomies to know what grants a role will inherit.
     * as you could've guessed, this request requires the 'grant_option' grant.
     */
    export class GetRoleTaxonomiesBelowHierarchyRequest {

        static endpoint = `${users_server}/roles/below-hierarchy`;

        /**
         * @param {number} hierarchy
         */
        constructor(hierarchy) {
            this.hierarchy = hierarchy;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<import('@models/Users').RoleTaxonomyParams[]>}
         */
        do = async () => {
            const request_url = `${GetRoleTaxonomiesBelowHierarchyRequest.endpoint}?hierarchy=${this.hierarchy}`;

            const response = await fetch(request_url);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    /**
     * Creates a new role from a given taxonomy that is not already in the system(otherwise it will fail). requires the 'grant_option' grant.
     */
    export class PostCreateRoleRequest {
        static endpoint = `${users_server}/roles/role`;

        /**
         * 
         * @param {import('@models/Users').RoleTaxonomy} role_taxonomy 
         */
        constructor(role_taxonomy) {
            this.role_label = role_taxonomy.RoleLabel;
            this.role_hierarchy = role_taxonomy.RoleHierarchy;
            this.role_grants = role_taxonomy.RoleGrants;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const response = await fetch(PostCreateRoleRequest.endpoint, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            let created = false;

            if (response.status === 201) {
                created = true;
            }

            return new HttpResponse(response, created);
        }
    }

    /**
     * Adds a user to role. requires the 'grant_option' grant.
     */    
    export class PatchUserRolesRequest {
        static endpoint = `${users_server}/users/role`;

        /**
         * @param {string} username
         * @param {string} role_label
         */
        constructor(username, role_label) {
            this.username = username;
            this.role_label = role_label;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = `${PatchUserRolesRequest.endpoint}?username=${this.username}&role_label=${this.role_label}`;

            const response = await fetch(request_url, {
                method: "PATCH"
            });

            let added = false;

            if (response.status === 200) {
                added = true;
            }

            return new HttpResponse(response, added);
        }
    }

    /**
     * Deletes the link between a user and a role. requires the 'grant_option' grant.
     */
    export class DeleteUserFromRoleRequest {
        static endpoint = `${users_server}/users/role`;

        /**
         * @param {string} username
         * @param {string} role_label
         */
        constructor(username, role_label) {
            this.username = username;
            this.role_label = role_label;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = `${DeleteUserFromRoleRequest.endpoint}?username=${this.username}&role_label=${this.role_label}`;

            const response = await fetch(request_url, {
                method: "DELETE"
            });

            let deleted = false;

            if (response.status === 200) {
                deleted = true;
            }

            return new HttpResponse(response, deleted);
        }
    }

    /**
     * Deletes a grant from a role. requires the 'grant_option' grant.
     */
    export class DeleteGrantFromRoleRequest {
        static endpoint = `${users_server}/roles/remove-grant`;

        /**
         * @param {string} role_label
         * @param {string} grant
         */
        constructor(role_label, grant) {
            this.role_label = role_label;
            this.grant = grant;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = `${DeleteGrantFromRoleRequest.endpoint}?role_label=${this.role_label}&grant=${this.grant}`;

            const response = await fetch(request_url, {
                method: "DELETE"
            });

            let deleted = false;

            if (response.status === 204) {
                deleted = true;
            }

            return new HttpResponse(response, deleted);
        }
    }

    /**
     * Gets all the role_labels of roles that a given username is part of. requires the 'grant_option' grant.
     */
    export class GetUserRolesRequest {
        static endpoint = `${users_server}/users/roles`;

        /**
         * @param {string} username
         */
        constructor(username) {
            this.username = username;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<string[]>}
         */
        do = async () => {
            const request_url = `${GetUserRolesRequest.endpoint}?username=${this.username}`;

            const response = await fetch(request_url);

            let data = null;

            if (response.ok) {
                data = await response.json();
            }

            return new HttpResponse(response, data);
        }
    }

    /**
     * Deletes a user account. requires 'delete_users' grant.
     */
    export class DeleteUserRequest {
        static endpoint = `${users_server}/users/user`;

        /**
         * @param {string} user_uuid
         */
        constructor(user_uuid) {
            this.user_uuid = user_uuid;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = `${DeleteUserRequest.endpoint}?user_uuid=${this.user_uuid}`;

            const response = await fetch(request_url, {
                method: "DELETE"
            });

            let deleted = false;

            if (response.status === 200) {
                deleted = true;
            }

            return new HttpResponse(response, deleted);
        }
    }

    /**
     * Deletes a grant from the system. requires 'grant_option' grant. Deleting a grant does not mean that actions that require that grant will stop requiring it.
     * but it will ensure that no user(not including the super admin and admins as they have the ALL_PRIVILEGES grant which includes every grant except for grant_option, which 
     * only the super admin has) will have that grant as it will be removed from all roles. Also until is readded, any attempts to link it to a role will fail.
     */ 
    export class DeleteGrantRequest {
        static endpoint = `${users_server}/roles/grant`;

        /**
         * @param {string} grant
         */
        constructor(grant) {
            this.grant = grant;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = `${DeleteGrantRequest.endpoint}?grant=${this.grant}`;

            const response = await fetch(request_url, {
                method: "DELETE"
            });

            let deleted = false;

            if (response.status === 204) {
                deleted = true;
            }

            return new HttpResponse(response, deleted);
        }
    }

    /**
     * Deletes a role from the system. all users will be unassociated from the role. a role with hierarchy 0 cannot be deleted. requires 'grant_option' grant.
     */
    export class DeleteRoleRequest {
        static endpoint = `${users_server}/roles/role`;

        /**
         * @param {string} role_label
         */
        constructor(role_label) {
            this.role_label = role_label;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = `${DeleteRoleRequest.endpoint}?role_label=${this.role_label}`;

            const response = await fetch(request_url, {
                method: "DELETE"
            });

            let deleted = false;

            if (response.status === 204) {
                deleted = true;
            }

            return new HttpResponse(response, deleted);
        }
    }

    /**
     * Changes a user account's password. requires the 'modify_users' grant. 
     */
    export class PutChangeUserPasswordRequest {
        static endpoint = `${users_server}/users/user/secret`;

        /**
         * @param {string} uuid
         * @param {string} username
         * @param {string} secret_hash
         */
        constructor(uuid, username, secret_hash) {
            this.uuid = uuid;
            this.username = username;
            this.secret_hash = secret_hash;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = PutChangeUserPasswordRequest.endpoint;

            const response = await fetch(request_url, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            let changed = false;

            if (response.status === 204) {
                changed = true;
            }

            return new HttpResponse(response, changed);
        }
    }

    /**
     * Changes a user account's username. requires the 'modify_users' grant.
     */
    export class PutChangeUsernameRequest {
        static endpoint = `${users_server}/users/user/username`;

        /**
         * @param {string} uuid
         * @param {string} username
         */
        constructor(uuid, username) {
            this.uuid = uuid;
            this.username = username;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>}
         */
        do = async () => {
            const request_url = PutChangeUsernameRequest.endpoint;

            const response = await fetch(request_url, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            let changed = false;

            if (response.status === 204) {
                changed = true;
            }

            return new HttpResponse(response, changed);
        }
    }

/*=====  End of Users  ======*/