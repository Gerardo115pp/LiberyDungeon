<script>
    import { ChanTrackedBoard } from "@models/4Chan";
    import { createEventDispatcher, onMount } from "svelte";
    import { layout_properties } from "@stores/layout";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
    /** @type {ChanTrackedBoard}*/
    export let tracked_board;

    /** @type {bool} is this the current board */
    export let is_current_board = false;

    const tracked_board_event_dispatcher = createEventDispatcher();
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        if (is_current_board) {
            ensureSelectedBoardVisible();
        }
    })
    
   
   /*=============================================
   =            Methods            =
   =============================================*/


        /**
         * Ensure selected board is visible
         */
        const ensureSelectedBoardVisible = () => {
            const selected_board = getSelectedBoard();

            if (selected_board != null) {
                selected_board.scrollIntoView({behavior: "smooth", block: "center"});
            }
        }

        /**
         * Returns the board with the class .selected-board
         * @returns {HTMLLIElement | null}
         */
        const getSelectedBoard = () => {
            return document.querySelector(".libery-4chan-tracked-board-li.selected-board");
        }
   
        const handleTrackedBoardClick = () => {
            tracked_board_event_dispatcher("new-board-selected", {board_name: tracked_board.board_name});
        }   


   
   /*=====  End of Methods  ======*/ 
    
</script>

<li on:click={handleTrackedBoardClick} class="libery-4chan-tracked-board-li" class:selected-board={is_current_board}>
    <h3>
        <strong class="l4tbl-board-name">{tracked_board.board_name}</strong>
        {!layout_properties.IS_MOBILE ? tracked_board.description : ''}
    </h3>
</li>

<style>
    .libery-4chan-tracked-board-li {
        cursor: default;
        padding: var(--vspacing-1) var(--vspacing-2);
        transition: all 0.3s ease-in-out;
    }

    @media (pointer: fine) {
        .libery-4chan-tracked-board-li:hover {
            background: var(--main-dark-color-8);
        }
    }

    .libery-4chan-tracked-board-li.selected-board {
        background: var(--main-dark);
    }
    
    .libery-4chan-tracked-board-li h3 {
        user-select: none;
        display: flex;
        flex-direction: row;
        align-items: center;
        column-gap: var(--vspacing-2);
        font-family: var(--font-read);
    }

    .libery-4chan-tracked-board-li.selected-board h3 {
        color: var(--grey);
    }
    
    .l4tbl-board-name {
        font-family: var(--font-decorative);
        width: 6ch;
        font-weight: bold;
        font-size: var(--font-size-2);
        color: var(--main);
    }

    .l4tbl-board-name::before {
        font-family: var(--font-titles);
        color: var(--grey-6);
        content: "/";
        margin-right: calc(var(--vspacing-1) * 0.5);
    }

    .libery-4chan-tracked-board-li.selected-board .l4tbl-board-name {
        color: var(--grey-1);
    }

    @media only screen and (max-width: 768px) {
        .libery-4chan-tracked-board-li {
            width: 100%;
            padding: var(--vspacing-1);
        }

        .l4tbl-board-name {
            width: 100%;
            text-align: center;
            font-size: calc(var(--font-size-4) * 1.2);
        }

    }
</style>