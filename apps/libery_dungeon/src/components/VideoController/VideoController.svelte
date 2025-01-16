<script>
    import { onMount, onDestroy, createEventDispatcher, tick } from "svelte";
    import { active_media_index, automute_enabled, previous_media_index } from "@stores/media_viewer";
    import { saveMediaWatchPoint, getMediaWatchPoint } from "@models/Metadata";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import { layout_properties } from "@stores/layout";
    import { browser } from "$app/environment";
    import { videoDurationToString, encodeVideoTime, decodeVideoTime } from "@libs/utils";
    import generateVideoControllerContext from "./video_controller_hotkeys";
    import VideoControllerSettings from "./stores/video_controller_settings";
    import { current_cluster } from "@stores/clusters";
    import VideoMomentCreator from "./sub-components/VideoMomentCreator.svelte";
    import { LabeledError } from "@libs/LiberyFeedback/lf_models";
    import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    import { confirmPlatformMessage, emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import Page from "@app/routes/+page.svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/

        /**
         * The video controller component hotkey context.
         * @type {import('@libs/LiberyHotkeys/hotkeys_context').ComponentHotkeyContext}
         */
        export let component_hotkey_context = generateVideoControllerContext();
    
        /** 
         * @type {HTMLVideoElement} the video element that will be controlled
         */
        export let the_video_element;

        /**
         * The progress track bar.
         * @type {HTMLDivElement}
         */
        let the_progress_bar_track;

        let global_hotkeys_manager = getHotkeysManager();

        /**
         * The media uuid of the video element
         * @type {string}
         */
        export let media_uuid;

        /**
         * Whether the video element has been unmounted since handleVideoElementChange was called.
         * @type {boolean}
         */
        let video_element_unmounted = true;

        /**
         * Whether the component should disappear when the mouse is not over it 
         * @type {boolean}
         */
        export let auto_hide = true;

        /**
         * The minimum time in milliseconds a video must have to save the watch progress.
         * @type {number}
         * @default 'Five minutes'
         */
        const SAVE_WATCH_PROGRESS_THRESHOLD = (5 * (60 * 1000)); // 5 minutes, written in an easy to edit format.

        /**
         * The keybinds for the video controller
         * @type {Object<string, import('@libs/LiberyHotkeys/hotkeys').HotkeyDataParams>}
         */
        const keybinds = {
            PAUSE_VIDEO: {
                key_combo: "space",
                handler: handlePausePlay,
                options: {
                    description: "Pause/Play video",
                }
            },
            TOGGLE_MUTE: {
                key_combo: "m",
                handler: handleMutedToggleHotkey,
                options: {
                    description: "Mute/Unmute video",
                }
            },
            VOLUMEN_UP: {
                key_combo: "shift+up",
                handler: handleVolumeUpHotkey,
                options: {
                    description: "Increase volume by 10%",
                }
            },
            VOLUMEN_DOWN: {
                key_combo: "shift+down",
                handler: handleVolumeDownHotkey,
                options: {
                    description: "Decrease volume by 10%",
                }
            },
            VIDEO_MOMENT_NEXT: {
                key_combo: "z d",
                handler: handleSeekNextVideoMomentHotkey,
                options: {
                    description: "Seeks to the next video moment."
                }
            },
            VIDEO_MOMENT_PREV: {
                key_combo: "z a",
                handler: handleSeekPrevVideoMomentHotkey,
                options: {
                    description: "Seeks to the previous video moment."
                }
            },
            VIDEO_MOMENT_CREATE: {
                key_combo: "z z",
                handler: handleCreateVideoMoment,
                options: {
                    description: "Create a video moment at the current time.",
                    mode: "keyup"
                }
            },
            VIDEO_MOMENT_DELETE: {
                key_combo: "z x",
                handler: handleDeleteVideoMomentHotkey,
                options: {
                    description: "Delete the current video moment.",
                }
            },
            VIDEO_MOMENT_EDIT: {
                key_combo: "z C",
                handler: handleEditVideoMomentTitleHotkey,
                options: {
                    description: "Edit the title of the current video moment.",
                    mode: "keyup"
                }
            },
            FORWARD_VIDEO: {
                key_combo: "shift+x",
                handler: handleVideoForwardPercentageHotkey,
                options: {
                    description: "Forward video 5% of the total duration(min of 5 seconds). if overflows, jumps to the start",
                }
            },
            BACKWARD_VIDEO: {
                key_combo: "x",
                handler: handleVideoBackwardPercentageHotkey,
                options: {
                    description: "Backward video 5% of the total duration(min of 5 seconds). No overflow",
                }
            },
            FORWARD_SECS_VIDEO: {
                key_combo: "shift+alt+x",
                handler: handleVideoForwardSecondsHotkey,
                options: {
                    description: "Forward video 5 seconds",
                }
            },
            BACKWARD_SECS_VIDEO: {
                key_combo: "alt+x",
                handler: handleVideoBackwardSecondsHotkey,
                options: {
                    description: "Backward video 5 seconds",
                }
            },
            SKIP_FRAME_FORWARD: {
                key_combo: "shift+`",
                handler: handleSkipFrameForwardHotkey,
                options: {
                    can_repeat: true,
                    description: "Skip frame forward",
                }
            },
            SKIP_FRAME_BACKWARD: {
                key_combo: "`",
                handler: handleSkipFrameBackwardHotkey,
                options: {
                    can_repeat: true,
                    description: "Skip frame backward",
                }
            },
            REPLAY_FROM_LAST_FRAME_SKIP: {
                key_combo: "!",
                handler: handleReplayFromLastFrameSkip,
                options: {
                    description: "Replay from the last frame skip",
                }
            },
            SEEK_MINUTE: {
                key_combo: "\\d m",
                handler: handleSeekMinute,
                options: {
                    description: "Seek to a minute in the video",
                }
            },
            SPEED_UP_VIDEO: {
                key_combo: ".",
                handler: handleSpeedUpVideoHotkey,
                options: {
                    description: "Speed up video",
                }
            },
            SLOW_DOWN_VIDEO: {
                key_combo: ",",
                handler: handleSlowDownVideoHotkey,
                options: {
                    description: "Slow down video",
                }
            },
            CAPTURE_FRAME: {
                key_combo: "v f",
                handler: handleCaptureVideoFrame,
                options: {
                    description: "Capture the current frame of the video and downloads it",
                }
            },
            PRINT_FRAME: {
                key_combo: "v i",
                handler: handlePrintVideoFrame,
                options: {
                    description: "Shows the current frame exact timestamp.",
                }
            },
            TOGGLE_AUTO_HIDE: {
                key_combo: "p",
                handler: handleTogglePlayerAutoHide,
                options: {
                    description: "Whether to keep the video controls always visible or auto hide",
                }
            }
        }
        
        /*=============================================
        =            State            =
        =============================================*/
        
            /**
             * Whether the parent component has position: absolute set, mainly used for styling 
             * @type {boolean}
             */
            let controller_floating = false;

            /**
             * Whether the controller is visible or not
             * @type {boolean}
             */
            let controller_visible = !auto_hide;

            /**
             * Whether the mouse is over the component or not
             * @type {boolean}
             */
            let mouse_over_controller = false;

            /**
             * The timeout id for the controller visibility timeout 
             * @type {number | null}
             */
            let controller_visibility_interval_id = null;

            /**
             * The opacity of the controller, managed by the controller visibility timeout 
             * @type {number}
             */
            let controller_opacity = auto_hide ? 0 : 0.8;

            /**
             * Whether the video is currently paused 
             * @type {boolean}
             */
            let video_paused = false;

            /**
             * Whether the video is muted 
             * @type {boolean}
             */
            let video_muted = true;

            /**
             * Video progress percentage 
             * @type {number}
             */    
            let video_progress = 0;

            /**
             * The video progress in a string with format "hh:mm:ss"
             * @type {string}
             */
            let video_progress_string = "00:00";

            /**
             * The video duration in a string with format "hh:mm:ss"
             * @type {string}
             */
            let video_duration_string = "00:00";

            /**
             * Whether the video metadata has been loaded or not
             * @type {boolean}
             */
            let video_metadata_loaded = false;

            /**
             * Whether the video is long enough to save watch progress.
             * @type {boolean}
             */
            let save_watch_progress = false;

            /**
             * The last timestamp that was seeked back by using the frame skip(in backward direction). used to replay from that point on user input.
             * @type {number}
             */
            let last_frame_skip_timestamp = 0;

            /**
             * Whether a MouseDown event has been triggered on the time scrubber but not a MouseUp
             * @type {boolean}
             */
            let time_scrubber_mouse_down = false;

            
            /*----------  Video Moments  ----------*/
            
                /**
                 * Whether video moment creation mode is enabled
                 * @type {boolean}
                 */ 
                let enable_video_moment_creation = false;

                /**
                 * The user saved video moments for the current video.
                 * @type {import('@models/Metadata').VideoMoment[]}
                 */
                let current_video_moments = [];

                /**
                 * A video moment to edit the name of
                 * @type {import('@models/Metadata').VideoMoment | null}
                 */
                let editing_video_moment = null;
        
        /*=====  End of State  ======*/

        let active_media_index_unsubscriber = () => {};
    
        const dispatch = createEventDispatcher();

        $: handleVideoElementChange(the_video_element);
        $: handleMediaUUIDChange(media_uuid);
    
    /*=====  End of Properties  ======*/
    
    onMount(() => {

        active_media_index_unsubscriber = active_media_index.subscribe(handleActiveMediaIndexChange);

        defineVideoControllerKeybinds();
    })

    onDestroy(() => {
        if (!browser) return;

        removeVideoElementListeners(the_video_element);

        active_media_index_unsubscriber();

        saveVideoWatchProgress();
        
        removeVideoControllerKeybinds();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /*=============================================
        =            Hotkeys            =
        =============================================*/
        
            const defineVideoControllerKeybinds = () => {
                if (global_hotkeys_manager == null) {
                    console.error("The global hotkeys manager is null");
                    return;
                }

                const parent_context = component_hotkey_context.ParentHotkeysContext;

                if (parent_context == null) return;

                if (!parent_context.HasGeneratedHotkeysContext()) {
                    console.warn("In video_controller/video_controller.defineVideoControllerKeybinds: parent has not generated a hotkey context yet. aborting");
                    return;
                }

                const parent_hotkey_context = /** @type {import('@libs/LiberyHotkeys/hotkeys_context').default} */ (parent_context.HotkeysContext);
                
                if ($layout_properties.IS_MOBILE || !browser || !global_hotkeys_manager.hasLoadedContext()) return;


                const video_controls_description_group = "<video_controls>";

                Object.values(keybinds).forEach(keybind => {
                    keybind.options.description = `${video_controls_description_group} ${keybind.options.description ?? "Empty description"}`;

                    if (!parent_hotkey_context.hasHotkeyTrigger(keybind.key_combo)) {
                        parent_hotkey_context.register(keybind.key_combo, keybind.handler, keybind.options);
                    } else {
                        console.warn(`In video_controller/video_controller.defineVideoControllerKeybinds: keybind ${keybind.key_combo} already exists in the parent context`);
                    }
                });

                if (global_hotkeys_manager.ContextName === parent_context.HotkeysContextName) {
                    global_hotkeys_manager.reloadCurrentContext();
                }
            }

            const removeVideoControllerKeybinds = () => {
                if (global_hotkeys_manager == null) {
                    console.error("The global hotkeys manager is null");
                    return;
                }


                const parent_context = component_hotkey_context.ParentHotkeysContext;

                if (parent_context == null) return;

                if (!parent_context.HasGeneratedHotkeysContext()) {
                    console.warn("In video_controller/video_controller.removeVideoControllerKeybinds: parent has not generated a hotkey context yet. aborting");
                    return;
                }

                const parent_hotkey_context = /** @type {import('@libs/LiberyHotkeys/hotkeys_context').default} */ (parent_context.HotkeysContext);


                if ($layout_properties.IS_MOBILE || !global_hotkeys_manager.hasLoadedContext()) return;

                Object.values(keybinds).forEach(keybind => {
                    if (parent_hotkey_context.hasHotkeyTrigger(keybind.key_combo)) {
                        parent_hotkey_context.unregister(keybind.key_combo, keybind.options.mode ?? "keydown");
                    } else {
                        console.warn(`In video_controller/video_controller.removeVideoControllerKeybinds: keybind ${keybind.key_combo} doesn't exist in the parent context`);
                    }
                });

                if (global_hotkeys_manager.ContextName === parent_context.HotkeysContextName) {
                    global_hotkeys_manager.reloadCurrentContext();
                }
            }

            function handleVideoBackwardPercentageHotkey() {
                let duration_skipped = skipVideoPercentage(false);

                let milliseconds_part = 0;

                if (duration_skipped < 1) {
                    milliseconds_part = Math.round((duration_skipped - Math.trunc(duration_skipped)) * 1000);
                }

                let duration_string = videoDurationToString(duration_skipped);

                let feedback_message = `-${duration_string}`;

                if (milliseconds_part > 0) {
                    feedback_message += `.${milliseconds_part}`;
                }

                setDiscreteFeedbackMessage(feedback_message);
            }

            function handleVideoBackwardSecondsHotkey() {
                skipVideoSeconds(-5);

                let feedback_message = "-5 seconds";

                setDiscreteFeedbackMessage(feedback_message);
            }

            function handleMutedToggleHotkey() {
                toggleMute();

                let feedback_message = "muted: ";

                feedback_message += video_muted ? "on" : "off";

                setDiscreteFeedbackMessage(feedback_message);
            }            

            function handlePausePlay() {
                togglePauseVideo();

                let feedback_message = video_paused ? "paused" : "playing";
                
                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * Replay from the last frame skip timestamp.
             * @requires last_frame_skip_timestamp
             * @param {KeyboardEvent} event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            function handleReplayFromLastFrameSkip(event, hotkey) {
                if (last_frame_skip_timestamp === 0) return;

                setPlaybackCurrentTime(last_frame_skip_timestamp);
            }

            function handleSkipFrameBackwardHotkey() {
                skipFrame(false);

                last_frame_skip_timestamp = the_video_element.currentTime;

                let feedback_message = "<<< frame";

                setDiscreteFeedbackMessage(feedback_message);
            }

            function handleSkipFrameForwardHotkey() {
                skipFrame(true);

                let feedback_message = " frame >>>";

                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            function handleSeekNextVideoMomentHotkey(event, hotkey) {
                const next_video_moment = getNextVideoMoment();

                if (next_video_moment === null) {
                    setDiscreteFeedbackMessage("This video has no video moments saved.");

                    return;
                }

                the_video_element.currentTime = next_video_moment.DecodedTime;

                let feedback_message = `Skipped to: ${next_video_moment.Title}`;

                setDiscreteFeedbackMessage(feedback_message);
            }

            /**
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            function handleSeekPrevVideoMomentHotkey(event, hotkey) {
                pauseVideo();

                const previous_video_moment = getPreviousVideoMoment();

                if (previous_video_moment === null) {
                    setDiscreteFeedbackMessage("This video has no video moments saved.");

                    return;
                }

                the_video_element.currentTime = previous_video_moment.DecodedTime;

                let feedback_message = `Skipped to: ${previous_video_moment.Title}`;

                setDiscreteFeedbackMessage(feedback_message);
            
            }

            /**
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            async function handleCreateVideoMoment(event, hotkey) {
                pauseVideo();
                toggleAutoHideMode(false);

                await tick();

                enable_video_moment_creation = true;
            }

            /**
             * Handles the delete video moment hotkey. Deletes the current video moment.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            async function handleDeleteVideoMomentHotkey(event, hotkey) {
                if (current_video_moments.length === 0) {
                    new LabeledError(
                        "In @components/VideoController/VideoController.svelte:handleDeleteVideoMomentHotkey",
                        "No video moments to delete",
                        lf_errors.ERR_HUMAN_ERROR,
                    ).alert();

                    return;
                }

                const current_moment = getCurrentVideoMoment();

                if (current_moment === null) {
                    const discrete_failure_message = "It's ambiguous which video moment to delete, please seek to a specific moment.";

                    setDiscreteFeedbackMessage(discrete_failure_message);

                    return;
                }

                let user_confirmation = await confirmPlatformMessage({
                    message_title: "Delete moment?",
                    question_message: `Are you sure you wanna delete '${current_moment.Title}'?`,
                    danger_level: 1
                });

                if (user_confirmation !== 1) {
                    return;
                }

                const success = await $current_cluster.deleteVideoMoment(current_moment);

                if (!success) {
                    new LabeledError(
                        "In @components/VideoController/VideoController.svelte:handleDeleteVideoMomentHotkey",
                        "Error deleting video moment",
                        lf_errors.ERR_PROCESSING_ERROR
                    ).alert();

                    return;
                }

                await loadVideoMoments(media_uuid);

                emitPlatformMessage("Deleted video moment");
            }

            /**
             * Handles the edit video moment title hotkey.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            async function handleEditVideoMomentTitleHotkey(event, hotkey) {
                const current_video_moment = getCurrentVideoMoment();

                if (current_video_moment === null) {
                    const discrete_failure_message = "It's ambiguous which video moment to edit, please seek to a specific moment.";

                    setDiscreteFeedbackMessage(discrete_failure_message);

                    return;
                }

                toggleAutoHideMode(false);

                editing_video_moment = current_video_moment;
            }
            
            function handleSpeedUpVideoHotkey() {
                let new_playback_rate = setVideoPlaybackRate(true);

                let feedback_message = `speed: ${new_playback_rate}x`;

                setDiscreteFeedbackMessage(feedback_message);
            }

            function handleSlowDownVideoHotkey() {
                let new_playback_rate = setVideoPlaybackRate(false);

                let feedback_message = `speed: ${new_playback_rate}x`;

                setDiscreteFeedbackMessage(feedback_message);
            }

            function handleCaptureVideoFrame() {
                emitCaptureVideoFrame();
            }

            function handlePrintVideoFrame() {
                let current_time = the_video_element.currentTime;

                let milliseconds_part = Math.round((current_time - Math.trunc(current_time)) * 1000);

                let feedback_message = `frame: ${videoDurationToString(current_time)}.${milliseconds_part}`;

                setDiscreteFeedbackMessage(feedback_message);
            }

            function handleTogglePlayerAutoHide() {
                let new_auto_hide_state = toggleAutoHideMode();

                let feedback_message = `auto-hide: ${new_auto_hide_state ? "on" : "off"}`;

                setDiscreteFeedbackMessage(feedback_message);
            }

            function handleVideoForwardPercentageHotkey() {
                let duration_skipped = skipVideoPercentage(true);

                let milliseconds_part = 0;

                if (duration_skipped < 1) {
                    milliseconds_part = Math.round((duration_skipped - Math.trunc(duration_skipped)) * 1000);
                }

                let duration_string = videoDurationToString(duration_skipped);

                let feedback_message = `+${duration_string}`;

                if (milliseconds_part > 0) {
                    feedback_message += `.${milliseconds_part}`;
                }

                setDiscreteFeedbackMessage(feedback_message);
            }

            function handleVideoForwardSecondsHotkey() {
                skipVideoSeconds(5);

                let feedback_message = "+5 seconds";

                setDiscreteFeedbackMessage(feedback_message);
            }
            
            function handleVolumeDownHotkey() {
                console.log("Changing volume down");
                changeVolumenBy(-0.1);
            }
    
            function handleVolumeUpHotkey() {
                changeVolumenBy(0.1);
            }
        
        /*=====  End of Hotkeys  ======*/
        
        /*=============================================
        =            Video moments            =
        =============================================*/
        
            /**
             * Creates a new video moment on the currentTime. returns whether the creation was successful
             * @param {string} moment_name
             * @returns {Promise<boolean>}
             */
            const createVideoMoment = async moment_name => {
                const video_time = encodeVideoTime(the_video_element.currentTime);

                const new_moment = await $current_cluster.createVideoMoment(
                    media_uuid,
                    video_time,
                    moment_name 
                )

                if (new_moment === null) {
                    new LabeledError(
                        "In @components/VideoController/VideoController.svelte:createVideoMoment",
                        "Error creating video moment",
                        lf_errors.ERR_PROCESSING_ERROR
                    ).alert();

                    return false;
                }

                await loadVideoMoments(media_uuid);

                return true;
            }
            
            /**
             * Gets the next video moment in relation to the currentTime of the video element. if no more video moments are found cycles and returns the first one 
             * @returns {import('@models/Metadata').VideoMoment | null}
             */
            const getNextVideoMoment = () => {
                if (current_video_moments.length === 0) return null;

                let next_video_moment = current_video_moments[0];

                const current_time = the_video_element.currentTime;

                for (let video_moment of current_video_moments) {
                    if (video_moment.isBefore(current_time)) {
                        next_video_moment = video_moment;
                        break;
                    }
                }

                return next_video_moment;
            }

            /**
             * Returns the previous video moment in relation to the currentTime of the video element. if no more video moments are found cycles and returns the last one 
             * @returns {import('@models/Metadata').VideoMoment | null}
             */
            const getPreviousVideoMoment = () => {
                if (current_video_moments.length === 0) return null;

                let previous_video_moment = current_video_moments[current_video_moments.length - 1];

                const current_time = the_video_element.currentTime;

                for (let h = current_video_moments.length - 1; h >= 0; h--) {
                    let video_moment = current_video_moments[h];
                    if (video_moment.isAfter(current_time)) {
                        previous_video_moment = video_moment;
                        break;
                    }
                }

                return previous_video_moment;
            }

            /**
             * Returns the current video moment or null if there is now video moment that is 
             * close enough to the currentTime of the_video_element.
             * @returns {import('@models/Metadata').VideoMoment | null}
             */
            const getCurrentVideoMoment = () => {
                /**
                 * @type {import('@models/Metadata').VideoMoment | null}
                 */
                let current_video_moment = null;

                const current_moment = the_video_element.currentTime;
                
                for (let video_moment of current_video_moments) {
                    if (video_moment.isEqual(current_moment, 0.02)) {
                        current_video_moment = video_moment;
                        break;
                    }
                }

                return current_video_moment;
            }

            /**
             * Handles the cancellation of the new moment creation
             * process
             * @returns {void}
             */
            const handleNewMomentCreationCancel = () => {
                enable_video_moment_creation = false;
            }           

            /**
             * Handles the name committed event from the VideoMomentCreator
             * @param {string} new_moment_name
             */
            const handleNewNameCommitted = async new_moment_name => {
                enable_video_moment_creation = false;
                toggleAutoHideMode(true);

                if (the_video_element.paused) {
                    playVideo();
                }

                const success = await createVideoMoment(new_moment_name);

                if (!success) return;

                emitPlatformMessage(`Created video moment '${new_moment_name}'`);
            }

            /**
             * Handles the cancellation of current video moment edition
             * process
             * @returns {void}
             */
            const handleCurrentMomentEditingCancel = () => {
                editing_video_moment = null;
            }           

            /**
             * Handles the new name of the current video moment.
             * @param {string} new_moment_name
             */
            const handleCurrentMomentNewName = async new_moment_name => {
                if (editing_video_moment === null) return;

                const success = await editing_video_moment.alterMoment(new_moment_name);

                if (!success) {
                    new LabeledError(
                        "In @components/VideoController/VideoController.svelte:handleCurrentMomentNewName",
                        "Error editing video moment",
                        lf_errors.ERR_PROCESSING_ERROR
                    ).alert();

                    return;
                }
                
                current_video_moments = [...current_video_moments];

                editing_video_moment = null;

                emitPlatformMessage(`Edited video moment name to '${new_moment_name}'`);            
            }

            /**
             * Handles the click event of a video moment.
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleVideoMomentClicked = event => {
                const event_target = event.currentTarget;

                if (!(event_target instanceof HTMLElement)) {
                    console.error("In @components/VideoController/VideoController.svelte:handleVideoMomentClicked: event_target is not an HTMLElement");
                    return;
                };

                const video_moment_index_str = event_target.dataset.videoMomentIndex;

                if (video_moment_index_str === undefined) {
                    console.error("In @components/VideoController/VideoController.svelte:handleVideoMomentClicked: video_moment_index_str is undefined");
                    console.log("eventTarget:", event_target);
                    console.log("currentTarget:", event.currentTarget);
                    return;
                };

                const moment_index = parseInt(video_moment_index_str);

                const moment = current_video_moments[moment_index];

                if (moment === undefined) return;

                event.stopPropagation();

                setPlaybackCurrentTime(moment.DecodedTime, false);
            }

            /**
             * Loads the video moments available for a given media uuid if any.
             * If there are any video moments, then is stores them on current_video_moments.
             * @param {string} media_uuid
             * @returns {Promise<void>}
             */
            const loadVideoMoments = async media_uuid => {
                const new_video_moments = await $current_cluster.getMediaVideoMoments(media_uuid);

                if (new_video_moments === null) return;

                console.log("Loaded video moments:", new_video_moments);

                current_video_moments = new_video_moments;
            }

            /**
             * Resets the video moments.
             * @returns {void}
             */
            const resetCurrentVideoMoments = () => {
                current_video_moments = [];
            }           


        /*=====  End of Video moments  ======*/

        /*=============================================
        =            UI Event handlers            =
        =============================================*/
        
            /**
             * Handles the video progress bar click event
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleProgressClick = (event) => {
                event.stopPropagation();

                const current_target = event.currentTarget;

                if (!(current_target instanceof Element)) return;
                
                const rect = current_target.getBoundingClientRect();
                const clickX = event.clientX - rect.left;
                const new_progress = (clickX / rect.width);
                let new_video_time =  new_progress * the_video_element.duration;
                
                setPlaybackCurrentTime(new_video_time);
            };       

            /**
             * Handles the mouse movement on the entire document.
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleMouseMovement = (event) => {
                setControllerHiddenTimeout();

                handleTimeScrubberDragging(event)
            }

            /**
             * handles the mouse up event on the controller element.
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleControllerMouseUp = event => {
                if (time_scrubber_mouse_down) {
                    handleTimeScrubberMouseUp(event);
                }
            }

            /**
             * Handles the controller mouseenter event.
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleControllerMouseEnter = event => {
                mouse_over_controller = true;
            }

            /**
             * Handles the controller mouseleave event.
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleControllerMouseLeave = event => {
                mouse_over_controller = false;
            }

            /**
             * Handles the mouse down event on the time scrubber.
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleTimeScrubberMouseDown = event => {
                event.stopPropagation();

                time_scrubber_mouse_down = true;

                pauseVideo();
            }

            /**
             * Handles the mouse up event on the time scrubber.
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleTimeScrubberMouseUp = event => {
                event.stopPropagation();

                time_scrubber_mouse_down = false;

                playVideo();
            }

            /**
             * Handles the dragging of the time scrubber. If the time_scrubber_mouse_down variable
             * is true. otherwise it returns immediately
             * @param {MouseEvent} event
             * @returns {void}
             */
            const handleTimeScrubberDragging = event => {
                if (!time_scrubber_mouse_down) return;

                if (!mouse_over_controller) {
                    time_scrubber_mouse_down = false;
                    return;
                }

                const time_line_point = getProgressBarTimelinePoint(event.clientX, false);
            
                setPlaybackCurrentTime(time_line_point, false);

                updateVideoProgress();
            }
        
        /*=====  End of UI Event handlers  ======*/

        /**
         * Attaches necessary event listeners to the_video_element.
         * @param {HTMLVideoElement} new_video_element
         * @returns {void}
         */
        const attachVideoElementEventListeners = (new_video_element) => {
            if (new_video_element === null) return;

            new_video_element.addEventListener("timeupdate", handleVideoTimeUpdate);
            new_video_element.addEventListener("durationchange", handleVideoDurationChange);
            new_video_element.addEventListener("loadedmetadata", handleVideoMetadataLoaded);
            new_video_element.addEventListener("playing", handleVideoPlaying);
            window.addEventListener("beforeunload", saveVideoWatchProgress);
        }

        /**
         * Changes the volume of the video by the given amount
         * @param {number} amount the amount to change the volume by
         */
        const changeVolumenBy = (amount) => {
            let new_volume = Math.min(1, Math.max(0, the_video_element.volume + amount));
            
            the_video_element.volume = new_volume;

            let feedback_message = `volume: ${Math.round(new_volume * 100)}%`;

            if (the_video_element.muted) {
                feedback_message += " (muted)";
            }

            setDiscreteFeedbackMessage(feedback_message);
        }

        /**
         * Convert timeline point -> percentage.
         * @param {number} timeline_point
         * @returns {number} 
         */
        const convertTimelinePointToPercentage = timeline_point => {
            return (timeline_point / the_video_element.duration) * 100;
        }

        /**
         * Fetches the watch progress of the video.
         * @this the_video_element
         * @returns {Promise<void>}
         */
        async function getWatchProgress() {
            if (media_uuid == null || this.duration * 1000 < SAVE_WATCH_PROGRESS_THRESHOLD) return;

            let watch_progress = await getMediaWatchPoint(media_uuid);

            if (watch_progress == null) return;

            this.currentTime = decodeVideoTime(watch_progress);
        }

        /**
         * Returns the video_element currentTime in a string format "hh:mm:ss"
         */
        function getVideoCurrentProgressString() {
            return videoDurationToString(the_video_element.currentTime);
        }

        /**
         * Returns the time represented by a point in the progress bar.
         * @param {number} x_point - the x coordinate of the point.
         * @param {boolean} normalized whether the x_point is already relative to the progress bar rect(true) or requires normalization(false)
         * @returns {number}
         */
        function getProgressBarTimelinePoint(x_point, normalized) {
            if (!normalized && the_progress_bar_track === undefined) {
                console.error("In @components/VideoController/VideoController.svelte:getProgressBarTimelinePoint: the_progress_bar_track is undefined and the given point is not normalized. Cannot proceed.");
                return 0;
            }

            const progress_bar_rect = the_progress_bar_track.getBoundingClientRect();

            if (!normalized) {
                x_point = Math.min(x_point, (progress_bar_rect.x + progress_bar_rect.width));

                x_point = Math.max(0, (x_point - progress_bar_rect.x));
            }

            const point_timeline_percentage = (x_point / progress_bar_rect.width);

            return the_video_element.duration * point_timeline_percentage;
        }

        const handleActiveMediaIndexChange = () => {
            video_metadata_loaded = false;

            saveVideoWatchProgress();
        }
        
        /**
         * Applies the video settings to the video element when ever the video changes.
         * @param {HTMLVideoElement} video_element
         * @returns {void}
         */
        const handleVideoElementLoaded = (video_element) => {
            if (VideoControllerSettings.shouldPreservePlaybackSpeed()) {
                video_element.playbackRate = VideoControllerSettings.getPlaybackSpeed();
            } else {
                VideoControllerSettings.setPlaybackSpeed(1);
            }

            video_paused = the_video_element.paused;
        
            loadVideoMoments(media_uuid)
        }

        /**
         * Handles the changes on the reference to the the_video_element variable.
         * @param {HTMLVideoElement | null} new_video_element
         * @returns {void}
         */
        function handleVideoElementChange(new_video_element) {
            if (new_video_element === null) {
                removeVideoElementListeners(null);
                video_element_unmounted = true;
                return;
            }

            if (video_element_unmounted) {
                attachVideoElementEventListeners(new_video_element);
            }

            video_muted = the_video_element.muted;
        }

        /**
         * handles the change of the media uuid.
         * @param {string} new_media_uuid
         * @returns {void}
         */
        function handleMediaUUIDChange(new_media_uuid) {
            resetCurrentVideoMoments();
        }

        /**
         * Handles the playing event from the video downloader
         * @returns {void}
         */
        function handleVideoPlaying() {
            video_paused = false;
        }

        /**
         * Handles the SeekMinute keybind.
         * @param {KeyboardEvent} event
         * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
         * @returns {void}
         */
        function handleSeekMinute(event, hotkey) {
            if (!hotkey.WithVimMotion && hotkey.HasMatch) return;

            let minutes_in_video = Math.trunc(the_video_element.duration / 60);
            if (isNaN(minutes_in_video)) return;

            let minute_to_seek = hotkey.MatchMetadata?.MotionMatches[0];

            if (minute_to_seek == null) return;

            minute_to_seek = Math.min(minutes_in_video, Math.max(0, minute_to_seek));

            the_video_element.currentTime = minute_to_seek * 60;
        }
        
        /**
         * Handles the video element DurationChange event
         * @this the_video_element
         * @returns {void}
         */
        function handleVideoDurationChange() {
            setSaveWatchProgress.call(this);
        }

        /**
         * Handles the video element timeupdate event   
         * @this the_video_element
         * @returns {void}
         */
        function handleVideoTimeUpdate() {
            video_progress_string = videoDurationToString(this.currentTime);

            updateVideoProgress.call(this);
        }

        /**
         * Handles the video element loadedmetadata event
         * @this the_video_element
         * @returns {void}
         */
        function handleVideoMetadataLoaded() {
            video_metadata_loaded = true;
            video_duration_string = videoDurationToString(this.duration);

            // @ts-ignore
            getWatchProgress.call(this);

            handleVideoElementLoaded(this);
        }

        // TODO: create a proper visibility controller for mobile, although this more or less works, it's due to pure black magic and also once the controller appears it doesn't disappear ever again on unless the media viewer 
        // unmounts and mounts again the video controller (like when the media changes to an from video to image and back to video)
        const handleControllerVisibility = () => {
            if (mouse_over_controller || !auto_hide) return;


            controller_opacity = Math.max(0, controller_opacity - 0.5);
            controller_visible = controller_opacity > 0;

            if (!controller_visible && controller_visibility_interval_id !== null) {
                window.clearInterval(controller_visibility_interval_id);
                controller_visibility_interval_id = null;
            }
        }

        // Hotfix to the controller not disappearing on mobile. we check if the mouse(the finger) touches the controller but not a button, if so we hide the controller
        /**
         * Handles the touch event on the controller
         * @param {TouchEvent} event
         */
        const handleControllerTouch = (event) => {
            if (event.target === event.currentTarget)

            mouse_over_controller = false;
        }



        function emitCaptureVideoFrame() {
            dispatch("capture-frame");            
        }

        function togglePauseVideo() {
            video_paused = !the_video_element.paused;
            
            if (video_paused) {
                the_video_element.pause();
            } else {
                the_video_element.play();
            }
        }

        function pauseVideo() {
            video_paused = true;
            the_video_element.pause();
        }

        function playVideo() {
            video_paused = false;
            the_video_element.play();
        }

        /**
         * Removes the event listeners from the given video element and the global context.
         * @param {HTMLVideoElement | null} new_video_element
         * @returns {void}
         */
        const removeVideoElementListeners = (new_video_element) => {
            if (new_video_element instanceof HTMLVideoElement) {

                the_video_element.removeEventListener("timeupdate", handleVideoTimeUpdate);
                the_video_element.removeEventListener("durationchange", handleVideoDurationChange);
                the_video_element.removeEventListener("loadedmetadata", handleVideoMetadataLoaded);
                the_video_element.removeEventListener("playing", handleVideoPlaying);
            }

            window.removeEventListener("beforeunload", saveVideoWatchProgress);
        }

        /**
         * Changes the video playback rate increasing(true) or decreasing(false) it by 0.25
         * depending on the value of `increase`. returns the new playback rate.
         * @param {boolean} increase whether to increase or decrease the playback rate
         * @returns {number}
         */
        function setVideoPlaybackRate(increase) {
            const step_diff = 0.25;

            let step = increase ? step_diff : -step_diff;

            let new_playback_rate = Math.min(2, Math.max(0.25, the_video_element.playbackRate + step));

            the_video_element.playbackRate = new_playback_rate;
            VideoControllerSettings.setPlaybackSpeed(new_playback_rate);

            return new_playback_rate;
        }

        /**
         * Sets the video.currentTime to the given time clamping it to the video duration and 0.
         * if overflow_allowed is true, and, for example the new_current_time is -5, the currentTime will be set to
         * the video duration - 5 seconds. but if the current_time is video.duration + 5, then the duration
         * will be set to 0, not 5.
         * @param {number} new_current_time the new duration(in seconds) to set the video to
         * @param {boolean} overflow_allowed whether the duration can overflow the video duration
         * @requires video_element
         * @returns {void}
         */
        function setPlaybackCurrentTime(new_current_time, overflow_allowed = false) {
            let clamped_current_time = Math.min(the_video_element.duration, Math.max(0, new_current_time));
            
            if (overflow_allowed && clamped_current_time !== new_current_time) {
                clamped_current_time = new_current_time < 0 ? the_video_element.duration + new_current_time : 0;
            }

            the_video_element.currentTime = clamped_current_time;
        }

        /**
         * Sets save_watch_progress based on the video duration.
         * @this the_video_element
         * @returns {void}
         */
        function setSaveWatchProgress() {
            let video_duration_ms = encodeVideoTime(this.duration);
            save_watch_progress = video_duration_ms >= SAVE_WATCH_PROGRESS_THRESHOLD;
        }

        /**
         * Saves the video watch progress
         * @returns {void}
         */
        const saveVideoWatchProgress = () => {
            if (!save_watch_progress || media_uuid == null) return;

            let watch_progress = encodeVideoTime(the_video_element.currentTime);

            saveMediaWatchPoint(media_uuid, watch_progress);
        }

        /**
         * Skips video by 5%, if forward is true, skips forward, if false, skips backward.
         * Returns the percentage skipped as a positive number regardless of the direction.
         * @param {boolean} forward whether to skip forward or backward
         * @returns {number}
         */
        function skipVideoPercentage(forward) {
            if (isNaN(the_video_element.duration)) return 0; // If the video duration metadata hasn't been loaded, video.duration will be NaN. so we can't skip in percentage.

            let step = forward ? 1 : -1;

            let video_percentage = the_video_element.duration * 0.05;

            let new_time = the_video_element.currentTime + (step * video_percentage);

            setPlaybackCurrentTime(new_time, forward);

            return video_percentage;
        }

        /**
         * Skips video by a given amount of seconds.
         * @param {number} seconds
         */
        function skipVideoSeconds(seconds) {
            let direction_forward = seconds > 0;
            let new_time = the_video_element.currentTime + seconds;

            setPlaybackCurrentTime(new_time, direction_forward);
        }

        /**
         * Skips one frame forward. because of video buffering and load times, it also needs to pause the video to get the best results, but even
         * then it's not perfect since there is no way to know the user-agent's frame rate, so it assumes a static value which can be modified. the
         * default value is 15fps. Larger values give greater control but can also result in the user not noticing any change.
         * @param {boolean} forward whether to skip forward or backward
         * @requires video_element
         */
        function skipFrame(forward, assumed_fps = 30) {
            let frame_duration = 1 / assumed_fps; 
            let step = forward ? frame_duration : -frame_duration;
            let new_time = the_video_element.currentTime + step;

            if (!the_video_element.paused) {
                the_video_element.pause();
            }

            setPlaybackCurrentTime(new_time, true);
        }

        /**
         * Hides the controller after a given amount of time
         * @param {number} [hide_delay]
         */
        const setControllerHiddenTimeout = (hide_delay = 300) => {
            if (controller_visibility_interval_id === null) {
                controller_visibility_interval_id = window.setInterval(handleControllerVisibility, hide_delay);
            }

            controller_opacity = 1;
            controller_visible = true;
        }

        function toggleMute() {
            video_muted = !the_video_element.muted;

            automute_enabled.set(video_muted); // Muted value of the video element is reactive to this store
        }

        /**
         * toggles the auto hide mode. returns it's current state.
         * @param {boolean} [force_state]
         * @returns {boolean}
         */
        const toggleAutoHideMode = force_state => {
            if (force_state === undefined) {
                force_state = !auto_hide;
            }

            auto_hide = force_state;

            if (auto_hide) {
                setControllerHiddenTimeout();
            } else {
                controller_visible = true;
                controller_opacity = 1;
            }

            return auto_hide;
        }

        /**
         * Updates the video progress percentage
         * @modifies {video_progress}
         * @returns {void}
         */
        const updateVideoProgress = () => {
            if (isNaN(the_video_element.duration)) return;

            video_progress = convertTimelinePointToPercentage(the_video_element.currentTime);
        }

    /*=====  End of Methods  ======*/
        
</script>

<svelte:document on:mousemove={handleMouseMovement} />
<div id="libery-video-controller" 
    role="group" 
    aria-label="Video controls" 
    class="libery-dungeon-window"
    class:adebug={false}  
    class:time-scrubber-dragging={time_scrubber_mouse_down}
    style:opacity={controller_opacity}
    style:visibility={controller_visible ? "visible" : "hidden"}
    on:mouseenter={handleControllerMouseEnter}
    on:mouseleave={handleControllerMouseLeave}
    on:mouseup={handleControllerMouseUp}
    on:touchstart={handleControllerTouch}
>
    {#if enable_video_moment_creation}
        <div id="lvc-new-video-moment-wrapper">
            <VideoMomentCreator 
                onNameCommitted={handleNewNameCommitted}
                onCancel={handleNewMomentCreationCancel}
            />
        </div>
    {:else if editing_video_moment !== null}
        <div id="lvc-new-video-moment-wrapper">
            <VideoMomentCreator 
                starting_name={editing_video_moment.Title}
                creator_label="New name"
                onNameCommitted={handleCurrentMomentNewName}
                onCancel={handleCurrentMomentEditingCancel}
            />
        </div>
    {/if}
    <div id="lvc-content-wrapper">
        <div id="lvc-duration-section">            
            <div id="lvc-progress-current-duration-label" class="lvc-time-label">
                {#if !!the_video_element && the_video_element.readyState > 0}
                    <p>{video_progress_string}</p>
                {/if}
            </div>
            <div id="lvc-progress-bar-track">
                <div id="lvc-pbt-track-bar"
                    bind:this={the_progress_bar_track}
                    on:click={handleProgressClick}
                >
                    <div id="lvc-pbt-tc-progress-wrapper">
                        <div id="lvc-pbt-tc-progress"
                            style:scale="{video_progress}% 1"
                        ></div>
                    </div>
                    {#if current_video_moments.length > 0} 
                        {#each current_video_moments as video_moment, h}
                            <div class="lvc-pbt-tc-video-moment"
                                data-video-moment-index="{h}"
                                style:--translate-x="{video_moment.getTimelineStartPoint(the_video_element.duration)}cqw"
                                on:click={handleVideoMomentClicked}
                            >
                                <div class="lvc-pbt-tc-vm-title-wrapper">
                                    <div class="lvc-pbt-tc-vm-title">
                                        <p class="lvc-pbt-tc-vm-title-label">
                                            {video_moment.Title}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    {/if}
                    <div id="lvc-pbt-tc-time-scrubber"
                        style:translate="{video_progress}cqw"
                        on:mousedown={handleTimeScrubberMouseDown}
                        on:mouseup={handleTimeScrubberMouseUp}
                    ></div>
                </div>
            </div>
            <div id="lvc-progress-total-duration-label" class="lvc-time-label">
                <p>{video_duration_string}</p>
            </div>
        </div>
        <fieldset id="video-controls">
            <button class="lvc-control-btn" id="lvc-toggle-mute-btn" on:click={toggleMute}>
                <svg viewBox="0 0 24 24">
                    <path class="outline-path thin" d="M1 18L1 6H5l5 -4V22l-5 -4Z"/>
                    {#if !video_muted && !the_video_element.muted}
                        <path class="outline-path thin" d="M14 6Q 20 12 14 18"/>
                        <path class="outline-path thin" d="M14 2C 24 2 24 20 14 22"/>
                    {:else}
                        <path class="outline-path" d="M14 6L23 18M14 18L23 6"/>
                    {/if}
                </svg>
            </button>
            <button class="lvc-control-btn" id="lvc-back-btn" on:click={() => skipVideoPercentage(false)}>
                <svg viewBox="0 0 100 100">
                    <path class="outline-path" d="M 30 80A 40 40 0 1 0 25 20m5 -0.25l-5 0.25l0 -5" />
                    <text x="50" y="50">-5%</text>
                </svg>
            </button>
            <button class="lvc-control-btn" id="lvc-pause-btn" aria-label="Pause video" on:click={togglePauseVideo}>
                <svg viewBox="0 0 100 100">
                    {#if !video_paused}
                        <path d="M 30 20 L 30 80 L 45 80 L 45 20 L 30 20 Z M 55 20 L 55 80 L 70 80 L 70 20 L 55 20 Z" />
                    {:else}
                        <path d="M20 20L74 50L20 80Z"/>
                    {/if}
                </svg>
            </button>
            <button class="lvc-control-btn" id="lvc-forward-btn" aria-label="go forward 5%" on:click={() => skipVideoPercentage(true)}>
                <svg viewBox="0 0 100 100">
                    <path d="M 70 80A 40 40 0 1 1 75 20m-5 0.25l5 -0.25l0 -5" />
                    <text x="50" y="50">+5%</text>
                </svg>
            </button>
            <button class="lvc-control-btn" id="lvc-playback-speed-btn" aria-label="slow down video" on:click={() => setVideoPlaybackRate(false)}>
                <svg viewBox="0 0 24 24" fill="none">
                        <path class="outline-path thin" d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2" /> 
                        <path class="outline-path thin" d="M12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2" stroke-dasharray="4 3"/>
                        <path class="outline-path thin" d="M15.4137 10.941C16.1954 11.4026 16.1954 12.5974 15.4137 13.059L10.6935 15.8458C9.93371 16.2944 9 15.7105 9 14.7868L9 9.21316C9 8.28947 9.93371 7.70561 10.6935 8.15419L15.4137 10.941Z" />
                </svg>
            </button>
        </fieldset>
    </div>
</div>

<style>
    @keyframes toggle-controller-visibility {
        0% {
            visibility: hidden;
            opacity: 0;
        }
        100% {
            visibility: visible;
            opacity: 1;
        }
    }

    #libery-video-controller {
        position: relative;
        display: flex;
        container-type: size;
        background-color: hsl(from var(--body-bg-color) h s l / 0.5);
        backdrop-filter: blur(10px);
        border: 0.2px solid hsl(from var(--main-dark-color-7) h s l / 0.7);
        justify-content: center;
        align-items: center;
        box-sizing: border-box;
        min-height: 100px;
        height: 100%;
        width: 100%;
        padding: var(--spacing-1);
        transition: all .46s ease-in-out, visibility 0s linear;
    }

    #libery-video-controller.time-scrubber-dragging {
        cursor: grabbing;
        user-select: none;
    }

    #libery-video-controller #lvc-content-wrapper {
        display: grid;
        box-sizing: border-box;
        width: min(100cqw, 800px);
        height: 100cqh;
        row-gap: 10%;
    }

    #lvc-duration-section {
        width: 100%;
        display: flex;
        column-gap: 5%;
    }

    fieldset#video-controls {
        display: flex;
        justify-content: space-between;
    }

    
    /*=============================================
    =            Video moments            =
    =============================================*/
    
        #lvc-new-video-moment-wrapper {
            position: absolute;
            top: -50cqh;
            left: 0;
        }

        .lvc-pbt-tc-video-moment {
            position: absolute;
            background: var(--success-3);
            top: 0;
            left: 0;
            width: 100cqh;
            aspect-ratio: 1 / 1;
            border-radius: 9999px;
            transform-origin: center center;
            translate: var(--translate-x) 0;
            transition: .1s scale ease-out, .3s translate ease-out;

            &:hover {
                background: var(--success);
                translate: var(--translate-x) -40cqh;
                scale: 1.3;
            }
        }

        .lvc-pbt-tc-vm-title-wrapper {
            position: relative;
            width: 100cqh;
            aspect-ratio: 1 / 1;

            & .lvc-pbt-tc-vm-title {
                position: absolute;
                inset: auto -10% 200cqh auto;
                background: var(--grey-8);
                padding: calc(var(--spacing-1) * 0.5) var(--spacing-1);
                border-radius: var(--rounded-box-border-radius);
                opacity: 0;
                scale: 0;
                translate: 100% -20%;
                transition: scale .2s ease-out, opacity 0.3s ease-out, translate 0.3s ease-out;
            }

            & .lvc-pbt-tc-vm-title  p {
                font-size: calc(.7 * var(--font-size-1));
                
            }
        }

        .lvc-pbt-tc-video-moment:hover {
            & .lvc-pbt-tc-vm-title { 
                opacity: 1;
                scale: 1;
                translate: 0;
            }
        }

        .lvc-pbt-tc-video-moment:nth-child(4) {
            & .lvc-pbt-tc-vm-title { 
                opacity: 1;
            }
        }
    
    /*=====  End of Video moments  ======*/
    
    
    
    /*----------  Time labels  ----------*/

        .lvc-time-label {
            width: 8em
        }

        .lvc-time-label > p {
            box-sizing: border-box;
            color: var(--main-7);
            font-size: var(--font-size-2);
            text-align: center;
            font-family: var(--font-read);  
        }

        #lvc-progress-current-duration-label {
            grid-area: vpt;        
        }

        #lvc-progress-total-duration-label {
            grid-area: vdt;
        }
    
    /*----------  Control buttons  ----------*/
        
        #libery-video-controller button.lvc-control-btn {
            width: 100%;
            display: flex;
            background: none;
            border: none;
            align-items: center;
            justify-content: center;
            transition: all .28s ease-in-out;
        }

        #libery-video-controller.time-scrubber-dragging button.lvc-control-btn {
            pointer-events: none;
        }
        
        #libery-video-controller .lvc-control-btn svg {
            width: min(16cqw, 50px);
        }

        #libery-video-controller .lvc-control-btn svg path {
            fill: var(--main-dark);
        }
        
        #libery-video-controller .lvc-control-btn svg path.outline-path {
            stroke: var(--main-dark);
            fill: none;
            stroke-width: 2px;
        }

        #libery-video-controller .lvc-control-btn svg path.outline-path.thin {
            stroke-width: 1px;
        }

        #libery-video-controller .lvc-control-btn svg text {
            font-size: var(--font-size-3);
            font-family: var(--font-read);
            fill: var(--main-dark);
            transform-box: fill-box;
            transform-origin: center center;
            transform: translateX(-50%); 
        }

        @media (pointer:fine) {
            #libery-video-controller button:hover {
                background: var(--grey-5);
            }
        }

        #lvc-pause-btn {
            grid-area: p;

            & svg {
                overflow: visible;
            }
        }
        
        #libery-video-controller #lvc-back-btn.lvc-control-btn {
            grid-area: vb;

            & svg path {
                stroke: var(--main-dark);
                stroke-width: 5px;
                fill: transparent;
            }
        }

        #libery-video-controller #lvc-forward-btn.lvc-control-btn {
            grid-area: vf;

            & svg path {
                stroke: var(--main-dark);
                stroke-width: 5px;
                fill: transparent;
            }
        }

        #lvc-playback-speed-btn {
            grid-area: e2;
        }

        #lvc-toggle-mute-btn {
            grid-area: e1;

            & svg path {
                stroke-linejoin: round;
            }
        }
    
    /*----------  ProgressBar  ----------*/
    
        #lvc-progress-bar-track {
            grid-area: vp;
            width: 100%;
            /* height: 20cqh; */
            display: flex;
            align-items: center;
            justify-content: center;
        }

        #lvc-pbt-track-bar {
            position: relative;
            container-type: size;
            width: 100%;
            height: 8px;
            border-radius: 4px;
            background: var(--grey-black);

            & #lvc-pbt-tc-progress-wrapper {
                width: 100%;
                height: 100%;
                border-radius: inherit;
                overflow: hidden;
            }

            & #lvc-pbt-tc-progress {
                width: 100%;
                height: 100%;
                background: var(--grey-8);
                transform-origin: left center;
                transition: scale 0.22s linear;
            }

            & #lvc-pbt-tc-time-scrubber {
                position: absolute;
                background: hsl(from var(--main-dark) h s l / 0.9);
                top: 0;
                left: 0;
                width: 100cqh;
                height: 100cqh;
                border-radius: 50%;
                transform-origin: center;
                scale: 2.3;
                transition: scale 0.3s ease-out, translate 0.22s linear;
            } 

            & #lvc-pbt-tc-time-scrubber:hover {
                scale: 3;
            }
        }

        .time-scrubber-dragging #lvc-pbt-track-bar {
            & #lvc-pbt-tc-progress {
                transition: scale 0s linear;
            }

            & #lvc-pbt-tc-time-scrubber {
                transition: scale 0.3s ease-out, translate 0s linear;
            }
        }

    @media only screen and (max-width: 612px) {
        #libery-video-controller {
            padding: var(--vspacing-3);
        }
    }

    /* @container (width > 1000px) {
        #libery-video-controller #lvc-content-wrapper {
            padding-inline: 10cqw;
        }        
    } */
</style>