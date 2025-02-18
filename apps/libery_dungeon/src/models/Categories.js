import { 
    GetCategoryTreeLeafRequest, 
    getShortCategoryTreeRequest, 
    DeleteCategoryRequest, 
    PostCreateCategoryRequest, 
    GetCategoryRequest, 
    GetCategoryByFullpathRequest,
    PatchRenameCategoryRequest,
    PatchMoveCategoryRequest,
    PatchCategoryThumbnailRequest
} from "@libs/DungeonsCommunication/services_requests/categories_requests";
import { Media, getMediaIdentityByUUID, MediaIdentity } from "./Medias";
import { 
    getAllCategoryIndexes, 
    addCategoryIndex, 
    getCategoryIndex, 
    updateCategoryIndex 
} from "@databases/category_cache";
import {
    GetCategoryConfigRequest, 
    PatchCategoryBillboardTagsRequest,
    PatchCategoryBillboardMediasRequest
} from "@libs/DungeonsCommunication/services_requests/metadata_requests/category_metadata_requests";
import { browser } from "$app/environment";
import { getRandomMediaUrl } from "@libs/DungeonsCommunication/services_requests/media_requests";

const ROOT_CATEGORY_PROXY_UUID = "main";

/**
 * Category data structure params.
* @typedef {Object} CategoryParams
 * @property {string} uuid the 40 character identifier of the category
 * @property {string} name the name of the category
 * @property {string} fullpath the full path of the category in the virtual file system
 * @property {string} parent the 40 character identifier of the parent category
 * @property {string} cluster the 40 character identifier of the cluster
 * @property {string} category_thumbnail the thumbnail of the category. if the category has no thumbnail, it will be set to an empty string otherwise it will be a media uuid.
*/

export class Category {
    /** 
     * @type {string} the 40 character identifier of the category 
     */
    #uuid;

    /** 
     * @type {string} the name of the category 
     */
    #name;

    /** 
     * @type {string} the 40 character identifier of the parent category 
     */
    #parent;

    /** 
     * @type {string} the full path of the category in the virtual file system 
     */
    #fullpath;

    /** 
     * @type {string} the id of the cluster the category belongs to 
     */
    #cluster;

    /**
     * the thumbnail of the category. if the category has no thumbnail, it will be set to an empty string
     * @type {string}
     */
    #category_thumbnail;


    /**
     * @param {CategoryParams} param0 
     */
    constructor({uuid, name, parent, fullpath, cluster, category_thumbnail}) {
        this.#uuid = uuid;
        this.#name = name;
        this.#parent = parent;
        this.#fullpath = fullpath;
        this.#cluster = cluster;
        this.#category_thumbnail = category_thumbnail;
    }

    get FullPath() {
        return this.#fullpath;
    }

    get Name() {
        return this.#name;
    }

    get Parent() {
        return this.#parent;
    }

    get UUID() {
        return this.#uuid;
    }

    get Cluster() {
        return this.#cluster;
    }

    /**
     * Converts a category into an InnerCategory object
     * @returns {InnerCategory}
     */
    toInnerCategory = () => {
        return new InnerCategory({
            name: this.#name,
            uuid: this.#uuid,
            fullpath: this.#fullpath,
            category_thumbnail: this.#category_thumbnail
        });
    }

    /**
     * Checks if a given category leaf is an ancestor of this category. the
     * check is based on the category full path.
     * @param {CategoryLeaf} category_leaf 
     * @returns {boolean}
     */
    leafIsAncestor = category_leaf => {
        if (!(category_leaf instanceof CategoryLeaf)) return false;
        console.debug(`Checking if '${category_leaf.FullPath}' is an ancestor of '${this.#fullpath}': ${this.#fullpath.startsWith(category_leaf.FullPath)}`);

        return this.#fullpath.startsWith(category_leaf.FullPath);
    }

    /**
     * Checks if a given category leaf is a child of this category. the check is based on the categorys full path.
     * @param {CategoryLeaf} category_leaf
     */
    leafIsChild = category_leaf => {
        if (!(category_leaf instanceof CategoryLeaf)) return false;

        return category_leaf.FullPath.startsWith(this.#fullpath);
    }
}

/**
 * Requests the server for the category data of the provided category id
 * @param {string} category_id
 * @returns {Promise<Category | null>}
 */
export const getCategory = async category_id => {
    let category = null;
    let category_request = new GetCategoryRequest(category_id);

    let response = await category_request.do();


    if (response.status >= 200 && response.status < 300 && response.data !== null) {
        category = new Category(response.data);
    }

    return category;    
}

/**
 * Requests the server for the category data of the provided category fullpath and cluster id.
 * @param {string} category_fullpath    
 * @param {string} cluster_id
 * @returns {Promise<Category | null>}
 */
export const getCategoryByPath = async (category_fullpath, cluster_id) => {
    if ((category_fullpath == null || category_fullpath === "") || (cluster_id == null || cluster_id === "")) {
        throw new Error("The category fullpath and cluster id must be set");
    }
    
    let category = null;
    let category_request = new GetCategoryByFullpathRequest(category_fullpath, cluster_id);

    let response = await category_request.do();

    if (response.Ok) {
        category = new Category(response.data);
    }

    return category;
}

/**
 * Renames a category and returns true if the category was renamed successfully.
 * @param {string} category_id the 40 character identifier of the category
 * @param {string} new_name the new name of the category
 * @returns {Promise<boolean>}
 * @async
 */
export const renameCategory = async (category_id, new_name) => {
    let renamed_successfully = false;

    let category_rename_request = new PatchRenameCategoryRequest(category_id, new_name);

    let response = await category_rename_request.do();

    if (response.Ok && response.data?.uuid === category_id) {
        renamed_successfully = response.data != null;
    }

    return renamed_successfully;
}

/** 
 * Moves a category and all of it's subcategories and children to a new parent category.
 * @param {string} moved_category the 40 character identifier of the category to move
 * @param {string} new_parent_category the 40 character identifier of the new parent category
 * @returns {Promise<Category | null>} If request fails, returns null
 * @async
 */
export const moveCategory = async (moved_category, new_parent_category) => {
    let category_move_request = new PatchMoveCategoryRequest(new_parent_category, moved_category);

    let response = await category_move_request.do();

    let new_category = null;

    if (response.Ok && response.data?.uuid === moved_category) {
        new_category = new Category(response.data);
    }

    return new_category;
}

/**
 * Sets a new media uuid as the thumbnail of a category
 * @param {string} category_id the 40 character identifier of the category
 * @param {string} media_id the 40 character identifier of the media
 * @returns {Promise<boolean>} true if the thumbnail was set successfully
 */
export const changeCategoryThumbnail = async (category_id, media_id) => {
    let category_thumbnail_request = new PatchCategoryThumbnailRequest(category_id, media_id);

    let response = await category_thumbnail_request.do();

    return response.data;
}

/*=============================================
=            Category Tree            =
=============================================*/

