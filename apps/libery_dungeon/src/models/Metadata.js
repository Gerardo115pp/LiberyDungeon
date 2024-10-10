import { GetWatchPointRequest, PostWatchPointRequest } from "@libs/HttpRequests";


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

