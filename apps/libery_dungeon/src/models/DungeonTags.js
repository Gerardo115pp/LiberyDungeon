
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