    /**
    * @typedef {Object} CategoryWeakIdentityParams
     * @property {string} category_uuid
     * @property {string} category_path
     * @property {string} cluster_uuid
     * @property {string} cluster_path
    */

    /**
    * @typedef {Object} InnerCategoryParams
     * @property {string} name the name of the category
     * @property {string} uuid the 40 character identifier of the category
     * @property {string} fullpath the full path of the category in the virtual file system
     * @property {string} category_thumbnail the thumbnail of the category. if the category has no thumbnail, it will be set to an empty string
    */

    export class CategoryWeakIdentity {

        /**
         * @param {CategoryWeakIdentityParams} param0
         */
        constructor({category_uuid, category_path, cluster_uuid, cluster_path}) {
            this.category_uuid = category_uuid;
            this.category_path = category_path;
            this.cluster_uuid = cluster_uuid;
            this.cluster_path = cluster_path;
        }

        /**
         * Converts a CategoryWeakIdentity to a CategoryWeakIdentityParams object   
         * @returns {CategoryWeakIdentityParams}
         */
        toParams = () => {
            return {
                category_uuid: this.category_uuid,
                category_path: this.category_path,
                cluster_uuid: this.cluster_uuid,
                cluster_path: this.cluster_path
            };
        }
    }

    export class InnerCategory {
        /**
         * The value of an empty Fullpath
         * @type {string}
         */
        static EMPTY_FULLPATH = "NO_PATH";

        /**
         * the name of the category  
         * @type {string} 
         */
        name;

        /**
         * the 40 character identifier of the category
         * @type {string}
         */
        uuid;

        /**
         * Optional. Could no be provided, and if so. it will be set to InnerCategory.EMPTY_FULLPATH
         * @type {string} the full path of the category in the virtual file system
         */
        fullpath;

        /**
         * the thumbnail of the category. if the category has no thumbnail, it will be set to an empty string
         * @type {string}
         */
        #category_thumbnail;

        /**
         * The thumbnail media identity. To set it, call the setThumbnail method.
         * @type {import('@models/Medias').MediaIdentity | null}
         */
        #the_thumbnail;

        /**
         * The configuration of the category. has to be loaded before usage.
         * @type {CategoryConfig | null}
         */
        #config;

        /**
         * @param {InnerCategoryParams} param0 
         */
        constructor ({name, uuid, fullpath, category_thumbnail}) {
            if (name === undefined || uuid === undefined) {
                throw new Error("The name and uuid of the category must be set");
            }

            this.name = name;

            this.uuid = uuid;
            
            this.fullpath = fullpath ?? InnerCategory.EMPTY_FULLPATH;

            this.#category_thumbnail = category_thumbnail;

            this.#the_thumbnail = null;

            this.#config = null;
        }

        /**
         * The category configuration. It has to be loaded first by calling loadCategoryConfig.
         * @type {CategoryConfig | null}
         */
        get Config() {
            return this.#config;
        }

        /**
         * Returns a url for a random media inside this category. The actual random selection is performed by the server.
         * @param {string} cluster_id - required by the api.
         * @param {number} [cache_seconds] - If passed as a positive integer, the server will instruct the browser to cache the random media for the passed amount of seconds.
         * @returns {string}
         */
        getRandomMediaURL = (cluster_id, cache_seconds) => {
            return getRandomMediaUrl(this.uuid, cluster_id, cache_seconds);
        }

        /**
         * Whether the InnerCategory has a category_thumbnail set. THIS DOES NOT MEAN THE THUMBNAIL HAS BEEN LOADED.
         * you will need to check thumbnailIsLoaded to see if the thumbnail has been loaded and if not, call the setThumbnail method which also returns whether the thumbnail could be loaded.
         * @returns {boolean}
         */
        hasThumbnail = () => {
            return this.#category_thumbnail !== "";
        }

