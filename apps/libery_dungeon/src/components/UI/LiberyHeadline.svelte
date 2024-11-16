<script>
    import TaggedText from "@components/Wrappers/TaggedText.svelte";
    import viewport from "@components/viewport_actions/useViewportActions";
    import { elasticOut } from "svelte/easing";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * @type {string} - the tag name for the headline
         */
        export let headline_tag = "h1";

        /**
         * @type {string} - the text to be displayed in the headline
         */
        export let headline_text;

        /**
         * @type {string} - extra props to be added to the headline tag as html
         */
        export let extra_props="";

        /**
         * @type {string} - the color of the headline
         * @default "var(--text-color-1)"
         */
        export let headline_color = "var(--text-color-1)";

        /**
         * A color for the wrapping tags.
         * @type {string}
         * @default "var(--text-color-2"
         */
        export let headline_tag_color = "var(--text-color-2)";

        export let spacing = "var(--spacing-1)";
        
        /**
         * @type {string} 
         */
        export let headline_font_size;

        /**
         * By default if the tag name is not h1, the bottom lines are not shown. if this is set to true, the bottom lines will be shown no matter the tag name
         * @type {boolean}
         */
        export let force_bottom_lines = false;

        export let animated = true;
    
        /*----------  Style  ----------*/

                export let text_transform = "uppercase";

        /*----------  Animation  ----------*/

            let visible = false;
            export let animation_duration = 400;
            export let animation_delay = 0;
    
    /*=====  End of Properties  ======*/


    /**
     * @param {Element} node
     * @param {fallingTransitionConfig} param1
     * @typedef {Object} fallingTransitionConfig
     * @property {number} [delay]
     * @property {number} [duration=800]
     * @property {number} [rotation_start_at=.8]
     * @property {number} [rotation=90]
     * @property {number} [fall_height=2000]
     * @property {boolean} [invert=false]
     * @property {(t: number) => number} easing
     * @returns {import('svelte/animate').AnimationConfig}
     */
    const fallingTransition = (node, {delay, duration=800, rotation_start_at=.8, rotation=90, fall_height=2000, invert=false, easing=t => t}) => {
        
        // const rotation_start_at = .7;
        const rotation_offset = 1 - rotation_start_at; // defined for readability

        return {
            duration,
            delay,
            css: t => {
                t = easing(t);
                let elapsed_time = 1 - t;
                let transform = `translateY(-${Math.max(0, fall_height - (fall_height*(t/rotation_start_at)))}%)`

                transform += t < rotation_start_at ? `rotate(${invert ? '' : '-'}${rotation}deg)` : ` rotate(${invert ? '' : '-'}${Math.trunc(rotation * (elapsed_time/rotation_offset))}deg)`


                return `
                    transform: ${transform};
                    opacity: ${t};
                `
            }
        }
    }

</script>

<TaggedText 
    spacing={spacing} tag_name={headline_tag} {extra_props}
    tag_color={headline_tag_color}
>
    <div class="headline-wrapper" on:viewportEnter={() => visible = true} use:viewport>
        <h1 class="libery-headline" 
            style:color="{headline_color}" 
            style:font-size={headline_font_size} 
            style:text-transform={text_transform}
        >
            {headline_text}
        </h1>
        {#if (headline_tag === "h1" && animated && visible) || force_bottom_lines}
            <div class="bottom-lines">
                <svg viewBox="0 0 275 23" fill="none" preserveAspectRatio="xMidYMax">
                    <path in:fallingTransition={{delay: animation_delay*1.2 ,duration:animation_duration, rotation: 35,rotation_start_at: .8, easing: elasticOut}} class="line-short" d="M1 16.5088L124.385 2"/>
                    <path in:fallingTransition={{delay: animation_delay ,duration:animation_duration, rotation: 35,rotation_start_at: .7, fall_height: 2000, easing: elasticOut, invert: true}} class="line-long" d="M95.8305 16.5088L274.28 20.5088"/>
                </svg>                
            </div>
        {:else if headline_tag === "h1" && !animated}
            <div class="bottom-lines">
                <svg viewBox="0 0 275 23" fill="none" preserveAspectRatio="xMidYMax">
                    <path class="line-short" d="M1 16.5088L124.385 2"/>
                    <path class="line-long" d="M95.8305 16.5088L274.28 20.5088"/>
                </svg>                
            </div>
        {/if}
    </div>
</TaggedText>

<style>
    .headline-wrapper {
        display: flex;
        width: 100%;
        flex-direction: column;
        row-gap: var(--vspacing-2);
    }

    h1.libery-headline {
        position: relative;
        text-align: left;
        font-size: 112px;
        font-weight: 400;
        letter-spacing: -3.36px;
        white-space: nowrap;
        line-height: .75;
        margin: 0;
        padding: 0;
        z-index: var(--z-index-1);
    }

    .bottom-lines {
        max-width: var(--vspacing-6);
        padding-left: var(--vspacing-4);
    }

    .bottom-lines svg {
        overflow: visible;
    }

    .bottom-lines svg path{
        position: relative;
        stroke: var(--main-dark-color-7);
        stroke-width: .7ex;
        transform-box: fill-box;
        z-index: var(--z-index-2);
    }

    .bottom-lines svg path.line-short {
        transform-origin: left center;
    }

    .bottom-lines svg path.line-long {
        transform-origin: right center;
    }

</style>