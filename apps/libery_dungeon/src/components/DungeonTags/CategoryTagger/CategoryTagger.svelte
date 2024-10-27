<script>
    /*=============================================
    =            Imports            =
    =============================================*/
        import { deleteTagTaxonomy, getEntityTaggings, getTaxonomyTagsByUUID, tagEntity, untagEntity } from "@models/DungeonTags";
        import { cluster_tags, last_cluster_domain, refreshClusterTagsNoCheck } from "@stores/dungeons_tags";
        import TagTaxonomyCreator from "../TagTaxonomyComponents/TagTaxonomyCreator.svelte";
        import { createEventDispatcher, onDestroy, onMount } from "svelte";
        import { current_cluster } from "@stores/clusters";
        import { current_category } from "@stores/categories_tree";
        import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
        import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
        import TaxonomyTags from "../TagTaxonomyComponents/TaxonomyTags.svelte";
        import ClusterPublicTags from "./sub-components/ClusterPublicTags.svelte";
        import CategoryTaggings from "./sub-components/CategoryTaggings.svelte";
        import { json } from "@sveltejs/kit";
        import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
        import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
        import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import { toggleHotkeysSheet } from "@stores/layout";   
    import { CATEGORY_ENTITY_TYPE } from "@app/config/dungeon_tags_config";
    /*=====  End of Imports  ======*/
    
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

            const hotkeys_context_name = "category-tagger-tool";
        
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
         * @type {HTMLDialogElement}
         */
        let the_category_tagger_tool = null;

        
        /*----------  Component metadata  ----------*/
        
            /**
             * The count of dtt-sections in the category tagger.
             * @type {number}
             */
            let dtt_sections_count;
        
        /*----------  Hotkey state  ----------*/
        
            /**
             * The focused section indicator.
             * @type {nummber}
             */ 
            let ct_focused_section = 0;

            /**
             * Section active.
             * @type {boolean}
             */
            let ct_section_active = false;

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
                }
            });
        } else {
            await verifyLoadedClusterTags();
        }

        current_category_unsubscriber = current_category.subscribe(handleCurrentCategoryChange)

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
             * Defines the tools hotkeys.
             */ 
            const defineDesktopKeybinds = () => {
                if (global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    global_hotkeys_manager.dropContext(hotkeys_context_name);
                }

                const hotkeys_context = new HotkeysContext();

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

                hotkeys_context.register(["?"], toggleHotkeysSheet, {
                    description: `<${HOTKEYS_GENERAL_GROUP}>Opens the hotkeys cheat sheet.`
                });

                global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);

                global_hotkeys_manager.loadContext(hotkeys_context_name);
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
                if (ct_section_active) return;

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

                ct_section_active = true;
            }

            /**
             * Recovers the hotkeys control and deactivates the active section.
             */
            const handleRecoverHotkeysControl = () => {
                ct_section_active = false;
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
         * Defines the component's content metadata.
         */
        const defineComponentMetadata = () => {
            let dtt_sections = getSectionNodes();

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
         * @returns {NodeListOf<HTMLElement>}
         */
        const getSectionNodes = () => {
            return the_category_tagger_tool.querySelectorAll(".dctt-section");
        }

        /**
         * Handles the change of the current category.
         * @param {import("@models/Categories").CategoryLeaf} new_category
         */
        const handleCurrentCategoryChange = async new_category => {
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
         * @param {CustomEvent<{item_id: string}>} event
         */
        const handleTagSelection = async (event) => {
            const tag_id = event.detail.tag_id;

            let tagging_id = await tagEntity($current_category.uuid, tag_id, CATEGORY_ENTITY_TYPE);

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
            let new_taggings = await getEntityTaggings($current_category.uuid, $current_category.ClusterUUID);

            const category_uuid = $current_category.uuid;

            console.log(`'${category_uuid}' had taggings:`, new_taggings);

            current_category_taggings = new_taggings;
        }

        /**
         * Removes a given tag from the current category.
         * @param {number} tag_id
         */
        const removeCategoryTag = async tag_id => {
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
    class:section-activated={ct_section_active}
    class="libery-dungeon-window"
>
    <section id="tag-taxonomy-creator-section" 
        class="dctt-section"
        class:focused-section={ct_focused_section === 0}
    >
        <TagTaxonomyCreator
            has_hotkey_control={ct_focused_section === 0 && ct_section_active}
            on:tag-taxonomy-created={handleTagTaxonomyCreated}
            on:drop-hotkeys-control={handleRecoverHotkeysControl}
        />
    </section>
    <article id="dctt-current-category-tags-wrapper"
        class="dungeon-scroll dctt-section"
        class:focused-section={ct_focused_section === 1}
    >
        <CategoryTaggings 
            has_hotkey_control={ct_focused_section === 1 && ct_section_active}
            current_category_taggings={current_category_taggings}
            on:remove-category-tag={handleRemoveCategoryTag}
            on:drop-hotkeys-control={handleRecoverHotkeysControl}
        />
    </article>
    <article id="dctt-cluster-user-tags"
        class="dctt-section"
        class:focused-section={ct_focused_section === 2}
    >
        {#if cluster_tags_checked}
            <ClusterPublicTags 
                has_hotkey_control={ct_focused_section === 2 && ct_section_active}
                on:tag-selected={handleTagSelection}
                on:delete-taxonomy={handleTagTaxonomyDeleted}
                on:taxonomy-content-change={handleTaxonomyContentChanged}
                on:drop-hotkeys-control={handleRecoverHotkeysControl}
            />
        {/if}
    </article>
</dialog>

<style>
    #dungeon-category-tagger-tool {
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