        /**
         * Whether the category configuration has been loaded.
         * @returns {this is {#config: CategoryConfig, Config: CategoryConfig}}
         */
        hasConfig = () => {
            return this.#config !== null;
        }

        /**
         * Handles type conversion from InnerCategory instances.
         * @param {string} hint
         * @returns {string}
         * @throws {TypeError}
         */
        [Symbol.toPrimitive](hint) {
            switch (hint) {
                case "string":
                    return this.name;
                default:
                    throw new TypeError(`Cannot convert InnerCategory to ${hint}`);
            }
        }

        /**
         * Sets the thumbnail of the category. Returns true if the InnerCategory now has a usable thumbnail.
         * @return {Promise<boolean>}
         */
        loadThumbnail = async () => {
            if (this.#category_thumbnail === "") {
                return false;
            }

            if (this.#the_thumbnail != null) return true;

            const media_identity = await getMediaIdentityByUUID(this.#category_thumbnail);

            if (media_identity === null) return false;

            this.#the_thumbnail = media_identity;

            return true;
        }

        /**
         * Loads the category configuration. returns whether the configuration is now loaded.
         * @returns {Promise<CategoryConfig>}
         */
        loadCategoryConfig = async () => {
            if (this.hasConfig()) {
                return this.Config;
            }

            const category_config = await getCategoryConfig(this.uuid);

            if (category_config === null) {
                throw new Error("In @models/Categories.InnerCategory.loadCategoryConfig: The category configuration couldn't be loaded");
            }

            this.#config = category_config;

            return category_config;
        }

        /**
         * Sets the category thumbnail only locally
         * @param {import('@models/Medias').MediaIdentity} new_thumbnail
         */
        setThumbnail = (new_thumbnail) => {
            if (!(new_thumbnail instanceof MediaIdentity)) {
                throw new Error(`In InnerCategory<${this.name}>.setThumbnail: Attempted to pass a non-MediaIdentity object as the thumbnail`);
            }

            this.#category_thumbnail = new_thumbnail.uuid;

            this.#the_thumbnail = new_thumbnail;
        }
            
        /**
         * Converts an InnerCategory to a Category Leaf. this requires fetch the category data 
         * from the server.
         * @returns {Promise<CategoryLeaf | null>}
         */
        toCategoryLeaf = async () => {
            return getCategoryLeaf(this.uuid);
        }

        /**
         * Whether the thumbnail of the category has been loaded and is usable.
         * @returns {boolean}
         */
        thumbnailIsLoaded = () => {
            return this.#the_thumbnail !== null;
        }

        /**
         * The category thumbnail media identity. If called when the thumbnail has not been loaded, it will panic.
         * @returns {import('@models/Medias').MediaIdentity}
         */
        get Thumbnail() {
            if (this.#the_thumbnail === null) {
                throw new Error(`In InnerCategory<${this.name}>: Attempted to access a null thumbnail object. Call the setThumbnail method first, and if it returns true, then you can access the thumbnail.`);
            }

            return this.#the_thumbnail;
        }

        /**
         * Returns the category thumbnail uuid, if the category has no thumbnail, then this property returns an empty string.
         * @type {string}
         */
        get ThumbnailUUID() {
            return this.#category_thumbnail;
        }
    }

    /**
    * @typedef {Object} CategoryLeafInnerCategoryParam
     * @property {string} name the name of the category
     * @property {string} uuid the 40 character identifier of the category
    */

    /**
    * @typedef {Object} CategoryLeafMediaItemParam
     * @property {string} name the name of the media resource 
     * @property {string} uuid the 40 character identifier of the media resource
     * @property {string} type the type of the media resource
     * @property {string} last_seen DEPRECATED. Do not use.
     * @property {string} main_category the 40 character identifier of the main category
     * @property {number} downloaded_from an identifier of the source of the media resource
    */

    /**
    * @typedef {Object} CategoryLeafParams
     * @property {string} uuid the 40 character identifier of the category
     * @property {string} name the name of the category
     * @property {string} parent the 40 character identifier of the parent category
     * @property {InnerCategoryParams[]} inner_categories JSON data of the inner categories of this category
     * @property {import('@models/Medias').MediaParams[]} content JSON data of the media resources of this category
     * @property {string} fullpath the full path of the category in the cluster virtual file system
     * @property {string} cluster the uuid4 identifier of the cluster
     * @property {string} category_thumbnail the thumbnail of the category. if the category has no thumbnail, it will be set to an empty string else it will be a media uuid.
    */

    export class CategoryLeaf {
        /** 
         * 
         * @type {string}
         */
        #fullpath;

        /**
         * the inner categories of this category 
         * @type {InnerCategory[]}  
         */
        #inner_categories;

        /** 
         * References to the inner categories of this category, if the category hasn't been loaded yet, the value will be null
         * @type {Object<string, CategoryLeaf | null>} 
         */
        #inner_categories_map;

        /**
         * The cluster the category belongs to
         * @type {string}
         */
        #cluster;

        /**
         * The category thumbnail. if the category has no thumbnail, it will be set to an empty string
         * @type {string}
         */
        #category_thumbnail;

        /**
         * The configuration of the category. This is not loaded by default.
         * @type {CategoryConfig | null}
         */
        #config;

        /**
         * @param {CategoryLeafParams} param0 
         */
        constructor({uuid, name, parent, inner_categories, content, fullpath, cluster, category_thumbnail}) {
            /** @type {string} the 40 character identifier of the category */
            this.uuid = uuid;

            /** @type {string} the name of the category */
            this.name = name;

            /** @type {string} the 40 character identifier of the parent category */
            this.parent = parent;

            /** @type {CategoryLeaf | null} reference to the parent category */
            this.parent_category = null;

            this.#inner_categories = inner_categories.map(inner_category => {
                return new InnerCategory(inner_category);
            });

            this.#fullpath = fullpath;

            /** @type {Media[]} the media resources of this category */
            this.content = content.map(media => {
                return new Media(media, this.#fullpath);
            });

            this.#cluster = cluster;

            this.#category_thumbnail = category_thumbnail;

            this.#inner_categories_map = {};

            this.#setInnerCategoriesMap();
            console.debug("Inner categories map: ", this.#inner_categories_map);

            this.#config = null;
        }

        /**
         * Adds a an inner category to the this category as child. doesn't fetch its data from the server
         * @param {InnerCategory} new_child
         * @returns {void}
         */
        addInnerCategory = new_child => {
            if (!(new_child instanceof InnerCategory)) return;

            this.#inner_categories.push(new_child);
            this.#inner_categories_map[new_child.uuid] = null;
        }

        /**
         * Generates an InnerCategory representation of this category.
         * @returns {InnerCategory}
         */
        asInnerCategory = () => {
            return new InnerCategory({
                name: this.name,
                uuid: this.uuid,
                fullpath: this.#fullpath,
                category_thumbnail: this.#category_thumbnail
            });
        }

        /**
         * Converts a CategoryLeaf to a Category.
         * @returns {Category}
         */
        asCategory = () => {
            return new Category({
                name: this.name,
                uuid: this.uuid,
                parent: this.parent,
                fullpath: this.#fullpath,
                cluster: this.#cluster,
                category_thumbnail: this.#category_thumbnail
            });
        }

        get ClusterUUID() {
            return this.#cluster;
        }

        /**
         * Returns the category configuration if it's loaded, otherwise returns null. To load the configuration, call the loadCategoryConfig once.
         * @type {CategoryConfig | null}
         */
        get Config() {
            return this.#config;
        }

        get FullPath() {
            return this.#fullpath;
        }
        /**
         * @returns {CategoryLeaf[]} the inner categories of this category that have been loaded
         */
        getLeafs = () => {
            const leafs = Object.values(this.#inner_categories_map).filter(category => {
                return category !== null;
            });

            return leafs;
        }

        /**
         * Returns the inner category data of the provided category id
         * @param {string} category_id
         * @returns {CategoryLeaf | null} returns null if the category has not been loaded yet. undefined if the category is not a child category of this category
         */
        getInnerCategoryLeafData = category_id => {
            let is_child = false;

            /**
             * @type {CategoryLeaf | null}
             */
            let leaf_data = null;

            this.#inner_categories.forEach(inner_category => {
                // if child is already true, than keep it as true. if not, is_child will be true if the inner category is the category we are looking for
                is_child = (inner_category.uuid === category_id) || is_child;
            });

