<script>
    import { confirmPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { createEventDispatcher } from "svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The given role's taxonomy.
         * @type {import("@models/Users").RoleTaxonomy}
         */ 
        export let the_role_taxonomy;

        
        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the Click event on the role delete button.
         * @param {MouseEvent} event 
         */
        const handleRoleDeleteClick = async (event) => {
            let user_choice = await confirmPlatformMessage({
                message_title: `Delete ${the_role_taxonomy.RoleLabel} role`,
                question_message: `Are you sure you want to delete the ${the_role_taxonomy.RoleLabel} role? All users with this role will lose it and if that's the only role they have, they will loose access to the platform.`,
                confirm_label: "Yes, delete role",
                cancel_label: "Cancel",
                danger_level: 2,
                auto_focus_cancel: true
            });

            if (user_choice === 1) {
                emitRoleDeletionRequested();
            }
        }

        /**
         * Emits the role-deletion-requested event.
         */
        const emitRoleDeletionRequested = () => {
            dispatch("role-deletion-requested");
        }
    
    /*=====  End of Methods  ======*/
     
</script>

<section id="role-editor-danger-zone">
    <header id="redz-header">
        <h3>
            Danger zone
        </h3>
    </header>
    <form id="redz-danger-zone-form"
        action="none" 
    >
        <fieldset id="dangerous-actions">
            <button class="dungeon-button-1 danger"
                type="button"
                on:click={handleRoleDeleteClick}
            >
                Delete role
            </button>
        </fieldset>
    </form>
</section>

<style>
    section#role-editor-danger-zone {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-1);
    }

    header#redz-header {
        & > h3:first-of-type {
            font-family: var(--font-read);
            font-weight: 600;
            color: var(--danger-7);
            font-size: var(--font-size-3);
        }
    }

    form#redz-danger-zone-form {
        display: flex;
        width: 100%;
        flex-direction: row-reverse;
        padding: var(--spacing-3);
        border: 1px solid var(--danger-9);
        border-radius: var(--border-radius);
    }
</style>