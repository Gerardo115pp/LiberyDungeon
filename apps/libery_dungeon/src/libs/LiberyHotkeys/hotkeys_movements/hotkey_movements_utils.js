/**
* @typedef {Object} WrappedValue
 * @property {number} value
 * @property {boolean} overflowed_max
 * @property {boolean} overflowed_min
*/

/**
 * wrap the cursor around a linear navigation, if the cursor reaches the maximum or minium position it will it
 * will wrap around to the other side
 * @param {number} current - where the cursor is currently located
 * @param {number} max - the maximum position where the cursor could be
 * @param {number} direction - the direction to move the cursor
 * @param {number} [min] - the minimum position where the cursor could be, by default 0
 * @returns {WrappedValue}
 */
export const linearCycleNavigationWrap = (current, max, direction, min=0) => {
    /** @type {WrappedValue} */
    const wrapped_value = {
        value: current,
        overflowed_max: false,
        overflowed_min: false
    }

    let next_value = current + direction;

    wrapped_value.value = next_value;

    if (next_value > max) {
        wrapped_value.value = min;
        wrapped_value.overflowed_max = true;
    } else if (next_value < min) {
        wrapped_value.value = max;
        wrapped_value.overflowed_min = true;
    }

    return wrapped_value;
}


/**
 * An Row in the GridRowSequence. hold references to the previous and next row and an array with all the indexes of the original sequence the row contains.
 */
class HM_GridRow {
    /**
     * The previous row in the sequence
     * @type {HM_GridRow}
     */
    #previous_row;

    /**
     * The next row in the sequence
     * @type {HM_GridRow}
     */
    #next_row;

    /**
     * The indexes of the original sequence that this row contains
     * @type {number[]}
     */
    #row_index_members;

    /**
     * @param {HM_GridRow} [previous_row] - the previous row in the sequence
     */
    constructor(previous_row=null) {
        if (!(previous_row instanceof HM_GridRow) && previous_row !== null) {
            throw new Error("Previous row must be an instance of HM_GridRow");
        }

        this.#row_index_members = [];
        this.#previous_row = previous_row;
        this.#next_row = null;
    }

    /**
     * Adds an index to the row that is exactly one unit larger than the last index in the row. if there are no indexes in the row, then it uses the greatest index in the previous row
     * as reference. if the previous row is null, then it sets 0 as the first index in the row. returns the length of the row after adding the index
     * @returns {number}
     */
    addIndex() {
        if (this.#previous_row === null) {
            this.#row_index_members = [0];
            return;
        }

        const last_index = this.#row_index_members[this.#row_index_members.length - 1] ?? this.#previous_row.MaxIndex;
        
        this.#row_index_members.push(last_index + 1);
    }

    /**
     * Returns a column index(not the value of said column) clamped to the min and max values of the row. 
     * @param {number} column_index
     * @returns {number}
     */
    clampColumnIndex(column_index) {
        return Math.min(0, Math.max(column_index, this.length - 1));
    }

    /**
     * Returns a wrapped value of the index in the given column index. so if the index is greater than the greatest index in the row, it will return the last index in the row and set overflowed_max to true
     * if the index is less than 0, it will return the smallest index in the row and set overflowed_min to true. to be clear, the column index has nothing to do with the sequence index. This method does not\
     * cycle the index, it will just clamp it to the min and max values of the row.
     * @param {number} column_index - the column index
     * @returns {WrappedValue}
     */
    getIndexInColumn(column_index) {
        let wrapped_value = {
            value: NaN,
            overflowed_max: false,
            overflowed_min: false
        }

        wrapped_value.overflowed_min = column_index < 0;
        wrapped_value.overflowed_max = column_index >= this.length;

        let clamped_index = this.clampColumnIndex(column_index);

        wrapped_value.value = this.#row_index_members[clamped_index];

        return wrapped_value;
    }

    /**
     * Returns the length of the row. 1-indexed
     * @returns {number}
     */
    get length() {
        return this.#row_index_members.length;
    }

    /**
     * Returns the greatest index in the row
     * @returns {number}
     */
    get MaxIndex() {
        return this.#row_index_members[this.#row_index_members.length - 1];
    }

    /**
     * Returns the minimum index in the row.
     * @type {number}
     */
    get MinIndex() {
        return this.#row_index_members[0];
    }

    /**
     * The next row in the sequence
     * @type {HM_GridRow | null}
     */
    get NextRow() {
        return this.#next_row;
    }

    /**
     * The previous row in the sequence
     * @type {HM_GridRow | null}
     */
    get PreviousRow() {
        return this.#previous_row;
    }

