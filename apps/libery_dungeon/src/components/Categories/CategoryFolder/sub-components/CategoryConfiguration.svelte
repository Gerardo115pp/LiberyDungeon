<script>
    import DungeonDataEntry from '@components/Informative/DataEntries/DungeonDataEntry.svelte';
    import SettingEntry from '@components/Informative/DataEntries/SettingEntry.svelte';
    import { slide } from 'svelte/transition';
    import { quartOut } from 'svelte/easing';

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The category to show the configuration for.
         * @type {import('@models/Categories').InnerCategory}
         */ 
        export let the_inner_category;

        const category_setting_ids = {
            CATEGORY_THUMBNAIL: 'category_thumbnail_uuid',
        }
    
    /*=====  End of Properties  ======*/
    
</script>

<dialog open 
    class="category-config-wrapper libery-dungeon-window"
    transition:slide={{axis: 'y', easing: quartOut}}
>
    <header class="cacow-configuration-header">
        <h4 class="cacow-category-identity">
            Settings for {the_inner_category.name}
        </h4>
    </header>
    <menu class="cacow-settings">
        <div class="cacow-data-entry">
            <DungeonDataEntry
                font_size="var(--cacow-font-size)"
                information_entry_label="uuid"
                information_entry_value={the_inner_category.uuid}
                paste_on_click
            /> 
        </div>
        <div class="cacow-setting-entry">
            <SettingEntry
                id_selector={category_setting_ids.CATEGORY_THUMBNAIL}
                font_size="var(--cacow-font-size)"
                information_entry_label="thumbnail"
                information_entry_value={the_inner_category.ThumbnailUUID}
            /> 
        </div>
    </menu>
</dialog>

<style>
    dialog.category-config-wrapper {
        --cacow-font-size: var(--font-size-1);

        display: flex;
        container-type: size;
        width: 100%;
        height: 100%;
        background: var(--grey-t);
        flex-direction: column;
        row-gap: var(--spacing-3);
        padding-inline: 0.8em;
        padding-block: calc(var(--spacing-3) + var(--spacing-1));
        border-color: hsl(from var(--grey-8) h s l / 0.7);
        backdrop-filter: var(--backdrop-filter-blur);
        z-index: var(--z-index-t-5);
    }

    header.cacow-configuration-header {
        font-size: calc(var(--cacow-font-size) * 1.2);
    }

    menu.cacow-settings {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-2);
        color: var(--grey-1);
    }

    @container (width < 260px) {
        header.cacow-configuration-header {
            --cacow-font-size: calc(var(--font-size-1) * 0.8)
        }         

        menu.cacow-settings {
            --cacow-font-size: calc(var(--font-size-1) * 0.8)
        } 
    }
</style>