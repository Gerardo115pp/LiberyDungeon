import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";
import generateClusterPublicTagsHotkeyContext, { cluster_public_tags_context_name } from "../TagTaxonomyComponents/cluster_public_tags_hotkeys";
import generateTagTaxonomyCreatorHotkeyContext, { tag_taxonomy_creator_context_name } from "../TagTaxonomyComponents/tag_taxonomy_creator_hotkeys";
import generateMediaTaggingsHotkeyContext, { media_taggings_hotkey_context_name } from "./sub-components/media_taggings_hotkeys";

/**
 * The name of the hotkey context for the media tagger component.
 * @type {string}
 */
export const media_tagger_tool_context_name = "media-tagger-tool";

/**
 * The actions the media tagger exposes. the media tagger may have more actions than this, but these are the ones that needs to coordinate with other components.
 */
export const media_tagger_actions = {
    WS_NAVIGATION: common_hotkey_actions.UP_DOWN_NAVIGATION,
    AD_NAVIGATION: common_hotkey_actions.LEFT_RIGHT_NAVIGATION,
    COPY_CURRENT_MEDIA_TAGS: Symbol("COPY_CURRENT_MEDIA_TAGS"),
    CHANGE_COPY_REGISTRY: Symbol("CHANGE_COPY_REGISTRY"),
    PASTE_DUNGEON_TAGS: Symbol("PASTE_DUNGEON_TAGS"),
}

/**
 * The child hotkey context of the media tagger component.
 */
export const media_tagger_child_contexts = { 
    TAG_TAXONOMY_CREATOR: tag_taxonomy_creator_context_name,
    MEDIA_TAGGINGS: media_taggings_hotkey_context_name,
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

    media_tagger_hotkeys.registerHotkeyAction(media_tagger_actions.COPY_CURRENT_MEDIA_TAGS, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_IGNORE,
        hotkey_register_params: {
            hotkey_triggers: global_hotkey_action_triggers.ITEM_YANKING,
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION
            }
        }
    });

    media_tagger_hotkeys.registerHotkeyAction(media_tagger_actions.CHANGE_COPY_REGISTRY, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_IGNORE,
        hotkey_register_params: {
            hotkey_triggers: global_hotkey_action_triggers.REGISTER_CHANGE,
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION,
            }
        }
    });

    media_tagger_hotkeys.registerHotkeyAction(media_tagger_actions.PASTE_DUNGEON_TAGS, {
        overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_IGNORE,
        hotkey_register_params: {
            hotkey_triggers: global_hotkey_action_triggers.ITEM_PASTING,
            callback: HOTKEY_NULLISH_HANDLER,
            options: {
                description: HOTKEY_NULL_DESCRIPTION
            }
        }
    });

    const tag_taxonomy_creator_context = generateTagTaxonomyCreatorHotkeyContext();
    const media_taggings_context = generateMediaTaggingsHotkeyContext();
    const cluster_public_tags_context = generateClusterPublicTagsHotkeyContext();

    tag_taxonomy_creator_context.ParentHotkeysContext = media_tagger_hotkeys;
    media_taggings_context.ParentHotkeysContext = media_tagger_hotkeys;
    cluster_public_tags_context.ParentHotkeysContext = media_tagger_hotkeys;

    media_tagger_hotkeys.addChildContext(tag_taxonomy_creator_context);
    media_tagger_hotkeys.addChildContext(media_taggings_context);
    media_tagger_hotkeys.addChildContext(cluster_public_tags_context);

    media_tagger_hotkeys.SetFinal();

    return media_tagger_hotkeys;
}

export default generateMediaTaggerHotkeyContext;