    /**
     * Sets the next row in the sequence. if the next row is already set, it panics.
     * @param {HM_GridRow} next_row - the next row in the sequence
     */
    setNextRow(next_row) {
        if (this.#next_row !== null) {
            throw new Error("Next row is already set");
        }

        this.#next_row = next_row;
    }
}

/**
 * A class to represent a sequence of rows with different lengths, each row represents an index of the sequence
 */
class HM_GridRowSequence {
    /**
     * The first row in the grid
     * @type {HM_GridRow}
     */
    #first_row;

    /**
     * The last row in the grid
     * @type {HM_GridRow}
     */
    #last_row;

    /**
     * The row from which indexes for index operations are been taken from.
     * @type {HM_GridRow}
     */
    #current_row;

    /**
     * The number of traversing operations that took place from the first row to get to the current row.
     * @type {number}
     */
    #current_rowid;

    /**
     * The last visited column index of the current row.
     * @type {number}
     */
    #current_column_index;

    /**
     * The count of rows in the grid.
     * @type {number}
     */
    #row_count;

    constructor() {
        this.#first_row = null;
        this.#last_row = null;
        this.#current_row = null;
        this.#current_rowid = NaN;
        this.#current_column_index = 0;
        this.#row_count = 0;
    }

    /**
     * Adds an index to current row that is exactly one unit larger than the last index in the row. if there are no indexes in the row, then it uses the greatest index in the previous row
     * as reference. if called with no rows, panics.
     */
    addIndex() {
        if (this.#current_row === null) {
            throw new Error("No rows in the sequence");
        }

        this.#current_row.addIndex();
    }

    /**
     * Appends a new row to the grid row sequence.
     */
    appendRow() {
        let new_row = new HM_GridRow(this.#last_row);

        if (this.#first_row === null) {
            this.#first_row = new_row;
            this.focusFirstRow();
        }

        if (this.#last_row !== null) {
            this.#last_row.setNextRow(new_row);
        }

        this.#last_row = new_row;
        this.#row_count++;
    }

    /**
     * Clears the grid row sequence.
     */
    clear() {
        this.#first_row = null;
        this.#last_row = null;
        this.#current_row = null;
        this.#current_rowid = NaN;
        this.#current_column_index = 0;
        this.#row_count = 0;
    }

    /**
     * Returns the index focused by the cursor.
     * @returns {number}
     */
    get Cursor() {
        if (this.#current_row == null) {
            throw new Error("In LiberyHotkeys/hotkeys_movement HM_GridRowSequence.Cursor: Trying to access a null cursor. The cursor is a combination of the current row and the current column index. But the current row was null.");
        }
        
        this.#current_row.getIndexInColumn(this.#current_column_index).value;
    }

    /**
     * Sets the row focus to the first row.
     */
    focusFirstRow() {
        this.#current_row = this.#first_row;
        this.#current_rowid = 0;
    }

    /**
     * Sets the row focus to the last row.
     */
    focusLastRow() {
        this.#current_row = this.#last_row;
        this.#current_rowid = this.#row_count - 1;
    }

    /**
     * Returns the row from which indexes for index operations are been taken from.
     * @returns {HM_GridRow}
     */
    getCurrentRow() {
        if (this.#current_row === null) {
            this.focusFirstRow();
        }

        return this.#current_row;
    }

    /**
     * The row count of the Grid.
     * @type {number}
     */
    get length() {
        return this.#current_row;
    }

    /**
     * Moves the cursor to the row above the current row and returns the new index of the sequence. 
     * @param {boolean} cycle - if true, and the current row is the first row, the cursor will be set to the last row.
     * @returns {number}
     */
    moveUp(cycle=false) {
        let previous_row = this.#traverseBackwards();

        if (previous_row == null) {

            if (!cycle) {
                return this.Cursor;
            };

            this.focusLastRow();
            previous_row = this.#current_row;
        }

        this.#current_column_index = previous_row.clampColumnIndex(this.#current_column_index);

        return this.Cursor;
    }

    /**
     * move the cursor to the next index in the current row. If overflows, does nothing unless cycle is set to true. In that case calls traverseForwards with the cycle true and sets the column index to 0.
     * Returns the sequence index in the resulting current index.
     * @param {boolean} cycle
     * @returns {number}
     */
    moveRight(cycle=false) {
        let new_index = this.#current_column_index + 1;

        if (new_index >= this.#current_row.length) {

            if (!cycle) {
                return this.#current_row.MaxIndex;
            }

            this.#traverseForward();
            new_index = 0;
        }

        this.#current_column_index = new_index;

        return this.Cursor;
    }

    /**
     * Moves the cursor to the next row in the sequence. returns the new current row. If there is no more rows, returns null and doesn't modify the current row.
     * on traversal, update the column index to keep it within bounds of the new row.
     * @param {boolean} cycle - if true, and the current row is the last row, the cursor will be set to the first row.
     * @returns {number}
     */
    moveDown(cycle=false) {
        if (this.#traverseForward() !== null) {
            return this.Cursor;
        };

        this.focusFirstRow();

        this.#current_column_index = this.#current_row.clampColumnIndex(this.#current_column_index);

        return this.Cursor;
    }

    /**
     * moves the cursor to the previous index in the current row. If overflows, does nothing unless cycle is set to true. In that case calls traverseBackwards with the cycle true and sets the column index to the last index in the row.
     * Returns the sequence index in the resulting current index.
     * @param {boolean} cycle
     * @returns {number}
     */
    moveLeft(cycle=false) {
        this.#current_column_index--;

        if (this.#current_column_index < 0) {
            this.moveUp();

            this.#current_column_index = this.#current_row.length - 1;

        }

        return this.Cursor;
    }

    /**
     * Sets a given row index as the current row.
     * @param {number} row_index
     */
    setCurrentRow(row_index) {
        if (row_index < 0 || row_index >= this.#row_count) {
            throw new Error("Row index out of bounds");
        }

        if (row_index < this.#current_rowid) {
            this.focusFirstRow();
        }

        let infinite_loop_guard = 0;

        while (this.#current_rowid < row_index) {
            this.#traverseForward();
            infinite_loop_guard++;

            if (infinite_loop_guard > this.#row_count) {
                throw new Error("Infinite loop detected");
            }
        }
    }

    /**
     * Traverses to the previous row in the sequence. returns the new current row. If there is no previous rows, returns null and doesn't modify the current row.
     * on traversal, update the column index to keep it within bounds of the new row.
     * @returns {HM_GridRow | null}
     */
    #traverseBackwards() {
        if (this.#current_row == null) {
            this.focusFirstRow();
            return null;
        }

        if (this.#current_row.PreviousRow == null) {
            return null;
        }

        this.#current_row = this.#current_row.PreviousRow;
        this.#current_rowid--;

        this.#current_column_index = this.#current_row.clampColumnIndex(this.#current_column_index);

        return this.#current_row; 
    }

