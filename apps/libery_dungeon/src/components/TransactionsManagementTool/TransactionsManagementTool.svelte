<script>
    import { tmt_hotkeys_context_name } from "./tmt_state";
    import { deleteTransaction, emptyTrashcan, getTrashcanEntries, restoreTransaction } from "@models/Trashcan";
    import { createEventDispatcher, onMount } from "svelte";
    import TrashcanTab from "./sub-components/TrashcanTab.svelte";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
    import { hotkeys_sheet_visible } from "@stores/layout";
    import { current_category } from "@stores/categories_tree";
    import { confirmPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import { current_user_identity } from "@stores/user";

    
    /*=============================================
    =            Properties            =
    =============================================*/
        
        /*=============================================
        =            Hotkeys            =
        =============================================*/
    
            /**
             * @type {Object<string, import('@libs/LiberyHotkeys/hotkeys').HotkeyDataParams>}
             */
            const keybinds ={
                CLOSE_TMT_TOOL: {
                    key_combo: ["q", "="],
                    handler: handleCloseTMTTool,
                    options: {
                        description: `<${HOTKEYS_GENERAL_GROUP}> Close the Transactions Management Tool`,
                    }
                },
                OPEN_HOTKEYS_SHEET: {
                    key_combo: "?",
                    handler: () => {
                        hotkeys_sheet_visible.set(!$hotkeys_sheet_visible)
                    },
                    options: {
                        description: `<${HOTKEYS_GENERAL_GROUP}> Toggle the hotkeys sheet`,
                    }
                }
            }
        
        /*=====  End of Hotkeys  ======*/
    
        /**
         * All the trashcan entries returned from the server. Only returns those that had affected medias.
         * @type {import('@models/Trashcan').TrashcanTransactionEntry[]}
         */
        let trashcan_transaction_entries = [];

        let global_hotkeys_manager = getHotkeysManager();

        /**
         * Whether the tool dialog is open or not.
         * @type {boolean}
         */
        export let transaction_management_tool_open = false;


        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        addPrivilegedHotkeys();

        registerTMTToolHotkeysContext();

         // Retrieving entries at this point forces svelte mount TransactionsManagementTool before TrashcanTab(or the other tabs in the future) which assures that the hotkey context will be set correctly before 
            // those components attempt to bind extra hotkeys to it.
        retrieveTrashcanEntries();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
        
        /*=============================================
        =            Keybinds            =
        =============================================*/
        
            /**
             * Fills up privileged hotkeys based on the current user's grants. Most be called before the hotkeys context is registered.
             * as it only changes the keybinds object.
             * @requires global_hotkeys_manager
             * @requires keybinds
             */ 
            const addPrivilegedHotkeys = () => {
                if ($current_user_identity.canEmptyTrashcan()) {
                    keybinds.EMPTY_TRASHCAN = {
                        key_combo: ["x a"],
                        handler: handleTrashcanEmpty,
                        options: {
                            description: "<trashcan_content> Deletes all files in the trashcan and closes the tool",
                            await_execution: false,
                        }
                    }
                }
            }
                
        
            /**
             * Closes the Transactions Management Tool dialog.
             */
            function handleCloseTMTTool() {
                closeTMTTool();
            } 

            /**
             * Register the transactions management tool hotkeys context.
             * @requires global_hotkeys_manager
             * @requires tmt_hotkeys_context_name
             * @requires keybinds
             */
            const registerTMTToolHotkeysContext = () => {
                if (!global_hotkeys_manager.hasContext(tmt_hotkeys_context_name)) {
                    const hotkeys_context = new HotkeysContext();

                    for (let keybind of Object.values(keybinds)) {
                        hotkeys_context.register(keybind.key_combo, keybind.handler, keybind.options);

                    }

                    global_hotkeys_manager.declareContext(tmt_hotkeys_context_name, hotkeys_context);
                }

                global_hotkeys_manager.loadContext(tmt_hotkeys_context_name, true);                
            }
        
        /*=====  End of Keybinds  ======*/
    
        
        /*=============================================
        =            Trashcan operations            =
        =============================================*/
        
            /**
             * Retrieves all the trashcan entries and sets them to the trashcan_transaction_entries property.
             * @requires getTrashcanEntries
             * @requires trashcan_transaction_entries
             */
            const retrieveTrashcanEntries = async () => {
                let new_trashcan_entries = await getTrashcanEntries();
                
                if (new_trashcan_entries == null || new_trashcan_entries.length === 0) return;
                
                console.log("New trashcan entries: ", new_trashcan_entries);    
                trashcan_transaction_entries = new_trashcan_entries;
            }

            /**
             * Handles the delete-transaction event from the TrashcanTab component.
             * @param {CustomEvent<import('@models/Trashcan').TrashcanTransactionEntry>} event
             */ 
            const handleDeleteTrashcanTransaction = async (event) => {
                let transaction_entry = event.detail;

                await deleteTrashcanTransactionEntry(transaction_entry);
            }

            /**
             * Deletes all the transactions in the trashcan and it's files.
             */
            async function handleTrashcanEmpty() {
                let user_choice = await confirmPlatformMessage({
                    message_title: "Delete all files in the trashcan",
                    question_message: "This will irreversibly delete all the files in that are currently in the trashcan. Are you sure you want to proceed?",
                    confirm_label: "Delete my files",
                    cancel_label: "Cancel",
                    auto_focus_cancel: true,
                    danger_level: 2
                }); 

                if (user_choice !== 1) return;  

                emptyTrashcan();

                trashcan_transaction_entries = [];  

                closeTMTTool();
            }


            /**
             * deletes a single transaction.
             * @param {import('@models/Trashcan').TrashcanTransactionEntry} transaction_entry
             */
            const deleteTrashcanTransactionEntry = async (transaction_entry) => {
                console.log("Deleting transaction: ", transaction_entry);

                let filtered_entries = trashcan_transaction_entries.filter(entry => entry.timestamp !== transaction_entry.timestamp);

                let deleted = await deleteTransaction(transaction_entry.timestamp);
                console.log("Deleted: ", deleted);

                if (!deleted) return;

                trashcan_transaction_entries = filtered_entries;
            }

            /**
             * Handles the restore-transaction event from the TrashcanTab component.
             * @param {CustomEvent<import('@models/Trashcan').TrashcanTransactionEntry>} event
             */
            const handleRestoreTrashcanTransaction = async (event) => {
                let transaction_entry = event.detail;

                restoreTrashcanTransactionEntry(transaction_entry);                
            }

            /**
             * Restores a single transaction. into the current category.
             * @param {import('@models/Trashcan').TrashcanTransactionEntry} transaction_entry
             * @requires trashcan_transaction_entries
             * @requires current_category
             */
            const restoreTrashcanTransactionEntry = async (transaction_entry) => {
                console.log("Restoring transaction: ", transaction_entry);

                let restored = await restoreTransaction(transaction_entry.timestamp, $current_category.uuid);
                console.log("Restored: ", restored);

                if (!restored) return;

                trashcan_transaction_entries = trashcan_transaction_entries.filter(entry => entry.timestamp !== transaction_entry.timestamp);
            }

            /**
             * Handles the restore-media event from the TrashcanTransaction component. the TrashcanTab handles the restoration, but if the transaction had only that media, we remove it from the ui.
             * @param {CustomEvent<RestoredMediaEventDetail>} event
             * @typedef {Object} RestoredMediaEventDetail
             * @property {import('@models/Trashcan').TrashedMedia} media
             * @property {import('@models/Trashcan').TrashcanTransaction} transaction
             */
            const handleRestoreMedia = async (event) => {
                let media = event.detail.media; 
                let transaction = event.detail.transaction;

                if (transaction.Content.length > 1) return;

                let filtered_transaction_entries = trashcan_transaction_entries.filter(entry => entry.timestamp !== transaction.TransactionID);

                trashcan_transaction_entries = filtered_transaction_entries;    
            }
        
        /*=====  End of Trashcan  ======*/

        const closeTMTTool = () => {
            global_hotkeys_manager.loadPreviousContext();
            global_hotkeys_manager.dropContext(tmt_hotkeys_context_name);

            dispatch("tool-closed");
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<div id="transaction-management-tool-wrapper">
    <dialog id="transaction-management-tool"
        class="libery-dungeon-window"
        open={transaction_management_tool_open}
    >
        <header id="tmt-header">
            <h3 id="tmt-header-title">
                Transactions Management Tool    
            </h3>
            <p id="tmt-header-instructions">
                Here you can view the medias deletion transactions that have not been deleted yet. In future releases, you will be able to manage Move and Download transactions as well.
            </p>
        </header>
        <div id="tmt-transactions-wrapper">
            <menu id="transaction-type-tabs">
                <li class="transaction-type-tab">
                    <button id="trashcan-transaction-tab-selector">
                        Deleted
                    </button>
                </li>
            </menu>
            <section id="tmt-transaction-list-wrapper">
                {#if trashcan_transaction_entries.length > 0 && transaction_management_tool_open}
                    <TrashcanTab
                        trashcan_transaction_entries={trashcan_transaction_entries}
                        on:delete-transaction={handleDeleteTrashcanTransaction} 
                        on:restore-transaction={handleRestoreTrashcanTransaction}   
                        on:restore-media={handleRestoreMedia}
                    />
                {/if}
            </section>
        </div>
    </dialog>
</div>

<style>
    #transaction-management-tool-wrapper {
        width: 100%;
        height: 100%;
        display: grid;
        container-type: size;
        place-items: center;
        backdrop-filter: var(--backdrop-filer-blur); 
        background: hsl(from var(--grey) h s l / 0.5);
    }

    dialog#transaction-management-tool {
        display: flex;
        width: min(100cqw, 800px);
        height:90cqh;
        flex-direction: column;
        padding: calc(var(--spacing-3) + var(--spacing-2)) var(--spacing-1);
        color: var(--grey-1);
        row-gap: var(--spacing-3);
        outline: none;
    }

    header#tmt-header {
        display: flex;  
        padding: 0 var(--spacing-2);
        flex-direction: column;
        row-gap: var(--spacing-2);

        & h3#tmt-header-title {
            text-align: center;
            color: var(--main);
            font-size: var(--font-size-2);
            line-height: 1;
        }

        & p#tmt-header-instructions {
            font-size: var(--font-size-1);
            color: var(--grey-2);
        }
    }

    div#tmt-transactions-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: var(--spacing-2);
    }       
    /*=============================================
    =            Tabs            =
    =============================================*/
    
        menu#transaction-type-tabs {
            display: flex;
            height: max(4cqh, 10px);
            background: var(--grey);
            align-items: center;
        }

        .transaction-type-tab {
            height: 100%;   
        }

        button#trashcan-transaction-tab-selector {
            --trashcan-tab-color: color-mix(in hsl, var(--danger), var(--main-dark) 60%);

            background: var(--trashcan-tab-color);
            color: var(--grey);
            font-weight: 450;
            padding: 0 var(--spacing-2);
            height: 100%;
            border-radius: 0;
        }

    /*=====  End of Tabs  ======*/

    
    /*=============================================
    =            transactions section            =
    =============================================*/
    
        section#tmt-transaction-list-wrapper {
            padding: 0 var(--spacing-3);
            height: 50cqh;
            container-type: size;
        }
        
    /*=====  End of transactions section  ======*/
    
    
    
    


</style>