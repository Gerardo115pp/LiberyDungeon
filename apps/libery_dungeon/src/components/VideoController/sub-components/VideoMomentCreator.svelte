<script>
    import { onMount } from "svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * Starting name. 
         * @type {string}
         */ 
        export let starting_name = "";

        /**
         * The labels text.
         * @type {string}
         * @default 'New moment'
         */
        export let creator_label = "New moment";

        /**
         * The video moment creator input
         * @type {HTMLInputElement}
         */
        let the_video_moment_creator;

        /**
         * Whether the name is valid.
         * @type {boolean}
         */
        let moment_name_is_valid = starting_name != "";

        
        /*----------  Event handlers  ----------*/
        
            /**
             * A callback that is triggered when ever the user commits to the new moment's name.
             * @type {(new_name: string) => void}
             */ 
            export let onNameCommitted = (new_name) => {};
            
            /**
             * A callback that is called if the user decides to cancel the moment creation.
             * @type {() => void}
             */
            export let onCancel = () => {};

    
    /*=====  End of Properties  ======*/
    
    onMount(() => {
        the_video_moment_creator.focus();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the keyup event from the moment creator input.
         * @param {KeyboardEvent} event
         */ 
        const handleKeyUp = event => {
            if (the_video_moment_creator === null) return;

            event.preventDefault();

            if (the_video_moment_creator.validationMessage !== "") {
                the_video_moment_creator.setCustomValidity("");
            }

            moment_name_is_valid = the_video_moment_creator.checkValidity();
        }

        /**
         * Handles the keydown event from the moment creator input.
         * @param {KeyboardEvent} event
         */
        const handleKeyDown = event => {
            if (the_video_moment_creator == null) return;

            if (event.key === "Escape") {
                the_video_moment_creator.blur();
                return;
            }

            if (event.key === "Enter") {
                event.preventDefault();

                if (!moment_name_is_valid) {
                    the_video_moment_creator.reportValidity();
                    return;
                }

                onNameCommitted(the_video_moment_creator.value);
            }
        }

        /**
         * Handles the blur event from the moment creator input.
         * @param {FocusEvent} event
         */
        const handleBlur = event => {
            onCancel();
        }
    
    /*=====  End of Methods  ======*/
    
    
    
</script>

<label id="lvc-new-video-moment-creator" class="dungeon-input">
    <span class="dungeon-label">{creator_label}</span>
    <input 
        id="lvc-new-video-moment-creator"
        bind:this={the_video_moment_creator}
        type="text"
        bind:value={starting_name}
        on:keydown={handleKeyDown}
        on:keyup={handleKeyUp}
        on:blur={handleBlur}
        spellcheck="true"
        autofocus
    >
</label>

<style>
    #lvc-new-video-moment-creator {
        background: none;
        border: none;

        & span.dungeon-label {
            font-weight: 600;
        }

        & input {
            color: var(--main);
        }
    }
</style>