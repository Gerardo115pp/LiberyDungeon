<script>
    import { current_user_identity, has_user_access, access_state_confirmed, userLogout } from "@stores/user";
    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the click event on the logout button. Requests the server to invalidate the current user's access token and redirects to the login page.
        */
        const handleNavUserLogoutClick = async () => {
            userLogout();
        }
    
    /*=====  End of Methods  ======*/
    
</script>

{#if $access_state_confirmed && $has_user_access}
    <div id="ldn-nav-user-section-wrapper">
        <button id="ldn-nusw-user-label">
            {#if $current_user_identity != null}
                <h3>
                    {$current_user_identity.Username}
                </h3>
            {:else}
                <h3 style:color="var(--danger)">
                    <!-- 
                        This shouldn't be possible if has_user_access is true
                     -->
                    I am Error
                </h3>
            {/if}
        </button>
        <menu id="ldn-nusw-menu">
            {#if $current_user_identity?.canModifyUsers()}
                <li class="ldn-nusw-menu-item">
                    <a href="/users-management/dashboard" id="ldn-nusw-menu-item-user-dashboard">
                        Users dashboard
                    </a>
                </li>
            {/if}
            <li class="ldn-nusw-menu-item">
                <a id="ldn-nusw-menu-item-logout" href="/login" on:click={handleNavUserLogoutClick}>
                    Logout
                </a>
            </li>
        </menu>
    </div>        
{/if}

<style>
    #ldn-nav-user-section-wrapper {
        position: relative;
        display: grid;
        place-items: center;
        height: 100%;
        transition: background 0.3s ease-out;
        
        &:hover button#ldn-nusw-user-label{
            background: var(--main-dark);
            color: var(--body-bg-color);
            z-index: var(--z-index-t-4);
        }

        &:hover menu#ldn-nusw-menu {
            visibility: visible;
            translate: 0;
        }
    }

    button#ldn-nusw-user-label {
        cursor: default;
        width: 100%;
        background: var(--body-bg-color);
        padding: var(--spacing-2) var(--spacing-4);
        border-radius: 0;
        transition: background 0.3s ease-out;
        color: var(--main-dark);

        & h3 {
            font-family: var(--font-decorative);
            font-size: var(--font-size-2);
            text-transform: none;
            line-height: 1;
            color: inherit;
        }
    }

    
    /*=============================================
    =            Settings menu            =
    =============================================*/
    
        menu#ldn-nusw-menu {
            position: absolute;
            background: hsl(from var(--grey-8) h s calc(l * 0.7) / 0.88);
            top: 100%;
            right: 0;
            width: 100%;
            translate: 0 -20%;
            visibility: hidden;
            transition: background 0.3s ease-out 0.2s, translate 0.2s ease-out;
            z-index: var(--z-index-t-3);
        }

        li.ldn-nusw-menu-item {
            line-height: 1;
            height: 48px;
            padding: 0 var(--spacing-1);

            &:not(:last-child) {
                border-bottom: 1px solid var(--grey-9);
            }

            &:hover {
                background: hsl(from var(--main-9) h s l / 0.1);
            }

            &:has(a) {
                cursor: pointer;
                padding: 0;
            }

            &:has(a):hover {
                color: var(--main-2);
            }

            & > a {
                display: flex;
                width: 100%;
                height: 100%;
                align-items: center;
                padding: 0 var(--spacing-1);
                transition: color 0.3s ease-out;
            }
        }

        li:has(a#ldn-nusw-menu-item-logout):hover {
            color: var(--main-dark);
        }
    
    
    /*=====  End of Settings menu  ======*/
    
    

</style>