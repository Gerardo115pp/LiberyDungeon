import {
    GetClusterTaxonomiesRequest,
    GetDungeonTagByIDRequest,
    GetDungeonTagByNameRequest,
    GetEntityTaggingsRequest,
    GetEntitiesWithTagsRequest,
    GetClusterTagsRequest,
    PostTagTaxonomyRequest,
    PostDungeonTagRequest,
    PostTagEntityRequest,
    DeleteUntagEntityRequest, 
    DeleteTagTaxonomyRequest, 
    DeleteDungeonTagRequest,
    PatchRenameTagTaxonomyRequest,
    PatchRenameDungeonTagRequest 
} from '@libs/DungeonsCommunication/services_requests/metadata_requests/dungeon_tags_requests'

/**
* @typedef {Object} TagTaxonomyParams
 * @property {string} uuid - The identifier of the taxonomy.
 * @property {string} name - The name of the taxonomy.
 * @property {string} cluster_domain - The cluster where this taxonomy is present. If it's '' then it's a global taxonomy. Otherwise the vaule should be a CategoryCluster.uuid
 * @property {string} is_internal - If the taxonomy is internal, then the user should not modify it directly.
*/

/**
* @typedef {Object} TaxonomyTagsParams
 * @property {TagTaxonomyParams} taxonomy - The taxonomy object.
 * @property {DungeonTagParams[]} tags - The tags that belong to the taxonomy.
*/

/**
* @typedef {Object} DungeonTagParams
 * @property {number} id - The identifier of the tag.
 * @property {string} name - the name of the tag. Usually the piece of data that the tag represents.
 * @property {string} taxonomy - The taxonomy the tag belongs to. It's a uuid4 string.
 * @property {string} name_taxonomy  - It's a SHA1 hash of the name and taxonomy. used to keep the tag name unique within the taxonomy.
*/

/**
* @typedef {Object} DungeonTaggingParams
 * @property {number} tagging_id - The identifier of the tagging.
 * @property {DungeonTagParams} tag - The tag object.
 * @property {string} tagged_entity_uuid - An identifier(not enforced) of a generic entity the tag is associated with.
*/

/**
 * A generic piece of data that can represent anything and be attached to any identifier string.
 */
export class DungeonTag {
    /**
     * The identifier of the tag.
     * @type {number}
     */
    #id

    /**
     * the name of the tag. Usually the piece of data that the tag represents.
     * @type {string}
     */
    #name

    /**
     * The taxonomy the tag belongs to. It's a uuid4 string.
     * @type {string}
     */
    #taxonomy

    /**
     * Name-Taxonomy. it's a SHA1 hash of the name and taxonomy. used to keep the tag name unique within the taxonomy.
     * @type {string}
     */
    #name_taxonomy

    /**
     * @param {DungeonTagParams} params
     */
    constructor(params) {
        this.#id = params.id
        this.#name = params.name
        this.#taxonomy = params.taxonomy
        this.#name_taxonomy = params.name_taxonomy
    }
}

/**
 * An attribute the tags within the taxonomy are giving context about. E.g A TagTaxonomy could be 'Color' and the tags could be 'Red', 'Blue', 'Green' and 'Orange'.
 * And there could be another TagTaxonomy called 'Fruits' with tags 'Apple', 'Banana', 'Orange' and 'Grapes'. In both cases the tag 'Orange' is present but it's a different tag. 
 * Tag name uniqueness is enforced only within it's taxonomy.
 */
export class TagTaxonomy {
    /**
     * The identifier of the taxonomy.
     * @type {string}
     */
    #uuid

    /**
     * The name of the taxonomy.
     * @type {string}
     */
    #name

    /**
     * The cluster where this taxonomy is present. If it's '' then it's a global taxonomy. Otherwise the vaule should be a CategoryCluster.uuid
     * @type {string}
     */
    #cluster_domain

    /**
     * If the taxonomy is internal, then the user should not modify it directly.
     * @type {string}
     */
    #is_internal

    /**
     * @param {TagTaxonomyParams} params
     */
    constructor(params) {
        this.#uuid = params.uuid
        this.#name = params.name
        this.#cluster_domain = params.cluster_domain
        this.#is_internal = params.is_internal
    }
}

/**
 * A composition object that represents a taxonomy and the tags that belong to it.
 */
