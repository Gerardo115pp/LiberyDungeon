<script>
    import { InnerCategory } from "@models/Categories";
    import { getHotkeysManager } from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    
    /*=============================================
    =            Properties            =
    =============================================*/

        let global_hotkeys_manager = getHotkeysManager();
    
        /** @type {InnerCategory[]} the recently used categories */
        export let used_categories;


        /** @type{number} the index of the selected category */
        export let quick_selected_category_index = 0;
    
    
    /*=====  End of Properties  ======*/

</script>

<div id="quick-media-movement-tool">
    {#if used_categories.length !== 0}
        <ol id="qmmt-used-categories">
            {#each used_categories as uc, h}
                <li class="qmmt-uc-category-item" class:qmmt-uc-category-selected={quick_selected_category_index === h}>
                    {uc.name}
                </li>
            {/each}
        </ol>
    {:else}
        <p>No used categories</p>
    {/if}
</div>

<style>
    #quick-media-movement-tool {
        position: absolute;
        width: 30vw;
        max-height: 400px;
        top: 50%;
        left: -150%;
        padding: var(--vspacing-3);
        background: var(--grey-7);
        border-radius: var(--border-radius);
        transform: translate(-50%, -50%);
    }

    #qmmt-used-categories {
        display: flex;
        list-style: none;
        flex-wrap: wrap;
        flex-shrink: 0;
        padding: 0;
        margin: 0;
        height: 100%;
        overflow: auto;
    }

    .qmmt-uc-category-item {
        padding: var(--vspacing-1) var(--vspacing-2);
        background: var(--grey-5);
        color: var(--grey-1);
        border-radius: var(--border-radius);
        margin: var(--vspacing-1);
    }

    .qmmt-uc-category-selected {
        color: var(--main);
        border: 1px solid var(--main);
    }

    #quick-media-movement-tool p {
        text-align: center;
    }
</style>