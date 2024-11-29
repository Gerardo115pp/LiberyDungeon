import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";

/**
 * The name of the hotkey context for the cluster public tags component.
 * @type {string}
 */
export const cluster_public_tags_context_name = "active-cluster-public-tags-tool";

/**
 * The actions the cluster public tags exposes. the cluster public tags may have more actions than this, but these are the ones that needs to coordinate with other components.
 */
export const cluster_public_tags_actions = {
    WS_NAVIGATION: common_hotkey_actions.UP_DOWN_NAVIGATION,
    AD_NAVIGATION: common_hotkey_actions.LEFT_RIGHT_NAVIGATION,
}

/**
 * Generates a component hotkey context for the cluster public tags component.
 * @returns {ComponentHotkeyContext}
 */
const generateClusterPublicTagsHotkeyContext = () => {
    const cluster_public_tags_hotkeys = new ComponentHotkeyContext(cluster_public_tags_context_name);

    cluster_public_tags_hotkeys.registerHotkeyAction(cluster_public_tags_actions.WS_NAVIGATION,  {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_IGNORE,
        hotkey_register_params: {
            hotkey_triggers: ["w", "s"],
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION
            }
        }
    });

    cluster_public_tags_hotkeys.registerHotkeyAction(cluster_public_tags_actions.AD_NAVIGATION, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE,
        hotkey_register_params: {
            hotkey_triggers: [...global_hotkey_action_triggers.NAVIGATION_LEFT, ...global_hotkey_action_triggers.NAVIGATION_RIGHT],
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
            }

        }
    });

    cluster_public_tags_hotkeys.SetFinal();

    return cluster_public_tags_hotkeys;
}

export default generateClusterPublicTagsHotkeyContext;