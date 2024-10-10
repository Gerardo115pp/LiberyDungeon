<script>
    import MediaExplorer from "@pages/MediaExplorer/MediaExplorer.svelte";

    /*=============================================
    =            Properties            =
    =============================================*/
   
        /** 
         * @type {import('./+page').MediaExplorerParams}
         */
        export let data;


        /**
         * Whether to display the current category as a gallery or just as the normal MediasIcon.
         * In any case, the MediaIcon is also displayed.
         * @type {boolean}
         * @default false - change to true to modify the gallery UI.
         */
        let was_media_display_as_gallery = false;

        /**
         * Snapshot of the media explorer ephimeral state. for it to be restored navigation has to occur by navigating 'back'.
         * @see {@link https://kit.svelte.dev/docs/snapshots}
         * @type {import('./$types').Snapshot<Object>} */
        export const snapshot = {
            capture: () => {
                return {
                    was_media_display_as_gallery: was_media_display_as_gallery
                };
            }, 
            restore: (media_explorer_snapshot) => {
                console.log("Restoring media explorer snapshot: ", media_explorer_snapshot);
                was_media_display_as_gallery = media_explorer_snapshot.was_media_display_as_gallery ?? was_media_display_as_gallery;
            }
        }

        $: console.log("Was media display as gallery: ", was_media_display_as_gallery);
    
    
    /*=====  End of Properties  ======*/
    
    
</script>

<MediaExplorer 
    category_id={data.category_id}
    bind:was_media_display_as_gallery={was_media_display_as_gallery}
/>