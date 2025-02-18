<script>
    import { dungeon_tags_content_type } from "@app/common/content_types";
    import { ui_core_dungeon_references } from "@app/common/ui_references/core_ui_references";
    import { ui_pandasworld_tag_references } from "@app/common/ui_references/dungeon_tags_references";
    import SettingEntry from "@components/Informative/DataEntries/SettingEntry.svelte";
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { ClipboardContent } from "@models/Clipboard";
    import { DungeonTagList, getDungeonTagByID, getDungeonTags } from "@models/DungeonTags";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * A list of dungeon tag ids used to retrieve medias for the billboard.
         * @type {number[]}
         */ 
        export let the_billboard_tags_ids = [];

        /**
         * The loaded list of dungeon tags matching the provided tag ids.
         * @type {import('@models/DungeonTags').DungeonTag[]}
         */
        let the_billboard_tags = [];

        /**
         * The tags aggregator.
         * @type {SettingEntry}
        */
        let the_tags_aggregator;

        $: if (the_billboard_tags_ids.length > 0) {
            loadDungeonTags(the_billboard_tags_ids);
        } 
    
        
        /*----------  Event handlers  ----------*/
        
            /**
             * Triggered whenever a dungeon tag is added to the billboard tags.
             * @type {import('../category_configuration').CategoryConfig_BillboardDungeonTagsAdded}
             */ 
            export let onTagsAdded = (dungeon_tags) => {};

            /**
             * Triggered whenever a dungeon tag is removed from the billboard tags.
             * @type {import('../category_configuration').CategoryConfig_BillboardDungeonTagsRemoved}
             */
            export let onTagsRemoved = (dungeon_tag_id) => {};

    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/


        /**
         * Generates the dungeon_tags_map from a given list of dungeon tag ids.
         * @param {number[]} tag_ids
         * @returns {Promise<void>}
         */
        async function loadDungeonTags(tag_ids) {
            if (tag_ids.length === 0) {
                return resetComponentState();
            }

            const new_dungeon_tags = await getDungeonTags(tag_ids);

            if (new_dungeon_tags.length === 0) {
                return resetComponentState();
            }

            the_billboard_tags = new_dungeon_tags;
        }

        /**
         * Generates a human readable label for a dungeon tag.
         * @param {import('@models/DungeonTags').DungeonTag} dungeon_tag
         * @returns {string}
         */
        const generateHRDungeonTagLabel = dungeon_tag => {
            return dungeon_tag.Name;
        }
    
        /**
         * Handles the aggregation of the billboard medias.
         * @type {import('@components/Informative/DataEntries/data_entries').SettingChangeHandler}
         */
        const handleNewBillboardTag = async (setting_id, new_value) => {

            let tag_id = parseInt(new_value);

            if (isNaN(tag_id)) {
                console.error("In @components/Categories/CategoryConfiguration/sub-components/Setting__BillboardTags.svelte: the provided tag id is not valid");
                the_tags_aggregator.clearSettingValue();
                return;
            }

            const new_dungeon_tag = await getDungeonTagByID(tag_id);

            the_tags_aggregator.clearSettingValue();

            if (new_dungeon_tag === null) {
                const variable_environment = new VariableEnvironmentContextError("In @components/Categories/CategoryConfiguration/sub-components/Setting__BillboardTags.handleNewBillboardTag");

                variable_environment.addVariable("tag_id", tag_id);

                const labeled_err = new LabeledError(
                    variable_environment,
                    `Seems like there is no ${ui_pandasworld_tag_references.TAG.EntityName} with the id ${new_value}`,
                    lf_errors.ERR_HUMAN_ERROR
                )

                labeled_err.alert();
                return;
            }

            the_billboard_tags = [new_dungeon_tag, ...the_billboard_tags];

            onTagsAdded([new_dungeon_tag]);

            const hr_label = generateHRDungeonTagLabel(new_dungeon_tag);

            emitPlatformMessage(`Added ${hr_label}`);
        }

        /**
         * Handles the delete button of the DeletableItem component wrapping the dungeon tags
         * @param {CustomEvent<{item_id: number}>} event
         */
        const handleRemoveBillboardTag = event => {
            const tag_id = event.detail.item_id;

            /**
             * @type {import('@models/DungeonTags').DungeonTag}
             */
            let removed_dungeon_tag;

            /**
             * @type {typeof the_billboard_tags}
             */
            const new_billboard_dungeon_tags = [];

            for (let h = 0; h < the_billboard_tags.length; h++) {
                const on_index_tag = the_billboard_tags[h];

                if (on_index_tag.Id === tag_id) {
                    removed_dungeon_tag = on_index_tag;
                    continue
                }

                new_billboard_dungeon_tags.push(on_index_tag);
            }

            // @ts-ignore - an unassigned variable is equals to undefined... stupid ts.
            if (removed_dungeon_tag === undefined) {
                return;
            }

            the_billboard_tags = new_billboard_dungeon_tags;

            onTagsRemoved(removed_dungeon_tag.Id);

            const hr_label = generateHRDungeonTagLabel(removed_dungeon_tag);

            emitPlatformMessage(`Removed ${hr_label}`);
        }

        /**
         * Handles the clipboard paste event on the tag aggregator.
         * @param {ClipboardEvent} event
         * @returns {Promise<void>}
         */
        const handleClipboardPaste = async event => {
            console.log("event:", event);
            
            const clipboard_content = await readClipboard();

            if (clipboard_content === null) return;

            event.preventDefault();

            let new_dungeon_tags = clipboard_content.Content.DungeonTags;

            the_tags_aggregator.clearSettingValue();

            const current_tags = new Set(the_billboard_tags);
            let new_tags = new Set(new_dungeon_tags);

            new_tags = new_tags.difference(current_tags);

            the_billboard_tags = Array.from(current_tags.union(new_tags));
            new_dungeon_tags = Array.from(new_tags);

            onTagsAdded(new_dungeon_tags);
        }

        /**
         * Attempts to read the content from the clipboard. If it can parse a ClipboardContent<DungeonTagList>, it will return it. otherwise it will return null.
         * @returns {Promise<ClipboardContent<DungeonTagList> | null>}
         */
        const readClipboard = async () => {
            const clipboard_content_any = await ClipboardContent.fromClipboard();

            if (clipboard_content_any === null) return null;

            const content_any = clipboard_content_any.Content;


            const dungeon_tags_list = DungeonTagList.fromUnknown(content_any);

            if (dungeon_tags_list === null) return null;


            const clipboard_content = new ClipboardContent(dungeon_tags_content_type, dungeon_tags_list);

            return clipboard_content;
        }

        /**
         * Resets the state of the dungeon component.
         * @returns {void}
         */
        const resetComponentState = () => {
            the_billboard_tags = [];
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<div id="cacow-billboard-tags-setting">
    <div id="cacow-billboard-tags-input-wrapper">
        <SettingEntry
            bind:this={the_tags_aggregator}
            id_selector="billboard-tags-aggregator"
            font_size="var(--cacow-font-size)"
            information_entry_label="{ui_pandasworld_tag_references.TAG.EntityNamePlural} ID(or just paste them)"
            on_setting_change={handleNewBillboardTag}
            onClipboardPaste={handleClipboardPaste}
        />
    </div>
    <ol id="cacow-current-billboard-tags"
        class="dungeon-tag-group dungeon-tag-list" 
    >
        {#each the_billboard_tags as dungeon_tag}
            <DeleteableItem
                item_id={dungeon_tag.Id}    
                id_selector="cacow-billboard-tag-{dungeon_tag.Id}"
                on:item-deleted={handleRemoveBillboardTag}
            >
                <p class="cacow-billboard-tag">
                    {dungeon_tag.Name}
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