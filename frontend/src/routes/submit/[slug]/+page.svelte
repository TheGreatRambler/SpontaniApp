<script lang="ts">
  import { Button, Carousel, Navbar } from 'flowbite-svelte';
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  let index = $state(0);

</script>

<Navbar class="bg-primary-200">

  <Button href="/">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
      <path fill-rule="evenodd" d="M7.72 12.53a.75.75 0 0 1 0-1.06l7.5-7.5a.75.75 0 1 1 1.06 1.06L9.31 12l6.97 6.97a.75.75 0 1 1-1.06 1.06l-7.5-7.5Z" clip-rule="evenodd" />
    </svg>
  </Button>

  <Button href="/form">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
      <path fill-rule="evenodd" d="M12 2.25c-5.385 0-9.75 4.365-9.75 9.75s4.365 9.75 9.75 9.75 9.75-4.365 9.75-9.75S17.385 2.25 12 2.25ZM12.75 9a.75.75 0 0 0-1.5 0v2.25H9a.75.75 0 0 0 0 1.5h2.25V15a.75.75 0 0 0 1.5 0v-2.25H15a.75.75 0 0 0 0-1.5h-2.25V9Z" clip-rule="evenodd" />
    </svg>
  </Button>

</Navbar>

<main class="w-dvw h-dvh flex justify-center place-items-center">
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
