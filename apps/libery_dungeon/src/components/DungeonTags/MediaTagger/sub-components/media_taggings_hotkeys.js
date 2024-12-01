import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";


/**
 * The name of the hotkey context for the media taggings component.
 * @type {string}
 */
export const media_taggings_hotkey_context_name = "media-taggings";

/**
 * Generates a component hotkey context for the media taggings.
 * @returns {ComponentHotkeyContext}
 */
const generateMediaTaggingsHotkeyContext = () => {
    const media_taggings_context = new ComponentHotkeyContext(media_taggings_hotkey_context_name);

    media_taggings_context.SetFinal();

    return media_taggings_context;
}

export default generateMediaTaggingsHotkeyContext;