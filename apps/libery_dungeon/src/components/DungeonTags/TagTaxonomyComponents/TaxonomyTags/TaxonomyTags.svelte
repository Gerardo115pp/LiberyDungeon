<script>
    import { createDungeonTag, deleteDungeonTag, renameDungeonTag } from "@models/DungeonTags";
    import TagGroup from "../../Tags/TagGroup.svelte";
    import { LabeledError,  VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { createEventDispatcher, onDestroy, onMount, tick } from "svelte";
    import { confirmPlatformMessage, emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import generateTaxonomyTagsHotkeysContext, { taxonomy_tags_actions } from "./taxonomy_tags_hotkeys";
    import { browser } from "$app/environment";
    import { CursorMovementWASD } from "@common/keybinds/CursorMovement";
    import { SearchResultsWrapper } from "@app/common/keybinds/CommonActionWrappers";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { wrapShowHotkeysTable } from "@app/common/keybinds/CommonActionWrappers";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import { last_keyboard_focused_tag } from "./taxonomy_tags_store";
    import { common_action_groups } from "@app/common/keybinds/CommonActionsName";
    import { writable } from "svelte/store";

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
             * The taxonomy tags component hotkey context.
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
             */
            export let component_hotkey_context = generateTaxonomyTagsHotkeysContext();

            /*=============================================
            =            Hotkeys state            =
            =============================================*/

                /**
                 * Whether the component has mounted or not.
                 * @type {boolean}
                 */
                let has_taxonomy_tag_mounted = false;

                /**
                 * Whether to enable the tag renamer for the focused tag.
                 * @type {boolean}
                 */
                let renaming_focused_tag = false;
            
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
         * A string used to filter dungeon tags search results.
         * @type {import('svelte/store').Writable<string>})}
         */
        let search_query = writable("");
    
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
        let the_tag_group_component;

        /**
         * The grid navigation wrapper.
         * @type {CursorMovementWASD | null}
         */
        let the_grid_navigation_wrapper = null;

        /**
         * The dungeon tag search results hotkey wrapper.
         * @type {SearchResultsWrapper<import('@models/DungeonTags').DungeonTag> | null}
         */
        let the_dungeon_tag_search_results_wrapper = null;
        $: if (the_dungeon_tag_search_results_wrapper != null && (taxonomy_tags?.Tags.length ?? 0) > 0) {
            the_dungeon_tag_search_results_wrapper.updateSearchPool(taxonomy_tags.Tags);
        }
        
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
        has_taxonomy_tag_mounted = true;
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

            /* ------------------------------ Hotkey Setup ------------------------------ */
        
                /**
                 * Defines the tools hotkeys.
                 */ 
                const defineDesktopKeybinds = () => {
                    if (global_hotkeys_manager == null) {
                        console.error("Global hotkeys manager is not available");
                        return;
                    }

                    if (global_hotkeys_manager.hasContext(component_hotkey_context.HotkeysContextName)) {
                        global_hotkeys_manager.dropContext(component_hotkey_context.HotkeysContextName);
                    }

                    if (component_hotkey_context.HasGeneratedHotkeysContext()) {
                        component_hotkey_context.dropHotkeysContext();
                    }

                    const hotkeys_context = preparePublicHotkeyActions(component_hotkey_context);

                    setGridNavigationWrapper(hotkeys_context);

                    setSearchResultsWrapper(hotkeys_context);

                    wrapShowHotkeysTable(hotkeys_context);

                    global_hotkeys_manager.declareContext(component_hotkey_context.HotkeysContextName, hotkeys_context);
                    
                    global_hotkeys_manager.loadContext(component_hotkey_context.HotkeysContextName);
                }

                /**
                 * Prepares the public hotkey actions to generate the hotkeys context.
                 * @param {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext} new_component_hotkey_context
                 * @returns {import('@libs/LiberyHotkeys/hotkeys_context').default}
                 */
                const preparePublicHotkeyActions = (new_component_hotkey_context) => {

                    /* -------------------------- Public hotkey actions ------------------------- */
                        /* --------------------------- DROP_HOTKEY_CONTEXT -------------------------- */

                            let drop_hotkey_context_action = new_component_hotkey_context.getHotkeyActionOrPanic(taxonomy_tags_actions.DROP_HOTKEY_CONTEXT);

                            let drop_hotkey_context_description = `<${HOTKEYS_GENERAL_GROUP}> Deselects the active ${ui_taxonomy_reference.EntityName}`;

                            if (drop_hotkey_context_action.HasNullishCallback()) {
                                drop_hotkey_context_action.Callback = handleHotkeyControlDrop;
                            }

                            if (drop_hotkey_context_action.HasNullishDescription()) {
                                drop_hotkey_context_action.overwriteDescription(drop_hotkey_context_description);
                            }

                        /* ---------------------------- FOCUS_TAG_CREATOR --------------------------- */

                            const focus_tag_creator_action = new_component_hotkey_context.getHotkeyActionOrPanic(taxonomy_tags_actions.FOCUS_TAG_CREATOR);

                            focus_tag_creator_action.Callback = handleFocusTagCreator;

                            if (focus_tag_creator_action.HasNullishDescription()) {
                                focus_tag_creator_action.overwriteDescription(`<content>Focuses the ${ui_tag_reference.EntityNamePlural} creator input.`);
                            }

                            if (focus_tag_creator_action.OverwrittenOptions) {
                                throw new Error(`Something overwrote the options of the focus tag creator action. this is not allowed.`);
                            }

                        /* ------------------------------ RENAME_FOCUSED_TAG ------------------------------ */

                            const rename_focused_tag_action = new_component_hotkey_context.getHotkeyActionOrPanic(taxonomy_tags_actions.RENAME_FOCUSED_TAG);

                            rename_focused_tag_action.Callback = handlesTagRenamerHotkey;

                            if (rename_focused_tag_action.HasNullishDescription()) {
                                rename_focused_tag_action.overwriteDescription(`<content>Renames the focused ${ui_tag_reference.EntityName}.`);
                            }

                            rename_focused_tag_action.panicIfOptionsOverwritten();

                        /* --------------------------- SELECT_FOCUSED_TAG --------------------------- */

                            const select_focused_tag_action = new_component_hotkey_context.getHotkeyActionOrPanic(taxonomy_tags_actions.SELECT_FOCUSED_TAG);

                            select_focused_tag_action.Callback = handleSelectFocusedTag;

                            if (select_focused_tag_action.HasNullishDescription()) {
                                select_focused_tag_action.overwriteDescription(`<content>Selects the focused ${ui_tag_reference.EntityName}.`);
                            }

                            select_focused_tag_action.panicIfOptionsOverwritten();

                        /* --------------------------- DELETE_FOCUSED_TAG --------------------------- */

                            const delete_focused_tag_action = component_hotkey_context.getHotkeyActionOrPanic(taxonomy_tags_actions.DELETE_FOCUSED_TAG);

                            delete_focused_tag_action.Callback = handleDeleteFocusedTag;

                            if (delete_focused_tag_action.HasNullishDescription()) {
                                delete_focused_tag_action.overwriteDescription(`<content>Deletes the focused ${ui_tag_reference.EntityName}.`);
                            }

                            delete_focused_tag_action.panicIfOptionsOverwritten();

                    /* -------------------------------------------------------------------------- */

                    const hotkeys_context = component_hotkey_context.generateHotkeysContext();

                    return hotkeys_context;
                }

            /* -------------------------------------------------------------------------- */

            /**
             * Drops the component hotkey context
             */
            const dropHotkeyContext = () => {
                if (global_hotkeys_manager == null) {
                    console.error("Global hotkeys manager is not available");
                    return;
                };

                if (!global_hotkeys_manager.hasContext(component_hotkey_context.HotkeysContextName)) return;

                global_hotkeys_manager.dropContext(component_hotkey_context.HotkeysContextName);
            }

            /**
             * Emits an event to drop the hotkeys context
             */
            const emitDropHotkeyContext = () => {
                dispatch("drop-hotkeys-control"); 
            }

            /**
             * Handles the Cursor update event emitted by the_grid_navigation_wrapper.
             * @type {import("@common/keybinds/CursorMovement").CursorPositionCallback}
             */
            const handleCursorUpdate = (cursor_wrapped_value) => {
                focused_tag_index = cursor_wrapped_value.value;

                /** @type {import('@models/DungeonTags').DungeonTag | undefined}*/
                let focused_tag = getFocusedTag();

                if (focused_tag == null) {
                    let variable_environment = new VariableEnvironmentContextError("In TaxonomyTags.handleCursorUpdate");

                    variable_environment.addVariable("cursor_wrapped_value", cursor_wrapped_value);
                    variable_environment.addVariable("focused_tag_index", focused_tag_index);

                    const user_message = `Failed to get focused ${ui_tag_reference.EntityName} from index '${focused_tag_index}', This is a programming error. Am very sorry for the inconvenience.`;

                    let labeled_err = new LabeledError(variable_environment, user_message, lf_errors.PROGRAMMING_ERROR__BROKEN_STATE);

                    labeled_err.alert();
                    return;
                }

                last_keyboard_focused_tag.set(focused_tag);

                return false;
            }

            /**
             * Handles the tag renamer hotkey.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handlesTagRenamerHotkey = (event, hotkey) => {
                if (!has_hotkey_control) return;

                if (taxonomy_tags.Tags.length === 0) return;

                setTagRenamerState(true);
            }

            /**
             * Handles the focus tag creator hotkey.
             */
            const handleFocusTagCreator = () => {
                if (!has_hotkey_control) return;

                the_tag_group_component.focusTagCreator();
            }

            /**
             * Handles the focus search match hotkey.
             * @type {import("@common/keybinds/CommonActionWrappers").SearchResultsUpdateCallback<import('@models/DungeonTags').DungeonTag>}
             */
            const handleFocusSearchMatch = (search_result) => {
                let tag_id = search_result.Id;

                for (let h = 0; h < taxonomy_tags.Tags.length; h++) { 
                    if (taxonomy_tags.Tags[h].Id === tag_id) {
                        the_grid_navigation_wrapper?.updateCursorPosition(h);
                        break;
                    }
                }
            }

            /**
             * Handles the update of the search query label.
             * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCaptureCallback}
             */
            const handleSearchQueryUpdate = (event, captured_string) => {
                if (!has_hotkey_control) return;

                search_query.set(captured_string);
            }

            /**
             * Handles the select focused tag hotkey.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleSelectFocusedTag = (event, hotkey) => {
                if (!has_hotkey_control) return;

                if (taxonomy_tags.Tags.length === 0) return;

                let focused_tag_id = taxonomy_tags.Tags[focused_tag_index]?.Id;

                if (focused_tag_id == null) return;

                the_tag_group_component.emitTagSelected(focused_tag_id);
            }

            /**
             * Handles the deletion of the focused tag.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleDeleteFocusedTag = (event, hotkey) => {
                if (!has_hotkey_control) return;

                if (taxonomy_tags.Tags.length === 0) return;

                let focused_tag_id = taxonomy_tags.Tags[focused_tag_index]?.Id;

                if (focused_tag_id == null) return;

                if (focused_tag_index === taxonomy_tags.Tags.length - 1) {
                    focused_tag_index--;
                }

                the_tag_group_component.emitTagDeleted(focused_tag_id);
            }

            /**
             * Emits an event to close the section and drops the hotkeys context.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleHotkeyControlDrop = (event, hotkey) => {
                dropSearchResultsState();
                dropGridNavigationWrapper();

                resetHotkeyContext();
                emitDropHotkeyContext();
            }

            /**
             * Drops the tools hotkey contexts and loads the previous context.
             */
            const resetHotkeyContext = () => {
                if (global_hotkeys_manager == null) {
                    console.error("Global hotkeys manager is not available");
                    return;
                }

                if (global_hotkeys_manager.ContextName !== component_hotkey_context.HotkeysContextName) return; 

                global_hotkeys_manager.loadPreviousContext();
            }

            /**
             * Sets the grid navigation wrapper required data.
             * @param {import("@libs/LiberyHotkeys/hotkeys_context").default} hotkeys_context
             */
            const setGridNavigationWrapper = (hotkeys_context) => {
                if (!browser) return;

                if (the_grid_navigation_wrapper != null) {
                    the_grid_navigation_wrapper.destroy();
                }


                const tags_parent_selector = `#taxonomy-tags-${taxonomy_tags.Taxonomy.UUID} ol#tag-group-${taxonomy_tags.Taxonomy.UUID}`;

                const matching_elements_count = document.querySelectorAll(tags_parent_selector).length;
                if (matching_elements_count !== 1) {
                    throw new Error(`tag parent selector '${tags_parent_selector}' returned ${matching_elements_count}, expected exactly 1`);
                }


                the_grid_navigation_wrapper = new CursorMovementWASD(tags_parent_selector, handleCursorUpdate, {
                    initial_cursor_position: focused_tag_index,
                    sequence_item_name: ui_tag_reference.EntityName,
                    sequence_item_name_plural: ui_tag_reference.EntityNamePlural,
                    grid_member_selector: 'li:not(:has(input))',
                });
                the_grid_navigation_wrapper.setup(hotkeys_context);
                

                // @ts-ignore
                globalThis.the_grid_navigation_wrapper = the_grid_navigation_wrapper; 
            }

            /**
             * Sets the search results wrapper required data.
             * @param {import("@libs/LiberyHotkeys/hotkeys_context").default} hotkeys_context
             */
            const setSearchResultsWrapper = (hotkeys_context) => {
                the_dungeon_tag_search_results_wrapper = new SearchResultsWrapper(hotkeys_context, taxonomy_tags.Tags, handleFocusSearchMatch, {
                    search_hotkey: ["f"],
                    ui_search_result_reference: ui_tag_reference,
                    search_hotkey_handler: handleSearchQueryUpdate,
                    minimum_similarity: 0.8,
                });
            }

        /*=====  End of Keybinds  ======*/

        /**
         * Drops the grid navigation wrapper if it exists.
         */
        const dropGridNavigationWrapper = () => {
            if (the_grid_navigation_wrapper != null) {
                the_grid_navigation_wrapper.destroy();
            }
        }

        /**
         * Drops the search results state.
         */
        const dropSearchResultsState = () => {
            if (the_dungeon_tag_search_results_wrapper != null) {
                the_dungeon_tag_search_results_wrapper.dropSearchState();
            }

            console.log("Dropped search results state");
            search_query.set("");
            console.log("Search query: ", search_query);
        }

        /**
         * Ensures that if the element is keyboard focused, it is visible in the scroll container.
         */
        const ensureTaxonomyTagVisible = async () => {
            await tick();

            if (!the_taxonomy_tags_section || !is_keyboard_focused) return;

            the_taxonomy_tags_section.scrollIntoView({behavior: "smooth", block: "center"});
        }

        /**
         * Emits an event that should be interpreted as 'the tag taxonomy content has changed'. The taxonomy emits an event with a detail.taxonomy, this
         * contains the tag taxonomy uuid. 
         */
        const emitTaxonomyContentChange = () => {
            dispatch("taxonomy-content-change", {taxonomy: taxonomy_tags.Taxonomy.UUID});
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
         * Returns the current focused tag.
         * @returns {import('@models/DungeonTags').DungeonTag | undefined}
         */
        const getFocusedTag = () => {
            return taxonomy_tags.Tags[focused_tag_index];
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

                const labeled_err = new LabeledError(variable_environment, `Failed to create ${ui_tag_reference.EntityName} '${event?.detail?.tag_name}'`, lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();
                return;
            }

            emitTaxonomyContentChange();
        }

        /**
         * Handles the tag deleted event.
         * @param {CustomEvent<{tag_id: number}>} event
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
                message_title: `Delete ${ui_tag_reference.EntityName} '${tag_name}'`,
                question_message: `Are you sure you want to delete the ${ui_tag_reference.EntityName} '${tag_name}'? it will be disassociated from all ${ui_entity_reference.EntityNamePlural} and any other entity it is related to.`,
                danger_level: 1,
                cancel_label: "cancel",
                confirm_label: "Delete it",
                auto_focus_cancel: true,
            });

            if (user_choice !== 1) return;

            let tag_deleted = await deleteDungeonTag(tag_id);

            if (!tag_deleted) {
                const labeled_err = new LabeledError("In TaxonomyTags.handleTagDeleted", `Failed to delete ${ui_tag_reference.EntityName} with id '${tag_id}'`, lf_errors.ERR_PROCESSING_ERROR);
                labeled_err.alert();
                return;
            }

            emitTaxonomyContentChange();
        }

        /**
         * Handles the tag-renamed event.
         * @param {CustomEvent<{tag_id: number, new_tag_name: string}>} event
         */
        const handleTagRenamed = async event => {
            const tag_id = event?.detail?.tag_id;
            const new_name = event?.detail?.new_tag_name;

            if (tag_id == null || new_name == null) return;

            const tag_name = getTagNameByID(tag_id);
            if (tag_name === "") {
                console.error(`Tag id: '${tag_id}' did not match an dungeon tag withing ${taxonomy_tags.Taxonomy.Name}`);
                return;
            }

            if (tag_name === new_name) {
                setTagRenamerState(false);
                return;
            }

            setTagRenamerState(false);

            console.log(`Renaming tag '${tag_name}' to '${new_name}'`);

            const tag_renamed = await renameDungeonTag(tag_id, new_name);

            if (tag_renamed) {
                emitPlatformMessage(`Renamed '${tag_name}' to '${new_name}'`);
            } if (!tag_renamed) {
                const labeled_err = new LabeledError("In TaxonomyTags.handleTagRenamed", `Failed to rename ${ui_tag_reference.EntityName} with id '${tag_id}'`, lf_errors.ERR_PROCESSING_ERROR);
                labeled_err.alert();
                return;
            }

            emitTaxonomyContentChange();
        }

        /**
         * Handles the tag-rename-cancelled event.
         */
        const handleTagRenameCancelled = () => {
            setTagRenamerState(false);
        }

        /**
         * Set tag renamer state.
         * @param {boolean} state
         */
        const setTagRenamerState = state => {
            renaming_focused_tag = state;
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<section class="dungeon-taxonomy-content"
    bind:this={the_taxonomy_tags_section}
    id="taxonomy-tags-{taxonomy_tags.Taxonomy.UUID}"
    class:is-keyboard-focused={is_keyboard_focused}
    class:hotkey-control={has_hotkey_control}
>
    <header class="taxonomy-header dungeon-properties">
        <h4>
            {taxonomy_tags.Taxonomy.Name}
        </h4>
        {#if $search_query !== "" && has_hotkey_control}
            <p class="dtc-search-query html-tag">/{$search_query}</p>
        {/if}
    </header>
    <TagGroup 
        bind:this={the_tag_group_component}
        tag_group_id="tag-group-{taxonomy_tags.Taxonomy.UUID}"
        dungeon_tags={taxonomy_tags.Tags}
        enable_tag_creator={enable_tag_creation}
        enable_keyboard_selection={has_hotkey_control}
        rename_focused_tag={renaming_focused_tag}
        focused_tag_index={focused_tag_index}
        ui_taxonomy_reference={ui_taxonomy_reference}
        ui_tag_reference={ui_tag_reference}
        on:tag-selected
        on:tag-created={handleTagCreated}
        on:tag-deleted={handleTagDeleted}
        on:tag-renamed={handleTagRenamed}
        on:tag-rename-cancelled={handleTagRenameCancelled}
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

    .taxonomy-header.dungeon-properties {
        gap: var(--spacing-1);

        & p.dtc-search-query {
            font-size: var(--font-size-1);
            color: var(--grey-1);
            line-height: 1;
        }
    }
</style>