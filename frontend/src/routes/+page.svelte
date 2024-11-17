<script lang='ts'>
  import DestinationCard from '$lib/destinationCard.svelte';
  import MapComponent from '$lib/map.svelte';

  let tmpImage = 'rocks.jpg';

  const destinationData = [
    {
      "title": "Rocks",
      "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magnam aliquam quaerat.",
      "lat": 0.0,
      "lng": 0.0,
      "start": 1731825677,
      "stop": 1731925677,
      "initial_img_id": "",
    },
    {
      "title": "Rocks",
      "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magnam aliquam quaerat.",
      "lat": 0.0,
      "lng": 0.0,
      "start": 1731825677,
      "stop": 1731925677,
      "initial_img_id": "",
    },
    {
      "title": "Rocks",
      "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magnam aliquam quaerat.",
      "lat": 0.0,
      "lng": 0.0,
      "start": 1731825677,
      "stop": 1731925677,
      "initial_img_id": "",
    },
    {
      "title": "Rocks",
      "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magnam aliquam quaerat.",
      "lat": 0.0,
      "lng": 0.0,
      "start": 1731825677,
      "stop": 1731925677,
      "initial_img_id": "",
    },
    {
      "title": "Rocks",
      "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magnam aliquam quaerat.",
      "lat": 0.0,
      "lng": 0.0,
      "start": 1731825677,
      "stop": 1731925677,
      "initial_img_id": "",
    },
    {
      "title": "Rocks",
      "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magnam aliquam quaerat.",
      "lat": 0.0,
      "lng": 0.0,
      "start": 1731825677,
      "stop": 1731925677,
      "initial_img_id": "",
    },
    {
      "title": "Rocks",
      "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magnam aliquam quaerat.",
      "lat": 0.0,
      "lng": 0.0,
      "start": 1731825677,
      "stop": 1731925677,
      "initial_img_id": "",
    },
  ];
</script>

<div class="flex flex-col items-center w-full py-24 bg-primary-300">
  <h1 class="block m-2 text-6xl text-black font-bold">spontani</h1>
  <p class="text-lg text-primary-900">unite through adventure</p>
</div>

<main class="m-4">

  <MapComponent />

  <div class="my-12">
    <header class="my-4">
      <h2 class="my-2 mt-2 text-3xl font-bold">current destinations</h2>
      <p>where will you go today?</p>
    </header>
    <div class="grid grid-flow-row grid-cols-1 grid-cols-[repeat(auto-fill,_minmax(250px,_1fr))] gap-4 place-content-center">
      {#each destinationData as dest}
        <DestinationCard description={dest.description} endDate={dest.stop * 1000} img={tmpImage} lat={dest.lat} lng={dest.lng} name={dest.title} />
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
        <DestinationCard description="foo bar eggs spam"  endDate={0} img={tmpImage} lat={0.0} lng={0.0} name="placeholder" />
      {/each}
    </div>
  </div>
</main>
