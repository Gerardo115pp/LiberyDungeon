<script>
    import MediaExplorerNp from "@pages/MediaExplorer/NavPanel/MediaExplorerNP.svelte"
    import MediaViewerNp from "@pages/MediaViewer/NavPanel/MediaViewerNP.svelte";
    import { page } from "$app/stores";
    import { onMount, onDestroy } from "svelte";
    import { browser } from "$app/environment";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        const route_names = {
            MEDIA_VIEWER: /\/media\-viewer\/[a-zA-Z\d]{40}(\/\d+)?/,
            MEDIA_EXPLORER: /\/dungeon\-explorer\/[a-zA-Z\d]{40}$/,
            DEFAULT: /^\//
        }

        /**
         * @type {import('svelte').ComponentType<import('svelte').SvelteComponent> | null}
         */
        let current_component = null;   

        let page_change_unsubscriber = () => {};
    
    /*=====  End of Properties  ======*/
    
    


    onMount(() => {
        page_change_unsubscriber = page.subscribe(determinePageComponent);
    });

    onDestroy(() => {
        if (!browser) return;

        page_change_unsubscriber();
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * determines the correct page navbar component to use based on the current page route.
         * @returns {void}
         */  
        const determinePageComponent = () => {
            let new_page_nav_component = null;

            if (route_names.MEDIA_VIEWER.test($page.url.pathname)) {
                new_page_nav_component = MediaViewerNp;
            }

            if (route_names.MEDIA_EXPLORER.test($page.url.pathname)) {
                new_page_nav_component = MediaExplorerNp;
            }

            current_component = new_page_nav_component;
        }
        
        

    /*=====  End of Methods  ======*/
    
    
    
</script>

<div id="page-navmenu-wrapper">
    {#if current_component !== null}
        <svelte:component this={current_component} />
    {/if}
</div>

<style>
    /* #page-navmenu-wrapper {
        font-size: var(--font-size-1);
    } */
</style>