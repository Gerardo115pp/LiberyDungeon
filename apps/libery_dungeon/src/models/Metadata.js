import { GetWatchPointRequest, PostWatchPointRequest } from "@libs/DungeonsCommunication/services_requests/metadata_requests/watch_points_requests";
import { GetVideoMomentsRequest, PostVideoMomentsRequest } from "@libs/DungeonsCommunication/services_requests/metadata_requests/video_moments_requests";


/*=============================================
=            Watch points            =
=============================================*/

/**
 * @typedef {Object} WatchPoint
 * @property {number} ord
 * @property {string} media_uuid
 * @property {number} start_time
 */

/**
 * Requests the last watch point of a media with a time dimension(video, audio, that kind of thing). 
 * Used to resume the watch time of a media.
 * @param {string} media_uuid
 * @returns {Promise<number>}
 */
export const getMediaWatchPoint = async media_uuid => {
    const request = new GetWatchPointRequest(media_uuid);
    let response = await request.do();

    let last_watch_point = 0;

    if (response.Ok) {
        last_watch_point = response.data.start_time;
    }

    return last_watch_point;
}

/**
 * Saves a watch point of a media.
 * @param {string} media_uuid
 * @param {number} start_time
 * @returns {Promise<boolean>}
 */
export const saveMediaWatchPoint = async (media_uuid, start_time) => {
    const request = new PostWatchPointRequest(media_uuid, start_time);
    let response = await request.do();

    return response.Created;
}


/*=====  End of Watch points  ======*/


/*=============================================
=            Video moments            =
=============================================*/

    /**
     * A saved video moment
    * @typedef {Object} VideoMomentParams
     * @property {number} id
     * @property {string} video_uuid
     * @property {string} moment_title
     * @property {number} moment_time
    */

    export class VideoMoment {

        /**
         * @type {number}
         */
        #id

        /**
         * the uuid of the video media.
         * @type {string}
         */
        #video_uuid;

        /**
         * The name the video moment was saved with.
         * @type {string}
         */
        #moment_title;

        /**
         * The time of the video that hold the moment.
         * @type {number}
         */
        #moment_time;

        /**
         * @param {VideoMomentParams} param0
         */
        constructor({ id, video_uuid, moment_title, moment_time }) {
            this.#id = id;
            this.#video_uuid = video_uuid;
            this.#moment_title = moment_title;
            this.#moment_time = moment_time;
        }

        /**
         * The identifier for this video moment
         * @type {number}
         */
        get ID() {
            return this.#id;
        }

        /**
         * The start time of the moment.
         * @type {number}
         */
        get StartTime() {
            return this.#moment_time;
        }

        /**
         * The name the video moment was saved with.
         * @type {string}
         */
        get Title() {
            return this.#moment_title;
        }

        /**
         * @type {string}
         */
        get VideoUUID() {
            return this.#video_uuid;
        }
    }

    /**
     * Returns the video moments of a give video.
     * @param {string} video_uuid
     * @param {string} video_cluster
     * @returns {Promise<VideoMoment[]>}
     */
    export const getVideoMoments = async (video_uuid, video_cluster) => {
        /**
         * @type {VideoMoment[]}
         */
        const video_moments = [];

        const request = new GetVideoMomentsRequest(video_uuid, video_cluster);

        const response = await request.do();

        if (response.Ok) {
            response.data.forEach(vmp => {
                const new_video_moment = new VideoMoment(vmp);

                video_moments.push(new_video_moment);
            });
        }

        return video_moments;
    }

    /**
     * Creates a new video moment from the given parameters and saves on the metadata server.
     * @param {string} video_uuid
     * @param {string} video_cluster
     * @param {number} moment_time
     * @param {string} moment_title
     * @returns {Promise<VideoMoment | null>}
     */
    export const createVideoMoment = async (video_uuid, video_cluster, moment_time, moment_title) => {
        /**
         * @type {VideoMoment | null}
         */
        let new_video_moment = null;
        
        const request = new PostVideoMomentsRequest(video_uuid, video_cluster, moment_time, moment_title);

        const response = await request.do();

        if (response.Created) {
            new_video_moment = new VideoMoment(response.data);
        }

        return new_video_moment;
    }

/*=====  End of Video moments  ======*/



