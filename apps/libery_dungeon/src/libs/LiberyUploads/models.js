import { getChunkedUploadTicket, getUploadStreamTicket, insertMediaFile, uploadChunk } from "./requests";
import { humanReadableSize } from "./utils";

    
/**
 * The max byte size of a media file that we will handle in a single upload.
 * @type {number}
 */
export const MAX_MEDIA_FILE_SIZE = 100 * 1024 * 1024;

export const io_errors = {
    ERR_EOF: "ERR_EOF",
    ERR_INVALID_SEEK: "ERR_INVALID_SEEK",
}

export const seek_whence = {
    SeekStart: 0,
    SeekCurrent: 1,
    SeekEnd: 2,
}

export class MediaFile {
    /** @type {File} */
    #file_object

    /**
     * The URL source of the media file. Setted by calling loadURL() 
     * @type {string}
     */
    #src

    /**
     * The file's seek position.
     * @type {number}
     */
    #seek_pointer = 0;

    /**
     * @param {File} file the name of the media file
     */
    constructor (file) {
        this.#src = "";
        this.name = file.name;
        this.type = file.type;
        this.size = file.size;
        this.#file_object = file;

        // Callbacks 
        
        /**
         * Called when the media file has been uploaded.
         * @type {() => void}
         */
        this.onUploaded = () => {};

        /**
         * Called when the media is been uploaded in chunks and there is a progress update.
         * @param {number} chunks_uploaded
         * @param {number} total_chunks
         */
        this.onChunkUploadProgress = (chunks_uploaded, total_chunks) => {};
    }

    /**
     * Returns the media file as a Blob
     * @returns {Blob}
     */
    get File() {
        return this.#file_object;
    }

    /**
     * Uploads the file in a single request. If the file is too large, it will throw an error.
     * Returns true if the file was uploaded successfully, false otherwise.
     * @param {string} category_id
     * @returns {Promise<boolean>}
     */
    async fastUpload() {
        if (this.#isFileTooLarge()) {
            throw new Error("File is too large to be uploaded in a single request");
        }

        let was_uploaded = await insertMediaFile(this);

        if (was_uploaded) {
            if (this.onUploaded?.constructor.name === "AsyncFunction" || this.onUploaded?.constructor.name === "Function") {
                this.onUploaded();
            } else {
                console.error("onUploaded is not a function");
            }
        }

        return was_uploaded;
    }

    /**
     * Returns the mime type of the media file
     * @returns {string}
     */
    getMimeType () {
        return this.type;
    }

    /**
     * Returns a human readable string of the size of the media file
     * @returns {string}
     */
    get HumanReadableSize() {   
        return humanReadableSize(this.size);
    }

    /**
     * Returns true if the media file is an image, false if it is a video and throws an error if it is neither.
     * @param {boolean} strict if false and the media file is neither an image nor a video, the function will return null instead of throwing an error
     * @returns {boolean}
     * @throws {Error} if the media file is neither an image nor a video
     */
    isImage (strict = true) {
        if (this.type.startsWith("image")) {
            return true;
        } else if (this.type.startsWith("video")) {
            return false;
        } else if (strict) {
            throw new Error(`MediaFile is neither an image nor a video: ${this.type}`);
        }

        return null;
    }

