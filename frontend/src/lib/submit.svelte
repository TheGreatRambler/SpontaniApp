<script lang="ts">
  import { Button, Fileupload, Modal } from "flowbite-svelte";

  let files: FileList | undefined = $state();

  let { openModal = $bindable(), taskId, name }: { openModal: boolean; taskId: number, name: string } =
    $props();

  async function onSubmit(event: Event) {
    event.preventDefault();

    openModal = false;

    let file = files[0];

    await fetch(
      `${import.meta.env.VITE_BASE_URL}/post?request_type=upload_image&task_id=${taskId}&caption=testing`,
      {
        method: "POST",
        body: file,
        headers: {
          "Content-Type": file.type,
        },
      },
    );
  }
</script>

<Modal bind:open={openModal} size="sm" outsideclose>
  <h3 class="mb-4 text-xl font-bold text-black dark:text-white">
    take a selfie at {name.toLowerCase()}
  </h3>

  <form onsubmit={onSubmit} class="w-full">
    <!-- whee duplication go brr -->
    <!-- if this wasn't a hackathon project, I woulda made the component -->
    <Fileupload bind:files name="file" id="with_helper" class="mb-2" />
    <Button type="submit" class="w-full mt-2">Upload</Button>
  </form>
</Modal>
