<script>
    import { cluster_tags } from "@stores/dungeons_tags";
    import { createEventDispatcher, onDestroy, onMount } from "svelte";
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { toggleHotkeysSheet } from "@stores/layout";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import { linearCycleNavigationWrap } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";
    import { CursorMovementWASD, GRID_MOVEMENT_ITEM_CLASS } from "@app/common/keybinds/CursorMovement";
    import { browser } from "$app/environment";
    
    /*=============================================
    =            Properties            =
    =============================================*/

        /*=============================================
        =            Hotkeys            =
        =============================================*/
    
            /**
             * @type {import('@libs/LiberyHotkeys/libery_hotkeys').HotkeyContextManager | null}
             */
            const global_hotkeys_manager = getHotkeysManager();

            const hotkeys_context_name = "category_attributes";

            /*=============================================
            =            Hotkeys state            =
            =============================================*/

                /**
                 * Whether the component has mounted or not.
                 * @type {boolean}
                 */
                let component_mounted = false;
            
                /**
                 * Whether it has hotkey control.
                 * @type {boolean}
                 */ 
                export let has_hotkey_control = false;
                $: if (has_hotkey_control && component_mounted) {
                    defineDesktopKeybinds();
                }
        
            /*=====  End of Hotkeys state  ======*/
            
            /*=============================================
            =            Hotkey movement                 = 
            =============================================*/
            
                /**
                 * The focused TagTaxonomy list index.
                 * @type {number}
                 */
                let focused_taxonomy_index = 0;

                /**
                 * The focused taxonomy's focused tag index.
                 * @type {number}
                 */
                let focused_tag_index = 0;

                /**
                 * The grid navigation WASD keybind setter.
                 * @type {CursorMovementWASD | null}
                 */
                let the_wasd_keybind_wrapper = null;

                /**
                 * A map to cache last cursor position when changing taxonomies tagging lists.
                 * @type {Map<string, number>}
                 */
                const last_cursor_positions = new Map();
                
            
            /*=====  End of Hotkey movements  ======*/

        /*=====  End of Hotkeys  ======*/        
    
        /**
         * The current category's taggings.
         * @type {import("@models/DungeonTags").DungeonTagging[]}
         */ 
        export let current_media_taggings = [];
        $: handleMediaTaggingsChange(current_media_taggings);

        /**
         * The currently active media.
         * @type {import('@models/Medias').Media}
         */
        export let the_active_media;
        $: if (the_active_media != null && the_wasd_keybind_wrapper != null) {
            updateWasdGridSelectors();
        }

        /**
         * A map of TagTaxonomy names -> DungeonTags.
         * @type {Map<string, import("@models/DungeonTags").DungeonTagging[]> | null}
         */
        let tag_taxonomy_map = null;

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        component_mounted = true;
    });

    onDestroy(() => {
        if (!browser) return;

        dropHotkeyContext();
        dropGridNavigationWrapper();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /*=============================================
        =            Keybinds            =
        =============================================*/
        
            /**
             * Defines the tools hotkeys.
             */ 
            const defineDesktopKeybinds = () => {
                if (global_hotkeys_manager == null) return;

                if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    const hotkeys_context = new HotkeysContext();

                    setGridNavigationWrapper(hotkeys_context);

                    hotkeys_context.register(["alt+w", "alt+s"], handleTagTaxonomyNavigation, {
                        description: "<navigation>move the focused attributed up and down.",
                    });

                    hotkeys_context.register(["x"], handleCategoryUnassignTag, {
                        description: "<navigation>Removes the focused attribute value from the current category.",
                    });
                    
                    hotkeys_context.register(["q", "t"], handleCloseCategoryTaggerTool, {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Closes the category tagger tool.`,
                        await_execution: false
                    });

                    hotkeys_context.register(["?"], toggleHotkeysSheet, {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Opens the hotkeys cheat sheet.`
                    });

                    global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
                }
                
                global_hotkeys_manager.loadContext(hotkeys_context_name);
            }

            /**
             * Changes the taxonomy grid used by the WASD keybind wrapper. called by handleTagTaxonomyNavigation.
             */
            const changeTaxonomyGrid = () => {
                if (the_wasd_keybind_wrapper == null) return;


                const old_grid_parent_selector = the_wasd_keybind_wrapper.MovementController.DomParentSelector();

                last_cursor_positions.set(old_grid_parent_selector, focused_tag_index);


                const grid_selectors = getFocusedTaggingsGridSelectors();

                if (grid_selectors == null) {
                    console.error("No grid selectors found for the focused taggings.");
                    return;
                };

                the_wasd_keybind_wrapper.changeGridContainer(grid_selectors.grid_parent_selector, grid_selectors.grid_member_selector);


                const last_cursor_position = last_cursor_positions.get(grid_selectors.grid_parent_selector);

                if (last_cursor_position != null) {
                    let settled = the_wasd_keybind_wrapper.MovementController.Grid.setCursor(last_cursor_position);

                    if (settled) {
                        focused_tag_index = last_cursor_position;
                    }
                } else {
                    focused_tag_index = the_wasd_keybind_wrapper.MovementController.Grid.clampSequenceIndex(focused_tag_index);

                    the_wasd_keybind_wrapper.MovementController.Grid.setCursor(focused_tag_index);
                }
            }

            /**
             * Drops the component hotkey context
             */
            const dropHotkeyContext = () => {
                if (global_hotkeys_manager == null || !global_hotkeys_manager.hasContext(hotkeys_context_name)) return;

                global_hotkeys_manager.dropContext(hotkeys_context_name);
            }

            const dropGridNavigationWrapper = () => {
                if (the_wasd_keybind_wrapper != null) {
                    the_wasd_keybind_wrapper.destroy();
                }
            }

            /**
             * Emits an event to drop the hotkeys context
             */
            const emitDropHotkeyContext = () => {
                dispatch("drop-hotkeys-control");
            }

            /**
             * Returns the focused tagging.
             * @returns {import("@models/DungeonTags").DungeonTagging | null}
             */
            const getFocusedTagging = () => {
                if (tag_taxonomy_map == null) return null;

                const taxonomies = Array.from(tag_taxonomy_map.keys());

                if (taxonomies.length <= 0) return null;

                let focused_taxonomy = taxonomies[focused_taxonomy_index];

                if (focused_taxonomy == null) return null;

                // @ts-ignore
                let focused_tagging = tag_taxonomy_map.get(focused_taxonomy)[focused_tag_index];

                return focused_tagging;
            }

            /**
             * Handles the TagTaxonomy navigation hotkey.
             * @param {KeyboardEvent} event 
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey 
             */
            const handleTagTaxonomyNavigation = (event, hotkey) => {
                if (tag_taxonomy_map == null || tag_taxonomy_map.size === 0) return;
                

                let new_focused_taxonomy_index = focused_taxonomy_index;

                let navigation_step = event.key.toLowerCase() === "w" ? -1 : 1;


                new_focused_taxonomy_index = linearCycleNavigationWrap(new_focused_taxonomy_index, (tag_taxonomy_map.size - 1), navigation_step).value;
                console.log(`new_focused_taxonomy_index: ${new_focused_taxonomy_index}`);

                focused_taxonomy_index = new_focused_taxonomy_index;

                changeTaxonomyGrid();
            }

            /**
             * Removes the focused tag from the current category.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleCategoryUnassignTag = async (event, hotkey) => {
                let focused_tagging = getFocusedTagging();

                console.log("focused_tagging:", focused_tagging);

                if (focused_tagging == null) return;

                emitRemoveCategoryTagEvent(focused_tagging.Tag.Id);

                focused_tag_index = Math.max(0, focused_tag_index - 1);
            }

            /**
             * Emits an event to close the category tagger tool and drops the hotkeys context.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleCloseCategoryTaggerTool = (event, hotkey) => {
                resetHotkeyContext();
                emitDropHotkeyContext();
            }

            /**
             * Handles the cursor update event from the grid navigation wrapper.
             * @param {import("@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils").GridWrappedValue} cursor
             */
            const handleCursorUpdate = cursor => {
                focused_tag_index = cursor.value;
            }

            /**
             * Drops the tools hotkey contexts and loads the previous context.
             */
            const resetHotkeyContext = () => {
                if (global_hotkeys_manager == null || global_hotkeys_manager.ContextName !== hotkeys_context_name) return; 

                global_hotkeys_manager.loadPreviousContext();
            }

        /*=====  End of Keybinds  ======*/

        /**
         * Returns the tag taxonomy map from the taxonomies contained on the current category taggings. 
         * @param {import("@models/DungeonTags").DungeonTagging[]} taggings
         * @returns {Map<string, import("@models/DungeonTags").DungeonTagging[]>}
         */ 
        const getTagTaxonomyMap = taggings => {
            let new_tag_taxonomy_map = new Map();

            const taxonomy_uuid_to_name_lookup = new Map();

            for (let tagging of taggings) {
                let tag_taxonomy_name = taxonomy_uuid_to_name_lookup.get(tagging.Tag.TaxonomyUUID);

                if (tag_taxonomy_name == null) {
                    tag_taxonomy_name = getTagTaxonomyNameByUUID(tagging.Tag.TaxonomyUUID);
                    taxonomy_uuid_to_name_lookup.set(tagging.Tag.TaxonomyUUID, tag_taxonomy_name);
                }

                if (tag_taxonomy_name == null) {
                    console.warn(`Tag taxonomy with UUID '${tagging.Tag.TaxonomyUUID}' not found.`);
                    continue;
                }

                let tag_taxonomy_members = new_tag_taxonomy_map.get(tag_taxonomy_name) ?? [];

                tag_taxonomy_members.push(tagging);

                new_tag_taxonomy_map.set(tag_taxonomy_name, tag_taxonomy_members);
            }

            return new_tag_taxonomy_map;
        }

        /**
         * Returns a the name of a given tag taxonomy or null if it's not among the taxonomies in cluster_tags.
         * @param {string} tag_taxonomy_uuid
         * @returns {string | null}
         */
        const getTagTaxonomyNameByUUID = tag_taxonomy_uuid => {
            let tag_taxonomy_name = null;

            for (let tag_taxonomy_tags of $cluster_tags) {
                if (tag_taxonomy_tags.Taxonomy.UUID === tag_taxonomy_uuid) {
                    tag_taxonomy_name = tag_taxonomy_tags.Taxonomy.Name;
                    break;
                }
            }

            return tag_taxonomy_name;
        }

        /**
         * Returns the focused Taggings list pair of selectors on a 2 sized array with index 0 being the grid-parent selector and the index 1 being the grid-member selector.
         * This is used for the parameters of the CursorMovementWASD either on creation or afterwards for it's changeGridContainer method.
         * @returns {import('@common/interfaces/common_actions').GridSelectors | null}
         */
        const getFocusedTaggingsGridSelectors = () => {
            
            if (tag_taxonomy_map == null) return null;

            const grid_selectors = {
                grid_parent_selector: "",
                grid_member_selector: "",
            }


            const taxonomy_name = getFocusedTagTaxonomyName();

            grid_selectors.grid_parent_selector = `#media-${the_active_media.uuid}-attribute-${taxonomy_name}`;
            grid_selectors.grid_member_selector = `.${taxonomy_name}-${GRID_MOVEMENT_ITEM_CLASS}`;

            return grid_selectors;
        }

        /**
         * Returns the taxonomy of the focused array of taggings.
         * @returns {import("@models/DungeonTags").TagTaxonomy | null}
         */
        const getFocusedTagTaxonomy = () => {
            if (tag_taxonomy_map == null) return null;

            const [taxonomy_name, current_category_taggings] = [...tag_taxonomy_map][focused_taxonomy_index];

            if (current_category_taggings.length === 0) {
                console.warn(`This should never happen, a taxonomy in the tag_taxonomy_map['${taxonomy_name}'] has not taggings for the current category. but members of the_taxonomy_map are generated from the taggins of the current category. Reeks of a programming error.`);
                return null;
            }

            const taxonomy_uuid = current_category_taggings[0].Tag.TaxonomyUUID;

            const focused_taxonomy = $cluster_tags.find(tag_taxonomy => tag_taxonomy.Taxonomy.UUID === taxonomy_uuid);

            return focused_taxonomy?.Taxonomy ?? null;
        }

        /**
         * Returns the focused tag taxonomy name.
         * @returns {string | null}
         */
        const getFocusedTagTaxonomyName = () => {
            if (tag_taxonomy_map == null) return null;

            const [taxonomy_name, current_category_taggings] = [...tag_taxonomy_map][focused_taxonomy_index];

            return taxonomy_name;
        }

        /**
         * Handles changes to the cluster tags.
         * @param {import("@models/DungeonTags").DungeonTagging[]} new_taggings
         */
        function handleMediaTaggingsChange(new_taggings) {
            if (new_taggings.length <= 0) return;

            const updated_tag_taxonomy_map = getTagTaxonomyMap(new_taggings);

            if (updated_tag_taxonomy_map.size <= 0) {
                console.warn("No tag taxonomies found for the current category.");
                return;
            }

            console.log("Updating taxonomy map:", updated_tag_taxonomy_map);

            tag_taxonomy_map = updated_tag_taxonomy_map;
        }

        /**
         * Handles the item-deleted event on one of the current category tags.
         * @param {CustomEvent<{item_id: number}>} event 
         */
        const handleCategoryTagDeleted = event => {
            let tag_id = event?.detail?.item_id

            emitRemoveCategoryTagEvent(tag_id);
        }

        /**
         * Emits the remove-category-tag event with the given tag id.
         * @param {number} tag_id
         */
        const emitRemoveCategoryTagEvent = tag_id => {
            dispatch("remove-category-tag", {removed_tag: tag_id});
        }

        /**
         * Sets the grid navigation wrapper required data.
         * @param {import("@libs/LiberyHotkeys/hotkeys_context").default} hotkeys_context
         */
        const setGridNavigationWrapper = (hotkeys_context) => {
            if (!browser) return;

            if (the_wasd_keybind_wrapper != null) {
                the_wasd_keybind_wrapper.destroy();
            }

            const grid_selectors = getFocusedTaggingsGridSelectors();

            if (grid_selectors == null) {
                console.error("No grid selectors found for the focused taggings.");
                return;
            };

            const matching_elements_count = document.querySelectorAll(grid_selectors.grid_parent_selector).length;
            if (matching_elements_count !== 1) {
                throw new Error(`tag parent selector '${grid_selectors.grid_parent_selector}' returned ${matching_elements_count}, expected exactly 1`);
            }


            the_wasd_keybind_wrapper = new CursorMovementWASD(grid_selectors.grid_parent_selector, handleCursorUpdate, {
                initial_cursor_position: focused_tag_index,
                sequence_item_name: "value",
                sequence_item_name_plural: "values",
                grid_member_selector: grid_selectors.grid_member_selector,
            });
            the_wasd_keybind_wrapper.setup(hotkeys_context);
            

            // @ts-ignore
            globalThis.the_grid_navigation_wrapper = the_wasd_keybind_wrapper; 
        }

        /**
         * Updates the selectors for the wasd cursor movement wrapper.
         */
        const updateWasdGridSelectors = () => {
            if (the_wasd_keybind_wrapper == null) return;

            const grid_selectors = getFocusedTaggingsGridSelectors();

            if (grid_selectors == null) {
                console.error("No grid selectors found for the focused taggings.");
                return;
            };

            the_wasd_keybind_wrapper.changeGridContainer(grid_selectors.grid_parent_selector, grid_selectors.grid_member_selector);
        }

    

    /*=====  End of Methods  ======*/
    
</script>

{#if the_active_media != null && tag_taxonomy_map != null && $cluster_tags.length > 0} 
    <div id="mpt-current-media-tags"
        class:hotkey-control={has_hotkey_control}
    >
        <header id="mpt-cmt-header">
            <h4>
                Attributes for <span>{the_active_media.name}</span>
            </h4>
        </header>
        {#each tag_taxonomy_map as [taxonomy_name, taxonomy_members], h}
            {@const is_taxonomy_keyboard_focused = focused_taxonomy_index === h && has_hotkey_control}
            {@const taxonomy_members_list_id_selector = `media-${the_active_media.uuid}-attribute-${taxonomy_name}`}
            <ol id={taxonomy_members_list_id_selector}
                class="current-media-attribute dungeon-tag-container"
                class:focused-attribute={is_taxonomy_keyboard_focused}
            >
                <p class="dungeons-field-label">{taxonomy_name}</p>
                {#each taxonomy_members as dungeon_tagging, k}
                    {@const is_tag_keyboard_focused = is_taxonomy_keyboard_focused && focused_tag_index === k}
                    <DeleteableItem
                        class_selector="{taxonomy_name}-{GRID_MOVEMENT_ITEM_CLASS}"
                        item_color={!is_tag_keyboard_focused ? "var(--grey)" : "var(--grey-8)"}
                        item_id={dungeon_tagging.Tag.Id}
                        on:item-deleted={handleCategoryTagDeleted}
                        squared_style
                    >
                        <p class="cma-attribute-name-wrapper">
                            <i>
                                {k + 1}
                            </i>
                            <span class="cma-attribute-name">
                                {dungeon_tagging.Tag.Name}
                            </span>
                        </p>
                    </DeleteableItem>
                {/each}
            </ol>
        {/each}
    </div>
{/if}

<style>
    #mpt-current-media-tags {
        display: flex;
        flex-direction: column;
        align-items: center;
        row-gap: var(--spacing-1);

        & > * {
            width: 50%;
        }
    }

    header#mpt-cmt-header {
        text-align: center;
        margin-bottom: var(--spacing-2);

        & > h4 {
            color: var(--grey-6);
        }

        & > h4 > span {
            color: var(--main);
        }        
    }

    ol.current-media-attribute {
        color: var(--grey-1);
        line-height: 1;
        align-items: center;

        & :not(p.dungeons-field-label) {
            font-size: var(--font-size-1);
        }

        & .cma-attribute-name-wrapper {
            display: flex;
            column-gap: var(--spacing-1);
        }

        & .cma-attribute-name-wrapper > i {
            visibility: hidden;
            color: var(--grey-6);
        }

        & span.cma-attribute-name {
            display: flex;
            align-items: center;
            line-height: 1;
        }

        & p.dungeons-field-label {
            text-transform: lowercase;
        }

        & p.dungeons-field-label::first-letter {
            text-transform: uppercase;
        }

        &.focused-attribute p.dungeons-field-label:first-child {
            color: var(--main);
        }
    }

    #mpt-current-media-tags.hotkey-control .cma-attribute-name-wrapper > i {
        visibility: visible;        
    }
</style>