/**
 * A json serialized ClipboardContent
 * @template T
* @typedef {Object} JsonClipboardContent
 * @property {string} content_type
 * @property {T} content
*/

/**
 * An interface for serializing and deserializing objects to and from the clipboard.
 * @template {{ toJsonSerializable(): Object }} T
 */
export class ClipboardContent {

    /**
     * Attempts to parse a ClipboardContent<Object> from the clipboard content. If it can't parse it, it will return null.
     * @returns {Promise<ClipboardContent<any> | null>}
     */
    static async fromClipboard() {

        const clipboard = globalThis.navigator?.clipboard;
        
        if (!clipboard) return null;

        let clipboard_content = "";

        /**
         * @type {JsonClipboardContent<Object>}
         */
        let clipboard_object;

        try {
            clipboard_content = await clipboard.readText();

            clipboard_object = JSON.parse(clipboard_content);

            if (!clipboard_object.content_type || !clipboard_object.content) {
                return null;
            }
        } catch (error) {
            console.error(error);
            return null;
        }

        // @ts-ignore
        clipboard_object.content.toJsonSerializable = () => clipboard_object.content;

        /**
         * We copy the only to properties we need from the clipboard object, this is to erase any other properties that might be present.
         * @type {JsonClipboardContent<{ toJsonSerializable(): Object }>}
         */
        const json_serializable = {
            content_type: clipboard_object.content_type,
            content: {
                toJsonSerializable: () => clipboard_object.content,
                ...clipboard_object.content
            },
        }

        return new ClipboardContent(json_serializable.content_type, json_serializable.content);
    }

    /**
     * The content type for the clipboard content. Can be used by the system to determine how to parse/handle the content.
     * @type {string}
     */
    #content_type;

    /**
     * The content of the clipboard.
     * @type {T}
     */
    #content;

    /**
     * @param {string} content_type
     * @param {T} content
     */
    constructor(content_type, content) {
        this.#content_type = content_type;
        this.#content = content;
    }

    /**
     * The content type for the clipboard content. Can be used by the system to determine how to parse/handle the content.
     * @type {string}
     */
    get ContentType() {
        return this.#content_type;
    }

    /**
     * The content.
     * @type {T}
     */
    get Content() {
        return this.#content;
    }

    /**
     * Copies this clipboard content object to the actual clipboard. Has to be called in a context with web api's available.
     * @return {Promise<void>}
     */
    async copy() {
        const clipboard = globalThis.navigator.clipboard;

        if (!clipboard) return;

        const json_serializable = this.#content.toJsonSerializable();

        /**
         * @type {JsonClipboardContent<Object>}
         */
        const json_clipboard_content = {
            content_type: this.#content_type,
            content: json_serializable
        }

        try {
            const json_content = JSON.stringify(json_clipboard_content);

            await clipboard.writeText(json_content);
        } catch (error) {
            console.log("Failed to copy Object to clipboard", json_serializable);
            console.error(error);
        }

        return;
    }
}