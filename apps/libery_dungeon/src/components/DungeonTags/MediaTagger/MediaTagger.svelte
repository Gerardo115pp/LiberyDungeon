<script>
    /*=============================================
    =            Imports            =
    =============================================*/
        import { deleteTagTaxonomy, getEntityTaggings, getTaxonomyTagsByUUID, multiTagMedia, tagCategory, tagCategoryContent, tagMedia, untagCategoryContent, untagEntity } from "@models/DungeonTags";
        import { cluster_tags, last_cluster_domain, refreshClusterTagsNoCheck } from "@stores/dungeons_tags";
        import { createEventDispatcher, onDestroy, onMount } from "svelte";
        import { current_cluster } from "@stores/clusters";
        import { LabeledError, UIReference, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
        import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
        import ClusterPublicTags from "../TagTaxonomyComponents/ClusterPublicTags.svelte";
        import TagTaxonomyCreator from "../TagTaxonomyComponents/TagTaxonomyCreator.svelte";
        import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
        import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
        import generateTaxonomyTagsHotkeysContext, { taxonomy_tags_actions } from "../TagTaxonomyComponents/TaxonomyTags/taxonomy_tags_hotkeys";
        import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import { wrapShowHotkeysTable } from "@app/common/keybinds/CommonActionWrappers";
        import { common_action_groups } from "@app/common/keybinds/CommonActionsName";
        import { ui_core_dungeon_references } from "@app/common/ui_references/core_ui_references";
        import { ui_pandasworld_tag_references } from "@app/common/ui_references/dungeon_tags_references";
        import MediaTaggings from "./sub-components/MediaTaggings.svelte";
        import generateMediaTaggerHotkeyContext, { media_tagger_actions, media_tagger_child_contexts } from "./media_tagger_hotkeys";
        import { cluster_public_tags_actions } from "../TagTaxonomyComponents/cluster_public_tags_hotkeys";
        import { linearCycleNavigationWrap } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";
        import dungeon_tags_clipboard from "./stores/dungeon_tags_clipboard";
    import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
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
             * The media tagger component hotkey context 
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
             */
            export let component_hotkey_context = generateMediaTaggerHotkeyContext();

            /* -------------------------- sub-component context -------------------------- */

                /**
                 * The component hotkey context for the cluster public tags component.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let cluster_public_tags_hotkeys_context = component_hotkey_context.ChildHotkeysContexts.get(media_tagger_child_contexts.CLUSTER_PUBLIC_TAGS) ?? null;
                if (cluster_public_tags_hotkeys_context == null) {
                    throw new Error("In MediaTagger, invalid component_hotkey_context: ClusterPublicTags hotkeys context was not defined as a child context of the MediaTagger context.");
                }

                /**
                 * The component hotkey context for the tag taxonomy creator.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let tag_taxonomy_creator_hotkeys_context = component_hotkey_context.ChildHotkeysContexts.get(media_tagger_child_contexts.TAG_TAXONOMY_CREATOR) ?? null;
                if (tag_taxonomy_creator_hotkeys_context == null) {
                    throw new Error("In MediaTagger, invalid component_hotkey_context: TagTaxonomyCreator hotkeys context was not defined as a child context of the MediaTagger context.");
                }

                /**
                 * The component hotkey context for the media taggings component.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let media_taggings_hotkeys_context = component_hotkey_context.ChildHotkeysContexts.get(media_tagger_child_contexts.MEDIA_TAGGINGS) ?? null;
                if (media_taggings_hotkeys_context == null) {
                    throw new Error("In MediaTagger, invalid component_hotkey_context: MediaTaggings hotkeys context was not defined as a child context of the MediaTagger context.");
                }

                /**
                 * The component hotkeys context for the child taxonomy tags components.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let taxonomy_tags_hotkeys_context = null;
        
        /*=====  End of Hotkeys  ======*/

        /**
         * The media to manage taggings for.
         * @type {import('@models/Medias').Media}
         */
        export let the_active_media;

    
        /**
         * Whether the cluster_tags correctness with respect to the current_cluster has been checked.
         * @type {boolean}
         */ 
        let cluster_tags_checked = false;
        $: if (the_active_media !== null && cluster_tags_checked) {
            handleActiveMediaChange(the_active_media);
        }

        /**
         * Current media taggings.
         * @type {import("@models/DungeonTags").DungeonTagging[]}
         */
        let current_media_taggings = [];

        /**
         * The category tagger tool node.
         * @type {HTMLDialogElement | null}
         */
        let the_media_tagger_tool = null;
        
        
        /*----------  Entities getters  ----------*/
        
            /**
             * A method that returns a list of taggable medias.
             * @type {() => import('@models/Medias').Media[]}
             */ 
            export let getTaggbleMedias;

        /*----------  Hotkey state  ----------*/
        
            /**
             * The focused section indicator.
             * @type {number}
             */ 
            let mt_focused_section = 0;

            /**
             * Sub-component sections.
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext[]}
             */
            const sub_component_sections = [
                tag_taxonomy_creator_hotkeys_context,
                media_taggings_hotkeys_context,
                cluster_public_tags_hotkeys_context
            ];       
       
        /*----------  Style  ----------*/
        
            /**
             * A number between 0 and 1 that is used for the alpha channel of the tool's background color.
             * @type {number}
             * @default 1
             */
            export let background_alpha = 1;
                


        /*----------  UI references  ----------*/
        
            /**
             * A UiReference object, to create ui messages about the taggable entity. used to pass down to general purpose components
             * @type {UIReference}
             */
            const ui_entity_reference = ui_core_dungeon_references.MEDIA;

            /**
             * A UiReference object, to create ui messages about the tag taxonomy. used to pass down to general purpose components.
             * @type {UIReference}
             */
            const ui_taxonomy_reference = ui_pandasworld_tag_references.TAG_TAXONOMY;

            /**
             * A UiReference object, to create ui messages about the dungeon-tag. used to pass down to general purpose components.
             * @type {UIReference}
             */
            const ui_tag_reference = ui_pandasworld_tag_references.TAG;
        
 
        /* ----------------------------------- */

        let current_cluster_unsubscriber = () => {};

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        if ($current_cluster === null) {
            current_cluster_unsubscriber = current_cluster.subscribe(async (value) => {
                if (value !== null) {
                    await verifyLoadedClusterTags();
                    current_cluster_unsubscriber();
                }
            });
        } else {
            await verifyLoadedClusterTags();
        }

        defineSubComponentsHotkeysContext();

        defineDesktopKeybinds();
    });

    onDestroy(() => {
        resetHotkeyContext();
    })
    
    /*=============================================
    =            Methods            =
    =============================================*/
        
        /*=============================================
        =            Keybinds            =
        =============================================*/
        
            /**
             * Defines the tool's hotkeys.
             */ 
            export function defineDesktopKeybinds() {
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

                const hotkeys_context = preparePublicHotkeysActions(component_hotkey_context);

                hotkeys_context.register(["q"], handleCloseMediasTaggerTool, {
                    description: `<${HOTKEYS_GENERAL_GROUP}>Closes the ${ui_entity_reference.EntityNamePlural} tagger tool.`,
                    await_execution: false
                });

                hotkeys_context.register(["e"], handleSectionSelection, {
                    description: `<navigation>Selects the current section.`,
                    await_execution: false
                });

                hotkeys_context.register(["y i"], handleShowRegistryInformation, {
                    description: `<registries>Shows the current registry information.`,
                });

                hotkeys_context.register(["y n"], handleClearAllRegistries, {
                    description: `<registries>Clears all the registries.`,
                });

                
               

                component_hotkey_context.applyExtraHotkeys();

                wrapShowHotkeysTable(hotkeys_context);

                global_hotkeys_manager.declareContext(component_hotkey_context.HotkeysContextName, hotkeys_context);

                global_hotkeys_manager.loadContext(component_hotkey_context.HotkeysContextName);
                component_hotkey_context.Active = true;
            }

            /**
             * Prepares the public hotkey cations to generate the hotkeys context.
             * @param {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext} new_component_hotkey_context
             * @returns {import('@libs/LiberyHotkeys/hotkeys_context').default}
             */
            const preparePublicHotkeysActions = (new_component_hotkey_context) => {

                /* --------------------------- Up down navigation --------------------------- */

                    const up_down_navigation = new_component_hotkey_context.getHotkeyActionOrPanic(media_tagger_actions.WS_NAVIGATION);

                    up_down_navigation.Options = {
                        description: `<navigation>Moves up and down through the tool's sections.`,
                        await_execution: false
                    }

                    up_down_navigation.Callback = handleSectionFocusNavigation

                /* --------------------------- Dungeon tags coping -------------------------- */
                    
                    const dungeon_tags_copy = new_component_hotkey_context.getHotkeyActionOrPanic(media_tagger_actions.COPY_CURRENT_MEDIA_TAGS);

                    dungeon_tags_copy.Options = {
                        description: `<registries>Copies the ${ui_taxonomy_reference.EntityNamePlural} ${ui_tag_reference.EntityNamePlural} of the current ${ui_core_dungeon_references.MEDIA.EntityName}`,
                    }

                    dungeon_tags_copy.Callback = handleCopyDungeonTags;

                /* -------------------------- Dungeon tags pasting -------------------------- */

                    const dungeon_tags_paste_action = new_component_hotkey_context.getHotkeyActionOrPanic(media_tagger_actions.PASTE_DUNGEON_TAGS);

                    dungeon_tags_paste_action.Options = {
                        description: `<registries>Pastes the ${ui_taxonomy_reference.EntityNamePlural} ${ui_tag_reference.EntityNamePlural} in the current registry to the current ${ui_core_dungeon_references.MEDIA.EntityName}`,
                    }

                    dungeon_tags_paste_action.Callback = handlePasteDungeonTags;

                /* ----------------------------- Registry change ---------------------------- */

                    const change_registry_action = new_component_hotkey_context.getHotkeyActionOrPanic(media_tagger_actions.CHANGE_COPY_REGISTRY);

                    change_registry_action.Options = {
                        description: `<registries>Changes registry where the ${ui_taxonomy_reference.EntityNamePlural} ${ui_tag_reference.EntityNamePlural} are been copied from and to.`,
                    }

                    change_registry_action.Callback = handleChangeDungeonTagsRegistry;

                /* -------------------------------------------------------------------------- */

                const hotkey_context = new_component_hotkey_context.generateHotkeysContext();

                return hotkey_context;
            }

            /**
             * Handles coping the current media dungeon tags.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleCopyDungeonTags = async (event, hotkey) => {
                const media_tags = current_media_taggings.map(tagging => tagging.Tag);

                dungeon_tags_clipboard.writeOnCurrentRegister(media_tags);

                const content = dungeon_tags_clipboard.readRegister();

                if (content != null) {
                    try {
                        await content.copy();
                    } catch (error) {
                        console.error("Failed to copy the tags to the clipboard.", error);
                    } 
                }
            }

            /**
             * Applies all the tags in the dungeon tag registry to the current active media. First filters out the tags that the current media doesn't have.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handlePasteDungeonTags = async (event, hotkey) => {
                let registry_content = dungeon_tags_clipboard.readRegister();

                if (registry_content == null) {
                    registry_content = await dungeon_tags_clipboard.readClipboard();

                    if (registry_content === null) {
                        emitPlatformMessage(`No ${ui_taxonomy_reference.EntityName} ${ui_tag_reference.EntityNamePlural} found in registry or the clipboard.`);
                        return;
                    }
                };

                const non_applied_tags = filterDuplicatedDungeonTags(registry_content.Content.DungeonTags);

                let must_refresh_tags = await multiTagMedia(the_active_media, convertDungeonTagsToIDList(non_applied_tags));

                if (must_refresh_tags) {
                    await refreshActiveMediaTaggings(the_active_media);
                }
            }

            /**
             * Handles the change of the current registry
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleChangeDungeonTagsRegistry = (event, hotkey) => {
                const new_registry_name = event.key;

                dungeon_tags_clipboard.changeCurrentRegister(new_registry_name);

                const new_registry_content = dungeon_tags_clipboard.readRegister();

                if (!new_registry_content) {
                    emitPlatformMessage(`Changed to clean registry: '${new_registry_name}'`);
                    return;
                }

                printCurrentDTClipboardRegistry();
            }

            /**
             * Shows the registry information.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleShowRegistryInformation = (event, hotkey) => {
                printCurrentDTClipboardRegistry();
            }

            /**
             * Clears all the registries.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleClearAllRegistries = (event, hotkey) => {
                dungeon_tags_clipboard.resetMediaTaggerDungeonTagsClipboard();
                emitPlatformMessage("All registries cleared.");
            
            }

            /**
             * Emits an event to close the medias tagger tool and drops the hotkeys context.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleCloseMediasTaggerTool = (event, hotkey) => {
                resetHotkeyContext();
                emitCloseMediasTagger();
            }

            /**
             * Handles the section focus navigation.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleSectionFocusNavigation = (event, hotkey) => {
                if (sub_component_sections[mt_focused_section].Active) return;

                let new_focused_section = mt_focused_section;

                const navigation_step = event.key === "w" ? -1 : 1;

                new_focused_section = linearCycleNavigationWrap(mt_focused_section, sub_component_sections.length - 1, navigation_step).value;

                mt_focused_section = new_focused_section;
            }

            /**
             * Handles the selection of tools sections.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleSectionSelection = (event, hotkey) => {
                event.preventDefault();

                sub_component_sections[mt_focused_section].Active = true;
            }

            /**
             * Recovers the hotkeys control and deactivates the active section.
             */
            const handleRecoverHotkeysControl = () => {
                sub_component_sections[mt_focused_section].Active = false;
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
                    
                    defineClusterPublicTagsHotkeysContext();

                    component_hotkey_context.inheritExtraHotkeys();
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
                        alt_select_focused_tag_action.overwriteDescription(`${common_action_groups.CONTENT}Applies the focused ${ui_tag_reference.EntityName} to the current ${ui_entity_reference.EntityName}.`);

                        alt_select_focused_tag_action.Callback = handleApplyFocusedTagToMedia;
                    }

                    const alt_delete_focused_tag_action = taxonomy_tags_context.getHotkeyAction(taxonomy_tags_actions.ALT_DELETE_FOCUSED_TAG);

                    if (alt_delete_focused_tag_action != null && alt_delete_focused_tag_action.OverwriteBehavior === ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE) {
                        alt_delete_focused_tag_action.overwriteDescription(`${common_action_groups.CONTENT}Removes the focused ${ui_tag_reference.EntityName} from the current ${ui_entity_reference.EntityName}.`);

                        alt_delete_focused_tag_action.Callback = handleRemoveFocusedTagFromContent;
                    }

                    return taxonomy_tags_context;
                }

                /**
                 * Defines the hotkeys context for the cluster public tags component.
                 * @returns {void}
                 */
                const defineClusterPublicTagsHotkeysContext = () => {
                    const cluster_public_tags_context = component_hotkey_context.ChildHotkeysContexts.get(media_tagger_child_contexts.CLUSTER_PUBLIC_TAGS);

                    if (cluster_public_tags_context == null) {
                        throw Error("In MediaTagger.defineClusterPublicTagsHotkeysContext: ClusterPublicTags context was not defined as a child context of the MediaTagger context.");
                    }

                    /* -------------------------- left right navigation ------------------------- */

                        const component_left_right_navigation = component_hotkey_context.getHotkeyActionOrPanic(media_tagger_actions.AD_NAVIGATION);

                       if (!component_left_right_navigation.HasNullishCallback() && !component_left_right_navigation.HasNullishDescription()) {
                           const child_left_right_navigation = cluster_public_tags_context.getHotkeyActionOrPanic(cluster_public_tags_actions.AD_NAVIGATION);

                           child_left_right_navigation.overwriteDescription(component_left_right_navigation.Options.description);

                           child_left_right_navigation.Callback = component_left_right_navigation.Callback;
                       }

                    /* -------------------------------------------------------------------------- */
                }

                /**
                 * Applies the focused tag on the active media.
                 * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback} 
                 */
                const handleApplyFocusedTagToMedia = async (event, hotkey) => {
                    if (getEntityTaggings == null) {
                        console.error("In MediaTagger.handleApplyFocusedTagToMedia: getEntityTaggings was null. Cannot determine the media list that should get the dungeon tag applied to with out this")
                    }
                }

                /**
                 * Removes the focused tag from the category content medias.
                 * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback} 
                 */
                const handleRemoveFocusedTagFromContent = async (event, hotkey) => {
                    // REFACTOR: If alt select is implemented, implement it's reverse/undo action here.
                    console.warn("Removing focused tag from media is not implemented.");
                }
        
        /*=====  End of Keybinds  ======*/
        
        /**
         * Converts a list of dungeon tags to a list of tag ids.
         * @param {import('@models/DungeonTags').DungeonTag[]} dungeon_tags
         * @returns {number[]}
         */
        const convertDungeonTagsToIDList = dungeon_tags => {
            return dungeon_tags.map(tag => tag.Id);
        }

        /**
         * Emits an event to the parent to close the medias tagger tool.
         */
        const emitCloseMediasTagger = () => {
            dispatch("close-medias-tagger");
        }

        /**
         * Filters out a given DungeonTag list by removing tags that are present in current_media_taggings.
         * @param {import('@models/DungeonTags').DungeonTag[]} new_tags
         * @returns {import('@models/DungeonTags').DungeonTag[]}
         */
        const filterDuplicatedDungeonTags = new_tags => {
            const current_media_tags = new Set(current_media_taggings.map(tagging => tagging.Tag.Id));

            /**
             * @type {import('@models/DungeonTags').DungeonTag[]}
             */
            const non_applied_tags = [];

            for (let new_tag of new_tags) {
                if (!current_media_tags.has(new_tag.Id)) {
                    non_applied_tags.push(new_tag)
                }
            }

            return non_applied_tags;
        }

        /**
         * Handles the change of the current category.
         * TODO: This behavior is not needed for the medias tagger tool. Find all it's branch effects and replace them for an 'active media' change handling instead.
         * @param {import("@models/Medias").Media | null} new_active_media
         */
        async function handleActiveMediaChange(new_active_media) {
            if (new_active_media === null) return;

            console.log("Refreshing taggings of:", new_active_media);

            if (new_active_media === null) return;

            await refreshActiveMediaTaggings(new_active_media);
        }

        /**
         * Handles the remove-media-tag event emitted by the MediaTaggings component.
         * @param {CustomEvent<{removed_tag: number}>} event
         */
        const handleRemoveMediaTag = event => {
            let tag_id = event?.detail?.removed_tag;

            if (tag_id == null) return;

            removeMediaTag(tag_id);
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
                const variable_environment = new VariableEnvironmentContextError("In MediasTagger.handleTagTaxonomyDeleted");

                variable_environment.addVariable("taxonomy_uuid", taxonomy_uuid);

                const labeled_error = new LabeledError(variable_environment, "Failed to delete the tag taxonomy.", lf_errors.ERR_PROCESSING_ERROR);

                labeled_error.alert();
                return;
            }

            await refreshClusterTagsNoCheck($current_cluster.UUID);
        }

        /**
         * Refreshes the content of a TagTaxonomy.
         * REFACTOR: This method(an the other similar tag taxonomy changes methods) are exactly the same as the ones in the CategoryTagger. If they need no changes at all, maybe they should be extracted to a common file.
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
                const variable_environment = new VariableEnvironmentContextError("In MediaTagger.handleTaxonomyContentChanged");

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

            await refreshActiveMediaTaggings(the_active_media);
        }
        
        /**
         * Handles the tag-selected event from the ClusterPublicTags component.
         * @param {CustomEvent<{tag_id: number}>} event
         */
        const handleTagSelection = async (event) => {
            const tag_id = event.detail.tag_id;

            let tagging_id = await tagMedia(the_active_media.uuid, tag_id);

            if (tagging_id == null) {
                const variable_environment = new VariableEnvironmentContextError("In MediaTagger.handleTagSelection");

                variable_environment.addVariable("tag_id", tag_id);
                variable_environment.addVariable("the_active_media.uuid", the_active_media.uuid);

                const labeled_error = new LabeledError(variable_environment, "Failed to tag the media.", lf_errors.ERR_PROCESSING_ERROR);

                labeled_error.alert();
                return;
            }

            await refreshActiveMediaTaggings(the_active_media);
        }

        /**
         * Prints the current dungeon_tags_clipboard registry as a platform message.
         * @return {void}
         */
        const printCurrentDTClipboardRegistry = () => {
            const registry_content = dungeon_tags_clipboard.readRegister();

            if (registry_content == null) {
                emitPlatformMessage("Registry clean.");
                return;
            }

            let feedback_message = `Registry '${dungeon_tags_clipboard.getCurrentRegister()}' with ${registry_content.Content.DungeonTags.length} tags: `;

            let h = 0;

            for (let tag of registry_content.Content.DungeonTags) {
                feedback_message += tag.Name

                if (h === registry_content.Content.DungeonTags.length - 1) {
                    feedback_message += ".";
                } else {
                    feedback_message += ", ";
                }

                h++;
            }

            emitPlatformMessage(feedback_message);
        }

        /**
         * Refreshes the current media taggings and sets them on active_media_taggings property.
         * @param {import('@models/Medias').Media} new_active_media
         */
        async function refreshActiveMediaTaggings(new_active_media) {
            if (new_active_media === null ) return;

            let new_taggings = await getEntityTaggings(new_active_media.uuid, $current_cluster.UUID);

            console.log(`Media ${new_active_media.uuid} had taggings:`, new_taggings);

            current_media_taggings = new_taggings;
        }

        /**
         * Removes a given tag from the active media.
         * @param {number} tag_id
         */
        const removeMediaTag = async (tag_id) => {
            if (tag_id < 0) {
                console.error(`Tag ids are always positive numbers. got ${tag_id}`);
                return;
            }

            const deleted = await untagEntity(the_active_media.uuid, tag_id);

            if (!deleted) {
                const variable_environment = new VariableEnvironmentContextError("In MediaTagger.removeMediaTag");

                variable_environment.addVariable("tag_id", tag_id);
                variable_environment.addVariable("the_active_media.uuid", the_active_media.uuid);

                const labeled_err = new LabeledError(variable_environment, "Could not remove attribute.", lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();
                return;
            }

            await refreshActiveMediaTaggings(the_active_media);
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

<dialog open id="dungeon-medias-tagger-tool"
    bind:this={the_media_tagger_tool}
    class:section-activated={sub_component_sections[mt_focused_section].Active}
    class="libery-dungeon-window"
    style:--background-alpha={background_alpha}
>
    <section id="dmtt-tag-taxonomy-creator-section" 
        class="dmtt-section"
        class:focused-section={mt_focused_section === 0}
    >
        <TagTaxonomyCreator
            component_hotkey_context={tag_taxonomy_creator_hotkeys_context}
            ui_entity_reference={ui_entity_reference}
            ui_taxonomy_reference={ui_taxonomy_reference}
            on:tag-taxonomy-created={handleTagTaxonomyCreated}
            on:drop-hotkeys-control={handleRecoverHotkeysControl}
        />
    </section>
    <article id="dmtt-current-media-tags-wrapper"
        class="dungeon-scroll dmtt-section"
        class:focused-section={mt_focused_section === 1}
    >
        <MediaTaggings 
            component_hotkey_context={media_taggings_hotkeys_context}
            the_active_media={the_active_media}
            current_media_taggings={current_media_taggings}
            has_hotkey_control={media_taggings_hotkeys_context.Active}
            on:drop-hotkeys-control={handleRecoverHotkeysControl}
            on:remove-category-tag={handleRemoveMediaTag}
        />
    </article>
    <article id="dmtt-cluster-user-tags"
        class="dmtt-section"
        class:focused-section={mt_focused_section === 2}
    >
        {#if cluster_tags_checked && taxonomy_tags_hotkeys_context != null && cluster_public_tags_hotkeys_context != null}
            <ClusterPublicTags 
                has_hotkey_control={cluster_public_tags_hotkeys_context.Active}
                ui_entity_reference={ui_entity_reference}
                ui_taxonomy_reference={ui_taxonomy_reference}
                ui_tag_reference={ui_tag_reference}
                taxonomy_tags_hotkeys_context={taxonomy_tags_hotkeys_context}
                component_hotkey_context={cluster_public_tags_hotkeys_context}
                on:tag-selected={handleTagSelection}
                on:delete-taxonomy={handleTagTaxonomyDeleted}
                on:taxonomy-content-change={handleTaxonomyContentChanged}
                on:drop-hotkeys-control={handleRecoverHotkeysControl}
            />
        {/if}
    </article>
</dialog>

<style>

    dialog#dungeon-medias-tagger-tool {
        position: static;
        box-sizing: border-box;
        display: flex;
        width: 100%;
        height: 100%;
        container-type: size;
        flex-direction: column;
        background: hsl(from var(--body-bg-color) h s l / var(--background-alpha));
        row-gap: calc(var(--spacing-2) + var(--spacing-1));
        padding: var(--spacing-1);
        z-index: var(--z-index-t-1);
        outline: none;

        & > .dmtt-section {
            padding: 0 var(--spacing-2);
            border-left-width: calc(var(--med-border-width) * 1.5);
            border-style: solid;
            border-left-color:  transparent;
        }
    
        & > .dmtt-section:not(:last-child) {
            border-bottom: var(--border-thin-grey-8);
        }

        & > .dmtt-section.focused-section {
            border-left-color: hsl(from var(--main-6) h s l / 0.7);
            transition: border-left-color 0.3s ease-out;
        }
    }

    #dungeon-medias-tagger-tool.section-activated {
        & > .dmtt-section.focused-section {
            border-left-width: calc(var(--med-border-width) * 2.5);
        }
    }

    article#dmtt-cluster-user-tags {
        height: 35cqh;
        container-type: size;
    }

    article#dmtt-current-media-tags-wrapper {
        height: 30cqh;
        overflow: auto;
    }
</style>