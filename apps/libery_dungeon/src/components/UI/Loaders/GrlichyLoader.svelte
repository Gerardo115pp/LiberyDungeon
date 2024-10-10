<script>
    
    /*=============================================
    =            Properties            =
    =============================================*/

        /**
         * The loader label
         * @type {string}
         */    
        export let label = "Loading...";
        
        /*----------  Styles  ----------*/
        
            export let loader_color = "var(--main-dark)";
    
    /*=====  End of Properties  ======*/
    
    
</script>

<div 
    class="loader" 
    style:--loader-color={loader_color}
    style:--label-text="'{label}'"
>
</div>



<style>
    .loader {
        font-weight: bold;
        font-family: var(--font-titles);
        color: var(--loader-color);
        display: inline-grid;
        font-size: var(--font-size-2)
    }
    
    .loader:before,
    .loader:after {
        content: var(--label-text);
        grid-area: 1/1;
        -webkit-mask-size: 100% calc(var(--font-size-2) * .16),100% 100%;
        -webkit-mask-repeat: no-repeat;
        -webkit-mask-composite: xor;
                mask-composite: exclude;
        animation: l35-1 1s infinite;
    }

    .loader:before {
        -webkit-mask-image:
            linear-gradient(#000 0 0),
            linear-gradient(#000 0 0);
    }

    .loader:after {
        -webkit-mask-image:
            linear-gradient(#000 0 0);
        animation:
            l35-1  1s infinite,
            l35-2 .2s infinite cubic-bezier(0.5,200,0.5,-200);
    }

    @keyframes l35-1{
        0%   {-webkit-mask-position:0 20px,0 0}
        20%  {-webkit-mask-position:0 8px ,0 0}
        40%  {-webkit-mask-position:0 100%,0 0}
        60%  {-webkit-mask-position:0 3px ,0 0}
        80%  {-webkit-mask-position:0 15px,0 0}
        100% {-webkit-mask-position:0 0   ,0 0}
    }

    @keyframes l35-2{
        100% {transform:translate(0.1px)} 
    }
</style>