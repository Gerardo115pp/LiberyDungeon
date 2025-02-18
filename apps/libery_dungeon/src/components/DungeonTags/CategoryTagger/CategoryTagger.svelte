<script>
    /*=============================================
    =            Imports            =
    =============================================*/
        import { deleteTagTaxonomy, getEntityTaggings, getTaxonomyTagsByUUID, tagCategory, tagCategoryContent, tagMedia, untagCategoryContent, untagEntity } from "@models/DungeonTags";
        import { cluster_tags, last_cluster_domain, refreshClusterTagsNoCheck } from "@stores/dungeons_tags";
        import TagTaxonomyCreator from "../TagTaxonomyComponents/TagTaxonomyCreator.svelte";
        import { createEventDispatcher, onDestroy, onMount } from "svelte";
        import { current_cluster } from "@stores/clusters";
        import { current_category } from "@stores/categories_tree";
        import { LabeledError, UIReference, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
        import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
        import ClusterPublicTags from "../TagTaxonomyComponents/ClusterPublicTags.svelte";
        import CategoryTaggings from "./sub-components/CategoryTaggings.svelte";
        import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
        import HotkeysContext, { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
        import generateTaxonomyTagsHotkeysContext, { taxonomy_tags_actions } from "../TagTaxonomyComponents/TaxonomyTags/taxonomy_tags_hotkeys";
        import { last_keyboard_focused_tag } from "../TagTaxonomyComponents/TaxonomyTags/taxonomy_tags_store";
        import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import { wrapShowHotkeysTable } from "@app/common/keybinds/CommonActionWrappers";
        import { common_action_groups } from "@app/common/keybinds/CommonActionsName";
        import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
        import generateCategoryTaggerHotkeyContext, { category_tagger_child_contexts } from "./category_tagger_hotkeys";
    /*=====  End of Imports  ======*/
    
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
             * the component hotkey context for the categories tagger component
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
             */
            export let component_hotkey_context = generateCategoryTaggerHotkeyContext();


            /* -------------------------- sub-component context -------------------------- */

                /**
                 * The component hotkeys context for the child taxonomy tags components.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let taxonomy_tags_hotkeys_context = null;

                /**
                 * The component hotkey context for the tag taxonomy creator.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let tag_taxonomy_creator_hotkeys_context = component_hotkey_context.ChildHotkeysContexts.get(category_tagger_child_contexts.TAG_TAXONOMY_CREATOR) ?? null;
                if (tag_taxonomy_creator_hotkeys_context == null) {
                    throw new Error("In CategoryTagger, invalid component_hotkey_context: TagTaxonomyCreator hotkeys context was not defined as a child context of the CategoryTagger context.");
                }

                /**
                 * The component hotkey context for the category taggings component.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let category_taggings_hotkeys_context = component_hotkey_context.ChildHotkeysContexts.get(category_tagger_child_contexts.CATEGORY_TAGGINGS) ?? null;
                if (category_taggings_hotkeys_context == null) {
                    throw new Error("In CategoryTagger, invalid component_hotkey_context: CategoryTaggings hotkeys context was not defined as a child context of the CategoryTagger context.");
                }

                /**
                 * The component hotkey context for the cluster public tags component.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let cluster_public_tags_hotkeys_context = component_hotkey_context.ChildHotkeysContexts.get(category_tagger_child_contexts.CLUSTER_PUBLIC_TAGS) ?? null;
                if (cluster_public_tags_hotkeys_context == null) {
                    throw new Error("In CategoryTagger, invalid component_hotkey_context: ClusterPublicTags hotkeys context was not defined as a child context of the CategoryTagger context.");
                }

        /*=====  End of Hotkeys  ======*/
    
        /**
         * Whether the cluster_tags correctness with respect to the current_cluster has been checked.
         * @type {boolean}
         */ 
        let cluster_tags_checked = false;

        /**
         * Current category taggings.
         * @type {import("@models/DungeonTags").DungeonTagging[]}
         */
        let current_category_taggings = [];

        /**
         * The category tagger tool node.
         * @type {HTMLDialogElement | null}
         */
        let the_category_tagger_tool = null;
        

        /*----------  UI references  ----------*/
        
            /**
             * A UiReference object, to create ui messages about the taggable entity. used to pass down to general purpose components
             * @type {UIReference}
             */
            const ui_entity_reference = new UIReference("category", "categories");

            /**
             * A UiReference object, to create ui messages about the tag taxonomy. used to pass down to general purpose components.
             * @type {UIReference}
             */
            const ui_taxonomy_reference = new UIReference("attribute", "attributes");

            /**
             * A UiReference object, to create ui messages about the dungeon-tag. used to pass down to general purpose components.
             * @type {UIReference}
             */
            const ui_tag_reference = new UIReference("value", "values");
        
        /*----------  Component metadata  ----------*/
        
            /**
             * The count of dtt-sections in the category tagger.
             * @type {number}
             */
            let dtt_sections_count;
        
        /*----------  Hotkey state  ----------*/
        
            /**
             * The focused section indicator.
             * @type {number}
             */ 
            let ct_focused_section = 0;

            /**
             * Sub-component sections.
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext[]}
             */
            const sub_component_sections = [
                tag_taxonomy_creator_hotkeys_context,
                category_taggings_hotkeys_context,
                cluster_public_tags_hotkeys_context
            ];

        let current_cluster_unsubscriber = () => {};
        let current_category_unsubscriber = () => {};

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        defineComponentMetadata();

        if ($current_cluster === null) {
            current_cluster_unsubscriber = current_cluster.subscribe(async (value) => {
                if (value !== null) {
                    await verifyLoadedClusterTags();

                    current_cluster_unsubscriber();

                    current_cluster_unsubscriber = () => {}
                }
            });
        } else {
            await verifyLoadedClusterTags();
        }

        current_category_unsubscriber = current_category.subscribe(handleCurrentCategoryChange)

        defineSubComponentsHotkeysContext();

        defineDesktopKeybinds();
    });

    onDestroy(() => {
        resetHotkeyContext();
        current_category_unsubscriber();
        current_cluster_unsubscriber();
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
                if (global_hotkeys_manager == null) {
                    console.error("Hotkeys manager not available.");
                    return;
                };

                if (global_hotkeys_manager.hasContext(component_hotkey_context.HotkeysContextName)) {
                    global_hotkeys_manager.dropContext(component_hotkey_context.HotkeysContextName);
                }

                if (component_hotkey_context.HasGeneratedHotkeysContext()) {
                    component_hotkey_context.dropHotkeysContext();
                }

                const hotkeys_context = component_hotkey_context.generateHotkeysContext();

                hotkeys_context.register(["q", "t"], handleCloseCategoryTaggerTool, {
                    description: `<${HOTKEYS_GENERAL_GROUP}>Closes the category tagger tool.`,
                    await_execution: false
                });

                hotkeys_context.register(["w", "s"], handleSectionFocusNavigation, {
                    description: `<navigation>Moves up and down through the tool's sections.`,
                    await_execution: false
                });

                hotkeys_context.register(["e"], handleSectionSelection, {
                    description: `<navigation>Selects the current section.`,
                    await_execution: false
                });

                wrapShowHotkeysTable(hotkeys_context);

                global_hotkeys_manager.declareContext(component_hotkey_context.HotkeysContextName, hotkeys_context);

                global_hotkeys_manager.loadContext(component_hotkey_context.HotkeysContextName);
            }

            /**
             * Emits an event to close the category tagger tool and drops the hotkeys context.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleCloseCategoryTaggerTool = (event, hotkey) => {
                resetHotkeyContext();
                emitCloseCategoryTagger();
            }

            /**
             * Handles the section focus navigation.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleSectionFocusNavigation = (event, hotkey) => {
                if (sub_component_sections[ct_focused_section].Active) return;

                let new_focused_section = ct_focused_section;

                const navigation_step = event.key === "w" ? -1 : 1;

                new_focused_section += navigation_step;

                if (new_focused_section < 0) {
                    new_focused_section = dtt_sections_count - 1;
                } else if (new_focused_section >= dtt_sections_count) {
                    new_focused_section = 0;
                }

                ct_focused_section = new_focused_section;
            }

            /**
             * Handles the selection of tools sections.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleSectionSelection = (event, hotkey) => {
                event.preventDefault();

                sub_component_sections[ct_focused_section].Active = true;
            }

            /**
             * Recovers the hotkeys control and deactivates the active section.
             */
            const handleRecoverHotkeysControl = () => {
                sub_component_sections[ct_focused_section].Active = false;
            }

            /**
             * Drops the tools hotkey contexts and loads the previous context.
             */
            const resetHotkeyContext = () => {
                if (global_hotkeys_manager == null) return;

                if (global_hotkeys_manager.ContextName !== component_hotkey_context.HotkeysContextName) return; 

                global_hotkeys_manager.loadPreviousContext();
            }

            /* --------------------- sub-components hotkeys context --------------------- */

                /**
                 * Defines the hotkeys context for sub-components.
                 */
                const defineSubComponentsHotkeysContext = () => {   
                    taxonomy_tags_hotkeys_context = defineTaxonomyTagsHotkeysContext();
                }

                /**
                 * defines the hotkeys context for the taxonomy tags components.
                 * @requires generateTaxonomyTagsHotkeysContext
                 * @returns {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
                 */
                const defineTaxonomyTagsHotkeysContext = () => {
                    const taxonomy_tags_context = generateTaxonomyTagsHotkeysContext();

                    const alt_select_focused_tag_action = taxonomy_tags_context.getHotkeyAction(taxonomy_tags_actions.ALT_SELECT_FOCUSED_TAG);

                    if (alt_select_focused_tag_action != null && alt_select_focused_tag_action.OverwriteBehavior === ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE) {
                        alt_select_focused_tag_action.overwriteDescription(`${common_action_groups.CONTENT}Applies the focused ${ui_tag_reference.EntityName} to the current ${ui_entity_reference.EntityName} media content.`);

                        alt_select_focused_tag_action.Callback = handleApplyFocusedTagToContent;
                    }

                    const alt_delete_focused_tag_action = taxonomy_tags_context.getHotkeyAction(taxonomy_tags_actions.ALT_DELETE_FOCUSED_TAG);

                    if (alt_delete_focused_tag_action != null && alt_delete_focused_tag_action.OverwriteBehavior === ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE) {
                        alt_delete_focused_tag_action.overwriteDescription(`${common_action_groups.CONTENT}Removes the focused ${ui_tag_reference.EntityName} from the current ${ui_entity_reference.EntityName} media content.`);

                        alt_delete_focused_tag_action.Callback = handleRemoveFocusedTagFromContent;
                    }

                    return taxonomy_tags_context;
                }

                /**
                 * Applies the focused tag on the category content medias.
                 * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback} 
                 */
                const handleApplyFocusedTagToContent = async (event, hotkey) => {
                    if ($current_category == null) {
                        console.error("In CategoryTagger.handleApplyFocusedTagToContent: No current category available while trying to apply the focused tag to the content.");
                        return;
                    }
                    
                    console.log("Applying focused tag.");

                    const dungeon_tag = $last_keyboard_focused_tag;

                    if (dungeon_tag == null) {
                        console.error("In CategoryTagger.handleApplyFocusedTagToContent: No focused tag available while trying to apply it to the content.");
                        return;
                    }

                    console.log("Focused tag: ", dungeon_tag);

                    const content_count = $current_category.content.length;

                    let successfully_applied = await tagCategoryContent($current_category.uuid, dungeon_tag.Id);

                    if (!successfully_applied) {
                        const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.handleApplyFocusedTagToContent");

                        variable_environment.addVariable("dungeon_tag.Name", dungeon_tag.Name);
                        variable_environment.addVariable("current_category.uuid", $current_category.uuid);
                        variable_environment.addVariable("content_count", content_count);
                        variable_environment.addVariable("dungeon_tag.Id", dungeon_tag.Id);

                        const labeled_error = new LabeledError(variable_environment, "Failed to apply the focused tag to the content.", lf_errors.ERR_PROCESSING_ERROR);

                        labeled_error.alert();
                        return;
                    }

                    let feedback_message =`Applied the ${ui_tag_reference.EntityName} '${dungeon_tag.Name}' to ${content_count} medias.`

                    emitPlatformMessage(feedback_message);
                }

                /**
                 * Removes the focused tag from the category content medias.
                 * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback} 
                 */
                const handleRemoveFocusedTagFromContent = async (event, hotkey) => {
                    if ($current_category == null) {
                        console.error("In CategoryTagger.handleRemoveFocusedTagFromContent: No current category available while trying to remove the focused tag from the content.");
                        return;
                    }
                    
                    console.log("Removing focused tag.");

                    const dungeon_tag = $last_keyboard_focused_tag;

                    if (dungeon_tag == null) {
                        console.error("In CategoryTagger.handleRemoveFocusedTagFromContent: No focused tag available while trying to remove it from the content.");
                        return;
                    }

                    console.log("Focused tag: ", dungeon_tag);

                    const content_count = $current_category.content.length;

                    let successfully_removed = await untagCategoryContent($current_category.uuid, dungeon_tag.Id);

                    if (!successfully_removed) {
                        const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.handleRemoveFocusedTagFromContent");

                        variable_environment.addVariable("dungeon_tag.Name", dungeon_tag.Name);
                        variable_environment.addVariable("current_category.uuid", $current_category.uuid);
                        variable_environment.addVariable("content_count", content_count);
                        variable_environment.addVariable("dungeon_tag.Id", dungeon_tag.Id);

                        const labeled_error = new LabeledError(variable_environment, "Failed to remove the focused tag from the content.", lf_errors.ERR_PROCESSING_ERROR);

                        labeled_error.alert();
                        return;
                    }

                    let feedback_message =`Removed the ${ui_tag_reference.EntityName} '${dungeon_tag.Name}' from ${content_count} medias.`

                    emitPlatformMessage(feedback_message);
                }
        
        /*=====  End of Keybinds  ======*/

        /**
         * Defines the component's content metadata.
         */
        const defineComponentMetadata = () => {
            let dtt_sections = getSectionNodes();

            if (dtt_sections == null) {
                console.error("Could not find the category tagger tool sections.");
                return;
            }

            dtt_sections_count = dtt_sections.length;
        }
        
        /**
         * Emits an event to the parent to close the category tagger tool.
         */
        const emitCloseCategoryTagger = () => {
            dispatch("close-category-tagger");
        }

        /**
         * Returns all the section nodes of the category tagger tool.
         * @returns {NodeListOf<HTMLElement> | undefined}
         */
        const getSectionNodes = () => {
            if (the_category_tagger_tool == null) return; 

            return the_category_tagger_tool.querySelectorAll(".dctt-section");
        }

        /**
         * Handles the change of the current category.
         * @param {import("@models/Categories").CategoryLeaf | null} new_category
         */
        const handleCurrentCategoryChange = async (new_category) => {
            if (new_category === null) return;

            console.log("Refreshing taggings of:", new_category);

            if (new_category === null) return;

            await refreshCurrentCategoryTaggings();
        }

        /**
         * Handles the remove-category-tag event emitted by the CategoryTaggings component.
         * @param {CustomEvent<{removed_tag: number}>} event
         */
        const handleRemoveCategoryTag = event => {
            let tag_id = event?.detail?.removed_tag;

            if (tag_id == null) return;

            removeCategoryTag(tag_id);
        }

        /**
         * Handles the tag-taxonomy-created event from the TagTaxonomyCreator component.
         */
        const handleTagTaxonomyCreated = async () => {
            await refreshClusterTagsNoCheck($current_cluster.UUID);
        }

        /**
         * Handles the delete-taxonomy event from the ClusterPublicTags component.
         * @param {CustomEvent<{taxonomy_uuid: string}>} event
         */
        const handleTagTaxonomyDeleted = async event => {
            const taxonomy_uuid = event?.detail?.taxonomy_uuid;

            if (taxonomy_uuid == null) return;

            let deleted = await deleteTagTaxonomy(taxonomy_uuid);

            if (!deleted) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.handleTagTaxonomyDeleted");

                variable_environment.addVariable("taxonomy_uuid", taxonomy_uuid);

                const labeled_error = new LabeledError(variable_environment, "Failed to delete the tag taxonomy.", lf_errors.ERR_PROCESSING_ERROR);

                labeled_error.alert();
                return;
            }

            await refreshClusterTagsNoCheck($current_cluster.UUID);
        }

        /**
         * Refreshes the content of a TagTaxonomy.
         * @param {CustomEvent<{taxonomy: string}>} event
         */
        const handleTaxonomyContentChanged = async event => {
            const taxonomy_uuid = event?.detail?.taxonomy;
            console.log("Refreshing: ", event.detail.taxonomy);
            if (taxonomy_uuid == null) return;

            const taxonomy_tags_index = $cluster_tags.findIndex(tag => tag.Taxonomy.UUID === taxonomy_uuid);
            console.log("Index: ", taxonomy_tags_index);

            let new_taxonomy_tags = await getTaxonomyTagsByUUID(taxonomy_uuid);

            if (new_taxonomy_tags === null) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.handleTaxonomyContentChanged");

                variable_environment.addVariable("taxonomy_tags_index", taxonomy_tags_index);
                variable_environment.addVariable("cluster_tags", $cluster_tags);

                const labeled_error = new LabeledError(variable_environment, "Failed to refresh the tag taxonomy content. Closing and opening the tool may solve the issue.", lf_errors.ERR_LOADING_ERROR);

                labeled_error.alert();
                return;
            }

            // Ensure the taxonomy is inserted at the same index as it was before.
            /**
             * @type {import("@models/DungeonTags").TaxonomyTags[]}
             */
            let new_cluster_tags = [];

            if (taxonomy_tags_index > 0) {
                new_cluster_tags = $cluster_tags.slice(0, taxonomy_tags_index);
            }

            new_cluster_tags.push(new_taxonomy_tags);

            if (taxonomy_tags_index < $cluster_tags.length - 1) {
                new_cluster_tags = new_cluster_tags.concat($cluster_tags.slice(taxonomy_tags_index + 1));
            }

            cluster_tags.set(new_cluster_tags);

            await refreshCurrentCategoryTaggings();
        }
        
        /**
         * Handles the tag-selected event from the ClusterPublicTags component.
         * @param {CustomEvent<{tag_id: number}>} event
         */
        const handleTagSelection = async (event) => {
            if ($current_category == null) {
                console.error("In CategoryTagger.handleTagSelection: No current category available while trying to tag it.");
                return;
            }
            
            const tag_id = event.detail.tag_id;

            let tagging_id = await tagCategory($current_category.uuid, tag_id);

            if (tagging_id == null) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.handleTagSelection");

                variable_environment.addVariable("tag_id", tag_id);
                variable_environment.addVariable("current_category.uuid", $current_category.uuid);

                const labeled_error = new LabeledError(variable_environment, "Failed to tag the entity. Duplicated tagging?", lf_errors.ERR_LOADING_ERROR);

                labeled_error.alert();
                return;
            }

            console.log("Tagging ID: ", tagging_id);

            await refreshCurrentCategoryTaggings();
        }

        /**
         * Refreshes the current category taggings and sets them on current_category_taggings property.
         */
        const refreshCurrentCategoryTaggings = async () => {
            if ($current_category == null) {
                console.error("In CategoryTagger.refreshCurrentCategoryTaggings: No current category available while trying to refresh its taggings.");
                return;
            }
            
            let new_taggings = await getEntityTaggings($current_category.uuid, $current_category.ClusterUUID);

            const category_uuid = $current_category.uuid;

            console.log(`'${category_uuid}' had taggings:`, new_taggings);

            current_category_taggings = new_taggings;
        }

        /**
         * Removes a given tag from the current category.
         * @param {number} tag_id
         */
        const removeCategoryTag = async (tag_id) => {
            if ($current_category == null) {
                console.error("In CategoryTagger.removeCategoryTag: No current category available while trying to remove a tag.");
                return;
            }
            
            if (tag_id < 0) {
                console.error(`Tag ids are always positive numbers. got ${tag_id}`);
                return;
            }

            const deleted = await untagEntity($current_category.uuid, tag_id);

            if (!deleted) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.removeCategoryTag");

                variable_environment.addVariable("tag_id", tag_id);
                variable_environment.addVariable("current_category.uuid", $current_category.uuid);

                const labeled_err = new LabeledError(variable_environment, "Could not remove attribute.", lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();
            }

            await refreshCurrentCategoryTaggings();
        }
    
        /**
         * Verifies the cluster_tags correctness with respect to the current_cluster. And if necessary, updates the cluster_tags.
         * @requires cluster_tags_checked
         */    
        const verifyLoadedClusterTags = async () => {
            if ($last_cluster_domain === $current_cluster.UUID) {
                cluster_tags_checked = true;
            }

            let loaded = await refreshClusterTagsNoCheck($current_cluster.UUID);

            if (!loaded) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.verifyLoadedClusterTags");

                variable_environment.addVariable("current_cluster.UUID", $current_cluster.UUID);
                variable_environment.addVariable("cluster_tags", $cluster_tags);
                variable_environment.addVariable("last_cluster_domain", $last_cluster_domain);

                const labeled_error = new LabeledError(variable_environment, `No attributes set for '${$current_cluster.Name}'. Try adding some.`, lf_errors.ERR_LOADING_ERROR);

                labeled_error.alert();
            }

            cluster_tags_checked = true;
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<dialog open id="dungeon-category-tagger-tool"
    bind:this={the_category_tagger_tool}
    class:section-activated={sub_component_sections[ct_focused_section].Active}
    class="libery-dungeon-window"
