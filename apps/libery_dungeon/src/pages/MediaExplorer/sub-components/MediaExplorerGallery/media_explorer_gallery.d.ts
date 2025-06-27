import type { OrderedMedia } from "@models/Medias";
import type { MediaChangeType } from "@models/WorkManagers";
import { meg_media_change_state_event_names } from "./me_gallery_state";

export type ObserveMEGalleryCallback = (item: Element, media_item: OrderedMedia) => void;
export type UnobserveMEGalleryCallback = (item: Element, media_item: OrderedMedia) => void;

declare global {
    interface HTMLElementEventMap {
        [meg_media_change_state_event_names.ALTERED_MEDIA_CHANGE_STATE]: CustomEvent<MediaChangeType>;
    }
}
