<script>
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { all_grants } from "../../../app_page_store";
    import { createGrant, deleteGrant, isSystemGrant } from "@models/Users";
    import { confirmPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The grant creation form.
         * @type {HTMLFormElement}
         */
        let the_grant_creation_form; 
        
        /**
         * The new grant label.
         * @type {string}
         */
        let new_grant_label = "";


        /**
         * Whether or not the grant is ready to be created.
         * @type {boolean}
         */
        let new_grant_value_ready = false;
    
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods             =
    =============================================*/
    
        /**
         * Handles the delete event from the DeleteableItem that wraps the grant.
         * @param {string} grant 
         */ 
        const handleGrantDelete = async (grant) => {
            if (isSystemGrant(grant)) return;

            let user_confirmation_choice = await confirmPlatformMessage({
                message_title: "Delete grant",
                question_message: `Are you sure you want to delete the grant "${grant}"? Be aware that this will not make actions that require this grant stop requiring it, but it will remove it from all user roles and will not be available to be assigned for new roles until it is created again.`,
                confirm_label: "Remove the role grant",
                cancel_label: "Cancel",
                auto_focus_cancel: true,
                danger_level: 2
            });

            if (user_confirmation_choice !== 1) return;

            let grant_deleted = false;

            grant_deleted = await deleteGrant(grant);

            if (!grant_deleted) return;

            let new_grants = $all_grants.filter(g => g !== grant);
            all_grants.set(new_grants);
        }

        /**
         * Handles the keydown event on the grant creation input.
         * @param {KeyboardEvent} event 
         */
        const handleGrantCreationInputKeydown = (event) => {
            new_grant_value_ready = checkGrantLabelValidity();
            console.log("new_grant_value_ready: ", new_grant_value_ready);

            if (event.key === "Enter" && new_grant_value_ready) {
                event.preventDefault();
                createCurrentGrantLabel();
            }
        }

        /**
         * Handles the click event on the grant creation button
         * @param {MouseEvent} event 
         */
        const handleGrantCreationButtonClick = (event) => {
            event.preventDefault();
            createCurrentGrantLabel();
        }

        /**
         * Checks if the grant label is valid an ready to be created.
         * @returns {boolean}
         */
        const checkGrantLabelValidity = () => {
            return the_grant_creation_form.checkValidity();
        }


        /**
         * Takes the value of the new grant label and creates a new grant with it.
         * If successful, resets the_grant_creation_form.
         * @returns {Promise<void>}
         */
        const createCurrentGrantLabel = async () => {
            let grant_created = await createGrant(new_grant_label);

            if (!grant_created) return;

            all_grants.set([...$all_grants, new_grant_label]);

            new_grant_label = "";

            the_grant_creation_form.reset();
        }
            
         
    
    /*=====  End of Methods   ======*/
    
</script>

<form id="grant-creation-form"
    action="none"
    bind:this={the_grant_creation_form}
>
    <label class="dungeon-input">
        <span class="dungeon-label">
            New grant
        </span>
        <input 
            type="text"
            bind:value={new_grant_label}
            on:keydown={handleGrantCreationInputKeydown}
            spellcheck="false"
            minlength="4"
            maxlength="64"
            pattern="{'[a-z_]{4,64}'}"
            required
        />
    </label>
    <button class="dungeon-button-1"
        disabled={!new_grant_value_ready}
        on:click={handleGrantCreationButtonClick}
        type="submit"
    >
        Add grant
    </button>
</form>
<div id="platform-grants-wrapper">
    <ul id="platform-grants">
        {#each $all_grants as grant}
            <DeleteableItem
                item_color="var(--grey-6)"
                on:item-deleted={() => handleGrantDelete(grant)}
            >
                {grant}
            </DeleteableItem>
        {/each}
    </ul>
</div>

<style>
    form#grant-creation-form {
        display: flex;
        column-gap: var(--spacing-2);

        & > label {
            flex-grow: 2;
        }
    }

    ul#platform-grants {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        background-color: var(--grey-8);
        padding: var(--spacing-2);
        border-radius: var(--border-radius);
        gap: var(--spacing-2);
    }
</style>