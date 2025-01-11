<script>
    import SettingEntry from "@components/Informative/DataEntries/SettingEntry.svelte";
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * the current list of billboard media uuids.
         * @type {string[]}
         */ 
        export let the_billboard_media_uuids = [];
    
    /*=====  End of Properties  ======*/
    
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the aggregation of the billboard medias.
         * @type {import('@components/Informative/DataEntries/data_entries').SettingChangeHandler}
         */
        const handleNewBillboardMedia = (setting_id, new_value) => {}
    
    /*=====  End of Methods  ======*/
    
    
</script>

<div id="cacow-billboard-medias-setting">
    <div id="cacow-billboard-medias-input-wrapper">
        <SettingEntry
            id_selector="billboard-medias-aggregator"
            font_size="var(--cacow-font-size)"
            information_entry_label="New media (uuid)"
            on_setting_change={handleNewBillboardMedia}
        />
    </div>
    <ol id="cacow-current-billboard-medias"
        class="dungeon-tag-group dungeon-tag-list" 
    >
        {#each the_billboard_media_uuids as media_uuid}
            <DeleteableItem
                item_id={media_uuid}    
                id_selector="cacow-billboard-media-{media_uuid}"
            >
                <p class="cacow-billboard-media">
                    {media_uuid}
                </p>
            </DeleteableItem>
        {/each}
    </ol>
</div>

<style>
    #cacow-billboard-medias-setting {
        display: contents;
    }

    ol#cacow-current-billboard-medias {
        & p {
            font-size: var(--cacow-font-size);
        }
    }
</style>