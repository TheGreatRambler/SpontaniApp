<script lang='ts'>
    import { onMount } from 'svelte';
    import { Card, Button, Label, Input, Datepicker, Fileupload, Helper } from 'flowbite-svelte';
    import MapComponent from '$lib/map.svelte';

    let loaded = $state(false);

    let form_data = {
      title: "",
      description: "",
      lat: 0,
      lng: 0,
      start: new Date(),
      stop: new Date(),
      initial_image_id: 0,
    };

    let lat = $state(0.0);
    let lng = $state(0.0);

    onMount(async () => {
      navigator.geolocation.getCurrentPosition((position: GeolocationPosition) => {
        lat = position.coords.latitude;
        lng = position.coords.longitude;
        loaded = true;
      });
    });

    let live_place_name_element: HTMLSpanElement;
    let last_update = 0;
    let map_center = async (map: google.maps.Map) => {
        let center = map.getCenter();
        lat = center!.lat();
        lng = center!.lng();
        
        let timestamp = Date.now();
        if (timestamp - last_update > 1000 * 1) {
            last_update = timestamp;
            live_place_name_element.innerHTML = await (await fetch(`${import.meta.env.VITE_BASE_URL}/get?request_type=location_to_place_name&lat=${lat}&lng=${lng}`)).text();
        }
    };

    let on_form_submit = () => {
        console.log({
            title: form_data.title,
            description: form_data.description,
            lat: lat,
            lng: lng,
            start: Math.floor(form_data.start.getTime() / 1000),
            stop: Math.floor(form_data.stop.getTime() / 1000),
            initial_image_id: 0,
          });

          /*
        const response = await fetch(`${import.meta.env.VITE_BASE_URL}/post?request_type=create_task`, {
          method: "POST",
          body: JSON.stringify({
            title: form_data.title,
            description: form_data.description,
            lat: lat,
            lng: lng,
            start: Math.floor(form_data.start.getTime() / 1000),
            stop: Math.floor(form_data.stop.getTime() / 1000),
            initial_image_id: 0,
          }),
        });
        */

    };
</script>

<div class="flex min-h-screen dark:bg-primary-300">
    <div class="w-full flex flex-col space-y-6 px-8 py-4">
        <form class="flex flex-col space-y-6 w-full" action="/">
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">Create a New Destination</h3>
            
            <Label class="space-y-2">
                <span>Title</span>
                <Input 
                    bind:value={form_data.title} type="text" name="title" placeholder="Day at the Beach" required class="bg-white dark:bg-primary-600 border dark-black dark:border-primary-500 text-gray-900 dark:text-white dark:placeholder-gray-400 rounded-lg p-2 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-300"/>
            </Label>
            
            <Label class="space-y-2">
                <span>Description</span>
                <textarea bind:value={form_data.description} id="description" name="description" placeholder="Create sand castles at the beach." required class="block w-full p-2 text-gray-900 border dark:border-primary-500 rounded-lg bg-white dark:bg-primary-600 dark:text-white dark:placeholder-gray-400 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-300"></textarea>
            </Label>

            <Label class="space-y-2">
                <span>Start Date</span>
                <div class="mb-4 md:w-1/2">
                    <Datepicker bind:value={form_data.start} inputClass="bg-white dark:bg-primary-600 border border-primary-300 dark:border-primary-500" />
                </div>
            </Label>

            <Label class="space-y-2">
                <span>End Date</span>
                <div class="mb-4 md:w-1/2">
                    <Datepicker bind:value={form_data.stop} inputClass="bg-white dark:bg-primary-600 border border-primary-300 dark:border-primary-500" />
                </div>
            </Label>
            
            <Label for="with_helper" class="space-y-2">Upload Image</Label>
            <Fileupload id="with_helper" class="mb-2" />
            <Helper>SVG, PNG, JPG or GIF (MAX. 800x400px).</Helper>

            <Button class="bg-primary-600" onclick={on_form_submit}>Submit</Button>
        </form>
    </div>


    <div class="w-full flex flex-col space-y-6 px-8 py-4">
        <form class="flex flex-col space-y-6 w-full" action="/">
                <Label class="space-y-2">
                    <h3 class="text-xl font-medium text-gray-900 dark:text-white">Location</h3>
                </Label>

                {#if loaded}
                    <MapComponent markers={[]} start_lat={lat} start_lng={lng} map_center={map_center}>
                        <div class="absolute top-1/2 left-1/2 w-4 h-4 bg-red-500 rounded-full transform -translate-x-1/2 -translate-y-1/2 pointer-events-none"></div>
                        <span bind:this={live_place_name_element} class="text-center text-xl font-medium text-gray-900 dark:text-white">Here</span>
                    </MapComponent>
                {/if}
        </form>
    </div>
</div>
