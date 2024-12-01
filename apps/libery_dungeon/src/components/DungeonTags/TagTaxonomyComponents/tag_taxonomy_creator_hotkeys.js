import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";

/**
 * The name of the hotkey context for the tag taxonomy creator component.
 * @type {string}
 */
export const tag_taxonomy_creator_context_name = 'tag-taxonomy-creator';

/**
 * Generates a component hotkey context for the media tagger component.
 * @returns {ComponentHotkeyContext}
 */
const generateTagTaxonomyCreatorHotkeyContext = () => {
    const tag_taxonomy_creator_hotkeys = new ComponentHotkeyContext(tag_taxonomy_creator_context_name);

    tag_taxonomy_creator_hotkeys.SetFinal();

    return tag_taxonomy_creator_hotkeys
}

export default generateTagTaxonomyCreatorHotkeyContext;