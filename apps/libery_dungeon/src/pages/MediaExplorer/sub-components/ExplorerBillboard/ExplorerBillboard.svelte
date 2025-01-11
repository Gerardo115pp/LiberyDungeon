<script>
    import { browser } from '$app/environment';
    import LiberyHeadline from '@components/UI/LiberyHeadline.svelte';
    import { CanvasImage } from '@libs/LiberyColors/image_models';
    import { linearCycleNavigationWrap } from '@libs/LiberyHotkeys/hotkeys_movements/hotkey_movements_utils';
    import { getTaggedMedias } from '@models/Medias';
    import { current_cluster } from '@stores/clusters';
    import { navbar_ethereal } from '@stores/layout';
    import { onMount, onDestroy } from 'svelte';

    /*=============================================
    =            Properties            =
    =============================================*/
    
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
         * Returns the html billboard element.
         * @returns {HTMLImageElement | HTMLVideoElement | null}
         */
        const getBillboardElement = () => {
            const billboard_element = document.getElementById("mexbill-billboard");

            if (billboard_element == null) return billboard_element;

            if (!(billboard_element instanceof HTMLImageElement || billboard_element instanceof HTMLVideoElement)) {
                throw new Error("In ExplorerBillboard.getBillboardElement: The element with id 'mexbill-billboard' is not an instance of HTMLImageElement or HTMLVideoElement.");
            }

            return billboard_element;
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
                handleChangeBillboardImage(true); 
            }

            if (was_current_media_null && current_billboard_media != null) {
                onvalid_media_change();
            };
        }
    
        /**
         * Changes the current billboard media.
         * @param {boolean} [avoid_same_image] if true and uses random navigation, it will avoid landing in the same index.
         */ 
        const handleChangeBillboardImage = (avoid_same_image) => {
            if (avoid_same_image && current_billboard_media != null && billboard_medias.length === 1) {
                resetBillboardMediaIteration();
                return;
            }

            if (random_billboard_media_iteration) {
                iterBillboardMediasRandom();
            } else {
                iterBillboardMediasSequential();
            }
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
        }

        /**
         * Handles the load event of the billboard Image.
         * @param {Event} event
         */
        const handleBillboardImageLoad = (event) => {
            const billboard_element = event.target;

            if (!(billboard_element instanceof HTMLImageElement)) {
                console.error("In ExplorerBillboard.handleBillboardImageLoad: The event target is not an instance of HTMLImageElement.");
                return;
            }

            let new_white_percentage = calculateMediaWhitePercentage(billboard_element);

            use_dark_theme = new_white_percentage > 18;

            console.log("White percentage of the billboard image: ", new_white_percentage);
            console.log("White overlay color theme: ", !use_dark_theme);
        }

        /**
         * Handles the load event of the billboard Video.
         * @param {Event} event
         */
        const handleBillboardVideoLoad = (event) => {
            const billboard_element = event.target;

            if (!(billboard_element instanceof HTMLVideoElement)) {
                console.error("In ExplorerBillboard.handleBillboardVideoLoad: The event target is not an instance of HTMLVideoElement.");
                return;
            }

            use_dark_theme = false;
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
         * Iters through medias using a random progression system.
         * @param {boolean} [avoid_same_image] if true, it will avoid landing in the same index.
         * @returns {void}
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
                    return;
                }
            }

            current_billboard_media = new_billboard_media;

            console.log("Changed billboard media to: ", current_billboard_media);
        }

        /**
         * iterates over the billboard medias using a sequential progression system. When it reaches the limit of medias, it cycles.
         */
        const iterBillboardMediasSequential = () => {
            if (billboard_medias.length === 0) return;

            const next_index = linearCycleNavigationWrap(billboard_iterator_index, billboard_medias.length - 1, 1).value;

            billboard_iterator_index = next_index;
            current_billboard_media = billboard_medias[billboard_iterator_index];
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

            if (category_configuration.BillboardDungeonTags.length !== 0 && false) {
                // @ts-ignore
                new_billboard_medias = await loadBillboardMediasFromTags(category_configuration.BillboardDungeonTags);
                configuration_had_medias = true;
            } else if (category_configuration.BillboardMediaUUIDs.length > 0) {
                new_billboard_medias = await loadBillboardMediasFromUuids(category_configuration.BillboardMediaUUIDs);
                configuration_had_medias = true;
            } else {
                new_billboard_medias = the_billboard_category.content;
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
         * Resets the billboard iteration.
         * @param {boolean} [use_random_navigation]
         * @returns {void}
         */
        const resetBillboardMediaIteration = (use_random_navigation) => {
            billboard_iterator_index = 0;
            random_billboard_media_iteration = use_random_navigation === true;

            if (billboard_medias.length === 0) return;

            current_billboard_media = billboard_medias[billboard_iterator_index];
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
    
    /*=====  End of Methods  ======*/
    
</script>

<svelte:window 
    on:scroll|passive={handleWindowScroll}
/>
<section id="media-explorer-billboard"
    class:loaded-billboard={current_billboard_media != null}
    class:dark-overlay-color-theme={use_dark_theme}
>
    <div id="mexbill-underlay-billboard-wrapper">
        {#if current_billboard_media != null}
            {#if current_billboard_media.isImage()}
                <img
                    id="mexbill-billboard"
                    decoding="async"
                    src="{current_billboard_media.Url}" 
                    on:load={handleBillboardImageLoad}
                    on:error={handleBillboardImageError}
                    alt=""
                >
            {:else if current_billboard_media.isVideo()}
                <video 
                    id="mexbill-billboard"
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

            & > img, video {
                width: 100%;
                height: 100%;
                object-fit: cover;
                object-position: center 25%;
                z-index: var(--z-index-b-6);
            }

            &::after {
                content: "";
                position: absolute;
                inset: auto auto -2% 0;
                width: 100%;
                /* background: var(--body-bg-color); */
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
</style>