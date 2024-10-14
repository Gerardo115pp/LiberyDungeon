<script>
    /*=============================================
    =            Imports            =
    =============================================*/
    
        import CategoriesClustersItem from "@components/CategoriesClusters/CategoriesClustersItem.svelte";
        import ClusterCreationTool from "./sub-components/ClusterCreationTool/ClusterCreationTool.svelte";
        import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
        import { getAllCategoriesClusters } from "@models/CategoriesClusters";
        import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
        import { toggleHotkeysSheet } from "@stores/layout";
        import { onDestroy, onMount } from "svelte";
        import { HOTKEYS_GENERAL_GROUP } from "@libs/LiberyHotkeys/hotkeys_consts";
        import { browser } from "$app/environment";
        import DungeonSelectionHint from "./sub-components/DungeonSelectionHint.svelte";
        import { confirmPlatformMessage, emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
        import { deleteClusterRecord } from "@models/CategoriesClusters";
        import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
        import { lf_errors } from "@libs/LiberyFeedback/lf_errors";
    
    /*=====  End of Imports  ======*/
    
    /*=============================================
    =            Properties            =
    =============================================*/

        /**
         * The hotkeys context name of the dungeons selector app.
         * @type {string}
        */
        const hotkeys_context_name = "dungeons-selector";
    
        /**
         * The categories clusters of the system.
         * @type {import("@models/CategoriesClusters").CategoriesCluster[]}
         */ 
        let categories_clusters = [];

        let global_hotkeys_manager = getHotkeysManager();

        /**
         * Whether the cluster creation tool is visible or not.
         * @type {boolean}
         * @default false
         */
        let show_cluster_creation_tool = false; 

        /**
         * The amount of cluster to render per row.
         * @type {number}
        */
        const cluster_per_row = 3

        
        /*----------  State  ----------*/
        
            /**
             * The keyboard selected cluster index.
             * @type {number}
             */
            let focused_cluster_index = 0;
        
        /*----------  Phrases to use on the headlines  ----------*/
        
            const gaming_phrases = {
                headlines: [
                    "Blessings of the Exalted Father",
                    "Praise the Sun!",
                    "The cake is a lie",
                    "The purity of light",
                    "The darkness of the abyss",
                    "Divest yourself of everything",
                    "A Frenzied Flame to burn away the curses",
                    "Burn all that divides and distinguishes",
                    "May Chaos take the world!",
                    "Burn it all with the yellow chaos flame",
                    "Until all is one again",
                    "The blessing of dispar",
                    "Wololoo!",
                    "Stop right there criminal scum!",
                    "I used to be an adventurer like you",
                    "The wind is howling",
                    "The night is dark and full of terrors",
                    "Hey!, Henry's come to see us",
                ],
                subheadlines: [
                    "Stay a while and listen",
                    "May you'r blood never boil",
                    "git gud, scrub",
                    "Is this really necessary?",
                    "What is a man, but a miserable little pile of secrets?",
                    "Take an arrow to the knee",
                    "Ah shi, here we go again",
                    "All of your bases are belong to us",
                    "We are the end of the world",
                    "I am Error",
                    "Not now, not later, not ever!",
                    "I need to meditate, or mas*****te, or both",
                    "I am the danger",
                    "I am the one who knocks",
                    "Jesus Christ be praised",
                    "War, war never changes",   
                ]
            }
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        loadCategoriesClusters();

        defineComponentHotkeys();
    });

    onDestroy(() => {
        if (browser) {
            dropHotkeysContext();
        }
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
        
        /*=============================================
        =            Hotkeys            =
        =============================================*/
            // Only call directly defineComponentHotkeys. All the other methods in this section should only be called by the global hotkeys manager.

            const defineComponentHotkeys = () => {
                if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    const hotkeys_context = new HotkeysContext();

                    hotkeys_context.register(["w", "a", "s", "d"], handleDungeonHotkeyMovement, {
                        description: `<dungeon_selection> Moves the focus between the dungeons.`,
                    });

                    hotkeys_context.register(["e"], handleClusterSelection, {
                        description: `<dungeon_selection> Selects the focused dungeon.`,
                    });

                    hotkeys_context.register(["del"], handleDeleteClusterData, {
                        description: `<dungeon_managment> Deletes the focused dungeon data from the database but leaves the directory tree intact.`,
                    });

                    hotkeys_context.register(["c c"], handleClusterRename, {
                        description: `<dungeon_managment> Renames the focused dungeon.`,
                    });

                    hotkeys_context.register(["n"], toggleClusterCreationTool, {
                        description: `<dungeon_managment> Opens the dungeon creation tool.`,
                    });

                    hotkeys_context.register(["?"], toggleHotkeysSheet, {
                        description: `<${HOTKEYS_GENERAL_GROUP}> Opens the hotkeys sheet.`,
                    });

                    global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
                }

                global_hotkeys_manager.loadContext(hotkeys_context_name);
            }

            /**
             * Handles the dungeon hotkey movement.
             * @param {KeyboardEvent} event
             * @param {import('@libs/LiberyHotkeys/hotkeys').HotkeyData} hotkey
             */
            const handleDungeonHotkeyMovement = (event, hotkey) => {
                let key_combo = hotkey.KeyCombo.toLowerCase();

                const cluster_count = categories_clusters.length;
                const row_count = Math.ceil(cluster_count / cluster_per_row);

                let new_focused_cluster_index = focused_cluster_index;

                switch (key_combo) {
                    case "w": // Up
                        new_focused_cluster_index -= cluster_per_row;
                        
                        new_focused_cluster_index = new_focused_cluster_index < 0 ? ((row_count - 1) * cluster_per_row) + focused_cluster_index : new_focused_cluster_index;
                        new_focused_cluster_index = new_focused_cluster_index >= cluster_count ? new_focused_cluster_index - cluster_count : new_focused_cluster_index;
                        break;
                    case "s": // Down
                        new_focused_cluster_index += cluster_per_row;
                        
                        new_focused_cluster_index = new_focused_cluster_index >= cluster_count ? focused_cluster_index - ((row_count - 1) * cluster_per_row) : new_focused_cluster_index;
                        new_focused_cluster_index = new_focused_cluster_index < 0 ? new_focused_cluster_index + cluster_count : new_focused_cluster_index;
                        break;  
                    case "a": // Left
                        new_focused_cluster_index -= 1;
                        
                        new_focused_cluster_index = new_focused_cluster_index < 0 ? cluster_count - 1 : new_focused_cluster_index;
                        break;
                    case "d": // Right
                        new_focused_cluster_index += 1;
                        
                        new_focused_cluster_index = new_focused_cluster_index >= cluster_count ? 0 : new_focused_cluster_index;
                        break;
                }

                focused_cluster_index = new_focused_cluster_index;
            }

            /**
             * handles the selection of a cluster via hotkeys.
             */
            const handleClusterSelection = () => {
                const selected_cluster_element = document.querySelector(`#dciw-cluster-item-${focused_cluster_index} > :first-child`);
                if (selected_cluster_element === null) {
                    console.error(`The selected cluster element with index ${focused_cluster_index} was on the dom.`);
                    return;
                }
                console.log('selected_cluster_element: ', selected_cluster_element);

                selected_cluster_element.click();
            }

            /**
             * Emits a cluster-rename event on the keyboard focused element, this element is found by the class name 'is-keyboard-focused'.
             * @returns {void}
             */
            const handleClusterRename = () => {
                const selected_cluster_element = document.querySelector("ul#dsp-libery-dungeons-grid .is-keyboard-focused");

                if (selected_cluster_element === null) {
                    console.error(`The selected cluster element with index ${focused_cluster_index} was on the dom.`);
                    return;
                }

                let custom_event = new CustomEvent("cluster-rename", {
                    bubbles: false,
                });

                selected_cluster_element.dispatchEvent(custom_event);
            }


            /**
             * Handles the deletion of a cluster data.
             * @returns {void}
            */
            const handleDeleteClusterData = async () => {
                const user_confirmation = await confirmPlatformMessage({
                    message_title: "Are you sure you want to delete all the cluster data?",
                    question_message: "This action will delete the cluster data from the database but will leave the directory tree intact. you will need to recreate the cluster to add it back to the system.",
                    cancel_label: "Keep cluster",
                    confirm_label: "Delete only data",
                    danger_level: 2,
                    auto_focus_cancel: true,
                });

                if (user_confirmation !== 1) {
                    return;
                }

                const cluster_to_delete = categories_clusters[focused_cluster_index];
                if (cluster_to_delete === undefined) {
                    console.error(`The cluster to delete with index ${focused_cluster_index} was not found.`);
                    return;
                }

                let deleted = await deleteClusterRecord(cluster_to_delete.UUID);

                if (deleted) {
                    emitPlatformMessage(`The cluster '${cluster_to_delete.Name}' data was deleted successfully.`)
                    loadCategoriesClusters();
                } else {
                    let variable_enviroment = new VariableEnvironmentContextError("In @models/CategoriesClusters.handleDeleteClusterData" );
                    variable_enviroment.addVariable("cluster_to_delete", cluster_to_delete.UUID);
                    variable_enviroment.addVariable("deleted", deleted);
                    variable_enviroment.addVariable("user_confirmation", user_confirmation);

                    let labeled_error = new LabeledError(variable_enviroment, `There was an error while trying to delete data from the cluster '${cluster_to_delete.UUID}'.`, lf_errors.ERR_PROCESSING_ERROR);

                    emitPlatformMessage(labeled_error);
                }


            }
        
        /*=====  End of Hotkeys  ======*/

        /**
         * Gets a list of phrases and returns a random one.
        */
        const chooseRandomPhrase = (phrases) => {
            return phrases[Math.floor(Math.random() * phrases.length)];
        }

        /**
         * Drops the hotkeys context.
         * @returns {void}
        */
        const dropHotkeysContext = () => {
            global_hotkeys_manager.dropContext(hotkeys_context_name);
        }

        /**
         * Handles the creation of a new cluster by the ClusterCreationTool.
         */
        const handleClusterCreated = () => {
            show_cluster_creation_tool = false;
            loadCategoriesClusters();
        }
    
        /**
         * Loads the categories clusters of the system.
         */
        const loadCategoriesClusters = async () => {
            categories_clusters = await getAllCategoriesClusters();
            console.debug(categories_clusters);
        }

        /**
         * Toggles the visibility of the cluster creation tool.
         */
        const toggleClusterCreationTool = () => {
            show_cluster_creation_tool = !show_cluster_creation_tool;
        }
    
    /*=====  End of Methods  ======*/
    
</script>

<main id="dungeons-selector-page">
    {#if show_cluster_creation_tool}
        <div id="cluster-creation-tool-component-wrapper">
            <ClusterCreationTool
                on:cluster-created={handleClusterCreated}
            />
        </div>
    {/if}
    <header id="dsp-dungeons-header">
        <hgroup id="dsp-dungeons-headlines">
            <h1 id="dsp-dungeons-headline">
                {chooseRandomPhrase(gaming_phrases.headlines)}
            </h1>
            <p id="dsp-dungeons-subheadline">
                Choose a dungeon. Then {chooseRandomPhrase(gaming_phrases.subheadlines)}.
            </p>
        </hgroup>
    </header>
    <ul id="dsp-libery-dungeons-grid" 
        class:dtwo={false}
        style:grid-template-columns={`repeat(${cluster_per_row}, 1fr)`}
    >
        {#each categories_clusters as dungeon_cluster, h}
            {@const is_keyboard_focused = h === focused_cluster_index}
            <li id="dciw-cluster-item-{h}"
                class="dungeon-cluster-item-wrapper"
            >
                <CategoriesClustersItem 
                    cluster={dungeon_cluster}
                    keyboard_focused={is_keyboard_focused}
                />
                {#if is_keyboard_focused}
                    <DungeonSelectionHint />
                {/if}
            </li>
        {/each}
    </ul>
</main>

<style>
    #dungeons-selector-page {
        position: relative;
        display: flex;
        height: calc(100dvh - var(--navbar-height));
        flex-direction: column;
        align-items: center;
        gap: calc(var(--vspacing-4) * 1.5);
    }

    
    /*=============================================
    =            Headlines            =
    =============================================*/
    
        #dsp-dungeons-header {
            display: flex;
            
            flex-direction: column;
            align-items: center;
            padding: var(--vspacing-5) 0 0 0;
        }

        #dsp-dungeons-header hgroup#dsp-dungeons-headlines {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: var(--vspacing-2);
            text-align: center;

            & h1#dsp-dungeons-headline {
                font-size: var(--font-size-h2);
                line-height: 1;
            }

            & p#dsp-dungeons-subheadline {
                font-size: var(--font-size-p);
                font-family: var(--font-decorative);
                color: var(--grey-3);
                line-height: 1;
            }
        }
    
    /*=====  End of Headlines  ======*/
    
    

    #dsp-libery-dungeons-grid {
        display: grid;
        width: max-content;
        gap: var(--vspacing-4);
        height: 56dvh;
        justify-content: center;
        overflow-y: auto;
        list-style: none;
        padding: 0 var(--spacing-2);
        scrollbar-width: thin;
        scrollbar-color: transparent transparent;

        & > li {
            display: flex;
            flex-direction: column;
            width: 170px;
            height: 315px;
            justify-content: space-between;
        }
    }

    
    /*=============================================
    =            Cluster creation tool            =
    =============================================*/
    
        #cluster-creation-tool-component-wrapper {
            --cct-top-position: 57%;
            --cct-left-position: 50%;

            width: 40%;
            height: 70dvh;
            position: absolute;
            top: var(--cct-top-position);
            left: var(--cct-left-position);
            translate: calc(var(--cct-top-position) * -1) calc(var(--cct-left-position) * -1);
            transform-origin: center;
            z-index: var(--z-index-t-1);
        }
    
    /*=====  End of Cluster creation tool  ======*/
    
    
</style>