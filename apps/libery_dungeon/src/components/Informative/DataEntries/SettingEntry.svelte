<script>
    
    /*=============================================
    =            Properties            =
    =============================================*/

    
        /**
         * An id to identify this element with. has to be an id-selector compatible string.
         * @type {string}
         */
        export let id_selector;
    
        /**
         * The label of the information entry
         * @type {string}
         * @default 'data'
         */
        export let information_entry_label;

        /**
         * The value of the information entry
         * @type {string}
         * @default ''
         */
        export let information_entry_value = '';

        /**
         * A regex pattern for the value.
         * @type {string | null}
         */
        export let value_pattern = null;

        /**
         * Whether the input field attached to the passed value should be configured as required
         * @type {true | null}
         */
        export let is_required = null;

        
        /*----------  Style  ----------*/
        
            /**
             * the component's default font-size.
             * @type {string}
             * @default "var(--font-size-1)"
             */ 
            export let font_size = "var(--font-size-1)";
        
        /*----------  References  ----------*/
        
            /**
             * The input element used to change the given setting.
             * @type {HTMLInputElement | undefined}
             */
            let the_setting_changer_input;
    
    /*=====  End of Properties  ======*/


    
    /*=============================================
    =            Methods            =
    =============================================*/
    

    
    /*=====  End of Methods  ======*/
    
    
    
</script>

<li class="dungeon-setting-entry" 
    style:font-size="{font_size}"
>
    <label class="dungeon-input">
        <span class="dungeon-label">
            {information_entry_label}
        </span>
        <input 
            id={id_selector}
            bind:this={the_setting_changer_input}
            type="text"
            bind:value={information_entry_value}
            spellcheck="false"
            pattern={value_pattern}
            required={is_required}
        >
    </label>
</li>

<style>
    li.dungeon-setting-entry {
        & label.dungeon-input {
            font-size: inherit;
            padding-inline: 1em;
            padding-block: .4em;
            line-height: 1;
        }
        
        & label.dungeon-input input {
            font-size: inherit;
            color: var(--grey-3);
        }
    }

    @supports (color: rgb( from white r g b / 1)) {
        li.dungeon-setting-entry label {
            padding: var(--vspacing-1) var(--vspacing-2);
            background: hsl(from var(--grey) h s l / 0.8);
        }
    }
</style>

