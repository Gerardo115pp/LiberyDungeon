<script>
    import { current_category } from "@stores/categories_tree";
    import { cluster_tags } from "@stores/dungeons_tags";
    import { getEntityTaggings } from "@models/DungeonTags";
    import { createEventDispatcher, onDestroy, onMount } from "svelte";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { toggleHotkeysSheet } from "@stores/layout";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import { linearCycleNavigationWrap } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";
    import { browser } from "$app/environment";
    
    /*=============================================
    =            Properties            =
    =============================================*/

        /*=============================================
        =            Hotkeys            =
        =============================================*/
    
            /**
             * @type {import('@libs/LiberyHotkeys/libery_hotkeys').HotkeyContextManager}
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
            
            /*=====  End of Hotkey movements  ======*/
            
            

        /*=====  End of Hotkeys  ======*/        
    
        /**
         * The current category's taggings.
         * @type {import("@models/DungeonTags").DungeonTagging[]}
         */ 
        export let current_category_taggings = [];
        $: globalThis.current_category_taggings = current_category_taggings;
        $: handleCategoryTaggingsChange(current_category_taggings);

        /**
         * A map of TagTaxonomy names -> DungeonTags.
         * @type {Map<string, import("@models/DungeonTags").DungeonTagging[]> | null}
         */
        let tag_taxonomy_map = null;
        $: globalThis.tag_taxonomy_map = tag_taxonomy_map;

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        component_mounted = true;
    });

    onDestroy(() => {
        if (!browser) return;

        dropHotkeyContext();
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
                if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    const hotkeys_context = new HotkeysContext();

                    hotkeys_context.register(["w", "s"], handleTagTaxonomyNavigation, {
                        description: "<navigation>move the focused attributed up and down.",
                    });

                    hotkeys_context.register(["a", "d"], handleTaxonomyTagNavigation, {
                        description: "<navigation>move the focused tag left and right.",
                    });

                    hotkeys_context.register(["\\d g"], handleTaxonomyTagGoto, {
                        description: "<navigation>go to the typed tag index.",
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
             * Drops the component hotkey context
             */
            const dropHotkeyContext = () => {
                if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) return;

                global_hotkeys_manager.dropContext(hotkeys_context_name);
            }

            /**
             * Emits an event to drop the hotkeys context
             */
            const emitDropHotkeyContext = () => {
                dispatch("drop-hotkeys-control");
            }

            /**
             * Handles the TagTaxonomy navigation hotkey.
             * @param {KeyboardEvent} event 
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey 
             */
            const handleTagTaxonomyNavigation = (event, hotkey) => {
                if (tag_taxonomy_map == null) return;

                const taxonomies = Array.from(tag_taxonomy_map.keys());

                if (taxonomies.length <= 0) return;
                

                let new_focused_taxonomy_index = focused_taxonomy_index;

                let navigation_step = event.key === "w" ? -1 : 1;

                new_focused_taxonomy_index = linearCycleNavigationWrap(new_focused_taxonomy_index, (taxonomies.length - 1), navigation_step).value;
                console.log(`new_focused_taxonomy_index: ${new_focused_taxonomy_index}`);

                let taxonomy_tag_count = tag_taxonomy_map.get(taxonomies[new_focused_taxonomy_index]).length;

                focused_tag_index = Math.max(0, Math.min((taxonomy_tag_count - 1), focused_tag_index));

                focused_taxonomy_index = new_focused_taxonomy_index;
            }

            /**
             * Handles the TaxonomyTag navigation hotkey.
             * @param {KeyboardEvent} event 
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey 
             */
            const handleTaxonomyTagNavigation = (event, hotkey) => {
                if (tag_taxonomy_map == null) return;

                const taxonomies = Array.from(tag_taxonomy_map.keys());

                if (taxonomies.length <= 0) return;

                let taxonomy_tag_count = tag_taxonomy_map.get(taxonomies[focused_taxonomy_index]).length;

                let new_focused_tag_index = focused_tag_index;

                let navigation_step = event.key === "a" ? -1 : 1;

                new_focused_tag_index = linearCycleNavigationWrap(new_focused_tag_index, (taxonomy_tag_count - 1), navigation_step).value;

                focused_tag_index = new_focused_tag_index;
            }

            /**
             * Handles the TaxonomyTag goto hotkey.
             * @param {KeyboardEvent} event 
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey 
             */
            const handleTaxonomyTagGoto = (event, hotkey) => {
                if (tag_taxonomy_map == null || !hotkey?.WithVimMotion) return;

                const taxonomies = Array.from(tag_taxonomy_map.keys());

                if (taxonomies.length <= 0 || !hotkey.HasMatch) return;

                let vim_motion_value = hotkey.MatchMetadata.MotionMatches[0];

                vim_motion_value--; // 1-based index to 0-based index

                vim_motion_value = Math.max(0, Math.min((tag_taxonomy_map.get(taxonomies[focused_taxonomy_index]).length - 1), vim_motion_value));

                focused_tag_index = vim_motion_value;
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
             * Drops the tools hotkey contexts and loads the previous context.
             */
            const resetHotkeyContext = () => {
                if (global_hotkeys_manager.ContextName !== hotkeys_context_name) return; 

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
         * Handles changes to the cluster tags.
         * @param {import("@models/DungeonTags").DungeonTagging[]} new_taggings
         */
        function handleCategoryTaggingsChange(new_taggings) {
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

    /*=====  End of Methods  ======*/
    
</script>

{#if tag_taxonomy_map != null && $cluster_tags.length > 0} 
    <div id="cpt-current-category-tags"
        class:hotkey-control={has_hotkey_control}
    >
        <header id="cpt-cct-header">
            <h4>
                Attributes for <span>{$current_category.name}</span>
            </h4>
        </header>
        {#each tag_taxonomy_map as [taxonomy_name, taxonomy_members], h}
            {@const is_taxonomy_keyboard_focused = focused_taxonomy_index === h && has_hotkey_control}
            <ol id="{$current_category.uuid}-attribute-{taxonomy_name}"
                class="current-category-attribute dungeon-tag-container"
                class:focused-attribute={is_taxonomy_keyboard_focused}
            >
                <p class="dungeons-field-label">{taxonomy_name}</p>
                {#each taxonomy_members as dungeon_tagging, k}
                    {@const is_tag_keyboard_focused = is_taxonomy_keyboard_focused && focused_tag_index === k}
                    <DeleteableItem
                        item_color={!is_tag_keyboard_focused ? "var(--grey)" : "var(--grey-8)"}
                        item_id={dungeon_tagging.Tag.Id}
                        on:item-deleted={handleCategoryTagDeleted}
                        squared_style
                    >
                        <p class="cca-attribute-name-wrapper">
                            <i>
                                {k + 1}
                            </i>
                            <span class="cca-attribute-name">
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
    #cpt-current-category-tags {
        display: flex;
        flex-direction: column;
        align-items: center;
        row-gap: var(--spacing-1);

        & > * {
            width: 50%;
        }
    }

    header#cpt-cct-header {
        text-align: center;
        margin-bottom: var(--spacing-2);

        & > h4 {
            color: var(--grey-6);
        }

        & > h4 > span {
            color: var(--main);
        }        
    }

    ol.current-category-attribute {
        color: var(--grey-1);
        line-height: 1;
        align-items: center;

        & :not(p.dungeons-field-label) {
            font-size: var(--font-size-1);
        }

        & .cca-attribute-name-wrapper {
            display: flex;
            column-gap: var(--spacing-1);
        }

        & .cca-attribute-name-wrapper > i {
            visibility: hidden;
            color: var(--grey-6);
        }

        & span.cca-attribute-name {
            display: flex;
            align-items: center;
            line-height: 1;
        }

        &.focused-attribute p.dungeons-field-label:first-child {
            color: var(--main);
        }
    }

    #cpt-current-category-tags.hotkey-control .cca-attribute-name-wrapper > i {
        visibility: visible;        
    }
</style>