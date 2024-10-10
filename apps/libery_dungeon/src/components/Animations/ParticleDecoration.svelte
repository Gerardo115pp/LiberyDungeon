<script>
    import { StaticParticle, gaussRandom } from "./animatables";
    import { isMobile } from "@libs/itz-utils";

    
    /*=============================================
    =            ParticlesGeneration            =
    =============================================*/
    
        const absolute_particles_min_size = 10;

        /**
         * the max size a particle can be is times_min_size_max * particles_min_size
         * @type {number}
         */
        const times_min_size_max = 4;

        /**
         * the color of the particles
         * @type {string}
         */
        export let particles_color = "blue";

        /**
         * the minimum size a particle can be
         * @type {number}
         */
        export let particles_min_size = 15; // in px
        particles_min_size = particles_min_size < absolute_particles_min_size ? absolute_particles_min_size : particles_min_size;

        /**
         * the standard diviation of the particles size
         * @type {number}
         */
        export let particles_standard_diviation = 8;

        /**
         * the minimum number of particles
         * @type {number}
         */
        export let particles_min_population = 4;

        /**
         * the maximum number of particles
         * @type {number}
         */
        export let particles_max_population = 8;

        /**
         * the minimum opacity a particle can have
         * @type {number}
         */
        export let particles_min_opacity = 0.5;

        /**
         * the maximum opacity a particle can have
         * @type {number}
         */
        export let particles_max_opacity = 1;

        /**
         * no other particle can be close a particles (2*particle.radius)*particles_repultion
         * @type {number}
         */
        export let particles_repulsion = 2; 
    
    
        /**
         * @typedef {Object} ParticlePosition
         * @property {number} x
         * @property {number} y
        */

        /**
         * @returns {StaticParticle[]}
         * @requires particles_color
         * @requires particles_min_size
         * @requires particles_min_population
         * @requires particles_max_population 
         */
        const getParticles = () => {
            let particles = [];
            let particles_population = Math.floor(Math.random() * (particles_max_population - particles_min_population + 1) + particles_min_population);
            const max_denied_times = 3;


            for (let i = 0; i < particles_population; i++) {
                let random_position = getRandomParticlePosition();
                let times_denied = 0;
                let is_position_valid = true;

                do {
                    particles.forEach(p_other => {
                        if (p_other.IsToClose(random_position) ) {
                            is_position_valid = false;
                            times_denied++;
                        }
                    });
                } while (!is_position_valid && times_denied < max_denied_times);
                
                let new_particle = new StaticParticle(random_position.x, random_position.y, getRandomParticleSize()/2, particles_color);
                new_particle.repulsion_factor = particles_repulsion;
                particles.push(new_particle);
            }

            return particles;
        }

        /**
         * @returns {number}
         */
        const getRandomParticleSize = () => {
            let random_size = Math.floor(gaussRandom(particles_min_size, particles_standard_diviation))
            random_size = random_size < absolute_particles_min_size ? absolute_particles_min_size : random_size;
            random_size = random_size > particles_min_size * times_min_size_max ? particles_min_size * times_min_size_max : random_size;

            console.log(random_size);

            return random_size;
        }

        /**
         * @returns {ParticlePosition}
         */
        const getRandomParticlePosition = () => {
            let random_x = Math.floor(Math.random() * 100); // in %
            let random_y = Math.floor(Math.random() * 100); // in %

            return {
                x: random_x,
                y: random_y
            }
        }

        const getRandomOpacity = () => {
            return Math.random() * (particles_max_opacity - particles_min_opacity) + particles_min_opacity;
        }

    /*=====  End of ParticlesGeneration  ======*/
    

    
    /*=============================================
    =            ParticlesMovement            =
    =============================================*/
    
        /**
         * the higher the number the slower the particles will move
         * @type {number}
         */
        export let movement_min_duration = 14;

        /**
         * the higher the greater the movement difference between particles
         * @type {number}
         */
        export let movement_stdev = 10;

        /**
         * the higher the longe is going to take before the particles start moving
         * @type {number}
         */
        export let delay_min = 0.2;

        /**
         * the higher the longer is going to take before the particles start moving
         * @type {number}
         */
        export let delay_stdev = 1;
    
    /*=====  End of ParticlesMovement  ======*/
    
    

    
    /*=============================================
    =            MouseTracking            =
    =============================================*/
    
        let mouse_x = 0; // in %
        let mouse_y = 0;

        /**
         * determines how much the particles will move in relation to the mouse
         * @type {number}
         */
        export let mouse_tracking_factor = 17;

        /**
         * @type {HTMLDivElement}
         */
        let a_particles_wrapper;

        const updateMousePosition = (event) => {
            if (a_particles_wrapper === undefined || is_mobile) return;

            let a_particles_wrapper_rect = a_particles_wrapper.getBoundingClientRect();

            let x_in_px = event.clientX - a_particles_wrapper_rect.left;
            let y_in_px = event.clientY - a_particles_wrapper_rect.top;

            mouse_x = (x_in_px / a_particles_wrapper_rect.width) * mouse_tracking_factor;
            mouse_y = (y_in_px / a_particles_wrapper_rect.height) * mouse_tracking_factor;
        }
    
    /*=====  End of MouseTracking  ======*/
    
    

    //SECTION STYLE
    const is_mobile = isMobile();

    export let box_padding = "0";
    




</script>

