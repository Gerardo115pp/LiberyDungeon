<script>
    import FieldData from '../../libs/FieldData';
    import { FieldStates } from '../../libs/FieldData';

    // Styles
    export let isClear = false;
    export let isSquared = false;
    export let modeColumn = false;
    export let isTextArea = false;

    export let field_data = new FieldData("generic-input", /[.\n]+/, "any"); 
    export let input_label;
    export let label_color = "var(--libery-input-dark-color)"
    export let show_placeholder = false;
    export let input_padding = "var(--spacing-2) var(--spacing-2)";
    export let input_background = "none";
    export let onKeypress;
    export let onEnterPressed;
    export let onBlur;

    /* CSS variables */
    export let input_color = "var(--primary-color)";
    export let input_dark_color = "var(--dark-light-color)";
    export let border_color = "var(--primary-color)";
    export let title_font = "var(--font-title)";
    export let text_font = "var(--font-text)";
    export let font_size = "var(--font-size-2)";
    export let placeholder_color = "var(--grey-1)";

    export let initial_value;
    export let min = 0;
    export let max;
    export let autocomplete = "off";
    export let autofocus = false;

    $: state_color = getBorderColor(field_data.state) ;
    isTextArea = field_data.type === "textarea";

    let button_pointer = undefined;

    const awaitKeys = e => {
        if (e.key.toLowerCase() === 'enter' && onEnterPressed !== undefined) {

            field_data.getField().blur();
            return onEnterPressed(e);
            
        } else if (onKeypress !== undefined) {
            
            return onKeypress(e);
        }
    }

    const handleOutsideClickDetected = e => {
        if(e.currentTarget === e.target) {
            e.target.getElementsByTagName("input")[0].focus()
        }
    }

    const getBorderColor = state => {
        switch(state) {
            case FieldStates.NORMAL:
                return border_color || "--libery-input-dark-color";
            case FieldStates.HAS_ERRORS:
                return "--danger";
            case FieldStates.READY:
                return "--ready";
            default:
                return "--libery-input-dark-color"
        }
    }

</script>


<div 
    bind:this={button_pointer} 
    on:click={handleOutsideClickDetected}
    style="padding: {input_padding};border-color: var({state_color}); background: {input_background};" 
    class="input-container"
    class:column-container={modeColumn}
    class:squared={isSquared}
    class:clear-input={isClear}

    style:--libery-input-color={input_color}
    style:--libery-input-dark-color={input_dark_color}
    style:--libery-input-border-color={border_color}
    style:--libery-input-title-font={title_font}
    style:--libery-input-text-font={text_font}
    style:--placeholder-color={placeholder_color}
>
    {#if input_label !== undefined}
        <label 
            style:font-size={font_size}
            style:color={label_color}
            for={field_data.id}
        >
            {input_label}
        </label>
    {/if}
    {#if isTextArea}
        <textarea id={field_data.id} 
            style="margin: {input_padding}; font-size: {font_size};"
            placeholder={input_label === undefined || show_placeholder ? field_data.placeholder : ""}
            on:keydown={awaitKeys}
            on:blur={onBlur}
            autofocus={autofocus}
        />
    {:else if field_data.type !== "number"}
        <input id={field_data.id} type={field_data.type}
            style="margin: {input_padding}; font-size: {font_size};"
            placeholder={input_label === undefined || show_placeholder ? field_data.placeholder : ""}
            autocomplete="off"
            autofocus={autofocus}
            on:keydown={awaitKeys}
            on:blur={onBlur}
        />
    {:else}
        <input id={field_data.id} type={field_data.type}
            placeholder={input_label === undefined ? field_data.placeholder : ""}
            on:keydown={awaitKeys}
            on:blur={onBlur}
            autocomplete="off"
            autofocus={autofocus}
            value={initial_value}
            style="font-size: {font_size};"
            min={min}
            max={max}
        />
    {/if}

</div>

<style>
    .input-container {
        box-sizing: border-box;
        cursor: text;
        width: 100%;
        display: flex;
        border: 2px solid var(--libery-input-border-color);
        border-radius: 9999px;
        align-items: center;
    }   

    .input-container.column-container {
        flex-direction: column;
        align-items: flex-start;
    }

    .input-container label {
        font-family: var(--libery-input-title-font);
        color: var(--libery-input-dark-color);
        width: max-content;
        margin-left: 1vw;
        text-transform: capitalize;
        white-space: nowrap;
    }

    .input-container input {
        font-family: var(--libery-input-text-font);
        width: 100%;
        display: flex;
        background: none;
        border: none;
        color: var(--libery-input-color);
        padding: 0;
        align-items: center;
        outline: none;
    }

    .input-container textarea {
        font-family: var(--libery-input-text-font);
        width: 100%;
        display: flex;
        background: none;
        border: none;
        color: var(--libery-input-color);
        padding: 2ch  1ch;
        align-items: center;
        line-height: 1.2;
        outline: none;
    }

    .input-container input::-webkit-calendar-picker-indicator {
        filter: invert(1);
    }

    .input-container.column-container input {
        padding: 2ch 2ch;
    }

    .input-container.column-container label {
        margin-left: 0;
    }

    .input-container input[type='number'] {
        width: 5rem;
        text-align: center;
        -moz-appearance: textfield;
        appearance: textfield;
    }

    .input-container input::placeholder, .input-container textarea::placeholder {
        font-family: var(--libery-input-text-font);
        color: var(--placeholder-color);
        text-transform: unset;
    }

    input::-webkit-outer-spin-button,
    input::-webkit-inner-spin-button {
        -webkit-appearance: none;
    }

    .input-container.clear-input {
        border: none;
        border-bottom: 2px solid var(--libery-input-dark-color);
        border-radius: 0 !important;
        flex-direction: column;
    }

    .input-container.clear-input > label {
        display: block;
        margin-left: -2vw;
        font-weight: regular;
        align-self: flex-start;
    }

    .input-container > input[type='password'] {
        letter-spacing: .5em;
    }

    .input-container > input[type='password']::placeholder {
        letter-spacing: 0.4px;
    }


    .input-container.squared {
        border-radius: 5px;
    }
</style>
