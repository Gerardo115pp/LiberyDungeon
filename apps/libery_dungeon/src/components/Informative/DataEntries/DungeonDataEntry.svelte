<script>

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The label of the information entry
         * @type {string}
         * @default 'data'
         */
        export let information_entry_label = 'data';

        /**
         * The value of the information entry
         * @type {string}
         * @default ''
         */
        export let information_entry_value = '';

        /*----------  Style  ----------*/
        
            /**
             * the component's default font-size.
             * @type {string}
             * @default "var(--font-size-1)"
             */ 
            export let font_size = "var(--font-size-1)";
        
        /*----------  Behavior  ----------*/
        
            /**
             * Whether the information value should be pasted on click.
             * @type {boolean} 
             * @default false
             */
            export let paste_on_click = false;
    
    /*=====  End of Properties  ======*/


    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles a click on one of the media information items. If there is clipboard-write permission, the value of the item is copied to the clipboard.
         * @param {MouseEvent} event The click event
         */
         const handleMediaInformationItemClick = event => {
            if (!paste_on_click) return;

            try {
                navigator.clipboard.writeText(information_entry_value);
            } catch (error) {
                console.error('Error writing to clipboard:', error);
                return;
            }   
        }
    
    /*=====  End of Methods  ======*/
    
    
    
</script>

<li class="mv-mip-media-information-item" 
    class:click-to-paste={paste_on_click} 
    style:font-size="{font_size}"
    on:click={handleMediaInformationItemClick}
>
    <p class="mv-mip-media-information-content">
        <strong class="mv-mip-media-information-item-name">
            {information_entry_label}
        </strong>
        <span class="mv-mip-media-information-item-value">
            {information_entry_value}
        </span>
    </p>
    <div class="clipboard-paste-feedback-graphic">
        <svg viewBox="0 0 32 32">
            <path d="M27.2,8.22H23.78V5.42A3.42,3.42,0,0,0,20.36,2H5.42A3.42,3.42,0,0,0,2,5.42V20.36a3.43,3.43,0,0,0,3.42,3.42h2.8V27.2A2.81,2.81,0,0,0,11,30H27.2A2.81,2.81,0,0,0,30,27.2V11A2.81,2.81,0,0,0,27.2,8.22ZM5.42,21.91a1.55,1.55,0,0,1-1.55-1.55V5.42A1.54,1.54,0,0,1,5.42,3.87H20.36a1.55,1.55,0,0,1,1.55,1.55v2.8H11A2.81,2.81,0,0,0,8.22,11V21.91ZM28.13,27.2a.93.93,0,0,1-.93.93H11a.93.93,0,0,1-.93-.93V11a.93.93,0,0,1,.93-.93H27.2a.93.93,0,0,1,.93.93Z"></path>
        </svg>
    </div>
</li>

<style>
    li.mv-mip-media-information-item {
        border-radius: var(--border-radius);
        overflow: hidden;

        & p {
            font-size: inherit;
            padding: var(--vspacing-1);
            line-height: 1;
        }
    }

    strong.mv-mip-media-information-item-name {
        font-weight: bold;
        
        &::after {
            content: ':';
        }
    }

    span.mv-mip-media-information-item-value {
        text-wrap: wrap;
        word-break: break-all;
        word-wrap: break-word;
    }

    
    /*=============================================
    =            Clipboard feedback            =
    =============================================*/
    
        li.mv-mip-media-information-item {
            position: relative;
        }    

        .clipboard-paste-feedback-graphic {
            position: absolute;
            display: flex;
            width: 100%;
            height: 100%;
            inset: 0;
            background: hsl(from var(--grey) h s l / 0.7);
            container-type: size;
            justify-content: center;
            align-items: center;
            opacity: 0;
            transition: all 0.3s ease-in;

            & svg {
                width: 70cqh;
                height: 70cqh;
                fill: var(--main);
            }            
        }

        li.mv-mip-media-information-item.click-to-paste .clipboard-paste-feedback-graphic:hover {
            opacity: 1;
        }
    /*=====  End of Clipboard feedback  ======*/
    
    

    @supports (color: rgb( from white r g b / 1)) {
        li.mv-mip-media-information-item p {
            padding: var(--vspacing-1) var(--vspacing-2);
            background: hsl(from var(--grey) h s l / 0.8);
        }
    }
</style>

