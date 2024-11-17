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

<div class="flex min-h-screen dark:bg-primary-300">
    <div class="w-full flex flex-col space-y-6 px-8 py-4">
        <form class="flex flex-col space-y-6 w-full" action="/">
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">Create a New Destination</h3>
            
            <Label class="space-y-2">
                <span>Title</span>
                <Input 
                    type="text" name="title" placeholder="Day at the Beach" required class="bg-white dark:bg-primary-600 border dark-black dark:border-primary-500 text-gray-900 dark:text-white dark:placeholder-gray-400 rounded-lg p-2 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-300"/>
            </Label>
            
            <Label class="space-y-2">
                <span>Description</span>
                <textarea id="description" name="description" placeholder="Create sand castles at the beach." required class="block w-full p-2 text-gray-900 border dark:border-primary-500 rounded-lg bg-white dark:bg-primary-600 dark:text-white dark:placeholder-gray-400 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-300"></textarea>
            </Label>

            <Label class="space-y-2">
                <span>Choose a Date</span>
                <div class="mb-4 md:w-1/2">
                    <Datepicker bind:value={selectedDate} class="bg-white dark:bg-primary-600 border border-primary-300 dark:border-primary-500" />
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
        </form>
    </div>
</div>
