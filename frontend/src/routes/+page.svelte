<script lang='ts'>
  import DestinationCard from '$lib/destinationCard.svelte';

  let tmpImage = 'rocks.jpg';
</script>

<div class="flex flex-col items-center w-full py-24 bg-primary-300">
  <h1 class="block m-2 text-6xl text-black font-bold">spontani</h1>
  <p class="text-lg text-primary-900">unite through adventure</p>
</div>

<main class="mx-4">
  <div class="my-12">
    <header class="my-4">
      <h2 class="my-2 mt-2 text-3xl font-bold">current destinations</h2>
      <p>where will you go today?</p>
    </header>
    <div class="grid grid-flow-row grid-cols-1 grid-cols-[repeat(auto-fill,_minmax(250px,_1fr))] gap-4 place-content-center">
      {#each Array(10) as _}
        <DestinationCard img={tmpImage} endDate={new Date(Date.now())} name="Rocks" />
      {/each}
    </div>
  </div>

  <div class="my-12">
    <header class="my-4">
      <h2 class="my-2 mt-2 text-3xl font-bold">past destinations</h2>
      <p>the deadline for going to these is over; now you can see everyone's pictures!</p>
    </header>
    <div class="grid grid-flow-row grid-cols-1 grid-cols-[repeat(auto-fill,_minmax(250px,_1fr))] gap-4 place-content-center">
      {#each Array(3) as _}
        <DestinationCard img={tmpImage} endDate={new Date(Date.now())} name="Rocks" />
      {/each}
    </div>
  </div>
</main>
