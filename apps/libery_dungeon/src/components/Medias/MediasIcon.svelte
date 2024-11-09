<script>
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {number} to display the number of medias contained in this category*/
        export let images_count;
        /** @type {boolean} whether the category is currently focused by the keyboard selection */
        export let keyboard_focused = false;
        /** @type {boolean} whether to enter the media viewer on content click or emit an event*/
    
    
    /*=====  End of Properties  ======*/

</script>

<li class="ce-inner-medias">
    <button id="medias-link-wrapper" class:keyboard-focused={keyboard_focused}>
        <svg class="ce-im-icon" viewBox="0 0 110 110">
            <path class="ce-im-icon-hand" transform="translate(8, 0)" d="M30 80L32 47Q35 37 38 45L38 48L39.5 44Q42.5 34 45.5 45L45.5 47L46.7 42.5Q51 32 52.6 46L54.4 28q3 -8 5.2 0L60 80"/>
            <path class="ce-im-icon-eye-big" d="M5 55Q55 7 105 55Q55 103 5 55"/>
            <path class="ce-im-icon-eye-small" d="M40 58.4c3 -4 7.5 -10 15 -10s10.7 8 12.7 10q-15.7 17 -27 0Z"/>
            <ellipse class="ce-im-icon-iris-outter" cx="54" cy="57" rx="4" ry="4"/>
            <ellipse class="ce-im-icon-iris-inner" cx="54" cy="57" rx="2" ry="2"/>
            <path class="ce-im-icon-eyebrows-outter" d="M11 50l-4 -3M19 43l-2 -3.3M30 36.6l-2 -4.3M43 31.6l-2 -3.8m33 7l2 -3.5m10 7.5l2 -3.5m10 11.5l2 -3M11 60l-4 3m11 3l-4 3m11 2l-4 3m57 0l3 3m7 -8l3 3m7 -9l3 3"/>
            <path class="ce-im-icon-eyebrows-inner" d="M45.1 64l-2 2M53 66.4l0 3.7m9.4 -5.5l3 2"/>
        </svg>
        <div class="ce-im-label">
            <h3>Medias <strong>{images_count}</strong></h3>
        </div>
    </button>
</li>

<style>
    .ce-inner-medias {
        width: 100%;
        transition: all .3s ease-in;
        padding: var(--vspacing-1);
        border: .5px solid var(--grey);
    }

    :has(> .keyboard-focused) {
        /* background: var(--grey-8); */
        border-color: var(--main);
    }
    
    .ce-inner-medias button {
        display: flex;
        flex-direction: column;
        align-items: center;
        background: transparent;
        border: none;
        text-decoration: none;
    }

    @media(pointer: fine) {
        .ce-inner-medias:hover {
            /* background: var(--grey-8); */
            border-color: var(--grey-8);
            backdrop-filter: brightness(1.3);
        }
    }

    .ce-im-icon {
        width: 100%;
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    .ce-im-icon path {
        stroke: var(--main-dark-color-6);
        fill: none;
        stroke-width: 1.4%;
        stroke-linecap: round;
        stroke-linejoin: round;
    }

    .ce-im-icon ellipse {
        stroke: var(--main-dark-color-6);
        fill: none;
        stroke-width: 1.2%;
    }

    .ce-im-icon .ce-im-icon-eye-big {
        stroke-width: 1.5%;
    }

    .ce-im-label {
        width: max-content;
    }

    .ce-im-label h3{
        width: max-content;
        font-family: var(--font-read);
        text-align: center;
        font-size: clamp(var(--font-size-1), 4%, var(--font-size-3));
    }

    .ce-im-label h3 strong {
        color: var(--main-dark-color-5);
    }
</style>