<script>
    import MediaExplorer from "@pages/MediaExplorer/MediaExplorer.svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
   
        /** 
         * @type {import('./+page').MediaExplorerParams}
         */
        export let data;
    
    /*=====  End of Properties  ======*/

</script>

<MediaExplorer 
    category_id={data.category_id}
/>