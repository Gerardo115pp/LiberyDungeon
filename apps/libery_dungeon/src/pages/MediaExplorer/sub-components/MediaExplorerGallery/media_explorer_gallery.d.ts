import type { OrderedMedia } from "@models/Medias";

export type ObserveMEGalleryCallback = (item: Element, media_item: OrderedMedia) => void;
export type UnobserveMEGalleryCallback = (item: Element, media_item: OrderedMedia) => void;

