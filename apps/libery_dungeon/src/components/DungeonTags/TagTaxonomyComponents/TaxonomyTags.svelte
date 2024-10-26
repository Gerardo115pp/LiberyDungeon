<script>
    import { createDungeonTag, deleteDungeonTag } from "@models/DungeonTags";
    import TagGroup from "../Tags/TagGroup.svelte";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { createEventDispatcher, onDestroy, onMount, tick } from "svelte";
    import { confirmPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import { toggleHotkeysSheet } from "@stores/layout";
    import { browser } from "$app/environment";
    import { linearCycleNavigationWrap } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";
    import AnchorsOne from "@components/UI/AnchorsOne.svelte";

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

            const hotkeys_context_name = "taxonomy_tags";

            /*=============================================
            =            Hotkeys state            =
            =============================================*/

                /**
                 * Whether the component has mounted or not.
                 * @type {boolean}
                 */
                let has_taxonomy_tag_mounted = false;
            
                /**
                 * Whether it has hotkey control.
                 * @type {boolean}
                 */ 
                export let has_hotkey_control = false;
                $: if (has_hotkey_control && has_taxonomy_tag_mounted) {
                    defineDesktopKeybinds();
                }
        
            /*=====  End of Hotkeys state  ======*/

            
            /*=============================================
            =            Hotkeys movement            =
            =============================================*/
            
                /**
                 * The index of the focused tag.
                 * @type {number}
                 */
                let focused_tag_index = 0;
            
            /*=====  End of Hotkeys movement  ======*/
            
            

            /**
             * Whether the TagTaxonomy component is been focused by a keyboard cursor. This does not mean that the
             * component has hotkey control.
             * @type {boolean}
             */
            export let is_keyboard_focused = false;
            $: if (is_keyboard_focused) {
                ensureTaxonomyTagVisible();
            }

        /*=====  End of Hotkeys  ======*/
    
        /**
         * The Taxonomy tags composition class.
         * @type {import("@models/DungeonTags").TaxonomyTags}
         */
        export let taxonomy_tags;
        $: console.log("TaxonomyTags: ", taxonomy_tags);

        /**
         * Whether to allow the user to create new tags.
         * @type {boolean}
         */
        export let enable_tag_creation = false;

        /**
         * The Taxonomy tags section
         * @type {HTMLElement}
         */
        let the_taxonomy_tags_section;

        /**
         * The tag group component.
         * @type {TagGroup}
         */
        let tag_group_component;

        const dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/

    onMount(() => {
        has_taxonomy_tag_mounted = true;
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
                if (global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    global_hotkeys_manager.dropContext(hotkeys_context_name);
                }
                console.log(`Defining hotkeys for taxonomy tags<${taxonomy_tags.Taxonomy.Name}>`);

                const hotkeys_context = new HotkeysContext(); 

                hotkeys_context.register(["q"], handleHotkeyControlDrop, {
                    description: `<${HOTKEYS_GENERAL_GROUP}>Deselects the selected attribute section.`, 
                    await_execution: false
                });

                hotkeys_context.register(["a", "d"], handleTagNavigation, {
                    description: "<navigation>Navigates the value.", 
                    await_execution: false
                });

                hotkeys_context.register(["\\d g"], handleTagIndexGoto, {
                    description: "<navigation>Navigates to the typed index value", 
                    await_execution: false
                });

                hotkeys_context.register(["c"], handleFocusTagCreator, {
                    description: "<navigation>Focuses the tag creator input.", 
                    await_execution: false,
                    mode: "keyup"
                });

                hotkeys_context.register(["?"], toggleHotkeysSheet, { 
                    description: `<${HOTKEYS_GENERAL_GROUP}>Opens the hotkeys cheat sheet.`
                });

                global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
                
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
             * Handles the navigation of the tags.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleTagNavigation = (event, hotkey) => {
                if (!has_hotkey_control || (event.key !== "a" && event.key !== "d")) return;

                let tag_count = taxonomy_tags.Tags.length;

                if (tag_count === 0) return;

                let navigation_step = event.key === "a" ? -1 : 1;

                focused_tag_index = linearCycleNavigationWrap(focused_tag_index, tag_count - 1, navigation_step).value;
            }

            /**
             * Handles the tag index goto hotkey.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleTagIndexGoto = (event, hotkey) => {
                if (!hotkey.WithVimMotion || !hotkey.HasMatch) return;

                let vim_motion_value = hotkey.MatchMetadata.MotionMatches[0];
                vim_motion_value--; // 1 based index to 0 based index

                let tag_count = taxonomy_tags.Tags.length;

                if (tag_count === 0) return;

                let new_focused_tag_index = Math.max(0, Math.min(tag_count -1, vim_motion_value));

                focused_tag_index = new_focused_tag_index;
            }

            /**
             * Handles the focus tag creator hotkey.
             */
            const handleFocusTagCreator = () => {
                if (!has_hotkey_control) return;

                tag_group_component.focusTagCreator();
            }

            /**
             * Emits an event to close the section and drops the hotkeys context.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleHotkeyControlDrop = (event, hotkey) => {
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
         * Ensures that if the element is keyboard focused, it is visible in the scroll container.
         */
        const ensureTaxonomyTagVisible = async () => {
            await tick();

            if (!the_taxonomy_tags_section || !is_keyboard_focused) return;

            the_taxonomy_tags_section.scrollIntoView({behavior: "smooth", block: "center"});
        }

        /**
         * Returns the name of a tag by it's given id or an empty string if not found.
         * @param {number} tag_id
         * @returns {string}
         */
        const getTagNameByID = tag_id => {
            let tag_name = "";

            for (let dungeon_tag of taxonomy_tags.Tags) {
                if (dungeon_tag.Id === tag_id) {
                    tag_name = dungeon_tag.Name;
                    break;
                }
            }

            return tag_name;
        }

        /**
         * Handles the tag created event.
         * @param {CustomEvent<{tag_name: string}>} event
         */ 
        const handleTagCreated = async event => {
            /**
             * @type {import('@models/DungeonTags').DungeonTag | null}
             */
            let new_dungeon_tag = await createDungeonTag(event.detail.tag_name, taxonomy_tags.Taxonomy.UUID);

            if (new_dungeon_tag === null) {
                const variable_environment = new VariableEnvironmentContextError("In TaxonomyTags.handleTagCreated")
                variable_environment.addVariable("triggering_event", event);
                variable_environment.addVariable("taxonomy_tags.Taxonomy.UUID", taxonomy_tags.Taxonomy.UUID);

                const labeled_err = new LabeledError(variable_environment, `Failed to create tag '${event?.detail?.tag_name}'`);

                labeled_err.alert();
                return;
            }

            emitTaxonomyContentChange();
        }

        /**
         * Handles the tag deleted event.
         * @param {CustomEvent<{tag_id: string}>} event
         */
        const handleTagDeleted = async event => {
            
            let tag_id = event?.detail?.tag_id;

            if (tag_id == null) return;

            const tag_name = getTagNameByID(tag_id);

            if (tag_name === "") {
                console.error(`Tag id: '${tag_id}' did not match an dungeon tag withing ${taxonomy_tags.Taxonomy.Name}`);
                return;
            }

            const user_choice = await confirmPlatformMessage({
                message_title: `Delete value '${tag_name}'`,
                question_message: `Are you sure you want to delete '${tag_name}'? it will be disassociated from all medias and categorys it is set to.`,
                danger_level: 1,
                cancel_label: "cancel",
                confirm_label: "Delete it",
                auto_focus_cancel: true,
            });

            if (user_choice !== 1) return;

            let tag_deleted = await deleteDungeonTag(tag_id);

            if (!tag_deleted) {
                const labeled_err = new LabeledError("In TaxonomyTags.handleTagDeleted", `Failed to delete tag with id '${tag_id}'`);
                labeled_err.alert();
                return;
            }

            emitTaxonomyContentChange();
        }

        /**
         * Emits an event that should be interpreted as 'the tag taxonomy content has changed'. The taxonomy emits an event with a detail.taxonomy, this
         * contains the tag taxonomy uuid. 
         */
        const emitTaxonomyContentChange = () => {
            dispatch("taxonomy-content-change", {taxonomy: taxonomy_tags.Taxonomy.UUID});
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<section class="dungeon-taxonomy-content"
    bind:this={the_taxonomy_tags_section}
    class:is-keyboard-focused={is_keyboard_focused}
    class:hotkey-control={has_hotkey_control}
>
    <header class="taxonomy-header">
        <h4>
            {taxonomy_tags.Taxonomy.Name}
        </h4>
    </header>
    <TagGroup 
        bind:this={tag_group_component}
        dungeon_tags={taxonomy_tags.Tags}
        enable_tag_creator={enable_tag_creation}
        enable_keyboard_selection={has_hotkey_control}
        focused_tag_index={focused_tag_index}
        on:tag-selected
        on:tag-created={handleTagCreated}
        on:tag-deleted={handleTagDeleted}
    />
</section>

<style>
    .dungeon-taxonomy-content {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-1);
    }

    .dungeon-taxonomy-content.is-keyboard-focused {
        & header.taxonomy-header h4 {
            color: var(--main);
        }
    }

    .dungeon-taxonomy-content.hotkey-control {
        & header.taxonomy-header > h4 {
            color: var(--main-dark);
        }
    }
</style>