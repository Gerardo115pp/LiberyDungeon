<script>
    export let tooltip_text;
    // export let tool_tip_position = "top";
    export let show_tool_tip = true;
    

    
    /*=============================================
    =            Styling            =
    =============================================*/
    
    export let tooltip_color = "white";
    export let tooltip_background_color = "black";
    
    
    /*=====  End of Styling  ======*/
    
    
</script>

<div 
    class="tooltip-container"
    style:visibility={show_tool_tip ? "visible" : "hidden"}
    style:color={tooltip_color}
    style:background={tooltip_background_color}
    style:border-color={tooltip_background_color}
>
{tooltip_text}
</div>

<style>
    :global(:has(> .tooltip-container)) {
        position: relative;
    }

    .tooltip-container {
        position: absolute;
        bottom: 150%;
        border-radius: 15px;
    }

    .tooltip-container::after {
        content: "";
        position: absolute;
        top: 100%;
        left: 50%;
        margin-left: -5px;
        border-width: 5px;
        border-style: solid;
        border-color: inherit transparent transparent transparent;
    }


</style>