import { HotkeyFragment } from "./hotkeys_matchers";
import { HOTKEYS_HIDDEN_GROUP, HOTKEYS_GENERAL_GROUP } from "./hotkeys_consts";
/**
* @typedef {Object} HotkeyRegisterOptions
 * @property {boolean} bind - If true the hotkey will be binded immediately. default is false
 * @property {?string} description - The hotkey's description   
 * @property {"keypress"|"keydown"|"keyup"} mode - The mode of the keypress event. Default is "keydown"
*/

/**
* @typedef {Object} HotkeyDataParams
 * @property {string} key_combo
 * @property {function} handler 
 * @property {"keydown"|"keyup"} mode
 * @property {?string} description
*/

/**
* The default hotkey register options
 * @type {HotkeyRegisterOptions}
 */
export const default_hotkey_register_options = {
    bind: false,
    description: null,
    mode: "keydown"
}

/**
 * The hotkey callback type
 * @callback HotkeyCallback
 * @param {KeyboardEvent} event
 * @param {HotkeyData} hotkey
 */

export class HotkeyData {
    /**
     * @type {string} the key's name e.g: 'a', 'esc', '-', etc
     */
    #key_combo

    /**
     * The key combo fragments
     * @type {HotkeyFragment[]}
     */
    #key_combo_fragments

    /**
     * @type {HotkeyCallback} the callback to be called when the key is pressed
     */
    #callback

    /**
     * @type {"keypress"|"keydown"|"keyup"} the mode of the keypress event
     */
    #mode

    /**
     * @type {string} the hotkey's description
     * @default "<General>No information available"
     */
    #description

    /**
     * Whether the passed key combo is valid or not.
     * @type {boolean}
     */
    #is_valid;

    /**
     * Whether the hotkey can be prepended by a Vim motion. This is, it has a numeric metakey as the first fragment.
     * A hotkey with vim motion can have some sort of numeric metadata associated with it, which can be useful for some commands like 'move nth times to the right' for example.
     * @type {boolean}
     */
    #has_vim_motion;

    /**
     * Whether the hotkey callback is currently running.
     * @type {boolean}
     */
    #hotkey_execution_mutex;

    /**
     * @param {string} name the key's name e.g: 'a', 'esc', '-', etc
     * @param {function} callback the callback to be called when the key is pressed
     * @param {HotkeyRegisterOptions} options
     * @constructor
     */
    constructor(name, callback, options) {
        /** @type {string} the key's name e.g: 'a', 'esc', '-', etc */
        this.#key_combo = name
        /** @type {function} the callback to be called when the key is pressed */   
        this.#callback = callback

        
        this.#unpackOptions(options)
        
        this.#is_valid = true;
        this.#has_vim_motion = false;
        this.#hotkey_execution_mutex = false;

        this.#splitFragments()
    }

    /**
     * The clean hotkey's description string
     * @returns {string}
     */
    get Description() {
        return this.#description.replace(/<.*>/, "");
    }

    /**
     * The hotkey's 'face'. meaning the first hotkey identity on the hotkey. E.g: 'ctrl+k shift+space' => 'k'
     */
    get Face() {
        let first_fragment = this.#key_combo_fragments[0];

        return first_fragment.Identity;
    }

    /**
     * The hotkey's group
     * @returns {string}
     */
    get Group() {
        let hotkey_group = HOTKEYS_GENERAL_GROUP;
        
        /** @type {RegExpMatchArray} */
        let group_matches = this.#description.match(/<(.*)>/);

        if (group_matches != null && group_matches[1] !== undefined) {
            hotkey_group = group_matches[1];
        }

        return hotkey_group;
    }

    /**
     * Whether the hotkey is a sequence or not.
     * @returns {boolean}
     */
    get IsSequence() {
        return this.#key_combo_fragments.length > 1;
    }

    /**
     * The length of the hotkey. that is, how many fragments it has.
     * @returns {number}
     */
    get Length() {
        return this.#key_combo_fragments.length;
    }

    /**
     * The hotkey's mode
     * @type {"keydown"|"keyup"}
     */
    get Mode() {
        return this.#mode;
    }

    /**
     * The hotkey's 'Many Faces'. meaning all the hotkey identities on the first fragment of the hotkey. 
     * E.g: 'a+x' => ['a', 'x'], 'shift+3' => ['3', '#']
     * And yes, the name is a reference to Game of Thrones.
     * @returns {string[]}
     */
    get ManyFaces() {
        let first_fragment = this.#key_combo_fragments[0];

        return first_fragment.Identities;
    }

    /**
     * The many trails of the hotkey. Is the same as ManyFaces but for the last fragment of the hotkey.
     * @returns {string[]}
     * @see ManyFaces
     */
    get ManyTrails() {
        let last_fragment = this.#key_combo_fragments[this.#key_combo_fragments.length - 1];

        return last_fragment.Identities;
    }

