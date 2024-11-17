<script lang="ts">
    import { onMount } from "svelte";
    import {
        Button,
        Label,
        Input,
        Datepicker,
        Fileupload,
    } from "flowbite-svelte";
    import MapComponent from "$lib/map.svelte";

    let loaded = $state(false);
    let selectedStartDate: Date | null = $state(null);
    let selectedEndDate: Date | null = $state(null);

    let form_data = {
        title: "",
        description: "",
        lat: 0,
        lng: 0,
        initial_img_id: 0,
    };

    let lat = $state(0.0);
    let lng = $state(0.0);

    onMount(async () => {
        navigator.geolocation.getCurrentPosition(
            (position: GeolocationPosition) => {
                lat = position.coords.latitude;
                lng = position.coords.longitude;
                loaded = true;
            },
        );
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
            live_place_name_element.innerHTML = await (
                await fetch(
                    `${import.meta.env.VITE_BASE_URL}/get?request_type=location_to_place_name&lat=${lat}&lng=${lng}`,
                )
            ).text();
        }
    };

    let initial_image_loaded = false;
    let files: FileList | undefined = $state();
    let on_file_upload_change = async () => {
        initial_image_loaded = false;
        let file = files![0];

        let res = await (
            await fetch(
                `${import.meta.env.VITE_BASE_URL}/post?request_type=upload_image`,
                {
                    method: "POST",
                    body: file,
                    headers: {
                        "Content-Type": file.type,
                    },
                },
            )
        ).json();

        form_data.initial_img_id = res.id;
        initial_image_loaded = true;
    };

    let on_form_submit = async () => {
        if (selectedStartDate && selectedEndDate) {
            let interval_id = setInterval(async () => {
                if (initial_image_loaded) {
                    let res = await (
                        await fetch(
                            `${import.meta.env.VITE_BASE_URL}/post?request_type=create_task`,
                            {
                                method: "POST",
                                body: JSON.stringify({
                                    title: form_data.title,
                                    description: form_data.description,
                                    lat: lat,
                                    lng: lng,
                                    start: Math.floor(
                                        selectedStartDate!.getTime() / 1000,
                                    ),
                                    stop: Math.floor(
                                        selectedEndDate!.getTime() / 1000,
                                    ),
                                    initial_img_id: form_data.initial_img_id,
                                }),
                            },
                        )
                    ).json();

                    // Update image's task ID
                    await fetch(
                        `${import.meta.env.VITE_BASE_URL}/post?request_type=update_image&id=${form_data.initial_img_id}&task_id=${res.id}`,
                        {
                            method: "POST",
                        },
                    );

                    clearInterval(interval_id);
                }
            }, 500);
        }
    };
</script>

<div class="flex min-h-screen dark:bg-primary-300">
    <div class="w-full flex flex-col space-y-6 px-8 py-4">
        <form class="flex flex-col space-y-6 w-full" action="/">
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">
                Create a New Destination
            </h3>

            <Label class="space-y-2">
                <span>Title</span>
                <Input
                    bind:value={form_data.title}
                    type="text"
                    name="title"
                    placeholder="Day at the Beach"
                    required
                    class="bg-white dark:bg-primary-600 border dark-black dark:border-primary-500 text-gray-900 dark:text-white dark:placeholder-gray-400 rounded-lg p-2 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-300"
                />
            </Label>

            <Label class="space-y-2">
                <span>Description</span>
                <textarea
                    bind:value={form_data.description}
                    id="description"
                    name="description"
                    placeholder="Create sand castles at the beach."
                    required
                    class="block w-full p-2 text-gray-900 border dark:border-primary-500 rounded-lg bg-white dark:bg-primary-600 dark:text-white dark:placeholder-gray-400 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-300"
                ></textarea>
            </Label>

            <Label class="space-y-2">
                <span>Start Date</span>
                <div class="mb-4 md:w-1/2">
                    <Datepicker
                        bind:value={selectedStartDate}
                        inputClass="bg-white dark:bg-primary-600 border black dark:border-primary-500"
                    />
                </div>
            </Label>

            <Label class="space-y-2">
                <span>End Date</span>
                <div class="mb-4 md:w-1/2">
                    <Datepicker
                        bind:value={selectedEndDate}
                        inputClass="bg-white dark:bg-primary-600 border black dark:border-primary-500"
                    />
                </div>
            </Label>

            <Label for="with_helper" class="space-y-2">Upload Image</Label>
            <Fileupload
                bind:files
                onchange={on_file_upload_change}
                id="with_helper"
                class="mb-2"
            />

            <div class="flex space-x-4 justify-center">
                <Button
                    class="bg-primary-600 w-1/4"
                    onclick={() => (window.location.href = "/")}>Cancel</Button
                >
                <Button class="bg-primary-600 w-1/4" onclick={on_form_submit}
                    >Submit</Button
                >
            </div>
        </form>
    </div>

    <div class="w-full flex flex-col space-y-6 px-8 py-4">
        <form class="flex flex-col space-y-6 w-full" action="/">
            <Label class="space-y-2">
                <h3 class="text-xl font-medium text-gray-900 dark:text-white">
                    Location
                </h3>
            </Label>

            {#if loaded}
                <MapComponent
                    markers={[]}
                    start_lat={lat}
                    start_lng={lng}
                    {map_center}
                >
                    <div
                        class="absolute top-1/2 left-1/2 w-4 h-4 bg-red-500 rounded-full transform -translate-x-1/2 -translate-y-1/2 pointer-events-none"
                    ></div>
                </MapComponent>
                <span
                    bind:this={live_place_name_element}
                    class="text-center text-xl font-medium text-gray-900 dark:text-white"
                    >Here</span
                >
            {/if}
        </form>
    </div>
</div>