            console.debug("is_child: ", is_child);

            if (is_child) {
                leaf_data = this.#inner_categories_map[category_id] ?? null;
            }

            return leaf_data;
        }

        /**
         * @returns {boolean} if the category has content
         */
        hasContent = () => {    
            return this.content.length > 0;
        }

        /**
         * @returns {boolean} if the category has inner categories
         */
        hasInnerCategories = () => {    
            return this.#inner_categories.length > 0;
        }

        /**
         * Returns true if the provided category is a parent of this category
         * @param {CategoryLeaf} category
         * @returns {boolean}
         */
        isParentCategory = category => {
            if (!(category instanceof CategoryLeaf)) return false;
            return this.parent === category.uuid && this.parent !== "";
        }

        /**
         * Whether a uuid is the parent of this category
         * @param {string} category_uuid
         */
        isParentCategoryUUID = category_uuid => {
            return this.parent === category_uuid && this.parent !== "";
        }

        /**
         * Returns true if the provided category is a child of this category
         * @param {CategoryLeaf} category
         * @returns {boolean}
         */
        isChildCategory = category => {
            if (!(category instanceof CategoryLeaf)) return false;
            
            return this.#inner_categories_map[category.uuid] !== undefined;
        }

        /**
         * Whether a uuid is a child of this category
         * @param {string} category_uuid
         */
        isChildCategoryUUID = category_uuid => {
            console.debug("Checking if ", category_uuid, " is a child of ", this.uuid, ": ", this.#inner_categories_map[category_uuid] !== undefined);
            console.debug("Inner categories map: ", this.#inner_categories_map);
            return this.#inner_categories_map[category_uuid] !== undefined;
        }

        /**
         * Returns true if the provided category path is a parent of this category
         * @param {string} category_path
         * @returns {boolean}
         */
        isParentCategoryPath = category_path => {
            return this.#fullpath.startsWith(category_path);
        }

        /**
         * Returns true if the provided category path is a child of this category
         * @param {string} category_path
         * @returns {boolean}
         */
        isChildCategoryPath = category_path => {
            return category_path.startsWith(this.#fullpath);
        }

        /**
         * Returns true if the category has no inner categories and no content
         * @returns {boolean}
        */
        isEmpty = () => {
            return this.#inner_categories.length === 0 && this.content.length === 0;
        }

        get InnerCategories() {
            return this.#inner_categories;
        }

        /**
         * Loads the category configuration. returns whether the configuration is now loaded.
         * @returns {Promise<CategoryConfig | null>}
         */
        loadCategoryConfig = async () => {
            const category_config = await getCategoryConfig(this.uuid);

            if (category_config === null) return null;

            this.setCategoryConfig(category_config);

            return category_config;
        }

        /**
         * @returns {CategoryLeaf | null} the parent category of this category
         */
        get ParentCategory() {
            return this.parent_category;
        }

        /**
         * Sets the parent category of this category.
         * @param {CategoryLeaf} category
         */
        set ParentCategory(category) {
            if (!(category instanceof CategoryLeaf)) return;

            this.parent_category = category;
        }

        /**
         * Removes a loaded leaf from the inner categories map
         * @param {string} category_id
         */
        removeLoadedLeaf = category_id => {
            delete this.#inner_categories_map[category_id];
            this.#inner_categories_map[category_id] = null;
        }

        /**
         * Adds a child category data to this category
         * @param {CategoryLeaf} category
         * @returns {boolean} true if the category was added, false otherwise
         */
        setChildCategory = category => {
            if (!(category instanceof CategoryLeaf) || !this.isChildCategory(category)) return false;

            this.#inner_categories_map[category.uuid] = category;
            return true;
        }

        /**
         * Sets the category configuration
         * @param {CategoryConfig} config 
         */
        setCategoryConfig = config => {
            this.#config = config;
        }
        
        /**
         * Structures all the uuids of the inner categories of this category in a map stored in the #inner_categories_map property
         * @date 10/21/2023 - 12:01:05 AM
         */
        #setInnerCategoriesMap = () => {
            console.debug("Called #setInnerCategoriesMap");
            this.#inner_categories.forEach(inner_category => {
                this.#inner_categories_map[inner_category.uuid] = null;
            });
        }

