import { GetWatchPointRequest, PostWatchPointRequest } from "@libs/DungeonsCommunication/services_requests/metadata_requests/watch_points_requests";
import { 
    GetVideoMomentsRequest,
    GetAllClusterVideoMomentsRequest,
    PostVideoMomentsRequest,
    PutVideoMomentDataRequest,
    DeleteVideoMomentRequest,
} from "@libs/DungeonsCommunication/services_requests/metadata_requests/video_moments_requests";
import { encodeVideoTime, decodeVideoTime } from "@libs/utils";


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
         * Alters either the title or the time of the this moment. the time is expected to be decoded.
         * Returns whether the operation was successful
         * @param {string} [moment_title]
         * @param {number} [moment_time]
         * @returns {Promise<boolean>}
         */
        alterMoment = async (moment_title, moment_time) => {
            if (moment_time === undefined && moment_title === undefined) {
                return true;
            }

            let new_moment_title = moment_title ?? this.#moment_title;
            let new_moment_time = moment_time ?? this.#moment_time;

            const success = await updateVideoMoment(this.#id, new_moment_title, new_moment_time);

            if (success) {
                this.#moment_time = new_moment_time;
                this.#moment_title = new_moment_title;
            }

            return success;
        }

        /**
         * The decoded moment time. this can be used to set to HTMLVideoElement.currentTime.
         * @type {number}
         */
        get DecodedTime() {
            return decodeVideoTime(this.#moment_time);
        }

        /**
         * Returns the percentage T of the given video duration at which this moment ocurres. Where 0 <= T <= 100.
         * @param {number} total_duration
         * @returns {number}
         */
        getTimelineStartPoint = total_duration => {
            const decoded_time = decodeVideoTime(this.#moment_time);

            return Math.max(0, Math.min(1, (decoded_time / total_duration))) * 100;
        }

        /**
         * The identifier for this video moment
         * @type {number}
         */
        get ID() {
            return this.#id;
        }

        /**
         * Returns whether the given time point is after this video moment.
         * @param {number} time_point
         * @returns {boolean}
         */
        isAfter = time_point => {
            const decoded_time = decodeVideoTime(this.#moment_time);

            return time_point > decoded_time;
        }

        /**
         * Returns whether the given time point ocurres before this video moment. If the video moment and the
         * time_point ocurre at the exact same time. it returns true as well.
         * @param {number} time_point
         * @returns {boolean}
         */
        isBefore = time_point => {
            const decoded_time = decodeVideoTime(this.#moment_time);

            return time_point < decoded_time && !this.isEqual(time_point);
        }

        /**
         * Returns whether a given time point is equal to this moment time but with 
         * some tolerance. meaning it will still be true if it's really close.
         * @param {number} time_point
         * @param {number} [custom_precision] - default is 0.01
         * @returns {boolean}
         */
        isEqual = (time_point, custom_precision) => {
            let precision = 0.01;

            if (custom_precision != undefined) {
                precision = custom_precision;
            }

            const decoded_time = decodeVideoTime(this.#moment_time);
            const difference = Math.abs(time_point - decoded_time)

            return difference < (decoded_time * precision);
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
     * This is a composition class for a video moment and it's corresponding media identity
     */
    export class VideoMomentIdentity {
        /**
         * The media identity.
         * @type {import('@models/Medias').MediaIdentity}
         */
        #media_identity;

        /**
         * The video moment.
         * @type {VideoMoment}
         */
        #video_moment;

        /**
         * @param {import('@models/Medias').MediaIdentity} media_identity
         * @param {VideoMoment} video_moment
         */
        constructor(media_identity, video_moment) {
            this.#media_identity = media_identity;
            this.#video_moment = video_moment;
        }

        /**
         * The media identity related to this Video moment.
         * @type {import('@models/Medias').MediaIdentity}
         */
        get MediaIdentity() {
            return this.#media_identity;
        }

        /**
         * The Video moment instance
         * @type {VideoMoment}
         */
        get VideoMoment() {
            return this.#video_moment;
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
     * Returns all video moments related to a cluster by it's uuid.
     * @param {string} cluster_uuid
     * @returns {Promise<VideoMoment[]>}
     */
    export const getAllClusterVideoMoments = async cluster_uuid => {
        /**
        * @type {VideoMoment[]}
         */
        const cluster_video_moments = [];

        const request = new GetAllClusterVideoMomentsRequest(cluster_uuid);

        const response = await request.do();

        if (response.Ok) {
            response.data.forEach(vmp => {
                cluster_video_moments.push(new VideoMoment(vmp));
            });
        }

        return cluster_video_moments;
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

    /**
     * Deletes a video moment by it's id.
     * @param {number} moment_id
     * @returns {Promise<boolean>}
     */
    export const deleteVideoMoment = async (moment_id) => {
        const request = new DeleteVideoMomentRequest(moment_id);

        const response = await request.do();

        return response.data;
    }

    /**
     * Updates the title and/or time of a video moment by it's id. returns whether the operation was successful.
     * @param {number} moment_id
     * @param {string} moment_title
     * @param {number} moment_time
     */
    export const updateVideoMoment = async (moment_id, moment_title, moment_time) => {
        const request = new PutVideoMomentDataRequest(moment_id,moment_time, moment_title);

        const response = await request.do();

        return response.data;
    }

/*=====  End of Video moments  ======*/



