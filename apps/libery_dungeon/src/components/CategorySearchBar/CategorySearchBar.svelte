<script>
    import Input from "@components/Input/Input.svelte";
    import { createEventDispatcher } from "svelte";
    import FieldData from "@libs/FieldData";
    import { searchCategories } from "@models/WorkManagers";
    import { onMount } from "svelte";
    import { current_category } from "@stores/categories_tree";
    import { current_cluster } from "@stores/clusters";
    import { browser } from "$app/environment";

    
    /*=============================================
    =            Properties            =
    =============================================*/

        /** @type {boolean} whether the search bar is focused or not */
        export let autofocus = false;    

        const search_field_data = new FieldData("category-search-bar-input", /.+/, "Category search");

        const search_event_dispatcher = createEventDispatcher();
    
    
    onMount(() => {
        if (autofocus && browser) {
            setTimeout(() => {
                focusSearchBar();
            }, 400);
        }
    });

    
    /*=============================================
    =            Methods            =
    =============================================*/

        const handleSearch = async () => {
            if ($current_category === null) {
                throw new Error("On components/CategorySearchBar.svelte: Attempted to search categories without a current category. which is required to retrieve the cluster id");
            }

            const search_bar = search_field_data.getField();

            search_bar.blur();

            const search_query = search_field_data.getFieldValue();

            if (search_query === "") {
                return;
            }

            const response = await searchCategories(search_query, $current_category.ClusterUUID, $current_cluster.DownloadCategoryID);

            const search_results = response.data;

            search_event_dispatcher("search-results", {
                results: search_results,
                search_query,
            });

            search_field_data.clear();
        }

        export const focusSearchBar = () => {
            let search_bar = search_field_data.getField();

            if (search_bar === null) return;

            const search_bar_style =  window.getComputedStyle(search_bar);

            if (search_bar_style.visibility === "visible") {
                search_bar.focus();
            }
        }

        const handleSearchBarCommands = e => {
            const search_bar = search_field_data.getField();

            if (e.key === "Escape") {
                e.preventDefault();
                search_bar.blur();
            }

            if (e.key === "e" && e.ctrlKey) {
                e.preventDefault();
                handleSearch();
            }
        }
    
    /*=====  End of Methods  ======*/

</script>

<dialog open class="category-search-bar-wrapper">
    <Input
        input_background="var(--grey)"
        placeholder_color="var(--main-1)"
        input_padding="calc(var(--vspacing-1) * .5) var(--vspacing-1)"
        input_color="var(--main-5)"
        border_color="var(--main)"
        field_data={search_field_data}
        onKeypress={handleSearchBarCommands}
        onEnterPressed={handleSearch}
        isSquared={true}
    />
</dialog>

<style>
    dialog.category-search-bar-wrapper {
        position: static;
        background-color: transparent;
        border: none;
        padding: 0;
    }
</style>