<script>
    import { current_category } from "@stores/categories_tree";
    import { CategoriesTree, getCategoryTree, CategoryLeaf } from "@models/Categories";
    import { writable } from "svelte/store";
    import CategoryFolder from "@components/Categories/CategoryFolder.svelte";
    import MediasIcon from "@components/Medias/MediasIcon.svelte";
    import LiberyHeadline from "@components/UI/LiberyHeadline.svelte";
    import { createEventDispatcher, onMount } from "svelte";
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {CategoriesTree} the categories tree for the gallery explorer */
        let local_categories_tree;

        /** @type {import('svelte/store').Writable<CategoryLeaf>} the local current category */
        let local_current_category = writable(null);

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        local_current_category.set($current_category);

        local_categories_tree = await getCategoryTree($current_category.uuid, local_current_category);
    });    

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        const handleGoToParent = () => {
            local_categories_tree.navigateToParent();
        }

        const handleCategorySelected = e => {
            let leaf_uuid = e.detail?.category?.uuid;

            if (leaf_uuid === undefined) return;

            local_categories_tree.navigateToLeaf(leaf_uuid);
        }

        const handleContentSelected = e => {
            dispatch("open-remote-gallery", {
                content: $local_current_category?.content ?? [],
                fullpath: $local_current_category?.FullPath ?? "",
                category_id: $local_current_category?.uuid ?? "",
                name: $local_current_category?.name ?? ""
            });
        }
    
    /*=====  End of Methods  ======*/
   
    
</script>

<aside id="medias-gallery-explorer">
    {#if $local_current_category !== null}
        <div id="mge-headline-wrapper" on:click={handleGoToParent}>
            <LiberyHeadline headline_text="{$local_current_category.name}"  headline_tag="h3" forced_font_size="var(--font-size-4)"/>
        </div>
        <ul id="mge-category-content" class="libery-scroll">
            {#each $local_current_category.InnerCategories as category}
                <CategoryFolder category_leaf={category} change_category_on_select={false} on:category-selected={handleCategorySelected} />
            {/each}
            {#if $local_current_category.hasContent()}
                <MediasIcon category_id={$local_current_category.uuid} images_count={$local_current_category.content.length} enter_media_viewer={false} on:content-selected={handleContentSelected}/>
            {/if}
        </ul>    
    {/if}
</aside>

<style>
    aside#medias-gallery-explorer {
        display: grid;
        height: 85cqh;
        grid-template-columns: 1fr;
        grid-template-rows: 20% 73%;
        row-gap: var(--vspacing-1);
        padding: var(--vspacing-2);
    }

    div#mge-headline-wrapper {
        cursor: default;
    }

    ul#mge-category-content {
        display: grid;
        grid-row: 2 / span 1;
        overflow-y: auto;
        grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
        list-style: none;
        padding: 0;
    }
</style>