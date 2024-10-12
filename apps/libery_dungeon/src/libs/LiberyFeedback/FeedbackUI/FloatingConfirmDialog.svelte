<script>
    import { onDestroy, onMount, tick } from "svelte";
    import { setConfirmResponse, confirm_message } from "../lf_utils";
    import { readable, writable, readonly } from "svelte/store";
    import { HOTKEYS_HIDDEN_GROUP, HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import global_hotkeys_manager from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { hotkeys_sheet_visible, layout_properties } from "@stores/layout";
    import Page from "@app/routes/+page.svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/

        
        /*=============================================
        =            Hotkeys            =
        =============================================*/
        
            /**
             * The hotkeys context name to identify the confirm dialog hotkeys context.
             * It also serves as it's title on the hotkeys cheat sheet. the hotkeys cheat sheet
             * will replace '_' with spaces.
             * @type {string}
             */
            const hotkeys_context_name = "confirm_dialog"; 
        
        /*=====  End of Hotkeys  ======*/
        
        /**
         * The message question to confirm.
         * @type {import('../lf_models').ConfirmMessage}
         * @default null
         */
        export let the_message = null;

        /**
         * The response of the confirm dialog.
         * 0 = cancel, 1 = confirm, -1 = undecided.
         * any other value means we are still waiting for the user to choose or if there was a timeout set, then waiting for the
         * user or the timeout to elapse.
         * @type {import('svelte/store').Writable<number>}
        */
        export let the_response = writable(-2);
        
        
        /*----------  State  ----------*/
        
            /**
             * Current focused choice index.
             * @type {number}
             */
            let focused_choice_index = 0;   

        
        
        /*----------  Style  ----------*/

            /**
             * Whether the dialog should be floating at the bottom or top of the screen.
             * @type {boolean}
             * @default false
             */
            export let float_at_top = false;
        
        

        let confirm_message_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        confirm_message_unsubscriber = confirm_message.subscribe(handleNewMessage);

        setupConfirmResponseGlobalStore();
    });


    onDestroy(() => {
        confirm_message_unsubscriber();
        setComponentHotkeysContext(false);
    });


    
    /*=============================================
    =            Methods            =
    =============================================*/
            
            /*=============================================
            =            Hotkeys            =
            =============================================*/
                // Only call directly defineComponentHotkeys. All the other methods in this section should only be called by the global hotkeys manager.

                /**
                 * Defines the hotkeys to interact with the floating confirm dialog.
                 */        
                const defineComponentHotkeys = () => {
                    if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                        const hotkeys_context = new HotkeysContext();

                        hotkeys_context.register(["a", "d"], handleChoiceMovement, {
                            description: `<choosing> Change the focused choice.`,
                        });

                        hotkeys_context.register(["2 2"], handleSelectConfirmChoice, {
                            description: `<choosing> Accepts the dialog question, continuing with the action.`,
                        });

                        hotkeys_context.register(["x x"], handleCancelConfirmChoice, {
                            description: `<choosing> Cancels/Rejects the dialog question and stops the current action.`,
                        });

                        hotkeys_context.register(["e"], handleSelectFocusedChoice, {
                            description: `<choosing> Confirms the focused choice.`,
                        });

                        hotkeys_context.register(["?"], () => hotkeys_sheet_visible.set(!$hotkeys_sheet_visible), {
                            description: `<${HOTKEYS_GENERAL_GROUP}> Toggle the hotkeys cheat sheet.`,
                        });
                        

                        global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
                    }

                    global_hotkeys_manager.loadContext(hotkeys_context_name);
                }

                /**
                 * Handles the movement of the focused choice.
                 * @param {KeyboardEvent} event
                 * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
                 */
                const handleChoiceMovement = (event, hotkey) => {
                    let key_combo = hotkey.key_combo.toLowerCase();

                    if (key_combo === "a") {
                        focused_choice_index = 0;
                    } else if (key_combo === "d") {
                        focused_choice_index = 1;
                    }
                }

                /**
                 * Handles the selection of the confirm choice.
                 */
                const handleSelectConfirmChoice = () => {
                    focused_choice_index = 1;
                    setFinalChoice(focused_choice_index);
                }

                /**
                 * Handles the cancel choice.
                 */ 
                const handleCancelConfirmChoice = () => {
                    focused_choice_index = 0;
                    setFinalChoice(focused_choice_index);
                }

                /**
                 * Handles the selection of the focused choice.
                 */
                const handleSelectFocusedChoice = () => {
                    setFinalChoice(focused_choice_index);
                }
            
            /*=====  End of Hotkeys  ======*/

            /**
             * Closes the confirm dialog by setting `the_message` to null.
             */
            const closeConfirmDialog = () => {
                the_message = null;
            }   
            

            /**
             * drops the hotkeys context.
             */
            const dropComponentHotkeysContext = () => {
                global_hotkeys_manager.dropContext(hotkeys_context_name);
            }

            /**
             * Handles new confirm message.
             * @param {import('../lf_models').ConfirmMessage} new_message
             */
            const handleNewMessage = (new_message) => {
                const is_message_null = new_message === null;
                
                if (is_message_null) return;

                setComponentHotkeysContext(true);

                focused_choice_index = new_message.AutoFocusCancel ? 0 : 1;
                the_message = new_message;
            }

            /**
             * Handles the click on a choice.
             * @param {MouseEvent} event
             */
            const handleChoiceClick = (event) => {
                let choice = event.currentTarget.dataset.choice;

                choice = parseInt(choice);  

                setFinalChoice(choice);
            }
    
            /**
             * initialize the global confirm store by creating a readable store from the 'the_response' store.
             * @requires the_response
             * @requires setConfirmResponse
             */
            const setupConfirmResponseGlobalStore = () => {
                let readable_response = readonly(the_response);

                setConfirmResponse(readable_response);
            }

            /**
             * Sets a choice as the final question answer.
             * true is accept, false is cancel and null is undecided.
             * @param {number} choice 
             */
            const setFinalChoice = async (choice) => {
                setComponentHotkeysContext(false);
                the_response.set(choice);
                closeConfirmDialog();
                await tick();
                the_response.set(-2);
            }

            /**
             * sets the component hotkeys context.
             * @param {boolean} enable
             */
            const setComponentHotkeysContext = enable => {                
                if (enable) {
                    defineComponentHotkeys(); // if not defined, defines it and it loads it regardless.
                    global_hotkeys_manager.lockContextControl();
                } else {
                    global_hotkeys_manager.loadPreviousContext();
                    global_hotkeys_manager.unlockContextControl();
                }
            }
                
    
    /*=====  End of Methods  ======*/

