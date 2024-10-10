/*=============================================
=            Hotkeys context events            =
=============================================*/

/**
 * Hotkeys context events
 * @readonly
 * @enum {string}
 */
export const hotkeys_context_events = {
    CONTEXT_CHANGED: 'context_changed',
    CONTEXT_ACTIVATED: 'context_activated', // Not used 
    CONTEXT_DEACTIVATED: 'context_deactivated', // Not used
    CONTEXT_DECLARED: 'context_declared', // Not used
    CONTEXT_REMOVED: 'context_removed', // Not used
    HOTKEYS_ACTIVATED: 'hotkeys_activated', // Not used
    HOTKEYS_DEACTIVATED: 'hotkeys_deactivated', // Not used
    HOTKEYS_DECLARED: 'hotkeys_declared', // Not used
    HOTKEYS_REMOVED: 'hotkeys_removed', // Not used
    HOTKEYS_RESET: 'hotkeys_reset', // Not used
    HOTKEYS_TRIGGERED: 'hotkeys_triggered', // Not used
}

/**
 * @typedef {Object} HotkeysContextEventDetail
 * @property {?string} context_name
 * @property {?string} hotkey_name
 */

/**
 * @callback HotkeysContextEventCallback
 * @param {CustomEvent<HotkeysContextEventDetail>} event
 */

/**
 * Suscribes to hotkeys context events and returns a function to unsubscriber similar to svelte.createEventDispatcher. Useful to avoid using the wrong event target.
 * @param {hotkeys_context_events} event
 * @param {HotkeysContextEventCallback} callback
 * @returns {Function} Call this function to unsuscribe the event
 */
export const suscribeHotkeysContextEvents = (event, callback) => {
    // IMPORTANT: Use the window event target, for only this object is guaranteed to be available.
    const event_unsubscriber = () => {
        window.removeEventListener(event, callback);
    }

    window.addEventListener(event, callback);

    return event_unsubscriber;
}

/**
 * Dispatches a hotkeys context event
 * @param {hotkeys_context_events} event
 * @param {HotkeysContextEventDetail} detail
 * @returns {void}
 */
export const dispatchHotkeysContextEvent = (event, detail) => {
    window.dispatchEvent(new CustomEvent(event, {detail}));
}

/*=====  End of Hotkeys context events   ======*/