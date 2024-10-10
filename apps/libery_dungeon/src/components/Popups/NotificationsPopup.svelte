<script>
    import { onMount, onDestroy } from 'svelte';
    import notification_events from './events';
    let show_notifications = false;
    let notification_text = "";

    onMount(() => {
       document.addEventListener(notification_events.NEW_NOTIFICATION, showNotification);
    });

    onDestroy(() => {
        document.removeEventListener(notification_events.NEW_NOTIFICATION, showNotification);
    });
    
    const showNotification = e => {
        const { message } = e.detail;
        notification_text = message;
        show_notifications = true;
    }

    const hideNotification = () => {
        show_notifications = false;
    }

</script>

<div id="notification-popup-wrapper">
    {#if show_notifications}
        <div id="notification-popup" class="slide-fwd-bottom">
            <span class="notification-content">{notification_text}</span>
            <button on:click={hideNotification} class="nb-close-btn">X</button>
        </div>
    {/if}
</div>

<style>

    #notification-popup-wrapper {
        box-sizing: border-box;
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        z-index: 5;
        display: grid;
        place-items: center;
    }
    
    #notification-popup-wrapper:has(#notification-popup) {
        padding: var(--spacing-3);
    }

    #notification-popup {
        box-sizing: border-box;
        display: flex;
        width: 50%;
        background: var(--theme-purple);
        border-radius: var(--boxes-roundness);
        padding: var(--spacing-2) var(--spacing-2);
        color: var(--clear-color);
        justify-content: space-between;
        align-items: center;
        box-shadow: var(--boxes-shadow);
    }

    .notification-content {
        cursor: default;
        /* user-select: none; */
        font-size: var(--font-size-1);
        text-transform: uppercase;
        font-weight: 500;
    }

    #notification-popup .nb-close-btn {
        background: transparent;
        border: none;
        color: var(--clear-color);
        font-size: var(--font-size-2);
        font-weight: 500;
        transition: all 0.2s ease-in;
    }

    #notification-popup .nb-close-btn:hover {
        color: var(--theme-pearl);
    }

    @media only screen and (max-width: 768px) {
        #notification-popup-wrapper {
            padding: var(--spacing-2) 0 ;
        }

        #notification-popup {
            width: 95%;
            padding: var(--spacing-2);
        }

        .notification-content {
            max-width: 80%;
            font-size: var(--font-size-1);
        }
    }

</style>
