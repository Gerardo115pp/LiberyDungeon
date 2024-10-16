<script>
    import '@app/app.css';
    import { defineLayout, hotkeys_modal_visible, layout_properties, navbar_hidden } from '@stores/layout';
    import Navbar from '@components/Navbar/Navbar.svelte';
    import HotkeysHelpModal from '@components/Informative/HotkeyHelpModal/HotkeysHelpModal.svelte';
    import FloatingErrorDialog from '@libs/LiberyFeedback/FeedbackUI/FloatingErrorDialog.svelte';
    import { onDestroy, onMount } from 'svelte';
    import { LabeledError } from '@libs/LiberyFeedback/lf_models';
    import { lf_errors } from '@libs/LiberyFeedback/lf_errors';
    import FloatingMessageDialog from '@libs/LiberyFeedback/FeedbackUI/FloatingMessageDialog.svelte';
    import FloatingConfirmDialog from '@libs/LiberyFeedback/FeedbackUI/FloatingConfirmDialog.svelte';
    import global_platform_events_manager from '@libs/LiberyEvents/libery_events';
    import { page } from '$app/stores';
    import { access_state_confirmed, current_user_identity, has_user_access } from '@stores/user';
    import { goto } from '$app/navigation';
    import { getCurrentUserIdentity, isUsersInInitialSetupMode, validateUserAccessToken } from '@models/Users';
    import { INITIAL_SETUP_PAGE_PATH, isPublicPage, LOGIN_PAGE_PATH } from '@config/pages_routes';
    import { setupHotkeysManager, getHotkeysManager } from '@libs/LiberyHotkeys/libery_hotkeys';
    
    /*=============================================
    =            Properties            =
    =============================================*/
    
        let global_hotkeys_manager = null;
    
    /*=====  End of Properties  ======*/

    onMount(() => {
        defineLayout();

        if (global_hotkeys_manager == null) {
            setupHotkeysManager();
            global_hotkeys_manager = getHotkeysManager();
        }
        console.log("Global hotkeys manager: ", global_hotkeys_manager);
        
        console.log("Starting global platform events manager");
        global_platform_events_manager.start();

        checkUserAccess();
    });

    onDestroy(() => {
        console.log("Stopping global platform events manager");
        global_platform_events_manager.stop();

        if (global_hotkeys_manager != null) {
            global_hotkeys_manager.destroy();
        }
    });
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Checks if the user is logged in, if not, check if the system is on initial setup mode if so, redirects to the initial setup page. if is
         * just not logged in, redirects to the login page.
         */ 
        const checkUserAccess = async () => {
            // DELETE THIS IF, THIS IS JUST FOR TESTING
            // if (true) {
            //     has_user_access.set(true);
            //     access_state_confirmed.set(true);
            //     return;
            // }

            if ($page.url.pathname === INITIAL_SETUP_PAGE_PATH || $page.url.pathname === LOGIN_PAGE_PATH) {
                has_user_access.set(false);
                current_user_identity.set(null);
                access_state_confirmed.set(true);
                
                return;
            }

            let user_access_validation = await validateUserAccessToken();
            has_user_access.set(user_access_validation);

            if (!user_access_validation) {
                let is_initial_setup = await isUsersInInitialSetupMode();
                current_user_identity.set(null);

                if (is_initial_setup) {
                    goto("/initial-setup");
                } else {
                    goto("/login");
                }
            } else {
                /**
                 * @type {import("@models/Users").UserIdentity | null}
                 */
                let user_identity = null;

                user_identity = await getCurrentUserIdentity();

                has_user_access.set(user_identity != null);
                current_user_identity.set(user_identity);
            }

            access_state_confirmed.set(true);
        }
    
    /*=====  End of Methods  ======*/
    
    
</script>

{#if global_hotkeys_manager != null} 
    <div id="libery-website-wrapper" class:navbarless={$navbar_hidden}>
        <Navbar />
        {#if $hotkeys_modal_visible}
            <HotkeysHelpModal />
        {/if}
        <div id="libery-dungeon-content">
            {#if $access_state_confirmed || isPublicPage($page.url.pathname)}
                <slot></slot>  
            {/if}
        </div>
        <FloatingErrorDialog 
            display_error_duration={7000}
        />
        <FloatingMessageDialog 
            display_message_duration={5000}
        />
        <FloatingConfirmDialog />
    </div>
{/if}

<style>
    #libery-website-wrapper {
        border-top: .1px solid var(--grey);
        position: relative;
        /* container-type: inline-size; */
    }

    #libery-dungeon-content {
        margin-top: var(--navbar-height);
    }

    .navbarless #libery-dungeon-content {
        margin-top: 0;
    }
</style>