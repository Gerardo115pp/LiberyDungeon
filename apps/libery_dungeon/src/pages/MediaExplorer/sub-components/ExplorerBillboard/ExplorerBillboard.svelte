<script>
    import { browser } from '$app/environment';
    import LiberyHeadline from '@components/UI/LiberyHeadline.svelte';
    import { CanvasImage } from '@libs/LiberyColors/image_models';
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

        
        /*----------  Style  ----------*/
        
            /**
             * The theme color to use fort the overlay content(the text). It's value is determined by the color pallette of the current_billboard_media.
             * @type {boolean}
             * @default false - White overlay color theme is best for MOST cases except when the media is extremely white.
             */ 
            let use_dark_theme = false;
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        if (the_billboard_category != null && the_billboard_category.content.length > 0) {
            handleChangeBillboardImage();
        }

        checkWindowScroll();
    });

    onDestroy(() => {
        if (!browser) return;

        hanldeBillboardDestroy();
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
         * Changes the current billboard media.
         */ 
        const handleChangeBillboardImage = () => {
            if (current_billboard_media != null && the_billboard_category.content.length === 1) return; // no other options but the one that already exists.

            let inifite_loop_guard = 0;

            /**
             * @type {import('@models/Medias').Media | null}
             */
            let new_billboard_media = null;

            while (new_billboard_media == null) {
                const rand_media_index = Math.trunc(Math.random() * the_billboard_category.content.length);

                new_billboard_media = the_billboard_category.content[rand_media_index];

                if (current_billboard_media != null && current_billboard_media.uuid === new_billboard_media.uuid) {
                    new_billboard_media = null
                }

                inifite_loop_guard++;

                if (inifite_loop_guard > the_billboard_category.content.length) {
                    console.error(`In ExplorerBillboard.handleChangeBillboardImage: Infinite loop detected after ${the_billboard_category.content.length} iterations. Exiting.`);
                    return;
                }
            }

            current_billboard_media = new_billboard_media;

            console.log("Changed billboard media to: ", current_billboard_media);
        
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
        const hanldeBillboardDestroy = () => {
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