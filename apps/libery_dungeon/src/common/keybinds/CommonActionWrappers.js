import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
import { toggleHotkeysSheet } from "@stores/layout";
import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
import { global_hotkey_action_triggers } from "@app/config/hotkeys_config";

/**
 * Adds the common show hotkeys table action to a given hotkeys context.
 * @param {HotkeysContext} hotkeys_context
 * @returns {void}
 */
export const wrapShowHotkeysTable = (hotkeys_context) => {
    hotkeys_context.register(global_hotkey_action_triggers.TOGGLE_HOTKEYS_SHEET, toggleHotkeysSheet, { 
        description: `<${HOTKEYS_GENERAL_GROUP}>Open/Close the hotkeys cheat sheet.`
    });
}