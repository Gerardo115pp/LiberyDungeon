<script>
    import { onMount, onDestroy } from "svelte";

    
    
    /*=============================================
    =            Properties             =
    =============================================*/
    
        /**
         * The Id this component will use in the DOM.
         * @type {string}
         */
        export let component_id = `cover-slide-${crypto.randomUUID()}`;

        /**
         * Whether to use position absolute 0,0 to attach the cover. If false, it will fall under the 
         * user to position it.
         * @type {boolean}
         * @default true
         */
        export let use_absolute_position = true;

        /**
         * Whether to listen for mouse enter and mouse leave events on the parent component. or ignore them.
         * @type {boolean}
        */
        export let ignore_parent_events = false;

        /**
         * True if the parent event listeners have been attached
         * @type {boolean}
         */
        let parent_listeners_attached = false;

        /**
         * Parent component's DOM node
         * @type {HTMLElement|null}
         */
        let parent_html_node;

        /**
         * Whether to show the cover
         * @type {boolean}
         */
        let show_cover_animation = false;

        
        /*----------  Animation parameters  ----------*/
        
            /**
             * The duration of the uncover animation in milliseconds
             * @type {number}
             * @default 1200
             */
            export let uncover_duration = 1200;

            /**
             * The delay of the animation in milliseconds
             * @type {number}
             * @default 880
             */
            export let animation_delay = 200;



    /*=====  End of Properties  ======*/

    onMount(() => {
        parent_html_node = getParentNode();

        attachParentListeners();
    });

    onDestroy(() => {
        removeParentListeners();
    });


    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Attaches mouse enter and mouse leave event listeners to the parent component.
         */
        const attachParentListeners = () => {
            if (parent_html_node && !parent_listeners_attached) {
                parent_html_node.addEventListener("mouseenter", handleMouseEnter);
                parent_html_node.addEventListener("mouseleave", handleMouseLeave);
                parent_listeners_attached = true;
            }
        }
    
        /**
         * Returns the parent component's DOM node. Uses component_id and the :has selector to find the parent.
         * @returns {HTMLElement}
         */
        const getParentNode = () => {
            return document.querySelector(`:has(> #${component_id})`);
        }

        /**
         * Resets the cover animation state
         * @param {MouseEvent} event
         */
        const handleMouseLeave = (event) => {
            show_cover_animation = false;
        }

        /**
         * Starts the cover animation.
         * @param {MouseEvent} event
         */
        const handleMouseEnter = (event) => {
            show_cover_animation = !ignore_parent_events;
        }

        /**
         * Removes event listeners from the parent component.
         */
        const removeParentListeners = () => {
            if (parent_html_node && parent_listeners_attached) {
                parent_html_node.removeEventListener("mouseenter", handleMouseEnter);
                parent_html_node.removeEventListener("mouseleave", handleMouseLeave);
                parent_listeners_attached = false;
            }
        }
    
    /*=====  End of Methods  ======*/
    
    
    
    
    
</script>

<div class="hover-animations-cover-slide" 
    id="{component_id}" 
    class:the-cult-of-the-absolute={use_absolute_position}
>
    {#if show_cover_animation}
        <div class="hacs-cover-frame"
            style:animation-duration="{uncover_duration}ms"
            style:animation-delay="{animation_delay}ms"
        >
            <div class="hacs-cover-eraser"
                style:animation-duration="{uncover_duration * 0.5}ms"
                style:animation-delay="{animation_delay}ms"
            ></div>
        </div>
    {/if}
</div>
    
<style>
    @keyframes cover-slide {
        0% {
            translate: 100% 0 0;
        }
        100% {
            translate: 0 0 0;
        }
    }

    @keyframes elasticEraser {
        0% {
            opacity: 0;
            scale: 1 1
        }
        5% {
            opacity: 1;
        }
        25% {
            scale: 2.8 1;
        }
        50% {
            scale: 3.3 1;
        }
        80% {
            scale: 2.2 1;
        }
        100% {
            /* opacity: 0; */
            scale: 0.1 1;
        }
    }

    .hover-animations-cover-slide {
        --animation-timing-function: cubic-bezier(0.075, 0.82, 0.165, 1);
        position: static;
        width: 100%;
        height: 100%;
        z-index: var(--z-index-t-1);
        overflow: hidden;
    }

    .hover-animations-cover-slide.the-cult-of-the-absolute {
        position: absolute;
        top: 0;
        left: 0;
    }

    .hacs-cover-frame {
        width: 100%;
        height: 100%;
        translate: 100% 0 0;
        background: var(--main-dark);
        animation-name: cover-slide;
        animation-timing-function: var(--animation-timing-function);
        animation-fill-mode: forwards;
        opacity: 0.6;
    }

    .hacs-cover-eraser {
        width: 6%;
        height: 100%;
        background: var(--main-dark-color-8);
        animation-name: elasticEraser;
        animation-timing-function: linear;
        animation-fill-mode: forwards;
        transform-origin: left;
    }

    @supports (color: rgb(from white r g b)) {
        .hacs-cover-frame {
            background: rgb(from var(--main-dark) r g b / 0.6);
            opacity: 1;
        }

        .hacs-cover-eraser {
            background: hsl(from var(--main-dark-color-8) h s l / 0.8);
        }
    }
</style>