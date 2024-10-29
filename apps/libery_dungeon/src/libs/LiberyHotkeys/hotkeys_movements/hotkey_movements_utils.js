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
    AddIndex() {
        if (this.#previous_row === null) {
            this.#row_index_members = [0];
        }

        const last_index = this.#row_index_members[this.#row_index_members.length - 1] ?? this.#previous_row.MaxIndex;
        
        this.#row_index_members.push(last_index + 1);
    }

    /**
     * Returns a wrapped value of the index in the given column index. so if the index is greater than the greatest index in the row, it will return the last index in the row and set overflowed_max to true
     * if the index is less than 0, it will return the smallest index in the row and set overflowed_min to true. to be clear, the column index has nothing to do with the sequence index. This method does not\
     * cycle the index, it will just clamp it to the min and max values of the row.
     * @param {number} column_index - the column index
     * @returns {WrappedValue}
     */
    GetIndexInColumn(column_index) {
        let wrapped_value = {
            value: NaN,
            overflowed_max: false,
            overflowed_min: false
        }

        wrapped_value.overflowed_min = column_index < 0;
        wrapped_value.overflowed_max = column_index >= this.length;

        let clamped_index = Math.max(0, Math.min(column_index, this.length - 1));

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
    SetNextRow(next_row) {
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
    AddIndex() {
        if (this.#current_row === null) {
            throw new Error("No rows in the sequence");
        }

        this.#current_row.AddIndex();
    }

    /**
     * Appends a new row to the grid row sequence.
     */
    AppendRow() {
        let new_row = new HM_GridRow(this.#last_row);

        if (this.#first_row === null) {
            this.#first_row = new_row;
            this.#restartGridRowsTraversal();
        }

        if (this.#last_row !== null) {
            this.#last_row.SetNextRow(new_row);
        }

        this.#last_row = new_row;
        this.#row_count++;
    }

    /**
     * Returns the row from which indexes for index operations are been taken from.
     * @returns {HM_GridRow}
     */
    GetCurrentRow() {
        if (this.#current_row === null) {
            this.#restartGridRowsTraversal();
        }

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
                return this.#current_column_index
            };

            previous_row = this.#last_row;
            this.#current_rowid = this.#row_count - 1;
        }

        this.#current_column_index = Math.max(0, Math.min(this.#current_column_index, previous_row.length - 1));

        return previous_row.GetIndexInColumn(this.#current_column_index).value;

        // TODO: TOMORROW: implement the rest of the move methods(moveDown, moveLeft, moveRight)
    }

    /**
     * Restarts the sequence traversal by setting setting the current row = first row. 
     */
    #restartGridRowsTraversal() {
        this.#current_row = this.#first_row;
        this.#current_rowid = 0;
    }

    /**
     * Sets a given row index as the current row.
     * @param {number} row_index
     */
    SetCurrentRow(row_index) {
        if (row_index < 0 || row_index >= this.#row_count) {
            throw new Error("Row index out of bounds");
        }

        if (row_index < this.#current_rowid) {
            this.#restartGridRowsTraversal();
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
     * @returns {HM_GridRow | null}
     */
    #traverseBackwards() {
        if (this.#current_row == null) {
            this.#restartGridRowsTraversal();
            return null;
        }

        if (this.#current_row.PreviousRow == null) {
            return null;
        }

        this.#current_row = this.#current_row.PreviousRow;
        this.#current_rowid--;

        return this.#current_row; 
    }

    /**
     * Traverses to the next row in the sequence. returns the new current row. If there is no more rows, returns null and doesn't modify the current row.
     * @returns {HM_GridRow | null}
     */
    #traverseForward() {
        if (this.#current_row === null) {
            return this.#restartGridRowsTraversal();
        }

        if (this.#current_row.NextRow === null) {
            return null;
        }

        this.#current_row = this.#current_row.NextRow;
        this.#current_rowid++;

        return this.#current_row;
    }
}

