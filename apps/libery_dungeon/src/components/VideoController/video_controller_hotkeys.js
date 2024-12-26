import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";

/**
 * The name of the hotkey context for the video controller component. We define this context name to create a Component hotkey context, but this component actually just extends it's parent hotkey context.
 * @type {string}
 */
export const video_controller_context_name = "video_controller";

/**
 * The actions the video controller exposes. the video controller may have more actions than this, but these are the ones that needs to coordinate with other components.
 */
export const video_controller_actions = {};

/**
 * Generates a video controller hotkey context.
 * @returns {ComponentHotkeyContext}
 */
const generateVideoControllerContext = () => {
    const video_controller_context = new ComponentHotkeyContext(video_controller_context_name);

    video_controller_context.SetFinal();

    return video_controller_context;
}

export default generateVideoControllerContext;