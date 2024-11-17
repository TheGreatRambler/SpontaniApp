<script lang="ts">
  import { Carousel, Thumbnails } from 'flowbite-svelte';
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  let index = $state(0);

</script>

<main class="flex justify-center items-center w-dvw h-dvh">
  <div class="w-1/2 h-1/2">
    <header class="my-2 text-center">
      <h2 class="my-2 mt-2 text-3xl font-bold">{data.destination.title}</h2>
      <p>{data.destination.description}</p>
    </header>

    <div class="flex justify-center">
      <div class="carousel-container w-full max-w-full space-y-4">
        <Carousel images={data.images} let:Indicators let:Controls bind:index class="w-full [height:_350px_!important]">
          <Controls />
          <Indicators />
        </Carousel>
        <div class=" dark:text-white p-2 text-center">
            {data.images[index].alt}
        </div>
      </div>
    </div>
  </div>
</main>
