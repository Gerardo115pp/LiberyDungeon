<script>
    import { onDestroy, onMount } from "svelte";
    import { discrete_feedback_message } from "../lf_utils";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** 
         * The amount of time in milliseconds the log will be visible for after a message is sent.
         * @type {number}
         * @default 3000
         */
        export let display_message_duration = 2000; 

        /**
         * The timestamp of when a message was last sent.
         * @type {number}
         */
        let last_message_timestamp = 0;

        /**
         * Whether the log is hidden.
         * @type {boolean}
         */
        let log_hidden = true;

        /**
         * The message currently being displayed.
         * @type {string}
         */
        let current_message = "";

        let discrete_message_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        discrete_message_unsubscriber = discrete_feedback_message.subscribe(handleDiscreteFeedbackMessages);
    })

    onDestroy(() => {
        discrete_message_unsubscriber();
    })
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles new discrete feedback messages.
         * @param {string} new_message
         */
        const handleDiscreteFeedbackMessages = new_message => {
            if (new_message == null || new_message === "") return;
                
            showLog(new_message, display_message_duration);
        } 

        /**
         * Sets a new timeout for the log to be hidden with a given duration
         * @param {number} duration
         */
        const setLogTimeout = duration => {
            last_message_timestamp = Date.now();

            setTimeout(() => {
                const current_timestamp = Date.now();

                if (current_timestamp - last_message_timestamp >= duration){ // Allow the log hiding to be cancelled if a new message is sent before the timeout is reached.
                    hideLog();
                }
            }, duration);
        }
        
        /**
         * Hides the log gracefully.
         */
        const hideLog = () => {
            log_hidden = true;
            discrete_feedback_message.set(""); // This will be skipped by handleDiscreteFeedbackMessages.
        }

        /**
         * Shows the log and sets a timeout for it to be hidden.
         * @param {string} new_message
         * @param {number} duration
         */
        const showLog = (new_message, duration) => {
            current_message = new_message;
            log_hidden = false;
            setLogTimeout(duration);
        }

    /*=====  End of Methods  ======*/
        
</script>

<div class="discrete-feedback-log-wrapper">
    <p class="feedback-log"
        class:log-hidden={log_hidden}
    >
        {current_message}
    </p>
</div>

<style>
    .discrete-feedback-log-wrapper {
        padding: var(--spacing-1);
        color: var(--main-dark);

        & p.feedback-log {
            font-family: var(--font-read);
            background: var(--grey-9);
            font-size: var(--font-size-1);
            line-height: 1;
            padding: var(--spacing-1);
            border-radius: var(--border-radius);
            box-shadow: var(--shadow-1);
        }
    }

    @supports (color: rgb(from white r g b)) {
        .discrete-feedback-log-wrapper p.feedback-log {
            background: hsl(from var(--grey-9) h s l / 0.9);
        }
    }

    .feedback-log.log-hidden {
        opacity: 0;
        transition: opacity 0.5s;
    }
</style>