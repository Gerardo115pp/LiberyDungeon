<script>
    import { categories_tree, yanked_category } from "@stores/categories_tree";
    import { CategoryLeaf, CategoriesTree, moveCategory, renameCategory } from "@models/Categories";
    import { layout_properties } from "@stores/layout";
    import { createEventDispatcher } from "svelte";
    import { browser } from "$app/environment";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {CategoryLeaf} */
        export let category_leaf = {};

        /** @type {import('@models/Categories').InnerCategory} */
        export let inner_category = {};

        /**
         * Whether the category is currently focused by the keyboard selection
         * @type {boolean}
         */
        export let category_keyboard_focused = false;
      
        /*----------  Category Editing  ----------*/
        
            /**
             * Whether the category is been renamed
             * @type {boolean}
             */
            let category_renaming = false;

            /**
             * Category drag state
             * @type {boolean}
             */
            let category_dragging = false;

            /**
             * Whether another category is been dragged over this category.
             * @type {boolean}
            */
            let dragged_category_hovering = false;

            /** @type {HTMLLIElement} */
            let category_element;

            const dispatch = createEventDispatcher();

            $: if (category_keyboard_focused && browser) {
                ensureCategoryFocuseVisibility();
            }
    
    
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Renames the current category
         * @param {string} new_name
         */
        const applyCategoryRename = async new_name => {
            if (category_leaf == null) return; // Refuse to rename a category with out a leaf. Because we can't cause a reactivity update.
            
            let renamed_successfully = await renameCategory(category_leaf.uuid, new_name);
            if (!renamed_successfully) {
                alert("Failed to rename category. Repeated name?");
                return;
            } else {
                $categories_tree.updateCurrentCategory();
    
                category_leaf.name = new_name;            
            }

            category_renaming = false;
        }
    
        const handleCategoryClick = () => {
            return emitCategorySelectedEvent();
        }

        const handleCategoryRenameRequested = () => {
            category_renaming = (category_leaf != null); // Refuse to rename a category with out a leaf. Because we can't cause a reactivity update.
        }

        
        /*=============================================
        =            Drag Handlers            =
        =============================================*/
            // The process to implement Drag features is a bit strict and in some cases and requires a lot of
            // arbitrary steps to make it work. See: https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/Drag_operations#the_draggable_attribute

            /**
             * Handles the drag start event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragStart = e => {
                console.debug(`${category_leaf.name}: drag start`);
                category_dragging = true;
                e.dataTransfer.effectAllowed = "move";
                e.dataTransfer.setData("text/plain", category_leaf.uuid);
            }

            /**
             * Handles the drag end event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragEnd = e => {
                console.debug(`${category_leaf.name}: drag end`);   
                category_dragging = false;
                console.debug(`drag effect: ${e.dataTransfer.dropEffect}`);
            }

            /**
             * handles the drag enter event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragEnter = e => {
                console.debug(`${category_leaf.name}: drag enter`);
                const dragged_category_uuid = e.dataTransfer.getData("text/plain");

                if (!category_dragging && category_leaf.uuid !== dragged_category_uuid && category_leaf != null) {
                    e.preventDefault();
                    dragged_category_hovering = true;
                }
            }

            /**
             * handles the drag leave event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragLeave = e => {
                console.debug(`${category_leaf.name}: drag leave`);
                
                dragged_category_hovering = false;
            }

            /**
             * Handles the drag over event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragOver = e => {
                // Calling the preventDefault() method during both the dragenter and dragover event will indicate that a drop is allowed at that location.
                // See: https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/Drag_operations#specifying_drop_targets

                const dragged_category_uuid = e.dataTransfer.getData("text/plain");
                
                if (dragged_category_uuid !== category_leaf.uuid && !category_dragging && category_leaf != null) {
                    e.preventDefault(); // Accept the drop
                }

                
                e.dataTransfer.dropEffect = "move";
                return true;
            }

            /**
             * Handles the drop event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDrop = async e => {
                e.stopPropagation();
                e.preventDefault();

                dragged_category_hovering = false;
                category_dragging = false;

                const dragged_category_uuid = e.dataTransfer.getData("text/plain"); 
                
                console.debug(`${category_leaf.uuid} <- '${dragged_category_uuid}'`);

                if (dragged_category_uuid === category_leaf.uuid) {
                    throw new Error("Cannot drop a category on itself");
                }

                let updated_category = await moveCategory(dragged_category_uuid, category_leaf.uuid);

                if (updated_category != null && updated_category.UUID === dragged_category_uuid) {
                    $categories_tree.updateCurrentCategory();
                } else {
                    alert("Failed to move category");
                }
            }
        
        
        
        /*=====  End of Drag Handlers  ======*/
        
        /**
         * Handle the Keyboard event of the rename input
         * @param {KeyboardEvent} e
         */
        const handleRenameInput = e => {
            if (e.target instanceof HTMLInputElement && e.key === "Enter") {
                applyCategoryRename(e.target.value);
            }

            if (e.key === "Escape") {
                category_renaming = false;
            }
        }

        const emitCategorySelectedEvent = () => {
            dispatch("category-selected", {
                category: category_leaf,
                inner_category,
            });
        }

        function ensureCategoryFocuseVisibility() {
            if (!category_keyboard_focused || category_element === undefined) return;

            // TODO: add a check to see if the category is already visible

            category_element.scrollIntoView({
                behavior: "smooth",
                block: "center",
                inline: "center",
            });
        }

    /*=====  End of Methods  ======*/
    