export class TaxonomyTags {
    /**
     * The taxonomy object.
     * @type {TagTaxonomy}
     */
    #taxonomy

    /**
     * The tags that belong to the taxonomy.
     * @type {DungeonTag[]}
     */
    #tags

    /**
     * @param {TaxonomyTagsParams} params
     */
    constructor(params) {
        this.#taxonomy = new TagTaxonomy(params.taxonomy);
        this.#tags = params.tags.map(tag => new DungeonTag(tag));
    }
}

/**
 * A tag-identifier association.
 */
export class DungeonTagging {
    /**
     * The identifier of the tagging.
     * @type {number}
     */
    #tagging_id

    /**
     * The tag object.
     * @type {DungeonTag}
     */
    #tag

    /**
     * An identifier(not enforced) of a generic entity the tag is associated with.
     * @type {string}
     */
    #tagged_entity_uuid

    /**
     * @param {DungeonTaggingParams} params
     */
    constructor(params) {
        this.#tagging_id = params.tagging_id
        this.#tag = new DungeonTag(params.tag)
        this.#tagged_entity_uuid = params.tagged_entity_uuid
    }
}

/**
 * Returns all the TagTaxonomies associated with the cluster.
 * @param {string} cluster_uuid
 * @returns {Promise<TaxonomyTags[]>}
 */
export const getClusterTaxonomies = async (cluster_uuid) => {
    /**
     * @type {TagTaxonomy[]}
     */
    let cluster_taxonomies = [];

    const request = new GetClusterTaxonomiesRequest(cluster_uuid);

    const response = await request.do();
    
    if (response.Ok) {
        cluster_taxonomies = response.data.map(taxonomy_params => new TaxonomyTags(taxonomy_params));
    }

    return cluster_taxonomies;
}

/**
 * Returns a dungeon tag by its id.
 * @param {number} id
 * @returns {Promise<DungeonTag | null>}
 */
export const getDungeonTagByID = async (id) => {
    /** @type {DungeonTag | null} */
    let dungeon_tag = null;
    
    const request = new GetDungeonTagByIDRequest(id);

    const response = await request.do();

    if (response.Ok) {
        dungeon_tag = new DungeonTag(response.data);        
    }

    return dungeon_tag;
}

/**
 * Returns a dungeon tag by its name.
 * @param {string} name
 * @param {string} taxonomy
 * @returns {Promise<DungeonTag | null>}
 */
export const getDungeonTagByName = async (name, taxonomy) => {
    /** @type {DungeonTag | null} */
    let dungeon_tag = null;
    
    const request = new GetDungeonTagByNameRequest(name, taxonomy);

    const response = await request.do();

    if (response.Ok) {
        dungeon_tag = new DungeonTag(response.data);        
    }

    return dungeon_tag;
}

/**
 * Returns all of the tags associated with an entity identifier on a specific cluster.
 * @param {string} entity
 * @param {string} cluster_domain
 * @returns {Promise<DungeonTagging[]>}
 */
export const getEntityTaggings = async (entity, cluster_domain) => {
    /** @type {DungeonTagging[]} */
    let entity_taggings = [];

    const request = new GetEntityTaggingsRequest(entity, cluster_domain);

    const response = await request.do();

    if (response.Ok) {
        entity_taggings = response.data.map(tagging_params => new DungeonTagging(tagging_params));
    }

    return entity_taggings;
}

/**
 * Returns all the entity identifiers that are associated with ALL the passed tag ids.
 * Tags IDs are unique and can only be associated with one taxonomy. Same thing between taxonomies and clusters.
 * Which is a way of saying that for this operation cluster containment is 'inherited' from the taxonomy.
 * @param {number[]} tag_ids
 * @returns {Promise<string[]>}
 */
export const getEntitiesWithTags = async (tag_ids) => {
    /** @type {string[]} */
    let entities = [];

    const request = new GetEntitiesWithTagsRequest(tag_ids);

    const response = await request.do();

    if (response.Ok) {
        entities = response.data;
    }

    return entities;
}

/**
 * Returns all the tags present and unique to a cluster, as an array of TaxonomyTags.
 * @param {string} cluster_uuid
 * @returns {Promise<TaxonomyTags[]>}
 */
