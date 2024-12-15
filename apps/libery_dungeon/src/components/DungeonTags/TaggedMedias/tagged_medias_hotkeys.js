import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";
import generateClusterPublicTagsHotkeyContext, { cluster_public_tags_context_name } from "../TagTaxonomyComponents/cluster_public_tags_hotkeys";

/**
 * The name of the hotkey context for the tagged medias component
 * @type {string}
 */
export const tagged_medias_tool_context_name = "tagged-medias-tool";

/**
 * The actions the media tagger exposes. the media tagger may have more actions than this, but these are the ones that needs to coordinate with other components.
 */
export const tagged_medias_actions = {};

/**
 * The child hotkey context of the media tagger component.
 */
export const tagged_medias_child_contexts = {
    CLUSTER_PUBLIC_TAGS: cluster_public_tags_context_name
}

/**
 * Generates a new Component Hotkey context for the TaggedMedias component.
 * @returns {ComponentHotkeyContext}
 */
const generateTaggedMediasHotkeyContext = () => {
    const tagged_medias_hotkeys = new ComponentHotkeyContext(tagged_medias_tool_context_name);


    const child_cluster_public_tags_context = generateClusterPublicTagsHotkeyContext();

    child_cluster_public_tags_context.ParentHotkeysContext = tagged_medias_hotkeys;

    tagged_medias_hotkeys.addChildContext(child_cluster_public_tags_context);


    tagged_medias_hotkeys.SetFinal();

    return tagged_medias_hotkeys
}

export default generateTaggedMediasHotkeyContext;