<script lang="ts">
  import { Button, Carousel, Navbar } from 'flowbite-svelte';
  import type { PageData } from './$types';

  import SubmitModal from "$lib/submit.svelte";

  let { data }: { data: PageData } = $props();

  let openModal: boolean = $state(false);
  let index = $state(0);
</script>

<Navbar class="bg-primary-200 fixed top-0 w-full z-10">
  <Button href="/">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
      <path fill-rule="evenodd" d="M7.72 12.53a.75.75 0 0 1 0-1.06l7.5-7.5a.75.75 0 1 1 1.06 1.06L9.31 12l6.97 6.97a.75.75 0 1 1-1.06 1.06l-7.5-7.5Z" clip-rule="evenodd" />
    </svg>
  </Button>

  <Button onclick={() => openModal=true}>
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
      <path fill-rule="evenodd" d="M12 2.25c-5.385 0-9.75 4.365-9.75 9.75s4.365 9.75 9.75 9.75 9.75-4.365 9.75-9.75S17.385 2.25 12 2.25ZM12.75 9a.75.75 0 0 0-1.5 0v2.25H9a.75.75 0 0 0 0 1.5h2.25V15a.75.75 0 0 0 1.5 0v-2.25H15a.75.75 0 0 0 0-1.5h-2.25V9Z" clip-rule="evenodd" />
    </svg>
  </Button>
</Navbar>

<main class="flex justify-center items-center h-[calc(100vh-56px)] w-dvw overflow-hidden pt-[56px]">
  <div class="w-1/2 h-3/4 flex flex-col justify-center">
    <header class="my-2 text-center">
      <h2 class="my-2 mt-2 text-3xl font-bold">{data.destination.title}</h2>
      <p>{data.destination.description}</p>
    </header>

    <div class="flex justify-center flex-grow">
      <div class="carousel-container w-full max-w-full space-y-4">
        <Carousel 
          images={data.images} 
          let:Indicators 
          let:Controls 
          bind:index 
          class="w-full [height:_350px_!important]"
        >
          <Controls />
          <Indicators />
        </Carousel>
        <div class="dark:text-white p-2 text-center">
          {data.images[index].alt}
        </div>
      </div>
    </div>
  </div>
</main>

<SubmitModal bind:openModal taskId={data.destination.id} name={data.destination.title} showSlidesLink={false} />

<style>
  body {
    margin: 0;
    padding: 0;
    overflow: hidden; /* Prevent scrolling */
  }

  main {
    height: calc(100vh - 56px); /* Subtract navbar height */
  }
</style>
