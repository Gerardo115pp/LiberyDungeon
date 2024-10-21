<script>
    import { current_cluster } from "@stores/clusters";
    
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
    
    /*=====  End of Properties  ======*/
    
</script>

<form action="none" class="tag-taxonomy-creator">
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
                spellcheck="true"
                minlength="1"
                maxlength="64"
                pattern="{'[a-z_]{1,64}'}"
                required
            >
        </label>
        <button
            type="button"
            class="dungeon-button-1"
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

