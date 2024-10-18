import { 
    HotkeyFragment, 
    IsNumeric,
    IsModifier
} from "./hotkeys_matchers";
import { HOTKEYS_HIDDEN_GROUP, HOTKEYS_GENERAL_GROUP } from "./hotkeys_consts";
import { 
    MAX_TIME_BETWEEN_SEQUENCE_KEYSTROKES,
    HOTKEY_SPECIFICITY_PRECEDENCE
} from "./hotkeys_consts";

/**
* @typedef {Object} HotkeyRegisterOptions
 * @property {boolean} bind - If true the hotkey will be binded immediately. default is false
 * @property {?string} description - The hotkey's description   
 * @property {"keypress"|"keydown"|"keyup"} mode - The mode of the keypress event. Default is "keydown"
 * @property {boolean} await_execution - Whether the execution of a callback should end before another hotkey can be triggered. Default is true
 * @property {boolean} consider_time_in_sequence - Whether the hotkey sequence should expire if they are to far apart in time. Default is false
 * @property {boolean} can_repeat - Whether the hotkey should be triggered if the trigger is repeating(holding down the key). Default is false
*/

/**
* @typedef {Object} HotkeyDataParams
 * @property {string} key_combo
 * @property {function} handler 
 * @property {HotkeyRegisterOptions} options
*/

/**
* The default hotkey register options
 * @type {HotkeyRegisterOptions}
 */
export const default_hotkey_register_options = {
    bind: false,
    description: null,
    mode: "keydown",
    await_execution: true,
    consider_time_in_sequence: false,
    can_repeat: false,
}

/**
 * The hotkey callback type
 * @callback HotkeyCallback
 * @param {KeyboardEvent} event
 * @param {HotkeyData} hotkey
 */

class HotkeyMatch {
    /**
     * Any vim motion numeric data found. they are stored in the order they were found.
     * @type {number[]}
     */
    #motion_matches;

    /**
     * Whether the match is has been determined as successful.
     * @type {boolean}
     */
    #is_successful;

    constructor() {
        this.#motion_matches = [];
        this.#is_successful = false;
    }

    /**
     * Adds a numeric string as a vim motion match. Returns true if the string was correctly parsed and added.
     * @param {string} numeric_string
     * @returns {boolean}
     */
    addMotionMatch(numeric_string) {
        this.#panicIfSuccessful();

        let number = parseInt(numeric_string);

        if (isNaN(number)) {
            return false;
        }

        this.#motion_matches.push(number);

        return true;
    }

    /**
     * Adds a numeric string just as addMotionMatch but reverses the string before parsing it.
     * returns true if the string was correctly parsed and added.
     * @param {string} numeric_string
     * @returns {boolean}
     */ 
    addReversedMotionMatch(numeric_string) {
        this.#panicIfSuccessful();

        let reversed_string = numeric_string.split("").reverse().join("");

        return this.addMotionMatch(reversed_string);
    }

    /**
     * Any vim motion numeric data found. they are stored in the order they were found.
     * @type {number[]}
     */
    get MotionMatches() {
        return this.#motion_matches;
    }

    /**
     * Panics if called and the match.#is_successful is true.
     * @returns {void}
     */
    #panicIfSuccessful() {
        if (this.#is_successful) {
            throw new Error("Match is already successful. Refusing to add modifications.");
        }
    }

    /**
     * Reverses the motion matches array. When matched from history events, the last match of a hotkey will
     * be the first in the event history. This means that the fragments forcefully have to be matched in reverse order than 
     * they were written in the hotkey combo. This method reverses the matches so that the fit how the were written in the hotkey combo.
     * ALWAYS CALL IT AFTER A SUCCESSFUL MATCH.
     * @returns {void}
     */
    #reverseMotionMatches() {
        if (this.#is_successful) {
            throw new Error("Match is already successful. Refusing to reverse it.");
        }

        this.#motion_matches.reverse(); 
    }

    /**
     * Set the match as successful. Modifications after this will panic.
     * @returns {void}
     */
    setSuccessful() {
        if (this.#is_successful) {
            throw new Error("Match is already successful. Refusing to set it again.");
        }

        this.#reverseMotionMatches();        
        this.#is_successful = true;
    }

    /**
     * Whether the match was successful or not.
     * @returns {boolean}
     */
    get Successful() {
        return this.#is_successful;
    }
}

export class HotkeyData {
    /**
     * @type {string} the key's name e.g: 'a', 'esc', '-', etc
     */
    #key_combo

