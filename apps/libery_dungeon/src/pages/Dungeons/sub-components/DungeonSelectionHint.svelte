<script>
    import { fade, fly } from "svelte/transition";
    import { bounceOut, cubicIn, cubicOut, quartOut, sineIn, sineInOut } from "svelte/easing";
</script>

<div class="dungeon-selection-hint-wrapper" transition:fade={{easing: cubicIn, duration: 200}}>
    <svg viewBox="0 0 70.664 13.185">
        <path class="dshw-line-left dshw-line" d="m.51691 7.8034h28.27" />
        <path class="dshw-line-right dshw-line" d="m42.908 7.7149h28.27" />
        <path d="m28.679 7.7924s.37623.0566.50599.24532c.12976.18873 3.1464 3.7123 6.5632 3.6416 3.4168-.0706 6.6147-3.689 6.6147-3.689s.18213-.27264.56234-.27264" />
        <ellipse cx="35.675" cy="6.6524" rx="3" ry="3" in:fly={{y: -10, opacity: 0.1, delay: 100, duration: 300, easing: bounceOut}} />
      </svg>      
</div>

<style>
    .dungeon-selection-hint-wrapper {
        width: 100%;
    }

    .dungeon-selection-hint-wrapper svg {
        stroke: var(--main-dark-color-7);
        width: 100%;
        overflow: visible;
        
        & path {
            fill: none;
            stroke-width: .7;
            stroke-linecap: round;
            stroke-linejoin: round;
        }

        & ellipse {
            fill: hsl(from var(--main-dark-color-4) h calc(s * 1.8) l);
            stroke: none;
        }
        
    }

    .dshw-line {
        scale: 1;
        transition: scale .5s ease-out;
        transform-box: fill-box;

        @starting-style {
            scale: .1 1;  
        }

        &.dshw-line-left {
            transform-origin: right center;
        }

        &.dshw-line-right {
            transform-origin: left center;
        }
    }
</style>
