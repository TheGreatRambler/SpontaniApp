<script lang='ts'>
  import { onMount } from 'svelte';
  import { Loader } from '@googlemaps/js-api-loader';
  import DestinationCard from '$lib/destinationCard.svelte';

  let tmpImage = 'rocks.jpg';

  let mapElement: HTMLElement;

  let map: google.maps.Map | undefined;

  onMount(async () => {

    const loader = new Loader({
      apiKey: await (await fetch('https://55ztt2t02i.execute-api.us-east-1.amazonaws.com/prod/get?request_type=get_google_maps_key')).json(),
      version: 'weekly',
      // ...additionalOptions,
    });

    let latitude: number = 0.0;
    let longitude: number = 0.0;
    navigator.geolocation.getCurrentPosition(
      (position: GeolocationPosition) => {
        latitude = position.coords.latitude;
        longitude = position.coords.longitude;
        if (map !== undefined) {
          map.setCenter({ lat: latitude, lng: longitude });
        }
      });

    const mapOptions = {
      zoom: 16
    };

    // Callback
    loader.loadCallback(e => {

      if (e) {
      } else {
        map = new google.maps.Map(mapElement, mapOptions);
      }
    });

  });
</script>

<div class="flex flex-col items-center w-full py-24 bg-primary-300">
  <h1 class="block m-2 text-6xl text-black font-bold">spontani</h1>
  <p class="text-lg text-primary-900">unite through adventure</p>
</div>

<main class="m-4">

  <div bind:this={mapElement} class="m-auto w-full max-w-[800px] aspect-video"></div>

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
