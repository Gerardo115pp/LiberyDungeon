<script>
    import { layout_properties, navbar_ethereal, navbar_solid } from "@stores/layout";
    import { writable } from "svelte/store";
    import MainLogo from "@components/UI/MainLogo.svelte";
    import BurgerBtn from "@components/UI/BurgerBTN.svelte";
    import { fade } from "svelte/transition";
    import PageNavMenu from "./sub-components/PageNavMenu.svelte";
    import DropMenu from "./sub-components/DropMenu.svelte";
    import { navbar_hidden } from "@stores/layout";
    import NavUserMenu from "./sub-components/user_nav_section/NavUserMenu.svelte";
    import { current_user_identity } from "@stores/user";
    import { onDestroy, onMount } from "svelte";
    import { browser } from "$app/environment";

    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        let menu_visible = writable(false);
        /** @type {string} is the value of the css property --navbar-height */

        /**
         * @typedef {Object} NavbarSections
         * @property {string} name
         * @property {string} href
         * @property {NavbarSections[]} options 
        */

        /** @type {NavbarSections[]} */
        const public_dropdown_sections = [
            {
                name: "Home",
                href: "/",
                options: []
            },
            {
                name: "Dungeon",
                href: "/dungeon-explorer",
                options: []
            }
        ];

        /** @type {NavbarSections[]} */
        let dropdown_sections = [
            ...public_dropdown_sections
        ];


        let user_identity_unsubscriber = () => {};

    /*=====  End of Properties  ======*/

    onMount(() => {
        if ($current_user_identity != null) {
            populateDropdownSections();
        }

        user_identity_unsubscriber = current_user_identity.subscribe(populateDropdownSections);
    });

    onDestroy(() => {
        if (!browser) return;

        user_identity_unsubscriber();
    })
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Adds the 4chan downloader if the user has the appropriate permissions.
         * @returns {void}
         */
        const addFourChanDownloader = () =>  {
            let has_download_grant = ($current_user_identity?.canDownloadFiles() ?? false);

            if (has_download_grant) {
                dropdown_sections = [
                    ...dropdown_sections,
                    {
                        name: "4chan Downloads",
                        href: "/4chan-downloads",
                        options: []
                    }
                ]
            }
        }

        /**
         * Populates the dropdown sections with the optional sections.
         * @returns {void}
         */
        const populateDropdownSections = () => {
            dropdown_sections = [
                ...public_dropdown_sections
            ];
        
            addFourChanDownloader();
        }
    
    /*=====  End of Methods  ======*/
    
    

</script>

<nav id="libery-dungeon-navbar"
    role="toolbar" 
    class:navbar-is-hidden={$navbar_hidden} 
    class:navbar-ethereal={$navbar_ethereal && !$menu_visible}
    class:opaque-background={$menu_visible || $navbar_solid} 
    class:adebug={false}
>
    <div id="ldn-content">
        <div id="ldn-left-content">
            <div class="ldn-rc-burger-wrapper">
                <BurgerBtn on:click={() => menu_visible.set(!$menu_visible)} is_opened={menu_visible} />
            </div>
            <div id="ldn-main-logo-wrapper" class:adebug={false}>
                <MainLogo 
                    fill_color="var(--main-dark)" 
                    stroke_color="var(--main-dark)" 
                />
            </div>
            <PageNavMenu />
        </div>
        <div id="ldn-right-content">
            <div class="ldn-rc-item">
                <NavUserMenu />
            </div>
        </div>
    </div>
    {#if $menu_visible}
        <DropMenu visible={menu_visible} sections={dropdown_sections} />
    {/if}
</nav>

<style>

    #libery-dungeon-navbar {
        --navoptions-gap: 81px;

        position: fixed;
        top: 0;
        display: grid;
        height: var(--navbar-height);
        background: var(--glass-gradient);
        container-type: size;
        width: 100%;
        place-items: center;
        z-index: var(--z-index-t-2);
        backdrop-filter: blur(10px);
        border-bottom: .5px solid var(--grey-9);
        transition: background 0.3s ease-out, opacity 0.3s ease-out, visibility 0.3s linear allow-discrete, backdrop-filter 0.25s linear allow-discrete;
    }

    #libery-dungeon-navbar.navbar-ethereal {
        background: linear-gradient(to bottom, hsl(from var(--body-bg-color) h s l / 0.85), transparent);
        border-bottom: none;
        backdrop-filter: none;
    }

    :global(.dark-mode nav#libery-dungeon-navbar) {
        background: var(--black-glass-gradient);
    }
    
    #libery-dungeon-navbar.opaque-background {
        background: var(--grey) !important;
        border-bottom: 1px solid var(--grey-1);
    }

    #libery-dungeon-navbar.navbar-is-hidden {
        visibility: hidden;
        opacity: 0;
        height: 0px;
    }

    #ldn-content {
        display: flex;
        width: var(--page-content-width);
        justify-content: space-between;
    }

    #ldn-left-content {
        display: flex;
        align-items: center;
        gap: var(--navoptions-gap);
    }

    .ldn-rc-burger-wrapper {
        --burger-btn-width: calc(1.1 * var(--vspacing-3));

        box-sizing: border-box;
        width: var(--burger-btn-width);
        height: calc(0.8666666666666667 * var(--burger-btn-width));
    }

    #ldn-main-logo-wrapper {
        width: clamp(150px, 11vw, 200px);
    }

    .ldn-rc-item {
        height: 100%;
    }

    @media only screen and (max-width: 767px) {
        #libery-dungeon-navbar {
            --navoptions-gap: 10vw;
        }

        #ldn-main-logo-wrapper {
            width: calc(var(--vspacing-7) * 0.2);
        }
    }
</style>