    /**
     * The key combo specificity. This is the precedence a matched hotkey has over another one if both match a sequence
     * @type {number}
     */
    #key_combo_specificity

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
     * Metadata produces after a successful match. Populated by `match` if the hotkey matches the event history. destroyed after one call to `run`.
     * @type {HotkeyMatch | null}
     */
    #match_metadata;

    /**
     * Whether the hotkey callback is currently running.
     * @type {boolean}
     */
    #hotkey_execution_mutex;

    /**
     * the hotkey's options
     * @type {HotkeyRegisterOptions}
     */
    #the_options;

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

        this.#the_options = {
            ...default_hotkey_register_options,
            ...options
        };

        
        this.#unpackOptions(options)
        
        this.#is_valid = true;
        this.#has_vim_motion = false;
        this.#match_metadata = null;
        this.#hotkey_execution_mutex = false;
        this.#key_combo_specificity = 0;

        this.#splitFragments()

        if (HOTKEY_SPECIFICITY_PRECEDENCE) {
            this.#calculateSpecificity()
        }
    }

    /**
     * Whether a caller to `run` should await the callback's execution
     * @returns {boolean}
     */
    get AwaitExecution() {
        return this.#the_options.await_execution;
    }

    /**
     * Creates and sets a new match metadata.
     * @returns {HotkeyMatch}
     */
    #createMatchMetadata() {
        this.#match_metadata = new HotkeyMatch();
    }

    /**
     * Checks if the given KeyboardEvent is expired against the given compare_time. 
     * automatically returns false if the hotkey is not time sensitive or if it's not a sequence.
     * also return false if the compare_time is a negative number.
     * @param {KeyboardEvent} event - this is the past we want to know if it's expired.
     * @param {number} compare_time - this is the more recent event time we want to compare the event against.
     * @returns {boolean}
     */
    #checkEventExpired(event, compare_time) {
        if (!this.ConsiderTimeInSequence || compare_time < 0 || !this.IsSequence) {
            return false;
        }

        const elapsed_time = compare_time - event.timeStamp;

        return elapsed_time > MAX_TIME_BETWEEN_SEQUENCE_KEYSTROKES;
    }

    /**
     * Calculate the specificity base on the number of keys in the hotkey. Each hotkey count 1
     * except for vim motions which count 5. This requires the fragments to be already known, so
     * call it after #splitFragments.
     * @returns {void}
     */
    #calculateSpecificity() {
        this.#key_combo_specificity = 0;

        for (let fragment of this.#key_combo_fragments) {
            let fragment_specificity = fragment.NumericMetakey ? 5 : 1;

            fragment_specificity += fragment.AltRequired ? 1 : 0;

            fragment_specificity += fragment.CtrlRequired ? 1 : 0;

            fragment_specificity += fragment.ShiftRequired ? 1 : 0;

            fragment_specificity += fragment.UppercaseExplicit ? 1 : 0;

            this.#key_combo_specificity += fragment_specificity;
        }
    }

    /**
     * Whether time between hotkey fragments should be considered in a sequence.
     * @returns {boolean}
     */
    get ConsiderTimeInSequence() {
        return this.#the_options.consider_time_in_sequence;
    }

    /**
     * Whether the hotkey should be triggered if the trigger is repeating(holding down the key).
     * @returns {boolean}
     */
    get CanRepeat() {
        return this.#the_options.can_repeat;
    }

    /**
     * Destroys the match metadata.
     * @returns {void}
     */
    #destroyMatchMetadata() {
        this.#match_metadata = null;
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
     * Whether there is a HotkeyMatch. meaning the match method was called and the hotkey produced a match. and the `run` method has not yet been called or has not ended.
     * hint: If the method `run` has been called but it is still running, the hotkey will be in a lock state.
     * @returns {boolean}
     */
    get HasMatch() {
        return this.#match_metadata != null;
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
     * Returns the length of the hotkey based on the number of keys it includes. this does take into account modifier keys. Vim motions are counted as 5 keys. which 
     * makes their matching much more specific.
     * @returns {number}
     */
    get Specificity() {
        return this.#key_combo_specificity;
    }

    /**
     * Whether the hotkey is locked. meaning the `run` method is currently running.
     * @returns {boolean}
     */
    get Locked() {
        return this.#hotkey_execution_mutex;
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
     * @param {import('@libs/utils').StackBuffer<KeyboardEvent>} event_history
     * @returns {boolean}
     */
    match(event_history) {
        if (!this.Valid || event_history.Size < this.Length) return false; // Now we are allowed to assume the hotkey has at least one fragment.

        this.#createMatchMetadata();

        /**
         * We start assuming the hotkey is a match, if any fragment doesn't match it's corresponding event, we set this to false.
         * @type {boolean}
         */
        let hotkey_matched = true;

        // Iterator indexes
        let fragment_h = 0;
        let event_k = 0;

        /**
         * If this is a match, the last fragment will be the first in the event history.
         * so we check the fragments from last to first so that they have a chance to match the history.
         */
        let history_fragments = [...this.#key_combo_fragments].reverse();

        let fragment = history_fragments[fragment_h];

        /** 
         * Whether the fragment matches it's corresponding event.
         * @type {boolean} 
         */
        let fragment_match;

        let last_sequence_time = -1;

        /**
         * Used to store numeric key when parsing a vim motion.
         * @type {string}
         */
        let motion_match_number = "";

        do {
            let event = event_history.PeekN(event_k);
            event_k++;

            if (event == null && motion_match_number === "") {
                hotkey_matched = false;
                break;
            }
            
            if (IsModifier(event?.key)) {
                console.log(`Modifier ${event.key} found. Skipping.`);
                continue;
            }

            if (event != null && this.#checkEventExpired(event, last_sequence_time)) {
                hotkey_matched = false;
                break;
            } 


            if (fragment.NumericMetakey) { // Parse vim motion. If matches, interrupts the flow in all cases.
                // console.log("Parsing vim motion");
                fragment_match = false;

                if (event != null) {
                    fragment_match = fragment.matchNumericMetakey(event);
                }

                if (fragment_match) {
                    motion_match_number += event.key;
                    continue;
                }

                if (motion_match_number === "") {
                    this.#destroyMatchMetadata();
                    console.error("This doesn't make sense to me. There is probably some kind of bug if this line is ever triggered.");
                    hotkey_matched = false;
                    continue;
                }

                this.#match_metadata.addReversedMotionMatch(motion_match_number);
                fragment_h++;
                fragment = history_fragments[fragment_h];
                continue;
            } // Anything after this line can safely assume that this if statement did not match.

            fragment_match = fragment.match(event);

            if (!fragment_match) {
                console.log(`Fragment ${fragment.Identity} did not match event ${event.key}`);
                hotkey_matched = false;
            }

            fragment_h++;
            fragment = history_fragments[fragment_h];
            last_sequence_time = event.timeStamp;

        } while (hotkey_matched && fragment != null); // If we run out of fragments and hotkey_matched is still true, that should mean a positive match.

        console.log(`Hotkey ${this.#key_combo} matched: ${hotkey_matched}`);

        if (!hotkey_matched) {
            this.#destroyMatchMetadata();
        } else {
            this.#match_metadata.setSuccessful();
        }

        return hotkey_matched;
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
     * @returns {string}
     */
    get KeyCombo() {
        return this.#key_combo;
    }

    /**
     * Runs the hotkey's callback
     * @param {KeyboardEvent} event
     * @returns {Promise<void>}
     */
    async run(event) {
        if (this.#hotkey_execution_mutex) {
            console.error(`Hotkey ${this.#key_combo} is already running. Skipping execution.`);
            return;
        }

        this.#hotkey_execution_mutex = true;
        if (this.#callback?.constructor.name === "Function") {
            this.#callback(event, this);
        } else if (this.#callback?.constructor.name === "AsyncFunction") {
            await this.#callback(event, this);
        }
        this.#hotkey_execution_mutex = false;

        this.#destroyMatchMetadata();
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
    #unpackOptions() {
        this.#mode = this.#the_options.mode;
        this.#description = this.#the_options.description;

        if (this.#mode === "keypress") {
            this.#mode = "keydown" // Keypress is deprecated and will(someday) be removed according to MDN
        }

        if (this.#the_options.await_execution === undefined) {
            this.#the_options.await_execution = default_hotkey_register_options.await_execution;
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
     * Returns the hotkey match metadata.
     * @returns {HotkeyMatch | null}
     */
    get MatchMetadata() {
        return this.#match_metadata;
    }

    /**
     * Whether the hotkey has a vim motion or not.
     * @returns {boolean}
     */
    get WithVimMotion() {
        return this.#has_vim_motion;
    }
}