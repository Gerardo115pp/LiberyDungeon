<script>
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import TaxonomyTags from "../../TagTaxonomyComponents/TaxonomyTags.svelte";
    import { cluster_tags } from "@stores/dungeons_tags";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import { toggleHotkeysSheet } from "@stores/layout";
    import { createEventDispatcher, onDestroy, onMount, tick } from "svelte";
    import { browser } from "$app/environment";
    import { linearCycleNavigationWrap } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";
    
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

                const hotkeys_context_name = "cluster_public_tags";

                /*=============================================
                =            Hotkeys state            =
                =============================================*/

                    /**
                     * Whether the component has mounted or not.
                     * @type {boolean}
                     */
                    let has_mounted = false;
                
                    /**
                     * Whether it has hotkey control.
                     * @type {boolean}
                     */ 
                    export let has_hotkey_control = false;
                    $: if (has_hotkey_control && has_mounted) {
                        defineDesktopKeybinds();
                    }
            
                /*=====  End of Hotkeys state  ======*/
                
                
                /*=============================================
                =            Hotkeys movement            =
                =============================================*/
                
                    /**
                     * The index of the focused tag_taxonomy.
                     * @type {number}
                     */
                    let cpt_focused_tag_taxonomy_index = 0;
                    $: if (cpt_focused_tag_taxonomy_index >= 0 && has_hotkey_control) {
                        ensureFocusedTagTaxonomyVisible();
                    }

                    /**
                     * Whether the focused tag taxonomy is active.
                     * @type {boolean}
                     */
                    let cpt_focused_tag_taxonomy_active = false;
                
                /*=====  End of Hotkeys movement  ======*/

            /*=====  End of Hotkeys  ======*/

            const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        has_mounted = true;
    });

    onDestroy(() => {
        if (!browser) return;

        dropHotkeyContext();
    })

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

                    hotkeys_context.register(["q", "t"], handleDropHotkeyControl, {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Deselects the Category tagger section.`,
                        await_execution: false
                    });

                    hotkeys_context.register(["w", "s"], handleTagTaxonomyNavigation, {
                        description: "<navigation>Moves the changes the focused attribute up/down.",
                    });

                    hotkeys_context.register(["e"], handleTagTaxonomySelection, {
                        description: "<navigation>Selects the focused attribute.",
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
             * Handles the selection of tag taxonomies sections
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleTagTaxonomySelection = (event, hotkey) => {
                event.preventDefault();

                cpt_focused_tag_taxonomy_active = true; 
            }

            /**
             * Handles the navigation of the tag taxonomies.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleTagTaxonomyNavigation = (event, hotkey) => {
                if ($cluster_tags?.length === 0 || (event.key !== "w" && event.key !== "s")) return;

                let new_cpt_focused_tag_taxonomy_index = cpt_focused_tag_taxonomy_index;

                let navigation_step = event.key === "w" ? -1 : 1;

                new_cpt_focused_tag_taxonomy_index = linearCycleNavigationWrap(new_cpt_focused_tag_taxonomy_index, $cluster_tags.length - 1, navigation_step).value;
                console.log("new_cpt_focused_tag_taxonomy_index:", new_cpt_focused_tag_taxonomy_index);
                console.log("$cluster_tags.length:", $cluster_tags.length);

                cpt_focused_tag_taxonomy_index = new_cpt_focused_tag_taxonomy_index;
            }

            /**
             * Recovers the hotkeys control and deactivates the active section.
             */
            const handleRecoverHotkeysControl = () => {
                cpt_focused_tag_taxonomy_active = false;
            }

            /**
             * Emits an event to close the category tagger tool and drops the hotkeys context.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleDropHotkeyControl = (event, hotkey) => {
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
         * Ensures that focused tag taxonomy item on the minimap are visible.
         */
        const ensureFocusedTagTaxonomyVisible = async () => {
            await tick();

            const focused_tag_taxonomy = document.querySelector(".cpt-ttm-taxonomy.keyboard-focused");

            if (!focused_tag_taxonomy) return;

            focused_tag_taxonomy.scrollIntoView({
                block: "center",
                inline: "center"
            });
        }

    /*=====  End of Methods  ======*/

</script>

<div id="cluster-public-tags">
    <aside class="cpt-taxonomy-tags-minimap">
        <header class="cpt-ttm-header">
            <h4>
                Attributes
            </h4>
        </header>
        <ol class="cpt-ttm-tag-taxonomies">
            {#each $cluster_tags as taxonomy_tags, h (taxonomy_tags.Taxonomy.UUID)}
                {@const is_keyboard_focused = cpt_focused_tag_taxonomy_index === h && has_hotkey_control}
                <li  class="cpt-ttm-taxonomy"
                    class:keyboard-focused={is_keyboard_focused}
                    class:has-hotkey-control={is_keyboard_focused && cpt_focused_tag_taxonomy_active}
                >
                    <p class="cpt-ttm-taxonomy-name">
                        {taxonomy_tags.Taxonomy.Name}
                    </p>
                </li>
            {/each}
        </ol>
    </aside>
    <div id="cpt-sections-wrapper"
        class="dungeon-scroll" 
    >
        {#each $cluster_tags as taxonomy_tags, h (taxonomy_tags.Taxonomy.UUID)}
            {@const is_keyboard_focused = cpt_focused_tag_taxonomy_index === h && has_hotkey_control}
            <TaxonomyTags 
                taxonomy_tags={taxonomy_tags}
                has_hotkey_control={is_keyboard_focused && cpt_focused_tag_taxonomy_active}
                is_keyboard_focused={is_keyboard_focused}
                enable_tag_creation
                on:tag-selected
                on:taxonomy-content-change
                on:drop-hotkeys-control={handleRecoverHotkeysControl}
            />
        {/each}
    </div>
</div>

<style>
    #cluster-public-tags {
        display: grid;
        grid-template-columns: max-content 1fr;
        box-sizing: border-box;
        width: 100%;
        column-gap: var(--spacing-2);
    }
    
    /*=============================================
    =            Taxonomy tags minimap            =
    =============================================*/
    
        aside.cpt-taxonomy-tags-minimap {
            display: flex;
            height: 100cqh;
            flex-direction: column;
            overflow-y: auto;
            border-bottom: 0.2px solid var(--main-9);
            row-gap: var(--spacing-1);
            scrollbar-width: none;
        } 

        aside.cpt-taxonomy-tags-minimap header.cpt-ttm-header {
            text-align: center;

            & h4 {
                font-family: var(--font-read);
                color: var(--grey-6);
                font-weight: 600;
            }
        }

        aside.cpt-taxonomy-tags-minimap ol.cpt-ttm-tag-taxonomies {
            display: flex;
            flex-direction: column;
            row-gap: calc(0.25 * var(--spacing-1));

            & li.cpt-ttm-taxonomy {
                text-align: center;
                color: var(--grey-2);
                background: var(--grey-8);
                padding: var(--spacing-1) var(--spacing-2);
                border-radius: var(--border-radius);
            }

            & li.cpt-ttm-taxonomy p {
                font-size: var(--font-size-1);
                line-height: 1;
                text-transform: lowercase;
            }

            & li.cpt-ttm-taxonomy.keyboard-focused {
                background: var(--main-7);
                transition: background .2s ease-out;
            }

            & li.cpt-ttm-taxonomy.has-hotkey-control {
                background: var(--main-dark);
            }

            & li.cpt-ttm-taxonomy.has-hotkey-control p.cpt-ttm-taxonomy-name {
                color: var(--grey-1);
            }
        }
    
    /*=====  End of Taxonomy tags minimap   ======*/

    #cpt-sections-wrapper {
        display: flex;
        flex-direction: column;
        height: 100cqh;
        row-gap: var(--spacing-2);
        overflow-y: auto;
    }
</style>