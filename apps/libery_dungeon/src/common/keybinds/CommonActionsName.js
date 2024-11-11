/**
* Defines common actions that many components have to implement.
 * @module libery_dungeon/src/common/keybinds/CommonActions
*/

import { HOTKEYS_GENERAL_GROUP, HOTKEYS_HIDDEN_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";

export const DROP_HOTKEY_CONTEXT = Symbol("DROP_HOTKEY_CONTEXT");
export const WASD_NAVIGATION = Symbol("WASD_NAVIGATION");
export const SHOW_HOTKEYS_TABLE = Symbol("SHOW_HOTKEYS_TABLE");

export const common_action_groups = {
    NAVIGATION: "<navigation>",
    CONTENT: "<content>",
    VIDEO: "<video>",
    EXPERIMENTAL: "<experimental>",
    GENERAL: `<${HOTKEYS_GENERAL_GROUP}>`,
    HIDDEN: `<${HOTKEYS_HIDDEN_GROUP}>`,
}