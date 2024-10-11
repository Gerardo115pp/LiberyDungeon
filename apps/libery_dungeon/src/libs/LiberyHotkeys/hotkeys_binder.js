import { StackBuffer } from "@libs/utils";
import { 
    MAX_PAST_EVENTS, 
    MAX_TIME_BETWEEN_SEQUENCE_KEYSTROKES
} from "./hotkeys_consts";


export class HotkeysBinder {

    /**
     * A stack with the last MAX_PAST_EVENTS keydown KeyboardEvents
     * @type {StackBuffer<KeyboardEvent>}
     */
    #keyboard_past_keydowns;

    /**
     * A stack with the last MAX_PAST_EVENTS keyup KeyboardEvents
     * @type {StackBuffer<KeyboardEvent>}
     */
    #keyboard_past_keyups;

    /**
     * keydown handler function
     * @type {function}
     */
    #bound_handleKeyDown;

    /**
     * keyup handler function
     * @type {function}
     */
    #bound_handleKeyUp;

    constructor() {
        if (globalThis.addEventListener === undefined) {
            throw new Error("This environment does not support event listeners")
        }

        this.#keyboard_past_keydowns = new StackBuffer(MAX_PAST_EVENTS);
        this.#keyboard_past_keyups = new StackBuffer(MAX_PAST_EVENTS);
        
        this.#bound_handleKeyDown = this.#handleKeyDown.bind(this);
        this.#bound_handleKeyUp = this.#handleKeyUp.bind(this);

        this.#setup()
    }

    /**
     * Cleans up the binder
     */
    Destroy() {
        globalThis.removeEventListener("keydown", this.#handleKeyDown);
        globalThis.removeEventListener("keyup", this.#handleKeyUp);
    }

    /**
     * Handles the keydown event
     * @param {KeyboardEvent} event
     */
    #handleKeyDown(event) {
        console.log("keydown", event) 

        this.#keyboard_past_keydowns.Add(event);
    }

    /**
     * Handles the keyup event
     * @param {KeyboardEvent} event
     */
    #handleKeyUp(event) {
        console.log("keyup", event)

        this.#keyboard_past_keyups.Add(event);
    }

    /**
     * The past keydown events stack.
     * @returns {StackBuffer<KeyboardEvent>}
     */
    get PastKeyDowns() {
        return this.#keyboard_past_keydowns;
    }

    /**
     * The past keyup events stack.
     * @returns {StackBuffer<KeyboardEvent>}
     */
    get PastKeyUps() {
        return this.#keyboard_past_keyups;
    }

    /**
     * Attaches the binder to the global this object
     */
    #setup() {
        globalThis.addEventListener("keydown", this.#bound_handleKeyDown);
        globalThis.addEventListener("keyup", this.#bound_handleKeyUp);
    }
}

let GlobalHotkeysBinder = null;

export const setupHotkeysBinder = () => {
    if (GlobalHotkeysBinder != null) {
        throw new Error("Hotkeys binder already exists. Call destroyHotkeysBinder before setting up a new one")
    }

    GlobalHotkeysBinder = new HotkeysBinder();
}

export const destroyHotkeysBinder = () => {
    GlobalHotkeysBinder.Destroy();

    GlobalHotkeysBinder = null;
}

/**
 * Returns the global hotkeys binder.
 * @returns {HotkeysBinder | null}
 */
export const getHotkeysBinder = () => {
    return GlobalHotkeysBinder;
}