    /**
     * Matches a sequence of Keyboard events with the hotkey. the sequence of keyboard events most be of the same length as the length of the hotkey.
     * @param {KeyboardEvent[]} events
     * @returns {boolean}
     */
    match(events) {
        if (events.length !== this.Length) return false;
        console.log(`Matching against: ${this.#key_combo}`, events);

        let matches = true;

        for (let h = 0; h < this.Length && matches; h++) {
            let fragment = this.#key_combo_fragments[h];
            let event = events[h];

            matches = fragment.match(event);            
        }

        return matches;
    }

    /**
     * The hotkey's Trigger. a hotkey trigger is the key that caught in an event, should cause a verification to see if the hotkey has been entered.
     * In the case of single fragment hotkeys, this is the first(and only) hotkey identity in a hotkey so for "a" -> 'a', "ctrl+shift+a" -> 'a', etc. In the case of a sequence (multi-fragment) hotkey, 
     * this is the last fragment identity in the sequence. So e.g: "ctrl+shift+a shift+s" -> 's'.
     * @type {string}
     */
    get Trigger() {
        let trigger = this.Face;

        if (this.IsSequence) {
            trigger = this.Trail;
        }

        return trigger;
    }

    /**
     * Returns all the triggers of the hotkey. same as Trigger but instead of returning the main identity of 
     * either the first or last fragment, it returns all the identities of the corresponding fragment.
     * @type {string[]}
     */
    get Triggers() {
        let triggers = this.ManyFaces;

        if (this.IsSequence) {
            triggers = this.ManyTrails;
        }

        return triggers;
    }


    /**
     * The hotkey's combo. 
     * NOTE: I wrote it on snake_case to avoid breaking the existing user space of this library. But this is not a correct nomenclature. don't repeat this.
     * @returns {string}
     */
    get key_combo() {
        return this.#key_combo;
    }

    /**
     * Registers the hotkey
     * TODO: Remove this method. HotkeyBinder will manage the hotkey binding directly taking HotkeyData as a parameter.
     */
    key_bind() {
        if (this.mode === "") {
            // Mousetrap.bind(this.name, this.callback)
            return
        }
        // Mousetrap.bind(this.name, this.callback, this.mode)
    }

    /**
     * Runs the hotkey's callback
     * @param {KeyboardEvent} event
     */
    run(event) {
        if (this.#hotkey_execution_mutex) {
            console.error(`Hotkey ${this.#key_combo} is already running. Skipping execution.`);
            return;
        }

        this.#hotkey_execution_mutex = true;
        this.#callback(event, this);
        this.#hotkey_execution_mutex = false;
    }

    /**
     * Splits the key combo into fragments. Fragments are split by ' '(aka \s). thats why the space key is represented as 'space' not ' '.
     * E.g: 'ctrl+k space' => [HotkeyFragment('ctrl+k'), HotkeyFragment('space')]
     * @returns {void}
     */
    #splitFragments() {
        let fragments = this.#key_combo.split(" ")

        try {
            this.#key_combo_fragments = fragments.map((fragment) => new HotkeyFragment(fragment)); // If the fragment parsing finds invalid members, this will panic.
        } catch (error) {
            console.error(`Error parsing hotkey: ${this.#key_combo}. Error: ${error}`);
            this.#is_valid = false;
        }

        if (this.#is_valid && this.#key_combo_fragments.length === 0) {
            this.#is_valid = false;
            return;
        }

        if (this.#key_combo_fragments[0].NumericMetakey) {
            this.#has_vim_motion = true;
        }
    }
    
    /**
     * The 'trail' of the hotkey. The trail is the hotkey's last identity. E.g: 'ctrl+k shift+space' => 'space'
     * This is useful to detect key sequences. On if the hotkey is a sequence then HotkeyBinder could(and probably will) register the Trail of the hotkey as the trigger and use the PastKeyDowns/PastKeyUps to check if the entire sequence was
     * pressed in a certain time frame. 
     * 
     * TODO: If this doesn't correctly reflect how we end up implementing the key sequences, then remove the comment. if it does, remove this line.
     * And if you find this line a year+ from now, go to the mirror, and laugh at yourself :P.
     * @returns {string}
     */
    get Trail() {
        let last_fragment = this.#key_combo_fragments[this.#key_combo_fragments.length - 1];

        return last_fragment.Identity;
    }

    /**
     * Unpacks the hotkey options into the hotkey data respective properties
     * @param {HotkeyRegisterOptions} options
     */
    #unpackOptions(options) {
        this.#mode = options.mode
        this.#description = options.description

        if (this.#mode === "keypress") {
            this.#mode = "keydown" // Keypress is deprecated and will(someday) be removed according to MDN
        }
    }

    /**
     * Whether the hotkey passed was valid and correctly parsed.
     * @returns {boolean}
     */
    get Valid() {
        return this.#is_valid;
    }

    /**
     * Whether the hotkey has a vim motion or not.
     * @returns {boolean}
     */
    get WithVimMotion() {
        return this.#has_vim_motion;
    }
}