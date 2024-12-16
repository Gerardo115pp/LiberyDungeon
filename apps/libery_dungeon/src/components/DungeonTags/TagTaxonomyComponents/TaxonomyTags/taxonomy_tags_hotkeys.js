import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";

/**
 * The actions that the taxonomy tags component can handle.
 */
export const taxonomy_tags_actions = {
    DROP_HOTKEY_CONTEXT: common_hotkey_actions.DROP_HOTKEY_CONTEXT,
    WASD_NAVIGATION: common_hotkey_actions.WASD_NAVIGATION,
    FOCUS_TAG_CREATOR: Symbol("FOCUS_TAG_CREATOR"),
    RENAME_FOCUSED_TAG: Symbol("RENAME_FOCUSED_TAG"),
    SELECT_FOCUSED_TAG: Symbol("SELECT_FOCUSED_TAG"),
    ALT_SELECT_FOCUSED_TAG: Symbol("ALT_SELECT_FOCUSED_TAG"),
    DELETE_FOCUSED_TAG: Symbol("DELETE_FOCUSED_TAG"),
    ALT_DELETE_FOCUSED_TAG: Symbol("ALT_DELETE_FOCUSED_TAG"),
}

/**
 * The name of the hotkeys context for the taxonomy tags component.
 * @type {string}
 */
export const taxonomy_tags_context_name = "taxonomy_tags";


/**
 * Generates a component hotkeys context for the taxonomy tags component.
 * @returns {ComponentHotkeyContext}
 */
const generateTaxonomyTagsHotkeysContext = () => {
    const taxonomy_tags_hotkeys = new ComponentHotkeyContext(taxonomy_tags_context_name);

    taxonomy_tags_hotkeys.registerHotkeyAction(taxonomy_tags_actions.DROP_HOTKEY_CONTEXT, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE,
        hotkey_register_params: {
            hotkey_triggers: global_hotkey_action_triggers.QUITE_CONTEXT,
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                await_execution: false,
                description: HOTKEY_NULL_DESCRIPTION,
            },
        }
    });

    taxonomy_tags_hotkeys.registerHotkeyAction(taxonomy_tags_actions.FOCUS_TAG_CREATOR, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_WRAP,
        hotkey_register_params: {
            hotkey_triggers: ["i"],
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
                await_execution: false,
                mode: "keyup"
            },
        }
    });

    taxonomy_tags_hotkeys.registerHotkeyAction(taxonomy_tags_actions.RENAME_FOCUSED_TAG, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_IGNORE,
        hotkey_register_params: {
            hotkey_triggers: global_hotkey_action_triggers.ITEM_RENAMING,
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
                mode: "keyup"
            }
        }
    });

    taxonomy_tags_hotkeys.registerHotkeyAction(taxonomy_tags_actions.SELECT_FOCUSED_TAG, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_IGNORE,
        hotkey_register_params: {
            hotkey_triggers: global_hotkey_action_triggers.ITEM_SELECTION,
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
            }
        }
    });

    taxonomy_tags_hotkeys.registerHotkeyAction(taxonomy_tags_actions.ALT_SELECT_FOCUSED_TAG, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE,
        hotkey_register_params: {
            hotkey_triggers: ["shift+e"],
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
            },
        }
    });

    taxonomy_tags_hotkeys.registerHotkeyAction(taxonomy_tags_actions.DELETE_FOCUSED_TAG, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE,
        hotkey_register_params: {
            hotkey_triggers: global_hotkey_action_triggers.ITEM_DELETION_NON_IMPERATIVE,
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
            },
        }
    });

    taxonomy_tags_hotkeys.registerHotkeyAction(taxonomy_tags_actions.ALT_DELETE_FOCUSED_TAG, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE,
        hotkey_register_params: {
            hotkey_triggers: ["shift+x"],
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
            },
        }
    });

    taxonomy_tags_hotkeys.SetFinal();

    return taxonomy_tags_hotkeys;
}

export default generateTaxonomyTagsHotkeysContext;