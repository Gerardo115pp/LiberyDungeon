<script>
    import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
    import { LabeledError, VariableEnvironmentContextError } from '@libs/LiberyFeedback/lf_models';
    import { emitPlatformMessage } from '@libs/LiberyFeedback/lf_utils';
    import { getCategoryByPath } from '@models/Categories';
    import { current_category } from '@stores/categories_tree';
    import { current_cluster } from '@stores/clusters';
    import { createEventDispatcher } from 'svelte';
    import { slide } from 'svelte/transition';

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The category to create breadcrumbs for
         * @type {import('@models/Categories').CategoryLeaf}
        */
        export let category_leaf_item;

        /**
         * The category path fragments with which to create the breadcrumbs
         * @type {BreadcrumbFragment[]}
         * @typedef {Object} BreadcrumbFragment
         * @property {string} fragment_category_name
         * @property {string} fragment_path
        */
        let category_path_fragments = [];
        $: category_path_fragments = generateCategoryPathFragments(category_leaf_item);

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Auto of a list of path fragments, creates a BreadcrumbFragment.
         * @param {string[]} path_fragments
         * @returns {BreadcrumbFragment | null}
         */
        function createBreadcrumbFragment(path_fragments) {
            if (path_fragments.length === 0) {
                return null;
            }

            let fragment_category_name = path_fragments[path_fragments.length - 1];
            let fragment_path = path_fragments.join('/');

            if (!fragment_path.endsWith('/')) {
                fragment_path += '/';
            }
            
            return {
                fragment_category_name,
                fragment_path
            };
        }

        /**
         * Returns a color lightness modifier for an element index
         * @param {number} element_index
         * @returns {number}
         */
        const getElementColorLightnessMod = (element_index) => {
            const step = 0.5;
            const min_lightness_mod = 0;
            const max_lightness_mod = 2;

            return 1 + Math.min(max_lightness_mod, min_lightness_mod + (element_index * step));
        }

        /**
         * Parses a given category's path and returns an array of breadcrumb fragments. which are objects
         * containing the path fragment and the the entire path leading up to that fragment.
         * @param {import('@models/Categories').CategoryLeaf} category
         * @returns {BreadcrumbFragment[]}
         */ 
        function generateCategoryPathFragments(category) {
            let path_fragments = category.FullPath.split('/');
            path_fragments = path_fragments.filter(fragment => fragment !== '');
            
            let breadcrumb_fragments = [];

            let root_category_fragment = {
                fragment_category_name: "main",
                fragment_path: "/"
            };

            for (let h = 0; h < path_fragments.length; h++) {
                let current_fragment_list = path_fragments.slice(0, h + 1);

                let breadcrumb_fragment = createBreadcrumbFragment(current_fragment_list);

                if (breadcrumb_fragment) {
                    breadcrumb_fragments.push(breadcrumb_fragment);
                }
            }

            breadcrumb_fragments = [root_category_fragment, ...breadcrumb_fragments];

            return breadcrumb_fragments;
        }

        /**
         * Handles the click event on a breadcrumb fragment.
         * @param {MouseEvent} event
         */
        const handleBreadcrumbFragmentClick = async event => {
            if (!event.currentTarget || !(event.currentTarget instanceof HTMLElement)) return;

            if ($current_category == null) {
                console.error("In CategoryNavigationBreadcrumbs.handleBreadcrumbFragmentClick: $current_category is null");
                return;
            }
            
            let breadcrumb_path = event.currentTarget.dataset.breadcrumbPath;
            if (breadcrumb_path === $current_category.FullPath) {
                return;
            }

            if (breadcrumb_path == null) {
                console.error("In CategoryNavigationBreadcrumbs.handleBreadcrumbFragmentClick: breadcrumb_path is null");
                return;
            }

            let ancestor_category = await getCategoryByPath(breadcrumb_path, $current_cluster.UUID);
            if (ancestor_category == null) {
                let variable_environment_error = new VariableEnvironmentContextError("In CategoryNavigationBreadcrumbs/handleBreadcrumbFragmentClick");
                variable_environment_error.addVariable("breadcrumb_path", breadcrumb_path);

                let labeled_err = new LabeledError(variable_environment_error, `Sorry, there is been an error getting in '${breadcrumb_path}'`, lf_errors.ERR_PROCESSING_ERROR);

                labeled_err.alert();
                return;
            }

            dispatch('breadcrumb-selected', ancestor_category);
        }   
    
    /*=====  End of Methods  ======*/
    
