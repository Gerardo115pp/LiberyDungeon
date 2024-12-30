<script>
    import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
    import { LabeledError } from '@libs/LiberyFeedback/lf_models';


    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The active media we would be sharing.
         * @type {import('@models/Medias').Media}
         */
        export let the_active_media;
    
    /*=====  End of Properties  ======*/
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Copies the sharable link of the media to the users clipboard.
         * @param {MouseEvent} event
         */
        const handleShareButtonClick = async event => {
            const shareable_link = await the_active_media.getSharedUrl()

            if (shareable_link === null) {
                const labeled_err = new LabeledError(
                    "In @pages/MediaViewer/sub-components/MediaInformation/sub-components/ShareMediaBTN.handleShareButtonClick",
                    "Failed to generate a shareable link for this media.",
                    lf_errors.ERR_PROCESSING_ERROR
                );

                labeled_err.alert();
                return;
            }

            navigator.clipboard.writeText(shareable_link);
        }

    
    /*=====  End of Methods  ======*/
    
    

</script>

<button 
    class="smb-share-media-btn dungeon-button-1"
    on:click={handleShareButtonClick}
>
    Share media
</button>

<style>
    button.smb-share-media-btn {
        line-height: 1;
        padding-block: 0.4em;
        padding-inline: 2em;
    }
</style>