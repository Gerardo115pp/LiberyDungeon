import { 
    PostMediaRequest, 
    GetUploadStreamTicketRequest,
    GetChunkedUploadTicketRequest,
    PostChunkedUploadRequest 
} from '@libs/HttpRequests';

/**
 * Inserts a media file into a category, returns true if the operation was successful, false otherwise
 * @param {import('./models').MediaFile} media_file
 * @returns {Promise<boolean>}
 */
export const insertMediaFile = async (media_file) => {
    let inserted = false;
    const request = new PostMediaRequest(media_file.name, media_file.File);

    const response = await request.do();

    // 201 = Means all media uploads that the ticket been used allowed for, have been used up.
    // 204 = Means the media file was uploaded successfully and the ticket continues to be valid.
    if (response.status === 201 || response.status === 204) {
        inserted = true;
    }
    
    return inserted;
}

/**
 * Returns a ticket to upload a file in chunks. If successful, returns the ticket which is a jwt token, otherwise null.
 * @param {import('@libs/HttpRequests').GetChunkedUploadTicketParams} ticket_params
 * @returns {Promise<string|null>}
 */
export const getChunkedUploadTicket = async (ticket_params) => {
    const request = new GetChunkedUploadTicketRequest(ticket_params);
    let ticket = null;

    const response = await request.do();

    if (response.status === 200) {
        ticket = response.data;
    }
    
    return ticket;
}

/**
 * Requests an upload stream ticket. This is a jwt that authorizes each upload request. it has embedded the number of file uploads allowed. and after that
 * limit is exceeded, insert file requests will fail unless a new ticket is requested. Requires a user to have the 'upload_files' or 'ALL_PRIVILEGES' grants.
 * The token validity is 10 minutes per file that will be uploaded. After each file is uploaded, the token's validity is recalculated with only
 * the remaining number of files been considered. The token is stored as an HttpOnly cookie, not directly returned to the client.
 * @param {int} number_of_uploads
 * @param {string} recipient_category_uuid
 * @returns {Promise<boolean>}
 */
export const getUploadStreamTicket = async (number_of_uploads, recipient_category_uuid) => {
    let upload_uuid = crypto.randomUUID();

    const request = new GetUploadStreamTicketRequest(upload_uuid, number_of_uploads, recipient_category_uuid);
    let ticket_set = false;

    const response = await request.do();

    ticket_set = response.data.response;

    return ticket_set;
}

/**
 * Uploads a chunk of a file to the server. Returns:
 * -1 = the chunk was not uploaded
 * 0 = the chunk was uploaded successfully but the upload is not complete
 * 1 = the chunk was uploaded successfully and the upload is complete
 * @param {string} ticket
 * @param {Blob} chunk
 * @param {number} chunk_number
 * @returns {Promise<number>}
 */
export const uploadChunk = async (ticket, chunk, chunk_number) => {
    let uploaded = -1;
    const request = new PostChunkedUploadRequest(ticket, chunk, chunk_number);

    const response = await request.do();

    if (response.status === 201) {
        uploaded = 1;
    } else if (response.status === 204) {
        uploaded = 0;
    }

    return uploaded;
}