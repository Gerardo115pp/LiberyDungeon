import { HttpResponse, attributesToJson } from "../base";
import { categories_server } from "../services";

/**
 * Gets a leaf of the categories tree using the leaf id which is provided by the parents inner categories attribute. if you want the root of the tree, pass 'main' as the leaf id.
 * @date 10/20/2023 - 10:38:16 PM
 *
 * @export
 * @class GetCategoryDataRequest
 */
export class GetCategoryTreeLeafRequest {
    
    /**
     * @param {string} category_id
     * @param {string} [cluster_id]
     */
    constructor(category_id, cluster_id) {
        this.category_id = category_id;
        this.cluster_id = cluster_id;   
    }

    toJson = attributesToJson.bind(this);

    
    /**
     * Sends the http request to the server and returns the response.
     * @returns {Promise<HttpResponse<import("@models/Categories").CategoryLeafParams>>}
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
    
    /**
     * @param {string} category_id 
     * @param {string} cluster_id 
     */
    constructor(category_id, cluster_id) {
        this.category_id = category_id;
        this.cluster_id = cluster_id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/Categories").InnerCategoryParams[]>>}
     */
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
    
    /**
     * @param {string} category_id 
     */
    constructor(category_id) {
        this.category_id = category_id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/Categories").CategoryParams>>}
     */
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

    /**
     * @param {string} category_path
     * @param {string} category_cluster
     */
    constructor(category_path, category_cluster) {
        this.category_path = encodeURI(category_path);
        this.category_cluster = category_cluster;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/Categories").CategoryParams>>}
     */
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

    /**
     * @param {string} category_id 
     * @param {import('@models/Medias').Media[]} rejected_medias 
     * @param {Object<string, import('@models/Medias').Media[]>} moved_medias 
     */
    constructor(category_id, rejected_medias, moved_medias) {
        this._category_id = category_id;
        this.rejected_medias = rejected_medias;
        this.moved_medias = moved_medias;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<Object>>}
     */
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

    /**
     * @param {string} category_id 
     * @param {string} new_name 
     */
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

    /**
     * @param {string} new_parent_category 
     * @param {string} moved_category 
     */
    constructor(new_parent_category, moved_category) {
        this.new_parent_category = new_parent_category;
        this.moved_category = moved_category;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('@models/Categories').CategoryParams>>}
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

    /**
     * @param {string} query 
     * @param {string} cluster_id 
     * @param {string} ignore 
     */
    constructor(query, cluster_id, ignore) {
        this.query = query;
        this.cluster_id = cluster_id;
        this.ignore = ignore;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/Categories").CategoryParams[] | null>>}
     */
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
    
    /**
     * @param {string} category_id 
     * @param {boolean} force 
     */
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

    /**
     * @param {string} category_name 
     * @param {string} parent_category_id 
     */
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

    /**
     * @param {string} category_name 
     * @param {string} parent_category_id 
     * @param {string} parent_path 
     * @param {string} cluster 
     */
    constructor(category_name, parent_category_id, parent_path, cluster) {
        this.name = category_name;
        this.parent = parent_category_id;
        this.parent_path = parent_path;
        this.cluster = cluster;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/Categories").CategoryParams>>}
     */
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