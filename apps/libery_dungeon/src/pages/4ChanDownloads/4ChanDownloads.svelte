<script>
    import TrackedBoardLI from "@components/4chan/Boards/TrackedBoardLI.svelte";
    import BoardCatalog from "./sub-components/BoardCatalog/BoardCatalog.svelte";
    import { app_context_manager } from "@libs/AppContext/AppContextManager";
    import { app_contexts } from "@libs/AppContext/app_contexts";
    import { resetDownloadsContextStore } from "@stores/downloads";
    import { ChanTrackedBoard, getTrackedBoards } from "@models/4Chan";
    import ThreadContent from "./sub-components/BoardCatalog/ThreadContent.svelte";
    import { onMount } from "svelte";
    import { layout_properties } from "@stores/layout";
    import global_hotkeys_manager from "@libs/LiberyHotkeys/libery_hotkeys";
    import HotkeysContext from "@libs/LiberyHotkeys/hotkeys_context";
    import { selected_thread, selected_thread_id } from "./app_page_store";
    import { browser } from "$app/environment";
    import { access_state_confirmed, current_user_identity } from "@stores/user";
    
    /*=============================================
    =            AppContext            =
    =============================================*/
    
        const app_context = app_context_manager.CurrentContext;
        const app_page_name = "4-chan-downloads";

        if (app_context !== app_contexts.FOURCHAN) {
            app_context_manager.setAppContext(app_contexts.FOURCHAN, app_page_name, resetDownloadsContextStore);
        } else {
            app_context_manager.addOnContextExit(app_page_name, resetDownloadsContextStore);
        }
    
    /*=====  End of AppContext  ======*/
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        /** @type {ChanTrackedBoard[]} */
        let tracked_boards = [];

        /** 
         * @type {string} the selected board's name
         * @default ""
         */
        let selected_board = ""; // default to /b but /b is nsfw so in demos use /wg instead

        const hotkey_context_name = "4chan-downloads";
    
    /*=====  End of Properties  ======*/

    onMount(async () => {
        selected_board = getDefaultSelectedBoard();

        tracked_boards = await getTrackedBoards();        

        if (!layout_properties.IS_MOBILE) {
            defineDesktopKeybinds();
        }
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/

        const defineDesktopKeybinds = () => {
            if (!global_hotkeys_manager.hasContext(hotkey_context_name)) {
                const hotkeys_context = new HotkeysContext();

                hotkeys_context.register(["a", "q"], () => console.log("a or q pressed"));

                global_hotkeys_manager.declareContext(hotkey_context_name, hotkeys_context);
            }

            global_hotkeys_manager.loadContext(hotkey_context_name);
        }

        /**
         * Defines the best default selected board. The most popular board is /b but that is NSFW. If the user has
         * access to private content, then default to /b otherwise default to a 'safe for work' board.
         * @requires current_user_identity
         * @returns {string}
         */
        const getDefaultSelectedBoard = () => {
            const SAFE_FOR_WORK_BOARD = "g";
            const BEST_BOARD = "b";
            let best_default_board = SAFE_FOR_WORK_BOARD;

            if ($current_user_identity == null) return best_default_board;

            if ($current_user_identity.canReadPrivateContent()) {
                best_default_board = BEST_BOARD;
            }

            return best_default_board;
        }


        const handleNewBoardSelected = (event) => {
            $selected_thread = null;

            let new_board_name = event.detail.board_name;

            if (new_board_name !== selected_board) {
                selected_thread_id.set(null);
            }

            selected_board = new_board_name;
        }
    
    /*=====  End of Methods  ======*/

</script>

<main id="libery-4chan-downloads" class:adebug={false}>
    {#if browser && $access_state_confirmed && selected_board !== ""}
        {#if $selected_thread === null}
            <article id="l4d-board-content">
                <BoardCatalog board_name={selected_board} />
            </article>
        {:else}
            <article id="l4d-thread-content">
                <ThreadContent board_name={$selected_thread.board_name} thread_id={$selected_thread.uuid} />
            </article>
        {/if}
    {/if}
    <ul id="l4d-tracked-boards">
        {#each tracked_boards as tb}
            <TrackedBoardLI on:new-board-selected={handleNewBoardSelected} tracked_board={tb} is_current_board={tb.board_name === selected_board} />
        {/each}
    </ul>
</main>

<style>
    #libery-4chan-downloads * {
        box-sizing: border-box;
    }

    #libery-4chan-downloads  {
        display: grid;
        width: 100dvw;
        height: calc(100dvh - var(--navbar-height));
        overflow: hidden;
        grid-template-columns: calc(80% - calc(.5 * var(--vspacing-2))) calc(20% - calc(var(--vspacing-2) * .5));
        grid-auto-rows: calc(100dvh - var(--navbar-height));
        column-gap: var(--vspacing-2);
    }

    ul#l4d-tracked-boards {
        background: var(--grey-9);
        width: 100%;
        overflow: auto;
        height: 100%;
        grid-column: 2 / span 1;
        margin: 0;
        padding: 0;
        list-style: none;
        scrollbar-width: thin;
        scrollbar-color: var(--grey-8) var(--grey-9);
    }

    @media only screen and (max-width: 768px) {
        #libery-4chan-downloads {
            /* grid-template-columns: 94% 4%; */
            column-gap: var(--vspacing-1);
        }
    }
</style>