<script>
    import { getProxyMediaUrl } from "@libs/HttpRequests";
    import { ChanThreadReply } from "@models/4Chan";
    import { createEventDispatcher } from "svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * @type {ChanThreadReply} 
        */
        export let thread_reply;

        // /**
        //  * @type {string} the id of the op's post
        //  */
        // export let op_post_id;
    
        const dispatcher = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        const emitSelectMediaPreview = () => {
            const media_file_url = thread_reply.file;
            dispatcher("select-media-preview", { media_file_url })
        }   
    
    /*=====  End of Methods  ======*/
    
</script>

<li class="thread-reply-wrapper {!thread_reply.hasImages() ? 'no-image-post' : ''}" id="p{thread_reply.uuid}">
    <div class="trw-reply-media image-interact" style:visibility={thread_reply.hasImages() ? "visible" : "hidden"}>
        {#if thread_reply.hasImages()}
            <img on:click={emitSelectMediaPreview} src="{getProxyMediaUrl(thread_reply.thumbnail_url)}" alt="thumbnail {thread_reply.uuid}" loading="lazy">
        {/if}        
    </div>
    <div class="trw-reply-content">
        <div class="trw-reply-status">
            <span class="trw-reply-post-id">
                <strong>post id: <i>{thread_reply.uuid.replace('pc', '')}</i></strong>
            </span>
            <span class="trw-rs-date">
                {thread_reply.getCreationDate()}
            </span>
        </div>
        <div class="trw-reply-message">
            <p class="trw-reply-message-content">
                {@html thread_reply.message}
            </p>
        </div>
    </div>
</li>


<style>
    .thread-reply-wrapper {
        display: grid;
        grid-template-columns: 1fr 3fr;
        background: var(--grey-9);
        border-radius: var(--border-radius);
    }

    .trw-reply-media {
        display: flex;
        justify-content: center;
        align-items: center;
        padding: var(--vspacing-2);
    }

    .trw-reply-media img {
        max-width: 100%;    
        max-height: 100%;
        border-radius: var(--border-radius);
        object-fit: contain;
    }

    .trw-reply-content {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-items: flex-start;
        row-gap: var(--vspacing-2);
        padding: var(--vspacing-3);
    }

    .no-image-post .trw-reply-content {
        grid-column: 1 / -1;
    }

    .trw-reply-message-content {
        row-gap: var(--vspacing-1);
    }

    .trw-reply-status {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .trw-reply-post-id {
        font-family: var(--font-decorative);
        font-size: var(--font-size-1);
        color: var(--grey-2);
    }

    .trw-reply-post-id i {
        color: var(--main-dark-color-4);
    }

    .trw-rs-date {
        font-family: var(--font-decorative);
        font-size: var(--font-size-1);
        color: var(--grey-6);
    }



    :global(.trw-reply-message-content a) {
        box-sizing: content-box;
        pointer-events: none;
        text-decoration: none;
        color: var(--main-dark);
        line-height: 2;
    }

    @container thread-content-component (max-width: 450px) {
        .thread-reply-wrapper {
            grid-template-columns: 1fr;
            grid-auto-rows: max-content;
        }

        .trw-reply-media {
            width: 100%;
            height: 25cqh;
        }

        .trw-reply-media img {
            max-width: 100%;    
            height: 100%;
            border-radius: var(--border-radius);
            object-fit: contain;
        }

        .no-image-post .trw-reply-media {
            display: none;
        }   

        .no-image-post .trw-reply-content{
            grid-row: 1 / span 2;
        }

    }

</style>