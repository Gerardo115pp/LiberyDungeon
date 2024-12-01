import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "@libs/LiberyHotkeys/hotkeys_consts";
import * as common_hotkey_actions  from "@common/keybinds/CommonActionsName";
import { global_hotkey_action_triggers, global_hotkey_movement_triggers } from "@config/hotkeys_config";
import generateClusterPublicTagsHotkeyContext, {cluster_public_tags_context_name} from "@components/DungeonTags/TagTaxonomyComponents/cluster_public_tags_hotkeys";
import generateTagTaxonomyCreatorHotkeyContext, {tag_taxonomy_creator_context_name} from "@components/DungeonTags/TagTaxonomyComponents/tag_taxonomy_creator_hotkeys";
import generateCategoryTaggingsHotkeyContext, { category_taggings_hotkey_context_name } from "./sub-components/category_taggings_hotkeys";

/**
 * The name of the hotkey context for the category tagger component.
 * @type {string}
 */
export const category_tagger_tool_context_name = "category-tagger-tool";

/**
 * The child hotkey context of the media tagger component.
 */
export const category_tagger_child_contexts = {
    TAG_TAXONOMY_CREATOR: tag_taxonomy_creator_context_name,
    CATEGORY_TAGGINGS: category_taggings_hotkey_context_name,
    CLUSTER_PUBLIC_TAGS: cluster_public_tags_context_name,
}

/**
 * Generates a component hotkey context for the category tagger component.
 * @returns {ComponentHotkeyContext}
 */
const generateCategoryTaggerHotkeyContext = () => {
    const category_tagger_hotkeys = new ComponentHotkeyContext(category_tagger_tool_context_name);

    const tag_taxonomy_creator_context = generateTagTaxonomyCreatorHotkeyContext();
    const category_taggings_context = generateCategoryTaggingsHotkeyContext();
    const cluster_public_tags_context = generateClusterPublicTagsHotkeyContext();

    tag_taxonomy_creator_context.ParentHotkeysContext = category_tagger_hotkeys;
    category_taggings_context.ParentHotkeysContext = category_tagger_hotkeys;
    cluster_public_tags_context.ParentHotkeysContext = category_tagger_hotkeys;

    category_tagger_hotkeys.addChildContext(tag_taxonomy_creator_context);
    category_tagger_hotkeys.addChildContext(category_taggings_context);
    category_tagger_hotkeys.addChildContext(cluster_public_tags_context);

    category_tagger_hotkeys.SetFinal();

    return category_tagger_hotkeys;
}

export default generateCategoryTaggerHotkeyContext;