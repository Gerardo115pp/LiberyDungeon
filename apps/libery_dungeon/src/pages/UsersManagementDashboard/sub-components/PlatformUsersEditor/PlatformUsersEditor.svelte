<script>
    import { role_mode_enabled } from "@pages/UsersManagementDashboard/app_page_store";
    import UsersEditor from "./sub-components/UsersEditor.svelte";
    import { createEventDispatcher } from "svelte";
    import UsersEditorDangerZone from "./sub-components/UsersEditorDangerZone.svelte";
    import { deleteUser } from "@models/Users";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The subject user account which this component will edit.
         * @type {import("@models/Users").UserAccount}
         */ 
        export let subject_user_account;

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
    /**
     * Handles the user-data-changed event emitted by the UsersEditor component.
     */
    const handleUserDataChanged = () => {
        dispatch("user-data-changed");
    }

    /**
     * Handles the user-deletion-requested event emitted by the UsersEditorDangerZone component.
     */
    const handleUserDeletionRequested = async () => {
        let deleted = await deleteUser(subject_user_account.UUID);

        if (deleted) {
            emitUserDeleted();
        }
    }

    /**
     * Emits the user-deleted event.
     */
    const emitUserDeleted = () => {
        dispatch("user-deleted");
    }
        
    
    /*=====  End of Methods  ======*/
    
</script>

{#if !$role_mode_enabled && subject_user_account != null}
    <article id="users-editor-tools">
        <header id="uet-header">
            <h2>
                Change {subject_user_account.Username} data
            </h2>
            <section id="uet-header-instructions">
                <p class="dungeon-instructions">
                    Here you can change the user's username, password, and remove or add roles. For security reasons, you cannot see the user's current password as it is not stored in human readble text(<a href="https://stackoverflow.com/questions/1197417/why-are-plain-text-passwords-bad-and-how-do-i-convince-my-boss-that-his-treasur" target="_blank" rel="noopener noreferrer">This stackoverflow post</a> explains why in detail). But you can set any new password you want if you forgot the current password. It is recommended to use a password manager to store your passwords securely. <a href="https://bitwarden.com/" target="_blank" rel="noopener noreferrer">Bitwarden</a> is my personal recommendation but there are many good ones out there.
                </p>
                <p class="dungeon-instructions">
                    <strong>Roles</strong> ― When you add a role to your user, the user automatically gets all the permissions that role has. But the permissions always live on the role exclusively. they are not associated to the user directly. This means that if you remove a role from a user, the user loses all the permissions that role had(Unless the user had another role that also supplied the same permissions). 
                </p>
            </section>
        </header>
        <div id="uet-user-editor-tool">
            <UsersEditor 
                on:user-data-changed={handleUserDataChanged}
                {subject_user_account}
            />
        </div>
        <div id="uet-users-editor-danger-zone-wrapper">
            <UsersEditorDangerZone
                on:user-deletion-requested={handleUserDeletionRequested}
                the_user_account={subject_user_account}
            />
        </div>
    </article>
{/if}

<style>
    article#users-editor-tools {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-4);
    }

    header#uet-header {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);

        & > h2 {
            font-size: var(--font-size-h3);
            line-height: 1;
            text-align: center;
        }

        & > section#uet-header-instructions {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-2);
            padding: var(--spacing-1);
        }
    }
</style>