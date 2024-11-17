<script lang='ts'>
    import { Carousel, Thumbnails } from 'flowbite-svelte';
  
    export const images = [
        {
            alt: '1989',
            src: 'rocks.jpg',
            title: 'Rocks'
        },
        {
            alt: '1990',
            src: 'mountains.jpg',
            title: 'Mountains'
        },
        {
            alt: '1991',
            src: 'beach.jpg',
            title: 'Beach'
        },
        {
            alt: '1992',
            src: 'raindrop.jpg',
            title: 'Raindrops'
        }
    ];

    let image = images[0];
    let index = 0;
</script>

<main class="mx-4">
    <div class="my-5">
        <header class="my-2 text-center">
            <h2 class="my-2 mt-2 text-3xl font-bold">Day at the Beach</h2>
            <p>Create sand castles at the beach.</p>
        </header>

        <div class="flex justify-center items-center">
            <div class="carousel-container max-w-3xl space-y-4">
                <Carousel {images} let:Indicators let:Controls bind:index on:change={({ detail }) => (image = detail)}>
                    <Controls />
                    <Indicators />
                </Carousel>
                <Thumbnails {images} bind:index />
            </div>
        </div>

        <div class=" dark:text-white p-2 my-2 text-center">
            {images[index].alt}
        </div>
    </div>
</main>