</script>

<menu class="category-navigation-breadcrumbs" role="menubar">
    {#each category_path_fragments as fragment, h}        
        {@const is_last_fragment = h === category_path_fragments.length - 1}
        {@const is_penultimate_fragment = h === category_path_fragments.length - 2}
        <li 
            class:current-fragment={is_last_fragment}
            class:is-penultimate-fragment={is_penultimate_fragment}
            class="cnb-fragment" 
            role="menuitem"
            style:--fragment-color-lightness-mod={!is_last_fragment ? getElementColorLightnessMod(h) : 1}
            style:--next-fragment-color-lightness-mod={!is_last_fragment ? getElementColorLightnessMod(h + 1) : 1}
            data-breadcrumb-path={fragment.fragment_path}
            in:slide={{axis: 'x', duration: is_last_fragment ? 300 : 0}}
            on:click={handleBreadcrumbFragmentClick}
        >
            <button class="cnb-fragment-navigation-btn">
                {fragment.fragment_category_name}
            </button>
        </li>
    {/each}
</menu>

<style>
    menu.category-navigation-breadcrumbs {
        --active-breadcrumb-color: var(--main-dark-color-7);    
        /* --menu-border-style: 1.25px solid var(--grey-9); */
        --menu-border-style: none;
     
        display: flex;
        gap: var(--spacing-1);
        border-top: var(--menu-border-style);
        border-bottom: var(--menu-border-style);
        height: min(3.1dvh, 30px);
        overflow-y: hidden;
    }

    .cnb-fragment {
        height: 100%;
        background-color: var(--grey-1);
        padding: 0 var(--spacing-3);
        line-height: 1;
        border-radius: 0;
        
        &.current-fragment {
            background-color: var(--active-breadcrumb-color);
        }

        & > button.cnb-fragment-navigation-btn {
            display: grid;
            height: 100%;
            font-family: var(--font-decorative);
            font-size: calc(var(--font-size-fineprint) * 1.2);
            padding: 0 var(--spacing-3);
            place-items: center;
            line-height: 1;
            vertical-align: middle;
            outline: none;
        }
    }
    
    @supports (color: rgb(from white r g b)) {
        /* 
            Enable the stepped color gradient for browsers that support relative color syntax.
        */
        menu.category-navigation-breadcrumbs {
            --active-breadcrumb-color: hsl(from var(--main-dark-color-7) h s calc(l * 1.1));
            gap: 0;
        }

        .cnb-fragment {
            background-color: hsl(from var(--grey) h s calc(l * var(--next-fragment-color-lightness-mod)));
            padding: 0;
            border-radius: 0;
        }

        .cnb-fragment.is-penultimate-fragment {
            background-color: var(--active-breadcrumb-color);
        }
        
        .cnb-fragment > button.cnb-fragment-navigation-btn {
            border-radius: 0 var(--border-radius) var(--border-radius) 0; 
            background-color: hsl(from var(--grey) h s calc(l * var(--fragment-color-lightness-mod)));
        }

        .cnb-fragment.current-fragment {
            border-radius: 0 var(--border-radius) var(--border-radius) 0;
        }

        .cnb-fragment.current-fragment > button.cnb-fragment-navigation-btn {
            background-color: var(--active-breadcrumb-color);
        }

        .cnb-fragment:not(.current-fragment):hover > button.cnb-fragment-navigation-btn {
            background-color: hsl(from var(--grey-8) h s calc(l * var(--fragment-color-lightness-mod)));
            transition: background-color .2s ease-out;
        }

        .cnb-fragment:not(.is-penultimate-fragment):has(+ .cnb-fragment:hover) {
            background-color: hsl(from var(--grey-8) h s calc(l * var(--next-fragment-color-lightness-mod)));
            transition: background-color .2s ease-out;
        }
    }
</style>