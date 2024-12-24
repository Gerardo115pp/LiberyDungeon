<script>
    import { hasClipboardWritePermission } from "@libs/utils";
    import { onMount } from "svelte";
    import InformationEntry from "./sub-components/InformationEntry.svelte";
    import { active_media_index } from "@stores/media_viewer";
    import { current_category } from "@stores/categories_tree";
    import CategoryThumbnailSetBtn from "./sub-components/CategoryThumbnailSetBtn.svelte";
    import { getCategory } from "@models/Categories";

    
    /*=============================================
    =            Properties            =
    =============================================*/


        /**
         * external categories loaded from the server. This is a map of category uuids -> category objects.
         * @type {Map<string, import("@models/Categories").Category>}
         */
        let external_categories_cache = new Map();
    
        /**
         * The media that is being viewed
         * @type {import("@models/Medias").Media}
         */
        export let current_media_information;


        /**
         * The category been viewed
         * @type {import("@models/Categories").Category | null}
         */
        let current_category_information = null;
        $: if (current_category_information?.UUID !== current_media_information.main_category) {
            current_category_information = null;
            updateCategoryInformation(current_media_information.main_category);
        }

        /**
         * The cluster been viewed
         * @type {import("@models/CategoriesClusters").CategoriesCluster}
         */
        export let current_cluster_information;

        
        
        /*----------  Behavior  ----------*/

            /**
             * Whether the provided media exists inside the current category.
             * @type {boolean}
             * @default true
             */
            export let media_in_current_category = true;
        
            /**
             * Whether or not the app/site has paste permission.
             * @type {boolean}
             */
            let can_clipboard_paste = false;
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        can_clipboard_paste = await hasClipboardWritePermission();
    });

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Loads the appropriate category information. If media_in_current_category is true and current_category is no null, then it uses that. Otherwise checks the
         * external_categories_cache for the category and if it's not there then it fetches it from the server.
         * @param {string} new_category_uuid
         */ 
        async function updateCategoryInformation(new_category_uuid) {
            if ($current_category != null && media_in_current_category) {
                current_category_information = $current_category.asCategory();
                return;
            }

            const cached_category = external_categories_cache.get(new_category_uuid);

            if (cached_category !== undefined) {
                current_category_information = cached_category;
                return;
            }

            const fetched_category = await getCategory(new_category_uuid);

            if (fetched_category !== null) {
                external_categories_cache.set(new_category_uuid, fetched_category);
            }

            current_category_information = fetched_category;
        }
    
    /*=====  End of Methods  ======*/
    
    

</script>

<!-- 
    Remember end user terminology:
    Cluster -> Dungeon
    Category -> Cell
    Media -> Media
 -->

<article id="mv-medias-information-panel"
    class="libery-dungeon-window"
>
    {#if current_category_information !== null}
        <div id="mv-mip-content-wrapper">
            <hgroup id="mv-mip-media-information" class="mv-mip-media-information-section">
                <h3 class="mv-mip-section-name">
                    The media
                </h3>
                <ul class="mv-mip-media-information-list">
                    <InformationEntry 
                        information_entry_label="Name" 
                        information_entry_value="{current_media_information.name}" 
                        paste_on_click={can_clipboard_paste}
                    />
                    <InformationEntry 
                        information_entry_label="UUID" 
                        information_entry_value="{current_media_information.uuid}" 
                        paste_on_click={can_clipboard_paste}
                    />
                    <InformationEntry 
                        information_entry_label="Media Path"
                        information_entry_value="{current_cluster_information.FSPath}/{current_category_information.FullPath}{current_media_information.name}"
                        paste_on_click={can_clipboard_paste}
                    />
                </ul>
            </hgroup>
            <hgroup id="mv-mip-category-information" class="mv-mip-media-information-section">
                <h3 class="mv-mip-section-name">
                    The cell
                </h3>
                <ul class="mv-mip-media-information-list">
                    <InformationEntry 
                        information_entry_label="Name" 
                        information_entry_value="{current_category_information.Name}" 
                        paste_on_click={can_clipboard_paste}
                    />
                    <InformationEntry 
                        information_entry_label="UUID" 
                        information_entry_value="{current_category_information.UUID}" 
                        paste_on_click={can_clipboard_paste}
                    />
                    <InformationEntry 
                        information_entry_label="Parent cell UUID" 
                        information_entry_value="{current_category_information.Parent}" 
                        paste_on_click={can_clipboard_paste}
                    />
                </ul>
                <menu id="mv-mip-category-actions">
                    <CategoryThumbnailSetBtn 
                        the_current_media={current_media_information}
                    />
                </menu>
            </hgroup>
            <hgroup id="mv-mip-cluster-information" class="mv-mip-media-information-section">
                <h3 class="mv-mip-section-name">
                    The dungeon
                </h3>
                <ul class="mv-mip-media-information-list">
                    <InformationEntry 
                        information_entry_label="Name" 
                        information_entry_value="{current_cluster_information.Name}" 
                        paste_on_click={can_clipboard_paste}
                    />
                    <InformationEntry 
                        information_entry_label="UUID" 
                        information_entry_value="{current_cluster_information.UUID}" 
                        paste_on_click={can_clipboard_paste}
                    />
                    <InformationEntry 
                        information_entry_label="Root cell UUID" 
                        information_entry_value="{current_cluster_information.RootCategoryID}" 
                        paste_on_click={can_clipboard_paste}
                    />
                    <InformationEntry 
                        information_entry_label="Dungeon download cell" 
                        information_entry_value="{current_cluster_information.DownloadCategoryID}" 
                        paste_on_click={can_clipboard_paste}
                    />
                </ul>
            </hgroup>
        </div>
    {/if}
</article>

<style>
    article#mv-medias-information-panel {
        --mv-mip-background-color: var(--body-bg-color);
        --mv-mip-titles-color: var(--main-dark);
        --mv-mip-content-color: var(--grey-1);

        width: 100cqw;
        height: 100cqh;
        background: var(--mv-mip-background-color);
        border-radius: var(--border-radius);
        padding: var(--vspacing-2);

        & ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }
    }

    #mv-mip-content-wrapper {
        display: flex;
        flex-direction: column;
        gap: var(--vspacing-3);
    }

    hgroup.mv-mip-media-information-section {
        display: flex;
        flex-direction: column;
        gap: var(--vspacing-2);
    }

    h3.mv-mip-section-name {
        font-family: var(--font-read);
        font-size: var(--font-size-1);
        color: var(--mv-mip-titles-color);
    }

    ul.mv-mip-media-information-list {
        display: flex;
        flex-direction: column;
        gap: var(--vspacing-1);
    }

    @supports (color: rgb( from white r g b / 1)) {
        article#mv-medias-information-panel {
            --mv-mip-background-color: hsl(from var(--body-bg-color) h s calc(l * 1.5) / 0.9);
            --mv-mip-titles-color: var(--main-dark);
            --mv-mip-content-color: var(--grey-1);
        }
    }

    
    /*=============================================
    =            Media section actions            =
    =============================================*/
    
        menu#mv-mip-category-actions {
            width: 100%;
            display: flex;
            justify-content: center;
        } 
    
    /*=====  End of Media section actions  ======*/
    
    
</style>