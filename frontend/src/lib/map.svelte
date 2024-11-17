<script lang='ts'>
  import { onMount } from 'svelte';
  import { Loader } from '@googlemaps/js-api-loader';

  let mapElement: HTMLElement;

  let map: google.maps.Map | undefined;

  onMount(async () => {

    const loader = new Loader({
      apiKey: await (await (await fetch('https://55ztt2t02i.execute-api.us-east-1.amazonaws.com/prod/get?request_type=get_google_maps_key')).blob()).text(),
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

<div bind:this={mapElement} class="m-auto w-full max-w-[800px] aspect-video"></div>
