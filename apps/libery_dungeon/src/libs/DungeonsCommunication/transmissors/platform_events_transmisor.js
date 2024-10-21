import { base_domain, jd_server } from "../services";
import { parseUnsecureJWT } from "@libs/utils";

export class PlatformEventsTransmisor {

    /**
     * The socket used for communication
     * @type {WebSocket}
     */
    #socket;

    /**
     * set by the caller, valid events received will be forwarded to this callback  
     * @param {PlatformEventMessage} event_message
     */
    #on_message_callback = event_message => console.log("Received event: ", event_message);
    
    constructor() {
        this.host = `wss://${base_domain}${jd_server}/platform-events/public/suscribe`;
        this.#socket = null;
    }

    connect = () => {
        this.#socket = new WebSocket(this.host);
        this.#socket.onopen = this.onOpen;
        this.#socket.onmessage = this.onMessage;
        this.#socket.onclose = this.onClose;
        this.#socket.onerror = this.onError;
    }

    disconnect = () => {
        if (this.#socket == null) return;
        
        this.#socket.close();
    }

    /**
     * @param {MessageEvent} message 
     */
    onMessage = message => {
        console.debug("received message from server: ", message)
        const data = JSON.parse(message.data);

        /**
         * @type {PlatformEventMessage}
         */
        let event_message;

        try {
            event_message = new PlatformEventMessage(data);
        } catch (error) {
            console.error(`Wrong event message format: ${error}`);
        }

        this.#on_message_callback(event_message);
    }

    onOpen = () => {}

    onClose = () => {}

    onError = error => {}

    /**
     * Sets the event that will handle new messages
     * @param {(event_message: PlatformEventMessage) => void} callback
     */
    setOnMessageCallback = callback => {
        if (callback?.constructor.name !== "Function" && callback?.constructor.name !== "AsyncFunction") {
            console.log("The callback must be a function, got: ", callback);
            throw new TypeError("The callback must be a function, got: ", callback);
        }

        this.#on_message_callback = callback;
    }
}

/**
 * @template T
 */
export class PlatformEventMessage {
    /**
     * Some endpoints require async processing, in which case they will return a UUID to track the event's status
     * @type {string}
     */
    #event_uuid;

    /**
     * Identify the type of effect this event has
     * @type {string}
     */
    #event_type;

    /**
     * Human readable message
     * @type {string}
     */
    #event_message;

    /**
     * Event payload encoded as a signed JWT token
     * @type {string}
     */
    #event_payload;

    /**
     * @param {PlatformEventMessageParams} param0
     * @typedef {Object} PlatformEventMessageParams
     * @property {string} uuid
     * @property {string} event_type
     * @property {string} event_message
     * @property {string} event_payload
     */
    constructor({uuid, event_type, event_message, event_payload}) {
        this.#event_uuid = uuid;
        this.#event_type = event_type;
        this.#event_message = event_message;
        this.#event_payload = event_payload;
    }

    get UUID() {
        return this.#event_uuid;
    }

    get EventType() {
        return this.#event_type;
    }

    get EventMessage() {
        return this.#event_message;
    }

    /**
     * Parses the event payload without verifying the signature
     * @returns {EventPayload<T>}
     * @template T
     * @typedef {Object} EventPayload
     * @property {Object} header
     * @property {T} payload
     */
    get EventPayload() {
        let parsed_payload = parseUnsecureJWT(this.#event_payload);

        return parsed_payload;
    }
}