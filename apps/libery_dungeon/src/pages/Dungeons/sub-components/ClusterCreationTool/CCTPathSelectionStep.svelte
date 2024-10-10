<script>
    import { createEventDispatcher, onMount } from "svelte";
    import { getClusterRootPath, getNewClusterDirectoryOptions, validateClusterDirectory } from "@models/CategoriesClusters";
    import { LabeledError } from "@libs/LiberyFeedback/lf_models";
    import { emitPlatformMessage } from "@libs/LiberyFeedback/lf_utils";
    import CctPathEligibilityRules from "./CCTPathEligibilityRules.svelte";


    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The default cluster root path.
         * @type {string}
         */    
        let default_cluster_root_path = null;

        /**
         * The selected path directory options to create a new cluster in.
         * @typedef {Object} DirectoryOption
         * @property {string} path - The path to create the cluster in
         * @property {string} name - The name of the directory
        */

        let dispatch = createEventDispatcher();
        
        
        /*----------  State  ----------*/
        
            /**
             * The path we are show directory options for.
             * @type {string}
             */
            let current_selected_path = "";

            /**
             * Is path safe. Whether the path came from one of the directory options or if it was typed by the user.
             * @type {boolean}
             */
            let is_path_safe = false;

            /**
             * @type {DirectoryOption[]}
             */
            let directory_options = [];

            /**
             * Whether to show the directory eligibility rules.
             * @type {boolean}
             */
            let show_directory_eligibility_rules = false;

    /*=====  End of Properties  ======*/

    onMount(async () => {
        default_cluster_root_path = await retrieveDefaultClusterRootPath();

        if (default_cluster_root_path !== null) {
            current_selected_path = default_cluster_root_path;
        }
        
        retrieveDirectoryOptions();
    });

    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Retrieves the default directory where clusters are created in the selected server.
         * @returns {Promise<string | null>}
         */
        const retrieveDefaultClusterRootPath = async () => {
            let root_path = null; 
            
            try {
                root_path = await getClusterRootPath();
            } catch (error) {
                if (error instanceof LabeledError) {
                    error.alert();
                } else {
                    console.error(error);
                }
            }

            return root_path;
        }
    
        /**
         * Handles the click event a directory option.
         * @param {MouseEvent} e - The click event.
        */
        const handleDirectoryOptionClick = async (e) => {
            /** @type {string} */
            let directory_option_index_value = e.currentTarget.dataset.directoryOptionIndex;

            let directory_option_index = parseInt(directory_option_index_value);

            if (directory_option_index >= 0 && directory_option_index < directory_options.length) {
                current_selected_path = directory_options[directory_option_index].path;
                let path_availability = await validatePath(current_selected_path);

                is_path_safe = path_availability.is_path_valid;
                console.log(`Path availability reason: ${path_availability.reason}`);
            }

            retrieveDirectoryOptions();
        }

        /**
         * Handles the click event of the create cluster button.
         * @param {MouseEvent} e - The click event.
        */
        const handleCreateClusterButtonClick = e => {
            dispatch("new-cluster-path-selected", { path: current_selected_path });
        }

        /**
         * Handles the keypress event on the custom path input.
         * @param {KeyboardEvent} e - The keypress event.
        */
        const handleCustomPathInputKeyPress = async e => {
            if (e.key === "Enter") {
                /**
                 * @type {import('@models/CategoriesClusters').PathAvailability}
                 */
                let path_availability = await validatePath(current_selected_path);
                
                if (path_availability.is_path_valid !== is_path_safe) {
                    is_path_safe = path_availability.is_path_valid;
                }

                if (!path_availability.is_path_valid) {
                    emitPlatformMessage(path_availability.reason);
                }
            } else {
                is_path_safe = false;
            }
        }

        /**
         * Handles the eligibility rules label MouseEnter event.
         * @param {MouseEvent} e
        */
        const handleEligibilityRulesLabelMouseEnter = e => {
            show_directory_eligibility_rules = true;
        }

        /**
         * Handles the Instructions Wrapper MouseLeave event.
         * @param {MouseEvent} e
        */
        const handleInstructionsWrapperMouseLeave = e => {
            show_directory_eligibility_rules = false;
        }
    
        /**
         * Calls getNewClusterDirectoryOptions to get a list of directory options from the `current_selected_path`.
         * @requires current_selected_path
         * @returns {promise<void>}
         */
        const retrieveDirectoryOptions = async () => {
            /** @type {DirectoryOption[]}*/
            let new_directory_options = await getNewClusterDirectoryOptions(current_selected_path);

            if (new_directory_options.length !== 0) {
                directory_options = new_directory_options;
            }
        }

        /**
         * Validates the path to see if it is safe.
         * @returns {Promise<import('@models/CategoriesClusters').PathAvailability>}
         */
        const validatePath = async requested_path => await validateClusterDirectory(requested_path);
    
    /*=====  End of Methods  ======*/
    
    
