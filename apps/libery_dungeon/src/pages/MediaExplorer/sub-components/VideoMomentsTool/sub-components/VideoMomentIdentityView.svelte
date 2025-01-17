<script>
    import { navigateToMediaViewer } from '@libs/NavigationUtils';

    /*=============================================
    =            Propertie            =
    =============================================*/
    
        /**
         * @type {import('@models/Metadata').VideoMomentIdentity}
         */ 
        export let the_video_moment_identity;
    
    /*=====  End of Propertie  ======*/

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the click event on the video moment item.
         * @param {MouseEvent} event
         */ 
        const handleVideoMomentClick = event => {

            navigateToMediaViewer(the_video_moment_identity.MediaIdentity.CategoryUUID, {
                media_uuid: the_video_moment_identity.MediaIdentity.Media.uuid,
                time: the_video_moment_identity.VideoMoment.DecodedTime
            });
        }
    
    /*=====  End of Methods  ======*/
    
    
</script>

<!-- svelte-ignore a11y-interactive-supports-focus -->
<div class="diviext-moment-identity" 
    role="button"
    on:click={handleVideoMomentClick}
>
    <p>
        <strong>{the_video_moment_identity.MediaIdentity.Media.MediaName}</strong> ―
        {the_video_moment_identity.VideoMoment.Title}
    </p>
</div>

<style>
    .diviext-moment-identity {
        color: var(--grey-1);
        background: hsl(from var(--grey-9) h s l / 0.8);
        padding-inline: 3em;
        padding-block: .75em;
        transition: background 0.3s ease-out;

        & > p {
            font-size: var(--font-size-1);
        }

        & > p strong {
            color: var(--main);
        }

        &:hover {
            cursor: pointer;
            background: hsl(from var(--grey-7) h s l / 0.9);
        }
    }
</style>

