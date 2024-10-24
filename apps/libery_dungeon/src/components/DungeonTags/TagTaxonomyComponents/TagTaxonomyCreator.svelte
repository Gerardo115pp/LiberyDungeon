<script>
    import { current_cluster } from "@stores/clusters";
    import { cluster_tags } from "@stores/dungeons_tags";
    import { createTagTaxonomy } from "@models/DungeonTags";
    import { LabeledError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { createEventDispatcher } from "svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/
        // NOTE: Call TagTaxonomy -> Attributes in the UI. It's a more user-friendly term.
    
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

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    /*=============================================
    =            Methods            =
    =============================================*/


        /**
         * Creates a new tag taxonomy.
         */
        const createNewTagTaxonomyIfValid = async () => {
            if (!new_tag_taxonomy_name_is_valid) return;
            
            new_attribute_name = new_attribute_name.trim();

            let is_available_and_valid = CheckNewTagTaxonomyName();

            if (!is_available_and_valid) return;

            let new_taxonomy = await createTagTaxonomy(new_attribute_name, $current_cluster.UUID, false);

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
        const CheckNewTagTaxonomyName = () => {
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
         * Handles the keydown event on the new_attribute_name input field.
         * @param {KeyboardEvent} event
         */
        const handleKeyDown = (event) => {
            if (event.key === "Enter") {
                event.preventDefault();
                createNewTagTaxonomyIfValid();
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

