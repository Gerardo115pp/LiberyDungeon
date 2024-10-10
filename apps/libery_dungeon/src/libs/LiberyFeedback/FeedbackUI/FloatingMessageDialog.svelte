<script>
    import { onDestroy, onMount } from "svelte";
    import { subscribeToPlatformMessages } from "../lf_utils";

    
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
         * The amount of time in milliseconds the dialog will be visible for after a messages is sent.
         * @type {number}
         * @default 3000
         */
        export let display_message_duration = 3000;

        /**
         * The message the will be displayed in the dialog.
         * @type {string}
         */
        let received_message;

        /**
         * Whether the dialog should be visible.
         * @type {boolean}
        */
        let dialog_visible = false;

        
        /*----------  Unsubscribers  ----------*/
        
            let platform_message_unsubscriber = () => {}; 
    
    /*=====  End of Properties  ======*/

    onMount(() => {
            platform_message_unsubscriber = subscribeToPlatformMessages(handlePlatformMessages)
    });

    onDestroy(() => {
        platform_message_unsubscriber();
    });

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Hanldes new platform messages.
         * @param {string} new_message
        */
       const handlePlatformMessages = new_message => {
            if (new_message == null || new_message === "") return;

            received_message = new_message;

            dialog_visible = true;

            setTimeout(() => {
                dialog_visible = false;
            }, display_message_duration);
       }
    
    /*=====  End of Methods  ======*/
    
</script>

<dialog open={dialog_visible} id="libery-feedback-floating-message-dialog" class:float-at-top={float_at_top}>
    <h3 id="lffed-message">
        {received_message}
    </h3>
</dialog>

<style>
    dialog#libery-feedback-floating-message-dialog[open] {
        scale: 1;

        @starting-style {
            scale: 0;
        }
    }

    dialog#libery-feedback-floating-message-dialog {
        --modal-position-modifier: 0.75;
        --modal-width: min(400px, 50vw);

        position: fixed;
        background: var(--grey-9);
        width: var(--modal-width);
        inset: calc(100vh * var(--modal-position-modifier)) auto auto calc(50dvw - calc(var(--modal-width) / 2));
        padding: var(--vspacing-3);
        border: 1px solid var(--main-dark-color-5);
        border-radius: var(--border-radius);
        scale: 0;
        transition: scale .3s ease-in, display .3s ease allow-discrete;
        z-index: var(--z-index-t-7);
    }

    h3#lffed-message {
        font-family: var(--font-read);
        font-size: var(--font-size-2);
        color: var(--grey-1);
        text-align: center;
    }
</style>