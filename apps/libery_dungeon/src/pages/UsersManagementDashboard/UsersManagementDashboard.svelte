<script>
    /*=============================================
    =            Imports            =
    =============================================*/
        import { goto } from "$app/navigation";
        import Page from "@app/routes/+page.svelte";
        import UserCreation from "@components/Users/UsersCreation/UserCreation.svelte";
        import { confirmPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
        import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
        import global_hotkeys_manager from "@libs/LiberyHotkeys/libery_hotkeys";
        import { hotkeys_sheet_visible } from "@stores/layout";
        import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import { 
            getAllUsers,
            getAllRoleLabels,
            getAllGrants,
            getUserRoles,
            getRoleTaxonomy,
            UserAccount,
            getUserAccountFromUserEntry,
        } from "@models/Users";
        import { 
            confirmAccessState,
            has_user_access,
            access_state_confirmed,
            current_user_identity,
            banCurrentUser,
        } from "@stores/user";
        import { onDestroy, onMount } from "svelte";
        import { browser } from "$app/environment";
        import PlatformRolesCreator from "./sub-components/PlatformRolesCreator/PlatformRolesCreator.svelte";
        import { 
            role_mode_enabled,
            all_users,
            all_roles,
            all_grants,
        } from "./app_page_store"
        import PlatformRolesEditor from "./sub-components/PlatformRolesEditor/PlatformRolesEditor.svelte";
    import PlatformUsersEditor from "./sub-components/PlatformUsersEditor/PlatformUsersEditor.svelte";
    
    /*=====  End of Imports  ======*/
    
    /*=============================================
    =            Properties            =
    =============================================*/

        
        /*=============================================
        =            Hotkeys            =
        =============================================*/
        
           const hotkeys_context_name = "user_management_dashboard";
        
        /*=====  End of Hotkeys  ======*/
    
        /**
         * Whether it's been confirmed that the user is a super admin. Only as super admin has the
         * 'grant_option' grant, which is necessary for must user management operations.
         * @type {boolean}
         */ 
        let super_admin_confirmed = false;

        /**
         * The keyboard movement index.
         * @type {number}
         */
        let keyboard_movement_index = 0;

        /**
         * The role label the user has selected for editing.
         * @type {string}
         */
        let selected_role_label = "";
        // let selected_role_label = "privileged_visitor"; // Only for test

        /**
         * The user account the user has selected for editing.
         * @type {import("@models/Users").UserAccount}
         */
        let selected_user_account;
    
    /*=====  End of Properties  ======*/
   
    onMount(async () => {
        
        if (!$access_state_confirmed) {
            let could_confirm = await confirmAccessState();
        }

        super_admin_confirmed = $current_user_identity.canGrant();

        if (!super_admin_confirmed) {
            informUnauthorizedUser();
            return;
        }

        defineInitialData();

        defineComponentHotkeys();
    });
    
    onDestroy(() => {
        if (!browser) return;

        global_hotkeys_manager.dropContext(hotkeys_context_name);
    });

    /*=============================================
    =            Methods            =
    =============================================*/
        
        /*=============================================
        =            Hotkeys            =
        =============================================*/
            // Only call directly defineComponentHotkeys. All the other methods in this section should only be called by the global hotkeys manager. 

            /**
             * Defines the hotkeys to interact with the users management dashboard.
             */
            const defineComponentHotkeys = () => {
                if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    const hotkeys_context = new HotkeysContext();

                    hotkeys_context.register(["w", "s"], handleSideListsMovement, {
                        description: "<navigation> Move within side list up('w') and down('s').",
                    });

                    hotkeys_context.register(["a", "d"], handleSideListSwapping, {
                        description: "<navigation> Swap focus between the roles and users list.",
                    });

                    hotkeys_context.register(["e"], handleKeyboardSelection, {
                        description: "<navigation> Select the currently focused item.",
                    });

                    hotkeys_context.register(["q"], handleGoBack, {
                        description: `<${HOTKEYS_GENERAL_GROUP}> Leave the users dashboard and return to the previous page.`,
                    });

                    hotkeys_context.register(["?"], () => hotkeys_sheet_visible.set(!$hotkeys_sheet_visible), {
                        description: `<${HOTKEYS_GENERAL_GROUP}> Toggle the hotkeys cheat sheet.`,
                    });

                    global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
                }

                global_hotkeys_manager.loadContext(hotkeys_context_name);
            }

            /**
             * Handles the movement of the keyboard within the focused side list.
             * @param {KeyboardEvent} event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleSideListsMovement = (event, hotkey) => {
                let key_combo = hotkey.key_combo;
                const movement_increase = key_combo === "s" ? 1 : -1;

                let new_index = keyboard_movement_index + movement_increase;

                new_index = clampMovementIndex(new_index);

                keyboard_movement_index = new_index;
            }

            /**
             * Handles changing from the users list to the roles list and vice versa.
             */
            const handleSideListSwapping = () => {
                role_mode_enabled.set(!$role_mode_enabled);
            }

            /**
             * Handles the selection of the currently focused item.
             */
            const handleKeyboardSelection = () => {
                if ($role_mode_enabled) {
                    selectFocusedRole();
                } else {
                    selectFocusedUser();
                }
            }

            /**
             * Returns to the previous page. If it cannot return to the previous page without leaving the webapp, it will redirect to '/'
             * @return {void}
             */
            const handleGoBack = () => {
                let can_go_back = window.navigation?.canGoBack ?? false // Navigation API doesn't work in Safari and Firefox(to the surprise of no one).

                if (can_go_back) {
                    history.back();
                } else {
                    goto("/");
                }
            }

            /**
             * Selects the focused user.
             */
            const selectFocusedUser = async () => {
                if ($role_mode_enabled) return;
                
                refreshSelectedUserAccount();
            }

            /**
             * Selects the focused role.
             */
            const selectFocusedRole = () => {
                if (!$role_mode_enabled) return;

                let focused_role = $all_roles[keyboard_movement_index];

                console.log("Focused role: ", focused_role);

                selected_role_label = focused_role;
            }
        
        /*=====  End of Keybinds  ======*/

        /**
         * Loads all the initial data, should only be called from onMount.
         * @return {void}
         */
        const defineInitialData = async () => {
            await refreshAllUsersList();
            await refreshAllRolesList();
            await refreshAllGrantsList();
        }

        /**
         * Clamps the given index to the range of the currently focused side list.
         * @param {number} unsafe_index
         * @return {number}
         */
        const clampMovementIndex = (unsafe_index) => {
            let list_length = $role_mode_enabled ? $all_roles.length : $all_users.length;
            let safe_index = Math.max(0, Math.min(unsafe_index, list_length - 1));

            return safe_index;
        }

        /**
         * Handle the user-created event from the UserCreation component.
         * @param {CustomEvent<UserCreatedEvent>} 
         * @typedef {Object} UserCreatedEvent
         * @property {string} username 
         */
        const handleUserCreated = async (event) => {
            let new_user = event.detail.username;
            console.log("User created: ", new_user);

            await refreshAllUsersList();
        }

        /**
         * Handles the platform-role-created event emitted by the PlatformRolesCreator component.
         */
        const handlePlatformRoleCreated = async () => {
            refreshAllRolesList();
        }

        /**
         * Handles the role-deleted event emitted by the PlatformRolesEditor component.
         */
        const handleRoleDeleted = async () => {
            selected_role_label = "";

            await refreshAllRolesList();

            clampMovementIndex(keyboard_movement_index);
        }

        /**
         * Handles the user-data-changed event emitted by the PlatformUsersEditor component. Refreshes the selected_user_account.
         */
        const handleUserDataChanged = async () => {
            let focused_user = $all_users[keyboard_movement_index];
            let current_keyboard_movement_index = keyboard_movement_index;

            await refreshAllUsersList();

            if (focused_user.uuid !== $all_users[current_keyboard_movement_index].uuid) {
                let new_index = 0;

                for (let h = 0; h < $all_users.length; h++) {
                    if ($all_users[h].uuid === focused_user.uuid) {
                        new_index = h;
                        break;
                    }
                }

                keyboard_movement_index = new_index;
            }

            await refreshSelectedUserAccount();
        }

        /**
         * Handles the user-deleted event emitted by the PlatformUsersEditor component.
         */
        const handleUserDeleted = async () => {
            selected_user_account = null;

            await refreshAllUsersList();

            clampMovementIndex(keyboard_movement_index);
        }
        

        /**
         * Informes an unauthorized user that they don't have the necessary grants to access the page and their attempt 
         * will be reported.
         * @return {void}
         */ 
        const informUnauthorizedUser = async () => {
            let user_choice =  await confirmPlatformMessage({
                message_title: "You shall not pass!",
                question_message: "You don't have the necessary privileges to access this page. This incident will be reported. If you enter anyway, you will be banned from the platform.",
                cancel_label: "Go back",
                confirm_label: "Enter anyway",
                danger_level: 2,
                auto_focus_cancel: true,
            }); 

            let user_is_stupid = user_choice === 1;

            if (user_is_stupid) {
                // This is a measure to differentiate between malicious(and stupid) users and those who honestly made a mistake.
                banCurrentUser();
                alert("You idiot...");
            } else {
                goto("/"); // banCurrentUser() will logout the user which will redirect them to the login page(which they will not be able to use).
            }
        }

        /**
         * Fetches the user's data by a given user entry
         * @param {import("@models/Users").UserEntry} user_entry
         * @returns {Promise<import("@models/Users").UserAccount}
         */
        const fetchUserData = async (user_entry) => {
            let user_account = await getUserAccountFromUserEntry(user_entry);

            return user_account;
        }


        /**
         * Loads in all the existing user accounts in the platform and stores them in the 'all_users' property.
         * @return {void}
         */
        const refreshAllUsersList = async () => {
            let fresh_all_users = await getAllUsers();

            fresh_all_users = fresh_all_users.filter(user => user.username !== $current_user_identity.Username);

            if (fresh_all_users.length > 0) {
                all_users.set(fresh_all_users);
            }
        }

        /**
         * Loads in all the existing roles in the platform and stores them in the 'all_roles' property.
         * @return {Promise<void>}
         */
        const refreshAllRolesList = async () => {
            let fresh_all_roles = await getAllRoleLabels();

            fresh_all_roles = fresh_all_roles.filter(role => role !== "super_admin");

            if (fresh_all_roles.length > 0) {
                all_roles.set(fresh_all_roles);
            }
        }

        /**
         * Gets all the existing grants from the server and stores them in the all_grants store.
         * @returns {Promise<void>}
         */
        const refreshAllGrantsList = async () => {
            let new_grants = await getAllGrants();

            if (new_grants) {
                all_grants.set(new_grants);
            }
        }

        /**
         * Refreshes the selected_user_account.
         */
        const refreshSelectedUserAccount = async () => {
            let focused_user = $all_users[keyboard_movement_index];

            if (focused_user == null) return;

            selected_user_account = await fetchUserData(focused_user);
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<main id="dungeon-user-dashboard">
    {#if super_admin_confirmed}
        <aside id="dud-existing-users-wrapper" 
            class="side-list"
            class:side-list-focused={!$role_mode_enabled}
        >
            <h2 class="dud-list-type-label">
                Users
            </h2>
            <ul id="dud-existing-users-container"
                class="side-list-container"
            >
                {#each $all_users as user, h}
                    {@const is_keyboard_focused = keyboard_movement_index === h && !$role_mode_enabled}
                    <li class="dud-user-entry" 
                        class:is-keyboard-focused={is_keyboard_focused}
                    >
                        <span class="dud-euc-us-label">
                            {user.username}
                        </span>
                    </li>
                {/each}
            </ul>
        </aside>
        <article id="dud-operations-section">
            <section id="upper-operations-section"
                class="dud-operations-section dungeon-scroll"
            >
                {#if !$role_mode_enabled}                        
                    <UserCreation
                        on:user-created={handleUserCreated}
                        check_initial_setup={false}
                    />
                {/if}
                <PlatformRolesCreator 
                    on:platform-role-created={handlePlatformRoleCreated}
                />
            </section>
            <section id="lower-operations-section"
                class="dud-operations-section dungeon-scroll"
            >
                <PlatformRolesEditor
                    subject_role_label={selected_role_label}
                    on:role-deleted={handleRoleDeleted}
                />
                <PlatformUsersEditor
                    on:user-data-changed={handleUserDataChanged}
                    on:user-deleted={handleUserDeleted}
                    subject_user_account={selected_user_account}
                />
            </section>
        </article>
        <aside id="existing-roles-wrapper" 
            class="side-list"
            class:side-list-focused={$role_mode_enabled}
        >
            <h2 class="dud-list-type-label">
                Roles
            </h2>
            <ul id="existing-roles-container"
                class="side-list-container"
            >
                {#each $all_roles as role, h}
                    {@const is_keyboard_focused = keyboard_movement_index === h && $role_mode_enabled}
                    <li class="dud-erc-role-entry"
                        class:is-keyboard-focused={is_keyboard_focused}
                    >
                        <span class="dud-erc-re-label">
                            {role}
                        </span>
                    </li>     
                {/each}
        </aside>
    {/if}
</main>

<style>
    #dungeon-user-dashboard {
        --dark-divisor-border: 1px solid var(--grey-6);

        display: flex;
        flex-direction: row;
        width: 100%;
        container-type: size;
        height: calc(calc(100dvh - var(--navbar-height)) - 3px);
    }

    .side-list {
        width: min(30cqw, 300px);
        height: 100cqh;
        color: var(--main);
        
        & > h2.dud-list-type-label {
            background: hsl(from var(--grey-9) h s calc(l * 1.3));
            font-family: var(--font-decorative);
            font-size: calc(var(--font-size-h3) * 0.8);
            padding: var(--spacing-2) 0;
            text-align: center;
            line-height: 1;
        }

        & > ul {
            padding: var(--spacing-1);
        }

        & > ul > li {
            display: grid;
            place-items: center;
            padding: var(--spacing-2) 0;
        }

        & > ul > li:not(:last-child) {
            border-bottom: var(--dark-divisor-border);
        }
    }

    .side-list.side-list-focused {
        & h2.dud-list-type-label {
            color: var(--accent);
        }    
    }

    .side-list.side-list-focused ul.side-list-container > li.is-keyboard-focused {
        color: var(--accent);

        &:not(:last-child) {
            border-bottom: 1px solid var(--accent-3);
        }
    }
    
    /*=============================================
    =            Operations section            =
    =============================================*/

        .dud-operations-section {
            padding: var(--spacing-2) calc(var(--spacing-3) * 1.2); 
        }

    
        #dud-operations-section {
            flex-grow: 2;
            flex-basis: 40%;
            height: 100cqh;
            border-left: var(--dark-divisor-border);
            border-right: var(--dark-divisor-border);
        }

        #upper-operations-section {
            height: 50cqh;
            border-bottom: var(--dark-divisor-border);
        }

        #lower-operations-section {
            height: 50cqh;
        }
    
    
    /*=====  End of Operations section  ======*/

</style>