    /**
     * Returns whether the media file can be uploaded in a single request
     * @returns {boolean}
     */
    #isFileTooLarge() {
        return this.size > MAX_MEDIA_FILE_SIZE;
    }

    /**
     * Returns whether the file SHOULD be uploaded in a single request
     * @returns {boolean}
     */
    get IsFileTooLarge() {
        return this.size > (MAX_MEDIA_FILE_SIZE * 0.8); // MAX_MEDIA_FILE_SIZE is the absolute max size, but it would be better to upload files that are 80% of the max size in chunks
    }

    /**
     * Loads the File object as a URL on src
     * @returns {Promise<void>}
     */
    async loadURL() {
        if (this.#isFileTooLarge()) {
            throw new Error("File is too large to be loaded as a URL");
        }
        
        const file_reader = new FileReader();

        return new Promise((resolve, reject) => {
            file_reader.onload = () => {
                this.#src = file_reader.result;
                resolve();
            }

            file_reader.onerror = () => {
                reject(file_reader.error);
            }

            file_reader.readAsDataURL(this.#file_object);
        });
    }

    /**
     * Reads a chunk of the media file equal to N and returns Uint8Array of the chunk. Out of bounds reads will return alway return undefined.
     * @param {number} N 
     * @returns {Promise<ReadResult>}
     * @typedef {Object} ReadResult
     * @property {Uint8Array | undefined} chunk
     * @property {string | null} error
     */
    async read(N) {
        /**
         * @type {ReadResult}
         */
        const read_result = {
            chunk: undefined,
            error: null,
        }

        if (N === 0) return read_result;

        const chunk = this.readBlob(N);
        if (chunk === undefined) {
            read_result.error = io_errors.ERR_EOF;
            return read_result;
        }
        const chunk_buffer = await chunk.arrayBuffer();

        read_result.chunk = new Uint8Array(chunk_buffer);

        return read_result;
    }

    /**
     * Reads a chunk of the media file equal to N and returns it as a Blob. Attempts to read out more than the file size will return the remaining bytes.
     * and a EOF error.
     * @param {number} N
     * @returns {Promise<ReadBlobResult>}
     * @typedef {Object} ReadBlobResult
     * @property {Blob | undefined} chunk
     * @property {string | null} error
     */
    async readBlob(N) {
        /**
         * @type {ReadBlobResult}
         */
        let read_blob_result = {
            chunk: undefined,
            error: null,
        }

        if (this.#seek_pointer >= this.size) {
            return read_blob_result;
        }

        if (N === 0) return read_blob_result;

        let read_end = this.#seek_pointer + N;

        if (read_end > this.size) {
            read_end = this.size;
            read_blob_result.error = io_errors.ERR_EOF;
        }

        read_blob_result.chunk  = this.#file_object.slice(this.#seek_pointer, read_end);
    
        this.#seek_pointer = read_end;

        return read_blob_result;
    }


    /**
     * Seek sets the offset for the next Read to offset, interpreted according to whence: SeekStart means relative to the start of
     * the file, SeekCurrent means relative to the current offset, and SeekEnd means relative to the end (for example, offset = -2 specifies the penultimate byte of the file).
     * Seek returns the new offset relative to the start of the file or an error, if any. 
     * 
     * Seeking to an offset before the start of the file is an error. Seeking to any positive offset may be allowed, but if the new offset exceeds the size of the underlying object, it will 
     * throw an EOF error.
     * @param {number} offset
     * @param {number} whence
     * @returns {SeekResult}
     * @typedef {number} SeekResult
     * @property {number} offset
     * @property {string | null} error
     */
    seek(offset, whence) {
        let new_offset = 0;
        let err = null;

        switch (whence) {
            case seek_whence.SeekStart:
                new_offset = offset;
                break;
            case seek_whence.SeekCurrent:
                new_offset = this.#seek_pointer + offset;
                break;
            case seek_whence.SeekEnd:
                new_offset = this.size + offset;
                break;
        }

        if (new_offset < 0) {
            err = io_errors.ERR_INVALID_SEEK;
        }

        if (new_offset > this.size) {
            err = io_errors.ERR_EOF;
        }

        this.#seek_pointer = new_offset;

        return {
            offset: new_offset,
            error: err,
        };
    }

    /**
     * Returns the Src of the media file. if loadURL() has not been called, it will return an empty string.
     * @returns {string}
     */
    get Src() {
        return this.#src;
    }
}

export class ChunkedMediaUpload {
    /**
     * The media file to be uploaded
     * @type {MediaFile}
     */
    #media_file

    /**
     * @param {category_id} category_id
     * @param {MediaFile} media_file
     */
    constructor(category_id, media_file) {
        this.category_id = category_id;
        this.#media_file = media_file;
    }

