<script lang='ts'>
  import { onMount } from 'svelte';
  import { Loader } from '@googlemaps/js-api-loader';

  let { markers }: { markers: {lat: number, lng: number, title: string}[] } = $props();

  let mapElement: HTMLElement;

  let map: google.maps.Map | undefined;

  onMount(async () => {

    const loader = new Loader({
      apiKey: await (await (await fetch('https://f007qjswdf.execute-api.us-east-1.amazonaws.com/prod/get?request_type=get_google_maps_key')).blob()).text(),
      version: 'weekly',
      libraries: ['marker']
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

    loader.loadCallback(e => {
      if (e) {
        console.error(e);
      } else {
        map = new google.maps.Map(mapElement, mapOptions);

        markers.forEach((d) => {
          console.log(d);
          let mrkr = new google.maps.marker.AdvancedMarkerElement({
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
