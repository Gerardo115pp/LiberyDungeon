
<div class="hr-divider-wrapper">
    <svg viewBox="0 0 100 1" xmlns="http://www.w3.org/2000/svg">
        <g id="divider-line" >
            <line class="hr-background" x1="0" y1=".5" x2="100" y2=".5"/>
            <line class="hr-line" x1="0" y1=".2" x2="100" y2=".2"/>
        </g>
    </svg>
    <i class="html-tag">
        &lt;/hr&gt;
    </i>
</div>

<style>
    .hr-divider-wrapper {
        position: relative;
        width: 100%;
        /* background: var(--grey); */
    }

    .hr-line {
        stroke: var(--grey-8);
        stroke-width: .05px;
    }

    .hr-background {
        stroke: var(--grey);
        stroke-width: 3px;
    }

    i.html-tag {
        position: absolute;
        top: 50%;
        left: 75%;
        transform: translate(-50%, -50%);
        color: var(--grey-8);
        font-family: var(--font-decorative);
        font-size: calc(1.4 * var(--font-size-4));
        background: var(--grey);
        line-height: 1;
        padding: 0 8px;
    }
</style>