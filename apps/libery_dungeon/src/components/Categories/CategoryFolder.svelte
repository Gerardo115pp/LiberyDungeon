<script>
    import { use_category_folder_thumbnails } from "@config/ui_design";
    import { categories_tree, yanked_category } from "@stores/categories_tree";
    import { moveCategory, renameCategory } from "@models/Categories";
    import { createEventDispatcher } from "svelte";
    import { browser } from "$app/environment";
    import { current_cluster } from "@stores/clusters";
    import viewport from "@components/viewport_actions/useViewportActions";

    
    /*=============================================
    =            Properties            =
    =============================================*/

        /** @type {import('@models/Categories').InnerCategory} */
        export let inner_category;

        /**
         * Whether the category icon is ephemeral. for example if it's part for a category search results that are not in the current category
         * @type {boolean}
         */
        export let is_ephemeral = false;

        /**
         * Whether the category folder should be highlighted
         * @type {boolean}
         * @default false
         */
        export let highlight_category = false;

        /**
         * Whether the category is currently focused by the keyboard selection
         * @type {boolean}
         */
        export let category_keyboard_focused = false;
      
        
        /*----------  State  ----------*/
        
                /**
                 * If use_category_folder_thumbnails is enabled, this defines whether the category thumbnail url has been correctly loaded.
                 * @type {boolean}
                 */
                let category_thumbnail_loaded = false;
                $:if (inner_category.uuid) {
                    category_thumbnail_loaded = false;
                }

                /**
                 * Whether the category icon is currently visible.
                 * @type {boolean}
                 */
                let category_thumbnail_visible = false;


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

            /**
             * An additional html class(es) to add to the category element
             * @type {string}
             */
            export let category_item_class = "";

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
            if (is_ephemeral) return; // Refuse to rename a category with out a leaf. Because we can't cause a reactivity update.
            
            let renamed_successfully = await renameCategory(inner_category.uuid, new_name);
            if (!renamed_successfully) {
                alert("Failed to rename category. Repeated name?");
                return;
            } else {
                $categories_tree.updateCurrentCategory();
    
                inner_category.name = new_name;            
            }

            category_renaming = false;
        }
    
        const handleCategoryClick = () => {
            return emitCategorySelectedEvent();
        }

        const handleCategoryRenameRequested = () => {
            category_renaming = (inner_category != null); // Refuse to rename a category with out a leaf. Because we can't cause a reactivity update.
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
                if (is_ephemeral) return;

                category_dragging = true;

                if (e.dataTransfer == null) return;

                e.dataTransfer.effectAllowed = "move";
                e.dataTransfer.setData("text/plain", inner_category.uuid);
            }

            /**
             * Handles the drag end event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragEnd = e => {
                category_dragging = false;
            }

            /**
             * handles the drag enter event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragEnter = e => {
                if (is_ephemeral) return;

                const dragged_category_uuid = e.dataTransfer?.getData("text/plain");

                if (!category_dragging && inner_category.uuid !== dragged_category_uuid) {
                    e.preventDefault();
                    dragged_category_hovering = true;
                }
            }

            /**
             * handles the drag leave event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragLeave = e => {
                if (is_ephemeral) return;
                
                dragged_category_hovering = false;
            }

            /**
             * Handles the drag over event of the category
             * @param {DragEvent} e
             */
            const handleCategoryDragOver = e => {
                // Calling the preventDefault() method during both the dragenter and dragover event will indicate that a drop is allowed at that location.
                // See: https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/Drag_operations#specifying_drop_targets
                if (is_ephemeral || !e.dataTransfer) return;


                const dragged_category_uuid = e.dataTransfer.getData("text/plain");
                
                if (dragged_category_uuid !== inner_category.uuid && !category_dragging) {
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
                if (is_ephemeral || !e.dataTransfer) return;
                
                e.stopPropagation();
                e.preventDefault();

                dragged_category_hovering = false;
                category_dragging = false;

                const dragged_category_uuid = e.dataTransfer.getData("text/plain"); 
                
                console.debug(`${inner_category.uuid} <- '${dragged_category_uuid}'`);

                if (dragged_category_uuid === inner_category.uuid) {
                    throw new Error("Cannot drop a category on itself");
                }

                let updated_category = await moveCategory(dragged_category_uuid, inner_category.uuid);

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

        /**
         * Handles the load event from the category thumbnail
         * @param {Event} event
         */
        const handleCategoryThumbnailLoaded = event => {
            console.log("Category thumbnail loaded:", event.target);
            category_thumbnail_loaded = true;
        }

        /**
         * Handles the viewportEnter event on the category thumbnail.
         * @type {import('@components/viewport_actions/useViewportActions').ViewportEventHandler}
         */
        const handlerCategoryThumbnailViewportEnter = (event) => {
            category_thumbnail_visible = true;
        }

        /**
         * Handles the viewportLeave event on the category thumbnail.
         * @type {import('@components/viewport_actions/useViewportActions').ViewportEventHandler}
         */
        const handlerCategoryThumbnailViewportLeave = (event) => {
            category_thumbnail_visible = false;
        }

        const emitCategorySelectedEvent = () => {
            // TODO: Update subscribers to not expect 'inner_category'
            dispatch("category-selected", {
                category: inner_category,
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

<div class="ce-inner-wrapper">
    <li class="ce-inner-category {category_item_class}"
        bind:this={category_element} 
        class:keyboard-focused={category_keyboard_focused} 
        class:debug={false} 
        draggable="true"
        class:dragging={category_dragging}
        class:catergory-drop-target={dragged_category_hovering}
        class:yanked-category={$yanked_category === inner_category?.uuid && $yanked_category !== ""}
        class:category-highlighted={highlight_category}
        class:category-with-thumbnail={$use_category_folder_thumbnails && inner_category}
        class:category-thumbnail-loaded={category_thumbnail_loaded}
        on:click={handleCategoryClick}
        on:dragstart={handleCategoryDragStart}
        on:dragend={handleCategoryDragEnd}
        on:dragenter={handleCategoryDragEnter}
        on:dragover={handleCategoryDragOver}
        on:dragleave={handleCategoryDragLeave}
        on:drop={handleCategoryDrop}
        on:rename-requested={handleCategoryRenameRequested}
    >
        {#key inner_category.uuid}
            {#if $use_category_folder_thumbnails && inner_category != null}
                <div class="category-thumbnail-wrapper" 
                    on:viewportEnter={handlerCategoryThumbnailViewportEnter} 
                    on:viewportLeave={handlerCategoryThumbnailViewportLeave} 
                    use:viewport
                >
                    {#if category_thumbnail_visible}
                        <img fetchpriority="low" decoding="async" loading="lazy" src="{inner_category.getRandomMediaURL($current_cluster.UUID, 10600)}" alt="" aria-hidden="true" on:load={handleCategoryThumbnailLoaded} >
                    {/if}
                </div>
                <div class="ce-ic-thumbnail-overlay"></div>
            {/if}
        {/key}
        <svg class="ce-ic-icon" viewBox="0 0 110 80">
            <path class="category-vector-top" d="M55 10L95 30L55 50L15 30Z" />
            <path class="category-vector-layer" d="M15 40L55 60L95 40" />
            <path class="category-vector-layer" d="M15 47L55 67L95 47" />
            <path class="category-vector-layer" d="M15 54L55 74L95 54" />
        </svg>
        <div class="ce-ic-label">
            {#if !category_renaming}
                <h3>{inner_category?.name}</h3>
            {:else}
                <input type="text" id="ce-ic-rename-input" on:keyup={handleRenameInput} autofocus value="{inner_category.name}"/>
            {/if}
        </div>
    </li>
</div>

<style>
    .ce-inner-wrapper {
        --timing: 400ms;
        --rotation: 20deg;
    }

    li.ce-inner-category {
        width: 100%;
        display: flex;
        height: 100%;
        background: var(--grey);
        container-type: size;
        flex-direction: column;
        align-items: center;
        padding: var(--spacing-1);
        border: .5px solid var(--grey);
        gap: var(--spacing-1);
        transition: all .3s ease-in, opacity .5s ease-out, border .2s ease-out, transform .2s linear;

        &.yanked-category {
            opacity: 0.2 !important;
        }
    }

    li.keyboard-focused {
        border-color: var(--main);
    }

    li.ce-inner-category.category-highlighted:not(.keyboard-focused) {
        background: hsl(from var(--main) h s l / 0.05);
    }

    /* @media(pointer: fine) {
        .ce-inner-category:not(.dragging):hover {
            
        }
    } */

    
    /*=============================================
    =            CategoryThumbnail            =
    =============================================*/
    
        li.ce-inner-category.category-with-thumbnail {
            position: relative;
            background: transparent;
            border-radius: var(--border-radius);
            box-shadow: var(--shadow-2);
            transition: all var(--timing) ease-out, opacity .5s ease-out, border 900ms ease-out, transform var(--timing) linear, box-shadow var(--timing) ease-out;


            & .ce-ic-icon, & .ce-ic-label {
                z-index: var(--z-index-1);
            }

            & .ce-ic-thumbnail-overlay {
                position: absolute;
                width: 100%;
                height: 100%;
                background-color: hsl(from var(--grey) h s l / 0.4);
                z-index: var(--z-index-b-1);
            }
            
            & .category-thumbnail-wrapper {
                position: absolute;
                width: 100cqw;
                height: 100cqh;
                z-index: var(--z-index-b-2);
                overflow: hidden;
            }

            & .category-thumbnail-wrapper img {
                opacity: 0;
                width: 100cqw;
                height: 100cqh;
                object-fit: cover;
                object-position: center;
            }

            &::before {
                content: "";
                position: absolute;
                inset: 0;
                z-index: 100;
                background-image: radial-gradient(circle, transparent 50cqw, black);
                opacity: 0;
                transition: opacity var(--timing);
            }

            &::after {
                content: "";
                position: absolute;
                inset: 80% 0.5rem 0.5rem;
                translate: 0;
                transform: translateZ(-100px);
                background: black;
                filter: blur(1rem);
                opacity: 0;
                z-index: 1;
                transition: rotate var(--timing), translate var(--timing)
            }
        }

        li.ce-inner-category.category-with-thumbnail.category-thumbnail-loaded {
            align-items: center;
            justify-content: flex-end;
            padding: 0;
            border-top: 1px solid var(--grey-3);
            
            & .ce-ic-icon {
                display: none;
            }

            & .ce-ic-label {
                position: absolute;
                bottom: -0.4em;

            }

            & .ce-ic-label h3 {
                font-family: var(--font-decorative);
                width: 100%;
                background: hsl(from var(--grey-7) h s l / 0.95);
                font-size: var(--font-size-1);
                text-align: center;
                line-height: 1;
                padding: 0.5em 2ex;
                border-radius: var(--rounded-box-border-radius);
                color: var(--grey-2);
            }

            & .category-thumbnail-wrapper img {
                
                opacity: 1;
            }
        }

        .ce-inner-wrapper:has(> li.ce-inner-category.category-with-thumbnail) {
            perspective: 1000px;
        }

        li.ce-inner-category.category-with-thumbnail.keyboard-focused {
            border-top: 1px solid var(--grey-3);
            border-left: none;
            border-right: none;
            border-bottom: none;
            box-shadow: var(--shadow-2-perspective);
            scale: 1.05;
            transform-style: preserve-3d;
            transform-origin: center;
            rotate: x 15deg;
            
            /* & .category-thumbnail-wrapper img {
            } */

            & .ce-ic-label {
                inset: auto 2em;
                transform: translateY(-2em) translateZ(20px);
                transition: 400ms;
            }

            & .ce-ic-label h3 {
                color: var(--grey);
                background: var(--main-dark);
            }            

            &::before {
                opacity: 1;
            }

            &::after {
                rotate: x calc(var(--rotation) * -1);
                translate: 0 60px;
                opacity: .8;
            }
        }


            
    
    /*=====  End of CategoryThumbnail  ======*/
    
    

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

