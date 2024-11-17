<script lang='ts'>
  import { Modal } from 'flowbite-svelte';

  let { img, openModal = $bindable(), name }: { img: any, openModal: boolean, name: string } = $props();
</script>

<Modal bind:open={openModal} outsideclose>
  <h3 class="mb-4 text-xl font-bold text-black dark:text-white">recreate the picture of {name.toLowerCase()}</h3>
  <p>align the camera as close as you can to the original and take a selfie!</p>
</Modal>
