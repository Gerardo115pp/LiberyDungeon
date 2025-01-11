<script>
    import DungeonDataEntry from '@components/Informative/DataEntries/DungeonDataEntry.svelte';
    import SettingEntry from '@components/Informative/DataEntries/SettingEntry.svelte';
    import { slide } from 'svelte/transition';
    import { quartOut } from 'svelte/easing';
    import { getMediaIdentityByUUID } from '@models/Medias';
    import { LabeledError, VariableEnvironmentContextError } from '@libs/LiberyFeedback/lf_models';
    import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
    import { changeCategoryThumbnail } from '@models/Categories';
    import SettingBillboardMedias from './sub-components/Setting__BillboardMedias.svelte';
    import SettingBillboardTags from './sub-components/Setting__BillboardTags.svelte';
    import { ui_core_dungeon_references } from '@app/common/ui_references/core_ui_references';
    import { ui_pandasworld_tag_references } from '@app/common/ui_references/dungeon_tags_references';

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

        /**
         * The current category's configuration.
         * @type {import('@models/Categories').CategoryConfig | null}
         */
        let category_config = null;

        $: handleCurrentCategoryChange(the_inner_category);
       
        
        /*----------  Sections state  ----------*/
        
            /**
             * Whether the billboard section should be enabled.
             * @type {boolean}
             * @default false
             */ 
            export let enable_billboard_config = false;
        /*----------  Event handlers  ----------*/
        
            /**
             * Handler for the thumbnail-change event.
             * @type {import("./category_configuration").CategoryConfig_ThumbnailChanged}
             */ 
            export let on_thumbnail_change = updated_category => {};

            /**
             * Called when ever billboard configuration change.
             * @type {import('./category_configuration').CategoryConfig_CategoryConfigChanged}
             */
            export let onBillboardConfigChanged = new_config => {};
    
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

        /**
         * Handles the billboard media aggregation.
         * @type {import('./category_configuration').CategoryConfig_BillboardMediaAdded}
         */
        const handleBillboardMediaAdded = async media_identity => {
            if (category_config == null) {
                console.error("In CategoryConfiguration.handleBillboardMediaAdded: category_config is null");
                return;
            }

            const new_billboard_medias_uuids = [media_identity.Media.uuid, ...category_config.BillboardMediaUUIDs];

            await category_config.updateBillboardMediaUUIDs(new_billboard_medias_uuids);

            onBillboardConfigChanged(category_config)
        }

        /**
         * Handles the billboard media removal.
         * @type {import('./category_configuration').CategoryConfig_BillboardMediaRemoved}
         */
        const handleBillboardMediaRemoved = async media_uuid => {
            if (category_config == null) {
                console.error("In CategoryConfiguration.handleBillboardMediaRemoved: category_config is null");
                return;
            }

            const new_billboard_medias_uuids = category_config.BillboardMediaUUIDs.filter(media_uuid_in_list => media_uuid_in_list !== media_uuid);

            await category_config.updateBillboardMediaUUIDs(new_billboard_medias_uuids);

            onBillboardConfigChanged(category_config)
        }

        /**
         * Handles the current category change.
         * @param {import('@models/Categories').InnerCategory} new_category
         */
        const handleCurrentCategoryChange = async (new_category) => {
            if (new_category == null) {
                return;
            }

            category_config = await loadCategoryConfig();
        }

        /**
         * Loads the category's configuration.
         * @returns {Promise<import('@models/Categories').CategoryConfig | null>}
         */
        const loadCategoryConfig = async () => {
            if (the_inner_category.hasConfig()) {
                return the_inner_category.Config;
            }

            let loaded_config = null;

            try {
                loaded_config = await the_inner_category.loadCategoryConfig();
            } catch (err) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryConfiguration.loadCategoryConfig")

                variable_environment.addVariable("category_uuid", the_inner_category.uuid);

                const labeled_err = new LabeledError(variable_environment, `Sorry, there was an error loading the configuration for category '${the_inner_category.uuid}'.`, lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();
            }

            return loaded_config;
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<dialog open 
    class="category-config-wrapper libery-dungeon-window dungeon-scroll"
    transition:slide={{axis: 'y', easing: quartOut}}
>
    <section id="cacow-general-settings" 
        class="cacow-settings-section"
    >
        <header class="cacow-configuration-header">
            <h4 class="cacow-category-identity">
                General settings
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
    </section>
    {#if enable_billboard_config}
        <section id="cacow-billboard-section"
            class="cacow-settings-section"          
        >
            <header class="cacow-configuration-header">
                <h4 class="cacow-category-identity">
                    Billboard settings
                </h4>
                <div id="cacow-billboard-medias-setting-description" 
                    class="dungeon-description cacow-header-instructtions" 
                >
                    <p 
                        class="cacow-setting-instructions "
                    >
                        Change the medias that appear in the billboard for this {ui_core_dungeon_references.CATEGORY.EntityName}. You have three options:
                    </p>
                    <p class="cacow-setting-instructions">
                        <b>First</b>: Add a set of {ui_pandasworld_tag_references.TAG_TAXONOMY.EntityName} {ui_pandasworld_tag_references.TAG.EntityNamePlural} in the appropriate section(you can paste a set of copied {ui_pandasworld_tag_references.TAG.EntityNamePlural} copied from the Viewer) and then all {ui_core_dungeon_references.MEDIA.EntityNamePlural} that match, will be displayed in the billboard.
                    </p>
                    <p class="cacow-setting-instructions">
                        <b>Second</b>: Add a list of {ui_core_dungeon_references.MEDIA.EntityName} uuids and only those will show in the billboard.
                    </p>
                    <p class="cacow-setting-instructions">
                        <b>Third</b>: If the {ui_core_dungeon_references.CATEGORY.EntityName} has any {ui_core_dungeon_references.MEDIA.EntityNamePlural} then one of those will be picked at random.
                    </p>
                </div>
            </header>
            <div class="cacow-group-wrapper">
                {#if category_config != null}
                    <SettingBillboardMedias 
                        the_billboard_media_uuids={category_config.BillboardMediaUUIDs}
                        onMediaAdded={handleBillboardMediaAdded}
                        onMediaRemoved={handleBillboardMediaRemoved}
                    />
                {/if}
            </div>
            <div class="cacow-group-wrapper">
                {#if category_config != null}
                    <SettingBillboardTags 
                        the_billboard_tags={category_config.BillboardDungeonTags}
                    />
                {/if}
            </div>
        </section>
    {/if}
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
        row-gap: var(--spacing-4);
        padding-inline: 0.8em;
        padding-block: calc(var(--spacing-3) + var(--spacing-1));
        border-color: hsl(from var(--grey-8) h s l / 0.7);
        backdrop-filter: var(--backdrop-filter-blur);
        overscroll-behavior: contain;
        z-index: var(--z-index-t-5);
    }

    section.cacow-settings-section {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);
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

    .cacow-setting-instructions {
        color: var(--grey-3);
    } 

    

    /*=============================================
    =            Billboard settings            =
    =============================================*/
        header.cacow-configuration-header:has(.cacow-header-instructtions) {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-2);
        }

        #cacow-billboard-medias-setting-description {
            display: flex;
            flex-direction: column;
            row-gap: calc(var(--spacing-1) * 0.7);

            & > .cacow-setting-instructions:first-of-type {
                font-weight: 600;
            }
        }
    
        .cacow-group-wrapper {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-1);
            padding-block: var(--spacing-1);
        } 


        .cacow-group-wrapper:not(:last-of-type) {
            border-bottom: 1px solid var(--grey-9);
        }
    
    /*=====  End of Billboard settings  ======*/

    @container (width < 260px) {
        header.cacow-configuration-header {
            --cacow-font-size: calc(var(--font-size-1) * 0.8)
        }         

        menu.cacow-settings {
            --cacow-font-size: calc(var(--font-size-1) * 0.8)
        } 
    }
</style>