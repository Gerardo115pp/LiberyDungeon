<script>
    import { LabeledError, VariableEnvironmentContextError } from '@libs/LiberyFeedback/lf_models';
    import { me_gallery_changes_manager, me_gallery_yanked_medias } from './me_gallery_state';
    import { media_change_types } from '@models/WorkManagers';
    import { onMount, onDestroy } from 'svelte';
    import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
    import viewport from '@components/viewport_actions/useViewportActions';
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The media item that will be displayed
         * @type {import('@models/Medias').OrderedMedia}
         */
        export let ordered_media;

        /**
         * This component's DOM node
         * @type {HTMLDivElement}         
        */
        let this_dom_element;

        /**
         * The thumbnails for the media item will be requested with a width equivalent to this percentage of the viewport width.
         * @type {number}
         * @default 0.2
         */
        export let media_viewport_percentage = 0.2;

        
        
        /*----------  State  ----------*/
        
            /**
             * True if this media is within the gallery keyboard focused cell. 
             * @type {boolean}
             */
            export let is_keyboard_focused = false;
            $: if (is_keyboard_focused && this_dom_element != null) {
                ensureElementIsVisible(this_dom_element);
            }

            /**
             * If true, and the media is a video or an animation. Instead of mounting a thumbnail, the video or animated image will be mounted.
             * @type {boolean}
             */
            export let enable_heavy_rendering = false;

            /**
             * Whether to show the video title on a tooltip(false) or show it on top of the media item(true).
             * @type {boolean}
             * @default false
             */
            export let enable_video_titles = false;

            /**
             * Whether to magnify the media item when it is keyboard focused.
             * @type {boolean}
             */
            export let enable_magnify_on_keyboard_focus = false;

            /**
             * Whether this_dom_element is inside the viewport.
             * @type {boolean}
             */
            export let media_inside_viewport = false;

            /**
             * Whether to only render the item as skeleton.
             * @type {boolean}
             * @default false
             */
            export let is_skeleton = false;
        
        /*----------  Style  ----------*/

            /**
             * Whether the parent container is using a masonry layout
             * @type {boolean}
             */
            export let use_masonry = false;

            /**
             * Whether the media items is selected.
             * @type {boolean}
             */
            export let is_media_selected = false;

            /**
             * Whether the media item is set as deleted.
             * @type {boolean}
             */
            export let is_media_deleted = false;

            /**
             * Whether the media item is among the yanked medias.
             * @type {boolean}
             */
            export let is_media_yanked = false;
        
        /*----------  Unsubscribers  ----------*/
        
            let yanked_medias_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        suscribeToMediaChanges();
        yanked_medias_unsubscriber = me_gallery_yanked_medias.subscribe(handleYankedMediasChange);
    });

    onDestroy(() => {
        unsuscribeFromMediaChanges();

        yanked_medias_unsubscriber();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Scrolls given element into view.
         * @param {HTMLElement} element
         */
        function ensureElementIsVisible(element) {
            element.scrollIntoView({ behavior: "smooth", block: "center", inline: "center" });
        }

        /**
         * Generates a media changes uuid based on the media item passed.
         * @param {import('@models/Medias').OrderedMedia} media_item
         * @returns {string}
        */
        const generateMediaChangesUuid = media_item => media_item.uuid + '_changes_callback';

        /**
         * Handles media change events emitted by the me_gallery_changes_manager. These are emitted every time there is any change for any media,
         * so we have to check if the change is for this media item.
         * @param {string} change_type
         * @param {string} media_uuid
         * @returns {void}
         */
        const handleMediaChanges = (change_type, media_uuid) => {
            if (media_uuid !== ordered_media.uuid) return;
            
            let trusted_change_type = $me_gallery_changes_manager.getMediaChangeType(media_uuid);
            if (trusted_change_type !== change_type) {
                console.warn(`Received a wrong change type '${change_type}' does not match the real change type registered in media changes manager: '${trusted_change_type}'`);
            }

            switch (trusted_change_type) {
                case media_change_types.DELETED:
                    handleMediaChangeToDeleted();
                    break;
                case media_change_types.MOVED:
                    handleMediaChangeToSelected();  
                    break;
                case media_change_types.NORMAL:
                    handleMediaChangeToNormal();
                    break;
                default:
                    let variable_environment_error = new VariableEnvironmentContextError("In MEGalleryDisplayItem.handleMediaChanges")
                    variable_environment_error.addVariable("change_type", change_type)
                    variable_environment_error.addVariable("trusted_change_type", trusted_change_type)
                    variable_environment_error.addVariable("this", this)
                    variable_environment_error.addVariable("media_item", ordered_media)
                    let labeled_err = new LabeledError(variable_environment_error, `Something weird happened. Received unknown change type: ${change_type}`, lf_errors.ERR_PROCESSING_ERROR);

                    labeled_err.alert();
                    break;
            }
        }

        /**
         * handles a media change that sets the media item normal.
         */
        const handleMediaChangeToNormal = () => {
            is_media_deleted = false;
            is_media_selected = false;
        }

        /**
         * handles a media change that sets the media item as deleted.
         */
        const handleMediaChangeToDeleted = () => {
            is_media_deleted = true;
            is_media_selected = false;  
        }

        /**
         * Handles a media change that sets the media item as selected.
         */
        const handleMediaChangeToSelected = () => {
            is_media_selected = true;
            is_media_deleted = false;   
        }

        /**
         * Hanldes viewport enter event emitted by the viewport action.
         * @requires viewport
         * @requires media_inside_viewport
         */
        const handleViewportEnter = () => {
            media_inside_viewport = true;
        }

        /**
         * Hanldes viewport leave event emitted by the viewport action.
         * @requires viewport
         * @requires media_inside_viewport
         */
        const handleViewportLeave = () => {
            media_inside_viewport = false;
        }
        

        /**
         * Checks if the media item is in the yanked media uuids list and sets the status accordingly.
         */
        const handleYankedMediasChange = () => {
            let new_is_media_yanked = false
            
            for (let h = 0; h < $me_gallery_yanked_medias.length && !new_is_media_yanked; h++) {
                new_is_media_yanked = $me_gallery_yanked_medias[h].uuid === ordered_media.uuid;
            }

            is_media_yanked = new_is_media_yanked;
        }

        /**
         * Suscribe to media changes events emitted by the me_gallery_changes_manager.
         */
         const suscribeToMediaChanges = () => {
            if ($me_gallery_changes_manager == null) return;

            let element_media_uuid = generateMediaChangesUuid(ordered_media);
            $me_gallery_changes_manager.suscribeToChanges(element_media_uuid, handleMediaChanges)
        }

        /**
         * Unsuscribe from media changes events emitted by the me_gallery_changes_manager.
         */
        const unsuscribeFromMediaChanges = () => {
            if ($me_gallery_changes_manager == null) return;

            let element_media_uuid = generateMediaChangesUuid(ordered_media);
            $me_gallery_changes_manager.unsubscribeToChanges(element_media_uuid, handleMediaChanges)   
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<div bind:this={this_dom_element} class="meg-display-item-wrapper"
    class:use-masonry={use_masonry}
    class:meg-di-keyboard-focused={is_keyboard_focused}
    class:status-magnified={enable_magnify_on_keyboard_focus && is_keyboard_focused}
    class:status-deleted={is_media_deleted}
    class:status-selected={is_media_selected}   
    class:status-yanked={is_media_yanked}
    class:status-titles-enabled={enable_video_titles}
    on:viewportEnter={handleViewportEnter}
    on:viewportLeave={handleViewportLeave}
    use:viewport
>
    <div class="media-status-overlay"></div>
    <div class="media-type-label-overlay">
        <h4 
            class="media-type-label"
            class:media-type-video={ordered_media.Media.isVideo()}
        >
            {ordered_media.FileExtension}
        </h4>
    </div>
    <div class="media-name-tool-tip">
        <p class="media-name-tool-tip-content">
            {ordered_media.MediaName}
        </p>
    </div>
    {#if !is_skeleton}
        {#if !ordered_media.Media.isVideo() || (!enable_heavy_rendering && !(enable_magnify_on_keyboard_focus && is_keyboard_focused))}
            <!-- If not a video or heavy rendering is not enabled, neither globaly or just for this item(like when keyboard focused and magnify enabled) -->
            {#if ordered_media.Media.isAnimated() && (enable_magnify_on_keyboard_focus && is_keyboard_focused || enable_heavy_rendering)}
                <!-- Is animated -->
                <img src="{ordered_media.Media.Url}" alt="{ordered_media.Media.MediaName}" loading="lazy">
            {:else}
                <!-- Is not animated -->
                <img src="{ordered_media.Media.getResizedUrl(media_viewport_percentage)}" alt="{ordered_media.Media.MediaName}" loading="lazy">
            {/if}
        {:else if ordered_media.Media.isVideo() && (enable_magnify_on_keyboard_focus && is_keyboard_focused || enable_heavy_rendering)}
            {#if media_inside_viewport || !enable_heavy_rendering}
                <video 
                    src="{ordered_media.Media.Url}" 
                    muted 
                    autoplay
                    preload="metadata"
                    loop
                >
                </video>
            {/if}
        {/if}
    {/if}        
</div>

<style>
    .meg-display-item-wrapper {
        --display-item-border-radius: var(--border-radius);
        position: relative;
        width: 100cqw;
        height: 110cqw;
        background: var(--main-dark-color-7);
        border-radius: var(--display-item-border-radius);
        transition: scale 0.2s ease-out, translate 0.4s ease-out, box-shadow 0.1s linear allow-discrete;
    }

    .meg-display-item-wrapper.status-magnified {
        scale: 1.5;
        translate: 0 -20%;
        border: 1.2px solid var(--grey-8);
        background-color: black;
        box-shadow: 0px 2cqh 24px 20px hsl(from black h s l / .5);
        border: none;
        

        & img {
            object-fit: contain;
        }

        & video {
            object-fit: contain;
        }
    }    
    
    /*=============================================
    =            Status overlay            =
    =============================================*/

    
        .meg-display-item-wrapper .media-status-overlay {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            opacity: 0;
            transition: opacity 0.2s ease-out, background 0.3s ease-out;
            user-select: none;  
            z-index: var(--z-index-t-5);
        }
    
        .meg-display-item-wrapper.status-deleted .media-status-overlay {
            background: hsl(from var(--danger-8) h s l / 0.7);
            opacity: 1;
        }

        .meg-display-item-wrapper.status-selected .media-status-overlay {
            background: hsl(from var(--main-4) h s l / 0.7);
            opacity: 1;
        }

        .meg-display-item-wrapper.status-yanked {
            opacity: 0.3;
            filter: grayscale(1);
        }
    
    /*=====  End of Status overlay  ======*/

    
    /*=============================================
    =            File name and type            =
    =============================================*/
    
        .meg-display-item-wrapper .media-type-label-overlay {
            box-sizing: border-box;
            position: absolute;
            display: flex;
            inset: auto 0 0 auto;
            width: 100%;
            height: max-content;
            user-select: none;
            justify-content: flex-end;
            padding: var(--spacing-2);
            z-index: var(--z-index-t-4);
        }

        .media-type-label-overlay h4.media-type-label {
            padding: var(--spacing-1);
            background-color: var(--main-dark);
            color: var(--grey-9);
            border-radius: var(--border-radius);
            line-height: 1;

            &.media-type-video {
                background-color: var(--accent);
            }
        }

        .status-magnified .media-type-label-overlay h4.media-type-label {
            opacity: 0;
        }

        @supports (color: rgb(from white r g b)) {
            .media-type-label-overlay h4.media-type-label:not(.media-type-video) {
                background-color: hsl(from var(--main-dark) h s l / 0.3);
            }
        }

        
        /*----------  name tooltip  ----------*/

            .meg-display-item-wrapper .media-name-tool-tip {
                position: absolute;
                display: none;
                inset: 0 auto auto 0;
                width: max-content;
                max-width: 100%;
                transition: display .3s linear allow-discrete;
            }

            .meg-display-item-wrapper.status-titles-enabled .media-name-tool-tip {
                display: block;

                & p.media-name-tool-tip-content {
                    width: 96cqw;
                    font-size: calc(var(--font-size-1) * 0.9);
                    color: white;
                    line-height: 1.2;
                    text-align: center;
                    translate: 2px 5%;
                    scale: 1;
                }
            }



            .meg-display-item-wrapper:not(.status-titles-enabled):hover .media-name-tool-tip {
                display: block;
            }

            .meg-display-item-wrapper .media-name-tool-tip p.media-name-tool-tip-content {
                color: var(--grey-1);
                font-size: var(--font-size-fineprint);
                line-height: 1;
                transform-origin: center 0;
                scale: 0;
                translate: 0 -100%;
                padding: var(--spacing-1);
                background-color: var(--grey-9);
                word-break: break-all;
                transition: scale .3s ease-out;
            }

            .meg-display-item-wrapper:hover .media-name-tool-tip p.media-name-tool-tip-content {
                scale: 1;
            }

            @supports (color: rgb(from white r g b)) {
                .meg-display-item-wrapper .media-name-tool-tip p.media-name-tool-tip-content {
                    background-color: hsl(from var(--grey) h s l / 0.9);
                }

                .meg-display-item-wrapper.status-titles-enabled .media-name-tool-tip p.media-name-tool-tip-content {
                    background-color: hsl(from var(--grey) h s l / .8);
                }
            }
        
    /*=====  End of File name and type  ======*/
    
    

    .meg-display-item-wrapper.use-masonry {
        height: 100%;
    }

    .meg-display-item-wrapper img, .meg-display-item-wrapper video {
        max-width: 150cqw;
        width: 100%;
        height: 100%;
        object-fit: cover;
        object-position: center;
        border-radius: var(--display-item-border-radius);
        user-select: none;
    }

    @media(pointer: fine) {
        .meg-display-item-wrapper {
            cursor: pointer;
        }
    }
</style>