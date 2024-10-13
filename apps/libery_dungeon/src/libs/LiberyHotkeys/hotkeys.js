import { HotkeyFragment, IsNumeric } from "./hotkeys_matchers";
import { HOTKEYS_HIDDEN_GROUP, HOTKEYS_GENERAL_GROUP } from "./hotkeys_consts";
/**
* @typedef {Object} HotkeyRegisterOptions
 * @property {boolean} bind - If true the hotkey will be binded immediately. default is false
 * @property {?string} description - The hotkey's description   
 * @property {"keypress"|"keydown"|"keyup"} mode - The mode of the keypress event. Default is "keydown"
 * @property {boolean} [await_execution]
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
    mode: "keydown",
    await_execution: true
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

        this.#the_options = options;

        
        this.#unpackOptions(options)
        
        this.#is_valid = true;
        this.#has_vim_motion = false;
        this.#match_metadata = null;
        this.#hotkey_execution_mutex = false;

        this.#splitFragments()
    }


    /**
     * Whether a caller to `run` should await the callback's execution
     * @returns {boolean}
     */
    get AwaitExecution() {
        return this.#the_options.await_execution;
    }

    /**
     * Takes an array of Keyboard events, takes the ones that match a vim motion and populates with them the vim motion metadata.
     * returns a new array with the remaining events.
     * @param {KeyboardEvent[]} events
     * @returns {KeyboardEvent[]} 
     * @deprecated
     */
    #collectVimMotionEvents(events) {
        if (events.length === 0) return events;

        this.#match_metadata = "";

        let look_up_index = 0;
        let matches_vim_motion = true;

        do {
            let event = events[look_up_index];
            
            matches_vim_motion = IsNumeric(event.key);

            if (matches_vim_motion) {
                this.#match_metadata += event.key;
                look_up_index++;
            }

            if (look_up_index > 100000) {
                throw new Error("Infinite loop detected. Aborting."); // This should only happen if our code is not well thought out. If it is, then this will never happen(unless the user is an idiot). TODO: if this is not triggered in a year, remove it.
            }

        } while (matches_vim_motion && look_up_index < events.length);

        let remaining_events = events.slice(look_up_index);

        this.#match_metadata = this.#match_metadata === "" ? "0" : this.#match_metadata;
        if (this.#match_metadata !== "0") {
            this.#match_metadata = this.#match_metadata.split("").reverse().join("");
        }


        return remaining_events;
    }

    /**
     * Creates and sets a new match metadata.
     * @returns {HotkeyMatch}
     */
    #createMatchMetadata() {
        this.#match_metadata = new HotkeyMatch();
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
        if (!this.Valid) return false; // Now we are allowed to assume the hotkey has at least one fragment.
        console.log("Matching hotkey: ", this.#key_combo);
        console.log("Event history: ", event_history);

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

        /**
         * Used to store numeric key when parsing a vim motion.
         * @type {string}
         */
        let motion_match_number = "";

        do {
            let event = event_history.PeekN(event_k);
            event_k++;

            console.log(`Fragment: ${fragment.Identity}, Event: ${event?.key}`);

            // TODO: root this out and put it in a separate function.
            if (fragment.NumericMetakey) { // Parse vim motion. If matches, interrupts the flow in all cases.
                console.log("Parsing vim motion from: ", event)
                fragment_match = false;

                if (event != undefined) {
                    fragment_match = fragment.matchNumericMetakey(event);
                }

                if (fragment_match) {
                    console.log("matches")
                    motion_match_number += event.key;
                    continue;
                }

                console.log(`'${event?.key}' did not match a vim motion. final numeric string: ${motion_match_number}`);

                if (motion_match_number === "") {
                    this.#destroyMatchMetadata();
                    console.error("This doesn't make sense to me. There is probably some kind of bug if this line is ever triggered.");
                    hotkey_matched = false;
                    continue;
                }

                console.log("Reversing string");
                this.#match_metadata.addReversedMotionMatch(motion_match_number);
                fragment_h++;
                fragment = history_fragments[fragment_h];
                continue;
            } // Anything after this line can safely assume that this if statement did not match.

            fragment_match = fragment.match(event);

            if (!fragment_match) {
                console.log(`Fragment<${fragment.Identity}> did not match event<${event?.key}>`);
                hotkey_matched = false;
            }

            fragment_h++;
            fragment = history_fragments[fragment_h];
            console.log("Next fragment: ", fragment);
            console.log("Matched: ", hotkey_matched);

        } while (hotkey_matched && fragment != null); // If we run out of fragments and hotkey_matched is still true, that should mean a positive match.

        if (!hotkey_matched) {
            console.log("Hotkey did not match.");
            this.#destroyMatchMetadata();
        } else {
            console.log("Hotkey matched.");
            this.#match_metadata.setSuccessful();
        }

        console.log("===================================MATCHING ENDED===================================");

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