<script>
    import star from "@svg/libery_star_full.svg"
    import blue_star from "@svg/libery_star_blue.svg"
    import libery_line from "@svg/libery_line.svg"
    import { layout_properties } from "@stores/layout";
    import { writable } from "svelte/store";
    import { onMount } from "svelte";
    

    const star_x_offset = layout_properties.IS_MOBILE ? 44 : 74; // px
    const star_y_offset = layout_properties.IS_MOBILE ? 0 : -17; // px
    const stars_separation = layout_properties.IS_MOBILE ? 12 : 37; // px

    onMount(() => {
        height_store.set(star_height);
    })


    export let height_store = writable(0);
    let star_height = 0;
</script>

<div id="libery-star-wrapper">
    <div bind:clientHeight={star_height} class="svg-container" id="lsw-star">
        <svg style:overflow="visible" id="libery-star-hero" viewBox="0 0 850 784.3438">
            <svg x="0" y="0" id="libery-star-blue" viewBox="0 0 678.23 625.76">
                <defs>
                  <style>
                    .libery-star-blue-cls-1 {
                      fill: #032140;
                    }
                  </style>
                </defs>
                <g id="libery-star-blue" class="anim-shake-v" style:animation-duration="2s" style:animation-delay=".02s">
                  <polygon class="libery-star-blue-cls-1" points="344.95 309.97 428.76 179.76 439.98 412.48 344.95 309.97"/>
                  <polygon class="libery-star-blue-cls-1" points="245.83 422.68 337.17 309.97 344.95 309.97 439.98 412.48 440.08 414.6 350.78 399.36 340.09 398.75 245.83 422.68"/>
                  <polygon class="libery-star-blue-cls-1" points="251.66 175.1 337.17 309.97 245.83 422.68 249.03 341.61 251.66 175.1"/>
                  <polygon class="libery-star-blue-cls-1" points="251.66 175.1 337.5 228.67 344.95 228.34 428.76 179.76 344.95 309.97 337.17 309.97 251.66 175.1"/>
                  <polygon class="libery-star-blue-cls-1" points="147.05 213.33 174.9 254.58 335.23 92.31 147.05 213.33"/>
                  <polygon class="libery-star-blue-cls-1" points="349.16 92.31 508.57 243.09 539.1 203.69 349.16 92.31"/>
                  <polygon class="libery-star-blue-cls-1" points="251.66 175.1 250.1 275.96 157.05 273.17 251.66 175.1"/>
                  <polygon class="libery-star-blue-cls-1" points="339.43 92.31 335.23 92.31 251.66 175.1 314.72 215.69 339.43 92.31"/>
                  <polygon class="libery-star-blue-cls-1" points="339.43 92.31 342.23 92.31 371.18 213.77 344.95 228.34 337.5 228.67 314.72 215.69 339.43 92.31"/>
                  <polygon class="libery-star-blue-cls-1" points="437.26 174.9 371.18 213.77 342.23 92.31 349.16 92.31 437.26 174.9"/>
                  <polygon class="libery-star-blue-cls-1" points="250.1 275.96 157.05 273.17 154.1 275.96 249.03 341.61 250.1 275.96"/>
                  <polygon class="libery-star-blue-cls-1" points="245.83 422.68 154.1 275.96 249.03 341.61 245.83 422.68"/>
                  <polygon class="libery-star-blue-cls-1" points="154.1 275.96 93.27 337.77 232.23 434.34 313.85 484.87 245.83 422.68 154.1 275.96"/>
                  <polygon class="libery-star-blue-cls-1" points="245.83 422.68 340.09 398.75 342.57 497.72 343.97 500.88 307.22 478.81 245.83 422.68"/>
                  <polygon class="libery-star-blue-cls-1" points="340.09 398.75 342.57 497.72 343.97 500.88 346.89 497.72 350.78 399.36 340.09 398.75"/>
                  <polygon class="libery-star-blue-cls-1" points="350.78 399.36 440.08 414.6 441.14 419.76 435.31 431.43 411.99 454.75 343.97 500.88 346.89 497.72 350.78 399.36"/>
                  <polygon class="libery-star-blue-cls-1" points="432.93 266.25 535.5 266.25 437.26 174.9 428.76 179.76 432.93 266.25"/>
                  <polygon class="libery-star-blue-cls-1" points="432.93 266.25 535.5 266.25 436.27 335.57 432.93 266.25"/>
                  <polygon class="libery-star-blue-cls-1" points="436.27 335.57 440.08 414.6 535.5 266.25 436.27 335.57"/>
                  <polygon class="libery-star-blue-cls-1" points="535.5 266.25 592.72 324.54 445.03 435.31 411.99 454.75 435.31 431.43 441.14 419.76 440.08 414.6 535.5 266.25"/>
                  <polygon class="libery-star-blue-cls-1" points="343.97 500.88 411.99 454.75 445.03 435.31 459.06 424.79 531.51 580.09 531.51 587.47 343.97 500.88"/>
                  <polygon class="libery-star-blue-cls-1" points="544.17 360.96 531.54 579.58 459.08 424.28 544.17 360.96"/>
                  <polygon class="libery-star-blue-cls-1" points="343.97 500.88 226.19 430.14 155.47 598.55 343.97 500.88"/>
                  <polygon class="libery-star-blue-cls-1" points="226.19 430.14 148.53 376.17 155.47 598.55 226.19 430.14"/>
                  <polygon class="libery-star-blue-cls-1" points="260.43 544.17 343.97 625.76 424.35 537.99 343.97 500.88 260.43 544.17"/>
                  <polygon class="libery-star-blue-cls-1" points="118.54 369.24 69.96 435.31 151.54 472.84 150.47 438.32 118.54 369.24"/>
                  <polygon class="libery-star-blue-cls-1" points="538.55 458.17 623.1 421.04 570.29 361.46 565.52 367.29 540.47 425.59 538.55 458.17"/>
                  <polygon class="libery-star-blue-cls-1" points="147.05 213.33 154.1 4.75 218.6 167.4 147.05 213.33"/>
                  <polygon class="libery-star-blue-cls-1" points="218.6 167.4 279.8 128.1 154.1 4.75 218.6 167.4"/>
                  <polygon class="libery-star-blue-cls-1" points="154.1 4.75 335.23 92.31 279.8 128.1 154.1 4.75"/>
                  <polygon class="libery-star-blue-cls-1" points="349.16 92.31 520.82 0 406.98 126.35 349.16 92.31"/>
                  <polygon class="libery-star-blue-cls-1" points="406.98 126.35 520.82 0 465.3 160.49 406.98 126.35"/>
                  <polygon class="libery-star-blue-cls-1" points="465.3 160.49 539.1 203.69 520.82 0 465.3 160.49"/>
                  <polygon class="libery-star-blue-cls-1" points="294.6 72.67 335.23 22.02 384.77 73.16 349.16 92.31 335.23 92.31 294.6 72.67"/>
                  <polygon class="libery-star-blue-cls-1" points="335.23 92.31 335.23 22.02 349.16 92.31 335.23 92.31"/>
                  <polygon class="libery-star-blue-cls-1" points="148.07 183.23 147.05 213.33 174.9 254.58 150.19 279.93 0 235.15 148.07 183.23"/>
                  <polygon class="libery-star-blue-cls-1" points="150.19 279.93 93.27 337.77 0 235.15 150.19 279.93"/>
                  <polygon class="libery-star-blue-cls-1" points="508.57 243.09 539.1 203.69 536.52 174.9 678.23 227.37 535.5 266.25 508.57 243.09"/>
                  <polygon class="libery-star-blue-cls-1" points="535.5 266.25 592.72 324.54 678.23 227.37 535.5 266.25"/>
                  <polygon class="libery-star-blue-cls-1" points="592.72 324.54 544.17 360.96 540.47 425.59 565.52 367.29 570.29 361.46 592.72 324.54"/>
                  <polygon class="libery-star-blue-cls-1" points="150.47 438.32 118.54 369.24 93.27 337.77 148.53 376.17 150.47 438.32"/>
                </g>
            </svg>
            <path class="blue-line anim-shake-h" style:animation-delay=".5s" style:animation-duration="1s" d="M50 200L950.0 100"/>
            <path class="blue-line anim-shake-h" style:animation-duration="0.8s" d="M0 392.1719L880.0 392.1719"/>
            <path class="blue-line anim-shake-h" style:animation-delay=".5s" style:animation-duration="1.1s" d="M50 584.3438L950.0 684.3438"/>
            <svg x="38" y="20" id="libery-star-full"  viewBox="0 0 681.49 628.85">
                <defs>
                    <style>
                        .libery-star-full-cls-1 {
                            fill: #e45d00;
                        }
    
                        .libery-star-full-cls-2 {
                            fill: #fe7d00;
                        }
    
                        .libery-star-full-cls-3 {
                            fill: #c52d00;
                        }
    
                        .libery-star-full-cls-4 {
                            fill: #f27400;
                        }
    
                        .libery-star-full-cls-5 {
                            fill: #db4c00;
                        }
    
                        .libery-star-full-cls-6 {
                            fill: #9a0000;
                        }
    
                        .libery-star-full-cls-7 {
                            fill: #feb800;
                        }
    
                        .libery-star-full-cls-8 {
                            fill: #fe9c00;
                        }
    
                        .libery-star-full-cls-9 {
                            fill: #fdcf00;
                        }
    
                        .libery-star-full-cls-10 {
                            fill: #f76400;
                        }
    
                        .libery-star-full-cls-11 {
                            fill: #ad0000;
                        }
    
                        .libery-star-full-cls-12 {
                            fill: #fe5900;
                        }
    
                        .libery-star-full-cls-13 {
                            fill: #febc00;
                        }
    
                        .libery-star-full-cls-14 {
                            fill: #fead00;
                        }
    
                        .libery-star-full-cls-15 {
                            fill: #fa5100;
                        }
    
                        .libery-star-full-cls-16 {
                            fill: #b10000;
                        }
    
                        .libery-star-full-cls-17 {
                            fill: #eb3500;
                        }
    
                        .libery-star-full-cls-18 {
                            fill: #f71a00;
                        }
    
                        .libery-star-full-cls-19 {
                            fill: #f81c00;
                        }
    
                        .libery-star-full-cls-20 {
                            fill: #fd8100;
                        }
    
                        .libery-star-full-cls-21 {
                            fill: #f05b00;
                        }
    
                        .libery-star-full-cls-22 {
                            fill: #d01a00;
                        }
    
                        .libery-star-full-cls-23 {
                            fill: #e50d00;
                        }
                    </style>
                </defs>
                <g id="libery-star-full" class="anim-shake-v" style:animation-duration="2s" >
                    <polygon class="libery-star-full-element libery-star-full-cls-4" points="346.57 312.33 430.38 182.13 441.6 414.85 346.57 312.33"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-3" points="247.46 425.05 338.79 312.33 346.57 312.33 441.6 414.85 441.7 416.96 352.4 401.73 341.71 401.12 247.46 425.05"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-1" points="253.29 177.46 338.79 312.33 247.46 425.05 250.65 343.97 253.29 177.46"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-9" points="253.29 177.46 339.12 231.03 346.57 230.71 430.38 182.13 346.57 312.33 338.79 312.33 253.29 177.46"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-23" points="148.67 215.7 176.52 256.95 336.85 94.68 148.67 215.7"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-23" points="350.78 94.68 510.19 245.45 540.73 206.05 350.78 94.68"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-19" points="253.29 177.46 251.72 278.33 158.67 275.54 253.29 177.46"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-19" points="341.05 94.68 336.85 94.68 253.29 177.46 316.34 218.06 341.05 94.68"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-12" points="341.05 94.68 343.85 94.68 372.8 216.14 346.57 230.71 339.12 231.03 316.34 218.06 341.05 94.68"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-18" points="438.88 177.27 372.8 216.14 343.85 94.68 350.78 94.68 438.88 177.27"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-14" points="251.72 278.33 158.67 275.54 155.72 278.33 250.65 343.97 251.72 278.33"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-10" points="247.46 425.05 155.72 278.33 250.65 343.97 247.46 425.05"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-17" points="155.72 278.33 94.89 340.14 233.85 436.71 315.47 487.23 247.46 425.05 155.72 278.33"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-16" points="247.46 425.05 341.71 401.12 344.19 500.08 345.59 503.24 308.84 481.17 247.46 425.05"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-5" points="341.71 401.12 344.19 500.08 345.59 503.24 348.51 500.08 352.4 401.73 341.71 401.12"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-16" points="352.4 401.73 441.7 416.96 442.76 422.13 436.93 433.79 413.61 457.11 345.59 503.24 348.51 500.08 352.4 401.73"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-18" points="434.55 268.62 537.12 268.62 438.88 177.27 430.38 182.13 434.55 268.62"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-14" points="434.55 268.62 537.12 268.62 437.89 337.94 434.55 268.62"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-10" points="437.89 337.94 441.7 416.96 537.12 268.62 437.89 337.94"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-17" points="537.12 268.62 594.34 326.91 446.65 437.68 413.61 457.11 436.93 433.79 442.76 422.13 441.7 416.96 537.12 268.62"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-22" points="345.59 503.24 413.61 457.11 446.65 437.68 460.68 427.16 533.13 582.46 533.13 589.84 345.59 503.24"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-21" points="545.79 363.32 533.16 581.95 460.71 426.65 545.79 363.32"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-22" points="345.59 503.24 227.81 432.51 157.09 600.92 345.59 503.24"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-21" points="227.81 432.51 150.15 378.54 157.09 600.92 227.81 432.51"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-6" points="262.05 546.53 345.59 628.13 425.97 540.36 345.59 503.24 262.05 546.53"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-6" points="120.17 371.6 71.58 437.68 153.16 475.2 152.09 440.69 120.17 371.6"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-6" points="540.17 460.54 624.72 423.4 571.91 363.83 567.14 369.66 542.09 427.96 540.17 460.54"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-2" points="148.67 215.7 155.72 7.12 220.22 169.76 148.67 215.7"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-7" points="220.22 169.76 281.42 130.47 155.72 7.12 220.22 169.76"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-8" points="155.72 7.12 336.85 94.68 281.42 130.47 155.72 7.12"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-8" points="350.78 94.68 522.44 2.37 408.6 128.72 350.78 94.68"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-7" points="408.6 128.72 522.44 2.37 466.92 162.86 408.6 128.72"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-2" points="466.92 162.86 540.73 206.05 522.44 2.37 466.92 162.86"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-13" points="296.22 75.04 336.85 24.39 386.39 75.53 350.78 94.68 336.85 94.68 296.22 75.04"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-8" points="336.85 94.68 336.85 24.39 350.78 94.68 336.85 94.68"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-20" points="149.69 185.59 148.67 215.7 176.52 256.95 151.82 282.3 1.62 237.51 149.69 185.59"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-15" points="151.82 282.3 94.89 340.14 1.62 237.51 151.82 282.3"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-20" points="510.19 245.45 540.73 206.05 538.14 177.27 679.85 229.74 537.12 268.62 510.19 245.45"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-15" points="537.12 268.62 594.34 326.91 679.85 229.74 537.12 268.62"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-11" points="594.34 326.91 545.79 363.32 542.09 427.96 567.14 369.66 571.91 363.83 594.34 326.91"/>
                    <polygon class="libery-star-full-element libery-star-full-cls-11" points="152.09 440.69 120.17 371.6 94.89 340.14 150.15 378.54 152.09 440.69"/>
                </g>
            </svg>
        </svg>
    </div>