<div on:mousemove={updateMousePosition} class="a-particles-wrapper" style:padding={box_padding} bind:this={a_particles_wrapper}>
    <div class="apw-container apw-container-1"
        style:top="calc(0% + {box_padding})"
        style:left="calc(0% + {box_padding})"
        style:transform="translate({-mouse_x}%, {-mouse_y}%)"
    >
        <div class="particles-wrapper">
            {#each getParticles() as particle}
                <div class="particle-anima"
                    style:top="{particle.y}%"
                    style:left="{particle.x}%"
                    style:width="{particle.radius * 2}px"
                    style:height="{particle.radius * 2}px"
                    style:background="{particle.color}"
                    style:opacity="{getRandomOpacity()}"
                    style:animation-duration="{Math.floor(Math.random() * movement_stdev) + movement_min_duration}s"
                    style:animation-delay="{Math.floor((Math.random() * delay_stdev) + delay_min)}s"
                ></div>
            {/each}
        </div>
    </div>
    <div class="apw-container apw-container-2"
        style:top="calc(0% + {box_padding})"
        style:right="calc(0% + {box_padding})"
        style:transform="translate({mouse_x}%, {-mouse_y}%)"
    >
        <div class="particles-wrapper">
            {#each getParticles() as particle}
                <div class="particle-anima"
                    style:top="{particle.y}%"
                    style:left="{particle.x}%"
                    style:width="{particle.radius * 2}px"
                    style:height="{particle.radius * 2}px"
                    style:background="{particle.color}"
                    style:opacity="{getRandomOpacity()}"
                    style:animation-duration="{Math.floor(Math.random() * movement_stdev) + movement_min_duration}s"
                    style:animation-delay="{Math.floor((Math.random() * delay_stdev) + delay_min)}s"
                ></div>
            {/each}
        </div>
    </div>
    <div class="apw-container apw-container-3"
        style:bottom="calc(0% + {box_padding})"
        style:left="calc(0% + {box_padding})"
        style:transform="translate({-mouse_x}%, {mouse_y}%)"
    >
        <div class="particles-wrapper">
            {#each getParticles() as particle}
                <div class="particle-anima"
                    style:top="{particle.y}%"
                    style:left="{particle.x}%"
                    style:width="{particle.radius * 2}px"
                    style:height="{particle.radius * 2}px"
                    style:background="{particle.color}"
                    style:opacity="{getRandomOpacity()}"
                    style:animation-duration="{Math.floor(Math.random() * movement_stdev) + movement_min_duration}s"
                    style:animation-delay="{Math.floor((Math.random() * delay_stdev) + delay_min)}s"
                ></div>
            {/each}
        </div>
    </div>
    <div class="apw-container apw-container-4"
        style:bottom="calc(0% + {box_padding})"
        style:right="calc(0% + {box_padding})"
        style:transform="translate({mouse_x}%, {mouse_y}%)"
    >
        <div class="particles-wrapper">
            {#each getParticles() as particle}
                <div class="particle-anima"
                    style:top="{particle.y}%"
                    style:left="{particle.x}%"
                    style:width="{particle.radius * 2}px"
                    style:height="{particle.radius * 2}px"
                    style:background="{particle.color}"
                    style:opacity="{getRandomOpacity()}"
                    style:animation-duration="{Math.floor(Math.random() * movement_stdev) + movement_min_duration}s"
                    style:animation-delay="{Math.floor((Math.random() * delay_stdev) + delay_min)}s"
                ></div>
            {/each}
        </div>
    </div>
</div>

<style>

    
    /*=============================================
    =            Animations            =
    =============================================*/
        
        @keyframes slide-around-pop-inout {
            0% {
                transform: scale(.2) translate(0, 0);
                opacity: 0;
            }
            10% {
                transform: scale(1) translate(0, 90%);
                opacity: .3;
            }
            25% {
                transform: scale(1.4) translate(55%, -10%);
                opacity: .5;
            }
            40% {
                transform: scale(.9) translate(-20%, 35%);
                opacity: 1;
            }
            50% {
                transform: scale(1.3) translate(70%, -12%);
                opacity: 0;
            }
            60% {
                transform: scale(.4) translate(10%, 30%);
                opacity: .8;
            }
            70% {
                transform: scale(1.8) translate(-20%, 120%);
                opacity: .4;
            }
            80% {
                transform: scale(.6) translate(46%, -30%);
                opacity: 1;
            }
            100% {
                transform: scale(1) translate(0, 0);
                opacity: .1;
            }

        }
        
    /*=====  End of Animations  ======*/
    
    

    :global(:has(.a-particles-wrapper)) {
        position: relative;
    }

    .a-particles-wrapper {
        box-sizing: border-box;
        position: absolute;
        width: 100%;
        height: 100%;
        overflow: hidden;
        z-index: 1;
    }

    .apw-container {
        position: absolute;
        width: 20%;
        height: 25%;
        /* border: 1px solid red; */
    }

    .particles-wrapper {
        position: relative;
        width: 100%;
        height: 100%;
        transition: transform 0.05s ease-in-out;
    }

    .particle-anima {
        position: absolute;
        border-radius: 50%;
        animation-name: slide-around-pop-inout;
        animation-fill-mode: both;
        animation-iteration-count: infinite;
    }

    .apw-container-1 {
        top: 0;
        left: 0;
    }

    .apw-container-2 {
        top: 0;
        right: 0;
    }

    .apw-container-3 {
        bottom: 0;
        left: 0;
    }

    .apw-container-4 {
        bottom: 0;
        right: 0;
    }

    @media only screen and (max-width: 768px) {
        .apw-container {
            width: 100%;
        }
        
        .apw-container-2, .apw-container-4 {
            display: none;
        }
    }
</style>