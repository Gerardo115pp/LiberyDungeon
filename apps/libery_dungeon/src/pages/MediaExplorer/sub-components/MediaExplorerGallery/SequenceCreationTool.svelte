<script>
    import HotkeysContext from '@libs/LiberyHotkeys/hotkeys_context';
    import { me_gallery_changes_manager, me_gallery_yanked_medias } from './me_gallery_state';
    import MeGalleryDisplayItem from './MEGalleryDisplayItem.svelte';
    import { HOTKEYS_HIDDEN_GROUP } from '@libs/LiberyHotkeys/hotkeys_consts';
    import { getHotkeysManager } from '@libs/LiberyHotkeys/libery_hotkeys';
    import { onMount } from 'svelte';
    
    
    /*=============================================
    =            Properties            =
    =============================================*/
        
        /*=============================================
        =            Hotkeys            =
        =============================================*/
        
            const global_hotkeys_manager = getHotkeysManager();

            const hotkeys_context_name = "sequence_creation_tool";
        
        /*=====  End of Hotkeys  ======*/
    
        /**
         * The media items from which the sequence will be created.
         * @type {import('@models/Medias').OrderedMedia[]}
         */ 
        export let unsequenced_medias = [];

        /**
         * Medias per grid row.
         * @type {number}
         */
        let medias_per_row = 6;
        
        /*----------  State  ----------*/
        
            /**
             * Whether to show the sequence as skeleton items.
             * @type {boolean}
             */ 
            let skeleton_sequence = true;

            /**
             * Media focused index.
             * @type {number}
             */
            let sct_focus_index = 0;

            /**
             * Whether the auto select mode has been enabled.
             * @type {boolean}
             */
            let auto_select_mode = false;

            /**
             * Whether auto select mode works for yanking or unyanking.
             * Is determined by whether the media focused at the time of
             * activation of the auto select mode was yanked or not.
             * @type {boolean}
             */
            let auto_select_yanks = false;

            /**
             * Last time the media selection key was pressed.
             * @type {number}
             */
            let last_media_selection_key_press = 0;

            /**
             * The time in milliseconds after which the auto select mode
             * will be enabled.
             * @type {number}
             */
            const AUTO_SELECT_MODE_TIME = 100;
    
    /*=====  End of Properties  ======*/

    onMount(() => {

        defineSCTHotkeys();
    })
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /*=============================================
        =            Hotkeys            =
        =============================================*/
        
            /**
             * defines the sequence creation tool hotkeys.
             */    
            const defineSCTHotkeys = () => {
                if (global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    global_hotkeys_manager.dropContext(hotkeys_context_name);
                }

                const hotkeys_context = new HotkeysContext();

                hotkeys_context.register(["w", "a", "s", "d"], handleSequenceGridMovement, {
                    description: `<navigation> Move the focus on the sequence grid.`,
                });

                hotkeys_context.register(["space"], handleMediaSelection, {
                    description: `<selection> Select the focused media item. If kept pressed long enough, enables auto select mode. The selector turns green when auto select mode is enabled.`,
                    can_repeat: true,
                });

                hotkeys_context.register(["space"], e => e.preventDefault(), {
                    description: `<${HOTKEYS_HIDDEN_GROUP}> hidden`,
                    mode: "keyup",
                });

                global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);

                global_hotkeys_manager.loadContext(hotkeys_context_name);
            }

            /**
             * Handles the sequence grid movement.
             * @param {KeyboardEvent} key_event 
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleSequenceGridMovement = (key_event, hotkey) => {
                const media_count = unsequenced_medias.length;
                const row_count = Math.ceil(media_count / medias_per_row);

                let new_sct_focus_index = sct_focus_index;

                switch (hotkey.KeyCombo) {
                    case "w":
                        new_sct_focus_index -= medias_per_row;

                        new_sct_focus_index = new_sct_focus_index < 0 ? ((row_count - 1) * medias_per_row) + sct_focus_index : new_sct_focus_index;
                        new_sct_focus_index = new_sct_focus_index >= media_count ? new_sct_focus_index - medias_per_row : new_sct_focus_index;
                        break;
                    case "a":
                        new_sct_focus_index--;

                        new_sct_focus_index = new_sct_focus_index < 0 ? media_count - 1 : new_sct_focus_index;
                        break;
                    case "s":
                        new_sct_focus_index += medias_per_row;

                        new_sct_focus_index = new_sct_focus_index >= media_count ? sct_focus_index - ((row_count - 1) * medias_per_row) : new_sct_focus_index;
                        new_sct_focus_index = new_sct_focus_index < 0 ? new_sct_focus_index + medias_per_row : new_sct_focus_index;
                        break;
                    case "d":
                        new_sct_focus_index++;

                        new_sct_focus_index = new_sct_focus_index >= media_count ? 0 : new_sct_focus_index;
                        break;
                }

                let old_sct_focus_index = sct_focus_index;
                sct_focus_index = new_sct_focus_index;
                
                if (auto_select_mode) {
                    autoSelectHandler(old_sct_focus_index, sct_focus_index);
                }
            }

            /**
             * Handles the media selection.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleMediaSelection = (key_event, hotkey) => {
                key_event.preventDefault();

                if (key_event.repeat) {
                    let activated_by_time = ((key_event.timeStamp - last_media_selection_key_press) > AUTO_SELECT_MODE_TIME);
                    let time_valid = last_media_selection_key_press > 0;
                    let not_active = !auto_select_mode;

                    let is_activated = activated_by_time && time_valid && not_active;

                    if ( is_activated ) {
                        activateAutoSelectMode();
                    }
                    return;
                }

                if (auto_select_mode) {
                    console.log("Turning off auto select mode");
                    auto_select_mode = false;
                    last_media_selection_key_press = 0;
                    return;
                }
                
                last_media_selection_key_press = key_event.timeStamp;

                let selected_media = unsequenced_medias[sct_focus_index];

                let is_yanked = isMediaYanked(sct_focus_index);
                if (is_yanked) {
                    unyankSTCMedia(selected_media);
                } else {
                    yankSTCMedia(selected_media);
                }
            }
        
        /*=====  End of Hotkeys  ======*/

        /**
         * Handles the auto select mode. it adds all medias in the range provided(inclusive) to the yanked medias.
         * @param {number} start_index
         * @param {number} end_index
         */
        const autoSelectHandler = (start_index, end_index) => {
            let start = Math.min(start_index, end_index);
            let end = Math.max(start_index, end_index);

            for (let h = start; h <= end; h++) {
                let media_in_index = unsequenced_medias[h];

                if (media_in_index == null) continue;

                let is_yanked = isMediaYanked(h);
                if (auto_select_yanks && !is_yanked) {
                    yankSTCMedia(media_in_index);
                } else if (!auto_select_yanks && is_yanked) {
                    unyankSTCMedia(media_in_index);
                }
            }
        }

        /**
         * Activates the auto select mode.
         */
        const activateAutoSelectMode = () => {
            auto_select_mode = true;

            let selected_media = unsequenced_medias[sct_focus_index];

            auto_select_yanks = isMediaYanked(sct_focus_index);
            console.log("Auto select mode activated, yanks: ", auto_select_yanks);
        }


        /**
         * Returns whether the media on the passed index is yanked.
         * @param {number} index
         * @returns {boolean}
         */
        const isMediaYanked = index => {
            let media = unsequenced_medias[index];

            if (media == null) return false;

            let is_yanked = $me_gallery_yanked_medias.some(ymedia => ymedia.uuid === media.uuid);

            return is_yanked;
        }            


        /**
         * Unyanks the passed media.
         * @param {import('@models/Medias').OrderedMedia} media
         */
        const unyankSTCMedia = (media) => {
            me_gallery_yanked_medias.set($me_gallery_yanked_medias.filter(ymedia => ymedia.uuid !== media.uuid));
        }
            
        

        /**
         * yanks the passed media. If the media is already yanked, it will be unyanked.
         * @param {import('@models/Medias').OrderedMedia} media
         */
        const yankSTCMedia = (media) => {
            me_gallery_yanked_medias.set([...$me_gallery_yanked_medias, media]);
        }
        

    /*=====  End of Methods  ======*/
    
