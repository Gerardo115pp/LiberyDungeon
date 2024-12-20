/*=============================================
=            KEY MAPS            =
=============================================*/

    import { DEFAULT_CAPTURE_ACCEPT_TERMINATOR, DEFAULT_CAPTURE_CANCEL_TERMINATOR } from "./hotkeys_consts";

    /*----------  Key values for keyboard events  ----------*/
    // See: https://developer.mozilla.org/en-US/docs/Web/API/UI_Events/Keyboard_event_key_values

    const letter_keys = new Set([
        "a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
        "k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
        "u", "v", "w", "x", "y", "z"
    ]);

    const upper_letter_hotkeys = new Set([
        "A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
        "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
        "U", "V", "W", "X", "Y", "Z"
    ]);

    const number_hotkeys = new Set([
        "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"
    ]);

    const non_letter_character_producing_keys = new Set([
        "!", "\"", "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", "-", ".", "/", ":", ";", "<", "=", ">", "?", "@", "[", "\\", "]", "^", "_", "`", "{", "|", "}", "~"
    ]);

    const whitespace_hotkeys = new Set([
        " ", "space", "Tab", "Enter"
    ]);

    const modifier_keys = new Set([
        "Alt", "AltGraph", "CapsLock", "Control", "Fn", 
        "FnLock", "Hyper", "Meta", "NumLock", "ScrollLock",
        "Shift", "Super", "Symbol", "SymbolLock"
    ]);

    const navigation_keys = new Set([
        "ArrowDown", "ArrowLeft", "ArrowRight", "ArrowUp",
        "End", "Home", "PageDown", "PageUp"
    ]);


    const editing_keys = new Set([
        "Backspace", "Clear", "Copy", "CrSel", "Cut", "Delete",
        "EraseEof", "ExSel", "Insert", "Paste", "Redo", "Undo"
    ]);

    const function_keys = new Set([
        "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10",
        "F11", "F12"
    ]);


    
    /*----------  Code values for keyboard events  ----------*/
        // The KeyboardEvent.code property represents a physical key on the keyboard (as opposed to the character generated by pressing the key). In other words, 
        // this property returns a value that isn't altered by keyboard layout or the state of the modifier keys
        // See: https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/code
    
        const letter_codes = new Set([
            "KeyA", "KeyB", "KeyC", "KeyD", "KeyE", "KeyF", "KeyG", "KeyH", "KeyI", "KeyJ",
            "KeyK", "KeyL", "KeyM", "KeyN", "KeyO", "KeyP", "KeyQ", "KeyR", "KeyS", "KeyT",
            "KeyU", "KeyV", "KeyW", "KeyX", "KeyY", "KeyZ"
        ]);

        const number_codes = new Set([
            "Digit0", "Digit1", "Digit2", "Digit3", "Digit4", "Digit5", "Digit6", "Digit7", "Digit8", "Digit9"
        ]);

        const whitespace_codes = new Set([
            "Space", "Tab", "Enter"
        ]);

    
    /*----------  Hotkey alias  ----------*/
        // These are the names we use to refer to specific keys and metakeys. these are the hotkey fragments that are expected on key combos. e.g: "ctrl+shift+a", "c c", "\d x", etc.

        const SHIFT_KEY = "shift";
        const CONTROL_KEY = "ctrl";
        const COMMAND_KEY = "cmd";
        const ALT_KEY = "alt";
        const OPTIONS_KEY = "opt";
        const ESCAPE_KEY = "esc";
        const ENTER_KEY = "enter";
        const SPACE_KEY = "space";
        const TAB_KEY = "tab";
        const BACKSPACE_KEY = "backspace";
        const DELETE_KEY = "del";
        const INSERT_KEY = "ins";
        const HOME_KEY = "home";
        const END_KEY = "end";
        const PAGE_UP_KEY = "pgup";
        const PAGE_DOWN_KEY = "pgdown";
        const ARROW_UP_KEY = "up";
        const ARROW_DOWN_KEY = "down";
        const ARROW_LEFT_KEY = "left";
        const ARROW_RIGHT_KEY = "right";
        const CAPS_LOCK_KEY = "caps";
        const NUM_LOCK_KEY = "num";
        const SCROLL_LOCK_KEY = "scroll";

        const F1_KEY = "f1";
        const F2_KEY = "f2";
        const F3_KEY = "f3";
        const F4_KEY = "f4";
        const F5_KEY = "f5";
        const F6_KEY = "f6";
        const F7_KEY = "f7";
        const F8_KEY = "f8";
        const F9_KEY = "f9";
        const F10_KEY = "f10";
        const F11_KEY = "f11";
        const F12_KEY = "f12";

        // Multi match keys
        const NUMERIC_METAKEY = "\\d";
        const LETTER_KEY = "\\l";
        const STRING_KEY = "\\s";
        const CHARACTER_KEY = "\\c";

        const valid_fragment_identities = new Set([
            ESCAPE_KEY, ENTER_KEY, SPACE_KEY, TAB_KEY, BACKSPACE_KEY, DELETE_KEY, INSERT_KEY,
            HOME_KEY, END_KEY, PAGE_UP_KEY, PAGE_DOWN_KEY, ARROW_UP_KEY, ARROW_DOWN_KEY, ARROW_LEFT_KEY,
            ARROW_RIGHT_KEY, CAPS_LOCK_KEY, NUM_LOCK_KEY, SCROLL_LOCK_KEY, NUMERIC_METAKEY, LETTER_KEY,
            F1_KEY, F2_KEY, F3_KEY, F4_KEY, F5_KEY, F6_KEY, F7_KEY, F8_KEY, F9_KEY, F10_KEY, F11_KEY, F12_KEY
        ]);

        const shift_mutations_us_keyboard = new Map([
            ["1", "!"], ["2", "@"], ["3", "#"], ["4", "$"], ["5", "%"], ["6", "^"], ["7", "&"], ["8", "*"], ["9", "("], ["0", ")"],
            ["`", "~"], ["-", "_"], ["=", "+"], ["[", "{"], ["]", "}"], ["\\", "|"], [";", ":"], ["'", "\""], [",", "<"], [".", ">"], ["/", "?"],
            ["a", "A"], ["b", "B"], ["c", "C"], ["d", "D"], ["e", "E"], ["f", "F"], ["g", "G"], ["h", "H"], ["i", "I"], ["j", "J"],
            ["k", "K"], ["l", "L"], ["m", "M"], ["n", "N"], ["o", "O"], ["p", "P"], ["q", "Q"], ["r", "R"], ["s", "S"], ["t", "T"],
            ["u", "U"], ["v", "V"], ["w", "W"], ["x", "X"], ["y", "Y"], ["z", "Z"]
        ]);

        const capture_metakeys = new Set([
            STRING_KEY,
            CHARACTER_KEY,
        ]);

