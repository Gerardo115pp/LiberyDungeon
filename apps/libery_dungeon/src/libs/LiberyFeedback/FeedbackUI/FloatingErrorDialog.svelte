<script>
    import { onDestroy, onMount } from "svelte";
    import { subscribeToPlatformErrors } from "../lf_utils";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * Whether the dialog should be floating at the bottom or top of the screen.
         * @type {boolean}
         * @default false
         */ 
        export let float_at_top = false;

        /**
         * The amount of time in milliseconds the dialog will be visible for after an error occurs.
         * @type {number}
         * @default 3000
         */
        export let display_error_duration = 3000;

        /**
         * The error that will be displayed in the dialog.
         * @type {import('../lf_models').LabeledError}
         */
        let labeled_error;

        /**
         * Whether the dialog should be visible.
         * @type {boolean}
        */
        let dialog_visible = false;

        
        /*----------  Unsuscribers  ----------*/
        
            let platform_error_unsuscriber = () => {}; 
    
    /*=====  End of Properties  ======*/

    onMount(() => {
            platform_error_unsuscriber = subscribeToPlatformErrors(handlePlatformError)
    });

    onDestroy(() => {
        platform_error_unsuscriber();
    });

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Hanldes new platform errors.
         * @param {import('../lf_models').LabeledError} new_labeled_error
        */
       const handlePlatformError = new_labeled_error => {
            if (new_labeled_error == null) return;

            labeled_error = new_labeled_error;

            console.error(new_labeled_error.toDevString());

            dialog_visible = true;

            setTimeout(() => {
                dialog_visible = false;
            }, display_error_duration);
       }
    
    /*=====  End of Methods  ======*/
    
</script>

<dialog open={dialog_visible} id="libery-feedback-floating-error-dialog" class:float-at-top={float_at_top}>
    <h3 id="lffed-error-message">
        {labeled_error?.toHumanString()}
    </h3>
</dialog>

<style>
    dialog#libery-feedback-floating-error-dialog[open] {
        scale: 1;

        @starting-style {
            scale: 0;
        }
    }

    dialog#libery-feedback-floating-error-dialog {
        --modal-position-modifier: 0.75;
        --modal-width: min(400px, 50vw);

        position: fixed;
        background: var(--grey-9);
        width: var(--modal-width);
        inset: calc(100vh * var(--modal-position-modifier)) auto auto calc(50dvw - calc(var(--modal-width) / 2));
        padding: var(--vspacing-3);
        border: 1px solid var(--danger-8);
        border-radius: var(--border-radius);
        scale: 0;
        transition: scale .3s ease-in, display .3s ease allow-discrete;
        z-index: var(--z-index-t-8);
    }

    h3#lffed-error-message {
        font-family: var(--font-read);
        font-size: var(--font-size-2);
        color: var(--grey-1);
        text-align: center;
    }
</style>