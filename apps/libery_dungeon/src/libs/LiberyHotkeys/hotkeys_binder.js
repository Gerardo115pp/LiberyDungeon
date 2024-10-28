import { StackBuffer } from "@libs/utils";
import { 
    MAX_PAST_EVENTS, 
    HOTKEY_SPECIFICITY_PRECEDENCE,
    MAX_PAST_HOTKEYS_TRIGGERED, 
    MIN_TIME_BETWEEN_HOTKEY_REPEATS,
    DISABLE_KEYPRESS_EVENTS,
} from "./hotkeys_consts";
import { IsNumeric } from "./hotkeys_matchers";

class HotkeyTriggeredEvent {
    /**
     * The hotkey that was triggered
     * @type {import('./hotkeys').HotkeyData}
     */
    #the_hotkey;

    /**
     * The event that triggered the hotkey
     * @type {KeyboardEvent}
     */
    #the_event;

    /**
     * @param {import('./hotkeys')} hotkey 
     * @param {KeyboardEvent} event 
     */
    constructor(hotkey, event) {
        this.#the_hotkey = hotkey;
        this.#the_event = event;
    }

    /**
     * The key combo that triggered the hotkey
     * @type {string}
     */
    get KeyCombo() {
        return this.#the_hotkey.KeyCombo;
    }

    /**
     * Returns the time at which the hotkey was triggered.
     * @type {number}
     */
    get TriggerTime() {
        return this.#the_event.timeStamp;
    }
}

export class HotkeysController {

    /**
     * A stack with the last MAX_PAST_EVENTS keydown KeyboardEvents
     * @type {StackBuffer<KeyboardEvent>}
     */
    #keyboard_past_keydowns;

    /**
     * The last N hotkeys combos triggered by a keydown event.
     * @type {StackBuffer<HotkeyTriggeredEvent>}
     */
    #last_keydown_hotkeys_triggered;

    /**
     * A stack with the last MAX_PAST_EVENTS keyup KeyboardEvents
     * @type {StackBuffer<KeyboardEvent>}
     */
    #keyboard_past_keyups;

    /**
     * The last N hotkeys combos triggered by a keyup event.
     * @type {StackBuffer<HotkeyTriggeredEvent>}
     */
    #last_keyup_hotkeys_triggered;

    /**
     * keydown handler function bound to the HotkeyBinder instance so addEventListener wont override the this context with the EventTarget.
     * @type {function}
     */
    #bound_handleKeyDown;

    /**
     * keyup handler function bound to the HotkeyBinder instance so addEventListener wont override the this context with the EventTarget.
     * @type {function}
     */
    #bound_handleKeyUp;

    /**
     * Handler for the deprecated keypress event. It's only purpose is to block the event.
     */
    #bound_handleKeyPress;

    /**
     * A map of Hotkey triggers -> HotkeyData for keydown events. 
     * @type {Map<string, import('./hotkeys').HotkeyData[]>}
     */
    #keydown_hotkey_triggers;

    /**
     * A map of Hotkey triggers -> HotkeyData[] for keyup events. 
     * @type {Map<string, import('./hotkeys').HotkeyData[]>}
     */
    #keyup_hotkey_triggers;

    /**
     * The current hotkey context.
     * @type {import('./hotkeys_context').default | null}
     * @default null
     */
    #current_hotkey_context;

    /**
     * Whether a hotkey that requires awaiting is currently running or not.
     * @type {boolean}
     */
    #locked_on_execution;

    /**
     * Whether the execution is paused or not.
     * @type {boolean}
     */
    #paused;

    constructor() {
        if (globalThis.addEventListener === undefined) {
            throw new Error("This environment does not support event listeners")
        }


        this.#keyboard_past_keydowns = new StackBuffer(MAX_PAST_EVENTS);
        this.#last_keydown_hotkeys_triggered = new StackBuffer(MAX_PAST_HOTKEYS_TRIGGERED);
        this.#keyboard_past_keyups = new StackBuffer(MAX_PAST_EVENTS);
        this.#last_keyup_hotkeys_triggered = new StackBuffer(MAX_PAST_HOTKEYS_TRIGGERED);

        this.#keydown_hotkey_triggers = new Map();
        this.#keyup_hotkey_triggers = new Map();
        
        this.#bound_handleKeyDown = this.#handleKeyDown.bind(this);
        this.#bound_handleKeyUp = this.#handleKeyUp.bind(this);
        this.#bound_handleKeyPress = this.#handleKeyPress.bind(this);

        this.#current_hotkey_context = null;

        this.#locked_on_execution = false;
        this.#paused = false;

        this.#setup()
    }

