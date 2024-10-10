<script>
    import { goto } from "$app/navigation";
    import UserCreation from "@components/Users/UsersCreation/UserCreation.svelte";

    
    /*=============================================
    =            Methods            =
    =============================================*/
    
        /**
         * Handles the Initial user creation event emitted by the UserCreation component.
         * @param {CustomEvent<UserCreatedEvent>}
         * @typedef {Object} UserCreatedEvent
         * @property {string} username 
         */ 
        const handleInitialUserCreation = (event) => {
            console.log("User created: ", event.detail.username);

            goto("/login"); // If you are looking where to implement further initial setup steps, this is it, overwrite this line and add you steps
        }
    
    /*=====  End of Methods  ======*/
    
    
</script>

<main id="libery-dungeon-initial-setup-page" class="form-page-wrapper">
    <div id="ldisp-user-creation-wrapper" class="form-wrapper">
        <UserCreation 
            on:user-created={handleInitialUserCreation}
        />
    </div>
</main>

<style>

    #ldisp-user-creation-wrapper {
        padding: var(--spacing-4) 0;
    }
        
</style>