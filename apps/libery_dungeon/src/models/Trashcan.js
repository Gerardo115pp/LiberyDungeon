import { CategoryWeakIdentity } from "./Categories";
import { Media } from "./Medias.js";
import { 
    GetTrashcanEntriesRequest,
    GetTrashcanTransactionRequest,
    PatchRestoreMediaRequest,
    PatchRestoreTransactionRequest,
    DeleteTrashcanTransactionRequest,
    DeleteEmptyTrashcanRequest,
} from "@libs/DungeonsCommunication/services_requests/transactions_requests";
import { getTrashcanMediaUrl } from "@libs/DungeonsCommunication/services_requests/media_requests";


/**
* @typedef {Object} TrashcanTransactionEntryParams
 * @property {string} timestamp - the timestamp of the entry also the identifier of transaction related to this entry
 * @property {number} affected_medias - the number of medias affected by this transaction
*/

/**
* @typedef {Object} TrashcanTransactionParams
 * @property {string} transaction_id
 * @property {import('./Medias.js').MediaParams[]} content
 * @property {import('./Categories.js').CategoryWeakIdentityParams} origin_identity
*/

export class TrashedMedia {
    /**
     * The media object that was trashed.
     * @type {Media}
     */
    #the_media;

    /**
     * The identity form which the media was trashed.
     * @type {CategoryWeakIdentity}
     */
    #rejected_from; 

    /**
     * @param {import('@models/Medias').MediaParams} media_params   
     * @param {CategoryWeakIdentity} category_identity
     */
    constructor(media_params, category_identity) {
        this.#the_media = new Media(media_params);  
        this.#rejected_from = category_identity
    }

    /**
     * The name of the media with file extension.
     * @type {string}
     */
    get Name() {
        return this.#the_media.name;
    }

    /** 
     * The media uuid.
     * @type {string}
     */
    get UUID() {
        return this.#the_media.uuid;
    }

    /**
     * A url for a thumbnail of the trashed media.
     * @type {string}
     */
    get ThumbnailURL() {
        let thumbnail_url = getTrashcanMediaUrl(this.#the_media.name);
        console.log(thumbnail_url);
        return getTrashcanMediaUrl(this.#the_media.name);
    }

    /**
     * Returns a MediaParams object of the trashed media.
     * @returns {import('@models/Medias').MediaParams}
     */
    toParams() {
        return this.#the_media.toParams();
    }
}

export class TrashcanTransactionEntry {
    /**
     * @param {TrashcanTransactionEntryParams} param0
     */
    constructor({timestamp, affected_medias}) {
        this.timestamp = timestamp;
        this.affected_medias = affected_medias;
    }
}

export class TrashcanTransaction {

    /**
     * The identifier of the transaction. also a human readable datetime string.
     * @type {string}
     */
    #transaction_id;

    /**
     * A list of medias that were trashed in this transaction.
     * @type {TrashedMedia[]}
     */
    #content;

    /**
     * The category identity of the transaction.
     * @type {CategoryWeakIdentity}
     */
    #category_identity;

    /**
     * @param {TrashcanTransactionParams} param0
     */
    constructor({transaction_id, content, origin_identity}) {
        this.#transaction_id = transaction_id;
        this.#category_identity = this.#loadCategoryIdentity(origin_identity);
        this.#content = this.#loadContent(content, this.#category_identity);
    }

    /**
     * The trashed content of the transaction.  
     * @type {TrashedMedia[]}
     */
    get Content() {
        return this.#content;   
    }

    /**
     * The category identity of the transaction.
     * @type {CategoryWeakIdentity}
     */
    get CategoryWeakIdentity() {
        return this.#category_identity;
    }

    /**
     * Remove media from content.
     * @param {string} media_uuid
     */
    removeMedia(media_uuid) {
        this.#content = this.#content.filter(media => media.UUID !== media_uuid);
    }

    /**
     * Returns an exact copy of the transaction.
     * @returns {TrashcanTransaction}
     */
    clone() {
        let clone_content = this.#content.map(media => media.toParams());       
        let clone_category_identity = this.#category_identity.toParams();

        return new TrashcanTransaction({
            transaction_id: this.#transaction_id,
            content: clone_content,
            origin_identity: clone_category_identity
        });


    }

    /**
     * The identifier of the transaction. also a human readable datetime string.    
     * @type {string}
     */
    get TransactionID() {       
        return this.#transaction_id;
    }       


    /**
     * Takes a list of media parameters and loads them into Media instances.
     * @param {import('./Medias.js').MediaParams[]} content
     * @param {CategoryWeakIdentity} category_identity
     * @returns {Media[]}
     */
    #loadContent(content, category_identity) {
        return content.map(media => new TrashedMedia(media, category_identity));
    }

    /**
     * Takes a category identity parameters and loads them into a CategoryWeakIdentity instance.
     * @param {import('./Categories.js').CategoryWeakIdentityParams} category_identity
     * @returns {CategoryWeakIdentity}
     */
    #loadCategoryIdentity(category_identity) {
        return new CategoryWeakIdentity(category_identity);
    }
}

/**
 * Retrieves all trashcan entries.
 * @returns {Promise<TrashcanTransactionEntry[]>}
 */
export const getTrashcanEntries = async () => {
    /** @type {TrashcanTransactionEntry[]} */
    const entries = [];

    const request = new GetTrashcanEntriesRequest();

    const response = await request.do();

    if (response.Ok) {
        response.data.forEach(entry => {
            entries.push(new TrashcanTransactionEntry(entry));
        });
    }

    return entries;
}

/**
 * Retrieves a trashcan transaction.
 * @param {string} transaction_id
 * @returns {Promise<TrashcanTransaction>}
 */
export const getTrashcanTransaction = async (transaction_id) => {
    const request = new GetTrashcanTransactionRequest(transaction_id);

    const response = await request.do();

    if (response.Ok) {
        return new TrashcanTransaction(response.data);
    }

    return null;
}

/**
 * Restores a media from the trashcan.
 * @param {string} media_uuid
 * @param {string} transaction_id
 * @param {string} main_category
 * @returns {Promise<boolean>}
 */
export const restoreMedia = async (media_uuid, transaction_id, main_category) => {
    const request = new PatchRestoreMediaRequest(media_uuid, transaction_id, main_category);

    const response = await request.do();    

    return response.Ok; 
}

/**
 * Restores all medias in a transaction.
 * @param {string} transaction_id   
 * @param {string} main_category
 * @returns {Promise<boolean>}
 */
export const restoreTransaction = async (transaction_id, main_category) => {
    const request = new PatchRestoreTransactionRequest(transaction_id, main_category);

    const response = await request.do();

    return response.Ok;
}

/**
 * Deletes a transaction from the trashcan.
 * @param {string} transaction_id
 * @returns {Promise<boolean>}
 */
export const deleteTransaction = async (transaction_id) => {
    const request = new DeleteTrashcanTransactionRequest(transaction_id);

    const response = await request.do();

    return response.Ok;
}

/** 
 * Empties the trashcan. This permanently deletes all files in the trashcan.
 * @returns {Promise<boolean>}
 */
export const emptyTrashcan = async () => {
    const request = new DeleteEmptyTrashcanRequest();

    const response = await request.do();

    return response.Ok;
}