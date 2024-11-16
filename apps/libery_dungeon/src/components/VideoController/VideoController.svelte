<script>
    import { onMount, onDestroy, createEventDispatcher } from "svelte";
    import { active_media_index, automute_enabled, previous_media_index } from "@stores/media_viewer";
    import { saveMediaWatchPoint, getMediaWatchPoint } from "@models/Metadata";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import { layout_properties } from "@stores/layout";
    import { browser } from "$app/environment";
    import { videoDurationToString } from "@libs/utils";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {HTMLVideoElement} the video element that will be controlled */
        export let the_video_element;

        let global_hotkeys_manager = getHotkeysManager();

        /**
         * The media uuid of the video element
         * @type {string}
        */
        export let media_uuid;


        /** @type {boolean} wheter the component should disappear when the mouse is not over it */
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
                handler: handlePauseHotkey,
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
        
            /** @type {boolean} whether the parent component has position: absolute set, mainly used for styling */
            let controller_floating = false;

            /** @type {boolean} whether the controller is visible or not*/
            let controller_visible = !auto_hide;

            /** @type {boolean} whether the mouse is over the component or not*/
            let mouse_over_controller = false;

            /** @type {number | null} the timeout id for the controller visibility timeout */
            let controller_visibility_interval_id = null;

            /** @type {number} the opacity of the controller, managed by the controller visibility timeout */
            let controller_opacity = auto_hide ? 0 : 0.8;

            /** @type {boolean} whether the video is currently paused */
            let video_paused = false;

            /** @type {boolean} whether the video is muted */
            let video_muted = true;

            /** @type {number} video progress percentage */    
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
        
        /*=====  End of State  ======*/

        let active_media_index_unsubscriber = () => {};
    
        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    onMount(() => {
        the_video_element.addEventListener("timeupdate", handleVideoTimeUpdate);
        the_video_element.addEventListener("durationchange", handleVideoDurationChange);
        the_video_element.addEventListener("loadedmetadata", handleVideoMetadataLoaded);
        window.addEventListener("beforeunload", saveVideoWatchProgress);

        active_media_index_unsubscriber = active_media_index.subscribe(handleActiveMediaIndexChange);

        defineVideoControllerKeybinds();
    })

    onDestroy(() => {
        if (!browser) return;

        the_video_element.removeEventListener("timeupdate", handleVideoTimeUpdate);
        the_video_element.removeEventListener("durationchange", handleVideoDurationChange);
        the_video_element.removeEventListener("loadedmetadata", handleVideoMetadataLoaded);
        window.removeEventListener("beforeunload", saveVideoWatchProgress);

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

                
                if ($layout_properties.IS_MOBILE || !browser || !global_hotkeys_manager.hasLoadedContext()) return;

                const video_controls_description_group = "<video_controls>";

                Object.values(keybinds).forEach(keybind => {
                    keybind.options.description = `${video_controls_description_group} ${keybind.options.description ?? "Empty description"}`;
                    global_hotkeys_manager.registerHotkeyOnContext(keybind.key_combo, keybind.handler, keybind.options);
                })
            }

            const removeVideoControllerKeybinds = () => {
                if (global_hotkeys_manager == null) {
                    console.error("The global hotkeys manager is null");
                    return;
                }

                if ($layout_properties.IS_MOBILE || !global_hotkeys_manager.hasLoadedContext()) return;

                Object.values(keybinds).forEach(keybind => {
                    global_hotkeys_manager.unregisterHotkeyFromContext(keybind.key_combo, "keydown");
                });
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

            function handlePauseHotkey() {
                pauseVideo();

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

                setVideoDuration(last_frame_skip_timestamp);
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
                auto_hide = !auto_hide;

                if (auto_hide) {
                    setControllerHiddenTimeout();
                } else {
                    controller_visible = true;
                    controller_opacity = 1;
                }

                let feedback_message = `auto-hide: ${auto_hide ? "on" : "off"}`;

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
         * Fetches the watch progress of the video.
         * @this the_video_element
         * @returns {Promise<void>}
         */
        async function getWatchProgress() {
            if (media_uuid == null || this.duration * 1000 < SAVE_WATCH_PROGRESS_THRESHOLD) return;

            let watch_progress = await getMediaWatchPoint(media_uuid);

            if (watch_progress == null) return;

            this.currentTime = watch_progress / 1000;
        }

        /**
         * Returns the video_element currentTime in a string format "hh:mm:ss"
         */
        function getVideoCurrentProgressString() {
            return videoDurationToString(the_video_element.currentTime);
        }

        const handleActiveMediaIndexChange = () => {
            video_metadata_loaded = false;

            saveVideoWatchProgress();
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
        }

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
            
            the_video_element.currentTime = new_video_time;
        };

        const handleMouseMovement = () => {
            setControllerHiddenTimeout();
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

        function pauseVideo() {
            video_paused = !the_video_element.paused;
            
            if (video_paused) {
                the_video_element.pause();
            } else {
                the_video_element.play();
            }
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

            return new_playback_rate;
        }

        /**
         * Sets the video.currentTime to the given time clamping it to the video duration and 0.
         * if overflow_allowed is true, and, for example the new_duration is -5, the duration will be set to
         * the video duration - 5 seconds. but if the new_duration is video.duration + 5, then the duration
         * will be set to 0, not 5.
         * @param {number} new_duration the new duration(in seconds) to set the video to
         * @param {boolean} overflow_allowed whether the duration can overflow the video duration
         * @requires video_element
         * @returns {void}
        */
        function setVideoDuration(new_duration, overflow_allowed = false) {
            let clamped_duration = Math.min(the_video_element.duration, Math.max(0, new_duration));
            
            if (overflow_allowed && clamped_duration !== new_duration) {
                clamped_duration = new_duration < 0 ? the_video_element.duration + new_duration : 0;
            }

            the_video_element.currentTime = clamped_duration;
        }

        /**
         * Sets save_watch_progress based on the video duration.
         * @this the_video_element
         * @returns {void}
         */
        function setSaveWatchProgress() {
            let video_duration_ms = this.duration * 1000;
            save_watch_progress = video_duration_ms >= SAVE_WATCH_PROGRESS_THRESHOLD;
        }

        /**
         * Saves the video watch progress
         * @returns {void}
         */
        const saveVideoWatchProgress = () => {
            if (!save_watch_progress || media_uuid == null) return;

            let watch_progress = Math.floor(the_video_element.currentTime * 1000);

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

            setVideoDuration(new_time, forward);

            return video_percentage;
        }

        /**
         * Skips video by a given amount of seconds.
         * @param {number} seconds
        */
        function skipVideoSeconds(seconds) {
            let direction_forward = seconds > 0;
            let new_time = the_video_element.currentTime + seconds;

            setVideoDuration(new_time, direction_forward);
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

            setVideoDuration(new_time, forward);
        }

        function toggleMute() {
            video_muted = !video_muted;

            automute_enabled.set(video_muted); // Muted value of the video element is reactive to this store
        }

        /**
         * Updates the video progress percentage
         * @modifies {video_progress}
         * @returns {void}
        */
        const updateVideoProgress = () => {
            if (isNaN(the_video_element.duration)) return;

            video_progress = (the_video_element.currentTime / the_video_element.duration) * 100;
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

    /*=====  End of Methods  ======*/
        
</script>

<svelte:document on:mousemove={handleMouseMovement} />
<div id="libery-video-controller" 
    class="libery-dungeon-window"
    role="group" 
    aria-label="Video controls" 
    class:adebug={false}  
    style:opacity={controller_opacity}
    style:visibility={controller_visible ? "visible" : "hidden"}
    on:mouseenter={() => mouse_over_controller = true}
    on:mouseleave={() => mouse_over_controller = false}
    on:touchstart={handleControllerTouch}
>
    <div id="lvc-content-wrapper">
        <div id="lvc-duration-section">            
            <div id="lvc-progress-current-duration-label" class="lvc-time-label">
                {#if !!the_video_element && the_video_element.readyState > 0}
                    <p>{video_progress_string}</p>
                {/if}
            </div>
            <div id="lvc-progress-bar-track">
                <!-- <svg id="lvc-pbt-track-bar" 
                    viewBox="0 0 102 5"
                    on:click={handleProgressClick}
                    preserveAspectRatio="none"
                >
                    <path id="lvc-pbt-tb-empty-track" d="M 2 2.5 L 100 2.5" />
                    <path id="lvc-pbt-tb-progress-track" d="M 2 2.5 L {video_progress + 2} 2.5" />
                    <circle id="lvc-pbt-tb-progress-indicator" cx="{video_progress + 3}" cy="2.5" r="2.5" />
                </svg> -->
                <div id="lvc-pbt-track-bar">
                    <div id="lvc-pbt-tc-progress-wrapper">
                        <div id="lvc-pbt-tc-progress"
                            style:scale="{video_progress}% 1"
                        ></div>
                    </div>
                    <div id="lvc-pbt-tc-time-scrubber"
                        style:translate="{video_progress}cqw"
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
                    {#if !video_muted}
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
            <button class="lvc-control-btn" id="lvc-pause-btn" aria-label="Pause video" on:click={pauseVideo}>
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
        transition: all .46s ease-in-out;
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
                background: var(--main-dark);
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