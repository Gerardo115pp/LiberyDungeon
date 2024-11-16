<script>
    import DungeonDataEntry from '@components/Informative/DataEntries/DungeonDataEntry.svelte';
    import SettingEntry from '@components/Informative/DataEntries/SettingEntry.svelte';
    import { slide } from 'svelte/transition';
    import { quartOut } from 'svelte/easing';
    import { getMediaIdentityByUUID } from '@models/Medias';
    import { LabeledError, VariableEnvironmentContextError } from '@libs/LiberyFeedback/lf_models';
    import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
    import { changeCategoryThumbnail } from '@models/Categories';

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

        
        /*----------  Event handlers  ----------*/
        
            /**
             * Handler for the thumbnail-change event.
             * @type {import('./category_folder_subs').CategoryConfig_ThumbnailChanged}
             */ 
            export let on_thumbnail_change = updated_category => {};
    
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles new setting changes.
         * @type {import('@components/Informative/DataEntries/data_entries').SettingChangeHandler}
         */
        const handleSettingChanges = (setting_id, new_value) => {
            switch (setting_id) {
            case category_setting_ids.CATEGORY_THUMBNAIL: 
                handleThumbnailChanges(new_value);
                return;
            default:
                console.error(`In CategoryConfiguration.handleSettingChanges: no setting called ${setting_id}`);
                return;
            }
        }

        /**
         * Handles thumbnail changes.
         * @param {string} new_thumbnail_uuid
         * @returns {Promise<void>}
         */
        const handleThumbnailChanges = async (new_thumbnail_uuid) => {

            const thumbnail_media = await getMediaIdentityByUUID(new_thumbnail_uuid);

            if (thumbnail_media == null) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryConfiguration.handleThumbnailChanges")

                variable_environment.addVariable("category_uuid", the_inner_category.uuid);
                variable_environment.addVariable("current_thumbnail_uuid", the_inner_category.ThumbnailUUID);
                variable_environment.addVariable("new_thumbnail_uuid", new_thumbnail_uuid);

                const labeled_err = new LabeledError(variable_environment, `Seems like media '${new_thumbnail_uuid}' doesn't actually exist.`, lf_errors.ERR_HUMAN_ERROR);

                labeled_err.alert();

                return;
            }

            const thumbnail_set = await changeCategoryThumbnail(the_inner_category.uuid, new_thumbnail_uuid);

            if (!thumbnail_set) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryConfiguration.handleThumbnailChanges")

                variable_environment.addVariable("category_uuid", the_inner_category.uuid);
                variable_environment.addVariable("current_thumbnail_uuid", the_inner_category.ThumbnailUUID);
                variable_environment.addVariable("new_thumbnail_uuid", new_thumbnail_uuid);

                const labeled_err = new LabeledError(variable_environment, `Sorry, there was an error setting the thumbnail to '${new_thumbnail_uuid}'.`, lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();

                return;

            }

            the_inner_category.setThumbnail(thumbnail_media);

            on_thumbnail_change(the_inner_category);
        }
    
    /*=====  End of Methods  ======*/
    
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
                on_setting_change={handleSettingChanges}
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