/*=====  End of KEY MAPS  ======*/

/*=============================================
=            Matchers            =
=============================================*/

    /**
     * Returns true if the passed key is a valid hotkey identity. E.g: imagine a hotkey like "ctrl+shift+a", it's valid hotkey identity would be the "a" 
     * because it's not a modifier key and it is present in normal keyboard layouts(by normal i mean it's not exclusive of virtual keyboards).
     * other examples would be: "ctrl+[" -> "[", "ctrl+shift+1" -> "1", "b" -> "b", "ctrl+shift+alt+space" -> "space"
     * @param {string} key
     * @returns {boolean}
     */
    const isHotkeyIdentity = (key) => {
        let is_modifier = modifier_keys.has(key);

        if (is_modifier) return false;

        let key_in_uppercase = key.toUpperCase(); // Have normalized versions of the key for better matching
        let key_in_lowercase = key.toLowerCase();

        let is_letter = letter_keys.has(key_in_lowercase);

        let is_number = number_hotkeys.has(key);

        let is_special_character = non_letter_character_producing_keys.has(key);

        let is_other_valid_identity = valid_fragment_identities.has(key);

        return is_letter || is_number || is_special_character || is_other_valid_identity;
    }

    /**
     * Returns whether a passed identity is an alias.
     * @param {string} identity
     * @returns {boolean}
     */
    const isAliasIdentity = (identity) => {
        let is_alias = false;

        switch (identity) {
            case SPACE_KEY:
                is_alias = true;
                break;
            case ARROW_UP_KEY:
            case ARROW_DOWN_KEY:
            case ARROW_LEFT_KEY:
            case ARROW_RIGHT_KEY:
                is_alias = true;
                break;
        }

        return is_alias;
    }

    /**
     * Returns whether a passed key is a metakey. a metakey is a key that matches more then one key
     * for example \s(which would have to be written as "\\l" in a string) would match any letter key. and \d would match any number key.
     * @param {string} key
     * @returns {boolean}
     */
    const isMetaKey = (key) => {
        let is_meta = false;

        switch (key) {
            case NUMERIC_METAKEY:
            case LETTER_KEY:
                is_meta = true;
                break;
        }

        return is_meta;
    }

    /**
     * Whether the passed key produces a character.
     * @param {string} key
     * @returns {boolean}
     */
    const IsCharacterProducingKey = (key) => {
        let is_character_producing = IsLetter(key);

        if (!is_character_producing && number_hotkeys.has(key)) {
            is_character_producing = true;
        }

        if (!is_character_producing && non_letter_character_producing_keys.has(key)) {
            is_character_producing = true;
        }

        is_character_producing = is_character_producing || key === " ";

        return is_character_producing;        
    }

    /**
     * Parses a meta key into the keys it represents.
     * @param {string} key
     * @returns {string[]}
     */
    const parseMetaKey = (key) => {
        /**
         * @type {string[]}
         */
        let keys = [];

        switch (key) {
            case NUMERIC_METAKEY:
                keys = Array.from(number_hotkeys);
                break;
            case LETTER_KEY:
                keys = [
                    ...Array.from(letter_keys), 
                    ...Array.from(upper_letter_hotkeys)
                ];
                break;
        }

        return keys;
    }

    /**
     * Whether the passed key is a modifier key. The matching is case sensitive. Pass
     * the key as it appears in a KeyboardEvent.key property.
     * @param {string} key
     * @returns {boolean}
     */
    export const IsModifier = key => {
        return modifier_keys.has(key);
    }

    /**
     * Whether the passed key is a numeric key.
     * @param {string} key
     * @returns {boolean}
     */
    export const IsNumeric = key => {
        return number_hotkeys.has(key);
    }

    /**
     * Whether the passed key is a letter key.
     * @param {string} key
     * @returns {boolean}
     */
    export const IsLetter = key => {
        return letter_keys.has(key) || upper_letter_hotkeys.has(key);
    }

    /**
     * Whether a given key combo includes a capture metakey.
     * @param {string} key_combo
     * @returns {boolean}
     */
    export const IncludesCaptureMetakey = key_combo => {
        let includes_metakey = false;

        for (let capture_metakey of capture_metakeys[Symbol.iterator]()) {
            if (key_combo.includes(capture_metakey)) {
                includes_metakey = true;
            }
        }

        return includes_metakey;
    }

