<script>
    import { browser } from "$app/environment";
    import viewport from "@components/viewport_actions/useViewportActions";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The element id for the lazy wrapper
         * @type {string | null}
         */
        export let id=null;

        /**
         * The class name(s) to be added to the lazy wrapper
         * @type {string}
         */
        export let className = "";

        /**
         * The image url to be loaded
         * @type {string}
         */
        export let image_url;

        /**
         * @type {HTMLImageElement}
         */
        let image_element;

        /**
         * The image src to be loaded
         * @type {string}
         */
        let image_src = "";

        /**
         * @type  {boolean} whether or not the image src has been loaded
         */
        let is_loaded = false;

        /**
         * @type  {boolean} whether or not the image src has entered the viewport
         */
        let entered_viewport = false;

        $: if(image_src !== image_url && entered_viewport && browser) {
            is_loaded = false;
            mountSrc();
        }
    
    /*=====  End of Properties  ======*/
   
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * @type {import("@components/viewport_actions/useViewportActions").ViewportEventHandler}
         */
        const handleViewportEnter = event => {
            if (entered_viewport && is_loaded) return;

            entered_viewport = true;
            mountSrc();
        };

        const mountSrc = () => {
            image_element = document.createElement('img');

            image_element.onload = () => {
                image_src = image_element.src;
                is_loaded = true;
            };

            image_element.src = image_url;
        }
    
    /*=====  End of Methods  ======*/

</script>

<div use:viewport on:viewportEnter={handleViewportEnter}  class={`lazy-wrapper ${className !== undefined ? className : ''}`} {id}>
    {#if is_loaded}
        <slot name="lazy-wrapper-image" {image_src}></slot>
    {/if}
</div>

<style>
    
    @keyframes placeholderShimmer {
        0% {
            background-position: 220% 0
        }
        100% {
            background-position:  -210% 0
        }
    }
    
    .lazy-wrapper:empty {
        animation-duration: 1s;
        animation-fill-mode: forwards;
        animation-iteration-count: infinite;
        animation-name: placeholderShimmer;
        animation-timing-function: linear;
        width: 100%;
        height: 100%;
        background: rgb(236,236,236);
        background: linear-gradient(98deg, rgba(236,236,236,1) 34%, rgba(255,255,255,1) 38%, rgba(255,255,255,1) 50%, rgba(235,235,235,1) 71%); 
        background-size: 130%;
    }
</style>