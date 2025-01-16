import { metadata_server } from "@libs/DungeonsCommunication/services";
import { HttpResponse, attributesToJson } from "@libs/DungeonsCommunication/base";

/**
 * Creates a new video moment for a give video media(uuid) which was to exist in the given category_cluster(the video_cluster param)
 */
export class PostVideoMomentsRequest {

    static endpoint = `${metadata_server}/video-moments/moments`

    /**
     * @param {string} video_uuid
     * @param {string} video_cluster
     * @param {number} moment_time
     * @param {string} moment_title
     */
    constructor(video_uuid, video_cluster, moment_time, moment_title) {
        this.video_uuid = video_uuid;
        this.video_cluster = video_cluster;
        this.moment_time = moment_time;
        this.moment_title = moment_title;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('@models/Metadata').VideoMomentParams>>}
     */
    do = async () => {
        /**
         * @type {import('@models/Metadata').VideoMomentParams}
         */
        let new_moment_params = {
            id: NaN,
            video_uuid: this.video_uuid,
            moment_title: this.moment_title,
            moment_time: this.moment_time
        }

        const url = new URL(PostVideoMomentsRequest.endpoint, globalThis.location.origin);

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            if (response.ok) {
                const response_body = await response.json();

                new_moment_params.id = response_body.response ?? NaN;
            }
        } catch (err) {
            console.error(`In @libs/DungeonsCommunication/services_requests/metadata_requests/video_moments_requests PostVideoMomentsHandler.do: ${err}`);
            throw err;
        }

        return new HttpResponse(response, new_moment_params);
    }
}

/**
 * Returns all the video moments that have been saved for a video.
 */
export class GetVideoMomentsRequest {

    static endpoint = `${metadata_server}/video-moments/video-moments`;

    /**
     * @param {string} video_uuid
     * @param {string} video_cluster
     */
    constructor(video_uuid, video_cluster) {
        this.video_uuid = video_uuid;
        this.video_cluster = video_cluster;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('@models/Metadata').VideoMomentParams[]>>}
     */
    do = async () => {
        /**
         * @type {import('@models/Metadata').VideoMomentParams[]}
         */
        let video_moments = [];

        /**
         * @type {Response}
         */
        let response;

        const url = new URL(GetVideoMomentsRequest.endpoint, globalThis.location.origin);

        url.searchParams.append("video_uuid", this.video_uuid);
        url.searchParams.append("video_cluster", this.video_cluster);

        try {
            response = await fetch(url);

            if (response.ok) {
                video_moments = await response.json();
            }
        } catch (err) {
            console.error(`In @libs/DungeonsCommunication/services_requests/metadata_requests/video_moments_requests GetVideoMomentsHandler.do: ${err}`);
            throw err;
        }

        return new HttpResponse(response, video_moments);
    }
}

/**
 * Updates the title and/or the time of a video moment by its id.
 */
export class PutVideoMomentDataRequest {

    static endpoint = `${metadata_server}/video-moments/moment`;

    /**
     * @param {number} id - the id of the video moment.
     * @param {number} moment_time
     * @param {string} moment_title
     */
    constructor(id, moment_time, moment_title) {
        this.id = id;
        this.moment_time = moment_time;
        this.moment_title = moment_title;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        /**
         * @type {boolean}
         */
        let successful_update = false;

        /**
         * @type {Response}
         */
        let response;

        const url = new URL(PutVideoMomentDataRequest.endpoint, globalThis.location.origin);

        try {
            response = await fetch(url, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

            successful_update = response.ok;
        } catch (err) {
            console.error(`In @libs/DungeonsCommunication/services_requests/metadata_requests/video_moments_requests PutVideoMomentDataHandler.do: ${err}`);
            throw err;
        }

        return new HttpResponse(response, successful_update);
    }
}

/**
 * Deletes a video moment by it's id.
 */
export class DeleteVideoMomentRequest {

    static endpoint = `${metadata_server}/video-moments/moment`;

    /**
     * @param {number} id - the moment's identifier
     */
    constructor(id) {
        this.id = id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {boolean}
         */
        let successful_deletion = false;

        const url = new URL(DeleteVideoMomentRequest.endpoint, globalThis.location.origin);

        url.searchParams.append("id", this.id.toString())

        try {
            response = await fetch(url, {
                method: "DELETE"
            });

            successful_deletion = response.ok;
        } catch (err) {
            console.error(`In @libs/DungeonsCommunication/services_requests/metadata_requests/video_moments_requests DeleteVideoMomentHandler.do: ${err}`);
            throw err;
        }

        return new HttpResponse(response, successful_deletion);
    }
}