/*=====  End of Matchers  ======*/


/**
 * Transforms a identity from an alias to how it would be represented in a KeyboardEvent.key property.
 * @param {string} identity
 * @returns {string}
 */
const transformAliasIdentity = (identity) => {
    let transformed_identity = identity;

    switch (identity) {
    case SPACE_KEY:
        transformed_identity = " ";
        break;
    case ARROW_UP_KEY:
        transformed_identity = "ArrowUp";
        break;
    case ARROW_DOWN_KEY:
        transformed_identity = "ArrowDown";
        break;
    case ARROW_LEFT_KEY:
        transformed_identity = "ArrowLeft";
        break;
    case ARROW_RIGHT_KEY:
        transformed_identity = "ArrowRight";
        break;
    }

    return transformed_identity;
}

/**
* For reference different types of hotkeys would be:
 * "shift+x" is a key fragment.
 * "c c" is a key sequence. which has two key fragments.
 * "gZ" is a key chord. and is never spelled as a combo. For this you spell a combo 'g' and bind register it as a Chord passing it not a handler function, but a different hotkey context, which as soon as any key is pressed will 
 *      either trigger the chord if the next key is in the context or reset the previous context that contained the chord if it is not.
*/

/**
 * Represents a piece of a hotkey, by this i mean if there is a hotkey such as 'ctrl+k ctrl+t' this entire hotkey will be composed from two fragments: 'ctrl+k' and 'ctrl+t'
 * TODO: Right now fragments like 'a+x' are not supported. But this is planned to have support as soon as i can replace at least the features Mousetrap provides. and from there i will start adding more features.
 */
