<script>
    import events from "@components/Popups/events";
    import { hasClipboardWritePermission } from "@libs/utils";
    import { onMount } from "svelte";
    import InformationEntry from "./sub-components/InformationEntry.svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The media that is being viewed
         * @type {import("@models/Medias").Media}
        */
        export let current_media_information;


        /**
         * The category been viewed
         * @type {import("@models/Categories").CategoryLeaf}
         */
        export let current_category_information;


        /**
         * The cluster been viewed
         * @type {import("@models/CategoriesClusters").CategoriesCluster}
         */
        export let current_cluster_information;


        
        /*----------  Behavior  ----------*/
        
            /**
             * Whether or not the app/site has paste permission.
             * @type {boolean}
             */
            let can_clipboard_paste = false;
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        can_clipboard_paste = hasClipboardWritePermission();
    });

</script>

<!-- 
    Remember end user terminology:
    Cluster -> Dungeon
    Category -> Cell
    Media -> Media
 -->

<article id="mv-medias-information-panel">
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
                    information_entry_value="{current_category_information.FullPath}{current_media_information.name}"
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
                    information_entry_value="{current_category_information.name}" 
                    paste_on_click={can_clipboard_paste}
                />
                <InformationEntry 
                    information_entry_label="UUID" 
                    information_entry_value="{current_category_information.uuid}" 
                    paste_on_click={can_clipboard_paste}
                />
                <InformationEntry 
                    information_entry_label="Parent cell UUID" 
                    information_entry_value="{current_category_information.parent}" 
                    paste_on_click={can_clipboard_paste}
                />
            </ul>
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
</article>

<style>
    article#mv-medias-information-panel {
        --mv-mip-background-color: var(--grey-8);
        --mv-mip-border-color: var(--main-8);
        --mv-mip-titles-color: var(--main-dark);
        --mv-mip-content-color: var(--grey-1);

        width: 100cqw;
        height: 100cqh;
        background: var(--mv-mip-background-color);
        border-radius: var(--border-radius);
        border: 0.2px solid var(--mv-mip-border-color);
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
            --mv-mip-background-color: hsl(from var(--grey) h s calc(l * 1.5) / 0.9);
            --mv-mip-border-color: hsl(from var(--main-8) h calc(s * 0.4) calc(l * 0.8) / 1);
            --mv-mip-titles-color: var(--main-dark);
            --mv-mip-content-color: var(--grey-1);
        }
    }
</style>