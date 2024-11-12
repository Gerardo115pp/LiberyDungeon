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
        $: if (ordered_media && $me_gallery_changes_manager != null) {
            refreshMediaChangesSubscription();
        }

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
        
        /**
         * Whether the item should be draggable or not.
         * @type {boolean}
         * @default false
         */
        export let enable_dragging = false;

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

            /*----------  Drag operations  ----------*/
                
                /**
                 * Whether another media item is being dragged over this media item.
                 * @type {boolean}
                 */
                let media_dragged_over = false;


                /**
                 * Whether this media item is being dragged.
                 * @type {boolean}
                 */
                let is_dragged = false;
        
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

        
            /*----------  container border modifiers.  ----------*/

                /**
                 * Is on the top of it's container.
                 * @type {boolean}
                 */
                let is_on_top = false;

                /**
                 * Is on the bottom of it's container.
                 * @type {boolean}
                 */
                let is_on_bottom = false;

                /**
                 * Is on the left of it's container.
                 * @type {boolean}
                 */
                let is_on_left = false;

                /**
                 * Is on the right of it's container.
                 * @type {boolean}
                 */
                let is_on_right = false;

                /**
                 * Whether to check if the media item is on bordering limits of it's container. if enabled and on
                 * the media will auto adjust it's position to be fully visible when magnified.
                 * @type {boolean}
                 * @default false
                 */
                export let check_container_limits = false;
                $: if (this_dom_element && check_container_limits && enable_magnify_on_keyboard_focus && is_keyboard_focused) {
                    defineElementPositionModifiers();
                }

                /**
                 * The selector to get the container element. 
                 * @type {string}
                 * @default ":has(> .meg-display-item-wrapper)"
                 */
                export let container_selector = ":has(> .meg-display-item-wrapper)";
        
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
        
        /*=============================================
        =            Drag handlers            =
        =============================================*/
        
            /**
             * Handles the drag start event.
             * @param {DragEvent} event
             */
            const handleMediaItemDragStart = event => {
                is_dragged = true;

                // @ts-ignore - According to the MDN docs, this is not possibly null in a DragEvent for dragstart.
                event.dataTransfer.effectAllowed = 'move';

                console.log("Dragged media: ", ordered_media.Order);
            } 

            /**
             * Handles the drag end event.
             * @param {DragEvent} event
             */
            const handleMediaItemDragEnd = event => {
                is_dragged = false;
                console.log(`Media item ${ordered_media.Order} drag end`);

                // @ts-ignore - According to the MDN docs, this is not possibly null in a DragEvent for dragend.
                let effect = event.dataTransfer.dropEffect;

                console.log(`Drop effect: ${effect}`);
            }

            /**
             * Handles the drag enter event.
             * @param {DragEvent} event
             */
            const handleMediaItemDragEnter = event => {
                if (is_dragged) return;
                event.preventDefault();
                media_dragged_over = true;
            }

            /**
             * Handles the drag over event.
             * @param {DragEvent} event
             */
            const handleMediaItemDragOver = event => {
                if (is_dragged) return;
                event.preventDefault();

                // @ts-ignore - According to the MDN docs, this is not possibly null in a DragEvent for dragover.
                event.dataTransfer.dropEffect = 'move';

                media_dragged_over

                return true;
            }

            /**
             * Handles the drag leave event.
             * @param {DragEvent} event
             */
            const handleMediaItemDragLeave = event => {
                media_dragged_over = false;
            }

            /**
             * Handles the drop event.
             * @param {DragEvent} event
             */
            const handleMediaItemDrop = event => {
                event.preventDefault();

                media_dragged_over = false;

                console.log("Inserting medias before media: ", ordered_media.Order);
            }
        
        /*=====  End of Drag handlers  ======*/

        /**
         * Defines the element position modifiers based on the container limits.
         * only runs if check_container_limits is set to true.
         */
        function defineElementPositionModifiers() {
            if (!check_container_limits) return;

            let container_element = this_dom_element.closest(container_selector);

            if (container_element == null) {
                console.warn("Could not find the container element");
                return;
            }

            let container_rect = container_element.getBoundingClientRect();
            let element_rect = this_dom_element.getBoundingClientRect();

            const width_distance_threshold = element_rect.width * 0.9;
            console.log("Width distance threshold: ", width_distance_threshold);
            const height_distance_threshold = element_rect.height * 0.9;
            console.log("Height distance threshold: ", height_distance_threshold);

            is_on_top = (element_rect.top - container_rect.top) <= height_distance_threshold;
            is_on_bottom = (container_rect.bottom + height_distance_threshold) <= element_rect.bottom;
            is_on_left = (container_rect.left + width_distance_threshold) >= element_rect.left;
            is_on_right = (container_rect.right - width_distance_threshold) <= element_rect.right;
        }
    
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

            if ($me_gallery_changes_manager == null) {
                throw new Error("In MEGalleryDisplayItem.handleMediaChanges: me_gallery_changes_manager is null");
            }
            
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
                    // @ts-ignore
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
         * Cancels current media changes suscription and creates a new one.
         * @requires me_gallery_changes_manager
         * @requires handleMediaChanges
         */
        const refreshMediaChangesSubscription = () => {
            unsuscribeFromMediaChanges();
            suscribeToMediaChanges();
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
            $me_gallery_changes_manager.unsubscribeToChanges(element_media_uuid);   
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<div bind:this={this_dom_element} class="meg-display-item-wrapper"
    class:use-masonry={use_masonry}
    class:meg-di-keyboard-focused={is_keyboard_focused}
    class:status-magnified={enable_magnify_on_keyboard_focus && is_keyboard_focused}
    class:is-on-top-limit={is_on_top}
    class:is-on-bottom-limit={is_on_bottom}
    class:is-on-left-limit={is_on_left}
    class:is-on-right-limit={is_on_right}
    class:status-deleted={is_media_deleted}
    class:status-selected={is_media_selected}   
    class:status-yanked={is_media_yanked}
    class:status-titles-enabled={enable_video_titles}
    class:status-media-hovering={media_dragged_over}
    class:status-media-dragging={is_dragged}
    draggable="{enable_dragging}"
    on:dragstart={handleMediaItemDragStart}
    on:dragend={handleMediaItemDragEnd}
    on:dragenter={handleMediaItemDragEnter}
    on:dragover={handleMediaItemDragOver}
    on:dragleave={handleMediaItemDragLeave}
    on:drop={handleMediaItemDrop}
    data-media-order={ordered_media.Order}
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
                <!-- @see: https://svelte.dev/docs/basic-markup#attributes-and-props
                    '''
                        Sometimes, the attribute order matters as Svelte sets attributes sequentially in JavaScript.....Another example 
                        is <img src="..." loading="lazy" />. Svelte will set the img src before making the img element loading="lazy", which is probably too late.
                        Change this to <img loading="lazy" src="..."> to make the image lazily loaded.
                    '''
                    So, the loading attribute should be set before the src attribute.
                    This gave an amazing performance boost to the gallery.
                -->
                <img loading="lazy" src="{ordered_media.Media.Url}" alt="{ordered_media.Media.MediaName}">
            {:else}
                <!-- Is not animated -->
                <img loading="lazy" src="{ordered_media.Media.getResizedUrl(media_viewport_percentage)}" alt="{ordered_media.Media.MediaName}">
            {/if}
        {:else if ordered_media.Media.isVideo() && (enable_magnify_on_keyboard_focus && is_keyboard_focused || enable_heavy_rendering)}
            {#if media_inside_viewport || !enable_heavy_rendering}
            <video 
                    preload="metadata"
                    src="{ordered_media.Media.Url}" 
                    muted 
                    autoplay
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
        background: var(--grey-8);
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
        z-index: var(--z-index-t-3);
        

        & img {
            object-fit: contain;
        }

        & video {
            object-fit: contain;
        }

        /* top limit */
        &.is-on-top-limit {
            translate: 0 30%;
        }

        /* bottom limit */
        &.is-on-bottom-limit {
            translate: 0 -30%;
        }

        /* left limit */
        &.is-on-left-limit {
            translate: 30% 0;
        }

        /* right limit */
        &.is-on-right-limit {
            translate: -30% 0;
        }

        /* top-left limit */
        &.is-on-top-limit.is-on-left-limit {
            translate: 30% 30%;
        }

        /* top-right limit */
        &.is-on-top-limit.is-on-right-limit {
            translate: -30% 30%;
        }

        /* bottom-left limit */
        &.is-on-bottom-limit.is-on-left-limit {
            translate: 30% -30%;
        }

        /* bottom-right limit */
        &.is-on-bottom-limit.is-on-right-limit {
            translate: -30% -30%;
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
    
    
    /*=============================================
    =            Dragging oprations feedback            =
    =============================================*/
    
        :global(:has(> .meg-display-item-wrapper:not(.status-media-dragging)) ~ :has(> .meg-display-item-wrapper.status-media-hovering)) {
            & img, & video {
                translate: 15% 0;
                transition: translate 0.3s ease-out;
            }
        }

        .meg-display-item-wrapper.status-media-dragging {
            opacity: 0;
            transition: opacity 0.3s ease-out;
        }
            
    
    /*=====  End of Dragging oprations feedback  ======*/
    
    

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