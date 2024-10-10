<script>
    import FieldData, { verifyFormFields } from "@libs/FieldData";
    import Input from "@components/Input/Input.svelte";
    import ContactMessage from "@models/Contact";
    import { show_contact_form } from "@stores/layout";
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";

    let new_contact_submission = [
        new FieldData("cfm-lc-name", /^.{3,}$/, "name"),
        new FieldData("cfm-lc-email", /^.{3,}$/, "email"),
        new FieldData("cfm-lc-company", /^.{3,}$/, "company"),
        new FieldData("cfm-lc-budget", /^.{3,}$/, "budget"),
        [
            new FieldData("cfm-lc-start-date", /^.{3,}$/, "Project start", "date", false),
            new FieldData("cfm-lc-end-date", /^.{3,}$/, "end", "date", false)
        ],
        new FieldData("cfm-lc-message", /^.{3,}$/, "message", "textarea", false)
    ]

    // Defining placeholders
        new_contact_submission[0].placeholder = "John Doe";
        new_contact_submission[1].placeholder = "hello@example.com";
        new_contact_submission[2].placeholder = "Google Inc.";
        new_contact_submission[3].placeholder = "$1000 - $5000";
        new_contact_submission[5].placeholder = "I heard you were the best in the business!";
    // End of defining placeholders

    let is_form_ready = false;
    let show_form = $show_contact_form;

    
    /*=============================================
    =            Animation settings            =
    =============================================*/

        /**
         * @type {HTMLDivElement}
         * Element that contains the content of the popup and the we will scroll to the bottom of
         */
        let first_visible_element;
    
        const initial_animation_duration = 600;
        const content_animation_duration = 500;
    
    /*=====  End of Animation settings  ======*/
    
    


    onMount(() => {
        show_contact_form.subscribe(toggleContactForm);
    })


    const submitContactForm = async () => {
        if (!is_form_ready) return;

        const message_body = {
            "name": new_contact_submission[0].getFieldValue(),
            "email": new_contact_submission[1].getFieldValue(),
            "company": new_contact_submission[2].getFieldValue(),
            "budget": new_contact_submission[3].getFieldValue(),
            "start_date": new_contact_submission[4][0].getFieldValue(),     
            "end_date": new_contact_submission[4][1].getFieldValue(),
            "subject": new_contact_submission[5].getFieldValue()
        }

        const contact_message = new ContactMessage(message_body);

        const message_recieved = await contact_message.send();

        if (!message_recieved) {
            console.log("Message not recieved");
        }
    }

    const toggleContactForm = (new_value) => {
        show_form = new_value;

        setTimeout(() => {
            if (show_form) {
                first_visible_element.scrollIntoView({behavior: "smooth"});
            }
        }, initial_animation_duration);
    }

    const verifyContactForm = () => {
        let is_valid = verifyFormFields(new_contact_submission);
        is_form_ready = is_valid;
        new_contact_submission = [...new_contact_submission]; 
    }
</script>

