<script lang='ts'>
  import { Button, Fileupload, Modal } from 'flowbite-svelte';

  let { openModal = $bindable(), name }: { openModal: boolean, name: string } = $props();

  async function onSubmit(event: Event) {
    event.preventDefault();

    const formData = new FormData(event.target);

    const bytes = await (formData.get("file")).arrayBuffer();

    await fetch(`${import.meta.env.VITE_BASE_URL}/post?request_type=upload_image&task_id=1234&caption=testing`, {
      method: 'POST',
      body: bytes,
    });
  }

</script>

<Modal bind:open={openModal} size="sm" outsideclose>
  <h3 class="mb-4 text-xl font-bold text-black dark:text-white">take a selfie at {name.toLowerCase()}</h3>

  <form onsubmit={onSubmit} class="w-full">
    <!-- whee duplication go brr -->
    <!-- if this wasn't a hackathon project, I woulda made the component -->
    <Fileupload name="file" id="with_helper" class="mb-2" />
    <Button type="submit" class="w-full mt-2">Upload</Button>
  </form>

</Modal>