        /**
         * Updates the category leaf data from another category leaf
         * @param {CategoryLeaf} other_category
         * @returns {boolean} true if the category was updated, false otherwise
         */
        updateCategory = other_category => {
            if (!(other_category instanceof CategoryLeaf) || other_category.uuid !== this.uuid) return false;

            this.name = other_category.name;
            this.parent = other_category.parent;
            this.#inner_categories = other_category.InnerCategories;
            this.content = other_category.content;
            this.#fullpath = other_category.FullPath;

            this.#inner_categories_map = {};
            this.#setInnerCategoriesMap();
            // Object.keys(this.#inner_categories_map).forEach(uuid => {
            //     if (!other_category.isChildCategoryUUID(uuid)) {
            //         delete this.#inner_categories_map[uuid];
            //     }
            // });

            return true;
        }
    }

    /**
     * A tree structure of the categories in the cluster. Manages navigation and management of the categories and their content.
     * @requires svelte/store
     */
    export class CategoriesTree {
        /**
         * @param {CategoryLeaf} root_category 
         * @param {import('svelte/store').Writable<CategoryLeaf | null>} current_category_store 
         */
        constructor(root_category, current_category_store) {
            /** @type {CategoryLeaf} the root category of the tree */
            this.root_category = root_category;
            /** @type {CategoryLeaf} the currently selected category */
            this.current_category = root_category;
            /** 
             * @type {import('svelte/store').Writable<CategoryLeaf | null>} A store used for reactivity
             */   
            this.current_category_store = current_category_store;

            this.setCurrentCategory(root_category);
        }


        /**
         * Deletes a category that is a direct child of the current category
         * @param {string} category_id the uuid of the category to delete
         * @param {boolean} force if true, the category will be deleted even if it has content
         * @returns {Promise<boolean>} true if the category was deleted, false otherwise
         * @async
         * @memberof CategoriesTree
         */
        deleteChildCategory = async (category_id, force=false) => {
            let is_child = false;

            this.current_category.InnerCategories.forEach(ic => {
                if (ic.uuid === category_id) {
                    is_child = true;
                }
            });

            if (!is_child) {
                throw new Error("The category is not a child of the current category");
            }

            let deleted = await deleteCategory(category_id, force);

            if (deleted) {
                await this.updateCurrentCategory();
            }

            return deleted;
        }

        /**
         * Changes the root category of the tree. this means loosing all the loaded and changing the current category to the new root category.
         * @param {string} new_root_category_uuid
         * @returns {Promise<void>}
         */
        changeRootCategory = async new_root_category_uuid => {
            let new_root_category = await getCategoryLeaf(new_root_category_uuid);
            if (new_root_category === null) {
                throw new Error("The category couldn't be fetched, likely because it doesn't exist");
            }

            this.root_category = new_root_category;
            this.current_category = new_root_category;
            this.current_category_store.set(new_root_category);
        }

        /**
         * looks for a category inside the loaded tree
         * @param {Category} target_category - the category to look for
         * @returns {CategoryLeaf | null}
         */
        getLoadedCategory = target_category => {
            if (target_category.UUID === this.root_category.uuid) {
                return this.root_category;
            }
            
            let category_known_ancestor = this.root_category;
            let loaded_category = null;
            let category_parent_found;

            console.debug("target_category: ", target_category);

            do {
                console.debug("Category known ancestor: ", category_known_ancestor);
                category_parent_found = false;
                let ancestor_leafs = category_known_ancestor.getLeafs();
                console.debug("Ancestor leafs: ", ancestor_leafs);

                for (let leaf of ancestor_leafs) {
                    if (leaf.uuid === target_category.UUID) {
                        loaded_category = leaf;
                        break;
                    }

                    console.debug(`Checking if '${target_category.FullPath}' starts with '${leaf.FullPath}': ${target_category.FullPath.startsWith(leaf.FullPath)}`);
                    if (target_category.FullPath.startsWith(leaf.FullPath)) {
                        console.debug("Found");
                        category_known_ancestor = leaf;
                        category_parent_found = true;
                        break;
                    }
                }
            } while (loaded_category === null && category_parent_found);

            return loaded_category;
        }
       
        /**
         * Creates a new category as a child of the current category.
         * @param {string} category_name the name of the new category
         * @param {string} cluster_id the uuid of the cluster the category belongs to
         * @returns {Promise<Category | null>} the created category
         */
        insertChildCategory = async (category_name, cluster_id) => {
            if (category_name === "" || cluster_id === "") {
                throw new Error("The category name and cluster id must be set");
            }
            
            let created = null;

            created = await createCategory(this.current_category.uuid, this.current_category.FullPath, category_name, cluster_id);
            
            if (created !== null) {
                await this.updateCurrentCategory();
            }

            return created;
        }

        /**
         * Description placeholder
         * @date 10/21/2023 - 12:01:05 AM
         * @param {CategoryLeaf} category_leaf
         */
        setCurrentCategory = (category_leaf) => {
            this.current_category = category_leaf;
            this.current_category_store.set(category_leaf);
        }

        /**
         * Navigates to a category that is a child of the current category. Returns an error if 
         * the category is not a child of the current category.
         * @param {string} category_id the uuid of the category to navigate to
         * @returns {Promise<Error | null>} an error if the category is not a child of the current category
         */
        navigateToLeaf = async category_id => {
            let err = null;

            let category = this.current_category.getInnerCategoryLeafData(category_id);

            console.debug("Inner category data: \n", category);

            if (category === null) {
                category = await getCategoryLeaf(category_id);

                if (category === null) {
                    err = new Error("The category couldn't be fetched, likely because it doesn't exist");
                    return err;
                }
                
                let added = this.current_category.setChildCategory(category);
        
                if (!added) {
                    console.warn("The category couldn't be added to the current category, likely because it's not a child category of the current category");
                }
            } else if (category === undefined) {
                err = new Error("The category is not a child of the current category");
                return err;
            }

            this.setCurrentCategory(category);

            return err;
        }

        /**
         * Navigates to the parent category of the current category. If the parent category is not fetched
         * but the parent id is not ""(which would mean this is the root category of the current cluster),
         * it will fetch the parent category from the server.
         * @returns {Promise<Error | null>} an error if the parent category is not fetched and the parent id is not ""
         */
        navigateToParent = async () => {
            let err = null;

            if (this.current_category.parent === "") {
                return new Error("This is the root category of the cluster");
            }

            let parent_category = this.current_category.parent_category;

            if (parent_category === null) {
                parent_category = await getCategoryLeaf(this.current_category.parent);

                if (parent_category === null) {
                    return new Error("The parent category couldn't be fetched, likely because it doesn't exist");
                }
                
                parent_category.setChildCategory(this.current_category);

                if (this.root_category.uuid === this.current_category.uuid) {
                    this.root_category = parent_category;
                }

                this.current_category.ParentCategory = parent_category;
            }

            this.setCurrentCategory(parent_category);

            return err;
        }

        /**
         * updates the contents of a loaded category, if the category is not loaded, it does nothing. if it manages to update the category, returns true, otherwise false.
         * does not make any updates to the current category. if you want to update the current category, use the updateCurrentCategory method instead.
         * @param {string} category_id 
         * @returns {Promise<boolean>}
         */
        updateCategory = async category_id => {
            let category_updated = false;

            let category_data = await getCategory(category_id);
            if (category_data === null) {
                console.error("Category data was null");
                return false;
            }

            let loaded_category = this.getLoadedCategory(category_data);
            if (loaded_category === null) {
                console.error("Loaded category was null");
                return false;
            }

            console.log("Loaded category: ", loaded_category);

            let new_category_leaf = await getCategoryLeaf(category_id);
            if (new_category_leaf === null) {
                console.warn("The category couldn't be updated, likely because it doesn't exist");
                return category_updated;
            }

            category_updated = loaded_category.updateCategory(new_category_leaf);
            
            return category_updated;
        }

        /**
         * Fetches a fresh copy of the current category from the server. Then it updates the local
         * copy from the new data.
         * @returns {Promise<Error | null>} an error if the fetched category was null.
         */
        updateCurrentCategory = async () => {
            let category = await getCategoryLeaf(this.current_category.uuid);

            if (category === null) {
                return new Error("The category couldn't be updated, likely because it doesn't exist");
            }

            let updated = this.current_category.updateCategory(category);

            if (!updated) {
                console.warn("The category couldn't be updated, likely because both categories have no differences");
                return null;
            }

            this.setCurrentCategory(this.current_category);

            return null;
        }

        /**
         * Looks for a loaded category in the tree and if exists, deletes it from the tree. this ensures that it's contents will be up to date when the user navigates to it again.
         * @param {Category} target_category 
         */
        deleteLoadedCategoryCache = target_category => {
            let target_leaf = this.getLoadedCategory(target_category);

            if (target_leaf === null) return;

            let parent_leaf = target_leaf.parent_category;

            if (parent_leaf === null) return;

            parent_leaf.removeLoadedLeaf(target_leaf.uuid);
        }
    }

    /**
     * Sends a category creation request to the server
     * @param {string} parent_category_uuid the uuid of the parent category
     * @param {string} parent_category_fullpath the full path of the parent category
     * @param {string} cluster_id the uuid of the cluster
     * @param {string} category_name the name of the new category
     * @returns {Promise<Category | null>} 
     */
    export const createCategory = async (parent_category_uuid, parent_category_fullpath, category_name, cluster_id) => {
        const request = new PostCreateCategoryRequest(category_name, parent_category_uuid, parent_category_fullpath, cluster_id);
        const response = await request.do();

        if (!response.Created) {
            return null;
        }

        return new Category(response.data);
    }

    /**
     * Description placeholder
     * @date 10/20/2023 - 11:06:22 PM
     * @param {string} category_id
     * @param {import('svelte/store').Writable<CategoryLeaf | null>} category_store
     * @param {string} [cluster_id] the uuid of the cluster. If category_id is set and not equal to the root category proxy id, cluster_id can be omitted. otherwise, cluster_id must be set.
     * @export
     * @returns {Promise<CategoriesTree | null>}
     */
    export const getCategoryTree = async (category_id, category_store, cluster_id) => {
        if ((category_id == null || category_id === ROOT_CATEGORY_PROXY_UUID) && cluster_id === undefined) {
            throw new Error("Category id or cluster id must be set. If you are still using only the main category proxy id. you must set the cluster id so the server can know from which cluster to get the root category from.");
        }

        /**
         * @type {CategoriesTree | null}
         */
        let tree = null;

        category_id = category_id || ROOT_CATEGORY_PROXY_UUID;

        const request = new GetCategoryTreeLeafRequest(category_id, cluster_id);
        const response = await request.do();

        if (response.status >= 200 && response.status < 300) {
            const data = response.data;
            const root_category = new CategoryLeaf(data);
            tree = new CategoriesTree(root_category, category_store);

        }

        return tree;
    }

    /**
     * Requests the server for the category data of the provided category id
     * @param {string} category_id
     * @returns {Promise<CategoryLeaf | null>}
     */
    const getCategoryLeaf = async category_id => {
        const request = new GetCategoryTreeLeafRequest(category_id);
        const response = await request.do();
        let category = null;

        if (response.Ok) {
            category = new CategoryLeaf(response.data);
        }


        return category;
    }

    /**
     * Requests the server for the category data of the provided category id. Returns a
     * list of inner categories.
     * @param {string} category_id
     * @param {string} cluster_id
     * @returns {Promise<InnerCategory[]>}
     */
    export const getShortCategoryTree = (category_id, cluster_id) => {
        if ((category_id == null || category_id === ROOT_CATEGORY_PROXY_UUID) && cluster_id === undefined) {
            throw new Error("Category id or cluster id must be set. If you are still using only the main category proxy id. you must set the cluster id so the server can know from which cluster to get the root category from.");
        }

        category_id = category_id || ROOT_CATEGORY_PROXY_UUID;

        return new Promise(async (resolve, reject) => {
            const request = new getShortCategoryTreeRequest(category_id, cluster_id);
            
            const response = await request.do();

            if (response.status >= 200 && response.status < 300) {
                const data = response.data;

                /** @type {InnerCategory[]} */
                const leafs = [];

                data.forEach(leaf => {
                    leafs.push(new InnerCategory(leaf));
                });
                

                resolve(leafs);
            } else {
                reject(response);
            }

        });
    }

    /**
     * @param {string} category_id 
     * @param {boolean} force 
     * @returns 
     */
    export const deleteCategory = async (category_id, force=false) => {
        return new Promise(async (resolve, reject) => {
            const request = new DeleteCategoryRequest(category_id, force);
            const response = await request.do();

            if (response.status === 205) {
                resolve(true);
            }

            resolve(false);
        });
    }


