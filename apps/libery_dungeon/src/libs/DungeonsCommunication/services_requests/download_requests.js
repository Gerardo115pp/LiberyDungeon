import { collect_server, downloads_server } from "../services";
import { HttpResponse, attributesToJson, attributesToJsonExclusive } from "../base";

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

    /**
     * @param {string} board_name 
     */
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

    /**
     * @param {string} thread_id 
     * @param {string} board_name 
     */
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
     * @param {string} cluster_uuid
     */
    constructor(thread_uuid, board_name, target_category_name, parent_uuid, cluster_uuid) {
        this.thread_uuid = thread_uuid;
        this.board_name = board_name;
        this.target_category_name = target_category_name;
        this.parent_uuid = parent_uuid;
        this.cluster_uuid = cluster_uuid;
    }

    toJson = attributesToJsonExclusive.bind(this);

    /**
     * @returns {Promise<HttpResponse<{ download_uuid: string }>>}
     */
    do = async () => {
        const response = await fetch(`${collect_server}/4chan-downloads/thread/images`, {   
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        }); 

        /**
         * @type {{download_uuid: string}}
         */
        let data = {
            download_uuid: ""
        };

        if (response.status <= 200 && response.status < 300) {
            data = await response.json();
        }

        return new HttpResponse(response, {
            download_uuid: data.download_uuid
        });
    }
}

export class GetDownloadRegisterRequest {

    /**
     * @param {string} download_uuid 
     */
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