</div>

<style>
    #libery-star-wrapper {
        --line-offset: 13%;
        --limit-lines-rotation: 5deg;
        --star-width: calc(var(--vspacing-8) * 0.7193212890625);

        position: relative;
        width: 100%;
        height: 100%;
        z-index: var(--z-index-5);
    }

    :global(#libery-star-wrapper > svg) {
        width: 100%;
    }

    #libery-star-hero > svg {
        transform-box: fill-box;
        transform-origin: center;
        backface-visibility: hidden;
        overflow: visible;
    }


    /*----------  star  ----------*/
    
    .libery-star-full-element {
        stroke: #c47562;
        stroke-miterlimit: 10;
    }

    #lsw-star {
        position: absolute;
        width: var(--star-width);
        top: 0;
        left: 0;
        z-index: 3;
    }

    :global(#lsw-star polygon.libery-star-full-element) {
        stroke: var(--grey) !important;
    }


    /*----------  blue star  ----------*/

        #libery-star-blue {
            stroke: var(--accent-8);
        }
    
    /*----------  Lines  ----------*/
    
    
    .blue-line {
        stroke: var(--accent);
        stroke-width: 4px;
        /* stroke-dasharray: 10px; */
        stroke-linecap: round;
        stroke-miterlimit: 10;
    }

    
    /*----------  Mobile  ----------*/
    
    @media only screen and (max-width: 765px) {
        #libery-star-wrapper {
            --star-width: calc(var(--vspacing-8) * 0.31932);
            --line-offset: 13%;
            
            width: 80%;
            height: max-content;
        }

        /* #lsw-star {
        } */
    }
</style>