/*=====  End of Category Tree  ======*/

/*=============================================
=            Category Cache            =
=============================================*/

    class CategoryCache {
        /** 
         * @type {Map<string, number>} a map of category uuids to the index of last viewed media on that category 
         */
        #cache

        constructor() {
            this.#cache = new Map();

            this.#loadCache();
        }

        /**
         * Adds a category index to the cache
         * @param {string} category_uuid the uuid of the category
         * @param {number} media_index the media index of the category
         * @async
         */
        addCategoryIndex(category_uuid, media_index) {
            if (this.#cache.has(category_uuid)) {
                updateCategoryIndex(category_uuid, media_index);
            } else {
                addCategoryIndex(category_uuid, media_index);
            }
            
            this.#cache.set(category_uuid, media_index);
        }

        /**
         * Gets the media index of a category from the cache, if the category is not in the cache, it returns 0, which is useful to set the media index at the beginning of the category
         * @param {string} category_uuid the uuid of the category
         * @returns {Promise<number | undefined>} the media index of the category
         * @async   
         */
        async getCategoryIndex(category_uuid) {
            if (this.#cache.has(category_uuid)) {
                return this.#cache.get(category_uuid);
            } else {
                let result = await getCategoryIndex(category_uuid); // in theorie, all the categories should be in the local cache, but if the #loadCache method has not been resolved yet, calling this method will await for it to be resolved
                if (result === -1) return 0; // category not in cache

                this.#cache.set(category_uuid, result);
                return result;
            }
        }

        /**
         * Loads the cache from the database
         * @async
        */
        async #loadCache() {
            let cache = await getAllCategoryIndexes();
            
            this.#cache = cache;
        }

        /**
         * Sets the category index to 0
         * @param {string} category_uuid the uuid of the category
        */
        resetCategoryIndex(category_uuid) {
            if (this.#cache.has(category_uuid)) {
                this.#cache.set(category_uuid, 0);
                updateCategoryIndex(category_uuid, 0);
            }
        }
    }

    let tmp = null;
    if (browser) {
        tmp = new CategoryCache();
    }

    export const category_cache = tmp;