</script>

<li class="ce-inner-category"
    bind:this={category_element} 
    class:keyboard-focused={category_keyboard_focused} 
    class:debug={false} 
    draggable="true"
    class:dragging={category_dragging}
    class:catergory-drop-target={dragged_category_hovering}
    class:yanked-category={$yanked_category === category_leaf?.uuid && $yanked_category !== ""}
    on:click={handleCategoryClick}
    on:dragstart={handleCategoryDragStart}
    on:dragend={handleCategoryDragEnd}
    on:dragenter={handleCategoryDragEnter}
    on:dragover={handleCategoryDragOver}
    on:dragleave={handleCategoryDragLeave}
    on:drop={handleCategoryDrop}
    on:rename-requested={handleCategoryRenameRequested}
>
    <svg class="ce-ic-icon" viewBox="0 0 110 80">
        <path class="category-vector-top" d="M55 10L95 30L55 50L15 30Z" />
        <path class="category-vector-layer" d="M15 40L55 60L95 40" />
        <path class="category-vector-layer" d="M15 47L55 67L95 47" />
        <path class="category-vector-layer" d="M15 54L55 74L95 54" />
    </svg>
    <div class="ce-ic-label">
        {#if !category_renaming}
            <h3>{category_leaf?.name ?? inner_category.name}</h3>
        {:else}
            <input type="text" id="ce-ic-rename-input" on:keyup={handleRenameInput} autofocus value="{category_leaf?.name ?? inner_category.name}"/>
        {/if}
    </div>
</li>

<style>
    li.ce-inner-category {
        width: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        transition: all .3s ease-in;
        padding: var(--vspacing-1);
        border: .5px solid var(--grey);
        transition: opacity .5s ease-out, border .5s ease-out, transform .2s linear;
        background: var(--grey);
        gap: var(--spacing-1);

        &.yanked-category {
            opacity: 0.2 !important;
        }
    }

    li.keyboard-focused {
        border-color: var(--main);
    }

    @media(pointer: fine) {
        .ce-inner-category:not(.dragging):hover {
            /* background: var(--grey-8); */
            border-color: var(--grey-8);
            backdrop-filter: brightness(1.3);
        }
    }

    svg.ce-ic-icon {
        width: 100%;
        pointer-events: none;
    }

    svg.ce-ic-icon path {
        stroke: var(--main-dark-color-6);
        fill: none;
        stroke-width: 1.2%;
        stroke-linecap: round;
        stroke-linejoin: round;
    }

    svg.ce-ic-icon path.category-vector-top {
        stroke: var(--main-dark-color-7);
        stroke-width: 1.6%;
    }

    
    /*----------  Category label  ----------*/
        
        .ce-ic-label {
            pointer-events: none;
        }

        .ce-ic-label h3{
            width: fit-content;
            font-family: var(--font-read);
            text-align: center;
            font-size: clamp(var(--font-size-1), 4%, var(--font-size-3));
        }

        .ce-ic-label input {
            font-family: var(--font-read);
            outline: none;
            background: transparent;    
            font-size: clamp(var(--font-size-1), 4%, var(--font-size-3));
            font-weight: 500;
            border: none;
            color: var(--main-4);
            width: 90%;
        }

    /*----------  Dragging state  ----------*/
        .ce-inner-category.dragging {
            opacity: 0.4;
        }

        :global(:has(> .ce-inner-category.dragging) .ce-inner-category:not(.dragging, .catergory-drop-target, :hover)) {
            scale: 1 0.98;
            opacity: 0.8;
            background-color: transparent;
            filter: grayscale(0.5);
        }
            

        .ce-inner-category.catergory-drop-target {
            border: 1px dotted var(--main);
        }



</style>

