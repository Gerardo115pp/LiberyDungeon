/**
 * Never import this from another store
 */
import { resetDungeonTagsStore } from "./dungeons_tags";
import dungeon_tags_clipboard from "@components/DungeonTags/MediaTagger/stores/dungeon_tags_clipboard";

/**
 * Resets dungeon tags related store.
 * @returns {void}
 */
const resetDungeonTagsRelatedStores = () => {
    resetDungeonTagsStore();
    dungeon_tags_clipboard.resetMediaTaggerDungeonTagsClipboard();
}

/**
 * Resets dungeon specific state.
 * @returns {void}
 */
const resetDungeonState = () => {
    resetDungeonTagsRelatedStores();
}

const state_cleaners = {
    resetDungeonState,
    resetDungeonTagsRelatedStores
}

export default state_cleaners;