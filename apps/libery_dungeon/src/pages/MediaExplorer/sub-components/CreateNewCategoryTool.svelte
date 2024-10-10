<script>
    import { isCategoryNameAvailable } from "@libs/CategoriesUtils";
    import { current_category, categories_tree } from "@stores/categories_tree";
    import { current_cluster } from "@stores/clusters";
    import { category_creation_tool_mounted } from "@pages/MediaExplorer/app_page_store";
    import { LabeledError, VariableEnvironmentContextError } from "@libs/LiberyFeedback/lf_models";
    import { err_categories } from "@errors/err_categories";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The name to create the new category with.
         * @type {string}
         */
        let new_category_name = "";

        /**
         * The input field where the user types the new category name.
         * @type {HTMLInputElement}
         */
        let category_name_input_element;

        /**
         * The form element that contains the input field.
         * @type {HTMLFormElement}
         */
        let create_category_form;

        /**
         * Has the fist keydown event been triggered on the input field? used to prevent the hotkey from been added as part of the category name.
         * @type {boolean}
         */ 
        let first_keydown_triggered = false;

    
    /*=====  End of Properties  ======*/
   
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the keydown event on the new category name input field.
         * @param {KeyboardEvent} event 
         */
        const handleNewCategoryNameInput = async (event) => {
            if (!first_keydown_triggered) {
                event.preventDefault();
                first_keydown_triggered = true;
                return;
            }
            
            if (event.key === "Enter") {
                event.preventDefault();
                event.stopPropagation();
                let is_name_valid = await verifyCategoryName();
                if (is_name_valid) {
                    handleCreateNewCategory();
                }
            }

            if (event.key === "Escape") {
                event.preventDefault();
                event.stopPropagation();
                category_creation_tool_mounted.set(false);
            }
        }

        /**
         * Prevents the the hotkey keystroke from being added to the category name.
         * @param event
         */
        const handleKeyPrevent = (event) => {
            if (!first_keydown_triggered) {
                event.preventDefault();
                first_keydown_triggered = true;
            }
        }
    
        const handleCreateNewCategory = async () => {
            let err;
            let created_category;

            window.queueMicrotask(() => {
                category_creation_tool_mounted.set(false);
            });

            try {
                created_category = await $categories_tree.insertChildCategory(new_category_name, $current_cluster.UUID);
            } catch (error) {
                if (error instanceof LabeledError) {
                    console.log("ALERTING ERROR: ", error);
                    error.alert();
                    return;
                } else if (!(error instanceof Error)) {
                    err = new Error(error);
                    console.log("failed to create category: ", err);
                }
            }

            console.log("Created category: ", created_category);

            if (created_category == null) {
                let error_context = new VariableEnvironmentContextError("In CreateNewCategoryTool.handleCreateNewCategory")

                error_context.addVariable("new_category_name", new_category_name);
                error_context.addVariable("$current_cluster", $current_cluster);
                error_context.addVariable("created_category", created_category);

                err = new LabeledError(error_context, `Could not create category '${new_category_name}'`, err_categories.ERR_COULD_NOT_CREATE);
                err.alert();
            }
        }

        /**
         * Verifies the validity and availability of the new category name.
         * @returns {Promise<boolean>}
         */
        const verifyCategoryName = async () => {
            if (category_name_input_element == null || new_category_name === "") return;

            let name_valid = category_name_input_element.checkValidity();



            if (name_valid) {
                name_valid = await isCategoryNameAvailable(new_category_name, $current_category.uuid);
            }

            return name_valid;
        }
    
    /*=====  End of Methods  ======*/

</script>
    
{#if $category_creation_tool_mounted}
    <aside id="create-category-tool-wrapper">
        <form bind:this={create_category_form} id="create-category-tool" class="libery-dungeon-window">
            <h3 id="cct-tool-title">
                Create new category
            </h3>
            <p id="cct-tool-instructions">
                Please note that the creating of a new category will fail if there is another category with the same name in the current category. Also charactes like '/', '\', ':', '*', '?', '"' are not allowed in category names.
            </p>
            <label id="cct-new-category-name-field">
                <span>
                    Category name
                </span>
                <input id="cct-ncn-input" 
                    bind:value={new_category_name}
                    bind:this={category_name_input_element}
                    type="text" 
                    maxlength="200"
                    pattern="{'[a-zA-Z\\d\\s_\\-]{2,}'}"
                    on:keydown={handleNewCategoryNameInput}
                    on:keypress={handleKeyPrevent}
                    autofocus
                    required
                >
            </label>
        </form>
    </aside>
{/if}

<style>
    #create-category-tool-wrapper {
        position: absolute;
        container-type: size;
        display: grid;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        place-items: center;
    }
    
    #create-category-tool {
        box-sizing: content-box;
        display: flex;
        flex-direction: column;
        width: 30cqw;
        height: 30cqh;
        justify-content: center;
        align-items: center;
        padding: 0 var(--vspacing-3);
        gap: var(--vspacing-2);
    }

    p#cct-tool-instructions {
        font-size: calc(var(--font-size-p) * 0.7);
        color: var(--grey-4);
    }

    label#cct-new-category-name-field {
        display: flex;
        width: 100%;
        background: var(--grey);
        border: 1px solid var(--main);
        padding: var(--vspacing-1);
        border-radius: var(--border-radius);

        & > span {
            font-size: var(--font-size-p);
            padding: 0 var(--vspacing-1);
        }

        & > span::after {
            content: ":";
        }

        & > input {
            flex-grow: 2;
            font-size: var(--font-size-p);
            border: none;
            background: transparent;
            color: var(--main-dark-color-4);
            outline: none;
        }
    }

    label#cct-new-category-name-field:has(input:valid) {
        border-color: var(--success);
    }
        

    @media only screen and (max-width: 767px) {
        #create-category-tool {
            width: 90vw;
            height: 30vh;
        }
    }
</style>