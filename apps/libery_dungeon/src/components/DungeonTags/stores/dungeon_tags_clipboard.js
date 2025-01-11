import { ClipboardContent } from '@models/Clipboard';
import { DungeonTag } from '@models/DungeonTags';

/*=============================================
=            Properties            =
=============================================*/

    /**
     * A class only intended to condense all the dungeon_tags in a list of DungeonTag objects as a json string to satisfy the ClipboardContent generic constraint.
     */
    class DungeonTagList {

        /**
         * Parses an object into a DungeonTagList.
         * @param {unknown} the_object
         * @returns {DungeonTagList | null}
         */
        static fromUnknown(the_object) {
            if (typeof the_object != 'object' || !the_object) return null;

            /**
             * @typedef {Object} DungeonTagListParams
             * @property {import('@models/DungeonTags').DungeonTagParams[]} dungeon_tags
             */

            const obj = /** @type {DungeonTagListParams} */(the_object); // sometimes i really hate typescript...

            if (!obj.dungeon_tags || !Array.isArray(obj.dungeon_tags)) return null;

            /**
             * @type {import('@models/DungeonTags').DungeonTag[]}
             */
            const dungeon_tags = [];

            for (let dungeon_tag_params of obj.dungeon_tags) {
                if (!dungeon_tag_params.id || !dungeon_tag_params.name || !dungeon_tag_params.name_taxonomy || !dungeon_tag_params.taxonomy) {
                    console.log(`In DungeonTagList.fromUnknown: after ${dungeon_tags.length} items. found an invalid dungeon tag params`, dungeon_tag_params);
                    console.warn("Stopping the parsing of the dungeon tags");
                    return null;
                }

                dungeon_tags.push(new DungeonTag(dungeon_tag_params))
            }

            return new DungeonTagList(dungeon_tags);
        }

        /**
         * A list of dungeon tags
         * @param {import('@models/DungeonTags').DungeonTag[]} dungeon_tags
         */
        #dungeon_tags;
        
        /**
         * @param {import('@models/DungeonTags').DungeonTag[]} dungeon_tags
         */
        constructor(dungeon_tags) {
            this.#dungeon_tags = dungeon_tags;
        }

        /**
         * The dungeon tags.
         * @type {import('@models/DungeonTags').DungeonTag[]}
         */
        get DungeonTags() {
            return this.#dungeon_tags;
        }

        /**
         * Converts the dungeon tags in a json string
         * @returns {Object}
         */
        toJsonSerializable() {
            /**
             * @type {import('@models/DungeonTags').DungeonTagParams[]}
             */
            const dungeon_tags_as_params = [];
            let human_readable_label = "";

            for (let h=0; h < this.#dungeon_tags.length; h++) {
                let dungeon_tag = this.#dungeon_tags[h];

                dungeon_tags_as_params.push(dungeon_tag.toParams());

                human_readable_label += String(dungeon_tag);
                if (h < this.#dungeon_tags.length - 1) {
                    human_readable_label += ",";
                }
            }

            const json_content = {
                readable_label: human_readable_label,
                dungeon_tags: dungeon_tags_as_params,
            }

            return json_content;
        }
    }

    /**
     * the content type for the clipboard items
     * @constant
     * @type {string}
     */
    const dungeon_tags_content_type = "dungeon_tags/tags"

    /**
     * An emulation of the vim registries. This are register names(keyboard keys) -> some content wrapped around an object to help identify the content.
     * The content in the case of the media tagger will be a list of dungeon tags.
     * @type {Map<string, import('@models/Clipboard').ClipboardContent<DungeonTagList>>}
     */
    const dungeon_tags_registries = new Map();

    /**
     * The default register.
     * @type {string}
     */
    const DEFAULT_DUNGEON_TAGS_CLIPBOARD_REGISTRY = "l"

    /**
     * The current active register.
     * @type {string}
     */
    let current_register = DEFAULT_DUNGEON_TAGS_CLIPBOARD_REGISTRY; 

/*=====  End of Properties  ======*/


/*=============================================
=            Methods            =
=============================================*/



    /**
     * Changes the current register.
     * @param {string} register_name
     * @returns {void}
     */
    const changeCurrentRegister = (register_name) => {
        current_register = register_name;
    }

    /**
     * Returns the name of the current registry
     * @returns {string}
     */
    const getCurrentRegister = () => {
        return current_register;
    }

    /**
     * Returns all the registries names in the dungeon tags clipboard.
     * @return {string[]}
     */
    const getAllRegistryNames = () => {
        /**
         * the names of all registries with content.
         * @type {string[]}
         */
        const registries_with_content = [];

        for (let registry_name of dungeon_tags_registries.keys()) {
            const content = dungeon_tags_registries.get(registry_name);

            if (content instanceof ClipboardContent) {
                registries_with_content.push(registry_name);
            }
        }

        return registries_with_content;
    }

    /**
     * Returns the contents of the current register
     * @returns {ClipboardContent<DungeonTagList> | null}
     */
    const readRegister = () => {
        return dungeon_tags_registries.get(current_register) || null;
    }

    /**
     * Returns the contents of a register by name
     * @param {string} register_name
     * @returns {ClipboardContent<DungeonTagList> | null}
     */
    const readRegisterByName = (register_name) => {
        return dungeon_tags_registries.get(register_name) || null;
    }

    /**
     * Attempts to read the content from the clipboard. If it can parse a ClipboardContent<DungeonTagList>, it will return it. otherwise it will return null.
     * @returns {Promise<ClipboardContent<DungeonTagList> | null>}
     */
    const readClipboard = async () => {
        const clipboard_content_any = await ClipboardContent.fromClipboard();

        if (clipboard_content_any === null) return null;

        const content_any = clipboard_content_any.Content


        const dungeon_tags_list = DungeonTagList.fromUnknown(content_any);

        if (dungeon_tags_list === null) return null;


        const clipboard_content = new ClipboardContent(dungeon_tags_content_type, dungeon_tags_list);

        dungeon_tags_registries.set(current_register, clipboard_content);


        return clipboard_content;
    }

    /**
     * resets the store state.
     * @returns {void}
     */
    const resetMediaTaggerDungeonTagsClipboard = () => {
        dungeon_tags_registries.clear();
        current_register = DEFAULT_DUNGEON_TAGS_CLIPBOARD_REGISTRY;
    }

    /**
     * Pastes content to the current register.
     * @param {import('@models/DungeonTags').DungeonTag[]} content
     * @returns {void}
     */
    const writeOnCurrentRegister = (content) => {
        /**
         * @type {import('@models/Clipboard').ClipboardContent<DungeonTagList>}
         */
        const wrapped_content = new ClipboardContent(dungeon_tags_content_type, new DungeonTagList(content));

        dungeon_tags_registries.set(current_register, wrapped_content);
    }

    /**
     * All the public methods of the dungeon tag clipboard store.
     */
    const dungeon_tags_clipboard = {
        changeCurrentRegister,
        getCurrentRegister,
        getAllRegistryNames,
        readRegisterByName,
        readRegister,
        readClipboard,
        resetMediaTaggerDungeonTagsClipboard,
        writeOnCurrentRegister,
    }

    /**
     * All the public methods of the dungeon tag clipboard store.
     */
    export default dungeon_tags_clipboard; 
/*=====  End of Methods  ======*/