<script>
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import TaxonomyTags from "../../TagTaxonomyComponents/TaxonomyTags.svelte";
    import { cluster_tags } from "@stores/dungeons_tags";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import { toggleHotkeysSheet } from "@stores/layout";
    import { createEventDispatcher, onDestroy, onMount } from "svelte";
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
                        description: "<navigation>Moves the focused tag taxonomy up or down.",
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
            const handleSectionSelection = (event, hotkey) => {
                event.preventDefault();

                cpt_focused_tag_taxonomy_active = true; 
            }

            /**
             * Handles the navigation of the tag taxonomies.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleTagTaxonomyNavigation = (event, hotkey) => {

            }

            /**
             * Recovers the hotkeys control and deactivates the active section.
             */
            const handleRecoverHotkeysControl = () => {
                ct_section_active = false;
            }

            /**
             * Emits an event to close the category tagger tool and drops the hotkeys context.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleDropHotkeyControl = (event, hotkey) => {
                resetHotkeyContext();
                emitDropHotkeyContext();
                console.log("Add you'r close function here");
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
         * Returns all the non-internal TaxonomyTag count.
         * @returns {number}
         */
        const getNonInternalTagTaxonomyCount = () => {
            let count = 0;

            $cluster_tags.forEach(taxonomy_tags => {
                if (!taxonomy_tags.Taxonomy.IsInternal) {
                    count++;
                }
            });

            return count;
        }

    /*=====  End of Methods  ======*/

</script>

<div id="cpt-wrapper">
    {#each $cluster_tags as taxonomy_tags, h (taxonomy_tags.Taxonomy.UUID)}
        {#if !taxonomy_tags.Taxonomy.IsInternal}
            {@const is_keyboard_focused = cpt_focused_tag_taxonomy_index === h && has_hotkey_control}
            <TaxonomyTags 
                taxonomy_tags={taxonomy_tags}
                has_hotkey_control={is_keyboard_focused && cpt_focused_tag_taxonomy_active}
                is_keyboard_focused={is_keyboard_focused}
                enable_tag_creation
                on:tag-selected
                on:taxonomy-content-change
                on:drop-hotkeys-control={handleDropHotkeyControl}
            />
        {/if} 
    {/each}
</div>

<style>
    #cpt-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-2);
    }
</style>