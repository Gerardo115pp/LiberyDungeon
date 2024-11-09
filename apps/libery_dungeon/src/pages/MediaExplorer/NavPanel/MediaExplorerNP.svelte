<script>
    import { layout_properties } from "@stores/layout";
    import { category_creation_tool_mounted, category_search_focused, media_upload_tool_mounted, category_search_results, category_search_term } from "../app_page_store";
    import { current_category, navigateToParentCategory, categories_tree } from "@stores/categories_tree";
    import CategorySearchBar from "@components/CategorySearchBar/CategorySearchBar.svelte";
    import { InnerCategory } from "@models/Categories";
    import { confirmPlatformMessage } from "@libs/LiberyFeedback/lf_utils";

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        const handleDeleteCategoryContent = async () => {
            if ($current_category === null || $current_category === undefined || $current_category.hasInnerCategories()) return;

            // const delete_category = confirm(`Are you sure you want to delete the category ${$current_category.name} and all its content? like, SUPEEER sure?`);
            const delete_category_choice = await confirmPlatformMessage({
                message_title: `Delete category ${$current_category.name}`,
                question_message: `Are you sure you want to delete the category ${$current_category.name} and it's ${$current_category.content.length} media files?`,
                cancel_label: "keep it",
                confirm_label: "delete it",
                danger_level: 2,
            })

            if (delete_category_choice === 1) {
                const category_to_delete = $current_category.uuid;
                await navigateToParentCategory();
                $categories_tree.deleteChildCategory(category_to_delete, true);
            }
                
        } 


        /**
         * Handle the search results from the category search bar
         * @param {CustomEvent<{ results: import('@models/Categories').Category[], search_query: string }>} event
         * @returns {void}
         */
        const handleCategoriesResults = (event) => {
            const search_results = event.detail.results;
            const search_query = event.detail.search_query;

            console.log("Search results: ", search_results);

            if (search_results == null || search_results.length === 0) return;

            category_search_focused.set(false);
            category_search_results.set(search_results);
            category_search_term.set(search_query);
        }
    /*=====  End of Methods  ======*/
    
    
</script>

<ul class="page-nav-menu" id="media-explorer-nav-menu-wrapper">
    <li id="menmw-upload-medias">
        <button on:click={() => media_upload_tool_mounted.set(!$media_upload_tool_mounted)} id="upload-medias-btn" class="sketch-btn">
            Upload
        </button>
    </li>
    {#if $current_category !== null && $current_category.hasContent() && !$current_category.hasInnerCategories()}
        <li id="menmw-upload-medias">
            <button on:click={handleDeleteCategoryContent} id="delete-content-btn" class="sketch-btn">
                Delete content
            </button>
        </li>
    {/if}
    {#if layout_properties.IS_MOBILE}
         <li id="menmw-new-category">
             <button on:click={() => category_creation_tool_mounted.set(!$category_creation_tool_mounted)} id="new-category-btn" class="sketch-btn">
                 New category
             </button>
         </li>
    {/if}
    {#if $category_search_focused}
        <li id="menmw-search-categories">
            <CategorySearchBar autofocus on:search-results={handleCategoriesResults}/>
        </li>
    {/if}

</ul>

<style>
    #media-explorer-nav-menu-wrapper {
        gap: var(--vspacing-1);    
    }
</style>