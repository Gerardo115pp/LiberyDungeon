/*
.##.........######..##.....##....###....##....##....##.....##..#######..########..########.##........######.
.##....##..##....##.##.....##...##.##...###...##....###...###.##.....##.##.....##.##.......##.......##....##
.##....##..##.......##.....##..##...##..####..##....####.####.##.....##.##.....##.##.......##.......##......
.##....##..##.......#########.##.....##.##.##.##....##.###.##.##.....##.##.....##.######...##........######.
.#########.##.......##.....##.#########.##..####....##.....##.##.....##.##.....##.##.......##.............##
.......##..##....##.##.....##.##.....##.##...###....##.....##.##.....##.##.....##.##.......##.......##....##
.......##...######..##.....##.##.....##.##....##....##.....##..#######..########..########.########..######.
*/
import { 
    GetChanBoardsRequest,
    GetChanCatalogRequest,
    GetChanThreadRequest,
    PostThreadDownloadRequest,
    GetDownloadRegisterRequest
} from "@libs/HttpRequests";

/**
 * @typedef {Object} DownloadRegister a download register is a record of a completed download. in the case of thread medias downloads, the download_uuid is usually the thread uuid.
 * @property {boolean} exists Whether or not the download exists.
 * @property {number} download_count The number of media files that were downloaded.
*/

export class ChanTrackedBoard {
    constructor(board_name, description) {
        this.board_name = board_name;
        this.description = description;
    }
}

export class ChanCatalogThread {
    constructor({uuid, date, file, responses, images, teaser, image_url, teaser_thumb_width, teaser_thumb_height, board_name}) {
        this.uuid = uuid;
        this.date = date;
        this.file = file;
        this.responses = responses;
        this.images = images;
        this.description = teaser;
        this.image_url = image_url;
        this.teaser_thumb_width = teaser_thumb_width;
        this.teaser_thumb_height = teaser_thumb_height;
        this.board_name = board_name;
    }

    /**
     * Converts the date property from a unix timestamp to a formatted datetime string.
     * @returns {string} Formatted datetime string.
    */
    getCreationDate() {
        const date = new Date(this.date * 1000);

        return date.toLocaleString();
    }

    /**
     * Sends a request to the server to download the thread.
     * @param {string} category_name The category to download the thread to.
     * @param {string} [parent_uuid] The parent category's uuid.
     * @returns {Promise<string>} The response from the server.
    */
    async download(category_name, parent_uuid, cluster_uuid) {
        const request = new PostThreadDownloadRequest(this.uuid, this.board_name, category_name, parent_uuid, cluster_uuid);

        let response = await request.do();

        let download_uuid = null;

        if (response.status <= 200 && response.status < 300) {
            download_uuid = response.data.download_uuid;
        }

        return download_uuid;
    }


    
    /**
     * Checks if the thread has been downloaded.
     * @returns {Promise<DownloadRegister>} The download registry.
    */
    downloadRegister = async () => {
        let download_register = {
            exists: false,
            download_count: 0
        };

        const request = new GetDownloadRegisterRequest(this.uuid);

        let response = await request.do();

        if (response.status >= 200 && response.status < 300) {
            download_register = response.data;
        }

        return download_register;
    }
    
}

export class ChanThreadReply {
    constructor({uuid, message, date, file, thumbnail_url}) {
        this.uuid = uuid;
        this.message = message;
        if (isNaN(+date)) {
            date = parseInt(date);
            date = isNaN(date) ? new Date().getTime() : date;
        }

        this.date = date;
        this.file = file;
        this.thumbnail_url = thumbnail_url;
    }

    /**
     * Converts the date property from a unix timestamp to a formatted datetime string.
     * @returns {string} Formatted datetime string.
    */
    getCreationDate() {
        const date = new Date(this.date * 1000);

        return date.toLocaleString();
    }       

    hasImages() {
        return this.file !== "no-file";
    }
}

export class ChanThread {
    constructor({uuid, date, file, title, description, cover_image_url}) {
        this.uuid = uuid;
        if (isNaN(+date)) {
            date = parseInt(date);
            date = isNaN(date) ? new Date().getTime() : date;
        }
        this.date = date;
        this.file = file;
        this.title = title;
        this.description = description;
        this.cover_image_url = cover_image_url;
        /**
         * @type {ChanThreadReply[]}
         */
        this.replies = [];
    }

    /**
     * Converts the date property from a unix timestamp to a formatted datetime string.
     * @returns {string} Formatted datetime string.
    */
    getCreationDate() {
        const date = new Date(this.date * 1000);

        return date.toLocaleString();
    }
}

export const getTrackedBoards = async () => {
    const request = new GetChanBoardsRequest();

    const response = await request.do();

    const boards = [];

    for (let board_name of Object.keys(response.data)) {
        let tracked_board = new ChanTrackedBoard(board_name, response.data[board_name]);

        boards.push(tracked_board);
    }

    return boards;
}

export const getBoardCatalog = async (board_name) => {
    const request = new GetChanCatalogRequest(board_name);

    const response = await request.do();

    const threads = [];

    for (let thread of response.data) {
        let catalog_thread = new ChanCatalogThread(thread);

        threads.push(catalog_thread);
    }

    return threads;
}

/**
 * Gets the content of a thread.
 * @param {string} board_name The board name.
 * @param {string} thread_uuid The thread uuid.
 * @returns {Promise<ChanThread>} The thread.
 */
export const getThreadContent = async (board_name, thread_uuid) => {
    const request = new GetChanThreadRequest(thread_uuid, board_name);

    const response = await request.do();
    
    if (response.data.replies === undefined) {
        return null;
    }

    /**
     * @type {ChanThreadReply[]}
     */
    const replies = [];

    for (let reply of response.data.replies) {
        let thread_reply = new ChanThreadReply(reply);

        replies.push(thread_reply);
    }

    const thread = new ChanThread(response.data);

    thread.replies = replies;

    return thread;
}