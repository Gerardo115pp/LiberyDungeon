import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";
import generateClusterPublicTagsHotkeyContext, { cluster_public_tags_context_name } from "../TagTaxonomyComponents/cluster_public_tags_hotkeys";

/**
 * The name of the hotkey context for the media tagger component.
 * @type {string}
 */
export const media_tagger_tool_context_name = "active-media-tagger-tool";

/**
 * The actions the media tagger exposes. the media tagger may have more actions than this, but these are the ones that needs to coordinate with other components.
 */
export const media_tagger_actions = {
    WS_NAVIGATION: common_hotkey_actions.UP_DOWN_NAVIGATION,
    AD_NAVIGATION: common_hotkey_actions.LEFT_RIGHT_NAVIGATION,
}

/**
 * The child hotkey context of the media tagger component.
 */
export const media_tagger_child_contexts = {
    CLUSTER_PUBLIC_TAGS: cluster_public_tags_context_name,
}

/**
 * Generates a component hotkey context for the media tagger component.
 * @returns {ComponentHotkeyContext}
 */
const generateMediaTaggerHotkeyContext = () => {
    const media_tagger_hotkeys = new ComponentHotkeyContext(media_tagger_tool_context_name);

    media_tagger_hotkeys.registerHotkeyAction(media_tagger_actions.WS_NAVIGATION,  {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_IGNORE,
        hotkey_register_params: {
            hotkey_triggers: ["w", "s"],
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION
            }
        }
    });

    media_tagger_hotkeys.registerHotkeyAction(media_tagger_actions.AD_NAVIGATION, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE,
        hotkey_register_params: {
            hotkey_triggers: [...global_hotkey_action_triggers.NAVIGATION_LEFT, ...global_hotkey_action_triggers.NAVIGATION_RIGHT],
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
            }

        }
    });

    const cluster_public_tags_context = generateClusterPublicTagsHotkeyContext();

    cluster_public_tags_context.ParentHotkeysContext = media_tagger_hotkeys;

    media_tagger_hotkeys.addChildContext(cluster_public_tags_context);

    media_tagger_hotkeys.SetFinal();

    return media_tagger_hotkeys;
}

export default generateMediaTaggerHotkeyContext;