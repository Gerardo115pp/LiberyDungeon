
/**
 * A null handler for hotkeys.
 * @type {import('./hotkeys').HotkeyCallback}
 */
export const HOTKEY_NULLISH_HANDLER = (event, hotkey) => console.error("Called a nullish hotkey handler");

/**
 * A null handler for hotkey capture events.
* @type {import('./hotkeys').HotkeyCaptureCallback}
*/
export const HOTKEY_NULLISH_CAPTURE_HANDLER = (event, captured_string) => {};



/*=============================================
=            Defaults            =
=============================================*/

    /**
     * The default keyboard event mode for hotkeys that don't specify one.
     * @type {"keydown" | "keyup"}
     */
    export const DEFAULT_KEYBOARD_EVENT_MODE = "keydown";

    /**
     * The default await execution value for hotkeys that don't specify one
     * @type {boolean}
     */
    export const DEFAULT_AWAIT_EXECUTION = false;

    /**
     * The default hotkey mode for hotkeys that don't specify one
     * @type {"keydown" | "keyup"}
     */
    export const DEFAULT_HOTKEY_MODE = "keydown";

    /**
     * A nullish hotkey description. All instances of this description should be replaced.
     * @type {string}
     */
    export const HOTKEY_NULL_DESCRIPTION = "REPLACE ME";

    /**
     * The default value for consider_time_in_sequence for hotkeys that don't specify one
     * @type {boolean}
     */
    export const DEFAULT_CONSIDER_TIME_IN_SEQUENCE = true;

    /**
     * The default value for can_repeat for hotkeys that don't specify one
     * @type {boolean}
     */
    export const DEFAULT_CAN_REPEAT = true;


/*=====  End of Defaults  ======*/




/*=============================================
=            Hotkeys information            =
=============================================*/

    export const HOTKEYS_HIDDEN_GROUP = "hidden";
    export const HOTKEYS_GENERAL_GROUP = "general";

/*=====  End of Hotkeys information  ======*/


/*=============================================
=            Hotkey Matching Config            =
=============================================*/

    export const MAX_PAST_EVENTS = 100;
    export const MAX_PAST_HOTKEYS_TRIGGERED = 10;
    export const MAX_TIME_BETWEEN_SEQUENCE_KEYSTROKES = 1200; // Milliseconds
    export const MIN_TIME_BETWEEN_HOTKEY_REPEATS = 20; // Milliseconds
    // Whether to check all hotkeys with the same trigger and pick the longest one(true) or just the first one that matches(false)
    export const DISABLE_KEYPRESS_EVENTS = true;
    
    export const HOTKEY_SPECIFICITY_PRECEDENCE = true; 
    export const HOTKEY_CAPTURE_BASE_SPECIFICITY = 10;


/*=====  End of Hotkey Matching Config  ======*/


/*=============================================
=            Hotkey capture            =
=============================================*/

    /**
     * The default accept terminator for capture hotkeys that don't specify one.
     * Typing the accept terminator will mark the captured string as complete.
     * @type {string}
     */
    export const DEFAULT_CAPTURE_ACCEPT_TERMINATOR = "enter";

    /**
     * The default cancel terminator for capture hotkeys that don't specify one.
     * Typing the cancel terminator will cancel the capture.
     * @type {string}
     */
    export const DEFAULT_CAPTURE_CANCEL_TERMINATOR = "esc";

/*=====  End of Hotkey capture  ======*/

