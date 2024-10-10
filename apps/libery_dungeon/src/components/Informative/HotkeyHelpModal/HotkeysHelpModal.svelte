<script>
    import { hotkeys_context_events, suscribeHotkeysContextEvents } from "@libs/LiberyHotkeys/hotkeys_events";
    import global_hotkeys_manager from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysTable from "./sub-components/HotkeysTable.svelte";
    import { hotkeys_sheet_visible } from "@stores/layout";
    import { onDestroy, onMount } from "svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/

        /**
         * @type {boolean} If set to true, the should remain visible unless explicitly specified by the user(e.g clicks again on a help button). This variable does
         * not govern whether other components can change the visibility of the modal.
         */
        let modal_static = false;

        /**
         * The list of hotkeys to display.
         * @type {import('@libs/LiberyHotkeys/hotkeys_context').HotkeyData[]}
         */
        let current_hotkeys = [];

        /**
         * The name of the current context.
         * @type {string}
         */
        let current_context_name;

        /**
         * @type {Function} Context change event unsubscriber.
         */
        let context_change_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        context_change_unsubscriber = suscribeHotkeysContextEvents(hotkeys_context_events.CONTEXT_CHANGED, handleHotkeysContextChange);

        updateHotkeysValues(global_hotkeys_manager.ContextName);
    });

    onDestroy(() => {
        context_change_unsubscriber();
    });

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the hotkeys context change event.
         * @type {import('@libs/LiberyHotkeys/hotkeys_events').HotkeysContextEventCallback}
        */
        const handleHotkeysContextChange = e => {
            if (!global_hotkeys_manager.hasLoadedContext()) {
                setEmptyContextState();
                return;
            }

            // If the context name is the same and the hotkeys count is the same, do nothing.
            if (e.detail.context_name === current_context_name && current_hotkeys.length === global_hotkeys_manager.Context?.hotkeys.length) {
                return;
            }

            updateHotkeysValues(e.detail.context_name);
        }

        /**
         * Toggles the visibility of the hotkeys context modal.
         * @param {MouseEvent} e The event that triggered the toggle.
         */
        const handleMouseEnterMouseLeave = e => {
            if (modal_static) {
                return;
            }

            hotkeys_sheet_visible.set(e.type === "mouseenter");
        }
        
        /**
         * Toggles the visibility of the hotkeys context modal and sets the modal to be static.
         * @param {MouseEvent} e The event that triggered the toggle.
         */
        const handleModalButtonClick = e => {
            modal_static = !modal_static;

            // we check first to prevent unnecessary reactivity updates.
            if (modal_static !== $hotkeys_sheet_visible) {
                hotkeys_sheet_visible.set(modal_static);
            }
        }

        /**
         * Sets 'current_hotkeys' and 'current_context_name' to the current hotkeys context values.
         * @param {string} context_name The name of the current context.
         */
        const updateHotkeysValues = context_name => {
            if (global_hotkeys_manager.Context == null || global_hotkeys_manager.Context.hotkeys.length === 0) {
                console.warn(`Context is null: ${global_hotkeys_manager.Context == null}\nHotkeys count: ${global_hotkeys_manager.Context?.hotkeys.length}`);
                return;
            }   

            current_hotkeys = [...global_hotkeys_manager.Context?.hotkeys]; // Change the reference for reactivity.
            current_context_name = context_name;
        }

        /**
         * Set empty context state
         */
        const setEmptyContextState = () => {
            current_hotkeys = [];
            current_context_name = "";
        }    
    /*=====  End of Methods  ======*/
    
    

</script>

<div id="hotkeys-context-info-wrapper">
    {#if $hotkeys_sheet_visible && current_hotkeys.length > 0}
        {#key current_hotkeys}
            <div id="hotkeys-table-wrapper">
                <HotkeysTable 
                    context_name={current_context_name} 
                    current_hotkeys={current_hotkeys}
                /> 
            </div>            
        {/key}
    {/if}
    <button id="hotkeys-context-info-toggle" on:click={handleModalButtonClick} on:mouseenter={handleMouseEnterMouseLeave} on:mouseleave={handleMouseEnterMouseLeave}>
        ?
    </button>
</div>

<style>
    #hotkeys-context-info-wrapper {
        --modal-position-modifier: 0.90;

        position: fixed;
        inset: calc(100vh * var(--modal-position-modifier)) auto auto calc(105dvw * var(--modal-position-modifier));
        z-index:  var(--z-index-t-7);
    }
    
    #hotkeys-table-wrapper {
        position: fixed;
        inset: 15% 5% auto auto;
        contain: content;
    }
    
    #hotkeys-context-info-wrapper button#hotkeys-context-info-toggle {
        --hover-color: color-mix(in srgb, var(--main-dark-color-9), var(--grey-8) 40%);
        --toggle-button-size: clamp(10px, 2.4em, 10vw);

        display: flex;
        font-family: var(--font-read);
        width: var(--toggle-button-size);
        height: var(--toggle-button-size);
        background: var(--grey-9);
        color: var(--grey-1);
        font-size: var(--font-size-3);
        align-items: center;
        justify-content: center;
        border: 1px solid var(--main-dark-color-7);
        padding: 0;
        line-height: 1;
        outline: none;
        transition: all 300ms ease;

        &:hover {
            background: var(--hover-color);
        }
    }
</style>