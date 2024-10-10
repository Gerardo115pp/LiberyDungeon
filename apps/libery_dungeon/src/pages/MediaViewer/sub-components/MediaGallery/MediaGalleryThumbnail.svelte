<script>
    import LazyLoader from "@components/LazyLoader/LazyLoader.svelte";
    import { getMediaUrl } from "@libs/HttpRequests";
    import { createEventDispatcher } from "svelte";
    import { media_changes_manager, active_media_change, active_media_index } from "@stores/media_viewer";
    import { media_change_types } from "@models/WorkManagers";
    import { current_category } from "@stores/categories_tree";
    import { onMount } from "svelte";


    const event_dispatcher = createEventDispatcher();

    /** @type{import('@models/Medias').Media} the media to be displayed */
    export let media;

    /** @type{string} the current category path */
    export let category_path;

    /** @type{boolean} whether the deletion mode is enabled or not */
    export let deletion_mode = false;

    /** @type{"Moved" | "Deleted" | "Normal"} the current media change type */
    let media_current_change = media_change_types.NORMAL;

    onMount(() => {
        media_current_change = $media_changes_manager.getMediaChangeType(media.uuid);
    })

    const deleteMedia = () => {
        if (media_current_change === media_change_types.DELETED) {
            $media_changes_manager.unstageMediaDeletion(media.uuid);
            media_current_change = media_change_types.NORMAL;
            return;
        }

        media_current_change = media_change_types.DELETED;
        $media_changes_manager.stageMediaDeletion(media);

        let active_media = $current_category?.content[$active_media_index];
        if (active_media?.uuid === media.uuid) {
            active_media_change.set(media_current_change);
        }
    };

    const emitImageSelected = e => {
        if (deletion_mode) {
            return deleteMedia();
        }

        event_dispatcher("image-selected", media);
    }
</script>

<div class="mg-thumbnail-wrapper" class:deleted-media={media_current_change === media_change_types.DELETED} on:click|stopPropagation={emitImageSelected}>
    <LazyLoader className="mg-gg-media-loader" image_url={getMediaUrl(category_path, media.name, true)}>
        <img slot="lazy-wrapper-image" let:image_src src={image_src} alt="media thumbnail"/>
    </LazyLoader>
</div>

<style>
    .mg-thumbnail-wrapper {
        width: 100%;
        height: 100%;
        position: relative;
        overflow: hidden;
    }

    .mg-thumbnail-wrapper.deleted-media {
        filter: sepia(60%) hue-rotate(320deg) saturate(6) brightness(1.2);
    }

    img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform .2s ease-in-out;
    }

    img:hover {
        transform: scale(1.1);
    }
</style>