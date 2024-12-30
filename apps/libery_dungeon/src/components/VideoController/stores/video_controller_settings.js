
/*=============================================
=            Properties            =
=============================================*/

    /**
     * The settings for the video controller.
     */
    const video_controller_settings = {
        playback_speed: 1,
        preserve_playback_speed: true,
    }

    /**
     * The minimum value the playback speed can have.
     * @type {number}
     */
    const MINIMUM_PLAYBACK_SPEED = 0.1

    /**
     * The maximum value the playback speed can have.
     * @type {number}
     */
    const MAXIMUM_PLAYBACK_SPEED = 2;

/*=====  End of Properties  ======*/


/*=============================================
=            Methods            =
=============================================*/

    /**
     * Returns the video playback speed.
     * @returns {number}
     */
    const getPlaybackSpeed = () => {
        return video_controller_settings.playback_speed;
    }

    /**
     * Sets the video controller speed. Returns the playback speed that it save to the configuration, this may differ from the value passed as the function clamps the playback speed to be in a certain range.
     * @param {number} new_playback_speed
     * @returns {number}
     */
    const setPlaybackSpeed = new_playback_speed => {
        new_playback_speed = Math.max(MINIMUM_PLAYBACK_SPEED, Math.min(MAXIMUM_PLAYBACK_SPEED, new_playback_speed));

        video_controller_settings.playback_speed = new_playback_speed;

        return new_playback_speed;
    }

    /**
     * Returns whether the playback speed should be preserved between videos.
     * @returns {boolean}
     */
    const shouldPreservePlaybackSpeed = () => {
        return video_controller_settings.preserve_playback_speed;
    }


/*=====  End of Methods  ======*/

/**
 * The store controllers exported by the store.
 */
const settings_controllers = {
    getPlaybackSpeed,
    setPlaybackSpeed,
    shouldPreservePlaybackSpeed
}

export default settings_controllers;