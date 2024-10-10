/*
.########...#######..##......##.##....##.##........#######.....###....########...######.
.##.....##.##.....##.##..##..##.###...##.##.......##.....##...##.##...##.....##.##....##
.##.....##.##.....##.##..##..##.####..##.##.......##.....##..##...##..##.....##.##......
.##.....##.##.....##.##..##..##.##.##.##.##.......##.....##.##.....##.##.....##..######.
.##.....##.##.....##.##..##..##.##..####.##.......##.....##.#########.##.....##.......##
.##.....##.##.....##.##..##..##.##...###.##.......##.....##.##.....##.##.....##.##....##
.########...#######...###..###..##....##.########..#######..##.....##.########...######.
*/
import { GetCurrentDownloadRequest } from "@libs/HttpRequests";


export class DownloadProgress {
    /**
     * @type {string} the uuid of the download
     */
    #download_uuid;
    /**
     * @type {number} the total number of files in the download
     */
    #total_files;
    /**
     * @type {number} the number of files that have been downloaded so far
     */
    #downloaded_files;
    /**
     * @type {boolean} whether the download is completed or not
     */
    #completed;
    constructor({download_uuid, total_files, downloaded_files, completed}) {
        this.#download_uuid = download_uuid;
        this.#total_files = total_files;
        this.#downloaded_files = downloaded_files;
        this.#completed = completed;
    }

    get DownloadUUID() {
        return this.#download_uuid;
    }

    get TotalFiles() {
        return this.#total_files;
    }

    get DownloadedFiles() {
        return this.#downloaded_files;
    }

    get Completed() {
        return this.#completed;
    }

    isEmpty() {
        return this.#download_uuid === "";
    }

    /**
     * @return {number} a percentage value between 0 and 100
     */
    percentComplete() {
        return this.#downloaded_files / this.#total_files * 100;
    }
}

export const EMPTY_DOWNLOAD_PROGRESS = new DownloadProgress({
    download_uuid: "",
    total_files: 0,
    downloaded_files: 0,
    completed: false
});


/*=============================================
=            Methods            =
=============================================*/

/**
 * Gets the current download uuid from the server.
 * @returns {Promise<string>} The response from the server.
 */
export const getCurrentDownloadUUID = async () => {
    const request = new GetCurrentDownloadRequest();

    const response = await request.do();

    return response.data.download_uuid;
}

/*=====  End of Methods  ======*/