{#if show_form}
     <aside id="contact-form-modal-wrapper">
         <div style:animation-duration="{initial_animation_duration}ms" id="contact-form-modal">
             <header in:fade={{delay: (initial_animation_duration * 1), duration: content_animation_duration}} id="cfm-form-header">
                 <h3 class="cfm-fh-headline">
                     Level up your business — get a free quote!
                 </h3>
                 <div class="close-btn-wrapper">
                     <svg on:click={() => show_contact_form.set(false)} tabindex="0" role="button" aria-label="close btn" id="cfm-close-btn" viewBox="0 0 50 50">
                         <path d="M0 0L50 50" class="cfm-close-line"/>
                         <path d="M0 50L50 0" class="cfm-close-line"/>
                     </svg>
                 </div>
             </header>
             <div  id="cfm-form-content" in:fade={{delay: (initial_animation_duration * 1), duration: content_animation_duration}}>
                 <form action="none" id="cfm-libery-contact">
                     {#each new_contact_submission as fd}
                         {#if fd instanceof FieldData}
                             <div class="cfm-input-wrapper">
                                 <Input 
                                     field_data={fd}
                                     border_color="transparent"
                                     input_color="var(--main-dark-color-1)"
                                     text_font="var(--font-decorative)"
                                     label_color="var(--main-dark-color-5)"
                                     input_label={fd.name}
                                     placeholder_color="var(--grey-6)"
                                     input_padding="calc(var(--vspacing-1) * 1) var(--vspacing-2)"
                                     onBlur={verifyContactForm}
                                     show_placeholder
                                     modeColumn
                                     isSquared
                                 />
                             </div>
                             {:else if fd instanceof Array}
                             <div class="input-groups">
                                 {#each fd as nested_fd}
                                     <div class="cfm-input-wrapper">
                                         <Input 
                                             field_data={nested_fd}
                                             border_color="transparent"
                                             input_color="var(--main-dark-color-1)"
                                             label_color="var(--main-dark-color-5)"
                                             input_label={nested_fd.name}
                                             placeholder_color="var(--grey-6)"
                                             input_padding="calc(var(--vspacing-1) * 1) var(--vspacing-2)"
                                             show_placeholder
                                             modeColumn
                                             isSquared
                                         />
                                     </div>
                                 {/each}
                             </div>
                         {/if}
                     {/each}
                 </form>
                 <div id="cfm-libery-controls">
                     <div class="cfm-input-wrapper">
                         <input type="checkbox" checked>
                         <h5>
                             I agree to be a nice and kind person!
                         </h5>
                     </div>
                     <div bind:this={first_visible_element} class="button-1-wrapper">
                         <button on:click={submitContactForm} disabled={!is_form_ready} class="button-1">Send Message</button>
                     </div>
                 </div>
             </div>
         </div>
     </aside>
{/if}

<style>
    @keyframes present-contact-form {
        0% {
            clip-path: circle(20%);
            transform: translateY(-100%);
        }
        50% {
            clip-path: circle(20%);
            transform: translateY(0%);
        }
        100% {
            clip-path: circle(100%);
            transform: translateY(0%);
        }
    }

    #contact-form-modal-wrapper {
        position: fixed;
        width: 100vw;
        height: 100vh;
        display: flex;
        top: 0;
        left: 0;
        justify-content: center;
        align-items: center;
        background-color: var(--grey-t);
        z-index: var(--z-index-t-3);
    }
    
    #contact-form-modal {
        border-radius: var(--border-radius);
        background: var(--grey-8);
        width: 700px;
        height: 80%;
        padding: var(--vspacing-3);
        clip-path: circle(20%);
        transform: translateY(-100%);

        animation-name: present-contact-form;
        animation-timing-function: ease-in-out;
        animation-fill-mode: forwards;
    }

    header#cfm-form-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        height: 50px; 
        padding: 0 0 var(--vspacing-3) 0;
        background: var(--grey-8);
    }

    #cfm-close-btn {
        width: 20px;
        stroke: var(--grey-2);
        stroke-linecap: round;
        stroke-width: 1.5px;
        transition: all .2s ease-in-out;
    }

    @media(pointer: fine) {
        #cfm-close-btn:hover {
            stroke: var(--grey-1);
            transform: scale(1.1);
        }
    }

    #cfm-form-content {
        height: 90%;
        overflow-y: auto;
    }

    #cfm-form-content::-webkit-scrollbar{
        background: var(--grey-4);
    }

    #cfm-form-content::-webkit-scrollbar-thumb{
        background-color: var(--grey-4);
        border-radius: calc(.27777 * var(--border-radius-2));
        border: 3px solid var(--grey-8);
    }

    #cfm-form-content::-webkit-scrollbar-track{
        background: var(--grey-8);
    }

    #cfm-libery-contact {
        display: flex;
        flex-direction: column;
        gap: var(--vspacing-2);
    }

    #cfm-libery-contact .input-groups {
        display: flex;
        gap: var(--vspacing-3);
    }

    #cfm-libery-controls {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding:  0 var(--vspacing-1) 0 0;
    }

    .cfm-input-wrapper {
        display: flex;
        align-items: center;
        gap: var(--vspacing-1);
    }

    .cfm-input-wrapper h5 {
        font-family: var(--font-read);
    }

    #cfm-libery-controls .button-1-wrapper:has(button:disabled) {
        filter: saturate(.1);
        pointer-events: none;
    }

    @media only screen and (max-width: 767px) {
        #contact-form-modal {
            width: 100%;
            height: 100%;
            border-radius: 0;
            padding: var(--vspacing-3) var(--vspacing-2);
        }

        #cfm-libery-contact .input-groups {
            flex-direction: column;
            gap: var(--vspacing-2);
        }
    }
</style>