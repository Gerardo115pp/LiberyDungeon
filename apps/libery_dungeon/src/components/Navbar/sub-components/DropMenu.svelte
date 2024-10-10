<script>
    import { slide } from "svelte/transition";
    import { current_cluster } from "@stores/clusters";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
   
        /** 
         * @type {Section[]} 
         * @typedef {Object} Section
         * @property {string} name
         * @property {string} href
         */
        export let sections = [];

        /** @type {import('svelte/store').Writable<boolean>}*/
        export let visible;

        const appearance_animation_duration = 600;
        const appearance_animation_delay = 50;

        /**
         * We only show the dungeon link if there is a cluster selected. and if there is, we have to add the root category id to the href.
         * @type {string}
         */
        const dungeon_href = "/dungeon-explorer";
    
    
    /*=====  End of Properties  ======*/
    
    
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * @description Handles the click event on the anchor tags of the submenu items.
         * @param {MouseEvent} e
         */
        const clickAnchorHandler = e => {
            visible.set(false);
        }
    
    
    /*=====  End of Methods  ======*/

</script>

<div id="submenu-wrapper" style:animation-duration="{appearance_animation_duration}ms" style:animation-delay="{appearance_animation_delay}ms">
    <ul id="submenu">
        {#each sections as s, h}
            {@const link_href = s.href === dungeon_href && $current_cluster != null ? `${s.href}/${$current_cluster.RootCategoryID}` : s.href}
            {#if s.href !== dungeon_href || $current_cluster != null}
                <li in:slide={{axis: 'y', delay: (appearance_animation_delay+appearance_animation_duration) + (100*h), duration: 300}} class="submenu-item">
                    <a on:click={clickAnchorHandler} href="{link_href}">
                        {s.name}
                    </a>
                </li>
            {/if}
        {/each}
    </ul>
</div>

<style>
    @keyframes DropMenuAppears {
        0% {
            transform: translateX(-90%);
        }
        39% {
            transform: translateX(0);
            background-color: var(--main-dark);
        }
        100% {
            background: var(--grey);
            transform: translateX(0);
        }
    }

    #submenu-wrapper {
        position: fixed;
        display: grid;
        top: var(--navbar-height);
        left: 0;
        width: 100%;
        height: calc(100vh - var(--navbar-height));
        background: var(--main-dark);
        place-items: center;
        transform: translateX(-90%);

        animation-name: DropMenuAppears;
        animation-fill-mode: forwards;
    }

    #submenu {
        display: flex;
        flex-direction: column;
        align-items: center;
        row-gap: var(--vspacing-2);
        list-style: none;
        padding: 0;
    }

    .submenu-item {
        transform-origin: center center;
        transition: all .2s ease-in;
    }

    .submenu-item:hover {
        transform: translate(0, -2px) scale(1.1);
    }

    .submenu-item a {
        --angle-brackets-color: var(--grey-6);

        font-family: var(--font-decorative);
        cursor: pointer;
        text-decoration: none;
        color: var(--grey-1);
        font-size: calc(var(--font-size-h1) * .75);
        line-height: 1;
        transition: all .3s ease-in-out;
    }

    .submenu-item a:hover {
        color: var(--main-5);
        transform: translate(0, -5px) scale(1.1);
    }

    .submenu-item a::before {
        content: "<";
        margin-right: var(--vspacing-1);
        color: var(--angle-brackets-color);
    }

    .submenu-item a::after {
        content: "/>";
        margin-left: var(--vspacing-1);
        color: var(--angle-brackets-color);
    }
</style>