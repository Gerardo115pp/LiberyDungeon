import { medias_server } from "../services";
import { attributesToJson, HttpResponse } from "../base";

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
             * @type {import("../base").SingleStringResponse}
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