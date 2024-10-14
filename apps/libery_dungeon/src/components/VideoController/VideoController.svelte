<script>
    import { onMount, onDestroy, createEventDispatcher } from "svelte";
    import { active_media_index, automute_enabled, previous_media_index } from "@stores/media_viewer";
    import { saveMediaWatchPoint, getMediaWatchPoint } from "@models/Metadata";
    import global_hotkeys_manager from "@libs/LiberyHotkeys/libery_hotkeys";
    import { layout_properties } from "@stores/layout";
    import { browser } from "$app/environment";
    import { videoDurationToString } from "@libs/utils";


    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {HTMLVideoElement} the video element that will be controlled */
        export let video_element;

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
                handler: pauseVideo,
                options: {
                    description: "Pause/Play video",
                }
            },
            TOGGLE_MUTE: {
                key_combo: "m",
                handler: toggleMute,
                options: {
                    description: "Mute/Unmute video",
                }
            },
            VOLUMEN_UP: {
                key_combo: "up",
                handler: handleVolumeUpHotkey,
                options: {
                    description: "Increase volume by 10%",
                }
            },
            VOLUMEN_DOWN: {
                key_combo: "down",
                handler: handleVolumeDownHotkey,
                options: {
                    description: "Decrease volume by 10%",
                }
            },
            FORWARD_VIDEO: {
                key_combo: "shift+x",
                handler: () => skipVideoPercentage(true),
                options: {
                    description: "Forward video 5% of the total duration(min of 5 seconds). if overflows, jumps to the start",
                }
            },
            BACKWARD_VIDEO: {
                key_combo: "x",
                handler: () => skipVideoPercentage(false),
                options: {
                    description: "Backward video 5% of the total duration(min of 5 seconds). No overflow",
                }
            },
            FORWARD_SECS_VIDEO: {
                key_combo: "shift+alt+x",
                handler: () => skipVideoSeconds(5),
                options: {
                    description: "Forward video 5 seconds",
                }
            },
            BACKWARD_SECS_VIDEO: {
                key_combo: "alt+x",
                handler: () => skipVideoSeconds(-5),
                options: {
                    description: "Backward video 5 seconds",
                }
            },
            SKIP_FRAME_FORWARD: {
                key_combo: "shift+`",
                handler: () => skipFrame(true),
                options: {
                    can_repeat: true,
                    description: "Skip frame forward",
                }
            },
            SKIP_FRAME_BACKWARD: {
                key_combo: "`",
                handler: () => skipFrame(false),
                options: {
                    can_repeat: true,
                    description: "Skip frame backward",
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
                key_combo: ",",
                handler: () => setVideoPlaybackRate(false),
                options: {
                    description: "Speed up video",
                }
            },
            SLOW_DOWN_VIDEO: {
                key_combo: ".",
                handler: () => setVideoPlaybackRate(true),
                options: {
                    description: "Slow down video",
                }
            },
            SCREENSHOT_VIDEO: {
                key_combo: "shift+s",
                handler: () => emitScreenshotVideo(),
                options: {
                    description: "Take a screenshot of the video",
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

            /** @type {number} the timeout id for the controller visibility timeout */
            let controller_visibility_interval_id = null;

            /** @type {number} the opacity of the controller, managed by the controller visibility timeout */
            let controller_opacity = auto_hide ? 0 : 0.8;

            /** @type {boolean} whether the video is currently paused */
            let video_paused = false;

            /** @type {number} whether the video is muted */
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

        
        
        /*=====  End of State  ======*/

        let active_media_index_unsubscriber = () => {};
    
        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    onMount(() => {
        video_element.addEventListener("timeupdate", handleVideoTimeUpdate);
        video_element.addEventListener("durationchange", handleVideoDurationChange);
        video_element.addEventListener("loadedmetadata", handleVideoMetadataLoaded);
        window.addEventListener("beforeunload", saveVideoWatchProgress);

        active_media_index_unsubscriber = active_media_index.subscribe(handleActiveMediaIndexChange);

        defineVideoControllerKeybinds();
    })

    onDestroy(() => {
        if (!browser) return;

        video_element.removeEventListener("timeupdate", handleVideoTimeUpdate);
        video_element.removeEventListener("durationchange", handleVideoDurationChange);
        video_element.removeEventListener("loadedmetadata", handleVideoMetadataLoaded);
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
                if (layout_properties.IS_MOBILE || !browser || !global_hotkeys_manager.hasLoadedContext()) return;

                const video_controls_description_group = "<video_controls>";

                Object.values(keybinds).forEach(keybind => {
                    keybind.options.description = `${video_controls_description_group} ${keybind.options.description ?? "Empty description"}`;
                    global_hotkeys_manager.registerHotkeyOnContext(keybind.key_combo, keybind.handler, keybind.options);
                })
            }

            const removeVideoControllerKeybinds = () => {
                if (layout_properties.IS_MOBILE || !global_hotkeys_manager.hasLoadedContext()) return;

                Object.values(keybinds).forEach(keybind => {
                    global_hotkeys_manager.unregisterHotkeyFromContext(keybind.key_combo, "keydown");
                });
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
            video_element.volume = Math.min(1, Math.max(0, video_element.volume + amount));
        }

        /**
         * Fetches the watch progress of the video.
         * @this video_element
         * @returns {void}
         */
        async function getWatchProgress() {
            if (media_uuid == null || this.duration * 1000 < SAVE_WATCH_PROGRESS_THRESHOLD) return;

            let watch_progress = await getMediaWatchPoint(media_uuid);

            if (watch_progress == null) return;

            this.currentTime = watch_progress / 1000;
        }

        function pauseVideo() {
            video_paused = !video_paused;
            
            if (video_paused) {
                video_element.pause();
            } else {
                video_element.play();
            }
        }

        function toggleMute() {
            video_muted = !video_muted;

            automute_enabled.set(video_muted); // Muted value of the video element is reactive to this store
        }

        const handleActiveMediaIndexChange = (new_active_media_index) => {
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
            if (!hotkey.WithVimMotion) return;

            let minutes_in_video = Math.trunc(video_element.duration / 60);
            if (isNaN(minutes_in_video)) return;

            let minute_to_seek = hotkey.MatchMetadata.MotionMatches[0];

            minute_to_seek = Math.min(minutes_in_video, Math.max(0, minute_to_seek));

            video_element.currentTime = minute_to_seek * 60;
        }
        
        /**
         * Handles the video element DurationChange event
         * @this video_element
         * @returns {void}
         */
        function handleVideoDurationChange() {
            setSaveWatchProgress.call(this);
        }

        /**
         * Handles the video element timeupdate event   
         * @this video_element
         * @returns {void}
         */
        function handleVideoTimeUpdate() {
            video_progress_string = videoDurationToString(this.currentTime);

            updateVideoProgress.call(this);
        }

        /**
         * Handles the video element loadedmetadata event
         * @this video_element
         * @returns {void}
         */
        function handleVideoMetadataLoaded() {
            video_metadata_loaded = true;
            video_duration_string = videoDurationToString(this.duration);

            getWatchProgress.call(this);
        }


        function emitScreenshotVideo() {
            console.log("emitting screenshot video");
            dispatch("screenshot-video");            
        }

        /**
         * Sets the video playback rate
         * @param {boolean} increase whether to increase or decrease the playback rate
        */
        function setVideoPlaybackRate(increase) {
            const step_diff = 0.25;

            let step = increase ? step_diff : -step_diff;

            video_element.playbackRate = Math.min(2, Math.max(0.25, video_element.playbackRate + step));
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
            let clamped_duration = Math.min(video_element.duration, Math.max(0, new_duration));
            
            if (overflow_allowed && clamped_duration !== new_duration) {
                clamped_duration = new_duration < 0 ? video_element.duration + new_duration : 0;
            }

            video_element.currentTime = clamped_duration;
        }

        /**
         * Sets save_watch_progress based on the video duration.
         * @this video_element
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

            let watch_progress = Math.floor(video_element.currentTime * 1000);

            saveMediaWatchPoint(media_uuid, watch_progress);
        }

        /**
         * Skips video by 5%, if forward is true, skips forward, if false, skips backward
         * @param {boolean} forward whether to skip forward or backward
         */
        function skipVideoPercentage(forward) {
            let step = forward ? 1 : -1;
            let new_time = video_element.currentTime + (step * Math.max(video_element.duration * 0.05));

            setVideoDuration(new_time, forward);
        }

        /**
         * Skips video by a given amount of seconds.
         * @param {number} seconds
        */
        function skipVideoSeconds(seconds) {
            let direction_forward = seconds > 0;
            let new_time = video_element.currentTime + seconds;

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
            let new_time = video_element.currentTime + step;

            if (!video_element.paused) {
                video_element.pause();
            }

            setVideoDuration(new_time, forward);
        }

        /**
         * Updates the video progress percentage
         * @modifies {video_progress}
         * @returns {void}
        */
        const updateVideoProgress = () => {
            video_progress = (video_element.currentTime / video_element.duration) * 100;
        }
            
        /**
         * Handles the video progress bar click event
         * @param {MouseEvent} event
         * @returns {void}
        */
        const handleProgressClick = (event) => {
            event.stopPropagation();
            const rect = event.currentTarget.getBoundingClientRect();
            const clickX = event.clientX - rect.left;
            const new_progress = (clickX / rect.width);
            let new_video_time =  new_progress * video_element.duration;
            
            video_element.currentTime = new_video_time;
        };

        const handleMouseMovement = () => {
            if (controller_visibility_interval_id === null) {
                controller_visibility_interval_id = window.setInterval(handleControllerVisibility, 300);
            }

            controller_opacity = 1;
            controller_visible = true;
        }

        // TODO: create a proper visibility controller for mobile, although this more or less works, it's due to pure black magic and also once the controller appears it doesn't disappear ever again on unless the media viewer 
        // unmounts and mounts again the video controller (like when the media changes to an from video to image and back to video)
        const handleControllerVisibility = () => {
            if (mouse_over_controller || !auto_hide) return;


            controller_opacity = Math.max(0, controller_opacity - 0.5);
            controller_visible = controller_opacity > 0;

            if (!controller_visible) {
                window.clearInterval(controller_visibility_interval_id);
                controller_visibility_interval_id = null;
            }
        }

        // Hotfix to the controller not disappearing on mobile. we check if the mouse(the finger) touches the controller but not a button, if so we hide the controller
        const handleControllerTouch = (event) => {
            if (event.target === event.currentTarget)

            mouse_over_controller = false;
        }

    /*=====  End of Methods  ======*/
        
</script>

<svelte:document on:mousemove={handleMouseMovement} />
<div 
    role="group" 
    id="libery-video-controller" 
    aria-label="Video controls" 
    class:adebug={false}  
    style:opacity={controller_opacity}
    style:visibility={controller_visible ? "visible" : "hidden"}
    on:mouseenter={() => mouse_over_controller = true}
    on:mouseleave={() => mouse_over_controller = false}
    on:touchstart={handleControllerTouch}
>
    <div id="lvc-progress-current-duration-label" class="lvc-time-label">
        {#if !!video_element && video_element.readyState > 0}
            <p>{video_progress_string}</p>
        {/if}
    </div>
    <div id="lvc-progress-total-duration-label" class="lvc-time-label">
        <p>{video_duration_string}</p>
    </div>
    <div id="lvc-progress-bar-track">
        <svg id="lvc-pbt-track-bar" viewBox="0 0 102 5" on:click={handleProgressClick} >
            <path id="lvc-pbt-tb-empty-track" d="M 2 2.5 L 100 2.5" />
            <path id="lvc-pbt-tb-progress-track" d="M 2 2.5 L {video_progress + 2} 2.5" />
            <circle id="lvc-pbt-tb-progress-indicator" cx="{video_progress + 3}" cy="2.5" r="1.5" />
        </svg>
    </div>
    <button id="lvc-back-btn" on:click={() => skipVideoPercentage(false)}>
        <svg viewBox="0 0 100 100">
            <path class="outline-path" d="M 30 80A 40 40 0 1 0 25 20m5 -0.25l-5 0.25l0 -5" />
            <text x="50" y="50">-5%</text>
        </svg>
    </button>
    <button id="lvc-pause-btn" aria-label="Pause video" on:click={pauseVideo}>
        <svg viewBox="0 0 100 100">
            {#if !video_paused}
                <path d="M 30 20 L 30 80 L 45 80 L 45 20 L 30 20 Z M 55 20 L 55 80 L 70 80 L 70 20 L 55 20 Z" />
            {:else}
                <path d="M20 20L74 50L20 80Z"/>
            {/if}
        </svg>
    </button>
    <button id="lvc-forward-btn" aria-label="go forward 5%" on:click={() => skipVideoPercentage(true)}>
        <svg viewBox="0 0 100 100">
            <path class="outline-path" d="M 70 80A 40 40 0 1 1 75 20m-5 0.25l5 -0.25l0 -5" />
            <text x="50" y="50">+5%</text>
        </svg>
    </button>
    <button id="lvc-playback-speed-btn" aria-label="slow down video" on:click={() => setVideoPlaybackRate(false)}>
        <svg viewBox="0 0 24 24" fill="none">
                <path class="outline-path thin" d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2" /> 
                <path class="outline-path thin" d="M12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2" stroke-dasharray="4 3"/>
                <path class="outline-path thin" d="M15.4137 10.941C16.1954 11.4026 16.1954 12.5974 15.4137 13.059L10.6935 15.8458C9.93371 16.2944 9 15.7105 9 14.7868L9 9.21316C9 8.28947 9.93371 7.70561 10.6935 8.15419L15.4137 10.941Z" />
        </svg>
    </button>
    <button id="lvc-toggle-mute-btn" on:click={toggleMute}>
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
        container-type: inline-size;
        box-sizing: border-box;
        min-height: 100px;
        height: 100%;
        width: 100%;
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        background: var(--grey-8);
        padding: var(--vspacing-1);
        transition: all .46s ease-in-out;
        grid-template-areas: 
            "vpt vp vp vp vp vp vdt"
            "rs e1 vb p vf e2 rs2";
    }

    #libery-video-controller button {
        width: 100%;
        display: flex;
        background: none;
        border: none;
        align-items: center;
        justify-content: center;
        transition: all .28s ease-in-out;
    }

    
    
    #libery-video-controller svg {
        width: 10cqw;
    }

    @media (pointer:fine) {
        #libery-video-controller button:hover {
            background: var(--grey-5);
        }
    }

    
    #libery-video-controller svg path {
        fill: var(--main-dark);
    }
    
    #libery-video-controller svg path.outline-path {
        stroke: var(--main-dark);
        fill: none;
        stroke-width: 2px;
    }

    #libery-video-controller svg path.outline-path.thin {
        stroke-width: 1px;
    }

    #libery-video-controller svg text {
        font-size: var(--font-size-3);
        font-family: var(--font-read);
        fill: var(--main-dark);
        transform-box: fill-box;
        transform-origin: center center;
        transform: translateX(-50%); 
    }

    .lvc-time-label > p {
        color: var(--main-7);
        font-size: var(--font-size-1);
        text-align: center;
        font-family: var(--font-read);  
    }

    #lvc-progress-current-duration-label {
        grid-area: vpt;        
    }

    #lvc-progress-total-duration-label {
        grid-area: vdt;
    }



    #lvc-pause-btn {
        grid-area: p;
    }
    
    #lvc-pause-btn svg {
        width: 16cqw;
        overflow: visible;
    }
    
    #lvc-forward-btn {
        grid-area: vf;
    }

    #lvc-back-btn {
        grid-area: vb;
        /* transform: rotateX(180deg); */
    }

    #lvc-playback-speed-btn {
        grid-area: e2;
    }

    #lvc-toggle-mute-btn {
        grid-area: e1;
    }

    
    /*----------  ProgressBar  ----------*/
    
    #lvc-progress-bar-track {
        grid-area: vp;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    
    svg#lvc-pbt-track-bar {
        width: 100%;
        height: 100%;
        overflow: visible;
    }

    path#lvc-pbt-tb-empty-track {
        stroke: var(--grey-5);
        stroke-width: 2px;
        stroke-linecap: round;
    }

    circle#lvc-pbt-tb-progress-indicator {
        fill: var(--main-dark);
        transform-box: fill-box;
        transform-origin: center center;
        transition: all .35s ease-in-out;
    }

    circle#lvc-pbt-tb-progress-indicator:hover {
        transform: scale(1.8);
        opacity: .5;
    }

    path#lvc-pbt-tb-progress-track {
        stroke: var(--main-dark-color-8);
        stroke-width: 2px;
        stroke-linecap: round;
    }

    @media only screen and (max-width: 612px) {
        #libery-video-controller {
            padding: var(--vspacing-3);
        }
    }

</style>