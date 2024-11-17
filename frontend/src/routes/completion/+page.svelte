<script lang='ts'>
    import { Carousel } from 'flowbite-svelte';
  
    export const images = [
        {
            alt: 'A picture of rocks',
            src: 'rocks.jpg',
            title: 'Rocks'
        },
        {
            alt: 'A picture of snowy mountains',
            src: 'mountains.jpg',
            title: 'Mountains'
        },
        {
            alt: 'A picture of a beach',
            src: 'beach.jpg',
            title: 'Beach'
        },
        {
            alt: 'A picture of raindrops',
            src: 'raindrop.jpg',
            title: 'Raindrop'
        }
    ];

    let image = images[0];
</script>

<main class="mx-4">
    <div class="my-5">
        <header class="my-4">
            <h2 class="my-2 mt-2 text-3xl font-bold">Day at the Beach</h2>
            <p>Create sand castles at the beach.</p>
        </header>

        <div class="max-w-4xl space-y-4">
            <Carousel {images} let:Indicators let:Controls on:change={({ detail }) => (image = detail)}>
              <Controls />
              <Indicators />
            </Carousel>
            <div class="rounded h-10 bg-gray-300 dark:bg-gray-700 dark:text-white p-2 my-2 text-center">
                {image.alt}
            </div>
        </div> 
    </div>
</main>