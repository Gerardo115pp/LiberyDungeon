import { HttpResponse, attributesToJson } from "../../base";
import { metadata_server } from "../../services"

/**
 * Returns all of the TagTaxonomies associated with a cluster
 */
export class GetClusterTaxonomiesRequest {

    static endpoint = `${metadata_server}/dungeon-tags/taxonomies/cluster`;

    /**
     * @param {string} cluster_uuid 
     */
    constructor(cluster_uuid) {
        this.cluster_uuid = cluster_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").TagTaxonomyParams[]>>}
     */
    do = async () => {
        const url = `${GetClusterTaxonomiesRequest.endpoint}?cluster_uuid=${this.cluster_uuid}`;

        /** @type {import("@models/DungeonTags").TagTaxonomyParams[]} */
        let cluster_taxonomies = [];

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                cluster_taxonomies = await response.json();
            }

        } catch (error) {
            console.error("Error getting cluster taxonomies: ", error);
        }

        return new HttpResponse(response, cluster_taxonomies);
    }

}

/**
 * Returns a dungeon tag by its id.
 */
export class GetDungeonTagByIDRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tags`;

    /**
     * @param {number} id
     */
    constructor(id) {
        this.id = id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").DungeonTagParams | null>>}
     */
    do = async () => {
        const url = `${GetDungeonTagByIDRequest.endpoint}?id=${this.id}`;

        /** @type {import("@models/DungeonTags").DungeonTagParams} */
        let dungeon_tag = null;

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                dungeon_tag = await response.json();
            }

        } catch (error) {
            console.error("Error getting dungeon tag: ", error);
        }

        return new HttpResponse(response, dungeon_tag);
    }
}

/**
 * Returns a dungeon tag by its name and taxonomy.
 */
export class GetDungeonTagByNameRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tags`;

    /**
     * @param {string} name
     * @param {string} taxonomy
     */
    constructor(name, taxonomy) {
        this.name = name;
        this.taxonomy = taxonomy;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").DungeonTagParams | null>>}
     */
    do = async () => {
        const url = `${GetDungeonTagByNameRequest.endpoint}?name=${this.name}&taxonomy=${this.taxonomy}`;

        /** @type {import("@models/DungeonTags").DungeonTagParams} */
        let dungeon_tag = null;

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                dungeon_tag = await response.json();
            }

        } catch (error) {
            console.error("Error getting dungeon tag: ", error);
        }

        return new HttpResponse(response, dungeon_tag);
    }
}

/**
 * Returns all of the tags associated with an entity identifier on a specific cluster.
 */
export class GetEntityTaggingsRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tags/entity`;

    /**
     * @param {string} entity
     * @param {string} cluster_domain
     */
    constructor(entity, cluster_domain) {
        this.entity = entity;
        this.cluster_domain = cluster_domain;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").DungeonTaggingParams[]>>}
     */
    do = async () => {
        const url = `${GetEntityTaggingsRequest.endpoint}?entity=${this.entity}&cluster_domain=${this.cluster_domain}`;

        /** @type {import("@models/DungeonTags").DungeonTaggingParams[]} */
        let entity_taggings = [];

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                entity_taggings = await response.json();
            }

        } catch (error) {
            console.error("Error getting entity taggings: ", error);
        }

        return new HttpResponse(response, entity_taggings);
    }
}

/**
 * Returns all the entity identifiers that are associated with ALL the passed tag ids.
 * Tags IDs are unique and can only be associated with one taxonomy. Same thing between taxonomies and clusters. 
 * Which is a way of saying that for this operation cluster containment is 'inherited' from the taxonomy.
 */
export class GetEntitiesWithTagsRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tags/matching-entities`;

    /**
     * @param {number[]} tag_ids
     */
    constructor(tag_ids) {
        if (tag_ids.length == 0) {
            throw new Error("At least one tag id must be passed");
        }

        this.tag_ids = `${tag_ids[0]}`;

        for (let h = 1; h < tag_ids.length; h++) {
            this.tag_ids += `,${tag_ids[h]}`;
        }
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<string[]>>}
     */
    do = async () => {
        const url = `${GetEntitiesWithTagsRequest.endpoint}?tags=${this.tag_ids}`;

        /** @type {string[]} */
        let entities = [];

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                entities = await response.json();
            }

        } catch (error) {
            console.error("Error getting entities with tags: ", error);
        }

        return new HttpResponse(response, entities);
    }
}

