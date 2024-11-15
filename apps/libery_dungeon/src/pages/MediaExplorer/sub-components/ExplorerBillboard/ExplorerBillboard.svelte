<script>
    import LiberyHeadline from '@components/UI/LiberyHeadline.svelte';
    import { onMount } from 'svelte';

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
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        handleChangeBillboardImage();
    });

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
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
    
    /*=====  End of Methods  ======*/
    
    
    
</script>

<section id="media-explorer-billboard">
    <div id="mexbill-underlay-billboard-wrapper"
        class="full-vw"
    >
        {#if current_billboard_media != null}
            {#if current_billboard_media.isImage()}
                <img
                    decoding="async"
                    src="{current_billboard_media.Url}" 
                    alt=""
                >
            {:else if current_billboard_media.isVideo()}
                <video 
                    src="{current_billboard_media.Url}"
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
                    forced_font_size="var(--font-size-h3)"
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
        
        /* position: relative; */
        display: flex;
        width: 100%;
        height: 50dvh;
        container-type: size;
        flex-direction: column;
        justify-content: flex-end;
        z-index: var(--z-index-b-1);
    }
    
    /*=============================================
    =            Billboard            =
    =============================================*/
    
        #mexbill-underlay-billboard-wrapper {
            position: absolute;
            inset: 0;
            /* width: 100dvw; */
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

        /* -------------------------------- Synopsis -------------------------------- */

        #mexbill-synopsis-panel {
            --synopsis-panel-bg: hsl(from var(--grey-1) h s l / 0.08);

            width: max-content;
            background: var(--synopsis-panel-bg);
            padding-block-start: var(--spacing-2);
            padding-block-end: var(--spacing-4);
            padding-inline: var(--spacing-5);
            translate: calc(-1 * var(--common-page-inline-padding));
            box-shadow: 0 -10px 36px 40px var(--synopsis-panel-bg);
        }
    
    /*=====  End of Billboard  ======*/
    
</style>