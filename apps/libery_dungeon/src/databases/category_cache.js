import { browser } from "$app/environment";

/**
 * @type {IDBDatabase} the database object
*/
let db;

let load_promise = null;

if (browser) {
    load_promise = new Promise((resolve, reject) => {
        console.log("LOADING CATEGORY CACHE");
        if (!browser) {
            reject("IndexedDB is not available in this environment");
        }
    
    
        /**
         * @type {IDBOpenDBRequest} the request to open the database
        */
        const request = indexedDB.open("Pandasworld", 1);
        
        
        /**
         * A handler in case the database is not found or is opened with a different version
         * @param {IDBVersionChangeEvent} e the event that triggered the handler
         * @returns {void}
        */
        request.onupgradeneeded = e => {
            db = e.target.result;
            const category_cache = db.createObjectStore("category_cache", { keyPath: "uuid"});
            resolve();
        }
        
        /**
         * A handler in case the database is opened successfully
         * @param {IDBVersionChangeEvent} e the event that triggered the handler
        */
        request.onsuccess = e => {
            db = e.target.result;
            resolve();
        }
        
        /**
         * A handler in case the database is not opened successfully
         * @param {IDBVersionChangeEvent} e the event that triggered the handler
        */
        request.onerror = e => {
            
            reject(e.target.error);
        }
    }).catch(e => console.warn("Error loading category cache: ", e));
}


/**
 * Adds a category index to the cache
 * @param {string} category_uuid the uuid of the category
 * @param {number} media_index the index of the media
 * @returns {void}
*/
export const addCategoryIndex = (category_uuid, media_index) => {
    if (!browser) return;

    load_promise.then(() => {
        const transaction = db.transaction(["category_cache"], "readwrite");
        const object_store = transaction.objectStore("category_cache");
        const request = object_store.add({ uuid: category_uuid, media_index: media_index });
    
        request.onerror = e => {
            console.error("Error adding category index to cache: ", e.target.error);
        }
    });
}

/**
 * Gets the media index of a category from the cache
 * @param {string} category_uuid the uuid of the category
 * @returns {Promise<number>} a promise that resolves to the media index of the category
 * @async
*/
export const getCategoryIndex = async (category_uuid) => {
    if (!browser) return 0;

    let real_promise = await load_promise.then(async () => {
        const transaction = db.transaction(["category_cache"], "readonly");
        const object_store = transaction.objectStore("category_cache");
        const request = object_store.get(category_uuid);
    
        let index_promise = new Promise((resolve, reject) => {
            request.onsuccess = e => {
                if (e.target.result === undefined) {
                    resolve(-1);
                } else {
                    resolve(e.target.result.media_index);
                }
            }
    
            request.onerror = e => {
                reject(e.target.error);
            }
        });
    
        let index = await index_promise;
    
        return index;
    });

    return real_promise;
}


/**
 * Updates the media index of a category in the cache
 * @param {string} category_uuid the uuid of the category
 * @param {number} media_index the index of the media
 * @returns {void}
*/
export const updateCategoryIndex = (category_uuid, media_index) => {
    if (!browser) return;

    load_promise.then(() => {
        const transaction = db.transaction(["category_cache"], "readwrite");
        const object_store = transaction.objectStore("category_cache");
        const request = object_store.put({ uuid: category_uuid, media_index: media_index });

        request.onerror = e => {
            console.error("Error updating category index in cache: ", e.target.error);
        }
    });
}


/** 
 * Returns all the category indexes in the cache
 * @returns {Promise<Map<string, number>>} a map of category uuids to media indexes
 * @async
*/
export const getAllCategoryIndexes = async () => {
    if (!browser) return new Map();

    let real_promise = await load_promise.then(async () => {
        const transaction = db.transaction(["category_cache"], "readonly");
        const object_store = transaction.objectStore("category_cache");
        const request = object_store.getAll();
    
        let result_promise = new Promise((resolve, reject) => {
            request.onsuccess = e => {
                let result = new Map();
    
                for (let i = 0; i < e.target.result.length; i++) {
                    result.set(e.target.result[i].uuid, e.target.result[i].media_index);
                }
    
                resolve(result);
            }
    
            request.onerror = e => {
                reject(e.target.error);
            }
        });

        let result = await result_promise;
    
        return result;
    });

    return real_promise;
}
