<script lang="ts">
  import { onMount } from "svelte";
  import type { Snippet } from "svelte";
  import { Loader } from "@googlemaps/js-api-loader";

  let props: {
    children?: Snippet;
    markers: { lat: number; lng: number; title: string; color: string; id: number }[];
    start_lat: number;
    start_lng: number;
    map_center?: (map: google.maps.Map) => void;
  } = $props();

  let mapElement: HTMLElement;

  let map: google.maps.Map | undefined;

  onMount(async () => {
    const loader = new Loader({
      apiKey: await (
        await (
          await fetch(
            `${import.meta.env.VITE_BASE_URL}/get?request_type=get_google_maps_key`,
          )
        ).blob()
      ).text(),
      version: "weekly",
      libraries: ["marker"],
    });

    const mapOptions = {
      zoom: 16,
    };

    loader.loadCallback((e) => {
      if (e) {
        console.error(e);
      } else {
        map = new google.maps.Map(mapElement, mapOptions);

        map.setCenter({ lat: props.start_lat, lng: props.start_lng });

        if (props.map_center) {
          map.addListener("center_changed", () => {
            props.map_center!(map!);
          });
        }
      }
    });
  });

  let active_markers: any[] = [];
  $effect(() => {
    for (let i = 0; i < active_markers.length; i++) {
      active_markers[i].setMap(null);
    }
    active_markers.length = 0;

    props.markers.forEach((d) => {
      let mrkr = new google.maps.Marker({
        map: map,
        position: { lat: d.lat, lng: d.lng },
        title: d.title,
        icon: {
          path: google.maps.SymbolPath.CIRCLE,
          scale: 10,
          fillColor: d.color,
          fillOpacity: 1,
          strokeColor: "#FFFFFF",
          strokeWeight: 2,
        },
      });

      google.maps.event.addListener(mrkr, "mouseover", () => {
        mrkr.setIcon({
          path: google.maps.SymbolPath.CIRCLE,
          scale: 12, // Increase size on hover
          fillColor: d.color,
          fillOpacity: 1,
          strokeColor: "#000000", // Change stroke color on hover
          strokeWeight: 3,
        });
      });

      google.maps.event.addListener(mrkr, "mouseout", () => {
        mrkr.setIcon({
          path: google.maps.SymbolPath.CIRCLE,
          scale: 10, // Revert to original size
          fillColor: d.color,
          fillOpacity: 1,
          strokeColor: "#FFFFFF", // Revert stroke color
          strokeWeight: 2,
        });
      });

      google.maps.event.addListener(mrkr, "click", () => {
        window.location.href = `/submit/${d.id}`;
      });

      active_markers.push(mrkr);
    });
  });
</script>

<div class="relative m-auto w-full max-w-[800px] aspect-video">
  <!-- The map container -->
  <div bind:this={mapElement} class="absolute inset-0"></div>

  <!-- Children rendered on top of the map -->
  <div class="absolute inset-0 pointer-events-none">
    {@render props.children?.()}
  </div>
</div>
