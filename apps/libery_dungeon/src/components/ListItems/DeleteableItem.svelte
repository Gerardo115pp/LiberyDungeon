<script>
    import { createEventDispatcher } from "svelte";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /**
         * The item color.
         * @type {string}
         * @default "var(--grey-5)"
         */ 
        export let item_color = "var(--grey-5)";

        /**
         * Whether the element is currently protected from deletion.
         * @type {boolean}
         */
        export let is_protected = false;

        const dispatch = createEventDispatcher();
    
    /*=====  End of Properties  ======*/
    
   
   /*=============================================
   =            Methods            =
   =============================================*/
   
        /**
         * Emits the item-deleted event.
         */
        const emitDeleteItem = () => {
            dispatch("item-deleted");
        }

        /**
         * Emits the item-selected event.
         */
        const emitSelectItem = () => {
            dispatch("item-selected");
        }

        /**
         * Handles the click event on the delete button.
         */
        const handleDeleteButtonClick = () => {
            emitDeleteItem();
            return false;
        }

        /**
         * Handles the click event on the item.
         */
        const handleSelectClick = () => {
            emitSelectItem();
        }
   
   /*=====  End of Methods  ======*/
    
</script>

<li class="deletable-item"
    style:background="{item_color}"
    on:click={handleSelectClick}
>
    <div class="deletable-item-content">
        <slot/>
    </div>
    {#if !is_protected}
        <button class="delete-item"
            type="button"
            on:click|preventDefault|stopPropagation={handleDeleteButtonClick}
        >
            <svg viewBox="0 0 50 50">
                <path d="M1 1L49 49M1 49L49 1"/>
            </svg>
        </button>
    {/if}
</li>

<style>
    li.deletable-item {
        display: flex;
        align-self: center;
        align-items: center;
        gap: var(--spacing-1);
        padding: var(--spacing-1) calc(1.2 * var(--spacing-1));
        border-radius: var(--border-radius-2);
        transition: background .3s ease-out;
        
        &:has(button.delete-item:hover) {
            background: var(--danger-7) !important;
        }
    }

    button.delete-item {
        --btn-size: 9px;

        box-sizing: content-box;
        display: grid;
        width: var(--btn-size);
        height: var(--btn-size);
        place-items: center;
        padding: 5px;
    }

    svg {
        stroke: var(--grey-1);
        stroke-width: 2;
        width: 100%;
        height: 100%;
    }
</style>