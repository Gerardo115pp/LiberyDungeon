<script>
    import { createEventDispatcher, onMount } from "svelte";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { goto } from "$app/navigation";
    import { loginPlatformUser } from "@models/Users";
    import { access_state_confirmed, current_user_identity, has_user_access, setUserLoggedIn } from "@stores/user";
    
    /*=============================================
    =            Properties            =
    =============================================*/
   
        /**
         * the username to attempt the login with
         * @type {string}
         */
        let alleged_username = "";

        /**
         * The password to attempt the login with.
         * @type {string}
         */
        let alleged_password = "";

        /**
         * The user login form.
         * @type {HTMLFormElement}
         */
        let the_user_login_form;

        /**
         * Whether or not all the user login data is ready.
         * @type {boolean}
         */
        let user_login_data_ready = false;

        const dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/
   
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Verifies that all the information needed to create a new user is ready. if
         * initial_setup mode is active, this will include the initial_setup_secret.
         * @returns {boolean}
         */
        const checkUserLoginDataReady = () => {
            let data_ready = the_user_login_form.checkValidity();

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
         * Handles the keypress event on the username input field. It will check if the user login data is ready.
         * @param {KeyboardEvent} event
         */
        const handleUsernameInputKeypress = event => {
            user_login_data_ready = checkUserLoginDataReady();
        }

        /**
         * Handles the keypress event on the password input field. It will check if the user login data is ready.
         * @param {KeyboardEvent} event
         */
        const handlePasswordInputKeypress = event => {
            user_login_data_ready = checkUserLoginDataReady();
        }

        /**
         * Handles the click event on the submit button. 
         * @param {MouseEvent} event
         */
        const handleUserLoginFormSubmit = async event => {
            event.preventDefault();

            if (!user_login_data_ready) return;

            let created = login();
        }

        /**
         * Registers a new user either with createInitialUser or createUser, depending on the value of is_initial_setup.
         * Returns a promise that to true if the user was created successfully and false otherwise.
         * @returns {Promise<boolean>}
         */
        const login = async () => {
            let fresh_user_identity = await loginPlatformUser(alleged_username, alleged_password);

            if (fresh_user_identity == null) {
                let error_context = new VariableEnvironmentContextError("In LoginPage.login");

                error_context.addVariable("alleged_username", alleged_username);
                error_context.addVariable("alleged_password", alleged_password);

                let err = new LabeledError(error_context, "There seems to be a problem with the service. try again later", lf_errors.ERR_PROCESSING_ERROR);
                err.alert();

                return false;
            }

            setUserLoggedIn(fresh_user_identity); 

            goto("/");
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<main id="user-login-page" class="form-page-wrapper">
    <div id="user-login-wrapper" class="form-wrapper">
        <header id="ulw-header">
            <h1 id="ulw-header-title">
                The doors of Durin, lord of Moria.<br/><span class="sub-headline">Speak, friend, and enter.</span>
            </h1>
            <p id="ulw-instructions">
                The user credential(that is username and password) are given by the system administrator(And no, they are not "friend"... or are they??). Users can only be created by administrators.
            </p>
        </header>
        <form id="ulw-login-form"
            action="none"
            bind:this={the_user_login_form} 
        >
            <label class="dungeon-input">
                <span>
                    Your username
                </span>
                <input id="ulw-username-input"
                    type="text"
                    bind:value={alleged_username}
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
                    Your secret
                </span>
                <input 
                    type="password"
                    bind:value={alleged_password}
                    on:keypress={handlePasswordInputKeypress}
                    minlength="5"
                    maxlength="255"
                    spellcheck="false"
                    required
                >
            </label>
            <fieldset id="ulw-login-form-actions">
                <button id="ulw-submit-button"
                    class="dungeon-button-1"
                    disabled={!user_login_data_ready}
                    on:click|preventDefault={handleUserLoginFormSubmit}
                >
                    login
                </button>                        
            </fieldset>
        </form>
    </div>
</main>

<style>

    #user-login-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-3);
        padding: var(--spacing-4) 0;
    }
    
    /*=============================================
    =            Header            =
    =============================================*/
    
        header#ulw-header {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-2);

            & h1#ulw-header-title {
                text-align: center;
                font-size: var(--font-size-h3);
                line-height: 1.1;
                color: var(--main);
            }
        } 

        p#ulw-instructions {
            font-size: var(--font-size-1);
            color: var(--grey-3);
        }
    
    /*=====  End of Header  ======*/
    
    /*=============================================
    =            Form            =
    =============================================*/
    
        form#ulw-login-form {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-3);
        }

        form#ulw-login-form label.dungeon-input {
            & span::after {
                content: ":";
            }

            & input {
                color: var(--main);
            }
        }

        fieldset#ulw-login-form-actions {
            display: flex;
            justify-content: space-between;
            flex-direction: row-reverse;
        }
    
    /*=====  End of Form  ======*/
    
</style>