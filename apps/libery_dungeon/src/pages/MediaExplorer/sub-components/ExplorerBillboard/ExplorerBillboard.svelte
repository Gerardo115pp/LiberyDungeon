<script>
    import { browser } from '$app/environment';
    import LiberyHeadline from '@components/UI/LiberyHeadline.svelte';
    import { CanvasImage } from '@libs/LiberyColors/image_models';
    import { linearCycleNavigationWrap } from '@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils';
    import { getTaggedMedias } from '@models/Medias';
    import { current_cluster } from '@stores/clusters';
    import { navbar_ethereal } from '@stores/layout';
    import { onMount, onDestroy, tick } from 'svelte';
    import { cubicIn } from 'svelte/easing';
    import { fade } from 'svelte/transition';

    /*=============================================
    =            Properties            =
    =============================================*/

        /**
        * @typedef {Object} MediaSize
         * @property {number} width
         * @property {number} height
        */
    
        /**
         * The category to create the billboard for.
         * @type {import('@models/Categories').CategoryLeaf}
         */ 
        export let the_billboard_category; 

        /**
         * The current media showned in the billboard.
         * @type {import('@models/Medias').Media | null}
         */
        let current_billboard_media = null;

        /**
         * A media object that is been prefetched.
         * Used to create blending transitions.
         * @type {import('@models/Medias').Media | null}
         */
        let prefetched_billboard_media = null;

        /**
         * A callback triggered when the current_billboard_media changes to a valid media.
         * @type {() => void}
         */
        export let onvalid_media_change = () => {};
        
        /*----------  Style  ----------*/
        
            /**
             * The theme color to use fort the overlay content(the text). It's value is determined by the color pallette of the current_billboard_media.
             * @type {boolean}
             * @default false - White overlay color theme is best for MOST cases except when the media is extremely white.
             */ 
            let use_dark_theme = false;

            /**
             * Whether the natural dimensions of the current b
             * billboard media are smaller than the viewport. 
             * This is used to apply different styles like blur 
             * to minimize artifacts.
             * @type {boolean}
             */
            let current_media_is_small = false;

            /**
             * How much a small image should be blurred. this is in
             * pixels.
             * @type {number}
             * @default 1.2
             */
            let small_image_blur = 1.2;

            /**
             * The time the fade  out animation of the 
             * prefetched media last for.
             * @type {number}
             */
            const PREFETCH_FADE_OUT_DURATION = 400;
        
        /*----------  State  ----------*/

            /**
             * A list of medias that can be displayed on the billboard and are related to the current category.
             * @type {import('@models/Medias').Media[]}
             */
            let billboard_medias = [];

            /**
             * The index iterator for the billboard_medias.
             * @type {number}
             * @default 0
             */
            let billboard_iterator_index = 0;

            /**
             * Whether to use random iteration for the billboard medias.
             * @type {boolean}
             * @default false
             */
            let random_billboard_media_iteration = false;
        
            /**
             * Whether the billboard has mounted.
             * @type {boolean}
             */ 
            let component_mounted = false;
            $: if (the_billboard_category != null && component_mounted) {
                handleCategoryChange(the_billboard_category);
            }

            /**
             * Whether the billboard media has a vertical aspect ratio.
             * @type {boolean}
             */
            let billboard_media_vertical = false;

            /**
             * Whether a media i been prefetched.
             * @type {boolean}
             */
            let media_prefetching = false;
            
            /*----------  Looping interval  ----------*/

                /**
                 * The interval that controls the looping of the billboard media.
                 * @type {number}
                 */
                let billboard_media_looping_interval = NaN;

                /**
                 * The time that has to elapse before the billboard media changes.
                 * @type {number}
                 */
                const LOOPING_INTERVAL_TIME = 120000;

                /**
                 * Whether the billboard media looping interval is enabled.
                 * @type {boolean}
                 */
                let billboard_media_looping_enabled = false;
                

    
    /*=====  End of Properties  ======*/

    onMount(() => {
        checkWindowScroll();

        component_mounted = true;
    });

    onDestroy(() => {
        if (!browser) return;

        handleBillboardDestroy();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Called by the billboard media looping interval when is invoked.
         * @returns {void}
         */
        const billboardMediaLoopingIntervalCallback = () => {
            if (billboard_medias.length === 0 || media_prefetching) return;

            console.log("Looping interval called.")

            changeBillboardImage(true);
        }

        /**
         * checks if the scrollY is below a threshold and if so, turns the navbar ethereal.
         */
        const checkWindowScroll = () => {
            if (globalThis.self == null) return; // not a windowed context.

            if (window.scrollY < 100) {
                navbar_ethereal.set(true);
            } else {
                navbar_ethereal.set(false);
            }
        }

        /**
         * Calculates the white percentage of the passed Media element(HTML).
         * @param {HTMLImageElement | HTMLVideoElement} media_element
         * @returns {number}
         */
        const calculateMediaWhitePercentage = (media_element) => {
            console.time("Calculating white percentage of the media element.");
            const canvas_image = new CanvasImage(media_element);
            console.timeEnd("Calculating white percentage of the media element.");

            return canvas_image.whitePercentage(3);
        }

        /**
         * Determines if the current category has any media that can be displayed on the billboard.
         * @param {import('@models/Categories').CategoryLeaf} new_category
         * @returns {Promise<boolean>}
         */
        const categoryHasBillboardMedia = async (new_category) => {
            if (new_category == null) return false;

            /**
             * @type {import('@models/Categories').CategoryConfig | null}
             */
            let category_configuration = new_category.Config;

            if (category_configuration === null) {
                category_configuration = await new_category.loadCategoryConfig();
            }

            if (category_configuration == null) {
                return new_category.content.length >= 1
            }

            let configuration_has_medias = category_configuration.BillboardMediaUUIDs.length > 0 || category_configuration.BillboardDungeonTags.length > 0;

            return configuration_has_medias || new_category.content.length > 0 
        }

        /**
         * Changes the current billboard media.
         * @param {boolean} [avoid_same_image] if true and uses random navigation, it will avoid landing in the same index.
         */ 
        const changeBillboardImage = (avoid_same_image) => {
            /**
             * @type {import('@models/Medias').Media | null}
             */
            let new_billboard_media = null;

            if (avoid_same_image && current_billboard_media != null && billboard_medias.length === 1) {
                resetBillboardMediaIteration();
                new_billboard_media = billboard_medias[billboard_iterator_index];

                prefetchBillboardMedia(new_billboard_media);
                return;
            }



            if (random_billboard_media_iteration) {
                new_billboard_media = iterBillboardMediasRandom();
            } else {
                new_billboard_media = iterBillboardMediasSequential();
            }

            if (new_billboard_media == null) return;

            prefetchBillboardMedia(new_billboard_media);
        }

        /**
         * Clears the billboard media looping interval.
         * @returns {void}
         */
        const disableBillboardMediaLooping = () => {
            billboard_media_looping_enabled = false;

            if (isNaN(billboard_media_looping_interval)) return;

            clearInterval(billboard_media_looping_interval);

            billboard_media_looping_interval = NaN;
        }

        /**
         * Enables the billboard media looping interval.
         * @returns {void}
         */
        const enableBillboardMediaLooping = () => {
            if (billboard_media_looping_enabled || !isNaN(billboard_media_looping_interval)) {
                disableBillboardMediaLooping();

                billboard_media_looping_enabled = true;
            }

            random_billboard_media_iteration = billboard_medias.length > 2;


            billboard_media_looping_interval = setInterval(billboardMediaLoopingIntervalCallback, LOOPING_INTERVAL_TIME);
        }

        /**
         * Returns the html billboard element.
         * @returns {HTMLImageElement | HTMLVideoElement | null}
         */
        const getBillboardElement = () => {
            if (current_billboard_media == null) return null;

            const billboard_element = document.getElementById(`mexbill-billboard-${current_billboard_media.isVideo() ? "video" : "image"}`);

            if (billboard_element == null) return billboard_element;

            if (!(billboard_element instanceof HTMLImageElement || billboard_element instanceof HTMLVideoElement)) {
                throw new Error("In ExplorerBillboard.getBillboardElement: The element with id 'mexbill-billboard' is not an instance of HTMLImageElement or HTMLVideoElement.");
            }

            return billboard_element;
        }

        /**
         * Returns the size of the current billboard media.
         * @param {HTMLImageElement | HTMLVideoElement} [billboard_element]
         * @returns {MediaSize}
         */
        const getBillboardMediaSize = (billboard_element) => {

            /** @type {MediaSize} */
            const media_size = {
                width: 0,
                height: 0
            }


            if (billboard_element == null) {
                // @ts-ignore - i really dont care whether it is null or undefined...................  fing typescript.
                billboard_element = getBillboardElement();

                if (billboard_element == null) return media_size;
            } 

            if (billboard_element instanceof HTMLImageElement) {
                media_size.width = billboard_element.naturalWidth;
                media_size.height = billboard_element.naturalHeight;
            } else if (billboard_element instanceof HTMLVideoElement) {
                media_size.width = billboard_element.videoWidth;
                media_size.height = billboard_element.videoHeight;
            }

            return media_size;
        }
        
        /**
         * handles the current category change.
         * @param {import('@models/Categories').CategoryLeaf} new_category
         */
        async function handleCategoryChange(new_category) {
            const was_current_media_null = current_billboard_media == null;

            if (new_category == null) return;

            if (await categoryHasBillboardMedia(the_billboard_category)) {
                await loadBillboardMedias(the_billboard_category.Config)
                changeBillboardImage(true); 

            } else {
                disableBillboardMediaLooping();
            }

            if (was_current_media_null && current_billboard_media != null) {
                onvalid_media_change();
            };
        }

        /**
         *  Handles an onscroll event listener on the document window.
         */
        const handleWindowScroll = () => {
            checkWindowScroll();
        }
    
        /**
         * Handles the destruction of the billboard component.
         */
        const handleBillboardDestroy = () => {
            navbar_ethereal.set(false);

            disableBillboardMediaLooping();
        }

        /**
         * Handles the load event of the billboard Image.
         * @param {Event} event
         */
        const handleBillboardImageLoad = (event) => {
            optimizeBillboardDisplay();
        }

        /**
         * Handles the load event of the billboard Video.
         * @param {Event} event
         */
        const handleBillboardVideoLoad = (event) => {
            optimizeBillboardDisplay();
        }

        /**
         * Handles the error event of the billboard Image.
         * @param {Event} event
         */
        const handleBillboardImageError = (event) => {
            console.log("Error loading billboard image: ", event);
        }

        /**
         * Handles the error event of the billboard Video.
         * @param {Event} event
         */
        const handleBillboardVideoError = (event) => {
            console.log("Error loading billboard video: ", event);
        }

        /**
         * Hides the prefecthing media.
         * @returns {void}
         */
        const hidePrefetchedMedia = () => {
            setTimeout(() => {
                media_prefetching = false;
                console.log("Hiding prefetched media.");
            
            }, PREFETCH_FADE_OUT_DURATION * 1.5);
        }

        /**
         * Iters through medias using a random progression system.
         * @param {boolean} [avoid_same_image] if true, it will avoid landing in the same index.
         * @returns {import('@models/Medias').Media | null}
         */
        const iterBillboardMediasRandom = (avoid_same_image) => {
            let inifite_loop_guard = 0;

            /**
             * @type {import('@models/Medias').Media | null}
             */
            let new_billboard_media = null;

            while (new_billboard_media == null) {
                const rand_media_index = Math.trunc(Math.random() * billboard_medias.length);

                new_billboard_media = billboard_medias[rand_media_index];

                if (current_billboard_media != null && current_billboard_media.uuid === new_billboard_media.uuid && avoid_same_image) {
                    new_billboard_media = null
                }

                inifite_loop_guard++;

                if (inifite_loop_guard > billboard_medias.length) {
                    console.error(`In ExplorerBillboard.handleChangeBillboardImage: Infinite loop detected after ${billboard_medias.length} iterations. Exiting.`);
                    return null;
                }
            }

            return new_billboard_media;
        }

        /**
         * iterates over the billboard medias using a sequential progression system. When it reaches the limit of medias, it cycles.
         * @returns {import('@models/Medias').Media | null}
         */
        const iterBillboardMediasSequential = () => {
            if (billboard_medias.length === 0) return null;

            const next_index = linearCycleNavigationWrap(billboard_iterator_index, billboard_medias.length - 1, 1).value;

            billboard_iterator_index = next_index;
            return billboard_medias[billboard_iterator_index];
        }

        /**
         * Returns whether a given set of width and height values describes a vertical aspect ratio.
         * @param {number} width
         * @param {number} height
         * @returns {boolean}
         */
        const isVerticalAspectRatio = (width, height) => {
            return height >= (width * 1.1);
        }

        /**
         * Loads the billboard medias from a category configuration.
         * @param {import('@models/Categories').CategoryConfig | null} category_configuration
         * @returns {Promise<void>}
         */
        const loadBillboardMedias = async (category_configuration) => {
            if (category_configuration == null) {
                billboard_medias = the_billboard_category.content;
                resetBillboardMediaIteration(true);
                return;
            }

            /**
             * @type {import("@models/Medias").Media[]}
             */
            let new_billboard_medias = [];
            let configuration_had_medias = false;

            if (category_configuration.BillboardDungeonTags.length !== 0) {
                // @ts-ignore
                new_billboard_medias = await loadBillboardMediasFromTags(category_configuration.BillboardDungeonTags);
                configuration_had_medias = true;
                enableBillboardMediaLooping();
            } else if (category_configuration.BillboardMediaUUIDs.length > 0) {
                new_billboard_medias = await loadBillboardMediasFromUuids(category_configuration.BillboardMediaUUIDs);
                configuration_had_medias = true;
                enableBillboardMediaLooping();
            } else {
                new_billboard_medias = the_billboard_category.content;
                disableBillboardMediaLooping()
            }

            billboard_medias = new_billboard_medias;
            resetBillboardMediaIteration(true); // Always start with random navigation.
        }

        /**
         * Called by loadBillboardMedias. Loads the billboard medias from a category configuration using the list of media uuids.
         * @param {string[]} media_uuids
         * @returns {Promise<import('@models/Medias').Media[]>}
         */
        const loadBillboardMediasFromUuids = async (media_uuids) => {
            if ($current_cluster == null || media_uuids.length === 0) return [];

            const billboard_media_identites = await $current_cluster.getClusterMedias(media_uuids);

            let new_billboard_medias = billboard_media_identites.map(mi => mi.Media);

            return new_billboard_medias;
        }

        /**
         * Called by loadBillboardMedias. Loads the billboard medias from a category configuration using the list of media tags.
         * @param {number[]} media_tags
         * @returns {Promise<import('@models/Medias').Media[]>}
         */
        const loadBillboardMediasFromTags = async (media_tags) => {
            const page_content = await getTaggedMedias(media_tags, 1, 200) // TODO: Move the limit and page_num to constants.

            if (page_content === null) {
                return [];
            }

            const new_billboard_medias = page_content.content.map(mi => mi.Media);

            return new_billboard_medias;
        }            

        /**
         * Optimizies the billboard display to render the image on the most astetically pleasing way possible.
         * @returns {void}
         */
        const optimizeBillboardDisplay = () => {
            const billboard_element = getBillboardElement();

            if (current_billboard_media == null || billboard_element == null) return;

            const media_size = getBillboardMediaSize();

            verifyBillboardMediaAspectRatio(media_size);            

            verifyBillboardMediaIsSmall(media_size, current_billboard_media.isVideo());

            verifyBillboardMediaWhitePercentage(billboard_element);
        }

        /**
         * Prefetches the given media. Don't call this directly. should only be called by 'changeBillboardImage'.
         * @param {import('@models/Medias').Media} next_media
         * @returns {Promise<void>}
         */
        const prefetchBillboardMedia = async (next_media) => {
            if (next_media == null) return;

            prefetched_billboard_media = current_billboard_media;

            current_billboard_media = null;
            media_prefetching = true;

            console.log("Prefetching media: ", next_media);

            await tick();

            if (next_media.isVideo()) {
                const preloader_video_element = document.createElement("video");

                preloader_video_element.src = next_media.Url;

                preloader_video_element.muted = true;

                preloader_video_element.autoplay = true;

                preloader_video_element.onloadeddata = () => {
                    current_billboard_media = next_media;
                    hidePrefetchedMedia();
                }
            } else if (next_media.isImage()) {
                const preloader_image_element = document.createElement("img");

                preloader_image_element.src = next_media.Url;

                preloader_image_element.onload = () => {
                    current_billboard_media = next_media;
                    hidePrefetchedMedia();
                }
            }
        }

        /**
         * Resets the billboard iteration.
         * @param {boolean} [use_random_navigation]
         * @returns {void}
         */
        const resetBillboardMediaIteration = (use_random_navigation) => {
            billboard_iterator_index = 0;
            random_billboard_media_iteration = use_random_navigation === true;
        }

        /**
         * sets the enables/disables random navigation of billboard medias.
         * @param {boolean} enable
         * @returns {void}
         */
        const setRandomBillboardMediaIteration = (enable) => {
            random_billboard_media_iteration = enable;
        }

        /**
         * Updates the current category configuration. meant to be used from a parent component.
         * @param {import('@models/Categories').CategoryConfig} new_config 
         */
        export async function updateCurrentCategoryConfig(new_config) {
            if (new_config.CategoryUUID != the_billboard_category.uuid) return;

            the_billboard_category.setCategoryConfig(new_config);

            if (await categoryHasBillboardMedia(the_billboard_category)) {
                await loadBillboardMedias(new_config);
                resetBillboardMediaIteration(false);
            }

            console.log("Updated category configuration for the billboard.");
        
        }

        /**
         * Verifies the current billboard media aspect ratio.
         * @param {MediaSize} media_size
         * @returns {void}
         */
        const verifyBillboardMediaAspectRatio = (media_size) => {
            billboard_media_vertical = isVerticalAspectRatio(media_size.width, media_size.height);
        }

        /**
         * Verifies the current billboard media size in relation to the
         * viewport and determine if it is to small to be displayed normally.
         * @param {MediaSize} media_size
         * @param {boolean} is_video
         * @returns {void}
         */
        const verifyBillboardMediaIsSmall = (media_size, is_video) => {
            current_media_is_small =  media_size.width < (window.innerWidth * 0.8);

            if (!current_media_is_small) {
                small_image_blur = 0;
                return;
            }

            const min_blur = 1.2;
            const max_blur = is_video ? 2 : 6;

            const media_to_viewport_ratio = media_size.width / window.innerWidth;

            switch (true) {
                case media_to_viewport_ratio <= 0.2:
                    small_image_blur = max_blur;
                    break;
                case media_to_viewport_ratio <= 0.3:
                    small_image_blur = Math.max(max_blur * 0.5, min_blur);
                    break;
                default:
                    small_image_blur = min_blur;
            }

            return;
        }

        /**
         * Verifies the white percentage of the current billboard media and sets the use_dark_theme property accordingly.
         * @param {HTMLImageElement | HTMLVideoElement} billboard_element
         * @returns {void}
         */
        const verifyBillboardMediaWhitePercentage = (billboard_element) => {
            if (billboard_element == null || billboard_element instanceof HTMLVideoElement) {
                use_dark_theme = false;
                return;
            }

            let new_white_percentage = calculateMediaWhitePercentage(billboard_element);

            use_dark_theme = new_white_percentage > 18;           
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<svelte:window 
    on:scroll|passive={handleWindowScroll}
/>
<section id="media-explorer-billboard"
    class:loaded-billboard={current_billboard_media != null}
    class:billboard-looping-enabled={billboard_media_looping_enabled}
    class:billboard-media-vertical={billboard_media_vertical}
    class:billboard-media-small={current_media_is_small}
    class:dark-overlay-color-theme={use_dark_theme}
>
    <div id="mexbill-underlay-billboard-wrapper"
        style:--small-image-blur="{small_image_blur}px"
    >
        {#if current_billboard_media != null}
            {#if current_billboard_media.isImage()}
                <img
                    id="mexbill-billboard-image"
                    class="billboard-media"
                    decoding="async"
                    src="{current_billboard_media.Url}" 
                    on:load={handleBillboardImageLoad}
                    on:error={handleBillboardImageError}
                    alt=""
                >
            {:else if current_billboard_media.isVideo()}
                <video 
                    id="mexbill-billboard-video"
                    class="billboard-media"
                    src="{current_billboard_media.Url}"
                    on:loadeddata={handleBillboardVideoLoad}
                    on:error={handleBillboardVideoError}
                    muted
                    autoplay
                    loop
                >
                    <track kind="captions">
                </video>
            {/if}
        {/if}
        {#if prefetched_billboard_media != null && media_prefetching}
            {#if prefetched_billboard_media.isImage()}
                <img
                    id="mexbill-prefeched-billboard"
                    class="billboard-media billboard-media-prefetched"
                    decoding="async"
                    src="{prefetched_billboard_media.Url}" 
                    style:animation-duration="{PREFETCH_FADE_OUT_DURATION}ms"
                    alt=""
                >
            {:else if prefetched_billboard_media.isVideo()}
                <video 
                    id="mexbill-prefeched-billboard"
                    class="billboard-media billboard-media-prefetched"
                    style:animation-duration="{PREFETCH_FADE_OUT_DURATION}ms"
                    src="{prefetched_billboard_media.Url}"
                    muted
                    autoplay
                    loop
                >
                    <track kind="captions">
                </video>
            {/if}
        {/if}
    </div>
    <div id="mexbill-synopsis-panel">
        <hgroup id="mexbill-synpa-headlines">
            <slot name="billboard-headline" >
                <LiberyHeadline 
                    headline_tag="name"
                    headline_color="var(--grey-1)"
                    headline_font_size="var(--font-size-h3)"
                    extra_props="style='awesome!'"
                    headline_text={the_billboard_category.name}
                    force_bottom_lines
                />
            </slot>
        </hgroup>
    </div>
</section>

<style>



    section#media-explorer-billboard {
        --shadows-color: hsl(from var(--body-bg-color) h s l / 0.85);
        /* Style custom-properites for the LiberyHeadline */
        --text-color-1: var(--grey-1);
        --text-color-2: var(--grey-3);
        --text-shadow: none;
        
        /* position: relative; */
        display: flex;
        width: 100%;
        height: 300px;
        container-type: size;
        flex-direction: column;
        justify-content: flex-end;
        z-index: var(--z-index-b-1);
    }
    
    section#media-explorer-billboard.loaded-billboard {
        height: 50dvh;        
    }

    section#media-explorer-billboard.dark-overlay-color-theme {
        --text-shadow: 0 -4px 58px hsl(from var(--body-bg-color) h s l / 0.8);
        --text-color-1: var(--main-5);
        --text-color-2: var(--grey);
    }
    
    /*=============================================
    =            Billboard            =
    =============================================*/
    
        #mexbill-underlay-billboard-wrapper {
            position: absolute;
            inset: 0;
            width: 100dvw;
            height: 200cqh;
            z-index: var(--z-index-b-5);
            /* overflow: hidden; */

            & .billboard-media {
                position: absolute;
                inset: 0;
                width: 100%;
                height: 100%;
                object-fit: cover;
                object-position: center 25%;
                z-index: var(--z-index-b-6);
            }

            & .billboard-media-prefetched.billboard-media {
                z-index: var(--z-index-b-4);
                animation-name: billboard-fade-out;
                animation-fill-mode: forwards;
            }

            &::after {
                content: "";
                position: absolute;
                inset: auto auto -2% 0;
                width: 100%;
                background: linear-gradient(to top, hsl(from var(--body-bg-color) h s l / 0.99) 10%, hsl(from var(--body-bg-color) h s l / 0.001)) ;
                height: 30cqh;
            }

            &::before {
                content: "";
                position: absolute;
                inset: -3% auto auto 0;
                width: 100%;
                /* background: var(--body-bg-color); */
                background: linear-gradient(to bottom, hsl(from var(--body-bg-color) h s l / 0.99) 40%, hsl(from var(--body-bg-color) h s l / 0.01) 90%,  transparent);
                height: 7.5cqh;
            }
        }  

        section#media-explorer-billboard.billboard-media-vertical #mexbill-underlay-billboard-wrapper {

            height: 250cqh;

            & .billboard-media {
                object-position: center 30%;
            }
        }

        section#media-explorer-billboard.billboard-media-small #mexbill-underlay-billboard-wrapper {
            &  .billboard-media {
                filter: blur(var(--small-image-blur));
            }
        }
        
    /*=====  End of Billboard  ======*/
        
    /* -------------------------------- Synopsis -------------------------------- */

        section#media-explorer-billboard.loaded-billboard #mexbill-synopsis-panel {
            padding-inline: var(--spacing-5);
        }
        
        #mexbill-synopsis-panel {
            /* --synopsis-panel-bg: hsl(from var(--grey-black) h s l / 0.25); */
            
            width: 100%;
            /* background: var(--synopsis-panel-bg); */
            padding-block-start: var(--spacing-2);
            padding-block-end: var(--spacing-4);
            padding-inline: var(--common-page-inline-padding);
            /* translate: calc(-1 * var(--common-page-inline-padding)); */
            box-shadow: 0 -10px 36px 40px var(--synopsis-panel-bg);
        }

    /* ------------------------------- Animations ------------------------------- */
</style>