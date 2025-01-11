<script>
    import SettingEntry from "@components/Informative/DataEntries/SettingEntry.svelte";
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * A list of dungeon tag ids used to retrieve medias for the billboard.
         * @type {number[]}
         */ 
        export let the_billboard_tags = [];
    
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

<div id="cacow-billboard-tags-setting">
    <div id="cacow-billboard-tags-input-wrapper">
        <SettingEntry
            id_selector="billboard-tags-aggregator"
            font_size="var(--cacow-font-size)"
            information_entry_label="New media (uuid)"
            on_setting_change={handleNewBillboardMedia}
        />
    </div>
    <ol id="cacow-current-billboard-tags"
        class="dungeon-tag-group dungeon-tag-list" 
    >
        {#each the_billboard_tags as tag_id}
            <DeleteableItem
                item_id={tag_id}    
                id_selector="cacow-billboard-tag-{tag_id}"
            >
                <p class="cacow-billboard-tag">
                    {tag_id}
                </p>
            </DeleteableItem>
        {/each}
    </ol>
</div>

<style>
    #cacow-billboard-tags-setting {
        display: contents;
    }

    ol#cacow-current-billboard-tags {
        & p {
            font-size: var(--cacow-font-size);
        }
    }
</style>