>
    <section id="tag-taxonomy-creator-section" 
        class="dctt-section"
        class:focused-section={ct_focused_section === 0}
    >
        <TagTaxonomyCreator
            component_hotkey_context={tag_taxonomy_creator_hotkeys_context}
            ui_entity_reference={ui_entity_reference}
            ui_taxonomy_reference={ui_taxonomy_reference}
            on:tag-taxonomy-created={handleTagTaxonomyCreated}
            on:drop-hotkeys-control={handleRecoverHotkeysControl}
        />
    </section>
    <article id="dctt-current-category-tags-wrapper"
        class="dungeon-scroll dctt-section"
        class:focused-section={ct_focused_section === 1}
    >
        <CategoryTaggings 
            component_hotkey_context={category_taggings_hotkeys_context}
            current_category_taggings={current_category_taggings}
            on:remove-category-tag={handleRemoveCategoryTag}
            on:drop-hotkeys-control={handleRecoverHotkeysControl}
        />
    </article>
    <article id="dctt-cluster-user-tags"
        class="dctt-section"
        class:focused-section={ct_focused_section === 2}
    >
        {#if cluster_tags_checked && taxonomy_tags_hotkeys_context != null}
            <ClusterPublicTags 
                component_hotkey_context={cluster_public_tags_hotkeys_context}
                has_hotkey_control={cluster_public_tags_hotkeys_context.Active}
                ui_entity_reference={ui_entity_reference}
                ui_taxonomy_reference={ui_taxonomy_reference}
                ui_tag_reference={ui_tag_reference}
                taxonomy_tags_hotkeys_context={taxonomy_tags_hotkeys_context}
                on:tag-selected={handleTagSelection}
                on:delete-taxonomy={handleTagTaxonomyDeleted}
                on:taxonomy-content-change={handleTaxonomyContentChanged}
                on:drop-hotkeys-control={handleRecoverHotkeysControl}
            />
        {/if}
    </article>
</dialog>

<style>

    dialog#dungeon-category-tagger-tool {
        position: static;
        display: flex;
        width: clamp(400px, 82dvw, 1800px);
        height: calc(calc(100dvh - var(--navbar-height)) * 0.9);
        container-type: size;
        flex-direction: column;
        row-gap: calc(var(--spacing-2) + var(--spacing-1));
        padding: var(--spacing-1);
        z-index: var(--z-index-t-1);
        outline: none;

        & > .dctt-section {
            padding: 0 var(--spacing-2);
            border-left-width: calc(var(--med-border-width) * 1.5);
            border-style: solid;
            border-left-color:  transparent;
        }
    
        & > .dctt-section:not(:last-child) {
            border-bottom: var(--border-thin-grey-8);
        }

        & > .dctt-section.focused-section {
            border-left-color: hsl(from var(--main-6) h s l / 0.7);
            transition: border-left-color 0.3s ease-out;
        }
    }

    #dungeon-category-tagger-tool.section-activated {
        & > .dctt-section.focused-section {
            border-left-width: calc(var(--med-border-width) * 2.5);
        }
    }

    article#dctt-cluster-user-tags {
        height: 35cqh;
        container-type: size;
    }

    article#dctt-current-category-tags-wrapper {
        height: 30cqh;
        overflow: auto;
    }
</style>