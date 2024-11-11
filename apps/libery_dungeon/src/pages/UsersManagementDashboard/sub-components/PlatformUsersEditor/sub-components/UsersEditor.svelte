<script>
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { confirmPlatformMessage, emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { addUserToRole, changeUserPassword, changeUserUsername, findHighestRoleHierarchy, getRoleTaxonomy, removeUserFromRole } from "@models/Users";
    import { all_roles } from "@pages/UsersManagementDashboard/app_page_store";
    import { createEventDispatcher } from "svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** 
         * The given user account the component will edit.
         * @type {import("@models/Users").UserAccount}
         */ 
        export let subject_user_account;
        $: if (subject_user_account != null) {
            loadUserInformation();
        }

        /**
         * The roles the user currently has. 
         * @type {import("@models/Users").RoleTaxonomy[]}
         */
        let user_existing_roles = [];

        /**
         * The user's new roles.
         * @type {import("@models/Users").RoleTaxonomy[]}
         */
        let user_new_roles = [];

        let account_highest_hierarchy = Infinity

        /**
         * The user's new username.
         * @type {string}
         */
        let user_new_username = "";


        /**
         * The user's new password.
         * @type {string}
         */
        let user_new_password = "";
        
        /**
         * Whether there is new user data to modify the user account in any way.
         * @type {boolean}
         */
        let account_modified = false;
        
        /*----------  References  ----------*/

            /**
             * The user's username editor input.
             * @type {HTMLInputElement}
             */
            let the_new_username_input;

            /**
             * The user's new password editor input.
             * @type {HTMLInputElement}
             */
            let the_new_password_input;

            /**
             * The user new roles assigner input.
             * @type {HTMLInputElement}
             */
            let the_new_roles_assigner_input;

            /**
             * The entire user editor form.
             * @type {HTMLFormElement}
             */
            let the_user_editor_form;

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Checks whether there is any new user data to modify the user account in any way. It does this by
         * checking the values of the user_new_username, user_new_password, and if there are any roles in the
         * user_new_roles array.
         * @returns {boolean}
         */
        const checkAccountModified = () => {
            let new_account_data = user_new_username !== "";

            new_account_data = new_account_data || user_new_password !== "";

            new_account_data = new_account_data || user_new_roles.length > 0;

            return new_account_data;
        }

        /**
         * Called when ever the user account gets a new role, if the new role has a higher role hierarchy than the
         * current highest role hierarchy will change.
         * @param {number} new_role_hierarchy 
         */
        const changeHighestRoleHierarchyIfHigher = (new_role_hierarchy) => {
            if (new_role_hierarchy < account_highest_hierarchy) {
                account_highest_hierarchy = new_role_hierarchy;
            }
        }

        /**
         * Returns whether the user account roles and the user_existing_roles are different.
         * @returns {boolean}
         */
        const existingRolesChanged = () => {
            let existing_roles_set = new Set(user_existing_roles.map(role_taxonomy => role_taxonomy.RoleLabel));
            let user_roles_set = new Set(subject_user_account.Roles.map(role => role.RoleLabel));

            let symmetric_difference = existing_roles_set.symmetricDifference(user_roles_set);
            console.log("Symmetric difference:", symmetric_difference);

            let changed = symmetric_difference.size > 0; 

            return changed;
        }

        /**
         * Alert component error.
         * @param {string} context_label
         * @param {string} human_readable_message
         * @param {Object | null} local_variables
         * @param {boolean} log_objects
         */
        const alertComponentError = (context_label, human_readable_message, local_variables, log_objects=false) => {
            let environment_variable = new VariableEnvironmentContextError(context_label);

            environment_variable.addVariable("New username", user_new_username);
            environment_variable.addVariable("New password", user_new_password);
            environment_variable.addVariable("Account modified", account_modified);

            let labeled_error = new LabeledError(environment_variable, human_readable_message, lf_errors.ERR_PROCESSING_ERROR);

            if (log_objects) {
                console.log("Component object variables:");
                
                console.log("subject_user_account:", subject_user_account);
                console.log("user_existing_roles:", user_existing_roles);
                console.log("user_new_roles:", user_new_roles);
                console.log("the_new_username_input:", the_new_username_input);
                console.log("the_new_password_input:", the_new_password_input);
                console.log("the_new_roles_assigner_input:", the_new_roles_assigner_input);
                console.log("the_user_editor_form:", the_user_editor_form);
            }

            if (local_variables != null) {
                console.log("LOCAL VARIABLES");

                for(let lv of Object.keys(local_variables)) {
                    // @ts-ignore
                    console.log(`${lv}`, local_variables[lv]);
                }
            }

            labeled_error.alert();
        }

        /**
         * Emits the user-data-changed event.
         */
        const emitUserDataChanged = () => {
            dispatch("user-data-changed");
        }

        /**
         * Returns the highest role hierarchy comparing the roles in the user_existing_roles and user_new_roles arrays.
         * @returns {number}
         */
        const getHighestUnofficialRoleHierarchy = () => {
            const official_hierarchy = findHighestRoleHierarchy(user_existing_roles);
            const unofficial_hierarchy = findHighestRoleHierarchy(user_new_roles);

            return Math.min(official_hierarchy, unofficial_hierarchy);
        }
    
        /**
         * Handles the keydown event on the user's new username input.
         * @param {KeyboardEvent} event 
         */
        const handleUsernameKeydown = async (event) => {
            let username_valid = the_new_username_input.checkValidity();

            account_modified = username_valid ? true : checkAccountModified();

            if (event.key === "Enter") {
                the_new_username_input.blur();
            }
        }

        /**
         * Handles the keydown event on the user's new password input.
         * @param {KeyboardEvent} event 
         */
        const handlePasswordKeydown = async (event) => {
            let password_valid = user_new_password.length >= 5 && the_new_password_input.checkValidity();

            account_modified = password_valid ? true : checkAccountModified();

            if (event.key === "Enter") {
                the_new_password_input.blur();
            }
        }

        /**
         * Handles the keydown event on the user's new roles assigner input.
         * @param {KeyboardEvent} event 
         */
        const handleRolesAssignerKeydown = async (event) => {
            if (event.key !== "Enter") return;

            let new_role_label = the_new_roles_assigner_input.value;

            if (subject_user_account.isInRole(new_role_label)) {
                emitPlatformMessage(`The user is already in the role '${new_role_label}'.`);
                return;
            }

            let new_role_taxonomy = await getRoleTaxonomy(new_role_label);

            if (new_role_taxonomy == null) {
                emitPlatformMessage(`The role '${new_role_label}' does not exist.`);
                return;
            }

            the_new_roles_assigner_input.value = ""; // only clear on success. in case of error, allow the user to correct a typo or something.
            
            account_modified = true;

            user_new_roles = [...user_new_roles, new_role_taxonomy];
            changeHighestRoleHierarchyIfHigher(new_role_taxonomy.RoleHierarchy);
        }

        /**
         * Handles an the removal of a role the user is currently in.
         * @param {string} role_label 
         */
        const handleExistingRoleDelete = async (role_label) => {
            let new_existing_roles = user_existing_roles.filter(role_taxonomy => role_taxonomy.RoleLabel !== role_label);
            
            user_existing_roles = new_existing_roles;
            
            account_modified = user_existing_roles.length !== subject_user_account.Roles.length || account_modified;

            updateHighestRoleHierarchy();
        }

        /**
         * Handles the deletion of a user new role. this just involves removing it from the user_new_roles array.
         * @param {string} role_label 
         */
        const handleNewRoleDelete = async (role_label) => {
            let new_new_roles = user_new_roles.filter(role_taxonomy => role_taxonomy.RoleLabel !== role_label);

            user_new_roles = new_new_roles;

            account_modified = user_new_roles.length !== 0 || checkAccountModified();

            updateHighestRoleHierarchy();
        }

        /**
         * Handles the click event on save account data button.
         * @param {MouseEvent} event 
         */
        const handleAccountDataSaveBtnClick = async (event) => {
            console.log("Account_modified:", account_modified);

            processAccountChanges();
        }

        /**
         * Loads the user information into the reactive variables from subject_user_account.
         */
        function loadUserInformation() {
            user_existing_roles = subject_user_account.Roles;
            user_new_username = subject_user_account.Username;
            user_new_roles = [];
            user_new_password = "";
            account_highest_hierarchy = subject_user_account.HighestRoleHierarchy;
            account_modified = false;
        }

        /**
         * Returns all the roles the user has lost.
         * @returns {import("@models/Users").RoleTaxonomy[]}
         */
        const lostRoles = () => {
            let account_roles = new Map();
            subject_user_account.Roles.forEach(role => account_roles.set(role.RoleLabel, role));

            let existing_roles_array_set = new Set(user_existing_roles.map(role_taxonomy => role_taxonomy.RoleLabel));
            let account_roles_set = new Set(subject_user_account.Roles.map(role => role.RoleLabel));

            let lost_role_labels = account_roles_set.difference(existing_roles_array_set);

            let lost_roles = [];

            for (let role_label of lost_role_labels) {
                lost_roles.push(account_roles.get(role_label));
            }

            return lost_roles;
        }
        
        /*=============================================
        =            Account changes processors            =
        =============================================*/
        
            /**
             * Process the account changes and save them to the server.
             */
            const processAccountChanges = async () => {
                if (!account_modified) return;

                if (user_existing_roles.length === 0) {
                    emitPlatformMessage("You cannot remove all roles from a user. If you want to remove the user, check out the danger zone down below.");
                    return;
                }

                let username_changed = user_new_username !== subject_user_account.Username;
                let password_changed = user_new_password !== "" && user_new_password.length >= 5;
                let roles_added = user_new_roles.length > 0;
                let roles_removed = existingRolesChanged();

                let err_local_variables = {
                    username_changed,
                    password_changed,
                    roles_added,
                    roles_removed
                };

                let changes_feedback_message = "Successfully ";
                let changes_feedback_message_parts = [];
                
                if (username_changed) {
                    let username_change_error = await processUsernameChange();

                    console.log("Username change error:", username_change_error);
                    if (username_change_error !== "") {
                        alertComponentError("In UserEditor.processAccountChanges while processing username change", username_change_error, err_local_variables, true);
                        return;
                    }

                    changes_feedback_message_parts.push("changed the username");
                }

                if (password_changed) {
                    let password_change_error = await processPasswordChange();

                    if (password_change_error !== "") {
                        alertComponentError("In UserEditor.processAccountChanges while processing password change", password_change_error, err_local_variables, true);
                        return;
                    }

                    changes_feedback_message_parts.push("changed the password");
                }

                if (roles_added) {
                    let new_roles_error = await processNewRoles();

                    if (new_roles_error !== "") {
                        alertComponentError("In UserEditor.processAccountChanges while processing new roles", new_roles_error, err_local_variables, true);
                        return;
                    }

                    changes_feedback_message_parts.push(`added ${user_new_roles.length} new roles`);
                }

                if (roles_removed) {
                    let lost_roles = lostRoles();
                    console.log("Lost roles:", lost_roles);
                    let lost_roles_error = await processLostRoles(lost_roles);

                    if (lost_roles_error !== "") {
                        alertComponentError("In UserEditor.processAccountChanges while processing lost roles", lost_roles_error, err_local_variables, true);
                        return;
                    }

                    changes_feedback_message_parts.push(`removed ${lost_roles.length} roles`);
                }

                changes_feedback_message += changes_feedback_message_parts.join(", ") + ".";

                resetComponentState();
                emitUserDataChanged();

                await confirmPlatformMessage({
                    message_title: "Account changes saved",
                    question_message: changes_feedback_message,
                    confirm_label: "Ok",
                    auto_focus_cancel: false,
                    danger_level: -1,
                });
            }

            /**
             * Only call from processAccountChanges!. 
             * Processes the username change. Returns an empty string if successful, or an error message if not.
             * @returns {Promise<string>}
             */
            const processUsernameChange = async () => {
                let username_valid = the_new_username_input.checkValidity();

                if (!username_valid) {
                    return "The username is invalid.";
                }

                let username_changed = await changeUserUsername(subject_user_account.UUID, user_new_username);

                if (!username_changed) {
                    return "Server error while changing the username.";
                }

                return "";
            }

            /**
             * Only call from processAccountChanges!. 
             * * Processes the password change. Returns an empty string if successful, or an error message if not.
             * @returns {Promise<string>}
             */
            const processPasswordChange = async () => {
                let password_valid = user_new_password.length >= 5 && the_new_password_input.checkValidity();

                if (!password_valid) {
                    return "The password is invalid.";
                }

                let password_changed = await changeUserPassword(subject_user_account.UUID, subject_user_account.Username, user_new_password);

                if (!password_changed) {
                    return "Server error while changing the password.";
                }

                return "";
            }

            /**
             * Only call from processAccountChanges!.
             * Processes new roles by adding them to the user account. Returns an empty string if successful, or an error message if not.
             * @returns {Promise<string>}
             */
            const processNewRoles = async () => {
                if (user_new_roles.length === 0) return "";

                for (let new_role_taxonomy of user_new_roles) {
                    let role_added = await addUserToRole(subject_user_account.Username, new_role_taxonomy.RoleLabel);

                    if (!role_added) {
                        return `Server error while adding the role '${new_role_taxonomy.RoleLabel}'.`;
                    }
                }

                return "";
            }

            /**
             * Only call from processAccountChanges!.
             * Processes lost roles by removing them from the user account. Returns an empty string if successful, or an error message if not.
             * @param {import("@models/Users").RoleTaxonomy[]} lost_roles
             * @returns {Promise<string>}
             */
            const processLostRoles = async lost_roles => {
                if (lost_roles.length === 0) return "";

                for (let lost_role of lost_roles) {
                    let role_removed = await removeUserFromRole(subject_user_account.Username, lost_role.RoleLabel);

                    if (!role_removed) {
                        return `Server error while removing the role '${lost_role.RoleLabel}'.`;
                    }
                }

                return "";
            }
        
        /*=====  End of Account changes processors  ======*/

        /**
         * Updates the value of account_highest_hierarchy to the highest role hierarchy in the user_existing_roles and user_new_roles arrays.
         */
        function updateHighestRoleHierarchy() {
            if (!existingRolesChanged()) return;

            account_highest_hierarchy = getHighestUnofficialRoleHierarchy();
        }

        /**
         * Resets the component state to an user data unchanged state.
         */
        const resetComponentState = () => {
            user_new_password = "";
            account_modified = false;
        }
    
    /*=====  End of Methods  ======*/

</script>

<form id="user-editor-form"
    action="none"
>
    <label class="dungeon-input">
        <span class="dungeon-label">
            Username
        </span>
        <input 
            type="text"
            bind:value={user_new_username}
            bind:this={the_new_username_input}
            on:keydown={handleUsernameKeydown}
            spellcheck="false"
            minlength="2"
            maxlength="64"
            pattern="{'[a-zA-Z0-9_\\-\\.@]{2,}'}"
            required
        />
    </label>
    <label class="dungeon-input">
        <span class="dungeon-label">
            New Password
        </span>
        <input 
            type="text"
            bind:value={user_new_password}
            bind:this={the_new_password_input}
            on:keydown={handlePasswordKeydown}
            spellcheck="false"
            minlength="5"
            maxlength="255"
            pattern="{'[a-zA-Z0-9_\\-\\.@]{2,}'}"
        />
    </label>
    <p id="user-highest-role-hierarchy">
        <strong>Highest role hierarchy:</strong> <span>{account_highest_hierarchy === Infinity ? "None" : account_highest_hierarchy}</span>
    </p>
    <ul id="user-roles-list"
        class="dungeon-tag-list"
    >
        {#each user_existing_roles as role_taxonomy}
            <DeleteableItem
                item_color="var(--grey-6)"
                on:item-deleted={() => handleExistingRoleDelete(role_taxonomy.RoleLabel)}
                is_protected={subject_user_account.Roles.length === 1}
            >
                {role_taxonomy.RoleLabel}
            </DeleteableItem>
        {/each}
        {#each user_new_roles as new_role_taxonomy}
            <DeleteableItem
                item_color="var(--main-dark)"
                on:item-deleted={() => handleNewRoleDelete(new_role_taxonomy.RoleLabel)}
            >
                {new_role_taxonomy.RoleLabel}
            </DeleteableItem>
        {/each}
        <DeleteableItem
            item_color="var(--grey-9)"
            is_protected
        >
            <input id="add-grant-input"
                class="tag-creator"
                bind:this={the_new_roles_assigner_input}
                type="text"
                on:keydown={handleRolesAssignerKeydown}
                placeholder="Additional role"
                list="roles-available-datalist"
                spellcheck="false"
                minlength="4"
                maxlength="64"
                pattern="{'[a-z_]{4,64}'}"
            >
        </DeleteableItem>
    </ul>
    <datalist id="roles-available-datalist">
        {#each $all_roles as role_label}
            {#if !subject_user_account.isInRole(role_label)}
                <option value={role_label}/>
            {/if}
        {/each}
    </datalist>
    <fieldset id="user-edition-controls">
        <button class="dungeon-button-1"
            disabled={!account_modified}
            on:click|preventDefault={handleAccountDataSaveBtnClick}
            type="button"
        >
            Save account
        </button>
    </fieldset>
</form>

<style>
    form#user-editor-form {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);
    }

    p#user-highest-role-hierarchy {
        font-size: var(--font-size-1);
        line-height: 1;
        text-transform: lowercase;

        & > strong {
            color: var(--main);
        }
    }

    fieldset#user-edition-controls {
        display: flex;
        flex-direction: row-reverse;
    }
</style>