/**
 * Returns all the tags present and unique to a cluster, as an array of TaxonomyTagsParams.
 */
export class GetClusterTagsRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tags/cluster`;

    /**
     * @param {string} cluster_uuid
     */
    constructor(cluster_uuid) {
        this.cluster_uuid = cluster_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").TaxonomyTagsParams[]>>}
     */
    do = async () => {
        const url = `${GetClusterTagsRequest.endpoint}?cluster_uuid=${this.cluster_uuid}`;

        /** @type {import("@models/DungeonTags").TaxonomyTagsParams[]} */
        let cluster_tags = [];

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                cluster_tags = await response.json();
            }

        } catch (error) {
            console.error("Error getting cluster tags: ", error);
        }

        return new HttpResponse(response, cluster_tags);
    }
}

/**
 * Returns all the tags present and unique to a cluster that were defined by the user(non-internal), as an array of TaxonomyTagsParams.
 */
export class GetClusterUserTagsRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tags/user-defined/cluster`;

    /**
     * @param {string} cluster_uuid
     */
    constructor(cluster_uuid) {
        this.cluster_uuid = cluster_uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").TaxonomyTagsParams[]>>}
     */
    do = async () => {
        const url = `${GetClusterUserTagsRequest.endpoint}?cluster_uuid=${this.cluster_uuid}`;

        /** @type {import("@models/DungeonTags").TaxonomyTagsParams[]} */
        let cluster_tags = [];

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                cluster_tags = await response.json();
            }

        } catch (error) {
            console.error("Error getting cluster tags: ", error);
        }

        return new HttpResponse(response, cluster_tags);
    }
}

/**
 * Returns a TaxonomyTags corresponding to the given taxonomy id
 */
export class GetTaxonomyTagsRequest {

    static endpoint = `${metadata_server}/dungeon-tags/taxonomies/tags`;

    /**
     * @param {string} taxonomy
     */
    constructor(taxonomy) {
        this.taxonomy = taxonomy;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").TaxonomyTagsParams | null>>}
     */
    do = async () => {
        const url = `${GetTaxonomyTagsRequest.endpoint}?taxonomy=${this.taxonomy}`;

        /** @type {import("@models/DungeonTags").TaxonomyTagsParams} */
        let taxonomy_tags = null;

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                taxonomy_tags = await response.json();
            }

        } catch (error) {
            console.error("Error getting taxonomy tags: ", error);
        }

        return new HttpResponse(response, taxonomy_tags);
    }
}

/**
 * Creates a TagTaxonomy.
 */
export class PostTagTaxonomyRequest {

    static endpoint = `${metadata_server}/dungeon-tags/taxonomies`;

    /**
     * @param {import("@models/DungeonTags").TagTaxonomyParams} param0
     */
    constructor({ uuid, name, cluster_domain, is_internal }) {
        this.uuid = uuid;
        this.name = name;
        this.cluster_domain = cluster_domain;
        this.is_internal = is_internal;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        const url = PostTagTaxonomyRequest.endpoint;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {boolean}
         */
        let created = false;

        try {
            response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

        } catch (error) {
            console.error("Error posting tag taxonomy: ", error);
        }

        created = response?.status === 201;

        return new HttpResponse(response, created);
    }
}

/**
 * Creates a DungeonTag.
 */
export class PostDungeonTagRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tags`;

    /**
     * @param {string} name
     * @param {string} taxonomy
     */
    constructor(name, taxonomy) {
        this.name = name;
        this.taxonomy = taxonomy;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").DungeonTagParams | null>>}
     */
    do = async () => {
        /** @type {import("@models/DungeonTags").DungeonTagParams | null} */
        let new_dungeon_tag_params = null;

        const url = PostDungeonTagRequest.endpoint;

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

        } catch (error) {
            console.error("Error posting dungeon tag: ", error);
        }

        if (response?.status === 201) {
            new_dungeon_tag_params = await response.json();
        }

        return new HttpResponse(response, new_dungeon_tag_params);
    }
}

/**
 * Tags a given entity identifier with a tag by the tag id.
 */
export class PostTagEntityRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tag-entity`;

    /**
     * @param {string} entity
     * @param {number} tag_id
     */
    constructor(entity, tag_id) {
        this.entity = entity;
        this.tag_id = tag_id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * Returns a single number response, with the response being the tagging id.
     * @returns {Promise<HttpResponse<import("../../base").SingleNumberResponse>>}
     */
    do = async () => {
        const url = `${PostTagEntityRequest.endpoint}?entity=${this.entity}&tag_id=${this.tag_id}`;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {import("../../base").SingleNumberResponse}
         */
        let single_number_response = { response: null };

        try {
            response = await fetch(url, {
                method: "POST",
            });

        } catch (error) {
            console.error("Error tagging entity: ", error);
        }

        if (response?.status === 201) {
            single_number_response = await response.json();
        }

        return new HttpResponse(response, single_number_response);
    }
}

