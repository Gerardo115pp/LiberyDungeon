import { HttpResponse, attributesToJson } from "../base";
import { categories_server } from "../services";

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