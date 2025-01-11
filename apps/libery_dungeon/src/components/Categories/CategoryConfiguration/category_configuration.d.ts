import { CategoryConfig, InnerCategory } from "@models/Categories"
import type { DungeonTag } from "@models/DungeonTags";
import { MediaIdentity } from "@models/Medias";

export type CategoryConfig_ThumbnailChanged = (updated_category: InnerCategory) => void;

export type CategoryConfig_BillboardMediaAdded = (media_identity: MediaIdentity) => void;

export type CategoryConfig_BillboardMediaRemoved = (media_uuid: string) => void;

export type CategoryConfig_BillboardDungeonTagsAdded = (dungeon_tag: DungeonTag) => void;

export type CategoryConfig_BillboardDungeonTagsRemoved = (tag_id: number) => void;

export type CategoryConfig_CategoryConfigChanged = (category_config: CategoryConfig) => void;