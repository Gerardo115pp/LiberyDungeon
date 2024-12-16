import { DungeonTag } from "@models/DungeonTags"

/**
 * A subset of tag from a specific tag taxonomy.
 */
export type DungeonTags_GroupedTags<T> = {
    taxonomy_name: string,
    taxonomy_uuid: string,
    tags: DungeonTag[],
    component_ref: T | null,
}