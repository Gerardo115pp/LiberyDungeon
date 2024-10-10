<script>
    import { browser } from "$app/environment";
    import MediaViewer from "@pages/MediaViewer/MediaViewer.svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** 
         * @type {import('./+page').MediaViewerQueryParams}
         */
        export let data;
    
    
    /*=====  End of Properties  ======*/
    
    
</script>

{#if browser}
    <MediaViewer 
        url_category_id={data.category_id}
        url_media_index={data.media_index}
    />
{/if}