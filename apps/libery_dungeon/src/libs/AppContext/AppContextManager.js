/**
 * @description What is an app context? - An app context is namespace for related app pages. when changing app contexts, we need some stores to be reset among other things.
 * Why not to use location? - because two app pages can share the same app context, and have different locations. which is the case for the media viewer and the media explorer.
 */
class AppContextManager {
    /**
     * @type {string} - the name of the app context that is currently active
     */
    #current_context_name;
    /**
     * @type {Object<string, Function>} - an object app page -> function that is called when the app context is changed
     */
    #on_context_exit;
    constructor() {
        this.#current_context_name = "";
        this.#on_context_exit = {};
    }

    /**
     * Adds an app page together with a function that is called when the app context is changed
     * @param {string} page_name - the name of the app page
     * @param {function} fn - a function that is called when the app context is changed
    */
    addOnContextExit(page_name, fn) {
        this.#on_context_exit[page_name] = fn;
    }

    /**
     * @returns {string} - the name of the app context that is currently active
     */
    get CurrentContext() {
        return this.#current_context_name;
    }

    /**
     * Executes all the context exit functions
     */
    #executeContextExit() {
        for (const fn of Object.values(this.#on_context_exit)) {
            fn();
        }
    }

    /**
     * @param {string} context_name - the name of the app context to change to
     * @param {string} page_name - the name of the app page
     * @param {function} on_context_exit - a function that is called when the app context is changed
     */
    setAppContext(context_name, page_name, on_context_exit) {
        if (this.#current_context_name === context_name) return;

        if (this.#current_context_name !== "") {
            this.#executeContextExit();
        }

        this.#current_context_name = context_name;
        this.#on_context_exit = { page_name: on_context_exit };
    }
}

export const app_context_manager = new AppContextManager();