<script>
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { createEventDispatcher } from "svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * A list of dungeon tags to display.
         * @type {import("@models/DungeonTags").DungeonTag[]}
         */ 
        export let dungeon_tags = [];

        /**
         * Whether to enable the tag creator element.
         * @type {boolean}
        */
        export let enable_tag_creator = false;

        
        /*----------  Style  ----------*/
            /**
             * The tag group's color.
             * @type {string}
             * @default "var(--grey-9)"
             */
            export let tag_group_color = "var(--grey-9)";

            /**
             * The color to use for the tag creator element.
             * @type {string}
             * @default "var(--grey-8)"
             */
            export let tag_creator_color = "var(--grey-5)";
        
            /**
             * Whether to display the tags indexes. These been determined by the order in which the tags were passed.
             * @type {boolean}
             * @default true
             */
            export let expose_indexes = true;
        
        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    
    /*=============================================
    =            Method            =
    =============================================*/
    
        /**
         * Handles the tag selection event.
         * @param {CustomEvent<{item_id: string}>} event
         */    
        const handleTagSelection = (event) => {
            emitTagSelected(event.detail.item_id);
        }

        /**
         * Emits the tag-selected event.
         * @param {number} tag_id
        */
        const emitTagSelected = (tag_id) => {
            dispatch("tag-selected", {tag_id});
        }
    
    /*=====  End of Method  ======*/
    
    
</script>

<ol class="dungeon-tag-group dungeon-tag-list">
    {#each dungeon_tags as tag, index (tag.Id)}
        <DeleteableItem 
            item_color={tag_group_color}
            item_id={tag.Id}
            on:item-selected={handleTagSelection}
        >
            <p class="dtg-tag-name taxonomy-member">
                {#if expose_indexes}
                    <i class="dtg-tag-index">
                        {index + 1}
                    </i>
                {/if}
                <span>
                    {tag.Name}
                </span>
            </p>
        </DeleteableItem>
    {/each}
    {#if enable_tag_creator}
        <DeleteableItem
            item_color={tag_creator_color}
            is_protected
        >
            <input class="tag-creator taxonomy-member" 
                type="text"
                placeholder="New tag"
                minlength="3"
                maxlength="64"
                pattern="{'[a-z_]{4,64}'}"
            />
        </DeleteableItem>
    {/if}
</ol>

<style>
    .taxonomy-member {
        display: flex;
        font-size: var(--font-size-1);
        line-height: 1;
        column-gap: var(--spacing-2);
    }

    p.dtg-tag-name {

        padding: 0 var(--spacing-1);

        & > i.dtg-tag-index {
            color: var(--grey-3);
            line-height: 1;
        }

        & > span {
            line-height: 1;
        }
    }
</style>