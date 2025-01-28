import { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
import generateVideoControllerContext, { video_controller_context_name } from "@components/VideoController/video_controller_hotkeys";

/**
 * The hotkey context name of the media viewer.
 * @type {string}
 */
export const media_viewer_context_name = "media_viewer";


/**
 * The actions the media viewer exposes. the media viewer may have more actions than this, but these are the ones that needs to coordinate with other components.
 */
export const media_viewer_actions = {};

/**
 * The child hotkey context of the media tagger component.
 */
export const media_viewer_child_contexts = { 
    VIDEO_CONTROLLER: video_controller_context_name,
}

/**
 * Generates a component hotkey context for the media viewer.
 * @returns {ComponentHotkeyContext}
 */
const generateMediaViewerContext = () => {
    const media_viewer_context = new ComponentHotkeyContext(media_viewer_context_name);

    const video_controller_context = generateVideoControllerContext();

    media_viewer_context.addChildContext(video_controller_context);

    media_viewer_context.SetFinal();

    return media_viewer_context;
}

export default generateMediaViewerContext;