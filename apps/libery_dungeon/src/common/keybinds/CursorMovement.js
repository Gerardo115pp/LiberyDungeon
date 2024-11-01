import { GridNavigationWrapper } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";

/** 
 * A class used by default by the GridMovement wrappers to determine the sequence members. it is appended to the passed grid container selector resulting in a selector of the
 * form "#the-container-selector .dungeon-grid-member-item"
 * @type {string} 
 */
export const GRID_MOVEMENT_ITEM_CLASS = ".dungeon-grid-member-item";

/**
* A callback function that is called when the cursor position index is updated. If the callback returns true(and only true, not any other truthy value), the cursor position will be rolled back to the previous position.
* @callback CursorPositionCallback
 * @param {import("@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils").GridWrappedValue}
 * @returns {boolean | void}
*/

/**
 * Optional parameters for the CursorMovementWASD wrapper. Setting any triggers to an empty string will disable the hotkey.
* @typedef {Object} MovementModelOptions
 * @property {string} grid_member_selector - The selector for the grid members. Defaults to `GRID_MOVEMENT_ITEM_CLASS`
 * @property {number} initial_cursor_position - The initial cursor position. Defaults to 0
 * @property {string} sequence_item_name - Used to generate descriptions for user feedback. Defaults to "item"
 * @property {string} sequence_item_name_plural - Used to generate descriptions for user feedback. Defaults to "items"
 * @property {string} up_trigger - The key that moves the cursor up. Defaults to "w"
 * @property {string} down_trigger - The key that moves the cursor down. Defaults to "s"
 * @property {string} left_trigger - The key that moves the cursor left. Defaults to "a"
 * @property {string} right_trigger - The key that moves the cursor right. Defaults to "d"
 * @property {string} goto_item_finalizer - The key that finalizes the goto vim motion. Defaults to "g"
 * @property {string} goto_row_finalizer - The key that finalizes the goto row vim motion. Defaults to "l"(for line)
 * @property {string} row_start_trigger - The key that moves the cursor to the start of the row. Defaults to "shift+a"
 * @property {string} row_end_trigger - The key that moves the cursor to the end of the row. Defaults to "shift+d"
 * @property {string} first_row_trigger - The key that moves the cursor to the first row. Defaults to "shift+w"
 * @property {string} last_row_trigger - The key that moves the cursor to the last row. Defaults to "shift+s"
*/

/**
 * @type {MovementModelOptions}
 */
const default_cursor_movement_wasd_options = {
    grid_member_selector: GRID_MOVEMENT_ITEM_CLASS,
    initial_cursor_position: 0,
    sequence_item_name: "item",
    sequence_item_name_plural: "items",
    up_trigger: "w",
    down_trigger: "s",
    left_trigger: "a",
    right_trigger: "d",
    goto_item_finalizer: "g",
    goto_row_finalizer: "l",
    row_start_trigger: "shift+a",
    row_end_trigger: "shift+d",
    first_row_trigger: "shift+w",
    last_row_trigger: "shift+s",
}

/**
 * A class that abstracts common wasd cursor movement through a list of items presented in the UI as a grid, sequence of `N` rows with `M_n` columns, where `M_n` is
 * different for each row. Uses the `GridNavigationWrapper` under the hood, which determines what elements belong to different rows by checking their .DOMRect.y values.
 * Registers a single callback which is expected to receive a new cursor position. 
 */
export class CursorMovementWASD {
    /**
     * The callback that is called when the cursor position is updated.
     * @type {CursorPositionCallback}
     */
    #cursor_position_callback

    /**
     * The GridNavigationWrapper instance that is used to navigate through the grid.
     * @type {GridNavigationWrapper}
     */
    #grid_navigation_wrapper

    /**
     * @type {MovementModelOptions}
     */
    #movement_options;

    /**
     * @param {string} grid_container_selector - The selector for the grid container.
     * @param {CursorPositionCallback} cursor_position_callback - The callback that is called when the cursor position is updated.
     * @param {MovementModelOptions} [options] - Optional parameters for the CursorMovementWASD wrapper.
     */
    constructor(grid_container_selector, cursor_position_callback, options) {
        if (options == null) {
            options = {};
        }

        this.#cursor_position_callback = cursor_position_callback;
        this.#movement_options = {
            ...default_cursor_movement_wasd_options,
            ...options
        }

