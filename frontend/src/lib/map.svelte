<script lang='ts'>
  import { onMount } from 'svelte';
  import { Loader } from '@googlemaps/js-api-loader';

  let props: {
    markers: { lat: number; lng: number; title: string }[];
    start_lat: number;
    start_lng: number;
  } = $props();

  let mapElement: HTMLElement;

  let map: google.maps.Map | undefined;

  onMount(async () => {
    const loader = new Loader({
      apiKey: await (await (await fetch('https://f007qjswdf.execute-api.us-east-1.amazonaws.com/prod/get?request_type=get_google_maps_key')).blob()).text(),
      version: 'weekly',
      libraries: ['marker']
    });

    const mapOptions = {
      zoom: 16
    };

    loader.loadCallback(e => {
      if (e) {
        console.error(e);
      } else {
        map = new google.maps.Map(mapElement, mapOptions);

        map.setCenter({lat: props.start_lat, lng: props.start_lng});

        props.markers.forEach((d) => {
          console.log(d);
          let mrkr = new google.maps.Marker({
            map: map,
            position: {lat: d.lat, lng: d.lng},
            title: d.title
          });
        });
      }
    });

  });
</script>

<div bind:this={mapElement} class="m-auto w-full max-w-[800px] aspect-video"></div>
