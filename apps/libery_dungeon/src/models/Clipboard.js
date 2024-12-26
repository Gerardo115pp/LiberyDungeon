/**
 * An interface for serializing and deserializing objects to and from the clipboard.
 * @template {{ toJson(): string }} T
 */
export class ClipboardContent {
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

        const json_content = this.#content.toJson();

        return clipboard.writeText(json_content);
    }
}