<script>
    import { getTrashcanTransaction, restoreMedia } from '@models/Trashcan';
    import TrashcanTransactionComponent from './TrashcanTransaction.svelte';
    import { TrashcanTransaction } from '@models/Trashcan';
    import global_hotkeys_manager from '@libs/LiberyHotkeys/libery_hotkeys';
    import { tmt_hotkeys_context_name } from '../tmt_state';
    import { onDestroy, onMount } from 'svelte';
    import { confirmPlatformMessage } from '@libs/LiberyFeedback/lf_utils';
    import { createEventDispatcher } from 'svelte';
    import { current_category } from '@stores/categories_tree';
    import { current_user_identity } from '@stores/user';
    
    /*=============================================
    =            Properties            =
    =============================================*/
        
        /*=============================================
        =            Hotkeys            =
        =============================================*/
        
            /**
             * Binds to add to the Transactions Management Tool hotkeys context.
             * @type {Object<string, import('@libs/LiberyHotkeys/hotkeys').HotkeyDataParams>}
             */ 
            const keybinds = {
                ENTRY_FOCUS_MOVEMENT: {
                    key_combo: ["w", "s"],
                    description: "<trashcan_navigation> Move the keyboard focus up and down the transaction entries",
                    handler: handleKeyboardMovement
                },
                GO_TO_OPENED_TRANSACTION: {
                    key_combo: "a",
                    description: "<trashcan_navigation> If focus is in a transaction's content, moves the focus back to the transaction entry",
                    handler: handleFocusOpenedTransaction
                },
                ENTRY_SELECTION: {
                    key_combo: ["e", "d"],
                    description: "<trashcan_navigation> Select the transaction entry",
                    handler: handleKeyboardSelection
                }
            }
        
        /*=====  End of Hotkeys  ======*/
    
        /** 
         * Trashcan entries that have affected medias.
         * @type {import('@models/Trashcan').TrashcanTransactionEntry[]}
         */
        export let trashcan_transaction_entries = [];

        /**
         * A map of entries timestamp to their respective transaction details.  
         * @type {Map<string, import('@models/Trashcan').TrashcanTransaction}
        */
        const loaded_trashcan_transactions = new Map();


        /**
         * Keyboard focused index.
         * @type {number}
         */
        let entry_keyboard_focused_index = 0;

        /**
         * The index of the media entry that is currently keyboard focused.
         * @type {number}
         */
        let media_keyboard_focused_index = -1;

        /**
         * Whether the keyboard movement is focused on the transaction entries(false) or the media entries(true).
         * @type {boolean}
         */
        let media_movement_enabled = false;

        /**
         * The last trashcan entry index that was selected.
         * @type {number}
         */
        let selected_trashcan_entry_index = -1;


        /**
         * Current selected trashcan transaction details.
         * @type {import('@models/Trashcan').TrashcanTransaction}
         */
        let selected_trashcan_transaction = null;

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        addPrivilegedHotkeys();

        registerTrashCanHotkeys();
    });

    onDestroy(() => {
        removeTrashcanHotkeys();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
        
        /*=============================================
        =            Keybinds            =
        =============================================*/

            /**
             * Registers privileged keybinds to the keybinds object, must be called before registering hotkeys on the global hotkeys manager as it
             * only modifies the keybinds object. The access decisions are made based on the current_user_identity grants.
             * @requires keybinds
             * @requires global_hotkeys_manager
             * @requires current_user_identity
             */
            const addPrivilegedHotkeys = () => {
                if (current_user_identity == null) return;

                if ($current_user_identity.canModifyTrashcan()) {
                    keybinds.RESTORE_FOCUSED_TRANSACTION = {
                        key_combo: "r r",
                        description: "<trashcan_content> Restores all medias in the focused transaction to the current category",
                        handler: handleRestoreFocusedTransaction
                    }
                }

                if ($current_user_identity.canEmptyTrashcan()) {
                    keybinds.DELETE_FOCUSED_TRANSACTION = {
                        key_combo: "x x",
                        description: "<trashcan_content> Permanently deletes all medias in the focused transaction",
                        handler: handleDeleteFocusedTransaction
                    }
                }
            }


            /**
             * Ensure keyboard focus is visible in it's scrolling context.
             * @param {number} new_keyboard_focused_index
             */
            const ensureKeyboardEntryFocusVisible = (new_keyboard_focused_index) => {
                let focused_entry = getFocusedEntryElement(new_keyboard_focused_index);

                if (focused_entry == null) return;

                focused_entry.scrollIntoView({
                    block: "nearest",
                    inline: "nearest"
                });
            }

            /**
             * Returns the keyboard focused entry's dom element.
             * @param {number} new_keyboard_focused_index   
             * @returns {HTMLElement | null}
             */
            const getFocusedEntryElement = (new_keyboard_focused_index) => document.querySelector(`.tmt-te-transaction-btn[data-entry-index="${new_keyboard_focused_index}"]`);

            /**
             * Hanldes the movement of the keyboard focus transaction entry.
             * @param {KeyboardEvent} event 
             * @param {string} key_combo
             */ 
            function handleKeyboardMovement(event, key_combo) {
                let is_direction_up = key_combo === "w";
                if (media_movement_enabled) {
                    return passKeyboardMovementToTransactionContent(is_direction_up);
                }
                
                let new_keyboard_focused_index = is_direction_up ? Math.max(0, entry_keyboard_focused_index - 1) : Math.min(trashcan_transaction_entries.length - 1, entry_keyboard_focused_index + 1);
                
                if ((isEntrySelected(entry_keyboard_focused_index) && !is_direction_up) || (isEntrySelected(new_keyboard_focused_index) && is_direction_up)) {
                    enableMediaKeyboardMovement(is_direction_up);
                    return; 
                }

                entry_keyboard_focused_index = new_keyboard_focused_index;  
                
                ensureKeyboardEntryFocusVisible(entry_keyboard_focused_index);
            }

            /**
             * If confirmed, deletes the currently selected transaction from the system.
             */
            async function handleDeleteFocusedTransaction() {
                let focused_transaction = trashcan_transaction_entries[entry_keyboard_focused_index];

                let user_choice = await confirmPlatformMessage({
                    message_title: "Delete all files in transaction",
                    question_message: `This action is irreversible. Are you sure you want to delete all ${focused_transaction.affected_medias} files in this transaction? Again, you WILL NOT BE ABLE TO UNDO THIS!`,
                    auto_focus_cancel: true,
                    cancel_label: "Keep",
                    confirm_label: "Delete FOREVER",
                    danger_level: 1,
                });

                if (isEntrySelected(entry_keyboard_focused_index)) {
                    focusOpenedTransaction();
                    closeSelectedTransaction(); 
                    entry_keyboard_focused_index = Math.max(0, Math.min(trashcan_transaction_entries.length - 1, entry_keyboard_focused_index));
                }

                if (user_choice === 1) {
                    dispatch("delete-transaction", focused_transaction);
                }
                
            }

            /**
             * Selects a transaction if media_movement_enabled false, otherwise selects a media entry.
             */
            function handleKeyboardSelection() {
                if (media_movement_enabled) {
                    handleMediaKeyboardSelection();
                } else {
                    handleKeyboardEntrySelection();
                }
            }

            /**
             * Handles keyboard selection of a transaction entry.
             * @requires selected_trashcan_entry_index
             * @requires selected_trashcan_transaction
             * @requires keyboard_focused_index
             * @requires trashcan_transaction_entries
            */
            async function handleKeyboardEntrySelection() {
                if (entry_keyboard_focused_index === selected_trashcan_entry_index) {
                    selected_trashcan_entry_index = -1;
                    selected_trashcan_transaction = null;
                    return;
                }

                if (entry_keyboard_focused_index < 0 || entry_keyboard_focused_index >= trashcan_transaction_entries.length) return;

                let transaction_id = trashcan_transaction_entries[entry_keyboard_focused_index].timestamp;
                
                selected_trashcan_transaction = await selectTrashcanTransaction(transaction_id);
                selected_trashcan_entry_index = entry_keyboard_focused_index;

                enableMediaKeyboardMovement();
            }

            /**
             * Closes the currently selected transaction.
             */
            function handleFocusOpenedTransaction() {
                focusOpenedTransaction();       
            }

            /**
             * Restores the currently focused transaction.
             */
            async function handleRestoreFocusedTransaction() {
                let focused_transaction = trashcan_transaction_entries[entry_keyboard_focused_index];   

                let user_choice = await confirmPlatformMessage({
                    message_title: `Restore all files in transaction to '${$current_category.name}'`,
                    question_message: `Are you sure you want to restore all ${focused_transaction.affected_medias} files in this transaction?`,
                    auto_focus_cancel: focused_transaction.affected_medias > 10,
                    cancel_label: "Cancel",
                    confirm_label: "Restore",
                    danger_level: 0,
                });

                if (user_choice !== 1) return;

                if (isEntrySelected(entry_keyboard_focused_index)) {
                    focusOpenedTransaction();
                    closeSelectedTransaction();
                    entry_keyboard_focused_index = Math.max(0, Math.min(trashcan_transaction_entries.length - 1, entry_keyboard_focused_index));
                }

                dispatch("restore-transaction", focused_transaction);
            }

            const focusOpenedTransaction = () => {
                media_keyboard_focused_index = 0;
                passKeyboardMovementToTransactionContent(true);
                setKeyboardFocusAtTop();
            }

            /**
             * 
            */
            async function handleMediaKeyboardSelection() {
                const media_focus_out_of_bounds = media_keyboard_focused_index < 0 || media_keyboard_focused_index >= selected_trashcan_transaction.Content.length; 

                if (!validTransactionOpened() || !media_movement_enabled || media_focus_out_of_bounds) return;

                let focused_media = selected_trashcan_transaction.Content[media_keyboard_focused_index];

                let user_choice = await confirmPlatformMessage({
                    message_title: `Restore media on ${$current_category.name}`,
                    question_message: `Load '${focused_media.Name}' on the current_category`,
                    auto_focus_cancel: false,
                    cancel_label: "No", 
                    confirm_label: "Yes",
                    danger_level: 0,
                });

                if (user_choice === 1) {
                    restoreMediaToCurrentCategory(focused_media);
                }   
            }

            /**
             * Called from handleKeyboardMovement to 
             * @param {boolan} is_direction_up
             */
            const passKeyboardMovementToTransactionContent = (is_direction_up) => {
                let new_media_index = is_direction_up ? media_keyboard_focused_index - 1 : media_keyboard_focused_index + 1;

                if (new_media_index < 0) {
                    entry_keyboard_focused_index = selected_trashcan_entry_index;
                    media_movement_enabled = false;
                    media_keyboard_focused_index = -1;  
                    return;
                } else if (new_media_index >= selected_trashcan_transaction.Content.length) {
                    entry_keyboard_focused_index = Math.min(trashcan_transaction_entries.length - 1, selected_trashcan_entry_index + 1);
                    media_movement_enabled = false;
                }

                media_keyboard_focused_index = new_media_index;
            }

            /**
             * Registers the TrashcanTab hotkeys to the Transactions Management Tool hotkeys context.
             * @requires global_hotkeys_manager
             * @requires tmt_hotkeys_context_name
             */
            const registerTrashCanHotkeys = () => {
                if (global_hotkeys_manager.ContextName !== tmt_hotkeys_context_name) return;

                for (let keybind of Object.values(keybinds)) {
                    global_hotkeys_manager.registerHotkeyOnContext(keybind.key_combo, keybind.handler, {
                        description: keybind.description
                    })
                }
            }

            /**
             * Removes the TrashcanTab hotkeys from the Transactions Management Tool hotkeys context.
             */
            const removeTrashcanHotkeys = () => {
                if (global_hotkeys_manager.ContextName !== tmt_hotkeys_context_name) return;

                for (let keybind of Object.values(keybinds)) {
                    global_hotkeys_manager.unregisterHotkeyFromContext(keybind.key_combo, keybind.mode ?? "keydown");
                }
            }

            /**
             * Sets the Keyboard focused index at the top of it's scrolling context.
             */
            const setKeyboardFocusAtTop = () => {
                let focused_entry = getFocusedEntryElement(entry_keyboard_focused_index);

                if (focused_entry == null) return;

                focused_entry.scrollIntoView({
                    behavior: "smooth",
                    block: "start",
                });                    
            }
                
            /**
             * Enables media keyboard movement either at the bottom or top of the transaction content.
             * @param {boolean} at_end
             */
            const enableMediaKeyboardMovement = (at_end=false) => {
                if(!at_end) setKeyboardFocusAtTop();
                media_movement_enabled = true;
                media_keyboard_focused_index = at_end ? selected_trashcan_transaction.Content.length - 1 : 0;
            }
            
        
        /*=====  End of Keybinds  ======*/

        /**
         * Returns whether a given index represents a selected transaction entry.
         * @param {number} check_index
         * @returns {boolean
         */
        const isEntrySelected = (check_index) => {
            return selected_trashcan_entry_index === check_index && selected_trashcan_transaction != null;
        }

        /**
         * Closes the currently selected transaction.
         */
        const closeSelectedTransaction = () => {
            selected_trashcan_entry_index = -1;
            selected_trashcan_transaction = null;
        }
    
        /**
         * Handles the click event on a transaction entry button.   
         * @param {MouseEvent} event
         * @requires selected_trashcan_entry_index
         */ 
        const handleTransactionEntryClick = async (event) => {
            let entry_index = event.target?.dataset?.entryIndex;

            if (entry_index == null) return;


            entry_index = parseInt(entry_index);
            console.log("Entry index: ", entry_index);
            
            if (isEntrySelected(entry_index)) {
                return closeSelectedTransaction()
            }

            if (entry_index < 0 || entry_index >= trashcan_transaction_entries.length) return;

            let transaction_id = trashcan_transaction_entries[entry_index].timestamp;
            
            selected_trashcan_transaction = await selectTrashcanTransaction(transaction_id);
            selected_trashcan_entry_index = entry_index;
            console.log("Selected trashcan transaction: ", selected_trashcan_transaction);
        }

        /**
         * Restores a given media to the current category.
         * @param {import('@models/Trashcan').TrashedMedia} focused_media
         */
        const restoreMediaToCurrentCategory = async (focused_media) => {

            let transaction_copy = selected_trashcan_transaction.clone();

            transaction_copy.removeMedia(focused_media.UUID);

            dispatch("restore-media", {
                media: focused_media,
                transaction: selected_trashcan_transaction
            });

            let restored = await restoreMedia(focused_media.UUID, transaction_copy.TransactionID, $current_category.uuid);

            if (!restored) return;

            if (transaction_copy.Content.length === 0) {
                closeSelectedTransaction();
                media_movement_enabled = false;
                media_keyboard_focused_index = -1;
                entry_keyboard_focused_index = Math.max(0, Math.min(trashcan_transaction_entries.length - 1, selected_trashcan_entry_index));
            }
            
            selected_trashcan_transaction = transaction_copy;
        }

        /**
         * Selects a new trashcan transaction to be displayed.
         * @param {string} transaction_id
         * @returns {import('@models/Trashcan').TrashcanTransaction}
         */
        const selectTrashcanTransaction = async (transaction_id) => {
            if (loaded_trashcan_transactions.has(transaction_id)) {
                return loaded_trashcan_transactions.get(transaction_id);
            }

            let new_transaction = await getTrashcanTransaction(transaction_id);

            if (new_transaction == null) return null;

            loaded_trashcan_transactions.set(transaction_id, new_transaction);

            return new_transaction;
        }

        /**
         * Whether or not there is a transaction opened.
         * @returns {boolean}
         */
        const validTransactionOpened = () => {
            return (selected_trashcan_transaction instanceof TrashcanTransaction) && selected_trashcan_entry_index >= 0 && selected_trashcan_entry_index < trashcan_transaction_entries.length;
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<menu id="tmt-transaction-type-tab-entries">
    {#each trashcan_transaction_entries as entry, h}
        {@const selected_transaction = selected_trashcan_entry_index === h && selected_trashcan_transaction != null}
        {@const is_keyboard_focused = entry_keyboard_focused_index === h}
        <li class="tmt-transaction-entry" 
            class:entires-out-of-focus={media_movement_enabled}
        >
            <button class="tmt-te-transaction-btn"
                class:is-selected-transaction={selected_transaction}
                class:is-keyboard-focused={is_keyboard_focused}
                on:click={handleTransactionEntryClick} 
                data-entry-index={h}
            >
                <span class="transaction-time">
                    {entry.timestamp}
                </span>
                <span class="transaction-deleted-medias">
                    medias deleted: <i>{entry.affected_medias}</i>
                </span>
            </button>
            {#if selected_transaction}
                <div class="tmt-transaction-selected-details">
                    <TrashcanTransactionComponent 
                        the_transaction={selected_trashcan_transaction}
                        focused_media_index={media_keyboard_focused_index}
                        highlight_focused_media={media_movement_enabled}
                    />
                </div>
            {/if}   
        </li>
    {/each}
</menu>
    
<style>
    #tmt-transaction-type-tab-entries {
        display: flex;
        width: 100%;
        height: 100%; 
        background: var(--grey);
        flex-direction: column;
        border-radius: calc(var(--border-radius) * 0.8);
        border: 1px solid var(--main-dark-color-9); 
        overflow-y: auto;
        scrollbar-width: thin;
        scrollbar-color: var(--main) var(--grey-9);

    }

    
    /*=============================================
    =            Entries            =
    =============================================*/
    
        .tmt-transaction-entry {
            background-color: var(--grey-9);

            &:not(:last-child) {
                border-bottom: 1px solid var(--grey-7);
            }
        } 

        button.tmt-te-transaction-btn {
            width: 100%;
            height: 50px;
            display: flex;
            align-items: center;
            padding: 0 var(--spacing-4);
            justify-content: space-between;
            transition: background .4s ease-out;
            border-radius: 0;
            outline: none;

            &:hover {
                background: var(--grey-8);
            }   

            & .transaction-time {
                font-size: var(--font-size-2);
                color: var(--grey-4);
            }

            & .transaction-deleted-medias i {
                color: var(--danger-7); 
                padding-inline: var(--spacing-1);
            }

            & * {
                pointer-events: none;
            }
        }

        button.tmt-te-transaction-btn.is-keyboard-focused {
                background: hsl(from var(--main-9) h s l / 0.6);
        }

        .entires-out-of-focus button.tmt-te-transaction-btn.is-keyboard-focused {
            background: hsl(from var(--main-9) h calc(s * 0.5) l / 0.2);
        }
    
    /*=====  End of Entries  ======*/
    
</style>