    /**
     * Returns the amount of chunks of size MAX_MEDIA_FILE_SIZE that will be required to upload the media file.
     * @returns {number}
     */
    get ChunkCount() {
        return Math.ceil(this.#media_file.size / MAX_MEDIA_FILE_SIZE);
    }

    /**
     * Fetches an upload ticket for a media chunked upload. 
     * @returns {Promise<string>}
     */
    async getUploadTicket() {
        const upload_uuid = crypto.randomUUID();

        return await getChunkedUploadTicket({ 
            upload_uuid,
            upload_filename: this.#media_file.name,
            upload_size: this.#media_file.size,
            upload_chunks: this.ChunkCount,
            category_uuid: this.category_id,
        })
    }

    /**
     * Starts the upload of the media file in chunks. Returns true if the upload was successful, false otherwise.
     * @returns {Promise<boolean>}
     */
    async startUpload() {
        let upload_ticket = await this.getUploadTicket();
        if (upload_ticket === null) {
            console.error("Failed to get upload ticket. the ticket is null");   
            return false;
        }

        let upload_successful = false;

        for (let h = 0; h < this.ChunkCount; h++) {
            const read_blob_result = await this.#media_file.readBlob(MAX_MEDIA_FILE_SIZE);
            if (read_blob_result.chunk === undefined) {
                console.error("Failed to read chunk");
                return false;
            }

            this.#media_file.onChunkUploadProgress(h + 1, this.ChunkCount);

            const upload_state = await uploadChunk(upload_ticket, read_blob_result.chunk, h);
            
            if (upload_state === -1) {
                console.error(`Failed to upload chunk ${h}`);
                return false;
            } else if (upload_state === 1) {
                upload_successful = true;
            }
        }

        return upload_successful;
    }
}

/**
 * Uploads a list of media files. Handles the case where the file is too large to be uploaded in a single request and
 * uploads it as a chunked request instead.
 */
export class MediaUploader {

    /**
     * Files to be uploaded
     * @type {import('./models').MediaFile[]}
     */
    #files = [];

    /**
     * @param {import('./models').MediaFile[]} files
     */
    constructor(files) {
        this.#files = files;

        // Callbacks

        /**
         * Called when all media files have been uploaded.
         * @type {() => void}
         */
        this.onAllUploaded = () => {};

        /**
         * Called when a media file has been uploaded.
         * @param {MediaFile} file
         * @param {number} index
         */
        this.onFileUploaded = (file, index) => console.log(`File ${file.name} uploaded`);

        /**
         * Called when a media file has failed to upload.
         * @param {MediaFile} file
         * @param {number} index
         * @param {string} error
         */
        this.onFileUploadFailed = (file, index, error) => console.error(`File ${file.name} failed to upload: ${error}`);

        /**
         * Called when an entire upload stream has failed.
         * @param {string} error
         */
        this.onUploadStreamFailed = (error) => console.error(`Upload stream failed: ${error}`);
    }

    /**
     * Uploads the media files. If a file is too large, it will be uploaded as a chunked request.
     * @param {string} category_id
     */
    async startUpload(category_id) {
        let upload_authorized = await getUploadStreamTicket(this.#files.length, category_id);

        if (!upload_authorized) {
            this.onUploadStreamFailed("Failed to get upload ticket");

            return;
        }
        
        for (let h = 0; h < this.#files.length; h++) {
            const file = this.#files[h];
            let was_uploaded = false;
            
            if (file.IsFileTooLarge) {
                const chunked_upload = new ChunkedMediaUpload(category_id, file);
                was_uploaded = await chunked_upload.startUpload();
            } else {
                was_uploaded = await file.fastUpload();
            }

            if (was_uploaded) {
                this.onFileUploaded(file, h);
                file.onUploaded();
            } else {
                this.onFileUploadFailed(file, h, "Failed to upload file");
            }
        }

        this.onAllUploaded();
    }
}