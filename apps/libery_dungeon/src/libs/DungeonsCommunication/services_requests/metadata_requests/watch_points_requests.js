import { metadata_server } from "../../services";
import { HttpResponse, attributesToJson } from "../../base";

export class GetWatchPointRequest {
    /**
     * @param {string} media_uuid 
     */
    constructor(media_uuid) {
        this.media_uuid = media_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<WatchPointTime>>}
     * @typedef {Object} WatchPointTime
     * @property {number} start_time 
     */
    do = async () => {
        const response = await fetch(`${metadata_server}/watch-points?media_uuid=${this.media_uuid}`);

        /**
         * @type {WatchPointTime}
         */
        let data = {
            start_time: NaN
        };
        let http_response = new HttpResponse(response, data);

        if (!(response.status <= 200 && response.status < 300) || !response.headers.has("Content-Type")) {
            console.error("Error getting watch point: ", response);
            return http_response;
        }

        switch (response.headers.get("Content-Type")) {
            case "application/json":
                data = await response.json();
                http_response = new HttpResponse(response, data);
                break;
            case "application/octet-stream":
                let response_data = {};
                let array_data = await response.arrayBuffer();
                let data_view = new DataView(array_data);
                response_data.start_time = data_view.getUint32(0, true);
                http_response = new HttpResponse(response, response_data);
                break;
        }

        return http_response;
    }
}

export class PostWatchPointRequest {

    /**
     * @param {string} media_uuid 
     * @param {number} start_time 
     */
    constructor(media_uuid, start_time) {
        this.media_uuid = media_uuid;
        this.start_time = start_time;
    }

    toJson = attributesToJson.bind(this);

    do = async () => {
        const response = await fetch(`${metadata_server}/watch-points`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: this.toJson()
        });

        return new HttpResponse(response, {});
    }
}