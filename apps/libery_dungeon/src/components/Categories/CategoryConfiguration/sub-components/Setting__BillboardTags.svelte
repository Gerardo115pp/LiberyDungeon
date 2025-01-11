<script>
    import { ui_core_dungeon_references } from "@app/common/ui_references/core_ui_references";
    import { ui_pandasworld_tag_references } from "@app/common/ui_references/dungeon_tags_references";
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
        const handleNewBillboardTag = (setting_id, new_value) => {
            console.log("setting_id:", setting_id);
            console.log("new_value:", new_value);
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<div id="cacow-billboard-tags-setting">
    <div id="cacow-billboard-tags-input-wrapper">
        <SettingEntry
            id_selector="billboard-tags-aggregator"
            font_size="var(--cacow-font-size)"
            information_entry_label="{ui_pandasworld_tag_references.TAG.EntityNamePlural} ID(or just paste them)"
            on_setting_change={handleNewBillboardTag}
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
        {:else}
            <p class="dungeon-generic-text">
                No filtering {ui_pandasworld_tag_references.TAG_TAXONOMY.EntityName} {ui_pandasworld_tag_references.TAG.EntityNamePlural} added for this {ui_core_dungeon_references.CATEGORY.EntityName}'s billboard
            </p>
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