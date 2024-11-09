<script>
    /*=============================================
    =            Imports            =
    =============================================*/    
        import InnerCategoryView from "./InnerCategoryView.svelte";
        import { getShortCategoryTree } from "@models/Categories";
        import { InnerCategory } from "@models/Categories";
        import { createEventDispatcher, onDestroy, onMount } from "svelte";
        import { create_subcategory } from "@pages/MediaViewer/app_page_store";
        import { getCategory, createCategory } from "@models/Categories";
        import { categories_tree } from "@stores/categories_tree";
        import { current_cluster } from "@stores/clusters";
    /*=====  End of Imports  ======*/
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type { InnerCategory } */
        export let category;
        
        /** @type { InnerCategory[] } */
        let inner_categories = [];
        
        /** @type { boolean } */
        let is_expanded = false;
        
        /** @type { boolean } if the keyboard selector is over this element */
        export let is_keyboard_selected = false;
        
        /** @type { boolean } whether to use a light or a dark color */
        export let is_light = false;
        
        /** @type {boolean} whether to show inner categories or not */
        export let show_inner_categories = true;

        /* -------------------------- sub-category creation ------------------------- */

            /** @type { boolean } if true, show the create subcategory input */
            let show_create_subcategory_input = false;

            /**
             * The input field used to create a new subcategory
             * @type {HTMLInputElement | null}
             */
            let the_subcategory_creation_input = null;

            /**
             * The subcategory new name.
             * @type {string}
             * @default ""
             */
            let subcategory_name = "";

            /**
             * Whether the new subcategory name is valid.
             * @type {boolean} 
             */
            let new_subcategory_name_valid = false;

            /**
             * The subcategory creation form.
             * @type {HTMLFormElement | null}
             */
            let the_subcategory_creation_form = null;

        

        const dispatch = createEventDispatcher();

        
        /*----------  Unsubscribers  ----------*/
        
        let new_category_name_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        new_category_name_unsubscriber = create_subcategory.subscribe(setSubCategoryCreationState);
    })

    onDestroy(() => {
        new_category_name_unsubscriber();
    })
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles new subcategory creation.
         * @param {string} category_name
         */
        const createSubCategory = async category_name => {
            let category_complete_data = await getCategory(category.uuid);
        
            let new_category = await createCategory(category_complete_data.UUID, category_complete_data.FullPath, category_name, $current_cluster.UUID);
            console.debug("New subcategory created from MediaMovementsTool/InnerCategoryView.svelte", new_category);

            $categories_tree.deleteLoadedCategoryCache(category_complete_data);

            dispatch("new-subcategory-created", {
                category: new_category
            });
        }

        /**
         * Sets the subcategory creation state.
         * @param {boolean | null} enable
         */
        const setSubCategoryCreationState = enable => {
            the_subcategory_creation_form?.reset();

            show_create_subcategory_input = is_keyboard_selected && (enable === true); // only if true, not count truethy values
        }    

        /**
         * Creates a new subcategory with the subcategory_name if new_subcategory_name_valid is true.
         * @returns {Promise<void>}
         */
        const createSubCategoryIfValid = async () => {
            if (!new_subcategory_name_valid) return;

            subcategory_name = subcategory_name.trim();

            createSubCategory(subcategory_name);

            return;
        }

        /**
         * @param {MouseEvent} e
         */
        const handleCategoryItemClicked = e => {
            e.stopPropagation();

            let shift_key = e.shiftKey;

            dispatch("category-item-selected", { category: category, shift_key: shift_key });
        }

        /**
         * @param {KeyboardEvent} e 
         */
        const handleCreateSubcategoryInputKeypress = e => {
        }

        /**
         * @param {KeyboardEvent} e 
         */
        const handleCreateSubcategoryInputKeydown = e => {
            if (e.key === "Enter" || (e.key === "e" && e.ctrlKey)) {
                e.preventDefault();

                if ($create_subcategory === false) return;

                create_subcategory.set(false);

                createSubCategoryIfValid();
            }

            if (e.key === "Escape") {
                the_subcategory_creation_input?.blur();
                show_create_subcategory_input = false;
                create_subcategory.set(false);
            }
        }

        /**
         * @param {KeyboardEvent} event 
         */
        const handleCreateSubcategoryInputKeyup = event => {
            if (the_subcategory_creation_input == null) return;

            event.preventDefault();

            if (the_subcategory_creation_input.validationMessage !== "") {
                the_subcategory_creation_input.setCustomValidity("");
            }

            new_subcategory_name_valid = the_subcategory_creation_input.checkValidity();
        }

        const toggleInnerCategories = async () => {
            if (inner_categories.length === 0 && show_inner_categories) {
                inner_categories = await getShortCategoryTree(category.uuid, $current_cluster.UUID);
            }

            is_expanded = !is_expanded;
        }
    
    /*=====  End of Methods  ======*/
    
