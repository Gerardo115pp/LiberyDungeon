<script>
    import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
    import { LabeledError } from '@libs/LiberyFeedback/lf_models';
    import { renameMedia } from '@models/Medias';
    import { onMount } from 'svelte';

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * the media that will be renamed.
         * @type {import('@models/Medias').Media}
         */ 
        export let the_media;

        /**
         * A callback to be triggered when the renaming process finalizes.
         * @type {(success: boolean) => void}
        */
        export let onRenameDone;

        /**
         * The renamer input.
         * @type {HTMLInputElement}
         */
        let the_renamer_input;
        
        /*----------  Behavior  ----------*/
        
            /**
             * whether the input should be auto-focused on keyup
             * @type {boolean}
             * @default false
             */ 
            export let should_autofocus = false;
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        if (should_autofocus && the_renamer_input != null) {
            focusRenamerInput();
        }
    })
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Applies the new name to the media.
         * @param {string} new_name
         */
        const applyMediaRename = async (new_name) => {
            if (new_name === "" || new_name === the_media.name) {
                onRenameDone(false);
                return;
            }

            let rename_successfully = await the_media.rename(new_name, true);
            if (!rename_successfully) {
                const labeled_err = new LabeledError(
                    "In @pages/MediaExplorer/sub-component/MediaExplorerGallery/MediaRenamingInput.applyMediaRename",
                    `Failed to rename media<${the_media.uuid}> to ${new_name}. Repeated name?`,
                    lf_errors.ERR_HUMAN_ERROR
                );

                labeled_err.alert();
            }

            onRenameDone(rename_successfully);
        }

        /**
         * Focused the renamer input node in the dom.
         * @returns {void}
         */
        const focusRenamerInput = () => {
            if (the_renamer_input == null) return;

            the_renamer_input.focus();
        }
    
        /**
         * Handle the keyup Keyboard event of the rename input
         * @param {KeyboardEvent} e
         */
        const handleRenameInput = e => {
            if (e.target instanceof HTMLInputElement && e.key === "Enter") {
                applyMediaRename(e.target.value);
            }

            if (e.key === "Escape") {
                onRenameDone(false);
            }
        }

        /**
         * Handles the blur event of the renamer input.
         * @returns {void}
         */
        const handleRenamerInputBlur = () => {
            onRenameDone(false);
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<input class="mri-media-renaming-input"
    bind:this={the_renamer_input}
    type="text"
    on:keyup={handleRenameInput}
    on:blur={handleRenamerInputBlur}
    value="{the_media.MediaName}"
>

<style>
    input.mri-media-renaming-input {
        width: 100%;
        background: transparent;
        border: none;
        outline: none;
    }
</style>