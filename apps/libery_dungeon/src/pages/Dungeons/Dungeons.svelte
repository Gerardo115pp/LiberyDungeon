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
        import { CursorMovementWASD } from "@app/common/keybinds/CursorMovement";
        import { ui_core_dungeon_references } from "@app/common/ui_references/core_ui_references";
        import { ensureElementVisible } from "@libs/utils";
    
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
         * The grid navigation wrapper.
         * @type {CursorMovementWASD | null}
         */
        let the_grid_navigation_wrapper = null;
        
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
            dropCursorMovementWASD();
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
                if (global_hotkeys_manager == null) {
                    console.error("In Dungeons.defineComponentHotkeys: The global hotkeys manager was not found.");
                    return;
                }
                
                if (!global_hotkeys_manager.hasContext(hotkeys_context_name)) {
                    const hotkeys_context = new HotkeysContext();

                    // hotkeys_context.register(["w", "a", "s", "d"], handleDungeonHotkeyMovement, {
                    //     description: `<dungeon_selection> Moves the focus between the dungeons.`,
                    // });

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

                    setGridNavigationWrapper(hotkeys_context);

                    global_hotkeys_manager.declareContext(hotkeys_context_name, hotkeys_context);
                }

                global_hotkeys_manager.loadContext(hotkeys_context_name);
            }

            /**
             * Sets the grid navigation wrapper required data.
             * @param {import("@libs/LiberyHotkeys/hotkeys_context").default} hotkeys_context
             */
            const setGridNavigationWrapper = (hotkeys_context) => {
                if (the_grid_navigation_wrapper != null) {
                    the_grid_navigation_wrapper.destroy(); // This allows this function to be used for updating the navigation grid.
                }

                const grid_selectors = getGridSelectors();

                const matching_elements_count = document.querySelectorAll(grid_selectors.grid_parent_selector).length;
                if (matching_elements_count !== 1) {
                    throw new Error(`tag parent selector '${grid_selectors.grid_parent_selector}' returned ${matching_elements_count}, expected exactly 1`);
                }


                the_grid_navigation_wrapper = new CursorMovementWASD(grid_selectors.grid_parent_selector, handleCursorUpdate, {
                    initial_cursor_position: focused_cluster_index,
                    sequence_item_name: ui_core_dungeon_references.CATEGORY_CLUSTER.EntityName,
                    sequence_item_name_plural: ui_core_dungeon_references.CATEGORY_CLUSTER.EntityNamePlural,
                    grid_member_selector: grid_selectors.grid_member_selector,
                });
                the_grid_navigation_wrapper.setup(hotkeys_context);
            }

            /**
             * handles the selection of a cluster via hotkeys.
             */
            const handleClusterSelection = () => {
                /**
                 * @type {HTMLElement | null}
                 */
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
             * @returns {Promise<void>}
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

                    labeled_error.alert();
                }


            }

            /**
             * Handles the Cursor update event emitted by the_grid_navigation_wrapper.
             * @type {import("@common/keybinds/CursorMovement").CursorPositionCallback}
             */
            const handleCursorUpdate = (cursor_wrapped_value) => {
                focused_cluster_index = cursor_wrapped_value.value;

                ensureCurrentClusterElementVisible();

                return false;
            }
        
        /*=====  End of Hotkeys  ======*/

        /**
         * Gets a list of phrases and returns a random one.
         * @param {string[]} phrases - The list of phrases to choose from.
         */
        const chooseRandomPhrase = (phrases) => {
            return phrases[Math.floor(Math.random() * phrases.length)];
        }

        /**
         * Ensures the current cluster element is visible.
         * @returns {void}
         */
        const ensureCurrentClusterElementVisible = () => {
            const cluster_element = getCurrentClusterHTMLElement();

            if (cluster_element === null) {
                console.warn("In @pages/Dungeons/Dungeons.ensureCurrentClusterElementVisible: The current cluster element was not found.");
                return;
            }

            ensureElementVisible(cluster_element);
        }

        /**
         * Drops the hotkeys context.
         * @returns {void}
         */
        const dropHotkeysContext = () => {
            if (global_hotkeys_manager == null) return;

            global_hotkeys_manager.dropContext(hotkeys_context_name);
        }

        /**
         * Frees the resources used by the CursorMovementWASD instance.
         * @returns {void}
         */
        const dropCursorMovementWASD = () => {
            if (the_grid_navigation_wrapper != null) {
                the_grid_navigation_wrapper.destroy();
            }
        }

        /**
         * Returns the grid selector for the TaxonomyTags component.
         * @returns {import('@common/interfaces/common_actions').GridSelectors}
         */
        const getGridSelectors = () => {
            const parent_selector = '#dsp-content-wrapper > ul#dsp-libery-dungeons-grid';
            const child_selector = `.dungeon-cluster-item-wrapper`

            return {
                grid_parent_selector: parent_selector,
                grid_member_selector: child_selector,
            }
        }

        /**
         * Returns the cluster item wrapper element for the given item index.
         * @param {number} item_index
         * @returns {HTMLElement | null}
         */
        const getClusterHTMLElmenet = item_index => {
            const cluster_element = document.getElementById(`dciw-cluster-item-${item_index}`)
            
            return cluster_element;
        }

        /**
         * Returns the current cluster element.
         * @returns {HTMLElement | null}
         */
        const getCurrentClusterHTMLElement = () => {
            return getClusterHTMLElmenet(focused_cluster_index);
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
    <div id="dsp-underlay-wrapper">
        <img
            aria-hidden="true"
            decoding="async" 
            class="dsp-underlay"
            src="https://libery-dungeon.com/medias-api/random-medias-fs?cluster_id=72739c6a-5c29-47ba-b2c1-f554750f6788&category_id=f491e6df7d83435f1900a74e9538b9b2714871bf" alt=""
        >
    </div>
    {#if show_cluster_creation_tool}
        <div id="cluster-creation-tool-component-wrapper">
            <ClusterCreationTool
                on:cluster-created={handleClusterCreated}
            />
        </div>
    {/if}
    <div id="dsp-content-wrapper">
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
    </div>
</main>

<style>
    #dungeons-selector-page {
        --dsp-content-width: 1000px;

        position: relative;
        height: calc(100dvh - var(--navbar-height));
        container-type: inline-size;
        display: flex;
        justify-content: center;
    }

    #dsp-content-wrapper {
        background: hsl(from var(--body-bg-color) h s l / 0.98);
        width: var(--dsp-content-width);
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: calc(var(--vspacing-4) * 1.5);

        &::before {
            content: "";
            position: absolute;
            background-image: radial-gradient(ellipse 90% 71% at 18%, hsl(from var(--body-bg-color) h s l / 0.3) 0%, hsl(from var(--body-bg-color) h s l / 0.5) 20%, hsl(from var(--body-bg-color) h s l / 0.7) 45%, hsl(from var(--body-bg-color) h s l / 0.90) 70%, hsl(from var(--body-bg-color) h s l / 0.98) 90%);
            inset: 0;
            width: calc(calc(100dvw - var(--dsp-content-width)) / 2);
            height: 100%;
        }

        &::after {
            content: "";
            position: absolute;
            background-image: radial-gradient(ellipse 90% 71% at 80%, transparent 0%, hsl(from var(--body-bg-color) h s l / 0.5) 20%, hsl(from var(--body-bg-color) h s l / 0.7) 45%, hsl(from var(--body-bg-color) h s l / 0.90) 70%, hsl(from var(--body-bg-color) h s l / 0.98) 90%);
            inset: 0 0 auto auto;
            width: calc(calc(100dvw - var(--dsp-content-width)) / 2);
            height: 100%;
            right: 0;
        }
    }
    
    /*=============================================
    =            Dungeon underlay            =
    =============================================*/
        
        #dsp-underlay-wrapper {
            position: absolute;
            top: 0;
            left: 0;
            width: 100dvw;
            height: calc(100dvh - var(--navbar-height));
            z-index: var(--z-index-b-5);
        }

        img.dsp-underlay {
            width: 100%;
            height: 100%;
            object-fit: cover;
        }        
    
        @container (width > 1950px) {
            #dsp-content-wrapper::before {
                background-image: radial-gradient(ellipse 80% 81% at 0%, hsl(from var(--body-bg-color) h s l / 0.3) 0%, hsl(from var(--body-bg-color) h s l / 0.5) 20%, hsl(from var(--body-bg-color) h s l / 0.7) 45%, hsl(from var(--body-bg-color) h s l / 0.90) 70%, hsl(from var(--body-bg-color) h s l / 0.98) 90%);
            }

            #dsp-content-wrapper::after {
                background-image: radial-gradient(ellipse 80% 81% at 100%, hsl(from var(--body-bg-color) h s l / 0.4) 0%, hsl(from var(--body-bg-color) h s l / 0.5) 20%, hsl(from var(--body-bg-color) h s l / 0.7) 45%, hsl(from var(--body-bg-color) h s l / 0.90) 70%, hsl(from var(--body-bg-color) h s l / 0.98) 90%);
            }
        }

    /*=====  End of Dungeon underlay  ======*/
    
    /*=============================================
    =            Headlines            =
    =============================================*/
    
        #dsp-dungeons-header {
            display: flex;
            flex-direction: column;
            align-items: center;
            padding: var(--spacing-5) 0 0 0;
        }

        #dsp-dungeons-header hgroup#dsp-dungeons-headlines {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: var(--spacing-2);
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
        --grid-items-width: 170px;
        --grid-items-per-row: 3;
        --grid-items-gaps: calc(var(--grid-items-per-row) - 1);
        --grid-items-gap: var(--spacing-4);
        --grid-inline-padding: var(--spacing-2);
        --grid-block-padding: 0;

        box-sizing: content-box;
        display: flex;
        width: calc(calc(var(--grid-items-per-row) * var(--grid-items-width)) + calc(var(--grid-items-gaps) * var(--grid-items-gap)) + calc(2 * var(--grid-inline-padding)));
        flex-wrap: wrap;
        gap: var(--grid-items-gap);
        height: 56dvh;
        justify-content: center;
        overflow-y: auto;
        list-style: none;
        padding-inline: var(--grid-inline-padding);
        padding-block: var(--grid-block-padding);
        scrollbar-width: thin;
        scrollbar-color: transparent transparent;

        & > li {
            display: flex;
            flex-direction: column;
            width: var(--grid-items-width);
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