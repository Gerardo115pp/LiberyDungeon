<script>
    import { use_category_folder_thumbnails } from "@config/ui_design";
    import { categories_tree, yanked_category } from "@stores/categories_tree";
    import { moveCategory, renameCategory } from "@models/Categories";
    import { createEventDispatcher, onMount } from "svelte";
    import { browser } from "$app/environment";
    import { current_cluster } from "@stores/clusters";
    import { avoid_heavy_resources } from "@stores/layout";
    import CategoryConfiguration from "./sub-components/CategoryConfiguration.svelte";
    import { ensureElementVisible, isOutTheViewport } from "@libs/utils";

    
    /*=============================================
    =            Properties            =
    =============================================*/

        /** 
         * @type {import('@models/Categories').InnerCategory} 
         */
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

        /**
         * Whether to show the category configuration popover when the category is keyboard_focused.
         * @type {boolean}
         */
        export let show_category_configuration = false;

        /**
         * The image to use as a thumbnail for the category.
         * @type {string | null}
         */
        let category_thumbnail_url = null;

        /**
         * Whether the component has been mounted.
         * @type {boolean}
         */
        let component_mounted = false;
        
        /*----------  State  ----------*/
        
                /**
                 * If use_category_folder_thumbnails is enabled, this defines whether the category thumbnail url has been correctly loaded.
                 * @type {boolean}
                 */
                let category_thumbnail_loaded = false;
                $:if (inner_category.uuid && browser && component_mounted) {
                    determineCategoryThumbnailURL();
                }

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
                ensureCategoryFocusedVisibility();
            }
    
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        component_mounted = true;
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        /**
         * Renames the current category
         * @param {string} new_name
         */
        const applyCategoryRename = async (new_name) => {
            if ($categories_tree == null) {
                console.error("Categories tree store is not available");
                return;
            }
            
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

        /**
         * Determines the appropriate category thumbnail url to use.
         * @returns {Promise<void>}
         */
        const determineCategoryThumbnailURL = async () => {
            if (!$use_category_folder_thumbnails) return;

            if ($current_cluster == null) {
                throw new Error("In CategoryFolder.determineCategoryThumbnailURL: current_cluster is not available");
            }

            category_thumbnail_loaded = false;

            if (!inner_category.hasThumbnail()) {
                category_thumbnail_url = inner_category.getRandomMediaURL($current_cluster.UUID, 10600);
            }

            if (!inner_category.thumbnailIsLoaded()) {
                const thumbnail_loaded = await inner_category.loadThumbnail();

                if (!thumbnail_loaded) {
                    category_thumbnail_url = inner_category.getRandomMediaURL($current_cluster.UUID, 10600);
                    return;
                }
            }

            const media_thumbnail = inner_category.Thumbnail.Media;
            
            category_thumbnail_url = inner_category.Thumbnail.Media.isAnimatedImage() ? media_thumbnail.Url : media_thumbnail.getResizedUrl(25); // 25%  of the current viewport.
        }
    
        const handleCategoryClick = () => {
            return emitCategorySelectedEvent();
        }

        const handleCategoryRenameRequested = () => {
            category_renaming = (inner_category != null); // Refuse to rename a category with out a leaf. Because we can't cause a reactivity update.
        }

        /**
         * Handles the thumbnail changed event emitted by the CategoryConfiguration sub-component.
         * @type {import('./sub-components/category_folder_subs').CategoryConfig_ThumbnailChanged}
         */
        const handleThumbnailChanged = updated_category => {
            if (!inner_category.hasThumbnail() || !inner_category.loadThumbnail()) return;

            category_thumbnail_url = inner_category.Thumbnail.Media.isAnimatedImage() ? inner_category.Thumbnail.Media.Url : inner_category.Thumbnail.Media.getResizedUrl(25);
        
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
             * @param {DragEvent} event
             */
            const handleCategoryDrop = async (event) => {
                if ($categories_tree == null) {
                    console.error("Categories tree store is not available");
                    return;
                }
                
                if (is_ephemeral || !event.dataTransfer) return;
                
                event.stopPropagation();
                event.preventDefault();

                dragged_category_hovering = false;
                category_dragging = false;

                const dragged_category_uuid = event.dataTransfer.getData("text/plain"); 
                
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
            category_thumbnail_loaded = true;
        }

        const emitCategorySelectedEvent = () => {
            // TODO: Update subscribers to not expect 'inner_category'
            dispatch("category-selected", {
                category: inner_category,
                inner_category,
            });
        }

        function ensureCategoryFocusedVisibility() {
            if (!category_keyboard_focused || category_element === undefined) return;

            ensureElementVisible(category_element);
        }

    /*=====  End of Methods  ======*/
    
</script>

<div class="ce-inner-wrapper {category_item_class}">
    {#if component_mounted && category_keyboard_focused && show_category_configuration && !is_ephemeral}
        <div class="ce-inwra-configuration-popover">
            <CategoryConfiguration 
                the_inner_category={inner_category}
                on_thumbnail_change={handleThumbnailChanged}
            /> 
        </div>
    {/if}
    <li class="ce-inner-category"
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
                <div class="category-thumbnail-wrapper">
                    {#if category_thumbnail_url != null && !$avoid_heavy_resources}
                        <img 
                            fetchpriority="low" 
                            decoding="async" 
                            loading="lazy" 
                            src="{category_thumbnail_url}" 
                            alt="" 
                            aria-hidden="true" 
                            on:load={handleCategoryThumbnailLoaded} 
                        >
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

        position: relative;
        z-index: var(--z-index-1);
    }

    .ce-inner-wrapper:has(> .ce-inner-category.keyboard-focused) {
        z-index: var(--z-index-2);
    }

    .ce-inwra-configuration-popover {
        position: absolute;
        width: 150%;
        height: 110%;
        inset: -5% auto auto 50%;
        translate: -50% 0;
        z-index: var(--z-index-t-2);
        box-shadow: var(--shadow-1);
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
   
    /*=============================================
    =            CategoryThumbnail            =
    =============================================*/
    
        li.ce-inner-category.category-with-thumbnail {
            position: relative;
            background: transparent;
            border-radius: var(--border-radius);
            box-shadow: var(--shadow-2);
            padding-block: 0;
            transition: all var(--timing) ease-out, opacity calc(1.25 * var(--timing)) ease-out, border calc(2.25 * var(--timing)) ease-out, rotate var(--timing) ease, box-shadow var(--timing) ease-out;


            & .ce-ic-icon, & .ce-ic-label {
                z-index: var(--z-index-3);
            }

            & .ce-ic-label {
                position: absolute;
                bottom: -0.4em;
            }

            & .ce-ic-label h3 {
                font-family: var(--font-read);
                width: 100%;
                background: hsl(from var(--grey-7) h s l / 0.75);
                font-size: calc(var(--font-size-1) * 1.05);
                font-weight: 550;
                text-align: center;
                line-height: 1;
                padding: 0.5em 2ex;
                border-radius: var(--rounded-box-border-radius);
                color: var(--grey-2);
            }

            & .ce-ic-thumbnail-overlay {
                position: absolute;
                width: 100%;
                height: 100%;
                background-color: hsl(from var(--grey) h s l / 0.4);
                z-index: var(--z-index-2);
            }

            &.category-highlighted:not(.keyboard-focused) .ce-ic-thumbnail-overlay  {
                background-color: hsl(from var(--success) h s l / 0.4);
            }
            
            & .category-thumbnail-wrapper {
                position: absolute;
                width: 100cqw;
                height: 100cqh;
                z-index: var(--z-index-1);
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
                background-image: radial-gradient(circle, transparent 50cqw, var(--grey-black));
                opacity: 0;
                transition: opacity var(--timing);
                z-index: var(--z-index-4);
            }

            &::after {
                content: "";
                position: absolute;
                inset: 80% 0.5rem 0.5rem;
                translate: 0;
                transform: translateZ(-33cqw);
                background: black;
                filter: blur(1rem);
                opacity: 0;
                z-index: var(--z-index-3);
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
            rotate: x var(--rotation);

            & .ce-ic-label {
                inset: auto auto 2em auto;
                transform: translateY(-2em) translateZ(20px);
                transition: var(--timing);
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
                translate: 0 11.2560cqh;
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