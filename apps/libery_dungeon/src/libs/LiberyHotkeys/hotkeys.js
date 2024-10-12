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
     * @type {function} the callback to be called when the key is pressed
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
     * The hotkey's mode
     * @type {"keydown"|"keyup"}
     */
    get Mode() {
        return this.#mode;
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
}