export class HotkeyFragment {
    /**
     * the hotkey fragment like "ctrl+shift+a"
     * @type {string} 
     */
    #fragment

    /**
     * The fragment members. if we use the same example "ctrl+shift+a" the members will be ["ctrl", "shift", "a"]
     * @type {string[]} 
     */
    #fragment_members

    /**
     * The fragment identity. this is the key that would be passed on a KeyboardEvent.key property. e.g: "a" in "ctrl+shift+a"
     */
    #fragment_identity

    /**
     * All the detectable identities this fragment(support for fragments like 'a+x') has. e.g: "ctrl+shift+a" has "a" but because of shift it also has "A". for a+x it would be "a" and "x". so for the first example you would have "a" on
     * #fragment_identity and this would be ["A"], and for the second example, #fragment_identity = "a" and this would be ["x"]. for "a+x+d" you would have #fragment_identity = "a" and this would be ["x", "d"]
     * @type {string[]}
     */
    #alternate_identities

    /**
     * Whether the key requires control modifier.
     * @type {boolean} 
     */
    #control_modifier

    /**
     * Whether the key requires shift modifier.
     * @type {boolean} 
     */
    #shift_modifier

    /**
     * Whether the key requires alt modifier.
     * @type {boolean} 
     */
    #alt_modifier

    /**
     * Whether the fragment identity is set explicitly in uppercase. This is useful
     * to prevent false positives. for example shift+a matching "A"
     * @type {boolean}
     */
    #uppercase_explicit

    /**
     * Whether the fragment represents is a numberic metakey.
     * @type {boolean}
     */
    #numeric_metakey

    /**
     * @param {string} fragment the hotkey fragment. e.g: "ctrl+shift+a", "a", "ctrl+a", 
     */
    constructor(fragment) {

        this.#fragment = fragment;
        this.#fragment_members = this.#splitFragment();

        this.#control_modifier = false;
        this.#shift_modifier = false;
        this.#alt_modifier = false;
        this.#uppercase_explicit = false;
        this.#numeric_metakey = false; // Overwritten by #parseIdentityMember

        this.#fragment_identity = "";
        this.#alternate_identities = [];

        this.#parseMembers();
    }

    /**
     * If shift required adds to the alternate identities the shifted versions of the fragment identity and current alternate identities.
     * This is only reliable for US keyboards. and will not work correctly for other keyboard layouts. getKeyboardLayout() can be used to get the current keyboard layout. but
     * It requires a SecureContext which would force users to either get a domain or use something like mkcert to install the platform or be limited to localhost which mostly defeats the purpose of LiberyDungeon.
     * See: https://developer.mozilla.org/en-US/docs/Web/API/Keyboard/getLayoutMap
     * Called from #parseModifiers only.
     */
    #addShiftMutations() {
        if (!this.ShiftRequired) return;

        let current_identities = [...this.#alternate_identities];

        let shifted_fragment_identity = shift_mutations_us_keyboard.get(this.#fragment_identity);

        if (shifted_fragment_identity != null) {
            this.#alternate_identities.push(shifted_fragment_identity);
        }

        for (let identity of current_identities) {
            let shifted_identity = shift_mutations_us_keyboard.get(identity);

            if (shifted_identity != null) {
                this.#alternate_identities.push(shifted_identity);
            }
        }
    }

    /**
     * Whether the fragment requires the alt modifier.
     * @type {boolean}
     */
    get AltRequired() {
        return this.#alt_modifier;
    }

    /**
     * Whether the fragment requires the control modifier.
     * @type {boolean}
     */
    get CtrlRequired() {
        return this.#control_modifier;
    }

    /**
     * Detects the role a member has in the fragment. 
     * if its a modifier it will set the corresponding modifier flag.
     * @param {string} member
     */
    #detectMemberRole(member) {
        let valid_member = false;
        let is_control = false;
        let is_shift = false;
        let is_alt = false;

        if (!this.#control_modifier) {
            this.#control_modifier = member === CONTROL_KEY || member === COMMAND_KEY;
            is_control = this.#control_modifier;
            valid_member = is_control;
        }

        if (!this.#shift_modifier && !valid_member) {
            this.#shift_modifier = member === SHIFT_KEY;
            is_shift = this.#shift_modifier;
            valid_member = is_shift;
        }

        if (!this.#alt_modifier && !valid_member) {
            this.#alt_modifier = this.#alt_modifier || member === ALT_KEY || member === OPTIONS_KEY;
            is_alt = this.#alt_modifier;
            valid_member = is_alt;
        }

        if (!is_control && !is_shift && !is_alt && !valid_member) {
            valid_member = this.#parseIdentityMember(member);
        }

        if (!valid_member) {
            throw new Error(`Invalid hotkey fragment<${this.#fragment}> has invalid member '${member}'`);
        }
    }