    /**
     * Activates a given hotkey.
     * @param {import('./hotkeys').HotkeyData} hotkey
     * @param {KeyboardEvent} event
     */
    async #activateHotkey(hotkey, event) {
        if (hotkey == null) {
            throw new Error("Hotkey is null")
        }

        if (this.#locked_on_execution) {
            console.error("Hotkey execution is locked. Ignoring hotkey activation");
            return;
        }

        if (event.repeat) {
            let block_repeat = this.#shouldBlockHotkeyRepeat(hotkey, event);
            if (block_repeat) {
                return;
            }
        }

        this.#registerTriggeredHotkey(hotkey, event);

        this.#locked_on_execution = hotkey.AwaitExecution;
        try {
            await hotkey.run(event);
        } catch (error) {
            if (hotkey.Locked) {
                hotkey.releaseExecutionMutex(); // Only needs to be in catch cause the run method doesn't locks execution and releases it on success. But if an error occurs, the lock should be released.
            }

            console.error(`Error executing hotkey ${hotkey.KeyCombo}`, error);
        } finally {
            this.#locked_on_execution = false;
        }
    }

    /**
     * Sets a new hotkey context. if there is already a context set, it will be disabled.
     * @param {import('./hotkeys_context').default} new_context
     */
    bindContext(new_context) {
        if (this.#current_hotkey_context != null) {
            this.dropContext();
        }

        this.#current_hotkey_context = new_context;

        this.#populateHotkeyTriggers(new_context);
    }

    /**
     * Drops the current hotkey context
     * @returns {void}
     */
    dropContext() {
        if (this.#current_hotkey_context == null) return;

        this.#keydown_hotkey_triggers.clear();
        this.#keyup_hotkey_triggers.clear();

        this.#keyboard_past_keydowns.Clear(); // Prevent key strokes issued in previous contexts from affecting behavior in the new context.
        this.#keyboard_past_keyups.Clear(); 

        this.#current_hotkey_context = null;
    }

    /**
     * Cleans up the binder
     */
    destroy() {
        globalThis.removeEventListener("keydown", this.#bound_handleKeyDown);
        globalThis.removeEventListener("keyup", this.#bound_handleKeyUp);
    }

    /**
     * Whether the binder is locked for any reason and should not trigger hotkeys.
     * @type {boolean}
     */
    get ExecutionsForbidden() {
        let locked = false;

        locked = this.#locked_on_execution; 

        locked = locked || this.#paused;

        return locked;
    }

    /**
     * Returns the last N events of the given type.
     * @param {"keydown"|"keyup"} event_type
     * @param {number} n
     * @returns {KeyboardEvent[]}
     */
    #getLastEvents(event_type, n) {
        if (n < 0) {
            throw new Error("n must be greater than 0")
        }

        let all_past_events = event_type === "keydown" ? this.#keyboard_past_keydowns : this.#keyboard_past_keyups;

        if (all_past_events.Size < n) {
            throw new Error(`There are not enough past ${event_type} events to return ${n} of them`);
        }

        let last_n_events = [];

        for (let h = 0; h < n; h++) {
            let event =all_past_events.PeekN(h);

            if (event == null) break;
            
            last_n_events.push(event);
        }

        return last_n_events;
    }

    /**
     * Returns all the previous events that match a vim motion.
     * @returns {KeyboardEvent[]}
     */
    #getVimMotionMatchingEvents() {
        let matching_events = [];

        let event_index = 0;
        let matches_vim_motion = true;

        do {
            let event = this.#keyboard_past_keydowns.PeekN(event_index);

            if (event == null) break;

            matches_vim_motion = IsNumeric(event.key);

            if (matches_vim_motion) {
                matching_events.push(event);
            }
           
            event_index++;
        } while (matches_vim_motion);

        return matching_events;
    }

    /**
     * Returns the candidate hotkeys for a given key. Meaning the hotkeys that could potentially be triggered by the key.
     * @param {string} key
     * @param {string} mode
     * @returns {import('./hotkeys').HotkeyData[]}
     */
    #getCandidateHotkeys(key, mode) {
        let candidate_hotkeys = [];

        let triggers = mode === "keydown" ? this.#keydown_hotkey_triggers : this.#keyup_hotkey_triggers;

        let lower_case_matching = triggers.get(key.toLowerCase()) ?? []; // This is done so, for example hotkeys like shift+a are able to
        let upper_case_matching = triggers.get(key.toUpperCase()) ?? [];// be matched, cause the event.key will return 'A' and not 'a'
        let exact_matching = [];

        if (key.length > 1) { // Allow matching hotkeys like 'up' which will match 'ArrowUp'. Same case for 'enter' that matches 'Enter', 'backspace' that matches 'Backspace' etc.
            exact_matching = triggers.get(key) ?? [];
        }

        candidate_hotkeys = [
            ...lower_case_matching,
            ...upper_case_matching,
            ...exact_matching
        ];

        return candidate_hotkeys;
    }

    /**
     * Whether the binder has an active hotkey context or not.
     * @type {boolean}
     */
    get HasContext() {
        return this.#current_hotkey_context != null;
    }

    /**
     * Handles the keydown event
     * @param {KeyboardEvent} event
     */
    #handleKeyDown(event) {
        if (this.#shouldIgnoreEvent(event)) return;
        // console.log("keydown: ", event);

        this.#keyboard_past_keydowns.Add(event);
        
        if (this.ExecutionsForbidden) return;

        let matching_hotkey = this.#matchHotkey(event);

        if (matching_hotkey != null) {
            this.#activateHotkey(matching_hotkey, event);
        }
    }

    /**
     * Handles the keyup event
     * @param {KeyboardEvent} event
     */
    #handleKeyUp(event) {
        if (this.#shouldIgnoreEvent(event)) return;

        this.#keyboard_past_keyups.Add(event);

        if (this.ExecutionsForbidden) return;

        let matching_hotkey = this.#matchHotkey(event);

        if (matching_hotkey != null) {
            this.#activateHotkey(matching_hotkey, event);
        }
    }

    /**
     * Handles the keypress event. It's only purpose is to block the event unless the event is produced on a Input like element.
     * @param {KeyboardEvent} event
     */
    #handleKeyPress(event) {
        if (this.#shouldIgnoreEvent(event)) return;

        event.preventDefault();
        event.stopPropagation();
    }

    /**
     * Matches a Keyboard event with registered hotkeys. The matching behavior is based on HOTKEY_LENGTH_PRECEDENCE value. If true
     * it checks all candidates, meaning all hotkeys that have the same trigger as the event's key, and picks the longest hotkey. if
     * false, it checks all candidates until it finds one that matches and immediately returns it without checking the rest.
     * @param {KeyboardEvent} event
     * @returns {import('./hotkeys').HotkeyData | null}
     */
    #matchHotkey(event) {
        if (this.#current_hotkey_context == null || this.ExecutionsForbidden) return null;

        let is_keydown = event.type === "keydown";
        let is_keyup = event.type === "keyup";

        if (!is_keydown && !is_keyup) {
            throw new Error("Event type not supported")
        }

        let past_events = is_keydown ? this.#keyboard_past_keydowns : this.#keyboard_past_keyups;

        let candidate_hotkeys = this.#getCandidateHotkeys(event.key, event.type);
        // console.log("candidate hotkeys: ", candidate_hotkeys);

        if (candidate_hotkeys.length === 0) return;

        /** @type {import('./hotkeys').HotkeyData | null} */
        let matching_hotkey = null;

        if (HOTKEY_SPECIFICITY_PRECEDENCE) {
            matching_hotkey = this.#matchValidHotkey_BySpecificity(past_events, candidate_hotkeys);
        } else {
            matching_hotkey = this.#matchFirstValidHotkey(past_events, candidate_hotkeys);
        }

        return matching_hotkey;
    }

    /**
     * Returns the first hotkey from the candidates that matches the past events.
     * @param {StackBuffer<KeyboardEvent>} past_events
     * @param {import('./hotkeys').HotkeyData[]} candidates
     * @returns {import('./hotkeys').HotkeyData | null}
     */
    #matchFirstValidHotkey(past_events, candidates) {
        let matching_hotkey = null;

        for (let hotkey of candidates) {

            let matches = hotkey.match(past_events);

            if (matches) {
                matching_hotkey = hotkey;
                break;
            }
        }

        return matching_hotkey;
    }

    /**
     * Returns the candidate hotkey with the largest specificity that matches the past events.
     * @param {StackBuffer<KeyboardEvent>} past_events
     * @param {import('./hotkeys').HotkeyData[]} candidates
     * @returns {import('./hotkeys').HotkeyData | null}
     */
    #matchValidHotkey_BySpecificity(past_events, candidates) {
        let matching_hotkey = null;

        for (let hotkey of candidates) {
            let matches = hotkey.match(past_events);

            if (matches && (matching_hotkey == null || hotkey.Specificity > matching_hotkey.Specificity)) {
                matching_hotkey = hotkey;
            }
        }

        return matching_hotkey;
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
     * Populates the hotkey triggers maps with the hotkeys from given context.
     * @param {import('./hotkeys_context').default} context
     * @returns {void}
     */
    #populateHotkeyTriggers(context) {
        let hotkeys = context.hotkeys;

        for (let hotkey of hotkeys) {
            if (hotkey.Mode !== "keydown" && hotkey.Mode !== "keyup") {
                console.error(`Hotkey mode ${hotkey.Mode} is not supported. Skipping hotkey ${hotkey.KeyCombo}`);
                continue;
            }

            let all_triggers = hotkey.Mode === "keydown" ? this.#keydown_hotkey_triggers : this.#keyup_hotkey_triggers;

            const hotkey_triggers = hotkey.Triggers;

            for (let trigger of hotkey_triggers) {
                let hotkeys_on_trigger = [];

                if (all_triggers.has(trigger)) {
                    hotkeys_on_trigger = all_triggers.get(trigger);
                }

                hotkeys_on_trigger.push(hotkey);

                all_triggers.set(trigger, hotkeys_on_trigger);
            }
        }
    }

    /**
     * Pauses the hotkey execution and matching.
     * @returns {void}
     */
    pause() {
        this.#paused = true;
    }

    /**
     * Resumes the hotkey execution and matching.
     * @returns {void}
     */
    resume() {
        this.#paused = false;
    }


    /**
     * Registers a triggered hotkey 'event'(not as in dom event but as something that happened) on the appropriate history stack.
     * @param {import('./hotkeys').HotkeyData} hotkey
     * @param {KeyboardEvent} event
     */
    #registerTriggeredHotkey(hotkey, event) {
        let hotkey_event = new HotkeyTriggeredEvent(hotkey, event);

        let history_stack = hotkey.Mode === "keyup" ? this.#last_keyup_hotkeys_triggered : this.#last_keydown_hotkeys_triggered;

        history_stack.Add(hotkey_event)
    }

    /**
     * Checks if the given hotkey is repeating, if not, returns false. if it is, checks if it can repeat, if it can
     * returns false only if the time between executions is at least MIN_TIME_BETWEEN_HOTKEY_REPEATS if not it returns true.
     * @param {import('./hotkeys').HotkeyData} hotkey
     * @param {KeyboardEvent} event
     * @returns {boolean}
     */
    #shouldBlockHotkeyRepeat(hotkey, event) {
        if (!event.repeat) return false;

        if (!hotkey.CanRepeat) return true;

        let history_stack = hotkey.Mode === "keyup" ? this.#last_keyup_hotkeys_triggered : this.#last_keydown_hotkeys_triggered;

        let last_hotkey_trigger = history_stack.Peek();

        if (last_hotkey_trigger == null || last_hotkey_trigger.KeyCombo !== hotkey.KeyCombo) return false;

        let time_since_last_trigger = event.timeStamp - last_hotkey_trigger.TriggerTime;

        return time_since_last_trigger <= MIN_TIME_BETWEEN_HOTKEY_REPEATS;
    }

    /**
     * Attaches the binder to the global this object
     */
    #setup() {
        globalThis.addEventListener("keydown", this.#bound_handleKeyDown);
        globalThis.addEventListener("keyup", this.#bound_handleKeyUp);

        if (DISABLE_KEYPRESS_EVENTS) {
            globalThis.addEventListener("keypress", this.#bound_handleKeyPress);
        }
    }

    /**
     * Whether a given KeyboardEvent was originated from a Input like element or not. In which case, the
     * event should not trigger any hotkeys.
     * @param {KeyboardEvent} event
     */
    #shouldIgnoreEvent(event) {
        let origin_element = event.target;

        if (origin_element == null) return false;

        let should_ignore = false; 

        if (origin_element instanceof HTMLInputElement) {
            should_ignore = true;
        }

        if (!should_ignore && origin_element instanceof HTMLTextAreaElement) {
            should_ignore = true;
        }

        return should_ignore;
    }
}

/**
 * The global hotkeys binder
 * @type {HotkeysController | null}
 */
let GlobalHotkeysBinder = null;

export const setupHotkeysBinder = () => {
    if (GlobalHotkeysBinder != null) {
        throw new Error("Hotkeys binder already exists. Call destroyHotkeysBinder before setting up a new one")
    }

    GlobalHotkeysBinder = new HotkeysController();
}

export const destroyHotkeysBinder = () => {
    if (GlobalHotkeysBinder == null) {
        throw new Error("No hotkeys binder exists. Call setupHotkeysBinder before destroying it")
    }

    GlobalHotkeysBinder.destroy();

    GlobalHotkeysBinder = null;
}

/**
 * Returns the global hotkeys binder.
 * @returns {HotkeysController | null}
 */
export const getHotkeysBinder = () => {
    return GlobalHotkeysBinder;
}