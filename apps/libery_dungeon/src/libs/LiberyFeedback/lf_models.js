import { emitLabeledError } from "./lf_utils";


/*=============================================
=            Errors            =
=============================================*/

    /**
     * @name toString
     * @function
     * @returns {string}
     */

    /**
     * @typedef {Object} ErrorContextInterface
     * @property {toString} toString
     */

    export class LabeledError extends Error {
        /**
         * The initial error context passed to the constructor.
         * @type {ErrorContextInterface}
         */
        #error_context;

        /**
         * A label used to identify the error. is recommended to use a const defined string to robustly compare it to other labels.
         * @type {string}
         */
        #label; 

        /**
         * Additional context to be added to the error message.
         * @type {ErrorContextInterface[]}
         */
        #additional_context;


        /**
         * A strongly contextualized error intended to provide a strong developer and user feedback.
         * @param {ErrorContextInterface} error_context 
         * @param {string} message - user intended message  
         * @param {string} label - a const defined string used to identify the error
         * @param  {...any} params 
         */
        constructor(error_context, message, label, ...params) {
            super(message, ...params);

            if (error_context?.toString == null || typeof error_context.toString !== "function") {
                throw new Error("error_context must have a toString method");
            }

            if (Error.captureStackTrace) {
                Error.captureStackTrace(this, LabeledError);
            }

            this.#additional_context = [];
            this.#error_context = error_context;
            this.#label = label;
        }

        /**
         * Emits the error if the environment is a window. Else it just logs it to the console.
         */
        alert() {
            if (globalThis.addEventListener == null) {
                console.error(this.toDevString());
                return;
            }

            emitLabeledError(this);
        }

        /**
         * Appends additional context to the error.
         * @param {ErrorContextInterface} context 
         */
        appendContext(context) {
            if (context?.toString == null || typeof context.toString !== "function") {
                throw new Error("context must have a toString method");
            }

            this.#additional_context.push(context);
        }

        /**
         * Prepends context to the error
         * @param {ErrorContextInterface} context
         */
        prependContext(context) {
            this.#additional_context = [context, ...this.#additional_context];
        }
            

        /**
         * Returns a message intended for the user.
         * @returns {string}
         */
        toHumanString() {
            return this.message;
        }

        /**
         * Returns a message intended for the developer.
         * @returns {string}
         */
        toDevString() {
            let dev_message = `Error: ${this.#label}\n`;
            dev_message += this.#error_context.toString() + " -> ";
            for (let h = 0; h < this.#additional_context.length; h++) {
                dev_message += this.#additional_context[h].toString();

                if (h < this.#additional_context.length - 1) {
                    dev_message += " -> ";
                }
            }
            dev_message += "\n"

            if (this.cause != null) {
                dev_message += `Cause: ${this.cause}\n`;
            }


            if (this.stack != null) {
                dev_message += `\n${this.stack}`;
            }

            return dev_message;
        }
    }

    export class VariableEnvironmentContextError {
        /**
         * A map of variable names to their values.
         * @type {Map<string, any>}
         */
        #variables;

        /**
         * The context error identifier.
         */
        #context_identifier;

        /**
         * @param {ErrorContextInterface} context_identifier 
         */
        constructor(context_identifier) {
            this.#variables = new Map();
            this.#context_identifier = context_identifier.toString();
        }

        /**
         * Adds a variable to the context.
         * @param {string} name 
         * @param {any} value 
         */
        addVariable(name, value) {
            this.#variables.set(name, value);
        }

        /**
         * Removes a variable from the context.
         * @param {string} name 
         */
        removeVariable(name) {
            this.#variables.delete(name);
        }

        /**
         * Returns a string representation of the context.
         * @returns {string}
         */
        toString() {
            let context_string = `Context: ${this.#context_identifier}\n`;

            let variable_count = this.#variables.size;
            let variables_read = 0;

            for (let [name, value] of this.#variables) {
                variables_read++;
                context_string += `${name}: ${value}`;

                if (variables_read < variable_count) {
                    context_string += ", ";
                }
            }

            context_string += "\n";

            return context_string;
        }

    }

/*=====  End of Errors  ======*/

/*=============================================
=            Messages            =
=============================================*/

    /**
     * @typedef {Object} ConfirmMessageParams
     * @property {string} question_message - the message that explains the decision to the user.
     * @property {string} message_title
     * @property {string} [confirm_label]
     * @property {string} [cancel_label]
     * @property {boolean} [auto_focus_cancel]
     * @property {number} [danger_level] - 0 is informational, 1 is warning, 2 is danger. and -1 is success.
     * @property {number} [timeout]
    */

    export class ConfirmMessage {
        /**
         * The message to be confirmed.
         * @type {string}
         */
        #question_message;

        /**
         * The message title.
         * @type {string}
         */
        #message_title;

        /**
         * The label of the confirm message.
         * @type {string}
         * @default "Accept"
         */
        #confirm_label;

        /**
         * The label of the cancel message. 
         * @type {string}
         * @default "Cancel"
         */
        #cancel_label;

        /**
         * Whether the cancel button should be auto focused.
         * @type {boolean}
         * @default true
         */
        #auto_focus_cancel;

        /**
         * The danger level of the message. 0 is informational, 1 is warning, 2 is danger. and -1 is success.
         * @type {number}
         */
        #danger_level;

        /**
         * The timeout of the message. -1 is infinite and no timeout should be set.
         * @type {number}
         * @default -1
         */
        #timeout;


        /**
         * @param {ConfirmMessageParams} param0
         */
        constructor({ question_message, message_title, confirm_label = "Accept", cancel_label = "Cancel", auto_focus_cancel = true, danger_level = 0, timeout = -1 }) {
            this.#question_message = question_message;
            this.#message_title = message_title;
            this.#confirm_label = confirm_label;
            this.#cancel_label = cancel_label;
            this.#auto_focus_cancel = auto_focus_cancel;
            this.#danger_level = danger_level;
            this.#timeout = timeout;
        }

        get QuestionMessage() {
            return this.#question_message;
        }

        get MessageTitle() {
            return this.#message_title;
        }

        get ConfirmLabel() {
            return this.#confirm_label;
        }
        
        get CancelLabel() {
            return this.#cancel_label;
        }

        get AutoFocusCancel() {
            return this.#auto_focus_cancel;
        }

        get DangerLevel() {
            return this.#danger_level;
        }

        get Timeout() {
            return this.#timeout;
        }
    }

    /**
     * The different types of user responses to a confirm message.
     * @enum {number}
     */
    export const confirm_question_responses = {
        CANCEL: 0,
        CONFIRM: 1,
        NO_CHOICE: -1
    }

/*=====  End of Messages  ======*/


