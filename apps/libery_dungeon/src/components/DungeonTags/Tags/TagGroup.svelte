<script>
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * A list of dungeon tags to display.
         * @type {import("@models/DungeonTags").DungeonTag[]}
         */ 
        export let dungeon_tags = [];

        
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
    
    /*=====  End of Properties  ======*/
    
</script>

<ol class="dungeon-tag-group dungeon-tag-list">
    {#each dungeon_tags as tag, index (tag.Id)}
        <DeleteableItem 
            item_color={tag_group_color}
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