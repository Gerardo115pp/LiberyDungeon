<script>
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { addGrantToRole, getRoleTaxonomiesBelowHierarchy, removeGrantFromRole } from "@models/Users";
    import { all_grants } from "@pages/UsersManagementDashboard/app_page_store";
    import { createEventDispatcher } from "svelte";

    /*=============================================
    =            Properties            =
    =============================================*/

        /**
         * The subject role new grants.
         * @type {string[]}
         */
        let subject_new_grants = [];

        /**
         * The grants the subject role already has.
         * @type {string[]}
         */
        let subject_existing_grants = [];

        /**
         * The subject role's taxonomy.
         * @type {import("@models/Users").RoleTaxonomy}
         */
        export let subject_role_taxonomy;
        $: if (subject_role_taxonomy != null) {
            loadExistingGrants();
        }
        
        /*----------  References  ----------*/
        
            /**
             * The role grant linker form.
             * @type {HTMLFormElement}
             */ 
            let the_role_grant_linker_form;

            /**
             * The grant adder input.
             * @type {HTMLInputElement}
             */
            let the_grant_adder_input; 
        
        const dispatch = createEventDispatcher();
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/


        /**
         * Adds a grant to the new role's grant list.
         * @param {string} grant 
         */
        const addGrant = (grant) => {
            if (subject_new_grants.includes(grant)) {
                emitPlatformMessage("Grant already added.");
                return;
            }

            subject_new_grants = [...subject_new_grants, grant];
        }

        /**
         * Checks whether all data needed to create a new role is ready.
         * @returns {boolean}
         */
        const checkRoleDataValidity = () => {
            let is_valid = the_role_grant_linker_form.checkValidity();

            return is_valid;
        }

        /**
         * Emits the role grants change event. this event shall trigger after at least one grant is added to the role.
         */
        const emitRoleGrantsChange = () => {
            dispatch("role-grants-changed");
        }

        /**
         * Handles the delete event from the DeleteableItem that wraps the grant.
         * @param {string} grant 
         */
        const handleNewGrantDelete = async (grant) => {
            let filtered_grants = subject_new_grants.filter(g => g !== grant);

            subject_new_grants = filtered_grants;
        }

        /**
         * Handles the delete event from the DeleteableItem that wraps an existing grant.
         * @param {string} grant 
         */
        const handleExistingGrantDelete = async (grant) => {
            if (subject_role_taxonomy.RoleHierarchy === 0) return; // Super admins grants are immutable

            let filtered_grants = subject_existing_grants.filter(g => g !== grant);

            subject_existing_grants = filtered_grants;

            let removed = await removeGrantFromRole(subject_role_taxonomy.RoleLabel, grant);

            if (removed) {
                emitPlatformMessage(`Removed grant ${grant} from ${subject_role_taxonomy.RoleLabel}`);
                emitRoleGrantsChange();
            }
        }

        /**
         * Handles the keydown event on the add grant input.
         * @param {KeyboardEvent} event 
         */
        const handleAddGrantInputKeydown = (event) => {
            if (event.key === "Enter") {
                event.preventDefault();
                let is_valid = the_grant_adder_input.checkValidity();

                if (!is_valid) return;

                addGrant(the_grant_adder_input.value);

                the_grant_adder_input.value = "";
            }
        }

        /**
         * Handles the click event on the add to role button.
         * @param {MouseEvent} event 
         */
        const handleLinkGrantsToRoleBtnClick = async (event) => {
            if (subject_new_grants?.length === 0 || !checkRoleDataValidity()) return;

            /** @type {string[]}*/
            let failed_to_link_grants = await linkAllNewGrants();

            if (failed_to_link_grants.length > 0) {
                emitPlatformMessage(`Failed to link the following grants: ${failed_to_link_grants.join(", ")}`);
            } else {
                emitPlatformMessage(`Added ${subject_new_grants.length} grants to ${subject_role_taxonomy.RoleLabel}`);
            }

            if (failed_to_link_grants.length !== subject_new_grants.length) {
                emitRoleGrantsChange();
            }

            resetGrantLinker();
        }

        /**
         * Takes the grants on subject_role_taxonomy and stores them in subject_existing_grants.
         * Called when the subject_role_taxonomy reference changes.
         */
        function loadExistingGrants() {
            if (subject_role_taxonomy == null) return;

            subject_existing_grants = subject_role_taxonomy.RoleGrants;
        }

        /**
         * Adds all the grants on the subject_new_grants list to the system role on the subject_role_taxonomy. It returns a list
         * with all the grants that could not be added. if the list is empty, everything was added successfully.
         * @returns {Promise<string[]>}
         */
        const linkAllNewGrants = async () => {
            /** @type {string[]} */
            let failed_to_link_grants = [];

            for (let new_grant of subject_new_grants) {
                let linked = await addGrantToRole(subject_role_taxonomy.RoleLabel, new_grant);

                if (!linked) {
                    failed_to_link_grants.push(new_grant);
                }
            }

            return failed_to_link_grants;
        }
        

        /**
         * Links a given grant to a given role label.
         * @param {string} role_label
         * @param {string} new_grant
         * @returns {Promise<boolean>}
         */
        const linkGrantToRole = async (role_label, new_grant) => {
            let linked = await addGrantToRole(role_label, new_grant);

            return linked;
        }

        /**
         * Resets the component's data.
         */
        const resetGrantLinker = () => {
            subject_new_grants = [];
            the_role_grant_linker_form.reset();
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<form id="role-grant-linker-form"
    bind:this={the_role_grant_linker_form}
    action="none"
>
    <header id="rglf-header">
        <ul id="rglf-role-details">
            <li class="role-detail-field">
                <p>
                    <strong>Role</strong><span>{subject_role_taxonomy.RoleLabel}</span>
                </p>
            </li>
            <li class="role-detail-field">
                <p>
                    <strong>Hierarchy</strong><span>{subject_role_taxonomy.RoleHierarchy}</span>
                </p>                    
            </li>
        </ul>
    </header>
    <ul id="rglf-role-grants-list"
        class="dungeon-tag-list"
    >
        <!-- Existing grants -->
        {#each subject_existing_grants as grant}
            <DeleteableItem
                item_color="var(--grey-4)"
                on:item-deleted={() => handleExistingGrantDelete(grant)}
            >
                {grant}
            </DeleteableItem>
        {/each}
        <!-- New grants -->
        {#each subject_new_grants as grant}
            <DeleteableItem
                item_color="var(--main-dark)"
                on:item-deleted={() => handleNewGrantDelete(grant)}
            >
                {grant}
            </DeleteableItem>
        {/each}
        <DeleteableItem
            item_color="var(--grey-9)"
            is_protected
        >
            <input id="rglf-add-grant-input"
                class="tag-creator"
                bind:this={the_grant_adder_input}
                type="text"
                on:keydown={handleAddGrantInputKeydown}
                placeholder="Add additional grants"
                list="rglf-grants-datalist"
                spellcheck="false"
                minlength="4"
                maxlength="64"
                pattern="{'[a-z_]{4,64}'}"
            >
        </DeleteableItem>
    </ul>
    <datalist id="rglf-grants-datalist">
        {#each $all_grants as grant}
            <option value={grant}/>
        {/each}
    </datalist>
    <fieldset id="rglf-role-editor-controls">
        <button class="dungeon-button-1"
            disabled={subject_new_grants.length === 0}
            on:click|preventDefault={handleLinkGrantsToRoleBtnClick}
            type="button"
        >
            Save grants
        </button>
    </fieldset>
</form>

<style>
    form#role-grant-linker-form {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);
    }
    
    /*=============================================
    =            Header            =
    =============================================*/
    
        ul#rglf-role-details {
            display: flex;
            column-gap: var(--spacing-2);
        }

        li.role-detail-field {
            
            & > p {
                font-size: var(--font-size-);
                display: flex;
                gap: var(--spacing-1);  
            };

            & strong {
                color: var(--main);
            }

            & strong::after {
                content: ":";
            }

            & span {
                color: var(--grey-2);
                font-weight: 500;
            }
        } 
    
    /*=====  End of Header  ======*/
    
    /*=============================================
    =            Role grants form            =
    =============================================*/

        input#rglf-add-grant-input {
            width: 28ch;
        }

        fieldset#rglf-role-editor-controls {
            display: flex;
            flex-direction: row-reverse;
        }
    
    
    /*=====  End of Role grants form  ======*/

</style>