/*=====  End of Category Cache  ======*/


/*=============================================
=            Categories configuration            =
=============================================*/

/**
 * An object that holds the configuration of a category
* @typedef {Object} CategoryConfigParams
 * @property {string} category_uuid the 40 character identifier of the category
 * @property {string[]} billboard_media_uuids a list of media uuids to be displayed explorer billboard. takes precedence over the category content.
 * @property {number[]} billboard_dungeon_tags a list of dungeon tags, the medias matching these tags will be displayed in the explorer billboard. takes precedence over the billboard_media_uuids or the category content
*/

export class CategoryConfig {
    /** @type {string} the 40 character identifier of the category */
    #category_uuid;

    /** @type {string[]} a list of media uuids to be displayed explorer billboard */
    #billboard_media_uuids;

    /** @type {number[]} a list of dungeon tags, the medias matching these tags will be displayed in the explorer billboard */
    #billboard_dungeon_tags;

    /**
     * @param {CategoryConfigParams} param0
     */
    constructor({category_uuid, billboard_media_uuids, billboard_dungeon_tags}) {
        this.#category_uuid = category_uuid;
        this.#billboard_media_uuids = billboard_media_uuids;
        this.#billboard_dungeon_tags = billboard_dungeon_tags;
    }

    /**
     * The 40 character identifier of the category this configuration is for.
     * @type {string}
     */
    get CategoryUUID() {
        return this.#category_uuid;
    }

