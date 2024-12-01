<script>
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import TaxonomyTags from "./TaxonomyTags/TaxonomyTags.svelte";
    import { cluster_tags } from "@stores/dungeons_tags";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import { toggleHotkeysSheet } from "@stores/layout";
    import { createEventDispatcher, onDestroy, onMount, tick } from "svelte";
    import { browser } from "$app/environment";
    import { linearCycleNavigationWrap } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";
    import { confirmPlatformMessage, emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { renameTagTaxonomy } from "@models/DungeonTags";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { LabeledError } from "@libs/LiberyFeedback/lf_models";
    import generateTaxonomyTagsHotkeysContext from "@components/DungeonTags/TagTaxonomyComponents/TaxonomyTags/taxonomy_tags_hotkeys";
    import generateClusterPublicTagsHotkeyContext, { cluster_public_tags_actions } from "./cluster_public_tags_hotkeys";
    import { common_action_groups } from "@app/common/keybinds/CommonActionsName";
    import { wrapShowHotkeysTable } from "@app/common/keybinds/CommonActionWrappers";
    
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

                /**
                 * The component hotkey context for the cluster public tags component.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
                 */
                export let component_hotkey_context = generateClusterPublicTagsHotkeyContext();

                /**
                 * Component hotkey context to pass down to the Tag taxonomy components.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
                 */
                export let taxonomy_tags_hotkeys_context = generateTaxonomyTagsHotkeysContext();
                
                /*=============================================
                =            Hotkeys state            =
                =============================================*/
                
                    /**
                     * Whether it has hotkey control.
                     * @type {boolean}
                     */ 
                    export let has_hotkey_control = false;
                    $: component_hotkey_context.Active = has_hotkey_control;
            
                /*=====  End of Hotkeys state  ======*/
                
                
                /*=============================================
                =            Hotkeys movement            =
                =============================================*/
                
                    /**
                     * The index of the focused tag_taxonomy.
                     * @type {number}
                     */
                    let cpt_focused_tag_taxonomy_index = 0;
                    $: if (cpt_focused_tag_taxonomy_index >= 0 && component_hotkey_context.Active) {
                        ensureFocusedTagTaxonomyVisible();
                    }
                
                /*=====  End of Hotkeys movement  ======*/

            /*=====  End of Hotkeys  ======*/

            /**
             * Whether the focused tag taxonomy is on rename mode.
             * @type {boolean}
             */
            let cpt_focused_tag_taxonomy_rename_mode = false;

            /**
             * The tag taxonomy renamer input element.
             * @type {HTMLInputElement}
             */
            let the_taxonomy_renamer_input;

            /**
             * Whether the tag renamer name is ready.
             * @type {boolean}
             */
            let taxonomy_renamer_is_ready = false;


        /*----------  User feedback  ----------*/
        
            /**
             * A UIReference object, to create ui messages about the taggable entity.
             * @type {import('@libs/LiberyFeedback/lf_models').UIReference}
             */
            export let ui_entity_reference;

            /**
             * A UIReference object, to create ui messages about the tag taxonomy.
             * @type {import('@libs/LiberyFeedback/lf_models').UIReference}
             */
            export let ui_taxonomy_reference;

            /**
             * A UIReference object, to create ui messages about the dungeons tags. used to pass down to general purpose components.
             * @type {import('@libs/LiberyFeedback/lf_models').UIReference}
             */
            export let ui_tag_reference;
        
        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        component_hotkey_context.onActiveChange(handleComponentActiveState);

        defineSubComponentsHotkeysContext();
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
                if (global_hotkeys_manager == null) {
                    console.error("The hotkeys manager is not available");
                    return;
                }

                if (!global_hotkeys_manager.hasContext(component_hotkey_context.HotkeysContextName)) {
                    const hotkeys_context = preparePublicHotkeyActions(component_hotkey_context);

                    hotkeys_context.register(["q"], handleDropHotkeyControl, {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Deselects the Category tagger section.`,
                        await_execution: false
                    });

                    hotkeys_context.register(["e"], handleTagTaxonomySelection, {
                        description: "<navigation>Selects the focused attribute.",
                    });

                    hotkeys_context.register(["x"], handleFocusedTagTaxonomyDeletion, {
                        description: "<content>Deletes the focused attribute.",
                        await_execution: false,
                    });

                    hotkeys_context.register(["c c"], handleRenameFocusedTagTaxonomyHotkey, {
                        description: "<content>Renames the focused attribute.",
                        await_execution: false,
                        mode: "keyup"
                    });

                    wrapShowHotkeysTable(hotkeys_context);

                    component_hotkey_context.applyExtraHotkeys();

                    global_hotkeys_manager.declareContext(component_hotkey_context.HotkeysContextName, hotkeys_context);
                }
                
                global_hotkeys_manager.loadContext(component_hotkey_context.HotkeysContextName);
            }

            /**
             * Prepares the public hotkey actions to generate the hotkeys context.
             * @param {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext} new_component_hotkey_context
             * @returns {import('@libs/LiberyHotkeys/hotkeys_context').default}
             */
            const preparePublicHotkeyActions = (new_component_hotkey_context) => {
                if (component_hotkey_context.HasGeneratedHotkeysContext()) {
                    component_hotkey_context.dropHotkeysContext();
                }

                /* --------------------------- up/down navigation --------------------------- */

                    const up_down_navigation = new_component_hotkey_context.getHotkeyActionOrPanic(cluster_public_tags_actions.WS_NAVIGATION);

                    up_down_navigation.Options = {
                        description: `${common_action_groups.NAVIGATION}Moves the changes the focused attribute up/down.`
                    }

                    up_down_navigation.Callback = handleTagTaxonomyNavigation;

                /* -------------------------------------------------------------------------- */

                const hotkeys_context = new_component_hotkey_context.generateHotkeysContext();

                return hotkeys_context
            }

            /**
             * Drops the component hotkey context
             */
            const dropHotkeyContext = () => {
                if (global_hotkeys_manager == null || !global_hotkeys_manager.hasContext(component_hotkey_context.HotkeysContextName)) return;

                global_hotkeys_manager.dropContext(component_hotkey_context.HotkeysContextName);
            }

            /**
             * Emits an event to drop the hotkeys context
             */
            const emitDropHotkeyContext = () => {
                dispatch("drop-hotkeys-control");
            }

            /**
             * Handles the change in active state from the component hotkey context.
             * @param {boolean} new_state
             */
            const handleComponentActiveState = new_state => {
                has_hotkey_control = new_state;

                if (new_state) {
                    defineDesktopKeybinds();
                }
            }

            /**
             * Handles the selection of tag taxonomies sections
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleTagTaxonomySelection = (event, hotkey) => {
                event.preventDefault();

                taxonomy_tags_hotkeys_context.Active = true;
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
             * Deletes the focused tag taxonomy.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleFocusedTagTaxonomyDeletion = async (event, hotkey) => {
                if (taxonomy_tags_hotkeys_context.Active) return;

                let focused_taxonomy = $cluster_tags[cpt_focused_tag_taxonomy_index];

                if (focused_taxonomy == null) return;

                let user_choice = await confirmPlatformMessage({
                    message_title: `Delete '${focused_taxonomy.Taxonomy.Name}'`,
                    question_message: `Are you sure you wish to delete '${focused_taxonomy.Taxonomy.Name}' this will make it unavailable in the entire dungeon, remove any associations it was with medias or categories, and remove any values( it has ${focused_taxonomy.Tags.length}) associated with the attribute.`,
                    auto_focus_cancel: true,
                    danger_level: 2,
                    confirm_label: "Deleted it",
                    cancel_label: "cancel",
                });

                if (user_choice !== 1) return;

                emitDeleteTagTaxonomy(focused_taxonomy.Taxonomy.UUID);

                cpt_focused_tag_taxonomy_index = Math.max(0, cpt_focused_tag_taxonomy_index - 1);
            }

            /**
             * Recovers the hotkeys control and deactivates the active section.
             */
            const handleRecoverHotkeysControl = () => {
                taxonomy_tags_hotkeys_context.Active = false;
            }

            /**
             * Handles the renaming of the focused tag taxonomy.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleRenameFocusedTagTaxonomyHotkey = (event, hotkey) => {
                if (taxonomy_tags_hotkeys_context.Active) return;

                cpt_focused_tag_taxonomy_rename_mode = true;
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
                if (global_hotkeys_manager == null || global_hotkeys_manager.ContextName !== component_hotkey_context.HotkeysContextName) return; 

                global_hotkeys_manager.loadPreviousContext();
            }

            /* --------------------- sub-components hotkey contexts --------------------- */

                /**
                 * Defines the hotkeys context for sub-components.
                 */
                const defineSubComponentsHotkeysContext = () => {   
                    component_hotkey_context.addChildContext(taxonomy_tags_hotkeys_context);

                    component_hotkey_context.inheritExtraHotkeys();
                }

        /*=====  End of Keybinds  ======*/

        /**
         * Checks if a Tag taxonomy exists by its name.
         * @param {string} tag_name
         * @returns {boolean}
         */
        const checkTaxonomyNameExists = (tag_name) => {
            return $cluster_tags.some(taxonomy_tags => taxonomy_tags.Taxonomy.Name === tag_name);
        }
        
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

        /**
         * Emits an event to delete a given tag taxonomy by its UUID.
         * @param {string} taxonomy_uuid
         */
        const emitDeleteTagTaxonomy = (taxonomy_uuid) => {
            dispatch("delete-taxonomy", {
                taxonomy_uuid
            });
        }

        /**
         * Emits the taxonomy-content-change event with the given taxonomy UUID.
         * @param {string} taxonomy_uuid
         */
        const emitTaxonomyContentChange = taxonomy_uuid => {
            dispatch("taxonomy-content-change", {taxonomy: taxonomy_uuid});
        }

        /**
         * Returns the focused TaxonomyTags instances.
         * @returns {import("@models/DungeonTags").TaxonomyTags | null}
         */
        const getFocusedTagTaxonomy = () => {
            return $cluster_tags[cpt_focused_tag_taxonomy_index];
        }

        /**
         * Handles the tag renamer keydown event.
         * @param {KeyboardEvent} event
         */
        const handleTagRenamerKeyDown = (event) => {
            if (event.key === "Enter") {
                event.preventDefault();
                renameFocusedTagTaxonomy();
                return;
            }

            if (event.key === "Escape") {
                resetTaxonomyRenamerState();
                return;
            }
        }

        /**
         * Handles the tag renamer keyup event.
         * @param {KeyboardEvent} event
         */
        const handleTagRenamerKeyUp = (event) => {
            event.preventDefault();

            if (the_taxonomy_renamer_input.validationMessage !== "") {
                the_taxonomy_renamer_input.setCustomValidity("");
            }

            taxonomy_renamer_is_ready = the_taxonomy_renamer_input.checkValidity() && the_taxonomy_renamer_input.value !== "";
            console.log("Is tag renamer name valid: ", taxonomy_renamer_is_ready);
        }

        /**
         * Renames the focused tag taxonomy with the value of the tag renamer input.
         * @returns {Promise<void>}
         */
        const renameFocusedTagTaxonomy = async () => {
            let focused_taxonomy_tags = getFocusedTagTaxonomy();

            if (focused_taxonomy_tags == null) {
                console.error("The focused tag index is out of bounds");
                return;
            }

            let new_tag_name = the_taxonomy_renamer_input.value.trim();
            
            if (new_tag_name === focused_taxonomy_tags.Taxonomy.Name) {
                resetTaxonomyRenamerState();
                return;
            }

            if (!taxonomy_renamer_is_ready) {
                the_taxonomy_renamer_input.reportValidity();
                return;
            }

            let tag_name_available = !checkTaxonomyNameExists(new_tag_name);

            if (!tag_name_available) {
                let labeled_error = new LabeledError("In TagGroup.handleTagRenamerKeyDown", `The value '${new_tag_name}' already exists for this attribute`, lf_errors.ERR_FORBIDDEN_DUPLICATE);

                labeled_error.alert();
                return;
            }

            const old_name = focused_taxonomy_tags.Taxonomy.Name;

            let renamed = await renameTagTaxonomy(focused_taxonomy_tags.Taxonomy.UUID, new_tag_name);

            if (!renamed) {
                let labeled_error = new LabeledError("In TagGroup.handleTagRenamerKeyDown", `An error occurred while renaming the attribute '${old_name}' to '${new_tag_name}'`, lf_errors.ERR_PROCESSING_ERROR);

                labeled_error.alert();
            } else {
                emitTaxonomyContentChange(focused_taxonomy_tags.Taxonomy.UUID);
                emitPlatformMessage(`The attribute '${old_name}' was successfully renamed to '${new_tag_name}'`);
            }

            resetTaxonomyRenamerState();
        }

        /**
         * Resets the tag renamer related state.
         */
        const resetTaxonomyRenamerState = () => {
            taxonomy_renamer_is_ready = false;
            setTagRenamerState(false);
        }

        /**
         * Set tag renamer state.
         * @param {boolean} state
         */
        const setTagRenamerState = state => {
            cpt_focused_tag_taxonomy_rename_mode = state;
        }

    /*=====  End of Methods  ======*/

</script>

<div id="cluster-public-tags">
    <aside class="cpt-taxonomy-tags-minimap">
        <header class="cpt-ttm-header">
            <h4>
                {ui_taxonomy_reference.EntityNamePlural}
            </h4>
        </header>
        <ol class="cpt-ttm-tag-taxonomies">
            {#each $cluster_tags as taxonomy_tags, h (taxonomy_tags.Taxonomy.UUID)}
                {@const is_keyboard_focused = cpt_focused_tag_taxonomy_index === h && component_hotkey_context.Active}
                <li  class="cpt-ttm-taxonomy"
                    class:keyboard-focused={is_keyboard_focused}
                    class:has-hotkey-control={is_keyboard_focused && taxonomy_tags_hotkeys_context.Active}
                >
                    {#if is_keyboard_focused && cpt_focused_tag_taxonomy_rename_mode}
                        <input class="rename-input cpt-ttm-taxonomy-name"
                            bind:this={the_taxonomy_renamer_input}
                            value="{taxonomy_tags.Taxonomy.Name}"
                            on:keydown={handleTagRenamerKeyDown}
                            on:keyup={handleTagRenamerKeyUp}
                            type="text"
                            minlength="2"
                            maxlength="64"
                            pattern="{'[A-z_][A-z_\\s]{2,64}'}"
                            spellcheck="true"
                            autofocus
                            required
                        />
                    {:else}
                        <p class="cpt-ttm-taxonomy-name">
                            {taxonomy_tags.Taxonomy.Name}
                        </p>
                    {/if}
                </li>
            {/each}
        </ol>
    </aside>
    <div id="cpt-sections-wrapper"
        class="dungeon-scroll" 
    >
        {#each $cluster_tags as taxonomy_tags, h (taxonomy_tags.Taxonomy.UUID)}
            {@const is_keyboard_focused = cpt_focused_tag_taxonomy_index === h && component_hotkey_context.Active}
            <TaxonomyTags 
                taxonomy_tags={taxonomy_tags}
                has_hotkey_control={is_keyboard_focused && taxonomy_tags_hotkeys_context.Active}
                is_keyboard_focused={is_keyboard_focused}
                component_hotkey_context={taxonomy_tags_hotkeys_context}
                ui_entity_reference={ui_entity_reference}
                ui_taxonomy_reference={ui_taxonomy_reference}
                ui_tag_reference={ui_tag_reference}
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
        grid-template-columns: 8em 1fr;
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

            & h4::first-letter {
                font-family: var(--font-read);
                text-transform: uppercase;
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

            & li.cpt-ttm-taxonomy p, & li.cpt-ttm-taxonomy input {
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

            & li.cpt-ttm-taxonomy input.rename-input {
                display: block;
                text-transform: none;
                width: 100%;
                line-height: 1;
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