    /**
     * The fragment's identity. this is the key that would be passed on a KeyboardEvent.key property. e.g: "a" in "ctrl+shift+a"
     * @returns {string}
     */
    get Identity() {
        return this.#fragment_identity;
    }

    /**
     * All the fragment's identities, fragment_identity + alternate_identities
     * @returns {string[]}
     */
    get Identities() {
        return [this.#fragment_identity, ...this.#alternate_identities];
    }

    /**
     * Matches the hotkey fragment with a given KeyboardEvent
     * @param {KeyboardEvent} event
     * @returns {boolean}
     */
    match(event) {
        if (event == null) return false;
        // TODO: Add matched to the beginning of all but the first if statements so that if it has been already unmatched 
        // it doens't go through the rest of the checking process unnecessarily.

        let matched = true;

        if (event.altKey && !this.#alt_modifier) {
            matched = false;
        }

        if (event.ctrlKey && !this.#control_modifier) {
            matched = false;
        } // Shift cannot have the same treatment as it can also alter the key in more ways than just changing the case. e.g: shift+2 -> "@".
        // But checking that ctrl and alt are not pressed if not required is needed to avoid overwriting browser shortcuts unintentionally.
        
        if (!this.#matchIdentity(event)) {
            console.log(`Fragment: ${this.#fragment} did not match '${event.key}'`);
            matched = false;
        }

        if (matched && this.#control_modifier) {
            // Ensures an event that explicitly requires a control modifier, example 'ctrl+a' will not match 'a'.
            matched = event.ctrlKey;
        }

        if (matched && this.#shift_modifier) {
            // Ensures an event that explicitly requires a shift modifier, example 'shift+a' will not match 'a'.
            matched = event.shiftKey;
        }

        if (matched && this.#alt_modifier) {
            // Ensures an event that explicitly requires an alt modifier, example 'alt+a' will not match 'a'.
            matched = event.altKey;
        }

        return matched;
    }

    /**
     * The fragment's members
     * @type {string[]}
     */
    get Members() {
        return this.#fragment_members;
    }

    /**
     * Whether the fragment requires any modifier.
     * @returns {boolean}
     */
    get Modifier() {
        return this.#control_modifier || this.#shift_modifier || this.#alt_modifier;
    }

    /**
     * Matches the identity of the fragment against the passed KeyboardEvent. considering how modifiers would affect the match.
     * @param {KeyboardEvent} event
     */
    #matchIdentity(event) {
        let key = event.key;
        let is_letter = IsLetter(key);

        if (this.#shift_modifier || !is_letter) {
            key = key.toLowerCase();
        }

        let is_match = false;
        let all_identities = this.Identities;

        for (let h = 0; h < all_identities.length && !is_match; h++) {
            let identity_candidate = all_identities[h];

            if (!is_letter) {
                identity_candidate = identity_candidate.toLowerCase();
            }

            is_match = key === identity_candidate;
        }

        return is_match;
    }

    /**
     * Matches a numeric metakey against the passed KeyboardEvent.
     * @param {KeyboardEvent} event
     * @returns {boolean}
     */
    matchNumericMetakey(event) {
        if (!this.NumericMetakey) return false;

        let key = event.key;
        let is_number = IsNumeric(key);

        return is_number;
    }

    /**
     * Matches a vim motion from a sequence of events.
     * @param {import("@libs/utils").StackBuffer<KeyboardEvent>} event_history
     * @returns {number}
     */
    matchVimMotion(event_history) {
        if (!this.NumericMetakey || event_history.IsEmpty()) return NaN;

        let matched_motion = NaN;

        /**
         * numeric keys found in sequence.
         * @type {string[]}
         */
        let numeric_keys = [];

        let has_match = false;

        /**
         * @type {KeyboardEvent | null}
         */
        let event = event_history.PeekTraversing();

        while (event != null && IsNumeric(event.key)) {
            has_match = true;

            let key = event.key;

            numeric_keys.unshift(key);

            event = event_history.Traverse();
        }

        if (has_match) {
            matched_motion = parseInt(numeric_keys.join(""));
        }

        return matched_motion;
    }


