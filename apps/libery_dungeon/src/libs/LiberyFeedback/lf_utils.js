import { writable } from 'svelte/store';
import { ConfirmMessage } from './lf_models';

const PLATFORM_ERROR_EVENT = "platform-error"
const PLATFORM_MESSAGE_EVENT = "platform-message"

/**
 * Emit a labeled error if executed in a window environment and returns true. Otherwise just passes it to console.error and returns false.
 * @param {import('./lf_models').LabeledError} labeled_error
 * @returns {boolean} - true if the error was emitted
 */
export const emitLabeledError = labeled_error => {
    if (globalThis.addEventListener == null) {
        console.error(labeled_error.toDevString());
        return false;
    }

    const event = new CustomEvent(PLATFORM_ERROR_EVENT, {
        detail: labeled_error
    });

    globalThis.dispatchEvent(event);

    return true;
}

/**
 * Subscribe to platform errors.
 * @param {LabeledErrorHandler} callback
 * @returns {() => void} - unsubscribe function
 * @callback LabeledErrorHandler
 * @param {import('./lf_models').LabeledError} labeled_error
 */
export const subscribeToPlatformErrors = callback => {
    if (globalThis.addEventListener == null) {
        throw new Error("subscribeToPlatformErrors must be called in a window environment");
    }

    const handler = (/** @type {{ detail: import("./lf_models").LabeledError; }} */ event) => {
        if (event.detail instanceof Error) {
            callback(event.detail);
        }
    }

    const unsubscribe = () => {
        // @ts-ignore
        globalThis.removeEventListener(PLATFORM_ERROR_EVENT, handler);
    }

    // @ts-ignore
    globalThis.addEventListener(PLATFORM_ERROR_EVENT, handler);

    return unsubscribe;
}

/**
 * @typedef {CustomEvent<string>} PlatformMessageEvent
 */

/**
 * Emit a message if executed in a window environment.
 * @param {string} message
 */ 
export const emitPlatformMessage = message => {
    if (globalThis.addEventListener == null) {
        console.log(message);
        return;
    }

    /**
     * @type {PlatformMessageEvent}
     */
    const event = new CustomEvent(PLATFORM_MESSAGE_EVENT, {
        detail: message
    });

    globalThis.dispatchEvent(event);
}

/**
 * Subscribe to platform messages.
 * @param {PlatformMessageHandler} callback
 * @returns {() => void} - unsubscribe function
 * @callback PlatformMessageHandler
 * @param {string} message
 */
export const subscribeToPlatformMessages = callback => {
    if (globalThis.addEventListener == null) {
        throw new Error("subscribeToPlatformMessages must be called in a window environment");
    }

    /**
     * @type {EventListener}
     * @param {PlatformMessageEvent} event 
     */
    // @ts-ignore
    const handler = event => {
        if (typeof event.detail === "string") {
            callback(event.detail);
        }
    }

    const unsubscribe = () => {
        globalThis.removeEventListener(PLATFORM_MESSAGE_EVENT, handler);
    }

    globalThis.addEventListener(PLATFORM_MESSAGE_EVENT, handler);

    return unsubscribe;
}

/*=============================================
=            Confirm Messages            =
=============================================*/

    /**
     * The message that needs to be confirmed.
     * @type {import('svelte/store').Writable<import('./lf_models').ConfirmMessage | null>}
     * @default null    
     */
    export const confirm_message = writable(null);

    /**
     * the confirmation response. 0 must be interpreted as a 'cancelation' choice and 1 as a 'confirmation' choice. The store will be null(no talking about the value, the reference)
     * unless there is a confirm dialog component handling the confirmation message, in which case that component is expected to set the value to a readable store. 
     * 0 -> cancelation, 
     * 1 -> confirmation,
     * -1 -> no choice made.
     * 
     * @type {import('svelte/store').Readable<number>}
     * @default null
     */
    // @ts-ignore
    export let confirm_response = null;

    /**
     * Sets the confirm response store.
     * @param {import('svelte/store').Readable<number>} response_store
     */
    export const setConfirmResponse = response_store => {
        confirm_response = response_store;
    }

    /**
     * sets the confirm message to be displayed.
     * @param {import('./lf_models').ConfirmMessageParams} message_params
     * @returns {void}
     */
    const setConfirmMessage = (message_params) => {
        let new_confirm_message = new ConfirmMessage(message_params)

        confirm_message.set(new_confirm_message);   
    }

    /**
     * Emits a confirm message with the given parameters. returns a promise that will resolve to 1 for
     * a confirmation and 0 for a cancelation. if an explicit user choice was not made, the promise will resolve to -1.
     * @description 0 === cancelation,
     * @description 1 === confirmation,
     * @description -1 === no choice made.
     * @param {import('./lf_models').ConfirmMessageParams} message_params
     * @returns {Promise<number>}
     */
    export const confirmPlatformMessage = message_params => {
        if (confirm_response === null) {
            console.warn("confirm response store is null. The confirmation message will not be displayed.");
            return Promise.resolve(-1);
        }

        // @ts-ignore
        return new Promise((resolve, reject) => {

            let response_unsubscriber = () => {
                console.log("response unsubscriber called, but holds the wrong function value");
            }

            response_unsubscriber = confirm_response.subscribe(new_response => {
                if (new_response < -1 || new_response > 1) return;

                response_unsubscriber();
                resolve(new_response);
            });

            setConfirmMessage(message_params);
        });
    }

/*=====  End of Confirm Messages  ======*/


/*=============================================
=            Discrete feedback messages            =
=============================================*/

    /**
     * The message been displayed in the discrete feedback component.
     * @type {import('svelte/store').Writable<string | null>}
     */
    export const discrete_feedback_message = writable("");

    /**
     * Sets the discrete feedback message to be displayed.
     * @param {string} message
     */
    export const setDiscreteFeedbackMessage = message => {
        discrete_feedback_message.set(message);
    }

    globalThis.setDiscreteFeedbackMessage = setDiscreteFeedbackMessage;

/*=====  End of Discrete feedback messages  ======*/