</script>

<div id="cct-cluster-path-selection-step">
    <header id="cct-instructions-wrapper" on:mouseleave={handleInstructionsWrapperMouseLeave}>
        <h2 id="cct-instruction-label">
            Create a new dungeon.
        </h2>
        <p id="cct-instructions">
            These are the available directories provided by your deployed Categories services. Select one to create a new dungeon in. These are the directory <b on:mouseenter={handleEligibilityRulesLabelMouseEnter}>eligibility rules</b>
        </p>
        <CctPathEligibilityRules 
            show_path_elegibility_rules={show_directory_eligibility_rules}
        />
    </header>
    <ul id="cct-directory-options-container">
        {#each directory_options as directory_option, h}
            <li class="cct-directory-option" on:click={handleDirectoryOptionClick} data-directory-option-index={h}>
                <span class="cct-directory-name">{directory_option.name}</span>
                <span class="cct-directory-path"><i>{directory_option.path}</i></span>
            </li>
        {/each}
    </ul>
    <section id="cct-controls-section">
        <label id="cct-custom-path-input">
            <p class="field-tool-tip">
                Hit Enter to confirm.
            </p>
            <span>
                Dungeon path:
            </span>
            <input id="cct-custom-path-input"
                type="text"
                class:is-path-safe={is_path_safe}
                on:keypress={handleCustomPathInputKeyPress} 
                bind:value={current_selected_path}
            >
        </label>
        <button id="cct-create-cluster-button" class="dungeon-button-1" on:click={handleCreateClusterButtonClick} disabled={!is_path_safe}>
            Create
        </button>
    </section>
</div>

<style>

    #cct-cluster-path-selection-step {
        display: contents;
    }

    
    /*=============================================
    =            Header            =
    =============================================*/
    
        header#cct-instructions-wrapper {
            position: relative;
            display: flex;
            flex-direction: column;
            height: 13cqw;
            align-items: center;
            gap: var(--vspacing-1);

            & h2#cct-instruction-label {
                font-size: var(--font-size-h4);
                color: var(--main);
            }

            & p#cct-instructions {
                font-size: calc(var(--font-size-p-small) * .9);
                color: var(--grey-3);
                text-align: center;
            }

            & p#cct-instructions > b {
                color: var(--main);
                text-decoration: underline;
                cursor: help;
            }
        }
    
    
    /*=====  End of Header  ======*/
    
    


    ul#cct-directory-options-container {
        display: flex;
        overflow-y: auto;
        height: 50cqw;
        background: var(--grey);
        flex-direction: column;
        padding: 0;
        gap: var(--vspacing-1);
        list-style: none;

        scrollbar-width: thin;
        scrollbar-color: var(--main) var(--grey-9);
    }

    li.cct-directory-option {
        display: flex;
        flex-direction: column;
        width: 70cqw;
        padding: var(--vspacing-2);
        background: var(--grey-8);
        
        & span.cct-directory-name {
            font-size: var(--font-size-2);
            color: var(--main-dark);
            font-weight: bold;
        }

        & span.cct-directory-path {
            font-size: var(--font-size-fineprint);
            color: var(--grey-2);
        }
    }

    section#cct-controls-section {
        width: 70cqw;
        display: flex;
        align-items: center;
        justify-content: space-between;

        & > label {
            position: relative;
        }

        & label > .field-tool-tip {
            visibility: hidden;
            position: absolute;
            bottom: 110%;
            background: var(--grey);
            color: var(--grey-1);
            font-size: var(--font-size-fineprint);
            text-wrap: pretty;
            transition: all 0.2s ease-out;
        }
        & label:hover > .field-tool-tip {
            visibility: visible;
        }

    }

    label#cct-custom-path-input {
        display: flex;
        background-color: var(--grey-6);
        width: 70%;
        height: max-content;
        border-radius: var(--border-radius);
        align-items: center;
        gap: calc(var(--vspacing-1) * 0.5);

        & > span:first-of-type {
            font-size: var(--font-size-p-small);
            line-height: 1;
            color: var(--grey-2);
            padding: 0 0 0 var(--vspacing-1);
        }

        & > input {
            font-size: var(--font-size-p-small);
            line-height: 1;
            flex-grow: 2;
            /* font-weight: lighter; */
            background: transparent;
            color: var(--main-dark);
            border: none;
            outline: none;
        }

        & > input.is-path-safe {
            color: var(--success);
        }
    }

    @media (pointer:fine) {
        li.cct-directory-option {
            transition: all 0.2s ease-out;
        }

        li.cct-directory-option > * {
            user-select: none;
        }

        li.cct-directory-option:hover {
            background: var(--grey-7);
        }
    }
</style>