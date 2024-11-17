<script lang="ts">
  import { Carousel, Thumbnails } from 'flowbite-svelte';
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  let index = $state(0);

</script>

<main class="mx-4">
  <div class="my-5">
    <header class="my-2 text-center">
      <h2 class="my-2 mt-2 text-3xl font-bold">{data.destination.name}</h2>
      <p>{data.destination.description}</p>
    </header>

    <div class="flex justify-center items-center">
      <div class="carousel-container max-w-3xl space-y-4">
        <Carousel images={data.images} let:Indicators let:Controls bind:index>
          <Controls />
          <Indicators />
        </Carousel>
        <Thumbnails images={data.images} bind:index />
      </div>
    </div>

  </div>
</main>
