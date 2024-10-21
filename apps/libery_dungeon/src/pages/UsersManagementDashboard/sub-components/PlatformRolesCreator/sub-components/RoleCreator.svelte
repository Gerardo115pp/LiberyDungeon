<script>
    import DeleteableItem from "@components/ListItems/DeleteableItem.svelte";
    import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { compileGrants, createRole, getRoleTaxonomiesBelowHierarchy } from "@models/Users";
    import { all_grants } from "@pages/UsersManagementDashboard/app_page_store";
    import { createEventDispatcher } from "svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The new role label.
         * @type {string}
         */
        let new_role_label = "";

        /**
         * The new role hierarchy. Cannot be 0. lower numbers mean higher hierarchy, 0 is a super admin. Creating super admins is not currently considered in the system's 
         * design, but it will be in the future.
         * @type {number}
         */
        let new_role_hierarchy;

        /**
         * The new role's grant list.
         * @type {string[]}
         */
        let new_role_grants = [];

        /**
         * Whether or not the role is ready to be created.
         * @type {boolean}
         */
        let new_role_data_ready = false;
        
        /*----------  References  ----------*/
        
            /**
             * The role creation form.
             * @type {HTMLFormElement}
             */ 
            let the_role_creation_form;

            /**
             * The role label input.
             * @type {HTMLInputElement}
             */
            let the_role_label_input;

            /**
             * The role hierarchy input.
             * @type {HTMLInputElement}
             */
            let the_role_hierarchy_input;

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
            if (new_role_grants.includes(grant)) {
                emitPlatformMessage("Grant already added.");
                return;
            }

            new_role_grants = [...new_role_grants, grant];
        }

        /**
         * Checks whether all data needed to create a new role is ready.
         * @returns {boolean}
         */
        const checkRoleDataValidity = () => {
            let is_valid = the_role_creation_form.checkValidity();

            return is_valid;
        }

        /**
         * Emits the role-created event.
         */
        const emitRoleCreated = () => {
            dispatch("role-created");
        }

        /**
         * Handles the delete event from the DeleteableItem that wraps the grant.
         * @param {string} grant 
         */
        const handleGrantDelete = async (grant) => {
            let filtered_grants = new_role_grants.filter(g => g !== grant);

            new_role_grants = filtered_grants;
        }

        /**
         * Handles the change event on the new role hierarchy input.
         * @param {Event} event 
         */
        const handleNewRoleHierarchyChange = (event) => {
            console.log("Hierarchy input changed.");
            let input_valid = the_role_hierarchy_input.checkValidity();
            console.log("Hierarchy input valid: ", input_valid);

            if (!input_valid) {
                console.log("Hierarchy input invalid.");
                new_role_data_ready = false;
                new_role_hierarchy = undefined;
                return;
            }

            new_role_data_ready = checkRoleDataValidity();

            let new_hierarchy = parseInt(the_role_hierarchy_input.value);
            console.log("New hierarchy: ", new_hierarchy);

            if (isNaN(new_hierarchy)) {
                console.log("Hierarchy input is not a number.");
                new_role_data_ready = false;
                new_role_hierarchy = undefined;
                return;
            }

            if (new_hierarchy === 0) {
                emitPlatformMessage("Hierarchy cannot be 0.");

                new_role_data_ready = false;
                new_role_hierarchy = undefined;
                the_role_hierarchy_input.setCustomValidity("Hierarchy cannot be 0.");
            }

            new_role_hierarchy = new_hierarchy;

            refreshInheritedGrants();
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
         * Handles the keydown event on the role label input.
         * @param {KeyboardEvent} event 
         */
        const handleRoleLabelKeydown = (event) => {
            new_role_data_ready = checkRoleDataValidity();

            if (event.key === "Enter") {
                event.preventDefault();
                the_role_hierarchy_input.focus();
            }
        }

        /**
         * Handles the click event on the role creation button.
         * @param {MouseEvent} event 
         */
        const handleRoleCreationButtonClick = async (event) => {
            let created = registerNewRole(new_role_label, new_role_hierarchy, new_role_grants);
            
            if (created) {
                emitPlatformMessage(`Role ${new_role_label} created successfully.`);
                resetRoleCreator();
            }
        }

        /**
         * Fetches all the grants a role with a hierarchy equal to 'new_role_hierarchy' would inherit. the grant inheritance rules
         * state that a role inherits all the grants of roles with a lower hierarchy.
         * @returns {Promise<void>}
         */
        const refreshInheritedGrants = async () => {
            const roles_below_hierarchy = await getRoleTaxonomiesBelowHierarchy(new_role_hierarchy);

            const inherited_grants = compileGrants(roles_below_hierarchy);

            if (inherited_grants != null && inherited_grants.length > 0) {
                new_role_grants = inherited_grants;
            }
        }

        /**
         * Registers a new role with the given data.
         * @param {string} label
         * @param {number} hierarchy
         * @param {string[]} grants
         * @returns {Promise<boolean>}
         */
        const registerNewRole = async (label, hierarchy, grants) => {
            let created = await createRole(label, hierarchy, grants);

            if (created) {
                emitRoleCreated();
            }

            return created;
        }

        /**
         * Resets the component's data.
         */
        const resetRoleCreator = () => {
            new_role_label = "";
            new_role_hierarchy = undefined;
            new_role_grants = [];
            new_role_data_ready = false;
            the_role_creation_form.reset();
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<form id="role-creation-form"
    bind:this={the_role_creation_form}
    action="none"
>
    <label class="dungeon-input">
        <span class="dungeon-label">
            New role name
        </span>
        <input 
            type="text"
            bind:value={new_role_label}
            bind:this={the_role_label_input}
            on:keydown={handleRoleLabelKeydown}
            spellcheck="true"
            minlength="1"
            maxlength="64"
            pattern="{'[a-z_]{1,64}'}"
            required
        />
    </label>   
    <label class="dungeon-input">
        <span class="dungeon-label">
            Hierarchy
        </span>
        <input 
            type="text"
            bind:this={the_role_hierarchy_input}
            on:change={handleNewRoleHierarchyChange}
            inputmode="numeric"
            pattern="{'\\d{1,3}'}"
            min="1"
            step="1"
            max="100"
            required
        />
    </label>   
    <ul id="role-grants-list">
        {#each new_role_grants as grant}
            <DeleteableItem
                item_color="var(--grey-6)"
                on:item-deleted={() => handleGrantDelete(grant)}
            >
                {grant}
            </DeleteableItem>
        {/each}
        <DeleteableItem
            item_color="var(--grey-9)"
            is_protected
        >
            <input id="add-grant-input"
                bind:this={the_grant_adder_input}
                type="text"
                on:keydown={handleAddGrantInputKeydown}
                placeholder="Add a grant"
                list="grants-datalist"
                spellcheck="false"
                minlength="4"
                maxlength="64"
                pattern="{'[a-z_]{4,64}'}"
            >
        </DeleteableItem>
    </ul>
    <datalist id="grants-datalist">
        {#each $all_grants as grant}
            <option value={grant}/>
        {/each}
    </datalist>
    <fieldset id="role-creation-controls">
        <button class="dungeon-button-1"
            disabled={!new_role_data_ready}
            on:click|preventDefault={handleRoleCreationButtonClick}
            type="button"
        >
            Create role
        </button>
    </fieldset>
</form>

<style>
    form#role-creation-form {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);
    }
    
    /*=============================================
    =            Role grants form            =
    =============================================*/
    
        ul#role-grants-list {
            display: flex;
            flex-wrap: wrap;
            align-items: start;
            gap: var(--spacing-1);
            background: var(--grey-8);
            border-radius: var(--border-radius);
            padding: var(--spacing-1);
        }

        input#add-grant-input {
            box-sizing: border-box;
            width: 18ch;
            padding: 0 var(--spacing-1);
            background: transparent;
            color: var(--grey-1);
            border: none;
            outline: none;
        }

        input#add-grant-input::placeholder {
            color: var(--grey-5);
        }

        fieldset#role-creation-controls {
            display: flex;
            flex-direction: row-reverse;
        }
    
    
    /*=====  End of Role grants form  ======*/
    
    

</style>