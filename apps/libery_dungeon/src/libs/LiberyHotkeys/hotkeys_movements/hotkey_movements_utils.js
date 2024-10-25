/**
 * wrap the cursor around a linear navigation, if the cursor reaches the maximum or minium position it will it
 * will wrap around to the other side
 * @param {number} current - where the cursor is currently located
 * @param {number} max - the maximum position where the cursor could be
 * @param {number} direction - the direction to move the cursor
 * @param {number} [min] - the minimum position where the cursor could be, by default 0
 * @returns {WrappedValue}
 * @typedef {Object} WrappedValue
 * @property {number} value
 * @property {boolean} overflowed_max
 * @property {boolean} overflowed_min
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