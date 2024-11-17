<script lang="ts">
  import { Button, Fileupload, Modal, Label } from "flowbite-svelte";

  let files: FileList | undefined = $state();

  let caption = "";

  let {
    openModal = $bindable(),
    taskId,
    name,
    showSlidesLink
  }: { openModal: boolean; taskId: number; name: string, showSlidesLink: boolean } = $props();

  async function onSubmit(event: Event) {
    event.preventDefault();

    openModal = false;

    let file = files![0];

    await fetch(
      `${import.meta.env.VITE_BASE_URL}/post?request_type=upload_image&task_id=${taskId}&caption=${caption}`,
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
  {#if showSlidesLink}
    <a class="block text-orange-600 underline mb-4 text-xl font-bold" href={`/submit/${taskId}`}>
      View submissions
    </a>
  {/if}

  <Label class="space-y-2">
    <span>Caption</span>
    <textarea
      bind:value={caption}
      id="caption"
      name="caption"
      placeholder="That BBQ was insane!"
      required
      class="block w-full p-2 text-gray-900 border dark:border-primary-500 rounded-lg bg-white dark:bg-primary-600 dark:text-white dark:placeholder-gray-400 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-300"
    ></textarea>
  </Label>

  <form onsubmit={onSubmit} class="w-full">
    <!-- whee duplication go brr -->
    <!-- if this wasn't a hackathon project, I woulda made the component -->
    <Fileupload bind:files name="file" id="with_helper" class="mb-2" />
    <Button type="submit" class="w-full mt-2">Upload</Button>
  </form>
</Modal>