</script>

<dialog id="libery-feedback-floating-confirm-message-dialog"
    open={the_message != null}  
    class:float-at-top={float_at_top}
    class:is-success={the_message?.DangerLevel === -1}
    class:is-informational={the_message?.DangerLevel === 0}
    class:is-warning={the_message?.DangerLevel === 1}   
    class:is-danger={the_message?.DangerLevel === 2}
>
    {#if the_message != null}
        <div class="lffcmd-the-content-wrapper">
            <header class="lffcmd-the-title">
                <h2>
                    {the_message.MessageTitle}
                </h2>
            </header>
            <p class="lffcmd-the-question">
                {the_message.QuestionMessage}
            </p>
            <menu class="lffcmd-the-awnsers">
                {#if the_message.CancelLabel !== ""}
                    <li class="lffcmd-an-answer" class:focused-choice={focused_choice_index === 0}>
                        <button data-choice={0} class="lffcmd-awnser-picking-btn" on:click={handleChoiceClick}>
                            {the_message.CancelLabel}
                        </button>
                    </li>
                {/if}
                <li class="lffcmd-an-answer" class:focused-choice={focused_choice_index === 1}>
                    <button data-choice="1" class="lffcmd-awnser-picking-btn" on:click={handleChoiceClick}>
                        {the_message.ConfirmLabel || "accept"}
                    </button>
                </li>
            </menu>
        </div>        
    {/if}
</dialog>

<style>
    dialog#libery-feedback-floating-confirm-message-dialog {
        display: none;
        width: 100dvw;
        height: 100dvh;
        position: fixed;
        inset: 0;
        place-items: center;
        backdrop-filter: var(--backdrop-filer-blur);
        z-index: var(--z-index-t-6);
    }

    dialog#libery-feedback-floating-confirm-message-dialog[open] {
        display: grid;
    }

    /*=============================================
    =            Danger levels            =
    =============================================*/
    
        dialog#libery-feedback-floating-confirm-message-dialog.is-informational {
            background: hsl(from var(--accent-8) h s l / 0.2);

            & .lffcmd-the-content-wrapper {
                border: 1px solid var(--main-dark-color-5);
            }

            & header.lffcmd-the-title h2 {
                color: var(--main-dark);
            }
        }

        dialog#libery-feedback-floating-confirm-message-dialog.is-success {
            background: hsl(from var(--success-9) h s 10% / 0.2);

            & .lffcmd-the-content-wrapper {
                border: 1px solid var(--success-3);
            }

            & header.lffcmd-the-title h2 {
                color: var(--success-1);
            }
        }

        dialog#libery-feedback-floating-confirm-message-dialog.is-warning {
            background: hsl(from var(--warning-8) h s l / 0.2);

            & .lffcmd-the-content-wrapper {
                border: 1px solid var(--warning);
            }

            & header.lffcmd-the-title h2 {
                color: var(--warning-3);
            }
        }

        dialog#libery-feedback-floating-confirm-message-dialog.is-danger {
            background: hsl(from var(--danger-8) h s l / 0.2);

            & .lffcmd-the-content-wrapper {
                border: 1px solid var(--danger-7);
            }

            & header.lffcmd-the-title h2 {
                color: var(--danger-6);
            }
        }
    
    
    /*=====  End of Danger levels  ======*/

    dialog#libery-feedback-floating-confirm-message-dialog[open] .lffcmd-the-content-wrapper {
        scale: 1;

        @starting-style {
            scale: 0;
        }
    }

    .lffcmd-the-content-wrapper {
        background: var(--grey-9);
        display: flex;
        flex-direction: column;
        align-items: center;
        width: max-content;
        min-width: 300px;
        max-width: 560px;
        padding: var(--vspacing-1) var(--vspacing-3);
        border-radius: var(--border-radius);
        scale: 0;
        transition: scale .3s ease-in, display .3s ease allow-discrete;
        gap: var(--spacing-1);
    }

    .lffcmd-the-title {
        padding: var(--spacing-2) 0;
    }

    .lffcmd-the-title h2 {
        font-family: var(--font-read);
        font-size: var(--font-size-3);
        text-align: center;
    }

    p.lffcmd-the-question {
        font-family: var(--font-read);
        font-size: var(--font-size-p-small);
        color: var(--grey-1);
        text-align: center;
    }
    

    menu.lffcmd-the-awnsers {
        display: flex;
        justify-content: center;
        gap: var(--spacing-3);
        padding: var(--spacing-1) 0;
    }

    li.lffcmd-an-answer button.lffcmd-awnser-picking-btn {
        font-family: var(--font-titles);
        color: var(--grey-1);
        transition: color .3s ease-out;

        &:hover {
            color: var(--success) !important; 
        }
    }

    li.lffcmd-an-answer.focused-choice button.lffcmd-awnser-picking-btn {
        color: var(--main-dark);
    }
</style>