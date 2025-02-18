import { HttpResponse, arrayToParam, attributesToJson } from "../../base";
import { metadata_server, categories_server } from "../../services"

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

            throw error;
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

        /** @type {import("@models/DungeonTags").DungeonTagParams | null} */
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

            throw error;
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

        /** @type {import("@models/DungeonTags").DungeonTagParams | null} */
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

            throw error;
        }

        return new HttpResponse(response, dungeon_tag);
    }
}

/**
 * Returns a list of dungeon tags matching a list of tag ids.
 */
export class GetDungeonTagsRequest {
    
    static endpoint = `${metadata_server}/dungeon-tags/tags/by-ids`;

    /**
     * @param {number[]} tags
     */
    constructor(tags) {
        this.tags = tags;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@models/DungeonTags").DungeonTagParams[]>>}
     */
    do = async () => {
        /**
         * @type {import("@models/DungeonTags").DungeonTagParams[]}
         */
        let matching_tags = [];

        let tags_param = this.tags.map(t => t.toString()).join(',')

        const url = new URL(GetDungeonTagsRequest.endpoint, globalThis.location.origin);

        url.searchParams.append("tags", tags_param);

        /**
         * @type {Response}
         */
        let response;

        try {
            response = await fetch(url);

            if (response.ok) {
                matching_tags = await response.json();
            }
        } catch (err) {
            console.error("Error getting dungeon tags: ", err);

            throw err;
        }

        return new HttpResponse(response, matching_tags);
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

            throw error;    
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

            throw error;
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

            throw error;
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

            throw error;
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

        /** @type {import("@models/DungeonTags").TaxonomyTagsParams | null} */
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

            throw error;
        }

        return new HttpResponse(response, taxonomy_tags);
    }
}

/**
 * Tags a given entity with a list of tags provided as a list of tag ids(number[])
 */
export class PostMultiTagEntity {

    static endpoint = `${metadata_server}/dungeon-tags/multi-tag-entity`;

    /**
     * @param {string} entity_uuid
     * @param {number[]} dungeon_tags
     * @param {string} entity_type
     */
    constructor(entity_uuid, dungeon_tags, entity_type) {
        this.entity_uuid = entity_uuid;
        this.dungeon_tags = dungeon_tags;
        this.entity_type = entity_type;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('../../base').BooleanResponse>>}
     */
    do = async () => {
        const url = PostMultiTagEntity.endpoint;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {import('../../base').BooleanResponse}
         */
        let tagged = { response: false };

        try {
            response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

        } catch (error) {
            console.error("Error tagging entity: ", error);

            throw error;
        }

        if (response.ok) {
            tagged = await response.json();
        }

        return new HttpResponse(response, tagged);
    }
}

/**
 * Tags a given list of entities of the same type, with a list of tags provided as a list of tag ids(number[])
 */
export class PostMultiTagEntities {

    static endpoint = `${metadata_server}/dungeon-tags/multi-tag-entities`;

    /**
     * @param {string[]} entity_uuids
     * @param {number[]} dungeon_tags
     * @param {string} entity_type
     */
    constructor(entity_uuids, dungeon_tags, entity_type) {
        this.entity_uuids = entity_uuids;
        this.dungeon_tags = dungeon_tags;
        this.entity_type = entity_type;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import('../../base').BooleanResponse>>}
     */
    do = async () => {
        const url = PostMultiTagEntities.endpoint;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {import('../../base').BooleanResponse}
         */
        let tagged = { response: false };

        try {
            response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

        } catch (error) {
            console.error("Error tagging entity: ", error);

            throw error;
        }

        if (response.ok) {
            tagged = await response.json();
        }

        return new HttpResponse(response, tagged);
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

            throw error;
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

            throw error;
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
     * @param {string} entity_type
     * @param {number} tag_id
     */
    constructor(entity, entity_type, tag_id) {
        if (entity === "" || entity_type == null) {
            throw new Error("Entity and entity_type must be defined");
        }
        this.entity = entity;

        if (entity_type === "" || entity_type == null) {
            throw new Error("Entity type must be defined");
        }
        this.entity_type = entity_type;

        if (tag_id == null || tag_id < 0) {
            throw new Error("Tag id must be defined and positive");
        }
        this.tag_id = tag_id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * Returns a single number response, with the response being the tagging id.
     * @returns {Promise<HttpResponse<import("../../base").SingleNumberResponse>>}
     */
    do = async () => {
        const url = `${PostTagEntityRequest.endpoint}?entity=${this.entity}&tag_id=${this.tag_id}&entity_type=${this.entity_type}`;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {import("../../base").SingleNumberResponse}
         */
        let single_number_response = { response: -1 };

        try {
            response = await fetch(url, {
                method: "POST",
            });

        } catch (error) {
            console.error("Error tagging entity: ", error);

            throw error;
        }

        if (response?.status === 201) {
            single_number_response = await response.json();
        }

        return new HttpResponse(response, single_number_response);
    }
}

/**
 * Applies a tag to a list of entities.
 */
export class PostTagEntitiesRequest {
    
    static endpoint = `${metadata_server}/dungeon-tags/tag-entities`;

    /**
     * @param {number} tag_id
     * @param {string} entity_type
     * @param {string[]} entities_uuids
     */
    constructor(tag_id, entity_type, entities_uuids) {
        this.tag_id = tag_id;
        this.entity_type = entity_type;
        this.entities_uuids = entities_uuids;
    }

    toJson = attributesToJson.bind(this);

    /**
     *  @returns {Promise<HttpResponse<import("../../base").BooleanResponse>>}
     */
    do = async () => {
        const url = PostTagEntitiesRequest.endpoint;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {import("../../base").BooleanResponse}
         */
        let tagged = { response: false };

        try {
            response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: this.toJson()
            });

        } catch (error) {
            console.error("Error tagging entities: ", error);

            throw error;
        }

        if (response.ok) {
            tagged = await response.json();
        }

        return new HttpResponse(response, tagged);
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

                throw error;
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

            throw error;
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

            throw error;
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
        const url = `${PatchRenameTagTaxonomyRequest.endpoint}?uuid=${this.uuid}&new_name=${this.new_name}`;

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
            });
        } catch (error) {
            console.error("Error renaming tag taxonomy: ", error);

            throw error;
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
        const url = `${PatchRenameDungeonTagRequest.endpoint}?id=${this.id}&new_name=${this.new_name}`;

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
            });
        } catch (error) {
            console.error("Error renaming dungeon tag: ", error);

            throw error;
        }

        renamed = response?.status === 204;

        return new HttpResponse(response, renamed);
    }
}

