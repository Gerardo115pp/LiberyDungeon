<script>
    import { current_cluster } from "@stores/clusters";
    import { cluster_tags } from "@stores/dungeons_tags";
    import { createTagTaxonomy } from "@models/DungeonTags";
    import { LabeledError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { createEventDispatcher, onDestroy, onMount } from "svelte";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { toggleHotkeysSheet } from "@stores/layout";
    import { browser } from "$app/environment";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    
    /*=============================================
    =            Properties            =
    =============================================*/
        // NOTE: Call TagTaxonomy -> Attributes in the UI. It's a more user-friendly term.
        
        /*=============================================
        =            Hotkeys            =
        =============================================*/
        
            /**
             * @type {import("@libs/LiberyHotkeys/libery_hotkeys").HotkeyContextManager}
             */ 
            const global_hotkeys_manager = getHotkeysManager();

            const hotkeys_context_name = "tag-taxonomy-creator";
        
        /*=====  End of Hotkeys  ======*/
        
        /**
         * The name of the would-be tag taxonomy.
         * @type {string}
         * @default ""
         */ 
        let new_attribute_name = "";

        /**
         * The input field where the user types the new tag taxonomy name.
         * @type {HTMLInputElement}
         */
        let the_tag_taxonomy_name_input;

        /**
         * The tag creator form.
         * @type {HTMLFormElement}
         */
        let the_tag_creator_form;

        /**
         * Whether the new name is valid and ready to be created.
         * @type {boolean}
         */
        let new_tag_taxonomy_name_is_valid = false;
        
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
                focusNewAttributeNameInput();
            }
        
        /*=====  End of Hotkeys state  ======*/

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        has_mounted = true;
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
             * Emits an event to close the category tagger tool and drops the hotkeys context.
             * @param {KeyboardEvent} event
             * @param {import("@libs/LiberyHotkeys/hotkeys").HotkeyData} hotkey
             */
            const handleCloseCategoryTaggerTool = (event, hotkey) => {
                resetHotkeyContext();
                emitDropHotkeyContext()
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
         * Creates a new tag taxonomy.
         */
        const createNewTagTaxonomyIfValid = async () => {
            if (!new_tag_taxonomy_name_is_valid) return;
            
            new_attribute_name = new_attribute_name.trim();

            let is_available_and_valid = checkNewTagTaxonomyName();

            if (!is_available_and_valid) return;

            let new_taxonomy = true;
            // let new_taxonomy = await createTagTaxonomy(new_attribute_name, $current_cluster.UUID, false);

            if (new_taxonomy === null) {
                let labeled_err = new LabeledError("In TagTaxonomyCreator.createTagTaxonomyIfValid", "The was an error creating the new tag taxonomy.", lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();
            }

            emitTagTaxonomyCreated();

            resetTagTaxonomyCreator();
        }

        /**
         * Checks whether the new tag taxonomy name is valid.
         * @returns {boolean}
         */         
        const checkNewTagTaxonomyName = () => {
            let is_valid = the_tag_taxonomy_name_input.checkValidity();

            if (!is_valid) return false;

            is_valid = !taxonomyNameExistsInCluster(new_attribute_name);

            if (!is_valid) {
                the_tag_taxonomy_name_input.setCustomValidity("The attribute name already exists.");
                the_tag_taxonomy_name_input.reportValidity();
            }

            return is_valid;
        }

        /**
         * Emits the tag taxonomy created event.
         */
        const emitTagTaxonomyCreated = () => {
            dispatch("tag-taxonomy-created");
        }

        /**
         * Focuses the new attribute name input field.
         */
        const focusNewAttributeNameInput = () => {
            the_tag_taxonomy_name_input.focus();
        }

        /**
         * Handles the keydown event on the new_attribute_name input field.
         * @param {KeyboardEvent} event
         */
        const handleKeyDown = (event) => {
            if (event.key === "Enter") {
                event.preventDefault();
                createNewTagTaxonomyIfValid();
            }

            if (event.key === "Escape") {
                the_tag_taxonomy_name_input.blur();
            }
        }

        /**
         * Handles the keyup event on the new_attribute_name input field.
         * @param {KeyboardEvent} event
         */
        const handleKeyUp = (event) => {
            event.preventDefault();

            if (the_tag_taxonomy_name_input.validationMessage !== "") {
                the_tag_taxonomy_name_input.setCustomValidity("");
            }

            new_tag_taxonomy_name_is_valid = the_tag_taxonomy_name_input.checkValidity();
        }

        /**
         * Resets the tag taxonomy creator.
         */
        const resetTagTaxonomyCreator = () => {
            the_tag_creator_form.reset();
            new_tag_taxonomy_name_is_valid = false;
        }

        /**
         * Returns whether the passed taxonomy name exists withing the the array of TaxonomyTags stored in cluster_tags.
         * @param {string} taxonomy_name
         * @returns {boolean}
         */ 
        const taxonomyNameExistsInCluster = (taxonomy_name) => {
            return $cluster_tags.some(taxonomy_tag => taxonomy_tag.Taxonomy.Name === taxonomy_name);
        }


    
    /*=====  End of Methods  ======*/

</script>

<form  class="tag-taxonomy-creator"
    bind:this={the_tag_creator_form}
    action="none"
>
    <header id="ttc-header">
        <p class="dungeon-instructions">
            Here you can create new category attributes. These will be available in all the categories within <strong>{$current_cluster?.Name ?? ""}</strong>
        </p>
    </header>
    <fieldset class="ttc-fields">
        <label class="dungeon-input">
            <span class="dungeon-label">
                Attribute
            </span>
            <input 
                bind:this={the_tag_taxonomy_name_input}
                bind:value={new_attribute_name}
                type="text"
                on:keydown={handleKeyDown}
                on:keyup={handleKeyUp}
                spellcheck="true"
                minlength="2"
                maxlength="64"
                pattern="{'[A-z_][A-z_\\s]{2,64}'}"
                required
            >
        </label>
        <button
            type="button"
            class="dungeon-button-1"
            disabled={!new_tag_taxonomy_name_is_valid}
            on:click={createNewTagTaxonomyIfValid}
        >
            Create
        </button>
    </fieldset>
</form>

<style>
    form.tag-taxonomy-creator {
        display: flex;
        height: 20cqh;
        flex-direction: column;
        gap: var(--spacing-2);
        justify-content: center;
        align-items: center;
        padding: 0 var(--spacing-3);
    }

    fieldset.ttc-fields {
        display: flex;
        width: 50cqw;
        gap: var(--spacing-2);
        padding: 0 var(--spacing-2);

        & > label.dungeon-input {
            flex-grow: 2;
        }
    }
</style>