    /**
     * Whether the fragment is a numeric metakey.
     * @type {boolean}
     */
    get NumericMetakey() {
        return this.#numeric_metakey;
    }

    /**
     * Parses a fragment identity member, if successful, returns true.
     * @param {string} member
     * @returns {boolean}
     */
    #parseIdentityMember(member) {
        let is_identity = isHotkeyIdentity(member);

        if (!is_identity) return false;

        let identity = member;

        if (isAliasIdentity(member)) {
            identity = transformAliasIdentity(member);
        }

        if (isMetaKey(identity)) {
            let meta_keys = parseMetaKey(identity);

            this.#numeric_metakey = identity === NUMERIC_METAKEY;

            identity = member;
            this.#alternate_identities = this.#alternate_identities.concat(meta_keys);
        }
        
        if (this.#fragment_identity === "") {
            this.#fragment_identity = identity;
            this.#uppercase_explicit = upper_letter_hotkeys.has(member);
        } else {
            this.#alternate_identities.push(member);
        }

        return is_identity;
    }

    #parseMembers() {
        for (const member of this.#fragment_members) {
            this.#detectMemberRole(member);
        }

        this.#addShiftMutations()
    }

    /**
     * Splits the fragment into its members
     * @returns {string[]} the fragment members
     */
    #splitFragment() {
        return this.#fragment.split("+");
    }

    /**
     * Whether the fragment requires the shift modifier.
     * @type {boolean}
     */
    get ShiftRequired() {
        return this.#shift_modifier;
    }

    /**
     * Whether the fragment's identity is set explicitly in uppercase.
     * @type {boolean}
     */
    get UppercaseExplicit() {
        return this.#uppercase_explicit;
    }
}

/**
 * Hotkey matcher for a capture hotkey type. This is a hotkey with up to three fragments. The first fragment is the capture initializer which has to be a specified key(not a metakey). The second fragment is the capture metakey
 * which can be '\\c' for one letter key. '\\s' for any amount of character producing keys. The third fragment is the capture terminator which also has to be a specified key. The terminator can be omitted only for '\\c'
 * as it is implied that the capture ends when a letter key is pressed. For '\\s' the terminator can be omitted but a default one is used. The default terminator is defined in the hotkeys_consts.js file.
 */
export class HotkeyCaptureMatcher {

    static CAPTURE_STATE_UNACTIVE = Symbol("CAPTURE_STATE_UNACTIVE");

    static CAPTURE_STATE_ACTIVE = Symbol("CAPTURE_STATE_ACTIVE");

    static CAPTURE_STATE_COMPLETE = Symbol("CAPTURE_STATE_COMPLETE");

    static CAPTURE_STATE_CANCELLED = Symbol("CAPTURE_STATE_CANCELLED");

    /**
     * The capture initializer fragment.
     * @type {HotkeyFragment}
     */
    #initializer_fragment

    /**
     * The captured string.
     * @type {string[]}
     */
    #captured_string

    /**
     * The hotkey fragment that ends the capture.
     * @type {HotkeyFragment}
     */
    #accept_terminator_fragment

    /**
     * The hotkey fragment that cancels the capture.
     * @type {HotkeyFragment}
     */
    #cancel_terminator_fragment

    /**
     * The max capture length.
     * @type {number}
     * @default -1
     */
    #max_capture_length

    /**
     * The capture state.
     * @type {Symbol}
     */
    #capture_state

