<script>
    import { ui_core_dungeon_references } from "@app/common/ui_references/core_ui_references";
    import { ui_pandasworld_tag_references } from "@app/common/ui_references/dungeon_tags_references";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import { cluster_tags, cluster_tags_checked, refreshClusterTagsNoCheck } from "@stores/dungeons_tags";
    import generateTaggedMediasHotkeyContext, { tagged_medias_child_contexts } from "./tagged_medias_hotkeys";
    import ClusterPublicTags from "../TagTaxonomyComponents/ClusterPublicTags.svelte";
    import generateTaxonomyTagsHotkeysContext, { taxonomy_tags_actions } from "../TagTaxonomyComponents/TaxonomyTags/taxonomy_tags_hotkeys";
    import { current_cluster } from "@stores/clusters";
    import { onMount } from "svelte";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import SelectedTags from "../Tags/SelectedTags.svelte";
    import { last_keyboard_focused_tag } from "../TagTaxonomyComponents/TaxonomyTags/taxonomy_tags_store";
    import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
    
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
             * The component hotkey context for the tagged medias component.
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
             */
            export let component_hotkey_context = generateTaggedMediasHotkeyContext();

            /* ------------------------- sub-components contexts ------------------------ */

                /**
                 * The component hotkey context for the cluster public tags component that is a child of this component.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let cluster_public_tags_hotkey_context = component_hotkey_context.ChildHotkeysContexts.get(tagged_medias_child_contexts.CLUSTER_PUBLIC_TAGS) ?? null;
                if (cluster_public_tags_hotkey_context == null) {
                    throw new Error("In TaggedMedias, invalid component_hotkey_context: ClusterPublicTags hotkey context was not defined as a child context of the TaggedMedias hotkey context.");
                }

                /**
                 * The component hotkeys context for the child taxonomy tags components.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
                 */
                let taxonomy_tags_hotkeys_context = null;

            /* ------------------------------ hotkeys state ----------------------------- */

                /**
                 * The index of the child hotkey context that is currently focused(!= to active.)
                 * @type {number}
                 */ 
                let focused_section_index = 0;

                /**
                 * Sub-component sections.
                 * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext[]}
                 */
                const sub_component_sections = [
                    cluster_public_tags_hotkey_context,
                ];
        
        /*=====  End of Hotkeys  ======*/
        
        /**
         * The list of tags that medias must be tagged with.
         * @type {import('@models/DungeonTags').DungeonTag[]}
         */
        let filtering_dungeon_tags = [];

        /**
         * A callback called when event the filtering tags change.
         * @type {import('./tagged_medias').MediaTagsChangedCallback}
         */
        export let onFilterTagsChange = (tags) => {};
        
        /*----------  Style  ----------*/
        
            /**
             * A number between 0 and 1 that is used for the alpha channel of the tool's background color.
             * @type {number}
             * @default 1
             */
            export let background_alpha = 1;
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        if ($current_cluster === null) {
            const current_cluster_unsubscriber = current_cluster.subscribe(async (value) => {
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
    
    /*=============================================
    =            Methods            =
    =============================================*/
        
        /*=============================================
        =            Hotkeys            =
        =============================================*/

            /**
             * Defines the tool's hotkeys.
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

                const hotkeys_context = preparePublicHotkeysActions(component_hotkey_context); // Not using the TaggedMedias HotkeyContext right now.

                cluster_public_tags_hotkey_context.Active = true;
            }

            /**
             * Handles the delete focused tag action from the TagTaxonomy component.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleTaxonomyTagsDeleteFocusedTag = (event, hotkey) => {
                const focused_tag = $last_keyboard_focused_tag;

                if (focused_tag === null) {
                    console.error("In TaggedMedias.handleTaxonomyTagsDeleteFocusedTag: Seems like no TaxonomyTags component has focus any tags.")
                    return;
                }

                removeFilteringTag(focused_tag.Id);
            }
            
            /**
             * Recovers the hotkeys control and deactivates the active section.
             */
            const handleRecoverHotkeysControl = () => {
                sub_component_sections[focused_section_index].Active = false;
            }

            /**
             * Prepares the public hotkey cations to generate the hotkeys context.
             * @param {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext} new_component_hotkey_context
             * @returns {import('@libs/LiberyHotkeys/hotkeys_context').default}
             */
            const preparePublicHotkeysActions = new_component_hotkey_context => {
                const hotkey_context = new_component_hotkey_context.generateHotkeysContext();

                return hotkey_context;
            }

            /* ------------------------- sub-components hotkeys ------------------------- */

                /**
                 * Defines the hotkeys context for sub-components.
                 */
                const defineSubComponentsHotkeysContext = () => {   
                    taxonomy_tags_hotkeys_context = defineTaxonomyTagsHotkeyContext();

                    component_hotkey_context.inheritExtraHotkeys();
                }

                /**
                 * defines the hotkeys context for the taxonomy tags components.
                 * @requires generateTaxonomyTagsHotkeysContext
                 * @returns {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
                 */
                const defineTaxonomyTagsHotkeyContext = () => {
                    const taxonomy_tags_context = generateTaxonomyTagsHotkeysContext();

                    /* ------------------------- public action overwrite ------------------------ */

                        const delete_focused_tag_action = taxonomy_tags_context.getHotkeyActionOrPanic(taxonomy_tags_actions.DELETE_FOCUSED_TAG);

                        if (delete_focused_tag_action.OverwriteBehavior !== ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE) {
                            throw new Error(`In TaggedMedias.defineTaxonomyTagsHotkeyContext: Overwrite behavior for DELETE_FOCUSED_TAG action of taxonomy_tags_context is not set to REPLACE.`);
                        }

                        delete_focused_tag_action.Callback = handleTaxonomyTagsDeleteFocusedTag;

                    /* -------------------------------------------------------------------------- */

                    return taxonomy_tags_context;
                }
        

        
        /*=====  End of Hotkeys  ======*/

        /**
         * Adds a tag from the cluster tags to the filtering tags.
         * @param {number} cluster_tag_id
         * @returns {void}
         */
        const addFilteringTag = cluster_tag_id => {
            /**
             * The new filtering dungeon tag
             * @type {import('@models/DungeonTags').DungeonTag | null}
             */
            let new_filtering_tag = null;

            for (let taxonomy_tags of $cluster_tags) {
                new_filtering_tag = taxonomy_tags.findTagByID(cluster_tag_id) ?? null;

                if (new_filtering_tag !== null) {
                    break
                }
            }

            if (new_filtering_tag === null) {
                let variable_environment = new VariableEnvironmentContextError(`In TaggedMedias.AddFilteringTag: cluster_tag_id<${cluster_tag_id} not found in cluster_tags.`)

                variable_environment.addVariable('cluster_tag_id', cluster_tag_id);

                const labeled_error = new LabeledError(
                    variable_environment,
                    `Could not find ${ui_pandasworld_tag_references.TAG.EntityName} for this ${ui_core_dungeon_references.CATEGORY_CLUSTER}`,
                    lf_errors.PROGRAMMING_ERROR__BROKEN_STATE
                )

                labeled_error.alert();
                return;
            }

            filtering_dungeon_tags = [...filtering_dungeon_tags, new_filtering_tag];

            onFilterTagsChange(filtering_dungeon_tags);
        }

        /**
         * Handles the tag-selected event from the ClusterPublicTags component.
         * @param {CustomEvent<{tag_id: number}>} event
         */
        const handleTagSelection = async (event) => {
            if (!(typeof event.detail.tag_id === "number")) return;

            addFilteringTag(event.detail.tag_id);
        }

        /**
         * Removes a filtering tag from the filtering tag list.
         * @param {number} tag_id
         * @returns {void}
         */
        const removeFilteringTag = tag_id => {
            filtering_dungeon_tags = filtering_dungeon_tags.filter(tag => tag.Id !== tag_id);

            onFilterTagsChange(filtering_dungeon_tags);
        }

        /**
         * Verifies the cluster_tags correctness with respect to the current_cluster. And if necessary, updates the cluster_tags.
         * @requires cluster_tags_checked
         */    
        const verifyLoadedClusterTags = async () => {
            if ($cluster_tags_checked) return;

            let loaded = await refreshClusterTagsNoCheck($current_cluster.UUID);

            if (!loaded) {
                const variable_environment = new VariableEnvironmentContextError("In CategoryTagger.verifyLoadedClusterTags");

                variable_environment.addVariable("current_cluster.UUID", $current_cluster.UUID);

                const labeled_error = new LabeledError(variable_environment, `No attributes set for '${$current_cluster.Name}'. Try adding some.`, lf_errors.ERR_LOADING_ERROR);

                labeled_error.alert();
            }
        }   
    /*=====  End of Methods  ======*/
    
</script>

<dialog open id="dungeon-tagged-medias-tool"
    class="libery-dungeon-window"
    class:section-activated={true}
    style:--background-alpha={background_alpha}
>
    <header id="dtmt-header">
        <h2 id="dtmt-header-headline">
            Filter {ui_core_dungeon_references.MEDIA.EntityNamePlural} by {ui_pandasworld_tag_references.TAG_TAXONOMY.EntityNamePlural} with multiple {ui_pandasworld_tag_references.TAG.EntityNamePlural}
        </h2>
    </header>
    <div id="dtmt-filering-medias-section"
     class="dtmt-section"
    >
        <SelectedTags 
            selected_tag_list={filtering_dungeon_tags}
        >
            <h3 class="fmt-fms-headline" slot="headline">
                Medias must include:
            </h3>
        </SelectedTags>
    </div>
    <div id=dtmt-cluster-public-tags-section
        class="dtmt-section"
    >
        {#if $cluster_tags_checked && taxonomy_tags_hotkeys_context != null}
            <ClusterPublicTags 
                has_hotkey_control={cluster_public_tags_hotkey_context.Active}
                ui_entity_reference={ui_core_dungeon_references.MEDIA}
                ui_taxonomy_reference={ui_pandasworld_tag_references.TAG_TAXONOMY}
                ui_tag_reference={ui_pandasworld_tag_references.TAG}
                taxonomy_tags_hotkeys_context={taxonomy_tags_hotkeys_context}
                component_hotkey_context={cluster_public_tags_hotkey_context}
                on:drop-hotkeys-control={handleRecoverHotkeysControl}
                on:tag-selected={handleTagSelection}
            />
        {/if}
    </div>
</dialog>

<style>
    dialog#dungeon-tagged-medias-tool {
        display: flex;
        width: 100%;
        height: 100%;
        container-type: size;
        background: hsl(from var(--body-bg-color) h s l / var(--background-alpha));
        flex-direction: column;
        padding: var(--spacing-1);
        row-gap: calc(var(--spacing-2) + var(--spacing-1));
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

    header#dtmt-header {
        display: flex;
        flex-direction: column;
        align-items: center;
        border-bottom: 1px solid var(--grey-8);

        & > h2:first-of-type {
            font-family: var(--font-read);
            color: var(--grey-3);
            font-weight: 520;
            line-height: 2;
        }
    }

    #dtmt-filering-medias-section {
        container: size;
        height: calc(95cqh - calc(2em + 70cqh));
        padding-inline-start: var(--spacing-1);

        & .fmt-fms-headline {
            font-family: var(--font-read);
            color: var(--main-6);
            font-weight: 500;
        }
    }

    #dtmt-cluster-public-tags-section {
        container-type: size;
        height: 70cqh;
    }

</style>