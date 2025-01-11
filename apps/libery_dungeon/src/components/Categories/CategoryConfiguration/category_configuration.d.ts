import { InnerCategory } from "@models/Categories"
import { MediaIdentity } from "@models/Medias";

export type CategoryConfig_ThumbnailChanged = (updated_category: InnerCategory) => void;

export type CategoryConfig_BillboardMediaAdded = (media_identity: MediaIdentity) => void;

export type CategoryConfig_BillboardMediaRemoved = (media_uuid: string) => void;