<script>
    import LazyLoader from "@components/LazyLoader/LazyLoader.svelte";
    import { getMediaUrl } from "@libs/DungeonsCommunication/services_requests/media_requests";
    import { createEventDispatcher } from "svelte";
    import { media_changes_manager, active_media_change, active_media_index } from "@pages/MediaViewer/app_page_store";
    import { media_change_types } from "@models/WorkManagers";
    import { current_category } from "@stores/categories_tree";
    import { current_cluster } from "@stores/clusters";
    import { onMount } from "svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/

        /**
         * @type {import('@models/Medias').Media} the media to be displayed
         */
        export let media;

        /**
         * @type {string} the current category path
         */
        export let category_path;

        /**
         * @type {boolean} whether the deletion mode is enabled or not 
         */
        export let deletion_mode = false;

        /** 
         * @type {import('@models/WorkManagers').MediaChangeType} the current media change type
         */
        let media_current_change = media_change_types.NORMAL;

        const event_dispatcher = createEventDispatcher();
    
    
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        media_current_change = $media_changes_manager.getMediaChangeType(media.uuid);
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
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

        /**
         * Handles the click event on the div.mg-thumbnail-wrapper element
         * @param {MouseEvent} event
         */
        const emitImageSelected = event => {
            if (deletion_mode) {
                return deleteMedia();
            }

            event_dispatcher("image-selected", {media});
        }
    
    
    /*=====  End of Methods  ======*/

</script>

<div class="mg-thumbnail-wrapper" class:deleted-media={media_current_change === media_change_types.DELETED} on:click|stopPropagation={emitImageSelected}>
    <LazyLoader className="mg-gg-media-loader" image_url={getMediaUrl(category_path, $current_cluster.UUID, media.name, true)}>
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