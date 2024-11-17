<script lang="ts">
  import { Card } from "flowbite-svelte";

  import SubmitModal from "$lib/submit.svelte";

  import type { Destination } from "$lib/destinationInterface.ts";

  let { description, endDate, img, lat, lng, name, id, active }: Destination =
    $props();
  let openModal: boolean = $state(false);

  const tzOffset = new Date().getTimezoneOffset() * 60000;

  const expired = Date.now() - endDate > 0;

  const handleClick = (endDate) => {
    if (endDate > Date.now()) {
      openModal = true;
    } else {
      window.location.href = `/submit/${id}`;
    }
  };
</script>

<Card
  onclick={() => {
    handleClick(endDate);
  }}
  {img}
  class="max-w-[600px] hover:bg-gray-200"
  imgClass="aspect-[3/2] object-cover"
>
  <h3 class="font-bold text-lg text-black">{name.toLowerCase()}</h3>
  <p class="text-sm text-gray-600">
    {expired ? "closed on" : "available until"}
    {new Date(endDate - tzOffset).toISOString().slice(0, 10)}
  </p>
  <p class="text-gray-900">{description}</p>
</Card>
<SubmitModal bind:openModal taskId={id} {name} showSlidesLink={true} />
