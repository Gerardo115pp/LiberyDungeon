import { PlatformEventsTransmisor } from "@libs/DungeonsCommunication/transmissors/platform_events_transmisor";
import { Stack, hasWindowContext } from "@libs/utils";

/**
* @callback PlatformEventHandler
 * @param {import("@libs/DungeonsCommunication/transmissors/platform_events_transmisor").PlatformEventMessage} event
 * @returns {void}
*/


export class PlatformEventContext {
    /**
     * A event type to handler map
     * @type {Map<string, PlatformEventHandler>}
     */
    #events_map;

    constructor() {
        this.#events_map = new Map();
    }

    /**
     * Add a event handler to the context
     * @param {string} event_type the event name
     * @param {PlatformEventHandler} handler the event handler
     */
    addEventHandler(event_type, handler) {
        if (handler.constructor.name !== "AsyncFunction" && handler.constructor.name !== "Function") {
            throw new Error("The handler must be a function");
        }

        this.#events_map.set(event_type, handler);
    }

    /**
     * Returns an event handler from the context. If the context does not have a handler for the event_type, it will throw an error
     * @param {string} event_type the event name
     * @returns {PlatformEventHandler}
     */
    getEventHandler(event_type) {
        if (!this.#events_map.has(event_type)) {
            throw new Error("The event type does not exist");
        }

        return this.#events_map.get(event_type);
    }

    /**
     * Whether the context has a handler for the event_type
     * @param {string} event_type
     */
    hasEventHandler(event_type) {
        return this.#events_map.has(event_type);
    }

    /**
     * Remove a event handler from the context
     * @param {string} event_type the event name
     */
    removeEventHandler(event_type) {
        this.#events_map.delete(event_type);
    }
    
}

export class PlatformEventsContextManager {
    /**
     * A map of declared contexts and their names
     * @type {Map<string, PlatformEventContext>}
     */
    #contexts;

    /**
     * The current loaded context
     * @type {PlatformEventContext}
     */
    #current_context;

    /**
     * The current loaded context name
     * @type {string}
     */
    #current_context_name;

    /**
     * The event transmisor
     * @type {PlatformEventsTransmisor}
     */
    #transmisor;

    constructor() {
        this.#contexts = new Map();
        this.#current_context = null;
        this.#current_context_name = "";
        this.#transmisor = new PlatformEventsTransmisor();
    }

    /**
     * Connects the transmisor to start receiving events
     */
    start() {
        this.#transmisor.setOnMessageCallback(this.#processNewEvent.bind(this));

        this.#transmisor.connect();
    }

    /**
     * Stops receiving events by disconnecting the transmisor
     * @returns {void}
     */
    stop() {
        this.#transmisor.disconnect();
    }

    /**
     * Returns the current loaded context name
     * @returns {string}
     */
    get ContextName() {
        return this.#current_context_name;
    }

    /**
     * Declares a new context with out loading it
     * @param {string} context_name 
     * @param {PlatformEventContext} context
     */
    declareContext(context_name, context) {
        this.#contexts.set(context_name, context);
    }

    /**
     * Drops a context from the manager. if it's the current context, it will be unloaded
     * @param {string} context_name 
     */
    dropContext(context_name) {
        if (this.#current_context_name === context_name) {
            this.unloadContext();
        }
        
        this.#contexts.delete(context_name);
    }

    /**
     * Whether the manager has a context name registered
     * @param {string} context_name
     * @returns {boolean}
     */
    hasContext(context_name) {
        return this.#contexts.has(context_name);
    }

    /**
     * Returns whether there is a context loaded
     * @returns {boolean}
     */
    hasLoadedContext() {
        return this.#current_context !== null;
    }

    /**
     * Activates a new context. It has to exist in already by calling declareContext
     * @param {string} context_name 
     */
    loadContext(context_name) {
        if (!this.#contexts.has(context_name)) {
            throw new Error("The context does not exist");
        }
        console.log("Loading context: ", context_name);

        this.#current_context = this.#contexts.get(context_name);
        this.#current_context_name = context_name;  
    }

    /**
     * Process all new events and determines the handler to call
     * @param {import("@libs/DungeonsCommunication/transmissors/platform_events_transmisor").PlatformEventMessage} new_event
     */
    #processNewEvent(new_event) {
        if (this.#current_context === null || !this.#current_context.hasEventHandler(new_event.EventType)) {
            return this.defaultEventHandler(new_event);
        }

        let handler = this.#current_context.getEventHandler(new_event.EventType);
        handler(new_event);
    }

    /**
     * the default handler for new event, does nothing aside from logging the event.
     * its safe to override this method.
     * @param {import("@libs/DungeonsCommunication/transmissors/platform_events_transmisor").PlatformEventMessage} new_event
     * @returns {void}
     */
    defaultEventHandler(new_event) {
        console.log("New event received: ", new_event);
    }

    /**
     * Unloads the current context
     */
    unloadContext() {
        this.#current_context = null;
        this.#current_context_name = "";
    }
}

const global_platform_events_manager = new PlatformEventsContextManager();

if (hasWindowContext()) {
    window.global_platform_events_manager = global_platform_events_manager;
}

export default global_platform_events_manager;