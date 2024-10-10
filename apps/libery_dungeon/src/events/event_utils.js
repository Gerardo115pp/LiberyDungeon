/**
 * @typedef {Object} LiberyDungeonEvent
 * @property {string} event_name
 * @property {function} suscribe
 * @property {function} unsubscribe
 */

const suscribeToEvent = callback => {
    const event_name = this.event_name;

    if (event_name === undefined) {
        throw new Error('Event name is not defined, you probably forgot to bind this object to the event. dont use this function directly.');
    }

    document.addEventListener(event_name, callback);
}

const unsubscribeToEvent = callback => {
    const event_name = this.event_name;

    if (event_name === undefined) {
        throw new Error('Event name is not defined, you probably forgot to bind this object to the event. dont use this function directly.');
    }

    document.removeEventListener(event_name, callback);
}

/**
 * @param {LiberyDungeonEvent[]} event_list 
 */
export const bindEventMethods = event_list => {
    Object.entries(event_list).forEach(([key, value]) => {
        value.suscribe = suscribeToEvent.bind(value);
        value.unsubscribe = unsubscribeToEvent.bind(value);
    });
}