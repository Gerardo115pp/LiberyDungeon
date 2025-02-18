<script>
    import HotkeysContext from '@libs/LiberyHotkeys/hotkeys_context';
    import { me_gallery_changes_manager, me_gallery_yanked_medias } from './me_gallery_state';
    import MeGalleryDisplayItem from './MEGalleryDisplayItem.svelte';
    import { HOTKEYS_HIDDEN_GROUP, HOTKEYS_GENERAL_GROUP } from '@libs/LiberyHotkeys/hotkeys_consts';
    import { getHotkeysManager } from '@libs/LiberyHotkeys/libery_hotkeys';
    import { createEventDispatcher, onMount } from 'svelte';
    import { LabeledError } from '@libs/LiberyFeedback/lf_models';
    import { current_category } from '@stores/categories_tree';
    import { cleanFilenameString } from '@libs/utils';
    import { sequenceRenameMedias } from '@models/Medias';
    import { hotkeys_sheet_visible } from '@stores/layout';
    import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
    
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
        let medias_per_row = (globalThis.innerWidth ?? 1400) > 1920 ? 9 : 6;
        
        /*----------  State  ----------*/
        
            /**
             * Whether to show the sequence as skeleton items.
             * @type {boolean}
             */ 
            let skeleton_sequence = true;

            /**
             * The the name prefix the sequence will be saved with. for a sequence e.g episode_01, episode_02, etc.
             * the prefix would be 'episode_'.
             * @type {string}
             */
            let sequence_prefix = "";

            /**
             * Sequence editor.
             * @type {HTMLInputElement}
             */
            let the_sequence_prefix_editor;

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

            /**
             * The last saved sequence state, by default is the way the unsequenced_medias were passed.
             * @type {import('@models/Medias').OrderedMedia[]}
             */
            let last_saved_sequence = [];

            /**
             * Whether the magnify mode is enabled.
             * @type {boolean}
             * @default false
             */
            let magnify_mode = false;

        const dispatch = createEventDispatcher();

    /*=====  End of Properties  ======*/

    onMount(() => {
        generateInitialSequencePrefix();

        defineSCTHotkeys();
    });
    
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
                if (global_hotkeys_manager == null) {
                    console.error("In SequenceCreationTool.defineSCTHotkeys, hotkeys manager is not available.");
                    return;
                }
                
                if (global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    global_hotkeys_manager.dropContext(hotkeys_context_name);
                }

                const hotkeys_context = new HotkeysContext();

                hotkeys_context.register(["q"], handleGalleryClose, {
                    description: `<${HOTKEYS_GENERAL_GROUP}> Close the sequence creation tool.`,
                    await_execution: false,
                    mode: "keyup",
                })

                hotkeys_context.register(["w", "a", "s", "d"], handleSequenceGridMovement, {
                    description: `<navigation> Move the focus on the sequence grid.`,
                });

                hotkeys_context.register(["space"], handleMediaSelection, {
                    description: `<reordering> Select the focused media item. If kept pressed long enough, enables auto select mode. The selector turns green when auto select mode is enabled.`,
                    can_repeat: true,
                });

                hotkeys_context.register(["space"], e => e.preventDefault(), {
                    description: `<${HOTKEYS_HIDDEN_GROUP}> hidden`,
                    mode: "keyup",
                });

                hotkeys_context.register(["i"], handleInsertYankedMediasHotkey, {
                    description: `<reordering> Insert the yanked medias before the focused media.`,
                });

                hotkeys_context.register(["u"], handleRestoreLastSequence, {
                    description: `<reordering> Restore the order you had before the last change.`,
                });

                hotkeys_context.register(["\\d g"], handleGotoMedia, {
                    description: `<navigation> Go to the media with the given index.`,
                });

                hotkeys_context.register(["shift+w"], handleGoToTop, {
                    description: `<navigation> Go to the first media in the sequence.`,
                });

                hotkeys_context.register(["shift+s"], handleGoToBottom, {
                    description: `<navigation> Go to the last media in the sequence.`, 
                });

                hotkeys_context.register(["t"], handleFocusOnSequencePrefix, {
                    description: `<editing> Focus on the sequence prefix editor. Esc exits, as with any other input.`,
                    mode: "keyup",
                });

                hotkeys_context.register(["n"], handleClearSCT, {
                    description: `<reordering> Clear the current selection.`,
                });

                hotkeys_context.register(["m"], handleToggleMagnifyMode, {
                    description: `<media> Toggle the magnify mode.`,
                });

                hotkeys_context.register(["; s"], handleSaveSequence, {
                    description: `<saving> Save the sequence.`,
                });

                hotkeys_context.register(["?"], e => hotkeys_sheet_visible.set(!$hotkeys_sheet_visible), {
                    description: `<${HOTKEYS_GENERAL_GROUP}>Toggle the hotkeys cheatsheet`,
                });

                global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);

                global_hotkeys_manager.loadContext(hotkeys_context_name);
            }

            /**
             * Clears the sequence creation tool. Called with a hotkey.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleClearSCT = (key_event, hotkey) => {
                me_gallery_yanked_medias.set([]);
            }

            /**
             * Handles the close gallery hotkey.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleGalleryClose = (key_event, hotkey) => {
                if (global_hotkeys_manager == null) {
                    console.error("In SequenceCreationTool.handleGalleryClose, hotkeys manager is not available.");
                    return;
                }
                
                global_hotkeys_manager.loadPreviousContext();

                global_hotkeys_manager.dropContext(hotkeys_context_name);

                emitCloseGallery();
            }
            
            
            /**
             * Focuses on the sequence prefix editor.
             */
            const handleFocusOnSequencePrefix = () => {
                the_sequence_prefix_editor.focus();
            }

            /**
             * Moves the focus to the typed media index.
             * @param {KeyboardEvent} key_event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleGotoMedia = (key_event, hotkey) => {
                if (!hotkey.WithVimMotion || !hotkey.HasMatch) return;

                // @ts-ignore - HasMatch === true, means hotkey.MatchMetadata is not null
                let media_index = hotkey.MatchMetadata.MotionMatches[0]

                media_index = Math.max(0, Math.min(media_index, unsequenced_medias.length - 1));

                let old_focus_index = sct_focus_index;
                sct_focus_index = media_index;

                if (auto_select_mode) {
                    autoSelectHandler(old_focus_index, sct_focus_index);
                }
            }

            /**
             * handles the go to bottom hotkey.
             * @param {KeyboardEvent} key_event
            */
            const handleGoToBottom = (key_event) => {
                sct_focus_index = unsequenced_medias.length - 1;
            }

            /**
             * Handles the go to top hotkey.
             * @param {KeyboardEvent} key_event
             */
            const handleGoToTop = (key_event) => {
                sct_focus_index = 0;
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

            /**
             * Inserts the yanked medias before the focused media. Called by a hotkey.
             * @type {import('@libs/LiberyHotkeys/hotkeys').HotkeyCallback}
             */
            const handleInsertYankedMediasHotkey = (key_event, hotkey) => {
                if ($me_gallery_yanked_medias.length === 0) return;

                if (isCurrentMediaYanked()) {
                    let user_error = new LabeledError("Unstable operation forbidden", "You attempted to insert before a media that would be moved by the insertion of the medias, this operation could have unpredictable results.", lf_errors.PROGRAMMING_ERROR__BROKEN_STATE);

                    user_error.alert();
                    return;
                }

                last_saved_sequence = unsequenced_medias;

                insertMediasAt($me_gallery_yanked_medias, sct_focus_index);

                me_gallery_yanked_medias.set([]);

                sct_focus_index = Math.min(sct_focus_index, unsequenced_medias.length - 1);
            } 

            /**
             * Restores the last saved sequence.
             */
            const handleRestoreLastSequence = () => {
                unsequenced_medias = last_saved_sequence;
                last_saved_sequence = [];
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
             * Saves the sequence.
             */
            const handleSaveSequence = async () => {
                if ($current_category == null) {
                    console.error("In SequenceCreationTool.handleSaveSequence, current category is not available.");
                    return;
                }
                
                // Show the would be parameters to test the api

                let sequence_map = createSequenceMap(unsequenced_medias, sequence_prefix);

                let renamed = await sequenceRenameMedias(sequence_map, $current_category.uuid);

                console.log("Renamed: ", renamed);
            }

            /**
             * Toggles the magnify mode.
             */
            const handleToggleMagnifyMode = () => {
                toggleMagnifyMode();
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
        }

        /**
         * Creates a sequence map from the given medias array and name prefix.
         * @param {import('@models/Medias').OrderedMedia[]} medias
         * @param {string} prefix
         * @returns {SequenceMap}
         * @typedef {Object<string, string>} SequenceMap - A map of the media uuids to the new sequence names without the extension.
         */
        const createSequenceMap = (medias, prefix) => {
                /**
                 * @type {Object<string, string>}
                 */
                let sequence_map = {};

                const zero_padding = Math.max(2, Math.ceil(Math.log10(medias.length + 1)));

                for (let h = 0; h < medias.length; h++) {
                    let media = medias[h];

                    let index_string = (h + 1).toString().padStart(zero_padding, "0");

                    let new_name = `${prefix}${index_string}`;

                    sequence_map[media.uuid] = new_name;
                }

                return sequence_map;
        }

        /**
         * Emits an event so the MediaExplorerGallery can close the
         * sequence creation tool.
         */
        const emitCloseGallery = () => {
            dispatch("close-sct");
        }

        /**
         * Returns the currently focused media.
         * @returns {import('@models/Medias').OrderedMedia}
         */
        const getFocusedMedia = () => {
            return unsequenced_medias[sct_focus_index];
        }

        /**
         * Returns a media on a given index.
         * @param {number} index
         * @returns {import('@models/Medias').OrderedMedia}
         */
        const getMediaAtIndex = index => {
            return unsequenced_medias[index];
        }

        /**
         * Generates an initial sequence prefix based on the current category name
         */
        const generateInitialSequencePrefix = () => {
            if ($current_category == null) {
                console.error("In SequenceCreationTool.generateInitialSequencePrefix, current category is not available.");
                return;
            }
            
            let new_prefix = $current_category.name + "_";

            new_prefix = cleanFilenameString(new_prefix);

            if (new_prefix === "") {
                new_prefix = "sequence_";
            }

            sequence_prefix = new_prefix;
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
         * Returns whether the current media is yanked.
         * @returns {boolean}
         */
        const isCurrentMediaYanked = () => {
            return isMediaYanked(sct_focus_index);
        }

        /**
         * Inserts given medias array at before the given index.
         * @param {import('@models/Medias').OrderedMedia[]} medias
         * @param {number} insert_point
         */
        const insertMediasAt = (medias, insert_point) => {

            let duplicate_medias_lookup = new Set(medias.map(media => media.uuid));

            const current_focused_media = getMediaAtIndex(insert_point);

            let filtered_sequence = unsequenced_medias.filter(media => !duplicate_medias_lookup.has(media.uuid));

            insert_point = filtered_sequence.indexOf(current_focused_media);

            let new_sequence_left = filtered_sequence.slice(0, insert_point);
            let new_sequence_right = filtered_sequence.slice(insert_point);

            unsequenced_medias = [...new_sequence_left, ...medias, ...new_sequence_right];
        }

        /**
         * Toggles the magnify mode.
         */
        const toggleMagnifyMode = () => {
            magnify_mode = !magnify_mode;
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
        class="dungeon-scroll"
        class:magnify-mode-enabled={magnify_mode}
        style:--sct-grid-item-size="{(100 / medias_per_row) * 0.94}cqw"
    >
        <div id="sequence-parameters"
            class="cristal-surface"
        >
            <label class="dungeon-input">
                <span class="dungeon-label">
                    Sequence prefix
                </span>
                <input 
                    id="sequence-prefix-editor"
                    bind:value={sequence_prefix}
                    bind:this={the_sequence_prefix_editor}
                    type="text"
                    minlength="1"
                    maxlength="50"
                    spellcheck="true"
                    autocomplete="on"
                    pattern="{'^[a-zA-Z0-9_\\- ]+$'}"
                    required
                >
            </label>
            <ul id="categorie-medias-properties"
                class="dungeon-properties"
            >
                <p class="sct-property dungeon-property">
                    <strong class="dungeons-field-label">Medias</strong> {unsequenced_medias.length}
                </p>
                <p class="sct-property dungeon-property">
                    <strong class="dungeons-field-label">Selected</strong> {$me_gallery_yanked_medias.length}
                </p>
            </ul>
        </div>
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
                        enable_magnify_on_keyboard_focus={magnify_mode}
                        container_selector="#sequence-members"
                        check_container_limits
                        enable_video_titles
                        enable_dragging
                    />
                </li>
            {/each}
        </ul>
    </div>
{/if}

<style>

    #sequence-creation-tool {
        position: relative;
        width: 100%;
        height: calc(90dvh - var(--navbar-height));
        overflow: auto;
    }
    
    /*=============================================
    =            Parameters            =
    =============================================*/
    
        #sequence-parameters {
            position: sticky;
            display: flex;
            width: 100%;
            height: var(--primary-toolbar-height);
            top: 0;
            align-items: center;
            column-gap: var(--spacing-4);
            border-bottom: .5px solid var(--grey-9);
            z-index: var(--z-index-t-1);
            will-change: transform; 
            /* Will change - transform is said to improve performance on stick elements inside scrolling containers
              See: https://developer.mozilla.org/en-US/docs/Web/CSS/position#performance_accessibility
            */ 
        }

        #sequence-creation-tool.magnify-mode-enabled #sequence-parameters {
            opacity: 0;
            transition: opacity 0.2s ease-out;
        }

        label.dungeon-input input {
            color: var(--main);
        }

        ul#categorie-medias-properties {
            & p.sct-property {
                font-size: var(--font-size-1);
            }
        }
    
    /*=====  End of Parameters  ======*/

    ul#sequence-members {
        display: grid;
        width: 100%;
        container-type: inline-size;
        background: var(--grey);
        row-gap: var(--spacing-2);
        padding: 4px;
        list-style: none;
        margin: 0;
    }
    
    /*=============================================
    =            Sequence members            =
    =============================================*/
    
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
    
    /*=====  End of Sequence members  ======*/
</style>