    /**
     * Traverses to the next row in the sequence. returns the new current row. If there is no more rows, returns null and doesn't modify the current row.
     * on traversal, update the column index to keep it within bounds of the new row.
     * @returns {HM_GridRow | null}
     */
    #traverseForward() {
        if (this.#current_row === null) {
            return this.focusFirstRow();
        }

        if (this.#current_row.NextRow === null) {
            return null;
        }

        this.#current_row = this.#current_row.NextRow;
        this.#current_rowid++;

        this.#current_column_index = this.#current_row.clampColumnIndex(this.#current_column_index);

        return this.#current_row;
    }
}

/**
 * wraps the cursor around a grid navigation. Takes two selectors, one for the grid parent, and one for valid grid members. By default cycles the cursor if overflow occurs in any direction, any of them can be
 * opt-out. Once setup is called, it will monitor changes into the parent element by using the MutationObserver. to stop monitoring, call destroy(). It checks the elements height to determine if they belong to the
 * same row.
 */
export class GridNavigationWrapper {
    /**
     * The grid row sequence
     * @type {HM_GridRowSequence}
     */
    #grid_sequence;

    /**
     * The parent element of the grid
     * @type {HTMLElement | null}
     */
    #grid_parent;

    /**
     * The selector for the grid parent
     * @type {string}
     */
    #grid_parent_selector;

    /**
     * The selector for the grid members
     * @type {string}
     */
    #grid_member_selector;

    /**
     * Whether to cycle on up overflow movement.
     * @type {boolean}
     * @default true
     */
    #cycle_up;

    /**
     * Whether to cycle on right overflow movement.
     * @type {boolean}
     * @default true
     */
    #cycle_right;

    /**
     * Whether to cycle on down overflow movement.
     * @type {boolean}
     * @default true
     */
    #cycle_down;

    /**
     * Whether to cycle on left overflow movement.
     * @type {boolean}
     * @default true
     */
    #cycle_left;

    /**
     * The MutationObserver instance
     * @type {MutationObserver | null}
     */
    #mutation_observer;

