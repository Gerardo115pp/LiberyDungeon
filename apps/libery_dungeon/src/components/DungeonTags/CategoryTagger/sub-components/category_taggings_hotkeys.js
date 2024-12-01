import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";


/**
 * The name of the hotkey context for the category taggings component.
 * @type {string}
 */
export const category_taggings_hotkey_context_name = "category-taggings";

/**
 * Generates a component hotkey context for the category taggings.
 * @returns {ComponentHotkeyContext}
 */
const generateCategoryTaggingsHotkeyContext = () => {
    const category_taggings_context = new ComponentHotkeyContext(category_taggings_hotkey_context_name);

    category_taggings_context.SetFinal();

    return category_taggings_context;
}

export default generateCategoryTaggingsHotkeyContext;