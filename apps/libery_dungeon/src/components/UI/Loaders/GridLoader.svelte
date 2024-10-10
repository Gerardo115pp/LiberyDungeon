<script>
    
    /*=============================================
    =            Properties            =
    =============================================*/
        
        /*----------  Styles  ----------*/
        
            export let loader_color = "var(--main-dark)";
    
    /*=====  End of Properties  ======*/
    
    
</script>

<div class="loader" style:--loader-color={loader_color}></div>

<style>
    .loader {
        width: 148px;
        aspect-ratio: 1;
        border: 2px solid var(--loader-color);
        --c:conic-gradient(from -90deg at calc(100% - 2px) calc(100% - 2px),#0000 0 90deg,var(--loader-color) 0);
        background: var(--c),var(--c);
        background-size: 16px 16px;
        background-position: 0 0;
        animation: l3 1s infinite;
    }

    @keyframes l3{
        100% {
            background-position: -16px -16px,16px 16px
        }
    }
</style>