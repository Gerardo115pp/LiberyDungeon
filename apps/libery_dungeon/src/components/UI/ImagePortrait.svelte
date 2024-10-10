<script>
    export let image_src;
</script>

<div class="image-portrait-wrapper">
    <div class="back-rect"></div>
    <img class="image-element" src="{image_src}" alt="portrait" />
</div>

<style>
    .image-portrait-wrapper{
        width: 41.50%;
        height: max-content;
        position: relative;
    }

    .back-rect {
        position: absolute;
        top: 0;
        left: 0;
        transform: rotate(-1.175deg);
        width: 100%;
        height: 100%;
        background-color: var(--main-dark-color-8);
        z-index: 1;
        border-radius: var(--border-radius);
    }

    .image-element {
        position: relative;
        width: 100%;
        border: .5px solid var(--main);
        border-radius: var(--border-radius);
        transform: rotate(0.936deg);
        z-index: 2;
    }
</style>