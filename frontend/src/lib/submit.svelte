<script lang='ts'>
  import { Button, Fileupload, Input, Modal } from 'flowbite-svelte';

  let { openModal = $bindable(), name }: { openModal: boolean, name: string } = $props();

</script>

<Modal bind:open={openModal} size="sm" outsideclose>
  <h3 class="mb-4 text-xl font-bold text-black dark:text-white">take a selfie at {name.toLowerCase()}</h3>

  <form class="w-full">
    <!-- whee duplication go brr -->
    <!-- if this wasn't a hackathon project, I woulda made the component -->
    <Fileupload id="with_helper" class="mb-2" />
    <Input type="text" placeholder="Enter a caption (optional)" />
    <Button type="submit" class="w-full mt-2">Upload</Button>
  </form>

</Modal>
