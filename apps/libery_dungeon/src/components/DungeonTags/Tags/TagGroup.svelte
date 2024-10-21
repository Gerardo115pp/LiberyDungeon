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
            export let tag_creator_color = "var(--grey-8)";
        
            /**
             * Whether to display the tags indexes. These been determined by the order in which the tags were passed.
             * @type {boolean}
             * @default true
             */
            export let expose_indexes = true;
    
    /*=====  End of Properties  ======*/
    
</script>

<ol class="dungeon-tag-group">
    {#each dungeon_tags as tag, index (tag.Id)}
        <DeleteableItem 
            item_color={tag_group_color}
        >
            <p class="dtg-tag-name">
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
        <input class="dtg-tag-creator" 
            type="text"
            placeholder="New tag"
            minlength="3"
            maxlength="64"
            pattern="{'[a-z_]{4,64}'}"
        />
    </DeleteableItem>
</ol>