/**
 * Tags a the medias inside a category with a given tag.
 */
export class PostTagCategoryContentRequest {

    static endpoint = `${categories_server}/categories/tags/content`;

    /**
     * @param {string} category_uuid
     * @param {number} tag_id
     */
    constructor(category_uuid, tag_id) {
        this.category_uuid = category_uuid;
        this.tag_id = tag_id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        const url = `${PostTagCategoryContentRequest.endpoint}?category_uuid=${this.category_uuid}&tag_id=${this.tag_id}`;

        /**
         * @type {Response}
         */
        let response;

        /**
         * @type {boolean}
         */
        let tagged = false;

        try {
            response = await fetch(url, {
                method: "POST",
            });

        } catch (error) {
            console.error("Error tagging category content: ", error);

            throw error;
        }

        tagged = response?.status === 200;

        return new HttpResponse(response, tagged);
    }
}

/**
 * Untags a the medias inside a category with a given tag.
 */
export class DeleteUntagCategoryContentRequest {

    static endpoint = `${categories_server}/categories/tags/content`;

    /**
     * @param {string} category_uuid
     * @param {number} tag_id
     */
    constructor(category_uuid, tag_id) {
        this.category_uuid = category_uuid;
        this.tag_id = tag_id;
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<boolean>>}
     */
    do = async () => {
        const url = `${DeleteUntagCategoryContentRequest.endpoint}?category_uuid=${this.category_uuid}&tag_id=${this.tag_id}`;

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
                method: "DELETE",
            });

        } catch (error) {
            console.error("Error untagging category content: ", error);

            throw error;
        }

        untagged = response?.status === 200;

        return new HttpResponse(response, untagged);
    }
}

/**
 * Returns all the medias that have been tagged with a list of tag ids
 */
export class GetContentTaggedRequest {

    static endpoint = `${categories_server}/categories/tags/content-tagged`;

    /**
     * @param {number[]} tag_ids
     * @param {number} [page]
     * @param {number} [page_size]
     */
    constructor(tag_ids, page, page_size) {
        if (page == null) {
            page = 1;
        }

        if (page_size == null) {
            page_size = 100;
        }

        this.page = page;
        this.page_size = page_size;
        this.tags = arrayToParam(tag_ids);
    }

    toJson = attributesToJson.bind(this);

    /**
     * @returns {Promise<HttpResponse<import("@libs/DungeonsCommunication/dungeon_communication").PaginatedResponse<import('@models/Medias').MediaIdentityParams> | null>>}
     */
    do = async () => {
        const url = new URL(`${GetContentTaggedRequest.endpoint}`, globalThis.location.origin);

        url.searchParams.append("tags", this.tags);
        url.searchParams.append("page", this.page.toString());
        url.searchParams.append("page_size", this.page_size.toString());

        /**
         * @type {import('@libs/DungeonsCommunication/dungeon_communication').PaginatedResponse<import('@models/Medias').MediaIdentityParams> | null}
         */
        let tagged_content = null;

        /**
         * @type {Response | null}
         */
        let response = null;

        try {
            response = await fetch(url);

            if (response.ok) {
                tagged_content = await response.json();
            }
        } catch (error) {
            console.error("Error getting tagged content: ", error);

            throw error;
        }

        return new HttpResponse(response, tagged_content);
    }
}