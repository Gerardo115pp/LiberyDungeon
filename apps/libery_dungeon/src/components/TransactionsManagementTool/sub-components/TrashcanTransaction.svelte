<script>
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The transaction to display.
         * @type {import('@models/Trashcan').TrashcanTransaction}
         */ 
        export let the_transaction = null;

        /**
         * The keyboard focused media index.
         * @type {number}
         */
        export let focused_media_index = 0;
        $: ensureFocusedMediaVisible(focused_media_index);

        /**
         * Whether to highlight the focused media or not.
         * @type {boolean}
         */
        export let highlight_focused_media = false;
    
    /*=====  End of Properties  ======*/

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
    async function ensureFocusedMediaVisible(focused_media_index) {
        if (!highlight_focused_media) return;

        let focused_media = document.querySelector(`.trashed-media.trashed-media-entry-${focused_media_index}`);

        if (focused_media === null) return;

        console.log(focused_media);

        focused_media.scrollIntoView({
            behavior: 'smooth',
            block: 'center',
        });
    } 
    
    /*=====  End of Methods  ======*/
    
    
</script>

<div id="tmt-trashcan-transaction-wrapper">
    <menu id="trashed-medias-container">
        {#if the_transaction != null}
             {#each the_transaction.Content as media, k}
                <div class="trashed-media trashed-media-entry-{k}"
                    class:keyboard-focused-media={highlight_focused_media && k === focused_media_index}
                >
                    <span class="trashed-media-name">{media.Name}</span>
                    <img 
                        src="{media.ThumbnailURL}" 
                        alt="{media.Name}"
                        loading="lazy"
                    >
                </div> 
            {/each} 
        {/if}
    </menu>
</div>

<style>
    menu#trashed-medias-container {
        width: 100%;
        background: var(--grey);
        overflow-y: auto;
        height: 80cqh;
        scrollbar-width: none;
    }

    .trashed-media {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: var(--spacing-1) var(--spacing-3);

        &:not(:last-child) {
            border-bottom: 1px solid var(--grey-8);
        }

        &.keyboard-focused-media {
            background: var(--grey-8);
        }
    }
</style>