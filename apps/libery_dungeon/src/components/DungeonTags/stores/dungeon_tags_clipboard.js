import { ClipboardContent } from '@models/Clipboard';
import { DungeonTag, DungeonTagList } from '@models/DungeonTags';
import { dungeon_tags_content_type } from '@app/common/content_types';

/*=============================================
=            Properties            =
=============================================*/

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