import { 
    CommitCategoryTreeChangesRequest,
    GetCategorySearchResultsRequest, 
} from "@libs/DungeonsCommunication/services_requests/categories_requests";
import { PutResyncClusterBranchRequest } from "@libs/DungeonsCommunication/services_requests/categories_cluster_requests";
import { Media } from "./Medias";
import { Category, InnerCategory } from "./Categories";

/*=============================================
=            Media Changes            =
=============================================*/

    /**
    * @typedef {"Moved" | "Deleted" | "Normal" | null} MediaChangeType
    */

    /**
     * @type {Object<string, MediaChangeType>} the types of changes that can be made to a media
     */
    export const media_change_types = {
        MOVED: "Moved",
        DELETED: "Deleted",
        NORMAL: "Normal"
    }

    export class MediaChangesManager {
        /** @type {Map<string, Media>} a set of medias that would be deleted on commit, the server only receives the medias as an array, but not the uuids */
        #deleted_medias

        /** @type {Map<string, string>} a map of medias that would be moved on commit, the key is the media uuid and the value is the new category uuid. This is NOT sent to the server on commit */
        #moved_medias

        /**
         * A map of the categories uuids that medias will be moved to, the key is the category uuid and the value is an array of medias that will be moved to that category, this
         * is the information that is sent to the server, not the #moved_medias map.
         * @type {Object<string, Media[]>} 
         * 
         */
        #moved_medias_data

        /** 
         * @type {InnerCategory[]} a list of inner categories(only uuid and name) that at least one media is moved to 
         */
        #used_categories

        constructor() {
            this.#deleted_medias = new Map();
            this.#moved_medias = new Map();
            this.#moved_medias_data = {};
            this.#used_categories = [];
        }

        /**
         * The amount of changes that are staged, both deleted and moved medias.
         * @returns {number} the amount of changes
         */
        get ChangesAmount() {
            return this.#deleted_medias.size + this.#moved_medias.size;
        }

        /**
         * a set that would be deleted on commit.
         * @returns {Map<string, Media>}
         */
        get DeletedMediasMap() {
            return this.#deleted_medias;
        }

        /**
         * a list of all the medias that will be deleted on commit.
         * @returns {Media[]}
         */
        get DeletedMedias() {
            let deleted_medias = Array.from(this.#deleted_medias.values());

            return deleted_medias;
        }

        get MovedMediasMap() {
            return this.#moved_medias;
        }

        /**
         * a list of all the medias that will be moved on commit.
         * @returns {Media[]}
         */
        get MovedMedias() {
            /**
             * @type {Media[]}
             */
            let moved_medias = [];

            for (let cu of Object.keys(this.#moved_medias_data)) {
                moved_medias = moved_medias.concat(this.#moved_medias_data[cu]);
            }

            return moved_medias;
        }

        get MovedMediasData() {
            return this.#moved_medias_data;
        }

        get UsedCategories() {
            let sorted_used_categories = this.#used_categories.sort((a, b) => {
                let medias_on_a = this.#moved_medias_data[a.uuid].length;
                let medias_on_b = this.#moved_medias_data[b.uuid].length;
                
                if (medias_on_a > medias_on_b) {
                    return -1;
                } else if (medias_on_a < medias_on_b) {
                    return 1;
                } else {
                    return 0;
                }
            });

            return sorted_used_categories;
        }
        
        /**
         * Commits the changes made to the medias
         * @param {string} category_id the uuid of the category
         * @returns {Promise<import('@libs/DungeonsCommunication/base').HttpResponse<Object> | null>} the response of the request
         * @see CommitCategoryTreeChangesRequest
         * @see HttpResponse
         */
        async commitChanges(category_id) {
            if (this.#deleted_medias.size === 0 && this.MovedMediasMap.size === 0) {
                return null;
            }

            const request = new CommitCategoryTreeChangesRequest(category_id, Array.from(this.#deleted_medias.values()), this.#moved_medias_data);
            return await request.do();
        }

        /**
         * Clears the media of any changes
         * @param {string} media_uuid the uuid of the media
         */
        clearMediaChanges(media_uuid) {
            this.unstageMediaDeletion(media_uuid);
            this.unstageMediaMove(media_uuid);
        }

        /**
         * Clears all the deletion changes.
         * @returns {void}
         */
        clearAllDeletionChanges() {
            this.#deleted_medias.clear();
        }

        /**
         * Clears all the move changes.
         * @returns {void}
         */
        clearAllMoveChanges() {
            this.#moved_medias.clear();
            this.#moved_medias_data = {};
            this.#used_categories = [];
        }

        /**
         * Clears all the changes.
         * @returns {void}
         */
        clearAllChanges() {
            this.clearAllDeletionChanges();
            this.clearAllMoveChanges();
        }

        /**
         * Gets the MediaChangeType of a media by its uuid
         * @param {string} media_uuid the uuid of the media
         * @returns {MediaChangeType} the change type of the media
        */
        getMediaChangeType = (media_uuid) => {
            if (media_uuid === undefined || media_uuid === null) return null;

            if (this.#deleted_medias.has(media_uuid)) {
                return media_change_types.DELETED;
            } else if (this.#moved_medias.has(media_uuid)) {
                return media_change_types.MOVED;
            } else {
                return media_change_types.NORMAL;
            }
        }

        /**
         * Get the new category of a media by its uuid
         * @param {string} media_uuid the uuid of the media
         * @returns {InnerCategory | null} the new category of the media
         */
        getMediaNewCategory = media_uuid => {
            if (!this.#moved_medias.has(media_uuid)) {
                return null;
            }

            const new_category_uuid = this.#moved_medias.get(media_uuid);

            const new_category = this.#used_categories.find(category => category.uuid === new_category_uuid);

            return new_category ?? null;
        }

        /**
         * Adds a media to the deleted medias set
         * @param {Media} media the media to be deleted
        */
        stageMediaDeletion(media) {
                // check if media is staged to be moved to another category, if so, unstage it from the move staging
                if (this.#moved_medias.has(media.uuid)) {
                    this.unstageMediaMove(media.uuid);
                }

                this.#deleted_medias.set(media.uuid, media);
        }

        /**
         * Adds a media to the moved medias map
         * @param {Media} media the media to be moved 
         * @param {InnerCategory} new_category the uuid of the new category
        */
        stageMediaMove(media, new_category) {
                // check if media is already staged for move to another category
                if (this.#moved_medias.has(media.uuid)) {
                    this.unstageMediaMove(media.uuid);
                }

                // check if media is staged for deletion, if so, remove it from the deleted medias set
                if (this.#deleted_medias.has(media.uuid)) {
                    this.unstageMediaDeletion(media.uuid);
                }


                this.#moved_medias.set(media.uuid, new_category.uuid);
                
                if (!this.#moved_medias_data[new_category.uuid]) {
                    this.#moved_medias_data[new_category.uuid] = [];
                }

                this.#moved_medias_data[new_category.uuid].push(media);
                if ((this.#used_categories.find(category => category.uuid === new_category.uuid)) === undefined) {
                    this.#used_categories.push(new_category);
                }
        }

        /**
         * Removes a media from the delete staging set
         * @param {string} media_uuid the uuid of the media
         */
        unstageMediaDeletion(media_uuid){
            this.#deleted_medias.delete(media_uuid);
        }

        /**
         * Removes a media from the move staging map and the move staging data
         * @param {string} media_uuid the uuid of the media
        */
        unstageMediaMove(media_uuid) {
            if (!this.#moved_medias.has(media_uuid)) {
                return;
            }

            const new_category_uuid = this.#moved_medias.get(media_uuid);

            if (new_category_uuid === undefined) return;

            this.#moved_medias.delete(media_uuid);

            this.#moved_medias_data[new_category_uuid] = this.#moved_medias_data[new_category_uuid].filter(media => media.uuid !== media_uuid);

            if (this.#moved_medias_data[new_category_uuid].length === 0) {
                delete this.#moved_medias_data[new_category_uuid];

                this.#used_categories = this.#used_categories.filter(category => category.uuid !== new_category_uuid);
            }
        }
    }

    export class MediaChangesEmitter extends MediaChangesManager {
        /**
        * @callback onChangesCallback
         * @param {string} change_type
         * @param {string} media_uuid
         * @returns {void}
        */

        /**
        * @callback beforeCommitCallback
         * @this {MediaChangesEmitter}
         * @returns {void | Promise<void>}
        */
       
        /**
         * A map of string ids associated with a callback function. The callback function is called when a change is made to a media.
         * @type {Object<MediaChangeType, onChangesCallback>} 
         */
        #on_changes_callbacks

        /**
         * A callback that is called before the changes are committed.
         * @type {beforeCommitCallback}
         */
        #before_commit_callback

        /**
         * Whether an error on the before commit callback should stop the commit process.
         * @type {boolean}
         * @default true
         */
        #panic_on_before_callback_error

        constructor() {
            super();
            this.#on_changes_callbacks = {};
            this.#before_commit_callback = () => {};
            this.#panic_on_before_callback_error = true;
        }


        /**
         * Sets a callback to be called before the changes are committed. further calls to this method will overwrite the previous callback.
         * @param {beforeCommitCallback} callback the callback function
         */
        set BeforeCommit(callback) {
            if (typeof callback !== "function") {
                throw new Error(`Before commit callback must be a function, received ${typeof callback}`);
            }

            this.#before_commit_callback = callback.bind(this);
        }

        /**
         * Commits the changes made to the medias. also removes all the callbacks.
         * @param {string} category_id the uuid of the category
         * @returns {Promise<import('@libs/DungeonsCommunication/base').HttpResponse<Object> | null>} the response of the request
         * @see CommitCategoryTreeChangesRequest
         * @see HttpResponse
         */
        async commitChanges(category_id) {
            if (typeof this.#before_commit_callback === "function") {
                try {
                    if (this.#before_commit_callback.constructor.name === "AsyncFunction") {
                        await this.#before_commit_callback();
                    } else {
                        this.#before_commit_callback();
                    }
                } catch (e) {
                    if (this.#panic_on_before_callback_error) {
                        throw e;
                    } else {
                        console.error("Error on before commit callback:", e);
                        console.warn("Continuing with commit as panic_on_before_callback_error is set to false");
                    }
                }   
            }
            
            let response = await super.commitChanges(category_id);
            
            for (let id in this.#on_changes_callbacks) {
                this.unsubscribeToChanges(id);
            }

            return response;
        }

        /**
         * Clears the media of any changes and calls the callbacks with the NORMAL change type.
         * @param {string} media_uuid the uuid of the media
         * @returns {void}
         */
        clearMediaChanges(media_uuid) {
            super.clearMediaChanges(media_uuid);

            this.#broadcastChange(media_change_types.NORMAL, media_uuid);
        }

        /**
         * Clears all the deletion changes and calls the callbacks with the NORMAL change type.
         * @returns {void}
         */
        clearAllDeletionChanges() { 
            let affected_media = this.DeletedMedias;

            super.clearAllDeletionChanges();

            for (let m of affected_media) {
                console.log(`Broadcasting to deleted media ${m}`);
                this.#broadcastChange(media_change_types.NORMAL, m.uuid);
            }
        }


        /**
         * Clears all the move changes and calls the callbacks with the NORMAL change type.
         * @returns {void}
         */
        clearAllMoveChanges() {
            let affected_media = this.MovedMedias

            super.clearAllMoveChanges();

            for (let m of affected_media) {
                this.#broadcastChange(media_change_types.NORMAL, m.uuid);
            }
        }

        /**
         * Clears all the changes and calls the callbacks with the NORMAL change type.
         * @returns {void}
         */
        clearAllChanges() {
            this.clearAllDeletionChanges();
            this.clearAllMoveChanges();
        }

        /**
         * Broadcasts a given change type and media uuid to all the callbacks
         * @param {MediaChangeType} change_type the type of change
         * @param {string} media_uuid the uuid of the media
         * @returns {void}
         */
        #broadcastChange = (change_type, media_uuid) => {
            console.log(`Broadcasting change ${change_type} on media ${media_uuid}`);
            for (let id in this.#on_changes_callbacks) {
                try {
                    this.#on_changes_callbacks[id](change_type, media_uuid);
                } catch (e) {
                    console.warn(`Error calling callback ${id}:`, e, "Removing callback");
                    this.unsubscribeToChanges(id);
                }
            }
        }

        /**
         * Adds a media to the deleted medias set and calls the callbacks with th DELETE change type.
         * @param {Media} media the media to be deleted
         */
        stageMediaDeletion(media) {
            super.stageMediaDeletion(media);
            let applied_change = this.getMediaChangeType(media.uuid);
            this.#broadcastChange(applied_change, media.uuid);
        }

        /**
         * Adds a media to the moved medias map and calls the callbacks with the MOVED change type.
         * @param {Media} media the media to be moved
         * @param {InnerCategory} new_category the new category of the media
         * @returns {void}
         */
        stageMediaMove(media, new_category) {
            super.stageMediaMove(media, new_category);
            let applied_change = this.getMediaChangeType(media.uuid);
            this.#broadcastChange(applied_change, media.uuid);
        }

        /**
         * Removes a media from the delete staging set and calls the callbacks with the NORMAL change type.
         * @param {string} media_uuid the uuid of the media
         * @returns {void}
         */
        unstageMediaDeletion(media_uuid) {
            super.unstageMediaDeletion(media_uuid);
            this.#broadcastChange(media_change_types.NORMAL, media_uuid);
        }


        /**
         * Removes a media from the move staging map and the move staging data and calls the callbacks with the NORMAL change type.
         * @param {string} media_uuid the uuid of the media
         */
        unstageMediaMove = media_uuid => {
            super.unstageMediaMove(media_uuid);
            this.#broadcastChange(media_change_types.NORMAL, media_uuid);
        }

        /**
         * Adds a callback function to be called when a change is made to a media
         * @param {string} id the id of the callback
         * @param {onChangesCallback} callback the callback function
         */
        suscribeToChanges = (id, callback) => {
            this.#on_changes_callbacks[id] = callback;
        }

        /**
         * Removes a callback function from the list of callbacks
         * @param {string} id the id of the callback
         */
        unsubscribeToChanges = id => {
            delete this.#on_changes_callbacks[id];
        }
    }

/*=====  End of Media Changes  ======*/


/*=============================================
=            Category Search            =
=============================================*/

/**
 * Searches for categories based on the query
 * @param {string} query the query to search for
 * @param {string} cluster_id the cluster id to search in
 * @param {string} [ignore] an optional parameter to ignore a category by its uuid
 * @returns {Promise<import('@models/Categories').Category[]>} the response of the request
 */
export const searchCategories = async (query, cluster_id, ignore="") => {
    /**
     * @type {import('@models/Categories').Category[]}
     */
    let category_search_results = [];

    const request = new GetCategorySearchResultsRequest(query, cluster_id, ignore);

    let response = await request.do();

    if (response.Ok && response.data !== null) {
        category_search_results = response.data.map(category_params => new Category(category_params))
    }

    return category_search_results;
}

/*=====  End of Category Search  ======*/

/**
 * Requests the resyncronization of a cluster branch by it's category uuid. if one wants to resync the entire cluster, the root category uuid should be used.
 * @param {string} category_uuid
 * @param {string} cluster_id
 * @returns {Promise<boolean>} whether the resync was successful
 */
export const resyncClusterBranch = async (category_uuid, cluster_id) => {
    let sync_success = false;
    const request = new PutResyncClusterBranchRequest(cluster_id, category_uuid);
    
    let response  = await request.do();

    if (response.Ok) {
        sync_success = response.data;
    }

    return sync_success;
}

