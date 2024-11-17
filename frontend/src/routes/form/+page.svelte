<script>
    import { Card, Button, Label, Input, Datepicker, Fileupload, Helper } from 'flowbite-svelte';
    let selectedDate = null;
    let map;

    const location = {
        lat: 32.9857,
        lng: 96.7502
    };

    import { onMount } from 'svelte';

    onMount(() => {
        if (typeof google !== 'undefined') {
            map = new google.maps.Map(document.getElementById('map'), {
                center: location,
                zoom: 8
            });
        }

        new google.maps.Marker({
            position: location,
            map: map,
            title: 'London'
        });
    });
</script>

<div class="flex min-h-screen dark:bg-primary-300 ">
    <div class="w-full flex flex-col space-y-6 px-8 py-4">
        <form class="flex flex-col space-y-6 w-full" action="/">
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">Create a New Destination</h3>
            
            <Label class="space-y-2 ">
                <span>Title</span>
                <Input type="text" name="title" placeholder="Day at the Beach" required />
            </Label>
            
            <Label class="space-y-2">
                <span>Description</span>
                <textarea id="description" name="description" placeholder="Create sand castles at the beach." required class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-base dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:bg-primary-300 dark:focus:bg-primary-300 placeholder:text-sm"></textarea>
            </Label>

            <Label class="space-y-2">
                <span>Choose a Date</span>
                <div class="mb-4 md:w-1/2">
                    <Datepicker bind:value={selectedDate} />
                </div>
            </Label>
            
            <Label for="with_helper" class="space-y-2">Upload Image</Label>
            <Fileupload id="with_helper" class="mb-2" />
            <Helper>SVG, PNG, JPG or GIF (MAX. 800x400px).</Helper>

            <Button class="bg-primary-600">Submit</Button>
        </form>
    </div>


    <div class="w-full flex flex-col space-y-6 px-8 py-4">
        <form class="flex flex-col space-y-6 w-full" action="/">
               
        <Label class="space-y-2">
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">Location</h3>
        </Label>
        <div id="map" class="h-96 w-full"></div>
    </div>
</div>