        this.#grid_navigation_wrapper = new GridNavigationWrapper(grid_container_selector, this.#movement_options.grid_member_selector);
    }

    /**
     * Changes the grid container selector.
     * @param {string} grid_container_selector - The selector for the grid container.
     * @param {string} [grid_member_selector] - a new selector for the grid members. If not passed, keeps the one in the options and prepends the grid container selector to it.
     */
    changeGridContainer(grid_container_selector, grid_member_selector) {
        if (grid_member_selector == null || grid_member_selector === "") {
            grid_member_selector = this.#movement_options.grid_member_selector;
        } else {
            this.#movement_options.grid_member_selector = grid_member_selector;
        }

        if (this.#grid_navigation_wrapper != null) {
            this.#grid_navigation_wrapper.destroy();
        }

        this.#grid_navigation_wrapper = new GridNavigationWrapper(grid_container_selector, grid_member_selector);
    }

    /**
     * Cleans resources used by the CursorMovementWASD wrapper.
     * @returns {void}
     */
    destroy() {
        this.#grid_navigation_wrapper.destroy();
    }

    /**
     * Defines the common wasd cursor movement hotkeys. mapped to common wasd keys/modifiers.
     * @param {import("@libs/LiberyHotkeys/hotkeys_context").default} hotkeys_context
     * @returns {void}
     */
    #defineWASDMovementHotkeys(hotkeys_context) {
        
        // WASD movement hotkeys.
        if (this.#movement_options.up_trigger !== "") {
            const wasd_triggers = this.#getWASDMovementTriggers();
            
            hotkeys_context.register(wasd_triggers, this.#wasdMovementHandler.bind(this), {
                description: `<navigation>Moves changes the focused ${this.#movement_options.sequence_item_name}.`,
            });
        }

        // Goto item hotkey.
        if (this.#movement_options.goto_item_finalizer !== "") {
            const goto_item_triggers = this.#getGotoItemTriggers();
    
            hotkeys_context.register(goto_item_triggers, this.#wasdGotoItemHandler.bind(this), {
                description: `<navigation>Goes to a specific ${this.#movement_options.sequence_item_name}.`,
            });
        }

        // Goto row hotkey.
        if (this.#movement_options.goto_row_finalizer !== "") {
            const goto_row_triggers = this.#getGotoRowTriggers();

            hotkeys_context.register(goto_row_triggers, this.#wasdGotoRowHandler.bind(this), {
                description: `<navigation>Goes to a specific line of ${this.#movement_options.sequence_item_name_plural}.`,
            });
        }

        // Row start hotkey.
        if (this.#movement_options.row_start_trigger !== "") {
            const row_start_triggers = this.#getRowStartTriggers();

            hotkeys_context.register(row_start_triggers, this.#wasdRowStartHandler.bind(this), {
                description: `<navigation>Moves the cursor to the start of the row.`,
            });
        }

        // Row end hotkey.
        if (this.#movement_options.row_end_trigger !== "") {
            const row_end_triggers = this.#getRowEndTriggers();

            hotkeys_context.register(row_end_triggers, this.#wasdRowEndHandler.bind(this), {
                description: `<navigation>Moves the cursor to the end of the row.`,
            });
        }

        // First row hotkey.
        if (this.#movement_options.first_row_trigger !== "") {
            const first_row_triggers = this.#getFirstRowTriggers();

            hotkeys_context.register(first_row_triggers, this.#wasdFirstRowHandler.bind(this), {
                description: `<navigation>Moves the cursor to the first row.`,
            });
        }

        // Last row hotkey.
        if (this.#movement_options.last_row_trigger !== "") {
            const last_row_triggers = this.#getLastRowTriggers();

            hotkeys_context.register(last_row_triggers, this.#wasdLastRowHandler.bind(this), {
                description: `<navigation>Moves the cursor to the last row.`,
            });
        }
    }

    /**
     * Returns an array with the wasd movement triggers.
     * @returns {string[]}
     */
    #getWASDMovementTriggers() {
        return [
            this.#movement_options.up_trigger,
            this.#movement_options.down_trigger,
            this.#movement_options.left_trigger,
            this.#movement_options.right_trigger
        ]   
    }

    /**
     * Returns an array with the goto triggers for signle elements.
     * @returns {string[]}
     */
    #getGotoItemTriggers() {
        return [ `\\d ${this.#movement_options.goto_item_finalizer}` ];
    }

    /**
     * Returns an array with the goto triggers for rows.
     * @returns {string[]}
     */
    #getGotoRowTriggers() {
        return [ `\\d ${this.#movement_options.goto_row_finalizer}` ];
    }

    /**
     * Returns an array with the triggers for row start.
     * @returns {string[]}
     */
    #getRowStartTriggers() {
        return [ this.#movement_options.row_start_trigger ];
    }

    /**
     * Returns an array with the triggers for row end.
     * @returns {string[]}
     */
    #getRowEndTriggers() {
        return [ this.#movement_options.row_end_trigger ];
    }   

    /**
     * Returns an array with the triggers for the first row.
     * @returns {string[]}
     */
    #getFirstRowTriggers() {
        return [ this.#movement_options.first_row_trigger ];
    }

    /**
     * Returns an array with the triggers for the last row.
     * @returns {string[]}
     */
    #getLastRowTriggers() {
        return [ this.#movement_options.last_row_trigger ];
    }

    /**
     * Initializes the CursorMovementWASD wrapper. Call it only when the grid parent is mounted.
     * @param {import("@libs/LiberyHotkeys/hotkeys_context").default} hotkeys_context
     * @returns {void}
     */
    setup(hotkeys_context) {
        if (globalThis.self == null) {
            throw new Error("The CursorMovementWASD wrapper can only be used in a window context.");
        }

        // Setup the GridNavigationWrapper. 
        try {
            this.#grid_navigation_wrapper.setup();
        } catch (e) {
            console.error("In common/keybinds/CursorMovement.CursorMovementWASD while calling this.#grid_navigation_wrapper.setup(), selector exists?");
            throw e;
        }

        // If an initial cursor position other than 0(which is the default) is set, try to set it or panic.
        if (typeof this.#movement_options.initial_cursor_position === "number" && this.#movement_options.initial_cursor_position !== 0) {
            let cursor_set = this.#grid_navigation_wrapper.Grid.setCursor(this.#movement_options.initial_cursor_position);

            if (!cursor_set) {
                throw new Error("The initial cursor position is out of bounds or there was a problem setting it.");
            }
        }

        this.#defineWASDMovementHotkeys(hotkeys_context);
    }

    /**
     * Handles the cursor basic wasd movement.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #wasdMovementHandler = (event, hotkey) => {
        let current_cursor_position = this.#grid_navigation_wrapper.Grid.Cursor;

        /**
         * @type {import("@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils").GridWrappedValue}
         */
        let cursor_position;

        /**
         * The function to call and rollback the cursor position if the callback returns true.
         * @type {() => void}
         */
        let backward_move_function;

        switch (event.key) {
            case this.#movement_options.up_trigger:
                cursor_position = this.#grid_navigation_wrapper.Grid.moveUp();
                backward_move_function = this.#grid_navigation_wrapper.Grid.moveDown.bind(this.#grid_navigation_wrapper.Grid);
                break;
            case this.#movement_options.down_trigger:
                cursor_position = this.#grid_navigation_wrapper.Grid.moveDown();
                backward_move_function = this.#grid_navigation_wrapper.Grid.moveUp.bind(this.#grid_navigation_wrapper.Grid);
                break;
            case this.#movement_options.left_trigger:
                cursor_position = this.#grid_navigation_wrapper.Grid.moveLeft();
                backward_move_function = this.#grid_navigation_wrapper.Grid.moveRight.bind(this.#grid_navigation_wrapper.Grid);
                break;
            case this.#movement_options.right_trigger:
                cursor_position = this.#grid_navigation_wrapper.Grid.moveRight();
                backward_move_function = this.#grid_navigation_wrapper.Grid.moveLeft.bind(this.#grid_navigation_wrapper.Grid);
                break;
        }

        let should_rollback = this.#cursor_position_callback(cursor_position);

        if (should_rollback === true) {
            backward_move_function();
        }
    }

    /**
     * Handles the cursor goto element movement.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #wasdGotoItemHandler = (event, hotkey) => {
        if (!hotkey.WithVimMotion || !hotkey.HasMatch) return;

        let starting_cursor_position = this.#grid_navigation_wrapper.Grid.Cursor;

        let requested_position = hotkey.MatchMetadata.MotionMatches[0];
        console.log("requested_position", requested_position);

        requested_position = requested_position - 1; // 1-based index to 0-based index.

        requested_position = this.#grid_navigation_wrapper.Grid.clampSequenceIndex(requested_position);
        

        this.#grid_navigation_wrapper.Grid.setCursor(requested_position);


        let wrapped_cursor_position = this.#grid_navigation_wrapper.Grid.CursorWrapped;

        let should_rollback = this.#cursor_position_callback(wrapped_cursor_position);
        if (should_rollback === true) {
            this.#grid_navigation_wrapper.Grid.setCursor(starting_cursor_position);
        }
    }

    /**
     * Handles the cursor goto row movement.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #wasdGotoRowHandler = (event, hotkey) => {
        if (!hotkey.WithVimMotion || !hotkey.HasMatch) return;

        let starting_cursor_position = this.#grid_navigation_wrapper.Grid.CursorRow;

        let requested_position = hotkey.MatchMetadata.MotionMatches[0];

        requested_position = requested_position - 1; // 1-based index to 0-based index.
        console.log("requested_position", requested_position);

        requested_position = Math.max(0, Math.min(requested_position, this.#grid_navigation_wrapper.Grid.length - 1));

        this.#grid_navigation_wrapper.Grid.setCurrentRow(requested_position);

        let wrapped_cursor_position = this.#grid_navigation_wrapper.Grid.CursorWrapped;

        let should_rollback = this.#cursor_position_callback(wrapped_cursor_position);
        if (should_rollback === true) {
            this.#grid_navigation_wrapper.Grid.setCurrentRow(starting_cursor_position);
        }
    }

    /**
     * Handles the cursor row start movement.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #wasdRowStartHandler = (event, hotkey) => {
        let current_column = this.#grid_navigation_wrapper.Grid.CursorColumn;

        let cursor_position = this.#grid_navigation_wrapper.Grid.moveRowStart();

        let should_rollback = this.#cursor_position_callback(cursor_position);

        if (should_rollback === true) {
            this.#grid_navigation_wrapper.Grid.setCurrentRowColumn(current_column);
        }
    }

    /**
     * Handles the cursor row end movement.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #wasdRowEndHandler = (event, hotkey) => {
        let current_column = this.#grid_navigation_wrapper.Grid.CursorColumn;

        let cursor_position = this.#grid_navigation_wrapper.Grid.moveRowEnd();

        let should_rollback = this.#cursor_position_callback(cursor_position);

        if (should_rollback === true) {
            this.#grid_navigation_wrapper.Grid.setCurrentRowColumn(current_column);
        }
    }

    /**
     * Handles the cursor first row movement.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #wasdFirstRowHandler = (event, hotkey) => {
        let current_cursor = this.#grid_navigation_wrapper.Grid.Cursor;

        this.#grid_navigation_wrapper.Grid.focusFirstRow();

        let cursor_position = this.#grid_navigation_wrapper.Grid.CursorWrapped;

        let should_rollback = this.#cursor_position_callback(cursor_position);

        if (should_rollback === true) {
            this.#grid_navigation_wrapper.Grid.setCursor(current_cursor);
        }
    }

    /**
     * Handles the cursor last row movement.
     * @type {import("@libs/LiberyHotkeys/hotkeys").HotkeyCallback}
     */
    #wasdLastRowHandler = (event, hotkey) => {
        let current_cursor = this.#grid_navigation_wrapper.Grid.Cursor;

        this.#grid_navigation_wrapper.Grid.focusLastRow();

        let cursor_position = this.#grid_navigation_wrapper.Grid.CursorWrapped;

        let should_rollback = this.#cursor_position_callback(cursor_position);

        if (should_rollback === true) {
            this.#grid_navigation_wrapper.Grid.setCursor(current_cursor);
        }
    }
}