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
        import FieldData from "@libs/FieldData";
        import Input from "@components/Input/Input.svelte";
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
        
        /** @type { boolean } if true, show the create subcategory input */
        let show_create_subcategory_input = false;

        /** @type { FieldData } data for the create subcategory input */
        let create_subcategory_field = new FieldData("create-subcategory-input", /.+/, "New subcategory");

        const dispatch = createEventDispatcher();

        
        /*----------  Unsubscribers  ----------*/
        
        let new_category_name_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        // PASS A GOD DAMN PROP TO CHANGE THIS STATE THEN MAKE A REACTIVE STATEMENT! this down below is horrible.
        new_category_name_unsubscriber = create_subcategory.subscribe(value => {
            if (value) {
                handleCreateSubcategory();
            }
        });
    })

    onDestroy(() => {
        new_category_name_unsubscriber();
    })
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles new subcategory creation.
         * @param category_name
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

        const handleCreateSubcategory = () => {
            create_subcategory_field.clear();

            show_create_subcategory_input = is_keyboard_selected;            
        }    

        const handleCategoryItemSelected = e => {
            e.stopPropagation();

            let propagate_uuid = category.uuid;
            let propagate_name = category.name;
            let shift_key = e.shiftKey;

            if (!(e instanceof MouseEvent)) {
                propagate_name = e.detail.category.name;
                propagate_uuid = e.detail.category.uuid;
                shift_key = e.detail.shift_key;
            }

            dispatch("category-item-selected", {
                category: {
                    uuid: propagate_uuid,
                    name: propagate_name
                },
                shift_key: shift_key
            });
        }

        /**
         * @param { KeyboardEvent } e 
         */
        const handleCreateSubcategoryInputCommands = e => {
            if (e.key === "Enter" || (e.key === "e" && e.ctrlKey)) {
                e.preventDefault();

                const new_subcategory_name = create_subcategory_field.getFieldValue();

                if (create_subcategory === "") return;

                show_create_subcategory_input = false;
                create_subcategory.set(null); // 🤮🤮🤮🤮🤮

                createSubCategory(new_subcategory_name);
            }
        }
            

        const toggleInnerCategories = async () => {
            if (inner_categories.length === 0 && show_inner_categories) {
                inner_categories = await getShortCategoryTree(category.uuid);
            }

            is_expanded = !is_expanded;
        }
    
    /*=====  End of Methods  ======*/
    
</script>


<li id="inner-category-{category.uuid}" class="inner-category" class:keyboard-selected={is_keyboard_selected}>
    <div class="inner-category-wrapper" class:is-light={is_light} on:click={toggleInnerCategories}>
        {#if show_create_subcategory_input}
            <!-- 
                I deprecated the Input component in the previous century! DO NOT USE IT! the input element API now a day does everything this component was created for.
             -->
            <Input 
                field_data={create_subcategory_field}
                input_background="var(--grey)"
                placeholder_color="var(--main-1)"
                input_padding="calc(var(--vspacing-1) * .3) var(--vspacing-1)"
                input_color="var(--main-5)"
                font_size="var(--font-size-1)"
                border_color="var(--main)"
                onKeypress={handleCreateSubcategoryInputCommands}
                isSquared={true}
                autofocus={true}
            />
        {:else}
            <strong class="ic-category-name" on:click={handleCategoryItemSelected}>
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
                <InnerCategoryView category={ic} is_light={!is_light} on:category-item-selected={handleCategoryItemSelected} />
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