    /**
     * The list of media uuids to be displayed in the explorer billboard. Takes precedence over the category content.
     * @type {string[]}
     */
    get BillboardMediaUUIDs() {
        return this.#billboard_media_uuids;
    }

    /**
     * The list of dungeon tag ids. The medias matching these tags will be displayed in the explorer billboard. Takes precedence over the billboard_media_uuids or the category content.
     * @type {number[]}
     */
    get BillboardDungeonTags() {
        return this.#billboard_dungeon_tags;
    }

    /**
     * Sets a new list of media uuids to be displayed in the explorer billboard and updates the configuration on the server.
     * @param {string[]} media_uuids 
     * @returns {Promise<boolean>}
     */
    updateBillboardMediaUUIDs = async media_uuids => {
        let success = await updateCategoryConfig__BillboardMediaUUIDs(this.#category_uuid, media_uuids);

        if (success) {
            this.#billboard_media_uuids = media_uuids;
        }

        return success;
    }

    /**
     * Sets a new list of dungeon tags for the category and updates the configuration on the server.
     * @param {number[]} dungeon_tags
     * @returns {Promise<boolean>}
     */
    updateBillboardDungeonTags = async dungeon_tags => {
        let success = await updateCategoryConfig__BillboardDungeonTags(this.#category_uuid, dungeon_tags);

        if (success) {
            this.#billboard_dungeon_tags = dungeon_tags;
        }

        return success;
    }

    /**
     * Converts a CategoryConfig to a CategoryConfigParams object
     * @returns {CategoryConfigParams}
     */
    toParams = () => {
        return {
            category_uuid: this.#category_uuid,
            billboard_media_uuids: this.#billboard_media_uuids,
            billboard_dungeon_tags: this.#billboard_dungeon_tags
        };
    }
}

/**
 * Returns a CategoryConfig object from the server. Requires the user to have public_content_read grant.
 * @param {string} category_uuid the 40 character identifier of the category
 * @returns {Promise<CategoryConfig | null>}
 */
export const getCategoryConfig = async category_uuid => {
    let request = new GetCategoryConfigRequest(category_uuid);

    /** 
     * @type {CategoryConfig | null}
     */
    let category_config = null;

    let response = await request.do();

    if (response.Ok && response.data !== null) {
        category_config = new CategoryConfig(response.data);
    }

    return category_config;
}

/**
 * Updates the list of media uuids to be displayed in the explorer billboard for a category. Requires the user to have content_alter grant.
 * @param {string} category_uuid the 40 character identifier of the category
 * @param {string[]} media_uuids a list of media uuids
 * @returns {Promise<boolean>}
 */
export const updateCategoryConfig__BillboardMediaUUIDs = async (category_uuid, media_uuids) => {
    let request = new PatchCategoryBillboardMediasRequest(category_uuid, media_uuids);

    let response = await request.do();

    let success = false;

    if (response.Ok) {
        success = response.data.response;
    }

    return success;
}

/**
 * Updates the list of dungeon tags for a category. Requires the user to have content_alter grant.
 * @param {string} category_uuid the 40 character identifier of the category
 * @param {number[]} dungeon_tags a list of dungeon tag ids
 * @returns {Promise<boolean>}
 */
export const updateCategoryConfig__BillboardDungeonTags = async (category_uuid, dungeon_tags) => {
    let request = new PatchCategoryBillboardTagsRequest(category_uuid, dungeon_tags);

    let response = await request.do();

    let success = false;

    if (response.Ok) {
        success = response.data.response;
    }

    return success;
}

/*=====  End of Categories configuration  ======*/