/**
 * Removes a tag from an entity identifier. Takes the same parameters as PostTagEntityRequest.
 */
export class DeleteUntagEntityRequest {

        static endpoint = `${metadata_server}/dungeon-tags/tag-entity`;

        /**
         * @param {string} entity
         * @param {number} tag_id
         */
        constructor(entity, tag_id) {
            this.entity = entity;
            this.tag_id = tag_id;
        }

        toJson = attributesToJson.bind(this);

        /**
         * @returns {Promise<HttpResponse<boolean>>}
         */
        do = async () => {
            const url = `${DeleteUntagEntityRequest.endpoint}?entity=${this.entity}&tag_id=${this.tag_id}`;

            /**
             * @type {Response}
             */
            let response;

            /**
             * @type {boolean}
             */
            let untagged = false;

            try {
                response = await fetch(url, {
                    method: "DELETE"
                });

            } catch (error) {
                console.error("Error untagging entity: ", error);
            }

            untagged = response?.status === 204;

            return new HttpResponse(response, untagged);
        }
}

/**
 * Deletes a TagTaxonomy.
 */
export class DeleteTagTaxonomyRequest {

    static endpoint = `${metadata_server}/dungeon-tags/taxonomies`;

    /**
     * @param {string} uuid
     */
    constructor(uuid) {
        this.uuid = uuid;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        const url = `${DeleteTagTaxonomyRequest.endpoint}?uuid=${this.uuid}`;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {boolean}
         */
        let deleted = false;

        try {
            response = await fetch(url, {
                method: "DELETE"
            });

        } catch (error) {
            console.error("Error deleting tag taxonomy: ", error);
        }

        deleted = response?.status === 204;

        return new HttpResponse(response, deleted);
    }
}

/**
 * Deletes a DungeonTag.
 */
export class DeleteDungeonTagRequest {

    static endpoint = `${metadata_server}/dungeon-tags/tags`;

    /**
     * @param {number} id
     */
    constructor(id) {
        this.id = id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        const url = `${DeleteDungeonTagRequest.endpoint}?id=${this.id}`;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {boolean}
         */
        let deleted = false;

        try {
            response = await fetch(url, {
                method: "DELETE"
            });

        } catch (error) {
            console.error("Error deleting dungeon tag: ", error);
        }

        deleted = response?.status === 204;

        return new HttpResponse(response, deleted);
    }
}

/**
 * Renames a TagTaxonomy.
 */
export class PatchRenameTagTaxonomyRequest {

    static endpoint = `${metadata_server}/dungeon-tags/taxonomies/name`;

    /**
     * @param {string} uuid
     * @param {string} new_name
     */
    constructor(uuid, new_name) {
        this.uuid = uuid;
        this.new_name = new_name;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        const url = PatchRenameTagTaxonomyRequest.endpoint;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {boolean}
         */
        let renamed = false;

        try {
            response = await fetch(url, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });
        } catch (error) {
            console.error("Error renaming tag taxonomy: ", error);
        }

        renamed = response?.status === 204;

        return new HttpResponse(response, renamed);
    }
}

/**
 * Renames a DungeonTag.
 */
export class PatchRenameDungeonTagRequest {
    
    static endpoint = `${metadata_server}/dungeon-tags/tags/name`;

    /**
     * @param {number} id
     * @param {string} new_name
     */
    constructor(id, new_name) {
        this.id = id;
        this.new_name = new_name;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        const url = PatchRenameDungeonTagRequest.endpoint;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {boolean}
         */
        let renamed = false;

        try {
            response = await fetch(url, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });
        } catch (error) {
            console.error("Error renaming dungeon tag: ", error);
        }

        renamed = response?.status === 204;

        return new HttpResponse(response, renamed);
    }
}