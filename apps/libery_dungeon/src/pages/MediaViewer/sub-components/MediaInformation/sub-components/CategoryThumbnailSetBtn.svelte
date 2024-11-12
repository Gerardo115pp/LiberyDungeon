<script>
    import { current_category } from "@stores/categories_tree";
    import { changeCategoryThumbnail } from "@models/Categories";
    import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * @type {import("@models/Medias").Media}
         */ 
        export let the_current_media;
    
    /*=====  End of Properties  ======*/
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Changes the current category thumbnail to the_current_media when the button is clicked.
         * @param {MouseEvent} event
         */
        const handleChangeCategoryThumbnailBTNClick = async (event) => {
            event.preventDefault();

            if ($current_category == null) {
                console.error("In CategoryThumbnailSetBtn.handleChangeCategoryThumbnailBTNClick: $current_category is null.");
                return;
            }
            
            let changed = await changeCategoryThumbnail($current_category.uuid, the_current_media.uuid);
            
            if (changed) {
                emitPlatformMessage(`The set the thumbnail of the category ${$current_category.name} to the media ${the_current_media.name}.`);
            } else {
                const variable_environment = new VariableEnvironmentContextError("In CategoryThumbnailSetBtn.handleChangeCategoryThumbnailBTNClick");

                variable_environment.addVariable("the_current_media", the_current_media.uuid);
                variable_environment.addVariable("$current_category", $current_category.uuid);

                const labeled_error = new LabeledError(variable_environment, "Failed to change the category thumbnail.", lf_errors.ERR_PROCESSING_ERROR);

                labeled_error.alert();
            }
        }
    
    /*=====  End of Methods  ======*/

</script>

<button 
    class="ctsb-category-thumbnail-change-btn dungeon-button-1"
    on:click={handleChangeCategoryThumbnailBTNClick}
>
    Use media as thumbnail
</button>

<style>
    button.ctsb-category-thumbnail-change-btn {
        line-height: 1;
        padding-block: 0.4em;
        padding-inline: 2em;
    }
</style>