    /**
     * @param {string} key_combo
     */
    constructor(key_combo) {
        const combo_fragments = key_combo.split(" ");
        this.#captured_string = [];
        this.#capture_state = HotkeyCaptureMatcher.CAPTURE_STATE_UNACTIVE;
        this.#max_capture_length = -1;


        if (combo_fragments.length < 2 || combo_fragments.length > 3) {
            throw new Error(`Invalid capture hotkey combo: ${key_combo}`);
        }

        this.#initializer_fragment = new HotkeyFragment(combo_fragments[0]);

        if (combo_fragments[1] !== CHARACTER_KEY && combo_fragments[1] !== STRING_KEY) {
            throw new Error(`Invalid capture hotkey combo: '${key_combo}'. Unknown capture metakey: '${combo_fragments[1]}'`);
        }

        if (combo_fragments[1] === CHARACTER_KEY) {
            this.#max_capture_length = 1;
        }

        if (combo_fragments.length === 3) {
            this.#accept_terminator_fragment = new HotkeyFragment(combo_fragments[2]);
        } else {
            this.#accept_terminator_fragment = new HotkeyFragment(DEFAULT_CAPTURE_ACCEPT_TERMINATOR);
        }

        this.#cancel_terminator_fragment = new HotkeyFragment(DEFAULT_CAPTURE_CANCEL_TERMINATOR);
    }

    /**
     * The fragment that ends the capture with a successful state.
     * @type {HotkeyFragment}
     */
    get AcceptTerminatorFragment() {
        return this.#accept_terminator_fragment;
    }

    /**
     * The captured string.
     * @type {string}
     */
    get CapturedString() {
        let captured_string = "";

        if (this.#capture_state === HotkeyCaptureMatcher.CAPTURE_STATE_COMPLETE) {
            captured_string = this.#captured_string.join("");
        }

        return captured_string;
    }

    /**
     * Incomplete capture string. Call this when the capture is still to get the current string.
     * @type {string}
     */
    get IncompleteCapturedString() {
        return this.#captured_string.join("");
    }

    /**
     * The fragment that ends the capture with a cancelled state.
     * @type {HotkeyFragment}
     */
    get CancelTerminatorFragment() {
        return this.#cancel_terminator_fragment;
    }

    /**
     * Captures the key of a given KeyboardEvent. Panics if the capture is not active. Returns whether the capture has ended.
     * @param {KeyboardEvent} event
     */
    capture(event) {
        if (this.#capture_state !== HotkeyCaptureMatcher.CAPTURE_STATE_ACTIVE) {
            throw new Error("In LiberyHotkeys/hotkeys_matchers.js HotkeyCaptureMatcher.capture: Attempted to capture a key stroke when the capture has not been triggered.");
        }

        let capture_accepted = this.#accept_terminator_fragment.match(event);
        let capture_cancelled = !capture_accepted && this.#cancel_terminator_fragment.match(event);

        if (!capture_accepted && !capture_cancelled) {

            if (event.key === "Backspace") {
                this.#captured_string.pop();
            } else if (IsCharacterProducingKey(event.key)) {
                this.#captured_string.push(event.key);
            }

            if (this.#max_capture_length > 0 && this.#captured_string.length >= this.#max_capture_length) {
                capture_accepted = true;
                capture_cancelled = false;
            }
        }

        this.#capture_state = capture_accepted ? HotkeyCaptureMatcher.CAPTURE_STATE_COMPLETE : this.#capture_state;

        this.#capture_state = capture_cancelled ? HotkeyCaptureMatcher.CAPTURE_STATE_CANCELLED : this.#capture_state;

        return capture_accepted || capture_cancelled;
    }

    /**
     * The fragment that triggers the capture.
     * @type {HotkeyFragment}
     */
    get InitializerFragment() {
        return this.#initializer_fragment;
    }

    /**
     * Resets the capture state.
     * @returns {void}
     */
    reset() {
        this.#capture_state = HotkeyCaptureMatcher.CAPTURE_STATE_UNACTIVE;
        this.#captured_string = [];
    }

    /**
     * The state of the capture.
     * @type {Symbol}
     */
    get State() {
        return this.#capture_state;
    }

    /**
     * Matches agains the HotkeyCapture trigger. returns whether the capture has started. if so the state is implicitly activated.
     * @param {KeyboardEvent} event
     * @returns {boolean} 
     */
    tryTrigger(event) {
        let triggered = this.#initializer_fragment.match(event);
        console.log(`Capture trigger: ${triggered}`);

        if (triggered) {
            console.log("Capture triggered");
            this.#capture_state = HotkeyCaptureMatcher.CAPTURE_STATE_ACTIVE;
        }

        return triggered;
    } 
}