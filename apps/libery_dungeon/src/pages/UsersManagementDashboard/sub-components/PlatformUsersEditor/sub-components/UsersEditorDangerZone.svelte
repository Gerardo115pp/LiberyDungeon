<script>
    import { confirmPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { createEventDispatcher } from "svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The user account.
         * @type {import("@models/Users").UserAccount}
         */ 
        export let the_user_account;

        
        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the Click event on the delete user button.
         * @param {MouseEvent} event 
         */
        const handleUserDeleteClick = async (event) => {
            let user_choice = await confirmPlatformMessage({
                message_title: `Delete the user account '${the_user_account.Username}'`,
                question_message: `Are you sure you want to delete the user account '${the_user_account.Username}'? This action cannot be undone.`,
                confirm_label: "Yes, bye bye user",
                cancel_label: "Cancel",
                danger_level: 2,
                auto_focus_cancel: true
            });

            if (user_choice === 1) {
                emitUserDeletionRequested();
            }
        }

        /**
         * Emits the user-deletion-requested event.
         */
        const emitUserDeletionRequested = () => {
            dispatch("user-deletion-requested");
        }
    
    /*=====  End of Methods  ======*/
     
</script>

<section id="user-editor-danger-zone">
    <header id="uedz-header">
        <h3>
            Danger zone
        </h3>
    </header>
    <form id="uedz-danger-zone-form"
        action="none" 
    >
        <fieldset id="uedz-dangerous-actions">
            <button class="dungeon-button-1 danger"
                type="button"
                on:click={handleUserDeleteClick}
            >
                Delete account
            </button>
        </fieldset>
    </form>
</section>

<style>
    section#user-editor-danger-zone {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-1);
    }

    header#uedz-header {
        & > h3:first-of-type {
            font-family: var(--font-read);
            font-weight: 600;
            color: var(--danger-7);
            font-size: var(--font-size-3);
        }
    }

    form#uedz-danger-zone-form {
        display: flex;
        width: 100%;
        flex-direction: row-reverse;
        padding: var(--spacing-3);
        border: 1px solid var(--danger-9);
        border-radius: var(--border-radius);
    }
</style>