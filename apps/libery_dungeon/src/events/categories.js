import { bindEventMethods, LiberyDungeonEvent } from "@events/event_utils";

/**
 * an object containing all the events related to the categories.
 * @type {Object<string, LiberyDungeonEvent>} 
 */
export const categories_events = {
    LOADED_EMPTY_CATEGORY: {
        event_name: 'loaded-empty-category',
        suscribe: (callback) => {},
        unsubscribe: (callback) => {}
    },
}

export const emitLoadedEmptyCategory = () => {
    const event = new CustomEvent(categories_events.LOADED_EMPTY_CATEGORY);

    document.dispatchEvent(event);
}

bindEventMethods(categories_events);