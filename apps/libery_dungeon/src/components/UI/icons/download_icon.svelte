<script>

    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The color of the icon.
         * @type {string}
         * @default var(--grey-1)
         */ 
        export let icon_color = "var(--grey-1)";

        /**
         * The size of the icon
         * @type {string}
         * @default "24px"
         */
        export let icon_size = "24px";
    
    /*=====  End of Properties  ======*/
    
</script>

<svg 
    class="dungeon-download-icon"
    viewBox="0 0 24 24"
    style:--icon-color="{icon_color}"
    style:width="{icon_size}"
>
        <path d="M5 12V18C5 18.5523 5.44772 19 6 19H18C18.5523 19 19 18.5523 19 18V12"></path> 
        <path d="M12 3L12 15M12 15L16 11M12 15L8 11"></path>
</svg>

<style>
    .dungeon-download-icon {
        aspect-ratio: 1 / 1;

        & path {
            stroke: var(--icon-color);
            stroke-linecap: round;
            stroke-linejoin: round;
            fill: none;
        }
    }
</style>