export const getClusterTags = async (cluster_uuid) => {
    /** @type {TaxonomyTags[]} */
    let cluster_tags = [];

    const request = new GetClusterTagsRequest(cluster_uuid);

    const response = await request.do();

    if (response.Ok) {
        cluster_tags = response.data.map(taxonomy_params => new TaxonomyTags(taxonomy_params));
    }

    return cluster_tags;
}

/**
 * Creates a new tag taxonomy.
 * @param {string} name
 * @param {string} cluster_domain
 * @param {boolean} is_internal
 * @returns {Promise<TagTaxonomy | null>}
 */
export const createTagTaxonomy = async (name, cluster_domain, is_internal) => {
    /** @type {TagTaxonomy | null} */
    let taxonomy = null;

    const new_taxonomy_uuid = crypto.randomUUID();

    /** @type {TagTaxonomyParams} */
    const new_taxonomy_params = {
        uuid: new_taxonomy_uuid,
        name,
        cluster_domain,
        is_internal,
    }

    const request = new PostTagTaxonomyRequest(new_taxonomy_params);

    const response = await request.do();

    if (response.Ok) {
        taxonomy = new TagTaxonomy(new_taxonomy_params);
    }

    return taxonomy;
}

/**
 * Creates a new dungeon tag.
 * @param {string} name -- The name/data of the tag.
 * @param {string} taxonomy -- The taxonomy uuid.
 * @returns {Promise<DungeonTag | null>}
 */
export const createDungeonTag = async (name, taxonomy) => {
    /** @type {DungeonTag | null} */
    let new_dungeon_tag = null;

    const request = new PostDungeonTagRequest(name, taxonomy);

    const response = await request.do();

    if (response.Ok) {
        new_dungeon_tag = new DungeonTag(response.data);
    }

    return new_dungeon_tag;
}

/**
 * Tags a given entity identifier with a tag by the tag id. Returns the id of the tagging created.
 * @param {string} entity
 * @param {number} tag_id
 * @returns {Promise<number | null>}
 */
export const tagEntity = async (entity, tag_id) => {
    /** @type {number | null} */
    let tagging_id = null;

    const request = new PostTagEntityRequest(entity, tag_id);

    const response = await request.do();

    if (response.Ok) {
        tagging_id = response.data.response;
    }

    return tagging_id;
}

/**
 * Removes a tag from an entity.
 * @param {string} entity
 * @param {number} tag_id
 * @returns {Promise<boolean>}
 */
export const untagEntity = async (entity, tag_id) => {
    /** @type {boolean} */
    let success = false;

    const request = new DeleteUntagEntityRequest(entity, tag_id);

    const response = await request.do();

    if (response.Ok) {
        success = response.data.response;
    }

    return success;
}

/**
 * Deletes a tag taxonomy.
 * @param {string} taxonomy_uuid
 * @returns {Promise<boolean>}
 */
export const deleteTagTaxonomy = async (taxonomy_uuid) => {
    /** @type {boolean} */
    let success = false;

    const request = new DeleteTagTaxonomyRequest(taxonomy_uuid);

    const response = await request.do();

    if (response.Ok) {
        success = response.data;
    }

    return success;
}

/**
 * Deletes a dungeon tag. This will also remove the tag from all entities and the tag will not be available for tagging anymore.
 * @param {number} id
 * @returns {Promise<boolean>}
 */
export const deleteDungeonTag = async (id) => {
    /** @type {boolean} */
    let success = false;

    const request = new DeleteDungeonTagRequest(id);

    const response = await request.do();

    if (response.Ok) {
        success = response.data;
    }

    return success;
}

/**
 * Renames a tag taxonomy.
 * @param {string} taxonomy_uuid
 * @param {string} new_name
 * @returns {Promise<boolean>}  
 */
export const renameTagTaxonomy = async (taxonomy_uuid, new_name) => {
    /** @type {boolean} */
    let success = false;

    const request = new PatchRenameTagTaxonomyRequest(taxonomy_uuid, new_name);

    const response = await request.do();

    if (response.Ok) {
        success = response.data;
    }

    return success;
}

/**
 * Renames a dungeon tag.
 * @param {number} id
 * @param {string} new_name
 * @returns {Promise<boolean>}  
 */
export const renameDungeonTag = async (id, new_name) => {
    /** @type {boolean} */
    let success = false;

    const request = new PatchRenameDungeonTagRequest(id, new_name);

    const response = await request.do();

    if (response.Ok) {
        success = response.data;
    }

    return success;
}