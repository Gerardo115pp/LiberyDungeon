<script>
    /*=============================================
    =            Imports            =
    =============================================*/
        import VideoController from "@components/VideoController/VideoController.svelte";
        import { MediaChangesManager, media_change_types } from "@models/WorkManagers";
        import MediaMovementsTool from "./sub-components/MediaMovementsTool/MediaMovementsTool.svelte";
        import { categories_tree, current_category } from "@stores/categories_tree";
        import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
        import { app_context_manager } from "@libs/AppContext/AppContextManager";
        import { getCategoryTree, category_cache } from "@models/Categories";
        import MediasGallery from "./sub-components/MediaGallery/MediasGallery.svelte";
        import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import HotkeysContext, { ComponentHotkeyContext } from "@libs/LiberyHotkeys/hotkeys_context";
        import { app_contexts } from "@libs/AppContext/app_contexts";
        import { hotkeys_sheet_visible, inCinemaMode, inDarkMode, layout_properties, navbar_hidden, toggleCinemaMode, toggleDarkMode } from "@stores/layout";
        import { onMount, onDestroy, tick } from "svelte";
        import { getMediaUrl } from "@libs/DungeonsCommunication/services_requests/media_requests";
        import { replaceState } from "$app/navigation";
        import { media_types } from "@models/Medias";
        import {  } from "./app_page_store";
        import { 
            active_media_index, 
            active_media_change,
            media_changes_manager,
            random_media_navigation,
            previous_medias,
            previous_media_index,
            skip_deleted_medias,
            auto_move_on,
            auto_move_category,
            shared_active_media,
            media_viewer_hotkeys_context_name,
            automute_enabled,
            static_next_medias,
            resetMediaViewerPageStore,
            tagged_medias_tool_mounted
        } from "@pages/MediaViewer/app_page_store";
        import { current_cluster, loadCluster } from "@stores/clusters";
        import MediaInformationPanel from "./sub-components/MediaInformation/MediaInformationPanel.svelte";
        import MediaTagger from "@components/DungeonTags/MediaTagger/MediaTagger.svelte";
        import TaggedMedias from "@components/DungeonTags/TaggedMedias/TaggedMedias.svelte";
        import { media_tagging_tool_mounted } from "./app_page_store";
        import { goto } from "$app/navigation";
        import { page } from "$app/stores";
        import { current_user_identity } from "@stores/user";
        import DiscreteFeedbackLog from "@libs/LiberyFeedback/FeedbackUI/DiscreteFeedbackLog.svelte";
        import { confirmPlatformMessage, setDiscreteFeedbackMessage } from "@libs/LiberyFeedback/lf_utils";
        import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
        import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
        import { MediaFile, MediaUploader } from "@libs/LiberyUploads/models";
        import { common_action_groups } from "@app/common/keybinds/CommonActionsName";
        import { ui_core_dungeon_references } from "@app/common/ui_references/core_ui_references";
        import { ui_pandasworld_tag_references } from "@app/common/ui_references/dungeon_tags_references";
        import generateMediaTaggerHotkeyContext from "@components/DungeonTags/MediaTagger/media_tagger_hotkeys";
        import { LEFT_RIGHT_NAVIGATION }  from "@common/keybinds/CommonActionsName";
        import {
            active_tag_content_media_index,
            tagMode_changeFilteringTags,
            tagMode_disableTagMode,
            mv_tag_mode_enabled,
            mv_tagged_content,
            mv_tag_mode_total_content,
            tagMode_setActiveMediaIndex,
            tagMode_resetTaggedContentMode,
            mv_filtering_tags,
        } from "@pages/MediaViewer/features_wrappers/media_viewer_tag_mode";
        import generateTaggedMediasHotkeyContext from "@components/DungeonTags/TaggedMedias/tagged_medias_hotkeys";
        import { linearCycleNavigationWrap } from "@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils";
        import generateMediaViewerContext, { media_viewer_child_contexts } from "./media_viewer_hotkeys";
        import { navigateToDungeonExplorer } from "@libs/NavigationUtils";
    /*=====  End of Imports  ======*/
     
    /*=============================================
    =            App Context            =
    =============================================*/
    
        const app_context = app_context_manager.CurrentContext;
        const app_page_name = "media-viewer";

        if (app_context !== app_contexts.DUNGEONS) {
            app_context_manager.setAppContext(app_contexts.DUNGEONS, app_page_name, onComponentExit);
        } else {
            app_context_manager.addOnContextExit(app_page_name, onComponentExit)
        }
    
    /*=====  End of App Context  ======*/
    
    /*=============================================
    =            Properties            =
    =============================================*/

        let global_hotkeys_manager = getHotkeysManager();
    
        /*----------  Exports  ----------*/

        /**
         * A category id passed through the url
         * @type {string}
         */
        export let url_category_id;


        /**
         * The index of the media to display
         * @type {string}
         */
        export let url_media_index;

        /*----------  References  ----------*/

            /**
             * @type {HTMLVideoElement} that displays video medias
            */
            let video_element;
        
        
        /*----------  hotkeys contexts  ----------*/

            /**
             * The media viewer component hotkey context.
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
             */
            const component_hotkey_context = generateMediaViewerContext();
        
            /**
             * The component hotkey context for the child Media tagger.
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
             */
            let media_tagger_hotkeys_context = null;

            /**
             * The component hotkey context for the child TaggedMedias.
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
             */
            let tagged_medias_hotkeys_context = null;

            /**
             * The component hotkey context for video controller component.
             * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext | null}
             */
            let vide_controller_context = component_hotkey_context.ChildHotkeysContexts.get(media_viewer_child_contexts.VIDEO_CONTROLLER) ?? null;
            if (vide_controller_context == null) {
                throw new Error("In MediaTagger, invalid component_hotkey_context: ClusterPublicTags hotkeys context was not defined as a child context of the MediaTagger context.");
            }

        /*=============================================
        =            State            =
        =============================================*/

            /**
             * The active media.
             * @type {import('@models/Medias').Media | null}
             */
            let the_active_media = null;

            /** 
             * Whether the value for active_media_index has been determined. if this is false, any value active_media_index cannot
             * be considered valid.
             * @type {boolean}
             * @default false
             */
            let active_media_index_determined = false;

            /**
             * the time query params request the current media be resumed to.
             * @type {number}
             */
            let video_resume_time = NaN;

            /**
             * @type {number}
             * on each movement event, the medias will move by media_movement_factor * media_height 
             */
            let media_movement_factor = 0.2;

            /**
             * @type {number}
             * the medias can only move an amount equal to media_movement_threshold * media_height
             */
            let media_movement_threshold = 2;

            let media_zoom = 1;
            let media_zoom_factor = 0.15;
            const MAX_MEDIA_ZOOM = 6.3;
            const MIN_MEDIA_ZOOM = 0.1;

            /**
             * Whether to keep the media zoom modifier on media change.
             * @type {boolean}
             */
            let keep_media_zoom = false;
            
            /**
             * Whether to keep the media position modifier on media change.
             * @type {boolean}
             */
            let keep_media_position = false;

            /**
             * Whether to invert the media position modifier on media change.
             * Sounds weird but it's aimed for PDF/Comic pages where when you are at the bottom
             * of the page and you go to the next page, you want to be at the top of the page.
             * @type {boolean}
             */
            let page_reader_mode = false;

            /**
             * @type {boolean}
             * whether to show the media movement manager 
             */
            let show_media_movement_manager = false;

            /**
             * whether to show the media gallery 
             * @type {boolean}
             */
            let show_media_gallery = false;

            /**
             * whether to show the media information panel 
             * @type {boolean}
             */
            let show_media_information_panel = false;

            /**
             * Whether the media viewer is currently closing. 
             *  @type {boolean}
             */
            let media_viewer_closing = false;

            /**
             * Whether to enable Cinema mode which disables user media transformations and enters a fullscreen mode also scaling the media to fit the entire available space.
             * @type {boolean} 
             */
            let cinema_mode = false;

            
            /*----------  Has Unsaved changes  ----------*/
            
                /**
                 * Whether the user has selected any medias or has any other unsaved changes.
                 * used to determine whether to attach a beforeunload.
                 *  @type {boolean}
                 */
                let has_unsaved_changes = false;
            
            /*----------  Sub components state  ----------*/
            
                
                /*----------  Media Tagger  ----------*/
                    /**
                     * The media viewer media tagger.
                     * @type {MediaTagger | null}
                     */ 
                    let the_media_tagger = null;

                    /**
                     * Whether the media tagger should be hidden. unlike media_tagging_tool_mounted, which defines whether the Media tagging tool has been mounted on them dom.
                     * This variable determines whether the already mounted media tagger, is visible and can be interacted with(pointer events).
                     * @type {boolean}
                     * @default true
                     */
                    let media_tagger_hidden = true;
                
                /*----------  Tagged Medias  ----------*/
                
                    /**
                     * The TaggedMedias component instance.
                     * @type {TaggedMedias | null}
                     */ 
                    let the_tagged_medias = null;

                    /**
                     * Whether the tagged medias component should be hidden.
                     * @type {boolean}
                     */
                    let tagged_medias_hidden = true;

                /*----------  Video Controller  ----------*/
                
                    /**
                     * Whether to auto hide the video controller.
                     * @type {boolean}
                     * @default true
                     */ 
                    let auto_hide_video_controller = true;
        
        /*=====  End of State  ======*/
        
        /*----------  Unsubscriber  ----------*/
        
            let active_media_index_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/
    
    onMount(async () => {
        if (!ensureCluster()) return;

        toggleDarkMode(true);

        if (!$layout_properties.IS_MOBILE) {
            defineDesktopKeybinds();

            defineSubComponentsHotkeyContext();

            // @ts-ignore
            window.media_viewer_hotkey_context = component_hotkey_context;
        }

        if ($current_category === null || ($current_category.uuid !== url_category_id && url_category_id !== "")) {
            let new_categories_tree = await getCategoryTree(url_category_id, current_category);
            
            if (new_categories_tree === null) {
                let variable_environment = new VariableEnvironmentContextError("In MediaViewer.onMount while trying to set a new category tree");

                variable_environment.addVariable("url_category_id", url_category_id);

                const labeled_err = new LabeledError(variable_environment, "Could not load dungeon explorer data. CategoryTree is missing.", lf_errors.PROGRAMMING_ERROR__BROKEN_STATE);

                labeled_err.alert();

                goto("/");
                return;
            }
            
            categories_tree.set(new_categories_tree);
        }

        let determined_index = await determineActiveMediaIndexInitialValue();
        setActiveMediaIndex(determined_index, false);
        active_media_index_determined = true;


        console.log("Opening media viewer")

        // Sets the active media change every time the active media index changes
        active_media_index_unsubscriber = active_media_index.subscribe(new_index => {
            if ($current_category === null) return;

            if (new_index >= $current_category.content.length || new_index < 0) {
                return;
            }

            let media_change = $media_changes_manager.getMediaChangeType($current_category.content[new_index].uuid);
            active_media_change.set(media_change);
        });

        parseQueryParams();

        $media_changes_manager.setOnChangesMade(determineHasUnsavedChanges);
        // tagged_medias_tool_mounted.set(true);
    });

    onDestroy(() => {
        onComponentExit();
        
        if (inCinemaMode()) {
            toggleMediaViewerCinemaMode();
        }
        
        toggleDarkMode(false);
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /*=============================================
        =            Keybinding            =
        =============================================*/
        
            const defineDesktopKeybinds = () => {
                if (global_hotkeys_manager == null) {
                    console.error("The global hotkeys manager is null");
                    return;
                }

                if (!global_hotkeys_manager.hasContext(component_hotkey_context.HotkeysContextName)) {
                    const hotkeys_context = component_hotkey_context.generateHotkeysContext();

                    hotkeys_context.register(["a", "d"], handleMediaNavigation, {
                        description: "<navigate>Navigate through the medias, A for previous, D for next", 
                        can_repeat: true,
                    });

                    hotkeys_context.register(["\\d g"], handleMediaPositionJump, {
                        description: "<navigate>Jump to a specific media position. example: '5 g' will jump to the 5th media."
                    });

                    hotkeys_context.register("r", toggleRandomMediaNavigationHotkey, {
                        description: "<navigation>Toggle random media navigation."
                    });

                    // Media modifiers and configurations
                        hotkeys_context.register(["w", "s"], handleMoveImageUpDown, {
                            description: "<media_modification>Move the image up or down, W for up, S for down."
                        });

                        hotkeys_context.register(["shift+a", "shift+d"], handleMediaZoom, {
                            description: "<media_modification>Zoom in or out the media, Shift+A for zoom out, Shift+D for zoom in.",
                            can_repeat: true,
                        });

                        hotkeys_context.register("alt+shift+d", e => {e.preventDefault(); applyMediaModifiersConfig(true)}, {
                            description: "<media_modification>Reset the media zoom and position."
                        });

                        hotkeys_context.register(["; z"], toggleKeepMediaZoomHotkey, {
                            description: "<viewer_configuration>Changes whether to keep the media zoom when navigating through medias."
                        });

                        hotkeys_context.register(["; p"], toggleKeepMediaPositionHotkey, {
                            description: "<viewer_configuration>Changes whether to keep the media position when navigating through medias."
                        });

                        hotkeys_context.register(["; r"], toggleReaderModeHotkey, {
                            description: "<viewer_configuration>Maintains the zoom and sets the media on the top after each media change."
                        });

                        hotkeys_context.register(["f"], handleCinemaModeToggle, {
                            description: "<viewer_configuration>Cinema mode.",
                        });

                        hotkeys_context.register(["; d"], handleDarkModeToggle, {
                            description: "<viewer_configuration>Toggle dark mode.",
                            consider_time_in_sequence: true,
                        });

                        hotkeys_context.register(["; t"], handleShowTaggedMediasTool, {
                            description: `${common_action_groups.CONTENT}Opens up the tool to filter ${ui_core_dungeon_references.MEDIA.EntityNamePlural} by ${ui_pandasworld_tag_references.TAG_TAXONOMY.EntityName} ${ui_pandasworld_tag_references.TAG.EntityNamePlural}.`
                        });
                    // End 

                    hotkeys_context.register(["q", "esc"], handleGoBack, {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Go back to the media explorer.`
                    });

                    hotkeys_context.register("i", handleShowMediaInformationPanel, {
                        description: "<media_modification>Toggle the media information panel."
                    });
                    
                    if ($current_user_identity?.canPublicContentAlter()) {
                        hotkeys_context.register(["f2"], handleMediaMovementToggle, {
                            description: "<media_modification>Toggle the media movement manager.",
                            mode: "keyup"
                        });

                        hotkeys_context.register("e", rejectMedia, {
                            description: "<media_filtration>Reject/Delete the current media."
                        });

                        hotkeys_context.register("shift+alt+e", e => skip_deleted_medias.set(!$skip_deleted_medias), {
                            description: "<navigation>Toggle skipping deleted medias."
                        });

                        hotkeys_context.register("t", handleShowMediaTaggerTool, {
                            description: `${common_action_groups.CONTENT}Opens up the tool to configure the current ${ui_core_dungeon_references.MEDIA.EntityName}'s ${ui_pandasworld_tag_references.TAG_TAXONOMY.EntityNamePlural}.`
                        });
                    }

                    hotkeys_context.register("shift+g", e => show_media_gallery = !show_media_gallery, {
                        description: "<navigation>Toggle the media gallery."
                    });

                    hotkeys_context.register("n", e => clearActiveMediaChanges(), {
                        description: "<media_modification>Clear the current media changes(delete, move).",
                    });

                    hotkeys_context.register("?", e => hotkeys_sheet_visible.set(!$hotkeys_sheet_visible), {
                        description: `<${HOTKEYS_GENERAL_GROUP}>Toggle the hotkeys cheatsheet.`,
                    });

                    global_hotkeys_manager.declareContext(component_hotkey_context.HotkeysContextName, hotkeys_context);

                    media_viewer_hotkeys_context_name.set(component_hotkey_context.HotkeysContextName);
                }

                global_hotkeys_manager.loadContext(component_hotkey_context.HotkeysContextName);            
            }
            
            const clearActiveMediaChanges = () => {
                if ($current_category == null) {
                    console.error("In MediaViewer.clearActiveMediaChanges: $current_category is null.");
                    return;
                }
                
                if ($active_media_change === media_change_types.NORMAL) return;

                const current_media = getActiveMedia();

                $media_changes_manager.clearMediaChanges(current_media.uuid);

                active_media_change.set(media_change_types.NORMAL);
            }

            /**
             * Returns the next media index that is not marked as deleted. optionally, it can skip the medias that are marked as moved if
             * the skip_moved parameter is set to true.
             * @param {number} from_index
             * @param {boolean} forward
             * @param {boolean} skip_moved
             * @returns {number}
             */
            const getNextNotDeletedMediaIndex = (from_index, forward, skip_moved=false) => {
                if ($current_category == null) {
                    throw Error("In MediaViewer.getNextNotDeletedMediaIndex: $current_category is null.");
                }
                
                let next_index = from_index;
                const displayed_medias = getDisplayedMedias();
                const media_count = displayed_medias.length;
                const max_index = media_count - 1;
                let media_change;
                let media_uuid;

                // if the current media is the last media and forward was requested, or if the current media is the first media and backward was requested, 
                // then STAY IN THE SAME INDEX
                if ((from_index === (media_count - 1) && forward) || (from_index === 0 && !forward)) return from_index;

                let stop_loop = false;
                let is_out_of_bounds = false;

                do {
                    next_index += forward ? 1 : -1;
                    
                    is_out_of_bounds = next_index < 0 || next_index >= (media_count - 1);
                    if (is_out_of_bounds) {
                        stop_loop = true;
                        break;
                    }

                    media_uuid = $current_category.content[next_index]?.uuid;
                    media_change = $media_changes_manager.getMediaChangeType(media_uuid);
                    
                    stop_loop = media_change !== media_change_types.DELETED && (!skip_moved || media_change !== media_change_types.MOVED);
                } while(!stop_loop);

                next_index = Math.max(0, Math.min(next_index, max_index));
                media_uuid = displayed_medias[next_index].uuid;
                media_change = $media_changes_manager.getMediaChangeType(media_uuid);

                return media_change === media_change_types.DELETED || (media_change === media_change_types.MOVED && skip_moved) ? from_index : next_index;
            }

            /**
             * Returns the translate transformation of the media display element considering 
             * only the translate values of the element. If for any reason the media element
             * is null then returns null.
             * @returns {TranslateTransformation | null}
             * @typedef TranslateTransformation
             * @property {number} x
             * @property {number} y
             */
            const getMediaDisplayTranslatePosition = () => {
                let position = {
                    x: 0,
                    y: 0
                }

                let media_element = getHTMLMediaElement();

                if (media_element === null) {
                    return null;
                }

                let translate_string =  media_element.style.translate;

                if (translate_string === "") {
                    return position;
                }

                if (translate_string.includes(" ")) { // Two values
                    let values = translate_string.split(" ");
                    position.x = parseInt(values[0]);
                    position.y = parseInt(values[1]);
                } else { // One value
                     // If for example, translate is "10%" that doesn't mean Y or X are not set
                     // that means they both have the value of 10%.
                    position.x = parseInt(translate_string);
                    position.y = position.x;
                }

                return position;
            }

            /**
             * Handles the navigation between medias. If the random media navigation is enabled
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMediaNavigation = async (key_event, hotkey) => {
                if ($current_category == null) {
                    console.error("In MediaViewer.handleMediaNavigation: $current_category is null.");
                    return;
                }
                
                let key_combo = hotkey.KeyCombo.toLowerCase();
                
                if ($random_media_navigation) {
                    return handleRandomMediaNavigation(key_event, key_combo);
                }

                const the_media_index = getActiveMediaIndex();
                let new_index = the_media_index;
                const max_index = getMaxMediaIndex();
                const index_increment = key_combo === "d" ? 1 : -1;

                new_index = linearCycleNavigationWrap(the_media_index, max_index, index_increment).value;


                const new_active_media = getDisplayedMediaByIndex(new_index);

                if(media_change_types.DELETED === $media_changes_manager.getMediaChangeType(new_active_media.uuid) && $skip_deleted_medias) {
                    let not_deleted_new_index = getNextNotDeletedMediaIndex(new_index, key_combo !== "a");
                    new_index = (not_deleted_new_index === new_index) ? the_media_index : not_deleted_new_index;
                }

                automoveMedia();

                await changeDisplayedMedia(new_index);

                await tick();
                
                applyMediaModifiersConfig();
            }

            /**
             * Handles the media position jump. 
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMediaPositionJump = async (key_event, hotkey) => {
                if ($current_category == null) {
                    console.error("In MediaViewer.handleMediaPositionJump: $current_category is null.");
                    return;
                }
                
                if (!hotkey.WithVimMotion) {
                    console.error("The hotkey did not contain vim motion data");
                    return;
                }

                let match_metadata = hotkey.MatchMetadata;

                if (match_metadata === null || match_metadata.MotionMatches.length < 1) {
                    console.error("The hotkey did not contain any motion matches");
                    return;
                }

                const motion_value = match_metadata.MotionMatches[0]

                let move_position = motion_value - 1; // convert to zreo based index

                const displayed_medias = getDisplayedMedias();

                move_position = Math.max(0, Math.min(move_position, displayed_medias.length - 1));

                await changeDisplayedMedia(move_position);

                await tick();

                applyMediaModifiersConfig();
            }

            /**
             * Handles the toggling of the media movement manager
             * @param {KeyboardEvent} event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMediaMovementToggle = (event, hotkey) => {
                if (global_hotkeys_manager == null) {
                    console.error("The global hotkeys manager is null");
                    return;
                }

                if ($mv_tag_mode_enabled) {
                    let labeled_err = new LabeledError("Accidental action protection", `Cannot move ${ui_core_dungeon_references.MEDIA.EntityNamePlural} to another ${ui_core_dungeon_references.CATEGORY.EntityName}. This is just to prevent an accidental action.`, lf_errors.ERR_ACCIDENTAL_ACTION_PROTECTION);
                    labeled_err.alert();

                    return;
                }

                event.preventDefault();
                show_media_movement_manager = !show_media_movement_manager;

                global_hotkeys_manager.Binder.pause();

                setTimeout(() => {
                    global_hotkeys_manager.Binder.resume();
                }, 400);
            }

            /**
             * Handles the random media navigation. If the key_combo is "a", then the previous media index will be used as the new index.
             * Called only from handleMediaNavigation
             * @param {KeyboardEvent} key_event
             * @param {string} key_combo
             */
            const handleRandomMediaNavigation = async (key_event, key_combo) => {
                if ($current_category == null) {
                    console.error("In MediaViewer.handleRandomMediaNavigation: $current_category is null.");
                    return;
                }

                const displayed_medias = getDisplayedMedias();

                const direction_forward = key_combo === "d";

                if (!direction_forward && $previous_medias.Size <= 1) {
                    let error_feedback = "Cannot go back to previous media. No previous media to go back to.";

                    const media_index = getActiveMediaIndex();

                    if (media_index) {
                        error_feedback += " Try disabling random navigation.";
                    }

                    let labeled_err = new LabeledError("Human error", error_feedback, lf_errors.ERR_HUMAN_ERROR);

                    labeled_err.alert();
                    return;
                }
                
                const current_index = getActiveMediaIndex();
                let new_index = current_index;

                // If the stack is not empty and the key_combo is "d" then has gone back to the previous medias and is now going forward again.
                const skip_random_generation = !$static_next_medias.IsEmpty() && direction_forward;

                if (skip_random_generation) {
                    console.log("Random media navigation: skipping random generation casue statci_next_medias stack is not empty.");
                    let last_media_uuid = /** @type {string} */ ($static_next_medias.Pop());

                    if (!$static_next_medias.IsEmpty()) {
                        $previous_medias.Add(last_media_uuid);
                    }
                    
                    // @ts-ignore - we just checked that the stack is not empty. so last_media_uuid is not null.
                    new_index = getMediaIndexByUUID(last_media_uuid);
                } else if (direction_forward) {
                    new_index = Math.floor(Math.random() * displayed_medias.length);

                    while (new_index === current_index) {
                        console.log("Random media navigation: same index, trying again");
                        new_index = Math.floor(Math.random() * displayed_medias.length);
                    }
                }

                const active_media = getActiveMedia();

                if (!direction_forward) { // Navigating back

                    if ($previous_medias.Size > 1) {
                        const top_media_uuid = $previous_medias.Pop();

                        if (top_media_uuid != null && top_media_uuid !== $static_next_medias.Peek()) {
                            $static_next_medias.Add(top_media_uuid);
                        }
                    }

                    let previous_media_uuid = $previous_medias.Peek();

                    if (previous_media_uuid !== null) {
                        new_index = getMediaIndexByUUID(previous_media_uuid);
                    }
                } else { // If navigating forward but not on the already static next medias
                    const next_media = getDisplayedMediaByIndex(new_index);

                    if ($previous_medias.IsEmpty()) { 
                        $previous_medias.Add(active_media.uuid);
                    }

                    if (next_media.uuid !== $previous_medias.Peek()) {
                        $previous_medias.Add(next_media.uuid);
                    }
                }

                await setActiveMediaIndex(new_index);


                // TODO: Remove these debug logs.
                // console.log("Random media navigation - previous medias: ", $previous_medias);
                // console.log(`previous_medias: ${$previous_medias.toString()}`);
                // // @ts-ignore - we can add antyhing we like to the window object. whether ts likes it or not.
                // window.previous_medias = $previous_medias;
                // console.log("Random media navigation - static next medias: ", $static_next_medias);
                // console.log(`static_next_medias: ${$static_next_medias.toString()}`);
                // // @ts-ignore - we can add antyhing we like to the window object. whether ts likes it or not.
                // window.static_next_medias = $static_next_medias;

                await tick();
                
                applyMediaModifiersConfig();
            }

            /**
             * Handles the movement of the image up or down. The movement is based on the media height and the media_movement_factor.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMoveImageUpDown = (key_event, hotkey) => {
                if (cinema_mode) return;

                let key_combo = key_event.key;
                let movement_direction = key_combo === "w" ? -1 : 1;

                let media_translate = getMediaDisplayTranslatePosition();

                if (media_translate === null) return;

                let new_Y_translate = media_translate.y + (movement_direction * (media_movement_factor * 100));

                setMediaPosition({x: 0, y: new_Y_translate});

                let feedback_message = `pos: ${new_Y_translate}%`;

                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * Handles the zooming in or out of the media. The zooming is based on the media_zoom_factor.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMediaZoom = (key_event, hotkey) => {
                if (cinema_mode) return;

                let key_combo = hotkey.KeyCombo.toLowerCase();

                let media_element = getHTMLMediaElement();

                if (media_element === null) {
                    return;
                }

                let new_media_zoom = media_zoom;

                new_media_zoom += key_combo === "shift+a" ? -media_zoom_factor : media_zoom_factor;

                if (new_media_zoom >= MIN_MEDIA_ZOOM && new_media_zoom <= MAX_MEDIA_ZOOM) {
                    media_zoom = new_media_zoom;
                    media_element.style.scale = `${media_zoom}`;
                }

                let feedback_message = `Media zoom: ${Math.trunc(media_zoom * 100)}%`;

                setDiscreteFeedbackMessage(feedback_message);
            }

            const handleShowMediaInformationPanel = () => { 
                show_media_information_panel = !show_media_information_panel;
            }

            /**
             * Toggles the media tagger tool. and if necessary loads the MediaViewer 
             * hotkey context.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleShowMediaTaggerTool = async (event, hotkey) => {
                toggleMediaTaggerTool();
            }

            /**
             * Toggles the tagged media tool.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleShowTaggedMediasTool = async (event, hotkey) => {
                toggleTaggedMediasTool();
            }

            const handleGoBack = async () => {
                if ($current_category == null) {
                    console.error("In MediaViewer.handleGoBack: $current_category is null.");
                    return;
                }

                if ($categories_tree == null) {
                    console.error("In MediaViewer.handleGoBack: $categories_tree is null.");
                    return;
                }
                
                if (global_hotkeys_manager == null) {
                    console.error("The global hotkeys manager is null");
                    return;
                }

                const changes_confirmation = await confirmMediaChanges();

                if (!changes_confirmation) {
                    console.warn("User cancelled the media viewer exit.");
                    return;
                }

                if (media_viewer_closing) {
                    console.error("Media viewer is already closing.");
                    return;
                }

                media_viewer_closing = true;

                global_hotkeys_manager.unregisterCurrentContext(); // Prevents users from using hotkeys while the page is closing

                const media_index = getActiveMediaIndex()
                const displayed_medias = getDisplayedMedias();

                // change the active media index so that onComponentExit will save the correct updated index(saves the current value of active_media_index)
                let new_active_media_index = resolveNewMediaIndex($media_changes_manager, media_index, displayed_medias);
                
                await commitCurreentMediaChanges();

                // wait to set the active media index to the new value after async operations so the user doens't see weird image switching
                await setActiveMediaIndex(new_active_media_index, false);
            
                navigateToDungeonExplorer($current_category.uuid);
            }

            /**
             * Hanldes the toggling of the dark mode
             */
            const handleDarkModeToggle = () => {
                toggleDarkMode(false); // parameter `force_enable` instead of toggling it turns it on despite its current state. so we set force_enable to false

                setDiscreteFeedbackMessage("changing lights...");
            }

            /**
             * Handles the toggling of the cinema mode. 
             */
            const handleCinemaModeToggle = () => {
                toggleMediaViewerCinemaMode();
            }

            const rejectMedia = () => {
                if ($current_category == null) {
                    console.error("In MediaViewer.rejectMedia: $current_category is null.");
                    return;
                }

                const current_media = getActiveMedia();

                if (current_media == null) {
                    console.error("In MediaViewer.rejectMedia: active_media is null.");
                    return;
                }

                if ($mv_tag_mode_enabled) {
                    let labeled_err = new LabeledError("Accidental action protection", "Cannot reject media while tagging mode is enabled.", lf_errors.ERR_ACCIDENTAL_ACTION_PROTECTION);

                    labeled_err.alert();
                    return;
                }

                if ($random_media_navigation) {
                    let labeled_err = new LabeledError("Accidental action protection", "Cannot reject media while random navigation is enabled.", lf_errors.ERR_ACCIDENTAL_ACTION_PROTECTION);

                    labeled_err.alert();
                    return;
                }

                let feedback_message;
                
                const current_media_index = getActiveMediaIndex();
                let not_deleted_media_index = current_media_index;

                if ($active_media_change !== media_change_types.DELETED) {
                    $media_changes_manager.stageMediaDeletion(current_media);
                    not_deleted_media_index = getNextNotDeletedMediaIndex(current_media_index, true);
                    feedback_message = "stage for deletion";
                } else {
                    $media_changes_manager.unstageMediaDeletion(current_media.uuid);
                    feedback_message = "unstaged for deletion";
                }
                
                if (not_deleted_media_index !== current_media_index) {
                    changeDisplayedMedia(not_deleted_media_index);
                } else {
                    active_media_change.set($media_changes_manager.getMediaChangeType(current_media.uuid));
                }

                applyMediaModifiersConfig();


                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * Toggles the keep_media_zoom configuration value.
             */
            const toggleKeepMediaZoomHotkey = () => {
                let current_media_zoom_config = keep_media_zoom;
                resetMediaZoomConfig();
                keep_media_zoom = !current_media_zoom_config;

                let feedback_message = "";

                if (keep_media_zoom) {
                    feedback_message = "Media zoom: MANTAIN";
                } else {
                    feedback_message = "Media zoom: RESET";
                }

                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * Toggles the keep_media_position configuration value.
             */
            const toggleKeepMediaPositionHotkey = () => {
                let current_media_position_config = keep_media_position;
                resetMediaPositionConfig();
                keep_media_position = !current_media_position_config;

                let feedback_message = "";

                if (keep_media_position) {
                    feedback_message = "Media position: MANTAIN";
                } else {
                    feedback_message = "Media position: RESET";
                }

                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * Toggles the reader mode which sets the media position to the top of the media after each media change.
             * And also keeps the zoom of the media.
             * @returns {void}
             */
            const toggleReaderModeHotkey = () => {
                let current_page_reader_mode = page_reader_mode;
                resetMediaPositionConfig();
                page_reader_mode = !current_page_reader_mode;

                let feedback_message = "";

                if (page_reader_mode) {
                    feedback_message = "page reader: on";
                } else {
                    feedback_message = "page reader: off";
                }

                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * Toggles random media navigation.
             * @param {KeyboardEvent} key_event
             */
            const toggleRandomMediaNavigationHotkey = key_event => {
                key_event.preventDefault();
                random_media_navigation.set(!$random_media_navigation);

                resetRandomNavigationState();

                let feedback_message = $random_media_navigation ? "random navigation: on" : "random navigation: off";

                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * Toggles the state of the media tagger and if necessary loads the
             * MediaViewer hotkey context. This last part is necessary cause the 
             * media tagger defines it's own HotkeyContext on mount and when onmounted
             * it just loads the previous context, not the media viewer specifically. 
             * so this just makes sure the media viewer recovers hotkey control.
             */
            const toggleMediaTaggerTool = () => {
                if (!media_tagger_hotkeys_context || !global_hotkeys_manager) return;

                const media_tagger_was_mounted = $media_tagging_tool_mounted;
                const media_tagger_was_hidden = media_tagger_hidden;

                media_tagger_hidden = !media_tagger_hidden;

                if (!media_tagger_was_mounted) {
                    media_tagging_tool_mounted.set(true);
                } else {

                    if (media_tagger_was_hidden && the_media_tagger != null && media_tagger_hotkeys_context) {
                        const media_tagger_last_active_context = media_tagger_hotkeys_context.getLastActiveChild();
                        let last_context_loaded = false;

                        if (media_tagger_last_active_context != null) {

                            if (media_tagger_last_active_context.hasBindingFunction()) {
                                media_tagger_last_active_context.bindContext()
                                last_context_loaded = true;
                            } else if (global_hotkeys_manager.hasContext(media_tagger_last_active_context.HotkeysContextName)) {
                                last_context_loaded = global_hotkeys_manager.loadPastContext(media_tagger_last_active_context.HotkeysContextName);
                            }

                        }

                        if (!last_context_loaded) {
                            the_media_tagger.defineDesktopKeybinds();
                        }
                    } else if (!media_tagger_was_hidden) {
                        if (global_hotkeys_manager != null && global_hotkeys_manager.ContextName != component_hotkey_context.HotkeysContextName) {
                            defineDesktopKeybinds();
                        }
                    }
                }
            }

            /**
             * Toggles the state of the tagged medias tool.
             */
            const toggleTaggedMediasTool = async () => {
                if (!$tagged_medias_tool_mounted) {
                    tagged_medias_tool_mounted.set(true);
                    tagged_medias_hidden = false;
                    return
                }

                tagged_medias_hidden = !tagged_medias_hidden;

                if (tagged_medias_hotkeys_context) {
                    tagged_medias_hotkeys_context.Active = !tagged_medias_hidden;
                }
            }

            /* ---------------------- sub-components hotkey contex ---------------------- */

                /**
                 * Defines hotkey contexts for sub components.
                 */
                const defineSubComponentsHotkeyContext = () => {
                    media_tagger_hotkeys_context = defineMediaTaggerHotkeyContext();
                    component_hotkey_context.ChildHotkeysContexts.set(media_tagger_hotkeys_context.HotkeysContextName, media_tagger_hotkeys_context);

                    tagged_medias_hotkeys_context = defineTaggedMediasHotkeyContext();
                    component_hotkey_context.ChildHotkeysContexts.set(tagged_medias_hotkeys_context.HotkeysContextName, tagged_medias_hotkeys_context);
                }

                /**
                 * defines the hotkey context for the media tagger tool component.
                 * @requires generateMediaTaggerHotkeyContext
                 * @returns {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
                 */
                const defineMediaTaggerHotkeyContext = () => {
                    const media_tagger_context = generateMediaTaggerHotkeyContext();

                    /* -------------------------- left right navigation ------------------------- */
                        const left_right_navigation = media_tagger_context.getHotkeyAction(LEFT_RIGHT_NAVIGATION);

                        if (left_right_navigation == null) {
                            throw Error("In MediaViewer.defineMediaTaggerHotkeyContext: the generated hotkey context for the media tagger component had no LEFT_RIGHT_NAVIGATION action.")
                        }

                        if (left_right_navigation.OverwriteBehavior === ComponentHotkeyContext.OVERRIDE_BEHAVIOR_REPLACE) {
                            left_right_navigation.overwriteDescription(`${common_action_groups.NAVIGATION}Navigate through the medias, A for previous, D for next`);

                            left_right_navigation.Callback = handleMediaNavigation;
                        }
                    /* ---------------------------- media tagger hide --------------------------- */
                        media_tagger_context.registerExtraHotkey({
                            hotkey_triggers: ["t"],
                            callback: handleShowMediaTaggerTool,
                            options: {
                                description: `${common_action_groups.GENERAL}Hides the media tagger without having to leave the current tagger section.`,
                            },
                        });
                    /* -------------------------------------------------------------------------- */

                    media_tagger_context.ParentHotkeysContext = component_hotkey_context;

                    return media_tagger_context;
                }

                /**
                 * defines the hotkey context for the tagged medias component.
                 * @requires generateTaggedMediasHotkeyContext
                 * @returns {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
                 */
                const defineTaggedMediasHotkeyContext = () => {
                    tagged_medias_hotkeys_context = generateTaggedMediasHotkeyContext();

                    tagged_medias_hotkeys_context.ParentHotkeysContext = component_hotkey_context;

                    return tagged_medias_hotkeys_context;
                }
        
        /*=====  End of Keybinding  ======*/

        /*=============================================
        =            Media modifiers            =
        =============================================*/

            /**
             * Applies the media modifiers configuration to the media element. 
             * If the force_reset parameter is set to true, then the media will be reset to its default position
             * and zoom despite the configuration.
             * @param {boolean} [force_reset] - whether to force the reset despite config
             * @returns {void} 
             */
            const applyMediaModifiersConfig = force_reset => {
                force_reset = force_reset ?? false;

                if (force_reset) {
                    resetMediaPositionConfig();
                    resetMediaZoomConfig();
                }

                let media_element = getHTMLMediaElement();

                if (media_element === null) {
                    console.error("Odd, could not find the media element");
                    return;
                }

                // position:
                if (!keep_media_position) {
                    media_element.style.translate = "0px 0px";
                }

                // zoom:
                if (!keep_media_zoom && !page_reader_mode) {
                    media_element.style.scale = "1";
                    media_zoom = 1;
                }

                if (page_reader_mode) {
                    toPageTop();
                }
            }

            /**
             * Sets a given media position.
             * @param {TranslateTransformation} position
            */
            const setMediaPosition = position => {
                let media_element = getHTMLMediaElement();

                if (media_element === null) {
                    console.error("Odd, could not find the media element");
                    return;
                }

                let translate_string = `${position.x}% ${position.y}%`;

                media_element.style.translate = translate_string;
            }

            /**
             * Resets only the position configurations of the media.
             */
            const resetMediaPositionConfig = () => {
                keep_media_position = false;
                page_reader_mode = false;
            }

            /**
             * Resets only the zoom configurations of the media.
             */
            const resetMediaZoomConfig = () => {
                keep_media_zoom = false;
            }

            /**
             * Sets the position of the media top of media viewer. 
             * @returns {void}
             */
            const toPageTop = () => {
                const TOP_TO_ZOOM_RATIO = 0.265640; // to get this i divided the top value by the zoom when they where at the desired position... And it worked! :D

                const approximate_top = (media_zoom * TOP_TO_ZOOM_RATIO) * 100;

                setMediaPosition({x: 0, y: approximate_top});
            }
        
        /*=====  End of Media modifiers  ======*/

        /*=============================================
        =            Media Changes            =
        =============================================*/
        
            /**
             * Commits the media changes in the `media_changes_manager`
             * updates the current category tree leaf and resets the
             * `media_changes_manager` to a new instance. 
             * ATTENTION: Also removes the beforeunload event listener 
             * @returns {Promise<void>}
             */
            const commitCurreentMediaChanges = async () => {
                if ($media_changes_manager !== null && $current_category !== null) {
                    await $media_changes_manager.commitChanges($current_category.uuid);
                    
                    if ($categories_tree !== null) {
                        await $categories_tree.updateCurrentCategory();
                    }
                }

                media_changes_manager.set(new MediaChangesManager());

                removeBeforeUnloadListener();
            }

            /*----------  Unsaved changes  ----------*/

                /**
                 * Attaches a beforeunload event listener to the window to
                 * warn the user about unsaved changes.
                 * @returns {void}
                 */
                const attachBeforeUnloadListener = () => {
                    globalThis.addEventListener('beforeunload', handleBeforeUnload);
                }

                /**
                 * Sets the has_unsaved_changes flag based on the meg_gallery_changes_manager 
                 * state.
                 * @returns {void}
                 */
                const determineHasUnsavedChanges = () => {
                    let new_unsaved_changes_state = false;

                    if ($media_changes_manager !== null) {
                        new_unsaved_changes_state = $media_changes_manager.ChangesAmount > 0;
                    }

                    if (has_unsaved_changes === new_unsaved_changes_state) return;

                    if (new_unsaved_changes_state) {
                        console.debug("In MediaViewer.setHasUnsavedChanges: Attaching beforeunload listener.");
                        attachBeforeUnloadListener();
                    } else {
                        console.debug("In MediaViewer.setHasUnsavedChanges: Removing beforeunload listener.");
                        removeBeforeUnloadListener();
                    }

                    has_unsaved_changes = new_unsaved_changes_state;
                }

                /**
                 * handles the beforeunload event to warn the user about unsaved changes.
                 * @param {BeforeUnloadEvent} event
                 */
                const handleBeforeUnload = (event) => {
                    if (!has_unsaved_changes) return;

                    event.preventDefault();
                    event.returnValue = '';  // Legacy support for some browsers.
                }
            
                /**
                 * Removes the beforeunload event listener from the window.
                 * @returns {void}
                 */
                const removeBeforeUnloadListener = () => {
                    if (globalThis.removeEventListener === undefined) return;

                    globalThis.removeEventListener('beforeunload', handleBeforeUnload);
                }
        
        /*=====  End of Media Changes  ======*/

        const automoveMedia = () => {
            if ($current_category == null) {
                console.error("In MediaViewer.automoveMedia: $current_category is null.");
                return;
            }
            
            if ($auto_move_on === false || $auto_move_category == null) return;

            let current_media = getActiveMedia();

            let current_media_change = $media_changes_manager.getMediaChangeType(current_media.uuid);

            if (current_media_change !== media_change_types.NORMAL) return;

            $media_changes_manager.stageMediaMove(current_media, $auto_move_category);
        }

        /**
         * Captures a frame from the video element as a webp and downloads it.
         */
        const captureVideoFrame = async () => {
            if (video_element == null) {
                return;
            };

            let canvas = document.createElement("canvas");
         
            canvas.width = video_element.videoWidth;
            canvas.height = video_element.videoHeight;

            let ctx = canvas.getContext("2d");

            if (ctx === null) {
                const variable_enviroment = new VariableEnvironmentContextError("In MediaViewer.captureVideoFrame while trying to get the 2d context of the canvas element");

                variable_enviroment.addVariable("canvas is HTMLCanvasElement", canvas instanceof HTMLCanvasElement);
                variable_enviroment.addVariable("canvas is null", canvas === null);
                variable_enviroment.addVariable("video_element is null", video_element === null);
                variable_enviroment.addVariable("video_element is HTMLVideoElement", video_element instanceof HTMLVideoElement);

                const labeled_err = new LabeledError(variable_enviroment, "Could not capture the video frame. Are you using an old browser?", lf_errors.ERR_UNSUPPORTED_BROWSER_FEATURE);

                labeled_err.alert();
                return;
            }

            ctx.drawImage(video_element, 0, 0, canvas.width, canvas.height);

            const active_media = getActiveMedia();

            const media_name = active_media.MediaName;
            const frame_image_name = `${media_name}_frame_${video_element.currentTime}.webp`;

            const user_wants_upload = await confirmPlatformMessage({
                message_title: "Screenshot taken",
                question_message: "Do you want to upload the screenshot?",
                auto_focus_cancel: false,
                cancel_label: "No",
                confirm_label: "Yes",
                danger_level: -1
            });

            const generated_media_mime_type = "image/webp";

            if (user_wants_upload === 1) {
                canvas.toBlob(blob => {
                    if (blob === null) {
                        const variable_enviroment = new VariableEnvironmentContextError("In MediaViewer.captureVideoFrame");

                        variable_enviroment.addVariable("blob is null", blob === null);
                        variable_enviroment.addVariable("canvas is HTMLCanvasElement", canvas instanceof HTMLCanvasElement);
                        console.dirxml(canvas);
                        console.dirxml(active_media);

                        const labeled_err = new LabeledError(variable_enviroment, "Could not capture the video frame. Are you using an old browser?", lf_errors.ERR_UNSUPPORTED_BROWSER_FEATURE);

                        labeled_err.alert();
                        return;
                    }

                    uploadGeneratedMedia(blob, frame_image_name, generated_media_mime_type);
                }, generated_media_mime_type, 1);

                return;
            }

            let dataURL = canvas.toDataURL(generated_media_mime_type);

            // download the image as a file just for testing purposes
            let a = document.createElement("a");
            a.href = dataURL;
            a.download = frame_image_name;
            a.click();

            setDiscreteFeedbackMessage(`Screenshot taken at ${video_element.currentTime} seconds.`);
        }

        /**
         * Confirms the media changes and returns whether the user has accepeted them.
         * @returns {Promise<boolean>}
         */
        const confirmMediaChanges = async () => {
            if ($media_changes_manager.ChangesAmount === 0) {
                return true;
            }

            let user_confirmation = await confirmPlatformMessage({
                message_title: "Media changes",
                question_message: "Do you want to save the changes you made to the medias?",
                auto_focus_cancel: false,
                cancel_label: "No",
                confirm_label: "Yes",
                danger_level: -1
            });

            return user_confirmation === 1;
        }

        /**
         * Determines what active_media_index should be set initially. Either the one from a url parameter, the one from the cache.
         * @requires url_category_id - to load the cached index from the cache
         * @requires url_media_index - if this is valid, then the cached index will be ignored
         * @requires $current_category
         * @requires category_cache
         * @returns {Promise<number>} - the determined index
         */
        const determineActiveMediaIndexInitialValue = async () => {
            if ($current_category == null) {
                throw Error("In MediaViewer.determineActiveMediaIndexInitialValue: $current_category is null.");
            }
            
            let determined_index = 0;
            let url_holds_index = url_media_index != null && url_media_index !== "";

            if (url_holds_index) {
                determined_index = parseInt(url_media_index);
                url_holds_index = !isNaN(determined_index) && determined_index >= 0 && determined_index < $current_category.content.length;
            } 
            
            if (!url_holds_index) {
                // Attempts to get a cached media index for the current_category, and if the result is different from the current active media index, then it sets the active media index to the cached one
                // the url media_index params has priority over the cached index, so if params.media_index is not null, then no cached index will be used

                if (category_cache === null) {
                    let variable_environment = new VariableEnvironmentContextError("In MediaViewer.determineActiveMediaIndexInitialValue while trying to get a cached media index");

                    variable_environment.addVariable("category_cache", category_cache);
                    variable_environment.addVariable("url_category_id", url_category_id);
                    variable_environment.addVariable("url_media_index", url_media_index);
                    variable_environment.addVariable("detemined_index", determined_index);
                    variable_environment.addVariable("url_holds_index", url_holds_index);

                    const labeled_err = new LabeledError(variable_environment, "Could not load dungeon explorer data. CategoryCache is missing.", lf_errors.PROGRAMMING_ERROR__BROKEN_STATE);

                    labeled_err.alert();

                    goto("/");
                    return NaN;
                }

                let cached_category_index = await category_cache.getCategoryIndex(url_category_id);

                if (cached_category_index !== determined_index && cached_category_index != null) {
                    cached_category_index = Math.max(0, Math.min(cached_category_index, $current_category.content.length - 1));
                    determined_index = cached_category_index;
                }
            }

            return determined_index;
        }

        const ensureCluster = () => {
            if ($current_cluster !== null) {
                return true;
            }

            const able_to_load_cluster = loadCluster();

            if (!able_to_load_cluster) {
                goto("/");
            }

            return able_to_load_cluster;
        }

        /**
         * Returns the display item for medias regardless of the type of media.
         * @returns {HTMLMediaElement | null}
         */
        const getHTMLMediaElement = () => {
            return document.querySelector(".mw-media-element-display");
        }
        
        /*=============================================
        =            Active media state modifiers            =
        =============================================*/

            /**
             * Changes the displayed media based on a new given index. 
             * @param {number} new_index
             */
            const changeDisplayedMedia = async new_index => {
                const displayed_media_index = getActiveMediaIndex();

                if (new_index === displayed_media_index) return;

                await setActiveMediaIndex(new_index);
            }
        
            /**
             * Returns the current active media.
             * @returns {import('@models/Medias').Media}
             */
            const getActiveMedia = () => {
                let media_index = getActiveMediaIndex();

                return getDisplayedMediaByIndex(media_index);
            }

            /**
             * Returns the medias been displayed in the media viewer.
             * @returns {import('@models/Medias').Media[]}
             */
            const getDisplayedMedias = () => {
                if ($mv_tag_mode_enabled) {
                    return $mv_tagged_content;
                }

                if ($current_category == null) {
                    throw Error("In MediaViewer.getDisplayedMedias: $current_category is null.");
                }  
                
                return $current_category.content;
            }

            /**
             * Returns a displayed media by it's index.
             * @param {number} index
             * @returns {import('@models/Medias').Media}
             */
            const getDisplayedMediaByIndex = index => {
                const displayed_medias = getDisplayedMedias();

                return displayed_medias[index];
            }

            /**
             * Returns a media index by its uuid.
             * @param {string} media_uuid
             * @returns {number}
             */
            const getMediaIndexByUUID = media_uuid => {
                const displayed_medias = getDisplayedMedias();

                return displayed_medias.findIndex(media => media.uuid === media_uuid);
            }

            /**
             * Returns the value of the active media index depending on the current display state.
             * @returns {number}
             */
            const getActiveMediaIndex = () => {
                return $mv_tag_mode_enabled ? $active_tag_content_media_index : $active_media_index;
            }

            /**
             * Returns the max value the active media index can have.
             * @returns {number}
             */
            const getMaxMediaIndex = () => {
                let max_media_index = 0;

                if ($mv_tag_mode_enabled) {
                    max_media_index = $mv_tag_mode_total_content - 1;
                } else if ($current_category != null) {
                    max_media_index = $current_category.content.length - 1;
                }

                return max_media_index;
            }

            /**
             * Recalculates the new active media index subtracting from it the amount of medias that are removed from the category(delete and moved medias) and are
             * before the current active media index. If the current active media index is one of the removed or moved medias, it will first update this index until it finds a valid one.
             * if the amount of removed medias is greater or equal to total medias, then the new active media index will be 0.
             * @param {MediaChangesManager} changes_manager
             * @param {number} current_active_index
             * @param {import('@models/Medias').Media[]} all_medias
             * @returns {number}
             */
            const resolveNewMediaIndex = (changes_manager, current_active_index, all_medias) => {
                let new_index = current_active_index;
                let removed_medias = changes_manager.DeletedMediasMap;
                let moved_medias = changes_manager.MovedMediasMap;

                let is_active_media_removed = changes_manager.getMediaChangeType(all_medias[current_active_index].uuid) !== media_change_types.NORMAL;

                if (is_active_media_removed) {
                    new_index = getNextNotDeletedMediaIndex(current_active_index, true, true);
                    if (new_index === current_active_index) {
                        return current_active_index;
                    }
                }

                for (let h = 0; h < current_active_index; h++) {
                    let media_uuid = all_medias[h].uuid;
                    let media_change = changes_manager.getMediaChangeType(media_uuid);

                    if (media_change !== media_change_types.NORMAL) {
                        new_index--;
                    }
                }

                return new_index;
            }

            const saveActiveMediaToRoute = () => {
                if ($current_category == null) {
                    console.error("In MediaViewer.saveActiveMediaToRoute: $current_category is null.");
                    return;
                }

                const media_index = getActiveMediaIndex();
                
                replaceState(`/media-viewer/${$current_category.uuid}/${media_index}`, $page.state);
            }

            /**
             * Sets the active media index.
             * @param {number} new_index
             * @param {boolean} update_route
             */
            const setActiveMediaIndex = async (new_index, update_route=true)=> {
                if (!$mv_tag_mode_enabled) {
                    active_media_index.set(new_index);

                    if (update_route) {
                        saveActiveMediaToRoute();
                    }
                } else {
                    console.log("Setting active media index in tag mode");
                    const was_set = await tagMode_setActiveMediaIndex(new_index);

                    if (!was_set) return;
                }

                updateActiveMedia();
            }

            /**
             * Sets the active media.
             * @param {import('@models/Medias').Media} new_media
             * @returns {void}
             */
            const setActiveMedia = new_media => {
                the_active_media = new_media;
                shared_active_media.set(the_active_media);
            }

            /**
             * Updates the active media reference.
             * @returns {void}
             */
            const updateActiveMedia = () => {
                setActiveMedia(getActiveMedia());
            }
        
        /*=====  End of Active media state modifiers  ======*/

        /**
         * @param {CustomEvent<import('@models/Medias').Media>} event
         */
        const handleThumbnailClick = event => {
            if ($current_category == null) {
                console.error("In MediaViewer.handleThumbnailClick: $current_category is null.");
                return;
            }
            
            const media_selected = event.detail;
            
            if (media_selected === undefined || media_selected === null) return;

            const media_index = $current_category.content.findIndex(media => media.uuid === media_selected.uuid);

            if (media_index === -1) return;

            setActiveMediaIndex(media_index);
            applyMediaModifiersConfig();
            show_media_gallery = false;
        }

        /**
         * Handles the clousure of the media tagger.
         */
        const handleMediaTaggerClose = () => {
            toggleMediaTaggerTool();
        }

        /**
         * Handles the close event of the tagged medias component.
         */
        const handleTaggedMediasClose = () => {
            tagged_medias_hidden = true;
        }

        /**
         * Handles the video ready event of the video controller
         * @param {HTMLVideoElement} a_video_element
         * @returns {void}
         */
        const handleVideoReady = a_video_element => {
            if (isNaN(video_resume_time) || !video_resume_time) return;

            a_video_element.currentTime = video_resume_time;

            video_resume_time = NaN;
        }

        /**
         * Handles the change of tag filters for medias content
         * @type {import('@components/DungeonTags/TaggedMedias/tagged_medias').MediaTagsChangedCallback}
         */
        const handleFilteringMediasChange = async tags => {
            const tag_mode_was_active = $mv_tag_mode_enabled;

            if (tags.length === 0) {
                tagMode_disableTagMode();
                updateActiveMedia();
                return;
            }
            
            const had_content = await tagMode_changeFilteringTags(tags);

            if (!had_content) {
                tagMode_disableTagMode(true);

                if (tag_mode_was_active) {
                    updateActiveMedia();
                }

                const labeled_err = new LabeledError("In MediaViewer.handleFilteringMediasChange", `No ${ui_core_dungeon_references.MEDIA.EntityNamePlural} found for the selected filter ${ui_pandasworld_tag_references.TAG.EntityNamePlural}`, lf_errors.ERR_NO_CONTENT);

                labeled_err.alert();

                return;
            }

            console.log("Filtered medias:", $mv_tagged_content);

            updateActiveMedia();
        }

        function onComponentExit() {        
            if (category_cache != null) {
                const media_index = getActiveMediaIndex();
                category_cache.addCategoryIndex(url_category_id, media_index);
            }

            resetComponentSettings();

            active_media_index_unsubscriber();

            if (global_hotkeys_manager != null) {
                global_hotkeys_manager.dropContext(component_hotkey_context.HotkeysContextName);
            }

            resetMediaViewerPageStore();
            resetRandomNavigationState();
            tagMode_resetTaggedContentMode();
        }

        /**
         * Modifies the media viewer state depending on the query params it identifies.
         */
        function parseQueryParams() {
            const query_params = new URLSearchParams(globalThis.location.search);

            const MEDIA_UUID_KEY = "media_uuid"
            const TIME_KEY = "time";
            
            const media_uuid = query_params.get(MEDIA_UUID_KEY);

            if (media_uuid != null) {
                let media_uuid_index = getMediaIndexByUUID(media_uuid);

                if (media_uuid_index === -1) {
                    console.error("In MediaViewer.parseQueryParams: media_uuid_index is -1.");
                    return;
                }

                setActiveMediaIndex(media_uuid_index);
            }

            const time_str = query_params.get(TIME_KEY);

            if (time_str !== null) {
                video_resume_time = parseFloat(time_str);                
            }
        }

        const resetComponentSettings = () => {
            setActiveMediaIndex(0, false);
            auto_move_on.set(false);
            auto_move_category.set(null);

            applyMediaModifiersConfig(true);
        }

        /**
         * Resest the random navigation state.
         * @returns {void}
         */
        const resetRandomNavigationState = () => {
            const media_index = getActiveMediaIndex()

            previous_media_index.set(media_index);
            $previous_medias.Clear();
            $static_next_medias.Clear();
        }

        /**
         * Enters or exists fullscreen mode.
         * @param {boolean} enter
         */
        const toggleFullscreen = enter => {
            if (!document.fullscreenEnabled) return; 


            if (enter) {
                if (document.fullscreenElement !== null) return;
                console.log("Entering fullscreen mode")
                document.documentElement.requestFullscreen({navigationUI: "hide"});
            } else {
                if (document.fullscreenElement === null) return;

                document.exitFullscreen();
            }
        }

        const toggleMediaViewerCinemaMode = () => {
            toggleCinemaMode();

            cinema_mode = inCinemaMode();

            navbar_hidden.set(cinema_mode);

            toggleFullscreen(cinema_mode);

            if (!inDarkMode()) {
                toggleDarkMode(true);
            }

            applyMediaModifiersConfig(true);
        }

        /**
         * Uploads a given blob to the medias service.
         * @param {Blob} media_blob
         * @param {string} media_name
         * @param {string} mime_type
         * @returns {Promise<void>}
         */
        const uploadGeneratedMedia = async (media_blob, media_name, mime_type) => {
            if ($current_category == null) {
                console.error("In MediaViewer.uploadGeneratedMedia: $current_category is null.");
                return;
            }

            if ($categories_tree == null) {
                console.error("In MediaViewer.uploadGeneratedMedia: $categories_tree is null.");
                return;
            }
            

            const blob_file = new File([media_blob], media_name, {
                type: mime_type,
                lastModified: new Date().getTime()
            });

            const media_file = new MediaFile(blob_file);
            
            const media_uploader = new MediaUploader([media_file]);

            media_uploader.onAllUploaded = async () => {
                await $categories_tree.updateCurrentCategory();
                setDiscreteFeedbackMessage(`Media<${media_name}> uploaded.`);
            }

            media_uploader.startUpload($current_category.uuid);
        }

    
    /*=====  End of Methods  ======*/

</script>

<main id="libery-dungeon-media-viewer"
    style:position="relative"
    class:cinema-mode={cinema_mode}
>
    {#if $current_category !== null && active_media_index_determined && the_active_media}
        <div id="media-wrapper">
            {#if the_active_media.type === media_types.IMAGE}
                <img
                    class="mw-media-element-display"
                    src="{the_active_media.Url}"
                    alt="displayed media"
                >
            {:else}
                <video 
                    class="mw-media-element-display"
                    bind:this={video_element}  
                    src="{the_active_media.Url}" 
                    muted={$automute_enabled}
                    autoplay 
                    loop
                >
                    <track kind="caption"/>
                </video>
            {/if}
        </div>
        {#if video_element !== undefined && video_element !== null}
            <div id="mw-video-controller-wrapper">
                <VideoController 
                    component_hotkey_context={vide_controller_context}
                    the_video_element={video_element} 
                    media_uuid={the_active_media.uuid}
                    onVideoReady={handleVideoReady}
                    bind:auto_hide={auto_hide_video_controller}
                    on:capture-frame={captureVideoFrame}
                />
            </div>
        {/if}
        <div id="ldmv-category-changes-manager">
            {#if $current_user_identity?.canPublicContentAlter()}
                <MediaMovementsTool bind:is_component_visible={show_media_movement_manager}/>
            {/if}
        </div>
        <div id="ldmv-media-information-panel-wrapper"
            class:media-viewer-tool-hidden={!show_media_information_panel}
        >
            {#if show_media_information_panel && the_active_media && $current_cluster}
                <MediaInformationPanel 
                    current_cluster_information={$current_cluster} 
                    current_media_information={the_active_media}
                    media_in_current_category={!$mv_tag_mode_enabled}
                />
            {/if}
        </div>
        <div id="ldmv-media-tagger-tool" 
            class:media-viewer-tool-hidden={media_tagger_hidden}
            class:media-tagger-hidden={media_tagger_hidden}
        >
            {#if $media_tagging_tool_mounted && media_tagger_hotkeys_context != null}
                <MediaTagger 
                    bind:this={the_media_tagger}
                    component_hotkey_context={media_tagger_hotkeys_context}
                    the_active_media={the_active_media} 
                    background_alpha={0.8}
                    getTaggableMedias={getDisplayedMedias}
                    on:close-medias-tagger={handleMediaTaggerClose}
                />
            {/if}
        </div>
        <div id="ldmv-tagged-medias-tool"
            class:media-viewer-tool-hidden={tagged_medias_hidden}
        >
            {#if $tagged_medias_tool_mounted && tagged_medias_hotkeys_context != null}
                <TaggedMedias
                    filtering_dungeon_tags={$mv_filtering_tags}
                    component_hotkey_context={tagged_medias_hotkeys_context}
                    background_alpha={0.5}
                    onFilterTagsChange={handleFilteringMediasChange}
                    onClose={handleTaggedMediasClose}
                />
            {/if}
        </div>
        {#if $current_category?.content.length > 0 && show_media_gallery}            
            <div id="ldmv-media-gallery-wrapper">
                <MediasGallery on:thumbnail-click={handleThumbnailClick} all_available_medias={$current_category.content} category_path={$current_category.FullPath}/>
            </div>
        {/if}
    {/if}
    <div id="feedback-log-positioner">
        <DiscreteFeedbackLog />
    </div>
</main>

<style>
    :global(#libery-dungeon-content:has(main#libery-dungeon-media-viewer) ) {
        margin-top: 0 !important;
    }

    #libery-dungeon-media-viewer {
        height: 100vh;
        overflow: hidden;
    }

    #media-wrapper {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        transition: all .28s ease-in-out;
    }
    
    .mw-media-element-display {
        /* max-width: 90vw; */
        max-height: 98vh;
        transform-origin: center center;
        transition: scale .33s ease-out, translate .5s ease-out;
    }

    /* -------------------------- Media viewer - tools -------------------------- */


        .media-viewer-tool-hidden {
            visibility: hidden;
            opacity: 0;
            pointer-events: none;
        }

        #mw-video-controller-wrapper {
            position: absolute;
            width: 50dvw;
            height: 110px;
            left: 50%;
            bottom: 5%;
            transform: translateX(-50%);
        }

        #ldmv-category-changes-manager {
            position: absolute;
            top: var(--navbar-height);
            right: 0;
            width: 20vw;
            height: calc(100vh - var(--navbar-height));
            z-index: var(--z-index-t-1);
        }

        #ldmv-media-information-panel-wrapper {
            position: absolute;
            container-type: size;
            top: calc(var(--navbar-height) + var(--vspacing-2));
            width: clamp(300px, 45em, 40vw);
            height: calc(100dvh - calc(var(--navbar-height) * 1.5));
            font-size: var(--font-size-1);
            left: var(--vspacing-3);
            z-index: var(--z-index-t-1);
        }

        /*=============================================
        =            Media tagger            =
        =============================================*/
        
            #ldmv-media-tagger-tool {
                position: absolute;    
                width: 40dvw;
                height: calc(calc(100dvh * 0.97) - var(--navbar-height));
                inset: var(--navbar-height) auto auto var(--common-page-inline-padding);
            }        

            #ldmv-media-tagger-tool.media-tagger-hidden {
                opacity: 0;
                pointer-events: none;
            }
        
        /*=====  End of Media tagger  ======*/

        
        /*=============================================
        =            Tagged Medias            =
        =============================================*/
        
            #ldmv-tagged-medias-tool {
                position: absolute;
                width: max(350px, 33dvw);
                height: calc(calc(100dvh * .99) - var(--navbar-height)); 
                inset: var(--navbar-height) 0 auto auto;
            }
        
        /*=====  End of Tagged Medias  ======*/
        
        
        
        

        #ldmv-media-gallery-wrapper {
            position: absolute;
            top: 30%;
            left: 15%;
            width: 35vw;
            height: 50vh;
            z-index: var(--z-index-t-2);
        }

    /* -------------------------------------------------------------------------- */
    
    /*=============================================
    =            Cinema mode            =
    =============================================*/
    
        main#libery-dungeon-media-viewer.cinema-mode {

            & #media-wrapper {
                position: static;
                display: grid;
                place-items: center;
                top: 0;
                transform: translate(0, 0);
                height: 100%;
                width: 100%;
            }                

            & video.mw-media-element-display, img.mw-media-element-display {
                width: 100%;
                height: auto;
                object-fit: contain;
                object-position: center;
                max-height: 100dvh;
            }
        }
    
    /*=====  End of Cinema mode  ======*/
    
    
    /*=============================================
    =            Feedback            =
    =============================================*/
    
        #feedback-log-positioner {
            position: fixed;
            bottom: 0;
            left: 0;
            z-index: var(--z-index-t-1);
        }
    
    
    /*=====  End of Feedback  ======*/
    
    
</style>