<script>
    import { onMount } from 'svelte';
    import { HOTKEYS_GENERAL_GROUP, HOTKEYS_HIDDEN_GROUP } from '@libs/LiberyHotkeys/hotkeys_consts';

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The list of hotkeys to display.
         * @type {import('@libs/LiberyHotkeys/hotkeys_context').HotkeyData[]}
        */
        export let current_hotkeys;

        /**
         * The name of the current context.
         * @type {string}
         */
        export let context_name = "";


        /**
         * @typedef {Object} HotkeyMetadata
         * @property {string} hotkey
         * @property {string} description
         * @property {string} group parsed from the description. e.g  if the description starts with `<group_name>` then the group is `group_name`
         */

        /**
         * @type {Object<string,HotkeyMetadata[]>} The metadata of the hotkeys to display. Grouped by the group name.
         */
        let hotkey_metadata = {};

        /**
         * The list of hotkeys to display in the general hotkeys group. which is the one that has no name in the ui.
         * @type {HotkeyMetadata[]}
         */
        let general_hotkeys = [];

        
        /*----------  Styling  ----------*/
        
            /**
             * The padding of the hotkeys table.
             * @type {string}
             */
            let static_padding_width = "var(--vspacing-2)";

            /**
             * The count of characters of the longest hotkey description, used to calculate the width of the table.
             * @type {number}
             */
            let max_description_length = 0;
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        console.log('current_hotkeys', current_hotkeys);

        hotkey_metadata = getHotkeysMetadata(current_hotkeys);

        general_hotkeys = hotkey_metadata[HOTKEYS_GENERAL_GROUP] ?? [];
        delete hotkey_metadata[HOTKEYS_GENERAL_GROUP];
    });

    
    /*=============================================
    =            Methods                =
    =============================================*/
    
        /**
         * Parses the hotkeys metadata from the hotkeys context.
         * @param {import('@libs/LiberyHotkeys/hotkeys_context').HotkeyData[]} new_hotkeys
         * @returns {Object<string,HotkeyMetadata[]>}
         */
        const getHotkeysMetadata = new_hotkeys => {
            console.log('new_hotkeys', new_hotkeys);
            const new_metadata = {};
            const same_description_lookup = new Map();
            let new_max_description_length = 0;

            for (let hotkey_data of new_hotkeys) {
                const metadata = parseHotkeyDataToMetadata(hotkey_data);

                if (metadata.group === HOTKEYS_HIDDEN_GROUP) {
                    continue;
                }

                if (!same_description_lookup.has(metadata.description)) {
                    new_max_description_length = Math.max(new_max_description_length, metadata.description.length);
                    new_metadata[metadata.group] = new_metadata[metadata.group] != null ? [...new_metadata[metadata.group], metadata] : [metadata];
                    same_description_lookup.set(metadata.description, metadata);
                } else {
                    let repeated_metadata = same_description_lookup.get(metadata.description);
                    repeated_metadata.hotkey += `, ${metadata.hotkey}`; 
                }

            }

            max_description_length = new_max_description_length;

            console.log('new_metadata', new_metadata);  

            return new_metadata;
        }

        /**
         * Converts a hotkey data to a hotkey metadata.
         * @param {import('@libs/LiberyHotkeys/hotkeys_context').HotkeyData} hotkey_data
         * @returns {HotkeyMetadata}
         */ 
        const parseHotkeyDataToMetadata = hotkey_data => {
            /** @type {HotkeyMetadata} */
            let hotkey_metadata = {
                hotkey: `'${hotkey_data.name}'`,
                description: hotkey_data.Description, // Getter that removes the group tag.
                group: hotkey_data.Group // Getter that returns an extracted group name or 'General' as default.
            }

            return hotkey_metadata;
        }
        
    
    /*=====  End of Methods  ======*/
    
</script>

{#if current_hotkeys != null}
    <div id="hotkeys-table-wrapper"
        style:width={`calc(calc(${max_description_length}ch * 0.9) + ${static_padding_width})`} 
    >
        <h5 id="hotkeys-context-name">
            <span>{context_name?.replaceAll('_', ' ')}</span> active hotkeys
        </h5>
        <ul id="hotkeys-table">
            <ul class="hotkeys-group-list">
                {#each general_hotkeys as hotkey_metadata}
                    <li class="hotkey-item">
                        <kbd class="hotkey-name">
                            {hotkey_metadata.hotkey}
                        </kbd>
                        <span class="hotkey-description">
                            {hotkey_metadata.description}
                        </span>
                    </li>
                {/each}
            </ul>
            {#each Object.keys(hotkey_metadata) as group_name}
                <li class="hotkeys-group">
                    <h6 class="hotkeys-group-name">
                        {group_name.replaceAll('_', ' ')}
                    </h6>
                    <ul class="hotkeys-group-list">
                        {#each hotkey_metadata[group_name] as hotkey_metadata}
                            <li class="hotkey-item">
                                <kbd class="hotkey-name">
                                    {hotkey_metadata.hotkey}
                                </kbd>
                                <span class="hotkey-description">
                                    {hotkey_metadata.description}
                                </span>
                            </li>
                        {/each}
                    </ul>
                </li>
            {/each}
        </ul>
    </div>
{/if}

<style>
    #hotkeys-table-wrapper {
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        max-width: 70dvw;
        container-type: inline-size;
        overflow: auto;
        max-height: 70dvh;
        row-gap: var(--vspacing-3);
        background: var(--grey-9);
        padding: var(--vspacing-1);
        scrollbar-width: thin;
        scrollbar-color: var(--grey-7) var(--grey-9);

        & ul {
            list-style-type: none;
            padding: 0;
            margin: 0;
        }
    }

    h5#hotkeys-context-name {
        font-family: var(--font-read);
        font-size: var(--font-size-3);
        text-transform: capitalize;

        & span {
            font-style: italic;
            color: var(--main);
        }
    }

    ul#hotkeys-table {
        display: flex;
        flex-direction: column;
        row-gap: var(--vspacing-2);
    }

    
    /*----------  Hotkey group  ----------*/
    
        li.hotkeys-group {
            display: flex;
            flex-direction: column;
            row-gap: calc(var(--vspacing-1) * 1.2);
        }

        h6.hotkeys-group-name {
            font-family: var(--font-read);
            font-size: var(--font-size-2);
        }

        ul.hotkeys-group-list {
            display: flex;
            flex-direction: column;
        }

    
    /*----------  Hotkey item  ----------*/

        li.hotkey-item {
            display: grid;
            width: 100%;
            grid-template-columns: 1fr 2fr;
            column-gap: var(--vspacing-1);
            align-items: start;
            padding: 0.5ch 0;
            justify-items: left;

            &:not(:last-child) {
                border-bottom: 1px dashed var(--grey-6);
            }
        }

        kbd.hotkey-name {
            background: var(--grey);
            color: var(--main-dark);
            padding: 0 1ex;
            border-radius: 0.4ex;
        }

        span.hotkey-description {
            width: 100%;
            
            word-break: break-all   ;
        }
</style>