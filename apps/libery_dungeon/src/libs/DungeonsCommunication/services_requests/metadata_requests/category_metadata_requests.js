import { HttpResponse, arrayToParam, attributesToJson } from "../../base";
import { metadata_server } from "../../services"

/**
 * Fetches a Category config object from the server.
 */
export class GetCategoryConfigRequest {

    static endpoint = `${metadata_server}/categories-metadata/category-config`

    /**
     * @param {string} category_uuid
     */
    constructor(category_uuid) {
        this.category_uuid = category_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/Categories").CategoryConfigParams>>}
     */
    do = async () => {
        if (!URL.canParse(GetCategoryConfigRequest.endpoint, globalThis.location.origin)) {
            throw new Error("In @libs/DungeonsCommunication/services_requests/metadata_requests/category_metadata_requests.GetCategoryConfigRequest.do: There is an issue with the URL syntax. Likely a typo.");
        }
        
        const url = new URL(GetCategoryConfigRequest.endpoint, globalThis.location.origin);

        url.searchParams.append("category_uuid", this.category_uuid);

        /**
         * @type {import("@models/Categories").CategoryConfigParams}
         */
        let response_config = {
            category_uuid: "",
            billboard_dungeon_tags: [],
            billboard_media_uuids: [],
        }

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                response_config = await response.json();
            }
        } catch (error) {
            throw new Error(`In @libs/DungeonsCommunication/services_requests/metadata_requests/category_metadata_requests.GetCategoryConfigRequest.do: server error\n\n${error}`);
        }

        return new HttpResponse(response, response_config);
    }
}

/**
 * Patches the category billboard media list with a new list of media uuids. Requires the user to have content_alter grant.
 */
export class PatchCategoryBillboardMediasRequest {

    static endpoint  = `${metadata_server}/categories-metadata/category-config/billboard-medias`;

    /**
     * @param {string} category_uuid
     * @param {string[]} billboard_media_uuids
     */
    constructor(category_uuid, billboard_media_uuids) {
        this.category_uuid = category_uuid;
        this.billboard_media_uuids = billboard_media_uuids;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("../../base").BooleanResponse>>}
     */
    do = async () => {
        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {import('../../base').BooleanResponse}
         */
        let boolean_response = { response: false };
        

        try {
            response = await fetch(PatchCategoryBillboardMediasRequest.endpoint, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                },
                body: this.toJson(),
            });

            boolean_response.response = response.ok;
        } catch (error) {
            throw new Error(`In @libs/DungeonsCommunication/services_requests/metadata_requests/category_metadata_requests.PatchCategoryBillboardMediasRequest.do: server error\n\n${error}`);
        }

        return new HttpResponse(response, boolean_response);
    }
}

/**
 * Patch the category billboard tags with a new list of tag ids. Requires the user to have content_alter grant.
 */
export class PatchCategoryBillboardTagsRequest {

    static endpoint = `${metadata_server}/categories-metadata/category-config/billboard-tags`;

    /**
     * @param {string} category_uuid
     * @param {number[]} billboard_dungeon_tags
     */
    constructor(category_uuid, billboard_dungeon_tags) {
        this.category_uuid = category_uuid;
        this.billboard_dungeon_tags = billboard_dungeon_tags;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("../../base").BooleanResponse>>}
     */
    do = async () => {
        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {import('../../base').BooleanResponse}
         */
        let boolean_response = { response: false };

        try {
            response = await fetch(PatchCategoryBillboardTagsRequest.endpoint, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                },
                body: this.toJson(),
            });

            boolean_response.response = response.ok;
        } catch (error) {
            throw new Error(`In @libs/DungeonsCommunication/services_requests/metadata_requests/category_metadata_requests.PatchCategoryBillboardTagsRequest.do: server error\n\n${error}`);
        }

        return new HttpResponse(response, boolean_response);
    }
}