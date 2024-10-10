<script>
    import { layout_properties } from "@stores/layout";
    import { draw, fade } from "svelte/transition";
    import { onMount } from "svelte";
    import { tweened } from "svelte/motion";

    export let svg_width = 697;
    export let svg_height = 102;
    let svg_height_width_ratio = svg_height / svg_width;
    export let control_point_offset = 70;
    let x_offset = control_point_offset;
    let y_offset = (control_point_offset * svg_height_width_ratio) / 2;

    const central_square_size = layout_properties.IS_MOBILE ? 9 : 6;

    /**
     * @typedef {Object} ControlPointDifference
     * @property {number} dx
     * @property {number} dy
     */

    let control_point_difference = tweened({ dx: 0, dy: 0})
    let control_point_difference_end = { dx: -231.3881, dy: 47.12554999999999}
    let central_square_position = {x: 289.0815, y: 50.015550000000005}

    let control_line = {
        x1: (central_square_position.x + $control_point_difference.dx),
        y1: (central_square_position.y + $control_point_difference.dy),
        x2: (central_square_position.x - $control_point_difference.dx),
        y2: (central_square_position.y - $control_point_difference.dy)
    }

    let visual_control_line = {
        x1: (central_square_position.x + $control_point_difference.dx) + x_offset,
        y1: (central_square_position.y + $control_point_difference.dy) - y_offset,
        x2: (central_square_position.x - $control_point_difference.dx) - x_offset,
        y2: (central_square_position.y - $control_point_difference.dy) + y_offset
    }
    
    const end_point = {x: 555.765, y: 74.1149}

    
    /*----------  Animation  ----------*/
        export let animation_delay = 1500;
        let svg_visible = false;
        const line_draw_duration = 1500;

    onMount(() => {
        setTimeout(() => {
            svg_visible = true;

            setTimeout(() => {
                setControlLine(control_point_difference_end);
            }, line_draw_duration);
        }, animation_delay);
    });

    /**
     * @param {ControlPointDifference} new_params
     */
    function setControlLine(new_params) {
        control_point_difference.subscribe(value => {
            control_line = {
                x1: (central_square_position.x + value.dx),
                y1: (central_square_position.y + value.dy),
                x2: (central_square_position.x - value.dx),
                y2: (central_square_position.y - value.dy)
            }
            visual_control_line = {
                x1: (central_square_position.x + value.dx) + x_offset,
                y1: (central_square_position.y + value.dy) - y_offset,
                x2: (central_square_position.x - value.dx) - x_offset,
                y2: (central_square_position.y - value.dy) + y_offset
            }
        })

        control_point_difference.set(new_params);
    }
    
</script>

<svg id="anchor-symbol" viewBox="0 0 {svg_width} {svg_height}" fill="none" preserveAspectRatio="none">
    {#if svg_visible}
        <g id="anchor-element">
            <!-- <path data-name="main-curve-line" class="anchors-1-cls-2" d="M1.93344 0.195312C1.93344 0.195312 57.6934 97.1411 272.036 61.6883C486.378 26.2355 638.165 43.3881 695.765 84.1149" stroke="#053160" stroke-miterlimit="10"/> -->
            <path in:draw={{duration: line_draw_duration, speed: 3}} data-name="main-curve-line" class="anchors-1-cls-2" d="M1.93344 0.195312Q{control_line.x1} {control_line.y1} {central_square_position.x} {central_square_position.y}T{end_point.x} {end_point.y}" stroke="#053160" stroke-miterlimit="10"/>
            <rect in:fade={{delay: line_draw_duration}} id="center-rect" x="{central_square_position.x}" y="{central_square_position.y}" width="{central_square_size}" height="{central_square_size}" transform="translate(-{central_square_size/2}, -{central_square_size/2})"/>
            <path in:fade={{delay: line_draw_duration*.92}} data-name="guide-line" class="anchors-1-cls-2" d="M{visual_control_line.x1} {visual_control_line.y1}L{visual_control_line.x2} {visual_control_line.y2}"  stroke-miterlimit="10"/>
            <circle in:fade={{delay: line_draw_duration}} class="guide-controls" cx="{visual_control_line.x1}" cy="{visual_control_line.y1}" r={(central_square_size/2) * 1.4}/>
            <circle in:fade={{delay: line_draw_duration}} class="guide-controls" cx="{visual_control_line.x2}" cy="{visual_control_line.y2}" r={(central_square_size/2) * 1.4}/>
        </g>
    {/if}
</svg>    

<style>
    svg#anchor-symbol {
        max-width: calc(0.5816494140625 * var(--vspacing-8));
        height: 102px;
    }

    g#anchor-element {
        width: 100%;
    }

    /* .anchors-1-cls-1 {
        fill: var(--accent-8);
    } */
        
    .anchors-1-cls-2 {
        fill: none;
        stroke: var(--accent-8);
        stroke-width: 1.3px;
        stroke-miterlimit: 10;
    }

    .guide-controls {
        fill: var(--accent-8);
        opacity: 0.8;
    }

    #center-rect {
        stroke: var(--accent-8);
        fill: var(--accent-8);
        opacity: 0.8;
    }

    @media only screen and (max-width: 768px) {
        svg#anchor-symbol {
                width: 100%;
                max-width: none;
                height: auto;
            }

        .anchors-1-cls-2 {
            fill: none;
            stroke: var(--accent-7);
            stroke-width: 2.3px;
            stroke-miterlimit: 10;
        }
    }
</style>