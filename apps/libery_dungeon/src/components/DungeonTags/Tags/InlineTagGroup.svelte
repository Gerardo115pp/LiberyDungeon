<script>
    import DeleteableItem from '@components/ListItems/DeleteableItem.svelte';
    import { cleanIdSelector } from '@libs/utils';
    import { GRID_MOVEMENT_ITEM_CLASS } from '@app/common/keybinds/CursorMovement';
    
    /*=============================================
    =            Properties            =
    =============================================*/


    
        /**
         * The grouped dungeon tags.
         * @type {import('../dungeon_tags').DungeonTags_GroupedTags<any>}
         */
        export let the_grouped_tags;

        /**
         * The index of the focused tag.
         * @type {number}
         */
        let focused_tag_index = 0;

        /**
         * Whether the inline group is active.
         * @type {boolean}
         */
        export let active_inline_group = false;

    /*=====  End of Properties  ======*/
   
    
    /*=============================================
    =            Methods            =
    =============================================*/

    /*=====  End of Methods  ======*/
    
</script>

<ol class="tag-inline-group dungeon-tag-container">
    <p class="dungeons-field-label">
        {the_grouped_tags.taxonomy_name}
    </p>
    {#each the_grouped_tags.tags as inline_grouped_tag, h}
        {@const is_tag_keyboard_focused = active_inline_group && focused_tag_index === h}
        <DeleteableItem
            class_selector="itg-{cleanIdSelector(the_grouped_tags.taxonomy_name)}-{GRID_MOVEMENT_ITEM_CLASS}"
            item_color={!is_tag_keyboard_focused ? "var(--grey)" : "var(--grey-8)"}
            item_id={inline_grouped_tag.Id}
            squared_style
        >
            <p class="cma-attribute-name-wrapper"
                class:focused-attribute={is_tag_keyboard_focused}
            >
                <span class="cma-attribute-name">
                    {inline_grouped_tag.Name}
                </span>
            </p>
        </DeleteableItem>
    {/each}
</ol>

<style>
    ol.tag-inline-group {
        --itc-font-size: var(--font-size-1);
        line-height: 1;
        padding: 0 var(--spacing-2);
        align-items: center;

        & p.dungeons-field-label {
            font-size: var(--itc-font-size);
            text-transform: lowercase;
        }

        & p.dungeons-field-label::first-letter {
            text-transform: uppercase;
        }
    }

    p.cma-attribute-name-wrapper {
        font-size: var(--itc-font-size);
    }
</style>