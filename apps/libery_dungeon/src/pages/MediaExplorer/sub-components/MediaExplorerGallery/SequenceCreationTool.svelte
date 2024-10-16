<script>
    
    /*=============================================
    =            Properties            =
    =============================================*/

    import { me_gallery_changes_manager } from './me_gallery_state';
    import MeGalleryDisplayItem from './MEGalleryDisplayItem.svelte';

    
        /**
         * The media items from which the sequence will be created.
         * @type {import('@models/Medias').OrderedMedia[]}
         */ 
        export let unsequenced_medias = [];

        /**
         * Medias per grid row.
         * @type {number}
         */
        let medias_per_row = 10;
        
        /*----------  State  ----------*/
        
            /**
             * Whether to show the sequence as skeleton items.
             * @type {boolean}
             */ 
            let skeleton_sequence = true;

            /**
             * Media focused index.
             * @type {number}
             */
            let media_focus_index = 0;
    
    /*=====  End of Properties  ======*/
    
</script>

{#if $me_gallery_changes_manager != null}
    <div id="sequence-creation-tool">
        <div id="sequence-parameters"></div>
        <ul id="sequence-members"
            style:grid-template-columns="repeat({medias_per_row}, minmax(300px, 1fr))"
        >
            {#each unsequenced_medias as umedia, h}
                <li class="sct-sm-member-item"
                    class:is-skeleton={skeleton_sequence}
                    class:keyboard-selected={h === 0}
                >
                    <MeGalleryDisplayItem
                        ordered_media={umedia}
                        enable_video_titles
                        use_masonry={false}
                    />
                </li>
            {/each}
        </ul>
    </div>
{/if}

<style>
    ul#sequence-members {
        display: grid;
        background: var(--grey);
        gap: var(--spacing-2);
        padding: 4px;
        list-style: none;
        margin: 0;
    }

    li.sct-sm-member-item {
        position: relative;
        cursor: pointer;
        container-type: inline-size;
        background: var(--grey-9);
        width: 100%;
        height: 100%;
        z-index: var(--z-index-1);
        overflow: hidden;
    }

    li.sct-sm-member-item.is-skeleton {
        height: 340px !important;
    }

    li.sct-sm-member-item.keyboard-selected {
        outline: var(--main) solid 2px;
        z-index: var(--z-index-2);
    }
</style>