</script>

{#if $me_gallery_changes_manager != null}
    <div id="sequence-creation-tool"
        class:auto-select-enabled={auto_select_mode}
    >
        <div id="sequence-parameters"></div>
        <ul id="sequence-members"
            style:grid-template-columns="repeat({medias_per_row}, minmax(var(--sct-grid-item-size), 1fr))"
        >
            {#each unsequenced_medias as umedia, h}
                {@const sct_keyboard_focused = sct_focus_index === h}
                <li class="sct-sm-member-item"
                    class:is-skeleton={skeleton_sequence}
                    class:keyboard-selected={sct_keyboard_focused}
                    data-member-index={h}
                >
                    <MeGalleryDisplayItem
                        ordered_media={umedia}
                        is_keyboard_focused={sct_keyboard_focused}
                        use_masonry={false}
                        enable_video_titles
                        enable_dragging
                    />
                </li>
            {/each}
        </ul>
    </div>
{/if}

<style>
    ul#sequence-members {
        --sct-grid-item-size: 400px;
        display: grid;
        grid-auto-rows: calc(var(--sct-grid-item-size) * 1.1);
        background: var(--grey);
        gap: var(--spacing-2);
        padding: 4px;
        list-style: none;
        margin: 0;
    }

    li.sct-sm-member-item {
        position: relative;
        cursor: pointer;
        container-type: inline-size;
        background: var(--grey-9);
        width: var(--sct-grid-item-size);
        height: calc(var(--sct-grid-item-size) * 1.1);
        z-index: var(--z-index-1);
    }

    li.sct-sm-member-item.keyboard-selected {
        /* display: none; */
        outline: var(--main) solid 2px;
        z-index: var(--z-index-2);
    }
    
    .auto-select-enabled li.sct-sm-member-item.keyboard-selected {
        outline: var(--success) solid 2px;
    }
</style>