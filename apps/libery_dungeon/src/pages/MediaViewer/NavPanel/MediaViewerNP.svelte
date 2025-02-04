<script>
    import { current_category } from "@stores/categories_tree";
    import { media_change_types } from "@models/WorkManagers";
    import { 
        active_media_index, 
        shared_active_media,
        active_media_change, 
        random_media_navigation, 
        skip_deleted_medias, 
        media_changes_manager,
        auto_move_category,
        auto_move_on
    } from "@pages/MediaViewer/app_page_store";
    import { 
        mv_tag_mode_enabled, 
        mv_tagged_content,
        active_tag_content_media_index,

        mv_tag_mode_total_content

    } from "@pages/MediaViewer/features_wrappers/media_viewer_tag_mode";
    import DownloadIcon from "@components/UI/icons/download_icon.svelte";

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * get the new category name if a media has been moved
         * @param {string} media_uuid
         * @param {import('@models/WorkManagers').MediaChangeType} current_change
         */
        const getMediaNewCategory = (media_uuid, current_change) => {
            if ($current_category == null) {
                console.error("In MediasViewerNP.getMediaNewCategory: current_category is null");
                return;
            }

            console.log("getMediaNewCategory called");
            if (current_change !== media_change_types.MOVED) return;

            let displayed_media = $current_category.content[$active_media_index];
            
            let new_category = $media_changes_manager.getMediaNewCategory(media_uuid);
            
            
            let category_name = new_category?.name || $current_category?.name;

            console.log(`media_uuid: ${media_uuid}\ncurrent_change: ${current_change}\ndisplayed_media: ${JSON.stringify(displayed_media)}\nactive_media_index: ${$active_media_index}\nactive_media_change: ${$active_media_change}\nnew_category: ${JSON.stringify(new_category)}\ncategory_name: ${category_name}`);

            return category_name;
        }
            
    /*=====  End of Methods  ======*/ 

</script>

<ul class="page-nav-menu" id="media-viewer-navmenu-wrapper" class:adebug={false}>
    <li class="mvnw-category-name">
        {#if $mv_tag_mode_enabled}
            <span class="category-name">
                filtered content
            </span>
        {:else if $current_category !== null}
            {#if $active_media_change !== media_change_types.MOVED}
                <span class="category-name">{$current_category.name}</span>
            {:else if $active_media_index < $current_category.content.length && $active_media_index >= 0}
                <!-- {@debug displayed_media} -->
                <span class="category-name">{getMediaNewCategory($current_category.content[$active_media_index].uuid, $active_media_change)}</span>
            {/if}
        {/if}
    </li>
    <li class:mv-np-disabled-field={$mv_tag_mode_enabled} class="mvnw-media-change">
        <span class="media-change" 
            class:media-change-moved={$active_media_change === media_change_types.MOVED}
            class:media-change-deleted={$active_media_change === media_change_types.DELETED}
            class:media-change-normal={$active_media_change === media_change_types.NORMAL} 
        >
            {$active_media_change}
        </span>
    </li>
    <li id="mvnw-random-navegation-state">
        <span class="mvnw-rns-label">random navigation: <span class="mvnw-rns-state">{$random_media_navigation ? 'on' : 'off'}</span></span>
    </li>
    <li class:mv-np-disabled-field={$mv_tag_mode_enabled} id="mvnw-skip-deleted-state">
        <span class="mvnw-sds-label">skip deleted: <span class="mvnw-sds-state">{$skip_deleted_medias ? 'on' : 'off'}</span></span>
    </li>
    <li class:mv-np-disabled-field={$mv_tag_mode_enabled} id="mvnw-auto-move-state">
        <span class="mvnw-ams-label">auto moving to: {$auto_move_on ? $auto_move_category?.name : "disabled"}</span>
    </li>
    <li id="mvnw-media-downloader">
        {#if $shared_active_media != null}
            <a href="{$shared_active_media.Url}" id="mvnw-media-downloader-anchor">
                <button class="dungeon-button thin">
                    <DownloadIcon 
                        icon_color="var(--main)" 
                        icon_size="calc(var(--spacing-3) * 0.8)"
                    />
                </button>
            </a>
        {/if}
    </li>
    <li id="mvnw-media-counter">
        {#if $current_category !== null}
            <h2><span>{$mv_tag_mode_enabled ? $active_tag_content_media_index+1 : $active_media_index+1}</span> ⋯ <span>{$mv_tag_mode_enabled ? $mv_tag_mode_total_content : $current_category.content.length}</span></h2>
        {/if}
    </li>
</ul>

<style>
    #media-viewer-navmenu-wrapper {
        gap: var(--vspacing-3);
        color: var(--grey-2);
        font-weight: 600;

        & > li button {
            font-size: var(--font-size-1);
            font-weight: normal;
        }
    }

    #mvnw-media-counter h2 {
        font-size: var(--font-size-1);
        font-family: var(--font-read);
        line-height: .6;
        color: var(--grey-3);
    }

    #mvnw-media-counter h2 span {
        color: var(--main-dark);
    }

    .media-change {
        color: var(--accent-2);
    }

    .media-change-moved {
        color: var(--success-4);
    }

    .media-change-deleted {
        color: var(--danger);
    }

    .mv-np-disabled-field {
        display: none;
    }
</style>