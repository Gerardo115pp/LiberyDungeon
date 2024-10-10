<script>
    import { createInitialUser, createNewUser, isUsersInInitialSetupMode } from "@models/Users";
    import InitialUserSetup from "./sub-components/InitialUserSetup.svelte";
    import { createEventDispatcher, onMount } from "svelte";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    
    /*=============================================
    =            Properties            =
    =============================================*/
   
        /**
         * Whether to check for initial setup mode or not.
         * @type {boolean}
        */
        export let check_initial_setup = true;
        
        /*----------  State  ----------*/

            /**
             * The new user's username.
             * @type {string}
             */
            let new_user_username = "";

            /**
             * The new user's password.
             * @type {string}
             */
            let new_user_password = "";

            /**
             * The user creation form.
             * @type {HTMLFormElement}
             */
            let the_user_creation_form;

            /**
             * Whether or not all the user creation data is ready.
             * @type {boolean}
             */
            let user_creation_data_ready = false;
        
        /*----------  Initial setup  ----------*/
                
            /**
             * Whether the system(at least the users service) is on initial_setup mode.
             * @type {boolean}
             * @default undefined
             */ 
            let is_initial_setup;

            /**
             * The secret used to authenticate the user on initial setup mode. 
             * @type {string}
             */
            let initial_setup_secret;
    

        const dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/
   
    onMount(async () => {
        if (check_initial_setup) {
            is_initial_setup = await isUsersInInitialSetupMode();
        } else {
            is_initial_setup = false;
        }
    }) 
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Verifies that all the information needed to create a new user is ready. if
         * initial_setup mode is active, this will include the initial_setup_secret.
         * @returns {boolean}
         */
        const checkUserCreationDataReady = () => {
            let data_ready = false;

            if (the_user_creation_form != null) {
                data_ready = the_user_creation_form.checkValidity();

                console.log("Form validity: ", data_ready);
            }

            if (data_ready && is_initial_setup) {
                data_ready = initial_setup_secret != null

                console.log("Secret validity: ", data_ready);
            }

            return data_ready;
        }
    
        /**
         * Handles the secret-selected event from the InitialUserSetup component.
         * @param {CustomEvent<InitialUserSetupSecretSelectedEvent>} event
         * @typedef {Object} InitialUserSetupSecretSelectedEvent
         * @property {string} secret 
         */ 
        const handleInitialUserSetupSecretSelected = (event) => {
            console.log("Secret selected: ", event.detail.secret);

            initial_setup_secret = event.detail.secret;
        }

        /**
         * Handles the keypress event on the username input field. It will check if the user creation data is ready.
         * @param {KeyboardEvent} event
         */
        const handleUsernameInputKeypress = event => {
            user_creation_data_ready = checkUserCreationDataReady();
        }

        /**
         * Handles the keypress event on the password input field. It will check if the user creation data is ready.
         * @param {KeyboardEvent} event
         */
        const handlePasswordInputKeypress = event => {
            user_creation_data_ready = checkUserCreationDataReady();
        }

        /**
         * Handles the click event on the submit button. 
         * @param {MouseEvent} event
         */
        const handleUserCreationFormSubmit = async event => {
            event.preventDefault();

            if (!user_creation_data_ready) return;

            let created = await registerNewUser();

            if (created) {
                dispatch("user-created", {
                    username: new_user_username
                });
            } else {
                dispatch("user-creation-failed");
            }

            the_user_creation_form.reset();
        }

        /**
         * Handles the click event on the reset initial setup secret button.
         * @param {MouseEvent} event
         */
        const handleResetInitialSetupSecret = event => {
            initial_setup_secret = undefined;
        }

        /**
         * Registers a new user either with createInitialUser or createUser, depending on the value of is_initial_setup.
         * Returns a promise that to true if the user was created successfully and false otherwise.
         * @returns {Promise<boolean>}
         */
        const registerNewUser = async () => {
            let user_created = false;

            if (is_initial_setup) {
                console.log("Creating initial user...");
                user_created = await createInitialUser(new_user_username, new_user_password, initial_setup_secret);
            } else {
                console.log("Creating user...");
                user_created = await createNewUser(new_user_username, new_user_password);
            }

            if (!user_created) {
                let variable_environment = new VariableEnvironmentContextError("In Users/UsersCreation/UserCreation.registerNewUser");

                variable_environment.addVariable("new_user_username", new_user_username);
                variable_environment.addVariable("new_user_password", new_user_password);
                variable_environment.addVariable("is_initial_setup", is_initial_setup);
                variable_environment.addVariable("user_creation_data_ready", user_creation_data_ready);

                let labeled_err = new LabeledError(variable_environment, "Failed to create user", lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();
            }

            return user_created
        }
    
    /*=====  End of Methods  ======*/
    
</script>

{#if is_initial_setup != null}
    {#if is_initial_setup && initial_setup_secret == null}
        <InitialUserSetup
            on:secret-selected={handleInitialUserSetupSecretSelected}
        />
    {:else}
        <div id="user-creation-wrapper">
            <header id="ucw-header">
                <h1 id="ucw-header-title">
                    {#if is_initial_setup}
                        You'r first user
                    {:else}
                        Add a user
                    {/if}
                </h1>
                <p id="ucw-instructions">
                    {#if is_initial_setup}
                        Type a strong password for the first user account as it will have complete authority over the system. the password strength is not enforced in anyway. This is your system, and you'r free to do with it as you please but you will be the only one responsible for it. With great power comes great responsibility.
                    {:else}
                        New users will be automatically added to the visitor role. by default(that is unless you change it), this role can view the content in cluster that are not private and that's it, they cannot do anything else. You can change this role's permissions on the user management dashboard available to roles with the 'modify_users' or 'ALL_PRIVILEGES' grants(e.g super admins).
                    {/if}
                </p>
            </header>
            <form id="ucw-creation-form"
                action="none"
                bind:this={the_user_creation_form} 
            >
                <label class="dungeon-input">
                    <span>
                        Username
                    </span>
                    <p class="ucw-cf-tooltip field-tooltip">
                        The username of the new user. It must be unique. Valid characters are ascii letters(a-z, A-Z), numbers, and the symbols: <code> '_', '-', '.', '@'</code> 
                    </p>
                    <input id="ucw-username-input"
                        type="text"
                        bind:value={new_user_username}
                        on:keypress={handleUsernameInputKeypress}
                        spellcheck="false"
                        minlength="2"
                        maxlength="64"
                        pattern="{'[a-zA-Z0-9_\\-\\.@]{2,}'}"
                        required
                    >
                </label>
                <label class="dungeon-input">
                    <span>
                        Password
                    </span>
                    <p class="ucw-cf-tooltip field-tooltip">
                        Is strongly recommended to use a strong password but this is not enforced. It has to be at least 5 and less than 255 characters long. Aside from that, do as you like. your system, your rules, your risk.
                    </p>
                    <input 
                        type="text"
                        bind:value={new_user_password}
                        on:keypress={handlePasswordInputKeypress}
                        minlength="5"
                        maxlength="255"
                        spellcheck="false"
                        required
                    >
                </label>
                <fieldset id="ucw-creation-form-actions">
                    {#if is_initial_setup && initial_setup_secret != null}
                        <button id="ucw-reset-initial-setup-secret-btn" 
                            class="dungeon-button-2"
                            on:click|preventDefault={handleResetInitialSetupSecret}
                        >
                            Change my secret
                        </button>
                    {/if}                    
                    <button id="ucw-submit-button"
                        class="dungeon-button-1"
                        disabled={!user_creation_data_ready}
                        type="submit"
                        on:click|preventDefault={handleUserCreationFormSubmit}
                    >
                        Create user
                    </button>                        
                </fieldset>
            </form>
        </div>
    {/if}
{/if}

<style>
    #user-creation-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);
    }
    
    /*=============================================
    =            Header            =
    =============================================*/
    
        header#ucw-header {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-2);

            & h1#ucw-header-title {
                text-align: center;
                font-size: var(--font-size-h3);
                color: var(--main);
            }
        } 

        p#ucw-instructions {
            font-size: var(--font-size-1);
            color: var(--grey-3);
        }
    
    /*=====  End of Header  ======*/
    
    /*=============================================
    =            Form            =
    =============================================*/
    
        form#ucw-creation-form {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-3);
        }

        form#ucw-creation-form label.dungeon-input {
            & span::after {
                content: ":";
            }

            & input {
                color: var(--main);
            }
        }

        fieldset#ucw-creation-form-actions {
            display: flex;
            justify-content: space-between;
        }
    
    /*=====  End of Form  ======*/
    
</style>