</script>


<li id="inner-category-{category.uuid}" class="inner-category" class:keyboard-selected={is_keyboard_selected}>
    <div class="inner-category-wrapper" class:is-light={is_light} on:click={toggleInnerCategories}>
        {#if show_create_subcategory_input}
            <form action="none" class="subcategory-creation-form"
                bind:this={the_subcategory_creation_form}
            >
                <label class="subcategory-creation-wrapper dungeon-input">
                    <span class="dungeon-label">
                        subcategory     
                    </span>
                    <input 
                        bind:this={the_subcategory_creation_input}
                        bind:value={subcategory_name}
                        type="text"
                        autocomplete="off"
                        on:keydown={handleCreateSubcategoryInputKeydown}
                        on:keyup={handleCreateSubcategoryInputKeyup}
                        on:keypress={handleCreateSubcategoryInputKeypress}
                        autofocus
                        required
                    >
                </label>
            </form>
        {:else}
            <strong class="ic-category-name" on:click={handleCategoryItemClicked}>
                {category.name}
            </strong>
            {#if category.fullpath != null && category.fullpath != InnerCategory.EMPTY_FULLPATH}
                <i class="ic-category-fullpath">
                    {category.fullpath}
                </i> 
            {/if}
        {/if}
    </div>
    {#if is_expanded && inner_categories.length > 0}
        <ul id="inner-categories-{category.uuid}">
            {#each inner_categories as ic}
                <svelte:component this={InnerCategoryView} 
                    category={ic} 
                    is_light={!is_light} 
                    on:category-item-selected
                />
            {/each}
        </ul>
    {/if}
</li>

<style>
    .inner-category-wrapper {
        width: 100%;
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        justify-content: flex-start;
        background: var(--grey-8);
        padding: var(--vspacing-2) var(--vspacing-3);
        gap: var(--vspacing-1);
        transition: background 0.2s ease-in;
    }

    .inner-category-wrapper.is-light {
        background: var(--grey-6);
    }

    .keyboard-selected {    
        outline: var(--main-4) dashed 2px;

        & .inner-category-wrapper {
            filter: brightness(1.2);
        }
    }

    strong.ic-category-name {
        cursor: default;
        display: inline-block;
        font-size: var(--font-size-2);
        color: var(--main-dark);
        transform-origin: center;
        padding: var(--vspacing-1);
        border-radius: var(--border-radius);
        transition-property: transform, background, opacity;
        transition-duration: 0.2s;
        transition-timing-function: ease-in-out;
    }

    i.ic-category-fullpath {
        display: block;
        font-size: var(--font-size-fineprint);
        font-style: normal;
        padding: 0 0 0 var(--vspacing-1);
        color: var(--main);
        text-wrap: pretty;
    }

    @media(pointer: fine) {
        .inner-category-wrapper:hover {
            background: var(--main-dark);
        }

        .inner-category-wrapper:hover strong.ic-category-name {
            color: var(--grey-1);
        }

        strong.ic-category-name:hover {
            transform: scale(1.2) translateX(10%);
            background: var(--danger-7);
            opacity: 0.8;
        }
    }

    .inner-category ul {
        list-style: none;
        padding: 0;
        margin: 0;
    }
</style>