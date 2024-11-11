<script>
    import { createEventDispatcher } from "svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The secret used to authenticate the user on initial setup mode. The user defines this secret on the settings.json file of the users service.
         * With this secret, and only when the system is on initial setup mode(means it has no registered users), the user can authenticate and create a new user 
         * which will be automatically set as a super admin. If this process is attempted while not in initial setup mode, user creation will fail
         * even with a valid user access token.
         * @type {string}
         */ 
        let initial_setup_secret;

        /**
         * The textarea element where the user will input the secret.
         * @type {HTMLTextAreaElement}
         */
        let the_initial_setup_secret_textarea;

        /**
         * Whether the secret form is ready to be submitted or not.
         * @type {boolean}
         */
        let secret_form_ready = false;

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the onkeypress event on the secret input field. On enter, it will emit the secret to the parent component.
         * @param {KeyboardEvent} event - The keypress event.
         */ 
        const handleSecretInputKeypress = event => {
            secret_form_ready = validateSecretTextareaValue();
            console.log("Secret form ready: ", secret_form_ready);

            if (event.key === "Enter" && secret_form_ready) {
                emitSecret(initial_setup_secret);
            }
        }

        /**
         * Handles the change event on the secret input field.
         * @param {Event} event - The change event.
         */
        const handleSecretInputChange = event => {
            secret_form_ready = validateSecretTextareaValue();
        }

        /**
         * Handles the click event on the submit button. It will emit the secret to the parent component.
         * @param {MouseEvent} event - The click event.
         */
        const handleInitialSetupSecretSelected = event => {
            if (!secret_form_ready) return;

            emitSecret(initial_setup_secret);
        }

        /**
         * Emits the secret to the parent component if secret_form_ready is true.
         * @param {string} secret_value
         */
        const emitSecret = secret_value => {
            if (!secret_form_ready) return;

            dispatch("secret-selected", {
                secret: secret_value
            });
        }

        /**
         * Validates the secret input field.
         * @returns {boolean}
         */
        const validateSecretTextareaValue = () => {
            return the_initial_setup_secret_textarea.checkValidity();
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<div id="initial-setup-secret-form-wrapper">
    <form action="none" id="ssf-form">
        <header>
            <h1 id="ssf-form-title">
                Initial setup
            </h1>
            <p id="ssf-form-instructions">
                It appears that the system is running for the first time. To create the first user of your system, please enter the secret you wrote on the settings.json file of the users service operational data folder(by default on <code>/var/www/libery_dungeon/server/users_od</code>). if you forgot the secret or have not changed the default value, fear not, you can go right ahead and change it, restart the users service and come back here. for further instructions, please refer to <a href="https://github.com/Gerardo115pp/libery_dungeon#initial-setup" target="_blank" rel="noopener noreferrer">the initial setup docs</a>
            </p>
        </header>
        <label id="ssf-form-secret-field">
            <span>
                Your super secret key
            </span>
            <textarea id="ssf-form-secret-input"
                class="dungeon-input"
                bind:this={the_initial_setup_secret_textarea}
                bind:value={initial_setup_secret}
                on:keypress={handleSecretInputKeypress}
                on:change={handleSecretInputChange}
                minlength="8"
                spellcheck="false"
                autocomplete="off"
                autofocus
                required
            />
        </label>
        <button id="ssf-form-submit-button"
            disabled={!secret_form_ready}
            class="dungeon-button-1" 
            on:click|preventDefault={handleInitialSetupSecretSelected}
        >
            This is my secret
        </button>
    </form>
</div>

<style>
    #initial-setup-secret-form-wrapper form#ssf-form {
        display: flex;    
        flex-direction: column;
        row-gap: var(--spacing-3);
    }
    
    /*=============================================
    =            Header            =
    =============================================*/
    
        form#ssf-form > header:first-of-type {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-2);

            & h1#ssf-form-title {
                text-align: center;
                font-size: var(--font-size-h3);
                color: var(--main);
            }
        }

        p#ssf-form-instructions {
            font-size: var(--font-size-1);
            color: var(--grey-3);

            & a {
                color: var(--accent-3);
                text-decoration: underline;
                font-weight: 600;
            }

            & code {
                color: var(--grey-2);
                background: var(--grey-9);
                padding: 0.35ch 0.55ch;
            }
        }
    
    
    /*=====  End of HEader  ======*/
    
    
    /*=============================================
    =            Secret key field            =
    =============================================*/
    
        label#ssf-form-secret-field {
            display: flex;
            flex-direction: column;
            row-gap: var(--spacing-1);

            & span {
                font-size: var(--font-size-1);
                font-weight: bolder;
                color: var(--grey-1);
            }
        }

        textarea#ssf-form-secret-input {
            resize: none;
            outline: none;
            font-size: var(--font-size-fineprint);
            height: 10em;

            &:user-invalid {
                border-color: var(--danger);
            }

            &:valid {
                border-color: var(--success);
            }
        }
    
    /*=====  End of Secret key field  ======*/
    
    

</style>