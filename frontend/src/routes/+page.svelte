<script lang="ts">
  import { onMount } from "svelte";
  import DestinationCard from "$lib/destinationCard.svelte";
  import MapComponent from "$lib/map.svelte";

  let tmpImage = "rocks.jpg";

  let destinationData = $state([] as any[]);
  let completedDestinationData = $state([] as any[]);

  let loaded = $state(false);

  let start_lat = $state(0.0);
  let start_lng = $state(0.0);

  let markers = $state([] as any[]);

  // { lat: 32.98599729543064, lng: -96.7508045889115, title: "hello" }

  let get_recent_tasks = async (query: string) => {
    let res = await fetch(`${import.meta.env.VITE_BASE_URL}${query}`);
    let taskData = await res.json();

    let newDestinationData: any[] = [];
    for (let task of taskData) {
      let initial_image_url = (
        await (
          await fetch(
            `${import.meta.env.VITE_BASE_URL}/post?request_type=get_presigned_url&id=${task.initial_img_id}`,
            { method: "POST" },
          )
        ).json()
      ).url;
      newDestinationData.push({
        ...task,
        initial_image_url: initial_image_url,
      });
    }

    return newDestinationData;
  };

  onMount(async () => {
    navigator.geolocation.getCurrentPosition(
      (position: GeolocationPosition) => {
        start_lat = position.coords.latitude;
        start_lng = position.coords.longitude;
        loaded = true;

        (async function () {
          destinationData = await get_recent_tasks(
            `/get?request_type=get_nearby_recent_tasks&lat=${start_lat}&lng=${start_lng}`,
          );
        })();

        (async function () {
          completedDestinationData = await get_recent_tasks(
            "/get?request_type=get_completed_tasks",
          );
        })();
      },
    );
  });

  let last_update = 0;
  let map_center = async (map: google.maps.Map) => {
    let center = map.getCenter();

    let timestamp = Date.now();
    if (timestamp - last_update > 1000 * 1) {
      last_update = timestamp;

      // Trigger new search
      destinationData = await get_recent_tasks(
        `/get?request_type=get_nearby_recent_tasks&lat=${center!.lat()}&lng=${center!.lng()}`,
      );
    }
  };
</script>

  <button
    class="absolute top-4 right-4 sm:top-6 sm:right-6 bg-primary-500 text-white font-bold rounded-full size-16 text-4xl flex items-center justify-center hover:bg-primary-600 focus:outline-none"
    aria-label="Add"
    onclick={() => (window.location.href = "/form")}
  >
    +
  </button>

<div class="flex flex-col items-center w-full py-24 bg-primary-300">
  <h1 class="block m-2 text-6xl text-black font-bold">spontani</h1>
  <p class="text-lg text-primary-900">unite through adventure</p>

</div>

<main class="m-4">
  {#if loaded}
    <MapComponent
      markers={destinationData
        .map((task) => {
          return {
            lat: task.lat,
            lng: task.lng,
            title: task.title,
            color: "#4ECDC4",
            id: task.id,
          };
        })
        .concat(
          completedDestinationData.map((task) => {
            return {
              lat: task.lat,
              lng: task.lng,
              title: task.title,
              color: "#FF6B6B",
              id: task.id,
            };
          }),
        )}
      {start_lat}
      {start_lng}
      {map_center}
    />
  {/if}

  <div class="my-12">
    <header class="my-4">
      <h2 class="my-2 mt-2 text-3xl font-bold">current destinations</h2>
      <p>where will you go today?</p>
    </header>
    {#if destinationData.length < 1}
      <p>no existing destinations - add one with the plus button!</p>
    {/if}
    <div
      class="grid grid-flow-row grid-cols-1 grid-cols-[repeat(auto-fill,_minmax(250px,_1fr))] gap-4 place-content-center"
    >
      {#each destinationData as dest}
        <DestinationCard
          description={dest.description}
          endDate={dest.stop * 1000}
          img={dest.initial_image_url}
          lat={dest.lat}
          lng={dest.lng}
          name={dest.title}
          id={dest.id}
        />
      {/each}
    </div>
  </div>

  <div class="my-12">
    <header class="my-4">
      <h2 class="my-2 mt-2 text-3xl font-bold">past destinations</h2>
      <p>
        the deadline for going to these is over; now you can see everyone's
        pictures!
      </p>
    </header>
    <div
      class="grid grid-flow-row grid-cols-1 grid-cols-[repeat(auto-fill,_minmax(250px,_1fr))] gap-4 place-content-center"
    >
      {#each completedDestinationData as dest}
        <DestinationCard
          description={dest.description}
          endDate={dest.stop * 1000}
          img={dest.initial_image_url}
          lat={dest.lat}
          lng={dest.lng}
          name={dest.title}
          id={dest.id}
        />
      {/each}
    </div>
  </div>
</main>
