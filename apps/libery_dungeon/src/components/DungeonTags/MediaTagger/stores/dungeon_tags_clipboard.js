
/*=============================================
=            Properties            =
=============================================*/

    /**
     * the content type for the clipboard items
     * @constant
     * @type {string}
     */
    const dungeon_tags_content_type = "dungeon_tags/tags"

    /**
     * An emulation of the vim registries. This are register names(keyboard keys) -> some content wrapped around an object to help identify the content.
     * The content in the case of the media tagger will be a list of dungeon tags.
     * @type {Map<string, import('@app/common/interfaces/clipboard_registries').ClipboardItem<import('@models/DungeonTags').DungeonTag[]>>}
     */
    const dungeon_tags_registries = new Map();

    /**
     * The current active register.
     * @type {string}
     */
    let current_register = "\""; 

/*=====  End of Properties  ======*/


/*=============================================
=            Methods            =
=============================================*/

    /**
     * Changes the current register.
     * @param {string} register_name
     * @returns {void}
     */
    export const changeCurrentRegister = (register_name) => {
        current_register = register_name;
    }

    /**
     * Returns the contents of the current register
     * @returns {import('@models/DungeonTags').DungeonTag[] | null}
     */
    export const readRegister = () => {
        const register_contents = dungeon_tags_registries.get(current_register);

        return register_contents != null ? register_contents.content : null;
    }

    /**
     * resets the store state.
     * @returns {void}
     */
    export const resetMediaTaggerDungeonTagsClipboard = () => {
        dungeon_tags_registries.clear();
        current_register = "\"";
    }

    /**
     * Pastes content to the current register.
     * @param {import('@models/DungeonTags').DungeonTag[]} content
     * @returns {void}
     */
    export const writeOnCurrentRegister = (content) => {
        /**
         * @type {import('@common/interfaces/clipboard_registries').ClipboardItem<import('@models/DungeonTags').DungeonTag[]>}
         */
        const wrapped_content = {
            content_type: dungeon_tags_content_type,
            content: content, 
        }

        dungeon_tags_registries.set(current_register, wrapped_content);
    }

/*=====  End of Methods  ======*/