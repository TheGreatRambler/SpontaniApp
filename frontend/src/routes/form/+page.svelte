<script>
    import { Card, Button, Label, Input, Datepicker, P } from 'flowbite-svelte';
    let selectedDate = null;
</script>

<div class="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-gray-900">
    <Card>
        <form class="flex flex-col space-y-6" action="/">
            <h3 class="text-xl font-medium text-gray-900 dark:text-white">Create a New Destination</h3>
            
            <Label class="space-y-2">
                <span>Title</span>
                <Input type="text" name="title" placeholder="Adventures at the Dinosaur Park" required />
            </Label>
            
            <Label class="space-y-2">
                <label for="description" class="block mb-2 text-sm font-small text-gray-900 dark:text-white">Description</label>
                <textarea id="description" name="description" placeholder="Go to the park and take photos on top of the dinosaur statues." required class="block w-full p-2 text-gray-900 border border-gray-300 rounded-lg bg-gray-50 text-base focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 placeholder:text-sm"></textarea>
            </Label>

            <Label class="space-y-2">
                <span>Choose a Date</span>
                <div class="mb-4 md:w-1/2">
                    <Datepicker bind:value={selectedDate} />
                    <!-- <P class="mt-4">Selected date: {selectedDate ? selectedDate.toLocaleDateString() : 'None'}</P> -->
                </div>
            </Label>

            <Label class="space-y-2">
                <!-- add location here -->
                <span>Location</span>
            </Label>
            
            <Button class="!bg-blue-500">Submit</Button>
        </form>
    </Card>
</div>