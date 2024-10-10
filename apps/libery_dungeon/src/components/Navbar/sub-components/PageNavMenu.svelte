<script>
    import MediaExplorerNp from "@pages/MediaExplorer/NavPanel/MediaExplorerNP.svelte"
    import MediaViewerNp from "@pages/MediaViewer/NavPanel/MediaViewerNP.svelte";
    import { page } from "$app/stores";


    const route_names = {
        MEDIA_VIEWER: /\/media\-viewer\/[a-zA-Z\d]{40}(\/\d+)?/,
        MEDIA_EXPLORER: /\/dungeon\-explorer\/[a-zA-Z\d]{40}$/,
        DEFAULT: /^\//
    }
    
</script>

<div id="page-navmenu-wrapper">
    {#if route_names.MEDIA_VIEWER.test($page.url.pathname)}
        <MediaViewerNp />
    {:else if route_names.MEDIA_EXPLORER.test($page.url.pathname)}
        <MediaExplorerNp />
    {:else if route_names.DEFAULT.test($page.url.pathname)}
        <MediaExplorerNp />
    {/if}
</div>

<style>
    /* #page-navmenu-wrapper {
        font-size: var(--font-size-1);
    } */
</style>