    /**
     * The MutationObserver configuration
     * @type {MutationObserverInit}
     */
    #mutation_config;

    /**
     * @param {string} grid_parent_selector - the selector for the grid parent
     * @param {string} grid_member_selector - the selector for the grid members
     */
    constructor(grid_parent_selector, grid_member_selector) {
        this.#grid_sequence = new HM_GridRowSequence();
        this.#grid_parent = null;
        this.#grid_parent_selector = grid_parent_selector;
        this.#grid_member_selector = grid_member_selector;

        this.#cycle_up = true;
        this.#cycle_right = true;
        this.#cycle_down = true;
        this.#cycle_left = true;

        this.#mutation_observer = null;
        this.#mutation_config = {
            childList: true,
        };
    }
    
    /*=============================================
    =            Cycle Getter and setters            =
    =============================================*/
    
        /**
         * Whether the wrapper cycles on up overflow movement.
         * @type {boolean}
         * @default true
         */
        get CycleUp() {
            return this.#cycle_up;
        }

        /**
         * Whether the wrapper cycles on up overflow movement.
         * @type {boolean}
         */
        set CycleUp(value) {
            this.#cycle_up = value;
        }

        /**
         * Whether the wrapper cycles on right overflow movement.
         * @type {boolean}
         * @default true
         */
        get CycleRight() {
            return this.#cycle_right;
        }

        /**
         * Whether the wrapper cycles on right overflow movement.
         * @type {boolean}
         */
        set CycleRight(value) {
            this.#cycle_right = value;
        }

        /**
         * Whether the wrapper cycles on down overflow movement.
         * @type {boolean}
         * @default true
         */
        get CycleDown() {
            return this.#cycle_down;
        }

        /**
         * Whether the wrapper cycles on down overflow movement.
         * @type {boolean}
         */
        set CycleDown(value) {
            this.#cycle_down = value;
        }

        /**
         * Whether the wrapper cycles on left overflow movement.
         * @type {boolean}
         * @default true
         */
        get CycleLeft() {
            return this.#cycle_left;
        }

        /**
         * Whether the wrapper cycles on left overflow movement.
         * @type {boolean}
         */
        set CycleLeft(value) {
            this.#cycle_left = value;
        }
    
    /*=====  End of Cycle Getter and setters  ======*/

    /**
     * Returns the dom element matching the parent selector.
     * @returns {HTMLElement}
     */
    getDomParentElement() {
        return document.querySelector(this.#grid_parent_selector);
    }

    /**
     * Whether the wrapper has been setup.
     * @type {boolean}
     */
    get IsSetup() {
        return this.#grid_parent !== null;
    }

    /**
     * Callback for the MutationObserver.
     * @param {MutationRecord[]} mutations
     * @param {MutationObserver} observer
     */
    #onMutation(mutations, observer) {
        console.log("Mutation detected", mutations);

        this.scanGridMembers();
    }

    /**
     * Sets up the wrapper. If the wrapper is already setup or is called outside a Window context, it panics.
     */
    setup() {
        if (this.IsSetup ) {
            throw new Error("Already setup");
        }

        if (globalThis.self == null) {
            throw new Error("Cannot setup outside a Window context");
        }

        this.#grid_parent = this.getDomParentElement();

        if (this.#grid_parent === null) {
            throw new Error("Grid parent not found");
        }

        if (this.#grid_parent.childElementCount > 0) {
            this.scanGridMembers();
        }

        this.#mutation_observer = new MutationObserver(this.#onMutation.bind(this));

        this.#mutation_observer.observe(this.#grid_parent, this.#mutation_config);
    }

    /**
     * Scans and creates the HM_GridRowSequence from the matching grid members. Determines which element belong to the same row by checking their .getBoundingClientRect().y value.
     */
    scanGridMembers() {

        if (!this.#grid_parent.hasChildNodes()) {
            return;
        }

        let grid_members = this.#grid_parent.querySelectorAll(this.#grid_member_selector);

        if (grid_members.length === 0) {
            console.warn(`No grid members found with selecto '${this.#grid_member_selector}'`);
            return;
        }
        this.#grid_sequence.clear();

        let previous_element_y = NaN;

        for (let h = 1; h < grid_members.length; h++) {
            let current_element_rect = grid_members[h].getBoundingClientRect();

            if (isNaN(previous_element_y) || current_element_rect.y > previous_element_y) {
                this.#grid_sequence.appendRow();
            }
                
            this.#grid_sequence.addIndex();           

            previous_element_y = current_element_rect.y;
        }
    }
}

